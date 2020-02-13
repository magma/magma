/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"log"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/eps_authentication"
	"magma/lte/cloud/go/services/eps_authentication/servicers"
	"magma/lte/cloud/go/services/eps_authentication/storage"
	"magma/orc8r/cloud/go/service"
)

// DEPRECATED -- eps_authentication service is temporarily deprecated and currently not in-use.
// NOT WORKING -- calling any of the service's endpoints will result in a handler panic.
// To un-deprecate:
// 	- Point storage at subscriber data in configurator
//	- Search in directory for all items denoted DEPRECATED (e.g. skipped tests)
func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(lte.ModuleName, eps_authentication.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	// Add servicers to the service
	store := storage.NewSubscriberDBStorage()
	servicer, err := servicers.NewEPSAuthServer(store)
	if err != nil {
		log.Fatalf("EPS Auth Servicer Initialization Error: %s", err)
	}
	protos.RegisterEPSAuthenticationServer(srv.GrpcServer, servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
