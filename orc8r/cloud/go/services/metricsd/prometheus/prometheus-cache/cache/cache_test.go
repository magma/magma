/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package cache

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

const (
	timestamp = 1559953047
)

var (
	testName   = "testName"
	testValue  = "testValue"
	testLabels = []*dto.LabelPair{{Name: &testName, Value: &testValue}}
)

func TestCacheMetrics(t *testing.T) {
	cacheSingleFamily(t, 1)
	cacheSingleFamily(t, 100)
	cacheSingleFamily(t, 10000)

	cacheMultipleFamilies(t)
	cacheMultipleSeries(t)

	assertTimestampsSortedProperly(t)
}

func cacheSingleFamily(t *testing.T, metricsInFamily int) {
	cache := NewMetricCache()
	mf := makeFamily(dto.MetricType_GAUGE, "metricA", metricsInFamily, testLabels, timestamp)
	metrics := map[string]*dto.MetricFamily{"metricA": mf}

	cache.cacheMetrics(metrics)
	// 1 family, 1 series with multiple datapoints
	assert.Equal(t, len(cache.familyMap), 1)
	for _, family := range cache.familyMap {
		assert.Equal(t, 1, len(family.metrics))
		for _, metric := range family.metrics {
			assert.Equal(t, metricsInFamily, len(metric))
		}
	}
}

func cacheMultipleFamilies(t *testing.T) {
	cache := NewMetricCache()
	mf1 := makeFamily(dto.MetricType_GAUGE, "mf1", 5, testLabels, timestamp)
	mf2 := makeFamily(dto.MetricType_GAUGE, "mf2", 10, testLabels, timestamp)
	metrics := map[string]*dto.MetricFamily{"mf1": mf1, "mf2": mf2}

	cache.cacheMetrics(metrics)
	// 2 families each with 1 series
	assert.Equal(t, len(cache.familyMap), 2)
	for familyName, family := range cache.familyMap {
		if strings.HasPrefix(familyName, "mf1") {
			assert.Equal(t, 1, len(family.metrics))
			for _, metric := range family.metrics {
				assert.Equal(t, 5, len(metric))
			}
		} else {
			assert.Equal(t, 1, len(family.metrics))
			for _, metric := range family.metrics {
				assert.Equal(t, 10, len(metric))
			}
		}
	}
}

func cacheMultipleSeries(t *testing.T) {
	cache := NewMetricCache()
	mf1 := makeFamily(dto.MetricType_GAUGE, "mf1", 1, testLabels, timestamp)
	mf2 := makeFamily(dto.MetricType_GAUGE, "mf1", 1, []*dto.LabelPair{}, timestamp)
	mf1Map := map[string]*dto.MetricFamily{"mf1": mf1}
	mf2Map := map[string]*dto.MetricFamily{"mf1": mf2}

	cache.cacheMetrics(mf1Map)
	cache.cacheMetrics(mf2Map)
	// 1 family with 2 unique series
	assert.Equal(t, len(cache.familyMap), 1)
	for _, family := range cache.familyMap {
		assert.Equal(t, 2, len(family.metrics))
	}
}

func assertTimestampsSortedProperly(t *testing.T) {
	cache := NewMetricCache()
	counterValues := []float64{123, 234, 456}
	counterTimes := []int64{1, 2, 3}
	counter1 := dto.Counter{
		Value: &counterValues[0],
	}
	counter2 := dto.Counter{
		Value: &counterValues[1],
	}
	counter3 := dto.Counter{
		Value: &counterValues[2],
	}
	familyName := "mf1"
	mf := dto.MetricFamily{
		Name: &familyName,
		Metric: []*dto.Metric{{
			Counter:     &counter3,
			TimestampMs: &counterTimes[2],
		},
			{
				Counter:     &counter1,
				TimestampMs: &counterTimes[0],
			},
			{
				Counter:     &counter2,
				TimestampMs: &counterTimes[1],
			},
		},
	}

	metrics := map[string]*dto.MetricFamily{"mf": &mf}
	cache.cacheMetrics(metrics)

	expectedExpositionText := `# TYPE mf1 counter
mf1 123 1
mf1 234 2
mf1 456 3
`
	assert.Equal(t, expectedExpositionText, cache.exposeMetrics(cache.familyMap))
}

func makeFamily(familyType dto.MetricType, familyName string, numMetrics int, labels []*dto.LabelPair, timestamp int64) *dto.MetricFamily {
	metrics := make([]*dto.Metric, 0)
	for i := 0; i < numMetrics; i++ {
		met := prometheus.NewGauge(prometheus.GaugeOpts{Name: familyName, Help: familyName})
		met.Set(float64(i))
		var dtoMetric dto.Metric
		met.Write(&dtoMetric)

		dtoMetric.Label = append(dtoMetric.Label, labels...)
		dtoMetric.TimestampMs = &timestamp
		metrics = append(metrics, &dtoMetric)
	}

	return &dto.MetricFamily{
		Name:   &familyName,
		Help:   &familyName,
		Type:   &familyType,
		Metric: metrics,
	}
}
