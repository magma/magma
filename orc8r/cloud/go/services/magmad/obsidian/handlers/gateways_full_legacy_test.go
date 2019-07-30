/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	magmadh "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory/mocks"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestGetViewsForNetworkLegacy(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "0")
	// Set up test
	mockStore := &mocks.FullGatewayViewFactory{}
	obsidian.TLS = false

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

func TestGetViewsForNetworkEmptyResponseLegacy(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "0")
	mockStore := &mocks.FullGatewayViewFactory{}
	obsidian.TLS = false

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

func TestGetGatewayViews_QueryType1Legacy(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "0")
	testGetGatewayViews(t, "gateway_ids=gw0,gw1,badgw")
}

func TestGetGatewayViews_QueryType2Legacy(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "0")
	testGetGatewayViews(t, "gateway_ids[0]=gw0&gateway_ids[1]=gw1&gateway_ids[2]=badgw")
}
