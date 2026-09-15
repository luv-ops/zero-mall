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

type StockFrozenConsumer struct {
	svc *svc.ServiceContext
}

func NewStockFrozenConsumer(svcCtx *svc.ServiceContext) *StockFrozenConsumer {
	return &StockFrozenConsumer{svc: svcCtx}
}
func (l *StockFrozenConsumer) Consume(ctx context.Context, mvs []*rmq_client.MessageView) error {
	for _, msg := range mvs {
		var stockMsg mq.FrozenStockMsg
		if err := json.Unmarshal(msg.GetBody(), &stockMsg); err != nil {
			fmt.Printf("unmarshal stock-return failed: %v, body=%s\n", err, string(msg.GetBody()))
			continue // 反序列化失败跳过，不重试
		}
		if err := l.stockFrozen(ctx, &stockMsg); err != nil {
			return err // 业务失败，让RocketMQ重试
		}
	}
	return nil
}

// 业务逻辑
func (l *StockFrozenConsumer) stockFrozen(ctx context.Context, stockMsg *mq.FrozenStockMsg) error {
	//直接操作db，redis已经在创建订单时预冻结
	num, err := l.svc.StockModel.BatchFrozenStock(ctx, stockMsg.List)
	if err != nil {
		return err
	}

	if num == 0 {
		logx.Infof("stockFrozen update success but no rows affected")
		return nil
	}
	return nil
}
