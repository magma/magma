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
	"magma/orc8r/lib/go/protos"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestGetSamplesForMetrics(t *testing.T) {
	for _, testCase := range cases {
		t.Run(testCase.name, testCase.RunTest)
	}
}

type getSamplesTestCase struct {
	name             string
	metricAndContext protos.MetricAndContext
	expectedSamples  []Sample
}

func (c *getSamplesTestCase) RunTest(t *testing.T) {
	samples := GetSamplesForMetrics(&c.metricAndContext, c.metricAndContext.Family.Metric[0])
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
			name: "Pushed Metric with GatewayID",
			metricAndContext: protos.MetricAndContext{
				Family: tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, []*dto.LabelPair{{Name: tests.MakeStringPointer(metrics.GatewayLabelName), Value: &testGateway}}),
				Context: &protos.MetricContext{
					MetricName: testMetricName,
					MetricOriginContext: &protos.MetricContext_PushedMetric{
						PushedMetric: &protos.PushedMetricContext{
							NetworkId: testNetwork,
						},
					},
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
			name: "Pushed Metric with no GatewayID",
			metricAndContext: protos.MetricAndContext{
				Family: tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, []*dto.LabelPair{}),
				Context: &protos.MetricContext{
					MetricName: testMetricName,
					MetricOriginContext: &protos.MetricContext_PushedMetric{
						PushedMetric: &protos.PushedMetricContext{
							NetworkId: testNetwork,
						},
					},
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
			name: "Gateway Metric",
			metricAndContext: protos.MetricAndContext{
				Family: tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, simpleLabels),
				Context: &protos.MetricContext{
					MetricName: testMetricName,
					MetricOriginContext: &protos.MetricContext_GatewayMetric{
						GatewayMetric: &protos.GatewayMetricContext{
							NetworkId: testNetwork,
							GatewayId: testGateway,
						},
					},
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
			name: "Cloud Metric",
			metricAndContext: protos.MetricAndContext{
				Family: tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, simpleLabels),
				Context: &protos.MetricContext{
					MetricName: testMetricName,
					MetricOriginContext: &protos.MetricContext_CloudMetric{
						CloudMetric: &protos.CloudMetricContext{
							CloudHost: testCloudHost,
						},
					},
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
