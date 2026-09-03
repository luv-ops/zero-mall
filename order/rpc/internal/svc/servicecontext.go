package svc

import (
	"zeromall/cart/rpc/cartPb"
	"zeromall/goods/rpc/goodsPb"
	"zeromall/order/rpc/internal/config"
	"zeromall/order/rpc/internal/model"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	OrderModel model.OrderModel
	Redis      *redis.Redis
	CartRpc    cartPb.CartClient
	GoodsRpc   goodsPb.GoodsClient
	UserRpc    userpb.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	rdb := redis.MustNewRedis(c.RedisConf)
	cartClient := zrpc.MustNewClient(c.CartRpcConf)
	goodsClient := zrpc.MustNewClient(c.GoodsRpcConf)
	userClient := zrpc.MustNewClient(c.UserRpcConf)
	return &ServiceContext{
		Config:     c,
		OrderModel: model.NewOrderModel(sqlConn),
		Redis:      rdb,
		CartRpc:    cartPb.NewCartClient(cartClient.Conn()),
		GoodsRpc:   goodsPb.NewGoodsClient(goodsClient.Conn()),
		UserRpc:    userpb.NewUserClient(userClient.Conn()),
	}
}
