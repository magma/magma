/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var (
	queryParamsTestCases = []queryParamTestCase{
		{
			name:           "no params will error",
			urlString:      "",
			expectsError:   true,
			expectedParams: eventQueryParams{},
		},
		{
			name:        "only stream_name",
			urlString:   "",
			paramNames:  []string{"stream_name"},
			paramValues: []string{"streamOne"},
			expectedParams: eventQueryParams{
				StreamName: "streamOne",
			},
		},
		{
			name:         "excess path params",
			urlString:    "",
			paramNames:   []string{"stream_name", "bad_param"},
			paramValues:  []string{"streamOne", "bad_value"},
			expectsError: false,
			expectedParams: eventQueryParams{
				StreamName: "streamOne",
			},
		},
		{
			name:        "all query params",
			urlString:   "?event_type=mock_subscriber_event&hardware_id=123&tag=critical",
			paramNames:  []string{"stream_name"},
			paramValues: []string{"streamOne"},
			expectedParams: eventQueryParams{
				StreamName: "streamOne",
				EventType:  "mock_subscriber_event",
				HardwareID: "123",
				Tag:        "critical",
			},
		},
	}
)

// TestGetQueryParams tests that parameters in the url are parsed correctly
func TestGetQueryParams(t *testing.T) {
	for _, test := range queryParamsTestCases {
		t.Run(test.name, func(t *testing.T) {
			runQueryParamTestCase(t, test)
		})
	}
}

type queryParamTestCase struct {
	name           string
	urlString      string
	paramNames     []string
	paramValues    []string
	expectsError   bool
	expectedParams eventQueryParams
}

func runQueryParamTestCase(t *testing.T, tc queryParamTestCase) {
	req := httptest.NewRequest(echo.GET, fmt.Sprintf("/%s", tc.urlString), nil)
	c := echo.New().NewContext(req, httptest.NewRecorder())
	c.SetParamNames(tc.paramNames...)
	c.SetParamValues(tc.paramValues...)

	params, err := getQueryParameters(c)
	if tc.expectsError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
	assert.Equal(t, tc.expectedParams, params)
}
