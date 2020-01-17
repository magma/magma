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
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/servicers"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/sqorc"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(lte.ModuleName, subscriberdb.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	// Init the Datastore
	store, err :=
		datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE, sqorc.GetSqlBuilder())
	if err != nil {
		log.Fatalf("Failed to initialize datastore: %s", err)
	}

	subscriberDBStore, err := storage.NewSubscriberDBStorage(store)
	if err != nil {
		log.Fatalf("Failed to initialize subscriberdb store: %s", err)
	}

	// Add servicers to the service
	servicer, err := servicers.NewSubscriberDBServer(subscriberDBStore)
	if err != nil {
		log.Fatalf("Subscriberdb Servicer Initialization Error: %s", err)
	}
	protos.RegisterSubscriberDBControllerServer(srv.GrpcServer, servicer)
	srv.GrpcServer.RegisterService(protos.GetLegacySubscriberdbDesc(), servicer)

	assignmentServicer := servicers.NewPolicyAssignmentServer()
	protos.RegisterPolicyAssignmentControllerServer(srv.GrpcServer, assignmentServicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
