package svc

import (
	"zeromall/goods/rpc/internal/config"
	"zeromall/goods/rpc/internal/model"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	GoodsModel model.GoodsModel
	UserRpc    userpb.UserClient
	Redis      *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	client := zrpc.MustNewClient(c.UserRpcConf)
	rdb := redis.MustNewRedis(c.RedisConf)
	return &ServiceContext{
		Config:     c,
		GoodsModel: model.NewGoodsModel(conn),
		UserRpc:    userpb.NewUserClient(client.Conn()),
		Redis:      rdb,
	}
}
