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

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	orc8rModels "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/storage"
	orc8rProtos "magma/orc8r/lib/go/protos"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/thoas/go-funk"
)

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
		FlowList:       m.FlowList,
		MonitoringKey:  m.MonitoringKey,
		Priority:       m.Priority,
		Qos:            m.Qos,
		RatingGroup:    m.RatingGroup,
		Redirect:       m.Redirect,
		TrackingType:   m.TrackingType,
		AppName:        m.AppName,
		AppServiceType: m.AppServiceType,
	}
}

func (m *PolicyRule) fillFromConfig(entConfig interface{}) *PolicyRule {
	if entConfig == nil {
		return m
	}
	cfg := entConfig.(*PolicyRuleConfig)
	monKey := cfg.MonitoringKey
	_, err := base64.StdEncoding.DecodeString(monKey)
	if err != nil { // if not base64 - encode it for future use
		monKey = base64.StdEncoding.EncodeToString([]byte(monKey))
	}
	m.FlowList = cfg.FlowList
	m.MonitoringKey = monKey
	m.Priority = cfg.Priority
	m.Qos = cfg.Qos
	m.RatingGroup = cfg.RatingGroup
	m.Redirect = cfg.Redirect
	m.TrackingType = cfg.TrackingType
	m.AppName = cfg.AppName
	m.AppServiceType = cfg.AppServiceType
	return m
}

func (m *PolicyRuleConfig) ToProto(id string) *protos.PolicyRule {
	var (
		protoMKey = []byte{}
		err       error
	)
	if len(m.MonitoringKey) > 0 {
		if protoMKey, err = base64.StdEncoding.DecodeString(m.MonitoringKey); err != nil {
			glog.Warningf("Can't decode Monitoring Key '%q' for rule ID '%s', will use as is. Err: %v",
				m.MonitoringKey, id, err)
			protoMKey = []byte(m.MonitoringKey)
		}
	}
	rule := &protos.PolicyRule{
		Id:             id,
		Priority:       swag.Uint32Value(m.Priority),
		RatingGroup:    m.RatingGroup,
		MonitoringKey:  protoMKey,
		TrackingType:   protos.PolicyRule_TrackingType(protos.PolicyRule_TrackingType_value[m.TrackingType]),
		AppName:        protos.PolicyRule_AppName(protos.PolicyRule_AppName_value[m.AppName]),
		AppServiceType: protos.PolicyRule_AppServiceType(protos.PolicyRule_AppServiceType_value[m.AppServiceType]),
		HardTimeout:    0,
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

func (m *RatingGroup) ToProto() *protos.RatingGroup {
	limit_type := protos.RatingGroup_FINITE
	switch limit := *m.LimitType; limit {
	case "INFINITE_METERED":
		limit_type = protos.RatingGroup_INFINITE_METERED
	case "INFINITE_UNMETERED":
		limit_type = protos.RatingGroup_INFINITE_UNMETERED
	}
	rule := &protos.RatingGroup{
		Id:        uint32(m.ID),
		LimitType: limit_type,
	}
	return rule
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
