package server

import (
	"log"
	"net/http"
	"time"
)

var serveMux *http.ServeMux = http.NewServeMux()

// 注册访问路由
func HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	serveMux.HandleFunc(pattern, handler)
}
func Handle(pattern string, handler http.Handler) {
	serveMux.Handle(pattern, handler)
}

// 启动 HTTP 服务
// start http server
func RunHttp(adrr string) {

	// 手工配置 http.Server 服务
	server := http.Server{
		Addr:              adrr,            // 监听地址和端口
		Handler:           serveMux,        // Handle
		ReadTimeout:       5 * time.Second, // 读超时
		WriteTimeout:      5 * time.Second, // 写超时
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}

	// 启动监听
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
