/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"fmt"
	"reflect"

	cellular_protos "magma/lte/cloud/go/services/cellular/protos"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config"
)

const (
	CellularNetworkType = "cellular_network"
	CellularGatewayType = "cellular_gateway"
)

type CellularNetworkConfigManager struct{}

func (*CellularNetworkConfigManager) GetDomain() string {
	return config.SerdeDomain
}

func (*CellularNetworkConfigManager) GetType() string {
	return CellularNetworkType
}

func (*CellularNetworkConfigManager) Serialize(config interface{}) ([]byte, error) {
	castedConfig, ok := config.(*cellular_protos.CellularNetworkConfig)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid cellular network config type. Expected *CellularNetworkConfig, received %s",
			reflect.TypeOf(config),
		)
	}
	if err := cellular_protos.ValidateNetworkConfig(castedConfig); err != nil {
		return nil, fmt.Errorf("Invalid cellular network config: %s", err)
	}
	return protos.MarshalIntern(castedConfig)
}

func (*CellularNetworkConfigManager) Deserialize(message []byte) (interface{}, error) {
	cfg := &cellular_protos.CellularNetworkConfig{}
	err := protos.Unmarshal(message, cfg)
	return cfg, err
}

type CellularGatewayConfigManager struct{}

func (*CellularGatewayConfigManager) GetDomain() string {
	return config.SerdeDomain
}

func (*CellularGatewayConfigManager) GetType() string {
	return CellularGatewayType
}

func (*CellularGatewayConfigManager) Serialize(config interface{}) ([]byte, error) {
	castedConfig, ok := config.(*cellular_protos.CellularGatewayConfig)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid cellular gateway config type. Expected *CellularGatewayConfig, received %s",
			reflect.TypeOf(config),
		)
	}
	if err := cellular_protos.ValidateGatewayConfig(castedConfig); err != nil {
		return nil, fmt.Errorf("Invalid cellular gateway config: %s", err)
	}
	return protos.MarshalIntern(castedConfig)
}

func (*CellularGatewayConfigManager) Deserialize(message []byte) (interface{}, error) {
	cfg := &cellular_protos.CellularGatewayConfig{}
	err := protos.Unmarshal(message, cfg)
	return cfg, err
}
