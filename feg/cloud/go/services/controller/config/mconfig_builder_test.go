/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config_test

import (
	"testing"

	fegplugin "magma/feg/cloud/go/plugin"
	"magma/feg/cloud/go/protos/mconfig"
	feg_config "magma/feg/cloud/go/services/controller/config"
	config_protos "magma/feg/cloud/go/services/controller/protos"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/services/config"
	config_test_init "magma/orc8r/cloud/go/services/config/test_init"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestControllerBuilder_Build(t *testing.T) {
	plugin.RegisterPluginForTests(t, &fegplugin.FegOrchestratorPlugin{})
	config_test_init.StartTestService(t)
	builder := &feg_config.Builder{}
	actual, err := builder.Build("network", "feg1")
	assert.NoError(t, err)
	assert.Equal(t, map[string]proto.Message{}, actual)

	defaultNetCfg := config_protos.NewDefaultNetworkConfig()
	err = config.CreateConfig("network", feg_config.FegNetworkType, "network", defaultNetCfg)
	assert.NoError(t, err)

	expected := map[string]proto.Message{
		"s6a_proxy": &mconfig.S6AConfig{
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
			SubProfiles: make(map[string]*mconfig.HSSConfig_SubscriptionProfile),
		},
		"session_proxy": &mconfig.SessionProxyConfig{
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
		},
		"health": &mconfig.GatewayHealthConfig{
			RequiredServices:          []string{"S6A_PROXY", "SESSION_PROXY"},
			UpdateIntervalSecs:        10,
			UpdateFailureThreshold:    3,
			CloudDisconnectPeriodSecs: 10,
			LocalDisconnectPeriodSecs: 1,
		},
	}

	actual, err = builder.Build("network", "feg1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	defaultGwCfg := proto.Clone(defaultNetCfg)
	defaultGwCfg.(*config_protos.Config).S6A.Server.Address = "127.0.0.1:5555"
	defaultGwCfg.(*config_protos.Config).S6A.Server.LocalAddress = ":56789"
	defaultGwCfg.(*config_protos.Config).Hss.Server.Address = "127.0.0.1:5555"
	defaultGwCfg.(*config_protos.Config).Hss.Server.LocalAddress = ":56789"
	defaultGwCfg.(*config_protos.Config).Gy.Server.DestHost = "ocs.mno.com"
	defaultGwCfg.(*config_protos.Config).Gy.Server.DestRealm = "mno.com"
	defaultGwCfg.(*config_protos.Config).Gx.Server.DestHost = "pcrf.mno.com"
	defaultGwCfg.(*config_protos.Config).Gx.Server.DestRealm = "mno.com"
	defaultGwCfg.(*config_protos.Config).Swx.Server.Address = "127.0.0.1:9999"
	defaultGwCfg.(*config_protos.Config).Swx.Server.LocalAddress = ":12123"
	defaultGwCfg.(*config_protos.Config).Health.UpdateFailureThreshold = 4
	expected["s6a_proxy"].(*mconfig.S6AConfig).Server.Address = "127.0.0.1:5555"
	expected["s6a_proxy"].(*mconfig.S6AConfig).Server.LocalAddress = ":56789"
	expected["hss"].(*mconfig.HSSConfig).Server.Address = "127.0.0.1:5555"
	expected["hss"].(*mconfig.HSSConfig).Server.LocalAddress = ":56789"
	expected["session_proxy"].(*mconfig.SessionProxyConfig).Gy.Server.DestHost = "ocs.mno.com"
	expected["session_proxy"].(*mconfig.SessionProxyConfig).Gy.Server.DestRealm = "mno.com"
	expected["session_proxy"].(*mconfig.SessionProxyConfig).Gx.Server.DestHost = "pcrf.mno.com"
	expected["session_proxy"].(*mconfig.SessionProxyConfig).Gx.Server.DestRealm = "mno.com"
	expected["swx_proxy"].(*mconfig.SwxConfig).Server.Address = "127.0.0.1:9999"
	expected["swx_proxy"].(*mconfig.SwxConfig).Server.LocalAddress = ":12123"
	expected["health"].(*mconfig.GatewayHealthConfig).UpdateFailureThreshold = 4

	err = config.CreateConfig("network", feg_config.FegGatewayType, "feg1", defaultGwCfg)
	assert.NoError(t, err)

	actual, err = builder.Build("network", "feg1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	err = config.DeleteConfig("network", feg_config.FegGatewayType, "feg1")
	assert.NoError(t, err)
	expected["s6a_proxy"].(*mconfig.S6AConfig).Server.Address = ""
	expected["s6a_proxy"].(*mconfig.S6AConfig).Server.LocalAddress = ""
	expected["hss"].(*mconfig.HSSConfig).Server.Address = ""
	expected["hss"].(*mconfig.HSSConfig).Server.LocalAddress = ""
	expected["session_proxy"].(*mconfig.SessionProxyConfig).Gy.Server.DestHost = ""
	expected["session_proxy"].(*mconfig.SessionProxyConfig).Gy.Server.DestRealm = ""
	expected["session_proxy"].(*mconfig.SessionProxyConfig).Gx.Server.DestHost = ""
	expected["session_proxy"].(*mconfig.SessionProxyConfig).Gx.Server.DestRealm = ""
	expected["swx_proxy"].(*mconfig.SwxConfig).Server.Address = ""
	expected["swx_proxy"].(*mconfig.SwxConfig).Server.LocalAddress = ""
	expected["health"].(*mconfig.GatewayHealthConfig).UpdateFailureThreshold = 3

	actual, err = builder.Build("network", "feg1")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
