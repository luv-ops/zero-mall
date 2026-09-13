package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
	"zeromall/common/constant"
	"zeromall/common/convert"
	"zeromall/common/mq"
	"zeromall/goods/rpc/internal/model"
	"zeromall/user/rpc/userpb"

	"zeromall/goods/rpc/goodsPb"
	"zeromall/goods/rpc/internal/svc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AddGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGoodsLogic {
	return &AddGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddGoodsLogic) AddGoods(in *goodsPb.AddGoodsReq) (*goodsPb.AddGoodsResp, error) {
	// todo: add your logic here and delete this line
	//查上架人员是否有售卖权
	res, err := l.svcCtx.UserRpc.GetSellPower(l.ctx, &userpb.GetSellPowerReq{
		UserId: in.OwnUserId,
	})
	if err != nil {
		l.Logger.Error("调用GetSellPower失败", err.Error())
		return nil, status.Error(codes.Internal, "检验售卖权失败")
	}
	if res.Ok != true {
		return nil, status.Error(codes.PermissionDenied, constant.PermissionSellError)
	}
	goodsId := uuid.NewString()
	goods := model.Goods{
		GoodsId:           goodsId,
		Name:              in.Name,
		Cover:             in.Cover,
		PriceCent:         convert.YuanStrToCents(in.Price),
		OriginalPriceCent: convert.YuanStrToCents(in.OriginalPrice),
		CategoryId:        in.CategoryId,
		OwnUserId:         sql.NullString{String: in.OwnUserId, Valid: true},
		Desc:              in.Desc,
	}
	result, err := l.svcCtx.GoodsModel.Insert(l.ctx, &goods)
	if err != nil {
		l.Logger.Errorf(constant.MysqlFailed, "addFoods", "insert", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	num, _ := result.RowsAffected()
	//库存已从goods表移除，发送一条消息，插入库存
	msg := &mq.InsertStockMsg{
		Stock:     in.Stock,
		GoodsId:   goodsId,
		TimeStamp: time.Now().Unix(),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		l.Logger.Errorf(constant.MarshalErr, "addFoods", "jsonMarshal", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	err = l.svcCtx.Producer.Send(l.ctx, l.svcCtx.Config.RocketMQConf.Topics.TopicInsertStock, data)
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "addFoods", "send", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	return &goodsPb.AddGoodsResp{
		Ok: num > 0,
	}, nil
}
