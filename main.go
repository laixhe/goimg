package main

import (
	"github.com/laixhe/goimg/server"
	"github.com/laixhe/goimg/route"
)

func main(){

	// 初始化路由
	route.InitRoute()

	// 监听 8101 端口
	server.RunHttp(":8101")
}