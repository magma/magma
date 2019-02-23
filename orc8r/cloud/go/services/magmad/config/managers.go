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

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
)

const (
	MagmadGatewayType = "magmad_gateway"
	MagmadNetworkType = "magmad_network"
)

// To be deprecated! Config service DB tables have been seeded with magmad
// network configs that were migrated from legacy magmad network configs,
// so this will stick around for a bit. We can delete this after deleting
// all magmad network config types from the config service (to come in a
// future migration)
type MagmadNetworkConfigManager struct{}

func (*MagmadNetworkConfigManager) GetConfigType() string {
	return MagmadNetworkType
}

func (*MagmadNetworkConfigManager) GetGatewayIdsForConfig(networkId string, configKey string) ([]string, error) {
	return magmad.ListGateways(networkId)
}

func (*MagmadNetworkConfigManager) MarshalConfig(config interface{}) ([]byte, error) {
	castedConfig, ok := config.(*magmad_protos.MagmadNetworkRecord)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid magmad network config type. Expected *MagmadNetworkRecord, received %s",
			reflect.TypeOf(config),
		)
	}
	if err := magmad_protos.ValidateNetworkConfig(castedConfig); err != nil {
		return nil, fmt.Errorf("Invalid network config: %s", err)
	}
	return protos.MarshalIntern(castedConfig)
}

func (*MagmadNetworkConfigManager) UnmarshalConfig(message []byte) (interface{}, error) {
	cfg := &magmad_protos.MagmadNetworkRecord{}
	err := protos.Unmarshal(message, cfg)
	return cfg, err
}

type MagmadGatewayConfigManager struct{}

func (*MagmadGatewayConfigManager) GetConfigType() string {
	return MagmadGatewayType
}

func (*MagmadGatewayConfigManager) GetGatewayIdsForConfig(networkId string, configKey string) ([]string, error) {
	return []string{configKey}, nil
}

func (*MagmadGatewayConfigManager) MarshalConfig(config interface{}) ([]byte, error) {
	castedConfig, ok := config.(*magmad_protos.MagmadGatewayConfig)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid magmad gateway config type. Expected *MagmadGatewayConfig, received %s",
			reflect.TypeOf(config),
		)
	}
	if err := magmad_protos.ValidateGatewayConfig(castedConfig); err != nil {
		return nil, fmt.Errorf("Invalid gateway config: %s", err)
	}
	return protos.MarshalIntern(castedConfig)
}

func (*MagmadGatewayConfigManager) UnmarshalConfig(message []byte) (interface{}, error) {
	cfg := &magmad_protos.MagmadGatewayConfig{}
	err := protos.Unmarshal(message, cfg)
	return cfg, err
}
