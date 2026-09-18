// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"zeromall/common/Regx"
	"zeromall/common/xerr"
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

func (l *AddReceiveAddressLogic) AddReceiveAddress(req *types.AddReceiveAddressReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	//校验手机号
	if !Regx.IsValidPhone(req.ReceivePhone) {
		return nil, xerr.NewCodeError(xerr.PhoneIllegal)
	}
	_, err = l.svcCtx.UserRpc.AddRecAddress(l.ctx, &userpb.AddReceiveAddressReq{
		UserId:       userId,
		ReceiveName:  req.ReceiveName,
		ReceivePhone: req.ReceivePhone,
		AddressId:    req.AddressId,
		Detail:       req.Detail,
		IsDefault:    req.IsDefault,
	})
	if err != nil {
		return nil, xerr.FromRpcError(err)
	}

	return nil, nil
}
