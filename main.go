package main

import (
	"github.com/laixhe/goimg/config"
	"github.com/laixhe/goimg/route"
	"github.com/laixhe/goimg/server"
)

func main() {

	// 初始化路由
	route.InitRoute()

	// 开始监听
	server.RunHttp(config.HttpAddr())
}
