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
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/receivers"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/receivers/mocks"

	"github.com/labstack/echo"
	"github.com/prometheus/alertmanager/config"
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
		WebhookConfigs: []*config.WebhookConfig{{
			NotifierConfig: config.NotifierConfig{
				VSendResolved: true,
			},
			URL: &config.URL{
				URL: testWebhookURL,
			},
		}},
	}

	sampleRoute = config.Route{
		Receiver: "testSlackReceiver",
		Match:    map[string]string{"networkID": testNID},
		Routes: []*config.Route{{
			Receiver: "childReceiver",
			Match:    map[string]string{"severity": "critical"},
		}},
	}
)

func TestGetReceiverPostHandler(t *testing.T) {
	client := &mocks.AlertmanagerClient{}
	postReceiver := GetReceiverPostHandler(client, "")
	client.On("CreateReceiver", testNID, sampleReceiver).Return(nil)

	c, rec := buildContext(sampleReceiver, http.MethodPost, "/", ReceiverPath, testNID)

	err := postReceiver(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertCalled(t, "CreateReceiver", testNID, sampleReceiver)

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("CreateReceiver", testNID, receivers.Receiver{}).Return(echo.NewHTTPError(http.StatusBadRequest, "error"))
	postReceiver = GetReceiverPostHandler(client, "")
	c, _ = buildContext(nil, http.MethodPost, "/", ReceiverPath, testNID)
	err = postReceiver(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	client.AssertCalled(t, "CreateReceiver", testNID, receivers.Receiver{})
}

func TestGetGetReceiversHandler(t *testing.T) {
	client := &mocks.AlertmanagerClient{}
	client.On("GetReceivers", testNID).Return([]receivers.Receiver{sampleReceiver}, nil)
	getReceivers := GetGetReceiversHandler(client)

	c, rec := buildContext(nil, http.MethodGet, "/", ReceiverPath, testNID)

	err := getReceivers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertCalled(t, "GetReceivers", testNID)

	var receiver []receivers.Receiver
	err = json.Unmarshal(rec.Body.Bytes(), &receiver)
	assert.Equal(t, 1, len(receiver))
	assert.Equal(t, sampleReceiver, receiver[0])

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("GetReceivers", testNID).Return([]receivers.Receiver{}, errors.New("error"))
	getReceivers = GetGetReceiversHandler(client)
	c, _ = buildContext(nil, http.MethodGet, "/", ReceiverPath, testNID)
	err = getReceivers(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	client.AssertCalled(t, "GetReceivers", testNID)
}

func TestGetUpdateReceiverHandler(t *testing.T) {
	client := &mocks.AlertmanagerClient{}
	client.On("UpdateReceiver", testNID, &sampleReceiver).Return(nil)
	updateReceiver := GetUpdateReceiverHandler(client, "")

	c, rec := buildContext(sampleReceiver, http.MethodPut, "/", ReceiverPath, testNID)

	err := updateReceiver(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertCalled(t, "UpdateReceiver", testNID, &sampleReceiver)

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("UpdateReceiver", testNID, &receivers.Receiver{}).Return(errors.New("error"))
	updateReceiver = GetUpdateReceiverHandler(client, "")
	c, _ = buildContext(nil, http.MethodPut, "/", ReceiverPath, testNID)

	err = updateReceiver(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	client.AssertCalled(t, "UpdateReceiver", testNID, &receivers.Receiver{})
}

func TestGetDeleteReceiverHandler(t *testing.T) {
	client := &mocks.AlertmanagerClient{}
	client.On("DeleteReceiver", testNID, sampleReceiver.Name).Return(nil)

	deleteReceiver := GetDeleteReceiverHandler(client, "")

	q := make(url.Values)
	q.Set(ReceiverNameQueryParam, sampleReceiver.Name)
	c, rec := buildContext(nil, http.MethodGet, "/?"+q.Encode(), ReceiverPath, testNID)

	err := deleteReceiver(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertCalled(t, "DeleteReceiver", testNID, sampleReceiver.Name)

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("DeleteReceiver", testNID, sampleReceiver.Name).Return(errors.New("error"))
	deleteReceiver = GetDeleteReceiverHandler(client, "")

	c, _ = buildContext(nil, http.MethodGet, "/?"+q.Encode(), ReceiverPath, testNID)

	err = deleteReceiver(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	client.AssertCalled(t, "DeleteReceiver", testNID, sampleReceiver.Name)
}

func TestGetGetRouteHandler(t *testing.T) {
	client := &mocks.AlertmanagerClient{}
	client.On("GetRoute", testNID).Return(&sampleRoute, nil)
	getRoute := GetGetRouteHandler(client)

	c, rec := buildContext(nil, http.MethodGet, "/", RoutePath, testNID)

	err := getRoute(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertCalled(t, "GetRoute", testNID)

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("GetRoute", testNID).Return(nil, errors.New("error"))
	getRoute = GetGetRouteHandler(client)
	c, _ = buildContext(nil, http.MethodGet, "/", RoutePath, testNID)

	err = getRoute(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	client.AssertCalled(t, "GetRoute", testNID)
}

func TestGetUpdateRouteHandler(t *testing.T) {
	client := &mocks.AlertmanagerClient{}
	client.On("ModifyNetworkRoute", testNID, &sampleRoute).Return(nil)
	updateRoute := GetUpdateRouteHandler(client, "")

	c, rec := buildContext(sampleRoute, http.MethodPost, "/", ReceiverPath, testNID)

	err := updateRoute(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertCalled(t, "ModifyNetworkRoute", testNID, &sampleRoute)

	// Client Error
	client = &mocks.AlertmanagerClient{}
	client.On("ModifyNetworkRoute", testNID, &sampleRoute).Return(errors.New("error"))
	updateRoute = GetUpdateRouteHandler(client, "")
	c, _ = buildContext(sampleRoute, http.MethodPost, "/", ReceiverPath, testNID)

	err = updateRoute(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	client.AssertCalled(t, "ModifyNetworkRoute", testNID, &sampleRoute)
}

func TestDecodeReceiverPostRequest(t *testing.T) {
	c, _ := buildContext(sampleReceiver, http.MethodPost, "/", ReceiverPath, testNID)

	conf, err := decodeReceiverPostRequest(c)
	assert.NoError(t, err)
	assert.Equal(t, sampleReceiver, conf)
}

func TestDecodeRoutePostRequest(t *testing.T) {
	c, _ := buildContext(sampleRoute, http.MethodPost, "/", ReceiverPath, testNID)

	conf, err := decodeRoutePostRequest(c)
	assert.NoError(t, err)
	assert.Equal(t, sampleRoute, conf)
}

func buildContext(body interface{}, method, target, path, networkID string) (echo.Context, *httptest.ResponseRecorder) {
	bytes, _ := json.Marshal(body)
	req := httptest.NewRequest(method, target, strings.NewReader(string(bytes)))
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetPath(path)
	c.SetParamNames("file_prefix")
	c.SetParamValues(networkID)
	return c, rec
}
