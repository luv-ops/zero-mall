// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"zeromall/common/xerr"
	"zeromall/goods/rpc/goodsPb"

	"zeromall/goods/api/internal/svc"
	"zeromall/goods/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGoodsLogic {
	return &AddGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddGoodsLogic) AddGoods(req *types.AddGoodsReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line
	_, err = l.svcCtx.GoodsRpc.AddGoods(l.ctx, &goodsPb.AddGoodsReq{
		Name:          req.Name,
		Cover:         req.Cover,
		OriginalPrice: req.OriginalPrice,
		Price:         req.Price,
		Stock:         req.Stock,
		CategoryId:    req.CategoryId,
		Desc:          req.Desc,
	})
	if err != nil {
		return nil, xerr.FromRpcError(err)
	}
	return nil, nil
}
