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

type ChangeInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeInfoLogic {
	return &ChangeInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeInfoLogic) ChangeInfo(req *types.ChangeInfoReq) error {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	res, err := l.svcCtx.UserRpc.ChangeInfo(l.ctx, &userpb.ChangeInfoReq{
		UserId:   userId,
		Username: req.UserName,
		Avatar:   req.Avatar,
		Age:      req.Age,
		Sex:      req.Sex,
		Region:   req.Region,
	})
	if err != nil {
		return err
	}
	if res.Ok == false {
		return errors.New("修改失败")
	}

	return nil
}
