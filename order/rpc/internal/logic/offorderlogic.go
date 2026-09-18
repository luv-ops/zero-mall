package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"zeromall/common/constant"
	"zeromall/common/mq"
	"zeromall/common/xerr"

	"zeromall/order/rpc/internal/svc"
	"zeromall/order/rpc/orderPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OffOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOffOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OffOrderLogic {
	return &OffOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OffOrderLogic) OffOrder(in *orderPb.OrderOffReq) (*orderPb.OrderOffResp, error) {
	// todo: add your logic here and delete this line
	orderNo, _ := strconv.ParseInt(in.OrderNo, 10, 64)
	order, err := l.svcCtx.OrderModel.FindOneByOrderNo(l.ctx, orderNo)
	if err != nil {
		logx.Info(constant.WhereFailed, "offOrder", err)
		return nil, xerr.Server()
	}
	if order.UserId != in.UserId {
		return nil, xerr.NewCodeError(xerr.PermissionDenied)
	}
	if order.Status != 0 {
		return nil, xerr.NewCodeError(xerr.OrderStatusChange)
	}
	if order.ExpireTime.Before(time.Now()) {
		return nil, xerr.NewCodeError(xerr.OrderNotCancel)
	}
	num, err := l.svcCtx.OrderModel.CloseOrderByUser(l.ctx, orderNo)
	if err != nil {
		logx.Info(constant.WhereFailed, "offOrder", err)
		return nil, xerr.Server()
	}
	if num == 0 {
		return nil, xerr.NewCodeError(xerr.OrderOffErr)
	}
	//发送普通消息归还库存
	items, err := l.svcCtx.OrderItemModel.FindGIdsByOrderNo(l.ctx, orderNo)
	if err != nil {
		logx.Info(constant.WhereFailed, "offOrder", err)
		return nil, xerr.Server()
	}
	var list []*mq.ReturnStockItem
	for _, item := range items {
		list = append(list, &mq.ReturnStockItem{
			GoodsId: item.GoodsId,
			Num:     item.Num,
		})
	}
	var msg mq.ReturnStockMsg
	msg.List = list
	data, err := json.Marshal(msg)
	if err != nil {
		l.Logger.Errorf(constant.MarshalErr, "OffOrder", err)
		return nil, xerr.Server()
	}
	err = l.svcCtx.Producer.Send(l.ctx, l.svcCtx.Config.RocketMqConf.Topics.TopicReturnStock, data)
	if err != nil {
		logx.Info("send msg error", err)
		return nil, xerr.Server()
	}
	return &orderPb.OrderOffResp{
		Ok: true,
	}, nil
}
