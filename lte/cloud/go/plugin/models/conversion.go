/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"encoding/base64"
	"fmt"
	"log"
	"sort"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	orc8rModels "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
	orc8rProtos "magma/orc8r/lib/go/protos"

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
		m.DNS = cfg.(*orc8rModels.NetworkDNSConfig)
	}
	if cfg := n.Configs[orc8r.NetworkFeaturesConfig]; cfg != nil {
		m.Features = cfg.(*orc8rModels.NetworkFeatures)
	}
	if cfg := n.Configs[lte.NetworkSubscriberConfigType]; cfg != nil {
		m.SubscriberConfig = cfg.(*NetworkSubscriberConfig)
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

func (m *NetworkSubscriberConfig) GetFromNetwork(network configurator.Network) interface{} {
	res := orc8rModels.GetNetworkConfig(network, lte.NetworkSubscriberConfigType)
	if res == nil {
		return &NetworkSubscriberConfig{}
	}
	return res
}

func (m *NetworkSubscriberConfig) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, lte.NetworkSubscriberConfigType, m), nil
}

func (m *RuleNames) GetFromNetwork(network configurator.Network) interface{} {
	iNetworkSubscriberConfig := orc8rModels.GetNetworkConfig(network, lte.NetworkSubscriberConfigType)
	if iNetworkSubscriberConfig == nil {
		return RuleNames{}
	}
	return iNetworkSubscriberConfig.(*NetworkSubscriberConfig).NetworkWideRuleNames
}

func (m *RuleNames) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iNetworkSubscriberConfig := orc8rModels.GetNetworkConfig(network, lte.NetworkSubscriberConfigType)
	if iNetworkSubscriberConfig == nil {
		// allow update even not previously defined
		iNetworkSubscriberConfig = &NetworkSubscriberConfig{}
	}
	iNetworkSubscriberConfig.(*NetworkSubscriberConfig).NetworkWideRuleNames = *m
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, lte.NetworkSubscriberConfigType, iNetworkSubscriberConfig), nil
}

func (m *BaseNames) GetFromNetwork(network configurator.Network) interface{} {
	iNetworkSubscriberConfig := orc8rModels.GetNetworkConfig(network, lte.NetworkSubscriberConfigType)
	if iNetworkSubscriberConfig == nil {
		return BaseNames{}
	}
	return iNetworkSubscriberConfig.(*NetworkSubscriberConfig).NetworkWideBaseNames
}

func (m *BaseNames) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	iNetworkSubscriberConfig := orc8rModels.GetNetworkConfig(network, lte.NetworkSubscriberConfigType)
	if iNetworkSubscriberConfig == nil {
		// allow update even not previously defined
		iNetworkSubscriberConfig = &NetworkSubscriberConfig{}
	}
	iNetworkSubscriberConfig.(*NetworkSubscriberConfig).NetworkWideBaseNames = *m
	return orc8rModels.GetNetworkConfigUpdateCriteria(network.ID, lte.NetworkSubscriberConfigType, iNetworkSubscriberConfig), nil
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
		Type:      lte.CellularEnodebType,
		Key:       m.Serial,
		NewName:   swag.String(m.Name),
		NewConfig: m.Config,
	}
}

func (m *Subscriber) FromBackendModels(ent configurator.NetworkEntity) *Subscriber {
	m.ID = SubscriberID(ent.Key)
	m.Lte = ent.Config.(*LteSubscription)
	// If no profile in backend, return "default"
	if m.Lte.SubProfile == "" {
		m.Lte.SubProfile = "default"
	}
	for _, tk := range ent.Associations {
		if tk.Type == lte.ApnEntityType {
			m.ActiveApns = append(m.ActiveApns, tk.Key)
		}
	}
	return m
}

func (m *SubProfile) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *BaseNameRecord) ToEntity() configurator.NetworkEntity {
	return configurator.NetworkEntity{
		Type:         lte.BaseNameEntityType,
		Key:          string(m.Name),
		Associations: m.getAssociations(),
	}
}

func (m *BaseNameRecord) FromEntity(ent configurator.NetworkEntity) *BaseNameRecord {
	m.Name = BaseName(ent.Key)
	for _, tk := range ent.Associations {
		if tk.Type == lte.PolicyRuleEntityType {
			m.RuleNames = append(m.RuleNames, tk.Key)
		} else if tk.Type == lte.SubscriberEntityType {
			m.AssignedSubscribers = append(m.AssignedSubscribers, SubscriberID(tk.Key))
		}
	}
	return m
}

func (m *BaseNameRecord) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	return configurator.EntityUpdateCriteria{
		Type:              lte.BaseNameEntityType,
		Key:               string(m.Name),
		AssociationsToSet: m.getAssociations(),
	}
}

func (m *BaseNameRecord) getAssociations() []storage.TypeAndKey {
	allAssocs := make([]storage.TypeAndKey, 0, len(m.RuleNames)+len(m.AssignedSubscribers))
	allAssocs = append(allAssocs, m.RuleNames.ToAssocs()...)
	for _, sid := range m.AssignedSubscribers {
		allAssocs = append(allAssocs, storage.TypeAndKey{Type: lte.SubscriberEntityType, Key: string(sid)})
	}
	return allAssocs
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
	ret := configurator.NetworkEntity{
		Type:   lte.PolicyRuleEntityType,
		Key:    string(m.ID),
		Config: m.getConfig(),
	}
	for _, sid := range m.AssignedSubscribers {
		ret.Associations = append(ret.Associations, storage.TypeAndKey{Type: lte.SubscriberEntityType, Key: string(sid)})
	}
	return ret
}

func (m *PolicyRule) FromEntity(ent configurator.NetworkEntity) *PolicyRule {
	m.ID = PolicyID(ent.Key)
	m.fillFromConfig(ent.Config)
	for _, assoc := range ent.Associations {
		if assoc.Type == lte.SubscriberEntityType {
			m.AssignedSubscribers = append(m.AssignedSubscribers, SubscriberID(assoc.Key))
		}
	}
	return m
}

func (m *PolicyRule) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type:      lte.PolicyRuleEntityType,
		Key:       string(m.ID),
		NewConfig: m.getConfig(),
	}
	for _, sid := range m.AssignedSubscribers {
		ret.AssociationsToSet = append(ret.AssociationsToSet, storage.TypeAndKey{Type: lte.SubscriberEntityType, Key: string(sid)})
	}
	return ret
}

func (m *PolicyRule) getConfig() *PolicyRuleConfig {
	return &PolicyRuleConfig{
		FlowList:      m.FlowList,
		MonitoringKey: m.MonitoringKey,
		Priority:      m.Priority,
		Qos:           m.Qos,
		RatingGroup:   m.RatingGroup,
		Redirect:      m.Redirect,
		TrackingType:  m.TrackingType,
	}
}

func (m *PolicyRule) fillFromConfig(entConfig interface{}) *PolicyRule {
	if entConfig == nil {
		return m
	}

	cfg := entConfig.(*PolicyRuleConfig)
	m.FlowList = cfg.FlowList
	m.MonitoringKey = cfg.MonitoringKey
	m.Priority = cfg.Priority
	m.Qos = cfg.Qos
	m.RatingGroup = cfg.RatingGroup
	m.Redirect = cfg.Redirect
	m.TrackingType = cfg.TrackingType
	return m
}

func (m *PolicyRuleConfig) ToProto(id string) *protos.PolicyRule {
	var (
		protoMKey = []byte{}
		err       error
	)
	if len(m.MonitoringKey) > 0 {
		if protoMKey, err = base64.StdEncoding.DecodeString(m.MonitoringKey); err != nil {
			log.Printf("Error decoding Monitoring Key '%q' for rule ID '%s', will use as is. Err: %v",
				m.MonitoringKey, id, err)
			protoMKey = []byte(m.MonitoringKey)
		}
	}
	rule := &protos.PolicyRule{
		Id:            id,
		Priority:      swag.Uint32Value(m.Priority),
		RatingGroup:   m.RatingGroup,
		MonitoringKey: protoMKey,
		TrackingType:  protos.PolicyRule_TrackingType(protos.PolicyRule_TrackingType_value[m.TrackingType]),
		HardTimeout:   0,
	}
	if m.Redirect != nil {
		rule.Redirect = m.Redirect.ToProto()
	}
	if m.Qos != nil {
		rule.Qos = m.Qos.ToProto()
	}
	if m.FlowList != nil {
		flowList := make([]*protos.FlowDescription, 0, len(m.FlowList))
		for _, flow := range m.FlowList {
			flowList = append(flowList, flow.ToProto())
		}
		rule.FlowList = flowList
	}
	return rule
}

func (m *RedirectInformation) ToProto() *protos.RedirectInformation {
	return &protos.RedirectInformation{
		Support:       protos.RedirectInformation_Support(protos.RedirectInformation_Support_value[swag.StringValue(m.Support)]),
		AddressType:   protos.RedirectInformation_AddressType(protos.RedirectInformation_AddressType_value[swag.StringValue(m.AddressType)]),
		ServerAddress: swag.StringValue(m.ServerAddress),
	}
}

func (m *FlowQos) ToProto() *protos.FlowQos {
	return &protos.FlowQos{
		MaxReqBwUl: swag.Uint32Value(m.MaxReqBwUl),
		MaxReqBwDl: swag.Uint32Value(m.MaxReqBwDl),
		// The following values haven't been exposed via the API yet
		GbrUl: 0,
		GbrDl: 0,
		Qci:   0,
		Arp:   nil,
	}
}

func (m *FlowDescription) ToProto() *protos.FlowDescription {
	flowDescription := &protos.FlowDescription{
		Action: protos.FlowDescription_Action(protos.FlowDescription_Action_value[swag.StringValue(m.Action)]),
	}
	orc8rProtos.FillIn(m, flowDescription)

	flowDescription.Match = &protos.FlowMatch{
		Direction: protos.FlowMatch_Direction(protos.FlowMatch_Direction_value[swag.StringValue(m.Match.Direction)]),
		IpProto:   protos.FlowMatch_IPProto(protos.FlowMatch_IPProto_value[*m.Match.IPProto]),
	}
	orc8rProtos.FillIn(m.Match, flowDescription.Match)
	return flowDescription
}

func (m *RatingGroup) ToEntity() configurator.NetworkEntity {
	ret := configurator.NetworkEntity{
		Type:   lte.RatingGroupEntityType,
		Key:    fmt.Sprint(uint32(m.ID)),
		Config: m,
	}
	return ret
}

func (m *RatingGroup) FromEntity(ent configurator.NetworkEntity) (*RatingGroup, error) {
	ratingGroupID, err := swag.ConvertUint32(ent.Key)
	if err != nil {
		return nil, err
	}
	m.ID = RatingGroupID(ratingGroupID)
	m = ent.Config.(*RatingGroup)
	return m, nil
}

func (m *MutableRatingGroup) ToEntityUpdateCriteria(id uint32) configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type:      lte.RatingGroupEntityType,
		Key:       fmt.Sprint(id),
		NewConfig: m.ToRatingGroup(id),
	}
	return ret
}

func (m *MutableRatingGroup) ToRatingGroup(id uint32) *RatingGroup {
	ratingGroup := &RatingGroup{}
	ratingGroup.ID = RatingGroupID(id)
	ratingGroup.LimitType = m.LimitType
	return ratingGroup
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
