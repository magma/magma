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
	"magma/lte/cloud/go/lte"
	lte_service "magma/lte/cloud/go/services/lte"
	lte_analytics "magma/lte/cloud/go/services/lte/analytics"
	"magma/lte/cloud/go/services/lte/obsidian/handlers"
	lte_protos "magma/lte/cloud/go/services/lte/protos"
	"magma/lte/cloud/go/services/lte/servicers"
	lte_storage "magma/lte/cloud/go/services/lte/storage"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/analytics"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	state_protos "magma/orc8r/cloud/go/services/state/protos"
	provider_protos "magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(lte.ModuleName, lte_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating lte service: %s", err)
	}

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers())

	builder_protos.RegisterMconfigBuilderServer(srv.GrpcServer, servicers.NewBuilderServicer())
	provider_protos.RegisterStreamProviderServer(srv.GrpcServer, servicers.NewProviderServicer())
	state_protos.RegisterIndexerServer(srv.GrpcServer, servicers.NewIndexerServicer())

	swagger_protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicerFromFile(lte_service.ServiceName))

	// Init storage
	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error opening db connection: %v", err)
	}
	enbStateStore := lte_storage.NewEnodebStateLookup(db, sqorc.GetSqlBuilder())
	if err := enbStateStore.Initialize(); err != nil {
		glog.Fatalf("Error initializing enodeb state lookup storage: %v", err)
	}
	lte_protos.RegisterEnodebStateLookupServer(srv.GrpcServer, servicers.NewLookupServicer(enbStateStore))

	var serviceConfig lte_service.Config
	_, _, err = config.GetStructuredServiceConfig(lte.ModuleName, lte_service.ServiceName, &serviceConfig)
	if err != nil {
		glog.Infof("Failed unmarshaling service config %v", err)
		return
	}

	// Initialize analytics
	// userStateExpr is a metric which enables us to compute the number of active users using the network
	promQLClient := analytics.GetPrometheusClient()
	userStateManager := calculations.NewUserStateManager(promQLClient, "ue_connected")
	calcs := lte_analytics.GetAnalyticsCalculations(&serviceConfig.Analytics)
	collectorServicer := analytics.NewCollectorServicer(&serviceConfig.Analytics, promQLClient, calcs, userStateManager)
	protos.RegisterAnalyticsCollectorServer(srv.GrpcServer, collectorServicer)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running lte service and echo server: %s", err)
	}
}
