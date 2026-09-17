// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"zeromall/order/rpc/orderPb"

	"zeromall/order/api/internal/svc"
	"zeromall/order/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OffOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOffOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OffOrderLogic {
	return &OffOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OffOrderLogic) OffOrder(req *types.OrderOffReq) error {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	res, err := l.svcCtx.OrderRpc.OffOrder(l.ctx, &orderPb.OrderOffReq{
		OrderNo: req.OrderNo,
		UserId:  userId,
	})
	if err != nil {
		return err
	}
	if res.Ok == false {
		return errors.New("取消订单失败")
	}
	return nil
}
