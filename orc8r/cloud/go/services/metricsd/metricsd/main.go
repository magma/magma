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
	"time"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/collection"
	"magma/orc8r/cloud/go/services/metricsd/obsidian/handlers"
	exporter_protos "magma/orc8r/cloud/go/services/metricsd/protos"
	"magma/orc8r/cloud/go/services/metricsd/servicers"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/prometheus/client_model/go"
)

const (
	CloudMetricsCollectInterval = time.Second * 20
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, metricsd.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating orc8r service for metricsd: %s", err)
	}

	controllerServicer := servicers.NewMetricsControllerServer()
	protos.RegisterMetricsControllerServer(srv.GrpcServer, controllerServicer)

	exporterServicer := servicers.NewPushExporterServicer(srv.Config.MustGetStrings(metricsd.PrometheusPushAddresses))
	exporter_protos.RegisterMetricsExporterServer(srv.GrpcServer, exporterServicer)

	// Initialize gatherers
	metricsCh := make(chan *io_prometheus_client.MetricFamily)
	gatherer, err := collection.NewMetricsGatherer(getCollectors(), CloudMetricsCollectInterval, metricsCh)
	if err != nil {
		glog.Fatalf("Error initializing MetricsGatherer: %s", err)
	}
	go controllerServicer.ConsumeCloudMetrics(metricsCh, os.Getenv("HOST_NAME"))
	gatherer.Run()

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetObsidianHandlers(srv.Config))
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running metricsd service: %s", err)
	}
}

// getCollectors returns the set of metrics collectors.
// Returned collectors include disk usage, process statistics, and
// per-service custom metrics.
func getCollectors() []collection.MetricCollector {
	services := registry.ListControllerServices()

	collectors := []collection.MetricCollector{
		&collection.DiskUsageMetricCollector{},
		&collection.ProcMetricsCollector{},
	}
	for _, s := range services {
		collectors = append(collectors, collection.NewCloudServiceMetricCollector(s))
	}

	return collectors
}
