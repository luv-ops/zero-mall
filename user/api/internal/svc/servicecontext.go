// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"zeromall/user/api/internal/config"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc userpb.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.UserRpcConf)
	return &ServiceContext{
		Config:  c,
		UserRpc: userpb.NewUserClient(client.Conn()),
	}
}
