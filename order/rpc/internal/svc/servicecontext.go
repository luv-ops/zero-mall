package svc

import (
	"zeromall/cart/rpc/cartPb"
	"zeromall/common/mq"
	"zeromall/goods/rpc/goodsPb"
	"zeromall/order/rpc/internal/config"
	"zeromall/order/rpc/internal/model"
	"zeromall/user/rpc/userpb"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	OrderModel     model.OrderModel
	OrderItemModel model.OrderItemModel
	Redis          *redis.Redis
	CartRpc        cartPb.CartClient
	GoodsRpc       goodsPb.GoodsClient
	UserRpc        userpb.UserClient
	Snow           *snowflake.Node
	Producer       *mq.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	rdb := redis.MustNewRedis(c.RedisConf)
	cartClient := zrpc.MustNewClient(c.CartRpcConf)
	goodsClient := zrpc.MustNewClient(c.GoodsRpcConf)
	userClient := zrpc.MustNewClient(c.UserRpcConf)
	node, err := snowflake.NewNode(c.Snowflake.NodeId)
	if err != nil {
		panic(err)
	}
	pg := mq.ProducerConfig{
		Endpoint: c.RocketMqConf.Endpoint,
		Topics:   []string{c.RocketMqConf.Topics.TopicOrderOff},
	}
	pro, err := mq.NewProducer(&pg)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:         c,
		OrderModel:     model.NewOrderModel(sqlConn),
		OrderItemModel: model.NewOrderItemModel(sqlConn),
		Redis:          rdb,
		CartRpc:        cartPb.NewCartClient(cartClient.Conn()),
		GoodsRpc:       goodsPb.NewGoodsClient(goodsClient.Conn()),
		UserRpc:        userpb.NewUserClient(userClient.Conn()),
		Snow:           node,
		Producer:       pro,
	}
}
