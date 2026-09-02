package main

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/gateway"
)

func main() {
	var c gateway.GatewayConf
	conf.MustLoad("gateway/gateway.yaml", &c) // 加载配置[reference:7]
	gw := gateway.MustNewServer(c)            // 创建网关服务
	defer gw.Stop()
	gw.Start()
}
