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

// Magma's Swx Proxy Service converts gRPC requests into Swx protocol over diameter
package main

import (
	"flag"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/swx_proxy/servicers"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.SWX_PROXY)
	if err != nil {
		glog.Fatalf("Error creating Swx Proxy service: %s", err)
	}

	// Create servicers
	servicer, err := servicers.NewSwxProxiesWithHealthAndDefaultMultiplexor(
		servicers.GetSwxProxyConfig())
	if err != nil {
		glog.Fatalf("Failed to create SwxProxy: %v", err)
	}

	// Register services
	protos.RegisterSwxProxyServer(srv.GrpcServer, servicer)
	protos.RegisterServiceHealthServer(srv.GrpcServer, servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
