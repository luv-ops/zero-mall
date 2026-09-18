// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"zeromall/common/xerr"
	"zeromall/order/rpc/orderPb"

	"zeromall/order/api/internal/svc"
	"zeromall/order/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.CreateOrderResp, err error) {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	res, err := l.svcCtx.OrderRpc.CreateOrder(l.ctx, &orderPb.CreateOrderReq{
		UserId:   userId,
		GoodsIds: req.GoodsIds,
	})
	if err != nil {
		return nil, xerr.FromRpcError(err)
	}
	return &types.CreateOrderResp{
		OrderNo: res.OrderNo,
	}, nil
}
