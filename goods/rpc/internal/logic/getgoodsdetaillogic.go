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

type GetGoodsDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGoodsDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsDetailLogic {
	return &GetGoodsDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGoodsDetailLogic) GetGoodsDetail(in *goodsPb.GoodsDetailReq) (*goodsPb.GoodsDetailResp, error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.GoodsModel.FindOneByGoodsId(l.ctx, in.GoodsId)
	if err != nil {
		l.Logger.Errorf(constant.MysqlFailed, "getGoodsDetail ", "select", err)
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	if res == nil {
		return nil, status.Error(codes.NotFound, constant.GoodsNotFound)
	}

	return &goodsPb.GoodsDetailResp{
		GoodsId:       res.GoodsId,
		Name:          res.Name,
		Cover:         res.Cover,
		Price:         convert.CentsToYuanStr(res.PriceCent),
		OriginalPrice: convert.CentsToYuanStr(res.OriginalPriceCent),
		Stock:         res.Stock,
		Sales:         res.Sales,
		Desc:          res.Desc,
		Status:        res.Status,
	}, nil
}
