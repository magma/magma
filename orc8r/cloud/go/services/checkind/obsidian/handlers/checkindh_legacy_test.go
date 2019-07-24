/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers_test

import (
	"os"
	"testing"

	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/protos"
	checkindTestInit "magma/orc8r/cloud/go/services/checkind/test_init"
	"magma/orc8r/cloud/go/services/checkind/test_utils"
	"magma/orc8r/cloud/go/services/magmad"
	magmadProtos "magma/orc8r/cloud/go/services/magmad/protos"
	magmadTestInit "magma/orc8r/cloud/go/services/magmad/test_init"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	stateTestUtils "magma/orc8r/cloud/go/services/state/test_utils"

	"github.com/stretchr/testify/assert"
)

// TestCheckind is Obsidian Gateway Status Integration Test intended to be run
// on cloud VM
func TestCheckindLegacy(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "0")
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	magmadTestInit.StartTestService(t)
	checkindTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	restPort := tests.StartObsidian(t)

	// create a test network with a single GW
	testNetworkID, err := magmad.RegisterNetwork(
		&magmadProtos.MagmadNetworkRecord{Name: "Test Network 1"},
		"checkind_obsidian_test_network")
	assert.NoError(t, err)

	t.Logf("New Registered Network: %s", testNetworkID)

	hwID := protos.AccessGatewayID{Id: testAgHwId}
	logicalID, err := magmad.RegisterGateway(testNetworkID, &magmadProtos.AccessGatewayRecord{HwId: &hwID, Name: "Test GW Name"})
	assert.NoError(t, err)
	assert.NotEqual(t, logicalID, "")

	ctx := stateTestUtils.GetContextWithCertificate(t, testAgHwId)

	// put one checkin state into state service
	gwStatus := test_utils.GetGatewayStatusSwaggerFixture(testAgHwId)
	stateTestUtils.ReportGatewayStatus(t, ctx, gwStatus)

	getGWStatusNoError(t, restPort, testNetworkID, logicalID)
	getGWStatusNotFoundError(t, restPort, testNetworkID)

	magmad.ForceRemoveNetwork(testNetworkID)
}
