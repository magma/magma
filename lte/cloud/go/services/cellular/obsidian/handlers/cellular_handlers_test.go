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

	"magma/lte/cloud/go/lte"
	lteplugin "magma/lte/cloud/go/plugin"
	models2 "magma/lte/cloud/go/plugin/models"
	"magma/lte/cloud/go/services/cellular/test_utils"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	configuratorTestUtils "magma/orc8r/cloud/go/services/configurator/test_utils"

	"github.com/stretchr/testify/assert"
)

func TestGetNetworkConfigs(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configuratorTestInit.StartTestService(t)
	restPort := tests.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "cellular_obsidian_test_network"
	configuratorTestUtils.RegisterNetwork(t, networkID, "Test Network 1")

	// Happy path
	expectedConfig := models2.NewDefaultTDDNetworkConfig()
	marshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)

	createConfigTestCase := tests.Testcase{
		Name:     "Create Cellular Network Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:  expected,
		Expected: fmt.Sprintf(`"%s"`, networkID),
	}
	tests.RunTest(t, createConfigTestCase)

	happyPathTestCase := tests.Testcase{
		Name:     "Get Cellular Network Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:  "",
		Expected: expected,
	}
	tests.RunTest(t, happyPathTestCase)

	// No good way to test invalid configs from datastore without dropping down
	// to raw magmad api/grpc or datastore fixtures, so let's skip that
	// for now
}

func TestSetTDDNetworkConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configuratorTestInit.StartTestService(t)
	testSetNetworkConfigs(t, models2.NewDefaultTDDNetworkConfig(), models2.NewDefaultTDDNetworkConfig())
}

func TestSetFDDNetworkConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configuratorTestInit.StartTestService(t)
	testSetNetworkConfigs(t, models2.NewDefaultFDDNetworkConfig(), models2.NewDefaultFDDNetworkConfig())
}

func TestGetGatewayConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configuratorTestInit.StartTestService(t)
	restPort := tests.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "cellular_obsidian_test_network"
	configuratorTestUtils.RegisterNetwork(t, networkID, "Test Network 1")
	gatewayID := "g1"
	configuratorTestUtils.RegisterGateway(t, networkID, gatewayID, nil)
	enodebID := "enb1"
	registerEnodeb(t, networkID, enodebID)

	// Happy path
	expectedConfig := test_utils.NewDefaultGatewayConfig()
	marshaledCfg, err := expectedConfig.MarshalBinary()
	assert.NoError(t, err)
	expected := string(marshaledCfg)

	createConfigTestCase := tests.Testcase{
		Name:     "Create Cellular Gateway Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:  expected,
		Expected: `"g1"`,
	}
	tests.RunTest(t, createConfigTestCase)

	happyPathTestCase := tests.Testcase{
		Name:     "Get Cellular Gateway Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:  "",
		Expected: expected,
	}
	tests.RunTest(t, happyPathTestCase)

	// No good way to test invalid configs from datastore without dropping down
	// to raw magmad api/grpc or datastore fixtures, so let's skip that
	// for now
}

func TestSetGatewayConfigs(t *testing.T) {
	plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configuratorTestInit.StartTestService(t)
	restPort := tests.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	networkID := "cellular_obsidian_test_network"
	configuratorTestUtils.RegisterNetwork(t, networkID, "Test Network 1")
	gatewayID := "g2"
	configuratorTestUtils.RegisterGateway(t, networkID, gatewayID, nil)
	enodebID := "enb1"
	registerEnodeb(t, networkID, enodebID)
	enodebID2 := "enb2"
	registerEnodeb(t, networkID, enodebID2)

	// Happy path
	gatewayConfig := test_utils.NewDefaultGatewayConfig()
	marshaledCfg, err := gatewayConfig.MarshalBinary()
	assert.NoError(t, err)
	configString := string(marshaledCfg)

	createConfigTestCase := tests.Testcase{
		Name:     "Create Cellular Gateway Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:  configString,
		Expected: `"g2"`,
	}
	tests.RunTest(t, createConfigTestCase)

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

	setConfigTestCase := tests.Testcase{
		Name:     "Set Cellular Gateway Config",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:  configString,
		Expected: "",
	}
	tests.RunTest(t, setConfigTestCase)

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

	getConfigTestCase := tests.Testcase{
		Name:     "Get Updated Cellular Gateway Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:  "",
		Expected: configString,
	}
	tests.RunTest(t, getConfigTestCase)

	// Set new configs and remove enodeb associations
	gatewayConfig.AttachedEnodebSerials = []string{}
	marshaledCfg, err = gatewayConfig.MarshalBinary()
	assert.NoError(t, err)
	configString = string(marshaledCfg)

	setConfigTestCase = tests.Testcase{
		Name:     "Set Cellular Gateway Config And Remove Enodeb Associations",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:  configString,
		Expected: "",
	}
	tests.RunTest(t, setConfigTestCase)

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

	setConfigTestCase = tests.Testcase{
		Name:                     "Set Invalid Cellular Gateway Config",
		Method:                   "PUT",
		Url:                      fmt.Sprintf("%s/%s/gateways/%s/configs/cellular", testUrlRoot, networkID, gatewayID),
		Payload:                  configString,
		Expected:                 `{"message":"Invalid config: Gateway RAN config is nil"}`,
		Expect_http_error_status: true,
	}
	status, _, err := tests.RunTest(t, setConfigTestCase)
	assert.Equal(t, 400, status)
}

func testSetNetworkConfigs(t *testing.T, config *models2.NetworkCellularConfigs, expectedConfig *models2.NetworkCellularConfigs) {
	restPort := tests.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)
	networkID := "cellular_obsidian_test_network"
	configuratorTestUtils.RegisterNetwork(t, networkID, "Test Network 1")

	// Happy path
	marshaledCfg, err := config.MarshalBinary()
	assert.NoError(t, err)
	configString := string(marshaledCfg)

	createConfigTestCase := tests.Testcase{
		Name:     "Create Cellular Network Config",
		Method:   "POST",
		Url:      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:  configString,
		Expected: fmt.Sprintf(`"%s"`, networkID),
	}
	tests.RunTest(t, createConfigTestCase)

	config.Epc.Mcc = "123"
	config.Epc.SubProfiles = make(map[string]models2.NetworkEpcConfigsSubProfilesAnon)
	config.Epc.SubProfiles["test"] =
		models2.NetworkEpcConfigsSubProfilesAnon{
			MaxUlBitRate: 100, MaxDlBitRate: 200,
		}
	config.Ran.BandwidthMhz = 15

	expectedConfig.Epc.Mcc = "123"
	expectedConfig.Epc.SubProfiles = make(map[string]models2.NetworkEpcConfigsSubProfilesAnon)
	expectedConfig.Epc.SubProfiles["test"] =
		models2.NetworkEpcConfigsSubProfilesAnon{
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

	setConfigTestCase := tests.Testcase{
		Name:     "Set Cellular Network Config",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:  configString,
		Expected: "",
	}
	tests.RunTest(t, setConfigTestCase)
	getConfigTestCase := tests.Testcase{
		Name:     "Get Updated Cellular Network Config",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s/configs/cellular", testUrlRoot, networkID),
		Payload:  "",
		Expected: exSwaggerConfigString,
	}
	tests.RunTest(t, getConfigTestCase)
}

func registerEnodeb(t *testing.T, networkID, enodebID string) {
	enodebEntity := configurator.NetworkEntity{
		Key:  enodebID,
		Type: lte.CellularEnodebType,
	}
	_, err := configurator.CreateEntity(networkID, enodebEntity)
	assert.NoError(t, err)
}
