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
	lteModels "magma/lte/cloud/go/plugin/models"
	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	orc8rModels "magma/orc8r/cloud/go/pluginimpl/models"
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
	network := configurator.Network{
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
	if m.SubscriberConfig != nil {
		network.Configs[lte.NetworkSubscriberConfigType] = m.SubscriberConfig
	}
	return network
}

func (m *FegNetwork) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	update := configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			feg.FegNetworkType:          m.Federation,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
	if m.SubscriberConfig != nil {
		update.ConfigsToAddOrUpdate[lte.NetworkSubscriberConfigType] = m.SubscriberConfig
	}
	return update
}

func (m *FegNetwork) FromConfiguratorNetwork(n configurator.Network) interface{} {
	m.ID = models.NetworkID(n.ID)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[feg.FegNetworkType]; cfg != nil {
		m.Federation = cfg.(*NetworkFederationConfigs)
	}
	if cfg := n.Configs[orc8r.DnsdNetworkType]; cfg != nil {
		m.DNS = cfg.(*orc8rModels.NetworkDNSConfig)
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*orc8rModels.NetworkFeatures)
	}
	if cfg := n.Configs[lte.NetworkSubscriberConfigType]; cfg != nil {
		m.SubscriberConfig = cfg.(*lteModels.NetworkSubscriberConfig)
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
		m.Cellular = cfg.(*lteModels.NetworkCellularConfigs)
	}
	if cfg := n.Configs[orc8r.DnsdNetworkType]; cfg != nil {
		m.DNS = cfg.(*orc8rModels.NetworkDNSConfig)
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*orc8rModels.NetworkFeatures)
	}
	return m
}

func (m *NetworkFederationConfigs) GetFromNetwork(network configurator.Network) interface{} {
	return orc8rModels.GetNetworkConfig(network, feg.FegNetworkType)
}

func (m *NetworkFederationConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, feg.FegNetworkType, m), nil
}

func (m *FederationGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *FederationGateway) FromBackendModels(
	magmadGateway, federationGateway configurator.NetworkEntity,
	device *orc8rModels.GatewayDevice,
	status *orc8rModels.GatewayStatus,
) handlers.GatewayModel {
	// delegate most of the fillin to magmad gateway struct
	mdGW := (&orc8rModels.MagmadGateway{}).FromBackendModels(magmadGateway, device, status)
	// TODO: we should change this to a reflection based shallow copy
	m.ID, m.Name, m.Description, m.Magmad, m.Tier, m.Device, m.Status = mdGW.ID, mdGW.Name, mdGW.Description, mdGW.Magmad, mdGW.Tier, mdGW.Device, mdGW.Status
	m.Federation = federationGateway.Config.(*GatewayFederationConfigs)
	return m
}

func (m *MutableFederationGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MutableFederationGateway) GetMagmadGateway() *orc8rModels.MagmadGateway {
	return &orc8rModels.MagmadGateway{
		Description: m.Description,
		Device:      m.Device,
		ID:          m.ID,
		Magmad:      m.Magmad,
		Name:        m.Name,
		Tier:        m.Tier,
	}
}

func (m *MutableFederationGateway) GetAdditionalWritesOnCreate() []configurator.EntityWriteOperation {
	return []configurator.EntityWriteOperation{
		configurator.NetworkEntity{
			Type:        feg.FegGatewayType,
			Key:         string(m.ID),
			Name:        string(m.Name),
			Description: string(m.Description),
			Config:      m.Federation,
		},
		configurator.EntityUpdateCriteria{
			Type:              orc8r.MagmadGatewayType,
			Key:               string(m.ID),
			AssociationsToAdd: []storage.TypeAndKey{{Type: feg.FegGatewayType, Key: string(m.ID)}},
		},
	}
}

func (m *MutableFederationGateway) GetAdditionalEntitiesToLoadOnUpdate(gatewayID string) []storage.TypeAndKey {
	return []storage.TypeAndKey{{Type: feg.FegGatewayType, Key: gatewayID}}
}

func (m *MutableFederationGateway) GetAdditionalWritesOnUpdate(
	gatewayID string,
	loadedEntities map[storage.TypeAndKey]configurator.NetworkEntity,
) ([]configurator.EntityWriteOperation, error) {
	ret := []configurator.EntityWriteOperation{}
	existingEnt, ok := loadedEntities[storage.TypeAndKey{Type: feg.FegGatewayType, Key: gatewayID}]
	if !ok {
		return ret, merrors.ErrNotFound
	}

	entUpdate := configurator.EntityUpdateCriteria{
		Type:      feg.FegGatewayType,
		Key:       string(m.ID),
		NewConfig: m.Federation,
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

func (m *FederatedNetworkConfigs) GetFromNetwork(network configurator.Network) interface{} {
	return orc8rModels.GetNetworkConfig(network, feg.FederatedNetworkType)
}

func (m *FederatedNetworkConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, feg.FederatedNetworkType, m), nil
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

func (m *DiameterClientConfigs) ToMconfig() *mconfig.DiamClientConfig {
	res := &mconfig.DiamClientConfig{}
	protos.FillIn(m, res)
	return res
}

func (m *DiameterServerConfigs) ToMconfig() *mconfig.DiamServerConfig {
	res := &mconfig.DiamServerConfig{}
	protos.FillIn(m, res)
	return res
}

func (m *SubscriptionProfile) ToMconfig() *mconfig.HSSConfig_SubscriptionProfile {
	res := &mconfig.HSSConfig_SubscriptionProfile{}
	protos.FillIn(m, res)
	return res
}
