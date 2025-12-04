package gnet

import (
	"fmt"
	"time"

	"github.com/panjf2000/gnet/v2"
)

var _ gnet.EventHandler = (*eventHandler)(nil)

// eventHandler implements gnet.EventHandler
type eventHandler struct {
	router *router
}

func newEventHandler(r *router) *eventHandler {
	return &eventHandler{router: r}
}

func (h *eventHandler) OnBoot(eng gnet.Engine) gnet.Action {
	fmt.Println("gnet engine boot")
	return gnet.None
}

func (h *eventHandler) OnShutdown(eng gnet.Engine) {
	fmt.Println("gnet engine shutdown")
}

func (h *eventHandler) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	s := NewSession(c)
	h.router.OnSessionOpen(s)
	return nil, gnet.None
}

func (h *eventHandler) OnClose(c gnet.Conn, err error) gnet.Action {
	if s := GetSession(c); s != nil {
		h.router.OnSessionClose(s)
	}
	return gnet.None
}

// OnTraffic is invoked whenever data arrives on the connection.
func (h *eventHandler) OnTraffic(c gnet.Conn) gnet.Action {
	sess := GetSession(c)
	if sess == nil {
		sess = NewSession(c)
	}

	buf := make([]byte, 4096) // 4KB，你可自定义
	n, err := c.Read(buf)
	if err != nil {
		return gnet.Close
	}

	if n > 0 {
		h.router.dispatch(sess, buf[:n])
	}

	return gnet.None
}

func (h *eventHandler) OnTick() (time.Duration, gnet.Action) {
	return time.Second, gnet.None
}
