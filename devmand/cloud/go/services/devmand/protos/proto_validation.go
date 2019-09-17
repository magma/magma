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
