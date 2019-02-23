/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package policydb_test

import (
	"testing"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb"
	policydb_test_init "magma/lte/cloud/go/services/policydb/test_init"
	orcprotos "magma/orc8r/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

const (
	testNetworkId = "network"
)

func TestPolicyDBControllerClientMethods(t *testing.T) {
	policydb_test_init.StartTestService(t)

	// Something that doesn't exist will throw an error
	_, err := policydb.GetRule(testNetworkId, "doesn't exist")
	assert.Error(t, err)

	// Add a rule
	initialRule := &protos.PolicyRule{
		Id: "test",
		FlowList: []*protos.FlowDescription{
			{
				Match: &protos.FlowMatch{TcpSrc: 0},
			},
		},
		Priority: 10,
	}
	err = policydb.AddRule(testNetworkId, initialRule)
	assert.NoError(t, err)

	// Get it back
	actualRule, err := policydb.GetRule(testNetworkId, "test")
	assert.NoError(t, err)
	assert.Equal(t, orcprotos.TestMarshal(initialRule), orcprotos.TestMarshal(actualRule))

	actualRuleDefSet, err := policydb.GetAllRules(testNetworkId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(actualRuleDefSet))
	assert.Equal(t, orcprotos.TestMarshal(actualRuleDefSet[0]), orcprotos.TestMarshal(actualRule))

	err = policydb.DeleteRule(testNetworkId, "test")
	assert.NoError(t, err)

	_, err = policydb.GetRule(testNetworkId, "test")
	assert.Error(t, err)

	actualRuleDefSet, err = policydb.GetAllRules(testNetworkId)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(actualRuleDefSet))

	oldList, err := policydb.AddBaseName(testNetworkId, "base_name1", []string{"rule1", "rule2", "rule3"})
	assert.NoError(t, err)
	assert.Nil(t, oldList)
	assert.Len(t, oldList, 0)

	oldList, err = policydb.AddBaseName(testNetworkId, "base_name1", []string{"rule11", "rule21"})
	assert.NoError(t, err)
	assert.NotNil(t, oldList)
	assert.Len(t, oldList, 3)

	oldList, err = policydb.AddBaseName(testNetworkId, "base_name2", []string{"rule11", "rule12", "rule13"})
	assert.NoError(t, err)
	assert.Nil(t, oldList)

	baseNames, err := policydb.ListBaseNames(testNetworkId)
	assert.NoError(t, err)
	assert.Len(t, baseNames, 2)

	err = policydb.DeleteBaseName(testNetworkId, "base_name2")
	assert.NoError(t, err)

	baseNames, err = policydb.ListBaseNames(testNetworkId)
	assert.NoError(t, err)
	assert.Len(t, baseNames, 1)

	ruleNames, err := policydb.GetBaseName(testNetworkId, "base_name2")
	assert.Error(t, err)
	assert.Nil(t, ruleNames)

	ruleNames, err = policydb.GetBaseName(testNetworkId, "base_name1")
	assert.NoError(t, err)
	assert.NotNil(t, ruleNames)
	assert.Len(t, ruleNames, 2)
	assert.Equal(t, []string{"rule11", "rule21"}, ruleNames)
}
