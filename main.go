package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/laixhe/goimg/config"
	"github.com/laixhe/goimg/route"
	"github.com/laixhe/goimg/server"
)

func main() {
	logrus.SetOutput(os.Stdout)                  // 以 stdout 为输出，代替默认的 stderr
	logrus.SetFormatter(&logrus.JSONFormatter{}) // 日志作为 JSON 而不是默认的 ASCII 格式器
	logrus.SetReportCaller(true)                 // 记录方法名称
	logrus.SetLevel(logrus.TraceLevel)           // 日志级别

	// 初始化配置
	config.Init("config.yaml")

	// 开始监听
	server.NewServer().
		Func(route.Init).          // 初始化路由
		HttpRun(config.HttpAddr()) // 开始监听
}
