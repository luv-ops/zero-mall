package logic

import (
	"context"
	"encoding/json"
	"zeromall/common/constant"
	"zeromall/common/xerr"

	"zeromall/cart/rpc/cartPb"
	"zeromall/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
		return nil, xerr.Server()
	}
	var itemList []*cartPb.PreviewItemVO
	for idx, goodsId := range in.GoodsIds {
		jsonStr := jsonArr[idx]
		if jsonStr == "" {
			return nil, xerr.NewCodeError(xerr.GoodsNotInCart)
		}
		var item cartPb.CartItemRedis
		err = json.Unmarshal([]byte(jsonStr), &item)
		if err != nil {
			return nil, xerr.Server()
		}
		if item.Selected != 1 {
			return nil, xerr.NewCodeError(xerr.GoodsNotSelected)
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
