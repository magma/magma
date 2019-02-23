/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"fmt"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/feg_relay"
	"magma/feg/cloud/go/services/feg_relay/gw_to_feg_relay"
	"magma/feg/cloud/go/services/feg_relay/servicers"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
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
	protos.RegisterS6AGatewayServiceServer(srv.GrpcServer, servicer)
	protos.RegisterCSFBGatewayServiceServer(srv.GrpcServer, servicer)
	lteprotos.RegisterSessionProxyResponderServer(srv.GrpcServer, servicer)
	// create and run GW_TO_FEG httpserver
	gwToFeGServer := gw_to_feg_relay.NewGatewayToFegServer()
	go gwToFeGServer.Run(fmt.Sprintf(":%d", GwToFeGServerPort))
	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
