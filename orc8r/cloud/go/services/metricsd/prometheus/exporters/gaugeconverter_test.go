/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package exporters

import (
	"testing"

	tests "magma/orc8r/cloud/go/services/metricsd/test_common"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestCounterToGauge(t *testing.T) {
	originalFamily := tests.MakeTestMetricFamily(dto.MetricType_COUNTER, 2, sampleLabels)
	convertedGauge := counterToGauge(originalFamily)
	assert.Equal(t, dto.MetricType_GAUGE, *convertedGauge.Type)
	assert.Equal(t, tests.CounterMetricName, convertedGauge.GetName())
}

func TestHistogramToGauge(t *testing.T) {
	originalFamily := tests.MakeTestMetricFamily(dto.MetricType_HISTOGRAM, 1, sampleLabels)
	convertedFams := histogramToGauges(originalFamily)
	assert.Equal(t, 3, len(convertedFams))
	for _, family := range convertedFams {
		assert.Equal(t, dto.MetricType_GAUGE, *family.Type)
		name := family.GetName()
		for _, metric := range family.Metric {
			if name == (tests.HistogramMetricName + bucketPostfix) {
				assert.True(t, tests.HasLabelName(metric.Label, histogramBucketLabelName))
			} else if name == (tests.HistogramMetricName + sumPostfix) {
				assert.False(t, tests.HasLabelName(metric.Label, histogramBucketLabelName))
			} else if name == (tests.HistogramMetricName + countPostfix) {
				assert.False(t, tests.HasLabelName(metric.Label, histogramBucketLabelName))
			} else {
				// Unexpected family name
				t.Fail()
			}
		}
	}
}

func TestHistogramToGaugeValues(t *testing.T) {
	// Expected counts in buckets:
	// 1: 3
	// 5: 5
	// 10: 9
	// Count: 15
	// Sum: 81.5
	observations := []float64{0.5, 0.5, 0.5, 1.5, 1.5, 5.5, 5.5, 5.5, 5.5, 11, 11, 11, 11, 11, 11}
	origMetric := tests.MakePromoHistogram([]float64{1, 5, 10}, observations)

	expectedSum := 0.0
	for _, obs := range observations {
		expectedSum += obs
	}

	metricType := dto.MetricType_HISTOGRAM
	famName := "hist"
	origFam := &dto.MetricFamily{
		Name:   &famName,
		Help:   makeStringPointer("testFamilyHelp"),
		Type:   &metricType,
		Metric: []*dto.Metric{&origMetric},
	}

	convertedFams := histogramToGauges(origFam)
	assert.Equal(t, 3, len(convertedFams))

	for _, family := range convertedFams {
		name := family.GetName()
		for _, metric := range family.Metric {
			if name == (famName + bucketPostfix) {
				if tests.HasLabel(metric.Label, histogramBucketLabelName, "1") {
					assert.Equal(t, 3.0, metric.Gauge.GetValue())
				} else if tests.HasLabel(metric.Label, histogramBucketLabelName, "5") {
					assert.Equal(t, 5.0, metric.Gauge.GetValue())
				} else if tests.HasLabel(metric.Label, histogramBucketLabelName, "10") {
					assert.Equal(t, 9.0, metric.Gauge.GetValue())
				}
			} else if name == (famName + sumPostfix) {
				assert.Equal(t, expectedSum, metric.Gauge.GetValue())
			} else if name == (famName + countPostfix) {
				assert.Equal(t, float64(len(observations)), metric.Gauge.GetValue())
			} else {
				// Unexpected family name
				t.Fail()
			}
		}
	}
}

func TestSummaryToGauge(t *testing.T) {
	originalFamily := tests.MakeTestMetricFamily(dto.MetricType_SUMMARY, 1, sampleLabels)
	convertedFams := summaryToGauges(originalFamily)
	assert.Equal(t, 3, len(convertedFams))
	for _, family := range convertedFams {
		assert.Equal(t, dto.MetricType_GAUGE, *family.Type)
		name := family.GetName()
		for _, metric := range family.Metric {
			if name == tests.SummaryMetricName {
				assert.True(t, tests.HasLabelName(metric.Label, summaryQuantileLabelName))
			} else if name == (tests.SummaryMetricName + sumPostfix) {
				assert.False(t, tests.HasLabelName(metric.Label, summaryQuantileLabelName))
			} else if name == (tests.SummaryMetricName + countPostfix) {
				assert.False(t, tests.HasLabelName(metric.Label, summaryQuantileLabelName))
			} else {
				// Unexpected family name
				t.Fail()
			}
		}

	}
}
