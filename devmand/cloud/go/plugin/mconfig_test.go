/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package plugin_test

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	orc8rplugin "magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	"orc8r/devmand/cloud/go/devmand"
	"orc8r/devmand/cloud/go/plugin"
	models2 "orc8r/devmand/cloud/go/plugin/models"
	"orc8r/devmand/cloud/go/protos/mconfig"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	orc8rplugin.RegisterPluginForTests(t, &plugin.DevmandOrchestratorPlugin{})

	nID := "n1"
	nw := configurator.Network{ID: nID}
	device := configurator.NetworkEntity{
		Type: devmand.SymphonyDeviceType, Key: "d1",
		Config:             models2.NewDefaultSymphonyDeviceConfig(),
		ParentAssociations: []storage.TypeAndKey{storage.TypeAndKey{Type: devmand.SymphonyAgentType, Key: "a1"}},
	}
	agent := configurator.NetworkEntity{
		Type: devmand.SymphonyAgentType, Key: "a1",
		Associations: []storage.TypeAndKey{
			{Type: devmand.SymphonyDeviceType, Key: "d1"},
		},
		ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "a1"}},
	}
	gateway := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "a1",
		Associations: []storage.TypeAndKey{{Type: devmand.SymphonyAgentType, Key: "a1"}},
	}

	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{gateway, agent, device},
		Edges: []configurator.GraphEdge{
			{From: gateway.GetTypeAndKey(), To: agent.GetTypeAndKey()},
			{From: agent.GetTypeAndKey(), To: device.GetTypeAndKey()},
		},
	}

	actual := map[string]proto.Message{}
	builder := plugin.Builder{}
	err := builder.Build("n1", "a1", graph, nw, actual)
	assert.NoError(t, err)

	expected := map[string]proto.Message{
		"devmand": &mconfig.DevmandGatewayConfig{
			ManagedDevices: map[string]*mconfig.ManagedDevice{
				"d1": {
					DeviceConfig: "{}",
					DeviceType:   []string{"device_type 1", "device_type 2"},
					Channels: &mconfig.Channels{
						SnmpChannel: &mconfig.SNMPChannel{
							Community: "snmp community",
							Version:   "1",
						},
					},
					Host:     "device_host",
					Platform: "device_platform",
				},
			},
		},
	}
	assert.Equal(t, expected, actual)
}
