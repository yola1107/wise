package http

import "net/http"

type Server struct {
	srv *http.Server
}

func NewServer(addr string) *Server {
	return &Server{
		srv: &http.Server{Addr: addr},
	}
}

func (s *Server) Start() error {
	go s.srv.ListenAndServe()
	return nil
}

func (s *Server) Stop() error {
	return s.srv.Close()
}

func (s *Server) RegisterService(desc interface{}, impl interface{}) {
	// 业务代码可以注册 HTTP handler
}
