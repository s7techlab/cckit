// Code generated by protoc-gen-go. DO NOT EDIT.
// source: payment.proto

/*
Package schema is a generated protocol buffer package.

It is generated from these files:
	payment.proto

It has these top-level messages:
	Payment
*/
package schema

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Payment struct {
	Type   string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Id     string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	Amount int32  `protobuf:"varint,3,opt,name=amount" json:"amount,omitempty"`
}

func (m *Payment) Reset()                    { *m = Payment{} }
func (m *Payment) String() string            { return proto.CompactTextString(m) }
func (*Payment) ProtoMessage()               {}
func (*Payment) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Payment) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Payment) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Payment) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func init() {
	proto.RegisterType((*Payment)(nil), "schema.Payment")
}

func init() { proto.RegisterFile("payment.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 107 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x48, 0xac, 0xcc,
	0x4d, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2b, 0x4e, 0xce, 0x48, 0xcd,
	0x4d, 0x54, 0x72, 0xe5, 0x62, 0x0f, 0x80, 0x48, 0x08, 0x09, 0x71, 0xb1, 0x94, 0x54, 0x16, 0xa4,
	0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x42, 0x7c, 0x5c, 0x4c, 0x99, 0x29, 0x12,
	0x4c, 0x60, 0x11, 0xa6, 0xcc, 0x14, 0x21, 0x31, 0x2e, 0xb6, 0xc4, 0xdc, 0xfc, 0xd2, 0xbc, 0x12,
	0x09, 0x66, 0x05, 0x46, 0x0d, 0xd6, 0x20, 0x28, 0x2f, 0x89, 0x0d, 0x6c, 0xaa, 0x31, 0x20, 0x00,
	0x00, 0xff, 0xff, 0xa2, 0xf4, 0x61, 0x36, 0x66, 0x00, 0x00, 0x00,
}
