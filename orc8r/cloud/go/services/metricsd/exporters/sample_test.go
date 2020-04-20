/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package exporters

import (
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/protos"
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
	name             string
	metricAndContext protos.ContextualizedMetric
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
	simpleLabels   = []*dto.LabelPair{{Name: tests.MakeStrPtr("testLabel"), Value: tests.MakeStrPtr("testValue")}}

	cases = []getSamplesTestCase{
		{
			name: "Pushed Metric with GatewayID",
			metricAndContext: protos.ContextualizedMetric{
				Family: tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, []*dto.LabelPair{{Name: tests.MakeStrPtr(metrics.GatewayLabelName), Value: &testGateway}}),
				Context: &protos.Context{
					MetricName: testMetricName,
					OriginContext: &protos.Context_PushedMetric{
						PushedMetric: &protos.PushedContext{
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
			metricAndContext: protos.ContextualizedMetric{
				Family: tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, []*dto.LabelPair{}),
				Context: &protos.Context{
					MetricName: testMetricName,
					OriginContext: &protos.Context_PushedMetric{
						PushedMetric: &protos.PushedContext{
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
			metricAndContext: protos.ContextualizedMetric{
				Family: tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, simpleLabels),
				Context: &protos.Context{
					MetricName: testMetricName,
					OriginContext: &protos.Context_GatewayMetric{
						GatewayMetric: &protos.GatewayContext{
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
			metricAndContext: protos.ContextualizedMetric{
				Family: tests.MakeTestMetricFamily(dto.MetricType_GAUGE, 1, simpleLabels),
				Context: &protos.Context{
					MetricName: testMetricName,
					OriginContext: &protos.Context_CloudMetric{
						CloudMetric: &protos.CloudContext{
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
