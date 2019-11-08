/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package models

import (
	"github.com/go-openapi/strfmt"
)

var formatsRegistry = strfmt.NewFormats()

// Config conversion for magmad

func (m *MagmadGatewayConfig) ValidateModel() error {
	if err := m.ValidateGatewayConfig(); err != nil {
		return err
	}
	return m.Validate(formatsRegistry)
}

func (m *NetworkRecord) ValidateModel() error {
	if err := m.ValidateNetworkRecord(); err != nil {
		return err
	}
	return m.Validate(formatsRegistry)
}
