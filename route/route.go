package route

import (
	"github.com/laixhe/goimg/server"
	"github.com/laixhe/goimg/controller"
)

// 注册访问路由
func InitRoute() {

	// 
	server.HandleFunc("/", controller.Index)
	server.HandleFunc("/upload", controller.Upload)

}
