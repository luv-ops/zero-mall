package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"zeromall/common/mq"
	"zeromall/stock/rpc/internal/svc"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type StockReturnConsumer struct {
	svc *svc.ServiceContext
}

func NewStockReturnConsumer(svcCtx *svc.ServiceContext) *StockReturnConsumer {
	return &StockReturnConsumer{svc: svcCtx}
}
func (l *StockReturnConsumer) Consume(ctx context.Context, mvs []*rmq_client.MessageView) error {
	for _, msg := range mvs {
		var stockMsg mq.ReturnStockMsg
		if err := json.Unmarshal(msg.GetBody(), &stockMsg); err != nil {
			fmt.Printf("unmarshal stock-insert failed: %v, body=%s\n", err, string(msg.GetBody()))
			continue // 反序列化失败跳过，不重试
		}
		if err := l.stockReturn(ctx, &stockMsg); err != nil {
			return err // 业务失败，让RocketMQ重试
		}
	}
	return nil
}

// 业务逻辑
func (l *StockReturnConsumer) stockReturn(ctx context.Context, stockMsg *mq.ReturnStockMsg) error {
	//TODO 引入redis lua 原子归还预热库存
	num, err := l.svc.StockModel.StockBatchReturn(ctx, stockMsg.List)
	if err != nil {
		return err
	}
	if num == 0 {
		logx.Infof("stockReturn update success but no rows affected")
		return nil
	}
	return nil
}
