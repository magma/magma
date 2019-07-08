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
	"magma/orc8r/cloud/go/protos"
	checkind_models "magma/orc8r/cloud/go/services/checkind/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/magmad/obsidian/models"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	stateh "magma/orc8r/cloud/go/services/state/obsidian/handlers"
	"magma/orc8r/cloud/go/storage"

	"github.com/thoas/go-funk"
)

// GatewayState represents the current state of a gateway, including
// information on configuration parameters, status, and record
type GatewayState struct {
	// ID of the gateway
	GatewayID string
	// Configuration of the gateway, represented as a map from configuration types
	// to configuration objects
	Config map[string]interface{}

	// Gateway record
	Record *models.AccessGatewayRecord
	// Status of the gateway
	Status *checkind_models.GatewayStatus

	// Gateway record
	LegacyRecord *magmadprotos.AccessGatewayRecord // Deprecated
	// Status of the gateway
	LegacyStatus *protos.GatewayStatus // Deprecated
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
	loadedGateways, _, err := configurator.LoadEntities(networkID, nil, nil, nil, gatewayTKs, configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true})
	if err != nil {
		return nil, err
	}
	for _, gateway := range loadedGateways {
		record, err := device.GetDevice(networkID, orc8r.AccessGatewayRecordType, gateway.PhysicalID)
		if err != nil {
			return nil, fmt.Errorf("Error loading record: %s", err)
		}

		status, err := stateh.GetGWStatus(networkID, gateway.PhysicalID)
		if err == magmaerrors.ErrNotFound {
			status = nil
		} else if err != nil {
			return nil, fmt.Errorf("Error loading status: %s", err)
		}

		ret[gateway.Key] = &GatewayState{
			GatewayID: gateway.Key,
			Record:    record.(*models.AccessGatewayRecord),
			Status:    status,
			Config:    map[string]interface{}{orc8r.MagmadGatewayType: gateway.Config},
		}
	}

	// load all associated config entities
	allAssociations := getAllAssociations(loadedGateways)
	loadedConfigs, _, err := configurator.LoadEntities(networkID, nil, nil, nil, allAssociations, configurator.EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return nil, err
	}
	for _, config := range loadedConfigs {
		ret[config.Key].Config[config.Type] = config.Config
	}
	return ret, nil
}

func getAllAssociations(gateways []configurator.NetworkEntity) []storage.TypeAndKey {
	ret := []storage.TypeAndKey{}
	for _, gateway := range gateways {
		ret = append(ret, gateway.Associations...)
	}
	return ret
}
