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

package collection_test

import (
	"testing"

	"magma/feg/gateway/services/radiusd/collection"

	"github.com/stretchr/testify/assert"
)

const VALID_PROMETHEUS_METRICS = `
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 19

# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.12.7"} 1

# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 11979

# HELP listener_init_success_count The number of time 'listener_init' has succeeded
# TYPE listener_init_success_count counter
listener_init_success_count{listener=""} 4
`

const INVALID_PROMETHEUS_METRICS = `
# TYPE go_info gauge
go_info{version="go1.12.7" asdfasdf} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 3.01536e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
`

const EXPECTED_ERROR_TEXT = "Error parsing metric families from text: text format parsing error in line 3: unexpected end of label value \"go1.12.7\""

func TestParsePrometheusText(t *testing.T) {
	_, err := collection.ParsePrometheusText(INVALID_PROMETHEUS_METRICS)
	assert.Error(t, err)
	assert.EqualError(t, err, EXPECTED_ERROR_TEXT)

	metricFamilies, err := collection.ParsePrometheusText(VALID_PROMETHEUS_METRICS)
	assert.NoError(t, err)
	assert.Len(t, metricFamilies, 4)

	// go_goroutines
	metricFamily, ok := metricFamilies["go_goroutines"]
	assert.True(t, ok)
	metricArr := metricFamily.GetMetric()
	assert.Len(t, metricArr, 1)
	for _, metricSample := range metricArr {
		inp := metricSample.GetGauge()
		val := inp.GetValue()
		expected := float64(19)
		assert.Equal(t, expected, val)
	}

	// go_info
	metricFamily, ok = metricFamilies["go_info"]
	assert.True(t, ok)
	metricArr = metricFamily.GetMetric()
	assert.Len(t, metricArr, 1)
	for _, metricSample := range metricArr {
		inp := metricSample.GetGauge()
		val := inp.GetValue()
		expected := float64(1)
		assert.Equal(t, expected, val)
	}

	// go_memstats_frees_total
	metricFamily, ok = metricFamilies["go_memstats_frees_total"]
	assert.True(t, ok)
	metricArr = metricFamily.GetMetric()
	assert.Len(t, metricArr, 1)
	for _, metricSample := range metricArr {
		inp := metricSample.GetCounter()
		val := inp.GetValue()
		expected := float64(11979)
		assert.Equal(t, expected, val)
	}

	// listener_init_success_count
	metricFamily, ok = metricFamilies["listener_init_success_count"]
	assert.True(t, ok)
	metricArr = metricFamily.GetMetric()
	assert.Len(t, metricArr, 1)
	for _, metricSample := range metricArr {
		inp := metricSample.GetCounter()
		val := inp.GetValue()
		expected := float64(4)
		assert.Equal(t, expected, val)
	}
}
