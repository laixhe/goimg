package main

import (
	"github.com/laixhe/goimg/server"
	"github.com/laixhe/goimg/route"
)

func main(){
	route.InitRoute()
	server.RunHttp(":8101")
}