/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package models includes definitions of swagger generated REST API model Go structures
package models

import (
	"magma/orc8r/cloud/go/protos"

	"github.com/go-openapi/strfmt"
)

// ConvertibleConfig for user-facing swagger models
var formatsRegistry = strfmt.NewFormats()

func (m *NetworkCarrierWifiConfigs) ValidateModel() error {
	return m.Validate(formatsRegistry)
}

func (m *NetworkCarrierWifiConfigs) ToServiceModel() (interface{}, error) {
	return m, nil
}

func (m *NetworkCarrierWifiConfigs) FromServiceModel(magmadModel interface{}) error {
	protos.FillIn(magmadModel, m)
	return nil
}
