/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
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
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	stateTestUtils "magma/orc8r/cloud/go/services/state/test_utils"

	"github.com/stretchr/testify/assert"
)

const testAgHwId = "Test-AGW-Hw-Id"

// TestCheckind is Obsidian Gateway Status Integration Test intended to be run
// on cloud VM
func TestCheckind(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	restPort := tests.StartObsidian(t)

	// create a test network with a single GW
	networkID := "checkind_obsidian_test_network"
	err := configurator.CreateNetwork(
		configurator.Network{
			Name: "Test Network 1",
			ID:   networkID,
		},
	)
	assert.NoError(t, err)

	_, err = configurator.CreateEntity(networkID, configurator.NetworkEntity{
		Key:        testAgHwId,
		Type:       "magmad_gateway",
		PhysicalID: testAgHwId,
	})
	assert.NoError(t, err)

	// put one checkin state into state service
	ctx := stateTestUtils.GetContextWithCertificate(t, testAgHwId)
	gwStatus := models.NewDefaultGatewayStatus(testAgHwId)
	stateTestUtils.ReportGatewayStatus(t, ctx, gwStatus)

	getGWStatusNoError(t, restPort, networkID, testAgHwId)
	getGWStatusNotFoundError(t, restPort, networkID)

	err = configurator.DeleteNetwork(networkID)
	assert.NoError(t, err)
}

func getURL(restPort int, networkID string, logicalID string) string {
	url := fmt.Sprintf(
		"http://localhost:%d%s/networks/%s/gateways/%s/status",
		restPort,
		obsidian.RestRoot,
		networkID,
		logicalID,
	)
	return url
}

func getGWStatusNoError(t *testing.T, restPort int, networkID string, logicalID string) {
	url := getURL(restPort, networkID, logicalID)
	stateTestUtils.GetGWStatusViaHTTPNoError(t, url, networkID, logicalID)
}

func getGWStatusNotFoundError(t *testing.T, restPort int, networkID string) {
	url := getURL(restPort, networkID, "should-not-exist")
	stateTestUtils.GetGWStatusExpectNotFound(t, url)
}
