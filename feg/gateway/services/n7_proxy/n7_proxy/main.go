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
	"flag"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/policydb"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/n7_proxy/servicers"

	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.N7_PROXY)
	if err != nil {
		glog.Fatalf("Error creating N7 Proxy service: %s", err)
	}

	config := servicers.GetN7ProxyConfig()
	if config == nil {
		glog.Fatalf("Error loading config for N7 Proxy service: %s", err)
	}

	controlParams, err := servicers.GetControllerParams(config)
	if err != nil {
		glog.Fatalf("Error creating N7 Proxy service: %s", err)
	}

	cloudReg := registry.Get()
	dbClient, err := policydb.NewRedisPolicyDBClient(cloudReg)
	if err != nil {
		glog.Fatalf("Error connecting to redis store from N7 Proxy: %s", err)
	}

	sessController, err := servicers.NewCentralSessionControllerDefaultMultiplex(controlParams, dbClient)
	if err != nil {
		glog.Fatalf("Error creating session controller in N7 Proxy: %s", err)
	}

	// Add GRPC handlers to the service
	lteprotos.RegisterCentralSessionControllerServer(srv.GrpcServer, sessController)
	protos.RegisterServiceHealthServer(srv.GrpcServer, sessController)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running N7 service: %s", err)
	}
}
