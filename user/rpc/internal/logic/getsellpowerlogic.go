package logic

import (
	"context"
	"zeromall/common/constant"

	"zeromall/user/rpc/internal/svc"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetSellPowerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSellPowerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSellPowerLogic {
	return &GetSellPowerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSellPowerLogic) GetSellPower(in *userpb.GetSellPowerReq) (*userpb.GetSellPowerResp, error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserModel.SelectOneByField(l.ctx, in.UserId, "is_seller")
	if err != nil {
		l.Logger.Errorf(constant.MysqlFailed, "sellPower", "select", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	return &userpb.GetSellPowerResp{
		Ok: user.IsSeller != 0,
	}, nil
}
