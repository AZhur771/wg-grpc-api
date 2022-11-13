// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: peer_service.proto

package peerpb

import (
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Peer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PublicKey           string               `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	HasPresharedKey     bool                 `protobuf:"varint,2,opt,name=has_preshared_key,json=hasPresharedKey,proto3" json:"has_preshared_key,omitempty"`
	Endpoint            string               `protobuf:"bytes,3,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	PersistentKeepAlive string               `protobuf:"bytes,4,opt,name=persistent_keep_alive,json=persistentKeepAlive,proto3" json:"persistent_keep_alive,omitempty"`
	LastHandshake       *timestamp.Timestamp `protobuf:"bytes,5,opt,name=last_handshake,json=lastHandshake,proto3" json:"last_handshake,omitempty"`
	ReceiveBytes        int64                `protobuf:"varint,6,opt,name=receive_bytes,json=receiveBytes,proto3" json:"receive_bytes,omitempty"`
	TransmitBytes       int64                `protobuf:"varint,7,opt,name=transmit_bytes,json=transmitBytes,proto3" json:"transmit_bytes,omitempty"`
	AllowedIps          []string             `protobuf:"bytes,8,rep,name=allowed_ips,json=allowedIps,proto3" json:"allowed_ips,omitempty"`
	ProtocolVersion     uint32               `protobuf:"varint,9,opt,name=protocol_version,json=protocolVersion,proto3" json:"protocol_version,omitempty"`
	IsEnabled           bool                 `protobuf:"varint,10,opt,name=is_enabled,json=isEnabled,proto3" json:"is_enabled,omitempty"`
	IsActive            bool                 `protobuf:"varint,11,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
}

func (x *Peer) Reset() {
	*x = Peer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Peer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Peer) ProtoMessage() {}

func (x *Peer) ProtoReflect() protoreflect.Message {
	mi := &file_peer_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Peer.ProtoReflect.Descriptor instead.
func (*Peer) Descriptor() ([]byte, []int) {
	return file_peer_service_proto_rawDescGZIP(), []int{0}
}

func (x *Peer) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

func (x *Peer) GetHasPresharedKey() bool {
	if x != nil {
		return x.HasPresharedKey
	}
	return false
}

func (x *Peer) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Peer) GetPersistentKeepAlive() string {
	if x != nil {
		return x.PersistentKeepAlive
	}
	return ""
}

func (x *Peer) GetLastHandshake() *timestamp.Timestamp {
	if x != nil {
		return x.LastHandshake
	}
	return nil
}

func (x *Peer) GetReceiveBytes() int64 {
	if x != nil {
		return x.ReceiveBytes
	}
	return 0
}

func (x *Peer) GetTransmitBytes() int64 {
	if x != nil {
		return x.TransmitBytes
	}
	return 0
}

func (x *Peer) GetAllowedIps() []string {
	if x != nil {
		return x.AllowedIps
	}
	return nil
}

func (x *Peer) GetProtocolVersion() uint32 {
	if x != nil {
		return x.ProtocolVersion
	}
	return 0
}

func (x *Peer) GetIsEnabled() bool {
	if x != nil {
		return x.IsEnabled
	}
	return false
}

func (x *Peer) GetIsActive() bool {
	if x != nil {
		return x.IsActive
	}
	return false
}

type AddPeerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddPeerRequest) Reset() {
	*x = AddPeerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddPeerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddPeerRequest) ProtoMessage() {}

func (x *AddPeerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_peer_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddPeerRequest.ProtoReflect.Descriptor instead.
func (*AddPeerRequest) Descriptor() ([]byte, []int) {
	return file_peer_service_proto_rawDescGZIP(), []int{1}
}

type AddPeerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddPeerResponse) Reset() {
	*x = AddPeerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddPeerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddPeerResponse) ProtoMessage() {}

func (x *AddPeerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_peer_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddPeerResponse.ProtoReflect.Descriptor instead.
func (*AddPeerResponse) Descriptor() ([]byte, []int) {
	return file_peer_service_proto_rawDescGZIP(), []int{2}
}

type RemovePeerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PublicKey string `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
}

func (x *RemovePeerRequest) Reset() {
	*x = RemovePeerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemovePeerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemovePeerRequest) ProtoMessage() {}

func (x *RemovePeerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_peer_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemovePeerRequest.ProtoReflect.Descriptor instead.
func (*RemovePeerRequest) Descriptor() ([]byte, []int) {
	return file_peer_service_proto_rawDescGZIP(), []int{3}
}

func (x *RemovePeerRequest) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

type RemovePeerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RemovePeerResponse) Reset() {
	*x = RemovePeerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemovePeerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemovePeerResponse) ProtoMessage() {}

func (x *RemovePeerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_peer_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemovePeerResponse.ProtoReflect.Descriptor instead.
func (*RemovePeerResponse) Descriptor() ([]byte, []int) {
	return file_peer_service_proto_rawDescGZIP(), []int{4}
}

type UpdatePeerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PublicKey string `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
}

func (x *UpdatePeerRequest) Reset() {
	*x = UpdatePeerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePeerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePeerRequest) ProtoMessage() {}

func (x *UpdatePeerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_peer_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePeerRequest.ProtoReflect.Descriptor instead.
func (*UpdatePeerRequest) Descriptor() ([]byte, []int) {
	return file_peer_service_proto_rawDescGZIP(), []int{5}
}

func (x *UpdatePeerRequest) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

type UpdatePeerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdatePeerResponse) Reset() {
	*x = UpdatePeerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePeerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePeerResponse) ProtoMessage() {}

func (x *UpdatePeerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_peer_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePeerResponse.ProtoReflect.Descriptor instead.
func (*UpdatePeerResponse) Descriptor() ([]byte, []int) {
	return file_peer_service_proto_rawDescGZIP(), []int{6}
}

type GetPeerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PublicKey string `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
}

func (x *GetPeerRequest) Reset() {
	*x = GetPeerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPeerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPeerRequest) ProtoMessage() {}

func (x *GetPeerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_peer_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPeerRequest.ProtoReflect.Descriptor instead.
func (*GetPeerRequest) Descriptor() ([]byte, []int) {
	return file_peer_service_proto_rawDescGZIP(), []int{7}
}

func (x *GetPeerRequest) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

type GetPeersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetPeersRequest) Reset() {
	*x = GetPeersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPeersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPeersRequest) ProtoMessage() {}

func (x *GetPeersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_peer_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPeersRequest.ProtoReflect.Descriptor instead.
func (*GetPeersRequest) Descriptor() ([]byte, []int) {
	return file_peer_service_proto_rawDescGZIP(), []int{8}
}

type GetPeersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Peers []*Peer `protobuf:"bytes,1,rep,name=peers,proto3" json:"peers,omitempty"`
}

func (x *GetPeersResponse) Reset() {
	*x = GetPeersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPeersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPeersResponse) ProtoMessage() {}

func (x *GetPeersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_peer_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPeersResponse.ProtoReflect.Descriptor instead.
func (*GetPeersResponse) Descriptor() ([]byte, []int) {
	return file_peer_service_proto_rawDescGZIP(), []int{9}
}

func (x *GetPeersResponse) GetPeers() []*Peer {
	if x != nil {
		return x.Peers
	}
	return nil
}

var File_peer_service_proto protoreflect.FileDescriptor

var file_peer_service_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70,
	0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb8, 0x03, 0x0a, 0x04,
	0x50, 0x65, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x4b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x11, 0x68, 0x61, 0x73, 0x5f, 0x70, 0x72, 0x65, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x64, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f,
	0x68, 0x61, 0x73, 0x50, 0x72, 0x65, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4b, 0x65, 0x79, 0x12,
	0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x32, 0x0a, 0x15, 0x70,
	0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x6b, 0x65, 0x65, 0x70, 0x5f, 0x61,
	0x6c, 0x69, 0x76, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x70, 0x65, 0x72, 0x73,
	0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x65, 0x70, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x12,
	0x41, 0x0a, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61,
	0x6b, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x5f, 0x62, 0x79,
	0x74, 0x65, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x6d, 0x69, 0x74, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6d, 0x69, 0x74, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x1f,
	0x0a, 0x0b, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x69, 0x70, 0x73, 0x18, 0x08, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x49, 0x70, 0x73, 0x12,
	0x29, 0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x5f, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73,
	0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09,
	0x69, 0x73, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f,
	0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x22, 0x10, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x50, 0x65, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x11, 0x0a, 0x0f, 0x41, 0x64, 0x64, 0x50,
	0x65, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x32, 0x0a, 0x11, 0x52,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x50, 0x65, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x22,
	0x14, 0x0a, 0x12, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x50, 0x65, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x32, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50,
	0x65, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x22, 0x14, 0x0a, 0x12, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x50, 0x65, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x2f, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x50, 0x65, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79,
	0x22, 0x11, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x50, 0x65, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x36, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x50, 0x65, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x50, 0x65, 0x65, 0x72, 0x52, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x32, 0xb5, 0x06, 0x0a, 0x0b,
	0x50, 0x65, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x82, 0x01, 0x0a, 0x07,
	0x41, 0x64, 0x64, 0x50, 0x65, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x41, 0x64, 0x64, 0x50, 0x65, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x50, 0x65, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x46, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f,
	0x22, 0x0a, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x65, 0x65, 0x72, 0x73, 0x3a, 0x01, 0x2a, 0x92,
	0x41, 0x2e, 0x0a, 0x05, 0x50, 0x65, 0x65, 0x72, 0x73, 0x12, 0x0a, 0x41, 0x64, 0x64, 0x20, 0x61,
	0x20, 0x70, 0x65, 0x65, 0x72, 0x1a, 0x19, 0x41, 0x64, 0x64, 0x20, 0x61, 0x20, 0x70, 0x65, 0x65,
	0x72, 0x20, 0x74, 0x6f, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x12, 0xbc, 0x01, 0x0a, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x50, 0x65, 0x65, 0x72, 0x12,
	0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x50,
	0x65, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x50, 0x65, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x77, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1c, 0x2a, 0x17,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x65, 0x65, 0x72, 0x73, 0x2f, 0x7b, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x7d, 0x3a, 0x01, 0x2a, 0x92, 0x41, 0x52, 0x0a, 0x05, 0x50,
	0x65, 0x65, 0x72, 0x73, 0x12, 0x1b, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x20, 0x61, 0x20, 0x70,
	0x65, 0x65, 0x72, 0x20, 0x62, 0x79, 0x20, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x20, 0x6b, 0x65,
	0x79, 0x1a, 0x2c, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x20, 0x61, 0x20, 0x70, 0x65, 0x65, 0x72,
	0x20, 0x62, 0x79, 0x20, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x20, 0x6b, 0x65, 0x79, 0x20, 0x66,
	0x72, 0x6f, 0x6d, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x12,
	0xba, 0x01, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x65, 0x65, 0x72, 0x12, 0x19,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x65,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x65, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x75, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1c, 0x1a, 0x17, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x70, 0x65, 0x65, 0x72, 0x73, 0x2f, 0x7b, 0x70, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x5f, 0x6b, 0x65, 0x79, 0x7d, 0x3a, 0x01, 0x2a, 0x92, 0x41, 0x50, 0x0a, 0x05, 0x50, 0x65,
	0x65, 0x72, 0x73, 0x12, 0x1b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x70, 0x65,
	0x65, 0x72, 0x20, 0x62, 0x79, 0x20, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x20, 0x6b, 0x65, 0x79,
	0x1a, 0x2a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x70, 0x65, 0x65, 0x72, 0x20,
	0x62, 0x79, 0x20, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x20, 0x6b, 0x65, 0x79, 0x20, 0x6f, 0x6e,
	0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x12, 0x9f, 0x01, 0x0a,
	0x07, 0x47, 0x65, 0x74, 0x50, 0x65, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x65, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x22, 0x6e,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x12, 0x17, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x65, 0x65,
	0x72, 0x73, 0x2f, 0x7b, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x7d, 0x92,
	0x41, 0x4c, 0x0a, 0x05, 0x50, 0x65, 0x65, 0x72, 0x73, 0x12, 0x18, 0x47, 0x65, 0x74, 0x20, 0x61,
	0x20, 0x70, 0x65, 0x65, 0x72, 0x20, 0x62, 0x79, 0x20, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x20,
	0x6b, 0x65, 0x79, 0x1a, 0x29, 0x47, 0x65, 0x74, 0x20, 0x61, 0x20, 0x70, 0x65, 0x65, 0x72, 0x20,
	0x62, 0x79, 0x20, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x20, 0x6b, 0x65, 0x79, 0x20, 0x66, 0x72,
	0x6f, 0x6d, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x12, 0x82,
	0x01, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x50, 0x65, 0x65, 0x72, 0x73, 0x12, 0x17, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x65, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x50, 0x65, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x43,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x65, 0x65,
	0x72, 0x73, 0x92, 0x41, 0x2e, 0x0a, 0x05, 0x50, 0x65, 0x65, 0x72, 0x73, 0x12, 0x09, 0x47, 0x65,
	0x74, 0x20, 0x70, 0x65, 0x65, 0x72, 0x73, 0x1a, 0x1a, 0x47, 0x65, 0x74, 0x20, 0x70, 0x65, 0x65,
	0x72, 0x73, 0x20, 0x66, 0x72, 0x6f, 0x6d, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x42, 0x74, 0x5a, 0x09, 0x2e, 0x2f, 0x3b, 0x70, 0x65, 0x65, 0x72, 0x70, 0x62,
	0x92, 0x41, 0x66, 0x12, 0x05, 0x32, 0x03, 0x31, 0x2e, 0x30, 0x2a, 0x01, 0x01, 0x72, 0x5a, 0x0a,
	0x2f, 0x67, 0x52, 0x50, 0x43, 0x2b, 0x48, 0x54, 0x54, 0x50, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x20, 0x77,
	0x69, 0x72, 0x65, 0x67, 0x75, 0x61, 0x72, 0x64, 0x20, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73,
	0x12, 0x27, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x41, 0x5a, 0x68, 0x75, 0x72, 0x37, 0x37, 0x31, 0x2f, 0x77, 0x67,
	0x2d, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_peer_service_proto_rawDescOnce sync.Once
	file_peer_service_proto_rawDescData = file_peer_service_proto_rawDesc
)

func file_peer_service_proto_rawDescGZIP() []byte {
	file_peer_service_proto_rawDescOnce.Do(func() {
		file_peer_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_peer_service_proto_rawDescData)
	})
	return file_peer_service_proto_rawDescData
}

var file_peer_service_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_peer_service_proto_goTypes = []interface{}{
	(*Peer)(nil),                // 0: api.v1.Peer
	(*AddPeerRequest)(nil),      // 1: api.v1.AddPeerRequest
	(*AddPeerResponse)(nil),     // 2: api.v1.AddPeerResponse
	(*RemovePeerRequest)(nil),   // 3: api.v1.RemovePeerRequest
	(*RemovePeerResponse)(nil),  // 4: api.v1.RemovePeerResponse
	(*UpdatePeerRequest)(nil),   // 5: api.v1.UpdatePeerRequest
	(*UpdatePeerResponse)(nil),  // 6: api.v1.UpdatePeerResponse
	(*GetPeerRequest)(nil),      // 7: api.v1.GetPeerRequest
	(*GetPeersRequest)(nil),     // 8: api.v1.GetPeersRequest
	(*GetPeersResponse)(nil),    // 9: api.v1.GetPeersResponse
	(*timestamp.Timestamp)(nil), // 10: google.protobuf.Timestamp
}
var file_peer_service_proto_depIdxs = []int32{
	10, // 0: api.v1.Peer.last_handshake:type_name -> google.protobuf.Timestamp
	0,  // 1: api.v1.GetPeersResponse.peers:type_name -> api.v1.Peer
	1,  // 2: api.v1.PeerService.AddPeer:input_type -> api.v1.AddPeerRequest
	3,  // 3: api.v1.PeerService.RemovePeer:input_type -> api.v1.RemovePeerRequest
	5,  // 4: api.v1.PeerService.UpdatePeer:input_type -> api.v1.UpdatePeerRequest
	7,  // 5: api.v1.PeerService.GetPeer:input_type -> api.v1.GetPeerRequest
	8,  // 6: api.v1.PeerService.GetPeers:input_type -> api.v1.GetPeersRequest
	2,  // 7: api.v1.PeerService.AddPeer:output_type -> api.v1.AddPeerResponse
	4,  // 8: api.v1.PeerService.RemovePeer:output_type -> api.v1.RemovePeerResponse
	6,  // 9: api.v1.PeerService.UpdatePeer:output_type -> api.v1.UpdatePeerResponse
	0,  // 10: api.v1.PeerService.GetPeer:output_type -> api.v1.Peer
	9,  // 11: api.v1.PeerService.GetPeers:output_type -> api.v1.GetPeersResponse
	7,  // [7:12] is the sub-list for method output_type
	2,  // [2:7] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_peer_service_proto_init() }
func file_peer_service_proto_init() {
	if File_peer_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_peer_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Peer); i {
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
		file_peer_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddPeerRequest); i {
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
		file_peer_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddPeerResponse); i {
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
		file_peer_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemovePeerRequest); i {
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
		file_peer_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemovePeerResponse); i {
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
		file_peer_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePeerRequest); i {
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
		file_peer_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePeerResponse); i {
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
		file_peer_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPeerRequest); i {
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
		file_peer_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPeersRequest); i {
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
		file_peer_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPeersResponse); i {
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
			RawDescriptor: file_peer_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_peer_service_proto_goTypes,
		DependencyIndexes: file_peer_service_proto_depIdxs,
		MessageInfos:      file_peer_service_proto_msgTypes,
	}.Build()
	File_peer_service_proto = out.File
	file_peer_service_proto_rawDesc = nil
	file_peer_service_proto_goTypes = nil
	file_peer_service_proto_depIdxs = nil
}
