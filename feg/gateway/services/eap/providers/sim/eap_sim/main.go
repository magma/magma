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

// Package main implements Magma EAP SIM Service
package main

import (
	"flag"

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/eap/protos"
	"magma/feg/gateway/services/eap/providers/sim"
	"magma/feg/gateway/services/eap/providers/sim/servicers"
	_ "magma/feg/gateway/services/eap/providers/sim/servicers/handlers"
	managed_configs "magma/gateway/mconfig"
	"magma/orc8r/lib/go/service"
)

func init() {
	flag.Parse()
}

func main() {
	// Create the EAP SIM Provider service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.EAP_SIM)
	if err != nil {
		glog.Fatalf("Error creating EAP SIM service: %s", err)
	}

	simConfigs := &mconfig.EapSimConfig{}
	err = managed_configs.GetServiceConfigs(sim.EapSimServiceName, simConfigs)
	if err != nil {
		glog.Errorf("Error getting EAP SIM service configs: %s", err)
		simConfigs = nil
	}
	servicer, err := servicers.NewEapSimService(simConfigs)
	if err != nil {
		glog.Fatalf("failed to create EAP SIM Service: %v", err)
		return
	}
	protos.RegisterEapServiceServer(srv.GrpcServer, servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running EAP SIM service: %s", err)
	}
}
