/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package device_test

import (
	"fmt"
	"strconv"
	"testing"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/device/protos"
	"magma/orc8r/cloud/go/services/device/test_init"

	"github.com/stretchr/testify/assert"
)

const (
	typeVal   = "type"
	networkID = "network1"
)

type idAndEntity struct {
	id     *protos.DeviceID
	entity *protos.PhysicalEntity
}

func makeIDAndEntity(deviceID string, typeVal string, value []byte) idAndEntity {
	id := protos.DeviceID{DeviceID: deviceID, Type: typeVal}
	entity := protos.PhysicalEntity{
		DeviceID: deviceID,
		Type:     typeVal,
		Info:     value,
	}
	return idAndEntity{id: &id, entity: &entity}
}

func TestDeviceService(t *testing.T) {
	testSerde := &Serde{}
	err := serde.RegisterSerdes(testSerde)
	assert.NoError(t, err)
	test_init.StartTestService(t)

	serialized1, err := testSerde.Serialize(1)
	serialized2, err := testSerde.Serialize(2)
	bundle1 := makeIDAndEntity("device1", typeVal, serialized1)
	bundle2 := makeIDAndEntity("device2", typeVal, serialized2)
	// Entities that should fail to register
	unregisteredSerdeBundle := makeIDAndEntity("device2", "unregistered", serialized2)
	unserializableBundle := makeIDAndEntity("device3", typeVal, []byte("(*_*)"))

	// Check contract for empty network
	assertDevicesNotRegistered(t, bundle1.id, bundle2.id)

	// Check contract for empty requests
	registerDevicesAssertError(t, networkID)
	registerDevicesAssertError(t, "", bundle1.entity)

	// Registering ill formatted device values should fail
	registerDevicesAssertError(t, networkID, unregisteredSerdeBundle.entity)
	registerDevicesAssertError(t, networkID, unserializableBundle.entity)

	// Register and retrieve devices
	registerDevicesAssertNoError(t, networkID, bundle1.entity, bundle2.entity)
	assertDevicesAreRegistered(t, bundle1, bundle2)

	// Test deletion
	err = device.DeleteDevices(networkID, []*protos.DeviceID{bundle1.id})
	assert.NoError(t, err)
	assertDevicesNotRegistered(t, bundle1.id)
	assertDevicesAreRegistered(t, bundle2)
}

func assertDevicesAreRegistered(t *testing.T, bundles ...idAndEntity) {
	deviceIDs := []*protos.DeviceID{}
	for _, bundle := range bundles {
		deviceIDs = append(deviceIDs, bundle.id)
	}

	deviceMap, err := device.GetDeviceInfo(networkID, deviceIDs)
	assert.NoError(t, err)
	assertDevicesInEntityMap(t, deviceMap, bundles)
}

func assertDevicesNotRegistered(t *testing.T, deviceIDs ...*protos.DeviceID) {
	deviceMap, err := device.GetDeviceInfo(networkID, deviceIDs)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(deviceMap))
}

func registerDevicesAssertNoError(t *testing.T, networkID string, entities ...*protos.PhysicalEntity) {
	err := device.RegisterDevices(networkID, entities)
	assert.NoError(t, err)
}

func registerDevicesAssertError(t *testing.T, networkID string, entities ...*protos.PhysicalEntity) {
	err := device.RegisterDevices(networkID, entities)
	assert.Error(t, err)
}

func assertDevicesInEntityMap(t *testing.T, deviceMap map[string]*protos.PhysicalEntity, bundles []idAndEntity) {
	for _, bundle := range bundles {
		entity := deviceMap[bundle.id.DeviceID]
		assert.NotNil(t, entity)
		assert.Equal(t, bundle.entity.Info, entity.Info)
	}
	assert.Equal(t, len(bundles), len(deviceMap))
}

type Serde struct {
}

func (*Serde) GetDomain() string {
	return device.SerdeDomain
}

func (*Serde) GetType() string {
	return typeVal
}

func (*Serde) Serialize(in interface{}) ([]byte, error) {
	return []byte(fmt.Sprintf("%d", in)), nil

}

func (*Serde) Deserialize(message []byte) (interface{}, error) {
	return strconv.Atoi(string(message))
}
