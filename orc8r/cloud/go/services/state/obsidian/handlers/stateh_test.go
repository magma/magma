/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */
package handlers_test

import (
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/models"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	configuratorTestUtils "magma/orc8r/cloud/go/services/configurator/test_utils"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
)

const testAgHwId = "Test-AGW-Hw-Id"

// TestState is Obsidian Gateway Status Integration Test intended to be run
// on cloud VM
func TestState(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	restPort := tests.StartObsidian(t)

	testNetworkID := "state_obsidian_test_network"
	configuratorTestUtils.RegisterNetwork(t, testNetworkID, "Test Network 1")
	configuratorTestUtils.RegisterGateway(t, testNetworkID, testAgHwId, &models.GatewayDevice{HardwareID: testAgHwId})

	// encode the appropriate certificate into context
	ctx := test_utils.GetContextWithCertificate(t, testAgHwId)

	// put one checkin state into state service
	gwStatus := models.NewDefaultGatewayStatus(testAgHwId)
	test_utils.ReportGatewayStatus(t, ctx, gwStatus)

	getStateNoError(t, restPort, testNetworkID)
	getStateNotFoundError(t, restPort, testNetworkID)
}

func getURL(restPort int, networkID string, hwID string) string {
	url := fmt.Sprintf(
		"http://localhost:%d%s/networks/%s/gateways/%s/gateway_status",
		restPort,
		obsidian.RestRoot,
		networkID,
		hwID,
	)
	return url
}

func getStateNoError(t *testing.T, restPort int, networkID string) {
	url := getURL(restPort, networkID, testAgHwId)
	test_utils.GetGWStatusViaHTTPNoError(t, url, networkID, testAgHwId)
}

func getStateNotFoundError(t *testing.T, restPort int, networkID string) {
	url := getURL(restPort, networkID, "should-not-exist")
	test_utils.GetGWStatusExpectNotFound(t, url)
}
