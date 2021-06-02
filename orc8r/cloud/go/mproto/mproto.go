/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mproto

import (
	"magma/orc8r/cloud/go/mproto/protos"

	"github.com/golang/protobuf/proto"
)

// EncodeProtosDeterministic deterministically encodes a slice of protobuf messages,
// indexed by a unique identifier string, utilizing encodeProtoDeterministic
func EncodeProtosDeterministic(protoSlice map[string]proto.Message) ([]byte, error) {
	protoBytes := map[string][]byte{}
	for id, proto := range protoSlice {
		bytes, err := EncodeProtoDeterministic(proto)
		if err != nil {
			return nil, err
		}
		protoBytes[id] = bytes
	}

	return EncodeProtoDeterministic(&protos.ProtoSlice{ProtoBytes: protoBytes})
}

// EncodeProtoDeterministic encodes protobuf while enforcing deterministic serialization.
// NOTE: deterministic != canonical, so do not expect this encoding to be
// equal across languages or even versions of golang/protobuf/proto.
// For further reading, see below.
// 	- https://developers.google.com/protocol-buffers/docs/encoding#implications
//	- https://gist.github.com/kchristidis/39c8b310fd9da43d515c4394c3cd9510
func EncodeProtoDeterministic(pb proto.Message) ([]byte, error) {
	buf := &proto.Buffer{}
	buf.SetDeterministic(true)

	err := buf.Marshal(pb)
	return buf.Bytes(), err
}
