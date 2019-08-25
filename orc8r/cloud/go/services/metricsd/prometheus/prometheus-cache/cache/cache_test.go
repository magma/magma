/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package cache

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
)

const (
	timestamp = 1559953047

	sampleReceiveString = `
# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total counter
http_requests_total{method="post",code="200"} 1027 1395066363410
http_requests_total{method="post",code="400"}    3 1395066363021
http_requests_total{method="post",code="400"}    3 1395066363010
http_requests_total{method="post",code="400"}    3 1395066363330
http_requests_total{method="post",code="400"}    3 1395066363000
# HELP cpu_usage The total CPU usage.
# TYPE cpu_usage gauge
cpu_usage{host="A"} 1027 1395066363000
cpu_usage{host="B"}    3 1395066363100
cpu_usage{host="B"}    3 1395066363030
cpu_usage{host="B"}    3 1395066363130
cpu_usage{host="B"}    3 1395066363040
# HELP memory_usage The total memory usage.
# TYPE memory_usage gauge
memory_usage{host="A"} 5 1395066363920
memory_usage{host="A"} 5 1395066363130
memory_usage{host="A"} 5 1395066363430
memory_usage{host="A"} 5 1395066363590
`
)

var (
	testName   = "testName"
	testValue  = "testValue"
	testLabels = []*dto.LabelPair{{Name: &testName, Value: &testValue}}
)

func TestReceiveMetrics(t *testing.T) {
	cache := NewMetricCache(0, 10)
	resp, err := receiveString(cache, sampleReceiveString)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Check Internal Metrics
	assert.Equal(t, 14, int(getGaugeValue(cache.internalMetrics[internalMetricCacheSize])))
	assert.Equal(t, 0, int(getGaugeValue(cache.internalMetrics[internalMetricCacheLimit])))
}

func TestReceiveOverLimit(t *testing.T) {
	cache := NewMetricCache(1, 10)
	resp, err := receiveString(cache, sampleReceiveString)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotAcceptable, resp.Code)

	assert.Equal(t, 0, int(getGaugeValue(cache.internalMetrics[internalMetricCacheSize])))
	assert.Equal(t, 1, int(getGaugeValue(cache.internalMetrics[internalMetricCacheLimit])))
}

func TestReceiveBadMetrics(t *testing.T) {
	cache := NewMetricCache(0, 10)
	resp, _ := receiveString(cache, "bad metric string")
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func receiveString(cache *MetricCache, receiveString string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(receiveString))
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := cache.Receive(c)
	return rec, err
}

func TestScrape(t *testing.T) {
	cache := NewMetricCache(0, 10)
	_, err := receiveString(cache, sampleReceiveString)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = cache.Scrape(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// parse the output to make sure it gives valid response
	var parser expfmt.TextParser
	parsedFamilies, err := parser.TextToMetricFamilies(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(parsedFamilies))

	// make sure all metrics are returned.
	// there are 2 extra internal metrics
	sum := 0
	for _, family := range parsedFamilies {
		sum += len(family.Metric)
	}
	assert.Equal(t, 16, sum)
}

func TestScrapeBadMetrics(t *testing.T) {
	// check that Scrape handles errors
	assertWorkerPoolHandlesError(t)
}

func TestDebugEndpoint(t *testing.T) {
	cache := NewMetricCache(20, 10)
	_, err := receiveString(cache, sampleReceiveString)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/debug", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = cache.Debug(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	assert.Equal(t, 5, cache.stats.currentCountSeries)
	assert.Equal(t, 3, cache.stats.currentCountFamilies)
	assert.Equal(t, 14, cache.stats.currentCountDatapoints)
	assert.Equal(t, 3, cache.stats.lastReceiveNumFamilies)
}

func TestCacheMetrics(t *testing.T) {
	cacheSingleFamily(t, 1)
	cacheSingleFamily(t, 100)
	cacheSingleFamily(t, 10000)

	cacheMultipleFamilies(t)
	cacheMultipleSeries(t)

	assertTimestampsSortedProperly(t)
}

func cacheSingleFamily(t *testing.T, metricsInFamily int) {
	cache := NewMetricCache(0, 10)
	mf := makeFamily(dto.MetricType_GAUGE, "metricA", metricsInFamily, testLabels, timestamp)
	metrics := map[string]*dto.MetricFamily{"metricA": mf}

	cache.cacheMetrics(metrics)
	// 1 family, 1 series with multiple datapoints
	assert.Equal(t, len(cache.metricFamiliesByName), 1)
	for _, family := range cache.metricFamiliesByName {
		assert.Equal(t, 1, len(family.metrics))
		for _, metric := range family.metrics {
			assert.Equal(t, metricsInFamily, len(metric))
		}
	}
}

func cacheMultipleFamilies(t *testing.T) {
	cache := NewMetricCache(0, 10)
	mf1 := makeFamily(dto.MetricType_GAUGE, "mf1", 5, testLabels, timestamp)
	mf2 := makeFamily(dto.MetricType_GAUGE, "mf2", 10, testLabels, timestamp)
	metrics := map[string]*dto.MetricFamily{"mf1": mf1, "mf2": mf2}

	cache.cacheMetrics(metrics)
	// 2 families each with 1 series
	assert.Equal(t, len(cache.metricFamiliesByName), 2)
	for familyName, family := range cache.metricFamiliesByName {
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
	cache := NewMetricCache(0, 10)
	mf1 := makeFamily(dto.MetricType_GAUGE, "mf1", 1, testLabels, timestamp)
	mf2 := makeFamily(dto.MetricType_GAUGE, "mf1", 1, []*dto.LabelPair{}, timestamp)
	mf1Map := map[string]*dto.MetricFamily{"mf1": mf1}
	mf2Map := map[string]*dto.MetricFamily{"mf1": mf2}

	cache.cacheMetrics(mf1Map)
	cache.cacheMetrics(mf2Map)
	// 1 family with 2 unique series
	assert.Equal(t, len(cache.metricFamiliesByName), 1)
	for _, family := range cache.metricFamiliesByName {
		assert.Equal(t, 2, len(family.metrics))
	}
}

func assertTimestampsSortedProperly(t *testing.T) {
	cache := NewMetricCache(0, 10)
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
	assert.Equal(t, expectedExpositionText, cache.exposeMetrics(cache.metricFamiliesByName, 1))
}

func assertWorkerPoolHandlesError(t *testing.T) {
	cache := NewMetricCache(0, 10)
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
	blankFamilyName := "" // for error

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

	errorFamily := dto.MetricFamily{
		Name: &blankFamilyName,
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

	metrics := map[string]*dto.MetricFamily{"mf": &mf, "errorFamily": &errorFamily}
	cache.cacheMetrics(metrics)

	expectedExpositionText := `# TYPE mf1 counter
mf1 123 1
mf1 234 2
mf1 456 3
`
	assert.Equal(t, expectedExpositionText, cache.exposeMetrics(cache.metricFamiliesByName, 5))
}

func getGaugeValue(gauge prometheus.Gauge) float64 {
	var dtoMetric dto.Metric
	gauge.Write(&dtoMetric)
	return *dtoMetric.Gauge.Value
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
