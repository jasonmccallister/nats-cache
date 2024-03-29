// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        (unknown)
// source: cache/v1/cache.proto

package cachev1

import (
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

type ExistsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Database *uint32  `protobuf:"varint,1,opt,name=database,proto3,oneof" json:"database,omitempty"`
	Keys     []string `protobuf:"bytes,2,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *ExistsRequest) Reset() {
	*x = ExistsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExistsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExistsRequest) ProtoMessage() {}

func (x *ExistsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExistsRequest.ProtoReflect.Descriptor instead.
func (*ExistsRequest) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{0}
}

func (x *ExistsRequest) GetDatabase() uint32 {
	if x != nil && x.Database != nil {
		return *x.Database
	}
	return 0
}

func (x *ExistsRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type ExistsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys  []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Count uint32   `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *ExistsResponse) Reset() {
	*x = ExistsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExistsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExistsResponse) ProtoMessage() {}

func (x *ExistsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExistsResponse.ProtoReflect.Descriptor instead.
func (*ExistsResponse) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{1}
}

func (x *ExistsResponse) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *ExistsResponse) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

// GetRequest is the request message for the Get method. The database is optional and
// used primarily as a prefix for the key. This supports applications that use the
// service as a mutex, sessions, or cache without making additional accounts.
// If the database is not specified, the default database of 0 is assumed.
type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Database *uint32 `protobuf:"varint,1,opt,name=database,proto3,oneof" json:"database,omitempty"`
	Key      string  `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{2}
}

func (x *GetRequest) GetDatabase() uint32 {
	if x != nil && x.Database != nil {
		return *x.Database
	}
	return 0
}

func (x *GetRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

// GetResponse is the response message for the Get method.
type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	// ttl is the time to live for the key. If the key is not found, the ttl will be 0. Otherwise,
	// it will be expiration in unix time.
	Ttl int64 `protobuf:"varint,3,opt,name=ttl,proto3" json:"ttl,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{3}
}

func (x *GetResponse) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *GetResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *GetResponse) GetTtl() int64 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

// SetGetRequest is the request message for the Set method. The database is optional and
// used primarily as a prefix for the key. This supports applications that use the
// service as a mutex, sessions, or cache without making additional accounts.
// If the database is not specified, the default database of 0 is used.
type SetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Database *uint32 `protobuf:"varint,1,opt,name=database,proto3,oneof" json:"database,omitempty"`
	Key      string  `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value    []byte  `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	Ttl      *uint32 `protobuf:"varint,4,opt,name=ttl,proto3,oneof" json:"ttl,omitempty"`
}

func (x *SetRequest) Reset() {
	*x = SetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetRequest) ProtoMessage() {}

func (x *SetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetRequest.ProtoReflect.Descriptor instead.
func (*SetRequest) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{4}
}

func (x *SetRequest) GetDatabase() uint32 {
	if x != nil && x.Database != nil {
		return *x.Database
	}
	return 0
}

func (x *SetRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetRequest) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *SetRequest) GetTtl() uint32 {
	if x != nil && x.Ttl != nil {
		return *x.Ttl
	}
	return 0
}

type SetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Ttl   int64  `protobuf:"varint,3,opt,name=ttl,proto3" json:"ttl,omitempty"`
}

func (x *SetResponse) Reset() {
	*x = SetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetResponse) ProtoMessage() {}

func (x *SetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetResponse.ProtoReflect.Descriptor instead.
func (*SetResponse) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{5}
}

func (x *SetResponse) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *SetResponse) GetTtl() int64 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

// DeleteRequest is the request message for the Delete method. The database is optional and
// used primarily as a prefix for the key. If the database is not specified, the default
// database of 0 is used.
type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Database *uint32 `protobuf:"varint,1,opt,name=database,proto3,oneof" json:"database,omitempty"`
	Key      string  `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteRequest) GetDatabase() uint32 {
	if x != nil && x.Database != nil {
		return *x.Database
	}
	return 0
}

func (x *DeleteRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type DeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Deleted bool `protobuf:"varint,1,opt,name=deleted,proto3" json:"deleted,omitempty"`
}

func (x *DeleteResponse) Reset() {
	*x = DeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteResponse) ProtoMessage() {}

func (x *DeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteResponse.ProtoReflect.Descriptor instead.
func (*DeleteResponse) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteResponse) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

// PurgeRequest is the request message for the Purge method.
type PurgeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Database *uint32 `protobuf:"varint,1,opt,name=database,proto3,oneof" json:"database,omitempty"`
	Prefix   *string `protobuf:"bytes,2,opt,name=prefix,proto3,oneof" json:"prefix,omitempty"`
}

func (x *PurgeRequest) Reset() {
	*x = PurgeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PurgeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PurgeRequest) ProtoMessage() {}

func (x *PurgeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PurgeRequest.ProtoReflect.Descriptor instead.
func (*PurgeRequest) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{8}
}

func (x *PurgeRequest) GetDatabase() uint32 {
	if x != nil && x.Database != nil {
		return *x.Database
	}
	return 0
}

func (x *PurgeRequest) GetPrefix() string {
	if x != nil && x.Prefix != nil {
		return *x.Prefix
	}
	return ""
}

type PurgeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Purged bool `protobuf:"varint,1,opt,name=purged,proto3" json:"purged,omitempty"`
}

func (x *PurgeResponse) Reset() {
	*x = PurgeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PurgeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PurgeResponse) ProtoMessage() {}

func (x *PurgeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PurgeResponse.ProtoReflect.Descriptor instead.
func (*PurgeResponse) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{9}
}

func (x *PurgeResponse) GetPurged() bool {
	if x != nil {
		return x.Purged
	}
	return false
}

type GetMultiRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Database *uint32  `protobuf:"varint,1,opt,name=database,proto3,oneof" json:"database,omitempty"`
	Keys     []string `protobuf:"bytes,2,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *GetMultiRequest) Reset() {
	*x = GetMultiRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMultiRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMultiRequest) ProtoMessage() {}

func (x *GetMultiRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMultiRequest.ProtoReflect.Descriptor instead.
func (*GetMultiRequest) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{10}
}

func (x *GetMultiRequest) GetDatabase() uint32 {
	if x != nil && x.Database != nil {
		return *x.Database
	}
	return 0
}

func (x *GetMultiRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type Item struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Ttl   int64  `protobuf:"varint,3,opt,name=ttl,proto3" json:"ttl,omitempty"`
}

func (x *Item) Reset() {
	*x = Item{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Item) ProtoMessage() {}

func (x *Item) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Item.ProtoReflect.Descriptor instead.
func (*Item) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{11}
}

func (x *Item) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Item) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Item) GetTtl() int64 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

type GetMultiResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*Item `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *GetMultiResponse) Reset() {
	*x = GetMultiResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_v1_cache_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMultiResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMultiResponse) ProtoMessage() {}

func (x *GetMultiResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cache_v1_cache_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMultiResponse.ProtoReflect.Descriptor instead.
func (*GetMultiResponse) Descriptor() ([]byte, []int) {
	return file_cache_v1_cache_proto_rawDescGZIP(), []int{12}
}

func (x *GetMultiResponse) GetItems() []*Item {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_cache_v1_cache_proto protoreflect.FileDescriptor

var file_cache_v1_cache_proto_rawDesc = []byte{
	0x0a, 0x14, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x63, 0x68, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31,
	0x22, 0x51, 0x0a, 0x0d, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1f, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x88,
	0x01, 0x01, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x62,
	0x61, 0x73, 0x65, 0x22, 0x3a, 0x0a, 0x0e, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x4c, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a,
	0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x48,
	0x00, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x88, 0x01, 0x01, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x22, 0x47, 0x0a,
	0x0b, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x03, 0x74, 0x74, 0x6c, 0x22, 0x81, 0x01, 0x0a, 0x0a, 0x53, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62,
	0x61, 0x73, 0x65, 0x88, 0x01, 0x01, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x15,
	0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x01, 0x52, 0x03, 0x74,
	0x74, 0x6c, 0x88, 0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61,
	0x73, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x74, 0x74, 0x6c, 0x22, 0x47, 0x0a, 0x0b, 0x53, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03,
	0x74, 0x74, 0x6c, 0x22, 0x4f, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61,
	0x73, 0x65, 0x88, 0x01, 0x01, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0x22, 0x2a, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64,
	0x22, 0x64, 0x0a, 0x0c, 0x50, 0x75, 0x72, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1f, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x48, 0x00, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x1b, 0x0a, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x01, 0x52, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x88, 0x01, 0x01, 0x42, 0x0b,
	0x0a, 0x09, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x5f,
	0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x22, 0x27, 0x0a, 0x0d, 0x50, 0x75, 0x72, 0x67, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x75, 0x72, 0x67, 0x65,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x70, 0x75, 0x72, 0x67, 0x65, 0x64, 0x22,
	0x53, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1f, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0x22, 0x40, 0x0a, 0x04, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x03, 0x74, 0x74, 0x6c, 0x22, 0x38, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4d, 0x75, 0x6c,
	0x74, 0x69, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x61, 0x63, 0x68,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x32, 0xf9, 0x03, 0x0a, 0x0c, 0x43, 0x61, 0x63, 0x68, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x3d, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x17, 0x2e, 0x63, 0x61,
	0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x3d, 0x0a, 0x06, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x12, 0x17, 0x2e, 0x63, 0x61, 0x63,
	0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x78, 0x69, 0x73, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x34, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x14, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x63,
	0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x43, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x4d, 0x75, 0x6c, 0x74,
	0x69, 0x12, 0x19, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x4d, 0x75, 0x6c, 0x74, 0x69, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x63,
	0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x75, 0x6c, 0x74, 0x69,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x09, 0x47, 0x65,
	0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x14, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e,
	0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x3a, 0x0a, 0x05, 0x50, 0x75,
	0x72, 0x67, 0x65, 0x12, 0x16, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x75, 0x72, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x63, 0x61,
	0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x72, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x09, 0x53, 0x65, 0x74, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x12, 0x14, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x63, 0x61, 0x63, 0x68,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x34, 0x0a, 0x03, 0x53, 0x65, 0x74, 0x12, 0x14, 0x2e,
	0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0xa1, 0x01, 0x0a,
	0x0c, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x43,
	0x61, 0x63, 0x68, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x44, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x61, 0x73, 0x6f, 0x6e, 0x6d, 0x63, 0x63,
	0x61, 0x6c, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x6e, 0x61, 0x74, 0x73, 0x2d, 0x63, 0x61,
	0x63, 0x68, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x65, 0x6e,
	0x2f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x61, 0x63, 0x68, 0x65, 0x76,
	0x31, 0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa, 0x02, 0x08, 0x43, 0x61, 0x63, 0x68, 0x65, 0x2e,
	0x56, 0x31, 0xca, 0x02, 0x08, 0x43, 0x61, 0x63, 0x68, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x14,
	0x43, 0x61, 0x63, 0x68, 0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x09, 0x43, 0x61, 0x63, 0x68, 0x65, 0x3a, 0x3a, 0x56, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cache_v1_cache_proto_rawDescOnce sync.Once
	file_cache_v1_cache_proto_rawDescData = file_cache_v1_cache_proto_rawDesc
)

func file_cache_v1_cache_proto_rawDescGZIP() []byte {
	file_cache_v1_cache_proto_rawDescOnce.Do(func() {
		file_cache_v1_cache_proto_rawDescData = protoimpl.X.CompressGZIP(file_cache_v1_cache_proto_rawDescData)
	})
	return file_cache_v1_cache_proto_rawDescData
}

var file_cache_v1_cache_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_cache_v1_cache_proto_goTypes = []interface{}{
	(*ExistsRequest)(nil),    // 0: cache.v1.ExistsRequest
	(*ExistsResponse)(nil),   // 1: cache.v1.ExistsResponse
	(*GetRequest)(nil),       // 2: cache.v1.GetRequest
	(*GetResponse)(nil),      // 3: cache.v1.GetResponse
	(*SetRequest)(nil),       // 4: cache.v1.SetRequest
	(*SetResponse)(nil),      // 5: cache.v1.SetResponse
	(*DeleteRequest)(nil),    // 6: cache.v1.DeleteRequest
	(*DeleteResponse)(nil),   // 7: cache.v1.DeleteResponse
	(*PurgeRequest)(nil),     // 8: cache.v1.PurgeRequest
	(*PurgeResponse)(nil),    // 9: cache.v1.PurgeResponse
	(*GetMultiRequest)(nil),  // 10: cache.v1.GetMultiRequest
	(*Item)(nil),             // 11: cache.v1.Item
	(*GetMultiResponse)(nil), // 12: cache.v1.GetMultiResponse
}
var file_cache_v1_cache_proto_depIdxs = []int32{
	11, // 0: cache.v1.GetMultiResponse.items:type_name -> cache.v1.Item
	6,  // 1: cache.v1.CacheService.Delete:input_type -> cache.v1.DeleteRequest
	0,  // 2: cache.v1.CacheService.Exists:input_type -> cache.v1.ExistsRequest
	2,  // 3: cache.v1.CacheService.Get:input_type -> cache.v1.GetRequest
	10, // 4: cache.v1.CacheService.GetMulti:input_type -> cache.v1.GetMultiRequest
	2,  // 5: cache.v1.CacheService.GetStream:input_type -> cache.v1.GetRequest
	8,  // 6: cache.v1.CacheService.Purge:input_type -> cache.v1.PurgeRequest
	4,  // 7: cache.v1.CacheService.SetStream:input_type -> cache.v1.SetRequest
	4,  // 8: cache.v1.CacheService.Set:input_type -> cache.v1.SetRequest
	7,  // 9: cache.v1.CacheService.Delete:output_type -> cache.v1.DeleteResponse
	1,  // 10: cache.v1.CacheService.Exists:output_type -> cache.v1.ExistsResponse
	3,  // 11: cache.v1.CacheService.Get:output_type -> cache.v1.GetResponse
	12, // 12: cache.v1.CacheService.GetMulti:output_type -> cache.v1.GetMultiResponse
	3,  // 13: cache.v1.CacheService.GetStream:output_type -> cache.v1.GetResponse
	9,  // 14: cache.v1.CacheService.Purge:output_type -> cache.v1.PurgeResponse
	5,  // 15: cache.v1.CacheService.SetStream:output_type -> cache.v1.SetResponse
	5,  // 16: cache.v1.CacheService.Set:output_type -> cache.v1.SetResponse
	9,  // [9:17] is the sub-list for method output_type
	1,  // [1:9] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_cache_v1_cache_proto_init() }
func file_cache_v1_cache_proto_init() {
	if File_cache_v1_cache_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cache_v1_cache_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExistsRequest); i {
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
		file_cache_v1_cache_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExistsResponse); i {
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
		file_cache_v1_cache_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
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
		file_cache_v1_cache_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResponse); i {
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
		file_cache_v1_cache_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetRequest); i {
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
		file_cache_v1_cache_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetResponse); i {
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
		file_cache_v1_cache_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
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
		file_cache_v1_cache_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteResponse); i {
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
		file_cache_v1_cache_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PurgeRequest); i {
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
		file_cache_v1_cache_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PurgeResponse); i {
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
		file_cache_v1_cache_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMultiRequest); i {
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
		file_cache_v1_cache_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Item); i {
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
		file_cache_v1_cache_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMultiResponse); i {
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
	file_cache_v1_cache_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_cache_v1_cache_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_cache_v1_cache_proto_msgTypes[4].OneofWrappers = []interface{}{}
	file_cache_v1_cache_proto_msgTypes[6].OneofWrappers = []interface{}{}
	file_cache_v1_cache_proto_msgTypes[8].OneofWrappers = []interface{}{}
	file_cache_v1_cache_proto_msgTypes[10].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cache_v1_cache_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cache_v1_cache_proto_goTypes,
		DependencyIndexes: file_cache_v1_cache_proto_depIdxs,
		MessageInfos:      file_cache_v1_cache_proto_msgTypes,
	}.Build()
	File_cache_v1_cache_proto = out.File
	file_cache_v1_cache_proto_rawDesc = nil
	file_cache_v1_cache_proto_goTypes = nil
	file_cache_v1_cache_proto_depIdxs = nil
}
