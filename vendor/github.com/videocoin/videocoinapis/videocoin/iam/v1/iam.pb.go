// Code generated by protoc-gen-go. DO NOT EDIT.
// source: videocoin/iam/v1/iam.proto

package v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type Key struct {
	// The resource id.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The private key data. Only provided in `CreateServiceAccountKey`
	// responses. Make sure to keep the private key data secure because it
	// allows for the assertion of the service account identity.
	PrivateKeyData []byte `protobuf:"bytes,2,opt,name=private_key_data,json=privateKeyData,proto3" json:"private_key_data,omitempty"`
	// The public key data. Only provided in `GetServiceAccountKey` responses.
	PublicKeyData []byte `protobuf:"bytes,3,opt,name=public_key_data,json=publicKeyData,proto3" json:"public_key_data,omitempty"`
	// The key can be used after this timestamp.
	ValidAfterTime *timestamp.Timestamp `protobuf:"bytes,4,opt,name=valid_after_time,json=validAfterTime,proto3" json:"valid_after_time,omitempty"`
	// The key can be used before this timestamp.
	ValidBeforeTime      *timestamp.Timestamp `protobuf:"bytes,5,opt,name=valid_before_time,json=validBeforeTime,proto3" json:"valid_before_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Key) Reset()         { *m = Key{} }
func (m *Key) String() string { return proto.CompactTextString(m) }
func (*Key) ProtoMessage()    {}
func (*Key) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d67148b8fa02d13, []int{0}
}

func (m *Key) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Key.Unmarshal(m, b)
}
func (m *Key) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Key.Marshal(b, m, deterministic)
}
func (m *Key) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Key.Merge(m, src)
}
func (m *Key) XXX_Size() int {
	return xxx_messageInfo_Key.Size(m)
}
func (m *Key) XXX_DiscardUnknown() {
	xxx_messageInfo_Key.DiscardUnknown(m)
}

var xxx_messageInfo_Key proto.InternalMessageInfo

func (m *Key) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Key) GetPrivateKeyData() []byte {
	if m != nil {
		return m.PrivateKeyData
	}
	return nil
}

func (m *Key) GetPublicKeyData() []byte {
	if m != nil {
		return m.PublicKeyData
	}
	return nil
}

func (m *Key) GetValidAfterTime() *timestamp.Timestamp {
	if m != nil {
		return m.ValidAfterTime
	}
	return nil
}

func (m *Key) GetValidBeforeTime() *timestamp.Timestamp {
	if m != nil {
		return m.ValidBeforeTime
	}
	return nil
}

// The keys list response.
type ListKeysRequest struct {
	// Optional limit on the number of roles to include in the response.
	PageSize int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// Optional pagination token returned in an earlier ListRolesResponse.
	PageToken            string   `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListKeysRequest) Reset()         { *m = ListKeysRequest{} }
func (m *ListKeysRequest) String() string { return proto.CompactTextString(m) }
func (*ListKeysRequest) ProtoMessage()    {}
func (*ListKeysRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d67148b8fa02d13, []int{1}
}

func (m *ListKeysRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListKeysRequest.Unmarshal(m, b)
}
func (m *ListKeysRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListKeysRequest.Marshal(b, m, deterministic)
}
func (m *ListKeysRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListKeysRequest.Merge(m, src)
}
func (m *ListKeysRequest) XXX_Size() int {
	return xxx_messageInfo_ListKeysRequest.Size(m)
}
func (m *ListKeysRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListKeysRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListKeysRequest proto.InternalMessageInfo

func (m *ListKeysRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ListKeysRequest) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

// The keys list response.
type ListKeysResponse struct {
	Keys                 []*Key   `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListKeysResponse) Reset()         { *m = ListKeysResponse{} }
func (m *ListKeysResponse) String() string { return proto.CompactTextString(m) }
func (*ListKeysResponse) ProtoMessage()    {}
func (*ListKeysResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d67148b8fa02d13, []int{2}
}

func (m *ListKeysResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListKeysResponse.Unmarshal(m, b)
}
func (m *ListKeysResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListKeysResponse.Marshal(b, m, deterministic)
}
func (m *ListKeysResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListKeysResponse.Merge(m, src)
}
func (m *ListKeysResponse) XXX_Size() int {
	return xxx_messageInfo_ListKeysResponse.Size(m)
}
func (m *ListKeysResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListKeysResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListKeysResponse proto.InternalMessageInfo

func (m *ListKeysResponse) GetKeys() []*Key {
	if m != nil {
		return m.Keys
	}
	return nil
}

// The key get by id request.
type GetKeyRequest struct {
	KeyId                string   `protobuf:"bytes,2,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetKeyRequest) Reset()         { *m = GetKeyRequest{} }
func (m *GetKeyRequest) String() string { return proto.CompactTextString(m) }
func (*GetKeyRequest) ProtoMessage()    {}
func (*GetKeyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d67148b8fa02d13, []int{3}
}

func (m *GetKeyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetKeyRequest.Unmarshal(m, b)
}
func (m *GetKeyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetKeyRequest.Marshal(b, m, deterministic)
}
func (m *GetKeyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetKeyRequest.Merge(m, src)
}
func (m *GetKeyRequest) XXX_Size() int {
	return xxx_messageInfo_GetKeyRequest.Size(m)
}
func (m *GetKeyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetKeyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetKeyRequest proto.InternalMessageInfo

func (m *GetKeyRequest) GetKeyId() string {
	if m != nil {
		return m.KeyId
	}
	return ""
}

// The key delete request.
type DeleteKeyRequest struct {
	KeyId                string   `protobuf:"bytes,1,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteKeyRequest) Reset()         { *m = DeleteKeyRequest{} }
func (m *DeleteKeyRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteKeyRequest) ProtoMessage()    {}
func (*DeleteKeyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d67148b8fa02d13, []int{4}
}

func (m *DeleteKeyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteKeyRequest.Unmarshal(m, b)
}
func (m *DeleteKeyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteKeyRequest.Marshal(b, m, deterministic)
}
func (m *DeleteKeyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteKeyRequest.Merge(m, src)
}
func (m *DeleteKeyRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteKeyRequest.Size(m)
}
func (m *DeleteKeyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteKeyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteKeyRequest proto.InternalMessageInfo

func (m *DeleteKeyRequest) GetKeyId() string {
	if m != nil {
		return m.KeyId
	}
	return ""
}

func init() {
	proto.RegisterType((*Key)(nil), "videocoin.iam.v1.Key")
	proto.RegisterType((*ListKeysRequest)(nil), "videocoin.iam.v1.ListKeysRequest")
	proto.RegisterType((*ListKeysResponse)(nil), "videocoin.iam.v1.ListKeysResponse")
	proto.RegisterType((*GetKeyRequest)(nil), "videocoin.iam.v1.GetKeyRequest")
	proto.RegisterType((*DeleteKeyRequest)(nil), "videocoin.iam.v1.DeleteKeyRequest")
}

func init() {
	proto.RegisterFile("videocoin/iam/v1/iam.proto", fileDescriptor_6d67148b8fa02d13)
}

var fileDescriptor_6d67148b8fa02d13 = []byte{
	// 620 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x86, 0x15, 0xa7, 0xad, 0x9a, 0x29, 0x6d, 0xc3, 0xa2, 0x42, 0xea, 0x82, 0x1a, 0x2c, 0x04,
	0xa1, 0x08, 0x5b, 0x29, 0xb7, 0x4a, 0x3d, 0xb8, 0x14, 0x50, 0x95, 0xf6, 0x62, 0x7a, 0xe2, 0x62,
	0x6d, 0xe2, 0x49, 0xbb, 0x8a, 0xed, 0x35, 0xf6, 0xda, 0x91, 0x8b, 0x7a, 0xe1, 0x15, 0x78, 0x02,
	0x5e, 0x87, 0x23, 0xdc, 0x90, 0x90, 0x38, 0x70, 0xe2, 0x29, 0xd0, 0xae, 0xed, 0x36, 0x24, 0x29,
	0x9c, 0x22, 0xcd, 0xff, 0xff, 0x33, 0x93, 0x6f, 0xc7, 0xa0, 0x67, 0xcc, 0x43, 0x3e, 0xe0, 0x2c,
	0xb4, 0x18, 0x0d, 0xac, 0xac, 0x2b, 0x7f, 0xcc, 0x28, 0xe6, 0x82, 0x93, 0xe6, 0x95, 0x66, 0xca,
	0x62, 0xd6, 0xd5, 0xb7, 0xcf, 0x38, 0x3f, 0xf3, 0xd1, 0x52, 0x7a, 0x3f, 0x1d, 0x5a, 0x82, 0x05,
	0x98, 0x08, 0x1a, 0x44, 0x45, 0x44, 0xdf, 0x9a, 0x36, 0x60, 0x10, 0x89, 0xbc, 0x14, 0xab, 0x34,
	0x8d, 0x98, 0x35, 0x64, 0xe8, 0x7b, 0x6e, 0x1f, 0xcf, 0x69, 0xc6, 0x78, 0x5c, 0x1a, 0x36, 0x27,
	0x0c, 0x31, 0x26, 0x3c, 0x8d, 0x07, 0x58, 0x4a, 0xf7, 0x27, 0x24, 0x1a, 0x86, 0x5c, 0x50, 0xc1,
	0x78, 0x98, 0x94, 0xea, 0xbd, 0x09, 0x75, 0xe0, 0x33, 0x0c, 0x45, 0x21, 0x18, 0x9f, 0x35, 0xa8,
	0xf7, 0x30, 0x27, 0x6b, 0xa0, 0x31, 0xaf, 0x55, 0x6b, 0xd7, 0x3a, 0x0d, 0x47, 0x63, 0x1e, 0xe9,
	0x40, 0x33, 0x8a, 0x59, 0x46, 0x05, 0xba, 0x23, 0xcc, 0x5d, 0x8f, 0x0a, 0xda, 0xd2, 0xda, 0xb5,
	0xce, 0x2d, 0x67, 0xad, 0xac, 0xf7, 0x30, 0x3f, 0xa4, 0x82, 0x92, 0xc7, 0xb0, 0x1e, 0xa5, 0x7d,
	0x9f, 0x0d, 0xae, 0x8d, 0x75, 0x65, 0x5c, 0x2d, 0xca, 0x95, 0xef, 0x10, 0x9a, 0x19, 0xf5, 0x99,
	0xe7, 0xd2, 0xa1, 0xc0, 0xd8, 0x95, 0x60, 0x5a, 0x0b, 0xed, 0x5a, 0x67, 0x65, 0x57, 0x37, 0x8b,
	0xed, 0xcc, 0x0a, 0x8a, 0x79, 0x5a, 0x51, 0x73, 0xd6, 0x54, 0xc6, 0x96, 0x11, 0x59, 0x24, 0xaf,
	0xe1, 0x76, 0xd1, 0xa5, 0x8f, 0x43, 0x1e, 0x63, 0xd1, 0x66, 0xf1, 0xbf, 0x6d, 0xd6, 0x55, 0xe8,
	0x40, 0x65, 0x64, 0x75, 0x6f, 0xe7, 0xb7, 0xfd, 0x04, 0x36, 0xd5, 0xab, 0x5d, 0xbd, 0x61, 0x88,
	0x62, 0xcc, 0xe3, 0x91, 0x25, 0x79, 0xc0, 0x08, 0xf3, 0xc4, 0xfa, 0x30, 0xc2, 0xfc, 0xd2, 0x38,
	0x81, 0xf5, 0x63, 0x96, 0x88, 0x1e, 0xe6, 0x89, 0x83, 0xef, 0x53, 0x4c, 0x04, 0xd9, 0x82, 0x46,
	0x44, 0xcf, 0xd0, 0x4d, 0xd8, 0x05, 0x2a, 0x6a, 0x8b, 0xce, 0xb2, 0x2c, 0xbc, 0x65, 0x17, 0x48,
	0x1e, 0x00, 0x28, 0x51, 0xf0, 0x11, 0x86, 0x8a, 0x5a, 0xc3, 0x51, 0xf6, 0x53, 0x59, 0x30, 0xf6,
	0xa1, 0x79, 0xdd, 0x2e, 0x89, 0x78, 0x98, 0x20, 0x79, 0x0a, 0x0b, 0x72, 0x60, 0xab, 0xd6, 0xae,
	0x77, 0x56, 0x76, 0x37, 0xcc, 0xe9, 0xc3, 0x32, 0x7b, 0x98, 0x3b, 0xca, 0x62, 0x3c, 0x83, 0xd5,
	0x37, 0x28, 0xd3, 0xd5, 0x2e, 0x3a, 0x2c, 0x49, 0xf2, 0xcc, 0x2b, 0x46, 0x1d, 0xd4, 0x7f, 0xda,
	0x9a, 0xb3, 0x38, 0xc2, 0xfc, 0xc8, 0x33, 0x4c, 0x68, 0x1e, 0xa2, 0x8f, 0xea, 0xb5, 0x66, 0xfd,
	0xb5, 0x69, 0xff, 0xee, 0x8f, 0x3a, 0xd4, 0x8f, 0xec, 0x13, 0x72, 0x0c, 0x8d, 0x97, 0x31, 0x16,
	0xaf, 0x4c, 0xee, 0xce, 0x80, 0x7d, 0x25, 0x8f, 0x56, 0x9f, 0xbf, 0xa6, 0xd1, 0xfc, 0xf8, 0xed,
	0xd7, 0x27, 0x0d, 0x8c, 0x65, 0xf9, 0xb1, 0xc8, 0x95, 0xc9, 0x00, 0x96, 0xab, 0x7f, 0x4c, 0x1e,
	0xce, 0x86, 0xa6, 0xe0, 0xea, 0xc6, 0xbf, 0x2c, 0x05, 0xb0, 0x6a, 0x08, 0xb9, 0x1e, 0x32, 0x84,
	0xa5, 0x82, 0x0b, 0xd9, 0x9e, 0xcd, 0xff, 0x45, 0xec, 0xa6, 0xc5, 0x1f, 0x7d, 0xb7, 0x4b, 0x32,
	0xaa, 0xf9, 0x06, 0xb9, 0x53, 0x35, 0x57, 0x87, 0xe0, 0x32, 0x6f, 0x7f, 0xe7, 0x92, 0x04, 0xd0,
	0xb8, 0x42, 0x4a, 0xe6, 0xac, 0x3a, 0xcd, 0x5b, 0xbf, 0x01, 0xdf, 0xcc, 0xb8, 0x9d, 0x79, 0xe3,
	0xf4, 0xe3, 0x2f, 0xf6, 0xc6, 0xdc, 0x33, 0xfd, 0x6a, 0x77, 0xcf, 0x85, 0x88, 0x92, 0x3d, 0xcb,
	0x1a, 0x8f, 0xc7, 0x73, 0xce, 0x98, 0xa6, 0xe2, 0xdc, 0x1a, 0xf8, 0x3c, 0xf5, 0x9e, 0x47, 0x3e,
	0x15, 0x43, 0x1e, 0x07, 0x07, 0x0b, 0xef, 0xb4, 0xac, 0xdb, 0x5f, 0x52, 0x9b, 0xbc, 0xf8, 0x13,
	0x00, 0x00, 0xff, 0xff, 0xb3, 0xd6, 0x24, 0x50, 0xdc, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// IAMClient is the client API for IAM service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IAMClient interface {
	// Creates a [Key][cloud.api.iam.v1.Key] and returns it.
	CreateKey(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Key, error)
	// Lists [Keys][cloud.api.iam.v1.Key].
	ListKeys(ctx context.Context, in *ListKeysRequest, opts ...grpc.CallOption) (*ListKeysResponse, error)
	// Gets the [Key][cloud.api.iam.v1.Key] by key id.
	GetKey(ctx context.Context, in *GetKeyRequest, opts ...grpc.CallOption) (*Key, error)
	// Deletes a [Key][cloud.api.iam.v1.Key] by key id..
	DeleteKey(ctx context.Context, in *DeleteKeyRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type iAMClient struct {
	cc grpc.ClientConnInterface
}

func NewIAMClient(cc grpc.ClientConnInterface) IAMClient {
	return &iAMClient{cc}
}

func (c *iAMClient) CreateKey(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Key, error) {
	out := new(Key)
	err := c.cc.Invoke(ctx, "/videocoin.iam.v1.IAM/CreateKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) ListKeys(ctx context.Context, in *ListKeysRequest, opts ...grpc.CallOption) (*ListKeysResponse, error) {
	out := new(ListKeysResponse)
	err := c.cc.Invoke(ctx, "/videocoin.iam.v1.IAM/ListKeys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) GetKey(ctx context.Context, in *GetKeyRequest, opts ...grpc.CallOption) (*Key, error) {
	out := new(Key)
	err := c.cc.Invoke(ctx, "/videocoin.iam.v1.IAM/GetKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMClient) DeleteKey(ctx context.Context, in *DeleteKeyRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/videocoin.iam.v1.IAM/DeleteKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IAMServer is the server API for IAM service.
type IAMServer interface {
	// Creates a [Key][cloud.api.iam.v1.Key] and returns it.
	CreateKey(context.Context, *empty.Empty) (*Key, error)
	// Lists [Keys][cloud.api.iam.v1.Key].
	ListKeys(context.Context, *ListKeysRequest) (*ListKeysResponse, error)
	// Gets the [Key][cloud.api.iam.v1.Key] by key id.
	GetKey(context.Context, *GetKeyRequest) (*Key, error)
	// Deletes a [Key][cloud.api.iam.v1.Key] by key id..
	DeleteKey(context.Context, *DeleteKeyRequest) (*empty.Empty, error)
}

// UnimplementedIAMServer can be embedded to have forward compatible implementations.
type UnimplementedIAMServer struct {
}

func (*UnimplementedIAMServer) CreateKey(ctx context.Context, req *empty.Empty) (*Key, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateKey not implemented")
}
func (*UnimplementedIAMServer) ListKeys(ctx context.Context, req *ListKeysRequest) (*ListKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListKeys not implemented")
}
func (*UnimplementedIAMServer) GetKey(ctx context.Context, req *GetKeyRequest) (*Key, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKey not implemented")
}
func (*UnimplementedIAMServer) DeleteKey(ctx context.Context, req *DeleteKeyRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteKey not implemented")
}

func RegisterIAMServer(s *grpc.Server, srv IAMServer) {
	s.RegisterService(&_IAM_serviceDesc, srv)
}

func _IAM_CreateKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).CreateKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/videocoin.iam.v1.IAM/CreateKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).CreateKey(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_ListKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).ListKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/videocoin.iam.v1.IAM/ListKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).ListKeys(ctx, req.(*ListKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_GetKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).GetKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/videocoin.iam.v1.IAM/GetKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).GetKey(ctx, req.(*GetKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAM_DeleteKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServer).DeleteKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/videocoin.iam.v1.IAM/DeleteKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServer).DeleteKey(ctx, req.(*DeleteKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _IAM_serviceDesc = grpc.ServiceDesc{
	ServiceName: "videocoin.iam.v1.IAM",
	HandlerType: (*IAMServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateKey",
			Handler:    _IAM_CreateKey_Handler,
		},
		{
			MethodName: "ListKeys",
			Handler:    _IAM_ListKeys_Handler,
		},
		{
			MethodName: "GetKey",
			Handler:    _IAM_GetKey_Handler,
		},
		{
			MethodName: "DeleteKey",
			Handler:    _IAM_DeleteKey_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "videocoin/iam/v1/iam.proto",
}
