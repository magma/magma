/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streamer

import (
	"sort"

	"magma/lte/cloud/go/lte"
	lteModels "magma/lte/cloud/go/plugin/models"
	lteProtos "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
)

const (
	policyStreamName   = "policydb"
	baseNameStreamName = "base_names"

	mappingsStreamName = "rule_mappings"
)

type PoliciesProvider struct{}

func (provider *PoliciesProvider) GetStreamName() string {
	return policyStreamName
}

func (provider *PoliciesProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	gwEnt, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{})
	if err != nil {
		return nil, err
	}

	ruleEnts, err := configurator.LoadAllEntitiesInNetwork(gwEnt.NetworkID, lte.PolicyRuleEntityType, configurator.EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return nil, err
	}

	ruleProtos := make([]*lteProtos.PolicyRule, 0, len(ruleEnts))
	for _, rule := range ruleEnts {
		ruleProtos = append(ruleProtos, createRuleProtoFromEnt(rule))
	}
	return rulesToUpdates(ruleProtos)
}

func createRuleProtoFromEnt(ruleEnt configurator.NetworkEntity) *lteProtos.PolicyRule {
	if ruleEnt.Config == nil {
		return &lteProtos.PolicyRule{Id: ruleEnt.Key}
	}

	cfg := ruleEnt.Config.(*lteModels.PolicyRuleConfig)
	return &lteProtos.PolicyRule{
		Id:            ruleEnt.Key,
		Priority:      swag.Uint32Value(cfg.Priority),
		RatingGroup:   cfg.RatingGroup,
		MonitoringKey: cfg.MonitoringKey,
		Redirect:      redirectInfoToProto(cfg.Redirect),
		FlowList:      flowListToProto(cfg.FlowList),
		Qos:           qosModelToProto(cfg.Qos),
		TrackingType:  lteProtos.PolicyRule_TrackingType(lteProtos.PolicyRule_TrackingType_value[cfg.TrackingType]),
		HardTimeout:   0,
	}
}

func flowListToProto(flowList []*lteModels.FlowDescription) []*lteProtos.FlowDescription {
	ret := make([]*lteProtos.FlowDescription, 0, len(flowList))
	for _, srcFlow := range flowList {
		protoFlow := &lteProtos.FlowDescription{
			Action: lteProtos.FlowDescription_Action(lteProtos.FlowDescription_Action_value[swag.StringValue(srcFlow.Action)]),
		}
		protos.FillIn(srcFlow, protoFlow)

		protoFlow.Match = &lteProtos.FlowMatch{
			Direction: lteProtos.FlowMatch_Direction(lteProtos.FlowMatch_Direction_value[swag.StringValue(srcFlow.Match.Direction)]),
			IpProto:   lteProtos.FlowMatch_IPProto(lteProtos.FlowMatch_IPProto_value[*srcFlow.Match.IPProto]),
		}
		protos.FillIn(srcFlow.Match, protoFlow.Match)

		ret = append(ret, protoFlow)
	}
	return ret
}

func redirectInfoToProto(redirectModel *lteModels.RedirectInformation) *lteProtos.RedirectInformation {
	if redirectModel == nil {
		return nil
	}

	return &lteProtos.RedirectInformation{
		Support:       lteProtos.RedirectInformation_Support(lteProtos.RedirectInformation_Support_value[swag.StringValue(redirectModel.Support)]),
		AddressType:   lteProtos.RedirectInformation_AddressType(lteProtos.RedirectInformation_AddressType_value[swag.StringValue(redirectModel.AddressType)]),
		ServerAddress: swag.StringValue(redirectModel.ServerAddress),
	}
}

func qosModelToProto(qosModel *lteModels.FlowQos) *lteProtos.FlowQos {
	if qosModel == nil {
		return nil
	}

	return &lteProtos.FlowQos{
		MaxReqBwUl: swag.Uint32Value(qosModel.MaxReqBwUl),
		MaxReqBwDl: swag.Uint32Value(qosModel.MaxReqBwDl),
		// The following values haven't been exposed via the API yet
		GbrUl: 0,
		GbrDl: 0,
		Qci:   0,
		Arp:   nil,
	}
}

func rulesToUpdates(rules []*lteProtos.PolicyRule) ([]*protos.DataUpdate, error) {
	ret := make([]*protos.DataUpdate, 0, len(rules))
	for _, policy := range rules {
		marshaledPolicy, err := proto.Marshal(policy)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &protos.DataUpdate{Key: policy.Id, Value: marshaledPolicy})
	}
	sortUpdates(ret)
	return ret, nil
}

type BaseNamesProvider struct{}

func (provider *BaseNamesProvider) GetStreamName() string {
	return baseNameStreamName
}

func (provider *BaseNamesProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	gwEnt, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{})
	if err != nil {
		return nil, err
	}

	bnEnts, err := configurator.LoadAllEntitiesInNetwork(
		gwEnt.NetworkID,
		lte.BaseNameEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
	)
	if err != nil {
		return nil, err
	}

	bnProtos := make([]*lteProtos.ChargingRuleBaseNameRecord, 0, len(bnEnts))
	for _, bn := range bnEnts {
		baseNameRecord := (&lteModels.BaseNameRecord{}).FromEntity(bn)
		bnProto := &lteProtos.ChargingRuleBaseNameRecord{
			Name:         string(baseNameRecord.Name),
			RuleNamesSet: &lteProtos.ChargingRuleNameSet{RuleNames: baseNameRecord.RuleNames},
		}
		bnProtos = append(bnProtos, bnProto)
	}
	return bnsToUpdates(bnProtos)
}

func bnsToUpdates(bns []*lteProtos.ChargingRuleBaseNameRecord) ([]*protos.DataUpdate, error) {
	ret := make([]*protos.DataUpdate, 0, len(bns))
	for _, bn := range bns {
		// We only send the rule names set here
		marshaledBN, err := proto.Marshal(bn.RuleNamesSet)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &protos.DataUpdate{Key: bn.Name, Value: marshaledBN})
	}
	sortUpdates(ret)
	return ret, nil
}

type RuleMappingsProvider struct {
	// Set DeterministicReturn to true to sort all returned collections (to
	// make testing easier for e.g.)
	DeterministicReturn bool
}

func (r *RuleMappingsProvider) GetStreamName() string {
	return mappingsStreamName
}

func (r *RuleMappingsProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	gwEnt, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{})
	if err != nil {
		return nil, err
	}

	loadCrit := configurator.EntityLoadCriteria{LoadAssocsFromThis: true}
	ruleEnts, err := configurator.LoadAllEntitiesInNetwork(gwEnt.NetworkID, lte.PolicyRuleEntityType, loadCrit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load policy rules")
	}
	bnEnts, err := configurator.LoadAllEntitiesInNetwork(gwEnt.NetworkID, lte.BaseNameEntityType, loadCrit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load base names")
	}

	policiesBySid, err := r.getActivePoliciesBySid(ruleEnts, bnEnts)
	if err != nil {
		return nil, err
	}

	ret := make([]*protos.DataUpdate, 0, len(policiesBySid))
	for sid, policies := range policiesBySid {
		marshaledPolicies, err := proto.Marshal(policies)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal active policies")
		}
		ret = append(ret, &protos.DataUpdate{Key: sid, Value: marshaledPolicies})
	}
	if r.DeterministicReturn {
		sortUpdates(ret)
	}
	return ret, nil
}

func (r *RuleMappingsProvider) getActivePoliciesBySid(policyRules []configurator.NetworkEntity, baseNames []configurator.NetworkEntity) (map[string]*lteProtos.ActivePolicies, error) {
	allEnts := make([]configurator.NetworkEntity, 0, len(policyRules)+len(baseNames))
	allEnts = append(allEnts, policyRules...)
	allEnts = append(allEnts, baseNames...)

	policiesBySid := map[string]*lteProtos.ActivePolicies{}
	for _, ent := range allEnts {
		for _, tk := range ent.Associations {
			switch tk.Type {
			case lte.SubscriberEntityType:
				policies, found := policiesBySid[tk.Key]
				if !found {
					policies = &lteProtos.ActivePolicies{}
					policiesBySid[tk.Key] = policies
				}

				switch ent.Type {
				case lte.PolicyRuleEntityType:
					policies.AssignedPolicies = append(policies.AssignedPolicies, ent.Key)
				case lte.BaseNameEntityType:
					policies.AssignedBaseNames = append(policies.AssignedBaseNames, ent.Key)
				default:
					return nil, errors.Errorf("loaded unexpected entity of type %s", ent.Type)
				}
			}
		}
	}

	if r.DeterministicReturn {
		for _, policies := range policiesBySid {
			sort.Strings(policies.AssignedBaseNames)
			sort.Strings(policies.AssignedPolicies)
		}
	}

	return policiesBySid, nil
}

func sortUpdates(updates []*protos.DataUpdate) {
	sort.Slice(updates, func(i, j int) bool { return updates[i].Key < updates[j].Key })
}
