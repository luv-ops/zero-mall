package logic

import (
	"context"

	"zeromall/stock/rpc/internal/svc"
	"zeromall/stock/rpc/stockPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStockLogic {
	return &GetStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetStockLogic) GetStock(in *stockPb.GetStockReq) (*stockPb.GetStockResp, error) {
	// todo: add your logic here and delete this line

	return &stockPb.GetStockResp{}, nil
}
