package logic

import (
	"context"
	"errors"
	"zeromall/common/constant"
	"zeromall/user/rpc/internal/model"
	"zeromall/user/rpc/internal/svc"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChangeInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeInfoLogic {
	return &ChangeInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeInfoLogic) ChangeInfo(in *userpb.ChangeInfoReq) (*userpb.ChangeInfoResp, error) {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Error(codes.NotFound, constant.UserNotFound)
		}
		l.Logger.Errorf(constant.MysqlFailed, "changeInfo", "select", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	updateMap := make(map[string]any)
	if in.Username != nil {
		updateMap["username"] = *in.Username
	}
	if in.Avatar != nil {
		updateMap["avatar"] = *in.Avatar
	}
	if in.Region != nil {
		updateMap["region"] = *in.Region
	}
	if in.Sex != nil {
		updateMap["sex"] = *in.Sex
	}
	if in.Age != nil {
		updateMap["age"] = *in.Age
	}
	affectedRow, err := l.svcCtx.UserModel.UpdateField(l.ctx, updateMap, in.UserId)
	if err != nil {
		l.Logger.Errorf(constant.MysqlFailed, "changeInfo", "update", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	return &userpb.ChangeInfoResp{
		Ok: affectedRow > 0,
	}, nil
}
