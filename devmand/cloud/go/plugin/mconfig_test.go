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

package plugin_test

import (
	"testing"

	"magma/devmand/cloud/go/devmand"
	"magma/devmand/cloud/go/plugin"
	"magma/devmand/cloud/go/protos/mconfig"
	models2 "magma/devmand/cloud/go/services/devmand/obsidian/models"
	"magma/orc8r/cloud/go/orc8r"
	orc8rplugin "magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	orc8rplugin.RegisterPluginForTests(t, &plugin.DevmandOrchestratorPlugin{})

	nID := "n1"
	nw := configurator.Network{ID: nID}
	device := configurator.NetworkEntity{
		Type: devmand.SymphonyDeviceType, Key: "d1",
		Config:             models2.NewDefaultSymphonyDeviceConfig(),
		ParentAssociations: []storage.TypeAndKey{storage.TypeAndKey{Type: devmand.SymphonyAgentType, Key: "a1"}},
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

	actual := map[string]proto.Message{}
	builder := plugin.Builder{}
	err := builder.Build("n1", "a1", graph, nw, actual)
	assert.NoError(t, err)

	expected := map[string]proto.Message{
		"devmand": &mconfig.DevmandGatewayConfig{
			ManagedDevices: map[string]*mconfig.ManagedDevice{
				"d1": {
					DeviceConfig: "{}",
					DeviceType:   []string{"device_type 1", "device_type 2"},
					Channels: &mconfig.Channels{
						SnmpChannel: &mconfig.SNMPChannel{
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
	assert.Equal(t, expected, actual)
}
