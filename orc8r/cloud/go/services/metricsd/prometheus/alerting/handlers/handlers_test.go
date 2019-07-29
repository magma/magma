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
	"strings"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alert"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alert/mocks"

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
