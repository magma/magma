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

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/dispatcher"
	syncRpcBroker "magma/orc8r/cloud/go/services/dispatcher/broker"
	"magma/orc8r/cloud/go/services/dispatcher/httpserver"
	"magma/orc8r/cloud/go/services/dispatcher/servicers"
	"magma/orc8r/lib/go/protos"
	platform_service "magma/orc8r/lib/go/service"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

const (
	HttpServerPort = 9080
)

func main() {
	// Set MaxConnectionAge to infinity so Sync RPC stream doesn't restart
	var keepaliveParams = platform_service.GetDefaultKeepaliveParameters()
	keepaliveParams.MaxConnectionAge = 0
	keepaliveParams.MaxConnectionAgeGrace = 0

	// Create the service
	srv, err := service.NewOrchestratorService(
		orc8r.ModuleName,
		dispatcher.ServiceName,
		grpc.KeepaliveParams(keepaliveParams),
	)
	if err != nil {
		glog.Fatalf("Error creating service: %+v", err)
	}

	// create a broker
	broker := syncRpcBroker.NewGatewayReqRespBroker()

	// get ec2 public host name
	hostName := service.MustGetHostname()
	glog.Infof("SyncRPC hostname is %s", hostName)

	// create servicer
	syncRpcServicer, err := servicers.NewSyncRPCService(hostName, broker)
	if err != nil {
		glog.Fatalf("Error initializing syncRPC service: %+v", err)
	}
	protos.RegisterSyncRPCServiceServer(srv.GrpcServer, syncRpcServicer)

	// create http server
	httpServer := httpserver.NewSyncRPCHttpServer(broker)
	go httpServer.Run(fmt.Sprintf(":%d", HttpServerPort))

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %+v", err)
	}
}
