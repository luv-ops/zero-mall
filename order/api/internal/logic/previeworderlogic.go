// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"zeromall/order/rpc/orderPb"

	"zeromall/order/api/internal/svc"
	"zeromall/order/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreviewOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPreviewOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewOrderLogic {
	return &PreviewOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreviewOrderLogic) PreviewOrder(req *types.OrderPreviewReq) (resp *types.OrderPreviewResp, err error) {
	// todo: add your logic here and delete this line
	if len(req.GoodsIds) == 0 || req.ReceiveAddressId == 0 {
		return nil, errors.New("param error")
	}
	userId := l.ctx.Value("userId").(string)
	res, err := l.svcCtx.OrderRpc.PreviewOrder(l.ctx, &orderPb.OrderPreviewReq{
		UserId:           userId,
		GoodsIds:         req.GoodsIds,
		ReceiveAddressId: req.ReceiveAddressId,
	})
	var itemList []*types.PreviewItemVO
	for _, v := range res.ItemList {
		itemList = append(itemList, &types.PreviewItemVO{
			GoodsId:       v.GoodsId,
			GoodsName:     v.GoodsName,
			GoodsPrice:    v.GoodsPrice,
			GoodsCover:    v.GoodsCover,
			GoodsNum:      v.GoodsNum,
			SubtotalPrice: v.SubtotalPrice,
		})
	}
	return &types.OrderPreviewResp{
		ReceiverPhone:   res.ReceiverPhone,
		ReceiverName:    res.ReceiverName,
		ReceiverAddress: res.ReceiverAddress,
		TotalAmount:     res.TotalAmount,
		PayAmount:       res.PayAmount,
		ItemList:        itemList,
	}, err
}
