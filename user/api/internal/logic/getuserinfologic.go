// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"zeromall/user/rpc/userpb"

	"zeromall/user/api/internal/svc"
	"zeromall/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line
	//解析token
	userId := l.ctx.Value("userId").(string)
	if userId == "" {
		return nil, errors.New("请输入正确的用户id")
	}
	res, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &userpb.UserInfoReq{UserId: userId})
	if err != nil {
		return nil, err
	}
	return &types.UserInfoResp{
		UserId:   res.UserId,
		UserName: res.Username,
		Avatar:   res.Avatar,
		Phone:    res.Phone,
		Age:      res.Age,
		Sex:      res.Sex,
	}, nil
}
