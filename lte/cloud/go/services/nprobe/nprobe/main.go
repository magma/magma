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
	"magma/lte/cloud/go/services/nprobe"
	"magma/lte/cloud/go/services/nprobe/obsidian/handlers"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	"magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

func main() {
	// Create service
	srv, err := service.NewOrchestratorService(lte.ModuleName, nprobe.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %v", err)
	}

	// Attach handlers
	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers())
	protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicerFromFile(nprobe.ServiceName))

	// Run service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %v", err)
	}

}
