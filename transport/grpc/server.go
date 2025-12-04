package grpc

import (
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	srv *grpc.Server
	lis net.Listener
}

func NewServer(addr string) *Server {
	return &Server{
		srv: grpc.NewServer(),
	}
}

func (s *Server) Start() error {
	var err error
	s.lis, err = net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}
	go s.srv.Serve(s.lis)
	return nil
}

func (s *Server) Stop() error {
	s.srv.GracefulStop()
	return nil
}

func (s *Server) RegisterService(desc interface{}, impl interface{}) {
	// 业务代码可以注册 gRPC service
}
