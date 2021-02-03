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

	"magma/feg/cloud/go/feg"
	feg_mconfig "magma/feg/cloud/go/protos/mconfig"
	"magma/feg/cloud/go/serdes"
	feg_service "magma/feg/cloud/go/services/feg"
	"magma/feg/cloud/go/services/feg/obsidian/models"
	feg_test_init "magma/feg/cloud/go/services/feg/test_init"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	feg_test_init.StartTestService(t)

	// Empty case: no feg associated to magmad gateway
	nw := configurator.Network{ID: "n1"}
	gw := configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "gw1"}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{gw},
	}

	expected := map[string]proto.Message{}
	actual, err := build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// No GW config but network config exists
	nw.Configs = map[string]interface{}{
		feg.FegNetworkType: defaultConfig,
	}
	fegw := configurator.NetworkEntity{
		Type:               feg.FegGatewayType,
		Key:                "gw1",
		ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "gw1"}},
	}
	gw.Associations = []storage.TypeAndKey{{Type: feg.FegGatewayType, Key: "gw1"}}
	graph = configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{gw, fegw},
		Edges: []configurator.GraphEdge{
			{From: storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "gw1"}, To: storage.TypeAndKey{Type: feg.FegGatewayType, Key: "gw1"}},
		},
	}

	expected = map[string]proto.Message{
		"s6a_proxy": &feg_mconfig.S6AConfig{
			LogLevel: 1,
			Server: &feg_mconfig.DiamClientConfig{
				Protocol:         "sctp",
				Address:          "",
				Retransmits:      0x3,
				WatchdogInterval: 0x1,
				RetryCount:       0x5,
				ProductName:      "magma",
				Realm:            "magma.com",
				Host:             "magma-fedgw.magma.com",
			},
			PlmnIds:                 []string{"123456"},
			RequestFailureThreshold: 0.50,
			MinimumRequestThreshold: 1,
		},
		"s8_proxy": &feg_mconfig.S8Config{
			LogLevel:     1,
			LocalAddress: "10.0.0.1",
			PgwAddress:   "10.0.0.2",
		},
		"hss": &feg_mconfig.HSSConfig{
			Server: &feg_mconfig.DiamServerConfig{
				Protocol:  "tcp",
				DestHost:  "magma.com",
				DestRealm: "magma.com",
			},
			LteAuthOp:  []byte("EREREREREREREREREREREQ=="),
			LteAuthAmf: []byte("gA"),
			DefaultSubProfile: &feg_mconfig.HSSConfig_SubscriptionProfile{
				MaxUlBitRate: 100000000, // 100 Mbps
				MaxDlBitRate: 200000000, // 200 Mbps
			},
			SubProfiles:       nil,
			StreamSubscribers: false,
		},
		"session_proxy": &feg_mconfig.SessionProxyConfig{
			LogLevel: 1,
			Gx: &feg_mconfig.GxConfig{
				DisableGx: true,
				Server: &feg_mconfig.DiamClientConfig{
					Protocol:         "tcp",
					Address:          "",
					Retransmits:      0x3,
					WatchdogInterval: 0x1,
					RetryCount:       0x5,
					ProductName:      "magma",
					Realm:            "magma.com",
					Host:             "magma-fedgw.magma.com",
				},
				// Expect 2, one coming from server and one from serverS
				Servers: []*feg_mconfig.DiamClientConfig{
					{
						Protocol:         "tcp",
						Address:          "",
						Retransmits:      0x3,
						WatchdogInterval: 0x1,
						RetryCount:       0x5,
						ProductName:      "magma",
						Realm:            "magma.com",
						Host:             "magma-fedgw.magma.com",
					},
					{
						Protocol:         "tcp",
						Address:          "",
						Retransmits:      0x3,
						WatchdogInterval: 0x1,
						RetryCount:       0x5,
						ProductName:      "gx.magma",
						Realm:            "gx.magma.com",
						Host:             "magma-fedgw.magma.com",
					},
				},
				OverwriteApn: "apnGx.magma-fedgw.magma.com",
				VirtualApnRules: []*feg_mconfig.VirtualApnRule{
					{
						ApnFilter:                     ".*",
						ChargingCharacteristicsFilter: "1*",
						ApnOverwrite:                  "vApnGx.magma-fedgw.magma.com",
					},
				},
			},
			Gy: &feg_mconfig.GyConfig{
				DisableGy: true,
				Server: &feg_mconfig.DiamClientConfig{
					Protocol:         "tcp",
					Address:          "",
					Retransmits:      0x3,
					WatchdogInterval: 0x1,
					RetryCount:       0x5,
					ProductName:      "magma",
					Realm:            "magma.com",
					Host:             "magma-fedgw.magma.com",
				},
				// Expect 2, one coming from server and one from serverS
				Servers: []*feg_mconfig.DiamClientConfig{
					{
						Protocol:         "tcp",
						Address:          "",
						Retransmits:      0x3,
						WatchdogInterval: 0x1,
						RetryCount:       0x5,
						ProductName:      "magma",
						Realm:            "magma.com",
						Host:             "magma-fedgw.magma.com",
					},
					{
						Protocol:         "tcp",
						Address:          "",
						Retransmits:      0x3,
						WatchdogInterval: 0x1,
						RetryCount:       0x5,
						ProductName:      "gy.magma",
						Realm:            "gy.magma.com",
						Host:             "magma-fedgw.magma.com",
					},
				},
				InitMethod:   feg_mconfig.GyInitMethod_PER_SESSION,
				OverwriteApn: "apnGy.magma-fedgw.magma.com",
				VirtualApnRules: []*feg_mconfig.VirtualApnRule{
					{
						ApnFilter:                     ".*",
						ChargingCharacteristicsFilter: "1*",
						ApnOverwrite:                  "vApnGy.magma-fedgw.magma.com",
					},
				},
			},
			RequestFailureThreshold: 0.50,
			MinimumRequestThreshold: 1,
		},
		"swx_proxy": &feg_mconfig.SwxConfig{
			LogLevel: 1,
			Server: &feg_mconfig.DiamClientConfig{
				Protocol:         "sctp",
				Address:          "",
				Retransmits:      0x3,
				WatchdogInterval: 0x1,
				RetryCount:       0x5,
				ProductName:      "magma",
				Realm:            "magma.com",
				Host:             "magma-fedgw.magma.com",
			},
			// Expect 2, one coming from server and one from serverS
			Servers: []*feg_mconfig.DiamClientConfig{
				{
					Protocol:         "sctp",
					Address:          "",
					Retransmits:      0x3,
					WatchdogInterval: 0x1,
					RetryCount:       0x5,
					ProductName:      "magma",
					Realm:            "magma.com",
					Host:             "magma-fedgw.magma.com",
				},
				{
					Protocol:         "sctp",
					Address:          "",
					Retransmits:      0x3,
					WatchdogInterval: 0x1,
					RetryCount:       0x5,
					ProductName:      "swx1.magma",
					Realm:            "swx1.magma.com",
					Host:             "magma-fedgw.magma.com",
				},
			},
			VerifyAuthorization: false,
			CacheTTLSeconds:     10800,
		},
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
		"health": &feg_mconfig.GatewayHealthConfig{
			RequiredServices:          []string{"SWX_PROXY", "SESSION_PROXY"},
			UpdateIntervalSecs:        10,
			UpdateFailureThreshold:    3,
			CloudDisconnectPeriodSecs: 10,
			LocalDisconnectPeriodSecs: 1,
		},

		"csfb": &feg_mconfig.CsfbConfig{LogLevel: 1,
			Client: &feg_mconfig.SCTPClientConfig{
				LocalAddress:  "",
				ServerAddress: ""},
		},
	}

	actual, err = build(&nw, &graph, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// Put a config on the gw, erase the network config
	nw.Configs = map[string]interface{}{}
	fegw.Config = (*models.GatewayFederationConfigs)(defaultConfig)
	graph.Entities = []configurator.NetworkEntity{fegw, gw}

	actual, err = build(&nw, &graph, "gw1")
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

	builder := mconfig.NewRemoteBuilder(feg_service.ServiceName)
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

func uint32Ptr(i uint32) *uint32 {
	return &i
}

var defaultConfig = &models.NetworkFederationConfigs{
	S6a: &models.S6a{
		Server: &models.DiameterClientConfigs{
			Protocol:         "sctp",
			Retransmits:      3,
			WatchdogInterval: 1,
			RetryCount:       5,
			ProductName:      "magma",
			Host:             "magma-fedgw.magma.com",
			Realm:            "magma.com",
		},
		PlmnIds: []string{"123456"},
	},
	S8: &models.S8{
		LocalAddress: "10.0.0.1",
		PgwAddress:   "10.0.0.2",
	},
	Gx: &models.Gx{
		DisableGx: swag.Bool(true),
		Server: &models.DiameterClientConfigs{
			Protocol:         "tcp",
			Retransmits:      3,
			WatchdogInterval: 1,
			RetryCount:       5,
			ProductName:      "magma",
			Host:             "magma-fedgw.magma.com",
			Realm:            "magma.com",
		},
		Servers: []*models.DiameterClientConfigs{
			{
				Protocol:         "tcp",
				Retransmits:      3,
				WatchdogInterval: 1,
				RetryCount:       5,
				ProductName:      "gx.magma",
				Host:             "magma-fedgw.magma.com",
				Realm:            "gx.magma.com",
			},
		},
		OverwriteApn: "apnGx.magma-fedgw.magma.com",
		VirtualApnRules: []*models.VirtualApnRule{
			{
				ApnFilter:                     ".*",
				ChargingCharacteristicsFilter: "1*",
				ApnOverwrite:                  "vApnGx.magma-fedgw.magma.com",
			},
		},
	},
	Gy: &models.Gy{
		DisableGy: swag.Bool(true),
		Server: &models.DiameterClientConfigs{
			Protocol:         "tcp",
			Retransmits:      3,
			WatchdogInterval: 1,
			RetryCount:       5,
			ProductName:      "magma",
			Host:             "magma-fedgw.magma.com",
			Realm:            "magma.com",
		},
		Servers: []*models.DiameterClientConfigs{
			{
				Protocol:         "tcp",
				Retransmits:      3,
				WatchdogInterval: 1,
				RetryCount:       5,
				ProductName:      "gy.magma",
				Host:             "magma-fedgw.magma.com",
				Realm:            "gy.magma.com",
			},
		},
		InitMethod:   uint32Ptr(1),
		OverwriteApn: "apnGy.magma-fedgw.magma.com",
		VirtualApnRules: []*models.VirtualApnRule{
			{
				ApnFilter:                     ".*",
				ChargingCharacteristicsFilter: "1*",
				ApnOverwrite:                  "vApnGy.magma-fedgw.magma.com",
			},
		},
	},
	Hss: &models.Hss{
		Server: &models.DiameterServerConfigs{
			Protocol:  "tcp",
			DestHost:  "magma.com",
			DestRealm: "magma.com",
		},
		LteAuthOp:  []byte("EREREREREREREREREREREQ=="),
		LteAuthAmf: []byte("gA"),
		DefaultSubProfile: &models.SubscriptionProfile{
			MaxUlBitRate: 100000000, // 100 Mbps
			MaxDlBitRate: 200000000, // 200 Mbps
		},
		SubProfiles:       nil,
		StreamSubscribers: false,
	},
	Swx: &models.Swx{
		Server: &models.DiameterClientConfigs{
			Protocol:         "sctp",
			Retransmits:      3,
			WatchdogInterval: 1,
			RetryCount:       5,
			ProductName:      "magma",
			Host:             "magma-fedgw.magma.com",
			Realm:            "magma.com",
		},
		Servers: []*models.DiameterClientConfigs{
			{
				Protocol:         "sctp",
				Retransmits:      3,
				WatchdogInterval: 1,
				RetryCount:       5,
				ProductName:      "swx1.magma",
				Host:             "magma-fedgw.magma.com",
				Realm:            "swx1.magma.com",
			},
		},
		VerifyAuthorization: false,
		CacheTTLSeconds:     10800,
	},
	EapAka: &models.EapAka{
		Timeout: &models.EapAkaTimeouts{
			ChallengeMs:            20000,
			ErrorNotificationMs:    10000,
			SessionMs:              43200000,
			SessionAuthenticatedMs: 5000,
		},
		PlmnIds: []string{},
	},
	AaaServer: &models.AaaServer{
		IdleSessionTimeoutMs: 21600000,
		AccountingEnabled:    false,
		CreateSessionOnAuth:  false,
	},
	ServedNetworkIds: []string{},
	Health: &models.Health{
		HealthServices:           []string{"SWX_PROXY", "SESSION_PROXY"},
		UpdateIntervalSecs:       10,
		CloudDisablePeriodSecs:   10,
		LocalDisablePeriodSecs:   1,
		UpdateFailureThreshold:   3,
		RequestFailureThreshold:  0.50,
		MinimumRequestThreshold:  1,
		CPUUtilizationThreshold:  0.75,
		MemoryAvailableThreshold: 0.90,
	},
	Csfb: &models.Csfb{
		Client: &models.SctpClientConfigs{
			LocalAddress:  "",
			ServerAddress: "",
		},
	},
}
