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

	cellularprotos "magma/lte/cloud/go/services/cellular/protos"
	"magma/orc8r/cloud/go/protos"

	"github.com/go-openapi/strfmt"
)

var formatsRegistry = strfmt.NewFormats()

func init() {
	// Echo encodes/decodes base64 encoded byte arrays, no verification needed
	b64 := strfmt.Base64([]byte(nil))
	formatsRegistry.Add("byte", &b64, func(_ string) bool { return true })
}

func (m *NetworkCellularConfigs) ValidateModel() error {
	return m.Validate(formatsRegistry)
}

func (m *NetworkCellularConfigs) ToServiceModel() (interface{}, error) {
	magmadConfig := &cellularprotos.CellularNetworkConfig{}
	protos.FillIn(m, magmadConfig)
	magmadConfig.FegNetworkId = m.FegNetworkID
	if err := m.networkServicesToServiceModel(magmadConfig); err != nil {
		return nil, err
	}
	if err := cellularprotos.ValidateNetworkConfig(magmadConfig); err != nil {
		return nil, err
	}
	magmadConfig.Epc.RelayEnabled = m.Epc.RelayEnabled
	return magmadConfig, nil
}

func (m *NetworkCellularConfigs) networkServicesToServiceModel(magmadConfig *cellularprotos.CellularNetworkConfig) error {
	for _, serviceName := range m.Epc.NetworkServices {
		serviceEnum, err := cellularprotos.GetNetworkServiceEnum(serviceName)
		if err != nil {
			return err
		}
		magmadConfig.Epc.NetworkServices = append(magmadConfig.Epc.NetworkServices, serviceEnum)
	}
	return nil
}

func (m *NetworkCellularConfigs) FromServiceModel(magmadModel interface{}) error {
	magmadConfig, ok := magmadModel.(*cellularprotos.CellularNetworkConfig)
	if !ok {
		return fmt.Errorf(
			"Invalid magmad config type to convert to. Expected *CellularNetworkConfig but got %s",
			reflect.TypeOf(magmadModel),
		)
	}
	protos.FillIn(magmadModel, m)
	m.FegNetworkID = magmadConfig.FegNetworkId
	m.Epc.RelayEnabled = magmadConfig.Epc.RelayEnabled
	if err := m.networkServicesFromServiceModel(magmadConfig); err != nil {
		return err
	}
	return nil
}

func (m *NetworkCellularConfigs) networkServicesFromServiceModel(magmadConfig *cellularprotos.CellularNetworkConfig) error {
	if magmadConfig.Epc.NetworkServices != nil {
		for _, serviceEnum := range magmadConfig.Epc.NetworkServices {
			serviceName, err := cellularprotos.GetNetworkServiceName(serviceEnum)
			if err != nil {
				return err
			}
			m.Epc.NetworkServices = append(m.Epc.NetworkServices, serviceName)
		}
	}
	return nil
}

func (m *GatewayCellularConfigs) ValidateModel() error {
	return m.Validate(formatsRegistry)
}

func (m *GatewayCellularConfigs) ToServiceModel() (interface{}, error) {
	magmadConfig := &cellularprotos.CellularGatewayConfig{}
	protos.FillIn(m, magmadConfig)
	if err := cellularprotos.ValidateGatewayConfig(magmadConfig); err != nil {
		return nil, err
	}
	return magmadConfig, nil
}

func (m *GatewayCellularConfigs) FromServiceModel(magmadModel interface{}) error {
	protos.FillIn(magmadModel, m)
	return nil
}

func (m *NetworkEnodebConfigs) ValidateModel() error {
	return nil
}

func (m *NetworkEnodebConfigs) ToServiceModel() (interface{}, error) {
	magmadConfig := &cellularprotos.CellularEnodebConfig{}
	protos.FillIn(m, magmadConfig)
	if err := cellularprotos.ValidateEnodebConfig(magmadConfig); err != nil {
		return nil, err
	}
	return magmadConfig, nil
}

func (m *NetworkEnodebConfigs) FromServiceModel(magmadModel interface{}) error {
	protos.FillIn(magmadModel, m)
	return nil
}
