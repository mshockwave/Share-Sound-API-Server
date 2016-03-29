// Code generated by protoc-gen-go.
// source: audio_attachment.proto
// DO NOT EDIT!

/*
Package protos is a generated protocol buffer package.

It is generated from these files:
	audio_attachment.proto
	image_attachment.proto
	location.proto
	story.proto

It has these top-level messages:
	AudioAttachment
	ImageAttachment
	Location
	Story
*/
package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type AudioAttachment struct {
	Mime      string    `protobuf:"bytes,1,opt,name=mime" json:"mime,omitempty"`
	Name      string    `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	TimeStamp uint64    `protobuf:"varint,3,opt,name=time_stamp" json:"time_stamp,omitempty"`
	Location  *Location `protobuf:"bytes,4,opt,name=location" json:"location,omitempty"`
	Content   []byte    `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`
}

func (m *AudioAttachment) Reset()                    { *m = AudioAttachment{} }
func (m *AudioAttachment) String() string            { return proto.CompactTextString(m) }
func (*AudioAttachment) ProtoMessage()               {}
func (*AudioAttachment) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *AudioAttachment) GetLocation() *Location {
	if m != nil {
		return m.Location
	}
	return nil
}

func init() {
	proto.RegisterType((*AudioAttachment)(nil), "protos.AudioAttachment")
}

var fileDescriptor0 = []byte{
	// 153 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x4b, 0x2c, 0x4d, 0xc9,
	0xcc, 0x8f, 0x4f, 0x2c, 0x29, 0x49, 0x4c, 0xce, 0xc8, 0x4d, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0x62, 0x03, 0x53, 0xc5, 0x52, 0x7c, 0x39, 0xf9, 0xc9, 0x89, 0x25, 0x99, 0xf9,
	0x79, 0x10, 0x71, 0xa5, 0x32, 0x2e, 0x7e, 0x47, 0x90, 0x0e, 0x47, 0xb8, 0x06, 0x21, 0x1e, 0x2e,
	0x96, 0xdc, 0xcc, 0xdc, 0x54, 0x09, 0x46, 0x05, 0x46, 0x0d, 0x4e, 0x10, 0x2f, 0x2f, 0x11, 0xc8,
	0x63, 0x02, 0xf3, 0x84, 0xb8, 0xb8, 0x4a, 0x80, 0x72, 0xf1, 0xc5, 0x25, 0x89, 0xb9, 0x05, 0x12,
	0xcc, 0x40, 0x31, 0x16, 0x21, 0x25, 0x2e, 0x0e, 0x98, 0xa1, 0x12, 0x2c, 0x40, 0x11, 0x6e, 0x23,
	0x01, 0x88, 0xe1, 0xc5, 0x7a, 0x3e, 0x50, 0x71, 0x21, 0x7e, 0x2e, 0xf6, 0xe4, 0xfc, 0xbc, 0x12,
	0xa0, 0xf1, 0x12, 0xac, 0x40, 0x25, 0x3c, 0x49, 0x10, 0xf7, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x85, 0xd7, 0xfc, 0xe1, 0xb0, 0x00, 0x00, 0x00,
}