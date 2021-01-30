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
	"os"
	"testing"

	"magma/fbinternal/cloud/go/services/fbinternal"
	"magma/fbinternal/cloud/go/services/fbinternal/servicers"
	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/metricsd/protos"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	exporterServicer := servicers.NewExporterServicer(
		os.Getenv("METRIC_EXPORT_URL"),
		os.Getenv("FACEBOOK_APP_ID"),
		os.Getenv("FACEBOOK_APP_SECRET"),
		"magma",
		os.Getenv("METRICS_PREFIX"),
		servicers.ODSMetricsQueueLength,
		servicers.ODSMetricsExportInterval,
	)

	specServicer := swagger.NewSpecServicer("swaggerSpec")
	StartTestServiceInternal(t, exporterServicer, specServicer)
}

func StartTestServiceInternal(t *testing.T, exporter protos.MetricsExporterServer, specServicer swagger_protos.SwaggerSpecServer) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, fbinternal.ServiceName)

	if exporter != nil {
		protos.RegisterMetricsExporterServer(srv.GrpcServer, exporter)
	}
	if specServicer != nil {
		swagger_protos.RegisterSwaggerSpecServer(srv.GrpcServer, specServicer)
	}

	go srv.RunTest(lis)
}
