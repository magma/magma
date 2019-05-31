/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package device

import (
	"context"

	"magma/orc8r/cloud/go/errors"
	magma_errors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/device/protos"

	"github.com/golang/glog"
)

func getDeviceClient() (protos.DeviceClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewDeviceClient(conn), err
}

func CreateOrUpdate(networkID, deviceType, deviceKey string, info interface{}) error {
	client, err := getDeviceClient()
	if err != nil {
		return err
	}

	serializedInfo, err := serde.Serialize(SerdeDomain, deviceType, info)
	if err != nil {
		return err
	}
	entity := &protos.PhysicalEntity{
		DeviceID: deviceKey,
		Type:     deviceType,
		Info:     serializedInfo,
	}
	req := &protos.RegisterDevicesRequest{
		NetworkID: networkID,
		Entities:  []*protos.PhysicalEntity{entity},
	}
	_, err = client.RegisterDevices(context.Background(), req)
	return err
}

func DeleteDevices(networkID string, deviceIDs []*protos.DeviceID) error {
	client, err := getDeviceClient()
	if err != nil {
		return err
	}

	req := &protos.DeleteDevicesRequest{NetworkID: networkID, DeviceIDs: deviceIDs}
	_, err = client.DeleteDevices(context.Background(), req)
	return err
}

func GetDevice(networkID, deviceType, deviceKey string) (interface{}, error) {
	client, err := getDeviceClient()
	if err != nil {
		return nil, err
	}
	deviceID := &protos.DeviceID{Type: deviceType, DeviceID: deviceKey}
	req := &protos.GetDeviceInfoRequest{NetworkID: networkID, DeviceIDs: []*protos.DeviceID{deviceID}}
	res, err := client.GetDeviceInfo(context.Background(), req)
	if err != nil {
		return nil, err
	}
	device, ok := res.DeviceMap[deviceKey]
	if !ok {
		return nil, magma_errors.ErrNotFound
	}
	return serde.Deserialize(SerdeDomain, deviceType, device.Info)
}

func DoesDeviceExist(networkID, deviceType, deviceID string) bool {
	_, err := GetDevice(networkID, deviceType, deviceID)
	if err != nil {
		return false
	}
	return true
}
