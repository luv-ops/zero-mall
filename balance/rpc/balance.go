package main

import (
	"flag"
	"fmt"
	"zeromall/balance/rpc/internal/logic/consumer"

	"zeromall/balance/rpc/balancePb"
	"zeromall/balance/rpc/internal/config"
	"zeromall/balance/rpc/internal/server"
	"zeromall/balance/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/balance.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		balancePb.RegisterBalanceServer(grpcServer, server.NewBalanceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	manager, err := consumer.NewConsumerManager(ctx)
	if err != nil {
		panic(err)
	}
	group := service.NewServiceGroup()
	group.Add(manager)
	group.Add(s)
	defer group.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	group.Start()
}
