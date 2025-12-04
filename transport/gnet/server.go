package gnet

import (
	"context"
	"fmt"
	"time"

	"github.com/panjf2000/gnet/v2"
)

// Server is the gnet transport wrapper.
// It supports RegisterService(desc, impl), Start, Stop.
type Server struct {
	addr    string
	r       *router
	handler *eventHandler
	// keep reference to started flag
	started bool
}

// NewServer create server listening address like ":9002"
func NewServer(addr string) *Server {
	r := newRouter()
	h := newEventHandler(r)
	return &Server{
		addr:    addr,
		r:       r,
		handler: h,
	}
}

// RegisterService injects ServiceDesc methods into router
func (s *Server) RegisterService(desc *ServiceDesc, impl interface{}) {
	// attach impl to each MethodDesc. (optional)
	for i := range desc.Methods {
		desc.Methods[i].Impl = impl
	}
	s.r.addService(desc, impl)
}

// Start launches gnet.Run in goroutine
func (s *Server) Start() error {
	if s.started {
		return nil
	}
	s.started = true
	fmt.Println("starting gnet server on", s.addr)
	go func() {
		// Run is blocking; we run in goroutine to not block caller.
		err := gnet.Run(s.handler, "tcp://"+s.addr,
			gnet.WithMulticore(true),
			gnet.WithTCPKeepAlive(5*time.Minute),
		)
		if err != nil {
			fmt.Println("gnet.Run error:", err)
		}
	}()
	return nil
}

// Stop will attempt to stop the gnet server by address.
// gnet.Stop expects context.Context and addr (v2).
func (s *Server) Stop(ctx context.Context) error {
	if !s.started {
		return nil
	}
	s.started = false
	fmt.Println("stopping gnet server on", s.addr)
	// gnet.Stop(ctx, addr) will stop loops bound to that addr.
	return gnet.Stop(ctx, "tcp://"+s.addr)
}
