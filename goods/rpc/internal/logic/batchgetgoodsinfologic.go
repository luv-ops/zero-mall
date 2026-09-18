package logic

import (
	"context"
	"zeromall/common/convert"
	"zeromall/common/xerr"

	"zeromall/goods/rpc/goodsPb"
	"zeromall/goods/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetGoodsInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetGoodsInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetGoodsInfoLogic {
	return &BatchGetGoodsInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetGoodsInfoLogic) BatchGetGoodsInfo(in *goodsPb.BatchGetGoodsInfoReq) (*goodsPb.BatchGetGoodsInfoResp, error) {
	// todo: add your logic here and delete this line
	//禁止缓存
	var list []*goodsPb.GoodsInfoItem
	resp, err := l.svcCtx.GoodsModel.FindRowsByGoodsId(l.ctx, in.GoodsIds)
	if err != nil {
		return nil, xerr.Server()
	}

	for _, v := range resp {
		if v.Status != 1 {
			return nil, xerr.NewCodeError(xerr.GoodsNotSell)
		}
		list = append(list, &goodsPb.GoodsInfoItem{
			GoodsId:       v.GoodsId,
			Name:          v.Name,
			Cover:         v.Cover,
			Price:         convert.CentsToYuanStr(v.PriceCent),
			OriginalPrice: convert.CentsToYuanStr(v.OriginalPriceCent),
		})
	}
	return &goodsPb.BatchGetGoodsInfoResp{
		List: list,
	}, nil
}
