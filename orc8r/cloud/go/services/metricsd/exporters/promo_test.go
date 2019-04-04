/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package exporters_test

import (
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/cloud/go/services/metricsd/exporters/mocks"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	defaultConfig     = exporters.PrometheusExporterConfig{UseHostLabel: true}
	pushgatewayConfig = exporters.PrometheusExporterConfig{UseHostLabel: false}
)

func TestPromoSubmitGauge(t *testing.T) {
	testSubmitPrometheusType(t, dto.MetricType_GAUGE, defaultConfig)
	testSubmitPrometheusType(t, dto.MetricType_GAUGE, pushgatewayConfig)
}

func TestPromoSubmitCounter(t *testing.T) {
	testSubmitPrometheusType(t, dto.MetricType_COUNTER, defaultConfig)
	testSubmitPrometheusType(t, dto.MetricType_COUNTER, pushgatewayConfig)
}

func TestPromoSubmitSummary(t *testing.T) {
	testSubmitPrometheusType(t, dto.MetricType_SUMMARY, defaultConfig)
	testSubmitPrometheusType(t, dto.MetricType_SUMMARY, pushgatewayConfig)
}

func TestPromoSubmitHistogram(t *testing.T) {
	testSubmitPrometheusType(t, dto.MetricType_HISTOGRAM, defaultConfig)
	testSubmitPrometheusType(t, dto.MetricType_HISTOGRAM, pushgatewayConfig)
}

func testSubmitPrometheusType(t *testing.T, metricType dto.MetricType, config exporters.PrometheusExporterConfig) {
	registry := &mocks.Registerer{}
	exporter := exporters.NewPrometheusExporter(config)
	exporter.(*exporters.PrometheusExporter).Registry = registry

	serviceLabelPair := dto.LabelPair{
		Name:  makeStringPointer(exporters.SERVICE_LABEL_NAME),
		Value: makeStringPointer("testService"),
	}
	family := makeTestMetricFamily(metricType, 1, []*dto.LabelPair{&serviceLabelPair})
	context := exporters.MetricsContext{
		NetworkID:         "nID",
		GatewayID:         "gID",
		OriginatingEntity: "entity",
		DecodedName:       "testName",
	}

	registry.On("Register", mock.Anything).Return(nil)
	// Registering a new metric should not throw error
	err := exporter.Submit([]exporters.MetricAndContext{{Family: family, Context: context}})
	assert.NoError(t, err)
	registry.AssertCalled(t, "Register", mock.Anything)

	var numExpectedCalls int
	switch metricType {
	case dto.MetricType_GAUGE:
		numExpectedCalls = 1
	case dto.MetricType_COUNTER:
		numExpectedCalls = 1
	case dto.MetricType_SUMMARY:
		// 2 calls: 2 for sum/count
		numExpectedCalls = 2
	case dto.MetricType_HISTOGRAM:
		// 8 calls: 2 for sum/count, 2 each for 3 buckets
		numExpectedCalls = 8
	}
	registry.AssertNumberOfCalls(t, "Register", numExpectedCalls)

	// Clear method calls from registry to restart counting
	registry.Calls = []mock.Call{}

	// Updating existing metric should not throw error
	err = exporter.Submit([]exporters.MetricAndContext{{Family: family, Context: context}})
	assert.NoError(t, err)
	// Register() should not have been called, should have updated instead
	registry.AssertNotCalled(t, "Register", mock.Anything)
}

func TestSanitizePrometheusNames(t *testing.T) {
	goodName := "metric1_submetric"
	badName1 := "metric1.submetric"
	badName2 := "metric1&submetric"
	badName3 := "metric1-submetric"

	sanitizedName := exporters.SanitizePrometheusNames(goodName)
	assert.Equal(t, goodName, sanitizedName)

	sanitizedName = exporters.SanitizePrometheusNames(badName1)
	assert.Equal(t, goodName, sanitizedName)

	sanitizedName = exporters.SanitizePrometheusNames(badName2)
	assert.Equal(t, goodName, sanitizedName)

	sanitizedName = exporters.SanitizePrometheusNames(badName3)
	assert.Equal(t, goodName, sanitizedName)
}
