// Code generated by protoc-gen-go. DO NOT EDIT.
// source: basepb.proto

package basepb

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// 基础结构，主要提供给tcp内部使用，目前客户端用不到
type Base struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Base) Reset()         { *m = Base{} }
func (m *Base) String() string { return proto.CompactTextString(m) }
func (*Base) ProtoMessage()    {}
func (*Base) Descriptor() ([]byte, []int) {
	return fileDescriptor_acce15d0fbf96fc9, []int{0}
}

func (m *Base) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Base.Unmarshal(m, b)
}
func (m *Base) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Base.Marshal(b, m, deterministic)
}
func (m *Base) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Base.Merge(m, src)
}
func (m *Base) XXX_Size() int {
	return xxx_messageInfo_Base.Size(m)
}
func (m *Base) XXX_DiscardUnknown() {
	xxx_messageInfo_Base.DiscardUnknown(m)
}

var xxx_messageInfo_Base proto.InternalMessageInfo

// rpc 无消息返回时可以返回这个值
type Base_Success struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Base_Success) Reset()         { *m = Base_Success{} }
func (m *Base_Success) String() string { return proto.CompactTextString(m) }
func (*Base_Success) ProtoMessage()    {}
func (*Base_Success) Descriptor() ([]byte, []int) {
	return fileDescriptor_acce15d0fbf96fc9, []int{0, 0}
}

func (m *Base_Success) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Base_Success.Unmarshal(m, b)
}
func (m *Base_Success) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Base_Success.Marshal(b, m, deterministic)
}
func (m *Base_Success) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Base_Success.Merge(m, src)
}
func (m *Base_Success) XXX_Size() int {
	return xxx_messageInfo_Base_Success.Size(m)
}
func (m *Base_Success) XXX_DiscardUnknown() {
	xxx_messageInfo_Base_Success.DiscardUnknown(m)
}

var xxx_messageInfo_Base_Success proto.InternalMessageInfo

// rpc 发生错误时可以返回这个值
type Base_Error struct {
	ErrorCode            int64    `protobuf:"varint,1,opt,name=error_code,json=errorCode,proto3" json:"error_code,omitempty"`
	ErrorMessage         string   `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	Fields               []string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Base_Error) Reset()         { *m = Base_Error{} }
func (m *Base_Error) String() string { return proto.CompactTextString(m) }
func (*Base_Error) ProtoMessage()    {}
func (*Base_Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_acce15d0fbf96fc9, []int{0, 1}
}

func (m *Base_Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Base_Error.Unmarshal(m, b)
}
func (m *Base_Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Base_Error.Marshal(b, m, deterministic)
}
func (m *Base_Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Base_Error.Merge(m, src)
}
func (m *Base_Error) XXX_Size() int {
	return xxx_messageInfo_Base_Error.Size(m)
}
func (m *Base_Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Base_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Base_Error proto.InternalMessageInfo

func (m *Base_Error) GetErrorCode() int64 {
	if m != nil {
		return m.ErrorCode
	}
	return 0
}

func (m *Base_Error) GetErrorMessage() string {
	if m != nil {
		return m.ErrorMessage
	}
	return ""
}

func (m *Base_Error) GetFields() []string {
	if m != nil {
		return m.Fields
	}
	return nil
}

// ping client 发给服务器
type Base_Ping struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Base_Ping) Reset()         { *m = Base_Ping{} }
func (m *Base_Ping) String() string { return proto.CompactTextString(m) }
func (*Base_Ping) ProtoMessage()    {}
func (*Base_Ping) Descriptor() ([]byte, []int) {
	return fileDescriptor_acce15d0fbf96fc9, []int{0, 2}
}

func (m *Base_Ping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Base_Ping.Unmarshal(m, b)
}
func (m *Base_Ping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Base_Ping.Marshal(b, m, deterministic)
}
func (m *Base_Ping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Base_Ping.Merge(m, src)
}
func (m *Base_Ping) XXX_Size() int {
	return xxx_messageInfo_Base_Ping.Size(m)
}
func (m *Base_Ping) XXX_DiscardUnknown() {
	xxx_messageInfo_Base_Ping.DiscardUnknown(m)
}

var xxx_messageInfo_Base_Ping proto.InternalMessageInfo

// Pong server 收到client 发来的ping时立马回复Pong
type Base_Pong struct {
	Now                  string   `protobuf:"bytes,1,opt,name=now,proto3" json:"now,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Base_Pong) Reset()         { *m = Base_Pong{} }
func (m *Base_Pong) String() string { return proto.CompactTextString(m) }
func (*Base_Pong) ProtoMessage()    {}
func (*Base_Pong) Descriptor() ([]byte, []int) {
	return fileDescriptor_acce15d0fbf96fc9, []int{0, 3}
}

func (m *Base_Pong) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Base_Pong.Unmarshal(m, b)
}
func (m *Base_Pong) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Base_Pong.Marshal(b, m, deterministic)
}
func (m *Base_Pong) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Base_Pong.Merge(m, src)
}
func (m *Base_Pong) XXX_Size() int {
	return xxx_messageInfo_Base_Pong.Size(m)
}
func (m *Base_Pong) XXX_DiscardUnknown() {
	xxx_messageInfo_Base_Pong.DiscardUnknown(m)
}

var xxx_messageInfo_Base_Pong proto.InternalMessageInfo

func (m *Base_Pong) GetNow() string {
	if m != nil {
		return m.Now
	}
	return ""
}

func init() {
	proto.RegisterType((*Base)(nil), "basepb.base")
	proto.RegisterType((*Base_Success)(nil), "basepb.base.Success")
	proto.RegisterType((*Base_Error)(nil), "basepb.base.Error")
	proto.RegisterType((*Base_Ping)(nil), "basepb.base.Ping")
	proto.RegisterType((*Base_Pong)(nil), "basepb.base.Pong")
}

func init() { proto.RegisterFile("basepb.proto", fileDescriptor_acce15d0fbf96fc9) }

var fileDescriptor_acce15d0fbf96fc9 = []byte{
	// 191 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x2c, 0x8f, 0x41, 0x4a, 0xc5, 0x30,
	0x10, 0x40, 0x89, 0xa9, 0x91, 0x0c, 0x15, 0x24, 0x0b, 0x09, 0x01, 0xa1, 0xe8, 0xa6, 0xab, 0x6c,
	0xbc, 0x81, 0xe2, 0x52, 0x28, 0xf1, 0x00, 0xd2, 0xa6, 0x63, 0x08, 0xb4, 0x9d, 0x92, 0x58, 0xbd,
	0x8a, 0xc7, 0x95, 0x26, 0x7f, 0xf7, 0xde, 0x9b, 0xc5, 0xcc, 0x40, 0x3b, 0x8d, 0x19, 0xf7, 0xc9,
	0xee, 0x89, 0xbe, 0x49, 0x89, 0x6a, 0x8f, 0x7f, 0x0c, 0x9a, 0x13, 0x8d, 0x84, 0x9b, 0x8f, 0xc3,
	0x7b, 0xcc, 0xd9, 0x78, 0xb8, 0x7e, 0x4b, 0x89, 0x92, 0x7a, 0x00, 0xc0, 0x13, 0x3e, 0x3d, 0xcd,
	0xa8, 0x59, 0xc7, 0x7a, 0xee, 0x64, 0x29, 0xaf, 0x34, 0xa3, 0x7a, 0x82, 0xdb, 0x3a, 0x5e, 0x31,
	0xe7, 0x31, 0xa0, 0xbe, 0xea, 0x58, 0x2f, 0x5d, 0x5b, 0xe2, 0x7b, 0x6d, 0xea, 0x1e, 0xc4, 0x57,
	0xc4, 0x65, 0xce, 0x9a, 0x77, 0xbc, 0x97, 0xee, 0x62, 0x46, 0x40, 0x33, 0xc4, 0x2d, 0x18, 0x0d,
	0xcd, 0x40, 0x5b, 0x50, 0x77, 0xc0, 0x37, 0xfa, 0x2d, 0x4b, 0xa4, 0x3b, 0xf1, 0xc5, 0x80, 0xf6,
	0xb4, 0xda, 0x9c, 0x8e, 0x68, 0xc3, 0xb8, 0xe2, 0x12, 0x7f, 0xd0, 0xd6, 0xb3, 0x27, 0x51, 0xbe,
	0x78, 0xfe, 0x0f, 0x00, 0x00, 0xff, 0xff, 0xfe, 0xe9, 0xc0, 0x9b, 0xd5, 0x00, 0x00, 0x00,
}
