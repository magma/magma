/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"fmt"

	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func (m *Network) ToConfiguratorNetwork() configurator.Network {
	return configurator.Network{
		ID:          string(m.ID),
		Type:        string(m.Type),
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func (m *Network) FromConfiguratorNetwork(n configurator.Network) *Network {
	m.ID = models.NetworkID(n.ID)
	m.Type = models.NetworkType(n.Type)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg, exists := n.Configs[orc8r.DnsdNetworkType]; exists {
		m.DNS = cfg.(*NetworkDNSConfig)
	}
	if cfg, exists := n.Configs[orc8r.NetworkFeaturesConfig]; exists {
		m.Features = cfg.(*NetworkFeatures)
	}
	return m
}

func (m *Network) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewType:        swag.String(string(m.Type)),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func (m *NetworkFeatures) GetFromNetwork(network configurator.Network) interface{} {
	return GetNetworkConfig(network, orc8r.NetworkFeaturesConfig)
}

func (m *NetworkFeatures) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return GetNetworkConfigUpdateCriteria(network.ID, orc8r.NetworkFeaturesConfig, m), nil
}

func (m *NetworkDNSConfig) GetFromNetwork(network configurator.Network) interface{} {
	return GetNetworkConfig(network, orc8r.DnsdNetworkType)
}

func (m *NetworkDNSConfig) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return GetNetworkConfigUpdateCriteria(network.ID, orc8r.DnsdNetworkType, m), nil
}

func (m NetworkDNSRecords) GetFromNetwork(network configurator.Network) interface{} {
	iNetworkDnsConfig := GetNetworkConfig(network, orc8r.DnsdNetworkType)
	if iNetworkDnsConfig == nil {
		return nil
	}
	return iNetworkDnsConfig.(*NetworkDNSConfig).Records
}

func (m NetworkDNSRecords) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iNetworkDnsConfig := GetNetworkConfig(network, orc8r.DnsdNetworkType)
	if iNetworkDnsConfig == nil {
		return configurator.NetworkUpdateCriteria{}, errors.New("No DNS Config registered for this network")
	}
	iNetworkDnsConfig.(*NetworkDNSConfig).Records = m
	return GetNetworkConfigUpdateCriteria(network.ID, orc8r.DnsdNetworkType, iNetworkDnsConfig), nil
}

func (m *MagmadGateway) ToConfiguratorEntities() []configurator.NetworkEntity {
	gatewayEnt := configurator.NetworkEntity{
		Type:        orc8r.MagmadGatewayType,
		Key:         string(m.ID),
		Name:        string(m.Name),
		Description: string(m.Description),
		Config:      m.Magmad,
		PhysicalID:  m.Device.HardwareID,
	}
	return []configurator.NetworkEntity{gatewayEnt}
}

func (m *MagmadGateway) FromBackendModels(ent configurator.NetworkEntity, device *GatewayDevice, status *GatewayStatus) *MagmadGateway {
	m.ID = models.GatewayID(ent.Key)
	m.Name = models.GatewayName(ent.Name)
	m.Description = models.GatewayDescription(ent.Description)
	m.Magmad = ent.Config.(*MagmadGatewayConfigs)
	m.Device = device
	m.Status = status
	tierTK, err := ent.GetFirstParentOfType(orc8r.UpgradeTierEntityType)
	if err == nil {
		m.Tier = TierID(tierTK.Key)
	}

	return m
}

func (m *MagmadGateway) ToEntityUpdateCriteria(existingEnt configurator.NetworkEntity) []configurator.EntityUpdateCriteria {
	ret := []configurator.EntityUpdateCriteria{}
	gatewayUpdate := configurator.EntityUpdateCriteria{
		Type:      orc8r.MagmadGatewayType,
		Key:       string(m.ID),
		NewConfig: m.Magmad,
	}

	if m.Device.HardwareID != existingEnt.PhysicalID {
		gatewayUpdate.NewPhysicalID = swag.String(m.Device.HardwareID)
	}
	if string(m.Name) != existingEnt.Name {
		gatewayUpdate.NewName = swag.String(string(m.Name))
	}
	if string(m.Description) != existingEnt.Description {
		gatewayUpdate.NewDescription = swag.String(string(m.Description))
	}

	oldTierTK, _ := existingEnt.GetFirstParentOfType(orc8r.UpgradeTierEntityType)
	if oldTierTK.Key != string(m.Tier) {
		ret = append(
			ret,
			configurator.EntityUpdateCriteria{
				Type: orc8r.UpgradeTierEntityType, Key: oldTierTK.Key,
				AssociationsToDelete: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: string(m.ID)}},
			},
		)

		ret = append(
			ret,
			configurator.EntityUpdateCriteria{
				Type: orc8r.UpgradeTierEntityType, Key: string(m.Tier),
				AssociationsToAdd: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: string(m.ID)}},
			},
		)
	}

	// do the tier update to delete the old assoc first
	ret = append(ret, gatewayUpdate)
	return ret
}

func (m *MagmadGatewayConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	exists, err := configurator.DoesEntityExist(networkID, orc8r.MagmadGatewayType, gatewayID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Gateway %s does not exist", gatewayID)
	}

	return []configurator.EntityUpdateCriteria{
		{
			Key:       gatewayID,
			Type:      orc8r.MagmadGatewayType,
			NewConfig: m,
		},
	}, nil
}

func (m *MagmadGatewayConfigs) FromBackendModels(networkID string, gatewayID string) error {
	config, err := configurator.LoadEntityConfig(networkID, orc8r.MagmadGatewayType, gatewayID)
	if err != nil {
		return err
	}
	*m = *config.(*MagmadGatewayConfigs)
	return nil
}

func (m *TierID) FromBackendModels(networkID string, gatewayID string) error {
	entity, err := configurator.LoadEntity(networkID, orc8r.MagmadGatewayType, gatewayID, configurator.EntityLoadCriteria{LoadAssocsToThis: true})
	if err != nil {
		return err
	}
	for _, parentAssoc := range entity.ParentAssociations {
		if parentAssoc.Type == orc8r.UpgradeTierEntityType {
			*m = TierID(parentAssoc.Key)
			return nil
		}
	}
	return nil
}

func (m *TierID) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	tierID := string(*m)
	updateCriteria := []configurator.EntityUpdateCriteria{}

	exists, err := configurator.DoesEntityExist(networkID, orc8r.UpgradeTierEntityType, tierID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to look up tier")
	}
	if !exists {
		return nil, fmt.Errorf("Tier %s does not exist", tierID)
	}

	// Remove association from old tier
	entity, err := configurator.LoadEntity(networkID, orc8r.MagmadGatewayType, gatewayID, configurator.EntityLoadCriteria{LoadAssocsToThis: true})
	if err != nil {
		return nil, err
	}

	tierTK, err := entity.GetFirstParentOfType(orc8r.UpgradeTierEntityType)
	if err != merrors.ErrNotFound {
		if tierTK.Key == tierID {
			// no change
			return []configurator.EntityUpdateCriteria{}, nil
		}
		deleteCurrentTierAssoc := configurator.EntityUpdateCriteria{
			Type:                 tierTK.Type,
			Key:                  tierTK.Key,
			AssociationsToDelete: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gatewayID}},
		}
		updateCriteria = append(updateCriteria, deleteCurrentTierAssoc)
	}

	// Add association to new tier
	addNewTierAssoc := configurator.EntityUpdateCriteria{
		Type:              orc8r.UpgradeTierEntityType,
		Key:               tierID,
		AssociationsToAdd: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gatewayID}},
	}
	updateCriteria = append(updateCriteria, addNewTierAssoc)
	return updateCriteria, nil
}

func GetNetworkConfig(network configurator.Network, key string) interface{} {
	if network.Configs == nil {
		return nil
	}
	config, exists := network.Configs[key]
	if !exists {
		return nil
	}
	return config
}

func GetNetworkConfigUpdateCriteria(networkID string, key string, iConfig interface{}) configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID: networkID,
		ConfigsToAddOrUpdate: map[string]interface{}{
			key: iConfig,
		},
	}
}

func (m *Tier) ToNetworkEntity() configurator.NetworkEntity {
	return configurator.NetworkEntity{
		Type: orc8r.UpgradeTierEntityType, Key: string(m.ID),
		Name:         m.Name,
		Config:       m,
		Associations: getGatewayTKs(m.Gateways),
	}
}

func (m *Tier) ToUpdateCriteria() configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type: orc8r.UpgradeTierEntityType, Key: string(m.ID),
		NewName:           swag.String(m.Name),
		NewConfig:         m,
		AssociationsToSet: getGatewayTKs(m.Gateways),
	}
}

func (m *Tier) FromBackendModel(entity configurator.NetworkEntity) *Tier {
	tier := entity.Config.(*Tier)
	tier.Name = entity.Name
	tier.Gateways = getGatewayIDs(entity.Associations)
	return tier
}

func getGatewayTKs(gateways []models.GatewayID) []storage.TypeAndKey {
	return funk.Map(
		gateways,
		func(gw models.GatewayID) storage.TypeAndKey {
			return storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: string(gw)}
		}).([]storage.TypeAndKey)
}

func getGatewayIDs(gatewayTKs []storage.TypeAndKey) []models.GatewayID {
	return funk.Map(
		gatewayTKs,
		func(tk storage.TypeAndKey) models.GatewayID {
			return models.GatewayID(tk.Key)
		}).([]models.GatewayID)
}
