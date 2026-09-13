package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zeromall/cart/rpc/internal/svc"
	"zeromall/common/constant"
	"zeromall/common/mq"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type CartDelConsumer struct {
	svc *svc.ServiceContext
}

func NewCartDelConsumer(svcCtx *svc.ServiceContext) *CartDelConsumer {
	return &CartDelConsumer{svc: svcCtx}
}
func (l *CartDelConsumer) Consume(ctx context.Context, mvs []*rmq_client.MessageView) error {
	for _, msg := range mvs {
		var cartMsg mq.CartDelMsg
		if err := json.Unmarshal(msg.GetBody(), &cartMsg); err != nil {
			fmt.Printf("unmarshal cart-del failed: %v, body=%s\n", err, string(msg.GetBody()))
			continue // 反序列化失败跳过，不重试
		}
		if err := l.CartDel(ctx, &cartMsg); err != nil {
			return err // 业务失败，让RocketMQ重试
		}
	}
	return nil
}

// 业务逻辑
func (l *CartDelConsumer) CartDel(ctx context.Context, cartMsg *mq.CartDelMsg) error {
	key := constant.CartKey + cartMsg.UserId
	_, err := l.svc.Redis.HdelCtx(ctx, key, cartMsg.GoodsIds...)
	if err != nil {
		return err
	}
	//生产消息
	msg := mq.CartChangeMsg{
		UserId:    cartMsg.UserId,
		TimeStamp: time.Now().Unix(),
	}
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		logx.Errorf("json marshell err %v", err)
		return err
	}
	err = l.svc.Producer.Send(ctx, l.svc.Config.RocketMqConf.Topics.TopicSyncFiling, jsonStr)
	if err != nil {
		logx.Errorf("send msg err %v", err)
	}
	return nil
}
