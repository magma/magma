/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Access Control Manager is a service which stores, manages and verifies
// operator Identity objects and their rights to access (read/write) Entities.
package main

import (
	"log"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/accessd"
	"magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/cloud/go/services/accessd/servicers"
	"magma/orc8r/cloud/go/services/accessd/storage"
	"magma/orc8r/cloud/go/sqorc"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, accessd.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	// Init the Datastore
	ds, err :=
		datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE, sqorc.GetSqlBuilder())
	if err != nil {
		log.Fatalf("Failed to initialize datastore: %s", err)
	}

	// Add servicers to the service
	accessdStore := storage.NewAccessdDatastore(ds)
	accessdServer := servicers.NewAccessdServer(accessdStore)
	protos.RegisterAccessControlManagerServer(srv.GrpcServer, accessdServer)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
