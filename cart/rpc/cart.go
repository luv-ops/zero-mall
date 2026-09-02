package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"zeromall/cart/rpc/cartPb"
	"zeromall/cart/rpc/internal/config"
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

	rootCtx, cancel := context.WithCancel(context.Background())

	consume, err := consumer.NewConsumer(ctx)
	if err != nil {
		log.Fatal("构造消费者失败", err)
	}
	//开始消费循环和计时器
	//TODO 如果多实例部署，ticker必须单例，因为如果非单例，多个ticker会同时spop redis，造成无效redis Io
	consume.StartTicker(rootCtx)
	err = consume.StartConsumer(ctx.Config.RocketMqConf.Consumer.TopicSyncFiling)
	if err != nil {
		log.Fatal("消费循环启动失败", err)
	}
	defer func() {
		_ = consume.StopConsumer()
		cancel()
		s.Stop()
		_ = ctx.Producer.Stop()
	}()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
