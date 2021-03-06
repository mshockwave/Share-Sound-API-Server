// Code generated by protoc-gen-go.
// source: image_attachment.proto
// DO NOT EDIT!

package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ImageAttachment struct {
	Mime    string `protobuf:"bytes,1,opt,name=mime" json:"mime,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Content []byte `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
}

func (m *ImageAttachment) Reset()                    { *m = ImageAttachment{} }
func (m *ImageAttachment) String() string            { return proto.CompactTextString(m) }
func (*ImageAttachment) ProtoMessage()               {}
func (*ImageAttachment) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func init() {
	proto.RegisterType((*ImageAttachment)(nil), "protos.ImageAttachment")
}

var fileDescriptor1 = []byte{
	// 103 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0xcb, 0xcc, 0x4d, 0x4c,
	0x4f, 0x8d, 0x4f, 0x2c, 0x29, 0x49, 0x4c, 0xce, 0xc8, 0x4d, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0x62, 0x03, 0x53, 0xc5, 0x4a, 0x76, 0x5c, 0xfc, 0x9e, 0x20, 0x15, 0x8e, 0x70,
	0x05, 0x42, 0x3c, 0x5c, 0x2c, 0xb9, 0x99, 0xb9, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x20,
	0x5e, 0x5e, 0x22, 0x90, 0xc7, 0x04, 0xe6, 0xf1, 0x73, 0xb1, 0x27, 0xe7, 0xe7, 0x95, 0x00, 0x95,
	0x49, 0x30, 0x03, 0x05, 0x78, 0x92, 0x20, 0xe6, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xac,
	0xcf, 0x8f, 0xe2, 0x68, 0x00, 0x00, 0x00,
}
