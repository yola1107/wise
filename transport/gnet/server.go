package gnet

import (
	"time"

	"github.com/panjf2000/gnet/v2"
)

type Server struct {
	addr   string
	router *Router
}

func NewServer(addr string) *Server {
	return &Server{
		addr:   addr,
		router: NewRouter(), // 最简单 router skeleton
	}
}

func (s *Server) Start() error {
	go gnet.Run(&eventHandler{router: s.router}, "tcp://"+s.addr,
		gnet.WithMulticore(true),
		gnet.WithTCPKeepAlive(time.Minute*5),
	)
	return nil
}

// Stop 停止 gnet server
func (s *Server) Stop() error {
	// gnet 的 Stop 需要保存 EventHandler 或者通过 Shutdown 实现
	return nil
}

// RegisterService 注册 service handler
func (s *Server) RegisterService(desc interface{}, impl interface{}) {
	//s.router.Register(desc, impl)
}
