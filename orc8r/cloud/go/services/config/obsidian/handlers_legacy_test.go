/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package obsidian_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/config"
	cfgObsidian "magma/orc8r/cloud/go/services/config/obsidian"
	config_test_init "magma/orc8r/cloud/go/services/config/test_init"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/magmad/obsidian/models"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestReadAllKeysConfigHandlerLegacy(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{config.SerdeDomain}, &convertErrConfigManager{config.SerdeDomain}, &errConfigManager{config.SerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models.MagmadGatewayConfig{}))
	assert.NoError(t, err)
	obsidian.TLS = false // To bypass access control

	config_test_init.StartTestService(t)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	// 404
	actual := &fooConfig{}
	err = cfgObsidian.GetReadConfigHandler("google.com", "foo", mockKeyGetter, actual).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)

	// Happy path
	a_config := &fooConfig{Foo: "foo", Bar: "bar"}
	err = config.CreateConfig("network1", "foo", "key", a_config)
	assert.NoError(t, err)

	actual_keys := &[]string{}
	expected := &[]string{"key"}
	err = cfgObsidian.GetReadAllKeysConfigHandler("google.com", "foo").HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), actual_keys)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual_keys)

	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
}

func TestGetConfigHandlerLegacy(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{config.SerdeDomain}, &convertErrConfigManager{config.SerdeDomain}, &errConfigManager{config.SerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models.MagmadGatewayConfig{}))
	assert.NoError(t, err)
	obsidian.TLS = false // To bypass access control

	config_test_init.StartTestService(t)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	// 404
	actual := &fooConfig{}
	err = cfgObsidian.GetReadConfigHandler("google.com", "foo", mockKeyGetter, actual).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)

	// Happy path
	expected := &fooConfig{Foo: "foo", Bar: "bar"}
	err = config.CreateConfig("network1", "foo", "key", expected)
	assert.NoError(t, err)

	err = cfgObsidian.GetReadConfigHandler("google.com", "foo", mockKeyGetter, actual).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// Convert error
	expectedConvertErrCfg := &convertErrConfig{Val: 1}
	err = config.CreateConfig("network1", "convertErr", "key", expectedConvertErrCfg)
	assert.NoError(t, err)

	actualConvertErr := &convertErrConfig{}
	err = cfgObsidian.GetReadConfigHandler("google.com", "convertErr", mockKeyGetter, actualConvertErr).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)

	// Config service error
	expectedUnmarshalErrCfg := &errConfig{ShouldErrorOnMarshal: "N", ShouldErrorOnUnmarshal: "Y"}
	err = config.CreateConfig("network1", "err", "key", expectedUnmarshalErrCfg)
	assert.NoError(t, err)

	actualUnmarshalErr := &errConfig{}
	err = cfgObsidian.GetReadConfigHandler("google.com", "err", mockKeyGetter, actualUnmarshalErr).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)

	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
}

func TestCreateConfigHandlerLegacy(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{config.SerdeDomain}, &convertErrConfigManager{config.SerdeDomain}, &errConfigManager{config.SerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models.MagmadGatewayConfig{}))
	assert.NoError(t, err)
	obsidian.TLS = false // To bypass access control

	configurator_test_init.StartTestService(t)
	config_test_init.StartTestService(t)

	e := echo.New()

	// Happy path
	post := `{"Foo": "foo", "Bar": "bar"}`
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = cfgObsidian.GetCreateConfigHandler("google.com", "foo", mockKeyGetter, &fooConfig{}).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, `"key"`, rec.Body.String())
	actual, err := config.GetConfig("network1", "foo", "key")
	assert.NoError(t, err)
	assert.Equal(t, &fooConfig{Foo: "foo", Bar: "bar"}, actual)

	glog.Errorf("IGNORE REST")
	// Validate (convert) error
	post = `{"Val": 1}`
	req = httptest.NewRequest(echo.PUT, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = cfgObsidian.GetCreateConfigHandler("google.com", "convertErr", mockKeyGetter, &convertErrConfig{}).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.Contains(t, err.Error(), "Validate error")

	// Config service error (creating duplicate config)
	post = `{"Foo": "bar", "Bar": "foo"}`
	req = httptest.NewRequest(echo.PUT, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = cfgObsidian.GetCreateConfigHandler("google.com", "foo", mockKeyGetter, &fooConfig{}).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.Contains(t, err.Error(), "Creating already existing config")

	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
}

func TestUpdateConfigHandlerLegacy(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{config.SerdeDomain}, &convertErrConfigManager{config.SerdeDomain}, &errConfigManager{config.SerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models.MagmadGatewayConfig{}))
	assert.NoError(t, err)
	obsidian.TLS = false // To bypass access control

	config_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	err = config.CreateConfig("network1", "foo", "key", &fooConfig{Foo: "foo", Bar: "bar"})
	assert.NoError(t, err)

	e := echo.New()

	// Happy path
	post := `{"Foo": "bar", "Bar": "foo"}`
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = cfgObsidian.GetUpdateConfigHandler("google.com", "foo", mockKeyGetter, &fooConfig{}).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	actualFoo, err := config.GetConfig("network1", "foo", "key")
	assert.Equal(t, &fooConfig{Foo: "bar", Bar: "foo"}, actualFoo)

	// Validate (convert) error
	post = `{"Value": 1}`
	req = httptest.NewRequest(echo.POST, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = cfgObsidian.GetUpdateConfigHandler("google.com", "convertErr", mockKeyGetter, &convertErrConfig{}).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.Contains(t, err.Error(), "Validate error")

	// Config service error (updating nonexistent config)
	post = `{"Foo": "baz"}`
	req = httptest.NewRequest(echo.POST, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = cfgObsidian.GetUpdateConfigHandler("google.com", "foo", func(ctx echo.Context) (string, *echo.HTTPError) { return "dne", nil }, &fooConfig{}).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.Contains(t, err.Error(), "Updating nonexistent config")

	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
}

func TestDeleteConfigHandlerLegacy(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{config.SerdeDomain}, &convertErrConfigManager{config.SerdeDomain}, &errConfigManager{config.SerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models.MagmadGatewayConfig{}))
	assert.NoError(t, err)
	obsidian.TLS = false // To bypass access control

	config_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	err = config.CreateConfig("network1", "foo", "key", &fooConfig{Foo: "foo", Bar: "bar"})
	assert.NoError(t, err)

	e := echo.New()

	// Happy path
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = cfgObsidian.GetDeleteConfigHandler("google.com", "foo", mockKeyGetter).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Config service error - deleting nonexistent config
	err = cfgObsidian.GetDeleteConfigHandler("google.com", "foo", mockKeyGetter).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.Contains(t, err.Error(), "Deleting nonexistent config")

	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
}
