/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collector

import (
	"strconv"
	"testing"

	"magma/lte/cloud/go/services/nprobe"

	"github.com/stretchr/testify/assert"
)

var (
	tcQueriesParams = []multiStreamEventQueryParams{
		{
			networkID: "test1",
			events:    nprobe.GetESEventTypes(),
			streams:   nprobe.GetESStreams(),
			tags:      []string{"test_subscriber1"},
			timestamp: 0,
		},
		{
			networkID: "test2",
			events:    nprobe.GetESEventTypes(),
			streams:   nprobe.GetESStreams(),
			tags:      []string{"test_subscriber2"},
			timestamp: 5,
		},
		{
			networkID: "test3",
			events:    nprobe.GetESEventTypes(),
			streams:   nprobe.GetESStreams(),
			tags:      []string{"test_subscriber3"},
			timestamp: 10,
		},
	}

	queryFilters = map[string]int{
		"event_tag.keyword":   1,
		"event_type.keyword":  len(nprobe.GetESEventTypes()),
		"stream_name.keyword": len(nprobe.GetESStreams()),
	}
)

// TestMultiStreamsQuery tests that various combinations of query parameters
// results in a valid and expected ElasticSearch Query
func TestMultiStreamsQuery(t *testing.T) {
	for _, tc := range tcQueriesParams {
		runMultiStreamsQueryTestCase(t, tc)
	}
}

func runMultiStreamsQueryTestCase(t *testing.T, params multiStreamEventQueryParams) {
	query := params.toElasticBoolQuery()
	source, err := query.Source()
	assert.NoError(t, err)

	bQuery, ok := source.(map[string]interface{})["bool"].(map[string]interface{})
	assert.True(t, ok)

	if params.timestamp > 0 {
		must, ok := bQuery["must"].(map[string]interface{})
		assert.True(t, ok)
		assert.Len(t, must, 1)

		rangeQuery, ok := must["range"].(map[string]interface{})
		assert.True(t, ok)

		timeQuery, ok := rangeQuery["@timestamp"].(map[string]interface{})
		assert.True(t, ok)

		assert.Equal(t, timeQuery["format"], "epoch_millis")
		assert.Equal(t, timeQuery["from"], strconv.FormatInt(params.timestamp, 10))
	} else {
		_, ok := bQuery["must"]
		assert.False(t, ok)
	}

	if filters, ok := bQuery["filter"].([]interface{}); ok {
		for expectedKey, expectedVal := range queryFilters {
			filterExists := false
			for _, filter := range filters {
				terms, ok := filter.(map[string]interface{})["terms"].(map[string]interface{})
				assert.True(t, ok)
				values, ok := terms[expectedKey].([]interface{})
				if ok && len(values) == expectedVal {
					filterExists = true
					break
				}
			}
			assert.True(t, filterExists)
		}
	}
}
