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

	"magma/lte/cloud/go/lte"
	lteplugin "magma/lte/cloud/go/plugin"
	"magma/lte/cloud/go/services/cellular/obsidian/models"
	"magma/lte/cloud/go/services/cellular/test_utils"
	"magma/orc8r/cloud/go/obsidian/handlers"
	obsidian_test "magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test_utils "magma/orc8r/cloud/go/services/configurator/test_utils"

	"github.com/stretchr/testify/assert"
)

func TestGetNetworkConfigs(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)

	networkID := "cellular_obsidian_test_network"
	configurator_test_utils.RegisterNetwork(t, networkID, "Test Network 1")

	// Happy path
	expectedConfig := test_utils.NewDefaultFDDNetworkConfig()
	marshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Cellular Network Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:  expected,
		Expected: fmt.Sprintf(`"%s"`, networkID),
	}
	obsidian_test.RunTest(t, createConfigTestCase)

	happyPathTestCase := obsidian_test.Testcase{
		Name:     "Get Cellular Network Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:  "",
		Expected: expected,
	}
	obsidian_test.RunTest(t, happyPathTestCase)

	// No good way to test invalid configs from datastore without dropping down
	// to raw magmad api/grpc or datastore fixtures, so let's skip that
	// for now
}

func TestSetTDDNetworkConfigs(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	testSetNetworkConfigs(t, test_utils.NewDefaultTDDNetworkConfig(), test_utils.NewDefaultTDDNetworkConfig())
}

func TestSetFDDNetworkConfigs(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	testSetNetworkConfigs(t, test_utils.NewDefaultFDDNetworkConfig(), test_utils.NewDefaultFDDNetworkConfig())
}

func TestSetOldTddNetworkConfigs(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	testSetNetworkConfigs(t, test_utils.OldTDDNetworkConfig(), test_utils.NewDefaultTDDNetworkConfig())
}

func TestSetOldFddNetworkConfigs(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	testSetNetworkConfigs(t, test_utils.OldFDDNetworkConfig(), test_utils.NewDefaultFDDNetworkConfig())
}

func TestSetBadNetworkConfigs(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)

	networkID := "cellular_obsidian_test_network"
	configurator_test_utils.RegisterNetwork(t, networkID, "Test Network 1")

	config := test_utils.NewDefaultTDDNetworkConfig()

	// Fail RAN config check
	config.Ran.FddConfig = &models.NetworkRanConfigsFddConfig{
		Earfcndl: 0,
		Earfcnul: 18000,
	}
	marshaledCfg, err := config.MarshalBinary()
	assert.NoError(t, err)
	configString := string(marshaledCfg)

	setConfigTestCase := obsidian_test.Testcase{
		Name:                     "Set Both TDD+FDD Network Config",
		Method:                   "PUT",
		Url:                      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:                  configString,
		Expected:                 `{"message":"Invalid config: Only one of TDD or FDD configs can be set"}`,
		Expect_http_error_status: true,
	}
	status, _, err := obsidian_test.RunTest(t, setConfigTestCase)
	assert.Equal(t, 400, status)

	// Fail swagger validation
	config.Epc.Mcc = "abc"
	config.Ran.BandwidthMhz = 15
	marshaledCfg, err = config.MarshalBinary()
	assert.NoError(t, err)
	configString = string(marshaledCfg)
	setConfigTestCase = obsidian_test.Testcase{
		Name:                     "Set Invalid Cellular Network Config",
		Method:                   "PUT",
		Url:                      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:                  configString,
		Expected:                 `{"message":"Invalid config: validation failure list:\nvalidation failure list:\nmcc in body should match '^(\\d{3})$'"}`,
		Expect_http_error_status: true,
	}
	status, _, err = obsidian_test.RunTest(t, setConfigTestCase)
	assert.Equal(t, 400, status)

	// Fail swagger validation
	config.Epc.Mcc = "123"
	config.Ran.BandwidthMhz = 16
	marshaledCfg, err = config.MarshalBinary()
	assert.NoError(t, err)
	configString = string(marshaledCfg)

	setConfigTestCase = obsidian_test.Testcase{
		Name:                     "Set Invalid Cellular Network Config 2",
		Method:                   "PUT",
		Url:                      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:                  configString,
		Expected:                 `{"message":"Invalid config: validation failure list:\nvalidation failure list:\nbandwidth_mhz in body should be one of [3 5 10 15 20]"}`,
		Expect_http_error_status: true,
	}
	status, _, err = obsidian_test.RunTest(t, setConfigTestCase)
	assert.Equal(t, 400, status)

	// Fail swagger validation
	config = test_utils.NewDefaultTDDNetworkConfig()
	config.Epc.NetworkServices = []string{"metering", "dpi", "bad"}
	marshaledCfg, err = config.MarshalBinary()
	assert.NoError(t, err)
	configString = string(marshaledCfg)

	setConfigTestCase = obsidian_test.Testcase{
		Name:                     "Set Invalid Network Service Name",
		Method:                   "PUT",
		Url:                      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:                  configString,
		Expected:                 `{"message":"Invalid config: validation failure list:\nvalidation failure list:\nnetwork_services.2 in body should be one of [metering dpi policy_enforcement]"}`,
		Expect_http_error_status: true,
	}
	status, _, err = obsidian_test.RunTest(t, setConfigTestCase)
	assert.Equal(t, 400, status)
}

func TestSetBadOldConfigs(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)

	networkID := "cellular_obsidian_test_network"
	configurator_test_utils.RegisterNetwork(t, networkID, "Test Network 1")

	// Fail RAN config check
	config := test_utils.OldTDDNetworkConfig()
	config.Ran.Earfcndl = 125000

	marshaledCfg, err := config.MarshalBinary()
	assert.NoError(t, err)
	configString := string(marshaledCfg)

	setConfigTestCase := obsidian_test.Testcase{
		Name:                     "Set Invalid Earcndl Config",
		Method:                   "POST",
		Url:                      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:                  configString,
		Expected:                 `{"message":"Invalid config: Invalid EARFCNDL: no matching band"}`,
		Expect_http_error_status: true,
	}
	status, _, err := obsidian_test.RunTest(t, setConfigTestCase)
	assert.Equal(t, 400, status)
}

func TestGetGatewayConfigs(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)

	networkID := "cellular_obsidian_test_network"
	configurator_test_utils.RegisterNetwork(t, networkID, "Test Network 1")
	gatewayID := "g1"
	configurator_test_utils.RegisterGateway(t, networkID, gatewayID, nil)
	enodebID := "enb1"
	registerEnodeb(t, networkID, enodebID)

	// Happy path
	expectedConfig := test_utils.NewDefaultGatewayConfig()
	marshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Cellular Gateway Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:  expected,
		Expected: `"g1"`,
	}
	obsidian_test.RunTest(t, createConfigTestCase)

	happyPathTestCase := obsidian_test.Testcase{
		Name:     "Get Cellular Gateway Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:  "",
		Expected: expected,
	}
	obsidian_test.RunTest(t, happyPathTestCase)

	// No good way to test invalid configs from datastore without dropping down
	// to raw magmad api/grpc or datastore fixtures, so let's skip that
	// for now
}

func TestSetGatewayConfigs(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)

	networkID := "cellular_obsidian_test_network"
	configurator_test_utils.RegisterNetwork(t, networkID, "Test Network 1")
	gatewayID := "g2"
	configurator_test_utils.RegisterGateway(t, networkID, gatewayID, nil)
	enodebID := "enb1"
	registerEnodeb(t, networkID, enodebID)
	enodebID2 := "enb2"
	registerEnodeb(t, networkID, enodebID2)

	// Happy path
	gatewayConfig := test_utils.NewDefaultGatewayConfig()
	marshaledCfg, err := gatewayConfig.MarshalBinary()
	assert.NoError(t, err)
	configString := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Cellular Gateway Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:  configString,
		Expected: `"g2"`,
	}
	obsidian_test.RunTest(t, createConfigTestCase)

	gatewayConfig.Ran.Pci = 456
	gatewayConfig.Epc.IPBlock = "192.168.80.10/24" // Make sure filling handles IP properly
	gatewayConfig.NonEpsService.CsfbMcc = "123"
	gatewayConfig.NonEpsService.CsfbMnc = "23"
	gatewayConfig.NonEpsService.Lac = 10
	gatewayConfig.NonEpsService.CsfbRat = 1
	gatewayConfig.NonEpsService.Arfcn2g = []uint32{1, 2, 3}
	gatewayConfig.NonEpsService.NonEpsServiceControl = 2

	marshaledCfg, err = gatewayConfig.MarshalBinary()
	assert.NoError(t, err)
	configString = string(marshaledCfg)

	setConfigTestCase := obsidian_test.Testcase{
		Name:     "Set Cellular Gateway Config",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:  configString,
		Expected: "",
	}
	obsidian_test.RunTest(t, setConfigTestCase)

	// gateway should have an association to lte gateway entity
	magmadGatewayEntity, err := configurator.LoadEntity(
		networkID,
		orc8r.MagmadGatewayType,
		gatewayID,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
	)
	assert.NoError(t, err)
	assert.Equal(t, gatewayID, magmadGatewayEntity.Associations[0].Key)
	assert.Equal(t, lte.CellularGatewayType, magmadGatewayEntity.Associations[0].Type)

	// lte gateway should have an association to enodeb
	lteGatewayEntity, err := configurator.LoadEntity(
		networkID,
		lte.CellularGatewayType,
		gatewayID,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
	)
	assert.NoError(t, err)
	assert.Equal(t, enodebID, lteGatewayEntity.Associations[0].Key)
	assert.Equal(t, lte.CellularEnodebType, lteGatewayEntity.Associations[0].Type)

	getConfigTestCase := obsidian_test.Testcase{
		Name:     "Get Updated Cellular Gateway Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:  "",
		Expected: configString,
	}
	obsidian_test.RunTest(t, getConfigTestCase)

	// Set new configs and remove enodeb associations
	gatewayConfig.AttachedEnodebSerials = []string{}
	marshaledCfg, err = gatewayConfig.MarshalBinary()
	assert.NoError(t, err)
	configString = string(marshaledCfg)

	setConfigTestCase = obsidian_test.Testcase{
		Name:     "Set Cellular Gateway Config And Remove Enodeb Associations",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:  configString,
		Expected: "",
	}
	obsidian_test.RunTest(t, setConfigTestCase)

	// lte gateway should have no association to enodeb
	lteGatewayEntity, err = configurator.LoadEntity(
		networkID,
		lte.CellularGatewayType,
		gatewayID,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
	)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(lteGatewayEntity.Associations))

	// Fail proto validation (no swagger validation on gateway configs)
	gatewayConfig.Ran = nil
	marshaledCfg, err = gatewayConfig.MarshalBinary()
	assert.NoError(t, err)
	configString = string(marshaledCfg)

	setConfigTestCase = obsidian_test.Testcase{
		Name:                     "Set Invalid Cellular Gateway Config",
		Method:                   "PUT",
		Url:                      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:                  configString,
		Expected:                 `{"message":"Invalid config: Gateway RAN config is nil"}`,
		Expect_http_error_status: true,
	}
	status, _, err := obsidian_test.RunTest(t, setConfigTestCase)
	assert.Equal(t, 400, status)
}

func testSetNetworkConfigs(t *testing.T, config *models.NetworkCellularConfigs, expectedConfig *models.NetworkCellularConfigs) {
	restPort := obsidian_test.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)
	networkID := "cellular_obsidian_test_network"
	configurator_test_utils.RegisterNetwork(t, networkID, "Test Network 1")

	// Happy path
	marshaledCfg, err := config.MarshalBinary()
	assert.NoError(t, err)
	configString := string(marshaledCfg)

	createConfigTestCase := obsidian_test.Testcase{
		Name:     "Create Cellular Network Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:  configString,
		Expected: fmt.Sprintf(`"%s"`, networkID),
	}
	obsidian_test.RunTest(t, createConfigTestCase)

	config.Epc.Mcc = "123"
	config.Epc.SubProfiles = make(map[string]models.NetworkEpcConfigsSubProfilesAnon)
	config.Epc.SubProfiles["test"] =
		models.NetworkEpcConfigsSubProfilesAnon{
			MaxUlBitRate: 100, MaxDlBitRate: 200,
		}
	config.Ran.BandwidthMhz = 15

	expectedConfig.Epc.Mcc = "123"
	expectedConfig.Epc.SubProfiles = make(map[string]models.NetworkEpcConfigsSubProfilesAnon)
	expectedConfig.Epc.SubProfiles["test"] =
		models.NetworkEpcConfigsSubProfilesAnon{
			MaxUlBitRate: 100, MaxDlBitRate: 200,
		}
	expectedConfig.Ran.BandwidthMhz = 15

	config.Epc.NetworkServices = []string{"metering", "dpi", "policy_enforcement"}

	expectedConfig.Epc.NetworkServices = []string{"metering", "dpi", "policy_enforcement"}

	marshaledCfg, err = config.MarshalBinary()
	assert.NoError(t, err)
	configString = string(marshaledCfg)

	exMarshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	exSwaggerConfigString := string(exMarshaledCfg)

	setConfigTestCase := obsidian_test.Testcase{
		Name:     "Set Cellular Network Config",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:  configString,
		Expected: "",
	}
	obsidian_test.RunTest(t, setConfigTestCase)
	getConfigTestCase := obsidian_test.Testcase{
		Name:     "Get Updated Cellular Network Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:  "",
		Expected: exSwaggerConfigString,
	}
	obsidian_test.RunTest(t, getConfigTestCase)
}

func registerEnodeb(t *testing.T, networkID, enodebID string) {
	enodebEntity := configurator.NetworkEntity{
		Key:  enodebID,
		Type: lte.CellularEnodebType,
	}
	_, err := configurator.CreateEntity(networkID, enodebEntity)
	assert.NoError(t, err)
}
