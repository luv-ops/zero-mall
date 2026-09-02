// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"zeromall/goods/api/internal/svc"
	"zeromall/goods/api/internal/types"
	"zeromall/goods/rpc/goodsPb"

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

func (l *AddGoodsLogic) AddGoods(req *types.AddGoodsReq) error {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	resp, err := l.svcCtx.GoodsRpc.AddGoods(l.ctx, &goodsPb.AddGoodsReq{
		Name:          req.Name,
		Cover:         req.Cover,
		OriginalPrice: req.OriginalPrice,
		Price:         req.Price,
		Stock:         req.Stock,
		CategoryId:    req.CategoryId,
		Desc:          req.Desc,
		OwnUserId:     userId,
	})
	if err != nil {
		return err
	}
	if resp.Ok != true {
		return errors.New("上架商品失败")
	}
	return nil
}
