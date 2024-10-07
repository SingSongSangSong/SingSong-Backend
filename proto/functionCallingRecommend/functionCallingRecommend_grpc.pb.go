// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.28.0
// source: proto/functionCallingRecommend/functionCallingRecommend.proto

package functionCallingRecommend

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

// FunctionCallingRecommendClient is the client API for FunctionCallingRecommend service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FunctionCallingRecommendClient interface {
	GetFunctionCallingRecommendation(ctx context.Context, in *FunctionCallingRequest, opts ...grpc.CallOption) (*FunctionCallingResponse, error)
}

type functionCallingRecommendClient struct {
	cc grpc.ClientConnInterface
}

func NewFunctionCallingRecommendClient(cc grpc.ClientConnInterface) FunctionCallingRecommendClient {
	return &functionCallingRecommendClient{cc}
}

func (c *functionCallingRecommendClient) GetFunctionCallingRecommendation(ctx context.Context, in *FunctionCallingRequest, opts ...grpc.CallOption) (*FunctionCallingResponse, error) {
	out := new(FunctionCallingResponse)
	err := c.cc.Invoke(ctx, "/functionCallingRecommend.functionCallingRecommend/GetFunctionCallingRecommendation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FunctionCallingRecommendServer is the server API for FunctionCallingRecommend service.
// All implementations must embed UnimplementedFunctionCallingRecommendServer
// for forward compatibility
type FunctionCallingRecommendServer interface {
	GetFunctionCallingRecommendation(context.Context, *FunctionCallingRequest) (*FunctionCallingResponse, error)
	mustEmbedUnimplementedFunctionCallingRecommendServer()
}

// UnimplementedFunctionCallingRecommendServer must be embedded to have forward compatible implementations.
type UnimplementedFunctionCallingRecommendServer struct {
}

func (UnimplementedFunctionCallingRecommendServer) GetFunctionCallingRecommendation(context.Context, *FunctionCallingRequest) (*FunctionCallingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFunctionCallingRecommendation not implemented")
}
func (UnimplementedFunctionCallingRecommendServer) mustEmbedUnimplementedFunctionCallingRecommendServer() {
}

// UnsafeFunctionCallingRecommendServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FunctionCallingRecommendServer will
// result in compilation errors.
type UnsafeFunctionCallingRecommendServer interface {
	mustEmbedUnimplementedFunctionCallingRecommendServer()
}

func RegisterFunctionCallingRecommendServer(s grpc.ServiceRegistrar, srv FunctionCallingRecommendServer) {
	s.RegisterService(&FunctionCallingRecommend_ServiceDesc, srv)
}

func _FunctionCallingRecommend_GetFunctionCallingRecommendation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FunctionCallingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FunctionCallingRecommendServer).GetFunctionCallingRecommendation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/functionCallingRecommend.functionCallingRecommend/GetFunctionCallingRecommendation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FunctionCallingRecommendServer).GetFunctionCallingRecommendation(ctx, req.(*FunctionCallingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FunctionCallingRecommend_ServiceDesc is the grpc.ServiceDesc for FunctionCallingRecommend service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FunctionCallingRecommend_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "functionCallingRecommend.functionCallingRecommend",
	HandlerType: (*FunctionCallingRecommendServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFunctionCallingRecommendation",
			Handler:    _FunctionCallingRecommend_GetFunctionCallingRecommendation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/functionCallingRecommend/functionCallingRecommend.proto",
}
