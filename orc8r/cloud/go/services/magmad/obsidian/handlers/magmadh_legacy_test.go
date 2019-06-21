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
	"os"
	"testing"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	config_test_init "magma/orc8r/cloud/go/services/config/test_init"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/magmad/obsidian/models"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/stretchr/testify/assert"
)

func TestMagmadLegacy(t *testing.T) {
	_ = os.Setenv(handlers.UseNewHandlersEnv, "0")
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	magmad_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	config_test_init.StartTestService(t)
	restPort := tests.StartObsidian(t)

	testUrlRoot := fmt.Sprintf(
		"http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)

	// Test List Networks
	listCloudsTestCase := tests.Testcase{
		Name:     "List Networks",
		Method:   "GET",
		Url:      testUrlRoot,
		Payload:  "",
		Expected: `[]`,
	}
	tests.RunTest(t, listCloudsTestCase)

	// Test Register Network with requestedId
	registerNetworkWithIdTestCase := tests.Testcase{
		Name:                      "Register Network with Requested Id",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=magmad_obsidian_test_network", testUrlRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
		Expected:                  `"magmad_obsidian_test_network"`,
	}
	tests.RunTest(t, registerNetworkWithIdTestCase)

	// Test Removal Of Empty Network
	removeNetworkTestCase := tests.Testcase{
		Name:     "Remove Empty Network",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s", testUrlRoot, "magmad_obsidian_test_network"),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, removeNetworkTestCase)

	// Test Register Network with invalid requestedId
	registerNetworkWithInvalidIdTestCase := tests.Testcase{
		Name:                      "Register Network with Invalid Requested Id",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=00*my_network", testUrlRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, registerNetworkWithInvalidIdTestCase)

	// Register network with uppercase requestedId
	registerNetworkWithInvalidIdTestCase = tests.Testcase{
		Name:                      "Register Network with Invalid Requested Id",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=Magmad_obsidian_test_network", testUrlRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, registerNetworkWithInvalidIdTestCase)

	// Test Register Network
	registerNetworkTestCase := tests.Testcase{
		Name:                      "Register Network",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=magmad_obsidian_test_network", testUrlRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
	}
	_, networkId, _ := tests.RunTest(t, registerNetworkTestCase)

	json.Unmarshal([]byte(networkId), &networkId)

	// Test Register AG with invalid requestedId
	registerAGWithInvalidIdTestCase := tests.Testcase{
		Name:   "Register AG with Invalid Requested Id",
		Method: "POST",
		Url: fmt.Sprintf(
			"%s/%s/gateways?requested_id=%s", testUrlRoot, networkId, "*00_bad_ag"),
		Payload:                   `{"hw_id":{"id":"TestAGHwId12345"}, "name": "Test AG Name", "key": {"key_type": "ECHO"}}`,
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, registerAGWithInvalidIdTestCase)

	// Test Register AG with requestedId
	requestedAGId := "my_gateway-1"
	registerAGWithIdTestCase := tests.Testcase{
		Name:   "Register AG with Requested Id",
		Method: "POST",
		Url: fmt.Sprintf(
			"%s/%s/gateways?requested_id=%s", testUrlRoot, networkId, requestedAGId),
		Payload:  `{"hw_id":{"id":"TestAGHwId00001"}, "name": "Test AG Name",  "key": {"key_type": "ECHO"}}`,
		Expected: fmt.Sprintf(`"%s"`, requestedAGId),
	}
	tests.RunTest(t, registerAGWithIdTestCase)

	// Test Register AG
	registerAGTestCase := tests.Testcase{
		Name:     "Register AG",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways", testUrlRoot, networkId),
		Payload:  `{"hw_id":{"id":"TestAGHwId00002"}, "name": "Test AG Name", "key": {"key_type": "SOFTWARE_ECDSA_SHA256", "key": "MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE+Lckvw/eeV8CemEOWpX30/5XhTHKx/mm6T9MpQWuIM8sOKforNm5UPbZrdOTPEBAtGwJB6Uk9crjCIveFe+sN0zw705L94Giza4ny/6ASBcctCm2JJxFccVsocJIraSC"}}`,
		Expected: `"TestAGHwId00002"`,
	}
	tests.RunTest(t, registerAGTestCase)

	// Test Register without key
	registerAGTestCaseNoKey := tests.Testcase{
		Name:                      "Register AG without Key",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s/%s/gateways", testUrlRoot, networkId),
		Payload:                   `{"hw_id":{"id":"TestAGHwId00003"}, "name": "Test AG Name", "key": {}}`,
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, registerAGTestCaseNoKey)

	// Test Register without key content
	registerAGTestCaseNoKeyContent := tests.Testcase{
		Name:                      "Register AG with Key but no Key Content",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s/%s/gateways", testUrlRoot, networkId),
		Payload:                   `{"hw_id":{"id":"TestAGHwId00003"}, "name": "Test AG Name", "key": {"key_type":  "SOFTWARE_ECDSA_SHA256"}}`,
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, registerAGTestCaseNoKeyContent)

	// Test Register with wrong key content
	registerAGTestCaseWrongKeyContent := tests.Testcase{
		Name:                      "Register AG with Key but Wrong Key Content",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s/%s/gateways", testUrlRoot, networkId),
		Payload:                   `{"hw_id":{"id":"TestAGHwId00003"}, "name": "Test AG Name", "key": {"key_type":  "SOFTWARE_ECDSA_SHA256", "key":"AAAAAAAAAAAAAAAAAAAAAA=="}}`,
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, registerAGTestCaseWrongKeyContent)

	// Test Getting AG record
	getAGRecordTestCase := tests.Testcase{
		Name:   "Get AG Record With Specified Name",
		Method: "GET",
		Url: fmt.Sprintf("%s/%s/gateways/%s",
			testUrlRoot, networkId, requestedAGId),
		Payload:  "",
		Expected: `{"hw_id":{"id":"TestAGHwId00001"},"key":{"key_type":"ECHO"},"name":"Test AG Name"}`,
	}
	tests.RunTest(t, getAGRecordTestCase)

	getAGRecordTestCase = tests.Testcase{
		Name:     "Get AG Record With Default Name",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/TestAGHwId00002", testUrlRoot, networkId),
		Payload:  "",
		Expected: `{"hw_id":{"id":"TestAGHwId00002"},"key":{"key":"MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE+Lckvw/eeV8CemEOWpX30/5XhTHKx/mm6T9MpQWuIM8sOKforNm5UPbZrdOTPEBAtGwJB6Uk9crjCIveFe+sN0zw705L94Giza4ny/6ASBcctCm2JJxFccVsocJIraSC","key_type":"SOFTWARE_ECDSA_SHA256"},"name":"Test AG Name"}`,
	}
	tests.RunTest(t, getAGRecordTestCase)

	// Test Updating AG record
	setAGRecordTestCase := tests.Testcase{
		Name:     "Update AG Record Name",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/gateways/TestAGHwId00002", testUrlRoot, networkId),
		Payload:  `{"name": "SoDoSoPaTown Tower", "key": {"key_type": "ECHO"}}`,
		Expected: "",
	}
	tests.RunTest(t, setAGRecordTestCase)

	// Test Getting AG record 2
	getAGRecordTestCase = tests.Testcase{
		Name:     "Get AG Record With Modified Name",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/TestAGHwId00002", testUrlRoot, networkId),
		Payload:  "",
		Expected: `{"hw_id":{"id":"TestAGHwId00002"}, "key": {"key_type": "ECHO"}, "name": "SoDoSoPaTown Tower"}`,
	}
	tests.RunTest(t, getAGRecordTestCase)

	// Test Listing All Registered AGs
	listAGsTestCase := tests.Testcase{
		Name:                      "List Registered AGs",
		Method:                    "GET",
		Url:                       fmt.Sprintf("%s/%s/gateways", testUrlRoot, networkId),
		Payload:                   "",
		Expected:                  "",
		Skip_payload_verification: true,
	}
	_, r, _ := tests.RunTest(t, listAGsTestCase)

	exp1 := fmt.Sprintf(`["TestAGHwId00002","%s"]`, requestedAGId)
	exp2 := fmt.Sprintf(`["%s","TestAGHwId00002"]`, requestedAGId)

	if r != exp1 && r != exp2 {
		t.Fatalf("***** %s test failed, received: %s\n***** Expected: %s or %s",
			listAGsTestCase.Name, r, exp1, exp2)
	}

	// Test Removal Of Non Empty Network
	removeNetworkTestCase = tests.Testcase{
		Name:                      "Remove Non Empty Network",
		Method:                    "DELETE",
		Url:                       fmt.Sprintf("%s/%s", testUrlRoot, networkId),
		Payload:                   "",
		Expected:                  "",
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, removeNetworkTestCase)

	// Test Forced Removal
	removeNetworkTestCase = tests.Testcase{
		Name:     "Force Remove Non Empty Network",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s?mode=force", testUrlRoot, networkId),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, removeNetworkTestCase)

	// Test Register Network
	registerNetworkTestCase = tests.Testcase{
		Name:                      "Register Network 2",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=magmad_obisidian_test_network2", testUrlRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
	}
	_, networkId, _ = tests.RunTest(t, registerNetworkTestCase)

	json.Unmarshal([]byte(networkId), &networkId)

	// Test Register AG
	registerAGTestCase = tests.Testcase{
		Name:     "Register AG 2",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways", testUrlRoot, networkId),
		Payload:  `{"hw_id":{"id":"TestAGHwId12345"}, "key": {"key_type": "ECHO"}}`,
		Expected: `"TestAGHwId12345"`,
	}
	tests.RunTest(t, registerAGTestCase)

	// Test Listing All Registered AGs
	listAGsTestCase = tests.Testcase{
		Name:     "List Registered AGs 2",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways", testUrlRoot, networkId),
		Payload:  "",
		Expected: `["TestAGHwId12345"]`,
	}
	tests.RunTest(t, listAGsTestCase)

	expCfg := &models.MagmadGatewayConfig{}
	err := expCfg.FromServiceModel(newDefaultGatewayConfig())
	assert.NoError(t, err)
	marshaledCfg, err := expCfg.MarshalBinary()
	assert.NoError(t, err)
	expectedCfgStr := string(marshaledCfg)

	// Test Getting AG Configs
	createAGConfigTestCase := tests.Testcase{
		Name:     "Create AG Configs",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways/TestAGHwId12345/configs", testUrlRoot, networkId),
		Payload:  expectedCfgStr,
		Expected: `"TestAGHwId12345"`,
	}
	tests.RunTest(t, createAGConfigTestCase)

	getAGConfigTestCase := tests.Testcase{
		Name:   "Get AG Configs",
		Method: "GET",
		Url: fmt.Sprintf("%s/%s/gateways/TestAGHwId12345/configs",
			testUrlRoot, networkId),
		Payload:  "",
		Expected: expectedCfgStr,
	}
	tests.RunTest(t, getAGConfigTestCase)

	expCfg.Tier = "changed"
	marshaledCfg, err = expCfg.MarshalBinary()
	assert.NoError(t, err)
	expectedCfgStr = string(marshaledCfg)

	// Test Setting (Updating) AG Configs
	setAGConfigTestCase := tests.Testcase{
		Name:   "Set AG Configs",
		Method: "PUT",
		Url: fmt.Sprintf("%s/%s/gateways/TestAGHwId12345/configs",
			testUrlRoot, networkId),
		Payload:  expectedCfgStr,
		Expected: "",
	}
	tests.RunTest(t, setAGConfigTestCase)

	// Test Getting AG Configs After Config Update
	getAGConfigTestCase2 := tests.Testcase{
		Name:   "Get AG Configs 2",
		Method: "GET",
		Url: fmt.Sprintf("%s/%s/gateways/TestAGHwId12345/configs",
			testUrlRoot, networkId),
		Payload:  "",
		Expected: expectedCfgStr,
	}
	tests.RunTest(t, getAGConfigTestCase2)

	// Update network wide property
	//
	// Get Current Network Record
	networkCfg := &models.NetworkRecord{Name: "This Is A Test Network Name"}
	marshaledCfg, err = networkCfg.MarshalBinary()
	assert.NoError(t, err)
	expectedCfgStr = string(marshaledCfg)

	getNetworkRecordTestCase := tests.Testcase{
		Name:     "Get Network Record",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s", testUrlRoot, networkId),
		Payload:  "",
		Expected: expectedCfgStr,
	}
	tests.RunTest(t, getNetworkRecordTestCase)

	networkCfg.Name = "Updated Network Name"
	marshaledCfg, err = networkCfg.MarshalBinary()
	assert.NoError(t, err)
	expectedCfgStr = string(marshaledCfg)

	updateNetworkRecordTestCase := tests.Testcase{
		Name:     "Update Network Record",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s", testUrlRoot, networkId),
		Payload:  expectedCfgStr,
		Expected: "",
	}
	tests.RunTest(t, updateNetworkRecordTestCase)

	getNetworkRecordTestCase2 := tests.Testcase{
		Name:     "Get Network Record after Update",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s", testUrlRoot, networkId),
		Payload:  "",
		Expected: expectedCfgStr,
	}
	tests.RunTest(t, getNetworkRecordTestCase2)

	// Test AG Unregister
	setAGUnregisterTestCase := tests.Testcase{
		Name:   "Unregister AG",
		Method: "DELETE",
		Url: fmt.Sprintf("%s/%s/gateways/TestAGHwId12345",
			testUrlRoot, networkId),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, setAGUnregisterTestCase)

	// Test Listing All Registered AGs after Removal Of AG
	listAGsTestCase2 := tests.Testcase{
		Name:     "List Registered AGs",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways", testUrlRoot, networkId),
		Payload:  "",
		Expected: `[]`, // should return an empty array
	}
	tests.RunTest(t, listAGsTestCase2)

	// Test List Networks
	listNetworksTestCase := tests.Testcase{
		Name:     "List Networks",
		Method:   "GET",
		Url:      testUrlRoot,
		Payload:  "",
		Expected: fmt.Sprintf(`["%s"]`, networkId),
	}
	tests.RunTest(t, listNetworksTestCase)

	// Test Removal Of Empty Network
	removeNetworkTestCase = tests.Testcase{
		Name:     "Remove Empty Network",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s", testUrlRoot, networkId),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, removeNetworkTestCase)

	// Test List Networks
	listNetworksTestCase = tests.Testcase{
		Name:     "List Networks Post Delete",
		Method:   "GET",
		Url:      testUrlRoot,
		Payload:  "",
		Expected: "[]",
	}
	tests.RunTest(t, listNetworksTestCase)
}
