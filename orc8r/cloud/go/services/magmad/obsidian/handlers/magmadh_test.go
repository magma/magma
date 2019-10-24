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

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	models2 "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/magmad/obsidian/models"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestMagmad(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	restPort := tests.StartObsidian(t)

	testURLRoot := fmt.Sprintf(
		"http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	// Test List Networks
	listCloudsTestCase := tests.Testcase{
		Name:     "List Networks",
		Method:   "GET",
		Url:      testURLRoot,
		Payload:  "",
		Expected: `[]`,
	}
	tests.RunTest(t, listCloudsTestCase)

	// Test Register Network with requestedId
	registerNetworkWithIDTestCase := tests.Testcase{
		Name:                      "Register Network with Requested ID",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=magmad_obsidian_test_network", testURLRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
		Expected:                  `"magmad_obsidian_test_network"`,
	}
	tests.RunTest(t, registerNetworkWithIDTestCase)

	// Test Removal Of Empty Network
	removeNetworkTestCase := tests.Testcase{
		Name:     "Remove Empty Network",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s", testURLRoot, "magmad_obsidian_test_network"),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, removeNetworkTestCase)

	// Test Register Network with invalid requestedId
	registerNetworkWithInvalidIDTestCase := tests.Testcase{
		Name:                      "Register Network with Invalid Requested ID",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=00*my_network", testURLRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, registerNetworkWithInvalidIDTestCase)

	// Register network with uppercase requestedId
	registerNetworkWithInvalidIDTestCase = tests.Testcase{
		Name:                      "Register Network with Invalid Requested ID",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=Magmad_obsidian_test_network", testURLRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, registerNetworkWithInvalidIDTestCase)

	// Test Register Network
	registerNetworkTestCase := tests.Testcase{
		Name:                      "Register Network",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=magmad_obsidian_test_network", testURLRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
	}
	_, networkId, _ := tests.RunTest(t, registerNetworkTestCase)

	json.Unmarshal([]byte(networkId), &networkId)

	// Test Register AG with invalid requestedId
	registerAGWithInvalidIDTestCase := tests.Testcase{
		Name:   "Register AG with Invalid Requested ID",
		Method: "POST",
		Url: fmt.Sprintf(
			"%s/%s/gateways?requested_id=%s", testURLRoot, networkId, "*00_bad_ag"),
		Payload:                   `{"hardware_id":"TestAGHwId12345", "key": {"key_type": "ECHO"}}`,
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, registerAGWithInvalidIDTestCase)

	// Test Register AG with requestedId
	requestedAGId := "my_gateway-1"
	registerAGWithIdTestCase := tests.Testcase{
		Name:   "Register AG with Requested ID",
		Method: "POST",
		Url: fmt.Sprintf(
			"%s/%s/gateways?requested_id=%s", testURLRoot, networkId, requestedAGId),
		Payload:  `{"hardware_id":"TestAGHwId00001", "key": {"key_type": "ECHO"}}`,
		Expected: fmt.Sprintf(`"%s"`, requestedAGId),
	}
	tests.RunTest(t, registerAGWithIdTestCase)

	// Test Register AG
	registerAGTestCase := tests.Testcase{
		Name:     "Register AG",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways", testURLRoot, networkId),
		Payload:  `{"hardware_id":"TestAGHwId00002", "key": {"key_type": "SOFTWARE_ECDSA_SHA256", "key": "MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE+Lckvw/eeV8CemEOWpX30/5XhTHKx/mm6T9MpQWuIM8sOKforNm5UPbZrdOTPEBAtGwJB6Uk9crjCIveFe+sN0zw705L94Giza4ny/6ASBcctCm2JJxFccVsocJIraSC"}}`,
		Expected: `"TestAGHwId00002"`,
	}
	tests.RunTest(t, registerAGTestCase)

	// Test Register without key
	registerAGTestCaseNoKey := tests.Testcase{
		Name:                      "Register AG without Key",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s/%s/gateways", testURLRoot, networkId),
		Payload:                   `{"hardware_id":"TestAGHwId00003", "key": {}}`,
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, registerAGTestCaseNoKey)

	// Test Register without key content
	registerAGTestCaseNoKeyContent := tests.Testcase{
		Name:                      "Register AG with Key but no Key Content",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s/%s/gateways", testURLRoot, networkId),
		Payload:                   `{"hardware_id":"TestAGHwId00003", "key": {"key_type":  "SOFTWARE_ECDSA_SHA256"}}`,
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, registerAGTestCaseNoKeyContent)

	// Test Register with wrong key content
	registerAGTestCaseWrongKeyContent := tests.Testcase{
		Name:                      "Register AG with Key but Wrong Key Content",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s/%s/gateways", testURLRoot, networkId),
		Payload:                   `{"hardware_id":"TestAGHwId00003", "key": {"key_type":  "SOFTWARE_ECDSA_SHA256", "key":"AAAAAAAAAAAAAAAAAAAAAA=="}}`,
		Skip_payload_verification: true,
		Expect_http_error_status:  true,
	}
	tests.RunTest(t, registerAGTestCaseWrongKeyContent)

	// Test Listing All Registered AGs
	listAGsTestCase := tests.Testcase{
		Name:                      "List Registered AGs",
		Method:                    "GET",
		Url:                       fmt.Sprintf("%s/%s/gateways", testURLRoot, networkId),
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

	// Test Forced Removal
	removeNetworkTestCase = tests.Testcase{
		Name:     "Force Remove Non Empty Network",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s?mode=force", testURLRoot, networkId),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, removeNetworkTestCase)

	// Test Register Network
	registerNetworkTestCase = tests.Testcase{
		Name:                      "Register Network 2",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=magmad_obisidian_test_network2", testURLRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
	}
	_, networkId, _ = tests.RunTest(t, registerNetworkTestCase)

	json.Unmarshal([]byte(networkId), &networkId)

	// Test Register AG
	registerAGTestCase = tests.Testcase{
		Name:     "Register AG 2",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways", testURLRoot, networkId),
		Payload:  `{"hardware_id":"TestAGHwId12345", "key": {"key_type": "ECHO"}}`,
		Expected: `"TestAGHwId12345"`,
	}
	tests.RunTest(t, registerAGTestCase)

	// Test Listing All Registered AGs
	listAGsTestCase = tests.Testcase{
		Name:     "List Registered AGs 2",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways", testURLRoot, networkId),
		Payload:  "",
		Expected: `["TestAGHwId12345"]`,
	}
	tests.RunTest(t, listAGsTestCase)

	expCfg := NewDefaultGatewayConfig()
	marshaledCfg, err := expCfg.MarshalBinary()
	assert.NoError(t, err)
	expectedCfgStr := string(marshaledCfg)
	_, err = configurator.CreateEntity(networkId, configurator.NetworkEntity{
		Type: orc8r.UpgradeTierEntityType,
		Key:  "default",
	})
	assert.NoError(t, err)

	// Test Getting AG Configs
	createAGConfigTestCase := tests.Testcase{
		Name:     "Create AG Configs",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways/TestAGHwId12345/configs", testURLRoot, networkId),
		Payload:  expectedCfgStr,
		Expected: `"TestAGHwId12345"`,
	}
	tests.RunTest(t, createAGConfigTestCase)

	getAGConfigTestCase := tests.Testcase{
		Name:   "Get AG Configs",
		Method: "GET",
		Url: fmt.Sprintf("%s/%s/gateways/TestAGHwId12345/configs",
			testURLRoot, networkId),
		Payload:  "",
		Expected: expectedCfgStr,
	}
	tests.RunTest(t, getAGConfigTestCase)

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
		Url:      fmt.Sprintf("%s/%s", testURLRoot, networkId),
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
		Url:      fmt.Sprintf("%s/%s", testURLRoot, networkId),
		Payload:  expectedCfgStr,
		Expected: "",
	}
	tests.RunTest(t, updateNetworkRecordTestCase)

	getNetworkRecordTestCase2 := tests.Testcase{
		Name:     "Get Network Record after Update",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s", testURLRoot, networkId),
		Payload:  "",
		Expected: expectedCfgStr,
	}
	tests.RunTest(t, getNetworkRecordTestCase2)

	// Test AG Unregister
	setAGUnregisterTestCase := tests.Testcase{
		Name:   "Unregister AG",
		Method: "DELETE",
		Url: fmt.Sprintf("%s/%s/gateways/TestAGHwId12345",
			testURLRoot, networkId),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, setAGUnregisterTestCase)

	// Test Listing All Registered AGs after Removal Of AG
	listAGsTestCase2 := tests.Testcase{
		Name:     "List Registered AGs",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways", testURLRoot, networkId),
		Payload:  "",
		Expected: `[]`, // should return an empty array
	}
	tests.RunTest(t, listAGsTestCase2)

	// Test List Networks
	listNetworksTestCase := tests.Testcase{
		Name:     "List Networks",
		Method:   "GET",
		Url:      testURLRoot,
		Payload:  "",
		Expected: fmt.Sprintf(`["%s"]`, networkId),
	}
	tests.RunTest(t, listNetworksTestCase)

	// Test Removal Of Empty Network
	removeNetworkTestCase = tests.Testcase{
		Name:     "Remove Empty Network",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s", testURLRoot, networkId),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, removeNetworkTestCase)

	// Test List Networks
	listNetworksTestCase = tests.Testcase{
		Name:     "List Networks Post Delete",
		Method:   "GET",
		Url:      testURLRoot,
		Payload:  "",
		Expected: "[]",
	}
	tests.RunTest(t, listNetworksTestCase)
}

// Default gateway config struct. Please DO NOT MODIFY this struct in-place
func NewDefaultGatewayConfig() *models2.MagmadGatewayConfigs {
	return &models2.MagmadGatewayConfigs{
		AutoupgradeEnabled:      swag.Bool(true),
		AutoupgradePollInterval: 300,
		CheckinInterval:         60,
		CheckinTimeout:          10,
		DynamicServices:         []string{},
	}
}
