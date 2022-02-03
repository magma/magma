/*
Copyright 2021 The Magma Authors.

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
	"log"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/n7_n40_proxy/n7"
	"magma/feg/gateway/services/testcore/pcf/servicers"
	"magma/orc8r/lib/go/service"
)

func main() {
	n7Config, err := n7.GetN7Config()
	if err != nil {
		log.Fatalf("Error reading N7 config: %s", err)
	}
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.MOCK_PCF)
	if err != nil {
		log.Fatalf("Error creating mock %s service: %s", registry.MOCK_PCF, err)
	}
	pcfServer, err := servicers.NewMockPCFServer(&n7Config.ClientConfig)
	if err != nil {
		log.Fatalf("Error creating Mock PCF server: %s", err)
	}
	protos.RegisterMockPCFServer(srv.GrpcServer, pcfServer)
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running mock %s service: %s", registry.MOCK_PCF, err)
	}
}
