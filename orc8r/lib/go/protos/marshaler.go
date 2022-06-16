/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package protos is protoc generated GRPC package
package protos

import (
	"bytes"
	"strings"

	"github.com/golang/glog"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func Marshal(msg proto.Message) ([]byte, error) {
	buff, err := (&jsonpb.MarshalOptions{
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
		UseProtoNames:   true}).Marshal(msg)

	return buff, err
}

func MarshalIntern(msg proto.Message) ([]byte, error) {
	buff, err := (&jsonpb.MarshalOptions{EmitUnpopulated: true, Indent: " "}).Marshal(msg)
	return buff, err
}

func MarshalJSON(msg proto.Message) ([]byte, error) {
	buff, err := (&jsonpb.MarshalOptions{Indent: " "}).Marshal(msg)
	return buff, err
}

func Unmarshal(bt []byte, msg proto.Message) error {
	return (&jsonpb.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(bt, msg)
}

func TestMarshal(msg proto.Message) string {
	res, _ := Marshal(msg)
	return string(res)
}

type mconfigAnyResolver struct{}

// Resolve - AnyResolver interface implementation, it'll resolve any unregistered Any types to Void instead of
// returning error
func (mconfigAnyResolver) Resolve(typeUrl string) (proto.Message, error) {
	// Only the part of typeUrl after the last slash is relevant.
	mname := typeUrl
	if slash := strings.LastIndex(mname, "/"); slash >= 0 {
		mname = mname[slash+1:]
	}
	mt, _ := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(mname))
	if mt == nil {
		glog.V(4).Infof("mconfigAnyResolver: unknown message type %q", mname)
		return new(Bytes), nil
	} else {
		return mt.New().Interface().(proto.Message), nil
	}
}

// MarshalMconfig is a special mconfig marshaler tolerant to unregistered Any types
func MarshalMconfig(msg proto.Message) ([]byte, error) {
	buff, err := marshalMconfigs(msg)
	return buff.Bytes(), err
}

// MarshalMconfigToString - same as MarshalMconfig but returns string
func MarshalMconfigToString(msg proto.Message) (string, error) {
	buff, err := marshalMconfigs(msg)
	return buff.String(), err
}

func marshalMconfigs(msg proto.Message) (*bytes.Buffer, error) {
	var buff bytes.Buffer
	b, err := (&jsonpb.MarshalOptions{EmitUnpopulated: true, Indent: " "}).Marshal(msg)
	buff.Write(b)
	return &buff, err
}

// UnmarshalMconfig is a special mconfig Unmarshaler tolerant to unregistered Any types
func UnmarshalMconfig(bt []byte, msg proto.Message) error {
	return (&jsonpb.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(bt, msg)
}

// MarshalJSONPB implements JSONPBMarshaler interface for Bytes type
func (bm *Bytes) MarshalJSONPB(_ *jsonpb.MarshalOptions) ([]byte, error) {
	if bm != nil {
		var b = make([]byte, len(bm.Val))
		copy(b, bm.Val)
		return b, nil
	}
	return []byte{}, nil
}

// UnmarshalJSONPB implements JSONPBUnmarshaler interface for Bytes type
func (bm *Bytes) UnmarshalJSONPB(_ *jsonpb.UnmarshalOptions, b []byte) error {
	if bm != nil {
		bm.Reset()
		bm.Val = make([]byte, len(b))
		copy(bm.Val, b)
	}
	return nil
}
