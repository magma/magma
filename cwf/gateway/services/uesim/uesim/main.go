/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

// This starts the user equipment (ue) service.
package main

import (
	"magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/registry"
	"magma/cwf/gateway/services/uesim/servicers"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.UeSim)
	if err != nil {
		glog.Fatalf("Error creating UeSim service: %s", err)
	}

	store := blobstore.NewMemoryBlobStorageFactory()
	servicer, err := servicers.NewUESimServer(store)
	if err != nil {
		glog.Fatalf("Error creating UE server: %s", err)
	}
	protos.RegisterUESimServer(srv.GrpcServer, servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running UE service: %s", err)
	}
}
