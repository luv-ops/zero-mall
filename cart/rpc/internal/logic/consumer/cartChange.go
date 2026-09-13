package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"zeromall/cart/rpc/internal/svc"
	"zeromall/common/constant"
	"zeromall/common/mq"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type CartChangeConsumer struct {
	svc *svc.ServiceContext
}

func NewCartChangeConsumer(svcCtx *svc.ServiceContext) *CartChangeConsumer {
	return &CartChangeConsumer{svc: svcCtx}
}
func (l *CartChangeConsumer) Consume(ctx context.Context, mvs []*rmq_client.MessageView) error {
	for _, msg := range mvs {
		logx.Info("收到消息:%v", msg)
		var cartMsg mq.CartChangeMsg
		if err := json.Unmarshal(msg.GetBody(), &cartMsg); err != nil {
			fmt.Printf("unmarshal cart-change failed: %v, body=%s\n", err, string(msg.GetBody()))
			continue // 反序列化失败跳过，不重试
		}
		if err := l.CartChange(ctx, &cartMsg); err != nil {
			return err // 业务失败，让RocketMQ重试
		}
	}
	return nil
}

// 业务逻辑
func (l *CartChangeConsumer) CartChange(ctx context.Context, cartMsg *mq.CartChangeMsg) error {
	_, err := l.svc.Redis.SaddCtx(ctx, constant.PendingSyncCartKey, cartMsg.UserId)
	return err
}
