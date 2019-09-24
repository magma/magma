/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"magma/cwf/cloud/go/cwf"
	"magma/feg/cloud/go/feg"
	models3 "magma/feg/cloud/go/plugin/models"
	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	models2 "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func (m *CwfNetwork) GetEmptyNetwork() handlers.NetworkModel {
	return &CwfNetwork{}
}

func (m *CwfNetwork) ToConfiguratorNetwork() configurator.Network {
	return configurator.Network{
		ID:          string(m.ID),
		Type:        cwf.CwfNetworkType,
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			cwf.CwfNetworkType:          m.CarrierWifi,
			feg.FederatedNetworkType:    m.Federation,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func (m *CwfNetwork) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			cwf.CwfNetworkType:          m.CarrierWifi,
			feg.FederatedNetworkType:    m.Federation,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func (m *CwfNetwork) FromConfiguratorNetwork(n configurator.Network) interface{} {
	m.ID = models.NetworkID(n.ID)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[cwf.CwfNetworkType]; cfg != nil {
		m.CarrierWifi = cfg.(*NetworkCarrierWifiConfigs)
	}
	if cfg := n.Configs[feg.FederatedNetworkType]; cfg != nil {
		m.Federation = cfg.(*models3.FederatedNetworkConfigs)
	}
	if cfg := n.Configs[orc8r.DnsdNetworkType]; cfg != nil {
		m.DNS = cfg.(*models2.NetworkDNSConfig)
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*models2.NetworkFeatures)
	}
	return m
}

func (m *CwfGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *CwfGateway) FromBackendModels(
	magmadGateway, cellularGateway configurator.NetworkEntity,
	device *models2.GatewayDevice,
	status *models2.GatewayStatus,
) handlers.GatewayModel {
	// delegate most of the fillin to magmad gateway struct
	mdGW := (&models2.MagmadGateway{}).FromBackendModels(magmadGateway, device, status)
	// TODO: we should change this to a reflection based shallow copy
	m.ID, m.Name, m.Description, m.Magmad, m.Tier, m.Device, m.Status = mdGW.ID, mdGW.Name, mdGW.Description, mdGW.Magmad, mdGW.Tier, mdGW.Device, mdGW.Status
	return m
}

func (m *MutableCwfGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MutableCwfGateway) GetMagmadGateway() *models2.MagmadGateway {
	return &models2.MagmadGateway{
		Description: m.Description,
		Device:      m.Device,
		ID:          m.ID,
		Magmad:      m.Magmad,
		Name:        m.Name,
		Tier:        m.Tier,
	}
}

func (m *MutableCwfGateway) GetAdditionalWritesOnCreate() []configurator.EntityWriteOperation {
	return []configurator.EntityWriteOperation{
		configurator.NetworkEntity{
			Type:        cwf.CwfGatewayType,
			Key:         string(m.ID),
			Name:        string(m.Name),
			Description: string(m.Description),
			Config:      nil,
		},
		configurator.EntityUpdateCriteria{
			Type:              orc8r.MagmadGatewayType,
			Key:               string(m.ID),
			AssociationsToAdd: []storage.TypeAndKey{{Type: cwf.CwfGatewayType, Key: string(m.ID)}},
		},
	}
}

func (m *MutableCwfGateway) GetAdditionalEntitiesToLoadOnUpdate(gatewayID string) []storage.TypeAndKey {
	return []storage.TypeAndKey{{Type: cwf.CwfGatewayType, Key: gatewayID}}
}

func (m *MutableCwfGateway) GetAdditionalWritesOnUpdate(
	gatewayID string,
	loadedEntities map[storage.TypeAndKey]configurator.NetworkEntity,
) ([]configurator.EntityWriteOperation, error) {
	ret := []configurator.EntityWriteOperation{}
	existingEnt, ok := loadedEntities[storage.TypeAndKey{Type: cwf.CwfGatewayType, Key: gatewayID}]
	if !ok {
		return ret, merrors.ErrNotFound
	}

	entUpdate := configurator.EntityUpdateCriteria{
		Type: cwf.CwfGatewayType,
		Key:  string(m.ID),
	}
	if string(m.Name) != existingEnt.Name {
		entUpdate.NewName = swag.String(string(m.Name))
	}
	if string(m.Description) != existingEnt.Description {
		entUpdate.NewDescription = swag.String(string(m.Description))
	}

	ret = append(ret, entUpdate)
	return ret, nil
}

func (m *NetworkCarrierWifiConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return models2.GetNetworkConfigUpdateCriteria(network.ID, cwf.CwfNetworkType, m), nil
}

func (m *NetworkCarrierWifiConfigs) GetFromNetwork(network configurator.Network) interface{} {
	return models2.GetNetworkConfig(network, cwf.CwfNetworkType)
}
