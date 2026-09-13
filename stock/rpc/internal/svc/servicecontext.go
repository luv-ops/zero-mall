package svc

import (
	"zeromall/stock/rpc/internal/config"
	"zeromall/stock/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	StockModel model.StockModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.MustNewConn(c.Mysql)
	return &ServiceContext{
		Config:     c,
		StockModel: model.NewStockModel(conn),
	}
}
