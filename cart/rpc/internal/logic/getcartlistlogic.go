package logic

import (
	"context"
	"encoding/json"
	"zeromall/cart/rpc/cartPb"
	"zeromall/cart/rpc/internal/svc"
	"zeromall/common/constant"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetCartListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCartListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCartListLogic {
	return &GetCartListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCartListLogic) GetCartList(in *cartPb.GetCartListReq) (*cartPb.GetCartListResp, error) {
	// todo: add your logic here and delete this line
	key := constant.CartKey + in.UserId
	//获取redis
	redisMap, err := l.svcCtx.Redis.HgetallCtx(l.ctx, key)
	if err != nil {
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	var list []*cartPb.CartItem
	if len(redisMap) == 0 {
		return &cartPb.GetCartListResp{
			List: list,
		}, nil
	}
	for goodsId, jsonStr := range redisMap {
		var item cartPb.CartItemRedis
		err = json.Unmarshal([]byte(jsonStr), &item)
		if err != nil {
			logc.Infof(l.ctx, "unmarshell err :%v in %s", err, "getCartList")
			continue
		}
		list = append(list, &cartPb.CartItem{
			GoodsId:     goodsId,
			Num:         item.Num,
			Selected:    item.Selected,
			Name:        item.Name,
			Cover:       item.Cover,
			Price:       item.Price,
			OriginPrice: item.OriginPrice,
		})
	}
	return &cartPb.GetCartListResp{
		List: list,
	}, nil
}
