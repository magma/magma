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

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	models2 "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/serde"
	cfgObsidian "magma/orc8r/cloud/go/services/config/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"

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
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models2.MagmadGatewayConfigs{}))
	assert.NoError(t, err)
	obsidian.TLS = false // To bypass access control

	configurator_test_init.StartTestService(t)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = configurator.CreateNetwork(configurator.Network{ID: "network1"})
	assert.NoError(t, err)

	// 404
	actual := &fooConfig{}
	err = cfgObsidian.GetReadConfigHandler("google.com", "foo", mockKeyGetter, actual).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)

	// Happy path
	a_config := &fooConfig{Foo: "foo", Bar: "bar"}
	_, err = configurator.CreateEntity("network1", configurator.NetworkEntity{Key: "key", Type: "foo", Config: a_config})
	assert.NoError(t, err)

	actual_keys := &[]string{}
	expected := &[]string{"key"}
	err = cfgObsidian.GetReadAllKeysConfigHandler("google.com", "foo").HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), actual_keys)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual_keys)
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
}

func TestGetConfigHandler(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models2.MagmadGatewayConfigs{}))
	assert.NoError(t, err)
	obsidian.TLS = false // To bypass access control

	configurator_test_init.StartTestService(t)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = configurator.CreateNetwork(configurator.Network{ID: "network1"})
	assert.NoError(t, err)

	// 404
	actual := &fooConfig{}
	err = cfgObsidian.GetReadConfigHandler("google.com", "foo", mockKeyGetter, actual).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)

	// Happy path
	expected := &fooConfig{Foo: "foo", Bar: "bar"}
	_, err = configurator.CreateEntity("network1", configurator.NetworkEntity{Key: "key", Type: "foo", Config: expected})
	assert.NoError(t, err)

	err = cfgObsidian.GetReadConfigHandler("google.com", "foo", mockKeyGetter, actual).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// Config service error
	expectedUnmarshalErrCfg := &errConfig{ShouldErrorOnMarshal: "N", ShouldErrorOnUnmarshal: "Y"}
	_, err = configurator.CreateEntity("network1", configurator.NetworkEntity{Key: "key", Type: "err", Config: expectedUnmarshalErrCfg})
	assert.Error(t, err)

	actualUnmarshalErr := &errConfig{}
	err = cfgObsidian.GetReadConfigHandler("google.com", "err", mockKeyGetter, actualUnmarshalErr).HandlerFunc(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)

	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
}

func TestCreateConfigHandler(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models2.MagmadGatewayConfigs{}))
	assert.NoError(t, err)
	obsidian.TLS = false // To bypass access control

	configurator_test_init.StartTestService(t)

	e := echo.New()

	// Happy path
	post := `{"Foo": "foo", "Bar": "bar"}`
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("network_id")
	c.SetParamValues("network1")

	err = configurator.CreateNetwork(configurator.Network{ID: "network1"})
	assert.NoError(t, err)

	err = cfgObsidian.GetCreateConfigHandler("google.com", "foo", mockKeyGetter, &fooConfig{}).HandlerFunc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, `"key"`, rec.Body.String())
	actual, err := configurator.LoadEntityConfig("network1", "foo", "key")
	assert.NoError(t, err)
	assert.Equal(t, &fooConfig{Foo: "foo", Bar: "bar"}, actual)

	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
}

func TestUpdateConfigHandler(t *testing.T) {
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(&fooConfigManager{configurator.NetworkEntitySerdeDomain}, &convertErrConfigManager{configurator.NetworkEntitySerdeDomain}, &errConfigManager{configurator.NetworkEntitySerdeDomain})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models2.MagmadGatewayConfigs{}), configurator.NewNetworkConfigSerde("foo_network", &fooConfig{}))
	assert.NoError(t, err)
	obsidian.TLS = false // To bypass access control

	configurator_test_init.StartTestService(t)
	err = configurator.CreateNetwork(configurator.Network{ID: "network1", Configs: map[string]interface{}{"foo_network": &fooConfig{Foo: "foo", Bar: "bar"}}})
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

	err = cfgObsidian.GetUpdateConfigHandler("google.com", "foo_network", mockKeyGetter, &fooConfig{}).HandlerFunc(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	network, err := configurator.LoadNetwork("network1", false, true)
	assert.Equal(t, &fooConfig{Foo: "bar", Bar: "foo"}, network.Configs["foo_network"])

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

	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
}

func TestDeleteConfigHandler(t *testing.T) {
	t.Skip()
	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
	err := serde.RegisterSerdes(
		&fooConfigManager{configurator.NetworkEntitySerdeDomain},
		&convertErrConfigManager{configurator.NetworkEntitySerdeDomain},
		&errConfigManager{configurator.NetworkEntitySerdeDomain},
	)
	assert.NoError(t, err)
	err = serde.RegisterSerdes(configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models2.MagmadGatewayConfigs{}))
	assert.NoError(t, err)
	obsidian.TLS = false // To bypass access control

	configurator_test_init.StartTestService(t)

	err = configurator.CreateNetwork(configurator.Network{ID: "network1"})
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("network1", configurator.NetworkEntity{Type: "foo", Key: "key", Config: &fooConfig{Foo: "foo", Bar: "bar"}})
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

	serde.UnregisterSerdesForDomain(t, configurator.NetworkEntitySerdeDomain)
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

func (foo *fooConfig) MarshalBinary() ([]byte, error) {
	return json.Marshal(foo)
}

func (foo *fooConfig) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, foo)
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

func (*convertErrConfig) MarshalBinary() ([]byte, error) {
	return nil, errors.New("MarshalBinary error")
}

func (*convertErrConfig) UnmarshalBinary(data []byte) error {
	return errors.New("UnmarshalBinary error")
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

func (c *errConfig) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *errConfig) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
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
