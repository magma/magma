/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/streaming/storage"
	"magma/orc8r/cloud/go/services/streamer/mconfig/factory"
)

type NetworkStoredMconfigs = map[string]*storage.StoredMconfig

func loadStoredMconfigsForNetworks(gatewayIdsByNetworkId map[string][]string) (map[string]NetworkStoredMconfigs, error) {
	ret := make(map[string]NetworkStoredMconfigs, len(gatewayIdsByNetworkId))
	for networkId, gatewayIds := range gatewayIdsByNetworkId {
		mconfigs, err := store.GetMconfigs(networkId, gatewayIds)
		if err != nil {
			return map[string]NetworkStoredMconfigs{}, fmt.Errorf("Error loading stored mconfigs for network %s: %s", networkId, err)
		}
		ret[networkId] = mconfigs
	}
	return ret, nil
}

type NetworkBuiltMconfigs = map[string]*protos.GatewayConfigs

func buildMconfigsForNetworks(gatewayIdsByNetworkId map[string][]string) (map[string]NetworkBuiltMconfigs, error) {
	emptyRet := map[string]NetworkBuiltMconfigs{}
	ret := make(map[string]NetworkBuiltMconfigs, len(gatewayIdsByNetworkId))

	for networkId, gatewayIds := range gatewayIdsByNetworkId {
		innerRet := make(NetworkBuiltMconfigs, len(gatewayIds))
		for _, gatewayId := range gatewayIds {
			mcfg, err := factory.CreateMconfig(networkId, gatewayId)
			if err != nil {
				return emptyRet, fmt.Errorf("Error while building mconfig for (%s, %s): %s", networkId, gatewayId, err)
			}
			innerRet[gatewayId] = mcfg
		}
		ret[networkId] = innerRet
	}

	return ret, nil
}
