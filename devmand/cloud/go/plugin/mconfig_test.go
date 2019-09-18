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
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	"orc8r/devmand/cloud/go/devmand"
	"orc8r/devmand/cloud/go/plugin"
	"orc8r/devmand/cloud/go/protos/mconfig"
	"orc8r/devmand/cloud/go/services/devmand/obsidian/models"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	nw := configurator.Network{ID: "n1"}
	devmandDevice1 := configurator.NetworkEntity{
		Type:   devmand.DeviceType,
		Key:    "device1",
		Config: models.NewDefaultManagedDeviceModel(),
	}
	devmandGW := configurator.NetworkEntity{
		Type:         devmand.DevmandGatewayType,
		Key:          "gw1",
		Config:       &models.GatewayDevmandConfigs{},
		Associations: []storage.TypeAndKey{devmandDevice1.GetTypeAndKey()},
	}
	magmadGW := configurator.NetworkEntity{
		Type:         orc8r.MagmadGatewayType,
		Key:          "gw1",
		Associations: []storage.TypeAndKey{devmandGW.GetTypeAndKey()},
	}

	devmandDevice1.ParentAssociations = []storage.TypeAndKey{devmandGW.GetTypeAndKey()}
	devmandGW.ParentAssociations = []storage.TypeAndKey{magmadGW.GetTypeAndKey()}

	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{magmadGW, devmandGW, devmandDevice1},
		Edges: []configurator.GraphEdge{
			{From: magmadGW.GetTypeAndKey(), To: devmandGW.GetTypeAndKey()},
			{From: devmandGW.GetTypeAndKey(), To: devmandDevice1.GetTypeAndKey()},
		},
	}

	actual := map[string]proto.Message{}
	builder := plugin.Builder{}
	err := builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)

	expected := map[string]proto.Message{
		"devmand": &mconfig.DevmandGatewayConfig{
			ManagedDevices: map[string]*mconfig.ManagedDevice{
				"device1": {
					DeviceConfig: "config_json",
					DeviceType:   []string{"type_descriptor_1"},
					Channels: &mconfig.Channels{
						FrinxChannel:   &mconfig.FrinxChannel{},
						CambiumChannel: &mconfig.CambiumChannel{},
						OtherChannel:   &mconfig.OtherChannel{ChannelProps: map[string]string{}},
						SnmpChannel:    &mconfig.SNMPChannel{},
					},
					Host:     "hostname",
					Platform: "platform_name",
				},
			},
		},
	}
	assert.Equal(t, expected, actual)
}
