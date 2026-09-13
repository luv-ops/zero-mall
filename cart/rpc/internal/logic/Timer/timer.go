package Timer

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

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type Timer struct {
	svc *svc.ServiceContext
}

func NewTimer(svc *svc.ServiceContext) *Timer {
	return &Timer{svc: svc}
}
func (t *Timer) StartTicker(ctx context.Context) {
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
				res, err := t.svc.Redis.EvalCtx(ctx, `return redis.call('SPOP', KEYS[1], 50) `, keys, nil)
				if err != nil {
					logc.Infof(ctx, "eval spop err: %v", err)
					continue
				}
				var userIds []string
				if arr, ok := res.([]interface{}); ok {
					for _, v := range arr {
						if str, ok2 := v.(string); ok2 {
							//消费时已经反序列化了
							userIds = append(userIds, str)
						}
					}
				}
				// 循环执行归档同步MySQL
				for _, uid := range userIds {
					_ = t.syncCartToMysql(ctx, uid)
				}
			}
		}
	}()
}
func (t *Timer) syncCartToMysql(ctx context.Context, userId string) error {
	logc.Infof(ctx, "syncCartToMysql userId=%s", userId)
	// 1. HGETALL cart:{userId}
	key := constant.CartKey + userId
	redisMap, err := t.svc.Redis.HgetallCtx(ctx, key)
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

		redisItemMap[goodsId] = &model.Cart{
			UserId:          userId,
			GoodsId:         goodsId,
			Num:             item.Num,
			Selected:        item.Selected,
			PriceCent:       convert.YuanStrToCents(item.Price),
			Name:            item.Name,
			Cover:           item.Cover,
			OriginPriceCent: convert.YuanStrToCents(item.OriginPrice),
		}
	}
	//获取mysql数据也
	res, err := t.svc.CartModel.FindCartsByUserId(ctx, userId)
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
	return t.svc.CartModel.TransactCtx(ctx, userId, toUpdate, toInsert, toDeleteIds)
}
