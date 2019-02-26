/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"magma/orc8r/cloud/go/obsidian/config"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory/mocks"
	"magma/orc8r/cloud/go/services/magmad/obsidian/models"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetViewsForNetwork(t *testing.T) {
	// Set up test
	mockStore := &mocks.FullGatewayViewFactory{}
	config.TLS = false

	// Generate input/output objects
	networkID := "net1"
	gatewayStates := map[string]*view_factory.GatewayState{
		"gw0": {GatewayID: "gw0"},
		"gw1": {GatewayID: "gw1"},
	}
	modelStates := []*models.GatewayStateType{
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
	err := handlers.ListFullGatewayViews(c, mockStore)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify results
	var actualModelStates []*models.GatewayStateType
	err = json.Unmarshal(rec.Body.Bytes(), &actualModelStates)
	assert.NoError(t, err)
	assert.ElementsMatch(t, modelStates, actualModelStates)
	mockStore.AssertNumberOfCalls(t, "GetGatewayViewsForNetwork", 1)
}

func TestGetViewsForNetworkEmptyResponse(t *testing.T) {
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

	err := handlers.ListFullGatewayViews(c, mockStore)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var actualModelStates []*models.GatewayStateType
	err = json.Unmarshal(rec.Body.Bytes(), &actualModelStates)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(actualModelStates))
	mockStore.AssertNumberOfCalls(t, "GetGatewayViewsForNetwork", 1)
}

func TestGetGatewayViews_QueryType1(t *testing.T) {
	testGetGatewayViews(t, "gateway_ids=gw0,gw1,badgw")
}

func TestGetGatewayViews_QueryType2(t *testing.T) {
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
	modelStates := []*models.GatewayStateType{
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

	err := handlers.ListFullGatewayViews(c, mockStore)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var actualModelStates []*models.GatewayStateType
	err = json.Unmarshal(rec.Body.Bytes(), &actualModelStates)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(actualModelStates))
	assert.ElementsMatch(t, modelStates, actualModelStates)
	mockStore.AssertNumberOfCalls(t, "GetGatewayViews", 1)
}
