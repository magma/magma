/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package handlers_test

import (
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/device"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/storage"
	"orc8r/devmand/cloud/go/devmand"
	plugin2 "orc8r/devmand/cloud/go/plugin"
	"orc8r/devmand/cloud/go/plugin/handlers"
	models2 "orc8r/devmand/cloud/go/plugin/models"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestListNetworks(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
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

	obsidianHandlers := handlers.GetHandlers()
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
				Name:        "agent_1",
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
				Description: "Network 1",
				ID:          "n1",
				Name:        "agent_1",
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
		Name:        "agent_1",
		Description: "Network 1",
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
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
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
		Description: "Network 1",
		ID:          "n1",
		Name:        "network_1",
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
		ExpectedError:  "network n2 is not a Symphony network",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateNetwork(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	updateNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/symphony/:network_id", obsidian.PUT).HandlerFunc

	// Test validation failure
	payloadN1 := &models2.SymphonyNetwork{
		ID:          "n1",
		Name:        "updated network_1",
		Description: "Updated Network 1",
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
		Name:        "updated network_1",
		Description: "Updated Network 1",
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
		ExpectedError:  "network n2 is not a Symphony network",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestDeleteNetwork(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
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
		ExpectedError:  "network n2 is not a Symphony network",
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.ListNetworkIDs()
	assert.NoError(t, err)
	assert.Equal(t, []string{"n2"}, actual)
}

func TestListAgents(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	baseAgentsUrl := "/magma/v1/symphony/:network_id/agents"
	obsidianHandlers := handlers.GetHandlers()
	listAgents := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, baseAgentsUrl, obsidian.GET).HandlerFunc

	// Test 500
	tc := tests.Test{
		Method:         "GET",
		URL:            baseAgentsUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listAgents,
		ExpectedError:  "Not found",
		ExpectedStatus: 500,
	}
	tests.RunUnitTest(t, e, tc)

	// Test network with no agents
	seedNetworks(t)
	tc = tests.Test{
		Method:         "GET",
		URL:            baseAgentsUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listAgents,
		ExpectedResult: tests.JSONMarshaler(map[string]models2.SymphonyAgent{}),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	expectedEnts := configurator.NetworkEntities{}
	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnts, actualEnts)

	// Test network with one agent
	expectedResult := models2.SymphonyAgent{
		Name:        "agent_1",
		Description: "agent 1",
		ID:          "a1",
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		ManagedDevices: []string{"d1", "d2"},
		Tier:           "t1",
	}

	seedAgents(t)
	tc = tests.Test{
		Method:         "GET",
		URL:            baseAgentsUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listAgents,
		ExpectedResult: tests.JSONMarshaler(map[string]models2.SymphonyAgent{"a1": expectedResult}),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	expectedEnts = configurator.NetworkEntities{
		{
			NetworkID:          "n1",
			Type:               devmand.SymphonyAgentType,
			Key:                "a1",
			GraphID:            "10",
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "a1"}},
			Associations: []storage.TypeAndKey{
				{Type: devmand.SymphonyDeviceType, Key: "d1"},
				{Type: devmand.SymphonyDeviceType, Key: "d2"},
			},
		},
		{
			NetworkID:   "n1",
			Type:        orc8r.MagmadGatewayType,
			Key:         "a1",
			Name:        "agent_1",
			Description: "agent 1",
			PhysicalID:  "hw1",
			GraphID:     "10",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Associations:       []storage.TypeAndKey{{Type: devmand.SymphonyAgentType, Key: "a1"}},
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
		},
	}
	actualEnts, _, err = configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "a1"},
			{Type: devmand.SymphonyAgentType, Key: "a1"},
		},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnts, actualEnts)
}

func TestCreateAgent(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	baseAgentsUrl := "/magma/v1/symphony/:network_id/agents"
	obsidianHandlers := handlers.GetHandlers()
	createAgent := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, baseAgentsUrl, obsidian.POST).HandlerFunc

	// Initially empty
	seedNetworks(t)
	expectedEnts := configurator.NetworkEntities{}
	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnts, actualEnts)

	// Test missing payload
	tc := tests.Test{
		Method:         "POST",
		URL:            baseAgentsUrl,
		Payload:        &models2.SymphonyAgent{},
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createAgent,
		ExpectedStatus: 400,
		ExpectedError:  "validation failure list:\ndescription in body should be at least 1 chars long\ndevice in body is required\nid in body should be at least 1 chars long\nmagmad in body is required\nmanaged_devices in body is required\nname in body should be at least 1 chars long\ntier in body should match '^[a-z][\\da-z_]+$'",
	}
	tests.RunUnitTest(t, e, tc)

	// Test post new agent
	seedPreAgent(t)
	payload := &models2.SymphonyAgent{
		Name:        "agent_1",
		Description: "agent 1",
		ID:          "a1",
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		ManagedDevices: []string{},
		Tier:           "t1",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            baseAgentsUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createAgent,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	expectedEnts = configurator.NetworkEntities{
		{
			NetworkID:          "n1",
			Type:               devmand.SymphonyAgentType,
			Key:                "a1",
			GraphID:            "2",
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "a1"}},
		},
		{
			NetworkID:   "n1",
			Type:        orc8r.MagmadGatewayType,
			Key:         "a1",
			Name:        "agent_1",
			Description: "agent 1",
			PhysicalID:  "hw1",
			GraphID:     "2",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Associations:       []storage.TypeAndKey{{Type: devmand.SymphonyAgentType, Key: "a1"}},
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			Version:            1,
		},
		{
			NetworkID: "n1",
			Type:      orc8r.UpgradeTierEntityType, Key: "t1",
			GraphID:      "2",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "a1"}},
			Version:      1,
		},
	}
	actualEnts, _, err = configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnts, actualEnts)

	// Can't register the same device
	tc = tests.Test{
		Method:         "POST",
		URL:            baseAgentsUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createAgent,
		ExpectedStatus: 400,
		ExpectedError:  "device hw1 is already mapped to gateway a1",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetAgent(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	agentUrl := "/magma/v1/symphony/:network_id/agents/:agent_id"
	obsidianHandlers := handlers.GetHandlers()
	getAgent := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, agentUrl, obsidian.GET).HandlerFunc

	agent_id := "a1"

	// Test with missing agent
	seedNetworks(t)
	tc := tests.Test{
		Method:         "GET",
		URL:            agentUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{"n1", agent_id},
		Handler:        getAgent,
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// Get agent correctly
	seedAgents(t)
	expectedResult := tests.JSONMarshaler(models2.SymphonyAgent{
		Name:        "agent_1",
		Description: "agent 1",
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key: &models.ChallengeKey{
				KeyType: "ECHO",
			},
		},
		ID: "a1",
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		ManagedDevices: []string{"d1", "d2"},
		Tier:           "t1",
	})
	tc = tests.Test{
		Method:         "GET",
		URL:            agentUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{"n1", agent_id},
		Handler:        getAgent,
		ExpectedStatus: 200,
		ExpectedResult: expectedResult,
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateAgent(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	agentUrl := "/magma/v1/symphony/:network_id/agents/:agent_id"
	obsidianHandlers := handlers.GetHandlers()
	updateAgent := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, agentUrl, obsidian.PUT).HandlerFunc

	agent_id := "a1"

	// Test with missing agent
	seedNetworks(t)
	payload := &models2.SymphonyAgent{
		Name:        "agent_1",
		Description: "UPDATED agent 1",
		ID:          "a1",
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 200,
			CheckinInterval:         20,
			CheckinTimeout:          20,
		},
		ManagedDevices: []string{},
		Tier:           "t1",
	}
	tc := tests.Test{
		Method:         "PUT",
		URL:            agentUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{"n1", agent_id},
		Handler:        updateAgent,
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// Test updating agent correctly
	seedAgents(t)
	tc = tests.Test{
		Method:         "PUT",
		URL:            agentUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{"n1", agent_id},
		Handler:        updateAgent,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expectedEnts := configurator.NetworkEntities{
		{
			NetworkID:          "n1",
			Type:               devmand.SymphonyAgentType,
			Key:                "a1",
			GraphID:            "10",
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "a1"}},
			Associations:       []storage.TypeAndKey{storage.TypeAndKey{Type: "symphony_device", Key: "d1"}, storage.TypeAndKey{Type: "symphony_device", Key: "d2"}},
			Version:            1,
		},
		{
			NetworkID:   "n1",
			Type:        orc8r.MagmadGatewayType,
			Key:         "a1",
			Name:        "agent_1",
			Description: "UPDATED agent 1",
			PhysicalID:  "hw1",
			GraphID:     "10",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 200,
				CheckinInterval:         20,
				CheckinTimeout:          20,
			},
			Associations:       []storage.TypeAndKey{{Type: devmand.SymphonyAgentType, Key: "a1"}},
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			Version:            1,
		},
		{
			NetworkID: "n1",
			Type:      devmand.SymphonyDeviceType,
			Key:       "d1", GraphID: "10",
			ParentAssociations: []storage.TypeAndKey{storage.TypeAndKey{Type: "agent", Key: "a1"}},
		},
		{
			NetworkID: "n1",
			Type:      devmand.SymphonyDeviceType,
			Key:       "d2", GraphID: "10",
			ParentAssociations: []storage.TypeAndKey{storage.TypeAndKey{Type: "agent", Key: "a1"}},
		},
		{
			NetworkID: "n1",
			Type:      orc8r.UpgradeTierEntityType, Key: "t1",
			GraphID:      "10",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "a1"}},
			Version:      0,
		},
		{
			NetworkID: "n1",
			Type:      orc8r.UpgradeTierEntityType, Key: "t2",
			GraphID: "12",
			Version: 0,
		},
	}
	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnts, actualEnts)
}

func TestPartialUpdateAndGetAgent(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	agentUrl := "/magma/v1/symphony/:network_id/agents/:agent_id"
	nameUrl := fmt.Sprintf("%s/name", agentUrl)
	descriptionUrl := fmt.Sprintf("%s/description", agentUrl)
	magmadUrl := fmt.Sprintf("%s/magmad", agentUrl)
	tierUrl := fmt.Sprintf("%s/tier", agentUrl)
	deviceUrl := fmt.Sprintf("%s/device", agentUrl)
	managedDevicesUrl := fmt.Sprintf("%s/managed_devices", agentUrl)
	obsidianHandlers := handlers.GetHandlers()
	updateName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, nameUrl, obsidian.PUT).HandlerFunc
	updateDescription := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, descriptionUrl, obsidian.PUT).HandlerFunc
	updateMagmad := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, magmadUrl, obsidian.PUT).HandlerFunc
	updateTier := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, tierUrl, obsidian.PUT).HandlerFunc
	updateDevice := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, deviceUrl, obsidian.PUT).HandlerFunc
	updateManagedDevices := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, managedDevicesUrl, obsidian.PUT).HandlerFunc
	getName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, nameUrl, obsidian.GET).HandlerFunc
	getDescription := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, descriptionUrl, obsidian.GET).HandlerFunc
	getMagmad := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, magmadUrl, obsidian.GET).HandlerFunc
	getTier := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, tierUrl, obsidian.GET).HandlerFunc
	getDevice := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, deviceUrl, obsidian.GET).HandlerFunc
	getManagedDevices := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, managedDevicesUrl, obsidian.GET).HandlerFunc

	networkId := "n1"
	agentId := "a1"

	seedNetworks(t)
	seedAgents(t)

	expectedEnts := map[string]*configurator.NetworkEntity{
		"agent": &configurator.NetworkEntity{
			NetworkID:          networkId,
			Type:               devmand.SymphonyAgentType,
			Key:                agentId,
			GraphID:            "10",
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "a1"}},
			Associations:       []storage.TypeAndKey{storage.TypeAndKey{Type: "symphony_device", Key: "d1"}, storage.TypeAndKey{Type: "symphony_device", Key: "d2"}},
			Version:            0,
		},
		"gateway": &configurator.NetworkEntity{
			NetworkID:   networkId,
			Type:        orc8r.MagmadGatewayType,
			Key:         "a1",
			Name:        "agent_1",
			Description: "agent 1",
			PhysicalID:  "hw1",
			GraphID:     "10",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Associations:       []storage.TypeAndKey{{Type: devmand.SymphonyAgentType, Key: "a1"}},
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			Version:            0,
		},
		"device1": &configurator.NetworkEntity{
			NetworkID:          networkId,
			Type:               devmand.SymphonyDeviceType,
			Key:                "d1",
			GraphID:            "10",
			ParentAssociations: []storage.TypeAndKey{storage.TypeAndKey{Type: "agent", Key: "a1"}},
		},
		"device2": &configurator.NetworkEntity{
			NetworkID:          networkId,
			Type:               devmand.SymphonyDeviceType,
			Key:                "d2",
			GraphID:            "10",
			ParentAssociations: []storage.TypeAndKey{storage.TypeAndKey{Type: "agent", Key: "a1"}},
		},
		"tier1": &configurator.NetworkEntity{
			NetworkID: networkId,
			Type:      orc8r.UpgradeTierEntityType, Key: "t1",
			GraphID:      "10",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "a1"}},
			Version:      0,
		},
		"tier2": &configurator.NetworkEntity{
			NetworkID: networkId,
			Type:      orc8r.UpgradeTierEntityType, Key: "t2",
			GraphID: "12",
			Version: 0,
		},
	}

	// Test updating agent name

	// Update the ents that we expect and then convert them into a list so we can
	// compare them to what we get from configurator.LoadEntities later
	updatedName := "updated_agent_name"
	expectedEnts["gateway"].Name = updatedName
	expectedEnts["gateway"].Version++
	expectedEntsVals := make(configurator.NetworkEntities, 0, len(expectedEnts))
	// Key order matters to compare later
	key_order := []string{"agent", "gateway", "device1", "device2", "tier1", "tier2"}
	for _, v := range key_order {
		expectedEntsVals = append(expectedEntsVals, *expectedEnts[v])
	}

	payload := tests.JSONMarshaler(updatedName)
	tc := tests.Test{
		Method:         "PUT",
		URL:            nameUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        updateName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err := configurator.LoadEntities(
		networkId, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	tc = tests.Test{
		Method:         "GET",
		URL:            nameUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        getName,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test updating agent description
	updatedDescription := "updated_description"
	expectedEnts["gateway"].Description = updatedDescription
	expectedEnts["gateway"].Version++
	expectedEntsVals = make(configurator.NetworkEntities, 0, len(expectedEnts))
	for _, v := range key_order {
		expectedEntsVals = append(expectedEntsVals, *expectedEnts[v])
	}

	payload = tests.JSONMarshaler(updatedDescription)
	tc = tests.Test{
		Method:         "PUT",
		URL:            descriptionUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        updateDescription,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err = configurator.LoadEntities(
		networkId, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	tc = tests.Test{
		Method:         "GET",
		URL:            descriptionUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        getDescription,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test updating agent magmad
	updatedMagmad := &models.MagmadGatewayConfigs{
		AutoupgradeEnabled:      swag.Bool(false),
		AutoupgradePollInterval: 100,
		CheckinInterval:         30,
		CheckinTimeout:          10,
	}
	expectedEnts["gateway"].Config = updatedMagmad
	expectedEnts["gateway"].Version++
	expectedEntsVals = make(configurator.NetworkEntities, 0, len(expectedEnts))
	for _, v := range key_order {
		expectedEntsVals = append(expectedEntsVals, *expectedEnts[v])
	}

	payload = tests.JSONMarshaler(updatedMagmad)
	tc = tests.Test{
		Method:         "PUT",
		URL:            magmadUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        updateMagmad,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err = configurator.LoadEntities(
		networkId, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	tc = tests.Test{
		Method:         "GET",
		URL:            magmadUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        getMagmad,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test updating agent tier
	updatedAgentTier := "t2"
	expectedEnts["gateway"].ParentAssociations = []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: updatedAgentTier}}
	expectedEnts["tier1"].Associations = nil
	expectedEnts["tier1"].GraphID = "13"
	expectedEnts["tier1"].Version++
	expectedEnts["tier2"].Associations = []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "a1"}}
	expectedEnts["tier2"].GraphID = "10"
	expectedEnts["tier2"].Version++
	expectedEntsVals = make(configurator.NetworkEntities, 0, len(expectedEnts))
	for _, v := range key_order {
		expectedEntsVals = append(expectedEntsVals, *expectedEnts[v])
	}

	payload = tests.JSONMarshaler(updatedAgentTier)
	tc = tests.Test{
		Method:         "PUT",
		URL:            tierUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        updateTier,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err = configurator.LoadEntities(
		networkId, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	tc = tests.Test{
		Method:         "GET",
		URL:            tierUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        getTier,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test updating agent device
	updatedDeviceHardwareId := "hw2"
	payload = tests.JSONMarshaler(&models.GatewayDevice{
		HardwareID: updatedDeviceHardwareId,
		Key: &models.ChallengeKey{
			KeyType: "ECHO",
		},
	})
	tc = tests.Test{
		Method:         "PUT",
		URL:            deviceUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        updateDevice,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Configurator ents should not have changed
	actualEnts, _, err = configurator.LoadEntities(
		networkId, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	// But the HardwareID of the physical device with id "hw1" should have updated
	expectedDevice := &models.GatewayDevice{HardwareID: "hw2", Key: &models.ChallengeKey{KeyType: "ECHO"}}
	actualDevice, err := device.GetDevice(networkId, orc8r.AccessGatewayRecordType, "hw1")
	assert.NoError(t, err)
	assert.Equal(t, expectedDevice, actualDevice)

	tc = tests.Test{
		Method:         "GET",
		URL:            deviceUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        getDevice,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test updating agent managed devices
	updatedManagedDevices := []string{"d1"}
	expectedEnts["agent"].Associations = []storage.TypeAndKey{storage.TypeAndKey{Type: "symphony_device", Key: "d1"}}
	expectedEnts["agent"].Version++
	expectedEnts["device2"].ParentAssociations = nil
	expectedEnts["device2"].GraphID = "14"
	expectedEntsVals = make(configurator.NetworkEntities, 0, len(expectedEnts))
	for _, v := range key_order {
		expectedEntsVals = append(expectedEntsVals, *expectedEnts[v])
	}

	payload = tests.JSONMarshaler(updatedManagedDevices)
	tc = tests.Test{
		Method:         "PUT",
		URL:            managedDevicesUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        updateManagedDevices,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err = configurator.LoadEntities(
		networkId, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	tc = tests.Test{
		Method:         "GET",
		URL:            managedDevicesUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        getManagedDevices,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
}

func TestDeleteAgent(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.DevmandOrchestratorPlugin{})
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	agentUrl := "/magma/v1/symphony/:network_id/agents/:agent_id"
	obsidianHandlers := handlers.GetHandlers()
	getAgent := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, agentUrl, obsidian.GET).HandlerFunc
	deleteAgent := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, agentUrl, obsidian.DELETE).HandlerFunc

	networkId := "n1"
	agentId := "a1"

	seedNetworks(t)
	seedAgents(t)

	expectedResult := tests.JSONMarshaler(models2.SymphonyAgent{
		Name:        "agent_1",
		Description: "agent 1",
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key: &models.ChallengeKey{
				KeyType: "ECHO",
			},
		},
		ID: "a1",
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		ManagedDevices: []string{"d1", "d2"},
		Tier:           "t1",
	})

	tc := tests.Test{
		Method:         "GET",
		URL:            agentUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        getAgent,
		ExpectedResult: expectedResult,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "DELETE",
		URL:            agentUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        deleteAgent,
		ExpectedResult: nil,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            agentUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id", "agent_id"},
		ParamValues:    []string{networkId, agentId},
		Handler:        getAgent,
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)
}

// n1 is a symphony network, n2 is not
func seedNetworks(t *testing.T) {

	gatewayRecord := &models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}}
	err := device.RegisterDevice("n1", orc8r.AccessGatewayRecordType, "hw1", gatewayRecord)
	assert.NoError(t, err)

	_, err = configurator.CreateNetworks(
		[]configurator.Network{
			models2.NewDefaultSymphonyNetwork().ToConfiguratorNetwork(),
			{
				ID:          "n2",
				Type:        "blah",
				Name:        "network_2",
				Description: "Network 2",
				Configs:     map[string]interface{}{},
			},
		},
	)
	assert.NoError(t, err)
}

func seedPreAgent(t *testing.T) {
	// Create Tier necessary for the Agent's gateway to be in
	_, err := configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
			},
		},
	)
	assert.NoError(t, err)
}

func seedAgents(t *testing.T) {
	nID := "n1"
	_, err := configurator.CreateEntities(
		nID,
		[]configurator.NetworkEntity{
			{
				Type: devmand.SymphonyDeviceType, Key: "d1",
			},
			{
				Type: devmand.SymphonyDeviceType, Key: "d2",
			},
			{
				Type: devmand.SymphonyAgentType, Key: "a1",
				Associations: []storage.TypeAndKey{
					{Type: devmand.SymphonyDeviceType, Key: "d1"},
					{Type: devmand.SymphonyDeviceType, Key: "d2"},
				},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: "a1",
				Name:        "agent_1",
				Description: "agent 1",
				PhysicalID:  "hw1",
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
				Associations: []storage.TypeAndKey{{Type: devmand.SymphonyAgentType, Key: "a1"}},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: "a1"},
				},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t2",
			},
		},
	)
	assert.NoError(t, err)
}
