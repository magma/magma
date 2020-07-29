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

package test_init

import (
	"context"
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/cloud/go/services/metricsd/protos"
	"magma/orc8r/cloud/go/test_utils"
)

type exporterServicer struct {
	exporter exporters.Exporter
}

// StartNewTestExporter starts a new metrics exporter service which forwards
// calls to the passed exporter.
func StartNewTestExporter(t *testing.T, exporter exporters.Exporter) {
	labels := map[string]string{
		orc8r.MetricsExporterLabel: "true",
	}
	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, "MOCK_EXPORTER_SERVICE", labels, nil)
	servicer := &exporterServicer{exporter: exporter}
	protos.RegisterMetricsExporterServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
}

func (e *exporterServicer) Submit(ctx context.Context, req *protos.SubmitMetricsRequest) (*protos.SubmitMetricsResponse, error) {
	err := e.exporter.Submit(exporters.MakeNativeMetrics(req.Metrics))
	return &protos.SubmitMetricsResponse{}, err
}
