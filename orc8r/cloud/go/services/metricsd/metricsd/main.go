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
	"time"

	"github.com/golang/glog"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"google.golang.org/grpc"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/collection"
	"magma/orc8r/cloud/go/services/metricsd/obsidian/handlers"
	"magma/orc8r/cloud/go/services/metricsd/servicers/protected"
	"magma/orc8r/cloud/go/services/metricsd/servicers/southbound"
	"magma/orc8r/cloud/go/services/obsidian"
	swagger_protos "magma/orc8r/cloud/go/services/obsidian/swagger/protos"
	swagger_servicers "magma/orc8r/cloud/go/services/obsidian/swagger/servicers/protected"
	"magma/orc8r/lib/go/protos"
)

const (
	CloudMetricsCollectInterval = time.Second * 20
	// Setting Max Received Message gRPC size to 50MB
	CloudMetricsCollectMaxMsgSize = 50 * 1024 * 1024
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName,
		metricsd.ServiceName,
		service.WithGrpcOptions(grpc.MaxRecvMsgSize(CloudMetricsCollectMaxMsgSize)))

	if err != nil {
		glog.Fatalf("Error creating orc8r service for metricsd: %s", err)
	}

	cloudControllerServicer := protected.NewCloudMetricsControllerServer()
	protos.RegisterCloudMetricsControllerServer(srv.ProtectedGrpcServer, cloudControllerServicer)

	controllerServicer := southbound.NewMetricsControllerServer()
	protos.RegisterMetricsControllerServer(srv.GrpcServer, controllerServicer)

	swagger_protos.RegisterSwaggerSpecServer(srv.ProtectedGrpcServer, swagger_servicers.NewSpecServicerFromFile(metricsd.ServiceName))

	// Initialize gatherers
	additionalCollectors := []collection.MetricCollector{
		&collection.DiskUsageMetricCollector{},
		&collection.ProcMetricsCollector{},
	}
	metricsCh := make(chan *io_prometheus_client.MetricFamily)
	gatherer, err := collection.NewMetricsGatherer(additionalCollectors, CloudMetricsCollectInterval, metricsCh)
	if err != nil {
		glog.Fatalf("Error initializing MetricsGatherer: %s", err)
	}
	go cloudControllerServicer.ConsumeCloudMetrics(metricsCh, service.MustGetHostname())
	gatherer.Run()

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetObsidianHandlers(srv.Config))
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running metricsd service: %s", err)
	}
}
