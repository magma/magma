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

	"magma/cwf/cloud/go/cwf"
	cwf_mconfig "magma/cwf/cloud/go/protos/mconfig"
	"magma/cwf/cloud/go/serdes"
	cwf_service "magma/cwf/cloud/go/services/cwf"
	"magma/cwf/cloud/go/services/cwf/obsidian/models"
	cwf_test_init "magma/cwf/cloud/go/services/cwf/test_init"
	feg_mconfig "magma/feg/cloud/go/protos/mconfig"
	fegmodels "magma/feg/cloud/go/services/feg/obsidian/models"
	lte_mconfig "magma/lte/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
	orc8r_mconfig "magma/orc8r/lib/go/protos/mconfig"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	cwf_test_init.StartTestService(t)

	t.Run("empty network config", func(t *testing.T) {
		nw := configurator.Network{ID: "n1"}
		gw := configurator.NetworkEntity{
			Type: orc8r.MagmadGatewayType, Key: "gw1",
			Associations: []storage.TypeAndKey{
				{Type: cwf.CwfGatewayType, Key: "gw1"},
			},
		}
		graph := configurator.EntityGraph{
			Entities: []configurator.NetworkEntity{gw},
		}

		expected := map[string]proto.Message{}

		actual, err := build(&nw, &graph, "gw1")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("network config exists", func(t *testing.T) {
		nw := configurator.Network{
			ID: "n1",
			Configs: map[string]interface{}{
				cwf.CwfNetworkType: defaultnwConfig,
			},
		}
		gw := configurator.NetworkEntity{
			Type: orc8r.MagmadGatewayType, Key: "gw1",
			Associations: []storage.TypeAndKey{
				{Type: cwf.CwfGatewayType, Key: "gw1"},
			},
		}
		cwfGW := configurator.NetworkEntity{
			Type: cwf.CwfGatewayType, Key: "gw1",
			Config:             defaultgwConfig,
			ParentAssociations: []storage.TypeAndKey{gw.GetTypeAndKey()},
		}
		haPair := configurator.NetworkEntity{
			Config: &models.CwfHaPairConfigs{TransportVirtualIP: "10.10.10.11"},
			Type:   cwf.CwfHAPairType,
			Key:    "pair1",
			Associations: []storage.TypeAndKey{
				{Type: cwf.CwfGatewayType, Key: "gw1"},
				{Type: cwf.CwfGatewayType, Key: "gw2"},
			},
		}
		graph := configurator.EntityGraph{
			Entities: []configurator.NetworkEntity{cwfGW, gw, haPair},
			Edges: []configurator.GraphEdge{
				{From: gw.GetTypeAndKey(), To: cwfGW.GetTypeAndKey()},
				{From: haPair.GetTypeAndKey(), To: cwfGW.GetTypeAndKey()},
			},
		}

		expected := map[string]proto.Message{
			"eap_aka": &feg_mconfig.EapAkaConfig{LogLevel: 1,
				Timeout: &feg_mconfig.EapAkaConfig_Timeouts{
					ChallengeMs:            20000,
					ErrorNotificationMs:    10000,
					SessionMs:              43200000,
					SessionAuthenticatedMs: 5000,
				},
				PlmnIds: nil,
			},
			"aaa_server": &feg_mconfig.AAAConfig{LogLevel: 1,
				IdleSessionTimeoutMs: 21600000,
				AccountingEnabled:    false,
				CreateSessionOnAuth:  false,
			},
			"pipelined": &lte_mconfig.PipelineD{
				LogLevel:      protos.LogLevel_INFO,
				UeIpBlock:     "192.168.128.0/24", // Unused by CWF
				NatEnabled:    false,
				DefaultRuleId: "",
				Services: []lte_mconfig.PipelineD_NetworkServices{
					lte_mconfig.PipelineD_DPI,
					lte_mconfig.PipelineD_ENFORCEMENT,
				},
				AllowedGrePeers: []*lte_mconfig.PipelineD_AllowedGrePeer{
					{Ip: "1.2.3.4/24"},
					{Ip: "1.1.1.1/24", Key: 111},
				},
				LiUes: &lte_mconfig.PipelineD_LiUes{
					Imsis:   []string{"IMSI001010000000013"},
					Ips:     []string{"192.16.8.1"},
					Macs:    []string{"00:33:bb:aa:cc:33"},
					Msisdns: []string{"57192831"},
				},
				IpdrExportDst: &lte_mconfig.PipelineD_IPDRExportDst{
					Ip:   "192.168.128.88",
					Port: 2040,
				},
			},
			"sessiond": &lte_mconfig.SessionD{
				LogLevel:         protos.LogLevel_INFO,
				GxGyRelayEnabled: true,
				WalletExhaustDetection: &lte_mconfig.WalletExhaustDetection{
					TerminateOnExhaust: true,
					Method:             lte_mconfig.WalletExhaustDetection_GxTrackedRules,
				},
			},
			"redirectd": &lte_mconfig.RedirectD{
				LogLevel: protos.LogLevel_INFO,
			},
			"directoryd": &orc8r_mconfig.DirectoryD{
				LogLevel: protos.LogLevel_INFO,
			},
			"health": &cwf_mconfig.CwfGatewayHealthConfig{
				CpuUtilThresholdPct: 0,
				MemUtilThresholdPct: 0,
				GreProbeInterval:    0,
				IcmpProbePktCount:   0,
				GrePeers: []*cwf_mconfig.CwfGatewayHealthConfigGrePeer{
					{Ip: "1.2.3.4/24"},
					{Ip: "1.1.1.1/24"},
				},
				ClusterVirtualIp: "10.10.10.11",
			},
		}

		actual, err := build(&nw, &graph, "gw1")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
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

	builder := mconfig.NewRemoteBuilder(cwf_service.ServiceName)
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

var defaultnwConfig = &models.NetworkCarrierWifiConfigs{
	EapAka: &fegmodels.EapAka{
		Timeout: &fegmodels.EapAkaTimeouts{
			ChallengeMs:            20000,
			ErrorNotificationMs:    10000,
			SessionMs:              43200000,
			SessionAuthenticatedMs: 5000,
		},
		PlmnIds: nil,
	},
	AaaServer: &fegmodels.AaaServer{
		IdleSessionTimeoutMs: 21600000,
		AccountingEnabled:    false,
		CreateSessionOnAuth:  false,
	},
	NetworkServices: []string{"dpi", "policy_enforcement"},
	LiUes: &models.LiUes{
		Imsis:   []string{"IMSI001010000000013"},
		Ips:     []string{"192.16.8.1"},
		Macs:    []string{"00:33:bb:aa:cc:33"},
		Msisdns: []string{"57192831"},
	},
	DefaultRuleID: swag.String(""),
}

var defaultgwConfig = &models.GatewayCwfConfigs{
	AllowedGrePeers: models.AllowedGrePeers{
		{IP: "1.2.3.4/24"},
		{IP: "1.1.1.1/24", Key: swag.Uint32(111)},
	},
	IpdrExportDst: &models.IpdrExportDst{
		IP:   "192.168.128.88",
		Port: 2040,
	},
}
