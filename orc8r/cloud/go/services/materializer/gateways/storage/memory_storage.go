/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"errors"
	"fmt"
	"strings"
)

// MemoryGatewayViewStorage is an in-memory implementation of the GatewayViewStorage interface
type MemoryGatewayViewStorage struct {
	Networks map[string]map[string]*GatewayState
}

func networkNotFound(networkID string) error {
	return fmt.Errorf("Network ID not found: %s", networkID)
}

func gatewayNotFound(gatewayID string) error {
	return fmt.Errorf("Gateway ID not found: %s", gatewayID)
}

// NewMemoryGatewayViewStorage creates a new gateway, initializing the map
func NewMemoryGatewayViewStorage() GatewayViewStorage {
	return &MemoryGatewayViewStorage{
		Networks: make(map[string]map[string]*GatewayState),
	}
}

func (view *MemoryGatewayViewStorage) InitTables() error {
	// No work needs to be done for in-memory storage
	return nil
}

func (view *MemoryGatewayViewStorage) GetGatewayViewsForNetwork(networkID string) (map[string]*GatewayState, error) {
	if network, ok := view.Networks[networkID]; ok {
		return network, nil
	}
	return nil, networkNotFound(networkID)
}

func (view *MemoryGatewayViewStorage) GetGatewayViews(networkID string, gatewayIDs []string) (map[string]*GatewayState, error) {
	network, ok := view.Networks[strings.ToLower(networkID)]
	if !ok {
		return nil, networkNotFound(networkID)
	}
	gatewayStates := make(map[string]*GatewayState)
	for _, gatewayID := range gatewayIDs {
		gateway, ok := network[gatewayID]
		if !ok {
			return nil, gatewayNotFound(gatewayID)
		}
		gatewayStates[gatewayID] = gateway
	}
	return gatewayStates, nil
}

func (view *MemoryGatewayViewStorage) UpdateOrCreateGatewayViews(networkID string, updates map[string]*GatewayUpdateParams) error {
	network, ok := view.Networks[networkID]
	if !ok {
		network = make(map[string]*GatewayState)
		view.Networks[networkID] = network
	}
	for gatewayID, params := range updates {
		gatewayState, ok := network[gatewayID]
		if !ok {
			gatewayState = &GatewayState{
				Config: make(map[string]interface{}),
			}
			network[gatewayID] = gatewayState
		}
		if params.Offset <= gatewayState.Offset {
			return errors.New("Update offset less than current state")
		}
		if params.NewConfig != nil {
			for key, value := range params.NewConfig {
				gatewayState.Config[key] = value
			}
		}
		if params.ConfigsToDelete != nil {
			for _, configType := range params.ConfigsToDelete {
				delete(gatewayState.Config, configType)
			}
		}
		if params.NewRecord != nil {
			gatewayState.Record = params.NewRecord
		}
		if params.NewStatus != nil {
			gatewayState.Status = params.NewStatus
		}
		gatewayState.Offset = params.Offset
	}
	return nil
}

func (view *MemoryGatewayViewStorage) DeleteGatewayViews(networkID string, gatewayIDs []string) error {
	network, ok := view.Networks[networkID]
	if !ok {
		return networkNotFound(networkID)
	}
	for _, gatewayID := range gatewayIDs {
		_, ok := network[gatewayID]
		if !ok {
			return gatewayNotFound(gatewayID)
		}
		delete(network, gatewayID)
	}
	return nil
}
