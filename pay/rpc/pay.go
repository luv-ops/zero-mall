package main

import (
	"flag"
	"fmt"
	"zeromall/common/intercepter"
	"zeromall/common/mq"
	"zeromall/pay/rpc/internal/logic/txProducer"

	"zeromall/pay/rpc/internal/config"
	"zeromall/pay/rpc/internal/server"
	"zeromall/pay/rpc/internal/svc"
	"zeromall/pay/rpc/payPb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/pay.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		payPb.RegisterPayServer(grpcServer, server.NewPayServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(intercepter.ErrorInterceptor)
	pg := mq.ProducerConfig{
		Endpoint: c.RocketMqConf.Endpoint,
		Topics:   []string{c.RocketMqConf.Topics.TopicTxPaySuccess},
	}
	checker := txProducer.NewPayTransactionChecker(ctx)
	txPro, err := mq.NewTxProducer(&pg, checker)
	if err != nil {
		panic(err)
	}
	ctx.TxProducer = txPro
	defer func() {
		s.Stop()
		_ = txPro.Stop()
	}()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
