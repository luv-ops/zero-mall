package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"zeromall/common/constant"
	"zeromall/common/mq"
	"zeromall/stock/rpc/internal/svc"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/redis/go-redis/v9"
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
	//先操作db，再redis
	num, err := l.svc.StockModel.StockBatchReturn(ctx, stockMsg.List)
	if err != nil {
		return err
	}
	if num == 0 {
		logx.Infof("stockReturn update success but no rows affected")
		return nil
	}
	//num>0
	//TODO 引入redis luaScript 原子归还预热库存,pipeline减少网络往返
	var cmds []*redis.Cmd
	err = l.svc.Redis.PipelinedCtx(ctx, func(pipeliner redis.Pipeliner) error {
		for _, item := range stockMsg.List {
			stockKey := constant.StockGoodsKey + item.GoodsId
			cmd := pipeliner.EvalSha(ctx, l.svc.StockReturnSha, []string{stockKey}, item.Num)
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
			logx.Errorf("redis解冻脚本执行失败 goodsId=%s, err=%v", item.GoodsId, err)
			continue
		}
		ret, ok := res.(int64)
		if !ok || ret != 1 {
			logx.Errorf("redis解冻校验失败，冻结库存不足 goodsId=%s", item.GoodsId)
		}
	}
	return nil
}
