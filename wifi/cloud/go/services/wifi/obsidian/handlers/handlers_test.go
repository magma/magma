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
	"fmt"
	"testing"

	models3 "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/device"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/storage"
	"magma/wifi/cloud/go/serdes"
	"magma/wifi/cloud/go/services/wifi/obsidian/handlers"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	models2 "magma/wifi/cloud/go/services/wifi/obsidian/models"
	"magma/wifi/cloud/go/wifi"

	"github.com/go-openapi/swag"
)

func TestListNetworks(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	listNetworksURL := "/magma/v1/wifi"
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
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	baseNetworksURL := "/magma/v1/wifi"
	createNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, baseNetworksURL, obsidian.POST).HandlerFunc

	payload := models2.NewDefaultWifiNetwork()
	payload.Description = ""

	// test validation
	tc := tests.Test{
		Method:         "POST",
		URL:            baseNetworksURL,
		Payload:        tests.JSONMarshaler(payload),
		Handler:        createNetwork,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"description in body should be at least 1 chars long",
	}
	tests.RunUnitTest(t, e, tc)

	payload = models2.NewDefaultWifiNetwork()
	tc = tests.Test{
		Method:         "POST",
		URL:            baseNetworksURL,
		Payload:        tests.JSONMarshaler(payload),
		Handler:        createNetwork,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	expected := configurator.Network{
		ID:          "n1",
		Type:        wifi.WifiNetworkType,
		Name:        "network_1",
		Description: "Network 1",
		Configs: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
			wifi.WifiNetworkType:        models2.NewDefaultWifiNetworkConfig(),
		},
	}
	actual, err := configurator.LoadNetwork("n1", true, true, serdes.Network)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetNetwork(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	networkURL := "/magma/v1/wifi/:network_id"
	getNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, networkURL, obsidian.GET).HandlerFunc
	nID := "n1"
	expectedN1 := models2.NewDefaultWifiNetwork()

	// Test 404
	tc := tests.Test{
		Method:         "GET",
		URL:            networkURL,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        getNetwork,
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	// Correctly get n1
	seedNetworks(t)
	tc = tests.Test{
		Method:         "GET",
		URL:            networkURL,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        getNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedN1),
	}
	tests.RunUnitTest(t, e, tc)

	// try to get a non-Wifi network
	tc = tests.Test{
		Method:         "GET",
		URL:            networkURL,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        getNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not a <wifi_network> network",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateNetwork(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	networkURL := "/magma/v1/wifi/:network_id"
	updateNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, networkURL, obsidian.PUT).HandlerFunc
	nID := "n1"

	payload := &models2.WifiNetwork{
		ID:          "n1",
		Name:        "updated network_1",
		Description: "Updated Network 1",
		Features: &models.NetworkFeatures{
			Features: map[string]string{
				"feature_1_key": "feature_1_val",
			},
		},
		Wifi: &models2.NetworkWifiConfigs{
			VlAuthServerAddr:         "192.168.1.1",
			VlAuthServerPort:         1234,
			VlAuthServerSharedSecret: "updated ssssh",
			PingHostList:             []string{"172.16.0.1", "www.thefacebook.com"},
			PingNumPackets:           10,
			PingTimeoutSecs:          15,
			XwfRadiusServer:          "radiusnow",
			XwfConfig:                "line 1a\nline 2b",
			XwfDhcpDns1:              "4.8.8.7",
			XwfDhcpDns2:              "8.8.3.3",
			XwfRadiusSharedSecret:    "1231",
			XwfRadiusAuthPort:        2812,
			XwfRadiusAcctPort:        2813,
			XwfUamSecret:             "1233",
			XwfPartnerName:           "xwffcfull",
			MgmtVpnEnabled:           true,
			MgmtVpnProto:             "cows",
			MgmtVpnRemote:            "are still yummy",
			OpenrEnabled:             true,
			AdditionalProps:          map[string]string{"prop1": "val1", "prop2": "val2"},
		},
	}

	// Test 404
	tc := tests.Test{
		Method:         "PUT",
		URL:            networkURL,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        updateNetwork,
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	// seed networks, try to update n1 again
	seedNetworks(t)
	tc.ExpectedStatus = 204
	tests.RunUnitTest(t, e, tc)

	actualN1, err := configurator.LoadNetwork("n1", true, true, serdes.Network)
	assert.NoError(t, err)
	expected := configurator.Network{
		ID:          "n1",
		Type:        wifi.WifiNetworkType,
		Name:        "updated network_1",
		Description: "Updated Network 1",
		Configs: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: payload.Features,
			wifi.WifiNetworkType:        payload.Wifi,
		},
		Version: 1,
	}
	assert.Equal(t, expected, actualN1)

	// update n2, should be 400
	tc = tests.Test{
		Method:         "PUT",
		URL:            networkURL,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        updateNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not a <wifi_network> network",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestDeleteNetwork(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	networkURL := "/magma/v1/wifi/:network_id"
	deleteNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, networkURL, obsidian.DELETE).HandlerFunc
	nID := "n1"

	// Test 404
	tc := tests.Test{
		Method:         "DELETE",
		URL:            networkURL,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        deleteNetwork,
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	// seed networks, delete n1 for real
	seedNetworks(t)
	tc.ExpectedStatus = 204
	tests.RunUnitTest(t, e, tc)

	// delete n1 again, should be 404
	tc.ExpectedStatus = 404
	tests.RunUnitTest(t, e, tc)

	// try to delete n2, should be 400 (not Wifi network)
	tc = tests.Test{
		Method:         "DELETE",
		URL:            networkURL,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        deleteNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not a <wifi_network> network",
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.ListNetworkIDs()
	assert.NoError(t, err)
	assert.Equal(t, []string{"n2"}, actual)
}

func TestPartialUpdateAndGetNetwork(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	baseNetworkURL := "/magma/v1/wifi/:network_id"
	nameURL := fmt.Sprintf("%s/name", baseNetworkURL)
	descriptionURL := fmt.Sprintf("%s/description", baseNetworkURL)
	featuresURL := fmt.Sprintf("%s/features", baseNetworkURL)
	wifiURL := fmt.Sprintf("%s/wifi", baseNetworkURL)
	updateName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, nameURL, obsidian.PUT).HandlerFunc
	updateDescription := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, descriptionURL, obsidian.PUT).HandlerFunc
	updateFeatures := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, featuresURL, obsidian.PUT).HandlerFunc
	updateWifi := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, wifiURL, obsidian.PUT).HandlerFunc
	getName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, nameURL, obsidian.GET).HandlerFunc
	getDescription := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, descriptionURL, obsidian.GET).HandlerFunc
	getFeatures := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, featuresURL, obsidian.GET).HandlerFunc
	getWifi := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, wifiURL, obsidian.GET).HandlerFunc
	nID := "n1"

	updatedName := "updated network_1"
	updatedDescription := "Updated Network 1"
	updatedFeatures := &models.NetworkFeatures{
		Features: map[string]string{
			"feature_1_key": "feature_1_val",
		},
	}
	updatedWifi := &models2.NetworkWifiConfigs{
		VlAuthServerAddr:         "192.168.1.1",
		VlAuthServerPort:         1234,
		VlAuthServerSharedSecret: "updated ssssh",
		PingHostList:             []string{"172.16.0.1", "www.thefacebook.com"},
		PingNumPackets:           10,
		PingTimeoutSecs:          15,
		XwfRadiusServer:          "radiusnow",
		XwfConfig:                "line 1a\nline 2b",
		XwfDhcpDns1:              "4.8.8.7",
		XwfDhcpDns2:              "8.8.3.3",
		XwfRadiusSharedSecret:    "1231",
		XwfRadiusAuthPort:        2812,
		XwfRadiusAcctPort:        2813,
		XwfUamSecret:             "1233",
		XwfPartnerName:           "xwffcfull",
		MgmtVpnEnabled:           true,
		MgmtVpnProto:             "cows",
		MgmtVpnRemote:            "are still yummy",
		OpenrEnabled:             true,
		AdditionalProps:          map[string]string{"prop1": "val1", "prop2": "val2"},
	}

	// Test 404
	tc := tests.Test{
		Method:         "PUT",
		URL:            nameURL,
		Payload:        tests.JSONMarshaler(updatedName),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        updateName,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// Update name
	seedNetworks(t)
	tc = tests.Test{
		Method:         "PUT",
		URL:            nameURL,
		Payload:        tests.JSONMarshaler(updatedName),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        updateName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	actual, err := configurator.LoadNetwork("n1", true, true, serdes.Network)
	assert.NoError(t, err)
	assert.Equal(t, updatedName, actual.Name)
	tc = tests.Test{
		Method:         "GET",
		URL:            nameURL,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        getName,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(updatedName),
	}
	tests.RunUnitTest(t, e, tc)

	// Update description
	tc = tests.Test{
		Method:         "PUT",
		URL:            descriptionURL,
		Payload:        tests.JSONMarshaler(updatedDescription),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        updateDescription,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	actual, err = configurator.LoadNetwork("n1", true, true, serdes.Network)
	assert.NoError(t, err)
	assert.Equal(t, updatedDescription, actual.Description)
	tc = tests.Test{
		Method:         "GET",
		URL:            descriptionURL,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        getDescription,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(updatedDescription),
	}
	tests.RunUnitTest(t, e, tc)

	// Update features
	tc = tests.Test{
		Method:         "PUT",
		URL:            featuresURL,
		Payload:        tests.JSONMarshaler(updatedFeatures),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        updateFeatures,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	actual, err = configurator.LoadNetwork("n1", true, true, serdes.Network)
	assert.NoError(t, err)
	assert.Equal(t, updatedFeatures, actual.Configs[orc8r.NetworkFeaturesConfig].(*models.NetworkFeatures))
	tc = tests.Test{
		Method:         "GET",
		URL:            featuresURL,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        getFeatures,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(updatedFeatures),
	}
	tests.RunUnitTest(t, e, tc)

	// Update wifi
	tc = tests.Test{
		Method:         "PUT",
		URL:            wifiURL,
		Payload:        tests.JSONMarshaler(updatedWifi),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        updateWifi,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	actual, err = configurator.LoadNetwork("n1", true, true, serdes.Network)
	assert.NoError(t, err)
	assert.Equal(t, updatedWifi, actual.Configs[wifi.WifiNetworkType].(*models2.NetworkWifiConfigs))
	tc = tests.Test{
		Method:         "GET",
		URL:            wifiURL,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        getWifi,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(updatedWifi),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestListGateways(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	baseGatewaysUrl := "/magma/v1/wifi/:network_id/gateways"
	obsidianHandlers := handlers.GetHandlers()
	listGateways := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, baseGatewaysUrl, obsidian.GET).HandlerFunc
	nID := "n1"
	gID := "g1"
	mID := "m1"

	// Test 500
	tc := tests.Test{
		Method:         "GET",
		URL:            baseGatewaysUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        listGateways,
		ExpectedError:  "Not found",
		ExpectedStatus: 500,
	}
	tests.RunUnitTest(t, e, tc)

	// Test network with no gateways
	seedNetworks(t)
	tc = tests.Test{
		Method:         "GET",
		URL:            baseGatewaysUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        listGateways,
		ExpectedResult: tests.JSONMarshaler(map[string]models2.WifiGateway{}),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Empty(t, actualEnts)

	// Test network with one gateway
	expectedResult := models2.WifiGateway{
		Name:        "gateway_1",
		Description: "gateway 1",
		ID:          models3.GatewayID(gID),
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
		Tier: "t1",
		Wifi: models2.NewDefaultWifiGatewayConfig(),
	}

	seedGatewaysAndMeshes(t)
	tc = tests.Test{
		Method:         "GET",
		URL:            baseGatewaysUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        listGateways,
		ExpectedResult: tests.JSONMarshaler(map[string]models2.WifiGateway{gID: expectedResult}),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	expectedEnts := configurator.NetworkEntities{
		{
			NetworkID:   nID,
			Type:        orc8r.MagmadGatewayType,
			Key:         gID,
			Name:        "gateway_1",
			Description: "gateway 1",
			PhysicalID:  "hw1",
			GraphID:     "10",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Associations: []storage.TypeAndKey{{Type: wifi.WifiGatewayType, Key: gID}},
			ParentAssociations: []storage.TypeAndKey{
				{Type: wifi.MeshEntityType, Key: mID},
				{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
			},
		},
		{
			NetworkID:          nID,
			Type:               wifi.WifiGatewayType,
			Key:                gID,
			Name:               "gateway_1",
			Description:        "gateway 1",
			Config:             models2.NewDefaultWifiGatewayConfig(),
			GraphID:            "10",
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
		},
	}
	actualEnts, _, err = configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: gID},
			{Type: wifi.WifiGatewayType, Key: gID},
		},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnts, actualEnts)
}

func TestCreateGateway(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	baseGatewaysUrl := "/magma/v1/wifi/:network_id/gateways"
	obsidianHandlers := handlers.GetHandlers()
	createGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, baseGatewaysUrl, obsidian.POST).HandlerFunc
	nID := "n1"
	gID := "g1"
	mID := "m1"

	// Initially empty
	seedNetworks(t)
	actualEnts, _, err := configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Empty(t, actualEnts)

	// Test missing payload
	tc := tests.Test{
		Method:         "POST",
		URL:            baseGatewaysUrl,
		Payload:        &models2.MutableWifiGateway{},
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        createGateway,
		ExpectedStatus: 400,
		ExpectedError:  "validation failure list:\ndescription in body should be at least 1 chars long\ndevice in body is required\nid in body should be at least 1 chars long\nmagmad in body is required\nname in body should be at least 1 chars long\ntier in body should match '^[a-z][\\da-z_]+$'\nwifi in body is required",
	}
	tests.RunUnitTest(t, e, tc)

	// Test post new gateway
	seedPreGateway(t)
	payload := models2.NewDefaultWifiGateway()
	tc = tests.Test{
		Method:         "POST",
		URL:            baseGatewaysUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        createGateway,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	expectedEnts := configurator.NetworkEntities{
		{
			NetworkID:   nID,
			Type:        orc8r.MagmadGatewayType,
			Key:         gID,
			Name:        "gateway_1",
			Description: "gateway 1",
			PhysicalID:  "hw1",
			GraphID:     "2",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Associations:       []storage.TypeAndKey{{Type: wifi.WifiGatewayType, Key: gID}},
			ParentAssociations: []storage.TypeAndKey{{Type: wifi.MeshEntityType, Key: mID}, {Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			Version:            1,
		},
		{
			NetworkID: nID,
			Type:      wifi.MeshEntityType, Key: mID,
			GraphID:      "2",
			Name:         "mesh_1",
			Config:       models2.NewDefaultMeshWifiConfigs(),
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:      1,
		},
		{
			NetworkID: nID,
			Type:      orc8r.UpgradeTierEntityType, Key: "t1",
			GraphID:      "2",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:      1,
		},
		{
			NetworkID:          nID,
			Type:               wifi.WifiGatewayType,
			Key:                gID,
			Name:               "gateway_1",
			Description:        "gateway 1",
			GraphID:            "2",
			Config:             models2.NewDefaultWifiGatewayConfig(),
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
		},
	}
	actualEnts, _, err = configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnts, actualEnts)

	// Can't register the same device
	tc = tests.Test{
		Method:         "POST",
		URL:            baseGatewaysUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        createGateway,
		ExpectedStatus: 400,
		ExpectedError:  "device hw1 is already mapped to gateway g1",
	}
	tests.RunUnitTest(t, e, tc)
	// Can't create a new gateway with a mesh that doesn't exist
	payload.Wifi.MeshID = "DNE"
	tc = tests.Test{
		Method:         "POST",
		URL:            baseGatewaysUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        createGateway,
		ExpectedStatus: 400,
		ExpectedError:  "device hw1 is already mapped to gateway g1",
	}
	tests.RunUnitTest(t, e, tc)
	// Or an empty mesh
	payload.Wifi.MeshID = ""
	tc = tests.Test{
		Method:         "POST",
		URL:            baseGatewaysUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        createGateway,
		ExpectedStatus: 400,
		ExpectedError:  "device hw1 is already mapped to gateway g1",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetGateways(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	getGatewaysUrl := "/magma/v1/wifi/:network_id/gateways/:gateway_id"
	obsidianHandlers := handlers.GetHandlers()
	getGateways := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, getGatewaysUrl, obsidian.GET).HandlerFunc
	nID := "n1"
	gID := "g1"

	// Test network with no gateways
	seedNetworks(t)
	tc := tests.Test{
		Method:         "GET",
		URL:            getGatewaysUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        getGateways,
		ExpectedError:  "Not Found",
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	// Test correct get
	seedGatewaysAndMeshes(t)
	expectedResult := models2.WifiGateway{
		Name:        "gateway_1",
		Description: "gateway 1",
		ID:          models3.GatewayID(gID),
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
		Tier: "t1",
		Wifi: models2.NewDefaultWifiGatewayConfig(),
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            getGatewaysUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        getGateways,
		ExpectedResult: tests.JSONMarshaler(expectedResult),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateGateway(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	e := echo.New()

	updateGatewaysUrl := "/magma/v1/wifi/:network_id/gateways/:gateway_id"
	handlers := handlers.GetHandlers()
	updateGateway := tests.GetHandlerByPathAndMethod(t, handlers, updateGatewaysUrl, obsidian.PUT).HandlerFunc
	nID := "n1"
	gID := "g1"
	mID := "m1"
	nmID := "not_m1"

	// Test network with no gateways
	seedNetworks(t)
	updatedName := "updated_gateway_1"
	updatedDescription := "updated gateway 1"
	updatedMagmad := &models.MagmadGatewayConfigs{
		AutoupgradeEnabled:      swag.Bool(true),
		AutoupgradePollInterval: 500,
		CheckinInterval:         30,
		CheckinTimeout:          10,
	}
	updatedConfig := &models2.GatewayWifiConfigs{
		AdditionalProps: map[string]string{
			"gwprop1": "gwvalue1",
			"gwprop2": "gwvalue2",
			"newprop": "newvalue",
		},
		ClientChannel:                 "12",
		Info:                          "UpdatedGatewayInfo",
		IsProduction:                  false,
		Latitude:                      -90.0000,
		Longitude:                     0.0000,
		MeshID:                        models2.MeshID(nmID),
		MeshRssiThreshold:             -81,
		OverridePassword:              "stillpassword",
		OverrideSsid:                  "StillSuperFastWifiNetwork",
		OverrideXwfConfig:             "updated xwf config",
		OverrideXwfDhcpDns1:           "8.8.8.8",
		OverrideXwfDhcpDns2:           "8.8.4.4",
		OverrideXwfEnabled:            false,
		OverrideXwfPartnerName:        "xwfcfull",
		OverrideXwfRadiusAcctPort:     1813,
		OverrideXwfRadiusAuthPort:     1812,
		OverrideXwfRadiusServer:       "gradius.example.com",
		OverrideXwfRadiusSharedSecret: "xwfisgood",
		OverrideXwfUamSecret:          "theuamsecret",
		UseOverrideSsid:               false,
		UseOverrideXwf:                false,
		WifiDisabled:                  false,
	}
	payload := models2.WifiGateway{
		Name:        models3.GatewayName(updatedName),
		Description: models3.GatewayDescription(updatedDescription),
		ID:          models3.GatewayID(gID),
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		Magmad: updatedMagmad,
		Tier:   "t1",
		Wifi:   updatedConfig,
	}
	tc := tests.Test{
		Method:         "PUT",
		URL:            updateGatewaysUrl,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        updateGateway,
		ExpectedError:  "Not Found",
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	// Can't assign gatewaay to missing mesh
	seedGatewaysAndMeshes(t)
	tc = tests.Test{
		Method:         "PUT",
		URL:            updateGatewaysUrl,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        updateGateway,
		ExpectedError:  "failed to load entity being updated: expected to load 1 ent for update, got 0",
		ExpectedStatus: 500,
	}
	tests.RunUnitTest(t, e, tc)

	// Correctly update gateway
	_, err := configurator.CreateEntities(
		nID,
		[]configurator.NetworkEntity{
			{
				Type: wifi.MeshEntityType, Key: nmID,
				Name:         "not_mesh_1",
				Config:       models2.NewDefaultMeshWifiConfigs(),
				Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	tc = tests.Test{
		Method:         "PUT",
		URL:            updateGatewaysUrl,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        updateGateway,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expectedEnts := configurator.NetworkEntities{
		{
			NetworkID:          nID,
			Type:               orc8r.MagmadGatewayType,
			Key:                gID,
			Name:               updatedName,
			Description:        updatedDescription,
			PhysicalID:         "hw1",
			GraphID:            "10",
			Config:             updatedMagmad,
			Associations:       []storage.TypeAndKey{{Type: wifi.WifiGatewayType, Key: gID}},
			ParentAssociations: []storage.TypeAndKey{{Type: wifi.MeshEntityType, Key: nmID}, {Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			Version:            1,
		},
		{
			NetworkID: nID,
			Type:      wifi.MeshEntityType, Key: mID,
			Name:    "mesh_1",
			GraphID: "14",
			Config:  models2.NewDefaultMeshWifiConfigs(),
			Version: 1,
		},
		{
			NetworkID: nID,
			Type:      wifi.MeshEntityType, Key: nmID,
			Name:         "not_mesh_1",
			GraphID:      "10",
			Config:       models2.NewDefaultMeshWifiConfigs(),
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:      1,
		},
		{
			NetworkID: nID,
			Type:      orc8r.UpgradeTierEntityType, Key: "t1",
			GraphID:      "10",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:      0,
		},
		{
			NetworkID: nID,
			Type:      orc8r.UpgradeTierEntityType, Key: "t2",
			GraphID: "8",
			Version: 0,
		},
		{
			NetworkID:          nID,
			Type:               wifi.WifiGatewayType,
			Key:                gID,
			Name:               updatedName,
			Description:        updatedDescription,
			GraphID:            "10",
			Config:             updatedConfig,
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:            1,
		},
	}
	actualEnts, _, err := configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnts, actualEnts)
}

func TestDeleteGateway(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	e := echo.New()

	deleteGatewayURL := "/magma/v1/wifi/:network_id/gateways/:gateway_id"
	handlers := handlers.GetHandlers()
	deleteGateway := tests.GetHandlerByPathAndMethod(t, handlers, deleteGatewayURL, obsidian.DELETE).HandlerFunc
	nID := "n1"
	gID := "g1"
	mID := "m1"

	// Test delete missing gateway - it succeeds on a missing gateway,
	// which makes sense
	seedNetworks(t)
	tc := tests.Test{
		Method:         "DELETE",
		URL:            deleteGatewayURL,
		Payload:        nil,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        deleteGateway,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Correctly delete gateway
	seedGatewaysAndMeshes(t)
	tc = tests.Test{
		Method:         "DELETE",
		URL:            deleteGatewayURL,
		Payload:        nil,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        deleteGateway,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: gID},
			{Type: wifi.WifiGatewayType, Key: gID},
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
			{Type: wifi.MeshEntityType, Key: mID},
		},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)

	_, err = device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hw1", serdes.Device)
	assert.EqualError(t, err, "Not found")

	expectedEnts := configurator.NetworkEntities{
		{
			NetworkID: nID,
			Type:      wifi.MeshEntityType, Key: mID,
			Name:    "mesh_1",
			GraphID: "10",
			Config:  models2.NewDefaultMeshWifiConfigs(),
			Version: 0,
		},
		{NetworkID: "n1", Type: orc8r.UpgradeTierEntityType, Key: "t1", GraphID: "12"},
	}
	assert.Equal(t, expectedEnts, actualEnts)
}

func TestPartialUpdateAndGetGateway(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	gatewayUrl := "/magma/v1/wifi/:network_id/gateways/:gateway_id"
	nameUrl := fmt.Sprintf("%s/name", gatewayUrl)
	descriptionUrl := fmt.Sprintf("%s/description", gatewayUrl)
	magmadUrl := fmt.Sprintf("%s/magmad", gatewayUrl)
	tierUrl := fmt.Sprintf("%s/tier", gatewayUrl)
	deviceUrl := fmt.Sprintf("%s/device", gatewayUrl)
	wifiUrl := fmt.Sprintf("%s/wifi", gatewayUrl)
	obsidianHandlers := handlers.GetHandlers()
	updateName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, nameUrl, obsidian.PUT).HandlerFunc
	updateDescription := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, descriptionUrl, obsidian.PUT).HandlerFunc
	updateMagmad := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, magmadUrl, obsidian.PUT).HandlerFunc
	updateTier := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, tierUrl, obsidian.PUT).HandlerFunc
	updateDevice := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, deviceUrl, obsidian.PUT).HandlerFunc
	updateWifi := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, wifiUrl, obsidian.PUT).HandlerFunc
	getName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, nameUrl, obsidian.GET).HandlerFunc
	getDescription := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, descriptionUrl, obsidian.GET).HandlerFunc
	getMagmad := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, magmadUrl, obsidian.GET).HandlerFunc
	getTier := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, tierUrl, obsidian.GET).HandlerFunc
	getDevice := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, deviceUrl, obsidian.GET).HandlerFunc
	getWifi := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, wifiUrl, obsidian.GET).HandlerFunc

	nID := "n1"
	gID := "g1"
	mID := "m1"
	nmID := "not_m1"

	seedNetworks(t)
	seedGatewaysAndMeshes(t)
	_, err := configurator.CreateEntities(
		nID,
		[]configurator.NetworkEntity{
			{
				Type: wifi.MeshEntityType, Key: nmID,
				Name:   "not_mesh_1",
				Config: models2.NewDefaultMeshWifiConfigs(),
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	expectedEnts := map[string]*configurator.NetworkEntity{
		"magmad": {
			NetworkID:   nID,
			Type:        orc8r.MagmadGatewayType,
			Key:         gID,
			Name:        "gateway_1",
			Description: "gateway 1",
			PhysicalID:  "hw1",
			GraphID:     "10",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Associations:       []storage.TypeAndKey{{Type: wifi.WifiGatewayType, Key: gID}},
			ParentAssociations: []storage.TypeAndKey{{Type: wifi.MeshEntityType, Key: mID}, {Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			Version:            0,
		},
		"mesh": {
			NetworkID: nID,
			Type:      wifi.MeshEntityType, Key: mID,
			Name:         "mesh_1",
			GraphID:      "10",
			Config:       models2.NewDefaultMeshWifiConfigs(),
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:      0,
		},
		"mesh2": {
			NetworkID: nID,
			Type:      wifi.MeshEntityType, Key: nmID,
			Name:    "not_mesh_1",
			GraphID: "12",
			Config:  models2.NewDefaultMeshWifiConfigs(),
			Version: 0,
		},
		"tier1": {
			NetworkID: nID,
			Type:      orc8r.UpgradeTierEntityType, Key: "t1",
			GraphID:      "10",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:      0,
		},
		"tier2": {
			NetworkID: nID,
			Type:      orc8r.UpgradeTierEntityType, Key: "t2",
			GraphID: "8",
			Version: 0,
		},
		"gateway": {
			NetworkID:          nID,
			Type:               wifi.WifiGatewayType,
			Key:                gID,
			Name:               "gateway_1",
			Description:        "gateway 1",
			GraphID:            "10",
			Config:             models2.NewDefaultWifiGatewayConfig(),
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:            0,
		},
	}

	// Test updating gateway name

	// Update the ents that we expect and then convert them into a list so we can
	// compare them to what we get from configurator.LoadEntities later
	updatedName := "updated_gateway_name"
	expectedEnts["magmad"].Name = updatedName
	expectedEnts["magmad"].Version++
	expectedEntsVals := make(configurator.NetworkEntities, 0, len(expectedEnts))
	// Key order matters to compare later
	key_order := []string{"magmad", "mesh", "mesh2", "tier1", "tier2", "gateway"}
	for _, v := range key_order {
		expectedEntsVals = append(expectedEntsVals, *expectedEnts[v])
	}

	payload := tests.JSONMarshaler(updatedName)
	tc := tests.Test{
		Method:         "PUT",
		URL:            nameUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        updateName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err := configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	tc = tests.Test{
		Method:         "GET",
		URL:            nameUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        getName,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test updating gateway description
	updatedDescription := "updated_description"
	// expectedEnts["gateway"].Description = updatedDescription
	// expectedEnts["gateway"].Version++
	expectedEnts["magmad"].Description = updatedDescription
	expectedEnts["magmad"].Version++
	expectedEntsVals = make(configurator.NetworkEntities, 0, len(expectedEnts))
	for _, v := range key_order {
		expectedEntsVals = append(expectedEntsVals, *expectedEnts[v])
	}

	payload = tests.JSONMarshaler(updatedDescription)
	tc = tests.Test{
		Method:         "PUT",
		URL:            descriptionUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        updateDescription,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err = configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	tc = tests.Test{
		Method:         "GET",
		URL:            descriptionUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        getDescription,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test updating gateway magmad
	updatedMagmad := &models.MagmadGatewayConfigs{
		AutoupgradeEnabled:      swag.Bool(false),
		AutoupgradePollInterval: 100,
		CheckinInterval:         30,
		CheckinTimeout:          10,
	}
	expectedEnts["magmad"].Config = updatedMagmad
	expectedEnts["magmad"].Version++
	expectedEntsVals = make(configurator.NetworkEntities, 0, len(expectedEnts))
	for _, v := range key_order {
		expectedEntsVals = append(expectedEntsVals, *expectedEnts[v])
	}

	payload = tests.JSONMarshaler(updatedMagmad)
	tc = tests.Test{
		Method:         "PUT",
		URL:            magmadUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        updateMagmad,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err = configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	tc = tests.Test{
		Method:         "GET",
		URL:            magmadUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        getMagmad,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test updating gateway tier
	updatedGatewayTier := "t2"
	expectedEnts["magmad"].ParentAssociations = []storage.TypeAndKey{{Type: wifi.MeshEntityType, Key: mID}, {Type: orc8r.UpgradeTierEntityType, Key: updatedGatewayTier}}
	expectedEnts["tier1"].Associations = nil
	expectedEnts["tier1"].GraphID = "13"
	expectedEnts["tier1"].Version++
	expectedEnts["tier2"].Associations = []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}}
	expectedEnts["tier2"].GraphID = "10"
	expectedEnts["tier2"].Version++
	expectedEntsVals = make(configurator.NetworkEntities, 0, len(expectedEnts))
	for _, v := range key_order {
		expectedEntsVals = append(expectedEntsVals, *expectedEnts[v])
	}

	payload = tests.JSONMarshaler(updatedGatewayTier)
	tc = tests.Test{
		Method:         "PUT",
		URL:            tierUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        updateTier,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err = configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	tc = tests.Test{
		Method:         "GET",
		URL:            tierUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        getTier,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test updating gateway device
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
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        updateDevice,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Configurator ents should not have changed
	actualEnts, _, err = configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	// But the HardwareID of the physical device with id "hw1" should have updated
	expectedDevice := &models.GatewayDevice{HardwareID: "hw2", Key: &models.ChallengeKey{KeyType: "ECHO"}}
	actualDevice, err := device.GetDevice(nID, orc8r.AccessGatewayRecordType, "hw1", serdes.Device)
	assert.NoError(t, err)
	assert.Equal(t, expectedDevice, actualDevice)

	tc = tests.Test{
		Method:         "GET",
		URL:            deviceUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        getDevice,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test updating gateway magmad
	updatedWifi := &models2.GatewayWifiConfigs{
		AdditionalProps: map[string]string{
			"gwprop1": "gwvalue1",
			"gwprop2": "gwvalue2",
			"newprop": "newvalue",
		},
		ClientChannel:                 "12",
		Info:                          "UpdatedGatewayInfo",
		IsProduction:                  false,
		Latitude:                      -90.0000,
		Longitude:                     0.0000,
		MeshID:                        "missing_mesh",
		MeshRssiThreshold:             -81,
		OverridePassword:              "stillpassword",
		OverrideSsid:                  "StillSuperFastWifiNetwork",
		OverrideXwfConfig:             "updated xwf config",
		OverrideXwfDhcpDns1:           "8.8.8.8",
		OverrideXwfDhcpDns2:           "8.8.4.4",
		OverrideXwfEnabled:            false,
		OverrideXwfPartnerName:        "xwfcfull",
		OverrideXwfRadiusAcctPort:     1813,
		OverrideXwfRadiusAuthPort:     1812,
		OverrideXwfRadiusServer:       "gradius.example.com",
		OverrideXwfRadiusSharedSecret: "xwfisgood",
		OverrideXwfUamSecret:          "theuamsecret",
		UseOverrideSsid:               false,
		UseOverrideXwf:                false,
		WifiDisabled:                  false,
	}
	expectedEnts["gateway"].Config = updatedWifi
	expectedEnts["gateway"].Version++
	expectedEnts["magmad"].ParentAssociations = []storage.TypeAndKey{{Type: wifi.MeshEntityType, Key: nmID}, {Type: orc8r.UpgradeTierEntityType, Key: updatedGatewayTier}}
	expectedEnts["mesh"].Associations = nil
	expectedEnts["mesh"].GraphID = "15"
	expectedEnts["mesh"].Version++
	expectedEnts["mesh2"].GraphID = "10"
	expectedEnts["mesh2"].Version++
	expectedEnts["mesh2"].Associations = []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}}
	expectedEntsVals = make(configurator.NetworkEntities, 0, len(expectedEnts))
	for _, v := range key_order {
		expectedEntsVals = append(expectedEntsVals, *expectedEnts[v])
	}
	// Can't update wifi configs to missing mesh
	payload = tests.JSONMarshaler(updatedWifi)
	tc = tests.Test{
		Method:         "PUT",
		URL:            wifiUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        updateWifi,
		ExpectedError:  "failed to load entity being updated: expected to load 1 ent for update, got 0",
		ExpectedStatus: 500,
	}
	tests.RunUnitTest(t, e, tc)

	// Correctly update wifi configs
	updatedWifi.MeshID = models2.MeshID(nmID)
	tc = tests.Test{
		Method:         "PUT",
		URL:            wifiUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        updateWifi,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err = configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	tc = tests.Test{
		Method:         "GET",
		URL:            wifiUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        getWifi,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
}

func TestListMeshes(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	baseMeshesUrl := "/magma/v1/wifi/:network_id/meshes"
	obsidianHandlers := handlers.GetHandlers()
	listMeshes := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, baseMeshesUrl, obsidian.GET).HandlerFunc
	nID := "n1"
	gID := "g1"
	mID := "m1"

	// Test 404
	tc := tests.Test{
		Method:         "GET",
		URL:            baseMeshesUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        listMeshes,
		ExpectedError:  "Not found",
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	// Test network with no gateways
	seedNetworks(t)
	tc = tests.Test{
		Method:         "GET",
		URL:            baseMeshesUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        listMeshes,
		ExpectedResult: tests.JSONMarshaler([]string{}),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Empty(t, actualEnts)

	// List the meshes correctly
	seedGatewaysAndMeshes(t)
	tc = tests.Test{
		Method:         "GET",
		URL:            baseMeshesUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        listMeshes,
		ExpectedResult: tests.JSONMarshaler([]string{mID}),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	expectedEnts := configurator.NetworkEntities{
		{
			NetworkID: nID,
			Type:      wifi.MeshEntityType, Key: mID,
			Name:         "mesh_1",
			Config:       models2.NewDefaultMeshWifiConfigs(),
			GraphID:      "10",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
		},
	}
	actualEnts, _, err = configurator.LoadEntities(
		"n1", swag.String(wifi.MeshEntityType), nil, nil,
		nil,
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnts, actualEnts)
}

func TestCreateMesh(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	baseMeshesUrl := "/magma/v1/wifi/:network_id/meshes"
	obsidianHandlers := handlers.GetHandlers()
	createMesh := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, baseMeshesUrl, obsidian.POST).HandlerFunc
	nID := "n1"
	gID := "g1"
	mID := "m1"

	// Initially empty
	seedNetworks(t)
	actualEnts, _, err := configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Empty(t, actualEnts)

	// Test missing payload
	tc := tests.Test{
		Method:         "POST",
		URL:            baseMeshesUrl,
		Payload:        &models2.WifiMesh{},
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        createMesh,
		ExpectedStatus: 400,
		ExpectedError:  "validation failure list:\nconfig in body is required\ngateway_ids in body is required\nid in body should be at least 1 chars long\nname in body should be at least 1 chars long",
	}
	tests.RunUnitTest(t, e, tc)

	// Test post new gateway
	seedPreMesh(t)
	payload := models2.NewDefaultWifiMesh()
	tc = tests.Test{
		Method:         "POST",
		URL:            baseMeshesUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        createMesh,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	expectedEnts := configurator.NetworkEntities{
		{
			NetworkID:   nID,
			Type:        orc8r.MagmadGatewayType,
			Key:         gID,
			Name:        "gateway_1",
			Description: "gateway 1",
			PhysicalID:  "hw1",
			GraphID:     "2",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Associations:       []storage.TypeAndKey{{Type: wifi.WifiGatewayType, Key: gID}},
			ParentAssociations: []storage.TypeAndKey{{Type: wifi.MeshEntityType, Key: mID}, {Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			Version:            0,
		},
		{
			NetworkID: nID,
			Type:      wifi.MeshEntityType, Key: mID,
			Name:         "mesh_1",
			Config:       models2.NewDefaultMeshWifiConfigs(),
			GraphID:      "2",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
		},
		{
			NetworkID: nID,
			Type:      orc8r.UpgradeTierEntityType, Key: "t1",
			GraphID:      "2",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:      0,
		},
		{
			NetworkID:          nID,
			Type:               wifi.WifiGatewayType,
			Key:                gID,
			Name:               "gateway_1",
			Description:        "gateway 1",
			GraphID:            "2",
			Config:             models2.NewDefaultWifiGatewayConfig(),
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
		},
	}
	actualEnts, _, err = configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnts, actualEnts)

	// Can't register the same device
	tc = tests.Test{
		Method:         "POST",
		URL:            baseMeshesUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{nID},
		Handler:        createMesh,
		ExpectedStatus: 500,
		ExpectedError:  "an entity 'mesh-m1' already exists",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetMeshes(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	getMeshesUrl := "/magma/v1/wifi/:network_id/meshes/:mesh_id"
	obsidianHandlers := handlers.GetHandlers()
	getMeshes := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, getMeshesUrl, obsidian.GET).HandlerFunc
	nID := "n1"
	mID := "m1"

	// Test network with no meshes
	seedNetworks(t)
	tc := tests.Test{
		Method:         "GET",
		URL:            getMeshesUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id", "mesh_id"},
		ParamValues:    []string{nID, mID},
		Handler:        getMeshes,
		ExpectedError:  "Not Found",
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	// Test correct get
	seedGatewaysAndMeshes(t)
	expectedResult := models2.WifiMesh{
		ID:         "m1",
		Name:       models2.MeshName("mesh_1"),
		Config:     models2.NewDefaultMeshWifiConfigs(),
		GatewayIds: []models3.GatewayID{"g1"},
	}

	tc = tests.Test{
		Method:         "GET",
		URL:            getMeshesUrl,
		Payload:        nil,
		ParamNames:     []string{"network_id", "mesh_id"},
		ParamValues:    []string{nID, mID},
		Handler:        getMeshes,
		ExpectedResult: tests.JSONMarshaler(expectedResult),
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateMesh(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	e := echo.New()

	updateMeshesUrl := "/magma/v1/wifi/:network_id/meshes/:mesh_id"
	handlers := handlers.GetHandlers()
	updateMesh := tests.GetHandlerByPathAndMethod(t, handlers, updateMeshesUrl, obsidian.PUT).HandlerFunc
	nID := "n1"
	gID := "g1"
	mID := "m1"

	// Test network with no meshes
	seedNetworks(t)
	updatedName := "updated_mesh_1"
	updatedConfig := &models2.MeshWifiConfigs{
		AdditionalProps: map[string]string{
			"updated_mesh_prop1": "updated_mesh_value1",
			"updated_mesh_prop2": "updated_mesh_value2",
		},
		MeshChannelType: "HT20",
		MeshFrequency:   5825,
		MeshSsid:        "mesh_ssid",
		Password:        "password",
		Ssid:            "ssid",
		VlSsid:          "vl_ssid",
		XwfEnabled:      false,
	}
	updatedGateways := []models3.GatewayID{}

	payload := models2.WifiMesh{
		Name:       models2.MeshName(updatedName),
		ID:         models2.MeshID(mID),
		Config:     updatedConfig,
		GatewayIds: updatedGateways,
	}
	tc := tests.Test{
		Method:         "PUT",
		URL:            updateMeshesUrl,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "mesh_id"},
		ParamValues:    []string{nID, mID},
		Handler:        updateMesh,
		ExpectedError:  "Not Found",
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	// Then make sure we can't update the gateways list
	seedGatewaysAndMeshes(t)
	tc = tests.Test{
		Method:         "PUT",
		URL:            updateMeshesUrl,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "mesh_id"},
		ParamValues:    []string{nID, mID},
		Handler:        updateMesh,
		ExpectedError:  "can't update gateways here! please update the individual gateways instead.",
		ExpectedStatus: 400,
	}
	tests.RunUnitTest(t, e, tc)

	// Then update the mesh correctly
	updatedGateways = []models3.GatewayID{models3.GatewayID(gID)}
	payload.GatewayIds = updatedGateways

	tc = tests.Test{
		Method:         "PUT",
		URL:            updateMeshesUrl,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "mesh_id"},
		ParamValues:    []string{nID, mID},
		Handler:        updateMesh,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expectedEnts := configurator.NetworkEntities{
		{
			NetworkID:   nID,
			Type:        orc8r.MagmadGatewayType,
			Key:         gID,
			Name:        "gateway_1",
			Description: "gateway 1",
			PhysicalID:  "hw1",
			GraphID:     "10",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Associations: []storage.TypeAndKey{{Type: wifi.WifiGatewayType, Key: gID}},
			ParentAssociations: []storage.TypeAndKey{
				{Type: wifi.MeshEntityType, Key: mID},
				{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
			},
			Version: 0,
		},
		{
			NetworkID: nID,
			Type:      wifi.MeshEntityType, Key: mID,
			Name:         updatedName,
			Config:       updatedConfig,
			GraphID:      "10",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:      1,
		},
		{
			NetworkID: nID,
			Type:      orc8r.UpgradeTierEntityType, Key: "t1",
			GraphID:      "10",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:      0,
		},
		{
			NetworkID: nID,
			Type:      orc8r.UpgradeTierEntityType, Key: "t2",
			GraphID: "8",
			Version: 0,
		},
		{
			NetworkID:          nID,
			Type:               wifi.WifiGatewayType,
			Key:                gID,
			Name:               "gateway_1",
			Description:        "gateway 1",
			GraphID:            "10",
			Config:             models2.NewDefaultWifiGatewayConfig(),
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
		},
	}
	actualEnts, _, err := configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnts, actualEnts)
}

func TestDeleteMesh(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	e := echo.New()

	deleteMeshURL := "/magma/v1/wifi/:network_id/meshes/:mesh_id"
	wifiUrl := "/magma/v1/wifi/:network_id/gateways/:gateway_id/wifi"
	handlers := handlers.GetHandlers()
	deleteMesh := tests.GetHandlerByPathAndMethod(t, handlers, deleteMeshURL, obsidian.DELETE).HandlerFunc
	updateWifi := tests.GetHandlerByPathAndMethod(t, handlers, wifiUrl, obsidian.PUT).HandlerFunc

	nID := "n1"
	gID := "g1"
	mID := "m1"
	nmID := "not_m1"

	// Test delete missing mesh
	seedNetworks(t)
	tc := tests.Test{
		Method:         "DELETE",
		URL:            deleteMeshURL,
		Payload:        nil,
		ParamNames:     []string{"network_id", "mesh_id"},
		ParamValues:    []string{nID, mID},
		Handler:        deleteMesh,
		ExpectedError:  "Not Found",
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	// Can't delete a mesh if it has gateways
	seedGatewaysAndMeshes(t)
	tc = tests.Test{
		Method:         "DELETE",
		URL:            deleteMeshURL,
		Payload:        nil,
		ParamNames:     []string{"network_id", "mesh_id"},
		ParamValues:    []string{nID, mID},
		Handler:        deleteMesh,
		ExpectedError:  "can't delete a mesh with gateways!",
		ExpectedStatus: 400,
	}
	tests.RunUnitTest(t, e, tc)

	// Disassociate gateway then delete mesh
	_, err := configurator.CreateEntities(
		nID,
		[]configurator.NetworkEntity{
			{
				Type: wifi.MeshEntityType, Key: nmID,
				Name:   "not_mesh_1",
				Config: models2.NewDefaultMeshWifiConfigs(),
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	payload := models2.GatewayWifiConfigs{
		AdditionalProps: map[string]string{
			"gwprop1": "gwvalue1",
			"gwprop2": "gwvalue2",
		},
		ClientChannel:                 "11",
		Info:                          "GatewayInfo",
		IsProduction:                  false,
		Latitude:                      37.48497,
		Longitude:                     -122.148284,
		MeshID:                        models2.MeshID(nmID),
		MeshRssiThreshold:             -80,
		OverridePassword:              "password",
		OverrideSsid:                  "SuperFastWifiNetwork",
		OverrideXwfConfig:             "xwf config",
		OverrideXwfDhcpDns1:           "8.8.8.8",
		OverrideXwfDhcpDns2:           "8.8.4.4",
		OverrideXwfEnabled:            false,
		OverrideXwfPartnerName:        "xwfcfull",
		OverrideXwfRadiusAcctPort:     1813,
		OverrideXwfRadiusAuthPort:     1812,
		OverrideXwfRadiusServer:       "gradius.example.com",
		OverrideXwfRadiusSharedSecret: "xwfisgood",
		OverrideXwfUamSecret:          "theuamsecret",
		UseOverrideSsid:               false,
		UseOverrideXwf:                false,
		WifiDisabled:                  false,
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            wifiUrl,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{nID, gID},
		Handler:        updateWifi,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "DELETE",
		URL:            deleteMeshURL,
		Payload:        nil,
		ParamNames:     []string{"network_id", "mesh_id"},
		ParamValues:    []string{nID, mID},
		Handler:        deleteMesh,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: gID},
			{Type: wifi.WifiGatewayType, Key: gID},
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
			{Type: wifi.MeshEntityType, Key: mID},
		},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)

	expectedEnts := configurator.NetworkEntities{
		{
			NetworkID:   nID,
			Type:        orc8r.MagmadGatewayType,
			Key:         gID,
			Name:        "gateway_1",
			Description: "gateway 1",
			PhysicalID:  "hw1",
			GraphID:     "10",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Associations:       []storage.TypeAndKey{{Type: wifi.WifiGatewayType, Key: gID}},
			ParentAssociations: []storage.TypeAndKey{{Type: wifi.MeshEntityType, Key: nmID}, {Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
		},
		{
			NetworkID: nID,
			Type:      orc8r.UpgradeTierEntityType, Key: "t1",
			GraphID:      "10",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:      0,
		},
		{
			NetworkID:          nID,
			Type:               wifi.WifiGatewayType,
			Key:                gID,
			Name:               "gateway_1",
			Description:        "gateway 1",
			GraphID:            "10",
			Config:             &payload,
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:            1,
		},
	}
	assert.Equal(t, expectedEnts, actualEnts)
}

func TestPartialUpdateAndGetMesh(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	e := echo.New()

	meshUrl := "/magma/v1/wifi/:network_id/meshes/:mesh_id"
	nameUrl := fmt.Sprintf("%s/name", meshUrl)
	configUrl := fmt.Sprintf("%s/config", meshUrl)
	obsidianHandlers := handlers.GetHandlers()
	updateName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, nameUrl, obsidian.PUT).HandlerFunc
	updateConfig := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, configUrl, obsidian.PUT).HandlerFunc
	getName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, nameUrl, obsidian.GET).HandlerFunc
	getConfig := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, configUrl, obsidian.GET).HandlerFunc

	nID := "n1"
	gID := "g1"
	mID := "m1"

	seedNetworks(t)
	seedGatewaysAndMeshes(t)

	expectedEnts := map[string]*configurator.NetworkEntity{
		"magmad": {
			NetworkID:   nID,
			Type:        orc8r.MagmadGatewayType,
			Key:         gID,
			Name:        "gateway_1",
			Description: "gateway 1",
			PhysicalID:  "hw1",
			GraphID:     "10",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Associations:       []storage.TypeAndKey{{Type: wifi.WifiGatewayType, Key: gID}},
			ParentAssociations: []storage.TypeAndKey{{Type: wifi.MeshEntityType, Key: mID}, {Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			Version:            0,
		},
		"mesh": {
			NetworkID: nID,
			Type:      wifi.MeshEntityType, Key: mID,
			Name:         "mesh_1",
			GraphID:      "10",
			Config:       models2.NewDefaultMeshWifiConfigs(),
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:      0,
		},
		"tier1": {
			NetworkID: nID,
			Type:      orc8r.UpgradeTierEntityType, Key: "t1",
			GraphID:      "10",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:      0,
		},
		"tier2": {
			NetworkID: nID,
			Type:      orc8r.UpgradeTierEntityType, Key: "t2",
			GraphID: "8",
			Version: 0,
		},
		"gateway": {
			NetworkID:          nID,
			Type:               wifi.WifiGatewayType,
			Key:                gID,
			Name:               "gateway_1",
			Description:        "gateway 1",
			GraphID:            "10",
			Config:             models2.NewDefaultWifiGatewayConfig(),
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			Version:            0,
		},
	}

	// Test updating mesh name

	// Update the ents that we expect and then convert them into a list so we can
	// compare them to what we get from configurator.LoadEntities later
	updatedName := "updated_mesh_name"
	expectedEnts["mesh"].Name = updatedName
	expectedEnts["mesh"].Version++
	expectedEntsVals := make(configurator.NetworkEntities, 0, len(expectedEnts))
	// Key order matters to compare later
	key_order := []string{"magmad", "mesh", "tier1", "tier2", "gateway"}
	for _, v := range key_order {
		expectedEntsVals = append(expectedEntsVals, *expectedEnts[v])
	}

	payload := tests.JSONMarshaler(updatedName)
	tc := tests.Test{
		Method:         "PUT",
		URL:            nameUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "mesh_id"},
		ParamValues:    []string{nID, mID},
		Handler:        updateName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err := configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	tc = tests.Test{
		Method:         "GET",
		URL:            nameUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "mesh_id"},
		ParamValues:    []string{nID, mID},
		Handler:        getName,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	// Test updating mesh magmad
	updatedConfig := &models2.MeshWifiConfigs{
		AdditionalProps: map[string]string{
			"updated_mesh_prop1": "updated_mesh_value1",
			"updated_mesh_prop2": "updated_mesh_value2",
		},
		MeshChannelType: "HT20",
		MeshFrequency:   5825,
		MeshSsid:        "mesh_ssid",
		Password:        "password",
		Ssid:            "ssid",
		VlSsid:          "vl_ssid",
		XwfEnabled:      false,
	}
	expectedEnts["mesh"].Config = updatedConfig
	expectedEnts["mesh"].Version++
	expectedEntsVals = make(configurator.NetworkEntities, 0, len(expectedEnts))
	for _, v := range key_order {
		expectedEntsVals = append(expectedEntsVals, *expectedEnts[v])
	}

	payload = tests.JSONMarshaler(updatedConfig)
	tc = tests.Test{
		Method:         "PUT",
		URL:            configUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "mesh_id"},
		ParamValues:    []string{nID, mID},
		Handler:        updateConfig,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err = configurator.LoadEntities(
		nID, nil, nil, nil,
		[]storage.TypeAndKey{},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntsVals, actualEnts)

	tc = tests.Test{
		Method:         "GET",
		URL:            configUrl,
		Payload:        payload,
		ParamNames:     []string{"network_id", "mesh_id"},
		ParamValues:    []string{nID, mID},
		Handler:        getConfig,
		ExpectedResult: payload,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
}

// n1 is a wifi network, n2 is not
func seedNetworks(t *testing.T) {
	gatewayRecord := &models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}}
	err := device.RegisterDevice("n1", orc8r.AccessGatewayRecordType, "hw1", gatewayRecord, serdes.Device)
	assert.NoError(t, err)

	_, err = configurator.CreateNetworks(
		[]configurator.Network{
			models2.NewDefaultWifiNetwork().ToConfiguratorNetwork(),
			{
				ID:          "n2",
				Type:        "blah",
				Name:        "network_2",
				Description: "Network 2",
				Configs:     map[string]interface{}{},
			},
		},
		serdes.Network,
	)
	assert.NoError(t, err)
}

func seedPreGateway(t *testing.T) {
	// Create Tier necessary for the gateway to be in
	_, err := configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
			},
			{
				Type: wifi.MeshEntityType, Key: "m1",
				Name:   "mesh_1",
				Config: models2.NewDefaultMeshWifiConfigs(),
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
}

func seedPreMesh(t *testing.T) {
	// Create what's necessary to create a mesh
	nID := "n1"
	gID := "g1"

	_, err := configurator.CreateEntities(
		nID,
		[]configurator.NetworkEntity{
			{
				Type: wifi.WifiGatewayType, Key: gID,
				Name:        "gateway_1",
				Description: "gateway 1",
				Config:      models2.NewDefaultWifiGatewayConfig(),
				ParentAssociations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: gID},
				},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: gID,
				Name:        "gateway_1",
				Description: "gateway 1",
				PhysicalID:  "hw1",
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
				Associations:       []storage.TypeAndKey{{Type: wifi.WifiGatewayType, Key: gID}},
				ParentAssociations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: gID},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
}

func seedGatewaysAndMeshes(t *testing.T) {
	nID := "n1"
	gID := "g1"
	mID := "m1"
	_, err := configurator.CreateEntities(
		nID,
		[]configurator.NetworkEntity{
			{
				Type: wifi.WifiGatewayType, Key: gID,
				Name:               "gateway_1",
				Description:        "gateway 1",
				Config:             models2.NewDefaultWifiGatewayConfig(),
				ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: gID,
				Name:        "gateway_1",
				Description: "gateway 1",
				PhysicalID:  "hw1",
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
				Associations: []storage.TypeAndKey{{Type: wifi.WifiGatewayType, Key: gID}},
				ParentAssociations: []storage.TypeAndKey{
					{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
					{Type: wifi.MeshEntityType, Key: mID},
				},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: gID},
				},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t2",
			},
			{
				Type: wifi.MeshEntityType, Key: mID,
				Name:         "mesh_1",
				Config:       models2.NewDefaultMeshWifiConfigs(),
				Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
}
