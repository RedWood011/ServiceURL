// protoc  --go_out=../../pkg/api --go_opt=paths=source_relative --go-grpc_out=../../pkg/api --go-grpc_opt=paths=source_relative shortener.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.3
// source: urls.proto

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

const (
	URL_GetURLByID_FullMethodName      = "/pb.URL/GetURLByID"
	URL_PostOneURL_FullMethodName      = "/pb.URL/PostOneURL"
	URL_GetUserURLs_FullMethodName     = "/pb.URL/GetUserURLs"
	URL_PostBatchURLs_FullMethodName   = "/pb.URL/PostBatchURLs"
	URL_DeleteBatchURLs_FullMethodName = "/pb.URL/DeleteBatchURLs"
	URL_GetStats_FullMethodName        = "/pb.URL/GetStats"
)

// URLClient is the client API for URL service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLClient interface {
	GetURLByID(ctx context.Context, in *RetrieveRequest, opts ...grpc.CallOption) (*RetrieveResponse, error)
	PostOneURL(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	GetUserURLs(ctx context.Context, in *GetUserURLsRequest, opts ...grpc.CallOption) (*GetUserURLsResponse, error)
	PostBatchURLs(ctx context.Context, in *CreateBatchRequest, opts ...grpc.CallOption) (*CreateBatchResponse, error)
	DeleteBatchURLs(ctx context.Context, in *DeleteBatchRequest, opts ...grpc.CallOption) (*DeleteBatchResponse, error)
	GetStats(ctx context.Context, in *GetStatsRequest, opts ...grpc.CallOption) (*GetStatsResponse, error)
}

type uRLClient struct {
	cc grpc.ClientConnInterface
}

func NewURLClient(cc grpc.ClientConnInterface) URLClient {
	return &uRLClient{cc}
}

func (c *uRLClient) GetURLByID(ctx context.Context, in *RetrieveRequest, opts ...grpc.CallOption) (*RetrieveResponse, error) {
	out := new(RetrieveResponse)
	err := c.cc.Invoke(ctx, URL_GetURLByID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLClient) PostOneURL(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, URL_PostOneURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLClient) GetUserURLs(ctx context.Context, in *GetUserURLsRequest, opts ...grpc.CallOption) (*GetUserURLsResponse, error) {
	out := new(GetUserURLsResponse)
	err := c.cc.Invoke(ctx, URL_GetUserURLs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLClient) PostBatchURLs(ctx context.Context, in *CreateBatchRequest, opts ...grpc.CallOption) (*CreateBatchResponse, error) {
	out := new(CreateBatchResponse)
	err := c.cc.Invoke(ctx, URL_PostBatchURLs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLClient) DeleteBatchURLs(ctx context.Context, in *DeleteBatchRequest, opts ...grpc.CallOption) (*DeleteBatchResponse, error) {
	out := new(DeleteBatchResponse)
	err := c.cc.Invoke(ctx, URL_DeleteBatchURLs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLClient) GetStats(ctx context.Context, in *GetStatsRequest, opts ...grpc.CallOption) (*GetStatsResponse, error) {
	out := new(GetStatsResponse)
	err := c.cc.Invoke(ctx, URL_GetStats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLServer is the server API for URL service.
// All implementations must embed UnimplementedURLServer
// for forward compatibility
type URLServer interface {
	GetURLByID(context.Context, *RetrieveRequest) (*RetrieveResponse, error)
	PostOneURL(context.Context, *CreateRequest) (*CreateResponse, error)
	GetUserURLs(context.Context, *GetUserURLsRequest) (*GetUserURLsResponse, error)
	PostBatchURLs(context.Context, *CreateBatchRequest) (*CreateBatchResponse, error)
	DeleteBatchURLs(context.Context, *DeleteBatchRequest) (*DeleteBatchResponse, error)
	GetStats(context.Context, *GetStatsRequest) (*GetStatsResponse, error)
	mustEmbedUnimplementedURLServer()
}

// UnimplementedURLServer must be embedded to have forward compatible implementations.
type UnimplementedURLServer struct {
}

func (UnimplementedURLServer) GetURLByID(context.Context, *RetrieveRequest) (*RetrieveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetURLByID not implemented")
}
func (UnimplementedURLServer) PostOneURL(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostOneURL not implemented")
}
func (UnimplementedURLServer) GetUserURLs(context.Context, *GetUserURLsRequest) (*GetUserURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserURLs not implemented")
}
func (UnimplementedURLServer) PostBatchURLs(context.Context, *CreateBatchRequest) (*CreateBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostBatchURLs not implemented")
}
func (UnimplementedURLServer) DeleteBatchURLs(context.Context, *DeleteBatchRequest) (*DeleteBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBatchURLs not implemented")
}
func (UnimplementedURLServer) GetStats(context.Context, *GetStatsRequest) (*GetStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStats not implemented")
}
func (UnimplementedURLServer) mustEmbedUnimplementedURLServer() {}

// UnsafeURLServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLServer will
// result in compilation errors.
type UnsafeURLServer interface {
	mustEmbedUnimplementedURLServer()
}

func RegisterURLServer(s grpc.ServiceRegistrar, srv URLServer) {
	s.RegisterService(&URL_ServiceDesc, srv)
}

func _URL_GetURLByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServer).GetURLByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URL_GetURLByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServer).GetURLByID(ctx, req.(*RetrieveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URL_PostOneURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServer).PostOneURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URL_PostOneURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServer).PostOneURL(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URL_GetUserURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServer).GetUserURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URL_GetUserURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServer).GetUserURLs(ctx, req.(*GetUserURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URL_PostBatchURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServer).PostBatchURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URL_PostBatchURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServer).PostBatchURLs(ctx, req.(*CreateBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URL_DeleteBatchURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServer).DeleteBatchURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URL_DeleteBatchURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServer).DeleteBatchURLs(ctx, req.(*DeleteBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URL_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URL_GetStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServer).GetStats(ctx, req.(*GetStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// URL_ServiceDesc is the grpc.ServiceDesc for URL service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URL_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.URL",
	HandlerType: (*URLServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetURLByID",
			Handler:    _URL_GetURLByID_Handler,
		},
		{
			MethodName: "PostOneURL",
			Handler:    _URL_PostOneURL_Handler,
		},
		{
			MethodName: "GetUserURLs",
			Handler:    _URL_GetUserURLs_Handler,
		},
		{
			MethodName: "PostBatchURLs",
			Handler:    _URL_PostBatchURLs_Handler,
		},
		{
			MethodName: "DeleteBatchURLs",
			Handler:    _URL_DeleteBatchURLs_Handler,
		},
		{
			MethodName: "GetStats",
			Handler:    _URL_GetStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "urls.proto",
}