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
	"magma/orc8r/cloud/go/registry"
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

func RegisterDevices(networkID string, entities []*protos.PhysicalEntity) error {
	client, err := getDeviceClient()
	if err != nil {
		return err
	}
	req := &protos.RegisterDevicesRequest{NetworkID: networkID, Entities: entities}
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

func GetDeviceInfo(networkID string, deviceIDs []*protos.DeviceID) (map[string]*protos.PhysicalEntity, error) {
	client, err := getDeviceClient()
	if err != nil {
		return nil, err
	}
	req := &protos.GetDeviceInfoRequest{NetworkID: networkID, DeviceIDs: deviceIDs}
	res, err := client.GetDeviceInfo(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res.DeviceMap, nil
}
