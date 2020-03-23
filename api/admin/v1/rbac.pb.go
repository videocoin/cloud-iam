// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/admin/v1/rbac.proto

package v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// A view for Role objects.
type RoleView int32

const (
	// Omits the `included_permissions` field.
	// This is the default value.
	RoleView_BASIC RoleView = 0
	// Returns all fields.
	RoleView_FULL RoleView = 1
)

var RoleView_name = map[int32]string{
	0: "BASIC",
	1: "FULL",
}

var RoleView_value = map[string]int32{
	"BASIC": 0,
	"FULL":  1,
}

func (x RoleView) String() string {
	return proto.EnumName(RoleView_name, int32(x))
}

func (RoleView) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8d7d9f8f370131d7, []int{0}
}

// A role in the Identity and Access Management API.
type Role struct {
	// The name of the role.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Optional A human-readable title for the role.  Typically this
	// is limited to 100 UTF-8 bytes.
	Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	// Optional A human-readable description for the role.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// The names of the permissions this role grants when bound in an IAM policy.
	IncludedPermissions  []string `protobuf:"bytes,4,rep,name=included_permissions,json=includedPermissions,proto3" json:"included_permissions,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Role) Reset()         { *m = Role{} }
func (m *Role) String() string { return proto.CompactTextString(m) }
func (*Role) ProtoMessage()    {}
func (*Role) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7d9f8f370131d7, []int{0}
}

func (m *Role) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Role.Unmarshal(m, b)
}
func (m *Role) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Role.Marshal(b, m, deterministic)
}
func (m *Role) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Role.Merge(m, src)
}
func (m *Role) XXX_Size() int {
	return xxx_messageInfo_Role.Size(m)
}
func (m *Role) XXX_DiscardUnknown() {
	xxx_messageInfo_Role.DiscardUnknown(m)
}

var xxx_messageInfo_Role proto.InternalMessageInfo

func (m *Role) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Role) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Role) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Role) GetIncludedPermissions() []string {
	if m != nil {
		return m.IncludedPermissions
	}
	return nil
}

// The request to get all roles defined under a resource.
type ListRolesRequest struct {
	// Optional limit on the number of roles to include in the response.
	PageSize int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// Optional pagination token returned in an earlier ListRolesResponse.
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	// Optional view for the returned Role objects. When `FULL` is specified,
	// the `includedPermissions` field is returned, which includes a list of all
	// permissions in the role. The default value is `BASIC`, which does not
	// return the `includedPermissions` field.
	View                 RoleView `protobuf:"varint,3,opt,name=view,proto3,enum=videocoin.iam.admin.v1.RoleView" json:"view,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRolesRequest) Reset()         { *m = ListRolesRequest{} }
func (m *ListRolesRequest) String() string { return proto.CompactTextString(m) }
func (*ListRolesRequest) ProtoMessage()    {}
func (*ListRolesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7d9f8f370131d7, []int{1}
}

func (m *ListRolesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRolesRequest.Unmarshal(m, b)
}
func (m *ListRolesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRolesRequest.Marshal(b, m, deterministic)
}
func (m *ListRolesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRolesRequest.Merge(m, src)
}
func (m *ListRolesRequest) XXX_Size() int {
	return xxx_messageInfo_ListRolesRequest.Size(m)
}
func (m *ListRolesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRolesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRolesRequest proto.InternalMessageInfo

func (m *ListRolesRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ListRolesRequest) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

func (m *ListRolesRequest) GetView() RoleView {
	if m != nil {
		return m.View
	}
	return RoleView_BASIC
}

// The response containing the roles defined under a resource.
type ListRolesResponse struct {
	// The list of predefined roles.
	Roles []*Role `protobuf:"bytes,1,rep,name=roles,proto3" json:"roles,omitempty"`
	// To retrieve the next page of results, set
	// `ListRolesRequest.page_token` to this value.
	NextPageToken        string   `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRolesResponse) Reset()         { *m = ListRolesResponse{} }
func (m *ListRolesResponse) String() string { return proto.CompactTextString(m) }
func (*ListRolesResponse) ProtoMessage()    {}
func (*ListRolesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7d9f8f370131d7, []int{2}
}

func (m *ListRolesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRolesResponse.Unmarshal(m, b)
}
func (m *ListRolesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRolesResponse.Marshal(b, m, deterministic)
}
func (m *ListRolesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRolesResponse.Merge(m, src)
}
func (m *ListRolesResponse) XXX_Size() int {
	return xxx_messageInfo_ListRolesResponse.Size(m)
}
func (m *ListRolesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRolesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListRolesResponse proto.InternalMessageInfo

func (m *ListRolesResponse) GetRoles() []*Role {
	if m != nil {
		return m.Roles
	}
	return nil
}

func (m *ListRolesResponse) GetNextPageToken() string {
	if m != nil {
		return m.NextPageToken
	}
	return ""
}

type GetRoleRequest struct {
	// The name of the role.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRoleRequest) Reset()         { *m = GetRoleRequest{} }
func (m *GetRoleRequest) String() string { return proto.CompactTextString(m) }
func (*GetRoleRequest) ProtoMessage()    {}
func (*GetRoleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7d9f8f370131d7, []int{3}
}

func (m *GetRoleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRoleRequest.Unmarshal(m, b)
}
func (m *GetRoleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRoleRequest.Marshal(b, m, deterministic)
}
func (m *GetRoleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRoleRequest.Merge(m, src)
}
func (m *GetRoleRequest) XXX_Size() int {
	return xxx_messageInfo_GetRoleRequest.Size(m)
}
func (m *GetRoleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRoleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRoleRequest proto.InternalMessageInfo

func (m *GetRoleRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// A permission belongs to one or more roles.
type Permission struct {
	// The name of this Permission.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The title of this Permission.
	Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	// A brief description of what this Permission is used for.
	Description          string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Permission) Reset()         { *m = Permission{} }
func (m *Permission) String() string { return proto.CompactTextString(m) }
func (*Permission) ProtoMessage()    {}
func (*Permission) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7d9f8f370131d7, []int{4}
}

func (m *Permission) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Permission.Unmarshal(m, b)
}
func (m *Permission) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Permission.Marshal(b, m, deterministic)
}
func (m *Permission) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Permission.Merge(m, src)
}
func (m *Permission) XXX_Size() int {
	return xxx_messageInfo_Permission.Size(m)
}
func (m *Permission) XXX_DiscardUnknown() {
	xxx_messageInfo_Permission.DiscardUnknown(m)
}

var xxx_messageInfo_Permission proto.InternalMessageInfo

func (m *Permission) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Permission) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Permission) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func init() {
	proto.RegisterEnum("videocoin.iam.admin.v1.RoleView", RoleView_name, RoleView_value)
	proto.RegisterType((*Role)(nil), "videocoin.iam.admin.v1.Role")
	proto.RegisterType((*ListRolesRequest)(nil), "videocoin.iam.admin.v1.ListRolesRequest")
	proto.RegisterType((*ListRolesResponse)(nil), "videocoin.iam.admin.v1.ListRolesResponse")
	proto.RegisterType((*GetRoleRequest)(nil), "videocoin.iam.admin.v1.GetRoleRequest")
	proto.RegisterType((*Permission)(nil), "videocoin.iam.admin.v1.Permission")
}

func init() { proto.RegisterFile("api/admin/v1/rbac.proto", fileDescriptor_8d7d9f8f370131d7) }

var fileDescriptor_8d7d9f8f370131d7 = []byte{
	// 470 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x14, 0xc4, 0x89, 0x03, 0xf1, 0x8b, 0x28, 0xe9, 0x23, 0xa2, 0x91, 0x29, 0xaa, 0xe5, 0x43, 0x95,
	0xf6, 0x60, 0x2b, 0x81, 0x2b, 0x87, 0xa4, 0x12, 0x08, 0x29, 0x87, 0xca, 0x05, 0x84, 0xb8, 0x44,
	0x9b, 0xf8, 0x11, 0x56, 0x75, 0x76, 0x8d, 0x77, 0xeb, 0xa2, 0x22, 0x2e, 0x88, 0x13, 0x57, 0x7e,
	0x1a, 0x7f, 0x80, 0x03, 0x07, 0x7e, 0x06, 0xda, 0x75, 0xd3, 0x86, 0x8f, 0x54, 0x48, 0x3d, 0x7a,
	0x66, 0xf6, 0xcd, 0x78, 0xde, 0x2e, 0x6c, 0xb1, 0x9c, 0xc7, 0x2c, 0x5d, 0x70, 0x11, 0x97, 0xfd,
	0xb8, 0x98, 0xb2, 0x59, 0x94, 0x17, 0x52, 0x4b, 0xbc, 0x57, 0xf2, 0x94, 0xe4, 0x4c, 0x72, 0x11,
	0x71, 0xb6, 0x88, 0xac, 0x24, 0x2a, 0xfb, 0xfe, 0xce, 0x5c, 0xca, 0x79, 0x46, 0xb1, 0x39, 0xf7,
	0x86, 0x53, 0x96, 0x4e, 0xa6, 0xf4, 0x96, 0x95, 0x5c, 0x16, 0xd5, 0x41, 0x7f, 0x7b, 0x45, 0xc0,
	0x84, 0x90, 0x9a, 0x69, 0x2e, 0x85, 0xaa, 0xd8, 0xf0, 0x8b, 0x03, 0x6e, 0x22, 0x33, 0xc2, 0x2d,
	0x70, 0x05, 0x5b, 0x50, 0xd7, 0x09, 0x9c, 0x9e, 0x37, 0xaa, 0x7f, 0x1f, 0xd6, 0x12, 0x0b, 0x60,
	0x07, 0x1a, 0x9a, 0xeb, 0x8c, 0xba, 0x35, 0xc3, 0x24, 0xd5, 0x07, 0x06, 0xd0, 0x4a, 0x49, 0xcd,
	0x0a, 0x9e, 0x9b, 0x69, 0xdd, 0xba, 0xe5, 0x56, 0x21, 0xec, 0x43, 0x87, 0x8b, 0x59, 0x76, 0x92,
	0x52, 0x3a, 0xc9, 0xa9, 0x58, 0x70, 0xa5, 0x8c, 0x6f, 0xd7, 0x0d, 0xea, 0x3d, 0x2f, 0xb9, 0xbb,
	0xe4, 0x0e, 0x2f, 0xa9, 0xf0, 0xb3, 0x03, 0xed, 0x31, 0x57, 0xda, 0x04, 0x52, 0x09, 0xbd, 0x3b,
	0x21, 0xa5, 0xf1, 0x3e, 0x78, 0x39, 0x9b, 0xd3, 0x44, 0xf1, 0xb3, 0x2a, 0x5d, 0x23, 0x69, 0x1a,
	0xe0, 0x88, 0x9f, 0x11, 0x3e, 0x00, 0xb0, 0xa4, 0x96, 0xc7, 0x24, 0xce, 0x13, 0x5a, 0xf9, 0x73,
	0x03, 0xe0, 0x23, 0x70, 0x4b, 0x4e, 0xa7, 0x36, 0xde, 0xc6, 0x20, 0x88, 0xfe, 0xdd, 0x61, 0x64,
	0xfc, 0x5e, 0x72, 0x3a, 0x4d, 0xac, 0x3a, 0x94, 0xb0, 0xb9, 0x92, 0x42, 0xe5, 0x52, 0x28, 0xc2,
	0x01, 0x34, 0x0a, 0x03, 0x74, 0x9d, 0xa0, 0xde, 0x6b, 0x0d, 0xb6, 0xaf, 0x9a, 0x95, 0x54, 0x52,
	0xdc, 0x85, 0x3b, 0x82, 0xde, 0xeb, 0xc9, 0x5f, 0x11, 0x6f, 0x1b, 0xf8, 0x70, 0x19, 0x33, 0xdc,
	0x83, 0x8d, 0xa7, 0x64, 0xfd, 0x96, 0x3f, 0xbd, 0x6e, 0x1b, 0xe1, 0x2b, 0x80, 0xcb, 0xc6, 0x10,
	0x57, 0x65, 0xd7, 0xdb, 0xd7, 0xfe, 0x0e, 0x34, 0x97, 0x3d, 0xa0, 0x07, 0x8d, 0xd1, 0xf0, 0xe8,
	0xd9, 0x41, 0xfb, 0x06, 0x36, 0xc1, 0x7d, 0xf2, 0x62, 0x3c, 0x6e, 0x3b, 0x83, 0x9f, 0xe6, 0xaa,
	0x8c, 0x86, 0x07, 0xa8, 0xc0, 0xbb, 0xe8, 0x07, 0x7b, 0xeb, 0x8a, 0xf8, 0x73, 0x91, 0xfe, 0xde,
	0x7f, 0x28, 0xab, 0xb2, 0xc3, 0xcd, 0x4f, 0xdf, 0x7e, 0x7c, 0xad, 0xb5, 0xd0, 0xb3, 0x8f, 0xc0,
	0xfa, 0x1c, 0xc3, 0xad, 0xf3, 0x8e, 0x70, 0x77, 0xdd, 0xa0, 0xdf, 0x4b, 0xf4, 0xaf, 0xdc, 0x51,
	0xe8, 0x5b, 0x8f, 0x0e, 0xe2, 0x85, 0x47, 0xfc, 0xc1, 0x14, 0xf8, 0x78, 0xff, 0xe3, 0xc8, 0x7d,
	0x5d, 0x2b, 0xfb, 0xd3, 0x9b, 0xf6, 0x89, 0x3c, 0xfc, 0x15, 0x00, 0x00, 0xff, 0xff, 0xca, 0x7e,
	0x5e, 0x22, 0x94, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RBACClient is the client API for RBAC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RBACClient interface {
	// Lists [Roles][cloud.api.iam.v1.Roles].
	ListRoles(ctx context.Context, in *ListRolesRequest, opts ...grpc.CallOption) (*ListRolesResponse, error)
	// Gets a [Role][cloud.api.iam.v1.Role].
	GetRole(ctx context.Context, in *GetRoleRequest, opts ...grpc.CallOption) (*Role, error)
}

type rBACClient struct {
	cc grpc.ClientConnInterface
}

func NewRBACClient(cc grpc.ClientConnInterface) RBACClient {
	return &rBACClient{cc}
}

func (c *rBACClient) ListRoles(ctx context.Context, in *ListRolesRequest, opts ...grpc.CallOption) (*ListRolesResponse, error) {
	out := new(ListRolesResponse)
	err := c.cc.Invoke(ctx, "/videocoin.iam.admin.v1.RBAC/ListRoles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rBACClient) GetRole(ctx context.Context, in *GetRoleRequest, opts ...grpc.CallOption) (*Role, error) {
	out := new(Role)
	err := c.cc.Invoke(ctx, "/videocoin.iam.admin.v1.RBAC/GetRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RBACServer is the server API for RBAC service.
type RBACServer interface {
	// Lists [Roles][cloud.api.iam.v1.Roles].
	ListRoles(context.Context, *ListRolesRequest) (*ListRolesResponse, error)
	// Gets a [Role][cloud.api.iam.v1.Role].
	GetRole(context.Context, *GetRoleRequest) (*Role, error)
}

// UnimplementedRBACServer can be embedded to have forward compatible implementations.
type UnimplementedRBACServer struct {
}

func (*UnimplementedRBACServer) ListRoles(ctx context.Context, req *ListRolesRequest) (*ListRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRoles not implemented")
}
func (*UnimplementedRBACServer) GetRole(ctx context.Context, req *GetRoleRequest) (*Role, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRole not implemented")
}

func RegisterRBACServer(s *grpc.Server, srv RBACServer) {
	s.RegisterService(&_RBAC_serviceDesc, srv)
}

func _RBAC_ListRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RBACServer).ListRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/videocoin.iam.admin.v1.RBAC/ListRoles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RBACServer).ListRoles(ctx, req.(*ListRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RBAC_GetRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RBACServer).GetRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/videocoin.iam.admin.v1.RBAC/GetRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RBACServer).GetRole(ctx, req.(*GetRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RBAC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "videocoin.iam.admin.v1.RBAC",
	HandlerType: (*RBACServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListRoles",
			Handler:    _RBAC_ListRoles_Handler,
		},
		{
			MethodName: "GetRole",
			Handler:    _RBAC_GetRole_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/admin/v1/rbac.proto",
}
