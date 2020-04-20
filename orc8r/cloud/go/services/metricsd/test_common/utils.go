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

package test_common

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

const (
	GaugeMetricName     = "testGauge"
	CounterMetricName   = "testCounter"
	HistogramMetricName = "testHistogram"
	SummaryMetricName   = "testSummary"
	UntypedMetricName   = "testUntyped"
)

func MakeTestMetricFamily(metricType dto.MetricType, count int, labels []*dto.LabelPair) *dto.MetricFamily {
	var testMetric dto.Metric
	var familyName string
	switch metricType {
	case dto.MetricType_COUNTER:
		testMetric = MakePromoCounter(0)
		familyName = CounterMetricName
	case dto.MetricType_SUMMARY:
		testMetric = MakePromoSummary(map[float64]float64{0.1: 0.01}, []float64{})
		familyName = SummaryMetricName
	case dto.MetricType_HISTOGRAM:
		testMetric = MakePromoHistogram([]float64{1, 5, 10}, []float64{})
		familyName = HistogramMetricName
	case dto.MetricType_UNTYPED:
		testMetric = MakePromoUntyped(0)
		familyName = UntypedMetricName
	default:
		testMetric = MakePromoGauge(0)
		familyName = GaugeMetricName
	}

	testMetric.Label = labels
	metrics := make([]*dto.Metric, 0, count)
	for i := 0; i < count; i++ {
		metrics = append(metrics, &testMetric)
	}
	return &dto.MetricFamily{
		Name:   MakeStrPtr(familyName),
		Help:   MakeStrPtr("testFamilyHelp"),
		Type:   MakeMetricTypePointer(metricType),
		Metric: metrics,
	}
}

func MakePromoGauge(value float64) dto.Metric {
	var metric dto.Metric
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{Name: GaugeMetricName, Help: "testGaugeHelp"})
	gauge.Set(value)
	_ = gauge.Write(&metric)
	return metric
}

func MakePromoCounter(value float64) dto.Metric {
	var metric dto.Metric
	counter := prometheus.NewCounter(prometheus.CounterOpts{Name: CounterMetricName, Help: "testCounterHelp"})
	counter.Add(value)
	_ = counter.Write(&metric)
	return metric
}

func MakePromoSummary(objectives map[float64]float64, observations []float64) dto.Metric {
	var metric dto.Metric
	summary := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       SummaryMetricName,
			Help:       "testSummaryHelp",
			Objectives: objectives,
		},
	)
	for _, obs := range observations {
		summary.Observe(obs)
	}
	_ = summary.Write(&metric)
	return metric
}

func MakePromoHistogram(buckets []float64, observations []float64) dto.Metric {
	var metric dto.Metric
	histogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    HistogramMetricName,
			Help:    "testHistogramHelp",
			Buckets: buckets,
		},
	)
	for _, obs := range observations {
		histogram.Observe(obs)
	}
	_ = histogram.Write(&metric)
	return metric
}

func MakePromoUntyped(value float64) dto.Metric {
	var metric dto.Metric
	untyped := prometheus.NewUntypedFunc(prometheus.UntypedOpts{Name: UntypedMetricName, Help: "testUntypedHelp"}, func() float64 { return value })
	_ = untyped.Write(&metric)
	return metric
}

func MakeStrPtr(s string) *string {
	return &s
}

func MakeMetricTypePointer(t dto.MetricType) *dto.MetricType {
	return &t
}

func HasLabelName(labels []*dto.LabelPair, name string) bool {
	for _, label := range labels {
		if label.GetName() == name {
			return true
		}
	}
	return false
}

func HasLabel(labels []*dto.LabelPair, name, value string) bool {
	for _, label := range labels {
		if label.GetName() == name {
			return label.GetValue() == value
		}
	}
	return false
}
