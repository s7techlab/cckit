// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: crosscc/cclocator_setting.proto

package crosscc

import (
	context "context"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Request: set service resolving setting
type ServiceLocatorSetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Service   string `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`     // service identifier
	Channel   string `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`     // channel id
	Chaincode string `protobuf:"bytes,3,opt,name=chaincode,proto3" json:"chaincode,omitempty"` // chaincode name
}

func (x *ServiceLocatorSetRequest) Reset() {
	*x = ServiceLocatorSetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crosscc_cclocator_setting_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceLocatorSetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceLocatorSetRequest) ProtoMessage() {}

func (x *ServiceLocatorSetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_crosscc_cclocator_setting_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceLocatorSetRequest.ProtoReflect.Descriptor instead.
func (*ServiceLocatorSetRequest) Descriptor() ([]byte, []int) {
	return file_crosscc_cclocator_setting_proto_rawDescGZIP(), []int{0}
}

func (x *ServiceLocatorSetRequest) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *ServiceLocatorSetRequest) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

func (x *ServiceLocatorSetRequest) GetChaincode() string {
	if x != nil {
		return x.Chaincode
	}
	return ""
}

// State: ervice resolving setting
type ServiceLocator struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Service   string `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`     // service identifier
	Channel   string `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`     // channel id
	Chaincode string `protobuf:"bytes,3,opt,name=chaincode,proto3" json:"chaincode,omitempty"` // chaincode name
}

func (x *ServiceLocator) Reset() {
	*x = ServiceLocator{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crosscc_cclocator_setting_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceLocator) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceLocator) ProtoMessage() {}

func (x *ServiceLocator) ProtoReflect() protoreflect.Message {
	mi := &file_crosscc_cclocator_setting_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceLocator.ProtoReflect.Descriptor instead.
func (*ServiceLocator) Descriptor() ([]byte, []int) {
	return file_crosscc_cclocator_setting_proto_rawDescGZIP(), []int{1}
}

func (x *ServiceLocator) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *ServiceLocator) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

func (x *ServiceLocator) GetChaincode() string {
	if x != nil {
		return x.Chaincode
	}
	return ""
}

// Id: service resolving setting identifier
type ServiceLocatorId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Service string `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"` // service identifier
}

func (x *ServiceLocatorId) Reset() {
	*x = ServiceLocatorId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crosscc_cclocator_setting_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceLocatorId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceLocatorId) ProtoMessage() {}

func (x *ServiceLocatorId) ProtoReflect() protoreflect.Message {
	mi := &file_crosscc_cclocator_setting_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceLocatorId.ProtoReflect.Descriptor instead.
func (*ServiceLocatorId) Descriptor() ([]byte, []int) {
	return file_crosscc_cclocator_setting_proto_rawDescGZIP(), []int{2}
}

func (x *ServiceLocatorId) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

// List: service resolving settings
type ServiceLocators struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*ServiceLocator `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *ServiceLocators) Reset() {
	*x = ServiceLocators{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crosscc_cclocator_setting_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceLocators) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceLocators) ProtoMessage() {}

func (x *ServiceLocators) ProtoReflect() protoreflect.Message {
	mi := &file_crosscc_cclocator_setting_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceLocators.ProtoReflect.Descriptor instead.
func (*ServiceLocators) Descriptor() ([]byte, []int) {
	return file_crosscc_cclocator_setting_proto_rawDescGZIP(), []int{3}
}

func (x *ServiceLocators) GetItems() []*ServiceLocator {
	if x != nil {
		return x.Items
	}
	return nil
}

// Event: service resolving settings was set
type ServiceLocatorSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Service   string `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`     // service identifier
	Channel   string `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`     // channel id
	Chaincode string `protobuf:"bytes,3,opt,name=chaincode,proto3" json:"chaincode,omitempty"` // chaincode name
}

func (x *ServiceLocatorSet) Reset() {
	*x = ServiceLocatorSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crosscc_cclocator_setting_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceLocatorSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceLocatorSet) ProtoMessage() {}

func (x *ServiceLocatorSet) ProtoReflect() protoreflect.Message {
	mi := &file_crosscc_cclocator_setting_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceLocatorSet.ProtoReflect.Descriptor instead.
func (*ServiceLocatorSet) Descriptor() ([]byte, []int) {
	return file_crosscc_cclocator_setting_proto_rawDescGZIP(), []int{4}
}

func (x *ServiceLocatorSet) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *ServiceLocatorSet) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

func (x *ServiceLocatorSet) GetChaincode() string {
	if x != nil {
		return x.Chaincode
	}
	return ""
}

type PingServiceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Locator *ServiceLocator `protobuf:"bytes,1,opt,name=locator,proto3" json:"locator,omitempty"`
	Error   string          `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *PingServiceResponse) Reset() {
	*x = PingServiceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crosscc_cclocator_setting_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingServiceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingServiceResponse) ProtoMessage() {}

func (x *PingServiceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_crosscc_cclocator_setting_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingServiceResponse.ProtoReflect.Descriptor instead.
func (*PingServiceResponse) Descriptor() ([]byte, []int) {
	return file_crosscc_cclocator_setting_proto_rawDescGZIP(), []int{5}
}

func (x *PingServiceResponse) GetLocator() *ServiceLocator {
	if x != nil {
		return x.Locator
	}
	return nil
}

func (x *PingServiceResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type PingServiceResponses struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Responses []*PingServiceResponse `protobuf:"bytes,1,rep,name=responses,proto3" json:"responses,omitempty"`
}

func (x *PingServiceResponses) Reset() {
	*x = PingServiceResponses{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crosscc_cclocator_setting_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingServiceResponses) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingServiceResponses) ProtoMessage() {}

func (x *PingServiceResponses) ProtoReflect() protoreflect.Message {
	mi := &file_crosscc_cclocator_setting_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingServiceResponses.ProtoReflect.Descriptor instead.
func (*PingServiceResponses) Descriptor() ([]byte, []int) {
	return file_crosscc_cclocator_setting_proto_rawDescGZIP(), []int{6}
}

func (x *PingServiceResponses) GetResponses() []*PingServiceResponse {
	if x != nil {
		return x.Responses
	}
	return nil
}

var File_crosscc_cclocator_setting_proto protoreflect.FileDescriptor

var file_crosscc_cclocator_setting_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x63, 0x63, 0x2f, 0x63, 0x63, 0x6c, 0x6f, 0x63, 0x61,
	0x74, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x07, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x63, 0x63, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6c, 0x0a, 0x18, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x63,
	0x6f, 0x64, 0x65, 0x22, 0x62, 0x0a, 0x0e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x2c, 0x0a, 0x10, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x40, 0x0a, 0x0f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x2d, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x63,
	0x63, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72,
	0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x65, 0x0a, 0x11, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x12, 0x1c, 0x0a, 0x09, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x5e,
	0x0a, 0x13, 0x50, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x07, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x63, 0x63,
	0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52,
	0x07, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x52,
	0x0a, 0x14, 0x50, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x12, 0x3a, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x72, 0x6f, 0x73,
	0x73, 0x63, 0x63, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x09, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x73, 0x32, 0x8a, 0x04, 0x0a, 0x0e, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6c, 0x0a, 0x11, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x74, 0x12, 0x21, 0x2e, 0x63, 0x72, 0x6f,
	0x73, 0x73, 0x63, 0x63, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x6f, 0x72, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e,
	0x63, 0x72, 0x6f, 0x73, 0x73, 0x63, 0x63, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x22, 0x10,
	0x2f, 0x63, 0x72, 0x6f, 0x73, 0x63, 0x63, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x3a, 0x01, 0x2a, 0x12, 0x62, 0x0a, 0x11, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x6f, 0x72, 0x47, 0x65, 0x74, 0x12, 0x19, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73,
	0x63, 0x63, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x6f,
	0x72, 0x49, 0x64, 0x1a, 0x17, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x63, 0x63, 0x2e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x19, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x13, 0x12, 0x11, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x63, 0x63, 0x2f, 0x7b, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x7d, 0x12, 0x61, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x18, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x63, 0x63,
	0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x73,
	0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x12, 0x10, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x63,
	0x63, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x66, 0x0a, 0x0b, 0x50, 0x69,
	0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x19, 0x2e, 0x63, 0x72, 0x6f, 0x73,
	0x73, 0x63, 0x63, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x6f, 0x72, 0x49, 0x64, 0x1a, 0x1c, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x63, 0x63, 0x2e, 0x50,
	0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x12, 0x16, 0x2f, 0x63, 0x72, 0x6f,
	0x73, 0x63, 0x63, 0x2f, 0x70, 0x69, 0x6e, 0x67, 0x2f, 0x7b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x7d, 0x12, 0x5b, 0x0a, 0x0c, 0x50, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1d, 0x2e, 0x63, 0x72, 0x6f,
	0x73, 0x73, 0x63, 0x63, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x0e, 0x12, 0x0c, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x63, 0x63, 0x2f, 0x70, 0x69, 0x6e, 0x67, 0x42,
	0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x37,
	0x74, 0x65, 0x63, 0x68, 0x6c, 0x61, 0x62, 0x2f, 0x63, 0x63, 0x6b, 0x69, 0x74, 0x2f, 0x65, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x63, 0x63,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_crosscc_cclocator_setting_proto_rawDescOnce sync.Once
	file_crosscc_cclocator_setting_proto_rawDescData = file_crosscc_cclocator_setting_proto_rawDesc
)

func file_crosscc_cclocator_setting_proto_rawDescGZIP() []byte {
	file_crosscc_cclocator_setting_proto_rawDescOnce.Do(func() {
		file_crosscc_cclocator_setting_proto_rawDescData = protoimpl.X.CompressGZIP(file_crosscc_cclocator_setting_proto_rawDescData)
	})
	return file_crosscc_cclocator_setting_proto_rawDescData
}

var file_crosscc_cclocator_setting_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_crosscc_cclocator_setting_proto_goTypes = []interface{}{
	(*ServiceLocatorSetRequest)(nil), // 0: crosscc.ServiceLocatorSetRequest
	(*ServiceLocator)(nil),           // 1: crosscc.ServiceLocator
	(*ServiceLocatorId)(nil),         // 2: crosscc.ServiceLocatorId
	(*ServiceLocators)(nil),          // 3: crosscc.ServiceLocators
	(*ServiceLocatorSet)(nil),        // 4: crosscc.ServiceLocatorSet
	(*PingServiceResponse)(nil),      // 5: crosscc.PingServiceResponse
	(*PingServiceResponses)(nil),     // 6: crosscc.PingServiceResponses
	(*emptypb.Empty)(nil),            // 7: google.protobuf.Empty
}
var file_crosscc_cclocator_setting_proto_depIdxs = []int32{
	1, // 0: crosscc.ServiceLocators.items:type_name -> crosscc.ServiceLocator
	1, // 1: crosscc.PingServiceResponse.locator:type_name -> crosscc.ServiceLocator
	5, // 2: crosscc.PingServiceResponses.responses:type_name -> crosscc.PingServiceResponse
	0, // 3: crosscc.SettingService.ServiceLocatorSet:input_type -> crosscc.ServiceLocatorSetRequest
	2, // 4: crosscc.SettingService.ServiceLocatorGet:input_type -> crosscc.ServiceLocatorId
	7, // 5: crosscc.SettingService.ListServiceLocators:input_type -> google.protobuf.Empty
	2, // 6: crosscc.SettingService.PingService:input_type -> crosscc.ServiceLocatorId
	7, // 7: crosscc.SettingService.PingServices:input_type -> google.protobuf.Empty
	1, // 8: crosscc.SettingService.ServiceLocatorSet:output_type -> crosscc.ServiceLocator
	1, // 9: crosscc.SettingService.ServiceLocatorGet:output_type -> crosscc.ServiceLocator
	3, // 10: crosscc.SettingService.ListServiceLocators:output_type -> crosscc.ServiceLocators
	5, // 11: crosscc.SettingService.PingService:output_type -> crosscc.PingServiceResponse
	6, // 12: crosscc.SettingService.PingServices:output_type -> crosscc.PingServiceResponses
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_crosscc_cclocator_setting_proto_init() }
func file_crosscc_cclocator_setting_proto_init() {
	if File_crosscc_cclocator_setting_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_crosscc_cclocator_setting_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceLocatorSetRequest); i {
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
		file_crosscc_cclocator_setting_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceLocator); i {
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
		file_crosscc_cclocator_setting_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceLocatorId); i {
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
		file_crosscc_cclocator_setting_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceLocators); i {
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
		file_crosscc_cclocator_setting_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceLocatorSet); i {
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
		file_crosscc_cclocator_setting_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingServiceResponse); i {
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
		file_crosscc_cclocator_setting_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingServiceResponses); i {
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
			RawDescriptor: file_crosscc_cclocator_setting_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_crosscc_cclocator_setting_proto_goTypes,
		DependencyIndexes: file_crosscc_cclocator_setting_proto_depIdxs,
		MessageInfos:      file_crosscc_cclocator_setting_proto_msgTypes,
	}.Build()
	File_crosscc_cclocator_setting_proto = out.File
	file_crosscc_cclocator_setting_proto_rawDesc = nil
	file_crosscc_cclocator_setting_proto_goTypes = nil
	file_crosscc_cclocator_setting_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SettingServiceClient is the client API for SettingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SettingServiceClient interface {
	ServiceLocatorSet(ctx context.Context, in *ServiceLocatorSetRequest, opts ...grpc.CallOption) (*ServiceLocator, error)
	ServiceLocatorGet(ctx context.Context, in *ServiceLocatorId, opts ...grpc.CallOption) (*ServiceLocator, error)
	ListServiceLocators(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ServiceLocators, error)
	// Try to query chaincodes from service chaincode settings
	PingService(ctx context.Context, in *ServiceLocatorId, opts ...grpc.CallOption) (*PingServiceResponse, error)
	PingServices(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingServiceResponses, error)
}

type settingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSettingServiceClient(cc grpc.ClientConnInterface) SettingServiceClient {
	return &settingServiceClient{cc}
}

func (c *settingServiceClient) ServiceLocatorSet(ctx context.Context, in *ServiceLocatorSetRequest, opts ...grpc.CallOption) (*ServiceLocator, error) {
	out := new(ServiceLocator)
	err := c.cc.Invoke(ctx, "/crosscc.SettingService/ServiceLocatorSet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *settingServiceClient) ServiceLocatorGet(ctx context.Context, in *ServiceLocatorId, opts ...grpc.CallOption) (*ServiceLocator, error) {
	out := new(ServiceLocator)
	err := c.cc.Invoke(ctx, "/crosscc.SettingService/ServiceLocatorGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *settingServiceClient) ListServiceLocators(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ServiceLocators, error) {
	out := new(ServiceLocators)
	err := c.cc.Invoke(ctx, "/crosscc.SettingService/ListServiceLocators", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *settingServiceClient) PingService(ctx context.Context, in *ServiceLocatorId, opts ...grpc.CallOption) (*PingServiceResponse, error) {
	out := new(PingServiceResponse)
	err := c.cc.Invoke(ctx, "/crosscc.SettingService/PingService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *settingServiceClient) PingServices(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingServiceResponses, error) {
	out := new(PingServiceResponses)
	err := c.cc.Invoke(ctx, "/crosscc.SettingService/PingServices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SettingServiceServer is the server API for SettingService service.
type SettingServiceServer interface {
	ServiceLocatorSet(context.Context, *ServiceLocatorSetRequest) (*ServiceLocator, error)
	ServiceLocatorGet(context.Context, *ServiceLocatorId) (*ServiceLocator, error)
	ListServiceLocators(context.Context, *emptypb.Empty) (*ServiceLocators, error)
	// Try to query chaincodes from service chaincode settings
	PingService(context.Context, *ServiceLocatorId) (*PingServiceResponse, error)
	PingServices(context.Context, *emptypb.Empty) (*PingServiceResponses, error)
}

// UnimplementedSettingServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSettingServiceServer struct {
}

func (*UnimplementedSettingServiceServer) ServiceLocatorSet(context.Context, *ServiceLocatorSetRequest) (*ServiceLocator, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServiceLocatorSet not implemented")
}
func (*UnimplementedSettingServiceServer) ServiceLocatorGet(context.Context, *ServiceLocatorId) (*ServiceLocator, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServiceLocatorGet not implemented")
}
func (*UnimplementedSettingServiceServer) ListServiceLocators(context.Context, *emptypb.Empty) (*ServiceLocators, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListServiceLocators not implemented")
}
func (*UnimplementedSettingServiceServer) PingService(context.Context, *ServiceLocatorId) (*PingServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingService not implemented")
}
func (*UnimplementedSettingServiceServer) PingServices(context.Context, *emptypb.Empty) (*PingServiceResponses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingServices not implemented")
}

func RegisterSettingServiceServer(s *grpc.Server, srv SettingServiceServer) {
	s.RegisterService(&_SettingService_serviceDesc, srv)
}

func _SettingService_ServiceLocatorSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceLocatorSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingServiceServer).ServiceLocatorSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crosscc.SettingService/ServiceLocatorSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingServiceServer).ServiceLocatorSet(ctx, req.(*ServiceLocatorSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SettingService_ServiceLocatorGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceLocatorId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingServiceServer).ServiceLocatorGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crosscc.SettingService/ServiceLocatorGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingServiceServer).ServiceLocatorGet(ctx, req.(*ServiceLocatorId))
	}
	return interceptor(ctx, in, info, handler)
}

func _SettingService_ListServiceLocators_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingServiceServer).ListServiceLocators(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crosscc.SettingService/ListServiceLocators",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingServiceServer).ListServiceLocators(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SettingService_PingService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceLocatorId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingServiceServer).PingService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crosscc.SettingService/PingService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingServiceServer).PingService(ctx, req.(*ServiceLocatorId))
	}
	return interceptor(ctx, in, info, handler)
}

func _SettingService_PingServices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingServiceServer).PingServices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crosscc.SettingService/PingServices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingServiceServer).PingServices(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _SettingService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "crosscc.SettingService",
	HandlerType: (*SettingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ServiceLocatorSet",
			Handler:    _SettingService_ServiceLocatorSet_Handler,
		},
		{
			MethodName: "ServiceLocatorGet",
			Handler:    _SettingService_ServiceLocatorGet_Handler,
		},
		{
			MethodName: "ListServiceLocators",
			Handler:    _SettingService_ListServiceLocators_Handler,
		},
		{
			MethodName: "PingService",
			Handler:    _SettingService_PingService_Handler,
		},
		{
			MethodName: "PingServices",
			Handler:    _SettingService_PingServices_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "crosscc/cclocator_setting.proto",
}
