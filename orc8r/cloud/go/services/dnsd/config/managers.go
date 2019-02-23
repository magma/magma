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
	dns_protos "magma/orc8r/cloud/go/services/dnsd/protos"
	"magma/orc8r/cloud/go/services/magmad"
)

const (
	DnsdNetworkType = "dnsd_network"
)

type DnsNetworkConfigManager struct{}

func (*DnsNetworkConfigManager) GetConfigType() string {
	return DnsdNetworkType
}

func (*DnsNetworkConfigManager) GetGatewayIdsForConfig(networkId string, configKey string) ([]string, error) {
	return magmad.ListGateways(networkId)
}

func (*DnsNetworkConfigManager) MarshalConfig(config interface{}) ([]byte, error) {
	castedConfig, ok := config.(*dns_protos.NetworkDNSConfig)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid config type. Expected *NetworkDNSConfig, received %s",
			reflect.TypeOf(config),
		)
	}
	if err := dns_protos.ValidateNetworkConfig(castedConfig); err != nil {
		return nil, fmt.Errorf("Invalid network dns config: %s", err)
	}
	return protos.MarshalIntern(castedConfig)
}

func (*DnsNetworkConfigManager) UnmarshalConfig(message []byte) (interface{}, error) {
	cfg := &dns_protos.NetworkDNSConfig{}
	err := protos.Unmarshal(message, cfg)
	return cfg, err
}
