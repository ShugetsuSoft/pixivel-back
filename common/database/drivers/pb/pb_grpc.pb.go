// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.2
// source: pb.proto

package pb

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

// NearDBServiceClient is the client API for NearDBService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NearDBServiceClient interface {
	Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*NoneResponse, error)
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*QueryResponse, error)
	QueryById(ctx context.Context, in *QueryByIdRequest, opts ...grpc.CallOption) (*QueryResponse, error)
	Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*NoneResponse, error)
}

type nearDBServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNearDBServiceClient(cc grpc.ClientConnInterface) NearDBServiceClient {
	return &nearDBServiceClient{cc}
}

func (c *nearDBServiceClient) Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*NoneResponse, error) {
	out := new(NoneResponse)
	err := c.cc.Invoke(ctx, "/neardbv2.pb.NearDBService/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nearDBServiceClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*QueryResponse, error) {
	out := new(QueryResponse)
	err := c.cc.Invoke(ctx, "/neardbv2.pb.NearDBService/Query", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nearDBServiceClient) QueryById(ctx context.Context, in *QueryByIdRequest, opts ...grpc.CallOption) (*QueryResponse, error) {
	out := new(QueryResponse)
	err := c.cc.Invoke(ctx, "/neardbv2.pb.NearDBService/QueryById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nearDBServiceClient) Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*NoneResponse, error) {
	out := new(NoneResponse)
	err := c.cc.Invoke(ctx, "/neardbv2.pb.NearDBService/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NearDBServiceServer is the server API for NearDBService service.
// All implementations must embed UnimplementedNearDBServiceServer
// for forward compatibility
type NearDBServiceServer interface {
	Add(context.Context, *AddRequest) (*NoneResponse, error)
	Query(context.Context, *QueryRequest) (*QueryResponse, error)
	QueryById(context.Context, *QueryByIdRequest) (*QueryResponse, error)
	Remove(context.Context, *RemoveRequest) (*NoneResponse, error)
	mustEmbedUnimplementedNearDBServiceServer()
}

// UnimplementedNearDBServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNearDBServiceServer struct {
}

func (UnimplementedNearDBServiceServer) Add(context.Context, *AddRequest) (*NoneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedNearDBServiceServer) Query(context.Context, *QueryRequest) (*QueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedNearDBServiceServer) QueryById(context.Context, *QueryByIdRequest) (*QueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryById not implemented")
}
func (UnimplementedNearDBServiceServer) Remove(context.Context, *RemoveRequest) (*NoneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (UnimplementedNearDBServiceServer) mustEmbedUnimplementedNearDBServiceServer() {}

// UnsafeNearDBServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NearDBServiceServer will
// result in compilation errors.
type UnsafeNearDBServiceServer interface {
	mustEmbedUnimplementedNearDBServiceServer()
}

func RegisterNearDBServiceServer(s grpc.ServiceRegistrar, srv NearDBServiceServer) {
	s.RegisterService(&NearDBService_ServiceDesc, srv)
}

func _NearDBService_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NearDBServiceServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neardbv2.pb.NearDBService/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NearDBServiceServer).Add(ctx, req.(*AddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NearDBService_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NearDBServiceServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neardbv2.pb.NearDBService/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NearDBServiceServer).Query(ctx, req.(*QueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NearDBService_QueryById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NearDBServiceServer).QueryById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neardbv2.pb.NearDBService/QueryById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NearDBServiceServer).QueryById(ctx, req.(*QueryByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NearDBService_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NearDBServiceServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/neardbv2.pb.NearDBService/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NearDBServiceServer).Remove(ctx, req.(*RemoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NearDBService_ServiceDesc is the grpc.ServiceDesc for NearDBService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NearDBService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "neardbv2.pb.NearDBService",
	HandlerType: (*NearDBServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _NearDBService_Add_Handler,
		},
		{
			MethodName: "Query",
			Handler:    _NearDBService_Query_Handler,
		},
		{
			MethodName: "QueryById",
			Handler:    _NearDBService_QueryById_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _NearDBService_Remove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb.proto",
}
