/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"magma/orc8r/cloud/go/obsidian/config"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/services/checkind/test_utils"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	magmadh "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory/mocks"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	state_test_utils "magma/orc8r/cloud/go/services/state/test_utils"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetViewsForNetwork(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	// Set up test
	mockStore := &mocks.FullGatewayViewFactory{}
	config.TLS = false

	// Generate input/output objects
	networkID := "net1"
	gatewayStates := map[string]*view_factory.GatewayState{
		"gw0": {GatewayID: "gw0"},
		"gw1": {GatewayID: "gw1"},
	}
	modelStates := []*view_factory.GatewayStateType{
		{GatewayID: "gw0"},
		{GatewayID: "gw1"},
	}

	// Set up mock and get request handler
	mockStore.On("GetGatewayViewsForNetwork", networkID).Return(gatewayStates, nil)

	// Generate http request
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues(networkID)

	// Execute test
	err := magmadh.ListFullGatewayViewsLegacy(c, mockStore)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify results
	var actualModelStates []*view_factory.GatewayStateType
	err = json.Unmarshal(rec.Body.Bytes(), &actualModelStates)
	assert.NoError(t, err)
	assert.ElementsMatch(t, modelStates, actualModelStates)
	mockStore.AssertNumberOfCalls(t, "GetGatewayViewsForNetwork", 1)
}

func TestGetViewsForNetworkEmptyResponse(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	mockStore := &mocks.FullGatewayViewFactory{}
	config.TLS = false

	networkID := "badid"

	mockStore.On("GetGatewayViewsForNetwork", networkID).Return(map[string]*view_factory.GatewayState{}, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues(networkID)

	err := magmadh.ListFullGatewayViewsLegacy(c, mockStore)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var actualModelStates []*view_factory.GatewayStateType
	err = json.Unmarshal(rec.Body.Bytes(), &actualModelStates)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(actualModelStates))
	mockStore.AssertNumberOfCalls(t, "GetGatewayViewsForNetwork", 1)
}

func TestGetViewsForNetwork_Full(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	restPort := tests.StartObsidian(t)

	testURLRoot := fmt.Sprintf(
		"http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)
	networkID := "magmad_obsidian_test_network"
	registerNetworkWithIDTestCase := tests.Testcase{
		Name:                      "Register Network with Requested ID",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=%s", testURLRoot, networkID),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
		Expected:                  `"magmad_obsidian_test_network"`,
	}
	tests.RunTest(t, registerNetworkWithIDTestCase)

	// Test Register AG with requestedId
	requestedAGID := "my_gateway-1"
	registerAGWithIDTestCase := tests.Testcase{
		Name:   "Register AG with Requested ID",
		Method: "POST",
		Url: fmt.Sprintf(
			"%s/%s/gateways?requested_id=%s", testURLRoot, networkID, requestedAGID),
		Payload:  `{"hw_id":{"ID":"TestAGHwId00001"}, "name": "Test AG Name",  "key": {"key_type": "ECHO"}}`,
		Expected: fmt.Sprintf(`"%s"`, requestedAGID),
	}
	tests.RunTest(t, registerAGWithIDTestCase)

	getGatewaysFullView := tests.Testcase{
		Name:   "Get Gateways Full View",
		Method: "GET",
		Url: fmt.Sprintf(
			"%s/%s/gateways?view=full", testURLRoot, networkID),
		Payload:  "",
		Expected: `[{"config":{"magmad_gateway":null},"gateway_id":"my_gateway-1","record":{"hw_id":{"id":"TestAGHwId00001"},"key":{"key_type":"ECHO"},"name":"Test AG Name"},"status":null}]`,
	}
	tests.RunTest(t, getGatewaysFullView)

	expCfg := NewDefaultGatewayConfig()
	marshaledCfg, err := expCfg.MarshalBinary()
	assert.NoError(t, err)
	expectedCfgStr := string(marshaledCfg)
	// Test Setting (Updating) AG Configs With An Unregistered Tier
	setAGConfigTestCase := tests.Testcase{
		Name:   "Set AG Configs With Tier",
		Method: "POST",
		Url: fmt.Sprintf("%s/%s/gateways/%s/configs",
			testURLRoot, networkID, requestedAGID),
		Payload:  expectedCfgStr,
		Expected: `"my_gateway-1"`,
	}
	tests.RunTest(t, setAGConfigTestCase)

	getGatewaysFullView = tests.Testcase{
		Name:   "Get Gateways Full View",
		Method: "GET",
		Url: fmt.Sprintf(
			"%s/%s/gateways?view=full", testURLRoot, networkID),
		Payload:  "",
		Expected: fmt.Sprintf(`[{"config":{"magmad_gateway":%s},"gateway_id":"my_gateway-1","record":{"hw_id":{"id":"TestAGHwId00001"},"key":{"key_type":"ECHO"},"name":"Test AG Name"},"status":null}]`, expectedCfgStr),
	}
	tests.RunTest(t, getGatewaysFullView)

	// Test Gateway Full View with state
	ctx := state_test_utils.GetContextWithCertificate(t, "TestAGHwId00001")
	gwStatus := test_utils.GetGatewayStatusSwaggerFixture("TestAGHwId00001")
	state_test_utils.ReportGatewayStatus(t, ctx, gwStatus)
	status, response, err := tests.SendHttpRequest("GET", fmt.Sprintf("%s/%s/gateways?view=full", testURLRoot, networkID), "")
	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	gatewayStatesAndConfigs := []*view_factory.GatewayStateType{}
	err = json.Unmarshal([]byte(response), &gatewayStatesAndConfigs)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(gatewayStatesAndConfigs))
	assert.NotNil(t, gatewayStatesAndConfigs[0].Status)
	// 0 out timestamp and cert expiration time
	gatewayStatesAndConfigs[0].Status.CheckinTime = 0
	gatewayStatesAndConfigs[0].Status.CertExpirationTime = 0
	assert.Equal(t, gwStatus, gatewayStatesAndConfigs[0].Status)
}

func TestGetGatewayViews_QueryType1(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	testGetGatewayViews(t, "gateway_ids=gw0,gw1,badgw")
}

func TestGetGatewayViews_QueryType2(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	testGetGatewayViews(t, "gateway_ids[0]=gw0&gateway_ids[1]=gw1&gateway_ids[2]=badgw")
}

func testGetGatewayViews(t *testing.T, queryString string) {
	mockStore := &mocks.FullGatewayViewFactory{}
	config.TLS = false

	networkID := "net1"
	gatewayIDs := []string{"gw0", "gw1", "badgw"}
	gatewayStates := map[string]*view_factory.GatewayState{
		"gw0": {GatewayID: "gw0"},
		"gw1": {GatewayID: "gw1"},
	}
	modelStates := []*view_factory.GatewayStateType{
		{GatewayID: "gw0"},
		{GatewayID: "gw1"},
	}

	mockStore.On("GetGatewayViews", networkID, mock.MatchedBy(func(input []string) bool {
		return assert.ElementsMatch(t, gatewayIDs, input)
	})).Return(gatewayStates, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.URL.RawQuery = queryString
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues(networkID)

	err := magmadh.ListFullGatewayViewsLegacy(c, mockStore)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var actualModelStates []*view_factory.GatewayStateType
	err = json.Unmarshal(rec.Body.Bytes(), &actualModelStates)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(actualModelStates))
	assert.ElementsMatch(t, modelStates, actualModelStates)
	mockStore.AssertNumberOfCalls(t, "GetGatewayViews", 1)
}
