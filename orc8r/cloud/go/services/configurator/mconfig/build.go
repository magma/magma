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
	"time"

	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pkg/errors"
)

// TODO(T71525030): restore CreateMconfig and remove CreateMconfigJSON

//func CreateMconfig(network *storage.Network, graph *storage.EntityGraph, gatewayID string) (*protos.GatewayConfigs, error) {
//	builders, err := GetBuilders()
//	if err != nil {
//		return nil, err
//	}
//
//	configs := ConfigsByKey{}
//	for _, b := range builders {
//		partialConfig, err := b.Build(network, graph, gatewayID)
//		if err != nil {
//			return nil, errors.Wrapf(err, "mconfig builder %+v error", b)
//		}
//		for key, config := range partialConfig {
//			_, ok := configs[key]
//			if ok {
//				return nil, fmt.Errorf("received partial config for key %v from multiple mconfig builders", key)
//			}
//			configs[key] = config
//		}
//	}
//
//	mconfig := &protos.GatewayConfigs{
//		Metadata: &protos.GatewayConfigsMetadata{
//			CreatedAt: uint64(time.Now().Unix()),
//			Digest:    &protos.GatewayConfigsDigest{},
//		},
//		ConfigsByKey: configs,
//	}
//	mconfig.Metadata.Digest.Md5HexDigest, err = mconfig.GetMconfigDigest()
//	if err != nil {
//		return nil, errors.Wrap(err, "generate mconfig digest")
//	}
//
//	return mconfig, nil
//}

func CreateMconfigJSON(network *storage.Network, graph *storage.EntityGraph, gatewayID string) (*protos.GatewayConfigs, error) {
	builders, err := GetBuilders()
	if err != nil {
		return nil, err
	}

	configs := map[string]*any.Any{}
	for _, b := range builders {
		partialConfig, err := b.Build(network, graph, gatewayID)
		if err != nil {
			return nil, errors.Wrapf(err, "mconfig builder %+v error", b)
		}
		for key, config := range partialConfig {
			_, ok := configs[key]
			if ok {
				return nil, fmt.Errorf("received partial config for key %v from multiple mconfig builders", key)
			}
			configBytes := &wrappers.BytesValue{Value: config}
			configAny, err := ptypes.MarshalAny(configBytes)
			if err != nil {
				return nil, err
			}
			configs[key] = configAny
		}
	}

	mconfig := &protos.GatewayConfigs{
		Metadata:     &protos.GatewayConfigsMetadata{CreatedAt: uint64(time.Now().Unix())},
		ConfigsByKey: configs,
	}

	return mconfig, nil
}
