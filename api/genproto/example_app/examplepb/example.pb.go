// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.3
// source: examplepb/example.proto

package examplepb

import (
	reflect "reflect"
	sync "sync"

	simplepb "github.com/IguoChan/go-project/api/genproto/demo_app/simplepb"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ExReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	S *simplepb.SimpleRequest `protobuf:"bytes,1,opt,name=s,proto3" json:"s,omitempty"`
}

func (x *ExReq) Reset() {
	*x = ExReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_example_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExReq) ProtoMessage() {}

func (x *ExReq) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_example_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExReq.ProtoReflect.Descriptor instead.
func (*ExReq) Descriptor() ([]byte, []int) {
	return file_examplepb_example_proto_rawDescGZIP(), []int{0}
}

func (x *ExReq) GetS() *simplepb.SimpleRequest {
	if x != nil {
		return x.S
	}
	return nil
}

type ExResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code  int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ExResp) Reset() {
	*x = ExResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_example_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExResp) ProtoMessage() {}

func (x *ExResp) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_example_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExResp.ProtoReflect.Descriptor instead.
func (*ExResp) Descriptor() ([]byte, []int) {
	return file_examplepb_example_proto_rawDescGZIP(), []int{1}
}

func (x *ExResp) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ExResp) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_examplepb_example_proto protoreflect.FileDescriptor

var file_examplepb_example_proto_rawDesc = []byte{
	0x0a, 0x17, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2f, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x70, 0x62, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1e, 0x64, 0x65, 0x6d, 0x6f, 0x5f, 0x61, 0x70, 0x70, 0x2f, 0x73, 0x69, 0x6d,
	0x70, 0x6c, 0x65, 0x70, 0x62, 0x2f, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x2e, 0x0a, 0x05, 0x45, 0x78, 0x52, 0x65, 0x71, 0x12, 0x25, 0x0a, 0x01, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x70,
	0x62, 0x2e, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52,
	0x01, 0x73, 0x22, 0x32, 0x0a, 0x06, 0x45, 0x78, 0x52, 0x65, 0x73, 0x70, 0x12, 0x12, 0x0a, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x32, 0x4f, 0x0a, 0x07, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x12, 0x44, 0x0a, 0x01, 0x52, 0x12, 0x10, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x70, 0x62, 0x2e, 0x45, 0x78, 0x52, 0x65, 0x71, 0x1a, 0x11, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x70, 0x62, 0x2e, 0x45, 0x78, 0x52, 0x65, 0x73, 0x70, 0x22, 0x1a, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x14, 0x22, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x2f, 0x63, 0x79, 0x67, 0x3a, 0x01, 0x2a, 0x42, 0x4d, 0x5a, 0x4b, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x49, 0x67, 0x75, 0x6f, 0x43, 0x68, 0x61, 0x6e, 0x2f, 0x67,
	0x6f, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65,
	0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x5f, 0x61,
	0x70, 0x70, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x3b, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_examplepb_example_proto_rawDescOnce sync.Once
	file_examplepb_example_proto_rawDescData = file_examplepb_example_proto_rawDesc
)

func file_examplepb_example_proto_rawDescGZIP() []byte {
	file_examplepb_example_proto_rawDescOnce.Do(func() {
		file_examplepb_example_proto_rawDescData = protoimpl.X.CompressGZIP(file_examplepb_example_proto_rawDescData)
	})
	return file_examplepb_example_proto_rawDescData
}

var file_examplepb_example_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_examplepb_example_proto_goTypes = []interface{}{
	(*ExReq)(nil),                  // 0: examplepb.ExReq
	(*ExResp)(nil),                 // 1: examplepb.ExResp
	(*simplepb.SimpleRequest)(nil), // 2: simplepb.SimpleRequest
}
var file_examplepb_example_proto_depIdxs = []int32{
	2, // 0: examplepb.ExReq.s:type_name -> simplepb.SimpleRequest
	0, // 1: examplepb.Example.R:input_type -> examplepb.ExReq
	1, // 2: examplepb.Example.R:output_type -> examplepb.ExResp
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_examplepb_example_proto_init() }
func file_examplepb_example_proto_init() {
	if File_examplepb_example_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_examplepb_example_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExReq); i {
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
		file_examplepb_example_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExResp); i {
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
			RawDescriptor: file_examplepb_example_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_examplepb_example_proto_goTypes,
		DependencyIndexes: file_examplepb_example_proto_depIdxs,
		MessageInfos:      file_examplepb_example_proto_msgTypes,
	}.Build()
	File_examplepb_example_proto = out.File
	file_examplepb_example_proto_rawDesc = nil
	file_examplepb_example_proto_goTypes = nil
	file_examplepb_example_proto_depIdxs = nil
}
