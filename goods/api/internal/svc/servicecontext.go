// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"zeromall/goods/api/internal/config"
	"zeromall/goods/rpc/goodsPb"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	GoodsRpc goodsPb.GoodsClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.GoodsRpcConf)
	return &ServiceContext{
		Config:   c,
		GoodsRpc: goodsPb.NewGoodsClient(client.Conn()),
	}
}
