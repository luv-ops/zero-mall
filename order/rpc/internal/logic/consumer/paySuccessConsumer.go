package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"zeromall/common/constant"
	"zeromall/common/mq"
	"zeromall/order/rpc/internal/svc"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type PaySuccessConsumer struct {
	svc *svc.ServiceContext
}

func NewPaySuccessConsumer(svcCtx *svc.ServiceContext) *PaySuccessConsumer {
	return &PaySuccessConsumer{svc: svcCtx}
}
func (l *PaySuccessConsumer) Consume(ctx context.Context, mvs []*rmq_client.MessageView) error {
	for _, msg := range mvs {
		var payMsg mq.PaySuccessMsg
		if err := json.Unmarshal(msg.GetBody(), &payMsg); err != nil {
			fmt.Printf("unmarshal paysuccess failed: %v, body=%s\n", err, string(msg.GetBody()))
			continue // 反序列化失败跳过，不重试
		}
		if err := l.paySuccess(ctx, &payMsg); err != nil {
			return err // 业务失败，让RocketMQ重试
		}
	}
	return nil
}

// 支付成功逻辑
func (l *PaySuccessConsumer) paySuccess(ctx context.Context, msg *mq.PaySuccessMsg) error {
	//修改订单状态
	num, err := l.svc.OrderModel.UpdateStatus(ctx, msg.OrderNo, 1)
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("修改订单状态失败")
	}
	//查order_item获取数量
	orderItems, err := l.svc.OrderItemModel.FindGIdsByOrderNo(ctx, msg.OrderNo)
	if err != nil {
		return err
	}
	if len(orderItems) == 0 {
		return errors.New("此订单不存在购买商品")
	}
	var items []*mq.DeductStockItem
	for _, v := range orderItems {
		items = append(items, &mq.DeductStockItem{
			OrderNo: msg.OrderNo,
			GoodsId: v.GoodsId,
			Num:     v.Num,
		})
	}
	var deductMsg mq.DeductStockMsg
	deductMsg.List = items
	data, err := json.Marshal(deductMsg)
	if err != nil {
		logx.Errorf(constant.MarshalErr, "paySuccessConsumer err:%v", err)
		return err
	}
	//发送消息真实扣减库存
	err = l.svc.Producer.Send(ctx, l.svc.Config.RocketMqConf.Topics.TopicDeductStock, data)
	if err != nil {
		logx.Errorf(constant.WhereFailed, "send paySuccessMsg err:%v", err)
		return err
	}
	return nil
}
