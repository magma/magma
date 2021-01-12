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
	"magma/lte/cloud/go/services/ha"
	"magma/lte/cloud/go/services/ha/servicers"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(lte.ModuleName, ha.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating had service: %s", err)
	}
	servicer := servicers.NewHAServicer()
	protos.RegisterHaServer(srv.GrpcServer, servicer)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running ha service and echo server: %s", err)
	}
}
