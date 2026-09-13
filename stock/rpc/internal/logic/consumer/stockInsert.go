package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"zeromall/common/mq"
	"zeromall/stock/rpc/internal/model"
	"zeromall/stock/rpc/internal/svc"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
)

type StockInsertConsumer struct {
	svc *svc.ServiceContext
}

func NewStockInsertConsumer(svcCtx *svc.ServiceContext) *StockInsertConsumer {
	return &StockInsertConsumer{svc: svcCtx}
}
func (l *StockInsertConsumer) Consume(ctx context.Context, mvs []*rmq_client.MessageView) error {
	for _, msg := range mvs {
		var stockMsg mq.InsertStockMsg
		if err := json.Unmarshal(msg.GetBody(), &stockMsg); err != nil {
			fmt.Printf("unmarshal stock-insert failed: %v, body=%s\n", err, string(msg.GetBody()))
			continue // 反序列化失败跳过，不重试
		}
		if err := l.stockInsert(ctx, &stockMsg); err != nil {
			return err // 业务失败，让RocketMQ重试
		}
	}
	return nil
}

// 业务逻辑
func (l *StockInsertConsumer) stockInsert(ctx context.Context, stockMsg *mq.InsertStockMsg) error {
	stock := model.Stock{
		GoodsId:        stockMsg.GoodsId,
		AvailableStock: stockMsg.Stock,
	}
	_, err := l.svc.StockModel.Insert(ctx, &stock)
	return err
}
