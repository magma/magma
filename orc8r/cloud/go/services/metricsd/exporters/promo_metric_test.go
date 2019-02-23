/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package exporters_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/exporters"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

var (
	testNetworkLabels = prometheus.Labels{
		exporters.NetworkLabelInstance: "testInstance",
		exporters.NetworkLabelGateway:  "testGateway",
		exporters.NetworkLabelService:  "testService",
		exporters.NetworkLabelHost:     "testHost",
	}
)

func TestPrometheusGauge_Register(t *testing.T) {
	exporter := exporters.NewPrometheusExporter().(*exporters.PrometheusExporter)
	g := exporters.NewPrometheusGauge()

	testGaugeValue := 123.0
	gauge := makePromoGauge(testGaugeValue)
	g.Register(&gauge, "testGauge", exporter, testNetworkLabels)

	metrics, err := exporter.Registry.(*prometheus.Registry).Gather()
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, testGaugeValue, metrics[0].Metric[0].Gauge.GetValue())
}

func TestPrometheusGauge_Update(t *testing.T) {
	exporter := exporters.NewPrometheusExporter().(*exporters.PrometheusExporter)
	g := exporters.NewPrometheusGauge()
	gauge := makePromoGauge(0.0)

	g.Register(&gauge, "testGauge", exporter, testNetworkLabels)

	updatedGaugeValue := 123.0
	gauge.GetGauge().Value = &updatedGaugeValue
	g.Update(&gauge, testNetworkLabels)

	metrics, err := exporter.Registry.(*prometheus.Registry).Gather()
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, updatedGaugeValue, metrics[0].Metric[0].Gauge.GetValue())
}

func TestPrometheusCounter_Register(t *testing.T) {
	exporter := exporters.NewPrometheusExporter().(*exporters.PrometheusExporter)
	c := exporters.NewPrometheusCounter(exporter)

	testValue := 123.0
	counter := makePromoCounter(testValue)
	c.Register(&counter, "testCounter", exporter, testNetworkLabels)

	metrics, err := exporter.Registry.(*prometheus.Registry).Gather()
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, testValue, metrics[0].Metric[0].Counter.GetValue())
}

func TestPrometheusCounter_Update(t *testing.T) {
	exporter := exporters.NewPrometheusExporter().(*exporters.PrometheusExporter)
	c := exporters.NewPrometheusCounter(exporter)

	counter := makePromoCounter(0.0)
	c.Register(&counter, "testCounter", exporter, testNetworkLabels)

	updatedValue := 123.0
	counter.Counter.Value = &updatedValue
	err := c.Update(&counter, testNetworkLabels)
	assert.NoError(t, err)

	metrics, err := exporter.Registry.(*prometheus.Registry).Gather()
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, updatedValue, metrics[0].Metric[0].Counter.GetValue())

	// Test decreasing Counter Value
	decreasedValue := updatedValue - 1.0
	counter.Counter.Value = &decreasedValue
	err = c.Update(&counter, testNetworkLabels)
	assert.NoError(t, err)

	updatedMetrics, err := exporter.Registry.(*prometheus.Registry).Gather()
	assert.NoError(t, err)
	assert.Len(t, updatedMetrics, 1)
	assert.Equal(t, decreasedValue, updatedMetrics[0].Metric[0].Counter.GetValue())
}

func TestPrometheusSummary_Register(t *testing.T) {
	exporter := exporters.NewPrometheusExporter().(*exporters.PrometheusExporter)
	s := exporters.NewPrometheusSummary()

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
	observations := []float64{0.5, 0.6, 0.7}
	summary := makePromoSummary(objectives, observations)

	metricName := "testSummary"
	s.Register(&summary, metricName, exporter, testNetworkLabels)

	metrics, err := exporter.Registry.(*prometheus.Registry).Gather()
	assert.NoError(t, err)
	checkSummaryResults(t, metricName, objectives, observations, metrics)
}

func TestPrometheusSummary_Update(t *testing.T) {
	exporter := exporters.NewPrometheusExporter().(*exporters.PrometheusExporter)
	s := exporters.NewPrometheusSummary()

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
	observations := []float64{0.5, 0.6, 0.7}
	summary := makePromoSummary(objectives, observations)

	metricName := "testSummary"
	s.Register(&summary, metricName, exporter, testNetworkLabels)

	newObservations := []float64{0.8, 0.9}
	for _, obs := range newObservations {
		observations = append(observations, obs)
	}
	updatedSummary := makePromoSummary(objectives, observations)
	s.Update(&updatedSummary, testNetworkLabels)

	metrics, err := exporter.Registry.(*prometheus.Registry).Gather()
	assert.NoError(t, err)
	checkSummaryResults(t, metricName, objectives, observations, metrics)
}

func checkSummaryResults(
	t *testing.T,
	metricName string,
	objectives map[float64]float64,
	observations []float64,
	metrics []*dto.MetricFamily,
) {
	assert.Len(t, metrics, 2)

	observationSum := 0.0
	for _, o := range observations {
		observationSum += o
	}

	for _, metric := range metrics {
		if metric.GetName() == metricName+exporters.MetricPostfixCount {
			assert.Equal(t, float64(len(observations)), metric.Metric[0].Gauge.GetValue())
		} else if metric.GetName() == metricName+exporters.MetricPostfixSum {
			assert.Equal(t, observationSum, metric.Metric[0].Gauge.GetValue())
		} else {
			t.Fail()
		}
	}
}

func TestPrometheusHistogram_Register(t *testing.T) {
	exporter := exporters.NewPrometheusExporter().(*exporters.PrometheusExporter)

	metricBaseName := "testBaseName"
	h := exporters.NewPrometheusHistogram(metricBaseName)

	buckets := []float64{1.0, 5.0, 10.0}
	observations := []float64{0.5, 0.8, 2.0, 7.2, 9.2}
	histogram := makePromoHistogram(buckets, observations)

	h.Register(&histogram, metricBaseName, exporter, testNetworkLabels)

	metrics, err := exporter.Registry.(*prometheus.Registry).Gather()
	assert.NoError(t, err)

	checkHistogramResults(t, metricBaseName, buckets, observations, metrics)
}

func TestPrometheusHistogram_Update(t *testing.T) {
	exporter := exporters.NewPrometheusExporter().(*exporters.PrometheusExporter)

	metricBaseName := "testBaseName"
	h := exporters.NewPrometheusHistogram(metricBaseName)

	buckets := []float64{1.0, 5.0, 10.0}
	observations := []float64{0.5, 0.8, 2.0, 7.2, 9.2}
	histogram := makePromoHistogram(buckets, observations)

	h.Register(&histogram, metricBaseName, exporter, testNetworkLabels)

	newObservations := []float64{0.4, 2.5, 8.0}
	for _, obs := range newObservations {
		observations = append(observations, obs)
	}
	updatedHistogram := makePromoHistogram(buckets, observations)

	h.Update(&updatedHistogram, testNetworkLabels)
	metrics, err := exporter.Registry.(*prometheus.Registry).Gather()
	assert.NoError(t, err)

	checkHistogramResults(t, metricBaseName, buckets, observations, metrics)
}

func checkHistogramResults(
	t *testing.T,
	metricBaseName string,
	buckets,
	observations []float64,
	metrics []*dto.MetricFamily,
) {
	numExpectedMetrics := 2*len(buckets) + 2
	assert.Len(t, metrics, numExpectedMetrics)

	observationSum := 0.0
	bucketSums := make([]float64, len(buckets))
	for _, obs := range observations {
		observationSum += obs
		for idx, bucket := range buckets {
			if obs < bucket {
				bucketSums[idx]++
			}
		}
	}

	bucketCountPattern := regexp.MustCompile(fmt.Sprintf(".*bucket_[0-9]+%s", exporters.MetricPostfixCount))
	bucketLEPattern := regexp.MustCompile(fmt.Sprintf(".*bucket_[0-9]+%s", exporters.MetricPostfixLE))
	numberPattern := regexp.MustCompile("[0-9]+")
	for _, metric := range metrics {
		gaugeVal := metric.Metric[0].Gauge.GetValue()
		name := metric.GetName()
		if name == metricBaseName+exporters.MetricPostfixCount {
			assert.Equal(t, float64(len(observations)), gaugeVal)
		} else if name == metricBaseName+exporters.MetricPostfixSum {
			assert.Equal(t, observationSum, gaugeVal)
		} else if bucketCountPattern.MatchString(name) {
			bucketIdx, _ := strconv.Atoi(numberPattern.FindAllString(name, 1)[0])
			assert.Equal(t, bucketSums[bucketIdx], gaugeVal)
		} else if bucketLEPattern.MatchString(name) {
			bucketIdx, _ := strconv.Atoi(numberPattern.FindAllString(name, 1)[0])
			assert.Equal(t, buckets[bucketIdx], gaugeVal)
		} else {
			t.Fail()
		}
	}
}
