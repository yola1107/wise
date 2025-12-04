package gnet

import (
	"context"
	"encoding/binary"
)

type Handler func(srv interface{}, ctx context.Context, s *Session, payload []byte)

type MethodDesc struct {
	Ops     uint32
	Handler Handler
	Impl    interface{}
}

type ServiceDesc struct {
	ServiceName string
	HandlerType interface{}
	Methods     []MethodDesc
}

type Router struct {
	methods map[uint32]MethodDesc
	//onOpen  func(*Session)
	//onClose func(*Session)
}

func NewRouter() *Router {
	return &Router{
		methods: make(map[uint32]MethodDesc),
	}
}

func (r *Router) AddService(desc *ServiceDesc, impl interface{}) {
	for _, m := range desc.Methods {
		m.Impl = impl
		r.methods[m.Ops] = m
	}
}

func (r *Router) Dispatch(s *Session, data []byte) {
	if len(data) < 4 {
		return
	}
	ops := binary.BigEndian.Uint32(data[:4])
	if m, ok := r.methods[ops]; ok {
		m.Handler(m.Impl, context.Background(), s, data[4:])
	}
}

func (r *Router) OnSessionOpen(s *Session) {
	if r.onOpen != nil {
		r.onOpen(s)
	}
}

func (r *Router) OnSessionClose(s *Session) {
	if r.onClose != nil {
		r.onClose(s)
	}
}

func (r *Router) SetOpenHook(fn func(*Session))  { r.onOpen = fn }
func (r *Router) SetCloseHook(fn func(*Session)) { r.onClose = fn }
