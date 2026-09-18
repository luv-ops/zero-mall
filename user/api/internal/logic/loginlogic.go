// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"zeromall/common/Regx"
	"zeromall/common/jwt"
	"zeromall/common/xerr"
	"zeromall/user/rpc/userpb"

	"zeromall/user/api/internal/svc"
	"zeromall/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	//验证手机号
	if !Regx.IsValidPhone(req.Phone) {
		return nil, xerr.NewCodeError(xerr.PhoneIllegal)
	}
	res, err := l.svcCtx.UserRpc.Login(l.ctx, &userpb.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, xerr.FromRpcError(err)
	}
	token, err := jwt.GenerateToken(res.UserId, l.svcCtx.Config.Auth.AccessExpire, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, xerr.Server()
	}
	return &types.LoginResp{Token: token}, nil
}
