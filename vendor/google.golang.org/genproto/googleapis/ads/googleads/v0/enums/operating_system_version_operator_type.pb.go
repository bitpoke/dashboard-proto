// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v0/enums/operating_system_version_operator_type.proto

package enums // import "google.golang.org/genproto/googleapis/ads/googleads/v0/enums"

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

// The type of operating system version.
type OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType int32

const (
	// Not specified.
	OperatingSystemVersionOperatorTypeEnum_UNSPECIFIED OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType = 0
	// Used for return value only. Represents value unknown in this version.
	OperatingSystemVersionOperatorTypeEnum_UNKNOWN OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType = 1
	// Equals to the specified version.
	OperatingSystemVersionOperatorTypeEnum_EQUALS_TO OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType = 2
	// Greater than or equals to the specified version.
	OperatingSystemVersionOperatorTypeEnum_GREATER_THAN_EQUALS_TO OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType = 4
)

var OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "UNKNOWN",
	2: "EQUALS_TO",
	4: "GREATER_THAN_EQUALS_TO",
}
var OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType_value = map[string]int32{
	"UNSPECIFIED":            0,
	"UNKNOWN":                1,
	"EQUALS_TO":              2,
	"GREATER_THAN_EQUALS_TO": 4,
}

func (x OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType) String() string {
	return proto.EnumName(OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType_name, int32(x))
}
func (OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_operating_system_version_operator_type_f746666ec051bec7, []int{0, 0}
}

// Container for enum describing the type of OS operators.
type OperatingSystemVersionOperatorTypeEnum struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OperatingSystemVersionOperatorTypeEnum) Reset() {
	*m = OperatingSystemVersionOperatorTypeEnum{}
}
func (m *OperatingSystemVersionOperatorTypeEnum) String() string { return proto.CompactTextString(m) }
func (*OperatingSystemVersionOperatorTypeEnum) ProtoMessage()    {}
func (*OperatingSystemVersionOperatorTypeEnum) Descriptor() ([]byte, []int) {
	return fileDescriptor_operating_system_version_operator_type_f746666ec051bec7, []int{0}
}
func (m *OperatingSystemVersionOperatorTypeEnum) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OperatingSystemVersionOperatorTypeEnum.Unmarshal(m, b)
}
func (m *OperatingSystemVersionOperatorTypeEnum) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OperatingSystemVersionOperatorTypeEnum.Marshal(b, m, deterministic)
}
func (dst *OperatingSystemVersionOperatorTypeEnum) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OperatingSystemVersionOperatorTypeEnum.Merge(dst, src)
}
func (m *OperatingSystemVersionOperatorTypeEnum) XXX_Size() int {
	return xxx_messageInfo_OperatingSystemVersionOperatorTypeEnum.Size(m)
}
func (m *OperatingSystemVersionOperatorTypeEnum) XXX_DiscardUnknown() {
	xxx_messageInfo_OperatingSystemVersionOperatorTypeEnum.DiscardUnknown(m)
}

var xxx_messageInfo_OperatingSystemVersionOperatorTypeEnum proto.InternalMessageInfo

func init() {
	proto.RegisterType((*OperatingSystemVersionOperatorTypeEnum)(nil), "google.ads.googleads.v0.enums.OperatingSystemVersionOperatorTypeEnum")
	proto.RegisterEnum("google.ads.googleads.v0.enums.OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType", OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType_name, OperatingSystemVersionOperatorTypeEnum_OperatingSystemVersionOperatorType_value)
}

func init() {
	proto.RegisterFile("google/ads/googleads/v0/enums/operating_system_version_operator_type.proto", fileDescriptor_operating_system_version_operator_type_f746666ec051bec7)
}

var fileDescriptor_operating_system_version_operator_type_f746666ec051bec7 = []byte{
	// 320 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0xc1, 0x4a, 0xc3, 0x30,
	0x1c, 0xc6, 0x6d, 0x15, 0xc5, 0x0c, 0x71, 0xf4, 0xe0, 0x41, 0xd8, 0x61, 0x3b, 0xe8, 0x2d, 0x2d,
	0x78, 0x8b, 0xa7, 0x4c, 0xe3, 0x9c, 0x4a, 0x3b, 0xb7, 0xae, 0x82, 0x14, 0x42, 0xb5, 0x21, 0x0c,
	0xd6, 0xa4, 0x34, 0xdd, 0xa0, 0x4f, 0xe2, 0xdd, 0xa3, 0x8f, 0xe2, 0xa3, 0xf8, 0x0c, 0x1e, 0xa4,
	0xc9, 0x5a, 0x4f, 0xba, 0x4b, 0xf8, 0xc8, 0xf7, 0xe7, 0xc7, 0xc7, 0x0f, 0xdc, 0x71, 0x29, 0xf9,
	0x92, 0xb9, 0x49, 0xaa, 0x5c, 0x13, 0xeb, 0xb4, 0xf6, 0x5c, 0x26, 0x56, 0x99, 0x72, 0x65, 0xce,
	0x8a, 0xa4, 0x5c, 0x08, 0x4e, 0x55, 0xa5, 0x4a, 0x96, 0xd1, 0x35, 0x2b, 0xd4, 0x42, 0x0a, 0x6a,
	0x0a, 0x59, 0xd0, 0xb2, 0xca, 0x19, 0xcc, 0x0b, 0x59, 0x4a, 0xa7, 0x67, 0x00, 0x30, 0x49, 0x15,
	0x6c, 0x59, 0x70, 0xed, 0x41, 0xcd, 0x1a, 0xbc, 0x59, 0xe0, 0x2c, 0x68, 0x78, 0x33, 0x8d, 0x8b,
	0x0c, 0x2d, 0xd8, 0xc0, 0xc2, 0x2a, 0x67, 0x44, 0xac, 0xb2, 0x41, 0x06, 0x06, 0xdb, 0x2f, 0x9d,
	0x63, 0xd0, 0x99, 0xfb, 0xb3, 0x09, 0xb9, 0x1a, 0xdf, 0x8c, 0xc9, 0x75, 0x77, 0xc7, 0xe9, 0x80,
	0x83, 0xb9, 0x7f, 0xef, 0x07, 0x4f, 0x7e, 0xd7, 0x72, 0x8e, 0xc0, 0x21, 0x79, 0x9c, 0xe3, 0x87,
	0x19, 0x0d, 0x83, 0xae, 0xed, 0x9c, 0x82, 0x93, 0xd1, 0x94, 0xe0, 0x90, 0x4c, 0x69, 0x78, 0x8b,
	0x7d, 0xfa, 0xdb, 0xed, 0x0d, 0xbf, 0x2d, 0xd0, 0x7f, 0x95, 0x19, 0xfc, 0x77, 0xff, 0xf0, 0x7c,
	0xfb, 0xa4, 0x49, 0xed, 0x61, 0x62, 0x3d, 0x0f, 0x37, 0x24, 0x2e, 0x97, 0x89, 0xe0, 0x50, 0x16,
	0xdc, 0xe5, 0x4c, 0x68, 0x4b, 0x8d, 0xe5, 0x7c, 0xa1, 0xfe, 0x90, 0x7e, 0xa9, 0xdf, 0x77, 0x7b,
	0x77, 0x84, 0xf1, 0x87, 0xdd, 0x1b, 0x19, 0x14, 0x4e, 0x15, 0x34, 0xb1, 0x4e, 0x91, 0x07, 0x6b,
	0x51, 0xea, 0xb3, 0xe9, 0x63, 0x9c, 0xaa, 0xb8, 0xed, 0xe3, 0xc8, 0x8b, 0x75, 0xff, 0x65, 0xf7,
	0xcd, 0x27, 0x42, 0x38, 0x55, 0x08, 0xb5, 0x17, 0x08, 0x45, 0x1e, 0x42, 0xfa, 0xe6, 0x65, 0x5f,
	0x0f, 0xbb, 0xf8, 0x09, 0x00, 0x00, 0xff, 0xff, 0xdc, 0x17, 0x07, 0xf6, 0x0c, 0x02, 0x00, 0x00,
}
