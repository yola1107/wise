package gnet

import (
	"sync"

	"github.com/panjf2000/gnet/v2"
)

// Session is a very small wrapper for business state.
// Keep it light — no goroutines, no heavy buffers.
type Session struct {
	Conn gnet.Conn
	data sync.Map
}

// NewSession creates Session and attaches it to conn.Context()
func NewSession(c gnet.Conn) *Session {
	s := &Session{Conn: c}
	c.SetContext(s)
	return s
}

// GetSession returns session saved on conn.Context()
func GetSession(c gnet.Conn) *Session {
	if c == nil {
		return nil
	}
	if v := c.Context(); v != nil {
		if s, ok := v.(*Session); ok {
			return s
		}
	}
	return nil
}

func (s *Session) Set(key string, val interface{}) {
	s.data.Store(key, val)
}

func (s *Session) Get(key string) (interface{}, bool) {
	return s.data.Load(key)
}

func (s *Session) Delete(key string) {
	s.data.Delete(key)
}

// Send writes bytes to client (async).
// Use AsyncWrite for non-blocking behavior in gnet event loops.
func (s *Session) Send(b []byte) error {
	if s == nil || s.Conn == nil {
		return nil
	}
	// AsyncWrite is recommended to avoid blocking.
	// AsyncCallback can be nil if we don't need to handle the result.
	return s.Conn.AsyncWrite(b, nil)
}

func (s *Session) Close() error {
	if s == nil || s.Conn == nil {
		return nil
	}
	return s.Conn.Close()
}
