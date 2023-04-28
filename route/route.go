package route

import (
	"github.com/laixhe/goimg/app"
	"github.com/laixhe/goimg/server"
)

// Init 注册访问路由
func Init(s *server.Server) {
	controller := app.NewController()
	// 路由处理绑定
	s.Handle("/", controller)
	// 获取图片信息
	s.HandleFunc("/info", controller.Info)
	// 测试上传
	s.HandleFunc("/test", controller.Test)
	// 获取状态码
	s.HandleFunc("/statuscode", controller.StatusCode)
}
