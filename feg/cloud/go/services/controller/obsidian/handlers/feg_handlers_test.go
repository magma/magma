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

	fegplugin "magma/feg/cloud/go/plugin"
	"magma/feg/cloud/go/services/controller/test_utils"
	"magma/orc8r/cloud/go/obsidian"
	obsidian_test "magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/models"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test_utils "magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"

	"github.com/stretchr/testify/assert"
)

func TestGetNetworkConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &fegplugin.FegOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)
	networkID := "feg_obsidian_test_network"
	registerNetworkWithDefaultConfig(t, "Test Network 1", networkID, restPort)

	// Happy path
	config := test_utils.NewDefaultNetworkConfig()
	marshaledConfig, err := config.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledConfig)
	happyPathTestCase := obsidian_test.Testcase{
		Name:     "Get FeG Network Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/federation", testUrlRoot, networkID),
		Payload:  "",
		Expected: expected,
	}
	obsidian_test.RunTest(t, happyPathTestCase)

	deleteConfigTestCase := obsidian_test.Testcase{
		Name:   "Delete Federation Network Config",
		Method: "DELETE",
		Url:    fmt.Sprintf("%s/%s/configs/federation", testUrlRoot, networkID),
	}
	_, _, err = obsidian_test.RunTest(t, deleteConfigTestCase)
	assert.NoError(t, err)

	// No good way to test invalid configs from datastore without dropping down
	// to raw magmad api/grpc or datastore fixtures, so let's skip that
	// for now
}

func TestSetNetworkConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &fegplugin.FegOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "feg_obsidian_test_network"
	registerNetworkWithDefaultConfig(t, "Test Network 1", networkID, restPort)

	// Happy path
	config := test_utils.NewDefaultNetworkConfig()
	config.S6a.Server.Address = "192.168.11.22:555"
	config.Gx.Server.DestHost = "pcrf.mno.com"
	config.Gy.Server.DestHost = "ocs.mno.com"
	config.ServedNetworkIds = []string{"lte_network_A", "lte_network_B"}
	marshaledConfig, err := config.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledConfig)

	setConfigTestCase := obsidian_test.Testcase{
		Name:     "Set Federation Network Config",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/configs/federation", testUrlRoot, networkID),
		Payload:  expected,
		Expected: "",
	}
	obsidian_test.RunTest(t, setConfigTestCase)
	getConfigTestCase := obsidian_test.Testcase{
		Name:     "Get Updated Federation Network Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/federation", testUrlRoot, networkID),
		Payload:  "",
		Expected: expected,
	}
	obsidian_test.RunTest(t, getConfigTestCase)

	// Fail swagger validation
	config.S6a.Server.Protocol = "foobar"
	marshaledConfig, err = config.MarshalBinary()
	assert.NoError(t, err)
	expected = string(marshaledConfig)

	setConfigTestCase = obsidian_test.Testcase{
		Name:                     "Set Invalid Federation Network Config",
		Method:                   "PUT",
		Url:                      fmt.Sprintf("%s/%s/configs/federation", testUrlRoot, networkID),
		Payload:                  expected,
		Expected:                 `{"message":"Invalid config: validation failure list:\nvalidation failure list:\nvalidation failure list:\nprotocol in body should be one of [tcp tcp4 tcp6 sctp sctp4 sctp6]"}`,
		Expect_http_error_status: true,
	}
	status, _, err := obsidian_test.RunTest(t, setConfigTestCase)
	assert.Equal(t, 400, status)

	deleteConfigTestCase := obsidian_test.Testcase{
		Name:   "Delete Federation Network Config",
		Method: "DELETE",
		Url:    fmt.Sprintf("%s/%s/configs/federation", testUrlRoot, networkID),
	}
	_, _, err = obsidian_test.RunTest(t, deleteConfigTestCase)
	assert.NoError(t, err)
}

func TestGetGatewayConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &fegplugin.FegOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "feg_obsidian_test_network"
	registerNetworkWithDefaultConfig(t, "Test Network 1", networkID, restPort)
	gatewayID := "g1"
	registerGatewayWithDefaultConfig(t, networkID, gatewayID, restPort)

	// Happy path
	config := test_utils.NewDefaultGatewayConfig()
	marshaledConfig, err := config.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledConfig)
	happyPathTestCase := obsidian_test.Testcase{
		Name:     "Get Federation Gateway Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/federation", testUrlRoot, networkID, gatewayID),
		Payload:  "",
		Expected: expected,
	}
	obsidian_test.RunTest(t, happyPathTestCase)

	deleteConfigTestCase := obsidian_test.Testcase{
		Name:   "Delete Federation Gateway Config",
		Method: "DELETE",
		Url:    fmt.Sprintf("%s/%s/gateways/%s/configs/federation", testUrlRoot, networkID, gatewayID),
	}
	_, _, err = obsidian_test.RunTest(t, deleteConfigTestCase)
	assert.NoError(t, err)

	// No good way to test invalid configs from datastore without dropping down
	// to raw magmad api/grpc or datastore fixtures, so let's skip that
	// for now
}

func TestSetGatewayConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &fegplugin.FegOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	device_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "feg_obsidian_test_network"
	registerNetworkWithDefaultConfig(t, "Test Network 1", networkID, restPort)
	gatewayID := "g2"
	registerGatewayWithDefaultConfig(t, networkID, gatewayID, restPort)

	// Happy path
	config := test_utils.NewDefaultGatewayConfig()
	config.S6a.Server.Address = "192.168.11.22:555"
	marshaledConfig, err := config.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledConfig)

	setConfigTestCase := obsidian_test.Testcase{
		Name:     "Set Federation Gateway Config",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/federation", testUrlRoot, networkID, gatewayID),
		Payload:  expected,
		Expected: "",
	}
	obsidian_test.RunTest(t, setConfigTestCase)
	getConfigTestCase := obsidian_test.Testcase{
		Name:     "Get Updated Federation Gateway Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/federation", testUrlRoot, networkID, gatewayID),
		Payload:  "",
		Expected: expected,
	}
	obsidian_test.RunTest(t, getConfigTestCase)

	deleteConfigTestCase := obsidian_test.Testcase{
		Name:   "Delete Federation Gateway Config",
		Method: "DELETE",
		Url:    fmt.Sprintf("%s/%s/gateways/%s/configs/federation", testUrlRoot, networkID, gatewayID),
	}
	_, _, err = obsidian_test.RunTest(t, deleteConfigTestCase)
	assert.NoError(t, err)
}

func registerNetworkWithDefaultConfig(t *testing.T, networkName string, networkID string, port int) {
	configurator_test_utils.RegisterNetwork(t, networkID, networkName)

	config := test_utils.NewDefaultNetworkConfig()
	marshaledConfig, err := config.MarshalBinary()
	assert.NoError(t, err)

	_, _, err = obsidian_test.RunTest(t, obsidian_test.Testcase{
		Name:   "Create Default Federation Network Config",
		Method: "POST",
		Url: fmt.Sprintf("http://localhost:%d%s/networks/%s/configs/federation",
			port, obsidian.RestRoot, networkID),
		Payload:  string(marshaledConfig),
		Expected: "\"" + networkID + "\"",
	})
	assert.NoError(t, err)
}

func registerGatewayWithDefaultConfig(t *testing.T, networkID string, gatewayID string, port int) {
	gatewayRecord := &models.GatewayDevice{HardwareID: gatewayID}
	configurator_test_utils.RegisterGateway(t, networkID, gatewayID, gatewayRecord)

	config := test_utils.NewDefaultGatewayConfig()
	marshaledConfig, err := config.MarshalBinary()
	assert.NoError(t, err)

	_, _, err = obsidian_test.RunTest(t, obsidian_test.Testcase{
		Name:   "Create Default Federation Gateway Config",
		Method: "POST",
		Url: fmt.Sprintf(
			"http://localhost:%d%s/networks/%s/gateways/%s/configs/federation",
			port, obsidian.RestRoot, networkID, gatewayID),
		Payload:  string(marshaledConfig),
		Expected: "\"" + gatewayID + "\"",
	})
	assert.NoError(t, err)
}
