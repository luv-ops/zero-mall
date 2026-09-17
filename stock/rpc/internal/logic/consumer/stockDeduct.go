package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"zeromall/common/constant"
	"zeromall/common/mq"
	"zeromall/stock/rpc/internal/svc"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type StockDeductConsumer struct {
	svc *svc.ServiceContext
}

func NewStockDeductConsumer(svcCtx *svc.ServiceContext) *StockDeductConsumer {
	return &StockDeductConsumer{svc: svcCtx}
}
func (l *StockDeductConsumer) Consume(ctx context.Context, mvs []*rmq_client.MessageView) error {
	for _, msg := range mvs {
		var stockMsg mq.DeductStockMsg
		if err := json.Unmarshal(msg.GetBody(), &stockMsg); err != nil {
			fmt.Printf("unmarshal stock-deduct failed: %v, body=%s\n", err, string(msg.GetBody()))
			continue // 反序列化失败跳过，不重试
		}
		if err := l.stockDeduct(ctx, &stockMsg); err != nil {
			return err // 业务失败，让RocketMQ重试
		}
	}
	return nil
}

// 业务逻辑
func (l *StockDeductConsumer) stockDeduct(ctx context.Context, stockMsg *mq.DeductStockMsg) error {
	//先扣减mysql库存，再修改redis
	num, err := l.svc.StockModel.BatchDeductStock(ctx, stockMsg.List)
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("真实扣减库存失败，没有行被修改")
	}
	//修改redis
	var cmds []*redis.Cmd
	err = l.svc.Redis.PipelinedCtx(ctx, func(pipeliner redis.Pipeliner) error {
		for _, v := range stockMsg.List {
			key := constant.StockGoodsKey + v.GoodsId
			cmd := pipeliner.EvalSha(ctx, l.svc.StockDeductSha, []string{key}, v.Num)
			cmds = append(cmds, cmd)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for i, cmd := range cmds {
		res, err := cmd.Result()
		item := stockMsg.List[i]
		if err != nil {
			logx.Errorf("redis真实扣减库存脚本执行失败 goodsId=%s, err=%v", item.GoodsId, err)
			continue
		}
		ret, ok := res.(int64)
		if !ok || ret != 1 {
			logx.Errorf("redis真实扣减库存校验失败，冻结库存不足 goodsId=%s", item.GoodsId)
		}
	}
	return nil
}
