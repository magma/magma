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

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorh "magma/orc8r/cloud/go/services/configurator/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	fooSerdeType = "foo"
	networkID    = "networkID"
)

func TestGetNetworkConfigCRUDHandlers(t *testing.T) {
	test_init.StartTestService(t)
	fooSerde := fooSerde{}
	err := serde.RegisterSerdes(&fooSerde)
	assert.NoError(t, err)
	restPort := tests.StartObsidian(t)
	e := echo.New()

	// Test GetCreateNetworkConfigHandler
	_, err = configurator.CreateNetworks(
		[]*protos.Network{{
			Id: networkID,
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
	addParametersToContext(c, networkID, fooSerdeType)

	// Success
	url := getURL(restPort, networkID, fooSerdeType)
	err = configuratorh.GetCreateNetworkConfigHandler(url, &fooSerde).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	assertConfigExists(t, networkID, fooSerdeType, config)

	// Test GetReadNetworkConfigsHandler
	req = httptest.NewRequest(echo.GET, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	addParametersToContext(c, networkID, fooSerdeType)

	// Success
	url = getURL(restPort, networkID, fooSerdeType)
	err = configuratorh.GetReadNetworkConfigHandler(url).HandlerFunc(c)
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
	addParametersToContext(c, networkID, fooSerdeType)

	// Success
	url = getURL(restPort, networkID, fooSerdeType)
	err = configuratorh.GetUpdateNetworkConfigHandler(url, &fooSerde).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	assertConfigExists(t, networkID, fooSerdeType, updatedConfig)

	// TestGetDeleteConfigHandler
	req = httptest.NewRequest(echo.DELETE, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	addParametersToContext(c, networkID, fooSerdeType)

	// Success
	url = getURL(restPort, networkID, fooSerdeType)
	err = configuratorh.GetDeleteNetworkConfigHandler(url).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	assertConfigDoesNotExist(t, networkID, fooSerdeType)
}

func addParametersToContext(c echo.Context, networkID string, configType string) echo.Context {
	c.SetParamNames("network_id", "config_type")
	c.SetParamValues(networkID, fooSerdeType)
	return c
}

func assertConfigExists(t *testing.T, networkID string, configType string, config interface{}) {
	networks, notFound, err := configurator.LoadNetworks([]string{networkID}, true, true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(notFound))
	retrievedConfig, err := serde.Deserialize(configurator.SerdeDomain, configType, networks[networkID].Configs[configType])
	assert.NoError(t, err)
	assert.Equal(t, config, retrievedConfig)
}

func assertConfigDoesNotExist(t *testing.T, networkID string, configType string) {
	networks, notFound, err := configurator.LoadNetworks([]string{networkID}, true, true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(notFound))
	_, ok := networks[networkID].Configs[configType]
	assert.Equal(t, false, ok)
}

func getURL(restPort int, networkID string, configType string) string {
	url := fmt.Sprintf(
		"http://localhost:%d%s/networks/%s/configs/%s",
		restPort,
		handlers.REST_ROOT,
		networkID,
		configType,
	)
	return url
}

type FooConfigs struct {
	ConfigStr string `json:"config_str"`
	ConfigNum int    `json:"config_num"`
}

type fooSerde struct{}

func (*fooSerde) GetType() string {
	return fooSerdeType
}

func (*fooSerde) GetDomain() string {
	return configurator.SerdeDomain
}

func (*fooSerde) Serialize(c interface{}) ([]byte, error) {
	return json.Marshal(c)
}

func (*fooSerde) Deserialize(message []byte) (interface{}, error) {
	res := FooConfigs{}
	err := json.Unmarshal(message, &res)
	return res, err
}
