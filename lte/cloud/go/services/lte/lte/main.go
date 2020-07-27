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
	"magma/lte/cloud/go/services/lte/obsidian/handlers"
	"magma/lte/cloud/go/services/lte/servicers"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/streamer/protos"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(lte.ModuleName, lte_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating LTE service: %s", err)
	}

	protos.RegisterStreamProviderServer(srv.GrpcServer, servicers.NewLTEStreamProviderServicer())
	obsidian.AttachHandlers(srv.EchoServer, handlers.GetHandlers())
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running LTE service and echo server: %s", err)
	}
}
