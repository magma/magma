/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	eventdC "magma/orc8r/cloud/go/services/eventd/eventd_client"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type queryParamTestCase struct {
	name           string
	urlString      string
	paramNames     []string
	paramValues    []string
	expectsError   bool
	expectedParams eventdC.EventQueryParams
}

type queryMultiStreamParamTestCase struct {
	name           string
	urlString      string
	paramNames     []string
	paramValues    []string
	expectsError   bool
	expectedParams eventdC.MultiStreamEventQueryParams
}

var (
	queryParamsTestCases = []queryParamTestCase{
		{
			name:           "no params will error",
			urlString:      "",
			expectsError:   true,
			expectedParams: eventdC.EventQueryParams{},
		},
		{
			name:         "only stream_name",
			urlString:    "",
			paramNames:   []string{"stream_name"},
			paramValues:  []string{"streamOne"},
			expectsError: true,
		},
		{
			name:         "excess path params",
			urlString:    "",
			paramNames:   []string{"stream_name", "network_id", "bad_param"},
			paramValues:  []string{"streamOne", "nw1", "bad_value"},
			expectsError: false,
			expectedParams: eventdC.EventQueryParams{
				StreamName: "streamOne",
				NetworkID:  "nw1",
			},
		},
		{
			name:        "all query params",
			urlString:   "/nwk1?event_type=mock_subscriber_event&hardware_id=123&tag=critical",
			paramNames:  []string{"stream_name", "network_id"},
			paramValues: []string{"streamOne", "nw1"},
			expectedParams: eventdC.EventQueryParams{
				StreamName: "streamOne",
				EventType:  "mock_subscriber_event",
				HardwareID: "123",
				NetworkID:  "nw1",
				Tag:        "critical",
			},
		},
	}

	multiStreamQueryParamsTestCases = []queryMultiStreamParamTestCase{
		{
			name:           "no params will error",
			urlString:      "",
			expectsError:   true,
			expectedParams: eventdC.MultiStreamEventQueryParams{},
		},
		{
			name:         "only stream names",
			urlString:    "",
			paramNames:   []string{"streams"},
			paramValues:  []string{"streamOne,streamTwo"},
			expectsError: true,
		},
		{
			name:         "excess path params",
			urlString:    "",
			paramNames:   []string{"network_id", "bad_param"},
			paramValues:  []string{"nw1", "bad_value"},
			expectsError: false,
			expectedParams: eventdC.MultiStreamEventQueryParams{
				Streams:     []string{},
				NetworkID:   "nw1",
				HardwareIDs: []string{},
				Events:      []string{},
				Tags:        []string{},
				Size:        50,
			},
		},
		{
			name:        "all query params",
			urlString:   "?events=mock_subscriber_event&tags=critical&streams=streamOne,streamTwo",
			paramNames:  []string{"network_id"},
			paramValues: []string{"nw1"},
			expectedParams: eventdC.MultiStreamEventQueryParams{
				Streams:     []string{"streamOne", "streamTwo"},
				Events:      []string{"mock_subscriber_event"},
				NetworkID:   "nw1",
				Tags:        []string{"critical"},
				HardwareIDs: []string{},
				Size:        50,
			},
		},
	}
)

func TestEventsPath(t *testing.T) {
	assert.Equal(t,
		"/magma/v1/events/:network_id/:stream_name",
		EventsPath)
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
		assert.Equal(t, tc.expectedParams, params)
	}
}

// TestGetQueryParams tests that parameters in the url are parsed correctly
func TestGetMultiStreamQueryParams(t *testing.T) {
	for _, test := range multiStreamQueryParamsTestCases {
		t.Run(test.name, func(t *testing.T) {
			runMultiStreamQueryParamTestCase(t, test)
		})
	}
}

func runMultiStreamQueryParamTestCase(t *testing.T, tc queryMultiStreamParamTestCase) {
	req := httptest.NewRequest(echo.GET, fmt.Sprintf("/%s", tc.urlString), nil)
	c := echo.New().NewContext(req, httptest.NewRecorder())
	c.SetParamNames(tc.paramNames...)
	c.SetParamValues(tc.paramValues...)

	params, err := getMultiStreamQueryParameters(c)
	if tc.expectsError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedParams, params)
	}
}
