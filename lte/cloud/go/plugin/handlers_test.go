/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin_test

import (
	"crypto/x509"
	"fmt"
	"testing"
	"time"

	"magma/lte/cloud/go/lte"
	plugin2 "magma/lte/cloud/go/plugin"
	models2 "magma/lte/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/security/key"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/device"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestListNetworks(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := plugin2.GetHandlers()
	listNetworks := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/ltenetworks", obsidian.GET).HandlerFunc

	// Test empty response
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/ltenetworks",
		Payload:        nil,
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{}),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	seedNetworks(t)

	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/ltenetworks",
		Payload:        nil,
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n1", "n3"}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestCreateNetwork(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := plugin2.GetHandlers()
	createNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/ltenetworks", obsidian.POST).HandlerFunc

	// test validation
	tc := tests.Test{
		Method: "POST",
		URL:    "/magma/v1/ltenetworks",
		Payload: tests.JSONMarshaler(
			&models2.LteNetwork{
				Cellular:    models2.NewDefaultTDDNetworkConfig(),
				Description: "",
				DNS:         models.NewDefaultDNSConfig(),
				Features:    models.NewDefaultFeaturesConfig(),
				ID:          "n1",
				Name:        "foobar",
			},
		),
		Handler:        createNetwork,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"description in body should be at least 1 chars long",
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method: "POST",
		URL:    "/magma/v1/ltenetworks",
		Payload: tests.JSONMarshaler(
			&models2.LteNetwork{
				Cellular:    models2.NewDefaultTDDNetworkConfig(),
				Description: "Foo Bar",
				DNS:         models.NewDefaultDNSConfig(),
				Features:    models.NewDefaultFeaturesConfig(),
				ID:          "n1",
				Name:        "foobar",
			},
		),
		Handler:        createNetwork,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadNetwork("n1", true, true)
	assert.NoError(t, err)
	expected := configurator.Network{
		ID:          "n1",
		Type:        lte.LteNetworkType,
		Name:        "foobar",
		Description: "Foo Bar",
		Configs: map[string]interface{}{
			lte.CellularNetworkType:     models2.NewDefaultTDDNetworkConfig(),
			orc8r.DnsdNetworkType:       models.NewDefaultDNSConfig(),
			orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
		},
	}
	assert.Equal(t, expected, actual)
}

func TestGetNetwork(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := plugin2.GetHandlers()
	getNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/ltenetworks/:network_id", obsidian.GET).HandlerFunc

	// Test 404
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/ltenetworks/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetwork,
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	seedNetworks(t)

	expectedN1 := &models2.LteNetwork{
		Cellular:    models2.NewDefaultTDDNetworkConfig(),
		Description: "Foo Bar",
		DNS:         models.NewDefaultDNSConfig(),
		Features:    models.NewDefaultFeaturesConfig(),
		ID:          "n1",
		Name:        "foobar",
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/ltenetworks/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedN1),
	}
	tests.RunUnitTest(t, e, tc)

	// get a non-LTE network
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/ltenetworks/n2",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        getNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not an LTE network",
	}
	tests.RunUnitTest(t, e, tc)

	// get a network without any configs (poorly formed data)
	expectedN3 := &models2.LteNetwork{
		Description: "Bar Foo",
		ID:          "n3",
		Name:        "barfoo",
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/ltenetworks/n3",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n3"},
		Handler:        getNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedN3),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateNetwork(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := plugin2.GetHandlers()
	updateNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/ltenetworks/:network_id", obsidian.PUT).HandlerFunc

	// Test validation failure
	payloadN1 := &models2.LteNetwork{
		ID:          "n1",
		Name:        "updated foobar",
		Description: "Updated Foo Bar",
		Cellular:    models2.NewDefaultFDDNetworkConfig(),
		Features: &models.NetworkFeatures{
			Features: map[string]string{
				"bar": "baz",
				"baz": "quz",
			},
		},
		DNS: &models.NetworkDNSConfig{
			EnableCaching: swag.Bool(true),
			LocalTTL:      swag.Uint32(120),
			Records: []*models.DNSConfigRecord{
				{
					Domain:     "foobar.com",
					ARecord:    []strfmt.IPv4{"asdf", "hjkl"},
					AaaaRecord: []strfmt.IPv6{"abcd", "efgh"},
				},
				{
					Domain:  "facebook.com",
					ARecord: []strfmt.IPv4{"google.com"},
				},
			},
		},
	}
	tc := tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/ltenetworks/n1",
		Payload:        payloadN1,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateNetwork,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"validation failure list:\n" +
			"validation failure list:\n" +
			"a_record.0 in body must be of type ipv4: \"asdf\"\n" +
			"aaaa_record.0 in body must be of type ipv6: \"abcd\"",
	}
	tests.RunUnitTest(t, e, tc)

	payloadN1.DNS.Records = []*models.DNSConfigRecord{
		{
			Domain:  "foobar.com",
			ARecord: []strfmt.IPv4{"127.0.0.1", "127.0.0.2"},
			AaaaRecord: []strfmt.IPv6{
				"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
				"1234:0db8:85a3:0000:0000:8a2e:0370:1234",
			},
		},
		{
			Domain:  "facebook.com",
			ARecord: []strfmt.IPv4{"127.0.0.3"},
		},
	}
	// Test 404
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/ltenetworks/n1",
		Payload:        payloadN1,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateNetwork,
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	// seed networks, update n1 again
	seedNetworks(t)

	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/ltenetworks/n1",
		Payload:        payloadN1,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateNetwork,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualN1, err := configurator.LoadNetwork("n1", true, true)
	assert.NoError(t, err)
	expected := configurator.Network{
		ID:          "n1",
		Type:        lte.LteNetworkType,
		Name:        "updated foobar",
		Description: "Updated Foo Bar",
		Configs: map[string]interface{}{
			lte.CellularNetworkType:     models2.NewDefaultFDDNetworkConfig(),
			orc8r.DnsdNetworkType:       payloadN1.DNS,
			orc8r.NetworkFeaturesConfig: payloadN1.Features,
		},
		Version: 1,
	}
	assert.Equal(t, expected, actualN1)

	// update n2, should be 400
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/ltenetworks/n2",
		Payload:        payloadN1,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        updateNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not an LTE network",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestDeleteNetwork(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := plugin2.GetHandlers()
	deleteNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/ltenetworks/:network_id", obsidian.DELETE).HandlerFunc

	// Test 404
	tc := tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/ltenetworks/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        deleteNetwork,
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	// seed networks, delete n1 again
	seedNetworks(t)
	tc.ExpectedStatus = 204
	tests.RunUnitTest(t, e, tc)

	// delete n1 again, should be 404
	tc.ExpectedStatus = 404
	tests.RunUnitTest(t, e, tc)

	// try to delete n2, should be 400 (not LTE network)
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/ltenetworks/n2",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        deleteNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not an LTE network",
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.ListNetworkIDs()
	assert.NoError(t, err)
	assert.Equal(t, []string{"n2", "n3"}, actual)
}

func TestCellularPartialGet(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks"

	seedNetworks(t)

	handlers := plugin2.GetHandlers()
	getCellular := tests.GetHandlerByPathAndMethod(t, handlers,
		fmt.Sprintf("%s/:network_id/cellular", testURLRoot), obsidian.GET).HandlerFunc
	getEpc := tests.GetHandlerByPathAndMethod(t, handlers,
		fmt.Sprintf("%s/:network_id/cellular/epc", testURLRoot), obsidian.GET).HandlerFunc
	getRan := tests.GetHandlerByPathAndMethod(t, handlers,
		fmt.Sprintf("%s/:network_id/cellular/ran", testURLRoot), obsidian.GET).HandlerFunc
	getFegNetworkID := tests.GetHandlerByPathAndMethod(t, handlers,
		fmt.Sprintf("%s/:network_id/cellular/feg_network_id", testURLRoot), obsidian.GET).HandlerFunc

	// happy path
	tc := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/cellular/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getCellular,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models2.NewDefaultTDDNetworkConfig()),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// 404
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/cellular/", testURLRoot, "n2"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        getCellular,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// happy path
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/cellular/epc/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getEpc,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models2.NewDefaultTDDNetworkConfig().Epc),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// 404
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/cellular/epc/", testURLRoot, "n2"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        getEpc,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// happy path
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/cellular/ran/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getRan,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models2.NewDefaultTDDNetworkConfig().Ran),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// 404
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/cellular/ran/", testURLRoot, "n2"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        getRan,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// add 'n2' as FegNetworkID to n1
	cellularConfig := models2.NewDefaultTDDNetworkConfig()
	cellularConfig.FegNetworkID = "n2"
	err := configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{
		{
			ID: "n1",
			ConfigsToAddOrUpdate: map[string]interface{}{
				lte.CellularNetworkType: cellularConfig,
			},
		},
	})
	assert.NoError(t, err)

	// happy case FegNetworkID from cellular config
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/cellular/feg_network_id/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getFegNetworkID,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler("n2"),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestCellularPartialUpdate(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks"

	seedNetworks(t)
	handlers := plugin2.GetHandlers()
	updateCellular := tests.GetHandlerByPathAndMethod(t, handlers,
		fmt.Sprintf("%s/:network_id/cellular", testURLRoot), obsidian.PUT).HandlerFunc
	updateEpc := tests.GetHandlerByPathAndMethod(t, handlers,
		fmt.Sprintf("%s/:network_id/cellular/epc", testURLRoot), obsidian.PUT).HandlerFunc
	updateRan := tests.GetHandlerByPathAndMethod(t, handlers,
		fmt.Sprintf("%s/:network_id/cellular/ran", testURLRoot), obsidian.PUT).HandlerFunc
	updateFegNetworkID := tests.GetHandlerByPathAndMethod(t, handlers,
		fmt.Sprintf("%s/:network_id/cellular/feg_network_id", testURLRoot), obsidian.PUT).HandlerFunc

	// happy path update cellular config
	tc := tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/cellular/", testURLRoot, "n2"),
		Payload:        models2.NewDefaultFDDNetworkConfig(),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        updateCellular,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualN2, err := configurator.LoadNetwork("n2", true, true)
	assert.NoError(t, err)
	expected := configurator.Network{
		ID:          "n2",
		Type:        "blah",
		Name:        "foobar",
		Description: "Foo Bar",
		Configs: map[string]interface{}{
			lte.CellularNetworkType: models2.NewDefaultFDDNetworkConfig(),
		},
		Version: 1,
	}
	assert.Equal(t, expected, actualN2)

	// Validation error (celullar config has both tdd and fdd config)
	badCellularConfig := models2.NewDefaultTDDNetworkConfig()
	badCellularConfig.Ran.FddConfig = &models2.NetworkRanConfigsFddConfig{
		Earfcndl: 1,
		Earfcnul: 18001,
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/cellular/", testURLRoot, "n2"),
		Payload:        badCellularConfig,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        updateCellular,
		ExpectedStatus: 400,
		ExpectedError:  "only one of TDD or FDD configs can be set",
	}
	tests.RunUnitTest(t, e, tc)

	// Fail to put epc config to a network without cellular network configs
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/cellular/epc/", testURLRoot, "n3"),
		Payload:        models2.NewDefaultTDDNetworkConfig().Epc,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n3"},
		Handler:        updateEpc,
		ExpectedStatus: 400,
		ExpectedError:  "No cellular network config found",
	}
	tests.RunUnitTest(t, e, tc)

	// happy path update epc config
	epcConfig := models2.NewDefaultTDDNetworkConfig().Epc
	epcConfig.RelayEnabled = swag.Bool(true)
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/cellular/epc/", testURLRoot, "n2"),
		Payload:        epcConfig,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        updateEpc,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualN2, err = configurator.LoadNetwork("n2", true, true)
	assert.NoError(t, err)
	expected.Configs[lte.CellularNetworkType].(*models2.NetworkCellularConfigs).Epc = epcConfig
	expected.Version = 2
	assert.Equal(t, expected, actualN2)

	// Fail to put epc config to a network without cellular network configs
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/cellular/ran/", testURLRoot, "n3"),
		Payload:        models2.NewDefaultTDDNetworkConfig().Ran,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n3"},
		Handler:        updateRan,
		ExpectedStatus: 400,
		ExpectedError:  "No cellular network config found",
	}
	tests.RunUnitTest(t, e, tc)

	// Validation error
	ranConfig := models2.NewDefaultTDDNetworkConfig().Ran
	ranConfig.FddConfig = models2.NewDefaultFDDNetworkConfig().Ran.FddConfig
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/cellular/ran/", testURLRoot, "n2"),
		Payload:        ranConfig,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        updateRan,
		ExpectedStatus: 400,
		ExpectedError:  "only one of TDD or FDD configs can be set",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case update ran config
	ranConfig = models2.NewDefaultFDDNetworkConfig().Ran
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/cellular/ran/", testURLRoot, "n2"),
		Payload:        ranConfig,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        updateRan,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	actualN2, err = configurator.LoadNetwork("n2", true, true)
	assert.NoError(t, err)
	expected.Configs[lte.CellularNetworkType].(*models2.NetworkCellularConfigs).Ran = ranConfig
	expected.Version = 3
	assert.Equal(t, expected, actualN2)

	// Validation Error (should not be able to add nonexistent networkID as fegNetworkID)
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/cellular/feg_network_id/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler("bad-network-id"),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateFegNetworkID,
		ExpectedStatus: 400,
		ExpectedError:  "Network: bad-network-id does not exist",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/cellular/feg_network_id/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler("n2"),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateFegNetworkID,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
}

func TestCellularDelete(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks"

	seedNetworks(t)

	handlers := plugin2.GetHandlers()
	deleteCellular := tests.GetHandlerByPathAndMethod(t, handlers,
		fmt.Sprintf("%s/:network_id/cellular", testURLRoot), obsidian.DELETE).HandlerFunc

	tc := tests.Test{
		Method:         "DELETE",
		URL:            fmt.Sprintf("%s/%s/cellular/", testURLRoot, "n1"),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        deleteCellular,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	_, err := configurator.LoadNetworkConfig("n1", lte.CellularNetworkType)
	assert.EqualError(t, err, "Not found")
}

func TestListAndGetGateways(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.GetUnfreezeClockDeferFunc(t)()

	test_init.StartTestService(t)
	stateTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/gateways"

	handlers := plugin2.GetHandlers()
	listGateways := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc
	getGateway := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/:gateway_id", testURLRoot), obsidian.GET).HandlerFunc

	// Create 2 gateways, 1 with state and device, the other without
	// g2 will associate to 2 enodebs
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.CellularEnodebType, Key: "enb1"},
			{Type: lte.CellularEnodebType, Key: "enb2"},
			{
				Type: lte.CellularGatewayType, Key: "g1",
				Config: &models2.GatewayCellularConfigs{
					Epc: &models2.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
					Ran: &models2.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
				},
			},
			{
				Type: lte.CellularGatewayType, Key: "g2",
				Config: &models2.GatewayCellularConfigs{
					Epc: &models2.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
					Ran: &models2.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
				},
				Associations: []storage.TypeAndKey{
					{Type: lte.CellularEnodebType, Key: "enb1"},
					{Type: lte.CellularEnodebType, Key: "enb2"},
				},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Name: "foobar", Description: "foo bar",
				PhysicalID: "hw1",
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: "g1"}},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: "g2",
				Name: "barfoo", Description: "bar foo",
				PhysicalID: "hw2",
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: "g2"}},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: "g1"},
					{Type: orc8r.MagmadGatewayType, Key: "g2"},
				},
			},
		},
	)
	assert.NoError(t, err)
	err = device.RegisterDevice("n1", orc8r.AccessGatewayRecordType, "hw1", &models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}})
	assert.NoError(t, err)
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportGatewayStatus(t, ctx, models.NewDefaultGatewayStatus("hw1"))

	expected := map[string]*models2.LteGateway{
		"g1": {
			ID: "g1",
			Device: &models.GatewayDevice{
				HardwareID: "hw1",
				Key:        &models.ChallengeKey{KeyType: "ECHO"},
			},
			Name: "foobar", Description: "foo bar",
			Tier: "t1",
			Magmad: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Cellular: &models2.GatewayCellularConfigs{
				Epc: &models2.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
				Ran: &models2.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
			},
			Status: models.NewDefaultGatewayStatus("hw1"),
		},
		"g2": {
			ID:   "g2",
			Name: "barfoo", Description: "bar foo",
			Tier: "t1",
			Magmad: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Cellular: &models2.GatewayCellularConfigs{
				Epc: &models2.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
				Ran: &models2.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
			},
			ConnectedEnodebSerials: []string{"enb1", "enb2"},
		},
	}
	expected["g1"].Status.CheckinTime = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
	expected["g1"].Status.CertExpirationTime = time.Unix(1000000, 0).Add(time.Hour * 4).Unix()

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listGateways,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expected),
	}
	tests.RunUnitTest(t, e, tc)

	expectedGet := &models2.LteGateway{
		ID: "g1",
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		Name: "foobar", Description: "foo bar",
		Tier: "t1",
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		Cellular: &models2.GatewayCellularConfigs{
			Epc: &models2.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
			Ran: &models2.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
		},
		Status: models.NewDefaultGatewayStatus("hw1"),
	}
	expectedGet.Status.CheckinTime = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
	expectedGet.Status.CertExpirationTime = time.Unix(1000000, 0).Add(time.Hour * 4).Unix()
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: expectedGet,
	}
	tests.RunUnitTest(t, e, tc)

	expectedGet = &models2.LteGateway{
		ID:   "g2",
		Name: "barfoo", Description: "bar foo",
		Tier: "t1",
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		Cellular: &models2.GatewayCellularConfigs{
			Epc: &models2.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
			Ran: &models2.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
		},
		ConnectedEnodebSerials: []string{"enb1", "enb2"},
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g2"},
		ExpectedStatus: 200,
		ExpectedResult: expectedGet,
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateGateway(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.GetUnfreezeClockDeferFunc(t)()

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/gateways/:gateway_id"
	handlers := plugin2.GetHandlers()
	updateGateway := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.PUT).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.CellularEnodebType, Key: "enb1"},
			{Type: lte.CellularEnodebType, Key: "enb2"},
			{Type: lte.CellularEnodebType, Key: "enb3"},
			{
				Type: lte.CellularGatewayType, Key: "g1",
				Config: &models2.GatewayCellularConfigs{
					Epc: &models2.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
					Ran: &models2.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
				},
				Associations: []storage.TypeAndKey{
					{Type: lte.CellularEnodebType, Key: "enb1"},
					{Type: lte.CellularEnodebType, Key: "enb2"},
				},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Name: "foobar", Description: "foo bar",
				PhysicalID: "hw1",
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: "g1"}},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: "g1"},
				},
			},
		},
	)
	assert.NoError(t, err)
	err = device.RegisterDevice("n1", orc8r.AccessGatewayRecordType, "hw1", &models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}})
	assert.NoError(t, err)

	// update everything
	privateKey, err := key.GenerateKey("P256", 0)
	assert.NoError(t, err)
	marshaledPubKey, err := x509.MarshalPKIXPublicKey(key.PublicKey(privateKey))
	assert.NoError(t, err)
	pubkeyB64 := strfmt.Base64(marshaledPubKey)
	payload := &models2.MutableLteGateway{
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "SOFTWARE_ECDSA_SHA256", Key: &pubkeyB64},
		},
		ID:          "g1",
		Name:        "barbaz",
		Description: "bar baz",
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         25,
			CheckinTimeout:          15,
			AutoupgradePollInterval: 200,
			AutoupgradeEnabled:      swag.Bool(false),
			FeatureFlags:            map[string]bool{"foo": false},
			DynamicServices:         []string{"d1", "d2"},
		},
		Tier: "t1",
		Cellular: &models2.GatewayCellularConfigs{
			Epc: &models2.GatewayEpcConfigs{NatEnabled: swag.Bool(false), IPBlock: "172.10.10.0/24"},
			Ran: &models2.GatewayRanConfigs{Pci: 123, TransmitEnabled: swag.Bool(false)},
		},
		ConnectedEnodebSerials: []string{"enb1", "enb3"},
	}

	tc := tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g1"},
			{Type: lte.CellularGatewayType, Key: "g1"},
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
		},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	actualDevice, err := device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hw1")
	assert.NoError(t, err)

	expectedEnts := configurator.NetworkEntities{
		{
			NetworkID: "n1", Type: lte.CellularGatewayType, Key: "g1",
			Name: string(payload.Name), Description: string(payload.Description),
			Config:             payload.Cellular,
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
			Associations: []storage.TypeAndKey{
				{Type: lte.CellularEnodebType, Key: "enb1"},
				{Type: lte.CellularEnodebType, Key: "enb3"},
			},
			GraphID: "10",
			Version: 1,
		},
		{
			NetworkID: "n1", Type: orc8r.MagmadGatewayType, Key: "g1",
			Name: string(payload.Name), Description: string(payload.Description),
			PhysicalID:         "hw1",
			Config:             payload.Magmad,
			Associations:       []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: "g1"}},
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			GraphID:            "10",
			Version:            1,
		},
		{
			NetworkID: "n1", Type: orc8r.UpgradeTierEntityType, Key: "t1",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
			GraphID:      "10",
		},
	}
	assert.Equal(t, expectedEnts, actualEnts)
	assert.Equal(t, payload.Device, actualDevice)
}

func TestDeleteGateway(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.GetUnfreezeClockDeferFunc(t)()

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/gateways/:gateway_id"
	handlers := plugin2.GetHandlers()
	deleteGateway := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.DELETE).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.CellularEnodebType, Key: "enb1"},
			{Type: lte.CellularEnodebType, Key: "enb2"},
			{
				Type: lte.CellularGatewayType, Key: "g1",
				Config: &models2.GatewayCellularConfigs{
					Epc: &models2.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
					Ran: &models2.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
				},
				Associations: []storage.TypeAndKey{
					{Type: lte.CellularEnodebType, Key: "enb1"},
					{Type: lte.CellularEnodebType, Key: "enb2"},
				},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Name: "foobar", Description: "foo bar",
				PhysicalID: "hw1",
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: "g1"}},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: "g1"},
				},
			},
		},
	)
	assert.NoError(t, err)
	err = device.RegisterDevice("n1", orc8r.AccessGatewayRecordType, "hw1", &models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}})
	assert.NoError(t, err)

	tc := tests.Test{
		Method:         "DELETE",
		URL:            testURLRoot,
		Handler:        deleteGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g1"},
			{Type: lte.CellularGatewayType, Key: "g1"},
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
		},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	actualDevice, err := device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hw1")
	assert.NoError(t, err)

	expectedEnts := configurator.NetworkEntities{
		{NetworkID: "n1", Type: orc8r.UpgradeTierEntityType, Key: "t1", GraphID: "11"},
	}
	assert.Equal(t, expectedEnts, actualEnts)
	assert.Equal(t, &models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}}, actualDevice)
}

func TestGetCellularGatewayConfig(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/gateways/:gateway_id"
	handlers := plugin2.GetHandlers()
	getCellular := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular", testURLRoot), obsidian.GET).HandlerFunc
	getEpc := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/epc", testURLRoot), obsidian.GET).HandlerFunc
	getRan := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/ran", testURLRoot), obsidian.GET).HandlerFunc
	getNonEps := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/non_eps", testURLRoot), obsidian.GET).HandlerFunc
	getEnodebs := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/connected_enodeb_serial", testURLRoot), obsidian.GET).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.CellularEnodebType, Key: "enb1"},
			{Type: lte.CellularEnodebType, Key: "enb2"},
			{
				Type: lte.CellularGatewayType, Key: "g1",
				Config: newDefaultGatewayConfig(),
				Associations: []storage.TypeAndKey{
					{Type: lte.CellularEnodebType, Key: "enb1"},
					{Type: lte.CellularEnodebType, Key: "enb2"},
				},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Name: "foobar", Description: "foo bar",
				PhysicalID:   "hw1",
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: "g1"}},
			},
		},
	)
	assert.NoError(t, err)

	// 404
	tc := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/cellular", testURLRoot),
		Handler:        getCellular,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g2"},
		ExpectedResult: newDefaultGatewayConfig(),
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/cellular", testURLRoot),
		Handler:        getCellular,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedResult: newDefaultGatewayConfig(),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/cellular/epc", testURLRoot),
		Handler:        getEpc,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedResult: newDefaultGatewayConfig().Epc,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/cellular/ran", testURLRoot),
		Handler:        getRan,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedResult: newDefaultGatewayConfig().Ran,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/cellular/non_eps", testURLRoot),
		Handler:        getNonEps,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedResult: newDefaultGatewayConfig().NonEpsService,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/cellular/connected_enodeb_serial", testURLRoot),
		Handler:        getEnodebs,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedResult: tests.JSONMarshaler([]string{"enb1", "enb2"}),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateCellularGatewayConfig(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/gateways/:gateway_id"
	handlers := plugin2.GetHandlers()
	updateCellular := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular", testURLRoot), obsidian.PUT).HandlerFunc
	updateEpc := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/epc", testURLRoot), obsidian.PUT).HandlerFunc
	updateRan := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/ran", testURLRoot), obsidian.PUT).HandlerFunc
	updateNonEps := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/non_eps", testURLRoot), obsidian.PUT).HandlerFunc
	updateEnodebs := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/connected_enodeb_serial", testURLRoot), obsidian.PUT).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.CellularEnodebType, Key: "enb1"},
			{Type: lte.CellularEnodebType, Key: "enb2"},
			{Type: lte.CellularGatewayType, Key: "g1"},
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: "g1"}},
			},
		},
	)
	assert.NoError(t, err)

	tc := tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/cellular", testURLRoot),
		Handler:        updateCellular,
		Payload:        newDefaultGatewayConfig(),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expected := map[storage.TypeAndKey]configurator.NetworkEntity{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: "g1"}},
			GraphID:      "6",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayType, Key: "g1",
			Config:  newDefaultGatewayConfig(),
			GraphID: "6",
			Version: 1,
		},
	}

	entities, _, err := configurator.LoadEntities("n1", nil, swag.String("g1"), nil, nil, configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true})
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.ToEntitiesByID())

	modifiedCellularConfig := newDefaultGatewayConfig()
	modifiedCellularConfig.Epc.NatEnabled = swag.Bool(false)
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/cellular/epc", testURLRoot),
		Handler:        updateEpc,
		Payload:        modifiedCellularConfig.Epc,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expected = map[storage.TypeAndKey]configurator.NetworkEntity{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: "g1"}},
			GraphID:      "6",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayType, Key: "g1",
			Config:  modifiedCellularConfig,
			GraphID: "6",
			Version: 2,
		},
	}
	entities, _, err = configurator.LoadEntities("n1", nil, swag.String("g1"), nil, nil, configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true})
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.ToEntitiesByID())

	modifiedCellularConfig.Ran.TransmitEnabled = swag.Bool(false)
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/cellular/ran", testURLRoot),
		Handler:        updateRan,
		Payload:        modifiedCellularConfig.Ran,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expected = map[storage.TypeAndKey]configurator.NetworkEntity{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: "g1"}},
			GraphID:      "6",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayType, Key: "g1",
			Config:  modifiedCellularConfig,
			GraphID: "6",
			Version: 3,
		},
	}
	entities, _, err = configurator.LoadEntities("n1", nil, swag.String("g1"), nil, nil, configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true})
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.ToEntitiesByID())

	// validation failure
	modifiedCellularConfig.NonEpsService.CsfbMcc = "0"
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/cellular/ran", testURLRoot),
		Handler:        updateNonEps,
		Payload:        modifiedCellularConfig.NonEpsService,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 400,
		ExpectedError:  "validation failure list:\ncsfb_mcc in body should match '^(\\d{3})$'",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case
	modifiedCellularConfig.NonEpsService.CsfbMcc = "123"
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/cellular/ran", testURLRoot),
		Handler:        updateNonEps,
		Payload:        modifiedCellularConfig.NonEpsService,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expected = map[storage.TypeAndKey]configurator.NetworkEntity{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: "g1"}},
			GraphID:      "6",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayType, Key: "g1",
			Config:  modifiedCellularConfig,
			GraphID: "6",
			Version: 4,
		},
	}
	entities, _, err = configurator.LoadEntities("n1", nil, swag.String("g1"), nil, nil, configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true})
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.ToEntitiesByID())

	// happy case
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/cellular/connected_enodeb_serial", testURLRoot),
		Handler:        updateEnodebs,
		Payload:        tests.JSONMarshaler([]string{"enb1", "enb2"}),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expected = map[storage.TypeAndKey]configurator.NetworkEntity{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: "g1"}},
			GraphID:      "2",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayType, Key: "g1",
			Config:  modifiedCellularConfig,
			GraphID: "2",
			Version: 5,
			Associations: []storage.TypeAndKey{
				{Type: lte.CellularEnodebType, Key: "enb1"},
				{Type: lte.CellularEnodebType, Key: "enb2"},
			},
		},
	}
	entities, _, err = configurator.LoadEntities("n1", nil, swag.String("g1"), nil, nil, configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true})
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.ToEntitiesByID())
}

func TestListAndGetEnodebs(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/enodebs"

	handlers := plugin2.GetHandlers()
	listEnodebs := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc
	getEnodeb := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/:enodeb_serial", testURLRoot), obsidian.GET).HandlerFunc

	_, err = configurator.CreateEntities("n1", []configurator.NetworkEntity{
		{
			Type:       lte.CellularEnodebType,
			Key:        "abcdefg",
			Name:       "abc enodeb",
			PhysicalID: "abcdefg",
			Config: &models2.EnodebConfiguration{
				BandwidthMhz:           20,
				CellID:                 swag.Uint32(1234),
				DeviceClass:            "Baicells Nova-233 G2 OD FDD",
				Earfcndl:               39450,
				Pci:                    260,
				SpecialSubframePattern: 7,
				SubframeAssignment:     2,
				Tac:                    1,
				TransmitEnabled:        swag.Bool(true),
			},
		},
		{
			Type:       lte.CellularEnodebType,
			Key:        "vwxyz",
			Name:       "xyz enodeb",
			PhysicalID: "vwxyz",
			Config: &models2.EnodebConfiguration{
				BandwidthMhz:           15,
				CellID:                 swag.Uint32(4321),
				DeviceClass:            "Baicells Nova-243 OD TDD",
				Earfcndl:               39550,
				Pci:                    261,
				SpecialSubframePattern: 8,
				SubframeAssignment:     3,
				Tac:                    2,
				TransmitEnabled:        swag.Bool(false),
			},
		},
		{
			Type: lte.CellularGatewayType, Key: "gw1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularEnodebType, Key: "abcdefg"}},
		},
	})
	assert.NoError(t, err)

	expected := map[string]*models2.Enodeb{
		"abcdefg": {
			AttachedGatewayID: "gw1",
			Config: &models2.EnodebConfiguration{
				BandwidthMhz:           20,
				CellID:                 swag.Uint32(1234),
				DeviceClass:            "Baicells Nova-233 G2 OD FDD",
				Earfcndl:               39450,
				Pci:                    260,
				SpecialSubframePattern: 7,
				SubframeAssignment:     2,
				Tac:                    1,
				TransmitEnabled:        swag.Bool(true),
			},
			Name:   "abc enodeb",
			Serial: "abcdefg",
		},
		"vwxyz": {
			Config: &models2.EnodebConfiguration{
				BandwidthMhz:           15,
				CellID:                 swag.Uint32(4321),
				DeviceClass:            "Baicells Nova-243 OD TDD",
				Earfcndl:               39550,
				Pci:                    261,
				SpecialSubframePattern: 8,
				SubframeAssignment:     3,
				Tac:                    2,
				TransmitEnabled:        swag.Bool(false),
			},
			Name:   "xyz enodeb",
			Serial: "vwxyz",
		},
	}
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listEnodebs,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expected),
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getEnodeb,
		ParamNames:     []string{"network_id", "enodeb_serial"},
		ParamValues:    []string{"n1", "abcdefg"},
		ExpectedStatus: 200,
		ExpectedResult: expected["abcdefg"],
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getEnodeb,
		ParamNames:     []string{"network_id", "enodeb_serial"},
		ParamValues:    []string{"n1", "vwxyz"},
		ExpectedStatus: 200,
		ExpectedResult: expected["vwxyz"],
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getEnodeb,
		ParamNames:     []string{"network_id", "enodeb_serial"},
		ParamValues:    []string{"n1", "hello"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestCreateEnodeb(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/enodebs"

	handlers := plugin2.GetHandlers()
	createEnodeb := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.POST).HandlerFunc

	tc := tests.Test{
		Method:  "POST",
		URL:     testURLRoot,
		Handler: createEnodeb,
		Payload: &models2.Enodeb{
			Config: &models2.EnodebConfiguration{
				BandwidthMhz:           15,
				CellID:                 swag.Uint32(4321),
				DeviceClass:            "Baicells Nova-243 OD TDD",
				Earfcndl:               39550,
				Pci:                    261,
				SpecialSubframePattern: 8,
				SubframeAssignment:     3,
				Tac:                    2,
				TransmitEnabled:        swag.Bool(false),
			},
			Name:   "foobar",
			Serial: "abcdef",
		},
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.CellularEnodebType, "abcdef", configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.CellularEnodebType, Key: "abcdef",
		Name:       "foobar",
		PhysicalID: "abcdef",
		GraphID:    "2",
		Config: &models2.EnodebConfiguration{
			BandwidthMhz:           15,
			CellID:                 swag.Uint32(4321),
			DeviceClass:            "Baicells Nova-243 OD TDD",
			Earfcndl:               39550,
			Pci:                    261,
			SpecialSubframePattern: 8,
			SubframeAssignment:     3,
			Tac:                    2,
			TransmitEnabled:        swag.Bool(false),
		},
	}
	assert.Equal(t, expected, actual)

	tc = tests.Test{
		Method:  "POST",
		URL:     testURLRoot,
		Handler: createEnodeb,
		Payload: &models2.Enodeb{
			Config: &models2.EnodebConfiguration{
				BandwidthMhz:           15,
				CellID:                 swag.Uint32(4321),
				DeviceClass:            "Baicells Nova-243 OD TDD",
				Earfcndl:               39550,
				Pci:                    261,
				SpecialSubframePattern: 8,
				SubframeAssignment:     3,
				Tac:                    2,
				TransmitEnabled:        swag.Bool(false),
			},
			Name:              "foobar",
			Serial:            "abcdef",
			AttachedGatewayID: "gw1",
		},
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 400,
		ExpectedError:  "attached_gateway_id is a read-only property",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateEnodeb(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/enodebs/:enodeb_serial"

	handlers := plugin2.GetHandlers()
	updateEnodeb := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.PUT).HandlerFunc

	_, err = configurator.CreateEntities("n1", []configurator.NetworkEntity{
		{
			Type:       lte.CellularEnodebType,
			Key:        "abcdefg",
			Name:       "abc enodeb",
			PhysicalID: "abcdefg",
			Config: &models2.EnodebConfiguration{
				BandwidthMhz:           20,
				CellID:                 swag.Uint32(1234),
				DeviceClass:            "Baicells Nova-233 G2 OD FDD",
				Earfcndl:               39450,
				Pci:                    260,
				SpecialSubframePattern: 7,
				SubframeAssignment:     2,
				Tac:                    1,
				TransmitEnabled:        swag.Bool(true),
			},
		},
	})
	assert.NoError(t, err)

	tc := tests.Test{
		Method:  "PUT",
		URL:     testURLRoot,
		Handler: updateEnodeb,
		Payload: &models2.Enodeb{
			Config: &models2.EnodebConfiguration{
				BandwidthMhz:           15,
				CellID:                 swag.Uint32(4321),
				DeviceClass:            "Baicells Nova-243 OD TDD",
				Earfcndl:               39550,
				Pci:                    261,
				SpecialSubframePattern: 8,
				SubframeAssignment:     3,
				Tac:                    2,
				TransmitEnabled:        swag.Bool(false),
			},
			Name:   "foobar",
			Serial: "abcdefg",
		},
		ParamNames:     []string{"network_id", "enodeb_serial"},
		ParamValues:    []string{"n1", "abcdefg"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.CellularEnodebType, "abcdefg", configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.CellularEnodebType, Key: "abcdefg",
		Name:       "foobar",
		PhysicalID: "abcdefg",
		GraphID:    "2",
		Config: &models2.EnodebConfiguration{
			BandwidthMhz:           15,
			CellID:                 swag.Uint32(4321),
			DeviceClass:            "Baicells Nova-243 OD TDD",
			Earfcndl:               39550,
			Pci:                    261,
			SpecialSubframePattern: 8,
			SubframeAssignment:     3,
			Tac:                    2,
			TransmitEnabled:        swag.Bool(false),
		},
		Version: 1,
	}
	assert.Equal(t, expected, actual)

	tc = tests.Test{
		Method:  "PUT",
		URL:     testURLRoot,
		Handler: updateEnodeb,
		Payload: &models2.Enodeb{
			Config: &models2.EnodebConfiguration{
				BandwidthMhz:           15,
				CellID:                 swag.Uint32(4321),
				DeviceClass:            "Baicells Nova-243 OD TDD",
				Earfcndl:               39550,
				Pci:                    261,
				SpecialSubframePattern: 8,
				SubframeAssignment:     3,
				Tac:                    2,
				TransmitEnabled:        swag.Bool(false),
			},
			Name:              "foobar",
			Serial:            "abcdef",
			AttachedGatewayID: "gw1",
		},
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 400,
		ExpectedError:  "attached_gateway_id is a read-only property",
	}
}

func TestDeleteEnodeb(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/enodebs/:enodeb_serial"

	handlers := plugin2.GetHandlers()
	deleteEnodeb := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.DELETE).HandlerFunc

	_, err = configurator.CreateEntities("n1", []configurator.NetworkEntity{
		{
			Type:       lte.CellularEnodebType,
			Key:        "abcdefg",
			Name:       "abc enodeb",
			PhysicalID: "abcdefg",
			Config: &models2.EnodebConfiguration{
				BandwidthMhz:           20,
				CellID:                 swag.Uint32(1234),
				DeviceClass:            "Baicells Nova-233 G2 OD FDD",
				Earfcndl:               39450,
				Pci:                    260,
				SpecialSubframePattern: 7,
				SubframeAssignment:     2,
				Tac:                    1,
				TransmitEnabled:        swag.Bool(true),
			},
		},
	})
	assert.NoError(t, err)

	tc := tests.Test{
		Method:         "DELETE",
		URL:            testURLRoot,
		Handler:        deleteEnodeb,
		ParamNames:     []string{"network_id", "enodeb_serial"},
		ParamValues:    []string{"n1", "abcdefg"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.LoadEntity("n1", lte.CellularEnodebType, "abcdefg", configurator.FullEntityLoadCriteria())
	assert.EqualError(t, err, "Not found")
}

func TestCreateSubscriber(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/subscribers"
	handlers := plugin2.GetHandlers()
	createSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.POST).HandlerFunc

	tc := tests.Test{
		Method: "POST",
		URL:    testURLRoot,
		Payload: &models2.Subscriber{
			ID: "IMSI1234567890",
			Lte: &models2.LteSubscription{
				AuthAlgo: "MILENAGE",
				AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:    "ACTIVE",
			},
		},
		Handler:        createSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.SubscriberEntityType,
		Key:       "IMSI1234567890",
		Config: &models2.LteSubscription{
			AuthAlgo: "MILENAGE",
			AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			State:    "ACTIVE",
		},
		GraphID: "2",
	}
	assert.Equal(t, expected, actual)

	// validation failure
	tc = tests.Test{
		Method: "POST",
		URL:    testURLRoot,
		Payload: &models2.Subscriber{
			ID: "IMSI1234567898",
			Lte: &models2.LteSubscription{
				AuthAlgo: "MILENAGE",
				AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:    "ACTIVE",
			},
		},
		Handler:        createSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 400,
		ExpectedError:  "expected lte auth key to be 16 bytes but got 15 bytes",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestListSubscribers(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/subscribers"
	handlers := plugin2.GetHandlers()
	listSubscribers := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models2.Subscriber{}),
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
				Config: &models2.LteSubscription{
					AuthAlgo: "MILENAGE",
					AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					State:    "ACTIVE",
				},
			},
			{
				Type: lte.SubscriberEntityType, Key: "IMSI0987654321",
				Config: &models2.LteSubscription{
					AuthAlgo: "MILENAGE",
					AuthKey:  []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					AuthOpc:  []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					State:    "ACTIVE",
				},
			},
		},
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models2.Subscriber{
			"IMSI1234567890": {
				ID: "IMSI1234567890",
				Lte: &models2.LteSubscription{
					AuthAlgo: "MILENAGE",
					AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					State:    "ACTIVE",
				},
			},
			"IMSI0987654321": {
				ID: "IMSI0987654321",
				Lte: &models2.LteSubscription{
					AuthAlgo: "MILENAGE",
					AuthKey:  []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					AuthOpc:  []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					State:    "ACTIVE",
				},
			},
		}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetSubscriber(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/subscribers/:subscriber_id"
	handlers := plugin2.GetHandlers()
	getSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
			Config: &models2.LteSubscription{
				AuthAlgo: "MILENAGE",
				AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:    "ACTIVE",
			},
		},
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 200,
		ExpectedResult: &models2.Subscriber{
			ID: "IMSI1234567890",
			Lte: &models2.LteSubscription{
				AuthAlgo: "MILENAGE",
				AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:    "ACTIVE",
			},
		},
	}
	tests.RunUnitTest(t, e, tc)

}

func TestUpdateSubscriber(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/subscribers/:subscriber_id"
	handlers := plugin2.GetHandlers()
	updateSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.PUT).HandlerFunc

	tc := tests.Test{
		Method:  "PUT",
		URL:     testURLRoot,
		Handler: updateSubscriber,
		Payload: &models2.Subscriber{
			ID: "IMSI1234567890",
			Lte: &models2.LteSubscription{
				AuthAlgo: "MILENAGE",
				AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:    "ACTIVE",
			},
		},
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
			Config: &models2.LteSubscription{
				AuthKey: []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc: []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:   "ACTIVE",
			},
		},
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:  "PUT",
		URL:     testURLRoot,
		Handler: updateSubscriber,
		Payload: &models2.Subscriber{
			ID: "IMSI1234567890",
			Lte: &models2.LteSubscription{
				AuthAlgo: "MILENAGE",
				AuthKey:  []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				AuthOpc:  []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				State:    "INACTIVE",
			},
		},
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.SubscriberEntityType,
		Key:       "IMSI1234567890",
		Config: &models2.LteSubscription{
			AuthAlgo: "MILENAGE",
			AuthKey:  []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			AuthOpc:  []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			State:    "INACTIVE",
		},
		GraphID: "2",
		Version: 1,
	}
	assert.Equal(t, expected, actual)
}

func TestDeleteSubscriber(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/ltenetworks/:network_id/subscribers/:subscriber_id"
	handlers := plugin2.GetHandlers()
	deleteSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.DELETE).HandlerFunc

	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
			Config: &models2.LteSubscription{
				AuthAlgo: "MILENAGE",
				AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:    "ACTIVE",
			},
		},
	)
	assert.NoError(t, err)

	tc := tests.Test{
		Method:         "DELETE",
		URL:            testURLRoot,
		Handler:        deleteSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadAllEntitiesInNetwork("n1", lte.SubscriberEntityType, configurator.EntityLoadCriteria{})
	assert.Equal(t, 0, len(actual))
}

// n1, n3 are lte networks, n2 is not
func seedNetworks(t *testing.T) {
	_, err := configurator.CreateNetworks(
		[]configurator.Network{
			{
				ID:          "n1",
				Type:        lte.LteNetworkType,
				Name:        "foobar",
				Description: "Foo Bar",
				Configs: map[string]interface{}{
					lte.CellularNetworkType:     models2.NewDefaultTDDNetworkConfig(),
					orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
					orc8r.DnsdNetworkType:       models.NewDefaultDNSConfig(),
				},
			},
			{
				ID:          "n2",
				Type:        "blah",
				Name:        "foobar",
				Description: "Foo Bar",
				Configs:     map[string]interface{}{},
			},
			{
				ID:          "n3",
				Type:        lte.LteNetworkType,
				Name:        "barfoo",
				Description: "Bar Foo",
				Configs:     map[string]interface{}{},
			},
		},
	)
	assert.NoError(t, err)
}
