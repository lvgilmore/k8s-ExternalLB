// Code generated by protoc-gen-go. DO NOT EDIT.
// source: externallb/externalLB.proto

package grpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Data struct {
	ServiceName          string   `protobuf:"bytes,1,opt,name=ServiceName" json:"ServiceName,omitempty"`
	Namespace            string   `protobuf:"bytes,2,opt,name=Namespace" json:"Namespace,omitempty"`
	ResourceVersion      string   `protobuf:"bytes,3,opt,name=ResourceVersion" json:"ResourceVersion,omitempty"`
	Protocol             string   `protobuf:"bytes,4,opt,name=Protocol" json:"Protocol,omitempty"`
	ExternalIPs          []string `protobuf:"bytes,5,rep,name=ExternalIPs" json:"ExternalIPs,omitempty"`
	RouterID             int32    `protobuf:"varint,6,opt,name=RouterID" json:"RouterID,omitempty"`
	SyncTime             int64    `protobuf:"varint,7,opt,name=SyncTime" json:"SyncTime,omitempty"`
	IsCreated            bool     `protobuf:"varint,8,opt,name=IsCreated" json:"IsCreated,omitempty"`
	Nodes                []string `protobuf:"bytes,9,rep,name=nodes" json:"nodes,omitempty"`
	Ports                []*Port  `protobuf:"bytes,10,rep,name=Ports" json:"Ports,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Data) Reset()         { *m = Data{} }
func (m *Data) String() string { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()    {}
func (*Data) Descriptor() ([]byte, []int) {
	return fileDescriptor_externalLB_3d625253ec2af4f9, []int{0}
}
func (m *Data) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Data.Unmarshal(m, b)
}
func (m *Data) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Data.Marshal(b, m, deterministic)
}
func (dst *Data) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Data.Merge(dst, src)
}
func (m *Data) XXX_Size() int {
	return xxx_messageInfo_Data.Size(m)
}
func (m *Data) XXX_DiscardUnknown() {
	xxx_messageInfo_Data.DiscardUnknown(m)
}

var xxx_messageInfo_Data proto.InternalMessageInfo

func (m *Data) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *Data) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *Data) GetResourceVersion() string {
	if m != nil {
		return m.ResourceVersion
	}
	return ""
}

func (m *Data) GetProtocol() string {
	if m != nil {
		return m.Protocol
	}
	return ""
}

func (m *Data) GetExternalIPs() []string {
	if m != nil {
		return m.ExternalIPs
	}
	return nil
}

func (m *Data) GetRouterID() int32 {
	if m != nil {
		return m.RouterID
	}
	return 0
}

func (m *Data) GetSyncTime() int64 {
	if m != nil {
		return m.SyncTime
	}
	return 0
}

func (m *Data) GetIsCreated() bool {
	if m != nil {
		return m.IsCreated
	}
	return false
}

func (m *Data) GetNodes() []string {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func (m *Data) GetPorts() []*Port {
	if m != nil {
		return m.Ports
	}
	return nil
}

type Port struct {
	Name                 string   `protobuf:"bytes,1,opt,name=Name" json:"Name,omitempty"`
	Port                 int32    `protobuf:"varint,2,opt,name=Port" json:"Port,omitempty"`
	NodePort             int32    `protobuf:"varint,3,opt,name=NodePort" json:"NodePort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Port) Reset()         { *m = Port{} }
func (m *Port) String() string { return proto.CompactTextString(m) }
func (*Port) ProtoMessage()    {}
func (*Port) Descriptor() ([]byte, []int) {
	return fileDescriptor_externalLB_3d625253ec2af4f9, []int{1}
}
func (m *Port) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Port.Unmarshal(m, b)
}
func (m *Port) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Port.Marshal(b, m, deterministic)
}
func (dst *Port) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Port.Merge(dst, src)
}
func (m *Port) XXX_Size() int {
	return xxx_messageInfo_Port.Size(m)
}
func (m *Port) XXX_DiscardUnknown() {
	xxx_messageInfo_Port.DiscardUnknown(m)
}

var xxx_messageInfo_Port proto.InternalMessageInfo

func (m *Port) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Port) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Port) GetNodePort() int32 {
	if m != nil {
		return m.NodePort
	}
	return 0
}

type Result struct {
	Addr                 string   `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Result) Reset()         { *m = Result{} }
func (m *Result) String() string { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()    {}
func (*Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_externalLB_3d625253ec2af4f9, []int{2}
}
func (m *Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Result.Unmarshal(m, b)
}
func (m *Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Result.Marshal(b, m, deterministic)
}
func (dst *Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Result.Merge(dst, src)
}
func (m *Result) XXX_Size() int {
	return xxx_messageInfo_Result.Size(m)
}
func (m *Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Result proto.InternalMessageInfo

func (m *Result) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type Nodes struct {
	List                 []string `protobuf:"bytes,1,rep,name=list" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Nodes) Reset()         { *m = Nodes{} }
func (m *Nodes) String() string { return proto.CompactTextString(m) }
func (*Nodes) ProtoMessage()    {}
func (*Nodes) Descriptor() ([]byte, []int) {
	return fileDescriptor_externalLB_3d625253ec2af4f9, []int{3}
}
func (m *Nodes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Nodes.Unmarshal(m, b)
}
func (m *Nodes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Nodes.Marshal(b, m, deterministic)
}
func (dst *Nodes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Nodes.Merge(dst, src)
}
func (m *Nodes) XXX_Size() int {
	return xxx_messageInfo_Nodes.Size(m)
}
func (m *Nodes) XXX_DiscardUnknown() {
	xxx_messageInfo_Nodes.DiscardUnknown(m)
}

var xxx_messageInfo_Nodes proto.InternalMessageInfo

func (m *Nodes) GetList() []string {
	if m != nil {
		return m.List
	}
	return nil
}

type UpdateService struct {
	Addr                 string   `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	IsSuccess            bool     `protobuf:"varint,2,opt,name=IsSuccess" json:"IsSuccess,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateService) Reset()         { *m = UpdateService{} }
func (m *UpdateService) String() string { return proto.CompactTextString(m) }
func (*UpdateService) ProtoMessage()    {}
func (*UpdateService) Descriptor() ([]byte, []int) {
	return fileDescriptor_externalLB_3d625253ec2af4f9, []int{4}
}
func (m *UpdateService) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateService.Unmarshal(m, b)
}
func (m *UpdateService) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateService.Marshal(b, m, deterministic)
}
func (dst *UpdateService) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateService.Merge(dst, src)
}
func (m *UpdateService) XXX_Size() int {
	return xxx_messageInfo_UpdateService.Size(m)
}
func (m *UpdateService) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateService.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateService proto.InternalMessageInfo

func (m *UpdateService) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *UpdateService) GetIsSuccess() bool {
	if m != nil {
		return m.IsSuccess
	}
	return false
}

func init() {
	proto.RegisterType((*Data)(nil), "grpc.Data")
	proto.RegisterType((*Port)(nil), "grpc.Port")
	proto.RegisterType((*Result)(nil), "grpc.Result")
	proto.RegisterType((*Nodes)(nil), "grpc.Nodes")
	proto.RegisterType((*UpdateService)(nil), "grpc.UpdateService")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ExternalLB service

type ExternalLBClient interface {
	Create(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Result, error)
	Update(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Result, error)
	Delete(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Result, error)
	NodesChange(ctx context.Context, in *Nodes, opts ...grpc.CallOption) (*UpdateService, error)
	Sync(ctx context.Context, opts ...grpc.CallOption) (ExternalLB_SyncClient, error)
}

type externalLBClient struct {
	cc *grpc.ClientConn
}

func NewExternalLBClient(cc *grpc.ClientConn) ExternalLBClient {
	return &externalLBClient{cc}
}

func (c *externalLBClient) Create(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/grpc.ExternalLB/Create", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalLBClient) Update(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/grpc.ExternalLB/Update", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalLBClient) Delete(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/grpc.ExternalLB/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalLBClient) NodesChange(ctx context.Context, in *Nodes, opts ...grpc.CallOption) (*UpdateService, error) {
	out := new(UpdateService)
	err := grpc.Invoke(ctx, "/grpc.ExternalLB/NodesChange", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalLBClient) Sync(ctx context.Context, opts ...grpc.CallOption) (ExternalLB_SyncClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_ExternalLB_serviceDesc.Streams[0], c.cc, "/grpc.ExternalLB/Sync", opts...)
	if err != nil {
		return nil, err
	}
	x := &externalLBSyncClient{stream}
	return x, nil
}

type ExternalLB_SyncClient interface {
	Send(*Data) error
	Recv() (*UpdateService, error)
	grpc.ClientStream
}

type externalLBSyncClient struct {
	grpc.ClientStream
}

func (x *externalLBSyncClient) Send(m *Data) error {
	return x.ClientStream.SendMsg(m)
}

func (x *externalLBSyncClient) Recv() (*UpdateService, error) {
	m := new(UpdateService)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for ExternalLB service

type ExternalLBServer interface {
	Create(context.Context, *Data) (*Result, error)
	Update(context.Context, *Data) (*Result, error)
	Delete(context.Context, *Data) (*Result, error)
	NodesChange(context.Context, *Nodes) (*UpdateService, error)
	Sync(ExternalLB_SyncServer) error
}

func RegisterExternalLBServer(s *grpc.Server, srv ExternalLBServer) {
	s.RegisterService(&_ExternalLB_serviceDesc, srv)
}

func _ExternalLB_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Data)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalLBServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.ExternalLB/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalLBServer).Create(ctx, req.(*Data))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalLB_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Data)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalLBServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.ExternalLB/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalLBServer).Update(ctx, req.(*Data))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalLB_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Data)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalLBServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.ExternalLB/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalLBServer).Delete(ctx, req.(*Data))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalLB_NodesChange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nodes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalLBServer).NodesChange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.ExternalLB/NodesChange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalLBServer).NodesChange(ctx, req.(*Nodes))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalLB_Sync_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExternalLBServer).Sync(&externalLBSyncServer{stream})
}

type ExternalLB_SyncServer interface {
	Send(*UpdateService) error
	Recv() (*Data, error)
	grpc.ServerStream
}

type externalLBSyncServer struct {
	grpc.ServerStream
}

func (x *externalLBSyncServer) Send(m *UpdateService) error {
	return x.ServerStream.SendMsg(m)
}

func (x *externalLBSyncServer) Recv() (*Data, error) {
	m := new(Data)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _ExternalLB_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.ExternalLB",
	HandlerType: (*ExternalLBServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _ExternalLB_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _ExternalLB_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ExternalLB_Delete_Handler,
		},
		{
			MethodName: "NodesChange",
			Handler:    _ExternalLB_NodesChange_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Sync",
			Handler:       _ExternalLB_Sync_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "externallb/externalLB.proto",
}

func init() {
	proto.RegisterFile("externallb/externalLB.proto", fileDescriptor_externalLB_3d625253ec2af4f9)
}

var fileDescriptor_externalLB_3d625253ec2af4f9 = []byte{
	// 414 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xc1, 0x6e, 0x13, 0x31,
	0x10, 0xad, 0xbb, 0xbb, 0x21, 0x99, 0x80, 0x90, 0x0c, 0x07, 0x2b, 0xed, 0xc1, 0x5a, 0x71, 0xd8,
	0x0b, 0x01, 0xca, 0x17, 0x40, 0xd3, 0x43, 0x10, 0x8a, 0x22, 0x07, 0xb8, 0xbb, 0xde, 0x51, 0x89,
	0xb4, 0x5d, 0x47, 0xb6, 0x83, 0xe0, 0xdf, 0xf8, 0x18, 0x3e, 0x05, 0xcd, 0x38, 0x49, 0xd3, 0xaa,
	0x52, 0x4f, 0x3b, 0xef, 0xcd, 0xdb, 0xf1, 0xf3, 0xf3, 0xc0, 0x19, 0xfe, 0x4e, 0x18, 0x7a, 0xdb,
	0x75, 0xd7, 0xef, 0xf6, 0xe5, 0xd7, 0xcf, 0xd3, 0x4d, 0xf0, 0xc9, 0xcb, 0xf2, 0x26, 0x6c, 0x5c,
	0xfd, 0xf7, 0x14, 0xca, 0x99, 0x4d, 0x56, 0x6a, 0x18, 0xaf, 0x30, 0xfc, 0x5a, 0x3b, 0x5c, 0xd8,
	0x5b, 0x54, 0x42, 0x8b, 0x66, 0x64, 0x8e, 0x29, 0x79, 0x0e, 0x23, 0xfa, 0xc6, 0x8d, 0x75, 0xa8,
	0x4e, 0xb9, 0x7f, 0x47, 0xc8, 0x06, 0x5e, 0x1a, 0x8c, 0x7e, 0x1b, 0x1c, 0xfe, 0xc0, 0x10, 0xd7,
	0xbe, 0x57, 0x05, 0x6b, 0x1e, 0xd2, 0x72, 0x02, 0xc3, 0x25, 0x39, 0x70, 0xbe, 0x53, 0x25, 0x4b,
	0x0e, 0x98, 0x5c, 0x5c, 0xed, 0x8c, 0xce, 0x97, 0x51, 0x55, 0xba, 0x20, 0x17, 0x47, 0x14, 0xfd,
	0x6d, 0xfc, 0x36, 0x61, 0x98, 0xcf, 0xd4, 0x40, 0x8b, 0xa6, 0x32, 0x07, 0x4c, 0xbd, 0xd5, 0x9f,
	0xde, 0x7d, 0x5b, 0xdf, 0xa2, 0x7a, 0xa6, 0x45, 0x53, 0x98, 0x03, 0x26, 0xf7, 0xf3, 0x78, 0x19,
	0xd0, 0x26, 0x6c, 0xd5, 0x50, 0x8b, 0x66, 0x68, 0xee, 0x08, 0xf9, 0x1a, 0xaa, 0xde, 0xb7, 0x18,
	0xd5, 0x88, 0x4f, 0xcc, 0x40, 0x6a, 0xa8, 0x96, 0x3e, 0xa4, 0xa8, 0x40, 0x17, 0xcd, 0xf8, 0x02,
	0xa6, 0x14, 0xd9, 0x94, 0x28, 0x93, 0x1b, 0xf5, 0x17, 0x28, 0xa9, 0x90, 0x12, 0xca, 0xa3, 0xd8,
	0xb8, 0x26, 0x8e, 0x7a, 0x1c, 0x55, 0x65, 0xb2, 0x6e, 0x02, 0xc3, 0x85, 0x6f, 0x91, 0xf9, 0x22,
	0xbb, 0xdf, 0xe3, 0xfa, 0x1c, 0x06, 0x06, 0xe3, 0xb6, 0xe3, 0x69, 0xb6, 0x6d, 0xc3, 0x7e, 0x1a,
	0xd5, 0xf5, 0x19, 0x54, 0x0b, 0x36, 0x25, 0xa1, 0xec, 0xd6, 0x31, 0x29, 0xc1, 0x4e, 0xb9, 0xae,
	0x3f, 0xc1, 0x8b, 0xef, 0x9b, 0xd6, 0x26, 0xdc, 0xbd, 0xd7, 0x63, 0x13, 0x72, 0x02, 0xab, 0xad,
	0x73, 0x18, 0x23, 0x9b, 0xe2, 0x04, 0x76, 0xc4, 0xc5, 0x3f, 0x01, 0x70, 0x75, 0xd8, 0x11, 0xf9,
	0x06, 0x06, 0x39, 0x1b, 0xb9, 0xbb, 0x35, 0x2d, 0xc9, 0xe4, 0x79, 0xae, 0xb3, 0xcd, 0xfa, 0x84,
	0x54, 0xf9, 0xdc, 0xa7, 0x54, 0x33, 0xec, 0xf0, 0x09, 0xd5, 0x07, 0x18, 0xf3, 0x05, 0x2f, 0x7f,
	0xda, 0xfe, 0x06, 0xe5, 0x38, 0xb7, 0x99, 0x9a, 0xbc, 0xca, 0xe0, 0xde, 0x1d, 0xeb, 0x13, 0xf9,
	0x16, 0x4a, 0x7a, 0xdf, 0x7b, 0x63, 0x1f, 0x97, 0x36, 0xe2, 0xbd, 0xb8, 0x1e, 0xf0, 0xe2, 0x7f,
	0xfc, 0x1f, 0x00, 0x00, 0xff, 0xff, 0x4b, 0x62, 0x96, 0x3c, 0x17, 0x03, 0x00, 0x00,
}
