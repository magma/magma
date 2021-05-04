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
	"strings"
	"testing"

	"magma/feg/gateway/services/radiusd/collection"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
)

const GAUGE_METRIC = `
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 19
`
const GAUGE_VEC_METRIC = `
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.12.7"} 1
`
const COUNTER_METRIC = `
# HELP radius_go_memstats_frees_total Total number of frees.
# TYPE radius_go_memstats_frees_total counter
radius_go_memstats_frees_total 11979
`
const COUNTER_VEC_METRIC = `
# HELP listener_init_success_count The number of time 'listener_init' has succeeded
# TYPE listener_init_success_count counter
listener_init_success_count{listener=""} 4
`
const SUMMARY_METRIC = `
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0.000215687
go_gc_duration_seconds{quantile="0.25"} 0.000215687
go_gc_duration_seconds{quantile="0.5"} 0.000290718
go_gc_duration_seconds{quantile="0.75"} 0.000290718
go_gc_duration_seconds{quantile="1"} 0.000290718
go_gc_duration_seconds_sum 0.000506405
go_gc_duration_seconds_count 2
`

const GAUGE_METRIC_UPDATED = `
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 3
`

const COUNTER_METRIC_UPDATED = `
# HELP radius_go_memstats_frees_total Total number of frees.
# TYPE radius_go_memstats_frees_total counter
radius_go_memstats_frees_total 1
`

// getMetricFamily reads in a single prometheus metric and returns a
// MetricFamily.
// If there are multiple, it returns one MetricFamily at random
func getMetricFamily(prometheusText string) *dto.MetricFamily {
	reader := strings.NewReader(prometheusText)
	parser := expfmt.TextParser{}
	metricFamilies, _ := parser.TextToMetricFamilies(reader)
	for _, metricFamily := range metricFamilies {
		return metricFamily
	}
	return nil
}

func newGaugeFamily() *dto.MetricFamily {
	return getMetricFamily(GAUGE_METRIC)
}

func newGaugeVecFamily() *dto.MetricFamily {
	return getMetricFamily(GAUGE_VEC_METRIC)
}

func newCounterFamily() *dto.MetricFamily {
	return getMetricFamily(COUNTER_METRIC)
}

func newCounterVecFamily() *dto.MetricFamily {
	return getMetricFamily(COUNTER_VEC_METRIC)
}

func newSummaryFamily() *dto.MetricFamily {
	return getMetricFamily(SUMMARY_METRIC)
}

func TestCreateMetricAggregate(t *testing.T) {
	metricName := "radius_go_gc_duration_seconds"
	metricFamily := newSummaryFamily()
	_, err := collection.CreateMetricAggregate(metricName, metricFamily)
	assert.EqualError(t, err, "Not building MetricAggregate for metric type SUMMARY")

	metricName = "radius_go_goroutines"
	metricFamily = newGaugeFamily()
	metricAggregate, err := collection.CreateMetricAggregate(metricName, metricFamily)
	assert.NoError(t, err)
	collector := metricAggregate.GetCollector()
	_, ok := collector.(prometheus.Gauge)
	assert.True(t, ok)

	metricName = "radius_go_info"
	metricFamily = newGaugeVecFamily()
	metricAggregate, err = collection.CreateMetricAggregate(metricName, metricFamily)
	assert.NoError(t, err)
	collector = metricAggregate.GetCollector()
	_, ok = collector.(prometheus.GaugeVec)
	assert.True(t, ok)

	metricName = "radius_go_memstats_frees_total"
	metricFamily = newCounterFamily()
	metricAggregate, err = collection.CreateMetricAggregate(metricName, metricFamily)
	assert.NoError(t, err)
	collector = metricAggregate.GetCollector()
	_, ok = collector.(prometheus.Gauge)
	assert.True(t, ok)

	metricName = "listener_init_success_count"
	metricFamily = newCounterVecFamily()
	metricAggregate, err = collection.CreateMetricAggregate(metricName, metricFamily)
	assert.NoError(t, err)
	collector = metricAggregate.GetCollector()
	_, ok = collector.(prometheus.GaugeVec)
	assert.True(t, ok)
}

func TestGaugeMetricAggregate_Update(t *testing.T) {
	metricName := "radius_go_goroutines"
	metricFamily := newGaugeFamily()
	metricAggregate, err := collection.CreateMetricAggregate(metricName, metricFamily)
	assert.NoError(t, err)

	metricFamily = getMetricFamily(GAUGE_METRIC_UPDATED)
	metricAggregate.Update(metricFamily)
	collector := metricAggregate.GetCollector()
	gauge, ok := collector.(prometheus.Gauge)
	assert.True(t, ok)

	prometheus.MustRegister(gauge)

	metricFamily, ok = getMetricFamilyFromDefaultGatherer(metricName)
	assert.True(t, ok)
	metricArr := metricFamily.GetMetric()
	assert.Len(t, metricArr, 1)
	for _, metricSample := range metricArr {
		inp := metricSample.GetGauge()
		val := inp.GetValue()
		expected := float64(3)
		assert.Equal(t, expected, val)
	}
}

func TestCounterMetricAggregate_Update(t *testing.T) {
	metricName := "radius_go_memstats_frees_total"
	metricFamily := newCounterFamily()
	metricAggregate, err := collection.CreateMetricAggregate(metricName, metricFamily)
	assert.NoError(t, err)

	metricFamily = getMetricFamily(COUNTER_METRIC_UPDATED)
	metricAggregate.Update(metricFamily)
	collector := metricAggregate.GetCollector()
	gauge, ok := collector.(prometheus.Gauge)
	assert.True(t, ok)

	prometheus.MustRegister(gauge)

	metricFamily, ok = getMetricFamilyFromDefaultGatherer(metricName)
	assert.True(t, ok)
	metricArr := metricFamily.GetMetric()
	assert.Len(t, metricArr, 1)
	for _, metricSample := range metricArr {
		inp := metricSample.GetGauge()
		val := inp.GetValue()
		expected := float64(1)
		assert.Equal(t, expected, val)
	}

}

func getMetricFamilyFromDefaultGatherer(metricName string) (*dto.MetricFamily, bool) {
	metricFamilyArr, _ := prometheus.DefaultGatherer.Gather()

	for _, metricFamily := range metricFamilyArr {
		name := metricFamily.GetName()
		if metricName == name {
			return metricFamily, true
		}
	}
	return nil, false
}
