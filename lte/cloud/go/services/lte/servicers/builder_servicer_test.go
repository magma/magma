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

	"magma/lte/cloud/go/lte"
	lte_plugin "magma/lte/cloud/go/plugin"
	lte_mconfig "magma/lte/cloud/go/protos/mconfig"
	"magma/lte/cloud/go/serdes"
	lte_service "magma/lte/cloud/go/services/lte"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{}))
	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{}))
	lte_test_init.StartTestService(t)

	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: lte_models.NewDefaultTDDNetworkConfig(),
			orc8r.DnsdNetworkType: &models.NetworkDNSConfig{
				EnableCaching: swag.Bool(true),
			},
		},
	}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "gw1",
		Associations: []storage.TypeAndKey{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config: newDefaultGatewayConfig(),
		Associations: []storage.TypeAndKey{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
		},
		ParentAssociations: []storage.TypeAndKey{gw.GetTypeAndKey()},
	}
	enb := configurator.NetworkEntity{
		Type: lte.CellularEnodebEntityType, Key: "enb1",
		Config:             newDefaultEnodebConfig(),
		ParentAssociations: []storage.TypeAndKey{lteGW.GetTypeAndKey()},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{enb, lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTypeAndKey(), To: lteGW.GetTypeAndKey()},
			{From: lteGW.GetTypeAndKey(), To: enb.GetTypeAndKey()},
		},
	}

	expected := map[string]proto.Message{
		"enodebd": &lte_mconfig.EnodebD{
			LogLevel: protos.LogLevel_INFO,
			Pci:      260,
			TddConfig: &lte_mconfig.EnodebD_TDDConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
			BandwidthMhz:        20,
			AllowEnodebTransmit: true,
			Tac:                 1,
			PlmnidList:          "00101",
			CsfbRat:             lte_mconfig.EnodebD_CSFBRAT_2G,
			Arfcn_2G:            nil,
			EnbConfigsBySerial: map[string]*lte_mconfig.EnodebD_EnodebConfig{
				"enb1": {
					Earfcndl:               39150,
					SubframeAssignment:     2,
					SpecialSubframePattern: 7,
					Pci:                    260,
					TransmitEnabled:        true,
					DeviceClass:            "Baicells ID TDD/FDD",
					BandwidthMhz:           20,
					Tac:                    15000,
					CellId:                 138777000,
				},
			},
		},
		"mobilityd": &lte_mconfig.MobilityD{
			LogLevel: protos.LogLevel_INFO,
			IpBlock:  "192.168.128.0/24",
		},
		"mme": &lte_mconfig.MME{
			LogLevel:                 protos.LogLevel_INFO,
			Mcc:                      "001",
			Mnc:                      "01",
			Tac:                      1,
			MmeCode:                  1,
			MmeGid:                   1,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			EnableDnsCaching:         false,
			AttachedEnodebTacs:       []int32{15000},
			NatEnabled:               true,
		},
		"pipelined": &lte_mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24",
			NatEnabled:    true,
			DefaultRuleId: "",
			Services: []lte_mconfig.PipelineD_NetworkServices{
				lte_mconfig.PipelineD_ENFORCEMENT,
			},
			SgiManagementIfaceVlan: "",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:        protos.LogLevel_INFO,
			LteAuthOp:       []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:      []byte("\x80\x00"),
			SubProfiles:     nil,
			HssRelayEnabled: false,
		},
		"policydb": &lte_mconfig.PolicyDB{
			LogLevel: protos.LogLevel_INFO,
		},
		"sessiond": &lte_mconfig.SessionD{
			LogLevel:         protos.LogLevel_INFO,
			GxGyRelayEnabled: false,
			WalletExhaustDetection: &lte_mconfig.WalletExhaustDetection{
				TerminateOnExhaust: false,
			},
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
	}

	// Happy path
	actual, err := build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// Break with non-allowed network service
	setEPCNetworkServices([]string{"0xdeadbeef"}, &nw)
	actual, err = build(&nw, &graph, "gw1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown network service name 0xdeadbeef")

	// Don't break with deprecated network services
	setEPCNetworkServices([]string{"metering"}, &nw)
	expected["pipelined"] = &lte_mconfig.PipelineD{
		LogLevel:      protos.LogLevel_INFO,
		UeIpBlock:     "192.168.128.0/24",
		NatEnabled:    true,
		DefaultRuleId: "",
		Services: []lte_mconfig.PipelineD_NetworkServices{
			lte_mconfig.PipelineD_METERING,
		},
	}
	actual, err = build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestBuilder_Build_NonNat(t *testing.T) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{}))
	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{}))
	lte_test_init.StartTestService(t)

	// No dnsd config, no enodebs
	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: lte_models.NewDefaultTDDNetworkConfig(),
		},
	}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "gw1",
		Associations: []storage.TypeAndKey{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             newGatewayConfigNonNat("", "", ""),
		ParentAssociations: []storage.TypeAndKey{gw.GetTypeAndKey()},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTypeAndKey(), To: lteGW.GetTypeAndKey()},
		},
	}

	expected := map[string]proto.Message{
		"enodebd": &lte_mconfig.EnodebD{
			LogLevel: protos.LogLevel_INFO,
			Pci:      260,
			TddConfig: &lte_mconfig.EnodebD_TDDConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
			BandwidthMhz:        20,
			AllowEnodebTransmit: true,
			Tac:                 1,
			PlmnidList:          "00101",
			CsfbRat:             lte_mconfig.EnodebD_CSFBRAT_2G,
			Arfcn_2G:            nil,
			EnbConfigsBySerial:  nil,
		},
		"mobilityd": &lte_mconfig.MobilityD{
			LogLevel:        protos.LogLevel_INFO,
			IpBlock:         "192.168.128.0/24",
			IpAllocatorType: lte_mconfig.MobilityD_IP_POOL,
			StaticIpEnabled: false,
			MultiApnIpAlloc: false,
		},
		"mme": &lte_mconfig.MME{
			LogLevel:                 protos.LogLevel_INFO,
			Mcc:                      "001",
			Mnc:                      "01",
			Tac:                      1,
			MmeCode:                  1,
			MmeGid:                   1,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			AttachedEnodebTacs:       nil,
			NatEnabled:               false,
		},
		"pipelined": &lte_mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24",
			NatEnabled:    false,
			DefaultRuleId: "",
			Services: []lte_mconfig.PipelineD_NetworkServices{
				lte_mconfig.PipelineD_ENFORCEMENT,
			},
			SgiManagementIfaceVlan: "",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:        protos.LogLevel_INFO,
			LteAuthOp:       []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:      []byte("\x80\x00"),
			SubProfiles:     nil,
			HssRelayEnabled: false,
		},
		"policydb": &lte_mconfig.PolicyDB{
			LogLevel: protos.LogLevel_INFO,
		},
		"sessiond": &lte_mconfig.SessionD{
			LogLevel:         protos.LogLevel_INFO,
			GxGyRelayEnabled: false,
			WalletExhaustDetection: &lte_mconfig.WalletExhaustDetection{
				TerminateOnExhaust: false,
			},
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
	}
	actual, err := build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	setEPCNetworkIPAllocator(&nw, lte_models.DHCPBroadcastAllocationMode, false, false)
	expected["mobilityd"] = &lte_mconfig.MobilityD{
		LogLevel:        protos.LogLevel_INFO,
		IpBlock:         "192.168.128.0/24",
		IpAllocatorType: lte_mconfig.MobilityD_DHCP,
		StaticIpEnabled: false,
	}
	actual, err = build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	setEPCNetworkIPAllocator(&nw, lte_models.NATAllocationMode, false, false)
	expected["mobilityd"] = &lte_mconfig.MobilityD{
		LogLevel:        protos.LogLevel_INFO,
		IpBlock:         "192.168.128.0/24",
		IpAllocatorType: lte_mconfig.MobilityD_IP_POOL,
		StaticIpEnabled: false,
	}
	actual, err = build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	setEPCNetworkIPAllocator(&nw, lte_models.NATAllocationMode, true, false)
	expected["mobilityd"] = &lte_mconfig.MobilityD{
		LogLevel:        protos.LogLevel_INFO,
		IpBlock:         "192.168.128.0/24",
		IpAllocatorType: lte_mconfig.MobilityD_IP_POOL,
		StaticIpEnabled: true,
	}
	actual, err = build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	setEPCNetworkIPAllocator(&nw, lte_models.DHCPBroadcastAllocationMode, true, false)
	expected["mobilityd"] = &lte_mconfig.MobilityD{
		LogLevel:        protos.LogLevel_INFO,
		IpBlock:         "192.168.128.0/24",
		IpAllocatorType: lte_mconfig.MobilityD_DHCP,
		StaticIpEnabled: true,
	}
	actual, err = build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// validate multi apn mconfig
	setEPCNetworkIPAllocator(&nw, lte_models.DHCPBroadcastAllocationMode, true, true)
	expected["mobilityd"] = &lte_mconfig.MobilityD{
		LogLevel:        protos.LogLevel_INFO,
		IpBlock:         "192.168.128.0/24",
		IpAllocatorType: lte_mconfig.MobilityD_DHCP,
		StaticIpEnabled: true,
		MultiApnIpAlloc: true,
	}
	actual, err = build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// validate SGi vlan tag mconfig
	lteGW = configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             newGatewayConfigNonNat("30", "", ""),
		ParentAssociations: []storage.TypeAndKey{gw.GetTypeAndKey()},
	}
	graph = configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTypeAndKey(), To: lteGW.GetTypeAndKey()},
		},
	}
	expected["pipelined"] = &lte_mconfig.PipelineD{
		LogLevel:      protos.LogLevel_INFO,
		UeIpBlock:     "192.168.128.0/24",
		NatEnabled:    false,
		DefaultRuleId: "",
		Services: []lte_mconfig.PipelineD_NetworkServices{
			lte_mconfig.PipelineD_ENFORCEMENT,
		},
		SgiManagementIfaceVlan: "30",
	}

	actual, err = build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// validate SGi ip address
	// validate SGi vlan tag mconfig
	lteGW = configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             newGatewayConfigNonNat("44", "1.2.3.4", ""),
		ParentAssociations: []storage.TypeAndKey{gw.GetTypeAndKey()},
	}
	graph = configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTypeAndKey(), To: lteGW.GetTypeAndKey()},
		},
	}
	expected["pipelined"] = &lte_mconfig.PipelineD{
		LogLevel:      protos.LogLevel_INFO,
		UeIpBlock:     "192.168.128.0/24",
		NatEnabled:    false,
		DefaultRuleId: "",
		Services: []lte_mconfig.PipelineD_NetworkServices{
			lte_mconfig.PipelineD_ENFORCEMENT,
		},
		SgiManagementIfaceVlan:   "44",
		SgiManagementIfaceIpAddr: "1.2.3.4",
	}

	actual, err = build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// validate SGi ip address and gateway
	// validate SGi vlan tag mconfig
	lteGW = configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             newGatewayConfigNonNat("55", "1.2.3.4/24", "1.2.3.1"),
		ParentAssociations: []storage.TypeAndKey{gw.GetTypeAndKey()},
	}
	graph = configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTypeAndKey(), To: lteGW.GetTypeAndKey()},
		},
	}
	expected["pipelined"] = &lte_mconfig.PipelineD{
		LogLevel:      protos.LogLevel_INFO,
		UeIpBlock:     "192.168.128.0/24",
		NatEnabled:    false,
		DefaultRuleId: "",
		Services: []lte_mconfig.PipelineD_NetworkServices{
			lte_mconfig.PipelineD_ENFORCEMENT,
		},
		SgiManagementIfaceVlan:   "55",
		SgiManagementIfaceIpAddr: "1.2.3.4/24",
		SgiManagementIfaceGw:     "1.2.3.1",
	}

	actual, err = build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

}

func TestBuilder_Build_BaseCase(t *testing.T) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{}))
	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{}))
	lte_test_init.StartTestService(t)

	// No dnsd config, no enodebs
	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: lte_models.NewDefaultTDDNetworkConfig(),
		},
	}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "gw1",
		Associations: []storage.TypeAndKey{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             newDefaultGatewayConfig(),
		ParentAssociations: []storage.TypeAndKey{gw.GetTypeAndKey()},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTypeAndKey(), To: lteGW.GetTypeAndKey()},
		},
	}

	expected := map[string]proto.Message{
		"enodebd": &lte_mconfig.EnodebD{
			LogLevel: protos.LogLevel_INFO,
			Pci:      260,
			TddConfig: &lte_mconfig.EnodebD_TDDConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
			BandwidthMhz:        20,
			AllowEnodebTransmit: true,
			Tac:                 1,
			PlmnidList:          "00101",
			CsfbRat:             lte_mconfig.EnodebD_CSFBRAT_2G,
			Arfcn_2G:            nil,
			EnbConfigsBySerial:  nil,
		},
		"mobilityd": &lte_mconfig.MobilityD{
			LogLevel: protos.LogLevel_INFO,
			IpBlock:  "192.168.128.0/24",
		},
		"mme": &lte_mconfig.MME{
			LogLevel:                 protos.LogLevel_INFO,
			Mcc:                      "001",
			Mnc:                      "01",
			Tac:                      1,
			MmeCode:                  1,
			MmeGid:                   1,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			AttachedEnodebTacs:       nil,
			NatEnabled:               true,
		},
		"pipelined": &lte_mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24",
			NatEnabled:    true,
			DefaultRuleId: "",
			Services: []lte_mconfig.PipelineD_NetworkServices{
				lte_mconfig.PipelineD_ENFORCEMENT,
			},
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:        protos.LogLevel_INFO,
			LteAuthOp:       []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:      []byte("\x80\x00"),
			SubProfiles:     nil,
			HssRelayEnabled: false,
		},
		"policydb": &lte_mconfig.PolicyDB{
			LogLevel: protos.LogLevel_INFO,
		},
		"sessiond": &lte_mconfig.SessionD{
			LogLevel:         protos.LogLevel_INFO,
			GxGyRelayEnabled: false,
			WalletExhaustDetection: &lte_mconfig.WalletExhaustDetection{
				TerminateOnExhaust: false,
			},
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
	}

	actual, err := build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

// Minimal configuration of enodeB, inherit rest of props from nw/gw configs
func TestBuilder_BuildInheritedProperties(t *testing.T) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{}))
	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{}))
	lte_test_init.StartTestService(t)

	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: lte_models.NewDefaultTDDNetworkConfig(),
			orc8r.DnsdNetworkType: &models.NetworkDNSConfig{
				EnableCaching: swag.Bool(true),
			},
		},
	}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "gw1",
		Associations: []storage.TypeAndKey{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config: newDefaultGatewayConfig(),
		Associations: []storage.TypeAndKey{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
		},
		ParentAssociations: []storage.TypeAndKey{gw.GetTypeAndKey()},
	}
	enb := configurator.NetworkEntity{
		Type: lte.CellularEnodebEntityType, Key: "enb1",
		Config: &lte_models.EnodebConfig{
			ConfigType: "MANAGED",
			ManagedConfig: &lte_models.EnodebConfiguration{
				CellID:          swag.Uint32(42),
				DeviceClass:     "Baicells ID TDD/FDD",
				TransmitEnabled: swag.Bool(true),
			},
		},
		ParentAssociations: []storage.TypeAndKey{lteGW.GetTypeAndKey()},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{enb, lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTypeAndKey(), To: lteGW.GetTypeAndKey()},
			{From: lteGW.GetTypeAndKey(), To: enb.GetTypeAndKey()},
		},
	}

	expected := map[string]proto.Message{
		"enodebd": &lte_mconfig.EnodebD{
			LogLevel: protos.LogLevel_INFO,
			Pci:      260,
			TddConfig: &lte_mconfig.EnodebD_TDDConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
			BandwidthMhz:        20,
			AllowEnodebTransmit: true,
			Tac:                 1,
			PlmnidList:          "00101",
			CsfbRat:             lte_mconfig.EnodebD_CSFBRAT_2G,
			Arfcn_2G:            nil,
			EnbConfigsBySerial: map[string]*lte_mconfig.EnodebD_EnodebConfig{
				"enb1": {
					Earfcndl:               44590,
					SubframeAssignment:     2,
					SpecialSubframePattern: 7,
					Pci:                    260,
					TransmitEnabled:        true,
					DeviceClass:            "Baicells ID TDD/FDD",
					BandwidthMhz:           20,
					Tac:                    1,
					CellId:                 42,
				},
			},
		},
		"mobilityd": &lte_mconfig.MobilityD{
			LogLevel: protos.LogLevel_INFO,
			IpBlock:  "192.168.128.0/24",
		},
		"mme": &lte_mconfig.MME{
			LogLevel:                 protos.LogLevel_INFO,
			Mcc:                      "001",
			Mnc:                      "01",
			Tac:                      1,
			MmeCode:                  1,
			MmeGid:                   1,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			EnableDnsCaching:         false,
			AttachedEnodebTacs:       []int32{1},
			NatEnabled:               true,
		},
		"pipelined": &lte_mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24",
			NatEnabled:    true,
			DefaultRuleId: "",
			Services: []lte_mconfig.PipelineD_NetworkServices{
				lte_mconfig.PipelineD_ENFORCEMENT,
			},
			SgiManagementIfaceVlan: "",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:        protos.LogLevel_INFO,
			LteAuthOp:       []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:      []byte("\x80\x00"),
			SubProfiles:     nil,
			HssRelayEnabled: false,
		},
		"policydb": &lte_mconfig.PolicyDB{
			LogLevel: protos.LogLevel_INFO,
		},
		"sessiond": &lte_mconfig.SessionD{
			LogLevel:         protos.LogLevel_INFO,
			GxGyRelayEnabled: false,
			WalletExhaustDetection: &lte_mconfig.WalletExhaustDetection{
				TerminateOnExhaust: false,
			},
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
	}

	actual, err := build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestBuilder_BuildUnmanagedEnbConfig(t *testing.T) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{}))
	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{}))
	lte_test_init.StartTestService(t)

	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: lte_models.NewDefaultTDDNetworkConfig(),
			orc8r.DnsdNetworkType: &models.NetworkDNSConfig{
				EnableCaching: swag.Bool(true),
			},
		},
	}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "gw1",
		Associations: []storage.TypeAndKey{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config: newDefaultGatewayConfig(),
		Associations: []storage.TypeAndKey{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
		},
		ParentAssociations: []storage.TypeAndKey{gw.GetTypeAndKey()},
	}
	enb := configurator.NetworkEntity{
		Type: lte.CellularEnodebEntityType, Key: "enb1",
		Config:             newDefaultUnmanagedEnodebConfig(),
		ParentAssociations: []storage.TypeAndKey{lteGW.GetTypeAndKey()},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{enb, lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTypeAndKey(), To: lteGW.GetTypeAndKey()},
			{From: lteGW.GetTypeAndKey(), To: enb.GetTypeAndKey()},
		},
	}

	expected := map[string]proto.Message{
		"enodebd": &lte_mconfig.EnodebD{
			LogLevel: protos.LogLevel_INFO,
			Pci:      260,
			TddConfig: &lte_mconfig.EnodebD_TDDConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
			BandwidthMhz:        20,
			AllowEnodebTransmit: true,
			Tac:                 1,
			PlmnidList:          "00101",
			CsfbRat:             lte_mconfig.EnodebD_CSFBRAT_2G,
			Arfcn_2G:            nil,
			EnbConfigsBySerial: map[string]*lte_mconfig.EnodebD_EnodebConfig{
				"enb1": {
					CellId: 138777000,
					Tac:    1,
				},
			},
		},
		"mobilityd": &lte_mconfig.MobilityD{
			LogLevel: protos.LogLevel_INFO,
			IpBlock:  "192.168.128.0/24",
		},
		"mme": &lte_mconfig.MME{
			LogLevel:                 protos.LogLevel_INFO,
			Mcc:                      "001",
			Mnc:                      "01",
			Tac:                      1,
			MmeCode:                  1,
			MmeGid:                   1,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			EnableDnsCaching:         false,
			AttachedEnodebTacs:       []int32{1},
			NatEnabled:               true,
		},
		"pipelined": &lte_mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24",
			NatEnabled:    true,
			DefaultRuleId: "",
			Services: []lte_mconfig.PipelineD_NetworkServices{
				lte_mconfig.PipelineD_ENFORCEMENT,
			},
			SgiManagementIfaceVlan: "",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:        protos.LogLevel_INFO,
			LteAuthOp:       []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:      []byte("\x80\x00"),
			SubProfiles:     nil,
			HssRelayEnabled: false,
		},
		"policydb": &lte_mconfig.PolicyDB{
			LogLevel: protos.LogLevel_INFO,
		},
		"sessiond": &lte_mconfig.SessionD{
			LogLevel:         protos.LogLevel_INFO,
			GxGyRelayEnabled: false,
			WalletExhaustDetection: &lte_mconfig.WalletExhaustDetection{
				TerminateOnExhaust: false,
			},
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
	}

	actual, err := build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func build(network *configurator.Network, graph *configurator.EntityGraph, gatewayID string) (map[string]proto.Message, error) {
	networkProto, err := network.ToProto(serdes.Network)
	if err != nil {
		return nil, err
	}
	graphProto, err := graph.ToProto(serdes.Entity)
	if err != nil {
		return nil, err
	}
	builder := mconfig.NewRemoteBuilder(lte_service.ServiceName)
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

func newDefaultGatewayConfig() *lte_models.GatewayCellularConfigs {
	return &lte_models.GatewayCellularConfigs{
		Ran: &lte_models.GatewayRanConfigs{
			Pci:             260,
			TransmitEnabled: swag.Bool(true),
		},
		Epc: &lte_models.GatewayEpcConfigs{
			NatEnabled: swag.Bool(true),
			IPBlock:    "192.168.128.0/24",
		},
		NonEpsService: &lte_models.GatewayNonEpsConfigs{
			CsfbMcc:              "001",
			CsfbMnc:              "01",
			Lac:                  swag.Uint32(1),
			CsfbRat:              swag.Uint32(0),
			Arfcn2g:              nil,
			NonEpsServiceControl: swag.Uint32(0),
		},
		DNS: &lte_models.GatewayDNSConfigs{
			DhcpServerEnabled: swag.Bool(true),
			EnableCaching:     swag.Bool(false),
			LocalTTL:          swag.Int32(0),
		},
	}
}

func newGatewayConfigNonNat(vlan string, sgi_ip string, sgi_gw string) *lte_models.GatewayCellularConfigs {
	return &lte_models.GatewayCellularConfigs{
		Ran: &lte_models.GatewayRanConfigs{
			Pci:             260,
			TransmitEnabled: swag.Bool(true),
		},
		Epc: &lte_models.GatewayEpcConfigs{
			NatEnabled:                 swag.Bool(false),
			IPBlock:                    "192.168.128.0/24",
			SgiManagementIfaceVlan:     vlan,
			SgiManagementIfaceStaticIP: sgi_ip,
			SgiManagementIfaceGw:       sgi_gw,
		},
		NonEpsService: &lte_models.GatewayNonEpsConfigs{
			CsfbMcc:              "001",
			CsfbMnc:              "01",
			Lac:                  swag.Uint32(1),
			CsfbRat:              swag.Uint32(0),
			Arfcn2g:              nil,
			NonEpsServiceControl: swag.Uint32(0),
		},
		DNS: &lte_models.GatewayDNSConfigs{
			DhcpServerEnabled: swag.Bool(true),
			EnableCaching:     swag.Bool(false),
			LocalTTL:          swag.Int32(0),
		},
	}
}

func newDefaultEnodebConfig() *lte_models.EnodebConfig {
	return &lte_models.EnodebConfig{
		ConfigType: "MANAGED",
		ManagedConfig: &lte_models.EnodebConfiguration{
			Earfcndl:               39150,
			SubframeAssignment:     2,
			SpecialSubframePattern: 7,
			Pci:                    260,
			CellID:                 swag.Uint32(138777000),
			Tac:                    15000,
			BandwidthMhz:           20,
			TransmitEnabled:        swag.Bool(true),
			DeviceClass:            "Baicells ID TDD/FDD",
		},
	}
}

func newDefaultUnmanagedEnodebConfig() *lte_models.EnodebConfig {
	ip := strfmt.IPv4("192.168.0.124")
	return &lte_models.EnodebConfig{
		ConfigType: "UNMANAGED",
		UnmanagedConfig: &lte_models.UnmanagedEnodebConfiguration{
			CellID:    swag.Uint32(138777000),
			Tac:       swag.Uint32(1),
			IPAddress: &ip,
		},
	}
}

func setEPCNetworkServices(services []string, nw *configurator.Network) {
	inwConfig := nw.Configs[lte.CellularNetworkConfigType]
	cellularNwConfig := inwConfig.(*lte_models.NetworkCellularConfigs)
	cellularNwConfig.Epc.NetworkServices = services

	nw.Configs[lte.CellularNetworkConfigType] = cellularNwConfig
}

func setEPCNetworkIPAllocator(nw *configurator.Network, mode string, static_ip bool,
	multi_apn bool) {
	inwConfig := nw.Configs[lte.CellularNetworkConfigType]
	cellularNwConfig := inwConfig.(*lte_models.NetworkCellularConfigs)
	cellularNwConfig.Epc.Mobility = &lte_models.NetworkEpcConfigsMobility{
		IPAllocationMode:           mode,
		EnableStaticIPAssignments:  static_ip,
		EnableMultiApnIPAllocation: multi_apn,
	}

	nw.Configs[lte.CellularNetworkConfigType] = cellularNwConfig
}
