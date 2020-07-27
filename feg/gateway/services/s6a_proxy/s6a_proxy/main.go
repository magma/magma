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

// Magma LTE S6a Proxy Service
// The service is only exposed to other cloud services and should not be
// accessible to outside clients
package main

import (
	"flag"
	"log"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/s6a_proxy/servicers"
	"magma/orc8r/lib/go/service"
)

func init() {
	flag.Parse()
}

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.S6A_PROXY)
	if err != nil {
		log.Fatalf("Error creating S6a Proxy service: %s", err)
	}

	servicer, err := servicers.NewS6aProxy(servicers.GetS6aProxyConfigs())
	if err != nil {
		log.Fatalf("failed to create S6aProxy: %v", err)
		return
	}
	protos.RegisterS6AProxyServer(srv.GrpcServer, servicer)
	protos.RegisterServiceHealthServer(srv.GrpcServer, servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
