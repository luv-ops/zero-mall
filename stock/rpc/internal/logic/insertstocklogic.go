package logic

import (
	"context"
	"zeromall/stock/rpc/internal/svc"
	"zeromall/stock/rpc/stockPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type InsertStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInsertStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InsertStockLogic {
	return &InsertStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InsertStockLogic) InsertStock(in *stockPb.InsertStockReq) (*stockPb.InsertStockResp, error) {
	// todo: add your logic here and delete this line

	return &stockPb.InsertStockResp{}, nil
}
