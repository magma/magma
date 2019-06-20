/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package exporters_test

import (
	"testing"

	mxd_exp "magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/cloud/go/services/metricsd/graphite/exporters"
	"magma/orc8r/cloud/go/services/metricsd/test_common"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestSubmitGraphiteGauge(t *testing.T) {
	testSubmitGraphiteType(t, dto.MetricType_GAUGE)
}

func TestSubmitGraphiteCounter(t *testing.T) {
	testSubmitGraphiteType(t, dto.MetricType_COUNTER)
}

func TestSubmitGraphiteSummary(t *testing.T) {
	testSubmitGraphiteType(t, dto.MetricType_SUMMARY)
}

func TestSubmitGraphiteHistogram(t *testing.T) {
	testSubmitGraphiteType(t, dto.MetricType_HISTOGRAM)
}

func testSubmitGraphiteType(t *testing.T, metricType dto.MetricType) {
	exporter := exporters.NewGraphiteExporter([]exporters.Address{{Host: "", Port: 0}})

	serviceLabelPair := dto.LabelPair{
		Name:  test_common.MakeStringPointer(mxd_exp.SERVICE_LABEL_NAME),
		Value: test_common.MakeStringPointer("testService"),
	}
	family := test_common.MakeTestMetricFamily(metricType, 1, []*dto.LabelPair{&serviceLabelPair})
	context := mxd_exp.MetricsContext{
		NetworkID:         "nID",
		GatewayID:         "gID",
		OriginatingEntity: "entity",
		DecodedName:       "testName",
	}

	// Registering a metric should not throw an error
	//err := exporter.Submit(family, context)
	err := exporter.Submit([]mxd_exp.MetricAndContext{{Family: family, Context: context}})
	assert.NoError(t, err)

	// Updating existing metric should not throw error
	err = exporter.Submit([]mxd_exp.MetricAndContext{{Family: family, Context: context}})
	assert.NoError(t, err)
}
