/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers_test

import (
	"fmt"
	"os"
	"testing"

	cwfplugin "magma/cwf/cloud/go/plugin"
	"magma/cwf/cloud/go/services/carrier_wifi/obsidian/models"
	"magma/orc8r/cloud/go/obsidian/handlers"
	obsidian_test "magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/magmad"

	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/stretchr/testify/assert"
)

func TestGetNetworkConfigs(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &cwfplugin.CwfOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	magmad_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)

	networkId := registerNetwork(t, "Test Network 1", "cwf_obsidian_test_network")

	// Happy path
	expectedConfig := newDefaultCwfNetworkConfig()
	marshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Carrier WiFi Network Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/configs/carrier_wifi", testUrlRoot, networkId),
		Payload:  expected,
		Expected: fmt.Sprintf(`"%s"`, networkId),
	}
	obsidian_test.RunTest(t, createConfigTestCase)

	happyPathTestCase := obsidian_test.Testcase{
		Name:     "Get Carrier WiFi Network Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/carrier_wifi", testUrlRoot, networkId),
		Payload:  "",
		Expected: expected,
	}
	obsidian_test.RunTest(t, happyPathTestCase)
}

func TestSetNetworkConfigs(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &cwfplugin.CwfOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	magmad_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)

	networkId := registerNetwork(t, "Test Network 1", "cellular_obsidian_test_network")

	// Happy path
	expectedConfig := newDefaultCwfNetworkConfig()
	marshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Carrier WiFi Network Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/configs/carrier_wifi", testUrlRoot, networkId),
		Payload:  expected,
		Expected: fmt.Sprintf(`"%s"`, networkId),
	}

	obsidian_test.RunTest(t, createConfigTestCase)

	updatedConfig := newDefaultCwfNetworkConfig()
	updatedConfig.EapAka.Timeout.SessionAuthenticatedMs = 1
	updatedConfig.EapAka.Timeout.ChallengeMs = 5
	updatedConfig.AaaServer.CreateSessionOnAuth = true
	updatedConfig.AaaServer.AccountingEnabled = true

	marshaledUpdatedCfg, err := updatedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected2 := string(marshaledUpdatedCfg)

	updateConfigTestCase := obsidian_test.Testcase{
		Name:     "Update Carrier WiFi Network Config",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/configs/carrier_wifi", testUrlRoot, networkId),
		Payload:  expected2,
		Expected: fmt.Sprintf(`"%s"`, networkId),
	}
	obsidian_test.RunTest(t, updateConfigTestCase)

	updateTestCaseResult := obsidian_test.Testcase{
		Name:     "Get Updated Carrier WiFi Network Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/carrier_wifi", testUrlRoot, networkId),
		Payload:  "",
		Expected: expected2,
	}

	obsidian_test.RunTest(t, updateTestCaseResult)
}

func newDefaultCwfNetworkConfig() *models.NetworkCarrierWifiConfigs {
	return &models.NetworkCarrierWifiConfigs{
		RelayEnabled: true,
		AaaServer: &models.NetworkCarrierWifiConfigsAaaServer{
			AccountingEnabled:    false,
			CreateSessionOnAuth:  false,
			IDLESessionTimeoutMs: 21600000,
		},
		EapAka: &models.NetworkCarrierWifiConfigsEapAka{
			PlmnIds: []string{},
			Timeout: &models.EapAkaTimeouts{
				ChallengeMs:            20000,
				ErrorNotificationMs:    10000,
				SessionMs:              43200000,
				SessionAuthenticatedMs: 5000,
			},
		},
		FegNetworkID: "feg_network",
	}
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
			&magmad_protos.MagmadNetworkRecord{Name: networkName},
			networkID)
		assert.NoError(t, err)
		return networkId
	}
}
