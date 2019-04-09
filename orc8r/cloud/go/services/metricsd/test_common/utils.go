/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package test_common

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func MakeTestMetricFamily(metricType dto.MetricType, count int, labels []*dto.LabelPair) *dto.MetricFamily {
	var testMetric dto.Metric
	switch metricType {
	case dto.MetricType_COUNTER:
		testMetric = MakePromoCounter(0)
	case dto.MetricType_SUMMARY:
		testMetric = MakePromoSummary(map[float64]float64{0.1: 0.01}, []float64{})
	case dto.MetricType_HISTOGRAM:
		testMetric = MakePromoHistogram([]float64{1, 5, 10}, []float64{})
	default:
		testMetric = MakePromoGauge(0)
	}

	testMetric.Label = labels
	metrics := make([]*dto.Metric, 0, count)
	for i := 0; i < count; i++ {
		metrics = append(metrics, &testMetric)
	}
	return &dto.MetricFamily{
		Name:   MakeStringPointer("testFamily"),
		Help:   MakeStringPointer("testFamilyHelp"),
		Type:   MakeMetricTypePointer(metricType),
		Metric: metrics,
	}
}

func MakePromoGauge(value float64) dto.Metric {
	var metric dto.Metric
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{Name: "testGauge", Help: "testGaugeHelp"})
	gauge.Set(value)
	gauge.Write(&metric)
	return metric
}

func MakePromoCounter(value float64) dto.Metric {
	var metric dto.Metric
	counter := prometheus.NewCounter(prometheus.CounterOpts{Name: "testCounter", Help: "testCounterHelp"})
	counter.Add(value)
	counter.Write(&metric)
	return metric
}

func MakePromoSummary(objectives map[float64]float64, observations []float64) dto.Metric {
	var metric dto.Metric
	summary := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "testSummary",
			Help:       "testSummaryHelp",
			Objectives: objectives,
		},
	)
	for _, obs := range observations {
		summary.Observe(obs)
	}
	summary.Write(&metric)
	return metric
}

func MakePromoHistogram(buckets []float64, observations []float64) dto.Metric {
	var metric dto.Metric
	histogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "testHistogram",
			Help:    "testHistogramHelp",
			Buckets: buckets,
		},
	)
	for _, obs := range observations {
		histogram.Observe(obs)
	}
	histogram.Write(&metric)
	return metric
}

func MakeStringPointer(s string) *string {
	return &s
}

func MakeMetricTypePointer(t dto.MetricType) *dto.MetricType {
	return &t
}
