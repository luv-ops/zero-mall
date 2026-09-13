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
	CartRpcConf  zrpc.RpcClientConf
	GoodsRpcConf zrpc.RpcClientConf
	UserRpcConf  zrpc.RpcClientConf
	Snowflake    struct {
		NodeId int64
	}
	RocketMqConf struct {
		Endpoint string
		Topics   struct {
			TopicOrderOff    string
			TopicDelCart     string
			TopicReturnStock string
		}
		Consumer struct {
			Group struct {
				TopicOffGroup string
			}
			AwaitDuration     time.Duration
			MaxMsgNum         int32
			InvisibleDuration time.Duration
		}
	}
}
