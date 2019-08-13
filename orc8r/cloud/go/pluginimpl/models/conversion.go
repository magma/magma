/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/swag"
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
	return fetchNetworkConfig(network, orc8r.NetworkFeaturesConfig)
}

func (m *NetworkFeatures) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return configurator.NetworkUpdateCriteria{
		ID: network.ID,
		ConfigsToAddOrUpdate: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: m,
		},
	}, nil
}

func (m *NetworkDNSConfig) GetFromNetwork(network configurator.Network) interface{} {
	return fetchNetworkConfig(network, orc8r.DnsdNetworkType)
}

func (m *NetworkDNSConfig) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return configurator.NetworkUpdateCriteria{
		ID: network.ID,
		ConfigsToAddOrUpdate: map[string]interface{}{
			orc8r.DnsdNetworkType: m,
		},
	}, nil
}

func (m *MagmadGateway) ToConfiguratorEntities() []configurator.NetworkEntity {
	gatewayEnt := configurator.NetworkEntity{
		Type:         orc8r.MagmadGatewayType,
		Key:          string(m.ID),
		Name:         string(m.Name),
		Description:  string(m.Description),
		Config:       m.Magmad,
		PhysicalID:   m.Device.HardwareID,
		Associations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: string(m.Tier)}},
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

func fetchNetworkConfig(network configurator.Network, key string) interface{} {
	if network.Configs == nil {
		return nil
	}
	config, exists := network.Configs[key]
	if !exists {
		return nil
	}
	return config
}
