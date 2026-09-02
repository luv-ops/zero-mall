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

type UpdateCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCartLogic {
	return &UpdateCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCartLogic) UpdateCart(req *types.UpdateCartReq) (resp *types.UpdateCartResp, err error) {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	res, err := l.svcCtx.CartRpc.UpdateCart(l.ctx, &cartPb.UpdateCartReq{
		UserId:   userId,
		GoodsId:  req.GoodsId,
		Num:      req.Num,
		Selected: req.Selected,
	})
	if err != nil {
		return nil, err
	}
	return &types.UpdateCartResp{
		GoodsId:  res.GoodsId,
		Num:      res.Num,
		Selected: res.Selected,
	}, nil
}
