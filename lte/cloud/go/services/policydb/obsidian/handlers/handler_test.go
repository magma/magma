/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers_test

import (
	"encoding/json"
	"fmt"
	"testing"

	lteplugin "magma/lte/cloud/go/plugin"
	"magma/lte/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

// Integration test for the migrated configurator-based handlers
func TestPolicyDBHandlers(t *testing.T) {
	err := plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	assert.NoError(t, err)
	err = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	assert.NoError(t, err)
	configuratorTestInit.StartTestService(t)
	restPort := tests.StartObsidian(t)

	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	// Register Network
	registerNetworkTestCase := tests.Testcase{
		Name:                      "Register Network",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=policydb_obsidian_test_network", testUrlRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
	}
	_, networkId, _ := tests.RunTest(t, registerNetworkTestCase)
	err = json.Unmarshal([]byte(networkId), &networkId)
	assert.NoError(t, err)

	// First run test cases on policy rules (we can't create base names without
	// rules to link to)

	// Test Listing All Policy Rules
	listRulesTestCase := tests.Testcase{
		Name:     "List All Policy Rules",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  "",
		Expected: "[]",
	}
	_, _, _ = tests.RunTest(t, listRulesTestCase)

	testRule := &models.PolicyRule{
		ID: swag.String("Test"),
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_ICMP"),
					IPV4Dst:   "42.42.42.42",
					IPV4Src:   "192.168.0.1/24",
					TCPDst:    2,
					TCPSrc:    1,
					UDPDst:    4,
					UDPSrc:    3,
				},
			},
		},
		Priority:     swag.Uint32(5),
		RatingGroup:  uint32(2),
		TrackingType: "ONLY_OCS",
	}
	marshaledTestRule, err := json.Marshal(testRule)
	assert.NoError(t, err)

	// Test Add Rule
	addRuleTestCase := tests.Testcase{
		Name:     "Add Policy Rule",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  string(marshaledTestRule),
		Expected: `"Test"`,
	}
	_, _, _ = tests.RunTest(t, addRuleTestCase)

	// Test Read Rule Using URL based ID
	getRuleTestCase1 := tests.Testcase{
		Name:     "Get Rule",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test", testUrlRoot, networkId),
		Payload:  ``,
		Expected: string(marshaledTestRule),
	}
	_, _, _ = tests.RunTest(t, getRuleTestCase1)
}
