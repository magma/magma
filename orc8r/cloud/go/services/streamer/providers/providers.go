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

package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/streamer"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pkg/errors"
)

// MconfigProvider provides streamer mconfigs (magma configs).
type MconfigProvider struct{}

func (p *MconfigProvider) GetStreamName() string {
	return definitions.MconfigStreamName
}

func (p *MconfigProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	res, err := configurator.GetMconfigFor(gatewayId)
	if err != nil {
		return nil, errors.Wrap(err, "get mconfig from configurator")
	}

	// TODO(8/5/20): revert below once we send proto descriptors from mconfig_builders

	//if extraArgs != nil {
	//	// Currently, only use of extraArgs is mconfig hash
	//	receivedDigest := &protos.GatewayConfigsDigest{}
	//	if err := ptypes.UnmarshalAny(extraArgs, receivedDigest); err == nil {
	//		glog.V(2).Infof(
	//			"Received, generated config digests: %v, %v\n",
	//			receivedDigest,
	//			res.Configs.Metadata.Digest.Md5HexDigest,
	//		)
	//		return mconfigToUpdate(res.Configs, res.LogicalID, receivedDigest.Md5HexDigest)
	//	}
	//}
	//
	//return mconfigToUpdate(res.Configs, res.LogicalID, "")

	return bytesMconfigToUpdate(res.Configs, res.LogicalID)
}

func mconfigToUpdate(configs *protos.GatewayConfigs, logicalID string, digest string) ([]*protos.DataUpdate, error) {
	// Early/empty return if gateway already has this config
	if digest == configs.Metadata.Digest.Md5HexDigest {
		return []*protos.DataUpdate{}, streamer.EAGAIN // do not close the stream, there were no changes in configs
	}
	marshaledConfig, err := protos.MarshalJSON(configs)
	if err != nil {
		return nil, errors.Wrap(err, "marshal gateway mconfig")
	}
	return []*protos.DataUpdate{{Key: logicalID, Value: marshaledConfig}}, nil
}

// TODO(8/5/20): revert below once we send proto descriptors from mconfig_builders

const gatewayConfigsTemplate = `
{
  "configsByKey": {
	  %s
  },
  "metadata": {
    "createdAt": "%v"
  }
}`

// bytesMconfigToUpdate creates a GatewayConfigs data update by manually
// constructing the JSON-marshaled bytes for the data update.
// This hack is a temporary workaround until we upgrade to Go's protobuf
// library APIv2, which has built-in support for dynamic any.Any resolution,
// which is required to marshal an any.Any proto to JSON.
func bytesMconfigToUpdate(configs *protos.GatewayConfigs, logicalID string) ([]*protos.DataUpdate, error) {
	configsByKey := map[string][]byte{}

	for k, v := range configs.ConfigsByKey {
		bytesVal := &wrappers.BytesValue{}
		err := ptypes.UnmarshalAny(v, bytesVal)
		if err != nil {
			return nil, err
		}
		configsByKey[k] = bytesVal.Value
	}
	jsonConfig, err := getJSONConfig(configsByKey)
	if err != nil {
		return nil, err
	}

	var marshaledConfig json.RawMessage = []byte(fmt.Sprintf(gatewayConfigsTemplate, jsonConfig, uint64(time.Now().Unix())))
	prettyConfig, err := json.MarshalIndent(marshaledConfig, "", "  ")
	if err != nil {
		return nil, err
	}

	return []*protos.DataUpdate{{Key: logicalID, Value: prettyConfig}}, nil
}

func getJSONConfig(jsonByKey map[string][]byte) ([]byte, error) {
	var configs [][]byte
	for k, v := range jsonByKey {
		configs = append(configs, []byte(fmt.Sprintf(`"%s": %s`, k, v)))
	}
	joined := bytes.Join(configs, []byte(",\n"))
	return joined, nil
}
