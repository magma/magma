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

package servicers_test

import (
	"testing"

	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	tests "magma/orc8r/cloud/go/services/metricsd/test_common"
	"magma/orc8r/cloud/go/services/metricsd/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/servicers"
	"magma/orc8r/cloud/go/services/orchestrator/servicers/mocks"

	prometheus_models "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/mock"
	assert "github.com/stretchr/testify/require"
)

type exporterTestCase struct {
	name               string
	metrics            []exporters.MetricAndContext
	assertExpectations func(t *testing.T, client *mocks.EdgeHubClient)
}

func (tc exporterTestCase) RunTest(t *testing.T) {
	// Set client return
	client := mocks.EdgeHubClient{}
	client.On("Collect", mock.Anything, mock.Anything).Return(nil, nil)

	exporter := makeTestGRPCPushExporter(t, &client)

	err := exporter.Submit(tc.metrics)
	assert.NoError(t, err)
	tc.assertExpectations(t, &client)
}

func TestGRPCExporter(t *testing.T) {
	tests := []exporterTestCase{
		{
			name:    "submit no metrics",
			metrics: nil,
			assertExpectations: func(t *testing.T, client *mocks.EdgeHubClient) {
				client.AssertNotCalled(t, "Collect")
			},
		},
		{
			name:    "submit gauge",
			metrics: []exporters.MetricAndContext{{Family: tests.MakeTestMetricFamily(prometheus_models.MetricType_GAUGE, 1, []*prometheus_models.LabelPair{})}},
			assertExpectations: func(t *testing.T, client *mocks.EdgeHubClient) {
				client.AssertCalled(t, "Collect", mock.Anything, mock.Anything)
				client.AssertNumberOfCalls(t, "Collect", 1)
			},
		},
		{
			name:    "submit counter",
			metrics: []exporters.MetricAndContext{{Family: tests.MakeTestMetricFamily(prometheus_models.MetricType_COUNTER, 1, []*prometheus_models.LabelPair{})}},
			assertExpectations: func(t *testing.T, client *mocks.EdgeHubClient) {
				client.AssertCalled(t, "Collect", mock.Anything, mock.Anything)
				client.AssertNumberOfCalls(t, "Collect", 1)
			},
		},
		{
			name:    "submit untyped",
			metrics: []exporters.MetricAndContext{{Family: tests.MakeTestMetricFamily(prometheus_models.MetricType_UNTYPED, 1, []*prometheus_models.LabelPair{})}},
			assertExpectations: func(t *testing.T, client *mocks.EdgeHubClient) {
				client.AssertCalled(t, "Collect", mock.Anything, mock.Anything)
				client.AssertNumberOfCalls(t, "Collect", 1)
			},
		},
		{
			name:    "submit histogram",
			metrics: []exporters.MetricAndContext{{Family: tests.MakeTestMetricFamily(prometheus_models.MetricType_HISTOGRAM, 1, []*prometheus_models.LabelPair{})}},
			assertExpectations: func(t *testing.T, client *mocks.EdgeHubClient) {
				client.AssertCalled(t, "Collect", mock.Anything, mock.Anything)
				client.AssertNumberOfCalls(t, "Collect", 1)
			},
		},
		{
			name:    "submit summary",
			metrics: []exporters.MetricAndContext{{Family: tests.MakeTestMetricFamily(prometheus_models.MetricType_SUMMARY, 1, []*prometheus_models.LabelPair{})}},
			assertExpectations: func(t *testing.T, client *mocks.EdgeHubClient) {
				client.AssertCalled(t, "Collect", mock.Anything, mock.Anything)
				client.AssertNumberOfCalls(t, "Collect", 1)
			},
		},
		{
			name:    "submit many",
			metrics: []exporters.MetricAndContext{{Family: tests.MakeTestMetricFamily(prometheus_models.MetricType_GAUGE, 10, []*prometheus_models.LabelPair{})}},
			assertExpectations: func(t *testing.T, client *mocks.EdgeHubClient) {
				client.AssertCalled(t, "Collect", mock.Anything, mock.Anything)
				client.AssertNumberOfCalls(t, "Collect", 1)
			},
		},
	}

	for _, tc := range tests {
		tc.RunTest(t)
	}
}

func makeTestGRPCPushExporter(t *testing.T, client servicers.EdgeHubClient) exporters.Exporter {
	srv := &servicers.GRPCPushExporterServicer{
		GrpcClient:  client,
		PushAddress: "test",
	}
	test_init.StartTestServiceInternal(t, srv)

	exporter := exporters.NewRemoteExporter(metricsd.ServiceName)
	return exporter
}
