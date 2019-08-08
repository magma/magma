/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin_test

import (
	"os"
	"testing"

	"magma/lte/cloud/go/lte"
	plugin2 "magma/lte/cloud/go/plugin"
	models2 "magma/lte/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestListNetworks(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	// Test empty response
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/ltenetworks",
		Payload:        nil,
		Handler:        plugin2.ListNetworks,
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
		Handler:        plugin2.ListNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n1", "n3"}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestCreateNetwork(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	tc := tests.Test{
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
		Handler:        plugin2.CreateNetwork,
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
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	// Test 404
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/ltenetworks/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        plugin2.GetNetwork,
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
		Handler:        plugin2.GetNetwork,
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
		Handler:        plugin2.GetNetwork,
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
		Handler:        plugin2.GetNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedN3),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateNetwork(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	// Test 404
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
					ARecord:    []string{"asdf", "hjkl"},
					AaaaRecord: []string{"abcd", "efgh"},
				},
				{
					Domain:  "facebook.com",
					ARecord: []string{"google.com"},
				},
			},
		},
	}
	tc := tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/ltenetworks/n1",
		Payload:        models2.NewDefaultTDDNetworkConfig(),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        plugin2.UpdateNetwork,
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
		Handler:        plugin2.UpdateNetwork,
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
		Handler:        plugin2.UpdateNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not an LTE network",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestDeleteNetwork(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	// Test 404
	tc := tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/ltenetworks/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        plugin2.DeleteNetwork,
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
		Handler:        plugin2.DeleteNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not an LTE network",
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.ListNetworkIDs()
	assert.NoError(t, err)
	assert.Equal(t, []string{"n2", "n3"}, actual)
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
