package logic

import (
	"context"
	"encoding/json"
	"time"
	"zeromall/common/constant"
	"zeromall/common/mq"
	"zeromall/common/xerr"

	"zeromall/cart/rpc/cartPb"
	"zeromall/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCartLogic {
	return &UpdateCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCartLogic) UpdateCart(in *cartPb.UpdateCartReq) (*cartPb.UpdateCartResp, error) {
	// todo: add your logic here and delete this line
	key := constant.CartKey + in.UserId
	//如果num为0，则删除购物车商品
	if in.Num == 0 {
		ok, err := l.svcCtx.Redis.HdelCtx(l.ctx, key, in.GoodsId)
		if err != nil {
			l.Logger.Errorf("redis hdel err:%v", err)
			return nil, xerr.Server()
		}
		if !ok {
			return nil, xerr.NewCodeError(xerr.DeleteCartErr)
		}
		return &cartPb.UpdateCartResp{
			GoodsId: in.GoodsId,
		}, nil
	}
	keys := []string{key}
	res, err := l.svcCtx.Redis.EvalShaCtx(l.ctx, l.svcCtx.UpdateCartSha, keys, in.GoodsId, in.Num, in.Selected)
	if err != nil {
		l.Logger.Errorf("redis eval err:%v", err)
		return nil, xerr.Server()
	}
	if res == -1 {
		l.Logger.Errorf("redis update err:%v", err)
		return nil, xerr.Server()
	}
	//生产消息
	msg := mq.CartChangeMsg{
		UserId:    in.UserId,
		TimeStamp: time.Now().Unix(),
	}

	jsonStr, err := json.Marshal(msg)
	if err != nil {
		l.Logger.Errorf("json marshell err %v in %v", err, "UpdateCart")
		return nil, xerr.Server()
	}
	err = l.svcCtx.Producer.Send(l.ctx, l.svcCtx.Config.RocketMqConf.Topics.TopicSyncFiling, jsonStr)
	if err != nil {
		logx.Errorf("send msg err %v", err)
		return nil, xerr.Server()
	}
	return &cartPb.UpdateCartResp{
		GoodsId:  in.GoodsId,
		Num:      in.Num,
		Selected: in.Selected,
	}, nil
}
