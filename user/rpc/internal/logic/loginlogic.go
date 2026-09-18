package logic

import (
	"context"
	"errors"
	"zeromall/common/constant"
	"zeromall/common/xerr"
	"zeromall/user/rpc/internal/model"

	"zeromall/user/rpc/internal/svc"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *userpb.LoginReq) (*userpb.LoginResp, error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NewCodeError(xerr.NotFound)
		}
		l.Logger.Errorf(constant.MysqlFailed, "login", "FindOneByPhone", err.Error())
		return nil, xerr.Server()
	}
	//判断密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		return nil, xerr.NewCodeError(xerr.PasswordErr)
	}
	//密码正确
	return &userpb.LoginResp{
		UserId: user.UserId,
	}, nil
}
