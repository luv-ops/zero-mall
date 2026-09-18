package main

import (
	"flag"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/gateway"
)

var configFile = flag.String("f", "gateway.yaml", "the config file")

func main() {
	var c gateway.GatewayConf
	conf.MustLoad(*configFile, &c) // 加载配置[reference:7]
	gw := gateway.MustNewServer(c) // 创建网关服务
	defer gw.Stop()
	gw.Start()
}
