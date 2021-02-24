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

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	tests "magma/orc8r/cloud/go/services/metricsd/test_common"
	"magma/orc8r/cloud/go/services/metricsd/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/servicers"
	"magma/orc8r/cloud/go/services/orchestrator/servicers/mocks"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/registry"

	edge_hub "github.com/facebookincubator/prometheus-edge-hub/grpc"
	prometheus_models "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	edgeControllerServiceName = "edge_controller_service"
)

type exporterTestCase struct {
	name               string
	metrics            []exporters.MetricAndContext
	assertExpectations func(t *testing.T, client *mocks.EdgeHubServer)
}

func (tc exporterTestCase) RunTest(t *testing.T) {
	// Set client return
	srv := &mocks.EdgeHubServer{}
	srv.On("Collect", mock.Anything, mock.Anything).Return(&edge_hub.Void{}, nil)

	exporter := makeExporter(t, srv)

	err := exporter.Submit(tc.metrics)
	assert.NoError(t, err)
	tc.assertExpectations(t, srv)
}

func TestGRPCExporter(t *testing.T) {
	tcs := []exporterTestCase{
		{
			name:    "submit no metrics",
			metrics: nil,
			assertExpectations: func(t *testing.T, srv *mocks.EdgeHubServer) {
				srv.AssertNotCalled(t, "Collect")
			},
		},
		{
			name:    "submit gauge",
			metrics: []exporters.MetricAndContext{{Family: tests.MakeTestMetricFamily(prometheus_models.MetricType_GAUGE, 1, []*prometheus_models.LabelPair{})}},
			assertExpectations: func(t *testing.T, srv *mocks.EdgeHubServer) {
				srv.AssertCalled(t, "Collect", mock.Anything, mock.Anything)
				srv.AssertNumberOfCalls(t, "Collect", 1)
			},
		},
		{
			name:    "submit counter",
			metrics: []exporters.MetricAndContext{{Family: tests.MakeTestMetricFamily(prometheus_models.MetricType_COUNTER, 1, []*prometheus_models.LabelPair{})}},
			assertExpectations: func(t *testing.T, srv *mocks.EdgeHubServer) {
				srv.AssertCalled(t, "Collect", mock.Anything, mock.Anything)
				srv.AssertNumberOfCalls(t, "Collect", 1)
			},
		},
		{
			name:    "submit untyped",
			metrics: []exporters.MetricAndContext{{Family: tests.MakeTestMetricFamily(prometheus_models.MetricType_UNTYPED, 1, []*prometheus_models.LabelPair{})}},
			assertExpectations: func(t *testing.T, srv *mocks.EdgeHubServer) {
				srv.AssertCalled(t, "Collect", mock.Anything, mock.Anything)
				srv.AssertNumberOfCalls(t, "Collect", 1)
			},
		},
		{
			name:    "submit histogram",
			metrics: []exporters.MetricAndContext{{Family: tests.MakeTestMetricFamily(prometheus_models.MetricType_HISTOGRAM, 1, []*prometheus_models.LabelPair{})}},
			assertExpectations: func(t *testing.T, srv *mocks.EdgeHubServer) {
				srv.AssertCalled(t, "Collect", mock.Anything, mock.Anything)
				srv.AssertNumberOfCalls(t, "Collect", 1)
			},
		},
		{
			name:    "submit summary",
			metrics: []exporters.MetricAndContext{{Family: tests.MakeTestMetricFamily(prometheus_models.MetricType_SUMMARY, 1, []*prometheus_models.LabelPair{})}},
			assertExpectations: func(t *testing.T, srv *mocks.EdgeHubServer) {
				srv.AssertCalled(t, "Collect", mock.Anything, mock.Anything)
				srv.AssertNumberOfCalls(t, "Collect", 1)
			},
		},
		{
			name:    "submit many",
			metrics: []exporters.MetricAndContext{{Family: tests.MakeTestMetricFamily(prometheus_models.MetricType_GAUGE, 10, []*prometheus_models.LabelPair{})}},
			assertExpectations: func(t *testing.T, srv *mocks.EdgeHubServer) {
				srv.AssertCalled(t, "Collect", mock.Anything, mock.Anything)
				srv.AssertNumberOfCalls(t, "Collect", 1)
			},
		},
	}

	for _, tc := range tcs {
		tc.RunTest(t)
	}
}

// makeExporter creates the following
//	- edge hub servicer (standalone service)
//	- grpc metrics exporter servicer (standalone service)
//
// The returned exporter forwards to the metrics exporter, which in turn
// forwards to the edge hub.
func makeExporter(t *testing.T, mockEdge edge_hub.MetricsControllerServer) exporters.Exporter {
	edgeSrv, lis := test_utils.NewTestService(t, orc8r.ModuleName, edgeControllerServiceName)
	edge_hub.RegisterMetricsControllerServer(edgeSrv.GrpcServer, mockEdge)
	go edgeSrv.RunTest(lis)

	edgeAddr, err := registry.GetServiceAddress(edgeControllerServiceName)
	assert.NoError(t, err)
	assert.NotEmpty(t, edgeAddr)

	srv := servicers.NewGRPCPushExporterServicer(edgeAddr)
	test_init.StartTestServiceInternal(t, srv)

	exporter := exporters.NewRemoteExporter(metricsd.ServiceName)
	return exporter
}
