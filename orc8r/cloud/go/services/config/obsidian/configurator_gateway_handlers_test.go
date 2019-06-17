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
	"os"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/config/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func commonSetupGateways(t *testing.T) {
	_ = os.Setenv(handlers.UseNewHandlersEnv, "1")
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(
		configurator.NewNetworkEntityConfigSerde("cfg_gateway", &configType{}),
		configurator.NewNetworkEntityConfigSerde("err_gateway", &errValidateType{}),
	)
	assert.NoError(t, err)
	test_init.StartTestService(t)
	err = configurator.CreateNetwork(configurator.Network{ID: "network1"})
	assert.NoError(t, err)

}

func TestConfiguratorGetGatewayConfig(t *testing.T) {
	commonSetupGateways(t)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	handler := obsidian.GetReadConfigHandler("google.com", "cfg_gateway", mockKeyGetter, &configType{})

	// 404
	err := handler.MigratedHandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)

	// Happy path
	expected := &configType{Foo: "foo", Bar: "bar"}
	_, err = configurator.CreateEntity("network1", configurator.NetworkEntity{
		Type:   "cfg_gateway",
		Key:    "key",
		Config: expected,
	})
	assert.NoError(t, err)
	err = handler.MigratedHandlerFunc(c)
	assert.NoError(t, err)

	actual := &configType{}
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
}

func TestConfiguratorCreateGatewayConfig(t *testing.T) {
	commonSetupGateways(t)

	e := echo.New()

	// Happy path
	post := `{"Foo": "foo", "Bar": "bar"}`
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	_, err := configurator.CreateEntity("network1", configurator.NetworkEntity{
		Type: "magmad_gateway",
		Key:  "key",
	})
	assert.NoError(t, err)

	handler := obsidian.GetCreateConfigHandler("google.com", "cfg_gateway", mockKeyGetter, &configType{})
	err = handler.MigratedHandlerFunc(c)
	assert.NoError(t, err)
	actual, err := configurator.LoadEntity("network1", "cfg_gateway", "key", configurator.EntityLoadCriteria{LoadConfig: true})
	assert.NoError(t, err)
	assert.Equal(t, &configType{Foo: "foo", Bar: "bar"}, actual.Config)

	// Validation error
	post = `{"Msg": "hello"}`
	req = httptest.NewRequest(echo.POST, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	handler = obsidian.GetCreateConfigHandler("google.com", "err_gateway", mockKeyGetter, &errValidateType{})
	err = handler.MigratedHandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.Equal(t, "hello", err.(*echo.HTTPError).Message)

	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
}

func TestConfiguratorUpdateGatewayConfig(t *testing.T) {
	commonSetupGateways(t)

	e := echo.New()

	// Happy path - create a config with the PUT
	post := `{"Foo": "foo", "Bar": "bar"}`
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	_, err := configurator.CreateEntity("network1", configurator.NetworkEntity{
		Type:   "cfg_gateway",
		Key:    "key",
		Config: &configType{Foo: "foo", Bar: "bar"},
	})
	assert.NoError(t, err)

	// Happy path - update a config with the PUT
	handler := obsidian.GetUpdateConfigHandler("google.com", "cfg_gateway", mockKeyGetter, &configType{})
	post = `{"Foo": "foo2", "Bar": "bar2"}`
	req = httptest.NewRequest(echo.PUT, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = handler.MigratedHandlerFunc(c)
	assert.NoError(t, err)
	actual, err := configurator.LoadEntity("network1", "cfg_gateway", "key", configurator.EntityLoadCriteria{LoadConfig: true})
	assert.NoError(t, err)
	assert.Equal(t, &configType{Foo: "foo2", Bar: "bar2"}, actual.Config)

	// Validation error
	post = `{"Msg": "hello"}`
	req = httptest.NewRequest(echo.POST, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	handler = obsidian.GetUpdateConfigHandler("google.com", "err_gateway", mockKeyGetter, &errValidateType{})
	err = handler.MigratedHandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.Equal(t, "hello", err.(*echo.HTTPError).Message)

	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
}

func TestConfiguratorDeleteGatewayConfig(t *testing.T) {
	commonSetupGateways(t)
	_, err := configurator.CreateEntity("network1", configurator.NetworkEntity{
		Type:   "cfg_gateway",
		Key:    "key",
		Config: &configType{Foo: "foo"},
	})
	assert.NoError(t, err)

	e := echo.New()

	// Happy path
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	handler := obsidian.GetDeleteConfigHandler("google.com", "cfg_gateway", mockKeyGetter)
	err = handler.MigratedHandlerFunc(c)
	assert.NoError(t, err)

	actual, err := configurator.LoadEntity("network1", "cfg_gateway", "key", configurator.EntityLoadCriteria{LoadConfig: true})
	assert.EqualError(t, err, errors.ErrNotFound.Error())
	assert.Equal(t, configurator.NetworkEntity{}, actual)

	// Double delete - should be no error
	err = handler.MigratedHandlerFunc(c)
	assert.NoError(t, err)

	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
}
