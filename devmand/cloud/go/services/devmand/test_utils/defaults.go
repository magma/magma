/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package test_utils

import (
	"orc8r/devmand/cloud/go/services/devmand/obsidian/models"
	"orc8r/devmand/cloud/go/services/devmand/protos"
)

// NewDefaultProtosManagedDevice generates a basic device config
func NewDefaultProtosManagedDevice() *protos.ManagedDevice {
	deviceConfigJSONString := `{
		"interfaces": {
			"interface": [
				{
					"config": {
						"description": "ACCESS PORT",
						"name": "g0/0/0",
						"type": "ethernetCsmacd"
					},
					"name": "g0/0/0"
				}
			]
		}
	}`
	return &protos.ManagedDevice{
		DeviceConfig: deviceConfigJSONString,
		Host:         "1.2.3.4",
		DeviceType:   []string{"wifi", "access_point"},
		Platform:     "Ubiquiti",
		Channels: &protos.Channels{
			SnmpChannel: &protos.SNMPChannel{
				Community: "public",
				Version:   "v1",
			},
		},
	}
}

// NewDefaultProtosGatewayConfig generates a basic devmand gateway config
func NewDefaultProtosGatewayConfig() *protos.DevmandGatewayConfig {
	return &protos.DevmandGatewayConfig{
		ManagedDevices: []string{"test_device_1", "test_device_2"},
	}
}

func NewDefaultManagedDevice() *models.ManagedDevice {
	deviceConfigJSONString := `{
		"interfaces": {
			"interface": [
				{
					"config": {
						"description": "ACCESS PORT",
						"name": "g0/0/0",
						"type": "ethernetCsmacd"
					},
					"name": "g0/0/0"
				}
			]
		}
	}`
	return &models.ManagedDevice{
		DeviceConfig: deviceConfigJSONString,
		Host:         "1.2.3.4",
		DeviceType:   []string{"wifi", "access_point"},
		Platform:     "Ubiquiti",
		Channels: &models.ManagedDeviceChannels{
			SnmpChannel: &models.SnmpChannel{
				Community: "public",
				Version:   "v1",
			},
		},
	}
}

func NewDefaultGatewayConfig() *models.GatewayDevmandConfigs {
	return &models.GatewayDevmandConfigs{
		ManagedDevices: []string{"test_device_1", "test_device_2"},
	}
}
