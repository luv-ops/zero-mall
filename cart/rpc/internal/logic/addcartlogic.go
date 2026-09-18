package logic

import (
	"context"
	"encoding/json"
	"time"
	"zeromall/cart/rpc/cartPb"
	"zeromall/cart/rpc/internal/svc"
	"zeromall/common/constant"
	"zeromall/common/mq"
	"zeromall/common/xerr"
	"zeromall/goods/rpc/goodsPb"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCartLogic {
	return &AddCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddCartLogic) AddCart(in *cartPb.AddCartReq) (*cartPb.AddCartResp, error) {
	// todo: add your logic here and delete this line
	if in.AddNum <= 0 {
		return nil, xerr.NewCodeError(xerr.ParamErr)
	}
	key := constant.CartKey + in.UserId

	keys := []string{key}
	//调用lua脚本
	res, err := l.svcCtx.Redis.EvalShaCtx(l.ctx, l.svcCtx.AddCartSha, keys, in.GoodsId, in.AddNum)
	if err != nil {
		return nil, xerr.Server()
	}
	val, ok := res.(int64)
	if !ok {
		return nil, xerr.Server()
	}
	switch val {
	case int64(1):
		//调用goodsRpc，并调用另一个lua脚本
		goodsRes, err1 := l.svcCtx.GoodsRpc.GetGoodsDetail(l.ctx, &goodsPb.GoodsDetailReq{
			GoodsId: in.GoodsId,
		})
		if err1 != nil {
			_, err = l.svcCtx.Redis.HdelCtx(l.ctx, key, in.GoodsId)
			if err != nil {
				l.Logger.Errorf("redis hdel err: %v:%v", "addCart", err.Error())
				return nil, xerr.Server()
			}
			//rpc1调用rpc2，调用者也需要像api层调用rpc一样使用xerr.FromRpcError透传
			return nil, xerr.FromRpcError(err1)
		}
		if goodsRes == nil {
			//goodsRes不存在，则删除key，不需要回填
			_, err = l.svcCtx.Redis.HdelCtx(l.ctx, key, in.GoodsId)
			if err != nil {
				l.Logger.Errorf("redis hdel err: %v:%v", "addCart", err.Error())
				return nil, xerr.Server()
			}

			return nil, xerr.NewCodeError(xerr.NotFound)
		}
		obj := Obj{
			Name:        goodsRes.Name,
			Cover:       goodsRes.Cover,
			Price:       goodsRes.Price,
			OriginPrice: goodsRes.OriginalPrice,
		}
		snapJson, _ := json.Marshal(obj)
		resp, err := l.svcCtx.Redis.EvalShaCtx(l.ctx, l.svcCtx.BackFillSha, keys, in.GoodsId, snapJson)
		if err != nil {
			logc.Infof(l.ctx, "redis backfill err: %v", err)
			return nil, xerr.Server()
		}
		if resp == -1 {
			logc.Infof(l.ctx, "redis backfill fail:%v", err)
			return nil, xerr.Server()
		}
	}
	//生产消息
	msg := mq.CartChangeMsg{
		UserId:    in.UserId,
		TimeStamp: time.Now().Unix(),
	}

	jsonStr, err := json.Marshal(msg)
	if err != nil {
		logc.Infof(l.ctx, "json marshell err %v", err)
		return nil, xerr.Server()
	}
	err = l.svcCtx.Producer.Send(l.ctx, l.svcCtx.Config.RocketMqConf.Topics.TopicSyncFiling, jsonStr)
	if err != nil {
		logc.Infof(l.ctx, "send msg err %v", err)
		return nil, xerr.Server()
	}
	return &cartPb.AddCartResp{
		Ok: true,
	}, nil
}

type Obj struct {
	Name        string `json:"name"`
	Cover       string `json:"cover"`
	Price       string `json:"price"`
	OriginPrice string `json:"originPrice"`
}
