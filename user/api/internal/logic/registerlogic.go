// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"zeromall/common/Regx"
	"zeromall/common/jwt"
	"zeromall/user/api/internal/svc"
	"zeromall/user/api/internal/types"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	if !Regx.IsValidPhone(req.Phone) {
		return nil, errors.New("手机号非法")
	}
	rpcRes, err := l.svcCtx.UserRpc.Register(l.ctx, &userpb.RegisterReq{
		Phone:    req.Phone,
		Password: req.Password,
		Captcha:  req.Captcha,
	})
	if err != nil {
		return nil, err
	}
	token, err := jwt.GenerateToken(rpcRes.UserId, l.svcCtx.Config.Auth.AccessExpire, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, err
	}
	return &types.LoginResp{Token: token}, err
}
