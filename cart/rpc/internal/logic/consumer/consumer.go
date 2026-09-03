package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"zeromall/cart/rpc/cartPb"
	"zeromall/cart/rpc/internal/model"
	"zeromall/cart/rpc/internal/svc"
	"zeromall/common/constant"
	"zeromall/common/convert"
	"zeromall/common/mq"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type Consumer struct {
	consumer rmq_client.SimpleConsumer
	closeCh  chan struct{}
	svc      *svc.ServiceContext
}

func NewConsumer(svc *svc.ServiceContext) (*Consumer, error) {
	config := &rmq_client.Config{
		Endpoint:      svc.Config.RocketMqConf.Endpoint,
		ConsumerGroup: svc.Config.RocketMqConf.Consumer.Group,
		Credentials:   nil,
	}
	con, err := rmq_client.NewSimpleConsumer(config,
		rmq_client.WithSimpleAwaitDuration(svc.Config.RocketMqConf.Consumer.AwaitDuration),
		rmq_client.WithSimpleSubscriptionExpressions(map[string]*rmq_client.FilterExpression{
			svc.Config.RocketMqConf.Consumer.TopicSyncFiling: rmq_client.SUB_ALL,
		}),
	)
	if err != nil {
		return nil, err
	}
	if err = con.Start(); err != nil {
		return nil, err
	}
	return &Consumer{
		consumer: con,
		closeCh:  make(chan struct{}),
		svc:      svc,
	}, nil

}

func (c *Consumer) StartConsumer(topic string) error {
	err := c.consumer.Subscribe(topic, rmq_client.NewFilterExpression("*"))
	if err != nil {
		return err
	}
	go c.LoopPull()
	return nil
}
func (c *Consumer) LoopPull() {
	for {
		select {
		case <-c.closeCh:
			log.Println("cart consumer loop exit")
			return
		default:
		}
		// 拉取消息，等待
		mvs, err := c.consumer.Receive(context.Background(),
			c.svc.Config.RocketMqConf.Consumer.MaxMsgNum,
			c.svc.Config.RocketMqConf.Consumer.LoopDuration,
		)
		if err != nil {
			// 没有消息、超时属于正常，sleep一小段继续
			time.Sleep(500 * time.Millisecond)
			continue
		}
		for _, msg := range mvs {
			userId := string(msg.GetBody())
			// 调用归档组件：userId放入redis待同步set
			err := c.HandleCartChangeMsg(context.Background(), userId)
			if err != nil {

				logx.Infof("HandleCartChangeMsg err userId=%s err=%v", userId, err)
				// ❗处理失败：不Ack，SimpleConsumer会重新投递
				continue
			}

			// ✅消费成功，手动ACK
			ackErr := c.consumer.Ack(context.Background(), msg)
			if ackErr != nil {
				logx.Infof("ack message failed msgId=%s err=%v", msg.GetMessageId(), ackErr)
			}
		}
	}
}
func (c *Consumer) StopConsumer() error {
	return c.consumer.GracefulStop()
}

// HandleCartChangeMsg MQ回调调用，把userId加入Redis待同步Set
func (c *Consumer) HandleCartChangeMsg(ctx context.Context, userId string) error {
	_, err := c.svc.Redis.SaddCtx(ctx, constant.PendingSyncCartKey, userId)
	return err
}
func (c *Consumer) StartTicker(ctx context.Context) {
	go func() {
		// ticker写在这里！！
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("cart sync archiver exit")
				return
			case <-ticker.C:

				// 每3秒原子取出待同步用户
				keys := []string{constant.PendingSyncCartKey}
				//为了支持一次取出n个用户，所以采用原生redis
				res, err := c.svc.Redis.EvalCtx(ctx, `return redis.call('SPOP', KEYS[1], 50) `, keys, nil)
				if err != nil {
					logc.Infof(ctx, "eval spop err: %v", err)
					continue
				}

				userIds := make([]string, 0)
				if arr, ok := res.([]interface{}); ok {
					for _, v := range arr {
						if str, ok2 := v.(string); ok2 {
							var msg mq.CartChangeMsg
							err = json.Unmarshal([]byte(str), &msg)
							if err != nil {
								logc.Infof(ctx, "Unmarshal err: %v", err)
							}
							userIds = append(userIds, msg.UserId)
						}
					}
				}
				// 循环执行归档同步MySQL
				for _, uid := range userIds {
					_ = c.syncCartToMysql(ctx, uid)
				}
			}
		}
	}()
}
func (c *Consumer) syncCartToMysql(ctx context.Context, userId string) error {
	logc.Infof(ctx, "syncCartToMysql userId=%s", userId)
	// 1. HGETALL cart:{userId}
	key := constant.CartKey + userId
	redisMap, err := c.svc.Redis.HgetallCtx(ctx, key)
	if err != nil {
		return err
	}
	redisItemMap := make(map[string]*model.Cart)
	for goodsId, jsonStr := range redisMap {
		var item cartPb.CartItem
		err = json.Unmarshal([]byte(jsonStr), &item)
		if err != nil {
			logx.Errorf("json unmarshal error:%v in %s ", err, "syncCartToMysql")
			continue
		}
		priceCent, err := convert.YuanStrToCents(item.Price)
		if err != nil {
			log.Println("price转化失败", err)
			continue
		}
		originPriceCent, _ := convert.YuanStrToCents(item.OriginPrice)
		redisItemMap[goodsId] = &model.Cart{
			UserId:          userId,
			GoodsId:         goodsId,
			Num:             item.Num,
			Selected:        item.Selected,
			PriceCent:       priceCent,
			Name:            item.Name,
			Cover:           item.Cover,
			OriginPriceCent: originPriceCent,
		}
	}
	//获取mysql数据也
	res, err := c.svc.CartModel.FindCartsByUserId(ctx, userId)
	if err != nil {
		return err
	}
	//转map[string]*cart类型
	mysqlItemMap := make(map[string]*model.Cart)
	for _, v := range res {
		mysqlItemMap[v.GoodsId] = v
	}
	// diff 三组：待新增、待更新、待删除
	var toInsert []*model.Cart
	var toUpdate []*model.Cart
	var toDeleteIds []string
	//3.判断redis中数据与mysql中不同,
	//遍历redis，看mysql是否存在，redis有，mysql没有，则插入mysql中redis有的数据
	for goodsId, redisItem := range redisItemMap {
		dbItem, exist := mysqlItemMap[goodsId]
		if !exist {
			//mysql不存在，则加入待新增
			toInsert = append(toInsert, redisItem)
		} else {
			//比较他们的各个字段
			if redisItem.Num != dbItem.Num || redisItem.Selected != dbItem.Selected ||
				redisItem.UserId != dbItem.UserId || redisItem.Cover != dbItem.Cover ||
				redisItem.PriceCent != dbItem.PriceCent || redisItem.OriginPriceCent != dbItem.OriginPriceCent {
				//有一个字段不一样，就加入待更新
				toUpdate = append(toUpdate, redisItem)
			}

		}
	}

	//遍历mysql
	for goodsId, dbItem := range mysqlItemMap {
		// redis没有，mysql有，则删除mysql中redis没有的数据
		_, exist := redisItemMap[goodsId]
		if !exist {
			toDeleteIds = append(toDeleteIds, dbItem.GoodsId)
		}

	}
	//开启事务更新
	return c.svc.CartModel.TransactCtx(ctx, userId, toUpdate, toInsert, toDeleteIds)
}
