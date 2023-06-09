// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: account.proto

package v1

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

// AccountClient is the client API for Account service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountClient interface {
	Create(ctx context.Context, in *UpdateAccountRequest, opts ...grpc.CallOption) (*AccountResponse, error)
	Update(ctx context.Context, in *UpdateAccountRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	LoginPassword(ctx context.Context, in *LoginPasswordRequest, opts ...grpc.CallOption) (*AccountResponse, error)
	LoginMobile(ctx context.Context, in *LoginMobileRequest, opts ...grpc.CallOption) (*AccountResponse, error)
	LoginPlatform(ctx context.Context, in *LoginPlatformRequest, opts ...grpc.CallOption) (*AccountResponse, error)
	BindPlatform(ctx context.Context, in *BindPlatformRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UnbindPlatform(ctx context.Context, in *BindPlatformRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type accountClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountClient(cc grpc.ClientConnInterface) AccountClient {
	return &accountClient{cc}
}

func (c *accountClient) Create(ctx context.Context, in *UpdateAccountRequest, opts ...grpc.CallOption) (*AccountResponse, error) {
	out := new(AccountResponse)
	err := c.cc.Invoke(ctx, "/userPb.v1.Account/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) Update(ctx context.Context, in *UpdateAccountRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/userPb.v1.Account/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) LoginPassword(ctx context.Context, in *LoginPasswordRequest, opts ...grpc.CallOption) (*AccountResponse, error) {
	out := new(AccountResponse)
	err := c.cc.Invoke(ctx, "/userPb.v1.Account/LoginPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) LoginMobile(ctx context.Context, in *LoginMobileRequest, opts ...grpc.CallOption) (*AccountResponse, error) {
	out := new(AccountResponse)
	err := c.cc.Invoke(ctx, "/userPb.v1.Account/LoginMobile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) LoginPlatform(ctx context.Context, in *LoginPlatformRequest, opts ...grpc.CallOption) (*AccountResponse, error) {
	out := new(AccountResponse)
	err := c.cc.Invoke(ctx, "/userPb.v1.Account/LoginPlatform", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) BindPlatform(ctx context.Context, in *BindPlatformRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/userPb.v1.Account/BindPlatform", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) UnbindPlatform(ctx context.Context, in *BindPlatformRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/userPb.v1.Account/UnbindPlatform", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServer is the server API for Account service.
// All implementations must embed UnimplementedAccountServer
// for forward compatibility
type AccountServer interface {
	Create(context.Context, *UpdateAccountRequest) (*AccountResponse, error)
	Update(context.Context, *UpdateAccountRequest) (*emptypb.Empty, error)
	LoginPassword(context.Context, *LoginPasswordRequest) (*AccountResponse, error)
	LoginMobile(context.Context, *LoginMobileRequest) (*AccountResponse, error)
	LoginPlatform(context.Context, *LoginPlatformRequest) (*AccountResponse, error)
	BindPlatform(context.Context, *BindPlatformRequest) (*emptypb.Empty, error)
	UnbindPlatform(context.Context, *BindPlatformRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedAccountServer()
}

// UnimplementedAccountServer must be embedded to have forward compatible implementations.
type UnimplementedAccountServer struct {
}

func (UnimplementedAccountServer) Create(context.Context, *UpdateAccountRequest) (*AccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedAccountServer) Update(context.Context, *UpdateAccountRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedAccountServer) LoginPassword(context.Context, *LoginPasswordRequest) (*AccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginPassword not implemented")
}
func (UnimplementedAccountServer) LoginMobile(context.Context, *LoginMobileRequest) (*AccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginMobile not implemented")
}
func (UnimplementedAccountServer) LoginPlatform(context.Context, *LoginPlatformRequest) (*AccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginPlatform not implemented")
}
func (UnimplementedAccountServer) BindPlatform(context.Context, *BindPlatformRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BindPlatform not implemented")
}
func (UnimplementedAccountServer) UnbindPlatform(context.Context, *BindPlatformRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnbindPlatform not implemented")
}
func (UnimplementedAccountServer) mustEmbedUnimplementedAccountServer() {}

// UnsafeAccountServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountServer will
// result in compilation errors.
type UnsafeAccountServer interface {
	mustEmbedUnimplementedAccountServer()
}

func RegisterAccountServer(s grpc.ServiceRegistrar, srv AccountServer) {
	s.RegisterService(&Account_ServiceDesc, srv)
}

func _Account_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userPb.v1.Account/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).Create(ctx, req.(*UpdateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userPb.v1.Account/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).Update(ctx, req.(*UpdateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_LoginPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).LoginPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userPb.v1.Account/LoginPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).LoginPassword(ctx, req.(*LoginPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_LoginMobile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginMobileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).LoginMobile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userPb.v1.Account/LoginMobile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).LoginMobile(ctx, req.(*LoginMobileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_LoginPlatform_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginPlatformRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).LoginPlatform(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userPb.v1.Account/LoginPlatform",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).LoginPlatform(ctx, req.(*LoginPlatformRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_BindPlatform_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BindPlatformRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).BindPlatform(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userPb.v1.Account/BindPlatform",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).BindPlatform(ctx, req.(*BindPlatformRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_UnbindPlatform_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BindPlatformRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).UnbindPlatform(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userPb.v1.Account/UnbindPlatform",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).UnbindPlatform(ctx, req.(*BindPlatformRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Account_ServiceDesc is the grpc.ServiceDesc for Account service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Account_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "userPb.v1.Account",
	HandlerType: (*AccountServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Account_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Account_Update_Handler,
		},
		{
			MethodName: "LoginPassword",
			Handler:    _Account_LoginPassword_Handler,
		},
		{
			MethodName: "LoginMobile",
			Handler:    _Account_LoginMobile_Handler,
		},
		{
			MethodName: "LoginPlatform",
			Handler:    _Account_LoginPlatform_Handler,
		},
		{
			MethodName: "BindPlatform",
			Handler:    _Account_BindPlatform_Handler,
		},
		{
			MethodName: "UnbindPlatform",
			Handler:    _Account_UnbindPlatform_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "account.proto",
}
