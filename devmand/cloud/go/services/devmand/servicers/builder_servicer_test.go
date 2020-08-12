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

package servicers_test

import (
	"testing"

	"magma/devmand/cloud/go/devmand"
	devmand_plugin "magma/devmand/cloud/go/plugin"
	devmand_mconfig "magma/devmand/cloud/go/protos/mconfig"
	devmand_service "magma/devmand/cloud/go/services/devmand"
	"magma/devmand/cloud/go/services/devmand/obsidian/models"
	devmand_test_init "magma/devmand/cloud/go/services/devmand/test_init"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &devmand_plugin.DevmandOrchestratorPlugin{}))
	devmand_test_init.StartTestService(t)

	nw := configurator.Network{ID: "n1"}
	device := configurator.NetworkEntity{
		Type: devmand.SymphonyDeviceType, Key: "d1",
		Config:             models.NewDefaultSymphonyDeviceConfig(),
		ParentAssociations: []storage.TypeAndKey{{Type: devmand.SymphonyAgentType, Key: "a1"}},
	}
	agent := configurator.NetworkEntity{
		Type: devmand.SymphonyAgentType, Key: "a1",
		Associations: []storage.TypeAndKey{
			{Type: devmand.SymphonyDeviceType, Key: "d1"},
		},
		ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "a1"}},
	}
	gateway := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "a1",
		Associations: []storage.TypeAndKey{{Type: devmand.SymphonyAgentType, Key: "a1"}},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{gateway, agent, device},
		Edges: []configurator.GraphEdge{
			{From: gateway.GetTypeAndKey(), To: agent.GetTypeAndKey()},
			{From: agent.GetTypeAndKey(), To: device.GetTypeAndKey()},
		},
	}

	expected := map[string]proto.Message{
		"devmand": &devmand_mconfig.DevmandGatewayConfig{
			ManagedDevices: map[string]*devmand_mconfig.ManagedDevice{
				"d1": {
					DeviceConfig: "{}",
					DeviceType:   []string{"device_type 1", "device_type 2"},
					Channels: &devmand_mconfig.Channels{
						SnmpChannel: &devmand_mconfig.SNMPChannel{
							Community: "snmp community",
							Version:   "1",
						},
					},
					Host:     "device_host",
					Platform: "device_platform",
				},
			},
		},
	}

	actual, err := build(&nw, &graph, "a1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func build(network *configurator.Network, graph *configurator.EntityGraph, gatewayID string) (map[string]proto.Message, error) {
	networkProto, err := network.ToStorageProto()
	if err != nil {
		return nil, err
	}
	graphProto, err := graph.ToStorageProto()
	if err != nil {
		return nil, err
	}

	builder := mconfig.NewRemoteBuilder(devmand_service.ServiceName)
	res, err := builder.Build(networkProto, graphProto, gatewayID)
	if err != nil {
		return nil, err
	}

	configs, err := mconfig.UnmarshalConfigs(res)
	if err != nil {
		return nil, err
	}

	return configs, nil
}
