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

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteLogic {
	return &BatchDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchDeleteLogic) BatchDelete(in *cartPb.BatchDeleteReq) (*cartPb.BatchDeleteResp, error) {
	// todo: add your logic here and delete this line
	key := constant.CartKey + in.UserId
	ok, err := l.svcCtx.Redis.HdelCtx(l.ctx, key, in.GoodsIds...)
	if err != nil {
		return nil, xerr.Server()
	}
	//生产消息
	msg := mq.CartChangeMsg{
		UserId:    in.UserId,
		TimeStamp: time.Now().Unix(),
	}
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		logx.Errorf("json marshell err %v", err)
		return nil, xerr.Server()
	}
	err = l.svcCtx.Producer.Send(l.ctx, l.svcCtx.Config.RocketMqConf.Topics.TopicSyncFiling, jsonStr)
	if err != nil {
		logx.Errorf("send msg err %v", err)
		return nil, xerr.Server()
	}
	return &cartPb.BatchDeleteResp{
		Ok: ok,
	}, nil
}
