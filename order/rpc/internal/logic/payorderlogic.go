package logic

import (
	"context"
	"strconv"
	"zeromall/common/constant"
	"zeromall/common/xerr"
	"zeromall/pay/rpc/payPb"

	"zeromall/order/rpc/internal/svc"
	"zeromall/order/rpc/orderPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayOrderLogic) PayOrder(in *orderPb.OrderPayReq) (*orderPb.OrderPayResp, error) {
	// todo: add your logic here and delete this line
	//先获取订单状态
	orderNo, _ := strconv.ParseInt(in.OrderNo, 10, 64)
	order, err := l.svcCtx.OrderModel.FindOneByOrderNo(l.ctx, orderNo)
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "payOrder", err.Error())
		return nil, xerr.Server()
	}
	if order.UserId != in.UserId {
		return nil, xerr.NewCodeError(xerr.PermissionDenied)
	}
	if order.Status != 0 {
		return nil, xerr.NewCodeError(xerr.OrderStatusChange)
	}
	//调用payRpc
	res, err := l.svcCtx.PayRpc.CreatePay(l.ctx, &payPb.CreatePayReq{
		UserId:  in.UserId,
		OrderNo: orderNo,
		PayCent: order.PayCent,
	})
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "PayOrder", err)
		return nil, xerr.FromRpcError(err)
	}
	return &orderPb.OrderPayResp{
		Ok: res.Ok,
	}, nil
}
