// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"zeromall/goods/rpc/goodsPb"

	"zeromall/goods/api/internal/svc"
	"zeromall/goods/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdminGoodsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdminGoodsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminGoodsListLogic {
	return &GetAdminGoodsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdminGoodsListLogic) GetAdminGoodsList() (resp *types.SelfGoodsListResp, err error) {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	res, err := l.svcCtx.GoodsRpc.GetAdminGoodsList(l.ctx, &goodsPb.AdminGoodsListReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	var list []*types.SelfGoodsItem
	for _, item := range res.List {
		list = append(list, &types.SelfGoodsItem{
			GoodsId:       item.GoodsId,
			Name:          item.Name,
			Price:         item.Price,
			OriginalPrice: item.OriginalPrice,
			Stock:         item.Stock,
			Sales:         item.Sales,
			CategoryId:    item.CategoryId,
			Status:        item.Status,
		})
	}
	return &types.SelfGoodsListResp{
		List:  list,
		Total: res.Total,
	}, nil
}
