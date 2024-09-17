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
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb"
	"magma/lte/cloud/go/services/policydb/obsidian/handlers"
	policydb_servicer "magma/lte/cloud/go/services/policydb/servicers/southbound"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/obsidian"
	swagger_protos "magma/orc8r/cloud/go/services/obsidian/swagger/protos"
	swaggger_servicers "magma/orc8r/cloud/go/services/obsidian/swagger/servicers/protected"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(lte.ModuleName, policydb.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %s", err)
	}
	assignmentServicer := policydb_servicer.NewPolicyAssignmentServer()
	protos.RegisterPolicyAssignmentControllerServer(srv.GrpcServer, assignmentServicer)

	swagger_protos.RegisterSwaggerSpecServer(srv.ProtectedGrpcServer, swaggger_servicers.NewSpecServicerFromFile(policydb.ServiceName))

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers())
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service and echo server: %s", err)
	}
}
