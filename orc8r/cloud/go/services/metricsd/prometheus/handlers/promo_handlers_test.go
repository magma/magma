/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"magma/orc8r/cloud/go/metrics"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestPreprocessQuery(t *testing.T) {
	testQuery := "up"
	networkID := "network1"
	preprocessedQuery, err := preprocessQuery(testQuery, networkID)
	assert.NoError(t, err)
	expectedQuery := fmt.Sprintf("%s{%s=\"%s\"}", testQuery, metrics.NetworkLabelName, networkID)
	assert.Equal(t, expectedQuery, preprocessedQuery)
}

type seriesHandlerTestCase struct {
	name            string
	inputURL        string
	nID             string
	expectedStrings []string
}

func (tc *seriesHandlerTestCase) RunTest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, tc.inputURL, nil)
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)
	params, err := getSeriesMatches(c, tc.nID)
	assert.NoError(t, err)
	assert.Equal(t, tc.expectedStrings, params)
}

func TestGetPrometheusSeriesHandler(t *testing.T) {

	testCases := []seriesHandlerTestCase{
		{
			name:            "single match",
			inputURL:        "/?match=up",
			nID:             "test",
			expectedStrings: []string{`up{networkID="test"}`},
		},
		{
			name:            "two match",
			inputURL:        "/?match=up%20down",
			nID:             "test",
			expectedStrings: []string{`up{networkID="test"}`, `down{networkID="test"}`},
		},
		{
			name:            "complicated match",
			inputURL:        "/?match=up%20down%20{gatewayID=\"gw1\"}",
			nID:             "test",
			expectedStrings: []string{`up{networkID="test"}`, `down{networkID="test"}`, `{gatewayID="gw1",networkID="test"}`},
		},
		{
			name:            "no match",
			inputURL:        "/",
			nID:             "test",
			expectedStrings: []string{`{networkID="test"}`},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.RunTest)
	}
}
