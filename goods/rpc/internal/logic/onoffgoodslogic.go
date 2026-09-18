package logic

import (
	"context"
	"zeromall/common/constant"
	"zeromall/common/xerr"

	"zeromall/goods/rpc/goodsPb"
	"zeromall/goods/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnOffGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOnOffGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnOffGoodsLogic {
	return &OnOffGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OnOffGoodsLogic) OnOffGoods(in *goodsPb.OnOffGoodsReq) (*goodsPb.OnOffGoodsResp, error) {
	// todo: add your logic here and delete this line
	dataMap := make(map[string]any)
	dataMap["status"] = in.Status
	num, err := l.svcCtx.GoodsModel.UpdateFields(l.ctx, in.GoodsId, dataMap)
	if err != nil {
		l.Logger.Errorf(constant.MysqlFailed, "onOffGoods", "update", err.Error())
		return nil, xerr.Server()
	}

	return &goodsPb.OnOffGoodsResp{
		Ok: num != 0,
	}, nil
}
