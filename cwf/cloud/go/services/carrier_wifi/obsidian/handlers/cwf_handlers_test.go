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
	"testing"

	cwfplugin "magma/cwf/cloud/go/plugin"
	"magma/cwf/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/obsidian"
	obsidian_test "magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test_utils "magma/orc8r/cloud/go/services/configurator/test_utils"

	"github.com/stretchr/testify/assert"
)

func TestGetNetworkConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &cwfplugin.CwfOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "cwf_obsidian_test_network"
	configurator_test_utils.RegisterNetwork(t, networkID, "Test Network 1")

	// Happy path
	expectedConfig := newDefaultCwfNetworkConfig()
	marshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Carrier WiFi Network Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/configs/carrier_wifi", testUrlRoot, networkID),
		Payload:  expected,
		Expected: fmt.Sprintf(`"%s"`, networkID),
	}
	obsidian_test.RunTest(t, createConfigTestCase)

	happyPathTestCase := obsidian_test.Testcase{
		Name:     "Get Carrier WiFi Network Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/carrier_wifi", testUrlRoot, networkID),
		Payload:  "",
		Expected: expected,
	}
	obsidian_test.RunTest(t, happyPathTestCase)
}

func TestSetNetworkConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &cwfplugin.CwfOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)

	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "cwf_obsidian_test_network"
	configurator_test_utils.RegisterNetwork(t, networkID, "Test Network 1")

	// Happy path
	expectedConfig := newDefaultCwfNetworkConfig()
	marshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Carrier WiFi Network Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/configs/carrier_wifi", testUrlRoot, networkID),
		Payload:  expected,
		Expected: fmt.Sprintf(`"%s"`, networkID),
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
		Url:      fmt.Sprintf("%s/%s/configs/carrier_wifi", testUrlRoot, networkID),
		Payload:  expected2,
		Expected: fmt.Sprintf(`"%s"`, networkID),
	}
	obsidian_test.RunTest(t, updateConfigTestCase)

	updateTestCaseResult := obsidian_test.Testcase{
		Name:     "Get Updated Carrier WiFi Network Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/carrier_wifi", testUrlRoot, networkID),
		Payload:  "",
		Expected: expected2,
	}

	obsidian_test.RunTest(t, updateTestCaseResult)
}

func newDefaultCwfNetworkConfig() *models.NetworkCarrierWifiConfigs {
	return &models.NetworkCarrierWifiConfigs{
		AaaServer: &models.AaaServer{
			AccountingEnabled:    false,
			CreateSessionOnAuth:  false,
			IDLESessionTimeoutMs: 21600000,
		},
		EapAka: &models.EapAka{
			PlmnIds: []string{},
			Timeout: &models.EapAkaTimeout{
				ChallengeMs:            20000,
				ErrorNotificationMs:    10000,
				SessionMs:              43200000,
				SessionAuthenticatedMs: 5000,
			},
		},
		DefaultRuleID: nil,
	}
}
