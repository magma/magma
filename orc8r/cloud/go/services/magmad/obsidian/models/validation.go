/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import "errors"

func (m *MagmadGatewayConfig) ValidateGatewayConfig() error {
	if m == nil {
		return errors.New("Gateway m is nil")
	}
	if m.Tier == "" {
		return errors.New("Tier ID must be specified")
	}
	return nil
}

func (m *NetworkRecord) ValidateNetworkRecord() error {
	if m == nil {
		return errors.New("Network m is nil")
	}
	if m.Name == "" {
		return errors.New("Network name must be specified")
	}
	return nil
}
