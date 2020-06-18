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

	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	models2 "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
	"orc8r/devmand/cloud/go/devmand"

	"github.com/go-openapi/swag"
)

func (m *SymphonyNetwork) GetEmptyNetwork() handlers.NetworkModel {
	return &SymphonyNetwork{}
}

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

func (m *SymphonyNetwork) FromConfiguratorNetwork(n configurator.Network) interface{} {
	m.ID = models.NetworkID(n.ID)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*models2.NetworkFeatures)
	}
	return m
}

func (m *MutableSymphonyAgent) GetMagmadGateway() *models2.MagmadGateway {
	return &models2.MagmadGateway{
		Description: m.Description,
		Device:      m.Device,
		ID:          models.GatewayID(m.ID),
		Magmad:      m.Magmad,
		Name:        m.Name,
		Tier:        m.Tier,
	}
}

func (m *MutableSymphonyAgent) GetAdditionalWritesOnCreate() []configurator.EntityWriteOperation {
	ent := configurator.NetworkEntity{
		Type: devmand.SymphonyAgentType,
		Key:  string(m.ID),
	}
	for _, managedDevice := range m.ManagedDevices {
		ent.Associations = append(ent.Associations, storage.TypeAndKey{Type: devmand.SymphonyDeviceType, Key: managedDevice})
	}
	writes := []configurator.EntityWriteOperation{
		ent,
		configurator.EntityUpdateCriteria{
			Type:              orc8r.MagmadGatewayType,
			Key:               string(m.ID),
			AssociationsToAdd: []storage.TypeAndKey{{Type: devmand.SymphonyAgentType, Key: string(m.ID)}},
		},
	}

	return writes
}

func (m *MutableSymphonyAgent) GetAdditionalEntitiesToLoadOnUpdate(agentID string) []storage.TypeAndKey {
	return []storage.TypeAndKey{{Type: devmand.SymphonyAgentType, Key: agentID}}
}

func (m *MutableSymphonyAgent) GetAdditionalWritesOnUpdate(
	agentID string,
	loadedEntities map[storage.TypeAndKey]configurator.NetworkEntity,
) ([]configurator.EntityWriteOperation, error) {
	ret := []configurator.EntityWriteOperation{}
	_, ok := loadedEntities[storage.TypeAndKey{Type: devmand.SymphonyAgentType, Key: agentID}]
	if !ok {
		return nil, merrors.ErrNotFound
	}

	update := configurator.EntityUpdateCriteria{
		Type: devmand.SymphonyAgentType, Key: agentID,
		AssociationsToSet: []storage.TypeAndKey{},
	}
	for _, dID := range m.ManagedDevices {
		update.AssociationsToSet = append(update.AssociationsToSet, storage.TypeAndKey{Type: devmand.SymphonyDeviceType, Key: dID})
	}
	ret = append(ret, update)

	return ret, nil
}

func (m *MutableSymphonyAgent) ToConfiguratorEntity() configurator.NetworkEntity {
	ret := configurator.NetworkEntity{
		Type: devmand.SymphonyAgentType,
		Key:  string(m.ID),
	}
	for _, managedDevice := range m.ManagedDevices {
		ret.Associations = append(ret.Associations, storage.TypeAndKey{Type: devmand.SymphonyDeviceType, Key: managedDevice})
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
	update := configurator.EntityUpdateCriteria{
		Type: devmand.SymphonyAgentType, Key: agentID,
		AssociationsToSet: []storage.TypeAndKey{},
	}
	for _, dID := range *m {
		update.AssociationsToSet = append(update.AssociationsToSet, storage.TypeAndKey{Type: devmand.SymphonyDeviceType, Key: dID})
	}
	return []configurator.EntityUpdateCriteria{update}, nil
}

func (m *SymphonyDevice) FromBackendModels(ent configurator.NetworkEntity) *SymphonyDevice {
	m.Name = SymphonyDeviceName(ent.Name)
	m.ID = SymphonyDeviceID(ent.Key)
	m.Config = ent.Config.(*SymphonyDeviceConfig)
	for _, tk := range ent.ParentAssociations {
		if tk.Type == devmand.SymphonyAgentType {
			m.ManagingAgent = SymphonyDeviceAgent(tk.Key)
		}
	}
	state, err := state.GetState(ent.NetworkID, devmand.SymphonyDeviceStateType, ent.Key)
	if err == nil {
		m.State = state.ReportedState.(*SymphonyDeviceState)
	}
	return m
}

func (m *MutableSymphonyDevice) FromBackendModels(ent configurator.NetworkEntity) *MutableSymphonyDevice {
	m.Name = SymphonyDeviceName(ent.Name)
	m.ID = SymphonyDeviceID(ent.Key)
	m.Config = ent.Config.(*SymphonyDeviceConfig)
	for _, tk := range ent.ParentAssociations {
		if tk.Type == devmand.SymphonyAgentType {
			m.ManagingAgent = SymphonyDeviceAgent(tk.Key)
		}
	}
	return m
}

func (m *MutableSymphonyDevice) ToEntityUpdateCriteria(nID string) ([]configurator.EntityUpdateCriteria, error) {
	updates := []configurator.EntityUpdateCriteria{
		configurator.EntityUpdateCriteria{
			Type:      devmand.SymphonyDeviceType,
			Key:       string(m.ID),
			NewName:   swag.String(string(m.Name)),
			NewConfig: m.Config,
		},
	}

	// Get updates for agents
	oldDevice, err := configurator.LoadEntity(nID, devmand.SymphonyDeviceType, string(m.ID), configurator.EntityLoadCriteria{LoadAssocsToThis: true})
	if err != nil {
		return nil, err
	}
	oldAgent := ""
	for _, tk := range oldDevice.ParentAssociations {
		if tk.Type == devmand.SymphonyAgentType {
			oldAgent = tk.Key
		}
	}
	newAgent := string(m.ManagingAgent)
	agentUpdates := GetAgentUpdates(string(m.ID), oldAgent, newAgent)

	updates = append(updates, agentUpdates...)
	return updates, nil
}

func (m *SymphonyDeviceAgent) FromBackendModels(networkID, deviceID string) error {
	symphonyDeviceEntity, err := configurator.LoadEntity(networkID, devmand.SymphonyDeviceType, deviceID, configurator.EntityLoadCriteria{LoadAssocsToThis: true})
	if err != nil {
		return err
	}
	managingAgent := ""
	for _, tk := range symphonyDeviceEntity.ParentAssociations {
		if tk.Type == devmand.SymphonyAgentType {
			managingAgent = tk.Key
		}
	}
	*m = SymphonyDeviceAgent(managingAgent)
	return nil
}

func (m *SymphonyDeviceAgent) ToUpdateCriteria(networkID, deviceID string) ([]configurator.EntityUpdateCriteria, error) {
	symphonyDeviceEntity, err := configurator.LoadEntity(networkID, devmand.SymphonyDeviceType, deviceID, configurator.EntityLoadCriteria{LoadAssocsToThis: true})
	if err != nil {
		return nil, err
	}
	oldManagingAgent := ""
	for _, tk := range symphonyDeviceEntity.ParentAssociations {
		if tk.Type == devmand.SymphonyAgentType {
			oldManagingAgent = tk.Key
		}
	}
	return GetAgentUpdates(deviceID, oldManagingAgent, string(*m)), nil
}

func (m *SymphonyDeviceName) FromBackendModels(networkID, deviceID string) error {
	symphonyDeviceEntity, err := configurator.LoadEntity(networkID, devmand.SymphonyDeviceType, deviceID, configurator.EntityLoadCriteria{LoadMetadata: true})
	if err != nil {
		return err
	}
	*m = SymphonyDeviceName(symphonyDeviceEntity.Name)
	return nil
}

func (m *SymphonyDeviceName) ToUpdateCriteria(networkID, deviceID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Type:    devmand.SymphonyDeviceType,
			Key:     deviceID,
			NewName: swag.String(string(*m)),
		},
	}, nil
}

func (m *SymphonyDeviceConfig) FromBackendModels(networkID, deviceID string) error {
	deviceEntityConfig, err := configurator.LoadEntityConfig(networkID, devmand.SymphonyDeviceType, deviceID)
	if err != nil {
		return err
	}
	*m = *deviceEntityConfig.(*SymphonyDeviceConfig)
	return nil
}

func (m *SymphonyDeviceConfig) ToUpdateCriteria(networkID, deviceID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		configurator.EntityUpdateCriteria{
			Type:      devmand.SymphonyDeviceType,
			Key:       deviceID,
			NewConfig: m,
		},
	}, nil
}
