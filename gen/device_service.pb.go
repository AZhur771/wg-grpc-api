// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: device_service.proto

package wgpb

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Device struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                  string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Description         string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Name                string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Type                string `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	PublicKey           string `protobuf:"bytes,5,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	FirewallMark        int32  `protobuf:"varint,6,opt,name=firewall_mark,json=firewallMark,proto3" json:"firewall_mark,omitempty"`
	MaxPeersCount       int32  `protobuf:"varint,7,opt,name=max_peers_count,json=maxPeersCount,proto3" json:"max_peers_count,omitempty"`
	CurrentPeersCount   int32  `protobuf:"varint,8,opt,name=current_peers_count,json=currentPeersCount,proto3" json:"current_peers_count,omitempty"`
	Endpoint            string `protobuf:"bytes,9,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Address             string `protobuf:"bytes,10,opt,name=address,proto3" json:"address,omitempty"`
	Mtu                 int32  `protobuf:"varint,11,opt,name=mtu,proto3" json:"mtu,omitempty"`
	Dns                 string `protobuf:"bytes,12,opt,name=dns,proto3" json:"dns,omitempty"`
	Table               string `protobuf:"bytes,13,opt,name=table,proto3" json:"table,omitempty"`
	PersistentKeepAlive int32  `protobuf:"varint,14,opt,name=persistent_keep_alive,json=persistentKeepAlive,proto3" json:"persistent_keep_alive,omitempty"`
	PreUp               string `protobuf:"bytes,15,opt,name=pre_up,json=preUp,proto3" json:"pre_up,omitempty"`
	PreDown             string `protobuf:"bytes,16,opt,name=pre_down,json=preDown,proto3" json:"pre_down,omitempty"`
	PostUp              string `protobuf:"bytes,17,opt,name=post_up,json=postUp,proto3" json:"post_up,omitempty"`
	PostDown            string `protobuf:"bytes,18,opt,name=post_down,json=postDown,proto3" json:"post_down,omitempty"`
}

func (x *Device) Reset() {
	*x = Device{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Device) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Device) ProtoMessage() {}

func (x *Device) ProtoReflect() protoreflect.Message {
	mi := &file_device_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Device.ProtoReflect.Descriptor instead.
func (*Device) Descriptor() ([]byte, []int) {
	return file_device_service_proto_rawDescGZIP(), []int{0}
}

func (x *Device) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Device) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Device) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Device) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Device) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

func (x *Device) GetFirewallMark() int32 {
	if x != nil {
		return x.FirewallMark
	}
	return 0
}

func (x *Device) GetMaxPeersCount() int32 {
	if x != nil {
		return x.MaxPeersCount
	}
	return 0
}

func (x *Device) GetCurrentPeersCount() int32 {
	if x != nil {
		return x.CurrentPeersCount
	}
	return 0
}

func (x *Device) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Device) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Device) GetMtu() int32 {
	if x != nil {
		return x.Mtu
	}
	return 0
}

func (x *Device) GetDns() string {
	if x != nil {
		return x.Dns
	}
	return ""
}

func (x *Device) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *Device) GetPersistentKeepAlive() int32 {
	if x != nil {
		return x.PersistentKeepAlive
	}
	return 0
}

func (x *Device) GetPreUp() string {
	if x != nil {
		return x.PreUp
	}
	return ""
}

func (x *Device) GetPreDown() string {
	if x != nil {
		return x.PreDown
	}
	return ""
}

func (x *Device) GetPostUp() string {
	if x != nil {
		return x.PostUp
	}
	return ""
}

func (x *Device) GetPostDown() string {
	if x != nil {
		return x.PostDown
	}
	return ""
}

type AddDeviceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description         string `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	FirewallMark        int32  `protobuf:"varint,2,opt,name=firewall_mark,json=firewallMark,proto3" json:"firewall_mark,omitempty"`
	Endpoint            string `protobuf:"bytes,3,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Address             string `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
	Mtu                 int32  `protobuf:"varint,5,opt,name=mtu,proto3" json:"mtu,omitempty"`
	Dns                 string `protobuf:"bytes,6,opt,name=dns,proto3" json:"dns,omitempty"`
	Table               string `protobuf:"bytes,7,opt,name=table,proto3" json:"table,omitempty"`
	PersistentKeepAlive int32  `protobuf:"varint,8,opt,name=persistent_keep_alive,json=persistentKeepAlive,proto3" json:"persistent_keep_alive,omitempty"`
	PreUp               string `protobuf:"bytes,9,opt,name=pre_up,json=preUp,proto3" json:"pre_up,omitempty"`
	PreDown             string `protobuf:"bytes,10,opt,name=pre_down,json=preDown,proto3" json:"pre_down,omitempty"`
	PostUp              string `protobuf:"bytes,11,opt,name=post_up,json=postUp,proto3" json:"post_up,omitempty"`
	PostDown            string `protobuf:"bytes,12,opt,name=post_down,json=postDown,proto3" json:"post_down,omitempty"`
}

func (x *AddDeviceRequest) Reset() {
	*x = AddDeviceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddDeviceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddDeviceRequest) ProtoMessage() {}

func (x *AddDeviceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_device_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddDeviceRequest.ProtoReflect.Descriptor instead.
func (*AddDeviceRequest) Descriptor() ([]byte, []int) {
	return file_device_service_proto_rawDescGZIP(), []int{1}
}

func (x *AddDeviceRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *AddDeviceRequest) GetFirewallMark() int32 {
	if x != nil {
		return x.FirewallMark
	}
	return 0
}

func (x *AddDeviceRequest) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *AddDeviceRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *AddDeviceRequest) GetMtu() int32 {
	if x != nil {
		return x.Mtu
	}
	return 0
}

func (x *AddDeviceRequest) GetDns() string {
	if x != nil {
		return x.Dns
	}
	return ""
}

func (x *AddDeviceRequest) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *AddDeviceRequest) GetPersistentKeepAlive() int32 {
	if x != nil {
		return x.PersistentKeepAlive
	}
	return 0
}

func (x *AddDeviceRequest) GetPreUp() string {
	if x != nil {
		return x.PreUp
	}
	return ""
}

func (x *AddDeviceRequest) GetPreDown() string {
	if x != nil {
		return x.PreDown
	}
	return ""
}

func (x *AddDeviceRequest) GetPostUp() string {
	if x != nil {
		return x.PostUp
	}
	return ""
}

func (x *AddDeviceRequest) GetPostDown() string {
	if x != nil {
		return x.PostDown
	}
	return ""
}

type UpdateDeviceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                  string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Description         string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	FirewallMark        int32  `protobuf:"varint,3,opt,name=firewall_mark,json=firewallMark,proto3" json:"firewall_mark,omitempty"`
	Endpoint            string `protobuf:"bytes,4,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Address             string `protobuf:"bytes,5,opt,name=address,proto3" json:"address,omitempty"`
	Mtu                 int32  `protobuf:"varint,6,opt,name=mtu,proto3" json:"mtu,omitempty"`
	Dns                 string `protobuf:"bytes,7,opt,name=dns,proto3" json:"dns,omitempty"`
	Table               string `protobuf:"bytes,8,opt,name=table,proto3" json:"table,omitempty"`
	PersistentKeepAlive int32  `protobuf:"varint,9,opt,name=persistent_keep_alive,json=persistentKeepAlive,proto3" json:"persistent_keep_alive,omitempty"`
	PreUp               string `protobuf:"bytes,10,opt,name=pre_up,json=preUp,proto3" json:"pre_up,omitempty"`
	PreDown             string `protobuf:"bytes,11,opt,name=pre_down,json=preDown,proto3" json:"pre_down,omitempty"`
	PostUp              string `protobuf:"bytes,12,opt,name=post_up,json=postUp,proto3" json:"post_up,omitempty"`
	PostDown            string `protobuf:"bytes,13,opt,name=post_down,json=postDown,proto3" json:"post_down,omitempty"`
}

func (x *UpdateDeviceRequest) Reset() {
	*x = UpdateDeviceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateDeviceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDeviceRequest) ProtoMessage() {}

func (x *UpdateDeviceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_device_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDeviceRequest.ProtoReflect.Descriptor instead.
func (*UpdateDeviceRequest) Descriptor() ([]byte, []int) {
	return file_device_service_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateDeviceRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateDeviceRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UpdateDeviceRequest) GetFirewallMark() int32 {
	if x != nil {
		return x.FirewallMark
	}
	return 0
}

func (x *UpdateDeviceRequest) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *UpdateDeviceRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *UpdateDeviceRequest) GetMtu() int32 {
	if x != nil {
		return x.Mtu
	}
	return 0
}

func (x *UpdateDeviceRequest) GetDns() string {
	if x != nil {
		return x.Dns
	}
	return ""
}

func (x *UpdateDeviceRequest) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *UpdateDeviceRequest) GetPersistentKeepAlive() int32 {
	if x != nil {
		return x.PersistentKeepAlive
	}
	return 0
}

func (x *UpdateDeviceRequest) GetPreUp() string {
	if x != nil {
		return x.PreUp
	}
	return ""
}

func (x *UpdateDeviceRequest) GetPreDown() string {
	if x != nil {
		return x.PreDown
	}
	return ""
}

func (x *UpdateDeviceRequest) GetPostUp() string {
	if x != nil {
		return x.PostUp
	}
	return ""
}

func (x *UpdateDeviceRequest) GetPostDown() string {
	if x != nil {
		return x.PostDown
	}
	return ""
}

type GetDevicesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Skip   int32  `protobuf:"varint,1,opt,name=skip,proto3" json:"skip,omitempty"`
	Limit  int32  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Search string `protobuf:"bytes,3,opt,name=search,proto3" json:"search,omitempty"`
}

func (x *GetDevicesRequest) Reset() {
	*x = GetDevicesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDevicesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDevicesRequest) ProtoMessage() {}

func (x *GetDevicesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_device_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDevicesRequest.ProtoReflect.Descriptor instead.
func (*GetDevicesRequest) Descriptor() ([]byte, []int) {
	return file_device_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetDevicesRequest) GetSkip() int32 {
	if x != nil {
		return x.Skip
	}
	return 0
}

func (x *GetDevicesRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *GetDevicesRequest) GetSearch() string {
	if x != nil {
		return x.Search
	}
	return ""
}

type GetDevicesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Devices []*Device `protobuf:"bytes,1,rep,name=devices,proto3" json:"devices,omitempty"`
	Total   int32     `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
	HasNext bool      `protobuf:"varint,3,opt,name=has_next,json=hasNext,proto3" json:"has_next,omitempty"`
}

func (x *GetDevicesResponse) Reset() {
	*x = GetDevicesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDevicesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDevicesResponse) ProtoMessage() {}

func (x *GetDevicesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_device_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDevicesResponse.ProtoReflect.Descriptor instead.
func (*GetDevicesResponse) Descriptor() ([]byte, []int) {
	return file_device_service_proto_rawDescGZIP(), []int{4}
}

func (x *GetDevicesResponse) GetDevices() []*Device {
	if x != nil {
		return x.Devices
	}
	return nil
}

func (x *GetDevicesResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *GetDevicesResponse) GetHasNext() bool {
	if x != nil {
		return x.HasNext
	}
	return false
}

var File_device_service_proto protoreflect.FileDescriptor

var file_device_service_proto_rawDesc = []byte{
	0x0a, 0x14, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d,
	0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x15, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69,
	0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8a, 0x04, 0x0a, 0x06, 0x44, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x23, 0x0a, 0x0d,
	0x66, 0x69, 0x72, 0x65, 0x77, 0x61, 0x6c, 0x6c, 0x5f, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0c, 0x66, 0x69, 0x72, 0x65, 0x77, 0x61, 0x6c, 0x6c, 0x4d, 0x61, 0x72,
	0x6b, 0x12, 0x26, 0x0a, 0x0f, 0x6d, 0x61, 0x78, 0x5f, 0x70, 0x65, 0x65, 0x72, 0x73, 0x5f, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x6d, 0x61, 0x78, 0x50,
	0x65, 0x65, 0x72, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x13, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x65, 0x65, 0x72, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50,
	0x65, 0x65, 0x72, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12,
	0x10, 0x0a, 0x03, 0x6d, 0x74, 0x75, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6d, 0x74,
	0x75, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x6e, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x64, 0x6e, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x0d, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x32, 0x0a, 0x15, 0x70, 0x65, 0x72,
	0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x6b, 0x65, 0x65, 0x70, 0x5f, 0x61, 0x6c, 0x69,
	0x76, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x05, 0x52, 0x13, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73,
	0x74, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x65, 0x70, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x12, 0x15, 0x0a,
	0x06, 0x70, 0x72, 0x65, 0x5f, 0x75, 0x70, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70,
	0x72, 0x65, 0x55, 0x70, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x72, 0x65, 0x5f, 0x64, 0x6f, 0x77, 0x6e,
	0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x65, 0x44, 0x6f, 0x77, 0x6e, 0x12,
	0x17, 0x0a, 0x07, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x75, 0x70, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x55, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6f, 0x73, 0x74,
	0x5f, 0x64, 0x6f, 0x77, 0x6e, 0x18, 0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6f, 0x73,
	0x74, 0x44, 0x6f, 0x77, 0x6e, 0x22, 0xe5, 0x02, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x44, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d,
	0x66, 0x69, 0x72, 0x65, 0x77, 0x61, 0x6c, 0x6c, 0x5f, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0c, 0x66, 0x69, 0x72, 0x65, 0x77, 0x61, 0x6c, 0x6c, 0x4d, 0x61, 0x72,
	0x6b, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x74, 0x75, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6d, 0x74, 0x75, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x6e, 0x73,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x64, 0x6e, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x12, 0x32, 0x0a, 0x15, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x5f,
	0x6b, 0x65, 0x65, 0x70, 0x5f, 0x61, 0x6c, 0x69, 0x76, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x13, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x65, 0x70,
	0x41, 0x6c, 0x69, 0x76, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x70, 0x72, 0x65, 0x5f, 0x75, 0x70, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x72, 0x65, 0x55, 0x70, 0x12, 0x19, 0x0a, 0x08,
	0x70, 0x72, 0x65, 0x5f, 0x64, 0x6f, 0x77, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x70, 0x72, 0x65, 0x44, 0x6f, 0x77, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x6f, 0x73, 0x74, 0x5f,
	0x75, 0x70, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x55, 0x70,
	0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x64, 0x6f, 0x77, 0x6e, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x74, 0x44, 0x6f, 0x77, 0x6e, 0x22, 0xf8, 0x02,
	0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x66, 0x69, 0x72, 0x65, 0x77,
	0x61, 0x6c, 0x6c, 0x5f, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c,
	0x66, 0x69, 0x72, 0x65, 0x77, 0x61, 0x6c, 0x6c, 0x4d, 0x61, 0x72, 0x6b, 0x12, 0x1a, 0x0a, 0x08,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x74, 0x75, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x6d, 0x74, 0x75, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x6e, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x64, 0x6e, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x32, 0x0a, 0x15,
	0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x6b, 0x65, 0x65, 0x70, 0x5f,
	0x61, 0x6c, 0x69, 0x76, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x13, 0x70, 0x65, 0x72,
	0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x65, 0x70, 0x41, 0x6c, 0x69, 0x76, 0x65,
	0x12, 0x15, 0x0a, 0x06, 0x70, 0x72, 0x65, 0x5f, 0x75, 0x70, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x70, 0x72, 0x65, 0x55, 0x70, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x72, 0x65, 0x5f, 0x64,
	0x6f, 0x77, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x65, 0x44, 0x6f,
	0x77, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x75, 0x70, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x55, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x70,
	0x6f, 0x73, 0x74, 0x5f, 0x64, 0x6f, 0x77, 0x6e, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x70, 0x6f, 0x73, 0x74, 0x44, 0x6f, 0x77, 0x6e, 0x22, 0x55, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x6b, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x6b, 0x69,
	0x70, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x22,
	0x68, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x07, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52,
	0x07, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x19,
	0x0a, 0x08, 0x68, 0x61, 0x73, 0x5f, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x68, 0x61, 0x73, 0x4e, 0x65, 0x78, 0x74, 0x32, 0xf0, 0x06, 0x0a, 0x0d, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x94, 0x01, 0x0a, 0x09,
	0x41, 0x64, 0x64, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x11, 0x2e, 0x41, 0x64, 0x64, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x62,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x22, 0x0c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x3a, 0x01, 0x2a, 0x92, 0x41, 0x48, 0x0a, 0x0d, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x0a, 0x41, 0x64, 0x64, 0x20, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x19, 0x41, 0x64, 0x64, 0x20, 0x64, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x20, 0x74, 0x6f, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x62, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x41, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x74, 0x68,
	0x12, 0x00, 0x12, 0xb4, 0x01, 0x0a, 0x0c, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x44, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x10, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x7a, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x15, 0x2a, 0x10, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x3a, 0x01, 0x2a, 0x92, 0x41, 0x5c, 0x0a, 0x0d, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x13, 0x52, 0x65,
	0x6d, 0x6f, 0x76, 0x65, 0x20, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x20, 0x62, 0x79, 0x20, 0x69,
	0x64, 0x1a, 0x24, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x20, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x20, 0x62, 0x79, 0x20, 0x69, 0x64, 0x20, 0x66, 0x72, 0x6f, 0x6d, 0x20, 0x74, 0x68, 0x65, 0x20,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x62, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x41, 0x70, 0x69,
	0x4b, 0x65, 0x79, 0x41, 0x75, 0x74, 0x68, 0x12, 0x00, 0x12, 0xb6, 0x01, 0x0a, 0x0c, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x14, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x78, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15,
	0x1a, 0x10, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x7b, 0x69,
	0x64, 0x7d, 0x3a, 0x01, 0x2a, 0x92, 0x41, 0x5a, 0x0a, 0x0d, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x20, 0x62, 0x79, 0x20, 0x69, 0x64, 0x1a, 0x22, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x20, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x20, 0x62, 0x79, 0x20, 0x69,
	0x64, 0x20, 0x6f, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x62, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x41, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x74, 0x68,
	0x12, 0x00, 0x12, 0x8d, 0x01, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x10, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x07, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x22, 0x65, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x12, 0x12, 0x10, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x92, 0x41, 0x4a, 0x0a, 0x0d, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x0a, 0x47, 0x65, 0x74, 0x20, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x1a, 0x1b, 0x47, 0x65, 0x74, 0x20, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x20,
	0x66, 0x72, 0x6f, 0x6d, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x62, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x41, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x74, 0x68,
	0x12, 0x00, 0x12, 0x9a, 0x01, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x12, 0x12, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x63, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x0e, 0x12, 0x0c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x92, 0x41, 0x4c, 0x0a, 0x0d, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x0b, 0x47, 0x65, 0x74, 0x20, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a,
	0x1c, 0x47, 0x65, 0x74, 0x20, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x20, 0x66, 0x72, 0x6f,
	0x6d, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x62, 0x10, 0x0a,
	0x0e, 0x0a, 0x0a, 0x41, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x74, 0x68, 0x12, 0x00, 0x1a,
	0x2b, 0x92, 0x41, 0x28, 0x12, 0x26, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x20, 0x74, 0x6f,
	0x20, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x20, 0x77, 0x69, 0x72, 0x65, 0x67,
	0x75, 0x61, 0x72, 0x64, 0x20, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x42, 0x09, 0x5a, 0x07,
	0x2e, 0x2f, 0x3b, 0x77, 0x67, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_device_service_proto_rawDescOnce sync.Once
	file_device_service_proto_rawDescData = file_device_service_proto_rawDesc
)

func file_device_service_proto_rawDescGZIP() []byte {
	file_device_service_proto_rawDescOnce.Do(func() {
		file_device_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_device_service_proto_rawDescData)
	})
	return file_device_service_proto_rawDescData
}

var file_device_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_device_service_proto_goTypes = []interface{}{
	(*Device)(nil),              // 0: Device
	(*AddDeviceRequest)(nil),    // 1: AddDeviceRequest
	(*UpdateDeviceRequest)(nil), // 2: UpdateDeviceRequest
	(*GetDevicesRequest)(nil),   // 3: GetDevicesRequest
	(*GetDevicesResponse)(nil),  // 4: GetDevicesResponse
	(*EntityIdRequest)(nil),     // 5: EntityIdRequest
	(*empty.Empty)(nil),         // 6: google.protobuf.Empty
}
var file_device_service_proto_depIdxs = []int32{
	0, // 0: GetDevicesResponse.devices:type_name -> Device
	1, // 1: DeviceService.AddDevice:input_type -> AddDeviceRequest
	5, // 2: DeviceService.RemoveDevice:input_type -> EntityIdRequest
	2, // 3: DeviceService.UpdateDevice:input_type -> UpdateDeviceRequest
	5, // 4: DeviceService.GetDevice:input_type -> EntityIdRequest
	3, // 5: DeviceService.GetDevices:input_type -> GetDevicesRequest
	5, // 6: DeviceService.AddDevice:output_type -> EntityIdRequest
	6, // 7: DeviceService.RemoveDevice:output_type -> google.protobuf.Empty
	6, // 8: DeviceService.UpdateDevice:output_type -> google.protobuf.Empty
	0, // 9: DeviceService.GetDevice:output_type -> Device
	4, // 10: DeviceService.GetDevices:output_type -> GetDevicesResponse
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_device_service_proto_init() }
func file_device_service_proto_init() {
	if File_device_service_proto != nil {
		return
	}
	file_common_entities_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_device_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Device); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_device_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddDeviceRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_device_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateDeviceRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_device_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDevicesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_device_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDevicesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_device_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_device_service_proto_goTypes,
		DependencyIndexes: file_device_service_proto_depIdxs,
		MessageInfos:      file_device_service_proto_msgTypes,
	}.Build()
	File_device_service_proto = out.File
	file_device_service_proto_rawDesc = nil
	file_device_service_proto_goTypes = nil
	file_device_service_proto_depIdxs = nil
}
