package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"
	"zeromall/cart/rpc/cartPb"
	"zeromall/common/constant"
	"zeromall/common/mq"

	"zeromall/goods/rpc/goodsPb"
	"zeromall/order/rpc/internal/model"
	"zeromall/order/rpc/internal/svc"
	"zeromall/order/rpc/orderPb"

	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *orderPb.CreateOrderReq) (*orderPb.CreateOrderResp, error) {
	// todo: add your logic here and delete this line
	//调用cartRpc主要是获取对应商品数量
	cartResp, err := l.svcCtx.CartRpc.BatchGetCart(l.ctx, &cartPb.BatchGetCartReq{
		UserId:   in.UserId,
		GoodsIds: in.GoodsIds,
	})
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "previewOrder", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	cartMap := make(map[string]*cartPb.PreviewItemVO)
	for _, v := range cartResp.ItemList {
		cartMap[v.GoodsId] = v
	}
	//计算金额,先通过goodsRpc拿到最新金额
	goodsRes, err := l.svcCtx.GoodsRpc.BatchGetGoodsInfo(l.ctx, &goodsPb.BatchGetGoodsInfoReq{
		GoodsIds: in.GoodsIds,
	})
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "previewOrder", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	goodsMap := make(map[string]*goodsPb.GoodsInfoItem)
	for _, v := range goodsRes.List {
		goodsMap[v.GoodsId] = v
	}

	var totalAmount, payAmount = decimal.Zero, decimal.Zero
	var (
		orderItemData []*model.OrderItem
	)
	//生成雪花id
	snowId := l.svcCtx.Snow.Generate().Int64()
	// map组装
	for goodsId, cartItem := range cartMap {
		goodsItem, _ := goodsMap[goodsId]

		originPrice, _ := decimal.NewFromString(goodsItem.OriginalPrice)
		price, _ := decimal.NewFromString(goodsItem.Price)
		num := decimal.NewFromInt(cartItem.GoodsNum)
		totalAmount = totalAmount.Add(originPrice.Mul(num))
		mulTotal := price.Mul(num)
		payAmount = payAmount.Add(mulTotal)
		orderItemData = append(orderItemData, &model.OrderItem{
			OrderNo:    snowId,
			GoodsId:    goodsId,
			GoodsName:  goodsItem.Name,
			GoodsCover: goodsItem.Cover,
			PriceCent:  price.Round(2).Mul(decimal.NewFromInt(100)).IntPart(),
			TotalCent:  payAmount.Round(2).Mul(decimal.NewFromInt(100)).IntPart(),
			Num:        cartItem.GoodsNum,
		})
	}

	orderData := &model.Order{
		Id:         snowId,
		UserId:     in.UserId,
		PayCent:    payAmount.Round(2).Mul(decimal.NewFromInt(100)).IntPart(),
		TotalCent:  totalAmount.Round(2).Mul(decimal.NewFromInt(100)).IntPart(),
		ExpireTime: time.Now().Add(30 * time.Minute),
		Remark:     sql.NullString{String: *in.Remark, Valid: true},
	}
	//插入两张表 order,order_item TODO 分布式事务
	sum, err := l.svcCtx.OrderModel.TxInsert(l.ctx, orderData, orderItemData)
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "previewOrder", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	if sum <= 0 {
		return nil, status.Error(codes.Internal, "创建订单失败")
	}
	//创建订单成功后删除购物车,发送消息异步删除
	cartMsg := mq.CartDelMsg{
		UserId:    in.UserId,
		GoodsIds:  in.GoodsIds,
		TimeStamp: time.Now().Unix(),
	}
	message, _ := json.Marshal(cartMsg)
	err = l.svcCtx.Producer.Send(l.ctx, l.svcCtx.Config.RocketMqConf.Topics.TopicDelCart, message)
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "createOrder send err", err.Error())
	}
	// 发送延迟消息，用于超时关闭订单
	msg := mq.OrderOffMessage{
		OrderNo:   strconv.FormatInt(snowId, 10),
		UserId:    in.UserId,
		TimeStamp: time.Now().Unix(),
	}
	data, _ := json.Marshal(msg)
	err = l.svcCtx.Producer.SendDelay(l.ctx, l.svcCtx.Config.RocketMqConf.Topics.TopicOrderOff, data, 30*time.Minute)
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "createOrder send delay", err.Error())
		//TODO 投递失败插入本地消息表，人工介入
	}
	return &orderPb.CreateOrderResp{
		OrderNo: strconv.FormatInt(snowId, 10),
	}, nil
}
