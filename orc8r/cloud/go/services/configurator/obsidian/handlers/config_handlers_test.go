/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorh "magma/orc8r/cloud/go/services/configurator/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	fooConfigType = "foo"
	networkID     = "networkID"
)

func TestGetNetworkConfigCRUDHandlers(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)
	fooNetworkSerde := configurator.NewNetworkConfigSerde(fooConfigType, &FooConfigs{})
	err := serde.RegisterSerdes(fooNetworkSerde)
	assert.NoError(t, err)
	restPort := tests.StartObsidian(t)
	e := echo.New()

	// Test GetCreateNetworkConfigHandler
	_, err = configurator.CreateNetworks(
		[]configurator.Network{{
			ID: networkID,
		}})
	assert.NoError(t, err)

	config := FooConfigs{
		ConfigNum: 100,
		ConfigStr: "hello!",
	}
	post, err := json.Marshal(config)
	assert.NoError(t, err)

	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(string(post)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	addParametersToContext(c, networkID, fooConfigType)

	// Success
	url := getURL(restPort, networkID, fooConfigType)
	err = configuratorh.GetCreateNetworkConfigHandler(url, fooConfigType, &FooConfigs{}).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	assertConfigExists(t, networkID, fooConfigType, &config)

	// Test GetReadNetworkConfigsHandler
	req = httptest.NewRequest(echo.GET, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	addParametersToContext(c, networkID, fooConfigType)

	// Success
	url = getURL(restPort, networkID, fooConfigType)
	err = configuratorh.GetReadNetworkConfigHandler(url, fooConfigType).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	actual := FooConfigs{}
	err = json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.NoError(t, err)
	assert.Equal(t, config, actual)

	// Test GetUpdateNetworkConfigsHandler
	updatedConfig := FooConfigs{
		ConfigNum: 32,
		ConfigStr: "goodbye!",
	}
	post, err = json.Marshal(updatedConfig)
	assert.NoError(t, err)

	req = httptest.NewRequest(echo.PUT, "/", strings.NewReader(string(post)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	addParametersToContext(c, networkID, fooConfigType)

	// Success
	url = getURL(restPort, networkID, fooConfigType)
	err = configuratorh.GetUpdateNetworkConfigHandler(url, fooConfigType, &FooConfigs{}).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	assertConfigExists(t, networkID, fooConfigType, &updatedConfig)

	// TestGetDeleteConfigHandler
	req = httptest.NewRequest(echo.DELETE, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	addParametersToContext(c, networkID, fooConfigType)

	// Success
	url = getURL(restPort, networkID, fooConfigType)
	err = configuratorh.GetDeleteNetworkConfigHandler(url, fooConfigType).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	assertConfigDoesNotExist(t, networkID, fooConfigType)
}

func addParametersToContext(c echo.Context, networkID string, configType string) echo.Context {
	c.SetParamNames("network_id", "config_type")
	c.SetParamValues(networkID, configType)
	return c
}

func assertConfigExists(t *testing.T, networkID string, configType string, config interface{}) {
	networks, notFound, err := configurator.LoadNetworks([]string{networkID}, false, true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(notFound))
	retrievedConfig := networks[0].Configs[configType]
	assert.Equal(t, config, retrievedConfig)
}

func assertConfigDoesNotExist(t *testing.T, networkID string, configType string) {
	networks, notFound, err := configurator.LoadNetworks([]string{networkID}, true, true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(notFound))
	_, ok := networks[0].Configs[configType]
	assert.Equal(t, false, ok)
}

func getURL(restPort int, networkID string, configType string) string {
	url := fmt.Sprintf(
		"http://localhost:%d%s/networks/%s/configs/%s",
		restPort,
		obsidian.RestRoot,
		networkID,
		configType,
	)
	return url
}

type FooConfigs struct {
	ConfigStr string `json:"config_str"`
	ConfigNum int    `json:"config_num"`
}

// MarshalBinary interface implementation
func (f *FooConfigs) MarshalBinary() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return swag.WriteJSON(f)
}

func (f *FooConfigs) UnmarshalBinary(b []byte) error {
	var res FooConfigs
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*f = res
	return nil
}
