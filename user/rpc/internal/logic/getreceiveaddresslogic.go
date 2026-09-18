package logic

import (
	"context"
	"errors"
	"zeromall/common/xerr"
	"zeromall/user/rpc/internal/model"

	"zeromall/user/rpc/internal/svc"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReceiveAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetReceiveAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReceiveAddressLogic {
	return &GetReceiveAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetReceiveAddressLogic) GetReceiveAddress(in *userpb.GetReceiveAddressReq) (*userpb.GetReceiveAddressResp, error) {
	// todo: add your logic here and delete this line
	list, err := l.svcCtx.RecAddressModel.FindAddressWithArea(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NewCodeError(xerr.NotFound)
		}
		return nil, xerr.Server()
	}
	var itemList []*userpb.AddressItem
	for _, v := range list {
		itemList = append(itemList, &userpb.AddressItem{
			ReceivePhone:     v.ReceiverPhone,
			ReceiveName:      v.ReceiverName,
			AddressMergeName: v.AddressMergeName,
			IsDefault:        v.IsDefault,
			Detail:           v.Detail,
		})
	}
	return &userpb.GetReceiveAddressResp{
		List: itemList,
	}, nil
}
