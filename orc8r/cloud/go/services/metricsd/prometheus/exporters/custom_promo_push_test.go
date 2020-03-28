/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package exporters

import (
	"regexp"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/exporters"
	tests "magma/orc8r/cloud/go/services/metricsd/test_common"
	"magma/orc8r/lib/go/metrics"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

var (
	sampleNetworkID  = "sampleNetwork"
	sampleGatewayID  = "sampleGateway"
	sampleMetricName = "metric_A"
	sampleLabels     = []*dto.LabelPair{
		{Name: tests.MakeStringPointer(metrics.NetworkLabelName), Value: tests.MakeStringPointer(sampleNetworkID)},
		{Name: tests.MakeStringPointer("testLabel"), Value: tests.MakeStringPointer("testValue")},
	}

	sampleGatewayContext = exporters.MetricsContext{
		MetricName: sampleMetricName,
		AdditionalContext: &exporters.GatewayMetricContext{
			NetworkID: sampleNetworkID,
			GatewayID: sampleGatewayID,
		},
	}
)

func TestCustomPushExporter_Submit(t *testing.T) {
	testSubmitCounter(t)
	testSubmitGauge(t)
	testSubmitHistogram(t)
	testSubmitSummary(t)
	testSubmitUntyped(t)

	testSubmitInvalidMetrics(t)
	testSubmitInvalidLabel(t)
	testSubmitInvalidName(t)
}

func TestNewCustomPushExporter(t *testing.T) {
	addrs := []string{"http://prometheus-cache:9091", "prometheus-cache:9091", "https://prometheus-cache:9091"}
	exp := NewCustomPushExporter(addrs).(*CustomPushExporter)
	protocolMatch := regexp.MustCompile("(http|https)://")
	for _, addr := range exp.pushAddresses {
		assert.True(t, protocolMatch.MatchString(addr))
	}
}

func testSubmitGauge(t *testing.T) {
	exp := makeTestCustomPushExporter()
	err := submitNewMetric(&exp, dto.MetricType_GAUGE, sampleGatewayContext)
	assert.NoError(t, err)
	assert.Equal(t, 1, totalMetricCount(&exp))

	err = submitNewMetric(&exp, dto.MetricType_GAUGE, sampleGatewayContext)
	assert.NoError(t, err)
	assert.Equal(t, 2, totalMetricCount(&exp))

	assert.Equal(t, len(exp.familiesByName), 1)
	for _, fam := range exp.familiesByName {
		assert.Equal(t, dto.MetricType_GAUGE, *fam.Type)
		for _, metric := range fam.Metric {
			assert.True(t, tests.HasLabel(metric.Label, "testLabel", "testValue"))
			assert.True(t, tests.HasLabel(metric.Label, metrics.NetworkLabelName, sampleNetworkID))
		}
	}
}

func testSubmitCounter(t *testing.T) {
	exp := makeTestCustomPushExporter()
	err := submitNewMetric(&exp, dto.MetricType_COUNTER, sampleGatewayContext)
	assert.NoError(t, err)
	assert.Equal(t, 1, totalMetricCount(&exp))

	err = submitNewMetric(&exp, dto.MetricType_COUNTER, sampleGatewayContext)
	assert.NoError(t, err)
	assert.Equal(t, 2, totalMetricCount(&exp))

	assert.Equal(t, len(exp.familiesByName), 1)
	for _, fam := range exp.familiesByName {
		assert.Equal(t, dto.MetricType_GAUGE, *fam.Type)
		for _, metric := range fam.Metric {
			assert.True(t, tests.HasLabel(metric.Label, "testLabel", "testValue"))
		}
	}
}

func testSubmitHistogram(t *testing.T) {
	exp := makeTestCustomPushExporter()
	err := submitNewMetric(&exp, dto.MetricType_HISTOGRAM, sampleGatewayContext)
	assert.NoError(t, err)
	assert.Equal(t, 5, totalMetricCount(&exp))

	err = submitNewMetric(&exp, dto.MetricType_HISTOGRAM, sampleGatewayContext)
	assert.NoError(t, err)
	assert.Equal(t, 10, totalMetricCount(&exp))

	assert.Equal(t, len(exp.familiesByName), 3)
	for name, fam := range exp.familiesByName {
		assert.Equal(t, dto.MetricType_GAUGE, *fam.Type)
		for _, metric := range fam.Metric {
			assert.True(t, tests.HasLabel(metric.Label, "testLabel", "testValue"))
			if strings.HasSuffix(name, bucketPostfix) {
				assert.True(t, tests.HasLabelName(metric.Label, histogramBucketLabelName))
			}
		}
	}
}

func testSubmitSummary(t *testing.T) {
	exp := makeTestCustomPushExporter()
	err := submitNewMetric(&exp, dto.MetricType_SUMMARY, sampleGatewayContext)
	assert.NoError(t, err)
	assert.Equal(t, 3, totalMetricCount(&exp))

	err = submitNewMetric(&exp, dto.MetricType_SUMMARY, sampleGatewayContext)
	assert.NoError(t, err)
	assert.Equal(t, 6, totalMetricCount(&exp))

	assert.Equal(t, len(exp.familiesByName), 3)
	for name, fam := range exp.familiesByName {
		assert.Equal(t, dto.MetricType_GAUGE, *fam.Type)
		for _, metric := range fam.Metric {
			assert.True(t, tests.HasLabel(metric.Label, "testLabel", "testValue"))
			if name == sampleMetricName {
				assert.True(t, tests.HasLabelName(metric.Label, summaryQuantileLabelName))
			}
		}
	}
}

func testSubmitUntyped(t *testing.T) {
	exp := makeTestCustomPushExporter()
	err := submitNewMetric(&exp, dto.MetricType_UNTYPED, sampleGatewayContext)
	assert.NoError(t, err)
	assert.Equal(t, 1, totalMetricCount(&exp))

	err = submitNewMetric(&exp, dto.MetricType_UNTYPED, sampleGatewayContext)
	assert.NoError(t, err)
	assert.Equal(t, 2, totalMetricCount(&exp))

	assert.Equal(t, len(exp.familiesByName), 1)
	for _, fam := range exp.familiesByName {
		assert.Equal(t, dto.MetricType_GAUGE, *fam.Type)
		for _, metric := range fam.Metric {
			assert.True(t, tests.HasLabel(metric.Label, "testLabel", "testValue"))
		}
	}

}

func testSubmitInvalidMetrics(t *testing.T) {
	// Submitting a metric family with 0 metrics should not register the family
	exp := makeTestCustomPushExporter()
	noMetricFamily := tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 0, sampleLabels)
	mc := exporters.MetricAndContext{
		Family:  noMetricFamily,
		Context: sampleGatewayContext,
	}
	metrics := []exporters.MetricAndContext{mc}

	err := exp.Submit(metrics)
	assert.NoError(t, err)
	assert.Equal(t, len(exp.familiesByName), 0)
}

func testSubmitInvalidName(t *testing.T) {
	// Submitting a metric with an invalid name should submit a renamed metric
	testInvalidName(t, "invalid metric name", "invalid_metric_name")
	testInvalidName(t, "0starts_with_number", "_0starts_with_number")
	testInvalidName(t, "bad?-/$chars", "bad____chars")
}

func testInvalidName(t *testing.T, inputName, expectedName string) {
	exp := makeTestCustomPushExporter()
	mf := tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, sampleLabels)

	mc := exporters.MetricAndContext{
		Family: mf,
		Context: exporters.MetricsContext{
			MetricName: inputName,
		},
	}
	metrics := []exporters.MetricAndContext{mc}

	err := exp.Submit(metrics)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(exp.familiesByName))
	for name := range exp.familiesByName {
		assert.Equal(t, expectedName, name)
	}
}

func testSubmitInvalidLabel(t *testing.T) {
	// Submitting a metric with invalid labelnames should not include that metric
	exp := makeTestCustomPushExporter()
	mf := tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 5, sampleLabels)
	extraMetric := tests.MakePromoGauge(10)
	mf.Metric[2] = &extraMetric
	mf.Metric[2].Label = append(mf.Metric[2].Label, &dto.LabelPair{Name: makeStringPointer("1"), Value: makeStringPointer("badLabelName")})

	mc := exporters.MetricAndContext{
		Family:  mf,
		Context: sampleGatewayContext,
	}
	metrics := []exporters.MetricAndContext{mc}

	err := exp.Submit(metrics)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(exp.familiesByName))
	for _, fam := range exp.familiesByName {
		assert.Equal(t, 4, len(fam.Metric))
	}

	// If all metrics are invalid, the family should not be submitted
	exp = makeTestCustomPushExporter()
	mf = tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, sampleLabels)
	badMetric := tests.MakePromoGauge(10)
	mf.Metric[0] = &badMetric
	mf.Metric[0].Label = append(mf.Metric[0].Label, &dto.LabelPair{Name: makeStringPointer("1"), Value: makeStringPointer("badLabelName")})

	mc = exporters.MetricAndContext{
		Family:  mf,
		Context: sampleGatewayContext,
	}
	metrics = []exporters.MetricAndContext{mc}

	err = exp.Submit(metrics)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(exp.familiesByName))
}

func totalMetricCount(exp *CustomPushExporter) int {
	total := 0
	for _, fam := range exp.familiesByName {
		total += len(fam.Metric)
	}
	return total
}

func submitNewMetric(exp *CustomPushExporter, mtype dto.MetricType, ctx exporters.MetricsContext) error {
	mc := exporters.MetricAndContext{
		Family:  tests.MakeTestMetricFamily(mtype, 1, sampleLabels),
		Context: ctx,
	}
	metrics := []exporters.MetricAndContext{mc}
	return exp.Submit(metrics)
}

func makeTestCustomPushExporter() CustomPushExporter {
	return CustomPushExporter{
		familiesByName: make(map[string]*dto.MetricFamily),
		exportInterval: pushInterval,
		pushAddresses:  []string{""},
	}
}
