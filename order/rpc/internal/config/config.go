package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql        sqlx.SqlConf
	RedisConf    redis.RedisConf
	CartRpcConf  zrpc.RpcClientConf
	GoodsRpcConf zrpc.RpcClientConf
	UserRpcConf  zrpc.RpcClientConf
}
