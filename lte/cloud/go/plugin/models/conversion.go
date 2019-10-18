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
	"reflect"
	"sort"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	models2 "magma/orc8r/cloud/go/pluginimpl/models"
	orcprotos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/thoas/go-funk"
)

func (m *LteNetwork) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	var res []error
	if err := m.Cellular.ValidateModel(); err != nil {
		res = append(res, err)
	}
	if err := m.DNS.ValidateModel(); err != nil {
		res = append(res, err)
	}
	if err := m.Features.ValidateModel(); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *LteNetwork) GetEmptyNetwork() handlers.NetworkModel {
	return &LteNetwork{}
}

func (m *LteNetwork) ToConfiguratorNetwork() configurator.Network {
	return configurator.Network{
		ID:          string(m.ID),
		Type:        lte.LteNetworkType,
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
		m.DNS = cfg.(*models2.NetworkDNSConfig)
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*models2.NetworkFeatures)
	}
	return m
}

func (m *NetworkCellularConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return models2.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkType, m), nil
}

func (m *NetworkCellularConfigs) GetFromNetwork(network configurator.Network) interface{} {
	return models2.GetNetworkConfig(network, lte.CellularNetworkType)
}

func (m FegNetworkID) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iCellularConfig := models2.GetNetworkConfig(network, lte.CellularNetworkType)
	if iCellularConfig == nil {
		return configurator.NetworkUpdateCriteria{}, fmt.Errorf("No cellular network config found")
	}
	iCellularConfig.(*NetworkCellularConfigs).FegNetworkID = m
	return models2.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkType, iCellularConfig), nil
}

func (m FegNetworkID) GetFromNetwork(network configurator.Network) interface{} {
	iCellularConfig := models2.GetNetworkConfig(network, lte.CellularNetworkType)
	if iCellularConfig == nil {
		return nil
	}
	return iCellularConfig.(*NetworkCellularConfigs).FegNetworkID
}

func (m *NetworkEpcConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iCellularConfig := models2.GetNetworkConfig(network, lte.CellularNetworkType)
	if iCellularConfig == nil {
		return configurator.NetworkUpdateCriteria{}, fmt.Errorf("No cellular network config found")
	}
	iCellularConfig.(*NetworkCellularConfigs).Epc = m
	return models2.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkType, iCellularConfig), nil
}

func (m *NetworkEpcConfigs) GetFromNetwork(network configurator.Network) interface{} {
	iCellularConfig := models2.GetNetworkConfig(network, lte.CellularNetworkType)
	if iCellularConfig == nil {
		return nil
	}
	return iCellularConfig.(*NetworkCellularConfigs).Epc
}

func (m *NetworkRanConfigs) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iCellularConfig := models2.GetNetworkConfig(network, lte.CellularNetworkType)
	if iCellularConfig == nil {
		return configurator.NetworkUpdateCriteria{}, fmt.Errorf("No cellular network config found")
	}
	iCellularConfig.(*NetworkCellularConfigs).Ran = m
	return models2.GetNetworkConfigUpdateCriteria(network.ID, lte.CellularNetworkType, iCellularConfig), nil
}

func (m *NetworkRanConfigs) GetFromNetwork(network configurator.Network) interface{} {
	iCellularConfig := models2.GetNetworkConfig(network, lte.CellularNetworkType)
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
	device *models2.GatewayDevice,
	status *models2.GatewayStatus,
) handlers.GatewayModel {
	// delegate most of the fillin to magmad gateway struct
	mdGW := (&models2.MagmadGateway{}).FromBackendModels(magmadGateway, device, status)
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

func (m *MutableLteGateway) GetMagmadGateway() *models2.MagmadGateway {
	return &models2.MagmadGateway{
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
	m.Serial = ent.Key
	m.Config = ent.Config.(*EnodebConfiguration)
	for _, tk := range ent.ParentAssociations {
		if tk.Type == lte.CellularGatewayType {
			m.AttachedGatewayID = tk.Key
		}
	}
	return m
}

func (m *Enodeb) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type:      lte.CellularEnodebType,
		Key:       m.Serial,
		NewName:   swag.String(m.Name),
		NewConfig: m.Config,
	}
}

func (m *Subscriber) FromBackendModels(ent configurator.NetworkEntity) *Subscriber {
	m.ID = ent.Key
	m.Lte = ent.Config.(*LteSubscription)
	// If no profile in backend, return "default"
	if m.Lte.SubProfile == "" {
		m.Lte.SubProfile = "default"
	}
	return m
}

func (m *SubProfile) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

var formatsRegistry = strfmt.NewFormats()

func (m *BaseNameRecord) ToEntity() configurator.NetworkEntity {
	return configurator.NetworkEntity{
		Type:         lte.BaseNameEntityType,
		Key:          string(m.Name),
		Associations: m.RuleNames.ToAssocs(),
	}
}

func (m *BaseNameRecord) FromEntity(ent configurator.NetworkEntity) *BaseNameRecord {
	m.Name = BaseName(ent.Key)
	for _, tk := range ent.Associations {
		if tk.Type == lte.PolicyRuleEntityType {
			m.RuleNames = append(m.RuleNames, tk.Key)
		}
	}
	return m
}

func (m RuleNames) ToAssocs() []storage.TypeAndKey {
	return funk.Map(
		m,
		func(rn string) storage.TypeAndKey {
			return storage.TypeAndKey{Type: lte.PolicyRuleEntityType, Key: rn}
		},
	).([]storage.TypeAndKey)
}

func (m *PolicyRule) ToEntity() configurator.NetworkEntity {
	return configurator.NetworkEntity{
		Type:   lte.PolicyRuleEntityType,
		Key:    *m.ID,
		Config: m,
	}
}

// PolicyRule's ToProto fills in passed protos.PolicyRule struct from
// receiver's protos.PolicyRule
func (policyRule *PolicyRule) ToProto(pfrm proto.Message) error {
	flowRuleProto, ok := pfrm.(*protos.PolicyRule)
	if !ok {
		return fmt.Errorf(
			"Invalid Destination Type %s, *protos.PolicyRule expected",
			reflect.TypeOf(pfrm))
	}
	if policyRule != nil || flowRuleProto != nil {
		orcprotos.FillIn(policyRule, flowRuleProto)
		trackingVal, ok := protos.PolicyRule_TrackingType_value[policyRule.TrackingType]
		if ok {
			flowRuleProto.TrackingType = protos.PolicyRule_TrackingType(trackingVal)
		}
		if flowRuleProto.FlowList == nil {
			flowRuleProto.FlowList = flowListToProto(policyRule.FlowList)
		}
		if policyRule.Redirect != nil {
			flowRuleProto.Redirect = redirectInfoToProto(policyRule.Redirect)
		}
		if policyRule.Priority != nil {
			flowRuleProto.Priority = *policyRule.Priority
		}
		flowRuleProto.Id = *policyRule.ID
		flowRuleProto.MonitoringKey = policyRule.MonitoringKey
		flowRuleProto.RatingGroup = policyRule.RatingGroup
	}
	return nil
}

func redirectInfoToProto(redirectModel *RedirectInformation) *protos.RedirectInformation {
	redirectProto := &protos.RedirectInformation{}
	orcprotos.FillIn(redirectModel, redirectProto)
	supportVal, ok := protos.RedirectInformation_Support_value[*redirectModel.Support]
	if ok {
		redirectProto.Support = protos.RedirectInformation_Support(supportVal)
	}
	addrTypeVal, ok := protos.RedirectInformation_AddressType_value[*redirectModel.AddressType]
	if ok {
		redirectProto.AddressType = protos.RedirectInformation_AddressType(addrTypeVal)
	}
	return redirectProto
}

// Fill protos.PolicyRule.FlowList From passed protos.PolicyRule.FlowList
func flowListToProto(flowList []*FlowDescription) []*protos.FlowDescription {
	var s []*protos.FlowDescription
	for i, flow := range flowList {
		s = append(s, new(protos.FlowDescription))
		orcprotos.FillIn(flow, s[i])
		match := flow.Match
		orcprotos.FillIn(match, s[i].Match)
		if match.IPProto != nil {
			protoVal, ok := protos.FlowMatch_IPProto_value[*match.IPProto]
			if ok {
				s[i].Match.IpProto = protos.FlowMatch_IPProto(protoVal)
			}
		}
		directionVal, ok := protos.FlowMatch_Direction_value[*match.Direction]
		if ok {
			s[i].Match.Direction = protos.FlowMatch_Direction(directionVal)
		}
		if flow.Action != nil {
			actionVal, ok := protos.FlowDescription_Action_value[*flow.Action]
			if ok {
				s[i].Action = protos.FlowDescription_Action(actionVal)
			}
		}
	}
	return s
}
