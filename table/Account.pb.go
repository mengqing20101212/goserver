// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: table/Account.proto

package table

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RoleShow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoleId          uint64 `protobuf:"varint,1,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
	ServerId        uint32 `protobuf:"varint,2,opt,name=server_id,json=serverId,proto3" json:"server_id,omitempty"`
	RoleName        string `protobuf:"bytes,3,opt,name=role_name,json=roleName,proto3" json:"role_name,omitempty"`
	Lv              uint32 `protobuf:"varint,4,opt,name=lv,proto3" json:"lv,omitempty"`
	LastLoginTimer  uint32 `protobuf:"varint,5,opt,name=last_login_timer,json=lastLoginTimer,proto3" json:"last_login_timer,omitempty"`
	LastLogoutTimer uint32 `protobuf:"varint,6,opt,name=last_logout_timer,json=lastLogoutTimer,proto3" json:"last_logout_timer,omitempty"`
}

func (x *RoleShow) Reset() {
	*x = RoleShow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_table_Account_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoleShow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleShow) ProtoMessage() {}

func (x *RoleShow) ProtoReflect() protoreflect.Message {
	mi := &file_table_Account_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoleShow.ProtoReflect.Descriptor instead.
func (*RoleShow) Descriptor() ([]byte, []int) {
	return file_table_Account_proto_rawDescGZIP(), []int{0}
}

func (x *RoleShow) GetRoleId() uint64 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

func (x *RoleShow) GetServerId() uint32 {
	if x != nil {
		return x.ServerId
	}
	return 0
}

func (x *RoleShow) GetRoleName() string {
	if x != nil {
		return x.RoleName
	}
	return ""
}

func (x *RoleShow) GetLv() uint32 {
	if x != nil {
		return x.Lv
	}
	return 0
}

func (x *RoleShow) GetLastLoginTimer() uint32 {
	if x != nil {
		return x.LastLoginTimer
	}
	return 0
}

func (x *RoleShow) GetLastLogoutTimer() uint32 {
	if x != nil {
		return x.LastLogoutTimer
	}
	return 0
}

type RoleShowList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role *RoleShow `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *RoleShowList) Reset() {
	*x = RoleShowList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_table_Account_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoleShowList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleShowList) ProtoMessage() {}

func (x *RoleShowList) ProtoReflect() protoreflect.Message {
	mi := &file_table_Account_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoleShowList.ProtoReflect.Descriptor instead.
func (*RoleShowList) Descriptor() ([]byte, []int) {
	return file_table_Account_proto_rawDescGZIP(), []int{1}
}

func (x *RoleShowList) GetRole() *RoleShow {
	if x != nil {
		return x.Role
	}
	return nil
}

type Account struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId   uint64        `protobuf:"varint,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	AccountName string        `protobuf:"bytes,2,opt,name=account_name,json=accountName,proto3" json:"account_name,omitempty"`
	CreateTimer uint32        `protobuf:"varint,3,opt,name=create_timer,json=createTimer,proto3" json:"create_timer,omitempty"`
	LoginTimer  uint32        `protobuf:"varint,4,opt,name=login_timer,json=loginTimer,proto3" json:"login_timer,omitempty"`
	LogoutTimer uint32        `protobuf:"varint,5,opt,name=logout_timer,json=logoutTimer,proto3" json:"logout_timer,omitempty"`
	Phone       string        `protobuf:"bytes,6,opt,name=phone,proto3" json:"phone,omitempty"`
	RoleList    *RoleShowList `protobuf:"bytes,7,opt,name=role_list,json=roleList,proto3" json:"role_list,omitempty"`
}

func (x *Account) Reset() {
	*x = Account{}
	if protoimpl.UnsafeEnabled {
		mi := &file_table_Account_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Account) ProtoMessage() {}

func (x *Account) ProtoReflect() protoreflect.Message {
	mi := &file_table_Account_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Account.ProtoReflect.Descriptor instead.
func (*Account) Descriptor() ([]byte, []int) {
	return file_table_Account_proto_rawDescGZIP(), []int{2}
}

func (x *Account) GetAccountId() uint64 {
	if x != nil {
		return x.AccountId
	}
	return 0
}

func (x *Account) GetAccountName() string {
	if x != nil {
		return x.AccountName
	}
	return ""
}

func (x *Account) GetCreateTimer() uint32 {
	if x != nil {
		return x.CreateTimer
	}
	return 0
}

func (x *Account) GetLoginTimer() uint32 {
	if x != nil {
		return x.LoginTimer
	}
	return 0
}

func (x *Account) GetLogoutTimer() uint32 {
	if x != nil {
		return x.LogoutTimer
	}
	return 0
}

func (x *Account) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *Account) GetRoleList() *RoleShowList {
	if x != nil {
		return x.RoleList
	}
	return nil
}

var File_table_Account_proto protoreflect.FileDescriptor

var file_table_Account_proto_rawDesc = []byte{
	0x0a, 0x13, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2f, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22,
	0xc3, 0x01, 0x0a, 0x08, 0x52, 0x6f, 0x6c, 0x65, 0x53, 0x68, 0x6f, 0x77, 0x12, 0x17, 0x0a, 0x07,
	0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x72,
	0x6f, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x6f, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x6c, 0x76, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x6c, 0x76, 0x12,
	0x28, 0x0a, 0x10, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x11, 0x6c, 0x61, 0x73,
	0x74, 0x5f, 0x6c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x72, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x72, 0x22, 0x36, 0x0a, 0x0c, 0x52, 0x6f, 0x6c, 0x65, 0x53, 0x68, 0x6f,
	0x77, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x52,
	0x6f, 0x6c, 0x65, 0x53, 0x68, 0x6f, 0x77, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0xfd, 0x01,
	0x0a, 0x07, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x72, 0x12, 0x1f,
	0x0a, 0x0b, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x72, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0a, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x72, 0x12,
	0x21, 0x0a, 0x0c, 0x6c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x72, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x6c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x33, 0x0a, 0x09, 0x72, 0x6f, 0x6c, 0x65,
	0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x53, 0x68, 0x6f, 0x77, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x08, 0x72, 0x6f, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x08, 0x5a,
	0x06, 0x2f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_table_Account_proto_rawDescOnce sync.Once
	file_table_Account_proto_rawDescData = file_table_Account_proto_rawDesc
)

func file_table_Account_proto_rawDescGZIP() []byte {
	file_table_Account_proto_rawDescOnce.Do(func() {
		file_table_Account_proto_rawDescData = protoimpl.X.CompressGZIP(file_table_Account_proto_rawDescData)
	})
	return file_table_Account_proto_rawDescData
}

var file_table_Account_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_table_Account_proto_goTypes = []interface{}{
	(*RoleShow)(nil),     // 0: protobuf.RoleShow
	(*RoleShowList)(nil), // 1: protobuf.RoleShowList
	(*Account)(nil),      // 2: protobuf.Account
}
var file_table_Account_proto_depIdxs = []int32{
	0, // 0: protobuf.RoleShowList.role:type_name -> protobuf.RoleShow
	1, // 1: protobuf.Account.role_list:type_name -> protobuf.RoleShowList
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_table_Account_proto_init() }
func file_table_Account_proto_init() {
	if File_table_Account_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_table_Account_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoleShow); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_table_Account_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoleShowList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_table_Account_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Account); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_table_Account_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_table_Account_proto_goTypes,
		DependencyIndexes: file_table_Account_proto_depIdxs,
		MessageInfos:      file_table_Account_proto_msgTypes,
	}.Build()
	File_table_Account_proto = out.File
	file_table_Account_proto_rawDesc = nil
	file_table_Account_proto_goTypes = nil
	file_table_Account_proto_depIdxs = nil
}
