/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/restrictor"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type seriesHandlerTestCase struct {
	name            string
	inputURL        string
	restrictor      restrictor.QueryRestrictor
	expectedStrings []string
}

func (tc *seriesHandlerTestCase) RunTest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, tc.inputURL, nil)
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)
	params, err := getSeriesMatches(c, tc.restrictor)
	assert.NoError(t, err)
	assert.Equal(t, tc.expectedStrings, params)
}

func TestGetPrometheusSeriesHandler(t *testing.T) {
	testCases := []seriesHandlerTestCase{
		{
			name:            "single match",
			inputURL:        "/?match=up",
			restrictor:      networkQueryRestrictorProvider("test"),
			expectedStrings: []string{`up{networkID="test"}`},
		},
		{
			name:            "two match",
			inputURL:        "/?match=up%20down",
			restrictor:      networkQueryRestrictorProvider("test"),
			expectedStrings: []string{`up{networkID="test"}`, `down{networkID="test"}`},
		},
		{
			name:            "complicated match",
			inputURL:        "/?match=up%20down%20{gatewayID=\"gw1\"}",
			restrictor:      networkQueryRestrictorProvider("test"),
			expectedStrings: []string{`up{networkID="test"}`, `down{networkID="test"}`, `{gatewayID="gw1",networkID="test"}`},
		},
		{
			name:            "no match",
			inputURL:        "/",
			restrictor:      networkQueryRestrictorProvider("test"),
			expectedStrings: []string{`{networkID="test"}`},
		},
		{
			name:            "tenant match",
			inputURL:        "/",
			restrictor:      *restrictor.NewQueryRestrictor(restrictor.Opts{ReplaceExistingLabel: false}).AddMatcher("networkID", "net1", "net2"),
			expectedStrings: []string{`{networkID=~"net1|net2"}`},
		},
		{
			name:            "tenant two match",
			inputURL:        "/?match=up%20down",
			restrictor:      *restrictor.NewQueryRestrictor(restrictor.Opts{ReplaceExistingLabel: false}).AddMatcher("networkID", "net1", "net2"),
			expectedStrings: []string{`up{networkID=~"net1|net2"}`, `down{networkID=~"net1|net2"}`},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.RunTest)
	}
}
