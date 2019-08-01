/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package pluginimpl_test

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/configurator"
	magmad_models "magma/orc8r/cloud/go/services/magmad/obsidian/models"
	upgrade_models "magma/orc8r/cloud/go/services/upgrade/obsidian/models"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBaseOrchestratorMconfigBuilder_Build(t *testing.T) {
	nw := configurator.Network{ID: "n1"}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType,
		Key:  "gw1",
		Config: &magmad_models.MagmadGatewayConfig{
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
	builder := &pluginimpl.BaseOrchestratorMconfigBuilder{}
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
		Config: &upgrade_models.Tier{
			Name:    "default",
			Version: "1.0.0-0",
			Images: []*upgrade_models.TierImagesItems0{
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

func TestDnsdMconfigBuilder_Build(t *testing.T) {
	nw := configurator.Network{ID: "n1"}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType,
		Key:  "gw1",
		Config: &magmad_models.MagmadGatewayConfig{
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

	actual := map[string]proto.Message{}
	builder := &pluginimpl.DnsdMconfigBuilder{}
	err := builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	expected := map[string]proto.Message{
		"dnsd": &mconfig.DnsD{},
	}
	assert.Equal(t, expected, actual)

	nw.Configs = map[string]interface{}{
		"dnsd_network": &models.NetworkDNSConfig{
			EnableCaching: swag.Bool(true),
			LocalTTL:      swag.Uint32(100),
			Records: []*models.DNSConfigRecord{
				{
					ARecord:     []string{"hello", "world"},
					AaaaRecord:  []string{"foo", "bar"},
					CnameRecord: []string{"baz"},
					Domain:      "facebook.com",
				},
				{
					ARecord: []string{"quz"},
				},
			},
		},
	}

	actual = map[string]proto.Message{}
	builder = &pluginimpl.DnsdMconfigBuilder{}
	err = builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	expected = map[string]proto.Message{
		"dnsd": &mconfig.DnsD{
			LogLevel:      protos.LogLevel_INFO,
			EnableCaching: true,
			LocalTTL:      100,
			Records: []*mconfig.NetworkDNSConfigRecordsItems{
				{
					ARecord:     []string{"hello", "world"},
					AaaaRecord:  []string{"foo", "bar"},
					CnameRecord: []string{"baz"},
					Domain:      "facebook.com",
				},
				{
					ARecord: []string{"quz"},
				},
			},
		},
	}
	assert.Equal(t, expected, actual)
}
