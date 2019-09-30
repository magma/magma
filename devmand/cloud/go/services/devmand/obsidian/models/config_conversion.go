/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package models

import (
	"fmt"
	"reflect"

	orcprotos "magma/orc8r/cloud/go/protos"
	"orc8r/devmand/cloud/go/services/devmand/protos"

	"github.com/go-openapi/strfmt"
)

var formatsRegistry = strfmt.NewFormats()

//ValidateModel - Validate a GatewayDevmandConfigs model
func (model *GatewayDevmandConfigs) ValidateModel() error {
	return model.Validate(formatsRegistry)
}

//ToServiceModel - Convert from GatewayDevmandConfigs to a DevmandModel
func (model *GatewayDevmandConfigs) ToServiceModel() (interface{}, error) {
	devmandConfig := &protos.DevmandGatewayConfig{}
	orcprotos.FillIn(model, devmandConfig)

	err := protos.ValidateGatewayConfig(devmandConfig)
	if err != nil {
		return nil, err
	}

	return devmandConfig, nil
}

//FromServiceModel - Convert from a GatewayDevmandConfigs to a DevmandModel
func (model *GatewayDevmandConfigs) FromServiceModel(devmandModel interface{}) error {
	_, ok := devmandModel.(*protos.DevmandGatewayConfig)
	if !ok {
		return fmt.Errorf(
			"Invalid devmandd config type to convert to. Expected *DevmandGatewayConfig but got %s",
			reflect.TypeOf(devmandModel),
		)
	}
	orcprotos.FillIn(devmandModel, model)
	return nil
}

//ValidateModel - Validate a ManagedDevice model
func (model *ManagedDevice) ValidateModel() error {
	return model.Validate(formatsRegistry)
}

//ToServiceModel - Convert from ManagedDevice to a DevmandModel
func (model *ManagedDevice) ToServiceModel() (interface{}, error) {
	managedDevice := &protos.ManagedDevice{}
	orcprotos.FillIn(model, managedDevice)

	err := protos.ValidateManagedDevice(managedDevice)
	if err != nil {
		return nil, err
	}

	return managedDevice, nil
}

//FromServiceModel - Convert from a ManagedDevice to a DevmandModel
func (model *ManagedDevice) FromServiceModel(managedDeviceModel interface{}) error {
	_, ok := managedDeviceModel.(*protos.ManagedDevice)
	if !ok {
		return fmt.Errorf(
			"Invalid device config type to convert to. Expected *ManagedDevice but got %s",
			reflect.TypeOf(managedDeviceModel),
		)
	}
	orcprotos.FillIn(managedDeviceModel, model)
	return nil
}
