/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"log"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/streamer"
	"magma/orc8r/cloud/go/services/streamer/servicers"
	"magma/orc8r/lib/go/protos"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, streamer.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	// Add servicers to the service
	servicer := &servicers.StreamingServer{}
	protos.RegisterStreamerServer(srv.GrpcServer, servicer)
	srv.GrpcServer.RegisterService(protos.GetLegacyStreamerDesc(), servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
