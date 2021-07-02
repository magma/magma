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
	"crypto/x509"
	"fmt"
	"testing"
	"time"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/lte/obsidian/handlers"
	lteModels "magma/lte/cloud/go/services/lte/obsidian/models"
	policyModels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/clock"
	models2 "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/device"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/security/key"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func init() {
	//_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func TestListNetworks(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	listNetworks := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte", obsidian.GET).HandlerFunc

	// Test empty response
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte",
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{}),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	seedNetworks(t)

	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte",
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n1", "n3"}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestCreateNetwork(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	createNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte", obsidian.POST).HandlerFunc

	// test validation - include TDD and FDD configs
	payload := &lteModels.LteNetwork{
		Cellular:    lteModels.NewDefaultTDDNetworkConfig(),
		Description: "blah",
		DNS:         models.NewDefaultDNSConfig(),
		Features:    models.NewDefaultFeaturesConfig(),
		ID:          "n1",
		Name:        "foobar",
	}
	payload.Cellular.Ran.FddConfig = &lteModels.NetworkRanConfigsFddConfig{
		Earfcndl: 17000,
		Earfcnul: 18000,
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte",
		Payload:        payload,
		Handler:        createNetwork,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"only one of TDD or FDD configs can be set",
	}
	tests.RunUnitTest(t, e, tc)

	// happy path
	payload = &lteModels.LteNetwork{
		Cellular:    lteModels.NewDefaultTDDNetworkConfig(),
		Description: "Foo Bar",
		DNS:         models.NewDefaultDNSConfig(),
		Features:    models.NewDefaultFeaturesConfig(),
		ID:          "n1",
		Name:        "foobar",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte",
		Payload:        payload,
		Handler:        createNetwork,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadNetwork("n1", true, true, serdes.Network)
	assert.NoError(t, err)
	expected := configurator.Network{
		ID:          "n1",
		Type:        lte.NetworkType,
		Name:        "foobar",
		Description: "Foo Bar",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: lteModels.NewDefaultTDDNetworkConfig(),
			orc8r.DnsdNetworkType:         models.NewDefaultDNSConfig(),
			orc8r.NetworkFeaturesConfig:   models.NewDefaultFeaturesConfig(),
		},
	}
	assert.Equal(t, expected, actual)
}

func TestGetNetwork(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	getNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id", obsidian.GET).HandlerFunc

	// Test 404
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n1",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetwork,
		ExpectedStatus: 404,
	}
	tests.RunUnitTest(t, e, tc)

	seedNetworks(t)

	expectedN1 := &lteModels.LteNetwork{
		Cellular:    lteModels.NewDefaultTDDNetworkConfig(),
		Description: "Foo Bar",
		DNS:         models.NewDefaultDNSConfig(),
		Features:    models.NewDefaultFeaturesConfig(),
		ID:          "n1",
		Name:        "foobar",
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n1",
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
		URL:            "/magma/v1/lte/n2",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        getNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not a <lte> network",
	}
	tests.RunUnitTest(t, e, tc)

	// get a network without any configs (poorly formed data)
	expectedN3 := &lteModels.LteNetwork{
		Description: "Bar Foo",
		ID:          "n3",
		Name:        "barfoo",
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n3",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n3"},
		Handler:        getNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedN3),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateNetwork(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	updateNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id", obsidian.PUT).HandlerFunc

	// Test validation failure
	payloadN1 := &lteModels.LteNetwork{
		ID:          "n1",
		Name:        "updated foobar",
		Description: "Updated Foo Bar",
		Cellular:    lteModels.NewDefaultFDDNetworkConfig(),
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
		URL:            "/magma/v1/lte/n1",
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
		URL:            "/magma/v1/lte/n1",
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
		URL:            "/magma/v1/lte/n1",
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
		Type:        lte.NetworkType,
		Name:        "updated foobar",
		Description: "Updated Foo Bar",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: lteModels.NewDefaultFDDNetworkConfig(),
			orc8r.DnsdNetworkType:         payloadN1.DNS,
			orc8r.NetworkFeaturesConfig:   payloadN1.Features,
		},
		Version: 1,
	}
	assert.Equal(t, expected, actualN1)

	// update n2, should be 400
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/lte/n2",
		Payload:        payloadN1,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        updateNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not a <lte> network",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestDeleteNetwork(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	deleteNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id", obsidian.DELETE).HandlerFunc

	// Test 404
	tc := tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/lte/n1",
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
		URL:            "/magma/v1/lte/n2",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        deleteNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "network n2 is not a <lte> network",
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.ListNetworkIDs()
	assert.NoError(t, err)
	assert.Equal(t, []string{"n2", "n3"}, actual)
}

func TestCellularPartialGet(t *testing.T) {
	configuratorTestInit.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/lte"

	seedNetworks(t)

	handlers := handlers.GetHandlers()
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
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getCellular,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(lteModels.NewDefaultTDDNetworkConfig()),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// 404
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/cellular/", testURLRoot, "n2"),
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
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getEpc,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(lteModels.NewDefaultTDDNetworkConfig().Epc),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// 404
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/cellular/epc/", testURLRoot, "n2"),
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
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getRan,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(lteModels.NewDefaultTDDNetworkConfig().Ran),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// 404
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/cellular/ran/", testURLRoot, "n2"),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        getRan,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// add 'n2' as FegNetworkID to n1
	cellularConfig := lteModels.NewDefaultTDDNetworkConfig()
	cellularConfig.FegNetworkID = "n2"
	err := configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{
		{
			ID: "n1",
			ConfigsToAddOrUpdate: map[string]interface{}{
				lte.CellularNetworkConfigType: cellularConfig,
			},
		},
	},
		serdes.Network,
	)
	assert.NoError(t, err)

	// happy case FegNetworkID from cellular config
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/cellular/feg_network_id/", testURLRoot, "n1"),
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
	configuratorTestInit.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/lte"

	seedNetworks(t)
	handlers := handlers.GetHandlers()
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
		Payload:        lteModels.NewDefaultFDDNetworkConfig(),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n2"},
		Handler:        updateCellular,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualN2, err := configurator.LoadNetwork("n2", true, true, serdes.Network)
	assert.NoError(t, err)
	expected := configurator.Network{
		ID:          "n2",
		Type:        "blah",
		Name:        "foobar",
		Description: "Foo Bar",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: lteModels.NewDefaultFDDNetworkConfig(),
		},
		Version: 1,
	}
	assert.Equal(t, expected, actualN2)

	// Validation error (cellular config has both tdd and fdd config)
	badCellularConfig := lteModels.NewDefaultTDDNetworkConfig()
	badCellularConfig.Ran.FddConfig = &lteModels.NetworkRanConfigsFddConfig{
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
		Payload:        lteModels.NewDefaultTDDNetworkConfig().Epc,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n3"},
		Handler:        updateEpc,
		ExpectedStatus: 400,
		ExpectedError:  "No cellular network config found",
	}
	tests.RunUnitTest(t, e, tc)

	// happy path update epc config
	epcConfig := lteModels.NewDefaultTDDNetworkConfig().Epc
	epcConfig.HssRelayEnabled = swag.Bool(true)
	epcConfig.GxGyRelayEnabled = swag.Bool(true)
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

	actualN2, err = configurator.LoadNetwork("n2", true, true, serdes.Network)
	assert.NoError(t, err)
	expected.Configs[lte.CellularNetworkConfigType].(*lteModels.NetworkCellularConfigs).Epc = epcConfig
	expected.Version = 2
	assert.Equal(t, expected, actualN2)

	// Fail to put epc config to a network without cellular network configs
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/cellular/ran/", testURLRoot, "n3"),
		Payload:        lteModels.NewDefaultTDDNetworkConfig().Ran,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n3"},
		Handler:        updateRan,
		ExpectedStatus: 400,
		ExpectedError:  "No cellular network config found",
	}
	tests.RunUnitTest(t, e, tc)

	// Validation error
	ranConfig := lteModels.NewDefaultTDDNetworkConfig().Ran
	ranConfig.FddConfig = lteModels.NewDefaultFDDNetworkConfig().Ran.FddConfig
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
	ranConfig = lteModels.NewDefaultFDDNetworkConfig().Ran
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
	actualN2, err = configurator.LoadNetwork("n2", true, true, serdes.Network)
	assert.NoError(t, err)
	expected.Configs[lte.CellularNetworkConfigType].(*lteModels.NetworkCellularConfigs).Ran = ranConfig
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
	configuratorTestInit.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/lte"

	seedNetworks(t)

	handlers := handlers.GetHandlers()
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

	_, err := configurator.LoadNetworkConfig("n1", lte.CellularNetworkConfigType, serdes.Network)
	assert.EqualError(t, err, "Not found")
}

func Test_GetNetworkSubscriberConfigHandlers(t *testing.T) {
	configuratorTestInit.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	seedNetworks(t)

	obsidianHandlers := handlers.GetHandlers()
	getSubscriberConfig := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/subscriber_config", obsidian.GET).HandlerFunc
	getRuleNames := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/subscriber_config/rule_names", obsidian.GET).HandlerFunc
	getBaseNames := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/subscriber_config/base_names", obsidian.GET).HandlerFunc

	// 404
	tc := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/", testURLRoot, "n1"),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getSubscriberConfig,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(&policyModels.NetworkSubscriberConfig{}),
	}
	tests.RunUnitTest(t, e, tc)

	subscriberConfig := &policyModels.NetworkSubscriberConfig{
		NetworkWideBaseNames: []policyModels.BaseName{"base1"},
		NetworkWideRuleNames: []string{"rule1"},
	}
	assert.NoError(t, configurator.UpdateNetworkConfig("n1", lte.NetworkSubscriberConfigType, subscriberConfig, serdes.Network))

	// happy case
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/", testURLRoot, "n1"),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getSubscriberConfig,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(subscriberConfig),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/base_names/", testURLRoot, "n1"),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getBaseNames,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(subscriberConfig.NetworkWideBaseNames),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/rule_names/", testURLRoot, "n1"),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getRuleNames,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(subscriberConfig.NetworkWideRuleNames),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_ModifyNetworkSubscriberConfigHandlers(t *testing.T) {
	configuratorTestInit.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	seedNetworks(t)

	obsidianHandlers := handlers.GetHandlers()
	putSubscriberConfig := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/subscriber_config", obsidian.PUT).HandlerFunc
	putRuleNames := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/subscriber_config/rule_names", obsidian.PUT).HandlerFunc
	putBaseNames := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/subscriber_config/base_names", obsidian.PUT).HandlerFunc
	postRuleName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/subscriber_config/rule_names/:rule_id", obsidian.POST).HandlerFunc
	postBaseName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/subscriber_config/base_names/:base_name", obsidian.POST).HandlerFunc
	deleteRuleName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/subscriber_config/rule_names/:rule_id", obsidian.DELETE).HandlerFunc
	deleteBaseName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/subscriber_config/base_names/:base_name", obsidian.DELETE).HandlerFunc

	subscriberConfig := &policyModels.NetworkSubscriberConfig{
		NetworkWideBaseNames: []policyModels.BaseName{"base1"},
		NetworkWideRuleNames: []string{"rule1"},
	}

	// non-existent network id
	tc := tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/base_names/", testURLRoot, "n32"),
		Payload:        tests.JSONMarshaler(subscriberConfig.NetworkWideBaseNames),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n32"},
		Handler:        putBaseNames,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/rule_names/", testURLRoot, "n32"),
		Payload:        tests.JSONMarshaler(subscriberConfig.NetworkWideRuleNames),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n32"},
		Handler:        putRuleNames,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// add to non existent config
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/base_names/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(subscriberConfig.NetworkWideBaseNames),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        putBaseNames,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/rule_names/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(subscriberConfig.NetworkWideRuleNames),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        putRuleNames,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	iSubscriberConfig, err := configurator.LoadNetworkConfig("n1", lte.NetworkSubscriberConfigType, serdes.Network)
	assert.NoError(t, err)
	assert.Equal(t, subscriberConfig, iSubscriberConfig.(*policyModels.NetworkSubscriberConfig))

	newRuleNames := []string{"rule2"}
	// happy case
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/rule_names/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(newRuleNames),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        putRuleNames,
		ExpectedStatus: 204,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	newBaseNames := []policyModels.BaseName{"base2"}
	// happy case
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/base_names/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(newBaseNames),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        putBaseNames,
		ExpectedStatus: 204,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	iSubscriberConfig, err = configurator.LoadNetworkConfig("n1", lte.NetworkSubscriberConfigType, serdes.Network)
	assert.NoError(t, err)
	actualSubscriberConfig := iSubscriberConfig.(*policyModels.NetworkSubscriberConfig)

	assert.ElementsMatch(t, newRuleNames, actualSubscriberConfig.NetworkWideRuleNames)
	assert.ElementsMatch(t, newBaseNames, actualSubscriberConfig.NetworkWideBaseNames)

	newSubscriberConfig := &policyModels.NetworkSubscriberConfig{
		NetworkWideBaseNames: []policyModels.BaseName{"base3"},
		NetworkWideRuleNames: []string{"rule3"},
	}
	// happy case
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(newSubscriberConfig),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        putSubscriberConfig,
		ExpectedStatus: 204,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	iSubscriberConfig, err = configurator.LoadNetworkConfig("n1", lte.NetworkSubscriberConfigType, serdes.Network)
	assert.NoError(t, err)
	actualSubscriberConfig = iSubscriberConfig.(*policyModels.NetworkSubscriberConfig)

	assert.Equal(t, newSubscriberConfig, actualSubscriberConfig)

	tc = tests.Test{
		Method:         "POST",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/rule_names/%s", testURLRoot, "n1", "rule4"),
		Payload:        tests.JSONMarshaler(newSubscriberConfig),
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "rule4"},
		Handler:        postRuleName,
		ExpectedStatus: 201,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// posting twice shouldn't affect anything
	tc = tests.Test{
		Method:         "POST",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/rule_names/%s", testURLRoot, "n1", "rule4"),
		Payload:        tests.JSONMarshaler(newSubscriberConfig),
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "rule4"},
		Handler:        postRuleName,
		ExpectedStatus: 201,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "POST",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/base_names/%s", testURLRoot, "n1", "base4"),
		Payload:        tests.JSONMarshaler(newSubscriberConfig),
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "base4"},
		Handler:        postBaseName,
		ExpectedStatus: 201,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)
	tc = tests.Test{
		Method:         "POST",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/base_names/%s", testURLRoot, "n1", "base4"),
		Payload:        tests.JSONMarshaler(newSubscriberConfig),
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "base4"},
		Handler:        postBaseName,
		ExpectedStatus: 201,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	newSubscriberConfig = &policyModels.NetworkSubscriberConfig{
		NetworkWideBaseNames: []policyModels.BaseName{"base3", "base4"},
		NetworkWideRuleNames: []string{"rule3", "rule4"},
	}
	iSubscriberConfig, err = configurator.LoadNetworkConfig("n1", lte.NetworkSubscriberConfigType, serdes.Network)
	assert.NoError(t, err)
	actualSubscriberConfig = iSubscriberConfig.(*policyModels.NetworkSubscriberConfig)
	assert.Equal(t, newSubscriberConfig, actualSubscriberConfig)

	tc = tests.Test{
		Method:         "DELETE",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/rule_names/%s", testURLRoot, "n1", "rule4"),
		Payload:        tests.JSONMarshaler(newSubscriberConfig),
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "rule4"},
		Handler:        deleteRuleName,
		ExpectedStatus: 204,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "DELETE",
		URL:            fmt.Sprintf("%s/%s/subscriber_config/base_names/%s", testURLRoot, "n1", "base4"),
		Payload:        tests.JSONMarshaler(newSubscriberConfig),
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "base4"},
		Handler:        deleteBaseName,
		ExpectedStatus: 204,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	newSubscriberConfig = &policyModels.NetworkSubscriberConfig{
		NetworkWideBaseNames: []policyModels.BaseName{"base3"},
		NetworkWideRuleNames: []string{"rule3"},
	}
	iSubscriberConfig, err = configurator.LoadNetworkConfig("n1", lte.NetworkSubscriberConfigType, serdes.Network)
	assert.NoError(t, err)
	actualSubscriberConfig = iSubscriberConfig.(*policyModels.NetworkSubscriberConfig)
	assert.Equal(t, newSubscriberConfig, actualSubscriberConfig)
}

func TestCreateGateway(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)

	// setup fixtures in backend
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	err = device.RegisterDevice(
		"n1", orc8r.AccessGatewayRecordType, "hw2",
		&models.GatewayDevice{
			HardwareID: "hw2",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		serdes.Device,
	)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/gateways"
	hands := handlers.GetHandlers()
	createGateway := tests.GetHandlerByPathAndMethod(t, hands, testURLRoot, obsidian.POST).HandlerFunc

	// happy path, no device
	payload := &lteModels.MutableLteGateway{
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		ID:          "g1",
		Name:        "foobar",
		Description: "foo bar",
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          10,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		Cellular:               newDefaultGatewayConfig(),
		ConnectedEnodebSerials: []string{"enb1"},
		Tier:                   "t1",
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Handler:        createGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g1"},
			{Type: lte.CellularGatewayEntityType, Key: "g1"},
		},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	actualDevice, err := device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hw1", serdes.Device)
	assert.NoError(t, err)

	expectedEnts := configurator.NetworkEntities{
		{
			NetworkID: "n1", Type: lte.CellularGatewayEntityType, Key: "g1",
			Name: string(payload.Name), Description: string(payload.Description),
			Config:             payload.Cellular,
			Associations:       []storage.TypeAndKey{{Type: lte.CellularEnodebEntityType, Key: "enb1"}},
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
			GraphID:            "2",
		},
		{
			NetworkID: "n1", Type: orc8r.MagmadGatewayType, Key: "g1",
			Name: string(payload.Name), Description: string(payload.Description),
			PhysicalID:         "hw1",
			Config:             payload.Magmad,
			Associations:       []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			GraphID:            "2",
			Version:            1,
		},
	}
	assert.Equal(t, expectedEnts, actualEnts)
	assert.Equal(t, payload.Device, actualDevice)

	// valid magmad gateway, invalid cellular - nothing should change on backend
	payload = &lteModels.MutableLteGateway{
		Device: &models.GatewayDevice{
			HardwareID: "hw2",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		ID:          "g3",
		Name:        "foobar",
		Description: "foo bar",
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          10,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		Cellular: newDefaultGatewayConfig(),
		// Invalid due to nonexistent enb
		ConnectedEnodebSerials: []string{"enb1", "dne"},
		Tier:                   "t1",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Handler:        createGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 500,
		ExpectedError:  "error creating gateway: rpc error: code = Internal desc = could not find entities matching [type:\"cellular_enodeb\" key:\"dne\" ]",
	}
	tests.RunUnitTest(t, e, tc)

	actualEnts, _, err = configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g3"},
			{Type: lte.CellularGatewayEntityType, Key: "g3"},
		},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	// the device should get created regardless
	actualDevice, err = device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hw2", serdes.Device)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(actualEnts))
	assert.Equal(t, payload.Device, actualDevice)

	// Some composite validation failures - bad device key, missing required
	// non-EPS control fields when non-EPS service control is on
	pubkeyB64 := strfmt.Base64("fake key")
	payload = &lteModels.MutableLteGateway{
		Device: &models.GatewayDevice{
			HardwareID: "foo-bar-baz-890",
			Key: &models.ChallengeKey{
				KeyType: "SOFTWARE_ECDSA_SHA256",
				Key:     &pubkeyB64,
			},
		},
		ID:          "g4",
		Name:        "foobar",
		Description: "foo bar",
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          10,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		Cellular:               newDefaultGatewayConfig(),
		ConnectedEnodebSerials: []string{},
		Tier:                   "t1",
	}
	payload.Cellular.NonEpsService = &lteModels.GatewayNonEpsConfigs{
		NonEpsServiceControl: swag.Uint32(1),
	}

	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Handler:        createGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"validation failure list:\n" +
			"arfcn_2g in body is required\n" +
			"csfb_mcc in body is required\n" +
			"csfb_mnc in body is required\n" +
			"csfb_rat in body is required\n" +
			"lac in body is required\n" +
			"Failed to parse key: asn1: structure error: tags don't match (16 vs {class:1 tag:6 length:97 isCompound:true}) {optional:false explicit:false application:false private:false defaultValue:<nil> tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} publicKeyInfo @2",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestListAndGetGateways(t *testing.T) {
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.UnfreezeClock(t)

	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/gateways"

	handlers := handlers.GetHandlers()
	listGateways := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc
	getGateway := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/:gateway_id", testURLRoot), obsidian.GET).HandlerFunc

	// Create 2 gateways, 1 with state and device, the other without
	// g2 will associate to 2 enodebs
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
			{Type: lte.CellularEnodebEntityType, Key: "enb2"},
			{
				Type: lte.CellularGatewayEntityType, Key: "g1",
				Config: &lteModels.GatewayCellularConfigs{
					Epc: &lteModels.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
					Ran: &lteModels.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
				},
			},
			{
				Type: lte.CellularGatewayEntityType, Key: "g2",
				Config: &lteModels.GatewayCellularConfigs{
					Epc: &lteModels.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
					Ran: &lteModels.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
				},
				Associations: []storage.TypeAndKey{
					{Type: lte.CellularEnodebEntityType, Key: "enb1"},
					{Type: lte.CellularEnodebEntityType, Key: "enb2"},
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
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
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
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g2"}},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: "g1"},
					{Type: orc8r.MagmadGatewayType, Key: "g2"},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	err = device.RegisterDevice(
		"n1", orc8r.AccessGatewayRecordType, "hw1",
		&models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}},
		serdes.Device,
	)
	assert.NoError(t, err)
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportGatewayStatus(t, ctx, models.NewDefaultGatewayStatus("hw1"))

	expected := map[string]*lteModels.LteGateway{
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
			Cellular: &lteModels.GatewayCellularConfigs{
				Epc: &lteModels.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
				Ran: &lteModels.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
			},
			Status:                 models.NewDefaultGatewayStatus("hw1"),
			ConnectedEnodebSerials: lteModels.EnodebSerials{},
			ApnResources:           lteModels.ApnResources{},
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
			Cellular: &lteModels.GatewayCellularConfigs{
				Epc: &lteModels.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
				Ran: &lteModels.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
			},
			ConnectedEnodebSerials: []string{"enb1", "enb2"},
			ApnResources:           lteModels.ApnResources{},
		},
	}
	expected["g1"].Status.CheckinTime = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))

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

	expectedGet := &lteModels.LteGateway{
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
		Cellular: &lteModels.GatewayCellularConfigs{
			Epc: &lteModels.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
			Ran: &lteModels.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
		},
		Status:                 models.NewDefaultGatewayStatus("hw1"),
		ConnectedEnodebSerials: lteModels.EnodebSerials{},
		ApnResources:           lteModels.ApnResources{},
	}
	expectedGet.Status.CheckinTime = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
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

	expectedGet = &lteModels.LteGateway{
		ID:   "g2",
		Name: "barfoo", Description: "bar foo",
		Tier: "t1",
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		Cellular: &lteModels.GatewayCellularConfigs{
			Epc: &lteModels.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
			Ran: &lteModels.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
		},
		ConnectedEnodebSerials: []string{"enb1", "enb2"},
		ApnResources:           lteModels.ApnResources{},
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
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.UnfreezeClock(t)

	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/gateways/:gateway_id"
	handlers := handlers.GetHandlers()
	updateGateway := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.PUT).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
			{Type: lte.CellularEnodebEntityType, Key: "enb2"},
			{Type: lte.CellularEnodebEntityType, Key: "enb3"},
			{
				Type: lte.CellularGatewayEntityType, Key: "g1",
				Config: &lteModels.GatewayCellularConfigs{
					Epc: &lteModels.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
					Ran: &lteModels.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
				},
				Associations: []storage.TypeAndKey{
					{Type: lte.CellularEnodebEntityType, Key: "enb1"},
					{Type: lte.CellularEnodebEntityType, Key: "enb2"},
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
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: "g1"},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	err = device.RegisterDevice(
		"n1", orc8r.AccessGatewayRecordType, "hw1",
		&models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}},
		serdes.Device,
	)
	assert.NoError(t, err)

	// update everything
	privateKey, err := key.GenerateKey("P256", 0)
	assert.NoError(t, err)
	marshaledPubKey, err := x509.MarshalPKIXPublicKey(key.PublicKey(privateKey))
	assert.NoError(t, err)
	pubkeyB64 := strfmt.Base64(marshaledPubKey)
	payload := &lteModels.MutableLteGateway{
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
		Cellular: &lteModels.GatewayCellularConfigs{
			Epc: &lteModels.GatewayEpcConfigs{NatEnabled: swag.Bool(false), IPBlock: "172.10.10.0/24"},
			Ran: &lteModels.GatewayRanConfigs{Pci: 123, TransmitEnabled: swag.Bool(false)},
		},
		ConnectedEnodebSerials: []string{"enb1", "enb3"},
		ApnResources:           lteModels.ApnResources{},
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
			{Type: lte.CellularGatewayEntityType, Key: "g1"},
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
		},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	actualDevice, err := device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hw1", serdes.Device)
	assert.NoError(t, err)

	expectedEnts := configurator.NetworkEntities{
		{
			NetworkID: "n1", Type: lte.CellularGatewayEntityType, Key: "g1",
			Name: string(payload.Name), Description: string(payload.Description),
			Config:             payload.Cellular,
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
			Associations: []storage.TypeAndKey{
				{Type: lte.CellularEnodebEntityType, Key: "enb1"},
				{Type: lte.CellularEnodebEntityType, Key: "enb3"},
			},
			GraphID: "10",
			Version: 1,
		},
		{
			NetworkID: "n1", Type: orc8r.MagmadGatewayType, Key: "g1",
			Name: string(payload.Name), Description: string(payload.Description),
			PhysicalID:         "hw1",
			Config:             payload.Magmad,
			Associations:       []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
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
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.UnfreezeClock(t)

	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/gateways/:gateway_id"
	handlers := handlers.GetHandlers()
	deleteGateway := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.DELETE).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
			{Type: lte.CellularEnodebEntityType, Key: "enb2"},
			{
				Type: lte.CellularGatewayEntityType, Key: "g1",
				Config: &lteModels.GatewayCellularConfigs{
					Epc: &lteModels.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
					Ran: &lteModels.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
				},
				Associations: []storage.TypeAndKey{
					{Type: lte.CellularEnodebEntityType, Key: "enb1"},
					{Type: lte.CellularEnodebEntityType, Key: "enb2"},
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
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: "g1"},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	err = device.RegisterDevice(
		"n1", orc8r.AccessGatewayRecordType, "hw1",
		&models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}},
		serdes.Device,
	)
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
			{Type: lte.CellularGatewayEntityType, Key: "g1"},
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
		},
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	actualDevice, err := device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hw1", serdes.Device)
	assert.Nil(t, actualDevice)
	assert.EqualError(t, err, "Not found")

	expectedEnts := configurator.NetworkEntities{
		{NetworkID: "n1", Type: orc8r.UpgradeTierEntityType, Key: "t1", GraphID: "11"},
	}
	assert.Equal(t, expectedEnts, actualEnts)
}

func TestGetCellularGatewayConfig(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/gateways/:gateway_id"
	handlers := handlers.GetHandlers()
	getCellular := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular", testURLRoot), obsidian.GET).HandlerFunc
	getEpc := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/epc", testURLRoot), obsidian.GET).HandlerFunc
	getRan := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/ran", testURLRoot), obsidian.GET).HandlerFunc
	getNonEps := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/non_eps", testURLRoot), obsidian.GET).HandlerFunc
	getEnodebs := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/connected_enodeb_serials", testURLRoot), obsidian.GET).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
			{Type: lte.CellularEnodebEntityType, Key: "enb2"},
			{
				Type: lte.CellularGatewayEntityType, Key: "g1",
				Config: newDefaultGatewayConfig(),
				Associations: []storage.TypeAndKey{
					{Type: lte.CellularEnodebEntityType, Key: "enb1"},
					{Type: lte.CellularEnodebEntityType, Key: "enb2"},
				},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Name: "foobar", Description: "foo bar",
				PhysicalID:   "hw1",
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			},
		},
		serdes.Entity,
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
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/gateways/:gateway_id"
	handlers := handlers.GetHandlers()
	updateCellular := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular", testURLRoot), obsidian.PUT).HandlerFunc
	updateEpc := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/epc", testURLRoot), obsidian.PUT).HandlerFunc
	updateRan := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/ran", testURLRoot), obsidian.PUT).HandlerFunc
	updateNonEps := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/cellular/non_eps", testURLRoot), obsidian.PUT).HandlerFunc
	updateEnodebs := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/connected_enodeb_serials", testURLRoot), obsidian.PUT).HandlerFunc
	postEnodeb := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/connected_enodeb_serials", testURLRoot), obsidian.POST).HandlerFunc
	deleteEnodeb := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/connected_enodeb_serials", testURLRoot), obsidian.DELETE).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
			{Type: lte.CellularEnodebEntityType, Key: "enb2"},
			{Type: lte.CellularGatewayEntityType, Key: "g1"},
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			},
		},
		serdes.Entity,
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

	expected := configurator.NetworkEntitiesByTK{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			GraphID:      "6",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayEntityType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayEntityType, Key: "g1",
			Config:  newDefaultGatewayConfig(),
			GraphID: "6",
			Version: 1,
		},
	}

	entities, _, err := configurator.LoadEntities(
		"n1", nil, swag.String("g1"), nil, nil,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.MakeByTK())

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

	expected = configurator.NetworkEntitiesByTK{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			GraphID:      "6",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayEntityType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayEntityType, Key: "g1",
			Config:  modifiedCellularConfig,
			GraphID: "6",
			Version: 2,
		},
	}
	entities, _, err = configurator.LoadEntities(
		"n1", nil, swag.String("g1"), nil, nil,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.MakeByTK())

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

	expected = configurator.NetworkEntitiesByTK{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			GraphID:      "6",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayEntityType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayEntityType, Key: "g1",
			Config:  modifiedCellularConfig,
			GraphID: "6",
			Version: 3,
		},
	}
	entities, _, err = configurator.LoadEntities(
		"n1", nil, swag.String("g1"), nil, nil,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.MakeByTK())

	// validation failure
	modifiedCellularConfig.NonEpsService.NonEpsServiceControl = swag.Uint32(1)
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

	expected = configurator.NetworkEntitiesByTK{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			GraphID:      "6",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayEntityType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayEntityType, Key: "g1",
			Config:  modifiedCellularConfig,
			GraphID: "6",
			Version: 4,
		},
	}
	entities, _, err = configurator.LoadEntities(
		"n1", nil, swag.String("g1"), nil, nil,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.MakeByTK())

	// connected enodeBs - happy case
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/connected_enodeb_serial", testURLRoot),
		Handler:        updateEnodebs,
		Payload:        tests.JSONMarshaler([]string{"enb1", "enb2"}),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expected = configurator.NetworkEntitiesByTK{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			GraphID:      "2",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayEntityType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayEntityType, Key: "g1",
			Config:  modifiedCellularConfig,
			GraphID: "2",
			Version: 5,
			Associations: []storage.TypeAndKey{
				{Type: lte.CellularEnodebEntityType, Key: "enb1"},
				{Type: lte.CellularEnodebEntityType, Key: "enb2"},
			},
		},
	}
	entities, _, err = configurator.LoadEntities(
		"n1", nil, swag.String("g1"), nil, nil,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.MakeByTK())

	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: lte.CellularEnodebEntityType, Key: "enb3"}, serdes.Entity)
	assert.NoError(t, err)

	// happy case
	tc = tests.Test{
		Method:         "POST",
		URL:            fmt.Sprintf("%s/connected_enodeb_serial", testURLRoot),
		Handler:        postEnodeb,
		Payload:        tests.JSONMarshaler("enb3"),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expected = configurator.NetworkEntitiesByTK{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			GraphID:      "10",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayEntityType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayEntityType, Key: "g1",
			Config:  modifiedCellularConfig,
			GraphID: "10",
			Version: 6,
			Associations: []storage.TypeAndKey{
				{Type: lte.CellularEnodebEntityType, Key: "enb1"},
				{Type: lte.CellularEnodebEntityType, Key: "enb2"},
				{Type: lte.CellularEnodebEntityType, Key: "enb3"},
			},
		},
	}
	entities, _, err = configurator.LoadEntities(
		"n1", nil, swag.String("g1"), nil, nil,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.MakeByTK())

	// happy case
	tc = tests.Test{
		Method:         "DELETE",
		URL:            fmt.Sprintf("%s/connected_enodeb_serial", testURLRoot),
		Handler:        deleteEnodeb,
		Payload:        tests.JSONMarshaler("enb3"),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expected = configurator.NetworkEntitiesByTK{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			GraphID:      "10",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayEntityType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayEntityType, Key: "g1",
			Config:  modifiedCellularConfig,
			GraphID: "10",
			Version: 7,
			Associations: []storage.TypeAndKey{
				{Type: lte.CellularEnodebEntityType, Key: "enb1"},
				{Type: lte.CellularEnodebEntityType, Key: "enb2"},
			},
		},
	}
	entities, _, err = configurator.LoadEntities(
		"n1", nil, swag.String("g1"), nil, nil,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.MakeByTK())

	// Clear enb serial list
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/connected_enodeb_serial", testURLRoot),
		Handler:        updateEnodebs,
		Payload:        tests.JSONMarshaler([]string{}),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expected = configurator.NetworkEntitiesByTK{
		storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      orc8r.MagmadGatewayType, Key: "g1",
			Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: "g1"}},
			GraphID:      "10",
			Version:      0,
		},
		storage.TypeAndKey{Type: lte.CellularGatewayEntityType, Key: "g1"}: {
			NetworkID: "n1",
			Type:      lte.CellularGatewayEntityType, Key: "g1",
			Config:  modifiedCellularConfig,
			GraphID: "10",
			Version: 8,
		},
	}
	entities, _, err = configurator.LoadEntities(
		"n1", nil, swag.String("g1"), nil, nil,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expected, entities.MakeByTK())
}

func TestListAndGetEnodebs(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/enodebs"

	handlers := handlers.GetHandlers()
	listEnodebs := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc
	getEnodeb := tests.GetHandlerByPathAndMethod(t, handlers, fmt.Sprintf("%s/:enodeb_serial", testURLRoot), obsidian.GET).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type:        lte.CellularEnodebEntityType,
				Key:         "abcdefg",
				Name:        "abc enodeb",
				Description: "abc enodeb description",
				PhysicalID:  "abcdefg",
				Config: &lteModels.EnodebConfig{
					ConfigType: "MANAGED",
					ManagedConfig: &lteModels.EnodebConfiguration{
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
			},
			{
				Type:        lte.CellularEnodebEntityType,
				Key:         "vwxyz",
				Name:        "xyz enodeb",
				Description: "xyz enodeb description",
				PhysicalID:  "vwxyz",
				Config: &lteModels.EnodebConfig{
					ConfigType: "MANAGED",
					ManagedConfig: &lteModels.EnodebConfiguration{
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
			},
			{
				Type: lte.CellularGatewayEntityType, Key: "gw1",
				Associations: []storage.TypeAndKey{{Type: lte.CellularEnodebEntityType, Key: "abcdefg"}},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	expected := map[string]*lteModels.Enodeb{
		"abcdefg": {
			AttachedGatewayID: "gw1",
			Config: &lteModels.EnodebConfiguration{
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
			EnodebConfig: &lteModels.EnodebConfig{
				ConfigType: "MANAGED",
				ManagedConfig: &lteModels.EnodebConfiguration{
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
			Name:        "abc enodeb",
			Description: "abc enodeb description",
			Serial:      "abcdefg",
		},
		"vwxyz": {
			Config: &lteModels.EnodebConfiguration{
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
			EnodebConfig: &lteModels.EnodebConfig{
				ConfigType: "MANAGED",
				ManagedConfig: &lteModels.EnodebConfiguration{
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
			Name:        "xyz enodeb",
			Description: "xyz enodeb description",
			Serial:      "vwxyz",
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
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/enodebs"

	handlers := handlers.GetHandlers()
	createEnodeb := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.POST).HandlerFunc

	tc := tests.Test{
		Method:  "POST",
		URL:     testURLRoot,
		Handler: createEnodeb,
		Payload: &lteModels.Enodeb{
			Config: &lteModels.EnodebConfiguration{
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
			EnodebConfig: &lteModels.EnodebConfig{
				ConfigType: "MANAGED",
				ManagedConfig: &lteModels.EnodebConfiguration{
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
			Name:        "foobar",
			Description: "foobar description",
			Serial:      "abcdef",
		},
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.CellularEnodebEntityType, "abcdef", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.CellularEnodebEntityType, Key: "abcdef",
		Name:        "foobar",
		Description: "foobar description",
		PhysicalID:  "abcdef",
		GraphID:     "2",
		Config: &lteModels.EnodebConfig{
			ConfigType: "MANAGED",
			ManagedConfig: &lteModels.EnodebConfiguration{
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
	}
	assert.Equal(t, expected, actual)

	tc = tests.Test{
		Method:  "POST",
		URL:     testURLRoot,
		Handler: createEnodeb,
		Payload: &lteModels.Enodeb{
			Config: &lteModels.EnodebConfiguration{
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
			EnodebConfig: &lteModels.EnodebConfig{
				ConfigType: "MANAGED",
				ManagedConfig: &lteModels.EnodebConfiguration{
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

	ip := strfmt.IPv4("192.168.0.124")
	tc = tests.Test{
		Method:  "POST",
		URL:     testURLRoot,
		Handler: createEnodeb,
		Payload: &lteModels.Enodeb{
			Config: &lteModels.EnodebConfiguration{
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
			EnodebConfig: &lteModels.EnodebConfig{
				ConfigType: "UNMANAGED",
				UnmanagedConfig: &lteModels.UnmanagedEnodebConfiguration{
					CellID:    swag.Uint32(1234),
					IPAddress: &ip,
					Tac:       swag.Uint32(1),
				},
			},
			Name:        "foobar",
			Description: "foobar description",
			Serial:      "unmanaged",
		},
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)
	_, err = configurator.LoadEntity("n1", lte.CellularEnodebEntityType, "unmanaged", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
}

func TestUpdateEnodeb(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/enodebs/:enodeb_serial"

	handlers := handlers.GetHandlers()
	updateEnodeb := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.PUT).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type:        lte.CellularEnodebEntityType,
				Key:         "abcdefg",
				Name:        "abc enodeb",
				Description: "abc enodeb description",
				PhysicalID:  "abcdefg",
				Config: &lteModels.EnodebConfig{
					ConfigType: "MANAGED",
					ManagedConfig: &lteModels.EnodebConfiguration{
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
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc := tests.Test{
		Method:  "PUT",
		URL:     testURLRoot,
		Handler: updateEnodeb,
		Payload: &lteModels.Enodeb{
			Config: &lteModels.EnodebConfiguration{
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
			EnodebConfig: &lteModels.EnodebConfig{
				ConfigType: "MANAGED",
				ManagedConfig: &lteModels.EnodebConfiguration{
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
			Name:        "foobar",
			Description: "new description",
			Serial:      "abcdefg",
		},
		ParamNames:     []string{"network_id", "enodeb_serial"},
		ParamValues:    []string{"n1", "abcdefg"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.CellularEnodebEntityType, "abcdefg", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.CellularEnodebEntityType, Key: "abcdefg",
		Name:        "foobar",
		Description: "new description",
		PhysicalID:  "abcdefg",
		GraphID:     "2",
		Config: &lteModels.EnodebConfig{
			ConfigType: "MANAGED",
			ManagedConfig: &lteModels.EnodebConfiguration{
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
		Version: 1,
	}
	assert.Equal(t, expected, actual)

	tc = tests.Test{
		Method:  "PUT",
		URL:     testURLRoot,
		Handler: updateEnodeb,
		Payload: &lteModels.Enodeb{
			Config: &lteModels.EnodebConfiguration{
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
			EnodebConfig: &lteModels.EnodebConfig{
				ConfigType: "MANAGED",
				ManagedConfig: &lteModels.EnodebConfiguration{
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
			Name:              "foobar",
			Serial:            "abcdefg",
			AttachedGatewayID: "gw1",
		},
		ParamNames:     []string{"network_id", "enodeb_serial"},
		ParamValues:    []string{"n1", "abcdefg"},
		ExpectedStatus: 400,
		ExpectedError:  "attached_gateway_id is a read-only property",
	}
	tests.RunUnitTest(t, e, tc)

	ip := strfmt.IPv4("192.168.0.124")
	tc = tests.Test{
		Method:  "PUT",
		URL:     testURLRoot,
		Handler: updateEnodeb,
		Payload: &lteModels.Enodeb{
			Config: &lteModels.EnodebConfiguration{
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
			EnodebConfig: &lteModels.EnodebConfig{
				ConfigType: "UNMANAGED",
				UnmanagedConfig: &lteModels.UnmanagedEnodebConfiguration{
					CellID:    swag.Uint32(1234),
					IPAddress: &ip,
					Tac:       swag.Uint32(1),
				},
			},
			Name:        "foobar",
			Description: "new description",
			Serial:      "abcdefg",
		},
		ParamNames:     []string{"network_id", "enodeb_serial"},
		ParamValues:    []string{"n1", "abcdefg"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
}

func TestDeleteEnodeb(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/enodebs/:enodeb_serial"

	handlers := handlers.GetHandlers()
	deleteEnodeb := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.DELETE).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type:       lte.CellularEnodebEntityType,
				Key:        "abcdefg",
				Name:       "abc enodeb",
				PhysicalID: "abcdefg",
				Config: &lteModels.EnodebConfig{
					ConfigType: "MANAGED",
					ManagedConfig: &lteModels.EnodebConfiguration{
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
			},
		},
		serdes.Entity,
	)
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

	_, err = configurator.LoadEntity("n1", lte.CellularEnodebEntityType, "abcdefg", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.EqualError(t, err, "Not found")
}

func TestGetEnodebState(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/enodebs/:enodeb_serial/state"

	handlers := handlers.GetHandlers()
	getEnodebState := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: lte.CellularEnodebEntityType, Key: "serial1",
				PhysicalID: "serial1",
			},
			{
				Type: orc8r.MagmadGatewayType, Key: "gw1",
				PhysicalID:   "hwid1",
				Associations: []storage.TypeAndKey{{Type: lte.CellularEnodebEntityType, Key: "serial1"}},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	// 404
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getEnodebState,
		ParamNames:     []string{"network_id", "enodeb_serial"},
		ParamValues:    []string{"n1", "serial1"},
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// report state
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.UnfreezeClock(t)

	// encode the appropriate certificate into context
	ctx := test_utils.GetContextWithCertificate(t, "hwid1")
	reportEnodebState(t, ctx, "serial1", lteModels.NewDefaultEnodebStatus())
	expected := lteModels.NewDefaultEnodebStatus()
	expected.TimeReported = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
	expected.ReportingGatewayID = "gw1"

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getEnodebState,
		ParamNames:     []string{"network_id", "enodeb_serial"},
		ParamValues:    []string{"n1", "serial1"},
		ExpectedStatus: 200,
		ExpectedResult: expected,
	}
	tests.RunUnitTest(t, e, tc)
}

func TestCreateApn(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/apns"
	handlers := handlers.GetHandlers()
	createApn := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.POST).HandlerFunc

	// default apn profile should always succeed
	payload := newAPN("foo")
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        payload,
		Handler:        createApn,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.APNEntityType, "foo", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.APNEntityType,
		Key:       "foo",
		Config:    payload.ApnConfiguration,
		GraphID:   "2",
	}
	assert.Equal(t, expected, actual)
}

func TestListApns(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/apns"
	handlers := handlers.GetHandlers()
	listApns := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listApns,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*lteModels.Apn{}),
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: lte.APNEntityType, Key: "oai.ipv4",
				Config: &lteModels.ApnConfiguration{
					Ambr: &lteModels.AggregatedMaximumBitrate{
						MaxBandwidthDl: swag.Uint32(200),
						MaxBandwidthUl: swag.Uint32(200),
					},
					QosProfile: &lteModels.QosProfile{
						ClassID:                 swag.Int32(9),
						PreemptionCapability:    swag.Bool(true),
						PreemptionVulnerability: swag.Bool(false),
						PriorityLevel:           swag.Uint32(15),
					},
				},
			},
			{
				Type: lte.APNEntityType, Key: "oai.ims",
				Config: &lteModels.ApnConfiguration{
					Ambr: &lteModels.AggregatedMaximumBitrate{
						MaxBandwidthDl: swag.Uint32(100),
						MaxBandwidthUl: swag.Uint32(100),
					},
					QosProfile: &lteModels.QosProfile{
						ClassID:                 swag.Int32(5),
						PreemptionCapability:    swag.Bool(true),
						PreemptionVulnerability: swag.Bool(false),
						PriorityLevel:           swag.Uint32(5),
					},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listApns,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*lteModels.Apn{
			"oai.ipv4": {
				ApnName: "oai.ipv4",
				ApnConfiguration: &lteModels.ApnConfiguration{
					Ambr: &lteModels.AggregatedMaximumBitrate{
						MaxBandwidthDl: swag.Uint32(200),
						MaxBandwidthUl: swag.Uint32(200),
					},
					QosProfile: &lteModels.QosProfile{
						ClassID:                 swag.Int32(9),
						PreemptionCapability:    swag.Bool(true),
						PreemptionVulnerability: swag.Bool(false),
						PriorityLevel:           swag.Uint32(15),
					},
				},
			},
			"oai.ims": {
				ApnName: "oai.ims",
				ApnConfiguration: &lteModels.ApnConfiguration{
					Ambr: &lteModels.AggregatedMaximumBitrate{
						MaxBandwidthDl: swag.Uint32(100),
						MaxBandwidthUl: swag.Uint32(100),
					},
					QosProfile: &lteModels.QosProfile{
						ClassID:                 swag.Int32(5),
						PreemptionCapability:    swag.Bool(true),
						PreemptionVulnerability: swag.Bool(false),
						PriorityLevel:           swag.Uint32(5),
					},
				},
			},
		}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetApn(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/apns/:apn_name"
	handlers := handlers.GetHandlers()
	getApn := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getApn,
		ParamNames:     []string{"network_id", "apn_name"},
		ParamValues:    []string{"n1", "oai.ipv4"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: lte.APNEntityType, Key: "oai.ipv4",
			Config: &lteModels.ApnConfiguration{
				Ambr: &lteModels.AggregatedMaximumBitrate{
					MaxBandwidthDl: swag.Uint32(200),
					MaxBandwidthUl: swag.Uint32(200),
				},
				QosProfile: &lteModels.QosProfile{
					ClassID:                 swag.Int32(9),
					PreemptionCapability:    swag.Bool(true),
					PreemptionVulnerability: swag.Bool(false),
					PriorityLevel:           swag.Uint32(15),
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getApn,
		ParamNames:     []string{"network_id", "apn_name"},
		ParamValues:    []string{"n1", "oai.ipv4"},
		ExpectedStatus: 200,
		ExpectedResult: &lteModels.Apn{
			ApnName: "oai.ipv4",
			ApnConfiguration: &lteModels.ApnConfiguration{
				Ambr: &lteModels.AggregatedMaximumBitrate{
					MaxBandwidthDl: swag.Uint32(200),
					MaxBandwidthUl: swag.Uint32(200),
				},
				QosProfile: &lteModels.QosProfile{
					ClassID:                 swag.Int32(9),
					PreemptionCapability:    swag.Bool(true),
					PreemptionVulnerability: swag.Bool(false),
					PriorityLevel:           swag.Uint32(15),
				},
			},
		},
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateApn(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/apns/:apn_name"
	handlers := handlers.GetHandlers()
	updateApn := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.PUT).HandlerFunc

	// 404
	payload := &lteModels.Apn{
		ApnName: "oai.ipv4",
		ApnConfiguration: &lteModels.ApnConfiguration{
			Ambr: &lteModels.AggregatedMaximumBitrate{
				MaxBandwidthDl: swag.Uint32(100),
				MaxBandwidthUl: swag.Uint32(100),
			},
			QosProfile: &lteModels.QosProfile{
				ClassID:                 swag.Int32(5),
				PreemptionCapability:    swag.Bool(true),
				PreemptionVulnerability: swag.Bool(false),
				PriorityLevel:           swag.Uint32(5),
			},
		},
	}

	tc := tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateApn,
		Payload:        payload,
		ParamNames:     []string{"network_id", "apn_name"},
		ParamValues:    []string{"n1", "oai.ipv4"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// Add the APN Configuration
	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: lte.APNEntityType, Key: "oai.ipv4",
			Config: &lteModels.ApnConfiguration{
				Ambr: &lteModels.AggregatedMaximumBitrate{
					MaxBandwidthDl: swag.Uint32(200),
					MaxBandwidthUl: swag.Uint32(200),
				},
				QosProfile: &lteModels.QosProfile{
					ClassID:                 swag.Int32(9),
					PreemptionCapability:    swag.Bool(true),
					PreemptionVulnerability: swag.Bool(false),
					PriorityLevel:           swag.Uint32(15),
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateApn,
		Payload:        payload,
		ParamNames:     []string{"network_id", "apn_name"},
		ParamValues:    []string{"n1", "oai.ipv4"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.APNEntityType, "oai.ipv4", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.APNEntityType,
		Key:       "oai.ipv4",
		Config:    payload.ApnConfiguration,
		GraphID:   "2",
		Version:   1,
	}
	assert.Equal(t, expected, actual)
}

func TestDeleteApn(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/apns/:apn_name"
	handlers := handlers.GetHandlers()
	deleteApn := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.DELETE).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: lte.APNEntityType, Key: "oai.ipv4",
				Config: &lteModels.ApnConfiguration{
					Ambr: &lteModels.AggregatedMaximumBitrate{
						MaxBandwidthDl: swag.Uint32(200),
						MaxBandwidthUl: swag.Uint32(200),
					},
					QosProfile: &lteModels.QosProfile{
						ClassID:                 swag.Int32(9),
						PreemptionCapability:    swag.Bool(true),
						PreemptionVulnerability: swag.Bool(false),
						PriorityLevel:           swag.Uint32(15),
					},
				},
			},
			{
				Type: lte.APNEntityType, Key: "oai.ims",
				Config: &lteModels.ApnConfiguration{
					Ambr: &lteModels.AggregatedMaximumBitrate{
						MaxBandwidthDl: swag.Uint32(100),
						MaxBandwidthUl: swag.Uint32(100),
					},
					QosProfile: &lteModels.QosProfile{
						ClassID:                 swag.Int32(5),
						PreemptionCapability:    swag.Bool(true),
						PreemptionVulnerability: swag.Bool(false),
						PriorityLevel:           swag.Uint32(5),
					},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc := tests.Test{
		Method:         "DELETE",
		URL:            testURLRoot,
		Handler:        deleteApn,
		ParamNames:     []string{"network_id", "apn_name"},
		ParamValues:    []string{"n1", "oai.ipv4"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, _, err := configurator.LoadAllEntitiesOfType("n1", lte.APNEntityType, configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(actual))
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.APNEntityType,
		Key:       "oai.ims",
		Config: &lteModels.ApnConfiguration{
			Ambr: &lteModels.AggregatedMaximumBitrate{
				MaxBandwidthDl: swag.Uint32(100),
				MaxBandwidthUl: swag.Uint32(100),
			},
			QosProfile: &lteModels.QosProfile{
				ClassID:                 swag.Int32(5),
				PreemptionCapability:    swag.Bool(true),
				PreemptionVulnerability: swag.Bool(false),
				PriorityLevel:           swag.Uint32(5),
			},
		},
		GraphID: "4",
		Version: 0,
	}
	assert.Equal(t, expected, actual[0])
}

func TestAPNResource(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n0", configurator.NetworkEntity{Type: orc8r.UpgradeTierEntityType, Key: "t0"}, serdes.Entity)
	assert.NoError(t, err)

	e := echo.New()
	urlBase := "/magma/v1/lte/:network_id/gateways"
	urlManage := urlBase + "/:gateway_id"
	lteHandlers := handlers.GetHandlers()
	getAllGateways := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlBase, obsidian.GET).HandlerFunc
	postGateway := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlBase, obsidian.POST).HandlerFunc
	putGateway := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlManage, obsidian.PUT).HandlerFunc
	getGateway := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlManage, obsidian.GET).HandlerFunc
	deleteGateway := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlManage, obsidian.DELETE).HandlerFunc

	postAPN := tests.GetHandlerByPathAndMethod(t, lteHandlers, "/magma/v1/lte/:network_id/apns", obsidian.POST).HandlerFunc
	deleteAPN := tests.GetHandlerByPathAndMethod(t, lteHandlers, "/magma/v1/lte/:network_id/apns/:apn_name", obsidian.DELETE).HandlerFunc

	gw := newMutableGateway("gw0")

	// Get all, initially empty
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways",
		Handler:        getAllGateways,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]lteModels.MutableLteGateway{}),
	}
	tests.RunUnitTest(t, e, tc)

	// Post err, APN names don't match
	gw.ApnResources = lteModels.ApnResources{"apn0": {ApnName: "apn1", ID: "res0", VlanID: 4}}
	tc = tests.Test{
		Method:                 "POST",
		URL:                    "/magma/v1/lte/n0/gateways",
		Payload:                gw,
		Handler:                postGateway,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n0"},
		ExpectedStatus:         400,
		ExpectedErrorSubstring: "APN resources key (apn0) and APN name (apn1) must match",
	}
	tests.RunUnitTest(t, e, tc)

	// Post err, APN doesn't exist
	gw.ApnResources = lteModels.ApnResources{"apn0": {ApnName: "apn0", ID: "res0", VlanID: 4}}
	tc = tests.Test{
		Method:                 "POST",
		URL:                    "/magma/v1/lte/n0/gateways",
		Payload:                gw,
		Handler:                postGateway,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n0"},
		ExpectedStatus:         500, // this would actually make more sense as a 400, but it's a non-trivial fix
		ExpectedErrorSubstring: `could not find entities matching [type:"apn" key:"apn0" ]`,
	}
	tests.RunUnitTest(t, e, tc)

	// Post APNs
	apn0 := newAPN("apn0")
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/:network_id/apns",
		Payload:        apn0,
		Handler:        postAPN,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)
	apn1 := newAPN("apn1")
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/:network_id/apns",
		Payload:        apn1,
		Handler:        postAPN,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)
	apn2 := newAPN("apn2")
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/:network_id/apns",
		Payload:        apn2,
		Handler:        postAPN,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Post, successful
	gw.ApnResources = lteModels.ApnResources{"apn0": {ApnName: "apn0", ID: "res0", VlanID: 4}}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n0/gateways",
		Payload:        gw,
		Handler:        postGateway,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Get all, posted gateway found
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways",
		Handler:        getAllGateways,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*lteModels.MutableLteGateway{"gw0": gw}),
	}
	tests.RunUnitTest(t, e, tc)

	// Put err, APN doesn't exist
	gw.ApnResources = lteModels.ApnResources{"apnXXX": {ApnName: "apnXXX", ID: "res0", VlanID: 4}}
	tc = tests.Test{
		Method:                 "PUT",
		URL:                    "/magma/v1/lte/n0/gateways/gw0",
		Payload:                gw,
		ParamNames:             []string{"network_id", "gateway_id"},
		ParamValues:            []string{"n0", "gw0"},
		Handler:                putGateway,
		ExpectedStatus:         500, // would make more sense as 400
		ExpectedErrorSubstring: `could not find entities matching [type:"apn" key:"apnXXX" ]`,
	}
	tests.RunUnitTest(t, e, tc)

	// Put err, mismatched APN names
	gw.ApnResources = lteModels.ApnResources{"apn1": {ApnName: "apnXXX", ID: "res0", VlanID: 4}}
	tc = tests.Test{
		Method:                 "PUT",
		URL:                    "/magma/v1/lte/n0/gateways/gw0",
		Payload:                gw,
		ParamNames:             []string{"network_id", "gateway_id"},
		ParamValues:            []string{"n0", "gw0"},
		Handler:                putGateway,
		ExpectedStatus:         400,
		ExpectedErrorSubstring: "APN resources key (apn1) and APN name (apnXXX) must match",
	}
	tests.RunUnitTest(t, e, tc)

	// Put err, request has duplicate resource IDs
	gw.ApnResources = lteModels.ApnResources{
		"apn0": {ApnName: "apn0", ID: "res0"},
		"apn1": {ApnName: "apn1", ID: "res0"},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		Payload:        gw,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        putGateway,
		ExpectedStatus: 400,
		ExpectedError:  "duplicate APN resource ID in request: res0",
	}
	tests.RunUnitTest(t, e, tc)

	// Put, point to new APN
	gw.ApnResources = lteModels.ApnResources{"apn1": {ApnName: "apn1", ID: "res0", VlanID: 4}}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		Payload:        gw,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        putGateway,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Post err, resource ID already exists
	gw1 := newMutableGateway("gw1")
	gw1.ApnResources = lteModels.ApnResources{"apn0": {ApnName: "apn0", ID: "res0"}}
	tc = tests.Test{
		Method:                 "POST",
		URL:                    "/magma/v1/lte/n0/gateways",
		Payload:                gw1,
		ParamNames:             []string{"network_id", "gateway_id"},
		ParamValues:            []string{"n0", "gw1"},
		Handler:                postGateway,
		ExpectedStatus:         500, // TODO(8/21/20): this should really be a 400
		ExpectedErrorSubstring: "an entity 'apn_resource-res0' already exists",
	}
	tests.RunUnitTest(t, e, tc)

	// Put, create new APN resource
	// TODO: make sure that this works
	gw.ApnResources = lteModels.ApnResources{
		"apn1": {ApnName: "apn1", ID: "res1", VlanID: 4},
		"apn2": {ApnName: "apn2", ID: "res2", VlanID: 4},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		Payload:        gw,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        putGateway,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Configurator confirms old APN resource was deleted
	exists, err := configurator.DoesEntityExist("n0", lte.APNResourceEntityType, "res0")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Get, changes are reflected
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        getGateway,
		ExpectedStatus: 200,
		ExpectedResult: gw,
	}
	tests.RunUnitTest(t, e, tc)

	// Delete
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        deleteGateway,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get err, not found
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        getGateway,
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// Configurator confirms APN resource was deleted
	exists, err = configurator.DoesEntityExist("n0", lte.APNResourceEntityType, "res1")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Configurator confirms all APN resources are now deleted
	ents, _, err := configurator.LoadAllEntitiesOfType("n0", lte.APNResourceEntityType, configurator.EntityLoadCriteria{}, serdes.Entity)
	assert.NoError(t, err)
	assert.Empty(t, ents)

	// Post, add gateway back
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n0/gateways",
		Payload:        gw,
		Handler:        postGateway,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Configurator confirms gw's APN resources exist again
	ents, _, err = configurator.LoadAllEntitiesOfType("n0", lte.APNResourceEntityType, configurator.EntityLoadCriteria{LoadConfig: true}, serdes.Entity)
	assert.NoError(t, err)
	assert.Len(t, ents, 2)
	assert.ElementsMatch(t, []string{"res1", "res2"}, []string{ents[0].Key, ents[1].Key})
	assert.ElementsMatch(t, []string{"res1", "res2"}, []string{(&lteModels.ApnResource{}).FromEntity(ents[0]).ID, (&lteModels.ApnResource{}).FromEntity(ents[1]).ID})

	// Delete linked APN
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/lte/n0/apns/apn1",
		Handler:        deleteAPN,
		ParamNames:     []string{"network_id", "apn_name"},
		ParamValues:    []string{"n0", "apn1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Configurator confirms gateway now has only 1 apn_resource assoc
	gwEnt, err := configurator.LoadEntity(
		"n0", lte.CellularGatewayEntityType, "gw0",
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Len(t, gwEnt.Associations.Filter(lte.APNResourceEntityType), 1)
	assert.Equal(t, "res2", gwEnt.Associations.Filter(lte.APNResourceEntityType).Keys()[0])

	// Configurator confirms APN resource was deleted due to cascading delete
	ents, _, err = configurator.LoadAllEntitiesOfType("n0", lte.APNResourceEntityType, configurator.EntityLoadCriteria{}, serdes.Entity)
	assert.NoError(t, err)
	assert.Len(t, ents, 1)
	assert.Equal(t, "res2", ents[0].Key)

	// Get, APN resource is gone
	gw.ApnResources = lteModels.ApnResources{"apn2": {ApnName: "apn2", ID: "res2", VlanID: 4}}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        getGateway,
		ExpectedStatus: 200,
		ExpectedResult: gw,
	}
	tests.RunUnitTest(t, e, tc)
}

// Regression test for issue #3088
func TestAPNResource_Regression_3088(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n0", configurator.NetworkEntity{Type: orc8r.UpgradeTierEntityType, Key: "t0"}, serdes.Entity)
	assert.NoError(t, err)

	e := echo.New()
	urlBase := "/magma/v1/lte/:network_id/gateways"
	urlManage := urlBase + "/:gateway_id"
	lteHandlers := handlers.GetHandlers()
	getAllGateways := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlBase, obsidian.GET).HandlerFunc
	postGateway := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlBase, obsidian.POST).HandlerFunc
	putGateway := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlManage, obsidian.PUT).HandlerFunc
	getGateway := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlManage, obsidian.GET).HandlerFunc

	postAPN := tests.GetHandlerByPathAndMethod(t, lteHandlers, "/magma/v1/lte/:network_id/apns", obsidian.POST).HandlerFunc

	gw0 := newMutableGateway("gw0")
	gw1 := newMutableGateway("gw1")

	// Get all, initially empty
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways",
		Handler:        getAllGateways,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]lteModels.MutableLteGateway{}),
	}
	tests.RunUnitTest(t, e, tc)

	// Post APN
	apn0 := newAPN("apn0")
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/:network_id/apns",
		Payload:        apn0,
		Handler:        postAPN,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Post gw0, successful
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n0/gateways",
		Payload:        gw0,
		Handler:        postGateway,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Post gw1, successful
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n0/gateways",
		Payload:        gw1,
		Handler:        postGateway,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Get all, posted gateway found
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways",
		Handler:        getAllGateways,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*lteModels.MutableLteGateway{"gw0": gw0, "gw1": gw1}),
	}
	tests.RunUnitTest(t, e, tc)

	// Put, add apn_resource to gw0
	gw0.ApnResources = lteModels.ApnResources{"apn0": {ApnName: "apn0", ID: "res0", VlanID: 4}}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		Payload:        gw0,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        putGateway,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get, changes are reflected
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        getGateway,
		ExpectedStatus: 200,
		ExpectedResult: gw0,
	}
	tests.RunUnitTest(t, e, tc)

	// Get all, only gw0 has an apn_resource
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways",
		Handler:        getAllGateways,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*lteModels.MutableLteGateway{"gw0": gw0, "gw1": gw1}),
	}
	tests.RunUnitTest(t, e, tc)
}

// Regression test for issue #3149
func TestAPNResource_Regression_3149(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n0", configurator.NetworkEntity{Type: orc8r.UpgradeTierEntityType, Key: "t0"}, serdes.Entity)
	assert.NoError(t, err)

	e := echo.New()
	urlBase := "/magma/v1/lte/:network_id/gateways"
	urlManage := urlBase + "/:gateway_id"
	lteHandlers := handlers.GetHandlers()
	getAllGateways := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlBase, obsidian.GET).HandlerFunc
	postGateway := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlBase, obsidian.POST).HandlerFunc
	putGateway := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlManage, obsidian.PUT).HandlerFunc
	getGateway := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlManage, obsidian.GET).HandlerFunc
	deleteGateway := tests.GetHandlerByPathAndMethod(t, lteHandlers, urlManage, obsidian.DELETE).HandlerFunc

	postAPN := tests.GetHandlerByPathAndMethod(t, lteHandlers, "/magma/v1/lte/:network_id/apns", obsidian.POST).HandlerFunc

	// Create enb0
	_, err = configurator.CreateEntities("n0", []configurator.NetworkEntity{{Type: lte.CellularEnodebEntityType, Key: "enb0"}}, serdes.Entity)
	assert.NoError(t, err)

	gw0 := newMutableGateway("gw0")

	// Get all, initially empty
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways",
		Handler:        getAllGateways,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]lteModels.MutableLteGateway{}),
	}
	tests.RunUnitTest(t, e, tc)

	// Post 2 APNs
	apnInternet := newAPN("internet")
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/:network_id/apns",
		Payload:        apnInternet,
		Handler:        postAPN,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)
	apnManagement := newAPN("management")
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/:network_id/apns",
		Payload:        apnManagement,
		Handler:        postAPN,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Post gw0 with 2 apn_resources, successful
	gw0.ApnResources = lteModels.ApnResources{
		"internet": {
			ApnName:    "internet",
			GatewayIP:  "192.168.10.1",
			GatewayMac: "e0:63:da:22:47:21",
			ID:         "internet_apn_resource",
			VlanID:     13,
		},
		"management": {
			ApnName:    "management",
			GatewayIP:  "192.168.9.1",
			GatewayMac: "e0:63:da:22:47:21",
			ID:         "management_apn_resource",
			VlanID:     12,
		},
	}
	gw0.ConnectedEnodebSerials = []string{"enb0"}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n0/gateways",
		Payload:        gw0,
		Handler:        postGateway,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Get all, posted gateway found
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways",
		Handler:        getAllGateways,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*lteModels.MutableLteGateway{"gw0": gw0}),
	}
	tests.RunUnitTest(t, e, tc)

	// Get posted gateway, found
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        getGateway,
		ExpectedStatus: 200,
		ExpectedResult: gw0,
	}
	tests.RunUnitTest(t, e, tc)

	// Delete
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        deleteGateway,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get, not found
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        getGateway,
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// Post gw0 with 0 apn_resources, successful
	gw0.ApnResources = lteModels.ApnResources{}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n0/gateways",
		Payload:        gw0,
		Handler:        postGateway,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Get posted gateway, found
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        getGateway,
		ExpectedStatus: 200,
		ExpectedResult: gw0,
	}
	tests.RunUnitTest(t, e, tc)

	// Put, add 2 apn_resources to gw0
	gw0.ApnResources = lteModels.ApnResources{
		"internet": {
			ApnName:    "internet",
			GatewayIP:  "192.168.10.1",
			GatewayMac: "e0:63:da:22:47:21",
			ID:         "internet_apn_resource",
			VlanID:     13,
		},
		"management": {
			ApnName:    "management",
			GatewayIP:  "192.168.9.1",
			GatewayMac: "e0:63:da:22:47:21",
			ID:         "management_apn_resource",
			VlanID:     12,
		},
	}
	gw0.ConnectedEnodebSerials = []string{"enb0"}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		Payload:        gw0,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        putGateway,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get, changes are reflected
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways/gw0",
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n0", "gw0"},
		Handler:        getGateway,
		ExpectedStatus: 200,
		ExpectedResult: gw0,
	}
	tests.RunUnitTest(t, e, tc)

	// Get all, changes are reflected
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/gateways",
		Handler:        getAllGateways,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*lteModels.MutableLteGateway{"gw0": gw0}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestHAGatewayPools(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)

	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	listHaPools := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/gateway_pools", obsidian.GET).HandlerFunc
	createHaPool := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/gateway_pools", obsidian.POST).HandlerFunc
	getHaPool := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/gateway_pools/:gateway_pool_id", obsidian.GET).HandlerFunc
	updateHaPool := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/gateway_pools/:gateway_pool_id", obsidian.PUT).HandlerFunc
	deleteHaPool := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/gateway_pools/:gateway_pool_id", obsidian.DELETE).HandlerFunc

	getPoolRecord := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/gateways/:gateway_id/cellular/pooling", obsidian.GET).HandlerFunc
	updatePoolRecord := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/gateways/:gateway_id/cellular/pooling", obsidian.PUT).HandlerFunc

	seedNetworks(t)

	// Test List HA Pairs empty
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n1/gateway_pools",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listHaPools,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*lteModels.CellularGatewayPool{}),
	}
	tests.RunUnitTest(t, e, tc)

	pool1 := &lteModels.MutableCellularGatewayPool{
		GatewayPoolID:   lteModels.GatewayPoolID("pool1"),
		GatewayPoolName: "pool 1",
		Config: &lteModels.CellularGatewayPoolConfigs{
			MmeGroupID: 1,
		},
	}

	// Create pool1
	gatewayPoolsURLRoot := "/magma/v1/lte/:network_id/gateway/gateway_pools"
	tc = tests.Test{
		Method:         "POST",
		URL:            gatewayPoolsURLRoot,
		Payload:        tests.JSONMarshaler(pool1),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createHaPool,
		ExpectedStatus: 201,
		ExpectedResult: tests.JSONMarshaler("pool1"),
	}
	tests.RunUnitTest(t, e, tc)

	seedTier(t, "n1")
	seedGateway(t, "n1", "g1")

	poolRecords := []lteModels.CellularGatewayPoolRecord{
		{
			GatewayPoolID:       "pool4",
			MmeCode:             1,
			MmeRelativeCapacity: 10,
		},
	}

	// Create fails as pool4 doesn't exist
	poolingURLRoot := "/magma/v1/lte/:network_id/gateways/:gateway_id/pooling"
	tc = tests.Test{
		Method:         "PUT",
		URL:            poolingURLRoot,
		Payload:        tests.JSONMarshaler(poolRecords),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		Handler:        updatePoolRecord,
		ExpectedStatus: 400,
		ExpectedError:  "Gateway pool pool4 does not exist",
	}
	tests.RunUnitTest(t, e, tc)

	// Create succeeds with pool1
	poolRecords[0].GatewayPoolID = "pool1"

	tc = tests.Test{
		Method:         "PUT",
		URL:            poolingURLRoot,
		Payload:        tests.JSONMarshaler(poolRecords),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		Handler:        updatePoolRecord,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get pool records
	tc = tests.Test{
		Method:         "GET",
		URL:            poolingURLRoot,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		Handler:        getPoolRecord,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(poolRecords),
	}
	tests.RunUnitTest(t, e, tc)

	// Get HA Pool
	expectedPool := &lteModels.CellularGatewayPool{
		GatewayPoolID:   lteModels.GatewayPoolID("pool1"),
		GatewayPoolName: "pool 1",
		Config: &lteModels.CellularGatewayPoolConfigs{
			MmeGroupID: 1,
		},
		GatewayIds: []models2.GatewayID{
			"g1",
		},
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/:gateway_pool_id", gatewayPoolsURLRoot),
		ParamNames:     []string{"network_id", "gateway_pool_id"},
		ParamValues:    []string{"n1", "pool1"},
		Handler:        getHaPool,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedPool),
	}
	tests.RunUnitTest(t, e, tc)

	// Update HA Pool
	expectedPool.GatewayPoolName = "pool 1 updated"
	expectedPool.Config = &lteModels.CellularGatewayPoolConfigs{MmeGroupID: 4}

	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/:gateway_pool_id", gatewayPoolsURLRoot),
		Payload:        tests.JSONMarshaler(expectedPool),
		ParamNames:     []string{"network_id", "gateway_pool_id"},
		ParamValues:    []string{"n1", "pool1"},
		Handler:        updateHaPool,
		ExpectedStatus: 201,
		ExpectedResult: tests.JSONMarshaler("pool1"),
	}
	tests.RunUnitTest(t, e, tc)

	// Ensure update succeeded
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/:gateway_pool_id", gatewayPoolsURLRoot),
		ParamNames:     []string{"network_id", "gateway_pool_id"},
		ParamValues:    []string{"n1", "pool1"},
		Handler:        getHaPool,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedPool),
	}
	tests.RunUnitTest(t, e, tc)

	// Create pool2
	pool2 := &lteModels.MutableCellularGatewayPool{
		GatewayPoolID:   lteModels.GatewayPoolID("pool2"),
		GatewayPoolName: "pool2",
		Config: &lteModels.CellularGatewayPoolConfigs{
			MmeGroupID: 1,
		},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            gatewayPoolsURLRoot,
		Payload:        tests.JSONMarshaler(pool2),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createHaPool,
		ExpectedStatus: 201,
		ExpectedResult: tests.JSONMarshaler("pool2"),
	}
	tests.RunUnitTest(t, e, tc)

	// Update g1 to reside in pool2
	poolRecords[0].GatewayPoolID = "pool2"
	tc = tests.Test{
		Method:         "PUT",
		URL:            poolingURLRoot,
		Payload:        tests.JSONMarshaler(poolRecords),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		Handler:        updatePoolRecord,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Ensure update moved the gateway properly
	expectedPool.GatewayIds = []models2.GatewayID{}
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/:gateway_pool_id", gatewayPoolsURLRoot),
		ParamNames:     []string{"network_id", "gateway_pool_id"},
		ParamValues:    []string{"n1", "pool1"},
		Handler:        getHaPool,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedPool),
	}
	tests.RunUnitTest(t, e, tc)

	expectedPool.GatewayIds = []models2.GatewayID{"g1"}
	expectedPool.GatewayPoolID = "pool2"
	expectedPool.GatewayPoolName = "pool2"
	expectedPool.Config.MmeGroupID = 1
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/:gateway_pool_id", gatewayPoolsURLRoot),
		ParamNames:     []string{"network_id", "gateway_pool_id"},
		ParamValues:    []string{"n1", "pool2"},
		Handler:        getHaPool,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedPool),
	}
	tests.RunUnitTest(t, e, tc)

	// Create pool3
	pool3 := &lteModels.MutableCellularGatewayPool{
		GatewayPoolID:   lteModels.GatewayPoolID("pool3"),
		GatewayPoolName: "pool3",
		Config: &lteModels.CellularGatewayPoolConfigs{
			MmeGroupID: 2,
		},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            gatewayPoolsURLRoot,
		Payload:        tests.JSONMarshaler(pool3),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createHaPool,
		ExpectedStatus: 201,
		ExpectedResult: tests.JSONMarshaler("pool3"),
	}
	tests.RunUnitTest(t, e, tc)

	// Adding gateway to pool3 should fail as it has a different MME GID
	pool3Record := lteModels.CellularGatewayPoolRecord{
		GatewayPoolID:       "pool3",
		MmeCode:             1,
		MmeRelativeCapacity: 10,
	}
	poolRecords = append(poolRecords, pool3Record)
	tc = tests.Test{
		Method:         "PUT",
		URL:            poolingURLRoot,
		Payload:        tests.JSONMarshaler(poolRecords),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		Handler:        updatePoolRecord,
		ExpectedStatus: 400,
		ExpectedError:  "Adding a gateway to pools with different MME group ID's (2), (1) is currently unsupported",
	}
	tests.RunUnitTest(t, e, tc)

	// Ensure pool deletion fails if gateway still resides in it
	tc = tests.Test{
		Method:         "DELETE",
		URL:            fmt.Sprintf("%s/:gateway_pool_id", gatewayPoolsURLRoot),
		ParamNames:     []string{"network_id", "gateway_pool_id"},
		ParamValues:    []string{"n1", "pool2"},
		Handler:        deleteHaPool,
		ExpectedStatus: 400,
		ExpectedError:  "Gateways [g1] still exist in pool pool2. All gateways must first be removed from the pool before it can be deleted",
	}
	tests.RunUnitTest(t, e, tc)

	// Remove pool record from gateway
	poolRecords = []lteModels.CellularGatewayPoolRecord{}
	tc = tests.Test{
		Method:         "PUT",
		URL:            poolingURLRoot,
		Payload:        tests.JSONMarshaler(poolRecords),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		Handler:        updatePoolRecord,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Now delete should succeed
	tc = tests.Test{
		Method:         "DELETE",
		URL:            fmt.Sprintf("%s/:gateway_pool_id", gatewayPoolsURLRoot),
		ParamNames:     []string{"network_id", "gateway_pool_id"},
		ParamValues:    []string{"n1", "pool2"},
		Handler:        deleteHaPool,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
}

func reportEnodebState(t *testing.T, ctx context.Context, enodebSerial string, req *lteModels.EnodebState) {
	client, err := state.GetStateClient()
	assert.NoError(t, err)

	serializedEnodebState, err := serde.Serialize(req, lte.EnodebStateType, serdes.State)
	assert.NoError(t, err)
	states := []*protos.State{
		{
			Type:     lte.EnodebStateType,
			DeviceID: enodebSerial,
			Value:    serializedEnodebState,
		},
	}
	_, err = client.ReportStates(
		ctx,
		&protos.ReportStatesRequest{States: states},
	)
	assert.NoError(t, err)
}

// n1, n3 are lte networks, n2 is not
func seedNetworks(t *testing.T) {
	_, err := configurator.CreateNetworks(
		[]configurator.Network{
			{
				ID:          "n1",
				Type:        lte.NetworkType,
				Name:        "foobar",
				Description: "Foo Bar",
				Configs: map[string]interface{}{
					lte.CellularNetworkConfigType: lteModels.NewDefaultTDDNetworkConfig(),
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
				Type:        lte.NetworkType,
				Name:        "barfoo",
				Description: "Bar Foo",
				Configs:     map[string]interface{}{},
			},
		},
		serdes.Network,
	)
	assert.NoError(t, err)
}

func seedGateway(t *testing.T, networkID string, gatewayID string) {
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	createGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/lte/:network_id/gateways", obsidian.POST).HandlerFunc

	gw := newMutableGateway(gatewayID)
	tc := tests.Test{
		Method:         "POST",
		URL:            fmt.Sprintf("/magma/v1/lte/%s/gateways", networkID),
		Handler:        createGateway,
		Payload:        gw,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

}

func seedTier(t *testing.T, networkID string) {
	// setup fixtures in backend
	_, err := configurator.CreateEntities(
		networkID,
		[]configurator.NetworkEntity{
			{Type: orc8r.UpgradeTierEntityType, Key: "t0"},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
}

func newDefaultGatewayConfig() *lteModels.GatewayCellularConfigs {
	return &lteModels.GatewayCellularConfigs{
		Ran: &lteModels.GatewayRanConfigs{
			Pci:             260,
			TransmitEnabled: swag.Bool(true),
		},
		Epc: &lteModels.GatewayEpcConfigs{
			NatEnabled: swag.Bool(true),
			IPBlock:    "192.168.128.0/24",
		},
		NonEpsService: &lteModels.GatewayNonEpsConfigs{
			CsfbMcc:              "001",
			CsfbMnc:              "01",
			Lac:                  swag.Uint32(1),
			CsfbRat:              swag.Uint32(0),
			Arfcn2g:              []uint32{},
			NonEpsServiceControl: swag.Uint32(0),
		},
		HeConfig: &lteModels.GatewayHeConfig{
			EnableHeaderEnrichment: swag.Bool(true),
			EnableEncryption:       swag.Bool(false),
			HeEncryptionAlgorithm:  lteModels.GatewayHeConfigHeEncryptionAlgorithmRC4,
			HeHashFunction:         lteModels.GatewayHeConfigHeHashFunctionMD5,
			HeEncodingType:         lteModels.GatewayHeConfigHeEncodingTypeBASE64,
		},
	}
}

func newAPN(name string) *lteModels.Apn {
	apn := &lteModels.Apn{
		ApnName: lteModels.ApnName(name),
		ApnConfiguration: &lteModels.ApnConfiguration{
			Ambr: &lteModels.AggregatedMaximumBitrate{
				MaxBandwidthDl: swag.Uint32(100),
				MaxBandwidthUl: swag.Uint32(100),
			},
			QosProfile: &lteModels.QosProfile{
				ClassID:                 swag.Int32(9),
				PreemptionCapability:    swag.Bool(true),
				PreemptionVulnerability: swag.Bool(false),
				PriorityLevel:           swag.Uint32(15),
			},
		},
	}

	return apn
}

func newMutableGateway(id string) *lteModels.MutableLteGateway {
	gw := &lteModels.MutableLteGateway{
		Device: &models.GatewayDevice{
			HardwareID: id + "_hwid",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		ID:          models2.GatewayID(id),
		Name:        "foobar",
		Description: "foo bar",
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          10,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		Cellular:               newDefaultGatewayConfig(),
		ConnectedEnodebSerials: []string{},
		Tier:                   "t0",
		ApnResources:           lteModels.ApnResources{},
	}
	return gw
}
