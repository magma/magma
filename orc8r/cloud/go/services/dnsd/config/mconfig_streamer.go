/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"fmt"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/config/streaming"
	dnsprotos "magma/orc8r/cloud/go/services/dnsd/protos"

	"github.com/golang/protobuf/ptypes"
)

type DnsdStreamer struct{}

func (*DnsdStreamer) GetSubscribedConfigTypes() []string {
	return []string{DnsdNetworkType}
}

func (*DnsdStreamer) SeedNewGatewayMconfig(
	networkId string,
	gatewayId string,
	mconfigOut *protos.GatewayConfigs, // output parameter
) error {
	dnsdNetworkConfig, err := config.GetConfig(networkId, DnsdNetworkType, networkId)
	if err != nil {
		return err
	}
	if dnsdNetworkConfig == nil {
		return nil
	}

	return applyDnsdNetworkConfig(dnsdNetworkConfig.(*dnsprotos.NetworkDNSConfig), mconfigOut)
}

func (*DnsdStreamer) ApplyMconfigUpdate(
	update *streaming.ConfigUpdate,
	oldMconfigsByGatewayId map[string]*protos.GatewayConfigs,
) (map[string]*protos.GatewayConfigs, error) {
	if update.ConfigType != DnsdNetworkType {
		return oldMconfigsByGatewayId, fmt.Errorf("Dnsd mconfig streamer received update to unsubscribed type %s", update.ConfigType)
	}
	switch update.Operation {
	case streaming.DeleteOperation:
		for _, mconfigValue := range oldMconfigsByGatewayId {
			delete(mconfigValue.ConfigsByKey, "dnsd")
		}
		return oldMconfigsByGatewayId, nil
	case streaming.CreateOperation, streaming.ReadOperation, streaming.UpdateOperation:
		if update.NewValue == nil {
			return oldMconfigsByGatewayId, nil
		}

		newValueCasted := update.NewValue.(*dnsprotos.NetworkDNSConfig)
		for _, mconfigValue := range oldMconfigsByGatewayId {
			if err := applyDnsdNetworkConfig(newValueCasted, mconfigValue); err != nil {
				return oldMconfigsByGatewayId, err
			}
		}
		return oldMconfigsByGatewayId, nil
	default:
		return oldMconfigsByGatewayId, fmt.Errorf("Unrecognized streaming operation: %s", update.Operation)
	}
}

func applyDnsdNetworkConfig(newConfig *dnsprotos.NetworkDNSConfig, mconfigOut *protos.GatewayConfigs) error {
	mconfigDnsd := &mconfig.DnsD{}
	protos.FillIn(newConfig, mconfigDnsd)
	mconfigDnsd.LogLevel = protos.LogLevel_INFO

	mconfigDnsdAny, err := ptypes.MarshalAny(mconfigDnsd)
	if err != nil {
		return err
	}
	mconfigOut.ConfigsByKey["dnsd"] = mconfigDnsdAny
	return nil
}
