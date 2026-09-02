package logic

import (
	"context"
	"encoding/json"
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
	key := constant.UserInfoKey + in.UserId
	jsonStr, err := l.svcCtx.Redis.GetCtx(l.ctx, key)
	var res userpb.UserInfoResp
	if err == nil {
		if jsonStr == constant.RedisEmptyValue {
			return nil, status.Error(codes.NotFound, constant.UserNotFound)
		} else if jsonStr != "" {
			err = json.Unmarshal([]byte(jsonStr), &res)
			if err != nil {
				l.Logger.Errorf(constant.UnmarshalErr, "userInfo", err.Error())
				return nil, status.Error(codes.Internal, constant.MiddlewareError)
			}
			return &res, nil
		}
	}
	//查mysql
	user, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			//缓存空值
			err = l.svcCtx.Redis.SetexCtx(l.ctx, key, constant.RedisEmptyValue, constant.ShortTTL)
			if err != nil {
				l.Logger.Errorf(constant.RedisFailed, "setEx", "userInfo", err.Error())
			}
			return nil, status.Error(codes.NotFound, constant.UserNotFound)
		}
		l.Logger.Errorf(constant.MysqlFailed, "userInfo", "FindOneByUserId", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	res = userpb.UserInfoResp{
		UserId:   in.UserId,
		Username: user.Username,
		Avatar:   user.Avatar.String,
		Phone:    user.Phone,
		Age:      user.Age,
		Sex:      user.Sex,
		Balance:  convert.CentsToYuanStr(user.BalanceCent),
	}
	//缓存
	str, err := json.Marshal(&res)
	if err != nil {
		l.Logger.Errorf(constant.MarshalErr, "userInfo", err.Error())
	}
	err = l.svcCtx.Redis.SetexCtx(l.ctx, key, string(str), constant.LongTTL)
	if err != nil {
		l.Logger.Errorf(constant.RedisFailed, "userInfo", err)
	}
	return &res, nil
}
