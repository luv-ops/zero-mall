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

type GetOtherInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOtherInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOtherInfoLogic {
	return &GetOtherInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOtherInfoLogic) GetOtherInfo(req *types.OthersBaseInfoReq) (resp *types.OthersBaseInfoResp, err error) {
	// todo: add your logic here and delete this line
	if req.UserId == "" {
		return nil, errors.New("参数不合法")
	}
	res, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &userpb.UserInfoReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &types.OthersBaseInfoResp{
		UserName: res.Username,
		UserId:   res.UserId,
		Avatar:   res.Avatar,
		Age:      res.Age,
		Sex:      res.Sex,
	}, nil
}
