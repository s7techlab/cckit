// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chaincode.proto

package service

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	peer "github.com/hyperledger/fabric-protos-go/peer"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type InvocationType int32

const (
	InvocationType_QUERY  InvocationType = 0
	InvocationType_INVOKE InvocationType = 1
)

var InvocationType_name = map[int32]string{
	0: "QUERY",
	1: "INVOKE",
}

var InvocationType_value = map[string]int32{
	"QUERY":  0,
	"INVOKE": 1,
}

func (x InvocationType) String() string {
	return proto.EnumName(InvocationType_name, int32(x))
}

func (InvocationType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_97136ef4b384cc22, []int{0}
}

type ChaincodeInput struct {
	// Chaincode name
	Chaincode string `protobuf:"bytes,1,opt,name=chaincode,proto3" json:"chaincode,omitempty"`
	// Channel name
	Channel string `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
	// Input contains the arguments for invocation.
	Args [][]byte `protobuf:"bytes,3,rep,name=args,proto3" json:"args,omitempty"`
	// TransientMap contains data (e.g. cryptographic material) that might be used
	// to implement some form of application-level confidentiality. The contents
	// of this field are supposed to always be omitted from the transaction and
	// excluded from the ledger.
	Transient            map[string][]byte `protobuf:"bytes,4,rep,name=transient,proto3" json:"transient,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ChaincodeInput) Reset()         { *m = ChaincodeInput{} }
func (m *ChaincodeInput) String() string { return proto.CompactTextString(m) }
func (*ChaincodeInput) ProtoMessage()    {}
func (*ChaincodeInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_97136ef4b384cc22, []int{0}
}

func (m *ChaincodeInput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeInput.Unmarshal(m, b)
}
func (m *ChaincodeInput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeInput.Marshal(b, m, deterministic)
}
func (m *ChaincodeInput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeInput.Merge(m, src)
}
func (m *ChaincodeInput) XXX_Size() int {
	return xxx_messageInfo_ChaincodeInput.Size(m)
}
func (m *ChaincodeInput) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeInput.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeInput proto.InternalMessageInfo

func (m *ChaincodeInput) GetChaincode() string {
	if m != nil {
		return m.Chaincode
	}
	return ""
}

func (m *ChaincodeInput) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

func (m *ChaincodeInput) GetArgs() [][]byte {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *ChaincodeInput) GetTransient() map[string][]byte {
	if m != nil {
		return m.Transient
	}
	return nil
}

type ChaincodeLocator struct {
	// Chaincode name
	Chaincode string `protobuf:"bytes,1,opt,name=chaincode,proto3" json:"chaincode,omitempty"`
	// Channel name
	Channel              string   `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaincodeLocator) Reset()         { *m = ChaincodeLocator{} }
func (m *ChaincodeLocator) String() string { return proto.CompactTextString(m) }
func (*ChaincodeLocator) ProtoMessage()    {}
func (*ChaincodeLocator) Descriptor() ([]byte, []int) {
	return fileDescriptor_97136ef4b384cc22, []int{1}
}

func (m *ChaincodeLocator) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeLocator.Unmarshal(m, b)
}
func (m *ChaincodeLocator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeLocator.Marshal(b, m, deterministic)
}
func (m *ChaincodeLocator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeLocator.Merge(m, src)
}
func (m *ChaincodeLocator) XXX_Size() int {
	return xxx_messageInfo_ChaincodeLocator.Size(m)
}
func (m *ChaincodeLocator) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeLocator.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeLocator proto.InternalMessageInfo

func (m *ChaincodeLocator) GetChaincode() string {
	if m != nil {
		return m.Chaincode
	}
	return ""
}

func (m *ChaincodeLocator) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

type ChaincodeExec struct {
	Type                 InvocationType  `protobuf:"varint,1,opt,name=type,proto3,enum=service.InvocationType" json:"type,omitempty"`
	Input                *ChaincodeInput `protobuf:"bytes,2,opt,name=input,proto3" json:"input,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ChaincodeExec) Reset()         { *m = ChaincodeExec{} }
func (m *ChaincodeExec) String() string { return proto.CompactTextString(m) }
func (*ChaincodeExec) ProtoMessage()    {}
func (*ChaincodeExec) Descriptor() ([]byte, []int) {
	return fileDescriptor_97136ef4b384cc22, []int{2}
}

func (m *ChaincodeExec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeExec.Unmarshal(m, b)
}
func (m *ChaincodeExec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeExec.Marshal(b, m, deterministic)
}
func (m *ChaincodeExec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeExec.Merge(m, src)
}
func (m *ChaincodeExec) XXX_Size() int {
	return xxx_messageInfo_ChaincodeExec.Size(m)
}
func (m *ChaincodeExec) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeExec.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeExec proto.InternalMessageInfo

func (m *ChaincodeExec) GetType() InvocationType {
	if m != nil {
		return m.Type
	}
	return InvocationType_QUERY
}

func (m *ChaincodeExec) GetInput() *ChaincodeInput {
	if m != nil {
		return m.Input
	}
	return nil
}

func init() {
	proto.RegisterEnum("service.InvocationType", InvocationType_name, InvocationType_value)
	proto.RegisterType((*ChaincodeInput)(nil), "service.ChaincodeInput")
	proto.RegisterMapType((map[string][]byte)(nil), "service.ChaincodeInput.TransientEntry")
	proto.RegisterType((*ChaincodeLocator)(nil), "service.ChaincodeLocator")
	proto.RegisterType((*ChaincodeExec)(nil), "service.ChaincodeExec")
}

func init() { proto.RegisterFile("chaincode.proto", fileDescriptor_97136ef4b384cc22) }

var fileDescriptor_97136ef4b384cc22 = []byte{
	// 424 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0x5d, 0x8b, 0x13, 0x31,
	0x14, 0x75, 0xb6, 0x1f, 0xcb, 0xdc, 0x5d, 0x6b, 0xb9, 0xc8, 0x3a, 0x0e, 0x3e, 0x94, 0x3e, 0x68,
	0x51, 0x9c, 0x4a, 0x7d, 0x59, 0x56, 0x45, 0x44, 0xe7, 0xa1, 0x2a, 0xea, 0x0e, 0xab, 0xe0, 0xd3,
	0x92, 0x66, 0xaf, 0xed, 0xd0, 0x9a, 0x84, 0x24, 0x2d, 0xce, 0xbf, 0xf5, 0x8f, 0x08, 0x92, 0xcc,
	0x74, 0xea, 0xb0, 0x28, 0xd2, 0xb7, 0x9b, 0xdc, 0x73, 0xce, 0x3d, 0x27, 0xb9, 0x70, 0x8b, 0x2f,
	0x58, 0x2e, 0xb8, 0xbc, 0xa2, 0x44, 0x69, 0x69, 0x25, 0x1e, 0x1a, 0xd2, 0x9b, 0x9c, 0x53, 0xfc,
	0x6a, 0x9e, 0xdb, 0xc5, 0x7a, 0x96, 0x70, 0xf9, 0x7d, 0xbc, 0x28, 0x14, 0xe9, 0x15, 0x5d, 0xcd,
	0x49, 0x8f, 0xbf, 0xb1, 0x99, 0xce, 0xf9, 0xd8, 0xa3, 0xcd, 0x58, 0x11, 0x69, 0x57, 0x2b, 0x69,
	0xd8, 0xea, 0x52, 0x93, 0x51, 0x52, 0x98, 0x4a, 0x2b, 0x7e, 0xf9, 0xff, 0x12, 0xb5, 0x8d, 0x4b,
	0xda, 0x90, 0xb0, 0xa5, 0xc0, 0xf0, 0x67, 0x00, 0xbd, 0xd7, 0xdb, 0xce, 0x54, 0xa8, 0xb5, 0xc5,
	0x7b, 0x10, 0xd6, 0xd8, 0x28, 0x18, 0x04, 0xa3, 0x30, 0xdb, 0x5d, 0x60, 0x04, 0x87, 0x7c, 0xc1,
	0x84, 0xa0, 0x55, 0x74, 0xe0, 0x7b, 0xdb, 0x23, 0x22, 0xb4, 0x99, 0x9e, 0x9b, 0xa8, 0x35, 0x68,
	0x8d, 0x8e, 0x33, 0x5f, 0xe3, 0x1b, 0x08, 0xad, 0x66, 0xc2, 0xe4, 0x24, 0x6c, 0xd4, 0x1e, 0xb4,
	0x46, 0x47, 0x93, 0xfb, 0x49, 0x95, 0x3f, 0x69, 0xce, 0x4d, 0x2e, 0xb6, 0xc0, 0x54, 0x58, 0x5d,
	0x64, 0x3b, 0x62, 0xfc, 0x1c, 0x7a, 0xcd, 0x26, 0xf6, 0xa1, 0xb5, 0xa4, 0xa2, 0x72, 0xe7, 0x4a,
	0xbc, 0x0d, 0x9d, 0x0d, 0x5b, 0xad, 0xc9, 0xbb, 0x3a, 0xce, 0xca, 0xc3, 0xd9, 0xc1, 0x69, 0x30,
	0x7c, 0x0b, 0xfd, 0x7a, 0xd2, 0x7b, 0xc9, 0x99, 0x95, 0x7a, 0xdf, 0x8c, 0xc3, 0x25, 0xdc, 0xac,
	0xb5, 0xd2, 0x1f, 0xc4, 0xf1, 0x11, 0xb4, 0x6d, 0xa1, 0x4a, 0x8d, 0xde, 0xe4, 0x4e, 0x9d, 0x6d,
	0x2a, 0x36, 0x6e, 0x54, 0x2e, 0xc5, 0x45, 0xa1, 0x28, 0xf3, 0x20, 0x7c, 0x0c, 0x9d, 0xdc, 0x45,
	0xf5, 0xaa, 0x47, 0x7f, 0xa0, 0x9b, 0x2f, 0x91, 0x95, 0xa8, 0x87, 0x0f, 0xa0, 0xd7, 0x94, 0xc1,
	0x10, 0x3a, 0xe7, 0x9f, 0xd3, 0xec, 0x6b, 0xff, 0x06, 0x02, 0x74, 0xa7, 0x1f, 0xbe, 0x7c, 0x7c,
	0x97, 0xf6, 0x83, 0xc9, 0xaf, 0x00, 0xc2, 0x5a, 0x02, 0x4f, 0xa1, 0xed, 0xad, 0x9d, 0x5c, 0x97,
	0x77, 0xf7, 0x71, 0x54, 0x7e, 0xbd, 0x49, 0x3e, 0x55, 0x4b, 0x95, 0x55, 0x3b, 0x85, 0x67, 0xd0,
	0x39, 0x5f, 0x93, 0x2e, 0xf0, 0x6f, 0xce, 0xfe, 0xc1, 0x7d, 0x06, 0x5d, 0x67, 0x76, 0x49, 0xfb,
	0x90, 0x5f, 0x40, 0x37, 0x75, 0x4b, 0x69, 0xf0, 0xee, 0x75, 0x72, 0xf5, 0x67, 0xf1, 0xc9, 0x96,
	0xbe, 0x8b, 0xe3, 0x38, 0x4f, 0x82, 0x59, 0xd7, 0x37, 0x9e, 0xfe, 0x0e, 0x00, 0x00, 0xff, 0xff,
	0x27, 0xc3, 0x64, 0xce, 0x6b, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ChaincodeClient is the client API for Chaincode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChaincodeClient interface {
	// Exec: Query or Invoke
	Exec(ctx context.Context, in *ChaincodeExec, opts ...grpc.CallOption) (*peer.ProposalResponse, error)
	// Query chaincode on home peer. Do NOT send to orderer.
	Query(ctx context.Context, in *ChaincodeInput, opts ...grpc.CallOption) (*peer.ProposalResponse, error)
	// Invoke chaincode on peers, according to endorsement policy and the SEND to orderer
	Invoke(ctx context.Context, in *ChaincodeInput, opts ...grpc.CallOption) (*peer.ProposalResponse, error)
	// Chaincode events stream
	Events(ctx context.Context, in *ChaincodeLocator, opts ...grpc.CallOption) (Chaincode_EventsClient, error)
}

type chaincodeClient struct {
	cc *grpc.ClientConn
}

func NewChaincodeClient(cc *grpc.ClientConn) ChaincodeClient {
	return &chaincodeClient{cc}
}

func (c *chaincodeClient) Exec(ctx context.Context, in *ChaincodeExec, opts ...grpc.CallOption) (*peer.ProposalResponse, error) {
	out := new(peer.ProposalResponse)
	err := c.cc.Invoke(ctx, "/service.Chaincode/Exec", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chaincodeClient) Query(ctx context.Context, in *ChaincodeInput, opts ...grpc.CallOption) (*peer.ProposalResponse, error) {
	out := new(peer.ProposalResponse)
	err := c.cc.Invoke(ctx, "/service.Chaincode/Query", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chaincodeClient) Invoke(ctx context.Context, in *ChaincodeInput, opts ...grpc.CallOption) (*peer.ProposalResponse, error) {
	out := new(peer.ProposalResponse)
	err := c.cc.Invoke(ctx, "/service.Chaincode/Invoke", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chaincodeClient) Events(ctx context.Context, in *ChaincodeLocator, opts ...grpc.CallOption) (Chaincode_EventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Chaincode_serviceDesc.Streams[0], "/service.Chaincode/Events", opts...)
	if err != nil {
		return nil, err
	}
	x := &chaincodeEventsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Chaincode_EventsClient interface {
	Recv() (*peer.ChaincodeEvent, error)
	grpc.ClientStream
}

type chaincodeEventsClient struct {
	grpc.ClientStream
}

func (x *chaincodeEventsClient) Recv() (*peer.ChaincodeEvent, error) {
	m := new(peer.ChaincodeEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChaincodeServer is the server API for Chaincode service.
type ChaincodeServer interface {
	// Exec: Query or Invoke
	Exec(context.Context, *ChaincodeExec) (*peer.ProposalResponse, error)
	// Query chaincode on home peer. Do NOT send to orderer.
	Query(context.Context, *ChaincodeInput) (*peer.ProposalResponse, error)
	// Invoke chaincode on peers, according to endorsement policy and the SEND to orderer
	Invoke(context.Context, *ChaincodeInput) (*peer.ProposalResponse, error)
	// Chaincode events stream
	Events(*ChaincodeLocator, Chaincode_EventsServer) error
}

func RegisterChaincodeServer(s *grpc.Server, srv ChaincodeServer) {
	s.RegisterService(&_Chaincode_serviceDesc, srv)
}

func _Chaincode_Exec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChaincodeExec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChaincodeServer).Exec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Chaincode/Exec",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChaincodeServer).Exec(ctx, req.(*ChaincodeExec))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chaincode_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChaincodeInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChaincodeServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Chaincode/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChaincodeServer).Query(ctx, req.(*ChaincodeInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chaincode_Invoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChaincodeInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChaincodeServer).Invoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Chaincode/Invoke",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChaincodeServer).Invoke(ctx, req.(*ChaincodeInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chaincode_Events_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ChaincodeLocator)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChaincodeServer).Events(m, &chaincodeEventsServer{stream})
}

type Chaincode_EventsServer interface {
	Send(*peer.ChaincodeEvent) error
	grpc.ServerStream
}

type chaincodeEventsServer struct {
	grpc.ServerStream
}

func (x *chaincodeEventsServer) Send(m *peer.ChaincodeEvent) error {
	return x.ServerStream.SendMsg(m)
}

var _Chaincode_serviceDesc = grpc.ServiceDesc{
	ServiceName: "service.Chaincode",
	HandlerType: (*ChaincodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Exec",
			Handler:    _Chaincode_Exec_Handler,
		},
		{
			MethodName: "Query",
			Handler:    _Chaincode_Query_Handler,
		},
		{
			MethodName: "Invoke",
			Handler:    _Chaincode_Invoke_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Events",
			Handler:       _Chaincode_Events_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chaincode.proto",
}
