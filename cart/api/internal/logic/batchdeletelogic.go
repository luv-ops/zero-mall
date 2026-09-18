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

type BatchDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteLogic {
	return &BatchDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteLogic) BatchDelete(req *types.BatchDeleteReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(string)
	_, err = l.svcCtx.CartRpc.BatchDelete(l.ctx, &cartPb.BatchDeleteReq{
		UserId:   userId,
		GoodsIds: req.GoodsIds,
	})
	if err != nil {
		return nil, xerr.FromRpcError(err)
	}
	return nil, nil
}
