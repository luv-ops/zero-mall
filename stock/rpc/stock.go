package main

import (
	"flag"
	"fmt"
	"zeromall/common/intercepter"
	"zeromall/stock/rpc/internal/logic/consumer"

	"zeromall/stock/rpc/internal/config"
	"zeromall/stock/rpc/internal/server"
	"zeromall/stock/rpc/internal/svc"
	"zeromall/stock/rpc/stockPb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/stock.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		stockPb.RegisterStockServer(grpcServer, server.NewStockServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(intercepter.ErrorInterceptor)
	group := service.NewServiceGroup()
	manager, err := consumer.NewConsumerManager(ctx)
	if err != nil {
		panic(err)
	}
	group.Add(s)
	group.Add(manager)
	defer group.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	group.Start()
}
