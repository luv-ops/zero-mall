package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"zeromall/common/constant"
	"zeromall/common/mq"

	"zeromall/order/rpc/internal/svc"
	"zeromall/order/rpc/orderPb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Error(codes.Internal, constant.UnmarshalErr)
	}
	if order.Status != 0 {
		return nil, status.Error(codes.Unavailable, "订单状态已变化不可取消")
	}
	if order.ExpireTime.Before(time.Now()) {
		return nil, status.Error(codes.Unavailable, "订单已超时关闭")
	}
	num, err := l.svcCtx.OrderModel.CloseOrderByUser(l.ctx, orderNo)
	if err != nil {
		logx.Info(constant.WhereFailed, "offOrder", err)
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	if num == 0 {
		return &orderPb.OrderOffResp{
			Ok: false,
		}, nil
	}
	//发送普通消息归还库存
	items, err := l.svcCtx.OrderItemModel.FindGIdsByOrderNo(l.ctx, orderNo)
	if err != nil {
		logx.Info(constant.WhereFailed, "offOrder", err)
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
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
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	err = l.svcCtx.Producer.Send(l.ctx, l.svcCtx.Config.RocketMqConf.Topics.TopicReturnStock, data)
	if err != nil {
		logx.Info("send msg error", err)
	}
	return &orderPb.OrderOffResp{
		Ok: true,
	}, nil
}
