/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package pluginimpl

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/magmad/obsidian/models"
	models2 "magma/orc8r/cloud/go/services/upgrade/obsidian/models"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBaseOrchestratorMconfigBuilder_Build(t *testing.T) {
	test_init.StartTestService(t)

	nw := configurator.Network{ID: "n1"}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType,
		Key:  "gw1",
		Config: &models.MagmadGatewayConfig{
			AutoupgradeEnabled:      true,
			AutoupgradePollInterval: 300,
			CheckinInterval:         60,
			CheckinTimeout:          10,
			DynamicServices:         []string{},
			FeatureFlags:            map[string]bool{},
		},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{gw},
	}

	// Make sure we get a 0.0.0-0 and no error if no tier
	builder := &BaseOrchestratorMconfigBuilder{}
	actual := map[string]proto.Message{}
	err := builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
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
			DynamicServices:         []string{},
			FeatureFlags:            map[string]bool{},
		},
		"metricsd": &mconfig.MetricsD{LogLevel: protos.LogLevel_INFO},
	}
	assert.Equal(t, expected, actual)

	// Put a tier in the graph
	tier := configurator.NetworkEntity{
		Type: orc8r.UpgradeTierEntityType,
		Key:  "default",
		Config: &models2.Tier{
			Name:    "default",
			Version: "1.0.0-0",
			Images: []*models2.TierImagesItems0{
				{Name: "Image1", Order: 42},
				{Name: "Image2", Order: 1},
			},
		},
	}
	graph = configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{gw, tier},
		Edges: []configurator.GraphEdge{
			{From: tier.GetTypeAndKey(), To: gw.GetTypeAndKey()},
		},
	}
	actual = map[string]proto.Message{}
	err = builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
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
			DynamicServices: []string{},
			FeatureFlags:    map[string]bool{},
		},
		"metricsd": &mconfig.MetricsD{LogLevel: protos.LogLevel_INFO},
	}
	assert.Equal(t, expected, actual)
}
