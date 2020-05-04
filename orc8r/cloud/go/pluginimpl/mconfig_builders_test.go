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
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/protos/mconfig"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBaseOrchestratorMconfigBuilder_Build(t *testing.T) {
	nw := configurator.Network{ID: "n1"}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType,
		Key:  "gw1",
		Config: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
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
	magmadSerde := configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models.MagmadGatewayConfigs{})
	upgradeReleaseSerde := configurator.NewNetworkEntityConfigSerde(orc8r.UpgradeReleaseChannelEntityType, &models.ReleaseChannel{})
	upgradeTierSerde := configurator.NewNetworkEntityConfigSerde(orc8r.UpgradeTierEntityType, &models.Tier{})
	err := serde.RegisterSerdes(magmadSerde, upgradeReleaseSerde, upgradeTierSerde)
	assert.NoError(t, err)

	// Make sure we get a 0.0.0-0 and no error if no tier
	builder := &pluginimpl.BaseOrchestratorMconfigBuilder{}
	actual := map[string]proto.Message{}
	err = builder.Build("n1", "gw1", graph, nw, actual)
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
			Images:                  nil,
			DynamicServices:         nil,
			FeatureFlags:            nil,
		},
		"metricsd": &mconfig.MetricsD{LogLevel: protos.LogLevel_INFO},
		"td-agent-bit": &mconfig.FluentBit{
			ExtraTags:        map[string]string{"network_id": "n1", "gateway_id": "gw1"},
			ThrottleRate:     1000,
			ThrottleWindow:   5,
			ThrottleInterval: "1m",
		},
	}
	assert.Equal(t, expected, actual)
	// Put a tier in the graph
	tier := configurator.NetworkEntity{
		Type: orc8r.UpgradeTierEntityType,
		Key:  "default",
		Config: &models.Tier{
			Name:    "default",
			Version: "1.0.0-0",
			Images: []*models.TierImage{
				{Name: swag.String("Image1"), Order: swag.Int64(42)},
				{Name: swag.String("Image2"), Order: swag.Int64(1)},
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
			DynamicServices: nil,
			FeatureFlags:    nil,
		},
		"metricsd": &mconfig.MetricsD{LogLevel: protos.LogLevel_INFO},
		"td-agent-bit": &mconfig.FluentBit{
			ExtraTags:        map[string]string{"network_id": "n1", "gateway_id": "gw1"},
			ThrottleRate:     1000,
			ThrottleWindow:   5,
			ThrottleInterval: "1m",
		},
	}
	assert.Equal(t, expected, actual)

	// Set list of files for log aggregation
	testThrottleInterval := "30h"
	testThrottleWindow := uint32(808)
	testThrottleRate := uint32(305)
	gw.Config = &models.MagmadGatewayConfigs{
		AutoupgradeEnabled:      swag.Bool(true),
		AutoupgradePollInterval: 300,
		CheckinInterval:         60,
		CheckinTimeout:          10,
		DynamicServices:         nil,
		FeatureFlags:            nil,
		Logging: &models.GatewayLoggingConfigs{
			Aggregation: &models.AggregationLoggingConfigs{
				TargetFilesByTag: map[string]string{
					"thing": "/var/log/thing.log",
					"blah":  "/some/directory/blah.log",
				},
				ThrottleRate:     &testThrottleRate,
				ThrottleWindow:   &testThrottleWindow,
				ThrottleInterval: &testThrottleInterval,
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
			DynamicServices: nil,
			FeatureFlags:    nil,
		},
		"metricsd": &mconfig.MetricsD{LogLevel: protos.LogLevel_INFO},
		"td-agent-bit": &mconfig.FluentBit{
			ExtraTags:        map[string]string{"network_id": "n1", "gateway_id": "gw1"},
			ThrottleRate:     305,
			ThrottleWindow:   808,
			ThrottleInterval: "30h",
			FilesByTag: map[string]string{
				"thing": "/var/log/thing.log",
				"blah":  "/some/directory/blah.log",
			},
		},
	}
	assert.Equal(t, expected, actual)

	// Check default values for log throttling
	gw.Config = &models.MagmadGatewayConfigs{
		AutoupgradeEnabled:      swag.Bool(true),
		AutoupgradePollInterval: 300,
		CheckinInterval:         60,
		CheckinTimeout:          10,
		DynamicServices:         nil,
		FeatureFlags:            nil,
		Logging: &models.GatewayLoggingConfigs{
			Aggregation: &models.AggregationLoggingConfigs{
				TargetFilesByTag: map[string]string{
					"thing": "/var/log/thing.log",
					"blah":  "/some/directory/blah.log",
				},
				// No throttle values
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
			DynamicServices: nil,
			FeatureFlags:    nil,
		},
		"metricsd": &mconfig.MetricsD{LogLevel: protos.LogLevel_INFO},
		"td-agent-bit": &mconfig.FluentBit{
			ExtraTags:        map[string]string{"network_id": "n1", "gateway_id": "gw1"},
			ThrottleRate:     1000,
			ThrottleWindow:   5,
			ThrottleInterval: "1m",
			FilesByTag: map[string]string{
				"thing": "/var/log/thing.log",
				"blah":  "/some/directory/blah.log",
			},
		},
	}
	assert.Equal(t, expected, actual)
}

func TestDnsdMconfigBuilder_Build(t *testing.T) {
	dnsdSerde := configurator.NewNetworkConfigSerde(orc8r.DnsdNetworkType, &models.NetworkDNSConfig{})
	err := serde.RegisterSerdes(dnsdSerde)
	assert.NoError(t, err)

	nw := configurator.Network{ID: "n1"}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType,
		Key:  "gw1",
		Config: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
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
	err = builder.Build("n1", "gw1", graph, nw, actual)
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
					ARecord:     []strfmt.IPv4{"127.0.0.1", "127.0.0.2"},
					AaaaRecord:  []strfmt.IPv6{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", "1234:0db8:85a3:0000:0000:8a2e:0370:1234"},
					CnameRecord: []string{"baz"},
					Domain:      "facebook.com",
				},
				{
					ARecord: []strfmt.IPv4{"quz"},
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
					ARecord:     []string{"127.0.0.1", "127.0.0.2"},
					AaaaRecord:  []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", "1234:0db8:85a3:0000:0000:8a2e:0370:1234"},
					CnameRecord: []string{"baz"},
					Domain:      "facebook.com",
				},
				{
					ARecord: []string{"quz"},
				},
			},
		},
	}
	assert.Equal(t, expected["dnsd"].String(), actual["dnsd"].String())
}
