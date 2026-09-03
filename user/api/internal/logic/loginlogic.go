// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"zeromall/common/Regx"
	"zeromall/common/constant"
	"zeromall/common/jwt"
	"zeromall/user/api/internal/svc"
	"zeromall/user/api/internal/types"
	"zeromall/user/rpc/userpb"

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
		return nil, errors.New(constant.PhoneIllegal)
	}
	res, err := l.svcCtx.UserRpc.Login(l.ctx, &userpb.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	token, err := jwt.GenerateToken(res.UserId, l.svcCtx.Config.Auth.AccessExpire, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, err
	}
	return &types.LoginResp{Token: token}, err
}
