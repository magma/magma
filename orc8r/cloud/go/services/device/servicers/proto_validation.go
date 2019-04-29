/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"fmt"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/device/protos"
)

func ValidateRegisterDevicesRequest(req *protos.RegisterDevicesRequest) error {
	if err := nonEmptyNetworkID(req.GetNetworkID()); err != nil {
		return err
	}
	entities := req.GetEntities()
	if err := nonEmptyEntities(entities); err != nil {
		return err
	}
	return deserializableWithSerde(entities)
}

func ValidateGetDeviceInfoRequest(req *protos.GetDeviceInfoRequest) error {
	return nonEmptyNetworkIDAndDeviceIDs(req.GetNetworkID(), req.GetDeviceIDs())
}

func ValidateDeleteDevicesRequest(req *protos.DeleteDevicesRequest) error {
	return nonEmptyNetworkIDAndDeviceIDs(req.GetNetworkID(), req.GetDeviceIDs())
}

func deserializableWithSerde(entities []*protos.PhysicalEntity) error {
	for _, entity := range entities {
		_, err := serde.Deserialize(device.SerdeDomain, entity.GetType(), entity.GetInfo())
		if err != nil {
			return err
		}
	}
	return nil
}

func nonEmptyNetworkID(networkID string) error {
	if len(networkID) == 0 {
		return fmt.Errorf("NetworkID must be non-empty")
	}
	return nil
}

func nonEmptyDeviceIDs(deviceIDs []*protos.DeviceID) error {
	if deviceIDs == nil || len(deviceIDs) == 0 {
		return fmt.Errorf("DeviceIDs field must be non-empty")
	}
	return nil
}

func nonEmptyEntities(entities []*protos.PhysicalEntity) error {
	if entities == nil || len(entities) == 0 {
		return fmt.Errorf("Entities field must be non-empty")
	}
	return nil
}

func nonEmptyNetworkIDAndDeviceIDs(networkID string, deviceIDs []*protos.DeviceID) error {
	if err := nonEmptyNetworkID(networkID); err != nil {
		return err
	}
	if err := nonEmptyDeviceIDs(deviceIDs); err != nil {
		return err
	}
	return nil
}
