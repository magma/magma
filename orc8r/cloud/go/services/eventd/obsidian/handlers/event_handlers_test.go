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
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

type elasticTestCase struct {
	name     string
	params   eventQueryParams
	expected *elastic.BoolQuery
}

type queryParamTestCase struct {
	name           string
	urlString      string
	paramNames     []string
	paramValues    []string
	expectsError   bool
	expectedParams eventQueryParams
}

type eventMapTestCase struct {
	name         string
	jsonSource   string
	expectedMaps []map[string]interface{}
	expectsError bool
}

var (
	elasticCases = []elasticTestCase{
		{
			name: "full query params",
			params: eventQueryParams{
				StreamName: "streamOne",
				EventType:  "an_event",
				HardwareID: "hardware-2",
				Tag:        "critical",
			},
			expected: elastic.NewBoolQuery().
				Filter(elastic.NewTermQuery("stream_name", "streamOne")).
				Filter(elastic.NewTermQuery("event_type", "an_event")).
				Filter(elastic.NewTermQuery("hardware_id", "hardware-2")).
				Filter(elastic.NewTermQuery("tag", "critical")),
		},
		{
			name: "partial query params",
			params: eventQueryParams{
				StreamName: "streamTwo",
				EventType:  "an_event",
			},
			expected: elastic.NewBoolQuery().
				Filter(elastic.NewTermQuery("stream_name", "streamTwo")).
				Filter(elastic.NewTermQuery("event_type", "an_event")),
		},
		{
			name: "only StreamName",
			params: eventQueryParams{
				StreamName: "streamThree",
			},
			expected: elastic.NewBoolQuery().
				Filter(elastic.NewTermQuery("stream_name", "streamThree")),
		},
	}

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

	eventMapTestCases = []eventMapTestCase{
		{
			name: "all fields",
			jsonSource: `{
				"stream_name": "a",
				"event_type": "b",
				"hw_id": "c",
				"tag": "d",
				"value":"{ \"some_property\": true }"
			}`,
			expectedMaps: []map[string]interface{}{
				{
					"stream_name":   "a",
					"event_type":    "b",
					"hardware_id":   "c",
					"tag":           "d",
					"some_property": true,
					"timestamp":     "",
				},
			},
		},
		{
			name: "partial fields with value present",
			jsonSource: `{
				"stream_name": "a",
				"event_type": "b",
				"value":"{}"
			}`,
			expectedMaps: []map[string]interface{}{
				{
					"stream_name": "a",
					"event_type":  "b",
					"hardware_id": "",
					"tag":         "",
					"timestamp":   "",
				},
			},
		},
		{
			name: "partial fields without a value",
			jsonSource: `{
				"stream_name": "a",
				"event_type": "b"
			}`,
			expectedMaps: []map[string]interface{}{},
			expectsError: true,
		},
	}
)

func TestElasticBoolQuery(t *testing.T) {
	for _, test := range elasticCases {
		t.Run(test.name, func(t *testing.T) {
			runToElasticBoolQueryTestCase(t, test)
		})
	}
}

func runToElasticBoolQueryTestCase(t *testing.T, tc elasticTestCase) {
	query := tc.params.ToElasticBoolQuery()
	assert.Equal(t, tc.expected, query)
}

// TestGetQueryParams tests that parameters in the url are parsed correctly
func TestGetQueryParams(t *testing.T) {
	for _, test := range queryParamsTestCases {
		t.Run(test.name, func(t *testing.T) {
			runQueryParamTestCase(t, test)
		})
	}
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

func TestGetEventMap(t *testing.T) {
	for _, test := range eventMapTestCases {
		t.Run(test.name, func(t *testing.T) {
			runEventMapTestCase(t, test)
		})
	}
}

func runEventMapTestCase(t *testing.T, tc eventMapTestCase) {
	hit := elastic.SearchHit{
		Source: []byte(tc.jsonSource),
	}
	maps, err := getEventMaps([]*elastic.SearchHit{&hit})
	if tc.expectsError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedMaps, maps)
	}
}
