package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"zeromall/common/constant"
	"zeromall/common/convert"
	"zeromall/goods/rpc/goodsPb"
	"zeromall/goods/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetGoodsDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGoodsDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsDetailLogic {
	return &GetGoodsDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGoodsDetailLogic) GetGoodsDetail(in *goodsPb.GoodsDetailReq) (*goodsPb.GoodsDetailResp, error) {
	// todo: add your logic here and delete this line
	//缓存穿透
	key := constant.GoodsInfoKey + in.GoodsId
	val, err := l.svcCtx.Redis.GetCtx(l.ctx, key)
	var detail goodsPb.GoodsDetailResp
	//查询到redis
	if err == nil {
		if val == constant.RedisEmptyValue {
			return nil, status.Error(codes.NotFound, "商品不存在")
		} else if val != "" {
			err = json.Unmarshal([]byte(val), &detail)
			if err != nil {
				l.Logger.Errorf(constant.UnmarshalErr, "getGoodsDetail", err)
				return nil, status.Error(codes.Internal, constant.MiddlewareError)
			}
			return &detail, nil
		}
	}
	//redis不存在查mysql
	res, err := l.svcCtx.GoodsModel.FindOneByGoodsId(l.ctx, in.GoodsId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//缓存空值，设置ttl 5分钟较短
			_ = l.svcCtx.Redis.SetexCtx(l.ctx, key, constant.RedisEmptyValue, constant.ShortTTL)

			return nil, status.Error(codes.NotFound, constant.GoodsNotFound)
		}
		l.Logger.Errorf(constant.MysqlFailed, "getGoodsDetail ", "select", err)
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}

	detail = goodsPb.GoodsDetailResp{
		GoodsId:       res.GoodsId,
		Name:          res.Name,
		Cover:         res.Cover,
		Price:         convert.CentsToYuanStr(res.PriceCent),
		OriginalPrice: convert.CentsToYuanStr(res.OriginalPriceCent),
		Stock:         res.Stock,
		Sales:         res.Sales,
		Desc:          res.Desc,
		Status:        res.Status,
	}
	jsonStr, err := json.Marshal(&detail)
	if err != nil {
		l.Logger.Errorf(constant.MarshalErr, "getGoodsDetail ", err)
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	//写回redis
	err = l.svcCtx.Redis.SetexCtx(l.ctx, key, string(jsonStr), constant.LongTTL)
	if err != nil {
		l.Logger.Errorf(constant.RedisFailed, "getGoodsDetail ", err)
	}

	return &detail, nil
}
