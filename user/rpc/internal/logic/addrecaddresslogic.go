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

type AddRecAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddRecAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRecAddressLogic {
	return &AddRecAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddRecAddressLogic) AddRecAddress(in *userpb.AddReceiveAddressReq) (*userpb.AddReceiveAddressResp, error) {
	// todo: add your logic here and delete this line
	//先根据addressId查地区是否是3级地区
	ok, err := l.svcCtx.AreaModel.AreaLevel3IsExist(l.ctx, in.AddressId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("非3级地区")
	}
	//先查询一次收货地址表，判断是否用户已经有收货地址，如果没有将第一条插入的is_default设置为1
	exist, err := l.svcCtx.RecAddressModel.ExistReceiveAddr(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	if !exist {
		in.IsDefault = 1
	}
	//插入收货地址表
	num, err := l.svcCtx.RecAddressModel.TxRecAddrInsert(l.ctx, &model.UserReceiveAddress{
		UserId:        in.UserId,
		ReceiverName:  in.ReceiveName,
		ReceiverPhone: in.ReceivePhone,
		AddressId:     in.AddressId,
		Detail:        in.Detail,
		IsDefault:     in.IsDefault,
	})
	if err != nil {
		l.Logger.Errorf("AddRecAddress Insert err:%v", err)
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	if in.IsDefault == 1 {
		key := constant.DefaultReceiveArea + in.UserId
		//删除旧的默认地址缓存
		_, err = l.svcCtx.Redis.DelCtx(l.ctx, key)
		if err != nil {
			l.Logger.Errorf("AddRecAddress DelCtx err:%v", err)
		}
	}
	return &userpb.AddReceiveAddressResp{
		Ok: num > 0,
	}, nil
}
