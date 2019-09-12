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
	"magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/storage"

	"github.com/stretchr/testify/assert"
)

const (
	typeVal   = "type"
	networkID = "network1"
)

type idAndInfo struct {
	deviceKey  string
	deviceType string
	info       interface{}
}

func TestDeviceService(t *testing.T) {
	testSerde := &Serde{}
	err := serde.RegisterSerdes(testSerde)
	assert.NoError(t, err)
	test_init.StartTestService(t)

	bundle1 := idAndInfo{deviceKey: "device1", deviceType: typeVal, info: 1}
	bundle2 := idAndInfo{deviceKey: "device2", deviceType: typeVal, info: 2}
	// Entities that should fail to register
	unregisteredSerdeBundle := idAndInfo{deviceKey: "device2", deviceType: "unregistered", info: 2}
	unserializableBundle := idAndInfo{deviceKey: "device2", deviceType: typeVal, info: "(*.*)"}

	// Check contract for empty network
	assertDevicesNotRegistered(t, bundle1, bundle2)

	// Check contract for empty requests
	registerDevicesAssertError(t, "", bundle1)

	// Registering ill formatted device values should fail
	registerDevicesAssertError(t, networkID, unregisteredSerdeBundle)
	registerDevicesAssertError(t, networkID, unserializableBundle)

	// Register and retrieve devices
	registerDevicesAssertNoError(t, networkID, bundle1)
	registerDevicesAssertNoError(t, networkID, bundle2)
	assertDevicesAreRegistered(t, bundle1, bundle2)

	// Registering a key already registered should fail
	registerDevicesAssertError(t, "network2", bundle1)

	// Update Devices
	bundle1.info = 5
	updateDevicesAssertNoError(t, networkID, bundle1)

	// Test deletion
	err = device.DeleteDevices(networkID, []storage.TypeAndKey{{Type: bundle1.deviceType, Key: bundle1.deviceKey}})
	assert.NoError(t, err)
	assertDevicesNotRegistered(t, bundle1)
	assertDevicesAreRegistered(t, bundle2)
}

func assertDevicesAreRegistered(t *testing.T, bundles ...idAndInfo) {
	for _, bundle := range bundles {
		actualInfo, err := device.GetDevice(networkID, bundle.deviceType, bundle.deviceKey)
		assert.NoError(t, err)
		assert.Equal(t, bundle.info, actualInfo)
	}
}

func assertDevicesNotRegistered(t *testing.T, bundles ...idAndInfo) {
	for _, bundle := range bundles {
		_, err := device.GetDevice(networkID, bundle.deviceType, bundle.deviceKey)
		assert.Error(t, err)
	}
}

func registerDevicesAssertNoError(t *testing.T, networkID string, bundle idAndInfo) {
	err := device.RegisterDevice(networkID, bundle.deviceType, bundle.deviceKey, bundle.info)
	assert.NoError(t, err)
}

func registerDevicesAssertError(t *testing.T, networkID string, bundle idAndInfo) {
	err := device.RegisterDevice(networkID, bundle.deviceType, bundle.deviceKey, bundle.info)
	assert.Error(t, err)
}

func updateDevicesAssertNoError(t *testing.T, networkID string, bundle idAndInfo) {
	err := device.UpdateDevice(networkID, bundle.deviceType, bundle.deviceKey, bundle.info)
	assert.NoError(t, err)
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
