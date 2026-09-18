// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"zeromall/common/Regx"
	"zeromall/common/xerr"
	"zeromall/user/rpc/userpb"

	"zeromall/user/api/internal/svc"
	"zeromall/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaLogic {
	return &CaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaLogic) Captcha(req *types.CaptchaReq) (resp *types.CaptchaResp, err error) {
	// todo: add your logic here and delete this line
	//验证手机号
	if !Regx.IsValidPhone(req.Phone) {
		return nil, xerr.NewCodeError(xerr.PhoneIllegal)
	}
	res, err := l.svcCtx.UserRpc.Captcha(l.ctx, &userpb.CaptchaReq{
		Phone: req.Phone,
	})
	if err != nil {
		return nil, xerr.FromRpcError(err)
	}
	return &types.CaptchaResp{Captcha: res.Captcha}, nil
}
