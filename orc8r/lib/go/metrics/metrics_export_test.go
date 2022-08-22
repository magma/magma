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

package metrics_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	prometheus_proto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/lib/go/metrics"
)

func GetMetricValue(families []*prometheus_proto.MetricFamily, familyName string) (*prometheus_proto.Metric, error) {
	for _, fam := range families {
		if familyName == *fam.Name {
			return fam.Metric[0], nil
		}
	}
	return nil, fmt.Errorf("Metric familyName %v not found\n", familyName)
}

func GetMetricValueByLabels(
	families []*prometheus_proto.MetricFamily,
	familyName string,
	labels []prometheus_proto.LabelPair,
) (*prometheus_proto.Metric, error) {
	for _, fam := range families {
		if familyName != *fam.Name {
			continue
		}
		for _, metric := range fam.Metric {
			if allLabelsPresent(*metric, labels) {
				return metric, nil
			}
		}
	}
	return nil, fmt.Errorf("No metric of %v with all given labels found\n", familyName)

}

func allLabelsPresent(metric prometheus_proto.Metric, labels []prometheus_proto.LabelPair) bool {
	metricLabels := metricLabelsAsMap(metric.Label)
	for _, requiredPair := range labels {
		val, ok := metricLabels[*requiredPair.Name]
		if !ok || val != *requiredPair.Value {
			return false
		}
	}
	return true
}

func metricLabelsAsMap(labels []*prometheus_proto.LabelPair) map[string]string {
	metricMap := make(map[string]string)
	for _, label := range labels {
		metricMap[*label.Name] = *label.Value
	}
	return metricMap
}

func TestGetMetrics(t *testing.T) {
	// setup
	labelKey := "some_words"
	labelVal1 := "blah1"
	labelVal2 := "blah2"
	testGauge1 := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "test_gauge_1",
		Help: "test gauge 1.",
	})
	testCounter1 := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_counter_1",
			Help: "test counter 1",
		},
		[]string{labelKey},
	)
	prometheus.MustRegister(testGauge1, testCounter1)
	testGauge1.Set(12.3)
	testCounter1.With(prometheus.Labels{labelKey: labelVal1}).Inc()
	testCounter1.With(prometheus.Labels{labelKey: labelVal2}).Add(2)

	// test GetMetrics
	metricFams, err := metrics.GetMetrics()
	assert.NoError(t, err)
	metric, err := GetMetricValue(metricFams, "test_gauge_1")
	assert.NoError(t, err)
	assert.Equal(t, 12.3, *metric.Gauge.Value)
	labelPair := []prometheus_proto.LabelPair{{Name: &labelKey, Value: &labelVal1}}
	metric, err = GetMetricValueByLabels(metricFams, "test_counter_1", labelPair)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(*metric.Counter.Value))
	labelPair = []prometheus_proto.LabelPair{{Name: &labelKey, Value: &labelVal2}}
	metric, err = GetMetricValueByLabels(metricFams, "test_counter_1", labelPair)
	assert.NoError(t, err)
	assert.Equal(t, 2, int(*metric.Counter.Value))

	// cleanup
	prometheus.Unregister(testGauge1)
	prometheus.Unregister(testCounter1)
}

func TestGetMetricsHistograms(t *testing.T) {
	// setup
	// example taken from
	// https://godoc.org/github.com/prometheus/client_golang/prometheus#Histogram
	testHistogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "test_histogram",
		Help:    "test histogram",
		Buckets: prometheus.LinearBuckets(20, 5, 5),
	})
	for i := 0; i < 1000; i++ {
		testHistogram.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
	}
	prometheus.MustRegister(testHistogram)

	// test GetMetrics
	metricFams, err := metrics.GetMetrics()
	assert.NoError(t, err)
	metric, err := GetMetricValue(metricFams, "test_histogram")
	assert.NoError(t, err)
	assert.Equal(t, 1000, int(*metric.Histogram.SampleCount))
	assert.Equal(t, 29969.50000000001, *metric.Histogram.SampleSum)
	assert.Equal(t, 5, len(metric.Histogram.Bucket))
	assert.Equal(t, 192, int(*metric.Histogram.Bucket[0].CumulativeCount))

	// cleanup
	prometheus.Unregister(testHistogram)
}

func TestGetMetricsSummaries(t *testing.T) {
	// setup
	// example taken from
	// https://godoc.org/github.com/prometheus/client_golang/prometheus#Summary
	testSummary := prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "test_summary",
		Help:       "test summary",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})

	for i := 0; i < 1000; i++ {
		testSummary.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
	}
	prometheus.MustRegister(testSummary)

	// test GetMetrics
	metricFams, err := metrics.GetMetrics()
	assert.NoError(t, err)
	metric, err := GetMetricValue(metricFams, "test_summary")
	assert.NoError(t, err)
	assert.Equal(t, 1000, int(*metric.Summary.SampleCount))
	assert.Equal(t, 29969.50000000001, *metric.Summary.SampleSum)
	assert.Equal(t, 3, len(metric.Summary.Quantile))
	assert.Equal(t, 31.1, *metric.Summary.Quantile[0].Value)

	// cleanup
	prometheus.Unregister(testSummary)
}
