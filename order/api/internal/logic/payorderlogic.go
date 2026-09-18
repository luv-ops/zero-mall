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

type PayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayOrderLogic) PayOrder(req *types.PayOrderReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	_, err = l.svcCtx.OrderRpc.PayOrder(l.ctx, &orderPb.OrderPayReq{
		UserId:  userId,
		OrderNo: req.OrderNo,
	})
	if err != nil {
		return nil, xerr.FromRpcError(err)
	}
	return nil, nil
}
