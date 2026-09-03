package logic

import (
	"context"
	"encoding/json"
	"zeromall/common/constant"

	"zeromall/cart/rpc/cartPb"
	"zeromall/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BatchGetCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetCartLogic {
	return &BatchGetCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetCartLogic) BatchGetCart(in *cartPb.BatchGetCartReq) (*cartPb.BatchGetCartResp, error) {
	// todo: add your logic here and delete this line
	key := constant.CartKey + in.UserId
	jsonArr, err := l.svcCtx.Redis.HmgetCtx(l.ctx, key, in.GoodsIds...)
	if err != nil {
		l.Logger.Errorf(constant.RedisFailed, "batchGetCart", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	var itemList []*cartPb.PreviewItemVO
	for idx, goodsId := range in.GoodsIds {
		jsonStr := jsonArr[idx]
		if jsonStr == "" {
			return nil, status.Error(codes.NotFound, "商品已不在购物车中")
		}
		var item cartPb.CartItemRedis
		err = json.Unmarshal([]byte(jsonStr), &item)
		if err != nil {
			return nil, status.Errorf(codes.Internal, constant.UnmarshalErr, "batchGetCart", err.Error())
		}
		if item.Selected != 1 {
			return nil, status.Error(codes.Unavailable, "商品未勾选")
		}
		//不拿金额，不信任redis存储的金额
		itemList = append(itemList, &cartPb.PreviewItemVO{
			GoodsId:  goodsId,
			GoodsNum: item.Num,
		})
	}
	return &cartPb.BatchGetCartResp{
		ItemList: itemList,
	}, nil
}
