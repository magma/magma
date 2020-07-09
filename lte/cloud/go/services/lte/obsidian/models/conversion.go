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
	"sort"

	"magma/lte/cloud/go/lte"
	policydbModels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	orc8rModels "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/thoas/go-funk"
)

func (m *LteNetwork) GetEmptyNetwork() handlers.NetworkModel {
	return &LteNetwork{}
}

func (m *LteNetwork) ToConfiguratorNetwork() configurator.Network {
	return configurator.Network{
		ID:          string(m.ID),
		Type:        lte.NetworkType,
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			lte.CellularNetworkType:     m.Cellular,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func (m *LteNetwork) ToUpdateCriteria() configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			lte.CellularNetworkType:     m.Cellular,
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func (m *LteNetwork) FromConfiguratorNetwork(n configurator.Network) interface{} {
	m.ID = models.NetworkID(n.ID)
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg := n.Configs[lte.CellularNetworkType]; cfg != nil {
		m.Cellular = cfg.(*NetworkCellularConfigs)
	}
	if cfg := n.Configs[orc8r.DnsdNetworkType]; cfg != nil {
		m.DNS = cfg.(*orc8rModels.NetworkDNSConfig)
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*orc8rModels.NetworkFeatures)
	}
	if cfg := n.Configs[lte.NetworkSubscriberConfigType]; cfg != nil {
		m.SubscriberConfig = cfg.(*policydbModels.NetworkSubscriberConfig)
	}
	return m
}

func (m *NetworkCellularConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkType, m), nil
}

func (m *NetworkCellularConfigs) GetFromNetwork(network configurator.Network) interface{} {
	return orc8rModels.GetNetworkConfig(network, lte.CellularNetworkType)
}

func (m FegNetworkID) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkType)
	if iCellularConfig == nil {
		return configurator.NetworkUpdateCriteria{}, fmt.Errorf("No cellular network config found")
	}
	iCellularConfig.(*NetworkCellularConfigs).FegNetworkID = m
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkType, iCellularConfig), nil
}

func (m FegNetworkID) GetFromNetwork(network configurator.Network) interface{} {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkType)
	if iCellularConfig == nil {
		return nil
	}
	return iCellularConfig.(*NetworkCellularConfigs).FegNetworkID
}

func (m *NetworkEpcConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkType)
	if iCellularConfig == nil {
		return configurator.NetworkUpdateCriteria{}, fmt.Errorf("No cellular network config found")
	}
	iCellularConfig.(*NetworkCellularConfigs).Epc = m
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkType, iCellularConfig), nil
}

func (m *NetworkEpcConfigs) GetFromNetwork(network configurator.Network) interface{} {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkType)
	if iCellularConfig == nil {
		return nil
	}
	return iCellularConfig.(*NetworkCellularConfigs).Epc
}

func (m *NetworkRanConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkType)
	if iCellularConfig == nil {
		return configurator.NetworkUpdateCriteria{}, fmt.Errorf("No cellular network config found")
	}
	iCellularConfig.(*NetworkCellularConfigs).Ran = m
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkType, iCellularConfig), nil
}

func (m *NetworkRanConfigs) GetFromNetwork(network configurator.Network) interface{} {
	iCellularConfig := orc8rModels.GetNetworkConfig(network, lte.CellularNetworkType)
	if iCellularConfig == nil {
		return nil
	}
	return iCellularConfig.(*NetworkCellularConfigs).Ran
}

func (m *LteGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *LteGateway) FromBackendModels(
	magmadGateway, cellularGateway configurator.NetworkEntity,
	device *orc8rModels.GatewayDevice,
	status *orc8rModels.GatewayStatus,
) handlers.GatewayModel {
	// delegate most of the fillin to magmad gateway struct
	mdGW := (&orc8rModels.MagmadGateway{}).FromBackendModels(magmadGateway, device, status)
	// TODO: we should change this to a reflection based shallow copy
	m.ID, m.Name, m.Description, m.Magmad, m.Tier, m.Device, m.Status = mdGW.ID, mdGW.Name, mdGW.Description, mdGW.Magmad, mdGW.Tier, mdGW.Device, mdGW.Status

	if cellularGateway.Config != nil {
		m.Cellular = cellularGateway.Config.(*GatewayCellularConfigs)
	}
	for _, tk := range cellularGateway.Associations {
		if tk.Type == lte.CellularEnodebType {
			m.ConnectedEnodebSerials = append(m.ConnectedEnodebSerials, tk.Key)
		}
	}
	sort.Strings(m.ConnectedEnodebSerials)

	return m
}

func (m *MutableLteGateway) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	// Custom validation only for cellular and device
	var res []error
	if err := m.Cellular.ValidateModel(); err != nil {
		res = append(res, err)
	}
	if err := m.Device.ValidateModel(); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MutableLteGateway) GetMagmadGateway() *orc8rModels.MagmadGateway {
	return &orc8rModels.MagmadGateway{
		Description: m.Description,
		Device:      m.Device,
		ID:          m.ID,
		Magmad:      m.Magmad,
		Name:        m.Name,
		Tier:        m.Tier,
	}
}

func (m *MutableLteGateway) GetAdditionalWritesOnCreate() []configurator.EntityWriteOperation {
	ent := configurator.NetworkEntity{
		Type:        lte.CellularGatewayType,
		Key:         string(m.ID),
		Name:        string(m.Name),
		Description: string(m.Description),
		Config:      m.Cellular,
	}
	for _, enbSerial := range m.ConnectedEnodebSerials {
		ent.Associations = append(ent.Associations, storage.TypeAndKey{Type: lte.CellularEnodebType, Key: enbSerial})
	}

	return []configurator.EntityWriteOperation{
		ent,
		configurator.EntityUpdateCriteria{
			Type:              orc8r.MagmadGatewayType,
			Key:               string(m.ID),
			AssociationsToAdd: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: string(m.ID)}},
		},
	}
}

func (m *MutableLteGateway) GetAdditionalEntitiesToLoadOnUpdate(gatewayID string) []storage.TypeAndKey {
	return []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: gatewayID}}
}

func (m *MutableLteGateway) GetAdditionalWritesOnUpdate(
	gatewayID string,
	loadedEntities map[storage.TypeAndKey]configurator.NetworkEntity,
) ([]configurator.EntityWriteOperation, error) {
	ret := []configurator.EntityWriteOperation{}
	existingEnt, ok := loadedEntities[storage.TypeAndKey{Type: lte.CellularGatewayType, Key: gatewayID}]
	if !ok {
		return ret, merrors.ErrNotFound
	}

	entUpdate := configurator.EntityUpdateCriteria{
		Type:      lte.CellularGatewayType,
		Key:       string(m.ID),
		NewConfig: m.Cellular,
	}
	if string(m.Name) != existingEnt.Name {
		entUpdate.NewName = swag.String(string(m.Name))
	}
	if string(m.Description) != existingEnt.Description {
		entUpdate.NewDescription = swag.String(string(m.Description))
	}

	for _, enbSerial := range m.ConnectedEnodebSerials {
		entUpdate.AssociationsToSet = append(entUpdate.AssociationsToSet, storage.TypeAndKey{Type: lte.CellularEnodebType, Key: enbSerial})
	}

	ret = append(ret, entUpdate)
	return ret, nil
}

func (m *GatewayCellularConfigs) FromBackendModels(networkID string, gatewayID string) error {
	cellularConfig, err := configurator.LoadEntityConfig(networkID, lte.CellularGatewayType, gatewayID)
	if err != nil {
		return err
	}
	*m = *cellularConfig.(*GatewayCellularConfigs)
	return nil
}

func (m *GatewayCellularConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Type: lte.CellularGatewayType, Key: gatewayID,
			NewConfig: m,
		},
	}, nil
}

func (m *GatewayEpcConfigs) FromBackendModels(networkID string, gatewayID string) error {
	gatewayConfig := &GatewayCellularConfigs{}
	err := gatewayConfig.FromBackendModels(networkID, gatewayID)
	if err != nil {
		return err
	}
	*m = *gatewayConfig.Epc
	return nil
}

func (m *GatewayEpcConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.Epc = m
	return cellularConfig.ToUpdateCriteria(networkID, gatewayID)
}

func (m *GatewayRanConfigs) FromBackendModels(networkID string, gatewayID string) error {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(networkID, gatewayID)
	if err != nil {
		return err
	}
	*m = *cellularConfig.Ran
	return nil
}

func (m *GatewayRanConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.Ran = m
	return cellularConfig.ToUpdateCriteria(networkID, gatewayID)
}

func (m *GatewayNonEpsConfigs) FromBackendModels(networkID string, gatewayID string) error {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(networkID, gatewayID)
	if err != nil {
		return err
	}
	*m = *cellularConfig.NonEpsService
	return nil
}

func (m *GatewayNonEpsConfigs) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	cellularConfig := &GatewayCellularConfigs{}
	err := cellularConfig.FromBackendModels(networkID, gatewayID)
	if err != nil {
		return nil, err
	}
	cellularConfig.NonEpsService = m
	return cellularConfig.ToUpdateCriteria(networkID, gatewayID)
}

func (m *EnodebSerials) FromBackendModels(networkID string, gatewayID string) error {
	cellularGatewayEntity, err := configurator.LoadEntity(networkID, lte.CellularGatewayType, gatewayID, configurator.EntityLoadCriteria{LoadAssocsFromThis: true})
	if err != nil {
		return err
	}
	enodebSerials := EnodebSerials{}
	for _, assoc := range cellularGatewayEntity.Associations {
		if assoc.Type == lte.CellularEnodebType {
			enodebSerials = append(enodebSerials, assoc.Key)
		}
	}
	*m = enodebSerials
	return nil
}

func (m *EnodebSerials) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	enodebSerials := []storage.TypeAndKey{}
	for _, enodebSerial := range *m {
		enodebSerials = append(enodebSerials, storage.TypeAndKey{Type: lte.CellularEnodebType, Key: enodebSerial})
	}
	return []configurator.EntityUpdateCriteria{
		{
			Type:              lte.CellularGatewayType,
			Key:               gatewayID,
			AssociationsToSet: enodebSerials,
		},
	}, nil
}

func (m *EnodebSerials) ToDeleteUpdateCriteria(networkID, gatewayID, enodebID string) configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type: lte.CellularGatewayType, Key: gatewayID,
		AssociationsToDelete: []storage.TypeAndKey{{Type: lte.CellularEnodebType, Key: enodebID}},
	}
}

func (m *EnodebSerials) ToCreateUpdateCriteria(networkID, gatewayID, enodebID string) configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type: lte.CellularGatewayType, Key: gatewayID,
		AssociationsToAdd: []storage.TypeAndKey{{Type: lte.CellularEnodebType, Key: enodebID}},
	}
}

func (m *Enodeb) FromBackendModels(ent configurator.NetworkEntity) *Enodeb {
	m.Name = ent.Name
	m.Description = ent.Description
	m.Serial = ent.Key
	if ent.Config != nil {
		m.Config = ent.Config.(*EnodebConfiguration)
	}
	for _, tk := range ent.ParentAssociations {
		if tk.Type == lte.CellularGatewayType {
			m.AttachedGatewayID = tk.Key
		}
	}
	return m
}

func (m *Enodeb) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type:           lte.CellularEnodebType,
		Key:            m.Serial,
		NewName:        swag.String(m.Name),
		NewDescription: swag.String(m.Description),
		NewConfig:      m.Config,
	}
}

func (m *Apn) FromBackendModels(ent configurator.NetworkEntity) *Apn {
	m.ApnName = ApnName(ent.Key)
	m.ApnConfiguration = ent.Config.(*ApnConfiguration)
	return m
}

func (m ApnList) ToAssocs() []storage.TypeAndKey {
	return funk.Map(
		m,
		func(rn string) storage.TypeAndKey {
			return storage.TypeAndKey{Type: lte.ApnEntityType, Key: rn}
		},
	).([]storage.TypeAndKey)
}
