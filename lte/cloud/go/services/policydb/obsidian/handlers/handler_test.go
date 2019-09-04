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
	"magma/lte/cloud/go/services/policydb/obsidian/models"
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
		ID: "Test",
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: "UPLINK",
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
		RatingGroup:  swag.Uint32(2),
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

	// Test Update Rule Using URL based ID
	testRule.FlowList = []*models.FlowDescription{
		{Match: &models.FlowMatch{IPProto: swag.String("IPPROTO_ICMP"), Direction: "DOWNLINK"}},
	}
	testRule.Priority, testRule.RatingGroup, testRule.TrackingType = swag.Uint32(10), swag.Uint32(3), "ONLY_OCS"
	marshaledTestRule, err = json.Marshal(testRule)
	assert.NoError(t, err)

	updateRuleUrlTestCase := tests.Testcase{
		Name:     "Update Rule",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test", testUrlRoot, networkId),
		Payload:  string(marshaledTestRule),
		Expected: ``,
	}
	_, _, _ = tests.RunTest(t, updateRuleUrlTestCase)

	// Verify update results
	getRuleTestCase2 := tests.Testcase{
		Name:     "Get Rule",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test", testUrlRoot, networkId),
		Payload:  ``,
		Expected: string(marshaledTestRule),
	}
	_, _, _ = tests.RunTest(t, getRuleTestCase2)

	// Get all rules
	getAllRulesTestCase := tests.Testcase{
		Name:     "Get All Rules",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `["Test"]`,
	}
	_, _, _ = tests.RunTest(t, getAllRulesTestCase)

	// Delete a rule
	deleteRuleTestCase := tests.Testcase{
		Name:     "Delete a Rule",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test", testUrlRoot, networkId),
		Payload:  ``,
		Expected: ``,
	}
	_, _, _ = tests.RunTest(t, deleteRuleTestCase)

	// Confirm delete
	getAllRulesTestCase = tests.Testcase{
		Name:     "Confirm Delete Rule",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `[]`,
	}
	_, _, _ = tests.RunTest(t, getAllRulesTestCase)

	// Test Add Rule Missing ID
	addRuleNegTestCase := tests.Testcase{
		Name:                      "Add Rule with a missing id",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:                   `{"flow_list": [{}],"priority": 0}`,
		Expect_http_error_status:  true,
		Skip_payload_verification: true,
	}
	_, _, _ = tests.RunTest(t, addRuleNegTestCase)

	// Update entry with bad IPProto
	updateRuleNegTestCase := tests.Testcase{
		Name:   "Update Rule with bad proto",
		Method: "PUT",
		Url: fmt.Sprintf("%s/%s/policies/rules/Test",
			testUrlRoot,
			networkId),
		Payload:                   `{"id":"Test","flow_list": [{"match": {"ip_proto": "IPPROTO_FOO", "direction": "DOWNLINK"}}],"priority": 10,"rating_group":3,"tracking_type": "ONLY_OCS"}`,
		Expect_http_error_status:  true,
		Skip_payload_verification: true,
	}
	_, _, _ = tests.RunTest(t, updateRuleNegTestCase)

	// Test Multi Match Add Rule
	testRule = &models.PolicyRule{
		ID: "Test_mult",
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("DENY"),
				Match: &models.FlowMatch{
					Direction: "UPLINK",
					IPProto:   swag.String("IPPROTO_TCP"),
					TCPDst:    2,
					TCPSrc:    1,
				},
			},
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: "UPLINK",
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
		RatingGroup:  swag.Uint32(2),
		TrackingType: "ONLY_OCS",
	}
	marshaledTestRule, err = json.Marshal(testRule)
	assert.NoError(t, err)

	addMultRuleTestCase := tests.Testcase{
		Name:     "Add Policy Rule (2 matches)",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  string(marshaledTestRule),
		Expected: `"Test_mult"`,
	}
	_, _, _ = tests.RunTest(t, addMultRuleTestCase)

	// Test Read Rule Using URL based ID
	getMultRuleTestCase1 := tests.Testcase{
		Name:     "Get Rule (2 matches)",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test_mult", testUrlRoot, networkId),
		Payload:  ``,
		Expected: string(marshaledTestRule),
	}
	_, _, _ = tests.RunTest(t, getMultRuleTestCase1)

	testRule = &models.PolicyRule{
		ID:           "Test_qos",
		Priority:     swag.Uint32(5),
		RatingGroup:  swag.Uint32(2),
		TrackingType: "ONLY_OCS",
		Qos: &models.FlowQos{
			MaxReqBwUl: 2000,
			MaxReqBwDl: 1000,
		},
	}
	marshaledTestRule, err = json.Marshal(testRule)
	assert.NoError(t, err)

	addQosRuleTestCase := tests.Testcase{
		Name:     "Add Policy Rule with QoS",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  string(marshaledTestRule),
		Expected: `"Test_qos"`,
	}
	_, _, _ = tests.RunTest(t, addQosRuleTestCase)

	getQosRuleTestCase := tests.Testcase{
		Name:     "Get Policy Rule with QoS",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test_qos", testUrlRoot, networkId),
		Payload:  ``,
		Expected: string(marshaledTestRule),
	}
	_, _, _ = tests.RunTest(t, getQosRuleTestCase)

	testRule = &models.PolicyRule{
		ID:           "Test_redirect",
		Priority:     swag.Uint32(5),
		RatingGroup:  swag.Uint32(2),
		TrackingType: "ONLY_OCS",
		Redirect: &models.RedirectInformation{
			Support:       "ENABLED",
			AddressType:   "URL",
			ServerAddress: "127.0.0.1",
		},
	}
	marshaledTestRule, err = json.Marshal(testRule)
	assert.NoError(t, err)

	addRedirectRuleTestCase := tests.Testcase{
		Name:     "Add Policy Rule with Redirect",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  string(marshaledTestRule),
		Expected: `"Test_redirect"`,
	}
	_, _, _ = tests.RunTest(t, addRedirectRuleTestCase)

	getRedirectRuleTestCase := tests.Testcase{
		Name:     "Get Policy Rule with Redirect",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test_redirect", testUrlRoot, networkId),
		Payload:  ``,
		Expected: string(marshaledTestRule),
	}
	_, _, _ = tests.RunTest(t, getRedirectRuleTestCase)

	// Now run base name test cases using the rules created above

	// Test Listing All Base Names
	listBaseNamesTestCase := tests.Testcase{
		Name:     "List All Base Names",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/base_names", testUrlRoot, networkId),
		Payload:  "",
		Expected: "[]",
	}
	_, _, _ = tests.RunTest(t, listBaseNamesTestCase)

	// Test Add BaseName
	addBaseNameTestCase := tests.Testcase{
		Name:     "Add Base Name",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/policies/base_names", testUrlRoot, networkId),
		Payload:  `{"name": "Test", "rule_names": ["Test_qos", "Test_redirect"]}`,
		Expected: `"Test"`,
	}
	_, _, _ = tests.RunTest(t, addBaseNameTestCase)

	// Test Read BaseName Using URL based name
	getBaseNameTestCase1 := tests.Testcase{
		Name:     "Get Base Name",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/base_names/Test", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `["Test_qos", "Test_redirect"]`,
	}
	_, _, _ = tests.RunTest(t, getBaseNameTestCase1)

	// Test Update BaseName Using URL based name
	updateBaseNameUrlTestCase := tests.Testcase{
		Name:     "Update BaseName",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/policies/base_names/Test", testUrlRoot, networkId),
		Payload:  `["Test_qos"]`,
		Expected: ``,
	}
	_, _, _ = tests.RunTest(t, updateBaseNameUrlTestCase)

	// Verify update BaseName
	getBaseNameTestCase2 := tests.Testcase{
		Name:     "Get BaseName",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/base_names/Test", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `["Test_qos"]`,
	}
	_, _, _ = tests.RunTest(t, getBaseNameTestCase2)

	// Get all BaseNames
	getAllBaseNameTestCase := tests.Testcase{
		Name:     "Get All Base Names",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/base_names", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `["Test"]`,
	}
	_, _, _ = tests.RunTest(t, getAllBaseNameTestCase)

	// Delete a BaseName
	deleteBaseNameTestCase := tests.Testcase{
		Name:     "Delete a BaseName",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s/policies/base_names/Test", testUrlRoot, networkId),
		Payload:  ``,
		Expected: ``,
	}
	_, _, _ = tests.RunTest(t, deleteBaseNameTestCase)

	// Confirm delete
	getAllBaseNameTestCase = tests.Testcase{
		Name:     "Confirm Delete BaseName",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/base_names", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `[]`,
	}
	_, _, _ = tests.RunTest(t, getAllBaseNameTestCase)
}
