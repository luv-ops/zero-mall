package logic

import (
	"context"
	"fmt"
	"math/rand"
	"zeromall/common/constant"

	"zeromall/user/rpc/internal/svc"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/core/logx"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaLogic {
	return &CaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CaptchaLogic) Captcha(in *userpb.CaptchaReq) (*userpb.CaptchaResp, error) {
	// todo: add your logic here and delete this line

	limitKey := "captcha:limit:" + in.Phone //验证码限流
	codeKey := "captcha:code:" + in.Phone   //验证码本身
	//尝试创建限流验证码key
	ok, err := l.svcCtx.Redis.SetnxExCtx(l.ctx, limitKey, "1", 60)
	if err != nil {
		l.Logger.Errorf(constant.RedisFailed, "captcha limit", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	if !ok {
		return nil, status.Error(codes.ResourceExhausted, constant.MsgSmsTooFrequently)
	}
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	err = l.svcCtx.Redis.SetexCtx(l.ctx, codeKey, code, 300)
	if err != nil {
		l.Logger.Errorf(constant.RedisFailed, "captcha code", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	return &userpb.CaptchaResp{
		Captcha: code,
	}, nil
}
