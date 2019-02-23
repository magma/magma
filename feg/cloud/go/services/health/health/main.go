/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"log"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/health"
	"magma/feg/cloud/go/services/health/servicers"
	"magma/feg/cloud/go/services/health/storage"
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/service"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(feg.ModuleName, health.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	// Init the Datastore
	healthDatastore, err := datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE)
	if err != nil {
		log.Fatalf("Failed to initialize datastore: %s", err)
	}

	healthStore, err := storage.NewHealthStore(healthDatastore)
	if err != nil {
		log.Fatalf("Failed to initialize health store: %s", err)
	}

	clusterDatastore, err := datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE)
	if err != nil {
		log.Fatalf("Failed to initialize datastore: %s", err)
	}
	clusterStore, err := storage.NewClusterStore(clusterDatastore)
	if err != nil {
		log.Fatalf("Failed to initialize cluster store: %s", err)
	}

	// Add servicers to the service
	healthServer := servicers.NewHealthServer(healthStore, clusterStore)
	if err != nil {
		log.Fatalf("Health Servicer Initialization Error: %s", err)
	}
	protos.RegisterHealthServer(srv.GrpcServer, healthServer)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
