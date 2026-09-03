// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"zeromall/order/api/internal/config"
	"zeromall/order/rpc/orderPb"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	OrderRpc orderPb.OrderClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.OrderRpcConf)
	return &ServiceContext{
		Config:   c,
		OrderRpc: orderPb.NewOrderClient(client.Conn()),
	}
}
