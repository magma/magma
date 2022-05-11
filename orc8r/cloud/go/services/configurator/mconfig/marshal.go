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

package mconfig

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"magma/orc8r/lib/go/protos"
)

func MarshalConfigs(configs map[string]proto.Message) (ConfigsByKey, error) {
	ret := ConfigsByKey{}
	for k, v := range configs {
		anyVal, err := ptypes.MarshalAny(v)
		if err != nil {
			return nil, err
		}
		bytesVal, err := protos.MarshalJSON(anyVal)
		if err != nil {
			return nil, err
		}
		ret[k] = bytesVal
	}
	return ret, nil
}

func UnmarshalConfigs(configs ConfigsByKey) (map[string]proto.Message, error) {
	ret := map[string]proto.Message{}
	for k, v := range configs {
		anyVal := &any.Any{}
		err := protos.Unmarshal(v, anyVal)
		if err != nil {
			return nil, fmt.Errorf("unmarshal mconfig from bytes to proto for key %s and bytes %v: %w", k, v, err)
		}
		msgVal, err := ptypes.Empty(anyVal)
		if err != nil {
			return nil, fmt.Errorf("create concrete proto.Message, for proto.Any %+v: %w", anyVal, err)
		}
		err = ptypes.UnmarshalAny(anyVal, msgVal)
		if err != nil {
			return nil, fmt.Errorf("unmarshal proto.Any into proto.Message, for proto.Any %+v: %w", anyVal, err)
		}
		ret[k] = msgVal
	}
	return ret, nil
}
