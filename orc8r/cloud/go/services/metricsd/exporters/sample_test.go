/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package exporters

import (
	"fmt"
	"testing"

	tests "magma/orc8r/cloud/go/services/metricsd/test_common"
	"magma/orc8r/lib/go/metrics"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestGetSamplesForMetrics(t *testing.T) {
	for _, testCase := range cases {
		t.Run(testCase.name, testCase.RunTest)
	}
}

type getSamplesTestCase struct {
	name            string
	family          dto.MetricFamily
	context         MetricsContext
	expectedSamples []Sample
}

func (c *getSamplesTestCase) RunTest(t *testing.T) {
	samples := GetSamplesForMetrics(MetricAndContext{Family: &c.family, Context: c.context}, c.family.Metric[0])
	assert.Equal(t, c.expectedSamples, samples)
}

var (
	testGateway    = "gw1"
	testNetwork    = "nw1"
	testMetricName = "metric1"
	testCloudHost  = "hostA"
	simpleLabels   = []*dto.LabelPair{{Name: tests.MakeStringPointer("testLabel"), Value: tests.MakeStringPointer("testValue")}}

	cases = []getSamplesTestCase{
		{
			name:   "Pushed Metric with GatewayID",
			family: *tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, []*dto.LabelPair{{Name: tests.MakeStringPointer(metrics.GatewayLabelName), Value: &testGateway}}),
			context: MetricsContext{
				MetricName: testMetricName,
				AdditionalContext: &PushedMetricContext{
					NetworkID: testNetwork,
				},
			},
			expectedSamples: []Sample{{
				entity:      fmt.Sprintf("%s.%s", testNetwork, testGateway),
				name:        testMetricName,
				value:       "0",
				timestampMs: 0,
				labels:      []*dto.LabelPair{},
			}},
		},
		{
			name:   "Pushed Metric with no GatewayID",
			family: *tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, []*dto.LabelPair{}),
			context: MetricsContext{
				MetricName: testMetricName,
				AdditionalContext: &PushedMetricContext{
					NetworkID: testNetwork,
				},
			},
			expectedSamples: []Sample{{
				entity:      testNetwork,
				name:        testMetricName,
				value:       "0",
				timestampMs: 0,
				labels:      []*dto.LabelPair{},
			}},
		},
		{
			name:   "Gateway Metric",
			family: *tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, simpleLabels),
			context: MetricsContext{
				MetricName: testMetricName,
				AdditionalContext: &GatewayMetricContext{
					NetworkID: testNetwork,
					GatewayID: testGateway,
				},
			},
			expectedSamples: []Sample{{
				entity:      fmt.Sprintf("%s.%s", testNetwork, testGateway),
				name:        testMetricName,
				value:       "0",
				timestampMs: 0,
				labels:      simpleLabels,
			}},
		},
		{
			name:   "Cloud Metric",
			family: *tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, simpleLabels),
			context: MetricsContext{
				MetricName: testMetricName,
				AdditionalContext: &CloudMetricContext{
					CloudHost: testCloudHost,
				},
			},
			expectedSamples: []Sample{{
				entity:      fmt.Sprintf("cloud.%s", testCloudHost),
				name:        testMetricName,
				value:       "0",
				timestampMs: 0,
				labels:      simpleLabels,
			}},
		},
	}
)
