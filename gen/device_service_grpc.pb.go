// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: device_service.proto

package wgpb

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DeviceServiceClient is the client API for DeviceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeviceServiceClient interface {
	AddDevice(ctx context.Context, in *AddDeviceRequest, opts ...grpc.CallOption) (*EntityIdRequest, error)
	RemoveDevice(ctx context.Context, in *EntityIdRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	UpdateDevice(ctx context.Context, in *UpdateDeviceRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	GetDevice(ctx context.Context, in *EntityIdRequest, opts ...grpc.CallOption) (*Device, error)
	GetDevices(ctx context.Context, in *GetDevicesRequest, opts ...grpc.CallOption) (*GetDevicesResponse, error)
}

type deviceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceServiceClient(cc grpc.ClientConnInterface) DeviceServiceClient {
	return &deviceServiceClient{cc}
}

func (c *deviceServiceClient) AddDevice(ctx context.Context, in *AddDeviceRequest, opts ...grpc.CallOption) (*EntityIdRequest, error) {
	out := new(EntityIdRequest)
	err := c.cc.Invoke(ctx, "/DeviceService/AddDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) RemoveDevice(ctx context.Context, in *EntityIdRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/DeviceService/RemoveDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) UpdateDevice(ctx context.Context, in *UpdateDeviceRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/DeviceService/UpdateDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) GetDevice(ctx context.Context, in *EntityIdRequest, opts ...grpc.CallOption) (*Device, error) {
	out := new(Device)
	err := c.cc.Invoke(ctx, "/DeviceService/GetDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) GetDevices(ctx context.Context, in *GetDevicesRequest, opts ...grpc.CallOption) (*GetDevicesResponse, error) {
	out := new(GetDevicesResponse)
	err := c.cc.Invoke(ctx, "/DeviceService/GetDevices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceServiceServer is the server API for DeviceService service.
// All implementations must embed UnimplementedDeviceServiceServer
// for forward compatibility
type DeviceServiceServer interface {
	AddDevice(context.Context, *AddDeviceRequest) (*EntityIdRequest, error)
	RemoveDevice(context.Context, *EntityIdRequest) (*empty.Empty, error)
	UpdateDevice(context.Context, *UpdateDeviceRequest) (*empty.Empty, error)
	GetDevice(context.Context, *EntityIdRequest) (*Device, error)
	GetDevices(context.Context, *GetDevicesRequest) (*GetDevicesResponse, error)
	mustEmbedUnimplementedDeviceServiceServer()
}

// UnimplementedDeviceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDeviceServiceServer struct {
}

func (UnimplementedDeviceServiceServer) AddDevice(context.Context, *AddDeviceRequest) (*EntityIdRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddDevice not implemented")
}
func (UnimplementedDeviceServiceServer) RemoveDevice(context.Context, *EntityIdRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveDevice not implemented")
}
func (UnimplementedDeviceServiceServer) UpdateDevice(context.Context, *UpdateDeviceRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDevice not implemented")
}
func (UnimplementedDeviceServiceServer) GetDevice(context.Context, *EntityIdRequest) (*Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDevice not implemented")
}
func (UnimplementedDeviceServiceServer) GetDevices(context.Context, *GetDevicesRequest) (*GetDevicesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDevices not implemented")
}
func (UnimplementedDeviceServiceServer) mustEmbedUnimplementedDeviceServiceServer() {}

// UnsafeDeviceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeviceServiceServer will
// result in compilation errors.
type UnsafeDeviceServiceServer interface {
	mustEmbedUnimplementedDeviceServiceServer()
}

func RegisterDeviceServiceServer(s grpc.ServiceRegistrar, srv DeviceServiceServer) {
	s.RegisterService(&DeviceService_ServiceDesc, srv)
}

func _DeviceService_AddDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).AddDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DeviceService/AddDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).AddDevice(ctx, req.(*AddDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_RemoveDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntityIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).RemoveDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DeviceService/RemoveDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).RemoveDevice(ctx, req.(*EntityIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_UpdateDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).UpdateDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DeviceService/UpdateDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).UpdateDevice(ctx, req.(*UpdateDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_GetDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntityIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).GetDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DeviceService/GetDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).GetDevice(ctx, req.(*EntityIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_GetDevices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDevicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).GetDevices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DeviceService/GetDevices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).GetDevices(ctx, req.(*GetDevicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DeviceService_ServiceDesc is the grpc.ServiceDesc for DeviceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeviceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "DeviceService",
	HandlerType: (*DeviceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddDevice",
			Handler:    _DeviceService_AddDevice_Handler,
		},
		{
			MethodName: "RemoveDevice",
			Handler:    _DeviceService_RemoveDevice_Handler,
		},
		{
			MethodName: "UpdateDevice",
			Handler:    _DeviceService_UpdateDevice_Handler,
		},
		{
			MethodName: "GetDevice",
			Handler:    _DeviceService_GetDevice_Handler,
		},
		{
			MethodName: "GetDevices",
			Handler:    _DeviceService_GetDevices_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "device_service.proto",
}
