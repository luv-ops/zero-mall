package logic

import (
	"context"
	"errors"
	"zeromall/common/constant"
	"zeromall/user/rpc/internal/model"
	"zeromall/user/rpc/internal/svc"
	"zeromall/user/rpc/userpb"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *userpb.RegisterReq) (*userpb.RegisterResp, error) {
	// todo: add your logic here and delete this line
	//校验redis验证码
	codeKey := "captcha:code:" + in.Phone
	code, err := l.svcCtx.Redis.GetCtx(l.ctx, codeKey)
	if err != nil {
		l.Logger.Errorf(constant.RedisFailed, "register", "Get", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	if in.Captcha != code {
		return nil, status.Error(codes.InvalidArgument, constant.CaptchaError)
	}
	//校验成功，删除验证码
	_, _ = l.svcCtx.Redis.DelCtx(l.ctx, codeKey)

	//校验手机号是否注册
	user, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		l.Logger.Errorf(constant.MysqlFailed, "register", "FindOneByPhone", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	if user != nil {
		return nil, status.Error(codes.AlreadyExists, constant.UserExist)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	//组装model
	userid := uuid.NewString()

	newUser := model.User{
		UserId:   userid,
		Username: uuid.NewString()[:8],
		Phone:    in.Phone,
		Password: string(hash),
	}
	_, err = l.svcCtx.UserModel.Insert(l.ctx, &newUser)
	if err != nil {
		l.Logger.Errorf(constant.MysqlFailed, "register", "insert", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	return &userpb.RegisterResp{UserId: userid}, nil
}
