// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"zeromall/goods/rpc/goodsPb"

	"zeromall/goods/api/internal/svc"
	"zeromall/goods/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGoodsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGoodsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsListLogic {
	return &GetGoodsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGoodsListLogic) GetGoodsList(req *types.GoodsListReq) (resp *types.GoodsListResp, err error) {
	// todo: add your logic here and delete this line
	if req.Page <= 0 {
		return nil, errors.New("参数page必须大于0")
	}
	res, err := l.svcCtx.GoodsRpc.GetGoodsList(l.ctx, &goodsPb.GoodsListReq{
		CategoryId: req.CategoryId,
		Page:       req.Page,
		PageSize:   req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	var list []*types.GoodsItem
	for _, item := range res.List {
		list = append(list, &types.GoodsItem{
			GoodsId:       item.GoodsId,
			Name:          item.Name,
			Price:         item.Price,
			Cover:         item.Cover,
			OriginalPrice: item.OriginalPrice,
			Sales:         item.Sales,
			Status:        item.Status,
			CategoryId:    item.CategoryId,
		})
	}
	return &types.GoodsListResp{
		List:  list,
		Total: res.Total,
	}, nil
}
