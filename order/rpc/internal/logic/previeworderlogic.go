package logic

import (
	"context"
	"zeromall/cart/rpc/cartPb"
	"zeromall/common/constant"
	"zeromall/goods/rpc/goodsPb"
	"zeromall/user/rpc/userpb"

	"zeromall/order/rpc/internal/svc"
	"zeromall/order/rpc/orderPb"

	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PreviewOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPreviewOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewOrderLogic {
	return &PreviewOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PreviewOrderLogic) PreviewOrder(in *orderPb.OrderPreviewReq) (*orderPb.OrderPreviewResp, error) {
	// todo: add your logic here and delete this line
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
	var items []*orderPb.PreviewItemVO
	var totalAmount, payAmount decimal.Decimal = decimal.Zero, decimal.Zero

	for goodsId, cartItem := range cartMap {
		goodsItem, _ := goodsMap[goodsId]

		originPrice, _ := decimal.NewFromString(goodsItem.OriginalPrice)
		price, _ := decimal.NewFromString(goodsItem.Price)
		num := decimal.NewFromInt(cartItem.GoodsNum)
		totalAmount = totalAmount.Add(originPrice.Mul(num))
		mulTotal := price.Mul(num)
		payAmount = payAmount.Add(mulTotal)

		items = append(items, &orderPb.PreviewItemVO{
			GoodsId:          goodsId,
			GoodsName:        goodsItem.Name,
			GoodsPrice:       goodsItem.Price,
			GoodsOriginPrice: goodsItem.OriginalPrice,
			GoodsCover:       goodsItem.Cover,
			GoodsNum:         cartItem.GoodsNum,
			SubtotalPrice:    mulTotal.String(),
		})
	}
	//TODO查询用户默认收获地址
	res, err := l.svcCtx.UserRpc.GetDefaultArea(l.ctx, &userpb.GetDefaultAreaReq{
		UserId: in.UserId,
	})
	if err != nil {
		l.Logger.Errorf(constant.RpcError, "previewOrder", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	return &orderPb.OrderPreviewResp{
		TotalAmount:     totalAmount.String(),
		PayAmount:       payAmount.String(),
		ItemList:        items,
		ReceiverAddress: res.ReceiverAddress,
		ReceiverPhone:   res.ReceiverPhone,
		ReceiverName:    res.ReceiverName,
	}, nil
}
