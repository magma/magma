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
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/service"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	wifi_service "magma/wifi/cloud/go/services/wifi"
	"magma/wifi/cloud/go/services/wifi/obsidian/handlers"
	"magma/wifi/cloud/go/services/wifi/servicers"
	"magma/wifi/cloud/go/wifi"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(wifi.ModuleName, wifi_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating wifi service %s", err)
	}

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers())

	builder_protos.RegisterMconfigBuilderServer(srv.GrpcServer, servicers.NewBuilderServicer())

	swagger_protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicerFromFile(wifi_service.ServiceName))

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %s", err)
	}
}
