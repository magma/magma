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
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alert"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alert/mocks"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/handlers"

	"github.com/labstack/echo"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	sampleDuration, _ = model.ParseDuration("5s")
	sampleAlert1      = rulefmt.Rule{
		Alert:       "testAlert1",
		For:         sampleDuration,
		Expr:        "up == 0",
		Labels:      map[string]string{"label": "value"},
		Annotations: map[string]string{"annotation": "value"},
	}
	sampleAlert2 = rulefmt.Rule{
		Alert:       "testAlert2",
		For:         sampleDuration,
		Expr:        "up == 0",
		Labels:      map[string]string{"label": "value"},
		Annotations: map[string]string{"annotation": "value"},
	}
)

func TestGetConfigureAlertHandler(t *testing.T) {
	client := &mocks.PrometheusAlertClient{}
	client.On("ValidateRule", sampleAlert1).Return(nil)
	client.On("RuleExists", testNID, sampleAlert1.Alert).Return(false)
	client.On("WriteRule", testNID, sampleAlert1).Return(nil)

	configureAlert := GetConfigureAlertHandler(client, "")

	c, rec := buildContext(sampleAlert1, http.MethodPost, "/", handlers.AlertConfigURL, testNID)

	err := configureAlert(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertCalled(t, "ValidateRule", sampleAlert1)
	client.AssertCalled(t, "RuleExists", testNID, sampleAlert1.Alert)
	client.AssertCalled(t, "WriteRule", testNID, sampleAlert1)
}

func TestGetRetrieveAlertHandler(t *testing.T) {
	client := &mocks.PrometheusAlertClient{}
	client.On("ReadRules", testNID, "").Return([]rulefmt.Rule{sampleAlert1}, nil)

	retrieveAlert := GetRetrieveAlertHandler(client)

	c, rec := buildContext(sampleAlert1, http.MethodPost, "/", handlers.AlertConfigURL, testNID)

	err := retrieveAlert(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertCalled(t, "ReadRules", testNID, "")
}

func TestGetDeleteAlertHandler(t *testing.T) {
	client := &mocks.PrometheusAlertClient{}
	client.On("DeleteRule", testNID, sampleAlert1.Alert).Return(nil)

	deleteAlert := GetDeleteAlertHandler(client, "")

	q := make(url.Values)
	q.Set(handlers.AlertNameQueryParam, sampleAlert1.Alert)
	c, rec := buildContext(nil, http.MethodDelete, "/?"+q.Encode(), handlers.AlertConfigURL, testNID)

	err := deleteAlert(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertCalled(t, "DeleteRule", testNID, sampleAlert1.Alert)

	// No alert name given
	c, _ = buildContext(nil, http.MethodDelete, "/", handlers.AlertConfigURL, testNID)
	err = deleteAlert(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	// assert DeleteRule hasn't been called again
	client.AssertNumberOfCalls(t, "DeleteRule", 1)
}

func TestUpdateAlertHandler(t *testing.T) {
	client := &mocks.PrometheusAlertClient{}
	client.On("RuleExists", testNID, sampleAlert1.Alert).Return(true)
	client.On("ValidateRule", sampleAlert1).Return(nil)
	client.On("UpdateRule", testNID, sampleAlert1).Return(nil)

	updateAlert := GetUpdateAlertHandler(client, "")

	c, rec := buildContext(sampleAlert1, http.MethodPut, "/", handlers.AlertConfigURL, testNID)
	c.SetParamNames("network_id", RuleNamePathParam)
	c.SetParamValues(testNID, sampleAlert1.Alert)

	err := updateAlert(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertCalled(t, "RuleExists", testNID, sampleAlert1.Alert)
	client.AssertCalled(t, "ValidateRule", sampleAlert1)
	client.AssertCalled(t, "UpdateRule", testNID, sampleAlert1)
}

func TestGetBulkAlertUpdateHandler(t *testing.T) {
	client := &mocks.PrometheusAlertClient{}
	bulkAlerts := []rulefmt.Rule{sampleAlert1, sampleAlert2}
	sampleUpdateResult := alert.BulkUpdateResults{
		Errors:   map[string]error{},
		Statuses: map[string]string{"testAlert1": "created", "testAlert2": "created"},
	}
	client.On("BulkUpdateRules", testNID, bulkAlerts).Return(sampleUpdateResult, nil)
	client.On("ValidateRule", mock.Anything).Return(nil)

	bulkUpdateFunc := GetBulkAlertUpdateHandler(client, "")

	bytes, err := json.Marshal([]rulefmt.Rule{sampleAlert1, sampleAlert2})
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(bytes)))
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)
	c.SetPath("/networks/:network_id/prometheus/alert_config/bulk")
	c.SetParamNames("network_id")
	c.SetParamValues(testNID)

	err = bulkUpdateFunc(c)
	assert.NoError(t, err)
	client.AssertCalled(t, "BulkUpdateRules", testNID, bulkAlerts)
	client.AssertCalled(t, "ValidateRule", mock.Anything)
	assert.Equal(t, http.StatusOK, rec.Code)

	var results alert.BulkUpdateResults
	err = json.Unmarshal(rec.Body.Bytes(), &results)
	assert.NoError(t, err)
	assert.Equal(t, sampleUpdateResult, results)
}
