package logic

import (
	"context"
	"encoding/json"
	"zeromall/common/constant"

	"zeromall/user/rpc/internal/svc"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDefaultAreaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDefaultAreaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDefaultAreaLogic {
	return &GetDefaultAreaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDefaultAreaLogic) GetDefaultArea(in *userpb.GetDefaultAreaReq) (*userpb.GetDefaultAreaResp, error) {
	// todo: add your logic here and delete this line
	//读取缓存
	key := constant.DefaultReceiveArea + in.UserId
	jsonStr, err := l.svcCtx.Redis.GetCtx(l.ctx, key)
	if err != nil {
		l.Logger.Errorf("GetDefaultArea Redis err:%v", err)
	}
	var item userpb.GetDefaultAreaResp
	if jsonStr != "" {
		err = json.Unmarshal([]byte(jsonStr), &item)
		if err != nil {
			l.Logger.Errorf(constant.UnmarshalErr, "getDefaultArea", err)
		}
		return &item, nil
	}
	//查询数据库
	res, err := l.svcCtx.RecAddressModel.FindOneByUIdWithDefault(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Errorf("getDefaultArea FindOneByUId err:%v", err)
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	item = userpb.GetDefaultAreaResp{
		ReceiverName:    res.ReceiverName,
		ReceiverPhone:   res.ReceiverPhone,
		ReceiverAddress: res.AddressMergeName,
	}
	//建立缓存
	bytes, err := json.Marshal(&item)
	if err != nil {
		l.Logger.Errorf(constant.MarshalErr, "getDefaultArea", err)
	}
	err = l.svcCtx.Redis.SetCtx(l.ctx, key, string(bytes))
	if err != nil {
		l.Logger.Errorf(constant.RedisFailed, "getDefaultArea", err)
	}
	return &item, nil
}
