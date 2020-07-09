/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package main

import (
	"magma/lte/cloud/go/lte"
	lte_service "magma/lte/cloud/go/services/lte"
	"magma/lte/cloud/go/services/lte/servicers"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/streamer/protos"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(lte.ModuleName, lte_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating LTE service: %s", err)
	}

	protos.RegisterStreamProviderServer(srv.GrpcServer, servicers.NewLTEStreamProviderServicer())

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running LTE service and echo server: %s", err)
	}
}
