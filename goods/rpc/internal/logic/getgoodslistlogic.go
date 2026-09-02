package logic

import (
	"context"
	"zeromall/common/constant"
	"zeromall/common/convert"
	"zeromall/goods/rpc/goodsPb"
	"zeromall/goods/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetGoodsListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGoodsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsListLogic {
	return &GetGoodsListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGoodsListLogic) GetGoodsList(in *goodsPb.GoodsListReq) (*goodsPb.GoodsListResp, error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.GoodsModel.PageBreakFind(l.ctx, in.CategoryId, in.Page, in.PageSize)
	if err != nil {
		l.Logger.Errorf(constant.MysqlFailed, "getGoodsList ", "pageBreak", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	var list []*goodsPb.GoodsItem

	for _, item := range res {
		list = append(list, &goodsPb.GoodsItem{
			GoodsId:       item.GoodsId,
			Name:          item.Name,
			Cover:         item.Cover,
			Price:         convert.CentsToYuanStr(item.PriceCent),
			OriginalPrice: convert.CentsToYuanStr(item.OriginalPriceCent),
			Sales:         item.Sales,
			CategoryId:    item.CategoryId,
		})
	}
	return &goodsPb.GoodsListResp{
		List:  list,
		Total: int64(len(res)),
	}, nil
}
