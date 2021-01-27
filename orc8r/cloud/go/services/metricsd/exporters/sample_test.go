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
	metricName      string
	metricType      dto.MetricType
	metric          *dto.Metric
	context         MetricContext
	expectedSamples []Sample
}

func (c *getSamplesTestCase) RunTest(t *testing.T) {
	samples := GetSamplesForMetrics(c.metricName, c.metricType, c.metric, c.context)
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
			name:       "Pushed Metric with GatewayID",
			metricName: testMetricName,
			metricType: dto.MetricType_GAUGE,
			metric:     makeTestMetric([]*dto.LabelPair{{Name: tests.MakeStrPtr(metrics.GatewayLabelName), Value: &testGateway}}),
			context:    &PushedMetricContext{NetworkID: testNetwork},
			expectedSamples: []Sample{{
				entity:      fmt.Sprintf("%s.%s", testNetwork, testGateway),
				name:        testMetricName,
				value:       "0",
				timestampMs: 0,
				labels:      []*dto.LabelPair{},
			}},
		},
		{
			name:       "Pushed Metric with no GatewayID",
			metricName: testMetricName,
			metricType: dto.MetricType_GAUGE,
			metric:     makeTestMetric([]*dto.LabelPair{}),
			context:    &PushedMetricContext{NetworkID: testNetwork},
			expectedSamples: []Sample{{
				entity:      testNetwork,
				name:        testMetricName,
				value:       "0",
				timestampMs: 0,
				labels:      []*dto.LabelPair{},
			}},
		},
		{
			name:       "Gateway Metric",
			metricName: testMetricName,
			metricType: dto.MetricType_GAUGE,
			metric:     makeTestMetric(simpleLabels),
			context: &GatewayMetricContext{
				NetworkID: testNetwork,
				GatewayID: testGateway,
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
			name:       "Cloud Metric",
			metricName: testMetricName,
			metricType: dto.MetricType_GAUGE,
			metric:     makeTestMetric(simpleLabels),
			context:    &CloudMetricContext{CloudHost: testCloudHost},
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

func makeTestMetric(labels []*dto.LabelPair) *dto.Metric {
	val := 0.0
	return &dto.Metric{
		Label: labels,
		Gauge: &dto.Gauge{
			Value: &val,
		},
	}
}
