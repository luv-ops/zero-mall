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

type OnOffGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnOffGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnOffGoodsLogic {
	return &OnOffGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnOffGoodsLogic) OnOffGoods(req *types.OnOffGoodsReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line
	_, err = l.svcCtx.GoodsRpc.OnOffGoods(l.ctx, &goodsPb.OnOffGoodsReq{
		GoodsId: req.GoodsId,
		Status:  req.Status,
	})
	if err != nil {
		return nil, xerr.FromRpcError(err)
	}
	return nil, nil
}
