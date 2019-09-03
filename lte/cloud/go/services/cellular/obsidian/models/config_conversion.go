/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Implementation of ConvertibleConfig for user-facing swagger models
package models

import (
	"github.com/go-openapi/strfmt"
)

var formatsRegistry = strfmt.NewFormats()

func init() {
	// Echo encodes/decodes base64 encoded byte arrays, no verification needed
	b64 := strfmt.Base64([]byte(nil))
	formatsRegistry.Add("byte", &b64, func(_ string) bool { return true })
}

func (m *NetworkCellularConfigs) ValidateModel() error {
	if err := m.Validate(formatsRegistry); err != nil {
		return err
	}
	return m.ValidateNetworkConfig()
}

func (m *GatewayCellularConfigs) ValidateModel() error {
	if err := m.Validate(formatsRegistry); err != nil {
		return err
	}
	return m.ValidateGatewayConfig()
}

func (m *NetworkEnodebConfigs) ValidateModel() error {
	return m.ValidateEnodebConfig()
}
