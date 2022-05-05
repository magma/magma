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

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"magma/feg/cloud/go/feg"
	feg_serdes "magma/feg/cloud/go/serdes"
	feg_models "magma/feg/cloud/go/services/feg/obsidian/models"
	"magma/lte/cloud/go/lte"
	lte_mconfig "magma/lte/cloud/go/protos/mconfig"
	"magma/lte/cloud/go/serdes"
	lte_service "magma/lte/cloud/go/services/lte"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	storage_configurator "magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

const (
	// randomizedInterval300 and similar constants are calculated by adding delta * (300/5) to 300 where delta is a
	// random fraction between 0 and 1 based on the FNV-1 hash of the string "gw1". These constants emulate return
	// values of builder_servicer.getRandomizedSyncInterval("gw1", 300) etc.
	randomizedInterval300 = 344
	randomizedInterval120 = 137
	randomizedInterval90  = 103
	randomizedInterval60  = 68
)

func TestBuilder_Build(t *testing.T) {
	lte_test_init.StartTestService(t)

	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: lte_models.NewDefaultTDDNetworkConfig(),
			orc8r.DnsdNetworkType: &models.NetworkDNSConfig{
				EnableCaching: swag.Bool(true),
			},
			orc8r.NetworkSentryConfig: &models.NetworkSentryConfig{
				SampleRate:   swag.Float32(0.75),
				UploadMmeLog: true,
				URLPython:    "https://www.example.com/v1/api",
				URLNative:    "https://www.example.com/v1/api",
			},
		},
	}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "gw1",
		Associations: storage.TKs{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config: newDefaultGatewayConfig(),
		Associations: storage.TKs{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
		},
		ParentAssociations: storage.TKs{gw.GetTK()},
	}
	enb := configurator.NetworkEntity{
		Type: lte.CellularEnodebEntityType, Key: "enb1",
		Config:             newDefaultEnodebConfig(),
		ParentAssociations: storage.TKs{lteGW.GetTK()},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{enb, lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
			{From: lteGW.GetTK(), To: enb.GetTK()},
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
			MmeRelativeCapacity:      10,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			EnableDnsCaching:         false,
			AttachedEnodebTacs:       []int32{15000},
			NatEnabled:               true,
			CongestionControlEnabled: true,
			SentryConfig: &lte_mconfig.SentryConfig{
				SampleRate:   0.75,
				UploadMmeLog: true,
				DsnPython:    "https://www.example.com/v1/api",
				DsnNative:    "https://www.example.com/v1/api",
			},
			Enable5GFeatures: false,
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
			HeConfig:               &lte_mconfig.PipelineD_HEConfig{},
			LiUes:                  &lte_mconfig.PipelineD_LiUes{},
			Enable5GFeatures:       false,
			UpfNodeIdentifier:      "192.168.200.1",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:         protos.LogLevel_INFO,
			LteAuthOp:        []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:       []byte("\x80\x00"),
			SubProfiles:      nil,
			HssRelayEnabled:  false,
			SyncInterval:     randomizedInterval300,
			Enable5GFeatures: false,
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
			SentryConfig: &lte_mconfig.SentryConfig{
				SampleRate:   0.75,
				UploadMmeLog: true,
				DsnPython:    "https://www.example.com/v1/api",
				DsnNative:    "https://www.example.com/v1/api",
			},
			Enable5GFeatures: false,
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
		"liagentd": &lte_mconfig.LIAgentD{
			LogLevel: protos.LogLevel_INFO,
		},
		"sctpd": &lte_mconfig.SctpD{
			LogLevel: protos.LogLevel_INFO,
		},
	}

	// Happy path
	actual, err := buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	// Break with non-allowed network service
	setEPCNetworkServices([]string{"0xdeadbeef"}, &nw)
	_, err = buildNonFederated(&nw, &graph, "gw1")
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
		HeConfig:          &lte_mconfig.PipelineD_HEConfig{},
		LiUes:             &lte_mconfig.PipelineD_LiUes{},
		Enable5GFeatures:  false,
		UpfNodeIdentifier: "192.168.200.1",
	}
	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	// verify restricted plmns
	setEpcNetworkRestrictedPlmns(&nw, []*lte_models.PlmnConfig{
		{
			Mcc: "100",
			Mnc: "010",
		},
		{
			Mcc: "110",
			Mnc: "210",
		},
	})
	mmeVals := expected["mme"].(*lte_mconfig.MME)
	mmeVals.RestrictedPlmns = []*lte_mconfig.MME_PlmnConfig{
		{
			Mcc: "100",
			Mnc: "010",
		},
		{
			Mcc: "110",
			Mnc: "210",
		},
	}

	actual, err = buildLTEFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	// verify restricted imei
	setEpcNetworkRestrictedImeis(&nw, []*lte_models.Imei{
		{
			Tac: "01300600",
			Snr: "176148",
		},
		{
			Tac: "01200200",
			Snr: "176222",
		},
	})
	mmeVals.RestrictedImeis = []*lte_mconfig.MME_ImeiConfig{
		{
			Tac: "01300600",
			Snr: "176148",
		},
		{
			Tac: "01200200",
			Snr: "176222",
		},
	}

	actual, err = buildLTEFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	// verify service area map
	setEpcNetworkServiceAreaMap(&nw, map[string]lte_models.TacList{
		"001": []lte_models.Tac{111, 112},
		"002": []lte_models.Tac{211, 122},
	})
	mmeVals.ServiceAreaMaps = map[string]*lte_mconfig.MME_TacList{
		"001": {
			Tac: []uint32{111, 112},
		},
		"002": {
			Tac: []uint32{211, 122},
		},
	}
	actual, err = buildLTEFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)
}

func TestBuilder_Build_NonNat(t *testing.T) {
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
		Associations: storage.TKs{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             newGatewayConfigNonNat("", "", "", "", ""),
		ParentAssociations: storage.TKs{gw.GetTK()},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
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
			MmeRelativeCapacity:      10,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			AttachedEnodebTacs:       nil,
			NatEnabled:               false,
			CongestionControlEnabled: true,
			Enable5GFeatures:         false,
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
			HeConfig:               &lte_mconfig.PipelineD_HEConfig{},
			LiUes:                  &lte_mconfig.PipelineD_LiUes{},
			Enable5GFeatures:       false,
			UpfNodeIdentifier:      "192.168.200.1",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:         protos.LogLevel_INFO,
			LteAuthOp:        []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:       []byte("\x80\x00"),
			SubProfiles:      nil,
			HssRelayEnabled:  false,
			SyncInterval:     randomizedInterval300,
			Enable5GFeatures: false,
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
			Enable5GFeatures: false,
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
		"liagentd": &lte_mconfig.LIAgentD{
			LogLevel: protos.LogLevel_INFO,
		},
		"sctpd": &lte_mconfig.SctpD{
			LogLevel: protos.LogLevel_INFO,
		},
	}
	actual, err := buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	setEPCNetworkIPAllocator(&nw, lte_models.DHCPBroadcastAllocationMode, false, false)
	expected["mobilityd"] = &lte_mconfig.MobilityD{
		LogLevel:        protos.LogLevel_INFO,
		IpBlock:         "192.168.128.0/24",
		IpAllocatorType: lte_mconfig.MobilityD_DHCP,
		StaticIpEnabled: false,
	}
	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	setEPCNetworkIPAllocator(&nw, lte_models.NATAllocationMode, false, false)
	expected["mobilityd"] = &lte_mconfig.MobilityD{
		LogLevel:        protos.LogLevel_INFO,
		IpBlock:         "192.168.128.0/24",
		IpAllocatorType: lte_mconfig.MobilityD_IP_POOL,
		StaticIpEnabled: false,
	}
	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	setEPCNetworkIPAllocator(&nw, lte_models.NATAllocationMode, true, false)
	expected["mobilityd"] = &lte_mconfig.MobilityD{
		LogLevel:        protos.LogLevel_INFO,
		IpBlock:         "192.168.128.0/24",
		IpAllocatorType: lte_mconfig.MobilityD_IP_POOL,
		StaticIpEnabled: true,
	}
	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	setEPCNetworkIPAllocator(&nw, lte_models.DHCPBroadcastAllocationMode, true, false)
	expected["mobilityd"] = &lte_mconfig.MobilityD{
		LogLevel:        protos.LogLevel_INFO,
		IpBlock:         "192.168.128.0/24",
		IpAllocatorType: lte_mconfig.MobilityD_DHCP,
		StaticIpEnabled: true,
	}
	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	// validate multi apn mconfig
	setEPCNetworkIPAllocator(&nw, lte_models.DHCPBroadcastAllocationMode, true, true)
	expected["mobilityd"] = &lte_mconfig.MobilityD{
		LogLevel:        protos.LogLevel_INFO,
		IpBlock:         "192.168.128.0/24",
		IpAllocatorType: lte_mconfig.MobilityD_DHCP,
		StaticIpEnabled: true,
		MultiApnIpAlloc: true,
	}
	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	// validate SGi vlan tag mconfig
	lteGW = configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             newGatewayConfigNonNat("30", "", "", "", ""),
		ParentAssociations: storage.TKs{gw.GetTK()},
	}
	graph = configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
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
		HeConfig:               &lte_mconfig.PipelineD_HEConfig{},
		LiUes:                  &lte_mconfig.PipelineD_LiUes{},
		Enable5GFeatures:       false,
		UpfNodeIdentifier:      "192.168.200.1",
	}

	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	// validate SGi ip address
	// validate SGi vlan tag mconfig
	lteGW = configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             newGatewayConfigNonNat("44", "1.2.3.4", "", "", ""),
		ParentAssociations: storage.TKs{gw.GetTK()},
	}
	graph = configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
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
		HeConfig:                 &lte_mconfig.PipelineD_HEConfig{},
		LiUes:                    &lte_mconfig.PipelineD_LiUes{},
		Enable5GFeatures:         false,
		UpfNodeIdentifier:        "192.168.200.1",
	}

	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	// validate SGi ip address and gateway
	// validate SGi vlan tag mconfig
	lteGW = configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             newGatewayConfigNonNat("55", "1.2.3.4/24", "1.2.3.1", "", ""),
		ParentAssociations: storage.TKs{gw.GetTK()},
	}
	graph = configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
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
		HeConfig:                 &lte_mconfig.PipelineD_HEConfig{},
		LiUes:                    &lte_mconfig.PipelineD_LiUes{},
		UpfNodeIdentifier:        "192.168.200.1",
		Enable5GFeatures:         false,
	}

	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	// validate SGi ipv6 address and gateway
	// validate SGi vlan tag mconfig
	lteGW = configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             newGatewayConfigNonNat("55", "", "", "2001:4860:4860:0:0:0:0:1111/64", "2a12:577:9941:f99c:0002:0001:c731:f114"),
		ParentAssociations: storage.TKs{gw.GetTK()},
	}
	graph = configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
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
		SgiManagementIfaceVlan:     "55",
		SgiManagementIfaceIpv6Addr: "2001:4860:4860:0:0:0:0:1111/64",
		SgiManagementIfaceIpv6Gw:   "2a12:577:9941:f99c:0002:0001:c731:f114",
		HeConfig:                   &lte_mconfig.PipelineD_HEConfig{},
		LiUes:                      &lte_mconfig.PipelineD_LiUes{},
		UpfNodeIdentifier:          "192.168.200.1",
		Enable5GFeatures:           false,
	}

	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)

	// validate SGi ipv4 and ipv6 address and gateway
	// validate SGi vlan tag mconfig
	lteGW = configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             newGatewayConfigNonNat("55", "1.2.3.4/24", "1.2.3.1", "2001:4860:4860:0:0:0:0:9999/96", "2a12:577:9941:f99c:0002:0001:c731:f114"),
		ParentAssociations: storage.TKs{gw.GetTK()},
	}
	graph = configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
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
		SgiManagementIfaceVlan:     "55",
		SgiManagementIfaceIpAddr:   "1.2.3.4/24",
		SgiManagementIfaceGw:       "1.2.3.1",
		SgiManagementIfaceIpv6Addr: "2001:4860:4860:0:0:0:0:9999/96",
		SgiManagementIfaceIpv6Gw:   "2a12:577:9941:f99c:0002:0001:c731:f114",
		HeConfig:                   &lte_mconfig.PipelineD_HEConfig{},
		LiUes:                      &lte_mconfig.PipelineD_LiUes{},
		UpfNodeIdentifier:          "192.168.200.1",
		Enable5GFeatures:           false,
	}

	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)
}

func TestBuilder_Build_NgcConfig(t *testing.T) {
	lte_test_init.StartTestService(t)

	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: lte_models.NewDefaultTDDNetworkConfig(),
		},
	}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "gw1",
		Associations: storage.TKs{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	gatewayConfig := newGatewayConfigWithNGC()
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             gatewayConfig,
		ParentAssociations: storage.TKs{gw.GetTK()},
	}

	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
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
			LogLevel:                      protos.LogLevel_INFO,
			Mcc:                           "001",
			Mnc:                           "01",
			Tac:                           1,
			MmeCode:                       1,
			MmeGid:                        1,
			MmeRelativeCapacity:           10,
			NonEpsServiceControl:          lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                       "001",
			CsfbMnc:                       "01",
			Lac:                           1,
			HssRelayEnabled:               false,
			CloudSubscriberdbEnabled:      false,
			AttachedEnodebTacs:            nil,
			NatEnabled:                    true,
			CongestionControlEnabled:      true,
			Enable5GFeatures:              false,
			AmfName:                       "amf.example.org",
			AmfSetId:                      "2A1",
			AmfRegionId:                   "C1",
			AmfPointer:                    "1F",
			AmfDefaultSliceServiceType:    25,
			AmfDefaultSliceDifferentiator: "AFAFAF",
		},
		"pipelined": &lte_mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24",
			NatEnabled:    true,
			DefaultRuleId: "",
			Services: []lte_mconfig.PipelineD_NetworkServices{
				lte_mconfig.PipelineD_ENFORCEMENT,
			},
			HeConfig:          &lte_mconfig.PipelineD_HEConfig{},
			LiUes:             &lte_mconfig.PipelineD_LiUes{},
			Enable5GFeatures:  false,
			UpfNodeIdentifier: "192.168.200.1",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:         protos.LogLevel_INFO,
			LteAuthOp:        []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:       []byte("\x80\x00"),
			SubProfiles:      nil,
			HssRelayEnabled:  false,
			SyncInterval:     randomizedInterval300,
			Enable5GFeatures: false,
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
			Enable5GFeatures: false,
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
		"liagentd": &lte_mconfig.LIAgentD{
			LogLevel: protos.LogLevel_INFO,
		},
		"sctpd": &lte_mconfig.SctpD{
			LogLevel: protos.LogLevel_INFO,
		},
	}

	actual, err := buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)
}

func TestBuilder_Build_BaseCase(t *testing.T) {
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
		Associations: storage.TKs{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	heConfig := &lte_models.GatewayHeConfig{
		EnableHeaderEnrichment: swag.Bool(true),
		EnableEncryption:       swag.Bool(true),
		HeEncryptionAlgorithm:  "RC4",
		HeHashFunction:         "MD5",
		HeEncodingType:         "BASE64",
		EncryptionKey:          "melting_the_core",
		HmacKey:                "magmamagma",
	}
	gatewayConfig := newDefaultGatewayConfig()
	gatewayConfig.HeConfig = heConfig
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             gatewayConfig,
		ParentAssociations: storage.TKs{gw.GetTK()},
	}

	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
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
			MmeRelativeCapacity:      10,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			AttachedEnodebTacs:       nil,
			NatEnabled:               true,
			CongestionControlEnabled: true,
			Enable5GFeatures:         false,
		},
		"pipelined": &lte_mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24",
			NatEnabled:    true,
			DefaultRuleId: "",
			Services: []lte_mconfig.PipelineD_NetworkServices{
				lte_mconfig.PipelineD_ENFORCEMENT,
			},
			HeConfig: &lte_mconfig.PipelineD_HEConfig{
				EnableHeaderEnrichment: true,
				EnableEncryption:       true,
				EncryptionAlgorithm:    lte_mconfig.PipelineD_HEConfig_RC4,
				HashFunction:           lte_mconfig.PipelineD_HEConfig_MD5,
				EncodingType:           lte_mconfig.PipelineD_HEConfig_BASE64,
				EncryptionKey:          "melting_the_core",
				HmacKey:                "magmamagma",
			},
			LiUes:             &lte_mconfig.PipelineD_LiUes{},
			Enable5GFeatures:  false,
			UpfNodeIdentifier: "192.168.200.1",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:         protos.LogLevel_INFO,
			LteAuthOp:        []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:       []byte("\x80\x00"),
			SubProfiles:      nil,
			HssRelayEnabled:  false,
			SyncInterval:     randomizedInterval300,
			Enable5GFeatures: false,
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
			Enable5GFeatures: false,
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
		"liagentd": &lte_mconfig.LIAgentD{
			LogLevel: protos.LogLevel_INFO,
		},
		"sctpd": &lte_mconfig.SctpD{
			LogLevel: protos.LogLevel_INFO,
		},
	}

	actual, err := buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)
}

func TestBuilder_Build_ConfigOverride(t *testing.T) {
	lte_test_init.StartTestService(t)

	nwConfig := lte_models.NewDefaultTDDNetworkConfig()
	// Change sync interval from the default 300
	nwConfig.Epc.SubscriberdbSyncInterval = lte_models.SubscriberdbSyncInterval(120)

	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: nwConfig,
		},
	}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "gw1",
		Associations: storage.TKs{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}

	gatewayConfig := newDefaultGatewayConfig()
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             gatewayConfig,
		ParentAssociations: storage.TKs{gw.GetTK()},
	}

	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
		},
	}

	// no override. nw-wide 120 expected
	expected := map[string]proto.Message{
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:        protos.LogLevel_INFO,
			LteAuthOp:       []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:      []byte("\x80\x00"),
			SubProfiles:     nil,
			HssRelayEnabled: false,
			SyncInterval:    randomizedInterval120,
		},
	}

	actual, err := buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMessagesEqual(t, proto.MessageV2(expected["subscriberdb"]), proto.MessageV2(actual["subscriberdb"]))

	gatewayConfig.Epc.SubscriberdbSyncInterval = lte_models.SubscriberdbSyncInterval(90)
	// override. gw-specific 90 expected
	expected["subscriberdb"].(*lte_mconfig.SubscriberDB).SyncInterval = randomizedInterval90

	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMessagesEqual(t, proto.MessageV2(expected["subscriberdb"]), proto.MessageV2(actual["subscriberdb"]))

	nwConfig.Epc.SubscriberdbSyncInterval = 0
	gatewayConfig.Epc.SubscriberdbSyncInterval = 0

	// nw-wide and gw-specific not set. Service-level default expected
	expected["subscriberdb"].(*lte_mconfig.SubscriberDB).SyncInterval = randomizedInterval300

	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMessagesEqual(t, proto.MessageV2(expected["subscriberdb"]), proto.MessageV2(actual["subscriberdb"]))

	lte_test_init.StartTestServiceWithConfig(t, lte_service.Config{DefaultSubscriberdbSyncInterval: 30})

	// nw-wide and gw-specific not set. Service-level default (30) too low. Enforced minimum (60) expected
	expected["subscriberdb"].(*lte_mconfig.SubscriberDB).SyncInterval = randomizedInterval60

	actual, err = buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMessagesEqual(t, proto.MessageV2(expected["subscriberdb"]), proto.MessageV2(actual["subscriberdb"]))
}

func TestBuilder_Build_FederatedBaseCase(t *testing.T) {
	lte_test_init.StartTestService(t)

	// create a network and add feg.FederatedNetworkType item
	cellularConfig := lte_models.NewDefaultTDDNetworkConfig()
	cellularConfig.FegNetworkID = "n1" // this matches with NewDefaultFederatedNetworkConfigs
	nw := configurator.Network{
		ID:   "n_lte_1", // use a different name so it is not hte same as federated
		Type: feg.FederatedLteNetworkType,
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: cellularConfig,
			orc8r.DnsdNetworkType: &models.NetworkDNSConfig{
				EnableCaching: swag.Bool(true),
			},
			feg.FederatedNetworkType: feg_models.NewDefaultFederatedNetworkConfigs(),
		},
	}

	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "gw1",
		Associations: storage.TKs{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	heConfig := &lte_models.GatewayHeConfig{
		EnableHeaderEnrichment: swag.Bool(true),
		EnableEncryption:       swag.Bool(true),
		HeEncryptionAlgorithm:  "RC4",
		HeHashFunction:         "MD5",
		HeEncodingType:         "BASE64",
		EncryptionKey:          "melting_the_core",
		HmacKey:                "magmamagma",
	}
	gatewayConfig := newDefaultGatewayConfig()
	gatewayConfig.HeConfig = heConfig
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             gatewayConfig,
		ParentAssociations: storage.TKs{gw.GetTK()},
	}

	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
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
			MmeRelativeCapacity:      10,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			AttachedEnodebTacs:       nil,
			NatEnabled:               true,
			CongestionControlEnabled: true,
			FederatedModeMap: &lte_mconfig.FederatedModeMap{
				Enabled: true,
				Mapping: []*lte_mconfig.ModeMapItem{
					{
						Mode:      lte_mconfig.ModeMapItem_S8_SUBSCRIBER,
						Apn:       "internet1",
						ImsiRange: "000000000000001",
						Plmn:      "00101",
					},
				},
			},
			Enable5GFeatures: false,
		},
		"pipelined": &lte_mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24",
			NatEnabled:    true,
			DefaultRuleId: "",
			Services: []lte_mconfig.PipelineD_NetworkServices{
				lte_mconfig.PipelineD_ENFORCEMENT,
			},
			HeConfig: &lte_mconfig.PipelineD_HEConfig{
				EnableHeaderEnrichment: true,
				EnableEncryption:       true,
				EncryptionAlgorithm:    lte_mconfig.PipelineD_HEConfig_RC4,
				HashFunction:           lte_mconfig.PipelineD_HEConfig_MD5,
				EncodingType:           lte_mconfig.PipelineD_HEConfig_BASE64,
				EncryptionKey:          "melting_the_core",
				HmacKey:                "magmamagma",
			},
			LiUes:             &lte_mconfig.PipelineD_LiUes{},
			Enable5GFeatures:  false,
			UpfNodeIdentifier: "192.168.200.1",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:         protos.LogLevel_INFO,
			LteAuthOp:        []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:       []byte("\x80\x00"),
			SubProfiles:      nil,
			HssRelayEnabled:  false,
			SyncInterval:     randomizedInterval300,
			Enable5GFeatures: false,
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
			Enable5GFeatures: false,
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
		"liagentd": &lte_mconfig.LIAgentD{
			LogLevel: protos.LogLevel_INFO,
		},
		"sctpd": &lte_mconfig.SctpD{
			LogLevel: protos.LogLevel_INFO,
		},
	}

	// Use LTE FEG NETWORK parser for this case
	actual, err := buildLTEFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)
}

// Minimal configuration of enodeB, inherit rest of props from nw/gw configs
func TestBuilder_BuildInheritedProperties(t *testing.T) {
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
		Associations: storage.TKs{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config: newDefaultGatewayConfig(),
		Associations: storage.TKs{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
		},
		ParentAssociations: storage.TKs{gw.GetTK()},
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
		ParentAssociations: storage.TKs{lteGW.GetTK()},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{enb, lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
			{From: lteGW.GetTK(), To: enb.GetTK()},
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
			MmeRelativeCapacity:      10,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			EnableDnsCaching:         false,
			AttachedEnodebTacs:       []int32{1},
			NatEnabled:               true,
			CongestionControlEnabled: true,
			Enable5GFeatures:         false,
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
			HeConfig:               &lte_mconfig.PipelineD_HEConfig{},
			LiUes:                  &lte_mconfig.PipelineD_LiUes{},
			Enable5GFeatures:       false,
			UpfNodeIdentifier:      "192.168.200.1",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:         protos.LogLevel_INFO,
			LteAuthOp:        []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:       []byte("\x80\x00"),
			SubProfiles:      nil,
			HssRelayEnabled:  false,
			SyncInterval:     randomizedInterval300,
			Enable5GFeatures: false,
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
			Enable5GFeatures: false,
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
		"liagentd": &lte_mconfig.LIAgentD{
			LogLevel: protos.LogLevel_INFO,
		},
		"sctpd": &lte_mconfig.SctpD{
			LogLevel: protos.LogLevel_INFO,
		},
	}

	actual, err := buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)
}

func TestBuilder_BuildUnmanagedEnbConfig(t *testing.T) {
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
		Associations: storage.TKs{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config: newDefaultGatewayConfig(),
		Associations: storage.TKs{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
		},
		ParentAssociations: storage.TKs{gw.GetTK()},
	}
	enb := configurator.NetworkEntity{
		Type: lte.CellularEnodebEntityType, Key: "enb1",
		Config:             newDefaultUnmanagedEnodebConfig(),
		ParentAssociations: storage.TKs{lteGW.GetTK()},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{enb, lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
			{From: lteGW.GetTK(), To: enb.GetTK()},
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
					CellId:    138777000,
					Tac:       1,
					IpAddress: "192.168.0.124",
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
			MmeRelativeCapacity:      10,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			EnableDnsCaching:         false,
			AttachedEnodebTacs:       []int32{1},
			NatEnabled:               true,
			CongestionControlEnabled: true,
			Enable5GFeatures:         false,
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
			HeConfig:               &lte_mconfig.PipelineD_HEConfig{},
			LiUes:                  &lte_mconfig.PipelineD_LiUes{},
			Enable5GFeatures:       false,
			UpfNodeIdentifier:      "192.168.200.1",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:         protos.LogLevel_INFO,
			LteAuthOp:        []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:       []byte("\x80\x00"),
			SubProfiles:      nil,
			HssRelayEnabled:  false,
			SyncInterval:     randomizedInterval300,
			Enable5GFeatures: false,
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
			Enable5GFeatures: false,
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
		"liagentd": &lte_mconfig.LIAgentD{
			LogLevel: protos.LogLevel_INFO,
		},
		"sctpd": &lte_mconfig.SctpD{
			LogLevel: protos.LogLevel_INFO,
		},
	}

	actual, err := buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)
}

func TestBuilder_BuildCongestionControlConfig(t *testing.T) {
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
		Associations: storage.TKs{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config: newDefaultGatewayConfig(),
		Associations: storage.TKs{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
		},
		ParentAssociations: storage.TKs{gw.GetTK()},
	}
	enb := configurator.NetworkEntity{
		Type: lte.CellularEnodebEntityType, Key: "enb1",
		Config:             newDefaultUnmanagedEnodebConfig(),
		ParentAssociations: storage.TKs{lteGW.GetTK()},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{enb, lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
			{From: lteGW.GetTK(), To: enb.GetTK()},
		},
	}

	// Gateway specific config overrides network level
	gwConfig := lteGW.Config
	cellularGwConfig := gwConfig.(*lte_models.GatewayCellularConfigs)
	cellularGwConfig.Epc.CongestionControlEnabled = swag.Bool(false)

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
					CellId:    138777000,
					Tac:       1,
					IpAddress: "192.168.0.124",
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
			MmeRelativeCapacity:      10,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			EnableDnsCaching:         false,
			// Gateway congestion control enabled should be false
			CongestionControlEnabled: false,
			AttachedEnodebTacs:       []int32{1},
			NatEnabled:               true,
			Enable5GFeatures:         false,
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
			HeConfig:               &lte_mconfig.PipelineD_HEConfig{},
			LiUes:                  &lte_mconfig.PipelineD_LiUes{},
			Enable5GFeatures:       false,
			UpfNodeIdentifier:      "192.168.200.1",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:         protos.LogLevel_INFO,
			LteAuthOp:        []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:       []byte("\x80\x00"),
			SubProfiles:      nil,
			HssRelayEnabled:  false,
			SyncInterval:     randomizedInterval300,
			Enable5GFeatures: false,
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
			Enable5GFeatures: false,
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
		"liagentd": &lte_mconfig.LIAgentD{
			LogLevel: protos.LogLevel_INFO,
		},
		"sctpd": &lte_mconfig.SctpD{
			LogLevel: protos.LogLevel_INFO,
		},
	}

	actual, err := buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)
}

func TestBuilder_Build_MMEPool(t *testing.T) {
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
		Associations: storage.TKs{
			{Type: lte.CellularGatewayEntityType, Key: "gw1"},
		},
	}
	lteGatewayPool := configurator.NetworkEntity{
		Type: lte.CellularGatewayPoolEntityType, Key: "pool1",
		Config: &lte_models.CellularGatewayPoolConfigs{
			MmeGroupID: 2,
		},
	}
	lteGatewayConfigs := newDefaultGatewayConfig()
	lteGatewayConfigs.Pooling = lte_models.CellularGatewayPoolRecords{
		{
			GatewayPoolID:       "pool1",
			MmeCode:             3,
			MmeRelativeCapacity: 255,
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayEntityType, Key: "gw1",
		Config:             lteGatewayConfigs,
		Associations:       storage.TKs{},
		ParentAssociations: storage.TKs{gw.GetTK(), lteGatewayPool.GetTK()},
	}
	lteGatewayPool.Associations = storage.TKs{lteGW.GetTK()}

	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGatewayPool, lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTK(), To: lteGW.GetTK()},
			{From: lteGatewayPool.GetTK(), To: lteGW.GetTK()},
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
			MmeCode:                  3,
			MmeGid:                   2,
			MmeRelativeCapacity:      255,
			NonEpsServiceControl:     lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			HssRelayEnabled:          false,
			CloudSubscriberdbEnabled: false,
			AttachedEnodebTacs:       nil,
			NatEnabled:               true,
			CongestionControlEnabled: true,
			Enable5GFeatures:         false,
		},
		"pipelined": &lte_mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24",
			NatEnabled:    true,
			DefaultRuleId: "",
			Services: []lte_mconfig.PipelineD_NetworkServices{
				lte_mconfig.PipelineD_ENFORCEMENT,
			},
			HeConfig:          &lte_mconfig.PipelineD_HEConfig{},
			LiUes:             &lte_mconfig.PipelineD_LiUes{},
			Enable5GFeatures:  false,
			UpfNodeIdentifier: "192.168.200.1",
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:         protos.LogLevel_INFO,
			LteAuthOp:        []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:       []byte("\x80\x00"),
			SubProfiles:      nil,
			HssRelayEnabled:  false,
			SyncInterval:     randomizedInterval300,
			Enable5GFeatures: false,
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
			Enable5GFeatures: false,
		},
		"dnsd": &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
		},
		"liagentd": &lte_mconfig.LIAgentD{
			LogLevel: protos.LogLevel_INFO,
		},
		"sctpd": &lte_mconfig.SctpD{
			LogLevel: protos.LogLevel_INFO,
		},
	}

	actual, err := buildNonFederated(&nw, &graph, "gw1")
	assert.NoError(t, err)
	test_utils.AssertMapsEqual(t, expected, actual)
}

// buildLTEFederated builds a Federated_LTE network that comes from swagger feg_lte_network model
func buildLTEFederated(network *configurator.Network, graph *configurator.EntityGraph, gatewayID string) (map[string]proto.Message, error) {
	// use federated serded (this is still an LTE network)
	networkProto, err := network.ToProto(feg_serdes.Network)
	if err != nil {
		return nil, err
	}
	return buildImpl(networkProto, graph, gatewayID)
}

// buildNonFederated builds an non federated LTE network that comes from swagger lte_networl
func buildNonFederated(network *configurator.Network, graph *configurator.EntityGraph, gatewayID string) (map[string]proto.Message, error) {
	// use NON federated serded
	networkProto, err := network.ToProto(serdes.Network)
	if err != nil {
		return nil, err
	}
	return buildImpl(networkProto, graph, gatewayID)
}

func buildImpl(networkProto *storage_configurator.Network, graph *configurator.EntityGraph, gatewayID string) (map[string]proto.Message, error) {
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
			NatEnabled:               swag.Bool(true),
			IPBlock:                  "192.168.128.0/24",
			CongestionControlEnabled: swag.Bool(true),
			NodeIdentifier:           "192.168.200.1",
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
		HeConfig: &lte_models.GatewayHeConfig{},
		Pooling:  lte_models.CellularGatewayPoolRecords{},
	}
}

// TODO: simplify function args.
func newGatewayConfigNonNat(vlan string, sgi_ip string, sgi_gw string, sgi_ipv6 string, sgi_ipv6_gw string) *lte_models.GatewayCellularConfigs {
	return &lte_models.GatewayCellularConfigs{
		Ran: &lte_models.GatewayRanConfigs{
			Pci:             260,
			TransmitEnabled: swag.Bool(true),
		},
		Epc: &lte_models.GatewayEpcConfigs{
			NatEnabled:                 swag.Bool(false),
			IPBlock:                    "192.168.128.0/24",
			NodeIdentifier:             "192.168.200.1",
			SgiManagementIfaceVlan:     vlan,
			SgiManagementIfaceStaticIP: sgi_ip,
			SgiManagementIfaceGw:       sgi_gw,
			SgiManagementIfaceIPV6Addr: sgi_ipv6,
			SgiManagementIfaceIPV6Gw:   strfmt.IPv6(sgi_ipv6_gw),
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
		HeConfig: &lte_models.GatewayHeConfig{},
	}
}

func newGatewayConfigWithNGC() *lte_models.GatewayCellularConfigs {
	return &lte_models.GatewayCellularConfigs{
		Ran: &lte_models.GatewayRanConfigs{
			Pci:             260,
			TransmitEnabled: swag.Bool(true),
		},
		Epc: &lte_models.GatewayEpcConfigs{
			NatEnabled:               swag.Bool(true),
			IPBlock:                  "192.168.128.0/24",
			CongestionControlEnabled: swag.Bool(true),
			NodeIdentifier:           "192.168.200.1",
		},
		Ngc: &lte_models.GatewayNgcConfigs{
			AmfDefaultSd:  "AFAFAF",
			AmfDefaultSst: 25,
			AmfName:       "amf.example.org",
			AmfPointer:    "1F",
			AmfRegionID:   "C1",
			AmfSetID:      "2A1",
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
		HeConfig: &lte_models.GatewayHeConfig{},
		Pooling:  lte_models.CellularGatewayPoolRecords{},
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

func setEpcNetworkRestrictedPlmns(nw *configurator.Network, restrictedPlmns []*lte_models.PlmnConfig) {
	inwConfig := nw.Configs[lte.CellularNetworkConfigType]
	cellularNwConfig := inwConfig.(*lte_models.NetworkCellularConfigs)
	cellularNwConfig.Epc.RestrictedPlmns = restrictedPlmns
	nw.Configs[lte.CellularNetworkConfigType] = cellularNwConfig
}

func setEpcNetworkRestrictedImeis(nw *configurator.Network, restrictedImeis []*lte_models.Imei) {
	inwConfig := nw.Configs[lte.CellularNetworkConfigType]
	cellularNwConfig := inwConfig.(*lte_models.NetworkCellularConfigs)
	cellularNwConfig.Epc.RestrictedImeis = restrictedImeis
	nw.Configs[lte.CellularNetworkConfigType] = cellularNwConfig
}

func setEpcNetworkServiceAreaMap(nw *configurator.Network, serviceAreaMaps map[string]lte_models.TacList) {
	inwConfig := nw.Configs[lte.CellularNetworkConfigType]
	cellularNwConfig := inwConfig.(*lte_models.NetworkCellularConfigs)
	cellularNwConfig.Epc.ServiceAreaMaps = serviceAreaMaps
	nw.Configs[lte.CellularNetworkConfigType] = cellularNwConfig
}
