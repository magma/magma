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
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/config"
	dns_protos "magma/orc8r/cloud/go/services/dnsd/protos"

	"github.com/golang/protobuf/proto"
)

type DnsdMconfigBuilder struct{}

func (builder *DnsdMconfigBuilder) Build(networkId string, gatewayId string) (map[string]proto.Message, error) {
	networkDNSconfig, err := GetNetworkDNSConfig(networkId)
	if err != nil {
		return nil, err
	}
	if networkDNSconfig == nil {
		return map[string]proto.Message{}, nil
	}

	mconfigDnsD := &mconfig.DnsD{}
	protos.FillIn(networkDNSconfig, mconfigDnsD)
	mconfigDnsD.LogLevel = protos.LogLevel_INFO

	return map[string]proto.Message{
		"dnsd": mconfigDnsD,
	}, nil
}

func GetNetworkDNSConfig(networkId string) (*dns_protos.NetworkDNSConfig, error) {
	iNetworkDNSconfigs, err := config.GetConfig(networkId, DnsdNetworkType, networkId)
	if err != nil || iNetworkDNSconfigs == nil {
		return nil, err
	}
	networkDNSconfig, ok := iNetworkDNSconfigs.(*dns_protos.NetworkDNSConfig)
	if !ok {
		return nil, fmt.Errorf(
			"Received unexpected type for network record. "+
				"Expected *NetworkDNSconfig but got %s",
			reflect.TypeOf(iNetworkDNSconfigs),
		)
	}
	return networkDNSconfig, nil
}
