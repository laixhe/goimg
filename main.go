package main

import (
	"github.com/laixhe/goimg/route"
	"github.com/laixhe/goimg/server"
	"github.com/laixhe/goimg/config"
)

func main() {

	// 初始化路由
	route.InitRoute()

	// 开始监听
	server.RunHttp(config.HttpAddr())
}
