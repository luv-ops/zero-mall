package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zeromall/common/constant"
	"zeromall/common/mq"
	"zeromall/order/rpc/internal/svc"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type OrderOffConsumer struct {
	svc *svc.ServiceContext
}

func NewOrderOffConsumer(svcCtx *svc.ServiceContext) *OrderOffConsumer {
	return &OrderOffConsumer{svc: svcCtx}
}
func (l *OrderOffConsumer) Consume(ctx context.Context, mvs []*rmq_client.MessageView) error {
	for _, msg := range mvs {
		var orderMsg mq.OrderOffMessage
		if err := json.Unmarshal(msg.GetBody(), &orderMsg); err != nil {
			fmt.Printf("unmarshal cart-change failed: %v, body=%s\n", err, string(msg.GetBody()))
			continue // 反序列化失败跳过，不重试
		}
		if err := l.orderOff(ctx, &orderMsg); err != nil {
			return err // 业务失败，让RocketMQ重试
		}
	}
	return nil
}

// 消费逻辑
func (l *OrderOffConsumer) orderOff(ctx context.Context, msg *mq.OrderOffMessage) error {
	//先查一下expireTime和status
	order, err := l.svc.OrderModel.FindOneByOrderNo(ctx, msg.OrderNo)
	if err != nil {
		return err
	}
	//因为不能消息延迟消息的延迟时间，可能会导致消息还没有到设置的延迟时间就被投递
	if order.Status != int64(0) {
		//订单状态已变化，不再归还库存，直接丢弃消息
		return nil
	}
	if order.ExpireTime.After(time.Now()) {
		//订单还没有过期，也直接丢弃
		return nil
	}

	//先查询订单状态是否为未支付，如果未支付才能取消订单
	num, err := l.svc.OrderModel.CloseOrderTimeOut(ctx, msg.OrderNo)
	if err != nil {
		return err
	}
	//发送消息，归还冻结库存
	if num > 0 {
		//查询该订单的所有商品id
		list, err := l.svc.OrderItemModel.FindGIdsByOrderNo(ctx, msg.OrderNo)
		if err != nil {
			logx.Errorf(constant.WhereFailed, "order_consumer_orderOff err", err)
			return err
		}
		if len(list) == 0 {
			logx.Infof("order_consumer_orderOff查询到0条记录")
			return nil
		}
		var stocksMsg mq.ReturnStockMsg
		for _, item := range list {
			stocksMsg.List = append(stocksMsg.List, &mq.ReturnStockItem{
				GoodsId: item.GoodsId,
				Num:     item.Num,
			})
		}
		data, err := json.Marshal(stocksMsg)
		if err != nil {
			logx.Errorf(constant.MarshalErr, "order_consumer_orderOff err", err)
			return err
		}
		return l.svc.Producer.Send(ctx, l.svc.Config.RocketMqConf.Topics.TopicReturnStock, data)
	}
	return nil
}
