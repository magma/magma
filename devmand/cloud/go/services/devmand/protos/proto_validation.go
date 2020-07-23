/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package protos

import (
	"errors"
)

// ValidateGatewayConfig - Validate a DevmandGatewayConfig
func ValidateGatewayConfig(config *DevmandGatewayConfig) error {
	if config == nil {
		return errors.New("Gateway config is nil")
	}
	return nil
}

// ValidateManagedDevice - validate a ManagedDevice
func ValidateManagedDevice(config *ManagedDevice) error {
	if config == nil {
		return errors.New("Device config is nil")
	}
	return nil
}
