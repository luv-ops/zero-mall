// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"zeromall/cart/rpc/cartPb"

	"zeromall/cart/api/internal/svc"
	"zeromall/cart/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCartListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCartListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCartListLogic {
	return &GetCartListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCartListLogic) GetCartList() (resp *types.CartListResp, err error) {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	res, err := l.svcCtx.CartRpc.GetCartList(l.ctx, &cartPb.GetCartListReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	var list []*types.CartItem
	if len(res.List) == 0 {
		return &types.CartListResp{
			CartList: list,
		}, nil
	}
	for _, item := range res.List {
		list = append(list, &types.CartItem{
			Num:      item.Num,
			Name:     item.Name,
			Selected: item.Selected,
			GoodsId:  item.GoodsId,
			Cover:    item.Cover,
			Price:    item.Price,
		})
	}
	return &types.CartListResp{
		CartList: list,
	}, nil
}
