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
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"
)

func TestBaseNames(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "0")
	_ = plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	magmad_test_init.StartTestService(t)
	policydb_test_init.StartTestService(t)
	restPort := tests.StartObsidian(t)

	testUrlRoot := fmt.Sprintf(
		"http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

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

	// Test Listing All Base Names
	listBaseNamesTestCase := tests.Testcase{
		Name:     "List All Base Names",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/base_names", testUrlRoot, networkId),
		Payload:  "",
		Expected: "[]",
	}
	tests.RunTest(t, listBaseNamesTestCase)

	// Test Add BaseName
	addBaseNameTestCase := tests.Testcase{
		Name:     "Add Base Name",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/policies/base_names", testUrlRoot, networkId),
		Payload:  `{"name": "Test", "rule_names": ["rule 1", "rule 2", "rule 3"]}`,
		Expected: `"Test"`,
	}
	tests.RunTest(t, addBaseNameTestCase)

	// Test Read BaseName Using URL based name
	getBaseNameTestCase1 := tests.Testcase{
		Name:     "Get Base Name",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/base_names/Test", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `["rule 1", "rule 2", "rule 3"]`,
	}
	tests.RunTest(t, getBaseNameTestCase1)

	// Test Update BaseName Using URL based name
	updateBaseNameUrlTestCase := tests.Testcase{
		Name:     "Update BaseName",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/policies/base_names/Test", testUrlRoot, networkId),
		Payload:  `["rule 11", "rule 12", "rule 13"]`,
		Expected: ``,
	}
	tests.RunTest(t, updateBaseNameUrlTestCase)

	// Verify update BaseName
	getBaseNameTestCase2 := tests.Testcase{
		Name:     "Get BaseName",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/base_names/Test", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `["rule 11", "rule 12", "rule 13"]`,
	}
	tests.RunTest(t, getBaseNameTestCase2)

	// Get all BaseNames
	getAllBaseNameTestCase := tests.Testcase{
		Name:     "Get All Base Names",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/base_names", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `["Test"]`,
	}
	tests.RunTest(t, getAllBaseNameTestCase)

	// Delete a BaseName
	deleteBaseNameTestCase := tests.Testcase{
		Name:     "Delete a BaseName",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s/policies/base_names/Test", testUrlRoot, networkId),
		Payload:  ``,
		Expected: ``,
	}
	tests.RunTest(t, deleteBaseNameTestCase)

	// Confirm delete
	getAllBaseNameTestCase = tests.Testcase{
		Name:     "Confirm Delete BaseName",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/policies/base_names", testUrlRoot, networkId),
		Payload:  ``,
		Expected: `[]`,
	}
	tests.RunTest(t, getAllBaseNameTestCase)
}
