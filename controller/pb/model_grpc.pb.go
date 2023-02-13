// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: controller/pb/model.proto

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

// ListServiceClient is the client API for ListService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ListServiceClient interface {
	CreateList(ctx context.Context, in *CreateListRequest, opts ...grpc.CallOption) (*CreateListReply, error)
	SetList(ctx context.Context, opts ...grpc.CallOption) (ListService_SetListClient, error)
	SetListVersion(ctx context.Context, in *SetListVersionRequest, opts ...grpc.CallOption) (*SetListVersionReply, error)
}

type listServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewListServiceClient(cc grpc.ClientConnInterface) ListServiceClient {
	return &listServiceClient{cc}
}

func (c *listServiceClient) CreateList(ctx context.Context, in *CreateListRequest, opts ...grpc.CallOption) (*CreateListReply, error) {
	out := new(CreateListReply)
	err := c.cc.Invoke(ctx, "/pb.ListService/CreateList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *listServiceClient) SetList(ctx context.Context, opts ...grpc.CallOption) (ListService_SetListClient, error) {
	stream, err := c.cc.NewStream(ctx, &ListService_ServiceDesc.Streams[0], "/pb.ListService/SetList", opts...)
	if err != nil {
		return nil, err
	}
	x := &listServiceSetListClient{stream}
	return x, nil
}

type ListService_SetListClient interface {
	Send(*SetListRequest) error
	Recv() (*SetListReply, error)
	grpc.ClientStream
}

type listServiceSetListClient struct {
	grpc.ClientStream
}

func (x *listServiceSetListClient) Send(m *SetListRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *listServiceSetListClient) Recv() (*SetListReply, error) {
	m := new(SetListReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *listServiceClient) SetListVersion(ctx context.Context, in *SetListVersionRequest, opts ...grpc.CallOption) (*SetListVersionReply, error) {
	out := new(SetListVersionReply)
	err := c.cc.Invoke(ctx, "/pb.ListService/SetListVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ListServiceServer is the server API for ListService service.
// All implementations must embed UnimplementedListServiceServer
// for forward compatibility
type ListServiceServer interface {
	CreateList(context.Context, *CreateListRequest) (*CreateListReply, error)
	SetList(ListService_SetListServer) error
	SetListVersion(context.Context, *SetListVersionRequest) (*SetListVersionReply, error)
	mustEmbedUnimplementedListServiceServer()
}

// UnimplementedListServiceServer must be embedded to have forward compatible implementations.
type UnimplementedListServiceServer struct {
}

func (UnimplementedListServiceServer) CreateList(context.Context, *CreateListRequest) (*CreateListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateList not implemented")
}
func (UnimplementedListServiceServer) SetList(ListService_SetListServer) error {
	return status.Errorf(codes.Unimplemented, "method SetList not implemented")
}
func (UnimplementedListServiceServer) SetListVersion(context.Context, *SetListVersionRequest) (*SetListVersionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetListVersion not implemented")
}
func (UnimplementedListServiceServer) mustEmbedUnimplementedListServiceServer() {}

// UnsafeListServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ListServiceServer will
// result in compilation errors.
type UnsafeListServiceServer interface {
	mustEmbedUnimplementedListServiceServer()
}

func RegisterListServiceServer(s grpc.ServiceRegistrar, srv ListServiceServer) {
	s.RegisterService(&ListService_ServiceDesc, srv)
}

func _ListService_CreateList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListServiceServer).CreateList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ListService/CreateList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListServiceServer).CreateList(ctx, req.(*CreateListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ListService_SetList_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ListServiceServer).SetList(&listServiceSetListServer{stream})
}

type ListService_SetListServer interface {
	Send(*SetListReply) error
	Recv() (*SetListRequest, error)
	grpc.ServerStream
}

type listServiceSetListServer struct {
	grpc.ServerStream
}

func (x *listServiceSetListServer) Send(m *SetListReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *listServiceSetListServer) Recv() (*SetListRequest, error) {
	m := new(SetListRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ListService_SetListVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetListVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListServiceServer).SetListVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ListService/SetListVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListServiceServer).SetListVersion(ctx, req.(*SetListVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ListService_ServiceDesc is the grpc.ServiceDesc for ListService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ListService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.ListService",
	HandlerType: (*ListServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateList",
			Handler:    _ListService_CreateList_Handler,
		},
		{
			MethodName: "SetListVersion",
			Handler:    _ListService_SetListVersion_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SetList",
			Handler:       _ListService_SetList_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "controller/pb/model.proto",
}
