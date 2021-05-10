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
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
)

func MarshalConfigs(configs map[string]proto.Message) (ConfigsByKey, error) {
	ret := ConfigsByKey{}
	for k, v := range configs {
		anyVal, err := ptypes.MarshalAny(v)
		if err != nil {
			return nil, errors.WithStack(err)
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
			return nil, errors.Wrapf(err, "unmarshal mconfig from bytes to proto for key %s and bytes %v", k, v)
		}
		msgVal, err := ptypes.Empty(anyVal)
		if err != nil {
			return nil, errors.Wrapf(err, "create concrete proto.Message, for proto.Any %+v", anyVal)
		}
		err = ptypes.UnmarshalAny(anyVal, msgVal)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal proto.Any into proto.Message, for proto.Any %+v", anyVal)
		}
		ret[k] = msgVal
	}
	return ret, nil
}
