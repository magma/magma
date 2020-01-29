/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/client/mocks"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/config"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/receivers"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/prometheus/alert"

	"github.com/labstack/echo"
	amconfig "github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
)

const (
	testNID = "test"
)

var (
	testWebhookURL, _ = url.Parse("http://test.com")
	sampleReceiver    = receivers.Receiver{
		Name: "testSlackReceiver",
		SlackConfigs: []*receivers.SlackConfig{{
			APIURL:   "http://slack.com/12345",
			Channel:  "test_channel",
			Username: "test_username",
		}},
		WebhookConfigs: []*receivers.WebhookConfig{{
			NotifierConfig: amconfig.NotifierConfig{
				VSendResolved: true,
			},
			URL: &amconfig.URL{
				URL: testWebhookURL,
			},
		}},
	}

	fiveSeconds, _ = model.ParseDuration("5s")

	sampleRoute = amconfig.Route{
		Receiver: "testSlackReceiver",
		Match:    map[string]string{"networkID": testNID},
		Routes: []*amconfig.Route{{
			Receiver: "childReceiver",
			Match:    map[string]string{"severity": "critical"},
		}},
		GroupWait:      &fiveSeconds,
		GroupInterval:  &fiveSeconds,
		RepeatInterval: &fiveSeconds,
	}

	sampleJSONRoute = receivers.RouteJSONWrapper{
		Receiver:  "testSlackReceiver",
		GroupWait: "5s",
	}
	convertedSampleJSONRoute = amconfig.Route{
		Receiver:  "testSlackReceiver",
		GroupWait: &fiveSeconds,
	}
)

func TestGetReceiverPostHandler(t *testing.T) {
	// Successful Post
	client := &mocks.AlertmanagerClient{}
	client.On("CreateReceiver", testNID, sampleReceiver).Return(nil)
	client.On("ReloadAlertmanager").Return(nil)
	c, rec := buildContext(sampleReceiver, http.MethodPost, "/", v1receiverPath, testNID)

	err := GetReceiverPostHandler(client)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertExpectations(t)

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("CreateReceiver", testNID, receivers.Receiver{}).Return(errors.New("error"))
	c, _ = buildContext(nil, http.MethodPost, "/", v1receiverPath, testNID)

	err = GetReceiverPostHandler(client)(c)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=400, message=error`)
	client.AssertExpectations(t)

	// Alertmanager Error
	client = &mocks.AlertmanagerClient{}
	client.On("ReloadAlertmanager").Return(errors.New("error"))
	client.On("CreateReceiver", testNID, receivers.Receiver{}).Return(nil)
	c, _ = buildContext(nil, http.MethodPut, "/", v1receiverPath, testNID)

	err = GetReceiverPostHandler(client)(c)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=500, message=error`)
	client.AssertExpectations(t)
}

func TestGetGetReceiversHandler(t *testing.T) {
	// Successful Get
	client := &mocks.AlertmanagerClient{}
	client.On("GetReceivers", testNID).Return([]receivers.Receiver{sampleReceiver}, nil)
	c, rec := buildContext(nil, http.MethodGet, "/", v1receiverPath, testNID)

	err := GetGetReceiversHandler(client)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertExpectations(t)

	var receiver []receivers.Receiver
	err = json.Unmarshal(rec.Body.Bytes(), &receiver)
	assert.Equal(t, 1, len(receiver))
	assert.Equal(t, sampleReceiver, receiver[0])

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("GetReceivers", testNID).Return([]receivers.Receiver{}, errors.New("error"))
	c, _ = buildContext(nil, http.MethodGet, "/", v1receiverPath, testNID)

	err = GetGetReceiversHandler(client)(c)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=500, message=error`)
	client.AssertExpectations(t)

	//Get specific alert
	client = &mocks.AlertmanagerClient{}
	client.On("GetReceivers", testNID).Return([]receivers.Receiver{sampleReceiver}, nil)
	c, rec = buildContext(nil, http.MethodGet, "/", v1receiverNamePath, testNID)
	c.SetParamNames(receiverNameParam)
	c.SetParamValues("testSlackReceiver")

	err = GetGetReceiversHandler(client)(c)
	assert.NoError(t, err)
	var r receivers.Receiver
	err = json.Unmarshal(rec.Body.Bytes(), &r)
	assert.NoError(t, err)
	assert.Equal(t, sampleReceiver, r)

	//Get nonexistent alert
	client = &mocks.AlertmanagerClient{}
	client.On("GetReceivers", testNID).Return([]receivers.Receiver{sampleReceiver}, nil)
	c, rec = buildContext(nil, http.MethodGet, "/", v1receiverNamePath, testNID)
	c.SetParamNames(receiverNameParam)
	c.SetParamValues("testNewReceiver")

	err = GetGetReceiversHandler(client)(c)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=400, message=Receiver testNewReceiver not found`)
}

func TestGetUpdateReceiverHandler(t *testing.T) {
	// Successful Update
	client := &mocks.AlertmanagerClient{}
	client.On("UpdateReceiver", testNID, sampleReceiver.Name, &sampleReceiver).Return(nil)
	client.On("ReloadAlertmanager").Return(nil)

	c, rec := buildContext(sampleReceiver, http.MethodPut, "/", v1receiverPath, testNID)
	c.SetParamNames(receiverNameParam)
	c.SetParamValues(sampleReceiver.Name)

	err := GetUpdateReceiverHandler(client, receiverNamePathProvider)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertExpectations(t)

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("UpdateReceiver", testNID, sampleReceiver.Name, &receivers.Receiver{}).Return(errors.New("error"))
	c, _ = buildContext(nil, http.MethodPut, "/", v1receiverPath, testNID)
	c.SetParamNames(receiverNameParam)
	c.SetParamValues(sampleReceiver.Name)

	err = GetUpdateReceiverHandler(client, receiverNamePathProvider)(c)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=400, message=error`)
	client.AssertExpectations(t)

	// Alertmanager Error
	client = &mocks.AlertmanagerClient{}
	client.On("UpdateReceiver", testNID, sampleReceiver.Name, &receivers.Receiver{}).Return(nil)
	client.On("ReloadAlertmanager").Return(errors.New("error"))
	c, _ = buildContext(nil, http.MethodPut, "/", v1receiverPath, testNID)
	c.SetParamNames(receiverNameParam)
	c.SetParamValues(sampleReceiver.Name)

	err = GetUpdateReceiverHandler(client, receiverNamePathProvider)(c)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=500, message=error`)
	client.AssertExpectations(t)
}

func TestGetDeleteReceiverHandler(t *testing.T) {
	// Successful Delete
	client := &mocks.AlertmanagerClient{}
	client.On("DeleteReceiver", testNID, sampleReceiver.Name).Return(nil)
	client.On("ReloadAlertmanager").Return(nil)

	c, rec := buildContext(nil, http.MethodGet, "/", v1receiverPath, testNID)
	c.SetParamNames(receiverNameParam)
	c.SetParamValues(sampleReceiver.Name)

	err := GetDeleteReceiverHandler(client, receiverNamePathProvider)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertExpectations(t)

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("DeleteReceiver", testNID, sampleReceiver.Name).Return(errors.New("error"))
	c, _ = buildContext(nil, http.MethodGet, "/", v1receiverPath, testNID)
	c.SetParamNames(receiverNameParam)
	c.SetParamValues(sampleReceiver.Name)

	err = GetDeleteReceiverHandler(client, receiverNamePathProvider)(c)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=400, message=error`)
	client.AssertExpectations(t)

	// Alertmanager Error
	client = &mocks.AlertmanagerClient{}
	client.On("DeleteReceiver", testNID, sampleReceiver.Name).Return(nil)
	client.On("ReloadAlertmanager").Return(errors.New("error"))
	c, _ = buildContext(nil, http.MethodGet, "/", v1receiverPath, testNID)
	c.SetParamNames(receiverNameParam)
	c.SetParamValues(sampleReceiver.Name)

	err = GetDeleteReceiverHandler(client, receiverNamePathProvider)(c)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=500, message=error`)
	client.AssertExpectations(t)
}

func TestGetGetRouteHandler(t *testing.T) {
	// Successful Get
	client := &mocks.AlertmanagerClient{}
	client.On("GetRoute", testNID).Return(&sampleRoute, nil)
	c, rec := buildContext(nil, http.MethodGet, "/", v1routePath, testNID)

	err := GetGetRouteHandler(client)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var retrievedRoute receivers.RouteJSONWrapper
	body, _ := ioutil.ReadAll(rec.Body)
	err = json.Unmarshal(body, &retrievedRoute)
	assert.NoError(t, err)
	assert.Equal(t, *receivers.NewRouteJSONWrapper(sampleRoute), retrievedRoute)
	client.AssertExpectations(t)

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("GetRoute", testNID).Return(nil, errors.New("error"))
	c, _ = buildContext(nil, http.MethodGet, "/", v1routePath, testNID)

	err = GetGetRouteHandler(client)(c)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=400, message=error`)
	client.AssertExpectations(t)
}

func TestGetUpdateRouteHandler(t *testing.T) {
	// Successful Update
	client := &mocks.AlertmanagerClient{}
	client.On("ModifyTenantRoute", testNID, &sampleRoute).Return(nil)
	client.On("ReloadAlertmanager").Return(nil)
	c, rec := buildContext(sampleRoute, http.MethodPost, "/", v1receiverPath, testNID)

	err := GetUpdateRouteHandler(client)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertExpectations(t)

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("ModifyTenantRoute", testNID, &sampleRoute).Return(errors.New("error"))
	c, _ = buildContext(sampleRoute, http.MethodPost, "/", v1receiverPath, testNID)

	err = GetUpdateRouteHandler(client)(c)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=400, message=error`)
	client.AssertExpectations(t)

	// Alertmanager Error
	client = &mocks.AlertmanagerClient{}
	client.On("ModifyTenantRoute", testNID, &sampleRoute).Return(nil)
	client.On("ReloadAlertmanager").Return(errors.New("error"))
	c, _ = buildContext(sampleRoute, http.MethodPost, "/", v1receiverPath, testNID)

	err = GetUpdateRouteHandler(client)(c)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=500, message=error`)
	client.AssertExpectations(t)
}

func TestGetGetGlobalConfigHandler(t *testing.T) {
	defaultConfig := config.DefaultGlobalConfig()
	// Successful Get
	client := &mocks.AlertmanagerClient{}
	client.On("GetGlobalConfig").Return(&defaultConfig, nil)
	c, rec := buildContext(nil, http.MethodGet, "/", v1routePath, testNID)

	err := GetGetGlobalConfigHandler(client)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var retrievedGlobalConfig config.GlobalConfig
	body, _ := ioutil.ReadAll(rec.Body)
	err = json.Unmarshal(body, &retrievedGlobalConfig)
	assert.NoError(t, err)
	assert.Equal(t, config.DefaultGlobalConfig(), retrievedGlobalConfig)
	client.AssertExpectations(t)

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("GetGlobalConfig").Return(nil, errors.New("error"))
	c, _ = buildContext(nil, http.MethodGet, "/", v1routePath, testNID)

	err = GetGetGlobalConfigHandler(client)(c)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=500, message=error`)
	client.AssertExpectations(t)
}

func TestGetUpdateGlobalConfigHandler(t *testing.T) {
	// Successful Update
	client := &mocks.AlertmanagerClient{}
	client.On("SetGlobalConfig", config.DefaultGlobalConfig()).Return(nil)
	client.On("ReloadAlertmanager").Return(nil)
	c, rec := buildContext(config.DefaultGlobalConfig(), http.MethodPost, "/", v1receiverPath, testNID)

	err := GetUpdateGlobalConfigHandler(client)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertExpectations(t)

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("SetGlobalConfig", config.DefaultGlobalConfig()).Return(errors.New("error"))
	c, _ = buildContext(config.DefaultGlobalConfig(), http.MethodPost, "/", v1receiverPath, testNID)

	err = GetUpdateGlobalConfigHandler(client)(c)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=400, message=error`)
	client.AssertExpectations(t)

	// Alertmanager Error
	client = &mocks.AlertmanagerClient{}
	client.On("SetGlobalConfig", config.DefaultGlobalConfig()).Return(nil)
	client.On("ReloadAlertmanager").Return(errors.New("error"))
	c, _ = buildContext(config.DefaultGlobalConfig(), http.MethodPost, "/", v1receiverPath, testNID)

	err = GetUpdateGlobalConfigHandler(client)(c)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	assert.EqualError(t, err, `code=500, message=error`)
	client.AssertExpectations(t)
}

func TestDecodeReceiverPostRequest(t *testing.T) {
	// Successful Decode
	c, _ := buildContext(sampleReceiver, http.MethodPost, "/", v1receiverPath, testNID)
	conf, err := decodeReceiverPostRequest(c)
	assert.NoError(t, err)
	assert.Equal(t, sampleReceiver, conf)

	// error decoding route
	c, _ = buildContext(struct {
		Name bool `json:"name"`
	}{false}, http.MethodPost, "/", v1receiverPath, testNID)
	conf, err = decodeReceiverPostRequest(c)
	assert.EqualError(t, err, `error unmarshalling payload: json: cannot unmarshal bool into Go struct field Receiver.name of type string`)
}

func TestDecodeRoutePostRequest(t *testing.T) {
	// Successful Decode
	c, _ := buildContext(sampleRoute, http.MethodPost, "/", v1receiverPath, testNID)
	conf, err := decodeRoutePostRequest(c)
	assert.NoError(t, err)
	assert.Equal(t, sampleRoute, conf)

	// Decode JSONWrapped Route
	c, _ = buildContext(sampleJSONRoute, http.MethodPost, "/", v1receiverPath, testNID)
	conf, err = decodeRoutePostRequest(c)
	assert.NoError(t, err)
	assert.Equal(t, convertedSampleJSONRoute, conf)

	// error decoding route
	c, _ = buildContext(struct {
		Receiver bool `json:"receiver"`
	}{false}, http.MethodPost, "/", v1receiverPath, testNID)
	conf, err = decodeRoutePostRequest(c)
	assert.EqualError(t, err, `error unmarshalling route: json: cannot unmarshal bool into Go struct field RouteJSONWrapper.receiver of type string`)
}

type tenancyTestCase struct {
	name           string
	client         *mocks.AlertmanagerClient
	tenantProvider paramProvider
	context        *echo.Context
	expectedTenant string
	expectedError  error
}

func TestTenancyMiddleware(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()

	plainReq := httptest.NewRequest(http.MethodGet, "/", nil)
	plainContext := e.NewContext(plainReq, rec)

	pathTenantContext := e.NewContext(plainReq, rec)
	pathTenantContext.SetParamNames(tenantIDParam)
	pathTenantContext.SetParamValues(testNID)

	mtClient := &mocks.AlertmanagerClient{}
	mtClient.On("Tenancy").Return(&alert.TenancyConfig{RestrictorLabel: testNID})

	tests := []tenancyTestCase{{
		name:           "multi-tenant with path provided tenant",
		client:         mtClient,
		tenantProvider: pathTenantProvider,
		context:        &pathTenantContext,
		expectedTenant: testNID,
	}, {
		name:           "multi-tenant without path provided tenant",
		client:         mtClient,
		tenantProvider: pathTenantProvider,
		context:        &plainContext,
		expectedError:  errors.New("code=400, message=Must provide tenant_id parameter"),
	}}

	for _, test := range tests {
		t.Run(test.name, test.runTest)
	}
}

func (tc *tenancyTestCase) runTest(t *testing.T) {
	handler := func(c echo.Context) error { return nil }

	tenancyFunc := tenancyMiddlewareProvider(tc.client, tc.tenantProvider)
	if tc.expectedError != nil {
		assert.EqualError(t, tenancyFunc(handler)(*tc.context), tc.expectedError.Error())
	} else {
		assert.NoError(t, tenancyFunc(handler)(*tc.context))
		assert.Equal(t, (*tc.context).Get(tenantIDParam), tc.expectedTenant)
	}
}

func buildContext(body interface{}, method, target, path, tenantID string) (echo.Context, *httptest.ResponseRecorder) {
	bytes, _ := json.Marshal(body)
	req := httptest.NewRequest(method, target, strings.NewReader(string(bytes)))
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetPath(path)
	c.SetParamNames("tenant_id")
	c.SetParamValues(tenantID)
	c.Set(tenantIDParam, tenantID)
	return c, rec
}
