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

package servicers

import (
	"fmt"

	prometheus_models "github.com/prometheus/client_model/go"
)

const (
	bucketPostfix = "_bucket"
	countPostfix  = "_count"
	sumPostfix    = "_sum"

	histogramBucketLabelName = "le"
	summaryQuantileLabelName = "quantile"
)

var (
	gaugeType = prometheus_models.MetricType_GAUGE
)

func convertFamilyToGauges(baseFamily *prometheus_models.MetricFamily) []*prometheus_models.MetricFamily {
	gaugeFamilies := make([]*prometheus_models.MetricFamily, 0)
	switch *baseFamily.Type {
	case prometheus_models.MetricType_GAUGE:
		return append(gaugeFamilies, baseFamily)
	case prometheus_models.MetricType_COUNTER:
		return append(gaugeFamilies, counterToGauge(baseFamily))
	case prometheus_models.MetricType_HISTOGRAM:
		return append(gaugeFamilies, histogramToGauges(baseFamily)...)
	case prometheus_models.MetricType_SUMMARY:
		return append(gaugeFamilies, summaryToGauges(baseFamily)...)
	default: // Untyped
		return append(gaugeFamilies, untypedToGauge(baseFamily))
	}
}

// counterToGauge takes a counter and converts it to a gauge with the
// same value
func counterToGauge(family *prometheus_models.MetricFamily) *prometheus_models.MetricFamily {
	counterFamily := prometheus_models.MetricFamily{
		Name: makeStringPointer(family.GetName()),
		Type: &gaugeType,
	}
	for _, metric := range family.Metric {
		if metric.Counter == nil {
			continue
		}
		counterValue := *metric.Counter.Value
		counterMetric := prometheus_models.Metric{
			Label: metric.Label,
			Gauge: &prometheus_models.Gauge{
				Value: &counterValue,
			},
		}
		counterFamily.Metric = append(counterFamily.Metric, &counterMetric)
	}
	return &counterFamily
}

// histogramToGauges converts a histogram into 3 families of gauges, one for the
// buckets, sum, and count each.
func histogramToGauges(family *prometheus_models.MetricFamily) []*prometheus_models.MetricFamily {
	baseName := family.GetName()
	bucketFamily := prometheus_models.MetricFamily{
		Name: makeStringPointer(baseName + bucketPostfix),
		Type: &gaugeType,
	}
	sumFamily := prometheus_models.MetricFamily{
		Name: makeStringPointer(baseName + sumPostfix),
		Type: &gaugeType,
	}
	countFamily := prometheus_models.MetricFamily{
		Name: makeStringPointer(baseName + countPostfix),
		Type: &gaugeType,
	}

	for _, metric := range family.Metric {
		if metric.Histogram == nil {
			continue
		}
		sumValue := *metric.Histogram.SampleSum
		sumMetric := prometheus_models.Metric{
			Label: metric.Label,
			Gauge: &prometheus_models.Gauge{
				Value: &sumValue,
			},
		}
		sumFamily.Metric = append(sumFamily.Metric, &sumMetric)

		countValue := float64(*metric.Histogram.SampleCount)
		countMetric := prometheus_models.Metric{
			Label: metric.Label,
			Gauge: &prometheus_models.Gauge{
				Value: &countValue,
			},
		}
		countFamily.Metric = append(countFamily.Metric, &countMetric)

		baseLabels := make([]*prometheus_models.LabelPair, len(metric.GetLabel()))
		for idx, l := range metric.GetLabel() {
			copyL := *l
			baseLabels[idx] = &copyL
		}

		for _, bucket := range metric.Histogram.Bucket {
			bucketValue := float64(*bucket.CumulativeCount)
			bucketMetric := prometheus_models.Metric{
				Label: append(baseLabels, &prometheus_models.LabelPair{
					Name:  makeStringPointer(histogramBucketLabelName),
					Value: makeStringPointer(fmt.Sprintf("%g", bucket.GetUpperBound())),
				}),
				Gauge: &prometheus_models.Gauge{
					Value: &bucketValue,
				},
			}
			bucketFamily.Metric = append(bucketFamily.Metric, &bucketMetric)
		}
	}
	return []*prometheus_models.MetricFamily{&bucketFamily, &sumFamily, &countFamily}
}

// summaryToGauges converts a summary to 3 gauge families, one for the quantiles,
// sum, and count each
func summaryToGauges(family *prometheus_models.MetricFamily) []*prometheus_models.MetricFamily {
	baseName := family.GetName()
	quantFamily := &prometheus_models.MetricFamily{
		Name: makeStringPointer(baseName),
		Type: &gaugeType,
	}
	sumFamily := &prometheus_models.MetricFamily{
		Name: makeStringPointer(baseName + sumPostfix),
		Type: &gaugeType,
	}
	countFamily := &prometheus_models.MetricFamily{
		Name: makeStringPointer(baseName + countPostfix),
		Type: &gaugeType,
	}

	for _, metric := range family.Metric {
		if metric.Summary == nil {
			continue
		}
		sumValue := *metric.Summary.SampleSum
		sumMetric := prometheus_models.Metric{
			Label: metric.Label,
			Gauge: &prometheus_models.Gauge{
				Value: &sumValue,
			},
		}
		sumFamily.Metric = append(sumFamily.Metric, &sumMetric)

		countValue := float64(*metric.Summary.SampleCount)
		countMetric := prometheus_models.Metric{
			Label: metric.Label,
			Gauge: &prometheus_models.Gauge{
				Value: &countValue,
			},
		}
		countFamily.Metric = append(countFamily.Metric, &countMetric)

		baseLabels := make([]*prometheus_models.LabelPair, len(metric.GetLabel()))
		for idx, l := range metric.GetLabel() {
			copyL := *l
			baseLabels[idx] = &copyL
		}

		for _, quant := range metric.Summary.GetQuantile() {
			quantValue := quant.GetValue()
			quantMetric := &prometheus_models.Metric{
				Label: append(baseLabels, &prometheus_models.LabelPair{
					Name:  makeStringPointer(summaryQuantileLabelName),
					Value: makeStringPointer(fmt.Sprintf("%g", quant.GetQuantile())),
				}),
				Gauge: &prometheus_models.Gauge{
					Value: &quantValue,
				},
			}
			quantFamily.Metric = append(quantFamily.Metric, quantMetric)
		}
	}
	return []*prometheus_models.MetricFamily{quantFamily, sumFamily, countFamily}
}

// untypedToGauge takes an untyped metric and converts it to a gauge with the
// same value
func untypedToGauge(family *prometheus_models.MetricFamily) *prometheus_models.MetricFamily {
	untypedFamily := prometheus_models.MetricFamily{
		Name: makeStringPointer(family.GetName()),
		Type: &gaugeType,
	}
	for _, metric := range family.Metric {
		if metric.Untyped == nil {
			continue
		}
		untypedValue := *metric.Untyped.Value
		untypedMetric := prometheus_models.Metric{
			Label: metric.Label,
			Gauge: &prometheus_models.Gauge{
				Value: &untypedValue,
			},
		}
		untypedFamily.Metric = append(untypedFamily.Metric, &untypedMetric)
	}
	return &untypedFamily
}
