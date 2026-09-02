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

type GetGoodsDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGoodsDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsDetailLogic {
	return &GetGoodsDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGoodsDetailLogic) GetGoodsDetail(req *types.GoodsDetailReq) (resp *types.GoodsDetailResp, err error) {
	// todo: add your logic here and delete this line
	if req.GoodsId == "" {
		return nil, errors.New("goods id is empty")
	}
	res, err := l.svcCtx.GoodsRpc.GetGoodsDetail(l.ctx, &goodsPb.GoodsDetailReq{
		GoodsId: req.GoodsId,
	})
	if err != nil {
		return nil, err
	}

	return &types.GoodsDetailResp{
		GoodsId:       res.GoodsId,
		Name:          res.Name,
		Cover:         res.Cover,
		Price:         res.Price,
		OriginalPrice: res.OriginalPrice,
		Sales:         res.Sales,
		Desc:          res.Desc,
		Status:        res.Status,
	}, nil
}
