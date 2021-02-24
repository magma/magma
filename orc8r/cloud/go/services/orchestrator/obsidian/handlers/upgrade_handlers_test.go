/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers_test

import (
	"testing"

	models1 "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/configurator/test_utils"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func Test_ListReleaseChannels(t *testing.T) {
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
		serdes.Entity,
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

	entity, err := configurator.LoadInternalEntity(
		orc8r.UpgradeReleaseChannelEntityType, "channel1", configurator.EntityLoadCriteria{LoadMetadata: true,
			LoadConfig: true},
		serdes.Entity,
	)
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
		serdes.Entity,
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

	channelEntity, err := configurator.LoadInternalEntity(
		orc8r.UpgradeReleaseChannelEntityType, "channel1",
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
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

	_, err = configurator.LoadInternalEntity(
		orc8r.UpgradeReleaseChannelEntityType, "channel1",
		configurator.EntityLoadCriteria{},
		serdes.Entity,
	)
	assert.Error(t, err)
}

func Test_Tiers(t *testing.T) {
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

	assert.NoError(t, configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network))

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
			"version in body should be at least 1 chars long",
	}
	tests.RunUnitTest(t, e, tc)

	// gateway does not exist
	tc = tests.Test{
		Method:         "POST",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Payload:        &models.Tier{ID: models.TierID("tier1"), Images: []*models.TierImage{}, Gateways: []models1.GatewayID{"g1"}, Version: "1.2.3.4"},
		URL:            tiersRoot,
		Handler:        createTier,
		ExpectedStatus: 500,
		ExpectedError:  "could not find entities matching [type:\"magmad_gateway\" key:\"g1\" ]",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case create
	test_utils.RegisterGateway(t, "n1", "g1", nil)
	tier := &models.Tier{ID: models.TierID("tier1"), Images: []*models.TierImage{}, Gateways: []models1.GatewayID{"g1"}, Version: "1.2.3.4"}
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

	entities, _, err := configurator.LoadEntities(
		"n1", swag.String(orc8r.UpgradeTierEntityType), nil, nil, nil,
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	expected := configurator.NetworkEntitiesByTK{
		storage.TypeAndKey{Type: orc8r.UpgradeTierEntityType, Key: "tier1"}: {
			NetworkID: "n1",
			Type:      orc8r.UpgradeTierEntityType, Key: "tier1",
			Name:         string(tier.Name),
			Config:       tier,
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
			GraphID:      "4",
			Version:      0,
		},
	}
	assert.Equal(t, expected, entities.MakeByTK())

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
	entities, _, err = configurator.LoadEntities(
		"n1", nil, nil, nil, entitiesToQuery,
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	expected = configurator.NetworkEntitiesByTK{
		storage.TypeAndKey{Type: orc8r.UpgradeTierEntityType, Key: "tier1"}: {
			NetworkID: "n1",
			Type:      orc8r.UpgradeTierEntityType, Key: "tier1",
			Name:         string(tier.Name),
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
	assert.Equal(t, expected, entities.MakeByTK())

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

	entities, _, err = configurator.LoadEntities(
		"n1", nil, nil, nil, entitiesToQuery,
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	expected = configurator.NetworkEntitiesByTK{
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
	assert.Equal(t, expected, entities.MakeByTK())
}

func TestPartialTierReads(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	tiersRoot := "/magma/v1/networks/:network_id/tiers"
	manageTiers := tiersRoot + "/:tier_id"

	// register a network, gateways and a tier
	test_utils.RegisterNetwork(t, "n1", "network 1")
	test_utils.RegisterGateway(t, "n1", "g1", nil)
	tier := &models.Tier{
		Gateways: models.TierGateways([]models1.GatewayID{"g1"}),
		ID:       models.TierID("tier1"),
		Images:   models.TierImages{{Name: swag.String("image1"), Order: swag.Int64(0)}},
		Name:     "tier 1",
		Version:  "1-1-1-1",
	}

	_, err := configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: orc8r.UpgradeTierEntityType, Key: "tier1",
			Name:         string(tier.Name),
			Config:       tier,
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	obsidianHandlers := handlers.GetObsidianHandlers()
	getTierName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers+"/name", obsidian.GET).HandlerFunc
	getTierVersion := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers+"/version", obsidian.GET).HandlerFunc
	getTierImages := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers+"/images", obsidian.GET).HandlerFunc
	getTierGateways := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers+"/gateways", obsidian.GET).HandlerFunc

	// happy case name
	tc := tests.Test{
		Method:         "GET",
		URL:            manageTiers + "/name",
		Handler:        getTierName,
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models.TierName("tier 1")),
	}
	tests.RunUnitTest(t, e, tc)

	// happy case version
	tc = tests.Test{
		Method:         "GET",
		URL:            manageTiers + "/version",
		Handler:        getTierVersion,
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models.TierVersion("1-1-1-1")),
	}
	tests.RunUnitTest(t, e, tc)

	// 404
	tc = tests.Test{
		Method:         "GET",
		URL:            manageTiers + "/version",
		Handler:        getTierVersion,
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier2"},
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case images
	tc = tests.Test{
		Method:         "GET",
		URL:            manageTiers + "/images",
		Handler:        getTierImages,
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(tier.Images),
	}
	tests.RunUnitTest(t, e, tc)

	// happy case gateways
	tc = tests.Test{
		Method:         "GET",
		URL:            manageTiers + "/gateways",
		Handler:        getTierGateways,
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(tier.Gateways),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestPartialTierUpdates(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	tiersRoot := "/magma/v1/networks/:network_id/tiers"
	manageTiers := tiersRoot + "/:tier_id"

	// register a network, gateways and a tier
	test_utils.RegisterNetwork(t, "n1", "network 1")
	test_utils.RegisterGateway(t, "n1", "g1", nil)
	test_utils.RegisterGateway(t, "n1", "g2", nil)
	test_utils.RegisterGateway(t, "n1", "g3", nil)
	tier := &models.Tier{
		Gateways: models.TierGateways([]models1.GatewayID{"g1"}),
		ID:       models.TierID("tier1"),
		Images:   models.TierImages{{Name: swag.String("image1"), Order: swag.Int64(0)}},
		Name:     "tier 1",
		Version:  "1-1-1-1",
	}

	_, err := configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: orc8r.UpgradeTierEntityType, Key: "tier1",
			Name:         string(tier.Name),
			Config:       tier,
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	obsidianHandlers := handlers.GetObsidianHandlers()
	updateTierName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers+"/name", obsidian.PUT).HandlerFunc
	updateTierVersion := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers+"/version", obsidian.PUT).HandlerFunc
	updateTierImages := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers+"/images", obsidian.PUT).HandlerFunc
	updateTierGateways := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers+"/gateways", obsidian.PUT).HandlerFunc
	createTierImage := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers+"/images", obsidian.POST).HandlerFunc
	deleteTierImage := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers+"/images/:image_name", obsidian.DELETE).HandlerFunc
	createTierGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers+"/gateways", obsidian.POST).HandlerFunc
	deleteTierGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, manageTiers+"/gateways/:gateway_id", obsidian.DELETE).HandlerFunc

	// happy case name
	tc := tests.Test{
		Method:         "PUT",
		URL:            manageTiers + "/name",
		Handler:        updateTierName,
		Payload:        tests.JSONMarshaler("new name 1"),
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expectedTier := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      orc8r.UpgradeTierEntityType, Key: "tier1",
		Name:         "new name 1",
		Config:       tier,
		Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
		GraphID:      "2",
		Version:      1,
	}
	actualTier, err := configurator.LoadEntity(
		"n1", orc8r.UpgradeTierEntityType, "tier1",
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadConfig: true, LoadMetadata: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedTier, actualTier)

	// happy case version
	tc = tests.Test{
		Method:         "PUT",
		URL:            manageTiers + "/version",
		Handler:        updateTierVersion,
		Payload:        tests.JSONMarshaler("2-2-2-2"),
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	tier.Version = "2-2-2-2"
	expectedTier = configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      orc8r.UpgradeTierEntityType, Key: "tier1",
		Name:         "new name 1",
		Config:       tier,
		Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
		GraphID:      "2",
		Version:      2,
	}
	actualTier, err = configurator.LoadEntity(
		"n1", orc8r.UpgradeTierEntityType, "tier1",
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadConfig: true, LoadMetadata: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedTier, actualTier)

	// happy case images
	tc = tests.Test{
		Method:         "PUT",
		URL:            manageTiers + "/images",
		Handler:        updateTierImages,
		Payload:        tests.JSONMarshaler(models.TierImages{}),
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	tier.Images = models.TierImages{}
	expectedTier = configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      orc8r.UpgradeTierEntityType, Key: "tier1",
		Name:         "new name 1",
		Config:       tier,
		Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
		GraphID:      "2",
		Version:      3,
	}
	actualTier, err = configurator.LoadEntity(
		"n1", orc8r.UpgradeTierEntityType, "tier1",
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadConfig: true, LoadMetadata: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedTier, actualTier)

	// happy case gateways
	tc = tests.Test{
		Method:         "PUT",
		URL:            manageTiers + "/images",
		Handler:        updateTierGateways,
		Payload:        tests.JSONMarshaler(models.TierGateways{"g1", "g2"}),
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expectedTier = configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      orc8r.UpgradeTierEntityType, Key: "tier1",
		Name:   "new name 1",
		Config: tier,
		Associations: []storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g1"},
			{Type: orc8r.MagmadGatewayType, Key: "g2"},
		},
		GraphID: "2",
		Version: 4,
	}
	actualTier, err = configurator.LoadEntity(
		"n1", orc8r.UpgradeTierEntityType, "tier1",
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadConfig: true, LoadMetadata: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedTier, actualTier)

	// happy case add a new image
	tc = tests.Test{
		Method:         "PUT",
		URL:            manageTiers + "/images",
		Handler:        createTierImage,
		Payload:        tests.JSONMarshaler(models.TierImage{Order: swag.Int64(1), Name: swag.String("image2")}),
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	tier.Images = models.TierImages{{Order: swag.Int64(1), Name: swag.String("image2")}}
	expectedTier = configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      orc8r.UpgradeTierEntityType, Key: "tier1",
		Name:   "new name 1",
		Config: tier,
		Associations: []storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g1"},
			{Type: orc8r.MagmadGatewayType, Key: "g2"},
		},
		GraphID: "2",
		Version: 5,
	}
	actualTier, err = configurator.LoadEntity(
		"n1", orc8r.UpgradeTierEntityType, "tier1",
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadConfig: true, LoadMetadata: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedTier, actualTier)

	// delete err
	tc = tests.Test{
		Method:         "DELETE",
		URL:            manageTiers + "/images/:image_name",
		Handler:        deleteTierImage,
		ParamNames:     []string{"network_id", "tier_id", "image_name"},
		ParamValues:    []string{"n1", "tier1", "image1"},
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case delete
	tc = tests.Test{
		Method:         "DELETE",
		URL:            manageTiers + "/images/:image_name",
		Handler:        deleteTierImage,
		ParamNames:     []string{"network_id", "tier_id", "image_name"},
		ParamValues:    []string{"n1", "tier1", "image2"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	tier.Images = models.TierImages{}
	expectedTier = configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      orc8r.UpgradeTierEntityType, Key: "tier1",
		Name:   "new name 1",
		Config: tier,
		Associations: []storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g1"},
			{Type: orc8r.MagmadGatewayType, Key: "g2"},
		},
		GraphID: "2",
		Version: 6,
	}
	actualTier, err = configurator.LoadEntity(
		"n1", orc8r.UpgradeTierEntityType, "tier1",
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadConfig: true, LoadMetadata: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedTier, actualTier)

	// fail to add a non-existent gw
	tc = tests.Test{
		Method:         "POST",
		URL:            manageTiers + "/gateways",
		Payload:        tests.JSONMarshaler("g4"),
		Handler:        createTierGateway,
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		ExpectedStatus: 500,
		ExpectedError:  "could not find entities matching [type:\"magmad_gateway\" key:\"g4\" ]",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case add gateway
	tc = tests.Test{
		Method:         "POST",
		URL:            manageTiers + "/gateways",
		Payload:        tests.JSONMarshaler("g3"),
		Handler:        createTierGateway,
		ParamNames:     []string{"network_id", "tier_id"},
		ParamValues:    []string{"n1", "tier1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expectedTier = configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      orc8r.UpgradeTierEntityType, Key: "tier1",
		Name:   "new name 1",
		Config: tier,
		Associations: []storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g1"},
			{Type: orc8r.MagmadGatewayType, Key: "g2"},
			{Type: orc8r.MagmadGatewayType, Key: "g3"},
		},
		GraphID: "2",
		Version: 7,
	}
	actualTier, err = configurator.LoadEntity(
		"n1", orc8r.UpgradeTierEntityType, "tier1",
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadConfig: true, LoadMetadata: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedTier, actualTier)

	// happy case remove gateway
	tc = tests.Test{
		Method:         "DELETE",
		URL:            manageTiers + "/gateway/:gateway_id",
		Handler:        deleteTierGateway,
		ParamNames:     []string{"network_id", "tier_id", "gateway_id"},
		ParamValues:    []string{"n1", "tier1", "g3"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expectedTier = configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      orc8r.UpgradeTierEntityType, Key: "tier1",
		Name:   "new name 1",
		Config: tier,
		Associations: []storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g1"},
			{Type: orc8r.MagmadGatewayType, Key: "g2"},
		},
		GraphID: "2",
		Version: 8,
	}
	actualTier, err = configurator.LoadEntity(
		"n1", orc8r.UpgradeTierEntityType, "tier1",
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadConfig: true, LoadMetadata: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedTier, actualTier)
}
