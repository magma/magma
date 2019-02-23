/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos

import "errors"

func ValidateGatewayConfig(config *MagmadGatewayConfig) error {
	if config == nil {
		return errors.New("Gateway config is nil")
	}
	if config.GetTier() == "" {
		return errors.New("Tier ID must be specified")
	}
	return nil
}

func ValidateNetworkConfig(config *MagmadNetworkRecord) error {
	if config == nil {
		return errors.New("Network config is nil")
	}
	if config.GetName() == "" {
		return errors.New("Network name must be specified")
	}
	return nil
}
