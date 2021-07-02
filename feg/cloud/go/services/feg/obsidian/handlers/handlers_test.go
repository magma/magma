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
	"context"
	"testing"
	"time"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/serdes"
	"magma/feg/cloud/go/services/feg/obsidian/handlers"
	models2 "magma/feg/cloud/go/services/feg/obsidian/models"
	healthTestInit "magma/feg/cloud/go/services/health/test_init"
	healthTestUtils "magma/feg/cloud/go/services/health/test_utils"
	"magma/lte/cloud/go/lte"
	models3 "magma/lte/cloud/go/services/lte/obsidian/models"
	policyModels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestFederationNetworks(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	testHealthServicer, err := healthTestInit.StartTestService(t)
	assert.NoError(t, err)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	listNetworks := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg", obsidian.GET).HandlerFunc
	createNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg", obsidian.POST).HandlerFunc
	getNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg/:network_id", obsidian.GET).HandlerFunc
	updateNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg/:network_id", obsidian.PUT).HandlerFunc
	deleteNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg/:network_id", obsidian.DELETE).HandlerFunc
	getNetworkFederationConfig := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg/:network_id/federation", obsidian.GET).HandlerFunc
	getNetworkFederationStatus := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg/:network_id/cluster_status", obsidian.GET).HandlerFunc

	// Test ListNetworks
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg",
		Payload:        nil,
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{}),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	seedFederationNetworks(t)

	// Test CreateNetwork
	tc = tests.Test{
		Method: "POST",
		URL:    "/magma/v1/feg",
		Payload: tests.JSONMarshaler(
			&models2.FegNetwork{
				Federation:  models2.NewDefaultNetworkFederationConfigs(),
				Description: "Foo Bar",
				DNS:         models.NewDefaultDNSConfig(),
				Features:    models.NewDefaultFeaturesConfig(),
				ID:          "n4",
				Name:        "foobar",
			},
		),
		Handler:        createNetwork,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Test ListNetworks
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg",
		Payload:        nil,
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n1", "n3", "n4"}),
	}
	tests.RunUnitTest(t, e, tc)

	// Test GetNetwork
	expectedN1 := &models2.FegNetwork{
		Federation:  models2.NewDefaultNetworkFederationConfigs(),
		Description: "Foo Bar",
		DNS:         models.NewDefaultDNSConfig(),
		Features:    models.NewDefaultFeaturesConfig(),
		ID:          "n1",
		Name:        "foobar",
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedN1),
	}
	tests.RunUnitTest(t, e, tc)

	// Test UpdateNetwork
	payloadN1 := &models2.FegNetwork{
		ID:          "n1",
		Name:        "updated foobar",
		Description: "Updated Foo Bar",
		Federation:  models2.NewDefaultModifiedNetworkFederationConfigs(),
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
			},
		},
	}

	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/feg/n1",
		Payload:        payloadN1,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateNetwork,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualN1, err := configurator.LoadNetwork("n1", true, true, serdes.Network)
	assert.NoError(t, err)
	expected := configurator.Network{
		ID:          "n1",
		Type:        feg.FederationNetworkType,
		Name:        "updated foobar",
		Description: "Updated Foo Bar",
		Configs: map[string]interface{}{
			feg.FegNetworkType:          models2.NewDefaultModifiedNetworkFederationConfigs(),
			orc8r.DnsdNetworkType:       payloadN1.DNS,
			orc8r.NetworkFeaturesConfig: payloadN1.Features,
		},
		Version: 1,
	}
	assert.Equal(t, expected, actualN1)

	// Test GetFederationPartialGet
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg/n1/federation",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetworkFederationConfig,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models2.NewDefaultModifiedNetworkFederationConfigs()),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// setup fixtures in backend
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	seedFederationGateway(t)

	ctx := protos.NewGatewayIdentity("hw1", "n1", "g1").NewContextWithIdentity(context.Background())
	req := healthTestUtils.GetHealthyRequest()
	_, err = testHealthServicer.UpdateHealth(ctx, req)
	assert.NoError(t, err)

	expectedRes := &models2.FederationNetworkClusterStatus{
		ActiveGateway: "g1",
	}

	// Test Get Network HA status
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg/n1/status",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetworkFederationStatus,
		ExpectedStatus: 200,
		ExpectedResult: expectedRes,
	}
	tests.RunUnitTest(t, e, tc)

	// Test DeleteNetwork
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/feg/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        deleteNetwork,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg",
		Payload:        nil,
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n3", "n4"}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestFederationGateways(t *testing.T) {
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.UnfreezeClock(t)
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	testHealthServicer, err := healthTestInit.StartTestService(t)
	assert.NoError(t, err)

	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	createGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg/:network_id/gateways", obsidian.POST).HandlerFunc
	listGateways := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg/:network_id/gateways", obsidian.GET).HandlerFunc
	getGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg/:network_id/gateways/:gateway_id", obsidian.GET).HandlerFunc
	updateGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg/:network_id/gateways/:gateway_id", obsidian.PUT).HandlerFunc
	deleteGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg/:network_id/gateways/:gateway_id", obsidian.DELETE).HandlerFunc
	getHealth := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg/:network_id/gateways/:gateway_id/health_status", obsidian.GET).HandlerFunc
	seedFederationNetworks(t)

	// setup fixtures in backend
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	// Test CreateGateway
	payload := &models2.MutableFederationGateway{
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		ID:          "g1",
		Name:        "foobar",
		Description: "foo bar",
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          5,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		Federation: models2.NewDefaultGatewayFederationConfig(),
		Tier:       "t1",
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/feg/n1/gateways",
		Handler:        createGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Test ListGateways
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportGatewayStatus(t, ctx, models.NewDefaultGatewayStatus("hw1"))

	expected := map[string]*models2.FederationGateway{
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
			Federation: models2.NewDefaultGatewayFederationConfig(),
			Status:     models.NewDefaultGatewayStatus("hw1"),
		},
	}
	expected["g1"].Status.CheckinTime = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg/n1/gateways",
		Handler:        listGateways,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expected),
	}
	tests.RunUnitTest(t, e, tc)

	// Test GetGateway
	expectedGet := &models2.FederationGateway{
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
		Federation: models2.NewDefaultGatewayFederationConfig(),
		Status:     models.NewDefaultGatewayStatus("hw1"),
	}
	expectedGet.Status.CheckinTime = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg/n1/gateways/g1",
		Handler:        getGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: expectedGet,
	}
	tests.RunUnitTest(t, e, tc)

	// Test UpdateGateway
	payload = &models2.MutableFederationGateway{
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		ID:          "g1",
		Name:        "newname",
		Description: "bar baz",
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		Tier:       "t1",
		Federation: models2.NewDefaultGatewayFederationConfig(),
	}
	payload.Federation.AaaServer.AccountingEnabled = true

	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/feg/n1/gateways/g1",
		Handler:        updateGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expectedGet.Federation.AaaServer.AccountingEnabled = true
	expectedGet.Name = "newname"
	expectedGet.Description = "bar baz"
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg/n1/gateways/g1",
		Handler:        getGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: expectedGet,
	}
	tests.RunUnitTest(t, e, tc)

	// Test Get Health Status
	ctx = protos.NewGatewayIdentity("hw1", "n1", "g1").NewContextWithIdentity(context.Background())
	req := healthTestUtils.GetHealthyRequest()
	_, err = testHealthServicer.UpdateHealth(ctx, req)
	assert.NoError(t, err)

	expectedRes := &models2.FederationGatewayHealthStatus{
		Status:      models2.FederationGatewayHealthStatusStatusHEALTHY,
		Description: "OK",
	}

	// Test Health Gateway
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg/n1/gateways/g1/health_status",
		Handler:        getHealth,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: expectedRes,
	}
	tests.RunUnitTest(t, e, tc)

	// Test DeleteGateway
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/feg/n1/gateways/g1",
		Handler:        deleteGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg/n1/gateways/g1",
		Handler:        getGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestFederatedLteNetworks(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	listNetworks := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg_lte", obsidian.GET).HandlerFunc
	createNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg_lte", obsidian.POST).HandlerFunc
	getNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg_lte/:network_id", obsidian.GET).HandlerFunc
	updateNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg_lte/:network_id", obsidian.PUT).HandlerFunc
	deleteNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg_lte/:network_id", obsidian.DELETE).HandlerFunc
	getNetworkFederationConfig := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg_lte/:network_id/federation", obsidian.GET).HandlerFunc

	// Test ListNetworks
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg_lte",
		Payload:        nil,
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{}),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	seedFederatedLteNetworks(t)

	// CreateNetwork
	tc = tests.Test{
		Method: "POST",
		URL:    "/magma/v1/feg_lte",
		Payload: tests.JSONMarshaler(
			&models2.FegLteNetwork{
				Cellular:    models3.NewDefaultTDDNetworkConfig(),
				Federation:  models2.NewDefaultFederatedNetworkConfigs(),
				Description: "Foo Bar",
				DNS:         models.NewDefaultDNSConfig(),
				Features:    models.NewDefaultFeaturesConfig(),
				ID:          "n4",
				Name:        "foobar",
			},
		),
		Handler:        createNetwork,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Test ListNetworks
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg_lte",
		Payload:        nil,
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n1", "n3", "n4"}),
	}
	tests.RunUnitTest(t, e, tc)

	// Test GetNetwork
	expectedN1 := &models2.FegLteNetwork{
		Cellular:    models3.NewDefaultTDDNetworkConfig(),
		Federation:  models2.NewDefaultFederatedNetworkConfigs(),
		Description: "Foo Bar",
		DNS:         models.NewDefaultDNSConfig(),
		Features:    models.NewDefaultFeaturesConfig(),
		ID:          "n1",
		Name:        "foobar",
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg_lte/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedN1),
	}
	tests.RunUnitTest(t, e, tc)

	// Test UpdateNetwork
	payloadN1 := &models2.FegLteNetwork{
		ID:          "n1",
		Name:        "updated foobar",
		Description: "Updated Foo Bar",
		Cellular:    models3.NewDefaultFDDNetworkConfig(),
		Federation:  models2.NewDefaultFederatedNetworkConfigs(),
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
			},
		},
	}

	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/feg_lte/n1",
		Payload:        payloadN1,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateNetwork,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualN1, err := configurator.LoadNetwork("n1", true, true, serdes.Network)
	assert.NoError(t, err)
	expected := configurator.Network{
		ID:          "n1",
		Type:        feg.FederatedLteNetworkType,
		Name:        "updated foobar",
		Description: "Updated Foo Bar",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: models3.NewDefaultFDDNetworkConfig(),
			feg.FederatedNetworkType:      models2.NewDefaultFederatedNetworkConfigs(),
			orc8r.DnsdNetworkType:         payloadN1.DNS,
			orc8r.NetworkFeaturesConfig:   payloadN1.Features,
		},
		Version: 1,
	}
	assert.Equal(t, expected, actualN1)

	// Test GetFederationPartialGet
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg_lte/n1/federation",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetworkFederationConfig,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models2.NewDefaultFederatedNetworkConfigs()),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// Test DeleteNetwork
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/feg_lte/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        deleteNetwork,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg_lte",
		Payload:        nil,
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n3", "n4"}),
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_GetNetworkSubscriberConfigHandlers(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	e := echo.New()

	handlers := handlers.GetHandlers()
	getSubscriberConfig := tests.GetHandlerByPathAndMethod(t, handlers, "/magma/v1/feg_lte/:network_id/subscriber_config", obsidian.GET).HandlerFunc
	postSubscriberConfig := tests.GetHandlerByPathAndMethod(t, handlers, "/magma/v1/feg_lte/:network_id/subscriber_config", obsidian.PUT).HandlerFunc
	createNetwork := tests.GetHandlerByPathAndMethod(t, handlers, "/magma/v1/feg_lte", obsidian.POST).HandlerFunc

	// Pre: create network
	tc := tests.Test{
		Method: "POST",
		URL:    "/magma/v1/feg_lte",
		Payload: tests.JSONMarshaler(
			&models2.FegLteNetwork{
				Cellular:    models3.NewDefaultTDDNetworkConfig(),
				Federation:  models2.NewDefaultFederatedNetworkConfigs(),
				Description: "Foo Bar",
				DNS:         models.NewDefaultDNSConfig(),
				Features:    models.NewDefaultFeaturesConfig(),
				ID:          "n1",
				Name:        "network_1",
			},
		),
		Handler:        createNetwork,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	subConfig := &policyModels.NetworkSubscriberConfig{
		NetworkWideBaseNames: []policyModels.BaseName{"base1"},
		NetworkWideRuleNames: []string{"rule1"},
	}

	// Pass: create subscriber config
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/feg_lte/n1/subscriber_config",
		Payload:        tests.JSONMarshaler(subConfig),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        postSubscriberConfig,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Pass: read subscriber config
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/feg_lte/n1/subscriber_config",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getSubscriberConfig,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(subConfig),
	}
	tests.RunUnitTest(t, e, tc)
}

// n1, n3 are feg networks, n2 is not
func seedFederationNetworks(t *testing.T) {
	_, err := configurator.CreateNetworks(
		[]configurator.Network{
			{
				ID:          "n1",
				Type:        feg.FederationNetworkType,
				Name:        "foobar",
				Description: "Foo Bar",
				Configs: map[string]interface{}{
					feg.FegNetworkType:          models2.NewDefaultNetworkFederationConfigs(),
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
				Type:        feg.FederationNetworkType,
				Name:        "barfoo",
				Description: "Bar Foo",
				Configs:     map[string]interface{}{},
			},
		},
		serdes.Network,
	)
	assert.NoError(t, err)
}

// n1, n3 are feg networks, n2 is not
func seedFederatedLteNetworks(t *testing.T) {
	_, err := configurator.CreateNetworks(
		[]configurator.Network{
			{
				ID:          "n1",
				Type:        feg.FederatedLteNetworkType,
				Name:        "foobar",
				Description: "Foo Bar",
				Configs: map[string]interface{}{
					feg.FederatedNetworkType:      models2.NewDefaultFederatedNetworkConfigs(),
					lte.CellularNetworkConfigType: models3.NewDefaultTDDNetworkConfig(),
					orc8r.NetworkFeaturesConfig:   models.NewDefaultFeaturesConfig(),
					orc8r.DnsdNetworkType:         models.NewDefaultDNSConfig(),
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
				Type:        feg.FederatedLteNetworkType,
				Name:        "barfoo",
				Description: "Bar Foo",
				Configs:     map[string]interface{}{},
			},
		},
		serdes.Network,
	)
	assert.NoError(t, err)
}

func seedFederationGateway(t *testing.T) {
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	createGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/feg/:network_id/gateways", obsidian.POST).HandlerFunc

	// Test CreateGateway
	payload := &models2.MutableFederationGateway{
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		ID:          "g1",
		Name:        "foobar",
		Description: "foo bar",
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          5,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		Federation: models2.NewDefaultGatewayFederationConfig(),
		Tier:       "t1",
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/feg/n1/gateways",
		Handler:        createGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)
}
