/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/lte/cloud/go/lte"
	ltemodels "magma/lte/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	models2 "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func (m *FegNetwork) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *FegNetwork) GetEmptyNetwork() handlers.NetworkModel {
	return &FegNetwork{}
}

func (m *FegNetwork) ToConfiguratorNetwork() configurator.Network {
	return configurator.Network{
		ID:          string(m.ID),
		Type:        feg.FederationNetworkType,
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			feg.FegNetworkType:          m.Federation,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func (m *FegNetwork) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			feg.FegNetworkType:          m.Federation,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func (m *FegNetwork) FromConfiguratorNetwork(n configurator.Network) interface{} {
	m.ID = models.NetworkID(n.ID)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[feg.FegNetworkType]; cfg != nil {
		m.Federation = cfg.(*NetworkFederationConfigs)
	}
	if cfg := n.Configs[orc8r.DnsdNetworkType]; cfg != nil {
		m.DNS = cfg.(*models2.NetworkDNSConfig)
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*models2.NetworkFeatures)
	}
	return m
}

func (m *FegLteNetwork) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *FegLteNetwork) GetEmptyNetwork() handlers.NetworkModel {
	return &FegLteNetwork{}
}

func (m *FegLteNetwork) ToConfiguratorNetwork() configurator.Network {
	return configurator.Network{
		ID:          string(m.ID),
		Type:        feg.FederatedLteNetworkType,
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			lte.CellularNetworkType:     m.Cellular,
			feg.FederatedNetworkType:    m.Federation,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func (m *FegLteNetwork) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			lte.CellularNetworkType:     m.Cellular,
			feg.FederatedNetworkType:    m.Federation,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func (m *FegLteNetwork) FromConfiguratorNetwork(n configurator.Network) interface{} {
	m.ID = models.NetworkID(n.ID)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[feg.FederatedNetworkType]; cfg != nil {
		m.Federation = cfg.(*FederatedNetworkConfigs)
	}
	if cfg := n.Configs[lte.CellularNetworkType]; cfg != nil {
		m.Cellular = cfg.(*ltemodels.NetworkCellularConfigs)
	}
	if cfg := n.Configs[orc8r.DnsdNetworkType]; cfg != nil {
		m.DNS = cfg.(*models2.NetworkDNSConfig)
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*models2.NetworkFeatures)
	}
	return m
}

func (m *NetworkFederationConfigs) GetFromNetwork(network configurator.Network) interface{} {
	return models2.GetNetworkConfig(network, feg.FederatedNetworkType)
}

func (m *NetworkFederationConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return models2.GetNetworkConfigUpdateCriteria(network.ID, feg.FederatedNetworkType, m), nil
}

func (m *FederationGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *FederationGateway) FromBackendModels(
	magmadGateway, federationGateway configurator.NetworkEntity,
	device *models2.GatewayDevice,
	status *models2.GatewayStatus,
) handlers.GatewayModel {
	// delegate most of the fillin to magmad gateway struct
	mdGW := (&models2.MagmadGateway{}).FromBackendModels(magmadGateway, device, status)
	// TODO: we should change this to a reflection based shallow copy
	m.ID, m.Name, m.Description, m.Magmad, m.Tier, m.Device, m.Status = mdGW.ID, mdGW.Name, mdGW.Description, mdGW.Magmad, mdGW.Tier, mdGW.Device, mdGW.Status
	m.Federation = federationGateway.Config.(*GatewayFederationConfigs)
	return m
}

func (m *MutableFederationGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MutableFederationGateway) GetEmptyGateway() handlers.MutableGatewayModel {
	return &MutableFederationGateway{}
}

func (m *MutableFederationGateway) GetMagmadGateway() *models2.MagmadGateway {
	return &models2.MagmadGateway{
		Description: m.Description,
		Device:      m.Device,
		ID:          m.ID,
		Magmad:      m.Magmad,
		Name:        m.Name,
		Tier:        m.Tier,
	}
}

func (m *FederatedNetworkConfigs) GetFromNetwork(network configurator.Network) interface{} {
	return models2.GetNetworkConfig(network, feg.FederatedNetworkType)
}

func (m *FederatedNetworkConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return models2.GetNetworkConfigUpdateCriteria(network.ID, feg.FederatedNetworkType, m), nil
}

func (m *MutableFederationGateway) ToConfiguratorEntity() configurator.NetworkEntity {
	ret := configurator.NetworkEntity{
		Type:        feg.FegGatewayType,
		Key:         string(m.ID),
		Name:        string(m.Name),
		Description: string(m.Description),
		Config:      m.Federation,
	}
	return ret
}

func (m *MutableFederationGateway) GetMagmadGatewayUpdateCriteria() configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type:              orc8r.MagmadGatewayType,
		Key:               string(m.ID),
		AssociationsToAdd: []storage.TypeAndKey{{Type: feg.FegGatewayType, Key: string(m.ID)}},
	}
}

func (m *MutableFederationGateway) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type:           feg.FegGatewayType,
		Key:            string(m.ID),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		NewConfig:      m.Federation,
	}
	return ret
}

func (m *GatewayFederationConfigs) FromBackendModels(networkID string, gatewayID string) error {
	federationConfig, err := configurator.LoadEntityConfig(networkID, feg.FegGatewayType, gatewayID)
	if err != nil {
		return err
	}
	*m = *federationConfig.(*GatewayFederationConfigs)
	return nil
}

func (m *GatewayFederationConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Type: feg.FegGatewayType, Key: gatewayID,
			NewConfig: m,
		},
	}, nil
}

func (config *DiameterClientConfigs) ToMconfig() *mconfig.DiamClientConfig {
	res := &mconfig.DiamClientConfig{}
	protos.FillIn(config, res)
	return res
}

func (config *DiameterServerConfigs) ToMconfig() *mconfig.DiamServerConfig {
	res := &mconfig.DiamServerConfig{}
	protos.FillIn(config, res)
	return res
}

func (profile *SubscriptionProfile) ToMconfig() *mconfig.HSSConfig_SubscriptionProfile {
	res := &mconfig.HSSConfig_SubscriptionProfile{}
	protos.FillIn(profile, res)
	return res
}
