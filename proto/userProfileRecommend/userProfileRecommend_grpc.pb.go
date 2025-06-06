// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.28.0
// source: proto/userProfileRecommend/userProfileRecommend.proto

package userProfileRecommend

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

// UserProfileClient is the client API for UserProfile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserProfileClient interface {
	CreateUserProfile(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*ProfileResponse, error)
}

type userProfileClient struct {
	cc grpc.ClientConnInterface
}

func NewUserProfileClient(cc grpc.ClientConnInterface) UserProfileClient {
	return &userProfileClient{cc}
}

func (c *userProfileClient) CreateUserProfile(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*ProfileResponse, error) {
	out := new(ProfileResponse)
	err := c.cc.Invoke(ctx, "/userProfileRecommend.UserProfile/CreateUserProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserProfileServer is the server API for UserProfile service.
// All implementations must embed UnimplementedUserProfileServer
// for forward compatibility
type UserProfileServer interface {
	CreateUserProfile(context.Context, *ProfileRequest) (*ProfileResponse, error)
	mustEmbedUnimplementedUserProfileServer()
}

// UnimplementedUserProfileServer must be embedded to have forward compatible implementations.
type UnimplementedUserProfileServer struct {
}

func (UnimplementedUserProfileServer) CreateUserProfile(context.Context, *ProfileRequest) (*ProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserProfile not implemented")
}
func (UnimplementedUserProfileServer) mustEmbedUnimplementedUserProfileServer() {}

// UnsafeUserProfileServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserProfileServer will
// result in compilation errors.
type UnsafeUserProfileServer interface {
	mustEmbedUnimplementedUserProfileServer()
}

func RegisterUserProfileServer(s grpc.ServiceRegistrar, srv UserProfileServer) {
	s.RegisterService(&UserProfile_ServiceDesc, srv)
}

func _UserProfile_CreateUserProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserProfileServer).CreateUserProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userProfileRecommend.UserProfile/CreateUserProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserProfileServer).CreateUserProfile(ctx, req.(*ProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserProfile_ServiceDesc is the grpc.ServiceDesc for UserProfile service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserProfile_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "userProfileRecommend.UserProfile",
	HandlerType: (*UserProfileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUserProfile",
			Handler:    _UserProfile_CreateUserProfile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/userProfileRecommend/userProfileRecommend.proto",
}
