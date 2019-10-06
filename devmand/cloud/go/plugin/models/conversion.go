/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package models

import (
	"sort"

	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	models2 "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	"orc8r/devmand/cloud/go/devmand"

	"github.com/go-openapi/swag"
)

func (m *SymphonyNetwork) ToConfiguratorNetwork() configurator.Network {
	return configurator.Network{
		ID:          string(m.ID),
		Type:        devmand.SymphonyNetworkType,
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func (m *SymphonyNetwork) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func (m *SymphonyNetwork) FromConfiguratorNetwork(n configurator.Network) *SymphonyNetwork {
	m.ID = models.NetworkID(n.ID)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*models2.NetworkFeatures)
	}
	return m
}

func (m *SymphonyAgent) GetMagmadGateway() *models2.MagmadGateway {
	return &models2.MagmadGateway{
		Description: m.Description,
		Device:      m.Device,
		ID:          models.GatewayID(m.ID),
		Magmad:      m.Magmad,
		Name:        m.Name,
		Tier:        m.Tier,
	}
}

func (m *SymphonyAgent) GetAdditionalWritesOnCreate() []configurator.EntityWriteOperation {
	ent := configurator.NetworkEntity{
		Type: devmand.SymphonyAgentType,
		Key:  string(m.ID),
	}

	for _, managedDevice := range m.ManagedDevices {
		ent.Associations = append(ent.Associations, storage.TypeAndKey{Type: devmand.SymphonyDeviceType, Key: managedDevice})
	}

	return []configurator.EntityWriteOperation{
		ent,
		configurator.EntityUpdateCriteria{
			Type:              orc8r.MagmadGatewayType,
			Key:               string(m.ID),
			AssociationsToAdd: []storage.TypeAndKey{{Type: devmand.SymphonyAgentType, Key: string(m.ID)}},
		},
	}
}

func (m *SymphonyAgent) GetAdditionalEntitiesToLoadOnUpdate(agentID string) []storage.TypeAndKey {
	return []storage.TypeAndKey{{Type: devmand.SymphonyAgentType, Key: agentID}}
}

func (m *SymphonyAgent) GetAdditionalWritesOnUpdate(
	agentID string,
	loadedEntities map[storage.TypeAndKey]configurator.NetworkEntity,
) ([]configurator.EntityWriteOperation, error) {
	ret := []configurator.EntityWriteOperation{}
	_, ok := loadedEntities[storage.TypeAndKey{Type: devmand.SymphonyAgentType, Key: agentID}]
	if !ok {
		return ret, merrors.ErrNotFound
	}

	entUpdate := configurator.EntityUpdateCriteria{
		Type: devmand.SymphonyAgentType,
		Key:  string(m.ID),
	}

	for _, managedDevice := range m.ManagedDevices {
		entUpdate.AssociationsToSet = append(entUpdate.AssociationsToSet, storage.TypeAndKey{Type: devmand.SymphonyDeviceType, Key: managedDevice})
	}

	ret = append(ret, entUpdate)
	return ret, nil
}

func (m *SymphonyAgent) ToConfiguratorEntity() configurator.NetworkEntity {
	ret := configurator.NetworkEntity{
		Type: devmand.SymphonyAgentType,
		Key:  string(m.ID),
	}
	for _, managedDevice := range m.ManagedDevices {
		ret.Associations = append(ret.Associations, storage.TypeAndKey{Type: devmand.SymphonyDeviceType, Key: managedDevice})
	}
	return ret
}

func (m *SymphonyAgent) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type: devmand.SymphonyAgentType,
		Key:  string(m.ID),
	}
	for _, managedDeviceID := range m.ManagedDevices {
		ret.AssociationsToSet = append(ret.AssociationsToSet, storage.TypeAndKey{Type: devmand.SymphonyDeviceType, Key: managedDeviceID})
	}
	return ret
}

func (m *SymphonyAgent) FromBackendModels(
	magmadEnt, agentEnt configurator.NetworkEntity,
	device *models2.GatewayDevice,
	status *models2.GatewayStatus,
) handlers.GatewayModel {
	mdGW := (&models2.MagmadGateway{}).FromBackendModels(magmadEnt, device, status)
	m.ID, m.Name, m.Description, m.Magmad, m.Tier, m.Device, m.Status = mdGW.ID, mdGW.Name, mdGW.Description, mdGW.Magmad, mdGW.Tier, mdGW.Device, mdGW.Status

	for _, tk := range agentEnt.Associations {
		if tk.Type == devmand.SymphonyDeviceType {
			m.ManagedDevices = append(m.ManagedDevices, tk.Key)
		}
	}
	sort.Strings(m.ManagedDevices)

	return m
}

func (m *ManagedDevices) FromBackendModels(networkID string, agentID string) error {
	symphonyAgentEntity, err := configurator.LoadEntity(networkID, devmand.SymphonyAgentType, agentID, configurator.EntityLoadCriteria{LoadAssocsFromThis: true})
	if err != nil {
		return err
	}
	managedDevices := ManagedDevices{}
	for _, assoc := range symphonyAgentEntity.Associations {
		if assoc.Type == devmand.SymphonyDeviceType {
			managedDevices = append(managedDevices, assoc.Key)
		}
	}
	*m = managedDevices
	return nil
}

func (m *ManagedDevices) ToUpdateCriteria(networkID string, agentID string) ([]configurator.EntityUpdateCriteria, error) {
	managedDevices := []storage.TypeAndKey{}
	for _, managedDevice := range *m {
		managedDevices = append(managedDevices, storage.TypeAndKey{Type: devmand.SymphonyDeviceType, Key: managedDevice})
	}
	return []configurator.EntityUpdateCriteria{
		{
			Type:              devmand.SymphonyAgentType,
			Key:               agentID,
			AssociationsToSet: managedDevices,
		},
	}, nil
}

func (m *ManagedDevices) ToDeleteUpdateCriteria(networkID, agentID, deviceID string) configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type: devmand.SymphonyAgentType, Key: agentID,
		AssociationsToDelete: []storage.TypeAndKey{{Type: devmand.SymphonyDeviceType, Key: deviceID}},
	}
}

func (m *ManagedDevices) ToCreateUpdateCriteria(networkID, agentID, deviceID string) configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type: devmand.SymphonyAgentType, Key: agentID,
		AssociationsToAdd: []storage.TypeAndKey{{Type: devmand.SymphonyDeviceType, Key: deviceID}},
	}
}
