package svc

import (
	"zeromall/balance/rpc/internal/config"
	"zeromall/balance/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	BalanceModel model.BalanceModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.MustNewConn(c.Mysql)
	return &ServiceContext{
		Config:       c,
		BalanceModel: model.NewBalanceModel(conn),
	}
}
