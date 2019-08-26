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
				assert.True(t, hasLabelName(metric.Label, histogramBucketLabelName))
			} else if name == (tests.HistogramMetricName + sumPostfix) {
				assert.False(t, hasLabelName(metric.Label, histogramBucketLabelName))
			} else if name == (tests.HistogramMetricName + countPostfix) {
				assert.False(t, hasLabelName(metric.Label, histogramBucketLabelName))
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
				assert.True(t, hasLabelName(metric.Label, summaryQuantileLabelName))
			} else if name == (tests.SummaryMetricName + sumPostfix) {
				assert.False(t, hasLabelName(metric.Label, summaryQuantileLabelName))
			} else if name == (tests.SummaryMetricName + countPostfix) {
				assert.False(t, hasLabelName(metric.Label, summaryQuantileLabelName))
			} else {
				// Unexpected family name
				t.Fail()
			}
		}

	}
}
