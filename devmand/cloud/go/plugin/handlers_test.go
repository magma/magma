/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package plugin_test

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"orc8r/devmand/cloud/go/devmand"
	plugin2 "orc8r/devmand/cloud/go/plugin"
	models2 "orc8r/devmand/cloud/go/plugin/models"
)

func TestListNetworks(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := plugin2.GetHandlers()
	listNetworksURL := "/magma/v1/symphony"
	listNetworks := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, listNetworksURL, obsidian.GET).HandlerFunc

	// Test empty response
	tc := tests.Test{
		Method:         "GET",
		URL:            listNetworksURL,
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
		URL:            listNetworksURL,
		Payload:        nil,
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n1"}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestCreateNetwork(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := plugin2.GetHandlers()
	baseNetworksURL := "/magma/v1/symphony"
	createNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, baseNetworksURL, obsidian.POST).HandlerFunc

	// test validation
	tc := tests.Test{
		Method: "POST",
		URL:    baseNetworksURL,
		Payload: tests.JSONMarshaler(
			&models2.SymphonyNetwork{
				Description: "",
				ID:          "n1",
				Name:        "foobar",
				Features:    models.NewDefaultFeaturesConfig(),
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
		URL:    baseNetworksURL,
		Payload: tests.JSONMarshaler(
			&models2.SymphonyNetwork{
				Description: "Foo Bar",
				ID:          "n1",
				Name:        "foobar",
				Features:    models.NewDefaultFeaturesConfig(),
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
		Type:        devmand.SymphonyNetworkType,
		Name:        "foobar",
		Description: "Foo Bar",
		Configs: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
		},
	}
	assert.Equal(t, expected, actual)
}

func TestGetNetwork(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := plugin2.GetHandlers()
	getNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/symphony/:network_id", obsidian.GET).HandlerFunc

	// Test 404
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/symphony/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetwork,
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	seedNetworks(t)

	expectedN1 := &models2.SymphonyNetwork{
		Description: "Foo Bar",
		ID:          "n1",
		Name:        "foobar",
		Features:    models.NewDefaultFeaturesConfig(),
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/symphony/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedN1),
	}
	tests.RunUnitTest(t, e, tc)

	// get a non-Symphony network
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/symphony/n2",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        getNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not an Symphony network",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateNetwork(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := plugin2.GetHandlers()
	updateNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/symphony/:network_id", obsidian.PUT).HandlerFunc

	// Test validation failure
	payloadN1 := &models2.SymphonyNetwork{
		ID:          "n1",
		Name:        "updated foobar",
		Description: "Updated Foo Bar",
		Features: &models.NetworkFeatures{
			Features: map[string]string{
				"feature_1_key": "feature_1_val",
			},
		},
	}
	// Test 404
	tc := tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/symphony/n1",
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
		URL:            "/magma/v1/symphony/n1",
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
		Type:        devmand.SymphonyNetworkType,
		Name:        "updated foobar",
		Description: "Updated Foo Bar",
		Configs: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: payloadN1.Features,
		},
		Version: 1,
	}
	assert.Equal(t, expected, actualN1)

	// update n2, should be 400
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/symphony/n2",
		Payload:        payloadN1,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        updateNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not an Symphony network",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestDeleteNetwork(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := plugin2.GetHandlers()
	deleteNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/symphony/:network_id", obsidian.DELETE).HandlerFunc

	// Test 404
	tc := tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/symphony/n1",
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

	// try to delete n2, should be 400 (not Symphony network)
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/symphony/n2",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        deleteNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not an Symphony network",
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.ListNetworkIDs()
	assert.NoError(t, err)
	assert.Equal(t, []string{"n2"}, actual)
}

// n1, n3 are symphony networks, n2 is not
func seedNetworks(t *testing.T) {
	_, err := configurator.CreateNetworks(
		[]configurator.Network{
			{
				ID:          "n1",
				Type:        devmand.SymphonyNetworkType,
				Name:        "foobar",
				Description: "Foo Bar",
				Configs: map[string]interface{}{
					orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
				},
			},
			{
				ID:          "n2",
				Type:        "blah",
				Name:        "foobar",
				Description: "Foo Bar",
				Configs:     map[string]interface{}{},
			},
		},
	)
	assert.NoError(t, err)
}
