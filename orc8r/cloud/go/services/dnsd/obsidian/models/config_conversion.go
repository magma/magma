/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Implementation of ConvertibleConfig for user-facing swagger models
package models

import (
	"fmt"
	"reflect"

	"magma/orc8r/cloud/go/protos"
	dnsdprotos "magma/orc8r/cloud/go/services/dnsd/protos"

	"github.com/go-openapi/strfmt"
)

var formatsRegistry strfmt.Registry = strfmt.NewFormats()

func (m *NetworkDNSConfig) ValidateModel() error {
	if err := m.ValidateNetworkConfig(); err != nil {
		return err
	}
	return m.Validate(formatsRegistry)
}

func (m *NetworkDNSConfig) ToServiceModel() (interface{}, error) {
	magmadConfig := &dnsdprotos.NetworkDNSConfig{}

	protos.FillIn(m, magmadConfig)
	if err := dnsdprotos.ValidateNetworkConfig(magmadConfig); err != nil {
		return nil, err
	}
	return magmadConfig, nil
}

func (m *NetworkDNSConfig) FromServiceModel(magmadModel interface{}) error {
	_, ok := magmadModel.(*dnsdprotos.NetworkDNSConfig)
	if !ok {
		return fmt.Errorf(
			"Invalid magmad config type to convert to. Expected *NetworkDNSConfig but got %s",
			reflect.TypeOf(magmadModel),
		)
	}
	protos.FillIn(magmadModel, m)
	return nil
}
