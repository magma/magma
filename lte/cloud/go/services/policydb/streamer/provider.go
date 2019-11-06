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

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

const (
	policyStreamName   = "policydb"
	baseNameStreamName = "base_names"
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
		ruleConfig := rule.Config.(*lteModels.PolicyRule)
		ruleProto := &lteProtos.PolicyRule{}
		err = ruleConfig.ToProto(ruleProto)
		if err != nil {
			return nil, err
		}
		ruleProtos = append(ruleProtos, ruleProto)
	}
	return rulesToUpdates(ruleProtos)
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
	sort.Slice(ret, func(i, j int) bool { return ret[i].Key < ret[j].Key })
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
		bnConfig := bn.Config.(*lteModels.BaseNameRecord)
		ruleNames := make([]string, 0, len(bn.Associations))
		for _, assoc := range bn.Associations {
			ruleNames = append(ruleNames, assoc.Key)
		}
		bnProto := &lteProtos.ChargingRuleBaseNameRecord{
			Name:         string(bnConfig.Name),
			RuleNamesSet: &lteProtos.ChargingRuleNameSet{RuleNames: ruleNames},
		}
		bnProtos = append(bnProtos, bnProto)
	}
	return bnsToUpdates(bnProtos)
}

func bnsToUpdates(bns []*lteProtos.ChargingRuleBaseNameRecord) ([]*protos.DataUpdate, error) {
	ret := make([]*protos.DataUpdate, 0, len(bns))
	for _, bn := range bns {
		marshaledBN, err := proto.Marshal(bn)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &protos.DataUpdate{Key: bn.Name, Value: marshaledBN})
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].Key < ret[j].Key })
	return ret, nil
}
