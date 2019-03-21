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

	config_protos "magma/feg/cloud/go/services/controller/protos"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config"
)

const (
	FegNetworkType = "federation_network"
	FegGatewayType = "federation_gateway"
)

type FegNetworkConfigManager struct{}

func (*FegNetworkConfigManager) GetDomain() string {
	return config.SerdeDomain
}

func (*FegNetworkConfigManager) GetType() string {
	return FegNetworkType
}

func (*FegNetworkConfigManager) Serialize(config interface{}) ([]byte, error) {
	castedConfig, ok := config.(*config_protos.Config)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid FeG network config type. Expected *Config, received %s",
			reflect.TypeOf(config),
		)
	}
	if err := config_protos.ValidateNetworkConfig(castedConfig); err != nil {
		return nil, fmt.Errorf("Invalid FeG network config: %s", err)
	}
	return protos.MarshalIntern(castedConfig)
}

func (*FegNetworkConfigManager) Deserialize(message []byte) (interface{}, error) {
	cfg := &config_protos.Config{}
	err := protos.Unmarshal(message, cfg)
	return cfg, err
}

type FegGatewayConfigManager struct{}

func (*FegGatewayConfigManager) GetDomain() string {
	return config.SerdeDomain
}

func (*FegGatewayConfigManager) GetType() string {
	return FegGatewayType
}

func (*FegGatewayConfigManager) Serialize(config interface{}) ([]byte, error) {
	castedConfig, ok := config.(*config_protos.Config)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid FeG gateway config type. Expected *Config, received %s",
			reflect.TypeOf(config),
		)
	}
	if err := config_protos.ValidateGatewayConfig(castedConfig); err != nil {
		return nil, fmt.Errorf("Invalid FeG gateway config: %s", err)
	}
	return protos.MarshalIntern(castedConfig)
}

func (*FegGatewayConfigManager) Deserialize(message []byte) (interface{}, error) {
	cfg := &config_protos.Config{}
	err := protos.Unmarshal(message, cfg)
	return cfg, err
}
