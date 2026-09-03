// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"zeromall/common/Regx"
	"zeromall/common/constant"
	"zeromall/user/rpc/userpb"

	"zeromall/user/api/internal/svc"
	"zeromall/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddReceiveAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddReceiveAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddReceiveAddressLogic {
	return &AddReceiveAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddReceiveAddressLogic) AddReceiveAddress(req *types.AddReceiveAddressReq) error {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	//校验手机号
	if !Regx.IsValidPhone(req.ReceivePhone) {
		return errors.New(constant.PhoneIllegal)
	}
	resp, err := l.svcCtx.UserRpc.AddRecAddress(l.ctx, &userpb.AddReceiveAddressReq{
		UserId:       userId,
		ReceiveName:  req.ReceiveName,
		ReceivePhone: req.ReceivePhone,
		AddressId:    req.AddressId,
		Detail:       req.Detail,
		IsDefault:    req.IsDefault,
	})
	if err != nil {
		l.Logger.Errorf("AddReceiveAddress err:%v", err)
		return err
	}
	if !resp.Ok {
		return errors.New("新增收货地址失败")
	}
	return nil
}
