package logic

import (
	"context"
	"zeromall/common/constant"
	"zeromall/common/xerr"

	"zeromall/stock/rpc/internal/svc"
	"zeromall/stock/rpc/stockPb"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type PreHeatStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPreHeatStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreHeatStockLogic {
	return &PreHeatStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PreHeatStockLogic) PreHeatStock(in *stockPb.PreHeatStockReq) (*stockPb.PreHeatStockResp, error) {
	// todo: add your logic here and delete this line
	//查db
	stocks, err := l.svcCtx.StockModel.GetStockByGoodsIds(l.ctx, in.GoodsId)
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "preheatStock ", err)
		return nil, xerr.Server()
	}
	var cmds []*redis.Cmd
	err = l.svcCtx.Redis.PipelinedCtx(l.ctx, func(pipeliner redis.Pipeliner) error {
		for _, v := range stocks {
			key := constant.StockGoodsKey + v.GoodsId
			cmd := pipeliner.EvalSha(l.ctx, l.svcCtx.PreHeatStockSha, []string{key}, v.AvailableStock, v.FrozenStock)
			cmds = append(cmds, cmd)
		}
		return nil
	})
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "pipeline err in preheat", err)
		return nil, xerr.Server()
	}
	ok := true
	for i, cmd := range cmds {
		res, err := cmd.Result()
		item := stocks[i]
		if err != nil {
			ok = false
			logx.Errorf("redis预热脚本执行失败 goodsId=%s, err=%v", item.GoodsId, err)
			continue
		}
		ret, ok := res.(int64)
		if !ok || ret != 1 {
			ok = false
			logx.Errorf("redis预热校验失败，冻结库存不足 goodsId=%s", item.GoodsId)
		}
	}
	return &stockPb.PreHeatStockResp{
		Ok: ok,
	}, nil
}
