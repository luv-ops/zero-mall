package config

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql          sqlx.SqlConf
	BalanceRpcConf zrpc.RpcClientConf
	RocketMqConf   struct {
		Endpoint string
		Topics   struct {
			TopicTxPaySuccess string
		}
	}
}
