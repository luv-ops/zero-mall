// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"zeromall/cart/rpc/cartPb"
	"zeromall/common/xerr"

	"zeromall/cart/api/internal/svc"
	"zeromall/cart/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCartLogic {
	return &AddCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCartLogic) AddCart(req *types.AddCartReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	_, err = l.svcCtx.CartRpc.AddCart(l.ctx, &cartPb.AddCartReq{
		UserId:  userId,
		GoodsId: req.GoodsId,
		AddNum:  req.AddNum,
	})
	if err != nil {
		return nil, xerr.FromRpcError(err)
	}
	return nil, nil
}
