/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"magma/cwf/cloud/go/cwf"
	cwf_service "magma/cwf/cloud/go/services/cwf"
	cwf_analytics "magma/cwf/cloud/go/services/cwf/analytics"
	"magma/cwf/cloud/go/services/cwf/obsidian/handlers"
	"magma/cwf/cloud/go/services/cwf/servicers"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/analytics"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(cwf.ModuleName, cwf_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating cwf service %s", err)
	}

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers())

	builder_protos.RegisterMconfigBuilderServer(srv.GrpcServer, servicers.NewBuilderServicer())

	swagger_protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicerFromFile(cwf_service.ServiceName))

	var serviceConfig cwf_service.Config
	_, _, err = config.GetStructuredServiceConfig(cwf.ModuleName, cwf_service.ServiceName, &serviceConfig)
	if err != nil {
		glog.Infof("Failed unmarshalling service config %v", err)
		return
	}
	promQLClient := analytics.GetPrometheusClient()
	calcs := cwf_analytics.GetAnalyticsCalculations(&serviceConfig.Analytics)
	userStateManager := calculations.NewUserStateManager(promQLClient, "active_sessions")
	collectorServicer := analytics.NewCollectorServicer(&serviceConfig.Analytics, promQLClient, calcs, userStateManager)
	protos.RegisterAnalyticsCollectorServer(srv.GrpcServer, collectorServicer)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %s", err)
	}
}
