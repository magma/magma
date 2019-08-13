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

	models1 "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func Test_GetNetworkHandlers(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// Test empty case
	listNetworks := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        handlers.ListNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{}),
	}
	tests.RunUnitTest(t, e, listNetworks)

	// Test 404
	getNetwork := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, "no_such_network"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"no_such_network"},
		Handler:        handlers.GetNetwork,
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, getNetwork)

	// register a network
	networkID1 := "test_network1"
	networkName1 := "network1"
	network1 := configurator.Network{
		ID:   networkID1,
		Name: networkName1,
	}
	err := configurator.CreateNetwork(network1)
	assert.NoError(t, err)

	listNetworks = tests.Test{
		Method:  "GET",
		URL:     testURLRoot,
		Payload: nil,

		Handler:        handlers.ListNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{networkID1}),
	}
	tests.RunUnitTest(t, e, listNetworks)

	expectedNetwork1 := models.Network{
		ID:   models1.NetworkID(networkID1),
		Name: models1.NetworkName(networkName1),
	}

	getNetwork = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, networkID1),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        handlers.GetNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedNetwork1),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, getNetwork)

	// add network features
	networkFeatures1 := models.NewDefaultFeaturesConfig()
	update1 := configurator.NetworkUpdateCriteria{
		ID:                   networkID1,
		ConfigsToAddOrUpdate: map[string]interface{}{orc8r.NetworkFeaturesConfig: networkFeatures1},
	}
	err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{update1})
	assert.NoError(t, err)

	expectedNetwork1 = models.Network{
		ID:       models1.NetworkID(networkID1),
		Name:     models1.NetworkName(networkName1),
		Features: networkFeatures1,
	}

	getNetwork = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, networkID1),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        handlers.GetNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedNetwork1),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, getNetwork)

	// add dnsd configs and a description
	dnsdConfig := models.NewDefaultDNSConfig()
	description1 := "A Network"
	update1 = configurator.NetworkUpdateCriteria{
		ID:                   networkID1,
		NewDescription:       &description1,
		ConfigsToAddOrUpdate: map[string]interface{}{orc8r.DnsdNetworkType: dnsdConfig},
	}
	err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{update1})
	assert.NoError(t, err)

	expectedNetwork1 = models.Network{
		ID:          models1.NetworkID(networkID1),
		Name:        models1.NetworkName(networkName1),
		Description: models1.NetworkDescription("A Network"),
		Features:    networkFeatures1,
		DNS:         dnsdConfig,
	}

	getNetwork = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, networkID1),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        handlers.GetNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedNetwork1),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, getNetwork)

	// register a second network
	networkID2 := "test_network2"
	networkName2 := "network2"
	network2 := configurator.Network{
		ID:   networkID2,
		Name: networkName2,
	}
	err = configurator.CreateNetwork(network2)
	assert.NoError(t, err)

	listNetworks = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        handlers.ListNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{networkID1, networkID2}),
	}
	tests.RunUnitTest(t, e, listNetworks)
}

func Test_PostNetworkHandlers(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// test empty name, description
	networkID1 := "test_network1"
	network1 := models.NewDefaultNetwork(networkID1, "", "")
	postNetwork := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        handlers.RegisterNetwork,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"description in body should be at least 1 chars long\n" +
			"name in body should be at least 1 chars long",
	}
	tests.RunUnitTest(t, e, postNetwork)

	// test bad networkID format
	network1 = models.NewDefaultNetwork("Network*1", "name", "desc")
	postNetwork = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        handlers.RegisterNetwork,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"id in body should match '^[a-z][\\da-z_]+$'",
	}
	tests.RunUnitTest(t, e, postNetwork)

	// test no DNSConfig
	network1 = models.NewDefaultNetwork(networkID1, "name", "desc")
	network1.DNS = nil
	postNetwork = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        handlers.RegisterNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "validation failure list:\ndns in body is required",
	}
	tests.RunUnitTest(t, e, postNetwork)

	// test bad DNSConfig - domain
	network1 = models.NewDefaultNetwork(networkID1, "name", "desc")
	network1.DNS.Records[0].Domain = ""
	postNetwork = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        handlers.RegisterNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "validation failure list:\nvalidation failure list:\nvalidation failure list:\ndomain in body is required",
	}
	tests.RunUnitTest(t, e, postNetwork)

	// test bad DNSConfig - ARecord
	network1 = models.NewDefaultNetwork(networkID1, "name", "desc")
	network1.DNS.Records[0].ARecord[0] = "not ipv4"
	postNetwork = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        handlers.RegisterNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "validation failure list:\nvalidation failure list:\nvalidation failure list:\na_record.0 in body must be of type ipv4: \"not ipv4\"",
	}
	tests.RunUnitTest(t, e, postNetwork)

	// test bad DNSConfig - AaaaRecord
	network1 = models.NewDefaultNetwork(networkID1, "name", "desc")
	network1.DNS.Records[0].AaaaRecord[0] = "not ipv6"
	postNetwork = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        handlers.RegisterNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "validation failure list:\nvalidation failure list:\nvalidation failure list:\naaaa_record.0 in body must be of type ipv6: \"not ipv6\"",
	}
	tests.RunUnitTest(t, e, postNetwork)

	// happy case
	network1 = models.NewDefaultNetwork(networkID1, "name", "desc")
	postNetwork = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        handlers.RegisterNetwork,
		ExpectedStatus: 201,
		ExpectedResult: tests.JSONMarshaler(networkID1),
	}
	tests.RunUnitTest(t, e, postNetwork)

	actualNetwork1, err := configurator.LoadNetwork(networkID1, true, true)
	assert.NoError(t, err)
	expectedNetwork1 := configurator.Network{
		ID:          string(network1.ID),
		Type:        string(network1.Type),
		Name:        string(network1.Name),
		Description: string(network1.Description),
		Configs: map[string]interface{}{
			orc8r.DnsdNetworkType:       models.NewDefaultDNSConfig(),
			orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
		},
	}
	assert.Equal(t, expectedNetwork1, actualNetwork1)

	getNetwork := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, networkID1),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        handlers.GetNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(network1),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, getNetwork)

	listNetworks := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        handlers.ListNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{networkID1}),
	}
	tests.RunUnitTest(t, e, listNetworks)
}

func Test_DeleteNetworkHandlers(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// add a couple of networks
	networkID1 := "test_network1"
	network1 := models.NewDefaultNetwork(networkID1, "name", "desc")
	postNetwork := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        handlers.RegisterNetwork,
		ExpectedStatus: 201,
		ExpectedResult: tests.JSONMarshaler(networkID1),
	}
	tests.RunUnitTest(t, e, postNetwork)

	networkID2 := "test_network2"
	network2 := models.NewDefaultNetwork(networkID2, "name", "desc")
	postNetwork = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network2),
		Handler:        handlers.RegisterNetwork,
		ExpectedStatus: 201,
		ExpectedResult: tests.JSONMarshaler(networkID2),
	}
	tests.RunUnitTest(t, e, postNetwork)

	listNetworks := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        handlers.ListNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{networkID1, networkID2}),
	}
	tests.RunUnitTest(t, e, listNetworks)

	// delete and get
	deleteNetwork := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, networkID1),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        handlers.DeleteNetwork,
		ExpectedStatus: 204,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, deleteNetwork)

	listNetworks = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        handlers.ListNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{networkID2}),
	}
	tests.RunUnitTest(t, e, listNetworks)

	deleteNetwork = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, networkID2),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID2},
		Handler:        handlers.DeleteNetwork,
		ExpectedStatus: 204,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, deleteNetwork)

	listNetworks = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        handlers.ListNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{}),
	}
	tests.RunUnitTest(t, e, listNetworks)
}

func Test_PutNetworkHandlers(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// happy path
	// add a network
	networkID1 := "test_network1"
	network1 := models.NewDefaultNetwork(networkID1, "name", "desc")
	postNetwork := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        handlers.RegisterNetwork,
		ExpectedStatus: 201,
		ExpectedResult: tests.JSONMarshaler(networkID1),
	}
	tests.RunUnitTest(t, e, postNetwork)

	// change meta data
	network1.Name = models1.NetworkName("name2")
	network1.Type = "wifi"
	network1.Description = models1.NetworkDescription("desc2")
	putNetwork := tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        handlers.UpdateNetwork,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, putNetwork)

	getNetwork := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, networkID1),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        handlers.GetNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(network1),
	}
	tests.RunUnitTest(t, e, getNetwork)

	// change configs
	network1.DNS.EnableCaching = swag.Bool(false)
	network1.Features.Features["new-feature"] = "foobar"
	putNetwork = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        handlers.UpdateNetwork,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, putNetwork)

	getNetwork = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, networkID1),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        handlers.GetNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(network1),
	}
	tests.RunUnitTest(t, e, getNetwork)

	// try do delete DNS config
	network1.DNS = nil
	putNetwork = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        handlers.UpdateNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "validation failure list:\ndns in body is required",
	}
	tests.RunUnitTest(t, e, putNetwork)
}

func Test_GetNetworkPartialHandlers(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// register a network
	networkID1 := "test_network1"
	networkName1 := "network1"
	networkDesc1 := "network 1"
	type1 := "type1"
	network1 := configurator.Network{
		ID:          networkID1,
		Name:        networkName1,
		Description: networkDesc1,
		Type:        type1,
		Configs: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
			orc8r.DnsdNetworkType:       models.NewDefaultDNSConfig(),
		},
	}
	err := configurator.CreateNetwork(network1)
	assert.NoError(t, err)

	getName := handlers.GetPartialReadNetworkHandler(handlers.ManageNetworkNamePath, new(models1.NetworkName)).HandlerFunc
	getType := handlers.GetPartialReadNetworkHandler(handlers.ManageNetworkTypePath, new(models1.NetworkType)).HandlerFunc
	getDesc := handlers.GetPartialReadNetworkHandler(handlers.ManageNetworkDescriptionPath, new(models1.NetworkDescription)).HandlerFunc
	getFeatures := handlers.GetPartialReadNetworkHandler(handlers.ManageNetworkFeaturesPath, &models.NetworkFeatures{}).HandlerFunc
	getDNS := handlers.GetPartialReadNetworkHandler(handlers.ManageNetworkDNSPath, &models.NetworkDNSConfig{}).HandlerFunc

	getNetworkName := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/name/", testURLRoot, networkID1),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        getName,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(networkName1),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, getNetworkName)

	getNetworkType := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/type/", testURLRoot, networkID1),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        getType,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(type1),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, getNetworkType)

	getNetworkDesc := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/description/", testURLRoot, networkID1),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        getDesc,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(networkDesc1),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, getNetworkDesc)

	getNetworkFeatures := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/features/", testURLRoot, networkID1),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        getFeatures,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models.NewDefaultFeaturesConfig()),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, getNetworkFeatures)

	getDNSConfig := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/features/", testURLRoot, networkID1),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        getDNS,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models.NewDefaultDNSConfig()),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, getDNSConfig)
}

func Test_PutNetworkPartialHandlers(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// register a network
	networkID1 := "test_network1"
	networkName1 := "network1"

	network1 := configurator.Network{
		ID:      networkID1,
		Name:    networkName1,
		Configs: map[string]interface{}{},
	}
	err := configurator.CreateNetwork(network1)
	assert.NoError(t, err)

	updateName := handlers.GetPartialUpdateNetworkHandler(handlers.ManageNetworkNamePath, new(models1.NetworkName)).HandlerFunc
	updateType := handlers.GetPartialUpdateNetworkHandler(handlers.ManageNetworkTypePath, new(models1.NetworkType)).HandlerFunc
	updateDesc := handlers.GetPartialUpdateNetworkHandler(handlers.ManageNetworkDescriptionPath, new(models1.NetworkDescription)).HandlerFunc
	updateFeatures := handlers.GetPartialUpdateNetworkHandler(handlers.ManageNetworkFeaturesPath, &models.NetworkFeatures{}).HandlerFunc
	updateDNS := handlers.GetPartialUpdateNetworkHandler(handlers.ManageNetworkDNSPath, &models.NetworkDNSConfig{}).HandlerFunc

	// check for validity
	network1.Name = ""
	putNetworkName := tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/name/", testURLRoot, networkID1),
		Payload:        tests.JSONMarshaler(network1.Name),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        updateName,
		ExpectedStatus: 400,
		ExpectedError:  " in body should be at least 1 chars long",
	}
	tests.RunUnitTest(t, e, putNetworkName)

	// happy case
	network1.Name = "new_name"
	network1.Version = 1
	putNetworkName = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/name/", testURLRoot, networkID1),
		Payload:        tests.JSONMarshaler(network1.Name),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        updateName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, putNetworkName)

	actualNetwork, err := configurator.LoadNetwork(networkID1, true, false)
	assert.NoError(t, err)
	assert.Equal(t, network1, actualNetwork)

	// happy case
	network1.Type = "new_type"
	network1.Version = 2
	putNetworkType := tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/type/", testURLRoot, networkID1),
		Payload:        tests.JSONMarshaler(network1.Type),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        updateType,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, putNetworkType)
	actualNetwork, err = configurator.LoadNetwork(networkID1, true, false)
	assert.NoError(t, err)
	assert.Equal(t, network1, actualNetwork)

	// happy case
	network1.Description = "new_name"
	network1.Version = 3
	putNetworkDesc := tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/description/", testURLRoot, networkID1),
		Payload:        tests.JSONMarshaler(network1.Description),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        updateDesc,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, putNetworkDesc)

	actualNetwork, err = configurator.LoadNetwork(networkID1, true, false)
	assert.NoError(t, err)
	assert.Equal(t, network1, actualNetwork)

	// update full feature happy case
	newFeatures := &models.NetworkFeatures{
		Features: map[string]string{
			"hello": "world!!",
		},
	}
	putNetworkFeatures := tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/features/", testURLRoot, networkID1),
		Payload:        tests.JSONMarshaler(newFeatures),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        updateFeatures,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, putNetworkFeatures)

	config, err := configurator.LoadNetworkConfig(networkID1, orc8r.NetworkFeaturesConfig)
	assert.NoError(t, err)
	assert.Equal(t, newFeatures, config)

	// update full dns validation failure
	newDNS := models.NewDefaultDNSConfig()
	newDNS.Records = []*models.DNSConfigRecord{
		{
			ARecord:     []strfmt.IPv4{"192-88-99-142"},
			AaaaRecord:  []strfmt.IPv6{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
			CnameRecord: []string{"facebook.com"},
			Domain:      "facebook.com",
		},
	}
	putDNS := tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/dns/", testURLRoot, networkID1),
		Payload:        tests.JSONMarshaler(newDNS),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        updateDNS,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"validation failure list:\n" +
			"a_record.0 in body must be of type ipv4: \"192-88-99-142\"",
	}
	tests.RunUnitTest(t, e, putDNS)

	// update full DNS happy case
	newDNS = models.NewDefaultDNSConfig()
	newDNS.Records = []*models.DNSConfigRecord{
		{
			ARecord:     []strfmt.IPv4{"192.88.99.142"},
			AaaaRecord:  []strfmt.IPv6{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
			CnameRecord: []string{"facebook.com"},
			Domain:      "facebook.com",
		},
	}
	putDNS = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/dns/", testURLRoot, networkID1),
		Payload:        tests.JSONMarshaler(newDNS),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID1},
		Handler:        updateDNS,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, putDNS)

	config, err = configurator.LoadNetworkConfig(networkID1, orc8r.DnsdNetworkType)
	assert.NoError(t, err)
	assert.Equal(t, newDNS, config)
}
