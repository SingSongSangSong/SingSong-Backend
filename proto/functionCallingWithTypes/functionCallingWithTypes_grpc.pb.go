// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.28.0
// source: proto/functionCallingWithTypes/functionCallingWithTypes.proto

package functionCallingWithTypes

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

// FunctionCallingWithTypesRecommendClient is the client API for FunctionCallingWithTypesRecommend service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FunctionCallingWithTypesRecommendClient interface {
	GetFunctionCallingWithTypesRecommendation(ctx context.Context, in *FunctionCallingWithTypesRequest, opts ...grpc.CallOption) (*FunctionCallingWithTypesResponse, error)
}

type functionCallingWithTypesRecommendClient struct {
	cc grpc.ClientConnInterface
}

func NewFunctionCallingWithTypesRecommendClient(cc grpc.ClientConnInterface) FunctionCallingWithTypesRecommendClient {
	return &functionCallingWithTypesRecommendClient{cc}
}

func (c *functionCallingWithTypesRecommendClient) GetFunctionCallingWithTypesRecommendation(ctx context.Context, in *FunctionCallingWithTypesRequest, opts ...grpc.CallOption) (*FunctionCallingWithTypesResponse, error) {
	out := new(FunctionCallingWithTypesResponse)
	err := c.cc.Invoke(ctx, "/functionCallingWithTypes.FunctionCallingWithTypesRecommend/GetFunctionCallingWithTypesRecommendation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FunctionCallingWithTypesRecommendServer is the server API for FunctionCallingWithTypesRecommend service.
// All implementations must embed UnimplementedFunctionCallingWithTypesRecommendServer
// for forward compatibility
type FunctionCallingWithTypesRecommendServer interface {
	GetFunctionCallingWithTypesRecommendation(context.Context, *FunctionCallingWithTypesRequest) (*FunctionCallingWithTypesResponse, error)
	mustEmbedUnimplementedFunctionCallingWithTypesRecommendServer()
}

// UnimplementedFunctionCallingWithTypesRecommendServer must be embedded to have forward compatible implementations.
type UnimplementedFunctionCallingWithTypesRecommendServer struct {
}

func (UnimplementedFunctionCallingWithTypesRecommendServer) GetFunctionCallingWithTypesRecommendation(context.Context, *FunctionCallingWithTypesRequest) (*FunctionCallingWithTypesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFunctionCallingWithTypesRecommendation not implemented")
}
func (UnimplementedFunctionCallingWithTypesRecommendServer) mustEmbedUnimplementedFunctionCallingWithTypesRecommendServer() {
}

// UnsafeFunctionCallingWithTypesRecommendServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FunctionCallingWithTypesRecommendServer will
// result in compilation errors.
type UnsafeFunctionCallingWithTypesRecommendServer interface {
	mustEmbedUnimplementedFunctionCallingWithTypesRecommendServer()
}

func RegisterFunctionCallingWithTypesRecommendServer(s grpc.ServiceRegistrar, srv FunctionCallingWithTypesRecommendServer) {
	s.RegisterService(&FunctionCallingWithTypesRecommend_ServiceDesc, srv)
}

func _FunctionCallingWithTypesRecommend_GetFunctionCallingWithTypesRecommendation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FunctionCallingWithTypesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FunctionCallingWithTypesRecommendServer).GetFunctionCallingWithTypesRecommendation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/functionCallingWithTypes.FunctionCallingWithTypesRecommend/GetFunctionCallingWithTypesRecommendation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FunctionCallingWithTypesRecommendServer).GetFunctionCallingWithTypesRecommendation(ctx, req.(*FunctionCallingWithTypesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FunctionCallingWithTypesRecommend_ServiceDesc is the grpc.ServiceDesc for FunctionCallingWithTypesRecommend service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FunctionCallingWithTypesRecommend_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "functionCallingWithTypes.FunctionCallingWithTypesRecommend",
	HandlerType: (*FunctionCallingWithTypesRecommendServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFunctionCallingWithTypesRecommendation",
			Handler:    _FunctionCallingWithTypesRecommend_GetFunctionCallingWithTypesRecommendation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/functionCallingWithTypes/functionCallingWithTypes.proto",
}
