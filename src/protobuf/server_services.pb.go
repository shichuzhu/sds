// Code generated by protoc-gen-go. DO NOT EDIT.
// source: server_services.proto

/*
Package serverservices is a generated protocol buffer package.

It is generated from these files:
	server_services.proto

It has these top-level messages:
	ConfigInfo
	StringMessage
	StringArray
	IntMessage
	DetectorMessage
	Member
	FullMembershipList
	UDPMessage
	FileTransMessage
*/
package serverservices

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

type ConfigInfo struct {
	LogLevel int32 `protobuf:"varint,1,opt,name=log_level,json=logLevel" json:"log_level,omitempty"`
	VmIndex  int32 `protobuf:"varint,2,opt,name=vm_index,json=vmIndex" json:"vm_index,omitempty"`
}

func (m *ConfigInfo) Reset()                    { *m = ConfigInfo{} }
func (m *ConfigInfo) String() string            { return proto.CompactTextString(m) }
func (*ConfigInfo) ProtoMessage()               {}
func (*ConfigInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ConfigInfo) GetLogLevel() int32 {
	if m != nil {
		return m.LogLevel
	}
	return 0
}

func (m *ConfigInfo) GetVmIndex() int32 {
	if m != nil {
		return m.VmIndex
	}
	return 0
}

type StringMessage struct {
	Mesg string `protobuf:"bytes,1,opt,name=mesg" json:"mesg,omitempty"`
}

func (m *StringMessage) Reset()                    { *m = StringMessage{} }
func (m *StringMessage) String() string            { return proto.CompactTextString(m) }
func (*StringMessage) ProtoMessage()               {}
func (*StringMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *StringMessage) GetMesg() string {
	if m != nil {
		return m.Mesg
	}
	return ""
}

type StringArray struct {
	Mesgs []string `protobuf:"bytes,1,rep,name=mesgs" json:"mesgs,omitempty"`
}

func (m *StringArray) Reset()                    { *m = StringArray{} }
func (m *StringArray) String() string            { return proto.CompactTextString(m) }
func (*StringArray) ProtoMessage()               {}
func (*StringArray) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *StringArray) GetMesgs() []string {
	if m != nil {
		return m.Mesgs
	}
	return nil
}

type IntMessage struct {
	Mesg int32 `protobuf:"varint,1,opt,name=mesg" json:"mesg,omitempty"`
}

func (m *IntMessage) Reset()                    { *m = IntMessage{} }
func (m *IntMessage) String() string            { return proto.CompactTextString(m) }
func (*IntMessage) ProtoMessage()               {}
func (*IntMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *IntMessage) GetMesg() int32 {
	if m != nil {
		return m.Mesg
	}
	return 0
}

type DetectorMessage struct {
	Header  string `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Addr    string `protobuf:"bytes,2,opt,name=addr" json:"addr,omitempty"`
	SessNum int32  `protobuf:"varint,3,opt,name=sess_num,json=sessNum" json:"sess_num,omitempty"`
	Ttl     int32  `protobuf:"varint,4,opt,name=ttl" json:"ttl,omitempty"`
	NodeId  int32  `protobuf:"varint,5,opt,name=node_id,json=nodeId" json:"node_id,omitempty"`
}

func (m *DetectorMessage) Reset()                    { *m = DetectorMessage{} }
func (m *DetectorMessage) String() string            { return proto.CompactTextString(m) }
func (*DetectorMessage) ProtoMessage()               {}
func (*DetectorMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *DetectorMessage) GetHeader() string {
	if m != nil {
		return m.Header
	}
	return ""
}

func (m *DetectorMessage) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *DetectorMessage) GetSessNum() int32 {
	if m != nil {
		return m.SessNum
	}
	return 0
}

func (m *DetectorMessage) GetTtl() int32 {
	if m != nil {
		return m.Ttl
	}
	return 0
}

func (m *DetectorMessage) GetNodeId() int32 {
	if m != nil {
		return m.NodeId
	}
	return 0
}

type Member struct {
	Addr    string `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	SessNum int32  `protobuf:"varint,2,opt,name=sess_num,json=sessNum" json:"sess_num,omitempty"`
	NodeId  int32  `protobuf:"varint,3,opt,name=node_id,json=nodeId" json:"node_id,omitempty"`
}

func (m *Member) Reset()                    { *m = Member{} }
func (m *Member) String() string            { return proto.CompactTextString(m) }
func (*Member) ProtoMessage()               {}
func (*Member) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Member) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *Member) GetSessNum() int32 {
	if m != nil {
		return m.SessNum
	}
	return 0
}

func (m *Member) GetNodeId() int32 {
	if m != nil {
		return m.NodeId
	}
	return 0
}

type FullMembershipList struct {
	List []*Member `protobuf:"bytes,1,rep,name=list" json:"list,omitempty"`
}

func (m *FullMembershipList) Reset()                    { *m = FullMembershipList{} }
func (m *FullMembershipList) String() string            { return proto.CompactTextString(m) }
func (*FullMembershipList) ProtoMessage()               {}
func (*FullMembershipList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *FullMembershipList) GetList() []*Member {
	if m != nil {
		return m.List
	}
	return nil
}

type UDPMessage struct {
	MessageType string              `protobuf:"bytes,1,opt,name=message_type,json=messageType" json:"message_type,omitempty"`
	Dm          *DetectorMessage    `protobuf:"bytes,2,opt,name=dm" json:"dm,omitempty"`
	Fm          *FullMembershipList `protobuf:"bytes,3,opt,name=fm" json:"fm,omitempty"`
}

func (m *UDPMessage) Reset()                    { *m = UDPMessage{} }
func (m *UDPMessage) String() string            { return proto.CompactTextString(m) }
func (*UDPMessage) ProtoMessage()               {}
func (*UDPMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *UDPMessage) GetMessageType() string {
	if m != nil {
		return m.MessageType
	}
	return ""
}

func (m *UDPMessage) GetDm() *DetectorMessage {
	if m != nil {
		return m.Dm
	}
	return nil
}

func (m *UDPMessage) GetFm() *FullMembershipList {
	if m != nil {
		return m.Fm
	}
	return nil
}

// SDFS related messages
type FileTransMessage struct {
	// Types that are valid to be assigned to FileTransMessage:
	//	*FileTransMessage_Chunk
	//	*FileTransMessage_Config_
	FileTransMessage isFileTransMessage_FileTransMessage `protobuf_oneof:"file_trans_message"`
}

func (m *FileTransMessage) Reset()                    { *m = FileTransMessage{} }
func (m *FileTransMessage) String() string            { return proto.CompactTextString(m) }
func (*FileTransMessage) ProtoMessage()               {}
func (*FileTransMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type isFileTransMessage_FileTransMessage interface {
	isFileTransMessage_FileTransMessage()
}

type FileTransMessage_Chunk struct {
	Chunk []byte `protobuf:"bytes,1,opt,name=chunk,proto3,oneof"`
}
type FileTransMessage_Config_ struct {
	Config *FileTransMessage_Config `protobuf:"bytes,2,opt,name=config,oneof"`
}

func (*FileTransMessage_Chunk) isFileTransMessage_FileTransMessage()   {}
func (*FileTransMessage_Config_) isFileTransMessage_FileTransMessage() {}

func (m *FileTransMessage) GetFileTransMessage() isFileTransMessage_FileTransMessage {
	if m != nil {
		return m.FileTransMessage
	}
	return nil
}

func (m *FileTransMessage) GetChunk() []byte {
	if x, ok := m.GetFileTransMessage().(*FileTransMessage_Chunk); ok {
		return x.Chunk
	}
	return nil
}

func (m *FileTransMessage) GetConfig() *FileTransMessage_Config {
	if x, ok := m.GetFileTransMessage().(*FileTransMessage_Config_); ok {
		return x.Config
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*FileTransMessage) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _FileTransMessage_OneofMarshaler, _FileTransMessage_OneofUnmarshaler, _FileTransMessage_OneofSizer, []interface{}{
		(*FileTransMessage_Chunk)(nil),
		(*FileTransMessage_Config_)(nil),
	}
}

func _FileTransMessage_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*FileTransMessage)
	// file_trans_message
	switch x := m.FileTransMessage.(type) {
	case *FileTransMessage_Chunk:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		b.EncodeRawBytes(x.Chunk)
	case *FileTransMessage_Config_:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Config); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("FileTransMessage.FileTransMessage has unexpected type %T", x)
	}
	return nil
}

func _FileTransMessage_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*FileTransMessage)
	switch tag {
	case 1: // file_trans_message.chunk
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeRawBytes(true)
		m.FileTransMessage = &FileTransMessage_Chunk{x}
		return true, err
	case 2: // file_trans_message.config
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FileTransMessage_Config)
		err := b.DecodeMessage(msg)
		m.FileTransMessage = &FileTransMessage_Config_{msg}
		return true, err
	default:
		return false, nil
	}
}

func _FileTransMessage_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*FileTransMessage)
	// file_trans_message
	switch x := m.FileTransMessage.(type) {
	case *FileTransMessage_Chunk:
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Chunk)))
		n += len(x.Chunk)
	case *FileTransMessage_Config_:
		s := proto.Size(x.Config)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type FileTransMessage_Config struct {
	RemoteFilepath string `protobuf:"bytes,1,opt,name=remote_filepath,json=remoteFilepath" json:"remote_filepath,omitempty"`
	RepNumber      int32  `protobuf:"varint,2,opt,name=rep_number,json=repNumber" json:"rep_number,omitempty"`
}

func (m *FileTransMessage_Config) Reset()                    { *m = FileTransMessage_Config{} }
func (m *FileTransMessage_Config) String() string            { return proto.CompactTextString(m) }
func (*FileTransMessage_Config) ProtoMessage()               {}
func (*FileTransMessage_Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8, 0} }

func (m *FileTransMessage_Config) GetRemoteFilepath() string {
	if m != nil {
		return m.RemoteFilepath
	}
	return ""
}

func (m *FileTransMessage_Config) GetRepNumber() int32 {
	if m != nil {
		return m.RepNumber
	}
	return 0
}

func init() {
	proto.RegisterType((*ConfigInfo)(nil), "serverservices.ConfigInfo")
	proto.RegisterType((*StringMessage)(nil), "serverservices.StringMessage")
	proto.RegisterType((*StringArray)(nil), "serverservices.StringArray")
	proto.RegisterType((*IntMessage)(nil), "serverservices.IntMessage")
	proto.RegisterType((*DetectorMessage)(nil), "serverservices.DetectorMessage")
	proto.RegisterType((*Member)(nil), "serverservices.Member")
	proto.RegisterType((*FullMembershipList)(nil), "serverservices.FullMembershipList")
	proto.RegisterType((*UDPMessage)(nil), "serverservices.UDPMessage")
	proto.RegisterType((*FileTransMessage)(nil), "serverservices.FileTransMessage")
	proto.RegisterType((*FileTransMessage_Config)(nil), "serverservices.FileTransMessage.Config")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ServerServices service

type ServerServicesClient interface {
	ServerConfig(ctx context.Context, in *ConfigInfo, opts ...grpc.CallOption) (*StringMessage, error)
	ReturnMatches(ctx context.Context, in *StringArray, opts ...grpc.CallOption) (ServerServices_ReturnMatchesClient, error)
	CloseServer(ctx context.Context, in *IntMessage, opts ...grpc.CallOption) (*StringMessage, error)
	ActMembership(ctx context.Context, in *StringArray, opts ...grpc.CallOption) (*StringMessage, error)
	// SDFS client stub
	SdfsCall(ctx context.Context, in *StringArray, opts ...grpc.CallOption) (*StringArray, error)
	// SDFS related functions
	TransferFiles(ctx context.Context, opts ...grpc.CallOption) (ServerServices_TransferFilesClient, error)
	PutFile(ctx context.Context, in *StringMessage, opts ...grpc.CallOption) (*IntMessage, error)
}

type serverServicesClient struct {
	cc *grpc.ClientConn
}

func NewServerServicesClient(cc *grpc.ClientConn) ServerServicesClient {
	return &serverServicesClient{cc}
}

func (c *serverServicesClient) ServerConfig(ctx context.Context, in *ConfigInfo, opts ...grpc.CallOption) (*StringMessage, error) {
	out := new(StringMessage)
	err := grpc.Invoke(ctx, "/serverservices.ServerServices/ServerConfig", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverServicesClient) ReturnMatches(ctx context.Context, in *StringArray, opts ...grpc.CallOption) (ServerServices_ReturnMatchesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_ServerServices_serviceDesc.Streams[0], c.cc, "/serverservices.ServerServices/ReturnMatches", opts...)
	if err != nil {
		return nil, err
	}
	x := &serverServicesReturnMatchesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ServerServices_ReturnMatchesClient interface {
	Recv() (*StringMessage, error)
	grpc.ClientStream
}

type serverServicesReturnMatchesClient struct {
	grpc.ClientStream
}

func (x *serverServicesReturnMatchesClient) Recv() (*StringMessage, error) {
	m := new(StringMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *serverServicesClient) CloseServer(ctx context.Context, in *IntMessage, opts ...grpc.CallOption) (*StringMessage, error) {
	out := new(StringMessage)
	err := grpc.Invoke(ctx, "/serverservices.ServerServices/CloseServer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverServicesClient) ActMembership(ctx context.Context, in *StringArray, opts ...grpc.CallOption) (*StringMessage, error) {
	out := new(StringMessage)
	err := grpc.Invoke(ctx, "/serverservices.ServerServices/ActMembership", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverServicesClient) SdfsCall(ctx context.Context, in *StringArray, opts ...grpc.CallOption) (*StringArray, error) {
	out := new(StringArray)
	err := grpc.Invoke(ctx, "/serverservices.ServerServices/SdfsCall", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverServicesClient) TransferFiles(ctx context.Context, opts ...grpc.CallOption) (ServerServices_TransferFilesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_ServerServices_serviceDesc.Streams[1], c.cc, "/serverservices.ServerServices/TransferFiles", opts...)
	if err != nil {
		return nil, err
	}
	x := &serverServicesTransferFilesClient{stream}
	return x, nil
}

type ServerServices_TransferFilesClient interface {
	Send(*FileTransMessage) error
	CloseAndRecv() (*IntMessage, error)
	grpc.ClientStream
}

type serverServicesTransferFilesClient struct {
	grpc.ClientStream
}

func (x *serverServicesTransferFilesClient) Send(m *FileTransMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *serverServicesTransferFilesClient) CloseAndRecv() (*IntMessage, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(IntMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *serverServicesClient) PutFile(ctx context.Context, in *StringMessage, opts ...grpc.CallOption) (*IntMessage, error) {
	out := new(IntMessage)
	err := grpc.Invoke(ctx, "/serverservices.ServerServices/PutFile", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ServerServices service

type ServerServicesServer interface {
	ServerConfig(context.Context, *ConfigInfo) (*StringMessage, error)
	ReturnMatches(*StringArray, ServerServices_ReturnMatchesServer) error
	CloseServer(context.Context, *IntMessage) (*StringMessage, error)
	ActMembership(context.Context, *StringArray) (*StringMessage, error)
	// SDFS client stub
	SdfsCall(context.Context, *StringArray) (*StringArray, error)
	// SDFS related functions
	TransferFiles(ServerServices_TransferFilesServer) error
	PutFile(context.Context, *StringMessage) (*IntMessage, error)
}

func RegisterServerServicesServer(s *grpc.Server, srv ServerServicesServer) {
	s.RegisterService(&_ServerServices_serviceDesc, srv)
}

func _ServerServices_ServerConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServicesServer).ServerConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/serverservices.ServerServices/ServerConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServicesServer).ServerConfig(ctx, req.(*ConfigInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerServices_ReturnMatches_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StringArray)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ServerServicesServer).ReturnMatches(m, &serverServicesReturnMatchesServer{stream})
}

type ServerServices_ReturnMatchesServer interface {
	Send(*StringMessage) error
	grpc.ServerStream
}

type serverServicesReturnMatchesServer struct {
	grpc.ServerStream
}

func (x *serverServicesReturnMatchesServer) Send(m *StringMessage) error {
	return x.ServerStream.SendMsg(m)
}

func _ServerServices_CloseServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IntMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServicesServer).CloseServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/serverservices.ServerServices/CloseServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServicesServer).CloseServer(ctx, req.(*IntMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerServices_ActMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StringArray)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServicesServer).ActMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/serverservices.ServerServices/ActMembership",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServicesServer).ActMembership(ctx, req.(*StringArray))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerServices_SdfsCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StringArray)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServicesServer).SdfsCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/serverservices.ServerServices/SdfsCall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServicesServer).SdfsCall(ctx, req.(*StringArray))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerServices_TransferFiles_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServerServicesServer).TransferFiles(&serverServicesTransferFilesServer{stream})
}

type ServerServices_TransferFilesServer interface {
	SendAndClose(*IntMessage) error
	Recv() (*FileTransMessage, error)
	grpc.ServerStream
}

type serverServicesTransferFilesServer struct {
	grpc.ServerStream
}

func (x *serverServicesTransferFilesServer) SendAndClose(m *IntMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *serverServicesTransferFilesServer) Recv() (*FileTransMessage, error) {
	m := new(FileTransMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ServerServices_PutFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StringMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServicesServer).PutFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/serverservices.ServerServices/PutFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServicesServer).PutFile(ctx, req.(*StringMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _ServerServices_serviceDesc = grpc.ServiceDesc{
	ServiceName: "serverservices.ServerServices",
	HandlerType: (*ServerServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ServerConfig",
			Handler:    _ServerServices_ServerConfig_Handler,
		},
		{
			MethodName: "CloseServer",
			Handler:    _ServerServices_CloseServer_Handler,
		},
		{
			MethodName: "ActMembership",
			Handler:    _ServerServices_ActMembership_Handler,
		},
		{
			MethodName: "SdfsCall",
			Handler:    _ServerServices_SdfsCall_Handler,
		},
		{
			MethodName: "PutFile",
			Handler:    _ServerServices_PutFile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReturnMatches",
			Handler:       _ServerServices_ReturnMatches_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "TransferFiles",
			Handler:       _ServerServices_TransferFiles_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "server_services.proto",
}

func init() { proto.RegisterFile("server_services.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 616 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xdb, 0x52, 0x1a, 0x4d,
	0x10, 0x66, 0x39, 0x29, 0x8d, 0xa8, 0xd5, 0xe5, 0xef, 0x4f, 0xb0, 0xac, 0x90, 0xf5, 0x42, 0x2b,
	0x17, 0x24, 0x45, 0x5e, 0x20, 0x46, 0x8b, 0x92, 0x44, 0x0c, 0x59, 0xcc, 0xf5, 0xd6, 0xca, 0xf6,
	0xc2, 0x56, 0x66, 0x0f, 0x35, 0x33, 0x50, 0xf1, 0x3a, 0xaf, 0x90, 0xc7, 0xc8, 0xeb, 0xe4, 0x7d,
	0x52, 0x73, 0xb0, 0x10, 0x34, 0x70, 0x91, 0xab, 0x9d, 0xfe, 0xa6, 0xfb, 0xeb, 0x6f, 0xfa, 0xb0,
	0xf0, 0x9f, 0x20, 0x3e, 0x27, 0xee, 0xab, 0x4f, 0x3c, 0x26, 0xd1, 0xc9, 0x79, 0x26, 0x33, 0xdc,
	0x35, 0xf0, 0x03, 0xea, 0x5e, 0x02, 0x5c, 0x64, 0x69, 0x14, 0x4f, 0xfa, 0x69, 0x94, 0xe1, 0x11,
	0xd4, 0x58, 0x36, 0xf1, 0x19, 0xcd, 0x89, 0x35, 0x9d, 0xb6, 0x73, 0x56, 0xf1, 0xb6, 0x59, 0x36,
	0xb9, 0x56, 0x36, 0xbe, 0x80, 0xed, 0x79, 0xe2, 0xc7, 0x69, 0x48, 0xdf, 0x9b, 0x45, 0x7d, 0xb7,
	0x35, 0x4f, 0xfa, 0xca, 0x74, 0x4f, 0xa0, 0x31, 0x92, 0x3c, 0x4e, 0x27, 0x03, 0x12, 0x22, 0x98,
	0x10, 0x22, 0x94, 0x13, 0x12, 0x13, 0xcd, 0x51, 0xf3, 0xf4, 0xd9, 0x3d, 0x81, 0xba, 0x71, 0x3a,
	0xe7, 0x3c, 0xb8, 0xc7, 0x03, 0xa8, 0x28, 0x58, 0x34, 0x9d, 0x76, 0xe9, 0xac, 0xe6, 0x19, 0xc3,
	0x6d, 0x03, 0xf4, 0x53, 0xf9, 0x1c, 0x4d, 0xc5, 0xd2, 0xfc, 0x70, 0x60, 0xef, 0x92, 0x24, 0x8d,
	0x65, 0xc6, 0x1f, 0xfc, 0x0e, 0xa1, 0x3a, 0xa5, 0x20, 0x24, 0x6e, 0x13, 0x5a, 0x4b, 0xc5, 0x07,
	0x61, 0xc8, 0xb5, 0xdc, 0x9a, 0xa7, 0xcf, 0xea, 0x19, 0x82, 0x84, 0xf0, 0xd3, 0x59, 0xd2, 0x2c,
	0x99, 0x67, 0x28, 0xfb, 0x66, 0x96, 0xe0, 0x3e, 0x94, 0xa4, 0x64, 0xcd, 0xb2, 0x46, 0xd5, 0x11,
	0xff, 0x87, 0xad, 0x34, 0x0b, 0xc9, 0x8f, 0xc3, 0x66, 0x45, 0xa3, 0x55, 0x65, 0xf6, 0x43, 0x77,
	0x08, 0xd5, 0x01, 0x25, 0x77, 0x8f, 0x72, 0x38, 0x7f, 0xc9, 0x51, 0x5c, 0xce, 0xf1, 0x88, 0xb1,
	0xb4, 0xc4, 0xf8, 0x1e, 0xb0, 0x37, 0x63, 0xcc, 0xb0, 0x8a, 0x69, 0x9c, 0x5f, 0xc7, 0x42, 0xe2,
	0x6b, 0x28, 0xb3, 0x58, 0x48, 0x5d, 0xa4, 0x7a, 0xf7, 0xb0, 0xb3, 0xdc, 0xbe, 0x8e, 0xf1, 0xf6,
	0xb4, 0x8f, 0xfb, 0xd3, 0x01, 0xf8, 0x7a, 0x39, 0x7c, 0x28, 0xca, 0x2b, 0xd8, 0x49, 0xcc, 0xd1,
	0x97, 0xf7, 0x39, 0x59, 0x81, 0x75, 0x8b, 0xdd, 0xde, 0xe7, 0x84, 0x6f, 0xa0, 0x18, 0x1a, 0x85,
	0xf5, 0xee, 0xcb, 0x55, 0xee, 0x95, 0x22, 0x7b, 0xc5, 0x30, 0xc1, 0x2e, 0x14, 0x23, 0x53, 0xb6,
	0x7a, 0xd7, 0x5d, 0x0d, 0x78, 0x2a, 0xdf, 0x2b, 0x46, 0x89, 0xfb, 0xdb, 0x81, 0xfd, 0x5e, 0xcc,
	0xe8, 0x96, 0x07, 0xa9, 0x58, 0x74, 0xac, 0x32, 0x9e, 0xce, 0xd2, 0x6f, 0x5a, 0xd5, 0xce, 0x55,
	0xc1, 0x33, 0x26, 0x9e, 0x43, 0x75, 0xac, 0xe7, 0xd1, 0xaa, 0x3a, 0x7d, 0x92, 0x64, 0x85, 0xa9,
	0x63, 0xc6, 0xf7, 0xaa, 0xe0, 0xd9, 0xc0, 0xd6, 0x10, 0xaa, 0x06, 0xc3, 0x53, 0xd8, 0xe3, 0x94,
	0x64, 0x92, 0xfc, 0x28, 0x66, 0x94, 0x07, 0x72, 0x6a, 0x8b, 0xb0, 0x6b, 0xe0, 0x9e, 0x45, 0xf1,
	0x18, 0x80, 0x53, 0xae, 0xda, 0x75, 0x47, 0xdc, 0x76, 0xac, 0xc6, 0x29, 0xbf, 0xd1, 0xc0, 0x87,
	0x03, 0x40, 0x45, 0xe0, 0x4b, 0x95, 0xd7, 0xb7, 0x05, 0xec, 0xfe, 0x2a, 0xc3, 0xee, 0x48, 0x8b,
	0x1b, 0x59, 0x71, 0xf8, 0x09, 0x76, 0x0c, 0x62, 0x05, 0xb4, 0x56, 0xd5, 0x2f, 0x76, 0xad, 0x75,
	0xbc, 0x7a, 0xb7, 0xb4, 0x41, 0x6e, 0x01, 0x3f, 0x43, 0xc3, 0x23, 0x39, 0xe3, 0xe9, 0x20, 0x90,
	0xe3, 0x29, 0x09, 0x3c, 0x7a, 0x3e, 0x42, 0xaf, 0xd3, 0x46, 0xba, 0xb7, 0x0e, 0x7e, 0x84, 0xfa,
	0x05, 0xcb, 0x04, 0x19, 0x89, 0x4f, 0xc5, 0x2d, 0x16, 0x6f, 0xb3, 0xb8, 0x01, 0x34, 0xce, 0xc7,
	0x72, 0xd1, 0xed, 0x7f, 0x13, 0x87, 0x3d, 0xd8, 0x1e, 0x85, 0x91, 0xb8, 0x08, 0x18, 0x5b, 0xcf,
	0xb4, 0xee, 0xd2, 0x2d, 0xe0, 0x17, 0x68, 0xe8, 0xe1, 0x88, 0x88, 0xab, 0xe6, 0x0a, 0x6c, 0x6f,
	0x9a, 0x9f, 0xd6, 0x9a, 0x32, 0xb8, 0x85, 0x33, 0x07, 0x7b, 0xb0, 0x35, 0x9c, 0x49, 0x15, 0x86,
	0xeb, 0x9f, 0xb1, 0x9e, 0xe9, 0xae, 0xaa, 0x7f, 0xc0, 0xef, 0xfe, 0x04, 0x00, 0x00, 0xff, 0xff,
	0x4d, 0xef, 0xae, 0x8d, 0x99, 0x05, 0x00, 0x00,
}
