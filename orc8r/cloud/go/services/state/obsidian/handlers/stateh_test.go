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
	"magma/orc8r/cloud/go/protos"
	checkindTestUtils "magma/orc8r/cloud/go/services/checkind/test_utils"
	"magma/orc8r/cloud/go/services/magmad"
	magmadProtos "magma/orc8r/cloud/go/services/magmad/protos"
	magmadTestInit "magma/orc8r/cloud/go/services/magmad/test_init"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"

	"github.com/stretchr/testify/assert"
)

const testAgHwId = "Test-AGW-Hw-Id"

// TestState is Obsidian Gateway Status Integration Test intended to be run
// on cloud VM
func TestState(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	magmadTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	restPort := tests.StartObsidian(t)

	// create a test network with a single GW
	testNetworkId, err := magmad.RegisterNetwork(
		&magmadProtos.MagmadNetworkRecord{Name: "Test Network 1"},
		"state_obsidian_test_network")
	assert.NoError(t, err)
	hwId := protos.AccessGatewayID{Id: testAgHwId}
	_, err = magmad.RegisterGateway(
		testNetworkId,
		&magmadProtos.AccessGatewayRecord{HwId: &hwId, Name: "Test GW Name"},
	)
	assert.NoError(t, err)

	// encode the appropriate certificate into context
	ctx := test_utils.GetContextWithCertificate(t, testAgHwId)

	// put one checkin state into state service
	gwStatus := checkindTestUtils.GetGatewayStatusSwaggerFixture(testAgHwId)
	test_utils.ReportGatewayStatus(t, ctx, gwStatus)

	getStateNoError(t, restPort, testNetworkId)
	getStateNotFoundError(t, restPort, testNetworkId)
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
