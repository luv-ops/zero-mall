package main

import (
	"flag"
	"fmt"
	"zeromall/order/rpc/internal/config"
	"zeromall/order/rpc/internal/logic/consumer"
	"zeromall/order/rpc/internal/server"
	"zeromall/order/rpc/internal/svc"
	"zeromall/order/rpc/orderPb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		orderPb.RegisterOrderServer(grpcServer, server.NewOrderServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	mgr, err := consumer.NewConsumerManager(ctx)
	if err != nil {
		panic(err)
	}
	group := service.NewServiceGroup()
	group.Add(s)
	group.Add(mgr)
	defer func() {
		group.Stop()
		_ = ctx.Producer.Stop()
	}()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	group.Start()
}
