// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.2
// source: pkg/client/api.proto

package client

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

// TravelerClient is the client API for Traveler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TravelerClient interface {
	Node(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	Delete(ctx context.Context, in *DockerRequest, opts ...grpc.CallOption) (*DockerResponse, error)
}

type travelerClient struct {
	cc grpc.ClientConnInterface
}

func NewTravelerClient(cc grpc.ClientConnInterface) TravelerClient {
	return &travelerClient{cc}
}

func (c *travelerClient) Node(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.Traveler/Node", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *travelerClient) Delete(ctx context.Context, in *DockerRequest, opts ...grpc.CallOption) (*DockerResponse, error) {
	out := new(DockerResponse)
	err := c.cc.Invoke(ctx, "/grpc.Traveler/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TravelerServer is the server API for Traveler service.
// All implementations must embed UnimplementedTravelerServer
// for forward compatibility
type TravelerServer interface {
	Node(context.Context, *Request) (*Response, error)
	Delete(context.Context, *DockerRequest) (*DockerResponse, error)
	mustEmbedUnimplementedTravelerServer()
}

// UnimplementedTravelerServer must be embedded to have forward compatible implementations.
type UnimplementedTravelerServer struct {
}

func (UnimplementedTravelerServer) Node(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Node not implemented")
}
func (UnimplementedTravelerServer) Delete(context.Context, *DockerRequest) (*DockerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedTravelerServer) mustEmbedUnimplementedTravelerServer() {}

// UnsafeTravelerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TravelerServer will
// result in compilation errors.
type UnsafeTravelerServer interface {
	mustEmbedUnimplementedTravelerServer()
}

func RegisterTravelerServer(s grpc.ServiceRegistrar, srv TravelerServer) {
	s.RegisterService(&Traveler_ServiceDesc, srv)
}

func _Traveler_Node_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TravelerServer).Node(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Traveler/Node",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TravelerServer).Node(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Traveler_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DockerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TravelerServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Traveler/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TravelerServer).Delete(ctx, req.(*DockerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Traveler_ServiceDesc is the grpc.ServiceDesc for Traveler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Traveler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Traveler",
	HandlerType: (*TravelerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Node",
			Handler:    _Traveler_Node_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Traveler_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/client/api.proto",
}
