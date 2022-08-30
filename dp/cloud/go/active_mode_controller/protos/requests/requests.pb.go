//
//Copyright 2022 The Magma Authors.
//
//This source code is licensed under the BSD-style license found in the
//LICENSE file in the root directory of this source tree.
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.10.0
// source: dp/protos/requests.proto

package requests

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type RequestPayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload string `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *RequestPayload) Reset() {
	*x = RequestPayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dp_protos_requests_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestPayload) ProtoMessage() {}

func (x *RequestPayload) ProtoReflect() protoreflect.Message {
	mi := &file_dp_protos_requests_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestPayload.ProtoReflect.Descriptor instead.
func (*RequestPayload) Descriptor() ([]byte, []int) {
	return file_dp_protos_requests_proto_rawDescGZIP(), []int{0}
}

func (x *RequestPayload) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

type RequestDbIds struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids []int64 `protobuf:"varint,1,rep,packed,name=ids,proto3" json:"ids,omitempty"`
}

func (x *RequestDbIds) Reset() {
	*x = RequestDbIds{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dp_protos_requests_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestDbIds) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestDbIds) ProtoMessage() {}

func (x *RequestDbIds) ProtoReflect() protoreflect.Message {
	mi := &file_dp_protos_requests_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestDbIds.ProtoReflect.Descriptor instead.
func (*RequestDbIds) Descriptor() ([]byte, []int) {
	return file_dp_protos_requests_proto_rawDescGZIP(), []int{1}
}

func (x *RequestDbIds) GetIds() []int64 {
	if x != nil {
		return x.Ids
	}
	return nil
}

var File_dp_protos_requests_proto protoreflect.FileDescriptor

var file_dp_protos_requests_proto_rawDesc = []byte{
	0x0a, 0x18, 0x64, 0x70, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2a, 0x0a, 0x0e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x20, 0x0a, 0x0c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x44, 0x62, 0x49, 0x64, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x03, 0x52, 0x03, 0x69, 0x64, 0x73, 0x32, 0x43, 0x0a, 0x0f, 0x52, 0x61, 0x64, 0x69,
	0x6f, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x12, 0x30, 0x0a, 0x0e, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x12, 0x0f, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x1a, 0x0d,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x62, 0x49, 0x64, 0x73, 0x42, 0x3a, 0x5a,
	0x38, 0x6d, 0x61, 0x67, 0x6d, 0x61, 0x2f, 0x64, 0x70, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f,
	0x67, 0x6f, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x5f, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_dp_protos_requests_proto_rawDescOnce sync.Once
	file_dp_protos_requests_proto_rawDescData = file_dp_protos_requests_proto_rawDesc
)

func file_dp_protos_requests_proto_rawDescGZIP() []byte {
	file_dp_protos_requests_proto_rawDescOnce.Do(func() {
		file_dp_protos_requests_proto_rawDescData = protoimpl.X.CompressGZIP(file_dp_protos_requests_proto_rawDescData)
	})
	return file_dp_protos_requests_proto_rawDescData
}

var file_dp_protos_requests_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_dp_protos_requests_proto_goTypes = []interface{}{
	(*RequestPayload)(nil), // 0: RequestPayload
	(*RequestDbIds)(nil),   // 1: RequestDbIds
}
var file_dp_protos_requests_proto_depIdxs = []int32{
	0, // 0: RadioController.UploadRequests:input_type -> RequestPayload
	1, // 1: RadioController.UploadRequests:output_type -> RequestDbIds
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_dp_protos_requests_proto_init() }
func file_dp_protos_requests_proto_init() {
	if File_dp_protos_requests_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dp_protos_requests_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestPayload); i {
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
		file_dp_protos_requests_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestDbIds); i {
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
			RawDescriptor: file_dp_protos_requests_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_dp_protos_requests_proto_goTypes,
		DependencyIndexes: file_dp_protos_requests_proto_depIdxs,
		MessageInfos:      file_dp_protos_requests_proto_msgTypes,
	}.Build()
	File_dp_protos_requests_proto = out.File
	file_dp_protos_requests_proto_rawDesc = nil
	file_dp_protos_requests_proto_goTypes = nil
	file_dp_protos_requests_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RadioControllerClient is the client API for RadioController service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RadioControllerClient interface {
	UploadRequests(ctx context.Context, in *RequestPayload, opts ...grpc.CallOption) (*RequestDbIds, error)
}

type radioControllerClient struct {
	cc grpc.ClientConnInterface
}

func NewRadioControllerClient(cc grpc.ClientConnInterface) RadioControllerClient {
	return &radioControllerClient{cc}
}

func (c *radioControllerClient) UploadRequests(ctx context.Context, in *RequestPayload, opts ...grpc.CallOption) (*RequestDbIds, error) {
	out := new(RequestDbIds)
	err := c.cc.Invoke(ctx, "/RadioController/UploadRequests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RadioControllerServer is the server API for RadioController service.
type RadioControllerServer interface {
	UploadRequests(context.Context, *RequestPayload) (*RequestDbIds, error)
}

// UnimplementedRadioControllerServer can be embedded to have forward compatible implementations.
type UnimplementedRadioControllerServer struct {
}

func (*UnimplementedRadioControllerServer) UploadRequests(context.Context, *RequestPayload) (*RequestDbIds, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadRequests not implemented")
}

func RegisterRadioControllerServer(s *grpc.Server, srv RadioControllerServer) {
	s.RegisterService(&_RadioController_serviceDesc, srv)
}

func _RadioController_UploadRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestPayload)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RadioControllerServer).UploadRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RadioController/UploadRequests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RadioControllerServer).UploadRequests(ctx, req.(*RequestPayload))
	}
	return interceptor(ctx, in, info, handler)
}

var _RadioController_serviceDesc = grpc.ServiceDesc{
	ServiceName: "RadioController",
	HandlerType: (*RadioControllerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UploadRequests",
			Handler:    _RadioController_UploadRequests_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dp/protos/requests.proto",
}