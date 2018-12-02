// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mesg.proto

package main

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Words struct {
	Words                []string `protobuf:"bytes,1,rep,name=words,proto3" json:"words,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Words) Reset()         { *m = Words{} }
func (m *Words) String() string { return proto.CompactTextString(m) }
func (*Words) ProtoMessage()    {}
func (*Words) Descriptor() ([]byte, []int) {
	return fileDescriptor_80f9ee673fba9935, []int{0}
}

func (m *Words) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Words.Unmarshal(m, b)
}
func (m *Words) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Words.Marshal(b, m, deterministic)
}
func (m *Words) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Words.Merge(m, src)
}
func (m *Words) XXX_Size() int {
	return xxx_messageInfo_Words.Size(m)
}
func (m *Words) XXX_DiscardUnknown() {
	xxx_messageInfo_Words.DiscardUnknown(m)
}

var xxx_messageInfo_Words proto.InternalMessageInfo

func (m *Words) GetWords() []string {
	if m != nil {
		return m.Words
	}
	return nil
}

func init() {
	proto.RegisterType((*Words)(nil), "main.Words")
}

func init() { proto.RegisterFile("mesg.proto", fileDescriptor_80f9ee673fba9935) }

var fileDescriptor_80f9ee673fba9935 = []byte{
	// 72 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xca, 0x4d, 0x2d, 0x4e,
	0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xc9, 0x4d, 0xcc, 0xcc, 0x53, 0x92, 0xe5, 0x62,
	0x0d, 0xcf, 0x2f, 0x4a, 0x29, 0x16, 0x12, 0xe1, 0x62, 0x2d, 0x07, 0x31, 0x24, 0x18, 0x15, 0x98,
	0x35, 0x38, 0x83, 0x20, 0x9c, 0x24, 0x36, 0xb0, 0x5a, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xeb, 0xb7, 0xab, 0xfc, 0x39, 0x00, 0x00, 0x00,
}