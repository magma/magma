/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package obsidian_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	obisidan_config "magma/orc8r/cloud/go/obsidian/config"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/config/obsidian"
	config_test_init "magma/orc8r/cloud/go/services/config/test_init"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	config2 "magma/orc8r/cloud/go/services/magmad/config"
	"magma/orc8r/cloud/go/services/magmad/obsidian/models"
	storage2 "magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// For the test, fooConfig will be both the user model and the service model
type (
	// For happy path
	fooConfig struct {
		Foo, Bar string
	}
	fooConfigManager struct {
		domain string
	}

	// To coerce errors in config conversion
	convertErrConfig struct {
		Val int
	}
	convertErrConfigManager struct {
		domain string
	}

	// To coerce errors in config service serialization/deserialization
	errConfig struct {
		ShouldErrorOnMarshal, ShouldErrorOnUnmarshal string // Y | N
	}
	errConfigManager struct {
		domain string
	}
)

func mockKeyGetter(_ echo.Context) (string, *echo.HTTPError) {
	return "key", nil
}

func TestReadAllKeysConfigHandler(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{config.SerdeDomain}, &convertErrConfigManager{config.SerdeDomain}, &errConfigManager{config.SerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(config2.MagmadGatewayType, &models.MagmadGatewayConfig{}))
	assert.NoError(t, err)
	obisidan_config.TLS = false // To bypass access control

	config_test_init.StartTestService(t)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	// 404
	actual := &fooConfig{}
	err = obsidian.GetReadConfigHandler("google.com", "foo", mockKeyGetter, actual).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)

	// Happy path
	a_config := &fooConfig{Foo: "foo", Bar: "bar"}
	err = config.CreateConfig("network1", "foo", "key", a_config)
	assert.NoError(t, err)

	actual_keys := &[]string{}
	expected := &[]string{"key"}
	err = obsidian.GetReadAllKeysConfigHandler("google.com", "foo").HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), actual_keys)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual_keys)

	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
}

func TestGetConfigHandler(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{config.SerdeDomain}, &convertErrConfigManager{config.SerdeDomain}, &errConfigManager{config.SerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(config2.MagmadGatewayType, &models.MagmadGatewayConfig{}))
	assert.NoError(t, err)
	obisidan_config.TLS = false // To bypass access control

	config_test_init.StartTestService(t)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	// 404
	actual := &fooConfig{}
	err = obsidian.GetReadConfigHandler("google.com", "foo", mockKeyGetter, actual).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)

	// Happy path
	expected := &fooConfig{Foo: "foo", Bar: "bar"}
	err = config.CreateConfig("network1", "foo", "key", expected)
	assert.NoError(t, err)

	err = obsidian.GetReadConfigHandler("google.com", "foo", mockKeyGetter, actual).HandlerFunc(c)
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
	err = obsidian.GetReadConfigHandler("google.com", "convertErr", mockKeyGetter, actualConvertErr).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)

	// Config service error
	expectedUnmarshalErrCfg := &errConfig{ShouldErrorOnMarshal: "N", ShouldErrorOnUnmarshal: "Y"}
	err = config.CreateConfig("network1", "err", "key", expectedUnmarshalErrCfg)
	assert.NoError(t, err)

	actualUnmarshalErr := &errConfig{}
	err = obsidian.GetReadConfigHandler("google.com", "err", mockKeyGetter, actualUnmarshalErr).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)

	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
}

func TestCreateConfigHandler(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{config.SerdeDomain}, &convertErrConfigManager{config.SerdeDomain}, &errConfigManager{config.SerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(config2.MagmadGatewayType, &models.MagmadGatewayConfig{}))
	assert.NoError(t, err)
	obisidan_config.TLS = false // To bypass access control

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

	err = obsidian.GetCreateConfigHandler("google.com", "foo", mockKeyGetter, &fooConfig{}).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, `"key"`, rec.Body.String())
	actual, err := config.GetConfig("network1", "foo", "key")
	assert.NoError(t, err)
	assert.Equal(t, &fooConfig{Foo: "foo", Bar: "bar"}, actual)

	// Test changes are properly reflected in configurator
	testEntityConfigsInConfigurator(t, "network1", "foo", "key", &fooConfig{Foo: "foo", Bar: "bar"})

	glog.Errorf("IGNORE REST")
	// Validate (convert) error
	post = `{"Val": 1}`
	req = httptest.NewRequest(echo.PUT, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = obsidian.GetCreateConfigHandler("google.com", "convertErr", mockKeyGetter, &convertErrConfig{}).HandlerFunc(c)
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

	err = obsidian.GetCreateConfigHandler("google.com", "foo", mockKeyGetter, &fooConfig{}).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.Contains(t, err.Error(), "Creating already existing config")

	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
}

func TestUpdateConfigHandler(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{config.SerdeDomain}, &convertErrConfigManager{config.SerdeDomain}, &errConfigManager{config.SerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(config2.MagmadGatewayType, &models.MagmadGatewayConfig{}))
	assert.NoError(t, err)
	obisidan_config.TLS = false // To bypass access control

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

	err = obsidian.GetUpdateConfigHandler("google.com", "foo", mockKeyGetter, &fooConfig{}).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	actualFoo, err := config.GetConfig("network1", "foo", "key")
	assert.Equal(t, &fooConfig{Foo: "bar", Bar: "foo"}, actualFoo)

	// Test changes are properly reflected in configurator
	testEntityConfigsInConfigurator(t, "network1", "foo", "key", &fooConfig{Foo: "bar", Bar: "foo"})

	// Validate (convert) error
	post = `{"Value": 1}`
	req = httptest.NewRequest(echo.POST, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = obsidian.GetUpdateConfigHandler("google.com", "convertErr", mockKeyGetter, &convertErrConfig{}).HandlerFunc(c)
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

	err = obsidian.GetUpdateConfigHandler("google.com", "foo", func(ctx echo.Context) (string, *echo.HTTPError) { return "dne", nil }, &fooConfig{}).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.Contains(t, err.Error(), "Updating nonexistent config")

	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
}

func TestDeleteConfigHandler(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{config.SerdeDomain}, &convertErrConfigManager{config.SerdeDomain}, &errConfigManager{config.SerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(config2.MagmadGatewayType, &models.MagmadGatewayConfig{}))
	assert.NoError(t, err)
	obisidan_config.TLS = false // To bypass access control

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

	err = obsidian.GetDeleteConfigHandler("google.com", "foo", mockKeyGetter).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	/// Test changes are properly reflected in configurator
	testEntityConfigsNotInConfigurator(t, "network1", "foo", "key")

	// Config service error - deleting nonexistent config
	err = obsidian.GetDeleteConfigHandler("google.com", "foo", mockKeyGetter).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.Contains(t, err.Error(), "Deleting nonexistent config")

	serde.UnregisterSerdesForDomain(t, config.SerdeDomain)
}

func testEntityConfigsInConfigurator(t *testing.T, networkID, entityType, entityID string, expectedConfig interface{}) {
	entities, entitiesNotFound, err := configurator.LoadEntities(
		networkID,
		nil,
		nil,
		[]storage2.TypeAndKey{{Type: entityType, Key: entityID}},
		configurator.EntityLoadCriteria{LoadConfig: true},
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(entities))
	assert.Equal(t, 0, len(entitiesNotFound))
	assert.Equal(t, expectedConfig, entities[0].Config)
}

func testEntityConfigsNotInConfigurator(t *testing.T, networkID, entityType, entityID string) {
	entityTK := storage2.TypeAndKey{Type: entityType, Key: entityID}
	entities, entitiesNotFound, err := configurator.LoadEntities(
		networkID,
		nil,
		nil,
		[]storage2.TypeAndKey{entityTK},
		configurator.EntityLoadCriteria{LoadConfig: true},
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(entities))
	assert.Equal(t, 0, len(entitiesNotFound))
	assert.Nil(t, entities[0].Config)
}

// Interface implementations for test configs
func (*fooConfig) ValidateModel() error {
	return nil
}

func (foo *fooConfig) ToServiceModel() (interface{}, error) {
	return foo, nil
}

func (foo *fooConfig) FromServiceModel(serviceModel interface{}) error {
	casted := serviceModel.(*fooConfig)
	foo.Foo = casted.Foo
	foo.Bar = casted.Bar
	return nil
}

func (f *fooConfigManager) GetDomain() string {
	return f.domain
}

func (*fooConfigManager) GetType() string {
	return "foo"
}

func (*fooConfigManager) Serialize(config interface{}) ([]byte, error) {
	cfgCasted := config.(*fooConfig)
	return []byte(fmt.Sprintf("%s|%s", cfgCasted.Foo, cfgCasted.Bar)), nil
}

func (*fooConfigManager) Deserialize(message []byte) (interface{}, error) {
	foobar := string(message)
	foobarSplit := strings.Split(foobar, "|")
	if len(foobarSplit) != 2 {
		return nil, fmt.Errorf("Expected 2 fields, got %d", len(foobarSplit))
	}
	return &fooConfig{Foo: foobarSplit[0], Bar: foobarSplit[1]}, nil
}

func (*convertErrConfig) ValidateModel() error {
	return errors.New("Validate error")
}

func (*convertErrConfig) ToServiceModel() (interface{}, error) {
	return nil, errors.New("ToServiceModel error")
}

func (*convertErrConfig) FromServiceModel(serviceModel interface{}) error {
	return errors.New("FromSerivceModel error")
}

func (c *convertErrConfigManager) GetDomain() string {
	return c.domain
}

func (*convertErrConfigManager) GetType() string {
	return "convertErr"
}

func (*convertErrConfigManager) Serialize(config interface{}) ([]byte, error) {
	return []byte("convertErr"), nil
}

func (*convertErrConfigManager) Deserialize(message []byte) (interface{}, error) {
	return &convertErrConfig{}, nil
}

func (*errConfig) ValidateModel() error {
	return nil
}

func (c *errConfig) ToServiceModel() (interface{}, error) {
	return c, nil
}

func (c *errConfig) FromServiceModel(serviceModel interface{}) error {
	castedModel := serviceModel.(*errConfig)
	c.ShouldErrorOnMarshal = castedModel.ShouldErrorOnMarshal
	c.ShouldErrorOnUnmarshal = castedModel.ShouldErrorOnUnmarshal
	return nil
}

func (*errConfigManager) GetType() string {
	return "err"
}

func (e *errConfigManager) GetDomain() string {
	return e.domain
}

func (*errConfigManager) Serialize(config interface{}) ([]byte, error) {
	castedConfig := config.(*errConfig)
	if castedConfig.ShouldErrorOnMarshal == "Y" {
		return nil, errors.New("Serialize error")
	}
	return []byte(fmt.Sprintf("%s|%s", castedConfig.ShouldErrorOnMarshal, castedConfig.ShouldErrorOnUnmarshal)), nil
}

func (*errConfigManager) Deserialize(message []byte) (interface{}, error) {
	msgString := string(message)
	msgSplit := strings.Split(msgString, "|")
	if msgSplit[1] == "Y" {
		return nil, errors.New("Deserialize error")
	}
	return &errConfig{ShouldErrorOnMarshal: msgSplit[0], ShouldErrorOnUnmarshal: msgSplit[1]}, nil
}
