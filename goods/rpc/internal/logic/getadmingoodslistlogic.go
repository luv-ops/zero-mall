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

type GetAdminGoodsListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAdminGoodsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminGoodsListLogic {
	return &GetAdminGoodsListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAdminGoodsListLogic) GetAdminGoodsList(in *goodsPb.AdminGoodsListReq) (*goodsPb.AdminGoodsListResp, error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.GoodsModel.FindByOwnId(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Errorf(constant.MysqlFailed, "getAdminGoodsList ", "insert", err)
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	var list []*goodsPb.AdminGoodsItem
	for _, v := range res {
		list = append(list, &goodsPb.AdminGoodsItem{
			GoodsId:       v.GoodsId,
			Name:          v.Name,
			Cover:         v.Cover,
			Price:         convert.CentsToYuanStr(v.PriceCent),
			OriginalPrice: convert.CentsToYuanStr(v.OriginalPriceCent),
			Stock:         v.Stock,
			Sales:         v.Sales,
			CategoryId:    v.CategoryId,
			Status:        v.Status,
		})
	}
	return &goodsPb.AdminGoodsListResp{
		List:  list,
		Total: int64(len(res)),
	}, nil
}
