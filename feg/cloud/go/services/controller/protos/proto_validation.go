/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos

import (
	"errors"
)

func ValidateGatewayConfig(config *Config) error {
	if config == nil {
		return errors.New("Gateway config is nil")
	}
	return nil
}

func ValidateNetworkConfig(config *Config) error {
	if config == nil {
		return errors.New("Network config is nil")
	}
	return nil
}
