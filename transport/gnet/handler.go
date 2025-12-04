package gnet

import (
	"time"

	"github.com/panjf2000/gnet/v2"
)

// eventHandler 是 gnet TCP 事件处理器
type eventHandler struct {
	router *Router
}

// OnBoot 在 gnet 引擎启动时调用
func (h *eventHandler) OnBoot(eng gnet.Engine) gnet.Action {
	return gnet.None
}

// OnShutdown 在 gnet 引擎关闭时调用
func (h *eventHandler) OnShutdown(eng gnet.Engine) {
	// 可在这里清理资源
}

// OnOpen 新连接打开时调用
func (h *eventHandler) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	session := NewSession(c)
	h.router.OnSessionOpen(session)
	return nil, gnet.None
}

// OnClose 连接关闭时调用
func (h *eventHandler) OnClose(c gnet.Conn, err error) gnet.Action {
	session := GetSession(c)
	if session != nil {
		h.router.OnSessionClose(session)
	}
	return gnet.None
}

// OnTraffic 收到数据时调用
func (h *eventHandler) OnTraffic(c gnet.Conn) gnet.Action {
	session := GetSession(c)
	if session == nil {
		session = NewSession(c)
	}

	// TCP 原始消息直接分发
	msg, err := c.Read()
	if err != nil {
		return gnet.Close
	}

	h.router.Dispatch(session, msg)
	return gnet.None
}

// OnTick 定时任务
func (h *eventHandler) OnTick() (time.Duration, gnet.Action) {
	// 可做心跳或定时任务
	return time.Second * 1, gnet.None
}

// OnInitComplete gnet 初始化完成时调用
func (h *eventHandler) OnInitComplete(srv gnet.Server) gnet.Action {
	return gnet.None
}
