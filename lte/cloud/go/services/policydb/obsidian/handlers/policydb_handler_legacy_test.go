/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Tests for Policy REST Endpoints
package handlers_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	lteplugin "magma/lte/cloud/go/plugin"
	policydb_test_init "magma/lte/cloud/go/services/policydb/test_init"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"
)

func TestPolicyRules(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "0")
	_ = plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	magmad_test_init.StartTestService(t)
	policydb_test_init.StartTestService(t)
	restPort := tests.StartObsidian(t)

	testUrlRoot := fmt.Sprintf(
		"http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)

	// Test Register Network
	registerNetworkTestCase := tests.Testcase{
		Name:                      "Register Network",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=policydb_obsidian_test_network", testUrlRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
	}
	_, networkId, _ := tests.RunTest(t, registerNetworkTestCase)
	json.Unmarshal([]byte(networkId), &networkId)

	// Test Listing All Policy Rules
	listRulesTestCase := tests.Testcase{
		Name:     "List All Policy Rules",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  "",
		Expected: "[]",
	}
	tests.RunTest(t, listRulesTestCase)

	// Test Add Rule
	addRuleTestCase := tests.Testcase{
		Name:     "Add Policy Rule",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  `{"id":"Test","flow_list":[{"action":"PERMIT", "match": {"ip_proto":"IPPROTO_ICMP","ipv4_dst":"42.42.42.42","ipv4_src":"192.168.0.1/24","tcp_dst":2,"tcp_src":1,"udp_dst":4,"udp_src":3,"direction":"UPLINK"}}],"priority":5,"rating_group":2,"tracking_type":"ONLY_OCS"}`,
		Expected: `"Test"`,
	}
	tests.RunTest(t, addRuleTestCase)

	// Test Read Rule Using URL based ID
	getRuleTestCase1 := tests.Testcase{
		Name:     "Get Rule",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `{"id":"Test","flow_list":[{"action":"PERMIT", "match": {"ip_proto":"IPPROTO_ICMP","ipv4_dst":"42.42.42.42","ipv4_src":"192.168.0.1/24","tcp_dst":2,"tcp_src":1,"udp_dst":4,"udp_src":3,"direction":"UPLINK"}}],"priority":5,"rating_group":2,"tracking_type":"ONLY_OCS","monitoring_key":""}`,
	}
	tests.RunTest(t, getRuleTestCase1)

	// Test Update Rule Using URL based ID
	updateRuleUrlTestCase := tests.Testcase{
		Name:     "Update Rule",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test", testUrlRoot, networkId),
		Payload:  `{"id":"Test", "flow_list": [{"match": {"ip_proto": "IPPROTO_ICMP", "direction": "DOWNLINK"}}],"priority": 10,"rating_group": 3,"rating_group":3,"tracking_type": "ONLY_OCS"}`,
		Expected: ``,
	}
	tests.RunTest(t, updateRuleUrlTestCase)

	// Verify update results
	getRuleTestCase2 := tests.Testcase{
		Name:     "Get Rule",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `{"id":"Test","flow_list": [{"action": "PERMIT", "match": {"ip_proto": "IPPROTO_ICMP", "direction": "DOWNLINK"}}],"priority": 10,"rating_group":3,"tracking_type": "ONLY_OCS","monitoring_key":""}`,
	}
	tests.RunTest(t, getRuleTestCase2)

	// Get all rules
	getAllRulesTestCase := tests.Testcase{
		Name:     "Get All Rules",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `["Test"]`,
	}
	tests.RunTest(t, getAllRulesTestCase)

	// Delete a rule
	deleteRuleTestCase := tests.Testcase{
		Name:     "Delete a Rule",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test", testUrlRoot, networkId),
		Payload:  ``,
		Expected: ``,
	}
	tests.RunTest(t, deleteRuleTestCase)

	// Confirm delete
	getAllRulesTestCase = tests.Testcase{
		Name:     "Confirm Delete Rule",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `[]`,
	}
	tests.RunTest(t, getAllRulesTestCase)

	// Test Add Rule Missing ID
	addRuleNegTestCase := tests.Testcase{
		Name:                      "Add Rule with a missing id",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:                   `{"flow_list": [{}],"priority": 0}`,
		Expect_http_error_status:  true,
		Skip_payload_verification: true,
	}
	tests.RunTest(t, addRuleNegTestCase)

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
	tests.RunTest(t, updateRuleNegTestCase)

	// Test Multi Match Add Rule
	addMultRuleTestCase := tests.Testcase{
		Name:     "Add Policy Rule (2 matches)",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  `{"id":"Test_mult","flow_list":[{"action":"DENY", "match":{"ip_proto":"IPPROTO_TCP","tcp_dst":2,"tcp_src":1,"direction":"UPLINK"}},{"action":"DENY", "match":{"ip_proto":"IPPROTO_ICMP","ipv4_dst":"42.42.42.42","ipv4_src":"192.168.0.1/24","udp_dst":4,"udp_src":3,"direction":"DOWNLINK"}}],"priority":5,"rating_group":3,"tracking_type":"ONLY_OCS"}`,
		Expected: `"Test_mult"`,
	}
	tests.RunTest(t, addMultRuleTestCase)

	// Test Read Rule Using URL based ID
	getMultRuleTestCase1 := tests.Testcase{
		Name:     "Get Rule (2 matches)",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test_mult", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `{"id":"Test_mult","flow_list":[{"action":"DENY", "match":{"ip_proto":"IPPROTO_TCP","tcp_dst":2,"tcp_src":1,"direction":"UPLINK"}},{"action":"DENY", "match":{"ip_proto":"IPPROTO_ICMP","ipv4_dst":"42.42.42.42","ipv4_src":"192.168.0.1/24","udp_dst":4,"udp_src":3,"direction":"DOWNLINK"}}],"priority":5,"rating_group":3,"tracking_type":"ONLY_OCS","monitoring_key":""}`,
	}
	tests.RunTest(t, getMultRuleTestCase1)

	addQosRuleTestCase := tests.Testcase{
		Name:     "Add Policy Rule with QoS",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  `{"id":"Test_qos","flow_list":[],"priority":5,"rating_group":3,"tracking_type":"ONLY_OCS","qos":{"max_req_bw_ul": 2000, "max_req_bw_dl": 1000}}`,
		Expected: `"Test_qos"`,
	}
	tests.RunTest(t, addQosRuleTestCase)

	getQosRuleTestCase := tests.Testcase{
		Name:     "Get Policy Rule with QoS",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test_qos", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `{"id":"Test_qos","flow_list":null,"priority":5,"rating_group":3,"tracking_type":"ONLY_OCS","monitoring_key":"","qos":{"max_req_bw_ul": 2000, "max_req_bw_dl": 1000}}`,
	}
	tests.RunTest(t, getQosRuleTestCase)

	addRedirectRuleTestCase := tests.Testcase{
		Name:     "Add Policy Rule with Redirect",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/policies/rules", testUrlRoot, networkId),
		Payload:  `{"id":"Test_redirect","flow_list":[],"priority":5,"rating_group":3,"tracking_type":"ONLY_OCS","redirect":{"support": "ENABLED", "address_type": "URL", "server_address": "127.0.0.1"}}`,
		Expected: `"Test_redirect"`,
	}
	tests.RunTest(t, addRedirectRuleTestCase)

	getRedirectRuleTestCase := tests.Testcase{
		Name:     "Get Policy Rule with Redirect",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/rules/Test_redirect", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `{"id":"Test_redirect","flow_list":null,"priority":5,"rating_group":3,"tracking_type":"ONLY_OCS","monitoring_key":"","redirect":{"support": "ENABLED", "address_type": "URL", "server_address": "127.0.0.1"}}`,
	}
	tests.RunTest(t, getRedirectRuleTestCase)
}
