/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers_test

import (
	"testing"

	models1 "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/configurator/test_utils"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func Test_ListReleaseChannels(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/channels"
	obsidianHandlers := handlers.GetObsidianHandlers()
	listChannels := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, testURLRoot, obsidian.GET).HandlerFunc

	// List channels when none exist
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listChannels,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{}),
	}
	tests.RunUnitTest(t, e, tc)

	// add a channel
	_, err := configurator.CreateInternalEntity(
		configurator.NetworkEntity{
			Type: orc8r.UpgradeReleaseChannelEntityType, Key: "channel1",
			Config: &models.ReleaseChannel{
				ID:                "channel1",
				Name:              "channel 1",
				SupportedVersions: []string{"1-1-1-1"},
			},
		},
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listChannels,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"channel1"}),
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_CreateReleaseChannel(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/channels"
	obsidianHandlers := handlers.GetObsidianHandlers()
	createChannel := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, testURLRoot, obsidian.POST).HandlerFunc

	// happy case
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(&models.ReleaseChannel{ID: "channel1", SupportedVersions: []string{"1-2-3-4"}}),
		Handler:        createChannel,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	entity, err := configurator.LoadInternalEntity(orc8r.UpgradeReleaseChannelEntityType, "channel1", configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true})
	assert.NoError(t, err)
	t.Logf("%v", entity)

	expected := configurator.NetworkEntity{
		NetworkID: "network_magma_internal",
		Type:      orc8r.UpgradeReleaseChannelEntityType, Key: "channel1",
		Config:  &models.ReleaseChannel{ID: "channel1", SupportedVersions: []string{"1-2-3-4"}},
		Version: 0,
		GraphID: "2",
	}
	assert.Equal(t, expected, entity)

	// validation error
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(&models.ReleaseChannel{ID: "", SupportedVersions: []string{"1-2-3-4"}}),
		Handler:        createChannel,
		ExpectedStatus: 400,
		ExpectedError:  "validation failure list:\nid in body should be at least 1 chars long",
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_ReleaseChannel(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/channels/:channel_id"
	obsidianHandlers := handlers.GetObsidianHandlers()
	updateChannel := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, testURLRoot, obsidian.PUT).HandlerFunc
	getChannel := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, testURLRoot, obsidian.GET).HandlerFunc
	deleteChannel := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, testURLRoot, obsidian.DELETE).HandlerFunc

	// 404
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		ParamNames:     []string{"channel_id"},
		ParamValues:    []string{"channel1"},
		Handler:        getChannel,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// add a channel
	_, err := configurator.CreateInternalEntity(
		configurator.NetworkEntity{
			Type: orc8r.UpgradeReleaseChannelEntityType, Key: "channel1",
			Config: &models.ReleaseChannel{
				ID:                "channel1",
				Name:              "channel 1",
				SupportedVersions: []string{"1-1-1-1"},
			},
		},
	)
	assert.NoError(t, err)

	// happy get
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		ParamNames:     []string{"channel_id"},
		ParamValues:    []string{"channel1"},
		Handler:        getChannel,
		ExpectedStatus: 200,
		ExpectedResult: &models.ReleaseChannel{ID: "channel1", Name: "channel 1", SupportedVersions: []string{"1-1-1-1"}},
	}
	tests.RunUnitTest(t, e, tc)

	// happy update
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        &models.ReleaseChannel{ID: "channel1", Name: "modified channel 1", SupportedVersions: []string{}},
		ParamNames:     []string{"channel_id"},
		ParamValues:    []string{"channel1"},
		Handler:        updateChannel,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	channelEntity, err := configurator.LoadInternalEntity(orc8r.UpgradeReleaseChannelEntityType, "channel1", configurator.EntityLoadCriteria{LoadConfig: true})
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "network_magma_internal",
		Type:      orc8r.UpgradeReleaseChannelEntityType, Key: "channel1",
		Config:  &models.ReleaseChannel{ID: "channel1", Name: "modified channel 1", SupportedVersions: []string{}},
		Version: 1,
		GraphID: "2",
	}
	assert.Equal(t, expected, channelEntity)

	// happy delete
	tc = tests.Test{
		Method:         "DELETE",
		URL:            testURLRoot,
		ParamNames:     []string{"channel_id"},
		ParamValues:    []string{"channel1"},
		Handler:        deleteChannel,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.LoadInternalEntity(orc8r.UpgradeReleaseChannelEntityType, "channel1", configurator.EntityLoadCriteria{})
	assert.Error(t, err)
}

func Test_Tiers(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)

	e := echo.New()
	tiersRoot := "/magma/v1/networks/:network_id/tiers"
	manageTiers := tiersRoot + "/:tier_id"
	obsidianHandlers := handlers.GetObsidianHandlers()
	listTiers := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, tiersRoot, obsidian.GET).HandlerFunc
	createTier := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, tiersRoot, obsidian.POST).HandlerFunc
	updateTier := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers, obsidian.PUT).HandlerFunc
	readTier := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers, obsidian.GET).HandlerFunc
	deleteTier := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers, obsidian.DELETE).HandlerFunc

	assert.NoError(t, configurator.CreateNetwork(configurator.Network{ID: "n1"}))

	// happy case list
	tc := tests.Test{
		Method:         "GET",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		URL:            tiersRoot,
		Handler:        listTiers,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{}),
	}
	tests.RunUnitTest(t, e, tc)

	// validation failure
	tc = tests.Test{
		Method:         "POST",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Payload:        &models.Tier{ID: models.TierID("t*i*e*r*1")},
		URL:            tiersRoot,
		Handler:        createTier,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"gateways in body is required\n" +
			"id in body should match '^[a-z][\\da-z_]+$'\n" +
			"images in body is required\n" +
			"version in body is required",
	}
	tests.RunUnitTest(t, e, tc)

	// gateway does not exist
	tc = tests.Test{
		Method:         "POST",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Payload:        &models.Tier{ID: models.TierID("tier1"), Images: []*models.TierImage{}, Gateways: []models1.GatewayID{"g1"}, Version: swag.String("1.2.3.4")},
		URL:            tiersRoot,
		Handler:        createTier,
		ExpectedStatus: 500,
		ExpectedError:  "could not find entities matching [type:\"magmad_gateway\" key:\"g1\" ]",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case create
	test_utils.RegisterGateway(t, "n1", "g1", nil)
	tier := &models.Tier{ID: models.TierID("tier1"), Images: []*models.TierImage{}, Gateways: []models1.GatewayID{"g1"}, Version: swag.String("1.2.3.4")}
	tc = tests.Test{
		Method:         "POST",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Payload:        tier,
		URL:            tiersRoot,
		Handler:        createTier,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	entities, _, err := configurator.LoadEntities("n1", swag.String(orc8r.UpgradeTierEntityType), nil, nil, nil, configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	expected := map[storage.TypeAndKey]configurator.NetworkEntity{
		storage.TypeAndKey{Type: orc8r.UpgradeTierEntityType, Key: "tier1"}: {
			NetworkID: "n1",
			Type:      orc8r.UpgradeTierEntityType, Key: "tier1",
			Name:         tier.Name,
			Config:       tier,
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
			GraphID:      "4",
			Version:      0,
		},
	}
	assert.Equal(t, expected, entities.ToEntitiesByID())

	// happy case update
	tier.Name = "new name!"
	// switch to a new gateway
	test_utils.RegisterGateway(t, "n1", "g2", nil)
	tier.Gateways = []models1.GatewayID{"g2"}
	tc = tests.Test{
		Method:         "PUT",
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		Payload:        tier,
		URL:            tiersRoot,
		Handler:        updateTier,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	entitiesToQuery := []storage.TypeAndKey{
		{Type: orc8r.MagmadGatewayType, Key: "g1"},
		{Type: orc8r.MagmadGatewayType, Key: "g2"},
		{Type: orc8r.UpgradeTierEntityType, Key: "tier1"},
	}
	entities, _, err = configurator.LoadEntities("n1", nil, nil, nil, entitiesToQuery, configurator.FullEntityLoadCriteria())
	expected = map[storage.TypeAndKey]configurator.NetworkEntity{
		storage.TypeAndKey{Type: orc8r.UpgradeTierEntityType, Key: "tier1"}: {
			NetworkID: "n1",
			Type:      orc8r.UpgradeTierEntityType, Key: "tier1",
			Name:         tier.Name,
			Config:       tier,
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g2"}},
			GraphID:      "4",
			Version:      1,
		},
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			GraphID: "9",
			Version: 0,
		},
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g2"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g2",
			GraphID:            "4",
			Version:            0,
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "tier1"}},
		},
	}
	assert.Equal(t, expected, entities.ToEntitiesByID())

	// happy case read
	tc = tests.Test{
		Method:         "GET",
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		Payload:        tier,
		URL:            tiersRoot,
		Handler:        readTier,
		ExpectedStatus: 200,
		ExpectedResult: tier,
	}
	tests.RunUnitTest(t, e, tc)

	// happy case list
	tc = tests.Test{
		Method:         "GET",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		URL:            tiersRoot,
		Handler:        listTiers,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"tier1"}),
	}
	tests.RunUnitTest(t, e, tc)

	// 404 read
	tc = tests.Test{
		Method:         "GET",
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier2"},
		Payload:        tier,
		URL:            tiersRoot,
		Handler:        readTier,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case delete
	tc = tests.Test{
		Method:         "DELETE",
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		Payload:        tier,
		URL:            tiersRoot,
		Handler:        deleteTier,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	entities, _, err = configurator.LoadEntities("n1", nil, nil, nil, entitiesToQuery, configurator.FullEntityLoadCriteria())
	expected = map[storage.TypeAndKey]configurator.NetworkEntity{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			GraphID: "9",
			Version: 0,
		},
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g2"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g2",
			GraphID: "4",
			Version: 0,
		},
	}
	assert.Equal(t, expected, entities.ToEntitiesByID())
}
