/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/dispatcher"
	sync_rpc_broker "magma/orc8r/cloud/go/services/dispatcher/broker"
	"magma/orc8r/cloud/go/services/dispatcher/httpserver"
	"magma/orc8r/cloud/go/services/dispatcher/servicers"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

const HTTP_SERVER_PORT = 9080

func main() {
	// Set MaxConnectionAge to infinity so Sync RPC stream doesn't restart
	var keepaliveParams = service.GetDefaultKeepaliveParameters()
	keepaliveParams.MaxConnectionAge = 0
	keepaliveParams.MaxConnectionAgeGrace = 0

	// Create the service
	srv, err := service.NewOrchestratorServiceWithOptions(
		orc8r.ModuleName,
		dispatcher.ServiceName,
		grpc.KeepaliveParams(keepaliveParams),
	)
	if err != nil {
		glog.Fatalf("Error creating service: %s", err)
	}

	// create a broker
	broker := sync_rpc_broker.NewGatewayReqRespBroker()

	// get ec2 public host name
	hostName := getHostName()
	glog.V(2).Infof("hostName is: %v\n", hostName)
	// create servicer
	syncRpcServicer, err := servicers.NewSyncRPCService(hostName, broker)
	if err != nil {
		glog.Fatalf("SyncRPCService Initialization Error: %s", err)
	}

	// create http server
	httpServer := httpserver.NewSyncRPCHttpServer(broker)

	protos.RegisterSyncRPCServiceServer(srv.GrpcServer, syncRpcServicer)
	srv.GrpcServer.RegisterService(protos.GetLegacyDispatcherDesc(), syncRpcServicer)

	// run http server
	go httpServer.Run(fmt.Sprintf(":%d", HTTP_SERVER_PORT))

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}

// getHostName of the current SyncRPCService instance
func getHostName() string {
	// If there is env variable override, use the env variable
	// This can be used in dev cloud
	hostName, exist := os.LookupEnv("SERVICE_HOST_NAME")
	if exist {
		return hostName
	}
	//see https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/public-hostname")
	if err != nil {
		glog.Fatalf("Cannot get public-hostname of the current service instance")
	}
	if resp.StatusCode != 200 {
		errMsg, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		glog.Fatalf("Failed to getHostName: status code %d: %s", resp.StatusCode, errMsg)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body)
}
