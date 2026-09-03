// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"zeromall/user/rpc/userpb"

	"zeromall/user/api/internal/svc"
	"zeromall/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReceiveAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetReceiveAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReceiveAddressLogic {
	return &GetReceiveAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReceiveAddressLogic) GetReceiveAddress() (resp *types.GetReceiveAddressResp, err error) {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	res, err := l.svcCtx.UserRpc.GetReceiveAddress(l.ctx, &userpb.GetReceiveAddressReq{UserId: userId})
	if err != nil {
		return nil, err
	}
	var list []*types.AddressItem
	for _, v := range res.List {
		list = append(list, &types.AddressItem{
			ReceiveName:      v.ReceiveName,
			ReceivePhone:     v.ReceivePhone,
			AddressMergeName: v.AddressMergeName,
			IsDefault:        v.IsDefault,
			Detail:           v.Detail,
		})
	}
	return &types.GetReceiveAddressResp{
		List: list,
	}, nil
}
