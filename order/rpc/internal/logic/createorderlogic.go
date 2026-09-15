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
	var remark string
	var remarkBool bool
	if in.Remark != nil {
		remark = *in.Remark
		remarkBool = true
	} else {
		remark = ""
		remarkBool = false
	}
	orderData := &model.Orders{
		OrderNo:    snowId,
		UserId:     in.UserId,
		PayCent:    payAmount.Round(2).Mul(decimal.NewFromInt(100)).IntPart(),
		TotalCent:  totalAmount.Round(2).Mul(decimal.NewFromInt(100)).IntPart(),
		ExpireTime: time.Now().Add(l.svcCtx.Config.RocketMqConf.DelayOffOrderDuration),
		Remark:     sql.NullString{String: remark, Valid: remarkBool},
	}
	//调用redis预扣减库存
	var successList []*model.OrderItem
	ok := true
	for _, v := range orderItemData {
		key := constant.StockGoodsKey + v.GoodsId
		res, err := l.svcCtx.Redis.EvalShaCtx(l.ctx, l.svcCtx.StockFrozenSha, []string{key}, v.Num)
		if err != nil {
			ok = false
			break
		}
		ret, _ := res.(int64)
		if ret != int64(1) {
			ok = false
			break
		}
		successList = append(successList, v)
	}
	//有一个商品扣减失败，回滚
	if !ok {
		for _, v := range successList {
			key := constant.StockGoodsKey + v.GoodsId
			_, err := l.svcCtx.Redis.EvalShaCtx(l.ctx, l.svcCtx.UnFrozenStockSha, []string{key}, v.Num)
			if err != nil {
				l.Logger.Errorf("回滚失败", "unFrozenStockSha", err.Error())
				continue
			}
			return nil, status.Error(codes.Internal, "库存不足")
		}
	}
	//插入三张表 order,order_item,transactionLog
	//开启消息事务

	//发送给库存服务的消息
	var stockItems []*mq.FrozenItem
	for _, v := range orderItemData {
		stockItems = append(stockItems, &mq.FrozenItem{
			OrderNo: v.OrderNo,
			GoodsId: v.GoodsId,
			Num:     v.Num,
		})
	}
	var stockMsg mq.FrozenStockMsg
	stockMsg.List = stockItems
	stockData, err := json.Marshal(stockMsg)
	if err != nil {
		l.Logger.Errorf(constant.MarshalErr, "createOrder", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	//开启消息事务
	tx := l.svcCtx.TxProducer.BeginTransaction()
	//发送半消息,预冻结库存
	receipt, err := l.svcCtx.TxProducer.SendWithTransaction(l.ctx, l.svcCtx.Config.RocketMqConf.Topics.TopicFrozenStock, stockData, tx)
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "createOrder", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	transactionData := &model.TransactionLog{
		TxId:    receipt.TransactionId,
		OrderNo: snowId,
	}
	//插入三张表 order,order_item,transactionLog
	num, err := l.svcCtx.OrderModel.TxInsert(l.ctx, orderData, orderItemData, transactionData)
	if err != nil {
		//回滚
		_ = tx.RollBack()
		l.Logger.Errorf(constant.WhereFailed, "createOrder", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	//提交后，消费者才能真正收到消息
	if err := tx.Commit(); err != nil {
		l.Logger.Errorf(constant.WhereFailed, "commit err", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	if num == 0 {
		l.Logger.Infof("createOrder没有出现错误，但是没有任何行被影响")
	}
	// 发送延迟消息，用于超时关闭订单
	msg := mq.OrderOffMessage{
		OrderNo:   snowId,
		UserId:    in.UserId,
		TimeStamp: time.Now().Unix(),
	}
	data, _ := json.Marshal(msg)
	err = l.svcCtx.Producer.SendDelay(l.ctx, l.svcCtx.Config.RocketMqConf.Topics.TopicDelayOrderOff, data, l.svcCtx.Config.RocketMqConf.DelayOffOrderDuration)
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "createOrder send delay", err.Error())
		//TODO 投递失败插入本地消息表，人工介入
	}
	//创建订单成功后删除购物车,发送消息异步删除
	go func() {
		cartMsg := mq.CartDelMsg{
			UserId:    in.UserId,
			GoodsIds:  in.GoodsIds,
			TimeStamp: time.Now().Unix(),
		}
		message, _ := json.Marshal(cartMsg)
		//协程中不能用l.ctx,因为预计执行时间可能请求链路已结束
		err = l.svcCtx.Producer.Send(context.TODO(), l.svcCtx.Config.RocketMqConf.Topics.TopicDelCart, message)
		if err != nil {
			l.Logger.Errorf(constant.WhereFailed, "createOrder send err", err.Error())
		}
	}()
	return &orderPb.CreateOrderResp{
		OrderNo: strconv.FormatInt(snowId, 10),
	}, nil
}
