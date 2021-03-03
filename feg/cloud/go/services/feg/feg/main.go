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
	"magma/feg/cloud/go/feg"
	feg_service "magma/feg/cloud/go/services/feg"
	"magma/feg/cloud/go/services/feg/obsidian/handlers"
	"magma/feg/cloud/go/services/feg/servicers"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/service"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(feg.ModuleName, feg_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating feg service %s", err)
	}

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers())

	builder_protos.RegisterMconfigBuilderServer(srv.GrpcServer, servicers.NewBuilderServicer())

	swagger_protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicerFromFile(feg_service.ServiceName))

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %s", err)
	}
}
