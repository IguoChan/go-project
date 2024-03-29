// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.3
// source: simplepb/simple.proto

package simplepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SimpleClient is the client API for Simple service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SimpleClient interface {
	Route(ctx context.Context, in *SimpleRequest, opts ...grpc.CallOption) (*SimpleResponse, error)
}

type simpleClient struct {
	cc grpc.ClientConnInterface
}

func NewSimpleClient(cc grpc.ClientConnInterface) SimpleClient {
	return &simpleClient{cc}
}

func (c *simpleClient) Route(ctx context.Context, in *SimpleRequest, opts ...grpc.CallOption) (*SimpleResponse, error) {
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, "/simplepb.Simple/Route", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SimpleServer is the server API for Simple service.
// All implementations must embed UnimplementedSimpleServer
// for forward compatibility
type SimpleServer interface {
	Route(context.Context, *SimpleRequest) (*SimpleResponse, error)
	mustEmbedUnimplementedSimpleServer()
}

// UnimplementedSimpleServer must be embedded to have forward compatible implementations.
type UnimplementedSimpleServer struct {
}

func (UnimplementedSimpleServer) Route(context.Context, *SimpleRequest) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Route not implemented")
}
func (UnimplementedSimpleServer) mustEmbedUnimplementedSimpleServer() {}

// UnsafeSimpleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SimpleServer will
// result in compilation errors.
type UnsafeSimpleServer interface {
	mustEmbedUnimplementedSimpleServer()
}

func RegisterSimpleServer(s grpc.ServiceRegistrar, srv SimpleServer) {
	s.RegisterService(&Simple_ServiceDesc, srv)
}

func _Simple_Route_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimpleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimpleServer).Route(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/simplepb.Simple/Route",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimpleServer).Route(ctx, req.(*SimpleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Simple_ServiceDesc is the grpc.ServiceDesc for Simple service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Simple_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "simplepb.Simple",
	HandlerType: (*SimpleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Route",
			Handler:    _Simple_Route_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "simplepb/simple.proto",
}

// AAClient is the client API for AA service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AAClient interface {
	Get(ctx context.Context, in *AAReq, opts ...grpc.CallOption) (*AAres, error)
}

type aAClient struct {
	cc grpc.ClientConnInterface
}

func NewAAClient(cc grpc.ClientConnInterface) AAClient {
	return &aAClient{cc}
}

func (c *aAClient) Get(ctx context.Context, in *AAReq, opts ...grpc.CallOption) (*AAres, error) {
	out := new(AAres)
	err := c.cc.Invoke(ctx, "/simplepb.AA/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AAServer is the server API for AA service.
// All implementations must embed UnimplementedAAServer
// for forward compatibility
type AAServer interface {
	Get(context.Context, *AAReq) (*AAres, error)
	mustEmbedUnimplementedAAServer()
}

// UnimplementedAAServer must be embedded to have forward compatible implementations.
type UnimplementedAAServer struct {
}

func (UnimplementedAAServer) Get(context.Context, *AAReq) (*AAres, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedAAServer) mustEmbedUnimplementedAAServer() {}

// UnsafeAAServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AAServer will
// result in compilation errors.
type UnsafeAAServer interface {
	mustEmbedUnimplementedAAServer()
}

func RegisterAAServer(s grpc.ServiceRegistrar, srv AAServer) {
	s.RegisterService(&AA_ServiceDesc, srv)
}

func _AA_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AAReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AAServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/simplepb.AA/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AAServer).Get(ctx, req.(*AAReq))
	}
	return interceptor(ctx, in, info, handler)
}

// AA_ServiceDesc is the grpc.ServiceDesc for AA service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AA_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "simplepb.AA",
	HandlerType: (*AAServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _AA_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "simplepb/simple.proto",
}
