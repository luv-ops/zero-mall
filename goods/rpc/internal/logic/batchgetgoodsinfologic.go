package logic

import (
	"context"
	"zeromall/common/constant"
	"zeromall/common/convert"

	"zeromall/goods/rpc/goodsPb"
	"zeromall/goods/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	resp, err := l.svcCtx.GoodsModel.FindRowsByGoodsId(l.ctx, in.GoodsIds)
	if err != nil {
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	var list []*goodsPb.GoodsInfoItem
	for _, v := range resp {
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
