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
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
)

func TestHandlers(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configuratorTestInit.StartTestService(t)

	restPort := tests.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	// Test Register Network
	registerNetworkTestCase := tests.Testcase{
		Name:                      "Register Network",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=subscriberdb_obsidian_test_network", testUrlRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
	}
	_, networkId, _ := tests.RunTest(t, registerNetworkTestCase)

	json.Unmarshal([]byte(networkId), &networkId)

	// Test Listing All Subscribers
	listSubscribersTestCase := tests.Testcase{
		Name:     "List All Registered Subscribers",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/subscribers", testUrlRoot, networkId),
		Payload:  "",
		Expected: "[]",
	}
	tests.RunTest(t, listSubscribersTestCase)

	// Test Add Subscriber
	addSubTestCase := tests.Testcase{
		Name:   "Add Subscriber",
		Method: "POST",
		Url:    fmt.Sprintf("%s/%s/subscribers", testUrlRoot, networkId),
		Payload: `{"id":"IMSI12333333333", "lte":{"state":"ACTIVE",
				   "auth_algo":"MILENAGE",
				   "auth_key":"AAAAAAAAAAAAAAAAAAAAAA==",
				   "auth_opc":"AAECAwQFBgcICQoLDA0ODw==", "sub_profile":"superfast"}}`,
		Expected: `"IMSI12333333333"`,
	}
	tests.RunTest(t, addSubTestCase)

	addSubTestCase2 := tests.Testcase{
		Name:   "Add Subscriber",
		Method: "POST",
		Url:    fmt.Sprintf("%s/%s/subscribers", testUrlRoot, networkId),
		Payload: `{"id":"IMSI15234363333", "lte":{"state":"ACTIVE",
				   "auth_algo":"MILENAGE",
				   "auth_key":"DtR1RRaOr+LDnAdYKae2Hw==",
				   "auth_opc":"DtR1RRaOr+LDnAdYKae2Hw==", "sub_profile":"superfast"}}`,
		Expected: `"IMSI15234363333"`,
	}
	tests.RunTest(t, addSubTestCase2)

	deleteSubscriberTestCase := tests.Testcase{
		Name:   "Delete Subscriber",
		Method: "DELETE",
		Url: fmt.Sprintf(
			"%s/%s/subscribers/IMSI15234363333", testUrlRoot, networkId),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, deleteSubscriberTestCase)

	// Test Add Subscriber Using URL based ID
	addSubUrlTestCase := tests.Testcase{
		Name:   "Add Subscriber",
		Method: "POST",
		Url: fmt.Sprintf("%s/%s/subscribers/IMSI12333344444",
			testUrlRoot,
			networkId),
		Payload:  `{"id": "IMSI12333344444", "lte":{"auth_algo":"MILENAGE", "auth_key":"AAAAAAAAAAAAAAAAAAAAAA==", "auth_opc": "AAAAAAAAAAAAAAAAAAAAAA==", "state":"ACTIVE", "sub_profile": "default"}}`,
		Expected: `"IMSI12333344444"`,
	}
	tests.RunTest(t, addSubUrlTestCase)

	// Test Updating Subscriber with omitted opc
	updateSubscriberTestCase := tests.Testcase{
		Name:   "Update Subscriber Data, omit OPC",
		Method: "PUT",
		Url: fmt.Sprintf(
			"%s/%s/subscribers/IMSI12333344444", testUrlRoot, networkId),
		Payload: `{"id": "IMSI12333344444", "lte":{"state":"ACTIVE", "auth_algo":"MILENAGE",
			"auth_key":"AAAAAAAAAAAAAAAAAAAAAA==", "sub_profile": "default"}}`,
		Expected: "",
	}
	tests.RunTest(t, updateSubscriberTestCase)

	// Test Getting Subscriber After Config Update w omitted OPC
	getSubscriberTestCaseOpc1 := tests.Testcase{
		Name:   "Get Updated Subscriber, default OPC",
		Method: "GET",
		Url: fmt.Sprintf(
			"%s/%s/subscribers/IMSI12333344444", testUrlRoot, networkId),
		Payload: "",
		Expected: `{"id":"IMSI12333344444", "lte":{"state":"ACTIVE",
			"auth_algo":"MILENAGE",
			"auth_key":"AAAAAAAAAAAAAAAAAAAAAA==", "sub_profile":"default"}}`,
	}
	tests.RunTest(t, getSubscriberTestCaseOpc1)

	// Test Updating Subscriber with null opc
	updateSubscriberTestCase = tests.Testcase{
		Name:   "Update Subscriber Data, set OPC to null",
		Method: "PUT",
		Url: fmt.Sprintf(
			"%s/%s/subscribers/IMSI12333344444", testUrlRoot, networkId),
		Payload:  `{"id": "IMSI12333344444", "lte":{"state":"ACTIVE", "auth_algo":"MILENAGE", "auth_key":"AAAAAAAAAAAAAAAAAAAAAA==", "auth_opc":"", "sub_profile":"default"}}`,
		Expected: "",
	}
	tests.RunTest(t, updateSubscriberTestCase)

	// Test Getting AG Configs After Config Update w null OPC
	getSubscriberTestCaseOpc2 := tests.Testcase{
		Name:   "Get Updated Subscriber, null OPC",
		Method: "GET",
		Url: fmt.Sprintf(
			"%s/%s/subscribers/IMSI12333344444", testUrlRoot, networkId),
		Payload: "",
		Expected: `{"id":"IMSI12333344444", "lte":{"state":"ACTIVE",
			"auth_algo":"MILENAGE",
			"auth_key":"AAAAAAAAAAAAAAAAAAAAAA==","sub_profile":"default"}}`,
	}
	tests.RunTest(t, getSubscriberTestCaseOpc2)

	// Test Listing All Subscribers
	listSubscribersTestCase = tests.Testcase{
		Name:     "List All Registered Subscribers",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/subscribers", testUrlRoot, networkId),
		Payload:  "",
		Expected: `["IMSI12333333333","IMSI12333344444"]`,
		// OR `["IMSI12333344444", "IMSI12333333333"]`
		Skip_payload_verification: true,
	}
	_, resp, _ := tests.RunTest(t, listSubscribersTestCase)
	if resp != `["IMSI12333344444","IMSI12333333333"]` &&
		resp != `["IMSI12333333333","IMSI12333344444"]` {
		t.Fatalf("Unexpected Response: %s, expected: %s",
			resp, listSubscribersTestCase.Expected)
	}
	// Test Getting Subsriber Data
	getSubscriberTestCase := tests.Testcase{
		Name:   "Get Subscriber Data",
		Method: "GET",
		Url: fmt.Sprintf("%s/%s/subscribers/%s",
			testUrlRoot, networkId, "IMSI12333333333"),
		Payload:  "",
		Expected: `{"id":"IMSI12333333333","lte":{"auth_algo":"MILENAGE","auth_key":"AAAAAAAAAAAAAAAAAAAAAA==","auth_opc":"AAECAwQFBgcICQoLDA0ODw==","state":"ACTIVE","sub_profile":"superfast"}}`,
	}
	tests.RunTest(t, getSubscriberTestCase)
	// Test getting all subscriber data
	getAllSubscribersTestCase := tests.Testcase{
		Name:    "Get all subscriber data",
		Method:  "GET",
		Url:     fmt.Sprintf("%s/%s/subscribers?fields=all", testUrlRoot, networkId),
		Payload: "",
		Expected: `{
				"IMSI12333333333": {
					"id": "IMSI12333333333",
					"lte": {
						"auth_algo":"MILENAGE",
						"auth_key":"AAAAAAAAAAAAAAAAAAAAAA==",
						"auth_opc":"AAECAwQFBgcICQoLDA0ODw==",
						"state":"ACTIVE",
						"sub_profile": "superfast"
					}
				},
				"IMSI12333344444": {
					"id": "IMSI12333344444",
					"lte": {
						"auth_algo":"MILENAGE",
						"auth_key":"AAAAAAAAAAAAAAAAAAAAAA==",
						"state":"ACTIVE",
						"sub_profile": "default"
					}
				}
			}
		`,
	}
	tests.RunTest(t, getAllSubscribersTestCase)

	// Test Setting (Updating) Subscriber
	updateSubscriberTestCase = tests.Testcase{
		Name:   "Update Subscriber Data",
		Method: "PUT",
		Url: fmt.Sprintf(
			"%s/%s/subscribers/IMSI12333333333", testUrlRoot, networkId),
		Payload: `{"id": "IMSI12333333333", "lte":{"state":"ACTIVE", "auth_algo":"MILENAGE",
			"auth_key":"AAAAAAAAAAAAAAAAAAAAAA==",
			"auth_opc":"AAAAAAAAAAAAAAAAAAAAAA==", "sub_profile":"default"}}`,
		Expected: "",
	}
	tests.RunTest(t, updateSubscriberTestCase)

	// Test Getting Subsriber Configs After Config Update
	getSubscriberTestCase2 := tests.Testcase{
		Name:   "Get Updated Subscriber",
		Method: "GET",
		Url: fmt.Sprintf(
			"%s/%s/subscribers/IMSI12333333333", testUrlRoot, networkId),
		Payload: "",
		Expected: `{"id":"IMSI12333333333", "lte":{"state":"ACTIVE",
			"auth_algo":"MILENAGE",
			"auth_key":"AAAAAAAAAAAAAAAAAAAAAA==",
			"auth_opc":"AAAAAAAAAAAAAAAAAAAAAA==","sub_profile": "default"}}`,
	}
	tests.RunTest(t, getSubscriberTestCase2)

	// Delete Subscriber Test
	deleteSubscriberTestCase = tests.Testcase{
		Name:   "Delete Subscriber",
		Method: "DELETE",
		Url: fmt.Sprintf(
			"%s/%s/subscribers/IMSI12333333333", testUrlRoot, networkId),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, deleteSubscriberTestCase)

	// Test Listing All Registered Subscribers after a subscriber's removal
	listSubscribersTestCase = tests.Testcase{
		Name:     "List All Registered Subscribers",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/subscribers", testUrlRoot, networkId),
		Payload:  "",
		Expected: `["IMSI12333344444"]`,
	}
	tests.RunTest(t, listSubscribersTestCase)

	deleteSubscriberTestCase = tests.Testcase{
		Name:   "Delete Subscriber",
		Method: "DELETE",
		Url: fmt.Sprintf(
			"%s/%s/subscribers/IMSI12333344444", testUrlRoot, networkId),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, deleteSubscriberTestCase)
}
