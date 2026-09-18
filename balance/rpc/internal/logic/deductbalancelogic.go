package logic

import (
	"context"
	"zeromall/common/constant"
	"zeromall/common/xerr"

	"zeromall/balance/rpc/balancePb"
	"zeromall/balance/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeductBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductBalanceLogic {
	return &DeductBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductBalanceLogic) DeductBalance(in *balancePb.DeductBalanceReq) (*balancePb.DeductBalanceResp, error) {
	// todo: add your logic here and delete this line
	//先查询余额和version，
	balance, err := l.svcCtx.BalanceModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Errorf(constant.MysqlFailed, "select", "deduct balance model error", err.Error())
		return nil, xerr.Server()
	}
	if balance.AvailableCent < in.PayCent {
		return nil, xerr.NewCodeError(xerr.BalanceNotEnough)
	}
	num, err := l.svcCtx.BalanceModel.DeductBalance(l.ctx, in.PayCent, balance.Version, in.UserId)
	if err != nil {
		l.Logger.Errorf(constant.MysqlFailed, "deduct balance model error", err.Error())
		return nil, xerr.Server()
	}
	return &balancePb.DeductBalanceResp{
		Ok: num > 0,
	}, nil
}
