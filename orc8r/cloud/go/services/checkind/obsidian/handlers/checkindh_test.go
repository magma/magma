/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers_test

import (
	"fmt"
	"os"
	"testing"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/protos"
	checkindTestInit "magma/orc8r/cloud/go/services/checkind/test_init"
	"magma/orc8r/cloud/go/services/checkind/test_utils"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/magmad"
	magmadProtos "magma/orc8r/cloud/go/services/magmad/protos"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	stateTestUtils "magma/orc8r/cloud/go/services/state/test_utils"

	"github.com/stretchr/testify/assert"
)

const testAgHwId = "Test-AGW-Hw-Id"

// TestCheckind is Obsidian Gateway Status Integration Test intended to be run
// on cloud VM
func TestCheckind(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configuratorTestInit.StartTestService(t)
	checkindTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	restPort := tests.StartObsidian(t)

	// create a test network with a single GW
	testNetworkID := registerNetwork(t, "Test Network 1", "checkind_obsidian_test_network")

	t.Logf("New Registered Network: %s", testNetworkID)

	logicalID := registerGateway(t, testNetworkID, testAgHwId, testAgHwId, "Test GW Name")

	ctx := stateTestUtils.GetContextWithCertificate(t, testAgHwId)

	// put one checkin state into state service
	gwStatus := test_utils.GetGatewayStatusSwaggerFixture(testAgHwId)

	stateTestUtils.ReportGatewayStatus(t, ctx, gwStatus)

	getGWStatusNoError(t, restPort, testNetworkID, logicalID)
	getGWStatusNotFoundError(t, restPort, testNetworkID)

	err := configurator.DeleteNetwork(testNetworkID)
	assert.NoError(t, err)
}

func getURL(restPort int, networkID string, logicalID string) string {
	url := fmt.Sprintf(
		"http://localhost:%d%s/networks/%s/gateways/%s/status",
		restPort,
		handlers.REST_ROOT,
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

func registerNetwork(t *testing.T, networkName string, networkID string) string {
	useNewHandler := os.Getenv(orc8r.UseConfiguratorEnv)
	if useNewHandler == "1" {
		err := configurator.CreateNetwork(
			configurator.Network{
				Name: networkName,
				ID:   networkID,
			},
		)
		assert.NoError(t, err)
		return networkID
	} else {
		networkId, err := magmad.RegisterNetwork(
			&magmadProtos.MagmadNetworkRecord{Name: networkName},
			networkID)
		assert.NoError(t, err)
		return networkId
	}
}

func registerGateway(t *testing.T, networkID string, gatewayID string, hwID string, name string) string {
	useNewHandler := os.Getenv(orc8r.UseConfiguratorEnv)
	if useNewHandler == "1" {
		_, err := configurator.CreateEntity(networkID, configurator.NetworkEntity{
			Key:        gatewayID,
			Type:       "magmad_gateway",
			PhysicalID: hwID,
		})
		assert.NoError(t, err)
		return gatewayID
	} else {
		gatewayRecord := &magmadProtos.AccessGatewayRecord{
			HwId: &protos.AccessGatewayID{Id: gatewayID},
			Name: name,
		}
		registeredId, err := magmad.RegisterGateway(networkID, gatewayRecord)
		assert.NoError(t, err)
		return registeredId
	}
}
