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
	"github.com/golang/glog"

	"magma/lte/cloud/go/lte"
	lte_service "magma/lte/cloud/go/services/lte"
	lte_analytics "magma/lte/cloud/go/services/lte/analytics"
	"magma/lte/cloud/go/services/lte/obsidian/handlers"
	lte_protos "magma/lte/cloud/go/services/lte/protos"
	"magma/lte/cloud/go/services/lte/servicers"
	protected_servicers "magma/lte/cloud/go/services/lte/servicers/protected"
	lte_storage "magma/lte/cloud/go/services/lte/storage"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/analytics"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	analytics_servicer "magma/orc8r/cloud/go/services/analytics/servicers/protected"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	"magma/orc8r/cloud/go/services/obsidian"
	swagger_protos "magma/orc8r/cloud/go/services/obsidian/swagger/protos"
	swagger_servicers "magma/orc8r/cloud/go/services/obsidian/swagger/servicers/protected"
	state_protos "magma/orc8r/cloud/go/services/state/protos"
	provider_protos "magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/service/config"
)

func main() {
	srv, err := service.NewOrchestratorService(lte.ModuleName, lte_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating lte service: %s", err)
	}

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers())

	var serviceConfig lte_service.Config
	config.MustGetStructuredServiceConfig(lte.ModuleName, lte_service.ServiceName, &serviceConfig)

	builder_protos.RegisterMconfigBuilderServer(srv.ProtectedGrpcServer, protected_servicers.NewBuilderServicer(serviceConfig))
	provider_protos.RegisterStreamProviderServer(srv.ProtectedGrpcServer, servicers.NewProviderServicer())
	state_protos.RegisterIndexerServer(srv.ProtectedGrpcServer, protected_servicers.NewIndexerServicer())

	swagger_protos.RegisterSwaggerSpecServer(srv.ProtectedGrpcServer, swagger_servicers.NewSpecServicerFromFile(lte_service.ServiceName))

	// Init storage
	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error opening db connection: %v", err)
	}
	enbStateStore := lte_storage.NewEnodebStateLookup(db, sqorc.GetSqlBuilder())
	if err := enbStateStore.Initialize(); err != nil {
		glog.Fatalf("Error initializing enodeb state lookup storage: %v", err)
	}
	lte_protos.RegisterEnodebStateLookupServer(srv.ProtectedGrpcServer, protected_servicers.NewLookupServicer(enbStateStore))

	// Initialize analytics
	// userStateExpr is a metric which enables us to compute the number of active users using the network
	promQLClient := analytics.GetPrometheusClient()
	userStateManager := calculations.NewUserStateManager(promQLClient, "ue_connected")
	calcs := lte_analytics.GetAnalyticsCalculations(&serviceConfig.Analytics)
	collectorServicer := analytics_servicer.NewCollectorServicer(&serviceConfig.Analytics, promQLClient, calcs, userStateManager)
	protos.RegisterAnalyticsCollectorServer(srv.ProtectedGrpcServer, collectorServicer)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running lte service and echo server: %s", err)
	}
}
