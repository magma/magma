/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streamer

import (
	"fmt"
	"sort"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
)

// TODO: need to stream down the infinite credit charging keys from here
type RatingGroupsProvider struct{}

func (p *RatingGroupsProvider) GetStreamName() string {
	return lte.RatingGroupStreamName
}

// GetUpdates implements GetUpdates for the policies stream provider
func (provider *RatingGroupsProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	gwEnt, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{})
	if err != nil {
		return nil, err
	}

	ratingGroupEnts, err := configurator.LoadAllEntitiesInNetwork(gwEnt.NetworkID, lte.RatingGroupEntityType, configurator.EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return nil, err
	}

	rgProtos := make([]*lte_protos.RatingGroup, 0, len(ratingGroupEnts))
	var rgProto *lte_protos.RatingGroup
	for _, ratingGroup := range ratingGroupEnts {
		rgProto, err = createRatingGroupProtoFromEnt(ratingGroup)
		if err != nil {
			return nil, err
		}
		rgProtos = append(rgProtos, rgProto)
	}
	return ratingGroupsToUpdates(rgProtos)
}

func createRatingGroupProtoFromEnt(ratingGroupEnt configurator.NetworkEntity) (*lte_protos.RatingGroup, error) {
	if ratingGroupEnt.Config == nil {
		return nil, fmt.Errorf("Failed to convert to RatingGroup proto")
	}
	cfg := ratingGroupEnt.Config.(*models.RatingGroup)
	return cfg.ToProto(), nil
}

func ratingGroupsToUpdates(ratingGroups []*lte_protos.RatingGroup) ([]*protos.DataUpdate, error) {
	ret := make([]*protos.DataUpdate, 0, len(ratingGroups))
	for _, rg := range ratingGroups {
		marshaledRatingGroup, err := proto.Marshal(rg)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &protos.DataUpdate{Key: fmt.Sprint(rg.Id), Value: marshaledRatingGroup})
	}
	sortUpdates(ret)
	return ret, nil
}

type PoliciesProvider struct{}

func (p *PoliciesProvider) GetStreamName() string {
	return lte.PolicyStreamName
}

func (p *PoliciesProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	gwEnt, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{})
	if err != nil {
		return nil, err
	}

	ruleEnts, err := configurator.LoadAllEntitiesInNetwork(gwEnt.NetworkID, lte.PolicyRuleEntityType, configurator.EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return nil, err
	}

	ruleProtos := make([]*lte_protos.PolicyRule, 0, len(ruleEnts))
	for _, rule := range ruleEnts {
		ruleProtos = append(ruleProtos, createRuleProtoFromEnt(rule))
	}
	return rulesToUpdates(ruleProtos)
}

func createRuleProtoFromEnt(ruleEnt configurator.NetworkEntity) *lte_protos.PolicyRule {
	if ruleEnt.Config == nil {
		return &lte_protos.PolicyRule{Id: ruleEnt.Key}
	}

	cfg := ruleEnt.Config.(*models.PolicyRuleConfig)
	return cfg.ToProto(ruleEnt.Key)
}

func rulesToUpdates(rules []*lte_protos.PolicyRule) ([]*protos.DataUpdate, error) {
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

func (p *BaseNamesProvider) GetStreamName() string {
	return lte.BaseNameStreamName
}

func (p *BaseNamesProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
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

	bnProtos := make([]*lte_protos.ChargingRuleBaseNameRecord, 0, len(bnEnts))
	for _, bn := range bnEnts {
		baseNameRecord := (&models.BaseNameRecord{}).FromEntity(bn)
		bnProto := &lte_protos.ChargingRuleBaseNameRecord{
			Name:         string(baseNameRecord.Name),
			RuleNamesSet: &lte_protos.ChargingRuleNameSet{RuleNames: baseNameRecord.RuleNames},
		}
		bnProtos = append(bnProtos, bnProto)
	}
	return bnsToUpdates(bnProtos)
}

func bnsToUpdates(bns []*lte_protos.ChargingRuleBaseNameRecord) ([]*protos.DataUpdate, error) {
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

type RuleMappingsProvider struct{}

func (p *RuleMappingsProvider) GetStreamName() string {
	return lte.MappingsStreamName
}

// GetUpdates implements GetUpdates for the rule mappings stream provider
func (p *RuleMappingsProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
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

	policiesBySid, err := p.getAssignedPoliciesBySid(ruleEnts, bnEnts)
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
	sortUpdates(ret)
	return ret, nil
}

func (p *RuleMappingsProvider) getAssignedPoliciesBySid(policyRules []configurator.NetworkEntity, baseNames []configurator.NetworkEntity) (map[string]*lte_protos.AssignedPolicies, error) {
	allEnts := make([]configurator.NetworkEntity, 0, len(policyRules)+len(baseNames))
	allEnts = append(allEnts, policyRules...)
	allEnts = append(allEnts, baseNames...)

	policiesBySid := map[string]*lte_protos.AssignedPolicies{}
	for _, ent := range allEnts {
		for _, tk := range ent.Associations {
			switch tk.Type {
			case lte.SubscriberEntityType:
				policies, found := policiesBySid[tk.Key]
				if !found {
					policies = &lte_protos.AssignedPolicies{}
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

	for _, policies := range policiesBySid {
		sort.Strings(policies.AssignedBaseNames)
		sort.Strings(policies.AssignedPolicies)
	}

	return policiesBySid, nil
}

func sortUpdates(updates []*protos.DataUpdate) {
	sort.Slice(updates, func(i, j int) bool { return updates[i].Key < updates[j].Key })
}

type NetworkWideRulesProvider struct{}

func (p *NetworkWideRulesProvider) GetStreamName() string {
	return lte.NetworkWideRulesStreamName
}

func (p *NetworkWideRulesProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	gwEnt, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{})
	if err != nil {
		return nil, err
	}
	iNetworkSubscriberConfig, err := configurator.LoadNetworkConfig(gwEnt.NetworkID, lte.NetworkSubscriberConfigType)
	if err == merrors.ErrNotFound {
		return []*protos.DataUpdate{}, nil
	}
	if err != nil {
		return nil, err
	}
	config, ok := iNetworkSubscriberConfig.(*models.NetworkSubscriberConfig)
	if !ok {
		return nil, fmt.Errorf("failed to convert to NetworkSubscriberConfig")
	}

	assignedPolicies := &lte_protos.AssignedPolicies{AssignedPolicies: config.NetworkWideRuleNames}
	for _, baseName := range config.NetworkWideBaseNames {
		assignedPolicies.AssignedBaseNames = append(assignedPolicies.AssignedBaseNames, string(baseName))
	}

	marshaledPolicies, err := proto.Marshal(assignedPolicies)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal active policies")
	}
	return []*protos.DataUpdate{{Key: "", Value: marshaledPolicies}}, nil
}
