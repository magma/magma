/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/directoryd/servicers"
	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

func main() {
	// Create Magma micro-service
	directoryService, err := service.NewOrchestratorService(orc8r.ModuleName, directoryd.ServiceName)
	if err != nil {
		glog.Errorf("Error creating directory service: %s", err)
	}
	glog.V(2).Info("Init Directory Service...")

	// Init Datastore
	db, err := datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE, sqorc.GetSqlBuilder())
	if err != nil {
		glog.Errorf("Failed to initialize datastore: %s", err)
	}
	store := storage.GetDirectorydPersistenceService(db)

	// Create directory gRPC servicer
	directorydServicer, err := servicers.NewDirectoryServicer(store)
	if err != nil {
		glog.Errorf("Error creating directory gRPC servicer: %s", err)
	}

	// Add gRPC servicer to the directory service gRPC server
	protos.RegisterDirectoryServiceServer(directoryService.GrpcServer, directorydServicer)
	directoryService.GrpcServer.RegisterService(protos.GetLegacyDirectorydDesc(), directorydServicer)

	// Run the service
	glog.V(2).Info("Starting Directory Service...")
	err = directoryService.Run()
	if err != nil {
		glog.Errorf("Error running directory service: %s", err)
	}
}
