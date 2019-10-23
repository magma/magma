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
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/swag"
	"github.com/thoas/go-funk"
)

// GatewayState represents the current state of a gateway, including
// information on configuration parameters, status, and record
type GatewayState struct {
	// ID of the gateway
	GatewayID string `json:"gateway_id"`
	// Configuration of the gateway, represented as a map from configuration types
	// to configuration objects
	Config map[string]interface{} `json:"config"`
	// Name of the gateway
	Name string `json:"name"`
	// Gateway record
	Record *models.GatewayDevice `json:"record"`
	// Status of the gateway
	Status *models.GatewayStatus `json:"status"`
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
	gatewayIDs, err := configurator.ListEntityKeys(networkID, orc8r.MagmadGatewayType)
	if err != nil {
		return map[string]*GatewayState{}, fmt.Errorf("Error loading gateway IDs for network view: %s", err)
	}
	return f.GetGatewayViews(networkID, gatewayIDs)
}

func (f *FullGatewayViewFactoryImpl) GetGatewayViews(networkID string, gatewayIDs []string) (map[string]*GatewayState, error) {
	ret := make(map[string]*GatewayState, len(gatewayIDs))
	gatewayTKs := funk.Map(gatewayIDs, func(id string) storage.TypeAndKey { return storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: id} }).([]storage.TypeAndKey)
	loadedGateways, _, err := configurator.LoadEntities(networkID, swag.String(orc8r.MagmadGatewayType), nil, nil, gatewayTKs, configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true, LoadMetadata: true})
	if err != nil {
		return nil, err
	}
	for _, gateway := range loadedGateways {
		record, err := device.GetDevice(networkID, orc8r.AccessGatewayRecordType, gateway.PhysicalID)
		if err != nil {
			return nil, fmt.Errorf("Error loading record: %s", err)
		}

		status, err := state.GetGatewayStatus(networkID, gateway.PhysicalID)
		if err == magmaerrors.ErrNotFound {
			status = nil
		} else if err != nil {
			return nil, fmt.Errorf("Error loading status: %s", err)
		}

		gatewayRecord := record.(*models.GatewayDevice)

		ret[gateway.Key] = &GatewayState{
			GatewayID: gateway.Key,
			Name:      gateway.Name,
			Record:    gatewayRecord,
			Status:    status,
			Config:    map[string]interface{}{orc8r.MagmadGatewayType: gateway.Config},
		}
	}

	// load all associated configEntity entities
	allAssociations := getAllAssociatedConfigEntities(loadedGateways)
	if len(allAssociations) == 0 {
		// if allAssociations is length 0 the call below will load all entities
		return ret, nil
	}
	loadedConfigs, _, err := configurator.LoadEntities(networkID, nil, nil, nil, allAssociations, configurator.EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return nil, err
	}
	for _, configEntity := range loadedConfigs {
		ret[configEntity.Key].Config[configEntity.Type] = configEntity.Config
	}
	return ret, nil
}

// Relies on config entities sharing its key with the parent gateway entity
func getAllAssociatedConfigEntities(queriedGateways []configurator.NetworkEntity) []storage.TypeAndKey {
	ret := []storage.TypeAndKey{}
	for _, gatewayEnt := range queriedGateways {
		for _, associatedEnt := range gatewayEnt.Associations {
			if associatedEnt.Key == gatewayEnt.Key {
				ret = append(ret, associatedEnt)
			}
		}
	}
	return ret
}
