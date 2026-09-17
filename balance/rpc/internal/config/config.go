package config

import (
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql        sqlx.SqlConf
	RocketMqConf struct {
		Endpoint string
		Topics   struct {
			TopicInsertBalance string
		}
		Consumer struct {
			Group struct {
				BalanceInsertGroup string
			}
			AwaitDuration     time.Duration
			MaxMsgNum         int32
			InvisibleDuration time.Duration
		}
	}
}
