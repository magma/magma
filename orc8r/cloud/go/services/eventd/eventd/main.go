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

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/eventd"
	"magma/orc8r/cloud/go/services/eventd/obsidian/handlers"
	"magma/orc8r/cloud/go/services/obsidian"
	swagger_protos "magma/orc8r/cloud/go/services/obsidian/swagger/protos"
	servicers "magma/orc8r/cloud/go/services/obsidian/swagger/servicers/protected"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, eventd.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %+v", err)
	}

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetObsidianHandlers())

	swagger_protos.RegisterSwaggerSpecServer(srv.ProtectedGrpcServer, servicers.NewSpecServicerFromFile(eventd.ServiceName))

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running eventd service: %+v", err)
	}
}
