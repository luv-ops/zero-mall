package logic

import (
	"context"

	"zeromall/balance/rpc/balancePb"
	"zeromall/balance/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RefundBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundBalanceLogic {
	return &RefundBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundBalanceLogic) RefundBalance(in *balancePb.RefundBalanceReq) (*balancePb.RefundBalanceResp, error) {
	// todo: add your logic here and delete this line
	num, err := l.svcCtx.BalanceModel.RefundBalance(l.ctx, in.PayCent, in.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "退款失败")
	}
	return &balancePb.RefundBalanceResp{
		Ok: num > 0,
	}, nil
}
