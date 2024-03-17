//
// New here? You may find these links useful:
//
// Protobuf:
//   https://developers.google.com/protocol-buffers/docs/proto3
//   https://developers.google.com/protocol-buffers/docs/gotutorial
//
// GRPC:
//   https://www.grpc.io/docs/what-is-grpc/introduction/
//
// After model was changed, do not forget:
//   make protobuf
//

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: api/v1/proto/remote.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Remote_PressKey_FullMethodName = "/com.github.kaedwen.v1.Remote/PressKey"
	Remote_ExecTask_FullMethodName = "/com.github.kaedwen.v1.Remote/ExecTask"
	Remote_ListTask_FullMethodName = "/com.github.kaedwen.v1.Remote/ListTask"
)

// RemoteClient is the client API for Remote service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RemoteClient interface {
	PressKey(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*KeyResponse, error)
	ExecTask(ctx context.Context, in *ExecTaskRequest, opts ...grpc.CallOption) (*ExecTaskResponse, error)
	ListTask(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListTaskResponse, error)
}

type remoteClient struct {
	cc grpc.ClientConnInterface
}

func NewRemoteClient(cc grpc.ClientConnInterface) RemoteClient {
	return &remoteClient{cc}
}

func (c *remoteClient) PressKey(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*KeyResponse, error) {
	out := new(KeyResponse)
	err := c.cc.Invoke(ctx, Remote_PressKey_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *remoteClient) ExecTask(ctx context.Context, in *ExecTaskRequest, opts ...grpc.CallOption) (*ExecTaskResponse, error) {
	out := new(ExecTaskResponse)
	err := c.cc.Invoke(ctx, Remote_ExecTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *remoteClient) ListTask(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListTaskResponse, error) {
	out := new(ListTaskResponse)
	err := c.cc.Invoke(ctx, Remote_ListTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RemoteServer is the server API for Remote service.
// All implementations must embed UnimplementedRemoteServer
// for forward compatibility
type RemoteServer interface {
	PressKey(context.Context, *KeyRequest) (*KeyResponse, error)
	ExecTask(context.Context, *ExecTaskRequest) (*ExecTaskResponse, error)
	ListTask(context.Context, *emptypb.Empty) (*ListTaskResponse, error)
	mustEmbedUnimplementedRemoteServer()
}

// UnimplementedRemoteServer must be embedded to have forward compatible implementations.
type UnimplementedRemoteServer struct {
}

func (UnimplementedRemoteServer) PressKey(context.Context, *KeyRequest) (*KeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PressKey not implemented")
}
func (UnimplementedRemoteServer) ExecTask(context.Context, *ExecTaskRequest) (*ExecTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecTask not implemented")
}
func (UnimplementedRemoteServer) ListTask(context.Context, *emptypb.Empty) (*ListTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTask not implemented")
}
func (UnimplementedRemoteServer) mustEmbedUnimplementedRemoteServer() {}

// UnsafeRemoteServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RemoteServer will
// result in compilation errors.
type UnsafeRemoteServer interface {
	mustEmbedUnimplementedRemoteServer()
}

func RegisterRemoteServer(s grpc.ServiceRegistrar, srv RemoteServer) {
	s.RegisterService(&Remote_ServiceDesc, srv)
}

func _Remote_PressKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteServer).PressKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Remote_PressKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteServer).PressKey(ctx, req.(*KeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Remote_ExecTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteServer).ExecTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Remote_ExecTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteServer).ExecTask(ctx, req.(*ExecTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Remote_ListTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteServer).ListTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Remote_ListTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteServer).ListTask(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Remote_ServiceDesc is the grpc.ServiceDesc for Remote service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Remote_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.github.kaedwen.v1.Remote",
	HandlerType: (*RemoteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PressKey",
			Handler:    _Remote_PressKey_Handler,
		},
		{
			MethodName: "ExecTask",
			Handler:    _Remote_ExecTask_Handler,
		},
		{
			MethodName: "ListTask",
			Handler:    _Remote_ListTask_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/proto/remote.proto",
}
