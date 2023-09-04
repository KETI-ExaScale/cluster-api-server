// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.2
// source: pkg/client/metric/api.proto

package metric

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

// MetricGathererClient is the client API for MetricGatherer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricGathererClient interface {
	Node(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	GPU(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type metricGathererClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricGathererClient(cc grpc.ClientConnInterface) MetricGathererClient {
	return &metricGathererClient{cc}
}

func (c *metricGathererClient) Node(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.MetricGatherer/Node", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricGathererClient) GPU(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.MetricGatherer/GPU", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricGathererServer is the server API for MetricGatherer service.
// All implementations must embed UnimplementedMetricGathererServer
// for forward compatibility
type MetricGathererServer interface {
	Node(context.Context, *Request) (*Response, error)
	GPU(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedMetricGathererServer()
}

// UnimplementedMetricGathererServer must be embedded to have forward compatible implementations.
type UnimplementedMetricGathererServer struct {
}

func (UnimplementedMetricGathererServer) Node(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Node not implemented")
}
func (UnimplementedMetricGathererServer) GPU(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GPU not implemented")
}
func (UnimplementedMetricGathererServer) mustEmbedUnimplementedMetricGathererServer() {}

// UnsafeMetricGathererServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricGathererServer will
// result in compilation errors.
type UnsafeMetricGathererServer interface {
	mustEmbedUnimplementedMetricGathererServer()
}

func RegisterMetricGathererServer(s grpc.ServiceRegistrar, srv MetricGathererServer) {
	s.RegisterService(&MetricGatherer_ServiceDesc, srv)
}

func _MetricGatherer_Node_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricGathererServer).Node(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.MetricGatherer/Node",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricGathererServer).Node(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricGatherer_GPU_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricGathererServer).GPU(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.MetricGatherer/GPU",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricGathererServer).GPU(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// MetricGatherer_ServiceDesc is the grpc.ServiceDesc for MetricGatherer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetricGatherer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.MetricGatherer",
	HandlerType: (*MetricGathererServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Node",
			Handler:    _MetricGatherer_Node_Handler,
		},
		{
			MethodName: "GPU",
			Handler:    _MetricGatherer_GPU_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/client/metric/api.proto",
}
