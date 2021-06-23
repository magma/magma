/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	orc8rModels "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/storage"
	orc8rProtos "magma/orc8r/lib/go/protos"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
)

// TODO(8/21/20): provide entity-wise namespacing support from configurator
// Configurator only provides network-level namespacing.
// This is good enough for now, as subscriber IDs are validated to not contain
// underscores.
var magicNamespaceSeparator = "___"

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
	ent := configurator.NetworkEntity{
		Type:         lte.BaseNameEntityType,
		Key:          string(m.Name),
		Associations: m.GetAssocs(),
	}
	return ent
}

func (m *BaseNameRecord) FromEntity(ent configurator.NetworkEntity) *BaseNameRecord {
	m.Name = BaseName(ent.Key)
	for _, tk := range ent.ParentAssociations {
		if tk.Type == lte.SubscriberEntityType {
			m.AssignedSubscribers = append(m.AssignedSubscribers, SubscriberID(tk.Key))
		}
	}
	for _, tk := range ent.Associations {
		if tk.Type == lte.PolicyRuleEntityType {
			m.RuleNames = append(m.RuleNames, tk.Key)
		}
	}
	return m
}

func (m *BaseNameRecord) ToUpdateCriteria() configurator.EntityUpdateCriteria {
	update := configurator.EntityUpdateCriteria{
		Type:              lte.BaseNameEntityType,
		Key:               string(m.Name),
		AssociationsToSet: m.GetAssocs(),
	}
	return update
}

func (m *BaseNameRecord) GetAssocs() storage.TKs {
	return m.RuleNames.ToTKs()
}

func (m *BaseNameRecord) GetParentAssocs() []storage.TypeAndKey {
	var parents storage.TKs
	for _, sid := range m.AssignedSubscribers {
		parents = append(parents, storage.TypeAndKey{Type: lte.SubscriberEntityType, Key: string(sid)})
	}
	return parents
}

func (m RuleNames) ToTKs() storage.TKs {
	return storage.MakeTKs(lte.PolicyRuleEntityType, m)
}

func (m *PolicyRule) ToEntity() configurator.NetworkEntity {
	ent := configurator.NetworkEntity{
		Type:         lte.PolicyRuleEntityType,
		Key:          string(m.ID),
		Config:       m.getConfig(),
		Associations: m.GetAssocs(),
	}
	return ent
}

func (m *PolicyRule) FromEntity(ent configurator.NetworkEntity) *PolicyRule {
	m.ID = PolicyID(ent.Key)
	m.fillFromConfig(ent.Config)

	for _, assoc := range ent.ParentAssociations.Filter(lte.SubscriberEntityType) {
		m.AssignedSubscribers = append(m.AssignedSubscribers, SubscriberID(assoc.Key))
	}
	qosProfile, err := ent.Associations.GetFirst(lte.PolicyQoSProfileEntityType)
	if err == nil {
		m.QosProfile = qosProfile.Key
	}

	return m
}

func (m *PolicyRule) ToEntityUpdateCriteria() configurator.EntityUpdateCriteria {
	update := configurator.EntityUpdateCriteria{
		Type:              lte.PolicyRuleEntityType,
		Key:               string(m.ID),
		NewConfig:         m.getConfig(),
		AssociationsToAdd: m.GetAssocs(),
	}
	return update
}

func (m *PolicyRule) GetParentAssocs() storage.TKs {
	var parents []storage.TypeAndKey
	for _, sid := range m.AssignedSubscribers {
		parents = append(parents, storage.TypeAndKey{Type: lte.SubscriberEntityType, Key: string(sid)})
	}
	return parents
}

func (m *PolicyRule) GetAssocs() storage.TKs {
	var children []storage.TypeAndKey
	if m.QosProfile != "" {
		children = append(children, storage.TypeAndKey{Type: lte.PolicyQoSProfileEntityType, Key: m.QosProfile})
	}
	return children
}

func (m *PolicyRule) getConfig() *PolicyRuleConfig {
	return &PolicyRuleConfig{
		FlowList:                m.FlowList,
		MonitoringKey:           m.MonitoringKey,
		Priority:                m.Priority,
		RatingGroup:             m.RatingGroup,
		ServiceIdentifier:       m.ServiceIdentifier,
		Redirect:                m.Redirect,
		TrackingType:            m.TrackingType,
		AppName:                 m.AppName,
		AppServiceType:          m.AppServiceType,
		HeaderEnrichmentTargets: m.HeaderEnrichmentTargets,
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
	m.RatingGroup = cfg.RatingGroup
	m.ServiceIdentifier = cfg.ServiceIdentifier
	m.Redirect = cfg.Redirect
	m.TrackingType = cfg.TrackingType
	m.AppName = cfg.AppName
	m.AppServiceType = cfg.AppServiceType
	m.HeaderEnrichmentTargets = cfg.HeaderEnrichmentTargets
	return m
}

func (m PolicyIdsByApn) ToTKs(subscriberID string) []storage.TypeAndKey {
	var tks []storage.TypeAndKey
	for apnName := range m {
		tks = append(tks, storage.TypeAndKey{Type: lte.APNPolicyProfileEntityType, Key: makeAPNPolicyKey(subscriberID, apnName)})
	}
	return tks
}

func (m PolicyIdsByApn) ToEntities(subscriberID string) []configurator.NetworkEntity {
	var ents []configurator.NetworkEntity
	for apnName, policyIDs := range m {
		// Each apn_policy_profile has 1 edge to an apn and n edges to a policy_rule
		ent := configurator.NetworkEntity{
			Type:         lte.APNPolicyProfileEntityType,
			Key:          makeAPNPolicyKey(subscriberID, apnName),
			Associations: getAPNPolicyAssocs(apnName, policyIDs),
		}
		ents = append(ents, ent)
	}
	return ents
}

func GetAPN(apnPolicyProfileKey string) (string, error) {
	if !strings.Contains(apnPolicyProfileKey, magicNamespaceSeparator) {
		return "", errors.New("incorrectly formatted APNPolicyProfile key")
	}
	return strings.Split(apnPolicyProfileKey, magicNamespaceSeparator)[1], nil
}

func makeAPNPolicyKey(subscriberID, apnName string) string {
	return subscriberID + magicNamespaceSeparator + apnName
}

func getAPNPolicyAssocs(apnName string, policyIDs PolicyIds) []storage.TypeAndKey {
	var assocs []storage.TypeAndKey
	assocs = append(assocs, storage.TypeAndKey{Type: lte.APNEntityType, Key: apnName})
	for _, policyID := range policyIDs {
		assocs = append(assocs, storage.TypeAndKey{Type: lte.PolicyRuleEntityType, Key: string(policyID)})
	}
	return assocs
}

func (m *PolicyRuleConfig) ToProto(id string, qos *protos.FlowQos) *protos.PolicyRule {
	var (
		protoMKey []byte
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
		Qos:            qos,
	}
	if m.ServiceIdentifier != 0 {
		rule.ServiceIdentifier = &protos.ServiceIdentifier{Value: m.ServiceIdentifier}
	}
	if m.Redirect != nil {
		rule.Redirect = m.Redirect.ToProto()
	}
	if m.FlowList != nil {
		flowList := make([]*protos.FlowDescription, 0, len(m.FlowList))
		for _, flow := range m.FlowList {
			flowList = append(flowList, flow.ToProto())
		}
		rule.FlowList = flowList
	}
	if len(m.HeaderEnrichmentTargets) != 0 {
		rule.He = &protos.HeaderEnrichment{Urls: m.HeaderEnrichmentTargets}
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

	// Backwards compatible for old flow match definition
	if m.Match.IPSrc != nil {
		flowDescription.Match.IpSrc = &protos.IPAddress{
			Version: protos.IPAddress_IPVersion(protos.IPAddress_IPVersion_value[m.Match.IPSrc.Version]),
			Address: []byte(m.Match.IPSrc.Address),
		}
		flowDescription.Match.Ipv4Src = m.Match.IPSrc.Address
	}
	if m.Match.IPDst != nil {
		flowDescription.Match.IpDst = &protos.IPAddress{
			Version: protos.IPAddress_IPVersion(protos.IPAddress_IPVersion_value[m.Match.IPDst.Version]),
			Address: []byte(m.Match.IPDst.Address),
		}
		flowDescription.Match.Ipv4Dst = m.Match.IPDst.Address
	}

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

func (m *RatingGroup) FromEntity(ent configurator.NetworkEntity) *RatingGroup {
	return ent.Config.(*RatingGroup)
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

func (m *PolicyQosProfile) FromBackendModels(networkID string, key string) error {
	config, err := configurator.LoadEntityConfig(networkID, lte.PolicyQoSProfileEntityType, key, EntitySerdes)
	if err != nil {
		return err
	}
	*m = *config.(*PolicyQosProfile)
	return nil
}

func (m *PolicyQosProfile) ToUpdateCriteria(networkID string, key string) ([]configurator.EntityUpdateCriteria, error) {
	if key != m.ID {
		return nil, errors.New("id field is read-only")
	}

	exists, err := configurator.DoesEntityExist(networkID, lte.PolicyQoSProfileEntityType, key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("profile does not exist")
	}

	updates := []configurator.EntityUpdateCriteria{
		{
			Type:      lte.PolicyQoSProfileEntityType,
			Key:       key,
			NewConfig: m,
		},
	}
	return updates, nil
}

func (m *PolicyQosProfile) ToEntity() configurator.NetworkEntity {
	ent := configurator.NetworkEntity{
		Type:   lte.PolicyQoSProfileEntityType,
		Key:    m.ID,
		Config: m,
	}
	return ent
}

func (m *PolicyQosProfile) FromEntity(ent configurator.NetworkEntity) *PolicyQosProfile {
	return ent.Config.(*PolicyQosProfile)
}

func (m *PolicyQosProfile) ToProto() *protos.FlowQos {
	proto := &protos.FlowQos{
		MaxReqBwUl: swag.Uint32Value(m.MaxReqBwUl),
		MaxReqBwDl: swag.Uint32Value(m.MaxReqBwDl),
		Qci:        protos.FlowQos_Qci(m.ClassID),
	}
	if m.Gbr != nil {
		proto.GbrUl = swag.Uint32Value(m.Gbr.Uplink)
		proto.GbrDl = swag.Uint32Value(m.Gbr.Downlink)
	}
	if m.Arp != nil {
		arp := &protos.QosArp{PriorityLevel: swag.Uint32Value(m.Arp.PriorityLevel)}
		if swag.BoolValue(m.Arp.PreemptionCapability) {
			arp.PreCapability = 1
		}
		if swag.BoolValue(m.Arp.PreemptionVulnerability) {
			arp.PreVulnerability = 1
		}
		proto.Arp = arp
	}
	return proto
}

func (m PolicyIds) ToTKs() []storage.TypeAndKey {
	var tks []storage.TypeAndKey
	for _, policyID := range m {
		tks = append(tks, storage.TypeAndKey{Type: lte.PolicyRuleEntityType, Key: string(policyID)})
	}
	return tks
}

func (m BaseNames) ToTKs() []storage.TypeAndKey {
	var tks []storage.TypeAndKey
	for _, baseName := range m {
		tks = append(tks, storage.TypeAndKey{Type: lte.BaseNameEntityType, Key: string(baseName)})
	}
	return tks
}
