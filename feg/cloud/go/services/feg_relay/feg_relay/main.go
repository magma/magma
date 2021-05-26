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
	"fmt"

	"github.com/golang/glog"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/feg_relay"
	"magma/feg/cloud/go/services/feg_relay/gw_to_feg_relay"
	nh_servicers "magma/feg/cloud/go/services/feg_relay/gw_to_feg_relay/servicers"
	"magma/feg/cloud/go/services/feg_relay/servicers"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/service"
)

const GwToFeGServerPort = 9079

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(feg.ModuleName, feg_relay.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating Feg Proxy service: %s", err)
	}
	servicer, err := servicers.NewFegToGwRelayServer()

	if err != nil {
		glog.Fatalf("Failed to create FegToGwRelayServer: %v", err)
		return
	}

	// Register responders FEG -> AGW
	lteprotos.RegisterSessionProxyResponderServer(srv.GrpcServer, servicer)
	lteprotos.RegisterAbortSessionResponderServer(srv.GrpcServer, servicer)
	protos.RegisterS8ProxyResponderServer(srv.GrpcServer, servicer)

	// Register services AGW -> FEG
	protos.RegisterS6AGatewayServiceServer(srv.GrpcServer, servicer)
	protos.RegisterCSFBGatewayServiceServer(srv.GrpcServer, servicer)
	protos.RegisterSwxGatewayServiceServer(srv.GrpcServer, servicer)

	// Register Neutral Host Routing services
	nhServicer := nh_servicers.NewRelayRouter()
	protos.RegisterS6AProxyServer(srv.GrpcServer, nhServicer)
	protos.RegisterSwxProxyServer(srv.GrpcServer, nhServicer)
	protos.RegisterHelloServer(srv.GrpcServer, nhServicer)
	lteprotos.RegisterCentralSessionControllerServer(srv.GrpcServer, nhServicer)

	// Register S8 Proxy Neutral Host Routing services
	s8nhServicer := nh_servicers.NewS8RelayRouter(&nhServicer.Router)
	protos.RegisterS8ProxyServer(srv.GrpcServer, s8nhServicer)

	// create and run GW_TO_FEG httpserver
	gwToFeGServer := gw_to_feg_relay.NewGatewayToFegServer()
	go gwToFeGServer.Run(fmt.Sprintf(":%d", GwToFeGServerPort))
	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
