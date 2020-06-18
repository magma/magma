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

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"

	"magma/lte/cloud/go/lte"
	lteModels "magma/lte/cloud/go/plugin/models"
	lteProtos "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
)

type PoliciesProvider struct{}

// GetUpdatesImpl implements GetUpdates for the policies stream provider
func (provider *PoliciesProvider) GetUpdatesImpl(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
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
	return cfg.ToProto(ruleEnt.Key)
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
	return lte.BaseNameStreamName
}

// GetUpdatesImpl implements GetUpdates for the base names stream provider
func (provider *BaseNamesProvider) GetUpdatesImpl(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
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

// GetUpdatesImpl implements GetUpdates for the rule mapppings stream provider
func (r *RuleMappingsProvider) GetUpdatesImpl(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
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

	policiesBySid, err := r.getAssignedPoliciesBySid(ruleEnts, bnEnts)
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

func (r *RuleMappingsProvider) getAssignedPoliciesBySid(policyRules []configurator.NetworkEntity, baseNames []configurator.NetworkEntity) (map[string]*lteProtos.AssignedPolicies, error) {
	allEnts := make([]configurator.NetworkEntity, 0, len(policyRules)+len(baseNames))
	allEnts = append(allEnts, policyRules...)
	allEnts = append(allEnts, baseNames...)

	policiesBySid := map[string]*lteProtos.AssignedPolicies{}
	for _, ent := range allEnts {
		for _, tk := range ent.Associations {
			switch tk.Type {
			case lte.SubscriberEntityType:
				policies, found := policiesBySid[tk.Key]
				if !found {
					policies = &lteProtos.AssignedPolicies{}
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

type NetworkWideRulesProvider struct{}

// GetUpdatesImpl implements GetUpdates for the network wide rules stream
// provider
func (r *NetworkWideRulesProvider) GetUpdatesImpl(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
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
	config, ok := iNetworkSubscriberConfig.(*lteModels.NetworkSubscriberConfig)
	if !ok {
		return nil, fmt.Errorf("Failed to convert to NetworkSubscriberConfig")
	}

	assignedPolicies := &lteProtos.AssignedPolicies{AssignedPolicies: config.NetworkWideRuleNames}
	for _, baseName := range config.NetworkWideBaseNames {
		assignedPolicies.AssignedBaseNames = append(assignedPolicies.AssignedBaseNames, string(baseName))
	}

	marshaledPolicies, err := proto.Marshal(assignedPolicies)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal active policies")
	}
	return []*protos.DataUpdate{{Key: "", Value: marshaledPolicies}}, nil
}
