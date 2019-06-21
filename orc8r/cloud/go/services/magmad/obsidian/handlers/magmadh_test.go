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
	"magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"
)

func TestMagmad(t *testing.T) {
	_ = os.Setenv(handlers.UseNewHandlersEnv, "1")
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
}

// Default gateway config struct. Please DO NOT MODIFY this struct in-place
func newDefaultGatewayConfig() *protos.MagmadGatewayConfig {
	return &protos.MagmadGatewayConfig{
		AutoupgradeEnabled:      true,
		AutoupgradePollInterval: 300,
		CheckinInterval:         60,
		CheckinTimeout:          10,
		Tier:                    "default",
		DynamicServices:         []string{},
	}
}
