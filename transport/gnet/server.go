package gnet

import (
	"time"

	"github.com/panjf2000/gnet/v2"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) Start() error {
	go gnet.Run(&eventHandler{}, "tcp://"+s.addr, gnet.WithMulticore(true), gnet.WithTCPKeepAlive(time.Minute*5))
	return nil
}

func (s *Server) Stop() error {
	return gnet.Stop(nil, s.addr)
}

func (s *Server) RegisterService(desc interface{}, impl interface{}) {
	// 业务代码可以注册 gnet handler
}
