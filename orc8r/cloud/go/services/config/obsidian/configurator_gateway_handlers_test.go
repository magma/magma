/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package obsidian_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/config/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func commonSetupGateways(t *testing.T) {
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
	testGetEntityConfig(t, "cfg_gateway")
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
	err = handler.HandlerFunc(c)
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
	err = handler.HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.Equal(t, "Invalid config: hello", err.(*echo.HTTPError).Message)

	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
}

func TestConfiguratorUpdateGatewayConfig(t *testing.T) {
	commonSetupGateways(t)
	testEntityUpdate(t, "cfg_gateway", "err_gateway")
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
	err = handler.HandlerFunc(c)
	assert.NoError(t, err)

	actual, err := configurator.LoadEntity("network1", "cfg_gateway", "key", configurator.EntityLoadCriteria{LoadConfig: true})
	assert.EqualError(t, err, errors.ErrNotFound.Error())
	assert.Equal(t, configurator.NetworkEntity{}, actual)

	// Double delete - should be no error
	err = handler.HandlerFunc(c)
	assert.NoError(t, err)

	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
}
