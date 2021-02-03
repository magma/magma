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
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb"
	"magma/lte/cloud/go/services/policydb/obsidian/handlers"
	"magma/lte/cloud/go/services/policydb/servicers"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/swagger"
	swagger_protos "magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(lte.ModuleName, policydb.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %s", err)
	}
	assignmentServicer := servicers.NewPolicyAssignmentServer()
	protos.RegisterPolicyAssignmentControllerServer(srv.GrpcServer, assignmentServicer)

	swagger_protos.RegisterSwaggerSpecServer(srv.GrpcServer, swagger.NewSpecServicerFromFile(policydb.ServiceName))

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers())
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %s", err)
	}
}
