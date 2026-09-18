package logic

import (
	"context"
	"encoding/json"
	"errors"
	"zeromall/common/constant"
	"zeromall/common/mq"
	"zeromall/common/xerr"
	"zeromall/user/rpc/internal/model"
	"zeromall/user/rpc/internal/svc"
	"zeromall/user/rpc/userpb"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
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
		return nil, xerr.Server()
	}
	if in.Captcha != code {
		return nil, xerr.NewCodeError(xerr.CaptchaMistake)
	}
	//校验成功，删除验证码
	_, _ = l.svcCtx.Redis.DelCtx(l.ctx, codeKey)

	//校验手机号是否注册
	user, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		l.Logger.Errorf(constant.MysqlFailed, "register", "FindOneByPhone", err.Error())
		return nil, xerr.Server()
	}
	if user != nil {
		return nil, xerr.NewCodeError(xerr.UserHasExist)
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
	res, err := l.svcCtx.UserModel.Insert(l.ctx, &newUser)
	if err != nil {
		l.Logger.Errorf(constant.MysqlFailed, "register", "insert", err.Error())
		return nil, xerr.Server()
	}
	num, _ := res.RowsAffected()
	if num == 0 {
		return nil, xerr.NewCodeError(xerr.AddUserErr)
	}
	//发送一条消息，插入余额
	var msg mq.InsertBalanceMsg
	msg.UserId = userid
	data, err := json.Marshal(msg)
	if err != nil {
		l.Logger.Errorf(constant.MarshalErr, "register", "Marshal", err.Error())
		return nil, xerr.Server()
	}
	err = l.svcCtx.Producer.Send(l.ctx, l.svcCtx.Config.RocketMqConf.Topics.TopicInsertBalance, data)
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "register", "Send", err.Error())
		return nil, xerr.Server()
	}
	return &userpb.RegisterResp{UserId: userid}, nil
}
