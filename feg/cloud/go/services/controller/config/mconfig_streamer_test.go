/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config_test

import (
	"testing"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/cloud/go/services/controller/config"
	configprotos "magma/feg/cloud/go/services/controller/protos"
	orc8rprotos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/streaming"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestFegStreamer_ApplyMconfigUpdate(t *testing.T) {
	s := &config.FegStreamer{}

	inputMconfigs := map[string]*orc8rprotos.GatewayConfigs{
		"gw1": {ConfigsByKey: map[string]*any.Any{}},
		"gw2": {ConfigsByKey: map[string]*any.Any{}},
	}
	update := &streaming.ConfigUpdate{
		NetworkId:  "network",
		ConfigType: config.FegNetworkType,
		ConfigKey:  "network",
		NewValue:   configprotos.NewDefaultNetworkConfig(),
		Operation:  streaming.CreateOperation,
	}

	_, err := s.ApplyMconfigUpdate(update, inputMconfigs)
	assert.NoError(t, err)

	expected := map[string]proto.Message{
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
		},
	}
	expectedMconfig := getExpectedMconfig(t, expected)
	assert.Equal(t, map[string]*orc8rprotos.GatewayConfigs{"gw1": expectedMconfig, "gw2": expectedMconfig}, inputMconfigs)
}

func getExpectedMconfig(t *testing.T, expected map[string]proto.Message) *orc8rprotos.GatewayConfigs {
	ret := &orc8rprotos.GatewayConfigs{ConfigsByKey: map[string]*any.Any{}}
	for k, v := range expected {
		vAny, err := ptypes.MarshalAny(v)
		assert.NoError(t, err)
		ret.ConfigsByKey[k] = vAny
	}
	return ret
}
