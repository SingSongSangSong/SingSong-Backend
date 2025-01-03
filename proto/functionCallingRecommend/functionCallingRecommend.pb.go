// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v5.28.0
// source: proto/functionCallingRecommend/functionCallingRecommend.proto

package functionCallingRecommend

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

// 요청 메시지: memberId를 받아서 유사한 아이템을 검색
type FunctionCallingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MemberId int64  `protobuf:"varint,1,opt,name=memberId,proto3" json:"memberId,omitempty"`
	Command  string `protobuf:"bytes,2,opt,name=command,proto3" json:"command,omitempty"`
}

func (x *FunctionCallingRequest) Reset() {
	*x = FunctionCallingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_functionCallingRecommend_functionCallingRecommend_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FunctionCallingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FunctionCallingRequest) ProtoMessage() {}

func (x *FunctionCallingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_functionCallingRecommend_functionCallingRecommend_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FunctionCallingRequest.ProtoReflect.Descriptor instead.
func (*FunctionCallingRequest) Descriptor() ([]byte, []int) {
	return file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDescGZIP(), []int{0}
}

func (x *FunctionCallingRequest) GetMemberId() int64 {
	if x != nil {
		return x.MemberId
	}
	return 0
}

func (x *FunctionCallingRequest) GetCommand() string {
	if x != nil {
		return x.Command
	}
	return ""
}

// 유사한 아이템들의 리스트를 반환하는 응답 메시지
type FunctionCallingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SongInfoId []int64 `protobuf:"varint,1,rep,packed,name=songInfoId,proto3" json:"songInfoId,omitempty"` // 유사한 아이템 목록
}

func (x *FunctionCallingResponse) Reset() {
	*x = FunctionCallingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_functionCallingRecommend_functionCallingRecommend_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FunctionCallingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FunctionCallingResponse) ProtoMessage() {}

func (x *FunctionCallingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_functionCallingRecommend_functionCallingRecommend_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FunctionCallingResponse.ProtoReflect.Descriptor instead.
func (*FunctionCallingResponse) Descriptor() ([]byte, []int) {
	return file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDescGZIP(), []int{1}
}

func (x *FunctionCallingResponse) GetSongInfoId() []int64 {
	if x != nil {
		return x.SongInfoId
	}
	return nil
}

var File_proto_functionCallingRecommend_functionCallingRecommend_proto protoreflect.FileDescriptor

var file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDesc = []byte{
	0x0a, 0x3d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x43, 0x61, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64,
	0x2f, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x61, 0x6c, 0x6c, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x18, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x61, 0x6c, 0x6c, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x22, 0x4e, 0x0a, 0x16, 0x46, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x61, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0x39, 0x0a, 0x17, 0x46, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x61, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x6f, 0x6e, 0x67, 0x49, 0x6e, 0x66, 0x6f,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x0a, 0x73, 0x6f, 0x6e, 0x67, 0x49, 0x6e,
	0x66, 0x6f, 0x49, 0x64, 0x32, 0xa6, 0x01, 0x0a, 0x18, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x43, 0x61, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x64, 0x12, 0x89, 0x01, 0x0a, 0x20, 0x47, 0x65, 0x74, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x43, 0x61, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x30, 0x2e, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x43, 0x61, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x64, 0x2e, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x61, 0x6c, 0x6c, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x66, 0x75, 0x6e, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x43, 0x61, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x64, 0x2e, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x61, 0x6c, 0x6c,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x20, 0x5a,
	0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43,
	0x61, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDescOnce sync.Once
	file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDescData = file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDesc
)

func file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDescGZIP() []byte {
	file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDescOnce.Do(func() {
		file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDescData)
	})
	return file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDescData
}

var file_proto_functionCallingRecommend_functionCallingRecommend_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_functionCallingRecommend_functionCallingRecommend_proto_goTypes = []interface{}{
	(*FunctionCallingRequest)(nil),  // 0: functionCallingRecommend.FunctionCallingRequest
	(*FunctionCallingResponse)(nil), // 1: functionCallingRecommend.FunctionCallingResponse
}
var file_proto_functionCallingRecommend_functionCallingRecommend_proto_depIdxs = []int32{
	0, // 0: functionCallingRecommend.functionCallingRecommend.GetFunctionCallingRecommendation:input_type -> functionCallingRecommend.FunctionCallingRequest
	1, // 1: functionCallingRecommend.functionCallingRecommend.GetFunctionCallingRecommendation:output_type -> functionCallingRecommend.FunctionCallingResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_functionCallingRecommend_functionCallingRecommend_proto_init() }
func file_proto_functionCallingRecommend_functionCallingRecommend_proto_init() {
	if File_proto_functionCallingRecommend_functionCallingRecommend_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_functionCallingRecommend_functionCallingRecommend_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FunctionCallingRequest); i {
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
		file_proto_functionCallingRecommend_functionCallingRecommend_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FunctionCallingResponse); i {
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
			RawDescriptor: file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_functionCallingRecommend_functionCallingRecommend_proto_goTypes,
		DependencyIndexes: file_proto_functionCallingRecommend_functionCallingRecommend_proto_depIdxs,
		MessageInfos:      file_proto_functionCallingRecommend_functionCallingRecommend_proto_msgTypes,
	}.Build()
	File_proto_functionCallingRecommend_functionCallingRecommend_proto = out.File
	file_proto_functionCallingRecommend_functionCallingRecommend_proto_rawDesc = nil
	file_proto_functionCallingRecommend_functionCallingRecommend_proto_goTypes = nil
	file_proto_functionCallingRecommend_functionCallingRecommend_proto_depIdxs = nil
}
