package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zeromall/cart/rpc/cartPb"
	"zeromall/cart/rpc/internal/svc"
	"zeromall/common/constant"
	"zeromall/common/mq"
	"zeromall/goods/rpc/goodsPb"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Error(codes.InvalidArgument, "数量错误")
	}
	key := constant.CartKey + in.UserId

	keys := []string{key}
	//调用lua脚本
	res, err := l.svcCtx.Redis.EvalShaCtx(l.ctx, l.svcCtx.AddCartSha, keys, in.GoodsId, in.AddNum)
	if err != nil {
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	val, ok := res.(int64)
	if !ok {
		return nil, fmt.Errorf("lua脚本返回类型异常, res=%v", res)
	}
	switch val {
	case int64(1):
		logc.Info(l.ctx, "开始调用goodsRpc")
		//调用goodsRpc，并调用另一个lua脚本
		goodsRes, err := l.svcCtx.GoodsRpc.GetGoodsDetail(l.ctx, &goodsPb.GoodsDetailReq{
			GoodsId: in.GoodsId,
		})
		if err != nil {
			return nil, err
		}
		obj := Obj{
			Name:  goodsRes.Name,
			Cover: goodsRes.Cover,
			Price: goodsRes.Price,
		}
		snapJson, _ := json.Marshal(obj)
		resp, err := l.svcCtx.Redis.EvalShaCtx(l.ctx, l.svcCtx.BackFillSha, keys, in.GoodsId, snapJson)
		if err != nil {
			logc.Infof(l.ctx, "redis backfill err: %v", err)
			return nil, status.Error(codes.Internal, constant.MiddlewareError)
		}
		if resp == -1 {
			logc.Infof(l.ctx, "redis backfill fail:%v", err)
			return nil, status.Error(codes.Internal, constant.MiddlewareError)
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
		return nil, err
	}
	err = l.svcCtx.Producer.Send(l.ctx, l.svcCtx.Config.RocketMqConf.Producer.TopicSyncFiling, jsonStr)
	if err != nil {
		logc.Infof(l.ctx, "send msg err %v", err)
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	l.Logger.Info("发送消息成功")
	return &cartPb.AddCartResp{
		Ok: true,
	}, nil
}

type Obj struct {
	Name  string `json:"name"`
	Cover string `json:"cover"`
	Price string `json:"price"`
}
