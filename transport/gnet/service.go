package gnet

import "context"

// MethodDesc used by generated code.
type MethodDesc struct {
	Ops     uint32
	Handler func(srv interface{}, ctx context.Context, sess *Session, payload []byte)
	Impl    interface{} // filled at RegisterService time
}

// ServiceDesc used by generated code.
type ServiceDesc struct {
	ServiceName string
	HandlerType interface{}
	Methods     []MethodDesc
}
