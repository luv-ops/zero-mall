package logic

import (
	"context"
	"errors"
	"zeromall/common/constant"
	"zeromall/common/convert"
	"zeromall/user/rpc/internal/model"

	"zeromall/user/rpc/internal/svc"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *userpb.UserInfoReq) (*userpb.UserInfoResp, error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Error(codes.NotFound, constant.UserNotFound)
		}
		l.Logger.Errorf(constant.MysqlFailed, "userInfo", "FindOneByUserId", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}

	return &userpb.UserInfoResp{
		UserId:   in.UserId,
		Username: user.Username,
		Avatar:   user.Avatar.String,
		Phone:    user.Phone,
		Age:      user.Age,
		Sex:      user.Sex,
		Balance:  convert.CentsToYuanStr(user.BalanceCent),
	}, nil
}
