package route

import (
	"github.com/laixhe/goimg/server"
	"github.com/laixhe/goimg/uphand"
)

// 注册访问路由
func InitRoute() {

	// 路由处理绑定
	server.Handle("/", uphand.Controller{})

	// 测试上传
	server.HandleFunc("/test", uphand.Test)
}
