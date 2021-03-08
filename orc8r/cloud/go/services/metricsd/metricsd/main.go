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

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/collection"
	"magma/orc8r/cloud/go/services/metricsd/obsidian/handlers"
	"magma/orc8r/cloud/go/services/metricsd/servicers"
	"magma/orc8r/lib/go/protos"

	"github.com/getsentry/sentry-go"
	"github.com/golang/glog"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"google.golang.org/grpc"
)

const (
	CloudMetricsCollectInterval = time.Second * 20
	// Setting Max Received Message gRPC size to 50MB
	CloudMetricsCollectMaxMsgSize = 50 * 1024 * 1024
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://f6a54d1a20134c258b1e0b227d4d0982@o529355.ingest.sentry.io/5667116",
	})
	if err != nil {
		glog.Fatalf("sentry.Init: %s", err)
	}
	srv, err := service.NewOrchestratorService(orc8r.ModuleName,
		metricsd.ServiceName,
		grpc.MaxRecvMsgSize(CloudMetricsCollectMaxMsgSize))

	if err != nil {
		glog.Fatalf("Error creating orc8r service for metricsd: %s", err)
	}

	controllerServicer := servicers.NewMetricsControllerServer()
	protos.RegisterMetricsControllerServer(srv.GrpcServer, controllerServicer)

	swagger_protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicerFromFile(metricsd.ServiceName))

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
	go controllerServicer.ConsumeCloudMetrics(metricsCh, service.MustGetHostname())
	gatherer.Run()

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetObsidianHandlers(srv.Config))
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running metricsd service: %s", err)
	}
}
