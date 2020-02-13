/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin_test

import (
	"testing"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/plugin"
	"magma/feg/cloud/go/plugin/models"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	builder := &plugin.Builder{}

	// empty case: no feg associated to magmad gateway
	nw := configurator.Network{ID: "n1"}
	gw := configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "gw1"}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{gw},
	}

	actual := map[string]proto.Message{}
	expected := map[string]proto.Message{}
	err := builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// no GW config but network config exists
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
	actual = map[string]proto.Message{}
	expected = map[string]proto.Message{
		"s6a_proxy": &mconfig.S6AConfig{
			LogLevel: 1,
			Server: &mconfig.DiamClientConfig{
				Protocol:         "sctp",
				Address:          "",
				Retransmits:      0x3,
				WatchdogInterval: 0x1,
				RetryCount:       0x5,
				ProductName:      "magma",
				Realm:            "magma.com",
				Host:             "magma-fedgw.magma.com",
			},
			RequestFailureThreshold: 0.50,
			MinimumRequestThreshold: 1,
		},
		"hss": &mconfig.HSSConfig{
			Server: &mconfig.DiamServerConfig{
				Protocol:  "tcp",
				DestHost:  "magma.com",
				DestRealm: "magma.com",
			},
			LteAuthOp:  []byte("EREREREREREREREREREREQ=="),
			LteAuthAmf: []byte("gA"),
			DefaultSubProfile: &mconfig.HSSConfig_SubscriptionProfile{
				MaxUlBitRate: 100000000, // 100 Mbps
				MaxDlBitRate: 200000000, // 200 Mbps
			},
			SubProfiles:       make(map[string]*mconfig.HSSConfig_SubscriptionProfile),
			StreamSubscribers: false,
		},
		"session_proxy": &mconfig.SessionProxyConfig{
			LogLevel: 1,
			Gx: &mconfig.GxConfig{
				Server: &mconfig.DiamClientConfig{
					Protocol:         "tcp",
					Address:          "",
					Retransmits:      0x3,
					WatchdogInterval: 0x1,
					RetryCount:       0x5,
					ProductName:      "magma",
					Realm:            "magma.com",
					Host:             "magma-fedgw.magma.com",
				},
			},
			Gy: &mconfig.GyConfig{
				Server: &mconfig.DiamClientConfig{
					Protocol:         "tcp",
					Address:          "",
					Retransmits:      0x3,
					WatchdogInterval: 0x1,
					RetryCount:       0x5,
					ProductName:      "magma",
					Realm:            "magma.com",
					Host:             "magma-fedgw.magma.com",
				},
				InitMethod: mconfig.GyInitMethod_PER_SESSION,
			},
			RequestFailureThreshold: 0.50,
			MinimumRequestThreshold: 1,
		},
		"swx_proxy": &mconfig.SwxConfig{
			LogLevel: 1,
			Server: &mconfig.DiamClientConfig{
				Protocol:         "sctp",
				Address:          "",
				Retransmits:      0x3,
				WatchdogInterval: 0x1,
				RetryCount:       0x5,
				ProductName:      "magma",
				Realm:            "magma.com",
				Host:             "magma-fedgw.magma.com",
			},
			VerifyAuthorization: false,
			CacheTTLSeconds:     10800,
		},
		"eap_aka": &mconfig.EapAkaConfig{LogLevel: 1,
			Timeout: &mconfig.EapAkaConfig_Timeouts{
				ChallengeMs:            20000,
				ErrorNotificationMs:    10000,
				SessionMs:              43200000,
				SessionAuthenticatedMs: 5000,
			},
			PlmnIds: []string{},
		},
		"aaa_server": &mconfig.AAAConfig{LogLevel: 1,
			IdleSessionTimeoutMs: 21600000,
			AccountingEnabled:    false,
			CreateSessionOnAuth:  false,
		},
		"health": &mconfig.GatewayHealthConfig{
			RequiredServices:          []string{"SWX_PROXY", "SESSION_PROXY"},
			UpdateIntervalSecs:        10,
			UpdateFailureThreshold:    3,
			CloudDisconnectPeriodSecs: 10,
			LocalDisconnectPeriodSecs: 1,
		},

		"csfb": &mconfig.CsfbConfig{LogLevel: 1,
			Client: &mconfig.SCTPClientConfig{
				LocalAddress:  "",
				ServerAddress: ""},
		},
	}

	err = builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// put a config on the gw, erase the network config
	nw.Configs = map[string]interface{}{}
	fegw.Config = (*models.GatewayFederationConfigs)(defaultConfig)
	graph.Entities = []configurator.NetworkEntity{fegw, gw}

	actual = map[string]proto.Message{}
	err = builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
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
	},
	Gx: &models.Gx{
		Server: &models.DiameterClientConfigs{
			Protocol:         "tcp",
			Retransmits:      3,
			WatchdogInterval: 1,
			RetryCount:       5,
			ProductName:      "magma",
			Host:             "magma-fedgw.magma.com",
			Realm:            "magma.com",
		},
	},
	Gy: &models.Gy{
		Server: &models.DiameterClientConfigs{
			Protocol:         "tcp",
			Retransmits:      3,
			WatchdogInterval: 1,
			RetryCount:       5,
			ProductName:      "magma",
			Host:             "magma-fedgw.magma.com",
			Realm:            "magma.com",
		},
		InitMethod: uint32Ptr(1),
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
		SubProfiles:       make(map[string]models.SubscriptionProfile),
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
		IDLESessionTimeoutMs: 21600000,
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

func uint32Ptr(i uint32) *uint32 {
	return &i
}
