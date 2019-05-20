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

	"magma/orc8r/cloud/go/services/metricsd/exporters"
	tests "magma/orc8r/cloud/go/services/metricsd/test_common"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

var (
	sampleNetworkID  = "sampleNetwork"
	sampleGatewayID  = "sampleGateway"
	sampleHardwareID = "12345"
	sampleEntity     = "sampleNetwork.sampleGateway"
	sampleMetricName = "metric_A"
	sampleLabels     = []*dto.LabelPair{{Name: tests.MakeStringPointer("networkID"), Value: tests.MakeStringPointer(sampleNetworkID)}}

	sampleContext = exporters.MetricsContext{
		NetworkID:         sampleNetworkID,
		GatewayID:         sampleGatewayID,
		HardwareID:        sampleHardwareID,
		OriginatingEntity: sampleEntity,
		DecodedName:       sampleMetricName,
		MetricName:        sampleMetricName,
	}

	testPushAddress = ""
)

func TestCustomPushExporter_Submit(t *testing.T) {
	testSubmitType(t, dto.MetricType_GAUGE)
	testSubmitType(t, dto.MetricType_COUNTER)
	testSubmitType(t, dto.MetricType_HISTOGRAM)
	testSubmitType(t, dto.MetricType_SUMMARY)

	testSubmitInvalidMetrics(t)
}

func testSubmitType(t *testing.T, mtype dto.MetricType) {
	exp := makeTestCustomPushExporter()
	mc := exporters.MetricAndContext{
		Family:  tests.MakeTestMetricFamily(mtype, 1, sampleLabels),
		Context: sampleContext,
	}
	metrics := []exporters.MetricAndContext{mc}

	err := exp.Submit(metrics)
	assert.NoError(t, err)
	assert.Equal(t, len(exp.familiesByName), 1)
}

func testSubmitInvalidMetrics(t *testing.T) {
	// Submitting a metric family with 0 metrics should not register the family
	exp := makeTestCustomPushExporter()
	noMetricFamily := tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 0, sampleLabels)
	mc := exporters.MetricAndContext{
		Family:  noMetricFamily,
		Context: sampleContext,
	}
	metrics := []exporters.MetricAndContext{mc}

	err := exp.Submit(metrics)
	assert.NoError(t, err)
	assert.Equal(t, len(exp.familiesByName), 0)

	// Submitting a metric with differing type than its family should not register
	// that metric
	exp = makeTestCustomPushExporter()
	badMetricFamily := tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 0, sampleLabels)
	counter := tests.MakePromoCounter(0)
	badMetricFamily.Metric = append(badMetricFamily.Metric, &counter)
	mc = exporters.MetricAndContext{
		Family:  badMetricFamily,
		Context: sampleContext,
	}
	metrics = []exporters.MetricAndContext{mc}

	err = exp.Submit(metrics)
	assert.NoError(t, err)
	assert.Equal(t, len(exp.familiesByName), 1)
	assert.Equal(t, len(exp.familiesByName[sampleMetricName].Metric), 1)
}

func makeTestCustomPushExporter() CustomPushExporter {
	return CustomPushExporter{
		familiesByName: make(map[string]*dto.MetricFamily),
		exportInterval: pushInterval,
		pushAddress:    testPushAddress,
	}
}
