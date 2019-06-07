/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package main implements WiFi AAA server
package main

import (
	"log"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/aaa/servicers"
	"magma/feg/gateway/services/aaa/store"
	"magma/orc8r/cloud/go/service"
)

func main() {
	// Create a shared Session Table
	sessions := store.NewMemorySessionTable()

	// Create the EAP AKA Provider service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.AAA)
	if err != nil {
		log.Fatalf("Error creating EAP service: %s", err)
	}
	auth, _ := servicers.NewEapAuthenticator(nil)
	protos.RegisterAuthenticatorServer(srv.GrpcServer, auth)

	acct, _ := servicers.NewAccountingService(sessions)
	protos.RegisterAccountingServer(srv.GrpcServer, acct)

	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running AAA service: %s", err)
	}
}
