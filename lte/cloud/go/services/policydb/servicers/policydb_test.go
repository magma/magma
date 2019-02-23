/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"testing"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/servicers"
	orcprotos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestPolicydb(t *testing.T) {
	ds := test_utils.NewMockDatastore()
	ctx := context.Background()

	networkId := orcprotos.NetworkID{Id: "test"}
	rule := protos.PolicyRule{Id: "12345", Priority: 1000}
	lookup := protos.PolicyRuleLookup{NetworkId: &networkId, RuleId: "12345"}
	ruleData := protos.PolicyRuleData{Rule: &rule, NetworkId: &networkId}

	srv := servicers.NewPolicyDBServer(ds)

	_, err := srv.AddRule(ctx, &ruleData)
	assert.NoError(t, err)
	ruleData.Rule.Id = "67890"
	_, err = srv.AddRule(ctx, &ruleData)
	assert.NoError(t, err)

	allRules, err := srv.ListRules(ctx, &networkId)
	assert.NoError(t, err)

	if len(allRules.Rules) != 2 {
		t.Fatalf("Got %d rules, 2 expected...", len(allRules.Rules))
	}

	_, err = srv.DeleteRule(ctx, &lookup)
	assert.NoError(t, err)
	lookup.RuleId = "67890"
	_, err = srv.DeleteRule(ctx, &lookup)
	assert.NoError(t, err)

	allRules, err = srv.ListRules(ctx, &networkId)
	assert.NoError(t, err)

	if len(allRules.Rules) != 0 {
		t.Fatalf("Got %d rules, 0 expected...", len(allRules.Rules))
	}

	// Add a rule
	_, err = srv.AddRule(ctx, &ruleData)
	assert.NoError(t, err)
	_, err = srv.AddRule(ctx, &ruleData)
	assert.Error(t, err) // duplicate addition

	res, err := srv.GetRule(ctx, &lookup)
	assert.NoError(t, err)
	assert.Equal(t, orcprotos.TestMarshal(ruleData.Rule), orcprotos.TestMarshal(res))

	// Delete the rule
	_, err = srv.DeleteRule(ctx, &lookup)
	assert.NoError(t, err)
	_, err = srv.GetRule(ctx, &lookup)
	assert.Error(t, err) // rule already removed

	ruleSet, err := srv.ListRules(ctx, &networkId)
	assert.NoError(t, err)
	assert.Equal(t, orcprotos.TestMarshal(ruleSet), orcprotos.TestMarshal(&protos.PolicyRuleSet{Rules: []*protos.PolicyRule{}}))

	//
	// Base Names Tests
	//
	// Add Base Name
	bnLookup1 := &protos.ChargingRuleBaseNameLookup{NetworkID: &networkId, Name: "test base name"}
	old, err := srv.AddBaseName(ctx, &protos.ChargingRuleBaseNameRequest{
		Lookup: bnLookup1,
		Record: &protos.ChargingRuleNameSet{RuleNames: []string{"rule1", "rule2", "rule3"}},
	})
	assert.NoError(t, err)
	assert.Len(t, old.RuleNames, 0)

	// Update Base Name
	old, err = srv.AddBaseName(ctx, &protos.ChargingRuleBaseNameRequest{
		Lookup: bnLookup1,
		Record: &protos.ChargingRuleNameSet{RuleNames: []string{"rule11", "rule12", "rule11", "rule13", "rule14"}},
	})
	assert.NoError(t, err)
	assert.Len(t, old.RuleNames, 3)

	// Add Another Base Name
	bnLookup2 := &protos.ChargingRuleBaseNameLookup{NetworkID: &networkId, Name: "test base name 2"}
	old, err = srv.AddBaseName(ctx, &protos.ChargingRuleBaseNameRequest{
		Lookup: bnLookup2,
		Record: &protos.ChargingRuleNameSet{RuleNames: []string{"rule21", "rule22", "rule23"}},
	})
	assert.NoError(t, err)
	assert.Len(t, old.RuleNames, 0)

	// List Base Names
	bnSet, err := srv.ListBaseNames(ctx, &networkId)
	assert.NoError(t, err)
	assert.Len(t, bnSet.RuleNames, 2)

	// Delete Base Name
	_, err = srv.DeleteBaseName(ctx, bnLookup2)
	assert.NoError(t, err)

	// List All Base Names
	bnSet, err = srv.ListBaseNames(ctx, &networkId)
	assert.NoError(t, err)
	assert.Len(t, bnSet.RuleNames, 1)

	// Get Non-existant Base Name
	_, err = srv.GetBaseName(ctx, bnLookup2)
	assert.Error(t, err)

	// Get Base Name
	bnRecord, err := srv.GetBaseName(ctx, bnLookup1)
	assert.NoError(t, err)
	assert.Len(t, bnRecord.RuleNames, 4)

	assert.Equal(t, "rule11", bnRecord.RuleNames[0])
	assert.Equal(t, "rule12", bnRecord.RuleNames[1])
	assert.Equal(t, "rule13", bnRecord.RuleNames[2])
	assert.Equal(t, "rule14", bnRecord.RuleNames[3])
}
