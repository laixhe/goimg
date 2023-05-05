package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Server struct {
	http *http.Server
	mux  *http.ServeMux
}

func NewServer() *Server {
	return &Server{
		http: &http.Server{},
		mux:  http.NewServeMux(),
	}
}

func (s *Server) Func(f func(s *Server)) *Server {
	f(s)
	return s
}

// HandleFunc 注册访问路由
func (s *Server) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

// Handle 注册访问路由
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) HttpRun(addr string) {
	s.http.Addr = addr
	s.http.Handler = s.mux

	logrus.Debugf("http listen %s", addr)

	// 启动监听
	err := s.http.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
