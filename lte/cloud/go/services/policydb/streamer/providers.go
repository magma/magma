/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package streamer

import (
	"fmt"
	"sort"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
)

// TODO: need to stream down the infinite credit charging keys from here

type RatingGroupsProvider struct{}

func (p *RatingGroupsProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	gwEnt, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{}, serdes.Entity)
	if err != nil {
		return nil, err
	}

	ratingGroupEnts, _, err := configurator.LoadAllEntitiesOfType(
		gwEnt.NetworkID, lte.RatingGroupEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
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
		return nil, fmt.Errorf("failed to convert to RatingGroup proto")
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

func (p *PoliciesProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	gw, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{}, serdes.Entity)
	if err != nil {
		return nil, err
	}

	rules, _, err := configurator.LoadAllEntitiesOfType(
		gw.NetworkID, lte.PolicyRuleEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err != nil {
		return nil, err
	}
	qosProfiles, err := loadQosProfiles(gw.NetworkID)
	if err != nil {
		return nil, err
	}

	ruleProtos := make([]*lte_protos.PolicyRule, 0, len(rules))
	for _, rule := range rules {
		ruleProtos = append(ruleProtos, createRuleProtoFromEnt(rule, qosProfiles[rule.Key]))
	}
	return rulesToUpdates(ruleProtos)
}

// loadQosProfiles returns all policy_qos_profile ents, keyed by the key of
// their parent policy rule ent, once for each parent.
func loadQosProfiles(networkID string) (map[string]configurator.NetworkEntity, error) {
	profiles, _, err := configurator.LoadAllEntitiesOfType(
		networkID, lte.PolicyQoSProfileEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsToThis: true},
		serdes.Entity,
	)
	if err != nil {
		return nil, err
	}

	qosByProfileID := map[string]configurator.NetworkEntity{}
	for _, qos := range profiles {
		for _, tk := range qos.ParentAssociations.Filter(lte.PolicyRuleEntityType) {
			qosByProfileID[tk.Key] = qos
		}
	}

	return qosByProfileID, nil
}

func createRuleProtoFromEnt(rule, qosProfile configurator.NetworkEntity) *lte_protos.PolicyRule {
	if rule.Config == nil {
		return &lte_protos.PolicyRule{Id: rule.Key}
	}

	cfg := rule.Config.(*models.PolicyRuleConfig)

	var qos *lte_protos.FlowQos
	if qosProfile.Config != nil {
		qos = (&models.PolicyQosProfile{}).FromEntity(qosProfile).ToProto()
	}
	return cfg.ToProto(rule.Key, qos)
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

func (p *BaseNamesProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	gwEnt, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{}, serdes.Entity)
	if err != nil {
		return nil, err
	}

	bnEnts, _, err := configurator.LoadAllEntitiesOfType(
		gwEnt.NetworkID, lte.BaseNameEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true, LoadAssocsToThis: true},
		serdes.Entity,
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

type ApnRuleMappingsProvider struct{}

// GetUpdates implements GetUpdates for the rule mappings stream provider
func (p *ApnRuleMappingsProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	gwEnt, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{}, serdes.Entity)
	if err != nil {
		return nil, err
	}

	loadCrit := configurator.EntityLoadCriteria{LoadAssocsFromThis: true}
	subEnts, _, err := configurator.LoadAllEntitiesOfType(gwEnt.NetworkID, lte.SubscriberEntityType, loadCrit, serdes.Entity)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load subscribers")
	}

	ret := make([]*protos.DataUpdate, 0, len(subEnts))

	for _, subEnt := range subEnts {
		subscriberPolicySet, err := getSubscriberPolicySet(gwEnt.NetworkID, subEnt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build subscriber policy sets")
		}
		marshaled, err := proto.Marshal(subscriberPolicySet)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal subscriber policy sets")
		}
		ret = append(ret, &protos.DataUpdate{Key: subEnt.Key, Value: marshaled})
	}
	return ret, nil
}

func getSubscriberPolicySet(networkID string, subscriberEnt configurator.NetworkEntity) (*lte_protos.SubscriberPolicySet, error) {
	apnPolicyProfileTks := []storage.TypeAndKey{}
	globalPolicies := []string{}
	globalBaseNames := []string{}

	// Get all the TKs of ApnPolicyProfile ents
	// And also get the global policies and base names
	for _, tk := range subscriberEnt.Associations {
		switch tk.Type {
		case lte.PolicyRuleEntityType:
			globalPolicies = append(globalPolicies, tk.Key)
		case lte.BaseNameEntityType:
			globalBaseNames = append(globalBaseNames, tk.Key)
		case lte.APNPolicyProfileEntityType:
			apnPolicyProfileTks = append(apnPolicyProfileTks, tk)
		}
	}

	// Load in all the ApnPolicyProfile ents, they only
	// have incoming/outogoing assocs
	apnPolicyProfileEnts, err := loadApnPolicyProfileEnts(networkID, apnPolicyProfileTks)
	if err != nil {
		return nil, err
	}

	// Fill in per-APN policies/base-names
	rulesPerApn := []*lte_protos.ApnPolicySet{}
	for _, ent := range apnPolicyProfileEnts {
		apnPolicySet, err := buildApnPolicySet(ent)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get SubscriberPolicySet")
		}
		rulesPerApn = append(rulesPerApn, apnPolicySet)
	}

	return &lte_protos.SubscriberPolicySet{
		GlobalPolicies:  globalPolicies,
		GlobalBaseNames: globalBaseNames,
		RulesPerApn:     rulesPerApn,
	}, nil
}

func loadApnPolicyProfileEnts(networkID string, tks []storage.TypeAndKey) (configurator.NetworkEntities, error) {
	if len(tks) == 0 {
		return configurator.NetworkEntities{}, nil
	}
	loadCrit := configurator.EntityLoadCriteria{LoadAssocsFromThis: true}
	typeFilter := lte.APNPolicyProfileEntityType
	apnPolicyProfileEnts, _, err := configurator.LoadEntities(networkID, &typeFilter, nil, nil, tks, loadCrit, serdes.Entity)
	if err != nil {
		return nil, err
	}
	return apnPolicyProfileEnts, nil
}

func buildApnPolicySet(apnPolicyProfileEnt configurator.NetworkEntity) (*lte_protos.ApnPolicySet, error) {
	policies := []string{}
	var apn string
	apn, err := models.GetAPN(apnPolicyProfileEnt.Key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build ApnPolicySet")
	}

	for _, tk := range apnPolicyProfileEnt.Associations {
		switch tk.Type {
		case lte.PolicyRuleEntityType:
			policies = append(policies, tk.Key)
		case lte.APNEntityType:
			apn = tk.Key
		}
	}
	return &lte_protos.ApnPolicySet{
		Apn:              apn,
		AssignedPolicies: policies,
	}, nil
}

func sortUpdates(updates []*protos.DataUpdate) {
	sort.Slice(updates, func(i, j int) bool { return updates[i].Key < updates[j].Key })
}

type NetworkWideRulesProvider struct{}

func (p *NetworkWideRulesProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	gwEnt, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{}, serdes.Entity)
	if err != nil {
		return nil, err
	}
	iNetworkSubscriberConfig, err := configurator.LoadNetworkConfig(gwEnt.NetworkID, lte.NetworkSubscriberConfigType, serdes.Network)
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
