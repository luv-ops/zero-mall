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
		Topics   struct {
			TopicInsertStock   string
			TopicReturnStock   string
			TopicTxFrozenStock string
			TopicDeductStock   string
		}
		Consumer struct {
			Group struct {
				StockInsertGroup string
				StockReturnGroup string
				StockFrozenGroup string
				StockDeductGroup string
			}
			AwaitDuration     time.Duration
			MaxMsgNum         int32
			InvisibleDuration time.Duration
		}
	}
}
