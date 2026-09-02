package config

import (
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql        sqlx.SqlConf
	RedisConf    redis.RedisConf
	RocketMqConf struct {
		Endpoint string
		Producer struct {
			TopicSyncFiling string
		}
		Consumer struct {
			TopicSyncFiling string
			Group           string
			AwaitDuration   time.Duration
			MaxMsgNum       int32
			LoopDuration    time.Duration
		}
	}
	GoodsRpcConf zrpc.RpcClientConf
}
