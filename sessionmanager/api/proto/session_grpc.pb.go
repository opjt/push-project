// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: sessionmanager/api/proto/session.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SessionService_Connect_FullMethodName     = "/session.SessionService/Connect"
	SessionService_PushMessage_FullMethodName = "/session.SessionService/PushMessage"
)

// SessionServiceClient is the client API for SessionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SessionServiceClient interface {
	Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ServerMessage], error)
	PushMessage(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*PushResponse, error)
}

type sessionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSessionServiceClient(cc grpc.ClientConnInterface) SessionServiceClient {
	return &sessionServiceClient{cc}
}

func (c *sessionServiceClient) Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ServerMessage], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &SessionService_ServiceDesc.Streams[0], SessionService_Connect_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ConnectRequest, ServerMessage]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SessionService_ConnectClient = grpc.ServerStreamingClient[ServerMessage]

func (c *sessionServiceClient) PushMessage(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*PushResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PushResponse)
	err := c.cc.Invoke(ctx, SessionService_PushMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SessionServiceServer is the server API for SessionService service.
// All implementations must embed UnimplementedSessionServiceServer
// for forward compatibility.
type SessionServiceServer interface {
	Connect(*ConnectRequest, grpc.ServerStreamingServer[ServerMessage]) error
	PushMessage(context.Context, *PushRequest) (*PushResponse, error)
	mustEmbedUnimplementedSessionServiceServer()
}

// UnimplementedSessionServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSessionServiceServer struct{}

func (UnimplementedSessionServiceServer) Connect(*ConnectRequest, grpc.ServerStreamingServer[ServerMessage]) error {
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedSessionServiceServer) PushMessage(context.Context, *PushRequest) (*PushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushMessage not implemented")
}
func (UnimplementedSessionServiceServer) mustEmbedUnimplementedSessionServiceServer() {}
func (UnimplementedSessionServiceServer) testEmbeddedByValue()                        {}

// UnsafeSessionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SessionServiceServer will
// result in compilation errors.
type UnsafeSessionServiceServer interface {
	mustEmbedUnimplementedSessionServiceServer()
}

func RegisterSessionServiceServer(s grpc.ServiceRegistrar, srv SessionServiceServer) {
	// If the following call pancis, it indicates UnimplementedSessionServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SessionService_ServiceDesc, srv)
}

func _SessionService_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConnectRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SessionServiceServer).Connect(m, &grpc.GenericServerStream[ConnectRequest, ServerMessage]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SessionService_ConnectServer = grpc.ServerStreamingServer[ServerMessage]

func _SessionService_PushMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionServiceServer).PushMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SessionService_PushMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionServiceServer).PushMessage(ctx, req.(*PushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SessionService_ServiceDesc is the grpc.ServiceDesc for SessionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SessionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "session.SessionService",
	HandlerType: (*SessionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushMessage",
			Handler:    _SessionService_PushMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Connect",
			Handler:       _SessionService_Connect_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "sessionmanager/api/proto/session.proto",
}
