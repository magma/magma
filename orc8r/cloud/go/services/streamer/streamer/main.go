/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/streamer"
	streamer_protos "magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/cloud/go/services/streamer/servicers"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, streamer.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating streamer service: %s", err)
	}

	servicer := servicers.NewStreamerServicer()
	protos.RegisterStreamerServer(srv.GrpcServer, servicer)
	streamer_protos.RegisterStreamProviderServer(srv.GrpcServer, servicers.NewBaseOrchestratorStreamProviderServicer())

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running streamer service: %s", err)
	}
}
