/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/services/checkind"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/config/blacklist"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
)

func loadGatewayStatesForNetwork(networkID string) (map[string]*storage.GatewayState, error) {
	gatewayIDs, err := magmad.ListGateways(networkID)
	if err != nil {
		return nil, fmt.Errorf("Error loading gateways for network %s: %s", networkID, err)
	}
	states := make(map[string]*storage.GatewayState)
	for _, gatewayID := range gatewayIDs {
		state, err := loadGatewayState(networkID, gatewayID)
		if err != nil {
			return nil, err
		}
		states[gatewayID] = state
	}
	return states, nil
}

func loadGatewayState(networkID string, gatewayID string) (*storage.GatewayState, error) {
	record, err := magmad.FindGatewayRecord(networkID, gatewayID)
	if err != nil {
		return nil, fmt.Errorf("Error loading gateway record for %s: %s", gatewayID, err)
	}
	status, err := checkind.GetStatus(networkID, gatewayID)
	if err != nil {
		if datastore.IsErrNotFound(err) {
			status = nil
		} else {
			return nil, fmt.Errorf("Error loading gateway status for %s: %s", gatewayID, err)
		}
	}
	configs, err := loadConfigs(networkID, gatewayID)
	if err != nil {
		return nil, fmt.Errorf("Error loading configs for %s: %s", gatewayID, err)
	}
	return &storage.GatewayState{
		GatewayID: gatewayID,
		Config:    configs,
		Status:    status,
		Record:    record,
	}, nil
}

func loadConfigs(networkID string, gatewayID string) (map[string]interface{}, error) {
	configs := make(map[string]interface{})
	gatewayConfigs, err := config.GetConfigsByKey(networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	for typeAndKey, configObj := range gatewayConfigs {
		if !blacklist.IsConfigBlacklisted(typeAndKey.Type) {
			configs[typeAndKey.Type] = configObj
		}
	}
	return configs, nil
}
