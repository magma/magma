/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config_test

import (
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/config/streaming"
	"magma/orc8r/cloud/go/services/magmad/config"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/upgrade"
	upgrade_protos "magma/orc8r/cloud/go/services/upgrade/protos"
	upgrade_test_init "magma/orc8r/cloud/go/services/upgrade/test_init"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestMagmadStreamer_ApplyMconfigUpdate(t *testing.T) {
	upgrade_test_init.StartTestService(t)

	s := &config.MagmadStreamer{}

	inputMconfigs := map[string]*protos.GatewayConfigs{
		"gw1": {ConfigsByKey: map[string]*any.Any{}},
		"gw2": {ConfigsByKey: map[string]*any.Any{}},
	}
	update := &streaming.ConfigUpdate{
		NetworkId:  "network",
		ConfigType: config.MagmadGatewayType,
		ConfigKey:  "gw1",
		NewValue: &magmadprotos.MagmadGatewayConfig{
			AutoupgradeEnabled:      true,
			AutoupgradePollInterval: 300,
			CheckinInterval:         60,
			CheckinTimeout:          10,
			Tier:                    "test_tier_1",
		},
		Operation: streaming.CreateOperation,
	}

	_, err := s.ApplyMconfigUpdate(update, inputMconfigs)
	assert.NoError(t, err)

	// Make sure we get a 0.0.0-0 if default tier DNE
	expected := map[string]proto.Message{
		"control_proxy": &mconfig.ControlProxy{LogLevel: protos.LogLevel_INFO},
		"magmad": &mconfig.MagmaD{
			LogLevel:                protos.LogLevel_INFO,
			CheckinInterval:         60,
			CheckinTimeout:          10,
			AutoupgradeEnabled:      true,
			AutoupgradePollInterval: 300,
			PackageVersion:          "0.0.0-0",
			Images:                  []*mconfig.ImageSpec{},
			TierId:                  "test_tier_1",
		},
		"metricsd": &mconfig.MetricsD{LogLevel: protos.LogLevel_INFO},
	}
	expectedMconfig := getExpectedMconfig(t, expected)
	assert.Equal(t, map[string]*protos.GatewayConfigs{"gw1": expectedMconfig, "gw2": expectedMconfig}, inputMconfigs)

	// Create the default tier
	err = upgrade.CreateTier(
		"network",
		"test_tier_1",
		&upgrade_protos.TierInfo{
			Name:    "test_tier_1",
			Version: "1.0.0-0",
			Images: []*upgrade_protos.ImageSpec{
				{Name: "Image1", Order: 42},
				{Name: "Image2", Order: 1},
			},
		},
	)
	assert.NoError(t, err)

	_, err = s.ApplyMconfigUpdate(update, inputMconfigs)
	assert.NoError(t, err)

	// Make sure we get a 0.0.0-0 if default tier DNE
	expected = map[string]proto.Message{
		"control_proxy": &mconfig.ControlProxy{LogLevel: protos.LogLevel_INFO},
		"magmad": &mconfig.MagmaD{
			LogLevel:                protos.LogLevel_INFO,
			CheckinInterval:         60,
			CheckinTimeout:          10,
			AutoupgradeEnabled:      true,
			AutoupgradePollInterval: 300,
			PackageVersion:          "1.0.0-0",
			Images: []*mconfig.ImageSpec{
				{Name: "Image1", Order: 42},
				{Name: "Image2", Order: 1},
			},
			TierId: "test_tier_1",
		},
		"metricsd": &mconfig.MetricsD{LogLevel: protos.LogLevel_INFO},
	}
	expectedMconfig = getExpectedMconfig(t, expected)
	assert.Equal(t, map[string]*protos.GatewayConfigs{"gw1": expectedMconfig, "gw2": expectedMconfig}, inputMconfigs)

	update = &streaming.ConfigUpdate{
		NetworkId:  "network",
		ConfigType: config.MagmadGatewayType,
		ConfigKey:  "gw1",
		NewValue:   nil,
		Operation:  streaming.DeleteOperation,
	}
	_, err = s.ApplyMconfigUpdate(update, inputMconfigs)
	assert.NoError(t, err)
	expectedMconfig = &protos.GatewayConfigs{ConfigsByKey: map[string]*any.Any{}}
	assert.Equal(t, map[string]*protos.GatewayConfigs{"gw1": expectedMconfig, "gw2": expectedMconfig}, inputMconfigs)
}

func getExpectedMconfig(t *testing.T, expected map[string]proto.Message) *protos.GatewayConfigs {
	ret := &protos.GatewayConfigs{ConfigsByKey: map[string]*any.Any{}}
	for k, v := range expected {
		vAny, err := ptypes.MarshalAny(v)
		assert.NoError(t, err)
		ret.ConfigsByKey[k] = vAny
	}
	return ret
}
