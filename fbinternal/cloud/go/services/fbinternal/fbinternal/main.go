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

package main

import (
	"os"

	"magma/fbinternal/cloud/go/fbinternal"
	fbinternal_service "magma/fbinternal/cloud/go/services/fbinternal"
	"magma/fbinternal/cloud/go/services/fbinternal/servicers"
	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/metricsd/protos"
	"magma/orc8r/lib/go/definitions"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

const (
	defaultCategoryID = "magma"
	// Set max msg received to 50MB
	DefaultMaxGRPCMsgRecvSize = 50 * 1024 * 1024
)

func main() {
	srv, err := service.NewOrchestratorService(
		fbinternal.ModuleName,
		fbinternal_service.ServiceName,
		grpc.MaxRecvMsgSize(DefaultMaxGRPCMsgRecvSize),
	)
	if err != nil {
		glog.Fatalf("Error creating orc8r service for fbinternal: %s", err)
	}
	categoryID := definitions.GetEnvWithDefault("FACEBOOK_APP_CATEGORY_ID", defaultCategoryID)
	exporterServicer := servicers.NewExporterServicer(
		os.Getenv("METRIC_EXPORT_URL"),
		os.Getenv("FACEBOOK_APP_ID"),
		os.Getenv("FACEBOOK_APP_SECRET"),
		categoryID,
		os.Getenv("METRICS_PREFIX"),
		servicers.ODSMetricsQueueLength,
		servicers.ODSMetricsExportInterval,
	)
	protos.RegisterMetricsExporterServer(srv.GrpcServer, exporterServicer)

	swagger_protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicerFromFile(fbinternal_service.ServiceName))

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
