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
	"magma/feg/cloud/go/services/controller/obsidian/models"
	feg_protos "magma/feg/cloud/go/services/controller/protos"
	"magma/orc8r/cloud/go/obsidian/handlers"
	obsidian_test "magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/stretchr/testify/assert"
)

func TestLegacyGetNetworkConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &fegplugin.FegOrchestratorPlugin{})
	magmad_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)
	networkId := registerNetwork(t, "Test Network 1", "feg_obsidian_test_network", restPort)

	// Happy path
	expectedConfig := &models.NetworkFederationConfigs{}
	expectedConfig.FromServiceModel(feg_protos.NewDefaultProtosGatewayConfig())
	marshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)
	happyPathTestCase := obsidian_test.Testcase{
		Name:     "Get FeG Network Config Legacy",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/federation", testUrlRoot, networkId),
		Payload:  "",
		Expected: expected,
	}
	obsidian_test.RunTest(t, happyPathTestCase)

	// No good way to test invalid configs from datastore without dropping down
	// to raw magmad api/grpc or datastore fixtures, so let's skip that
	// for now
}

func TestLegacySetNetworkConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &fegplugin.FegOrchestratorPlugin{})
	magmad_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)

	networkId := registerNetwork(t, "Test Network 1", "feg_obsidian_test_network", restPort)

	// Happy path
	config := feg_protos.NewDefaultProtosNetworkConfig()
	config.S6A.Server.Address = "192.168.11.22:555"
	config.Gx.Server.DestHost = "pcrf.mno.com"
	config.Gy.Server.DestHost = "ocs.mno.com"
	config.ServedNetworkIds = []string{"lte_network_A", "lte_network_B"}
	swaggerConfig := &models.NetworkFederationConfigs{}
	swaggerConfig.FromServiceModel(config)
	assert.Len(t, swaggerConfig.ServedNetworkIds, 2)
	assert.Subset(t, swaggerConfig.ServedNetworkIds, config.ServedNetworkIds)
	marshaledCfg, err := swaggerConfig.MarshalBinary()
	assert.NoError(t, err)
	swaggerConfigString := string(marshaledCfg)

	setConfigTestCase := obsidian_test.Testcase{
		Name:     "Set Federation Network Config Legacy",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/configs/federation", testUrlRoot, networkId),
		Payload:  swaggerConfigString,
		Expected: "",
	}
	obsidian_test.RunTest(t, setConfigTestCase)
	getConfigTestCase := obsidian_test.Testcase{
		Name:     "Get Updated Federation Network Config Legacy",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/federation", testUrlRoot, networkId),
		Payload:  "",
		Expected: swaggerConfigString,
	}
	obsidian_test.RunTest(t, getConfigTestCase)

	// Fail swagger validation
	config.S6A.Server.Protocol = "foobar"
	swaggerConfig.FromServiceModel(config)
	marshaledCfg, err = swaggerConfig.MarshalBinary()
	assert.NoError(t, err)
	swaggerConfigString = string(marshaledCfg)

	setConfigTestCase = obsidian_test.Testcase{
		Name:                     "Set Invalid Federation Network Config Legacy",
		Method:                   "PUT",
		Url:                      fmt.Sprintf("%s/%s/configs/federation", testUrlRoot, networkId),
		Payload:                  swaggerConfigString,
		Expected:                 `{"message":"Invalid config: validation failure list:\nvalidation failure list:\nvalidation failure list:\nprotocol in body should be one of [tcp tcp4 tcp6 sctp sctp4 sctp6]"}`,
		Expect_http_error_status: true,
	}
	status, _, err := obsidian_test.RunTest(t, setConfigTestCase)
	assert.Equal(t, 400, status)

}

func TestLegacyGetGatewayConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &fegplugin.FegOrchestratorPlugin{})
	magmad_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)

	networkId := registerNetwork(t, "Test Network 1", "feg_obsidian_test_network", restPort)
	gatewayId := registerGateway(t, networkId, "g1", restPort)

	// Happy path
	expectedConfig := &models.GatewayFegConfigs{}
	expectedConfig.FromServiceModel(feg_protos.NewDefaultProtosGatewayConfig())
	marshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)
	happyPathTestCase := obsidian_test.Testcase{
		Name:     "Get Federation Gateway Config Legacy",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/federation", testUrlRoot, networkId, gatewayId),
		Payload:  "",
		Expected: expected,
	}
	obsidian_test.RunTest(t, happyPathTestCase)

	// No good way to test invalid configs from datastore without dropping down
	// to raw magmad api/grpc or datastore fixtures, so let's skip that
	// for now
}

func TestLegacySetGatewayConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &fegplugin.FegOrchestratorPlugin{})
	magmad_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)

	networkId := registerNetwork(t, "Test Network 1", "feg_obsidian_test_network", restPort)
	gatewayId := registerGateway(t, networkId, "g2", restPort)

	// Happy path
	gatewayConfig := feg_protos.NewDefaultProtosGatewayConfig()
	gatewayConfig.S6A.Server.Address = "192.168.11.22:555"
	swaggerConfig := &models.GatewayFegConfigs{}
	swaggerConfig.FromServiceModel(gatewayConfig)

	assert.Equal(t, gatewayConfig.S6A.Server.Address, swaggerConfig.S6a.Server.Address)

	marshaledCfg, err := swaggerConfig.MarshalBinary()
	assert.NoError(t, err)
	swaggerConfigString := string(marshaledCfg)

	setConfigTestCase := obsidian_test.Testcase{
		Name:     "Set Federation Gateway Config Legacy",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/federation", testUrlRoot, networkId, gatewayId),
		Payload:  swaggerConfigString,
		Expected: "",
	}
	obsidian_test.RunTest(t, setConfigTestCase)
	getConfigTestCase := obsidian_test.Testcase{
		Name:     "Get Updated Federation Gateway Config Legacy",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/federation", testUrlRoot, networkId, gatewayId),
		Payload:  "",
		Expected: swaggerConfigString,
	}
	obsidian_test.RunTest(t, getConfigTestCase)
}

func registerNetwork(t *testing.T, networkName string, networkId string, port int) string {
	networkId, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: networkName},
		networkId)
	assert.NoError(t, err)

	config := feg_protos.NewDefaultProtosNetworkConfig()
	swaggerConfig := &models.NetworkFederationConfigs{}
	err = swaggerConfig.FromServiceModel(config)
	assert.NoError(t, err)
	marshaledCfg, err := swaggerConfig.MarshalBinary()
	assert.NoError(t, err)
	swaggerConfigString := string(marshaledCfg)

	obsidian_test.RunTest(t, obsidian_test.Testcase{
		Name:   "Create Default Federation Network Config Legacy",
		Method: "POST",
		Url: fmt.Sprintf("http://localhost:%d%s/networks/%s/configs/federation",
			port, handlers.REST_ROOT, networkId),
		Payload:  swaggerConfigString,
		Expected: "\"" + networkId + "\"",
	})
	return networkId
}

func registerGateway(t *testing.T, networkId string, gatewayId string, port int) string {
	gatewayRecord := &magmad_protos.AccessGatewayRecord{
		HwId: &protos.AccessGatewayID{Id: gatewayId},
	}
	registeredId, err := magmad.RegisterGateway(networkId, gatewayRecord)
	assert.NoError(t, err)

	config := feg_protos.NewDefaultProtosGatewayConfig()
	swaggerConfig := &models.GatewayFegConfigs{}
	err = swaggerConfig.FromServiceModel(config)
	assert.NoError(t, err)
	marshaledCfg, err := swaggerConfig.MarshalBinary()
	assert.NoError(t, err)
	swaggerConfigString := string(marshaledCfg)

	obsidian_test.RunTest(t, obsidian_test.Testcase{
		Name:   "Create Default Federation Gateway Config Legacy",
		Method: "POST",
		Url: fmt.Sprintf(
			"http://localhost:%d%s/networks/%s/gateways/%s/configs/federation",
			port, handlers.REST_ROOT, networkId, registeredId),
		Payload:  swaggerConfigString,
		Expected: "\"" + registeredId + "\"",
	})
	return registeredId
}
