/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package test_utils

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/magmad/obsidian/models"

	"github.com/stretchr/testify/assert"
)

func RegisterNetwork(t *testing.T, networkID string, networkName string) {
	err := configurator.CreateNetwork(
		configurator.Network{
			ID:   networkID,
			Name: networkName,
		})
	assert.NoError(t, err)
}

func RegisterGateway(t *testing.T, networkID string, gatewayID string, record *models.AccessGatewayRecord) {
	var gwEntity configurator.NetworkEntity
	if record != nil {
		if device.DoesDeviceExist(networkID, orc8r.AccessGatewayRecordType, record.HwID.ID) {
			t.Fatalf("Hwid is already registered %s", record.HwID.ID)
		}
		// write into device
		err := device.CreateOrUpdate(networkID, orc8r.AccessGatewayRecordType, record.HwID.ID, record)
		assert.NoError(t, err)

		gwEntity = configurator.NetworkEntity{
			Name:       record.Name,
			Type:       orc8r.MagmadGatewayType,
			Key:        gatewayID,
			PhysicalID: record.HwID.ID,
		}
	} else {
		gwEntity = configurator.NetworkEntity{
			Type: orc8r.MagmadGatewayType,
			Key:  gatewayID,
		}
	}
	_, err := configurator.CreateEntity(networkID, gwEntity)
	assert.NoError(t, err)
}

// RemoveGateway assumes there is a device entity corresponding to the
// configurator entity
func RemoveGateway(t *testing.T, networkID, gatewayID string) {
	physicalID, err := configurator.GetPhysicalIDOfEntity(networkID, orc8r.MagmadGatewayType, gatewayID)
	assert.NoError(t, err)
	assert.NoError(t, device.DeleteDevice(networkID, orc8r.AccessGatewayRecordType, physicalID))
	assert.NoError(t, configurator.DeleteEntity(networkID, orc8r.MagmadGatewayType, gatewayID))
}
