package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"zeromall/balance/rpc/internal/model"
	"zeromall/balance/rpc/internal/svc"
	"zeromall/common/mq"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
)

type BalanceInsertConsumer struct {
	svc *svc.ServiceContext
}

func NewBalanceInsertConsumer(svcCtx *svc.ServiceContext) *BalanceInsertConsumer {
	return &BalanceInsertConsumer{svc: svcCtx}
}
func (l *BalanceInsertConsumer) Consume(ctx context.Context, mvs []*rmq_client.MessageView) error {
	for _, msg := range mvs {
		var balanceMsg mq.InsertBalanceMsg
		if err := json.Unmarshal(msg.GetBody(), &balanceMsg); err != nil {
			fmt.Printf("unmarshal stock-insert failed: %v, body=%s\n", err, string(msg.GetBody()))
			continue // 反序列化失败跳过，不重试
		}
		if err := l.balanceInsert(ctx, &balanceMsg); err != nil {
			return err // 业务失败，让RocketMQ重试
		}
	}
	return nil
}

// 业务逻辑
func (l *BalanceInsertConsumer) balanceInsert(ctx context.Context, msg *mq.InsertBalanceMsg) error {
	balance := model.Balance{
		UserId:        msg.UserId,
		AvailableCent: 1000000,
		Version:       1,
	}
	res, err := l.svc.BalanceModel.Insert(ctx, &balance)
	if err != nil {
		return err
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return errors.New("insert balance failed")
	}
	return nil
}
