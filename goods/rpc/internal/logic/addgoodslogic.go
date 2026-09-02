package logic

import (
	"context"
	"database/sql"
	"zeromall/common/constant"
	"zeromall/common/convert"
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
	//金额转换
	price, err := convert.YuanStrToCents(in.Price)
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "addGoods", "YuanStrToCents", err.Error())
		return nil, status.Error(codes.InvalidArgument, constant.GoodsArgError)
	}
	original, err := convert.YuanStrToCents(in.OriginalPrice)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, constant.GoodsArgError)
	}
	goods := model.Goods{
		GoodsId:           uuid.NewString(),
		Name:              in.Name,
		Cover:             in.Cover,
		PriceCent:         price,
		OriginalPriceCent: original,
		Stock:             in.Stock,
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

	return &goodsPb.AddGoodsResp{
		Ok: num > 0,
	}, nil
}
