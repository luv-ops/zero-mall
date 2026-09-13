package main

import (
	"context"
	"flag"
	"fmt"
	"zeromall/cart/rpc/cartPb"
	"zeromall/cart/rpc/internal/config"
	"zeromall/cart/rpc/internal/logic/Timer"
	"zeromall/cart/rpc/internal/logic/consumer"
	"zeromall/cart/rpc/internal/server"
	"zeromall/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/cart.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		cartPb.RegisterCartServer(grpcServer, server.NewCartServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// 组装消费者
	mgr, err := consumer.NewConsumerManager(ctx)
	if err != nil {
		panic(err)
	}
	//统一管理API服务和消费者管理
	group := service.NewServiceGroup()
	group.Add(s)
	group.Add(mgr)

	//开始定时器
	rootCtx, cancel := context.WithCancel(context.Background())
	timer := Timer.NewTimer(ctx)
	timer.StartTicker(rootCtx)

	defer func() {
		cancel()
		group.Stop()
		_ = ctx.Producer.Stop()

	}()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	group.Start()
}
