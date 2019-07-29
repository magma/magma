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

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/receivers"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/receivers/mocks"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	testNID = "test"
)

var (
	sampleReceiver = receivers.Receiver{
		Name: "testSlackReceiver",
		SlackConfigs: []*receivers.SlackConfig{{
			APIURL:   "http://slack.com/12345",
			Channel:  "test_channel",
			Username: "test_username",
		}},
	}
)

func TestGetReceiverPostHandler(t *testing.T) {
	client := &mocks.AlertmanagerClient{}
	postReceiver := GetReceiverPostHandler(client, "")
	client.On("CreateReceiver", testNID, sampleReceiver).Return(nil)

	bytes, err := json.Marshal(sampleReceiver)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(bytes)))
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)
	c.SetPath("/:network_id/receiver")
	c.SetParamNames("network_id")
	c.SetParamValues(testNID)

	err = postReceiver(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	client.AssertCalled(t, "CreateReceiver", testNID, sampleReceiver)
}

func TestGetGetReceiversHandler(t *testing.T) {
	client := &mocks.AlertmanagerClient{}
	client.On("GetReceivers", testNID).Return([]receivers.Receiver{sampleReceiver}, nil)

	getReceivers := GetGetReceiversHandler(client)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)
	c.SetPath("/:network_id/receiver")
	c.SetParamNames("network_id")
	c.SetParamValues(testNID)

	err := getReceivers(c)
	assert.NoError(t, err)
	client.AssertCalled(t, "GetReceivers", testNID)

	var receiver []receivers.Receiver
	err = json.Unmarshal(rec.Body.Bytes(), &receiver)
	assert.Equal(t, 1, len(receiver))
	assert.Equal(t, sampleReceiver, receiver[0])
}

func TestGetDeleteReceiverHandler(t *testing.T) {
	client := &mocks.AlertmanagerClient{}
	client.On("DeleteReceiver", testNID, sampleReceiver.Name).Return(nil)

	deleteReceiver := GetDeleteReceiverHandler(client, "")

	q := make(url.Values)
	q.Set(ReceiverNameQueryParam, sampleReceiver.Name)
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)
	c.SetPath("/:network_id/receiver")
	c.SetParamNames("network_id")
	c.SetParamValues(testNID)

	err := deleteReceiver(c)
	assert.NoError(t, err)
	client.AssertCalled(t, "DeleteReceiver", testNID, sampleReceiver.Name)
}
