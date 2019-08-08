/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"log"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/upgrade"
	"magma/orc8r/cloud/go/services/upgrade/protos"
	"magma/orc8r/cloud/go/services/upgrade/servicers"
	"magma/orc8r/cloud/go/sqorc"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, upgrade.ServiceName)
	if err != nil {
		log.Fatalf("error creating service: %s", err)
	}

	// Init the Datastore
	store, err := datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE, sqorc.GetSqlBuilder())
	if err != nil {
		log.Fatalf("Failed to initialize datastore: %s", err)
	}

	// Add servicers to the service
	servicer := servicers.NewUpgradeService(store)
	protos.RegisterUpgradeServiceServer(srv.GrpcServer, servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
