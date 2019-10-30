/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package obsidian_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/config/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func commonSetupNetworks(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, configurator.NetworkConfigSerdeDomain)
	err := serde.RegisterSerdes(
		configurator.NewNetworkConfigSerde("cfg_network", &configType{}),
		configurator.NewNetworkConfigSerde("err_network", &errValidateType{}),
	)
	assert.NoError(t, err)
	test_init.StartTestService(t)
	err = configurator.CreateNetwork(configurator.Network{ID: "network1"})
	assert.NoError(t, err)
}

func TestConfiguratorGetNetworkConfig(t *testing.T) {
	commonSetupNetworks(t)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	handler := obsidian.GetReadConfigHandler("google.com", "cfg_network", mockKeyGetter, &configType{})

	// 404
	err := handler.HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)

	// Happy path
	expected := &configType{Foo: "foo", Bar: "bar"}
	err = configurator.UpdateNetworkConfig("network1", "cfg_network", expected)
	assert.NoError(t, err)
	err = handler.HandlerFunc(c)
	assert.NoError(t, err)

	actual := &configType{}
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	serde.UnregisterSerdesForDomain(t, configurator.NetworkConfigSerdeDomain)
}

func TestConfiguratorCreateNetworkConfig(t *testing.T) {
	commonSetupNetworks(t)

	e := echo.New()

	// Happy path
	post := `{"Foo": "foo", "Bar": "bar"}`
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	handler := obsidian.GetCreateConfigHandler("google.com", "cfg_network", mockKeyGetter, &configType{})
	err := handler.HandlerFunc(c)
	assert.NoError(t, err)
	actual, err := configurator.GetNetworkConfigsByType("network1", "cfg_network")
	assert.NoError(t, err)
	assert.Equal(t, &configType{Foo: "foo", Bar: "bar"}, actual)

	// Validation error
	post = `{"Msg": "hello"}`
	req = httptest.NewRequest(echo.POST, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	handler = obsidian.GetCreateConfigHandler("google.com", "err_network", mockKeyGetter, &errValidateType{})
	err = handler.HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.Equal(t, "Invalid config: hello", err.(*echo.HTTPError).Message)

	serde.UnregisterSerdesForDomain(t, configurator.NetworkConfigSerdeDomain)
}

func TestConfiguratorUpdateNetworkConfig(t *testing.T) {
	commonSetupNetworks(t)

	e := echo.New()

	// Happy path - create a config with the PUT
	post := `{"Foo": "foo", "Bar": "bar"}`
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	handler := obsidian.GetUpdateConfigHandler("google.com", "cfg_network", mockKeyGetter, &configType{})
	err := handler.HandlerFunc(c)
	assert.NoError(t, err)
	actual, err := configurator.GetNetworkConfigsByType("network1", "cfg_network")
	assert.NoError(t, err)
	assert.Equal(t, &configType{Foo: "foo", Bar: "bar"}, actual)

	// Happy path - update a config with the PUT
	post = `{"Foo": "foo2", "Bar": "bar2"}`
	req = httptest.NewRequest(echo.PUT, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = handler.HandlerFunc(c)
	assert.NoError(t, err)
	actual, err = configurator.GetNetworkConfigsByType("network1", "cfg_network")
	assert.NoError(t, err)
	assert.Equal(t, &configType{Foo: "foo2", Bar: "bar2"}, actual)

	// Validation error
	post = `{"Msg": "hello"}`
	req = httptest.NewRequest(echo.POST, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	handler = obsidian.GetUpdateConfigHandler("google.com", "err_network", mockKeyGetter, &errValidateType{})
	err = handler.HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.Equal(t, "Invalid config: hello", err.(*echo.HTTPError).Message)

	serde.UnregisterSerdesForDomain(t, configurator.NetworkConfigSerdeDomain)
}

func TestConfiguratorDeleteNetworkConfig(t *testing.T) {
	commonSetupNetworks(t)
	err := configurator.UpdateNetworkConfig("network1", "cfg_network", &configType{Foo: "foo", Bar: "bar"})
	assert.NoError(t, err)

	e := echo.New()

	// Happy path
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	handler := obsidian.GetDeleteConfigHandler("google.com", "cfg_network", mockKeyGetter)
	err = handler.HandlerFunc(c)
	assert.NoError(t, err)

	actual, err := configurator.GetNetworkConfigsByType("network1", "cfg_network")
	assert.NoError(t, err)
	assert.Nil(t, actual)

	// Double delete - should be no error
	err = handler.HandlerFunc(c)
	assert.NoError(t, err)

	serde.UnregisterSerdesForDomain(t, configurator.NetworkConfigSerdeDomain)
}
