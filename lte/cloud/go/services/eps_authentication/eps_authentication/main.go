/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"log"

	"magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/services/eps_authentication"
	"magma/lte/cloud/go/services/eps_authentication/servicers"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/sql_utils"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(lte.ModuleName, eps_authentication.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	// Init the Datastore
	store, err := datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE, sql_utils.GetSqlBuilder())
	if err != nil {
		log.Fatalf("Failed to initialize datastore: %s", err)
	}

	subscriberDBStore, err := storage.NewSubscriberDBStorage(store)
	if err != nil {
		log.Fatalf("Failed to initialize subscriberdb store: %s", err)
	}

	// Add servicers to the service
	servicer, err := servicers.NewEPSAuthServer(subscriberDBStore)
	if err != nil {
		log.Fatalf("EPS Auth Servicer Initialization Error: %s", err)
	}
	protos.RegisterS6AProxyServer(srv.GrpcServer, servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
