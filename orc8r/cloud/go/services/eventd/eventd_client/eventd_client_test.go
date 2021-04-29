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

package eventd_client

import (
	"testing"

	"magma/orc8r/cloud/go/services/eventd/obsidian/models"

	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

type elasticTestCase struct {
	name     string
	params   EventQueryParams
	expected *elastic.BoolQuery
}

type eventResultTestCase struct {
	name            string
	jsonSource      string
	expectedResults []models.Event
	expectsError    bool
}

var (
	elasticCases = []elasticTestCase{
		{
			name: "full query params",
			params: EventQueryParams{
				StreamName: "streamOne",
				EventType:  "an_event",
				HardwareID: "hardware-2",
				NetworkID:  "test_network",
				Tag:        "critical",
			},
			expected: elastic.NewBoolQuery().
				Filter(elastic.NewTermQuery("stream_name.keyword", "streamOne")).
				Filter(elastic.NewTermQuery("network_id.keyword", "test_network")).
				Filter(elastic.NewTermQuery("event_type.keyword", "an_event")).
				Filter(elastic.NewTermQuery("hw_id.keyword", "hardware-2")).
				Filter(elastic.NewTermQuery("event_tag.keyword", "critical")),
		},
		{
			name: "partial query params",
			params: EventQueryParams{
				StreamName: "streamTwo",
				NetworkID:  "test_network_two",
				EventType:  "an_event",
			},
			expected: elastic.NewBoolQuery().
				Filter(elastic.NewTermQuery("stream_name.keyword", "streamTwo")).
				Filter(elastic.NewTermQuery("network_id.keyword", "test_network_two")).
				Filter(elastic.NewTermQuery("event_type.keyword", "an_event")),
		},
		{
			name: "only required Path params",
			params: EventQueryParams{
				StreamName: "streamThree",
				NetworkID:  "test_network_three",
			},
			expected: elastic.NewBoolQuery().
				Filter(elastic.NewTermQuery("stream_name.keyword", "streamThree")).
				Filter(elastic.NewTermQuery("network_id.keyword", "test_network_three")),
		},
	}

	eventResultTestCases = []eventResultTestCase{
		{
			name: "all fields",
			jsonSource: `{
				"stream_name": "a",
				"event_type": "b",
				"hw_id": "c",
				"event_tag": "d",
				"value":"{ \"some_property\": true }"
			}`,
			expectedResults: []models.Event{
				{
					StreamName: "a",
					EventType:  "b",
					HardwareID: "c",
					Tag:        "d",
					Value:      map[string]interface{}{"some_property": true},
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
			expectedResults: []models.Event{
				{
					StreamName: "a",
					EventType:  "b",
					Value:      map[string]interface{}{},
				},
			},
		},
		{
			name: "partial fields without a value",
			jsonSource: `{
				"stream_name": "a",
				"event_type": "b"
			}`,
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
	query := tc.params.toElasticBoolQuery()
	assert.Equal(t, tc.expected, query)
}

func TestGetEventResults(t *testing.T) {
	for _, test := range eventResultTestCases {
		t.Run(test.name, func(t *testing.T) {
			runEventResultTestCase(t, test)
		})
	}
}

func runEventResultTestCase(t *testing.T, tc eventResultTestCase) {
	hit := elastic.SearchHit{
		Source: []byte(tc.jsonSource),
	}
	results, err := GetEventResults([]*elastic.SearchHit{&hit})
	if tc.expectsError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedResults, results)
	}
}
