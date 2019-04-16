/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package view_factory

import (
	"fmt"

	magmaerrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/checkind"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/magmad"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/storage"
)

// GatewayState represents the current state of a gateway, including
// information on configuration parameters, status, and record
type GatewayState struct {
	// ID of the gateway
	GatewayID string
	// Configuration of the gateway, represented as a map from configuration types
	// to configuration objects
	Config map[string]interface{}
	// Status of the gateway
	Status *protos.GatewayStatus
	// Gateway record
	Record *magmadprotos.AccessGatewayRecord
}

// GatewayUpdateParams contains information from an update to a gateway state. Each
// parameter is nullable, and only non-null parameters will be used to update the gateway
// state
type GatewayUpdateParams struct {
	// Only the keys in this NewConfig map will be used to update the config of the GatewayState
	NewConfig map[string]interface{}
	// Configurations to delete
	ConfigsToDelete []string
	// New status of the gateway
	NewStatus *protos.GatewayStatus
	// New gateway record
	NewRecord *magmadprotos.AccessGatewayRecord
}

// FullGatewayViewFactory constructs `GatewayState`s for specified gateways
// within a network.
type FullGatewayViewFactory interface {
	// Get the states of all gateways in this network
	GetGatewayViewsForNetwork(networkID string) (map[string]*GatewayState, error)
	// Get the state of specific gateways
	GetGatewayViews(networkID string, gatewayIDs []string) (map[string]*GatewayState, error)
}

// FullGatewayViewFactoryImpl is the default implementation of
// FullGatewayViewFactory which uses service client APIs to construct
// `GatewayState`s
type FullGatewayViewFactoryImpl struct{}

func (f *FullGatewayViewFactoryImpl) GetGatewayViewsForNetwork(networkID string) (map[string]*GatewayState, error) {
	gatewayIDs, err := magmad.ListGateways(networkID)
	if err != nil {
		return map[string]*GatewayState{}, fmt.Errorf("Error loading gateway IDs for network view: %s", err)
	}
	return f.GetGatewayViews(networkID, gatewayIDs)
}

func (f *FullGatewayViewFactoryImpl) GetGatewayViews(networkID string, gatewayIDs []string) (map[string]*GatewayState, error) {
	ret := make(map[string]*GatewayState, len(gatewayIDs))
	for _, gatewayID := range gatewayIDs {
		state, err := loadGatewayView(networkID, gatewayID)
		if err != nil {
			return map[string]*GatewayState{}, fmt.Errorf("Error loading gateway %s view: %s", gatewayID, err)
		}
		ret[gatewayID] = state
	}
	return ret, nil
}

func loadGatewayView(networkID string, gatewayID string) (*GatewayState, error) {
	record, err := magmad.FindGatewayRecord(networkID, gatewayID)
	if err != nil {
		return nil, fmt.Errorf("Error loading record: %s", err)
	}
	configs, err := config.GetConfigsByKey(networkID, gatewayID)
	if err != nil {
		return nil, fmt.Errorf("Error loading configs: %s", err)
	}
	status, err := loadStatus(networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	return createGatewayView(gatewayID, record, configs, status), nil
}

func loadStatus(networkID string, gatewayID string) (*protos.GatewayStatus, error) {
	status, err := checkind.GetStatus(networkID, gatewayID)
	if err == magmaerrors.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Error loading status: %s", err)
	}
	return status, nil
}

type gatewayConfigs map[storage.TypeAndKey]interface{}

func createGatewayView(gatewayID string, record *magmadprotos.AccessGatewayRecord, configs gatewayConfigs, status *protos.GatewayStatus) *GatewayState {
	ret := &GatewayState{
		GatewayID: gatewayID,
		Status:    status,
		Record:    record,
		Config:    make(map[string]interface{}, len(configs)),
	}
	for typeAndKey, configVal := range configs {
		ret.Config[typeAndKey.Type] = configVal
	}
	return ret
}
