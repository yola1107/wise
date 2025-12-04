package gnet

import (
	"context"
	"encoding/binary"
	"sync"
)

// internal router mapping ops->handler func
type router struct {
	mu       sync.RWMutex
	handlers map[uint32]func(sess *Session, payload []byte)
	// optional hooks
	onOpen  func(*Session)
	onClose func(*Session)
}

func newRouter() *router {
	return &router{
		handlers: make(map[uint32]func(sess *Session, payload []byte)),
	}
}

func (r *router) register(ops uint32, h func(sess *Session, payload []byte)) {
	r.mu.Lock()
	r.handlers[ops] = h
	r.mu.Unlock()
}

func (r *router) dispatch(sess *Session, data []byte) {
	if len(data) < 4 {
		// ignore short frames
		return
	}
	ops := binary.BigEndian.Uint32(data[:4])
	r.mu.RLock()
	h, ok := r.handlers[ops]
	r.mu.RUnlock()
	if ok && h != nil {
		h(sess, data[4:])
	}
}

func (r *router) OnSessionOpen(s *Session) {
	if r.onOpen != nil {
		r.onOpen(s)
	}
}

func (r *router) OnSessionClose(s *Session) {
	if r.onClose != nil {
		r.onClose(s)
	}
}

func (r *router) SetOpenHook(fn func(*Session))  { r.onOpen = fn }
func (r *router) SetCloseHook(fn func(*Session)) { r.onClose = fn }

// helper to register ServiceDesc (called by Server.RegisterService)
func (r *router) addService(desc *ServiceDesc, impl interface{}) {
	for _, m := range desc.Methods {
		ops := m.Ops
		// copy local for closure
		handler := m.Handler
		// wrap to call generated handler function
		r.register(ops, func(sess *Session, payload []byte) {
			// call generated wrapper with impl
			if handler != nil {
				handler(impl, context.Background(), sess, payload)
			}
		})
	}
}
