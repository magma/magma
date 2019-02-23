/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config_test

import (
	"testing"

	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/protos"
	mconfig_protos "magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/config"
	config_test_init "magma/orc8r/cloud/go/services/config/test_init"
	magmad_config "magma/orc8r/cloud/go/services/magmad/config"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/upgrade"
	upgrade_protos "magma/orc8r/cloud/go/services/upgrade/protos"
	upgrade_test_init "magma/orc8r/cloud/go/services/upgrade/test_init"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestMagmadMconfigBuilder_Build(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	config_test_init.StartTestService(t)
	upgrade_test_init.StartTestService(t)

	builder := &magmad_config.MagmadMconfigBuilder{}
	actual, err := builder.Build("network", "gw1")
	assert.NoError(t, err)
	assert.Equal(t, map[string]proto.Message{}, actual)

	err = config.CreateConfig("network", magmad_config.MagmadGatewayType, "gw1", &magmadprotos.MagmadGatewayConfig{
		AutoupgradeEnabled:      true,
		AutoupgradePollInterval: 300,
		CheckinInterval:         60,
		CheckinTimeout:          10,
		DynamicServices:         []string{},
		Tier:                    "default",
	})
	assert.NoError(t, err)

	// Make sure we get a 0.0.0-0 and no error if default tier DNE
	actual, err = builder.Build("network", "gw1")
	expected := map[string]proto.Message{
		"control_proxy": &mconfig_protos.ControlProxy{LogLevel: protos.LogLevel_INFO},
		"magmad": &mconfig_protos.MagmaD{
			LogLevel:                protos.LogLevel_INFO,
			CheckinInterval:         60,
			CheckinTimeout:          10,
			AutoupgradeEnabled:      true,
			AutoupgradePollInterval: 300,
			PackageVersion:          "0.0.0-0",
			Images:                  []*mconfig_protos.ImageSpec{},
			TierId:                  "default",
			DynamicServices:         []string{},
			FeatureFlags:            map[string]bool{},
		},
		"metricsd": &mconfig_protos.MetricsD{LogLevel: protos.LogLevel_INFO},
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	err = upgrade.CreateTier(
		"network",
		"default",
		&upgrade_protos.TierInfo{
			Name:    "default",
			Version: "1.0.0-0",
			Images: []*upgrade_protos.ImageSpec{
				{Name: "Image1", Order: 42},
				{Name: "Image2", Order: 1},
			},
		},
	)
	assert.NoError(t, err)

	actual, err = builder.Build("network", "gw1")
	expected = map[string]proto.Message{
		"control_proxy": &mconfig_protos.ControlProxy{LogLevel: protos.LogLevel_INFO},
		"magmad": &mconfig_protos.MagmaD{
			LogLevel:                protos.LogLevel_INFO,
			CheckinInterval:         60,
			CheckinTimeout:          10,
			AutoupgradeEnabled:      true,
			AutoupgradePollInterval: 300,
			PackageVersion:          "1.0.0-0",
			Images: []*mconfig_protos.ImageSpec{
				{Name: "Image1", Order: 42},
				{Name: "Image2", Order: 1},
			},
			TierId:          "default",
			DynamicServices: []string{},
			FeatureFlags:    map[string]bool{},
		},
		"metricsd": &mconfig_protos.MetricsD{LogLevel: protos.LogLevel_INFO},
	}
	assert.Equal(t, expected, actual)
}
