/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"log"
	"time"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/health"
	"magma/feg/cloud/go/services/health/reporter"
	"magma/feg/cloud/go/services/health/servicers"
	"magma/feg/cloud/go/services/health/storage"
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/sqorc"
)

const (
	NETWORK_HEALTH_STATUS_REPORT_INTERVAL = time.Second * 60
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(feg.ModuleName, health.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	// Init the Datastore
	healthDatastore, err := datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE, sqorc.GetSqlBuilder())
	if err != nil {
		log.Fatalf("Failed to initialize datastore: %s", err)
	}

	healthStore, err := storage.NewHealthStore(healthDatastore)
	if err != nil {
		log.Fatalf("Failed to initialize health store: %s", err)
	}

	clusterDatastore, err := datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE, sqorc.GetSqlBuilder())
	if err != nil {
		log.Fatalf("Failed to initialize datastore: %s", err)
	}
	clusterStore, err := storage.NewClusterStore(clusterDatastore)
	if err != nil {
		log.Fatalf("Failed to initialize cluster store: %s", err)
	}

	// Add servicers to the service
	healthServer := servicers.NewHealthServer(healthStore, clusterStore)
	protos.RegisterHealthServer(srv.GrpcServer, healthServer)

	// create a networkHealthStatusReporter to monitor and periodically log metrics
	// on if all gateways in a network are unhealthy
	healthStatusReporter, err := reporter.NewNetworkHealthStatusReporter(healthStore)
	if err != nil {
		log.Fatalf("NetworkHealthStatusReporter Initialization Error: %s\n", err)
	}
	go healthStatusReporter.ReportHealthStatus(NETWORK_HEALTH_STATUS_REPORT_INTERVAL)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
