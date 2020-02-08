/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"crypto/rsa"
	"flag"
	"log"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/security/key"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/services/bootstrapper/servicers"
	"magma/orc8r/lib/go/protos"
)

var (
	keyFile = flag.String("cak", "bootstrapper.key.pem", "Bootstrapper's Private Key file")
)

func main() {
	// Create the service, flag will be parsed inside this function
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, bootstrapper.ServiceName)
	if err != nil {
		log.Fatalf("Error creating bootstrapper service: %s", err)
	}

	// Add servicers to the service
	privKey, err := key.ReadKey(*keyFile)
	if err != nil {
		log.Fatalf("Failed to read private key: %s", err)
	}
	servicer, err := servicers.NewBootstrapperServer(privKey.(*rsa.PrivateKey))
	if err != nil {
		log.Fatalf("Failed to create bootstrapper servicer: %s", err)
	}
	protos.RegisterBootstrapperServer(srv.GrpcServer, servicer)
	srv.GrpcServer.RegisterService(protos.GetLegacyBootstrapperDesc(), servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
