// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"zeromall/cart/api/internal/config"
	"zeromall/cart/rpc/cartPb"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	CartRpc cartPb.CartClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.CartRpcConf)
	return &ServiceContext{
		Config:  c,
		CartRpc: cartPb.NewCartClient(client.Conn()),
	}
}
