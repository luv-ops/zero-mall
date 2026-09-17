package svc

import (
	"zeromall/balance/rpc/balancePb"
	"zeromall/common/mq"
	"zeromall/pay/rpc/internal/config"
	"zeromall/pay/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	PayModel   model.PayModel
	BalanceRpc balancePb.BalanceClient
	TxProducer *mq.TxProducer
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.MustNewConn(c.Mysql)
	client := zrpc.MustNewClient(c.BalanceRpcConf)
	return &ServiceContext{
		Config:     c,
		PayModel:   model.NewPayModel(conn),
		BalanceRpc: balancePb.NewBalanceClient(client.Conn()),
		TxProducer: nil,
	}
}
