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
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/accessd"
	"magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/cloud/go/services/accessd/servicers"
	"magma/orc8r/cloud/go/services/accessd/storage"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/golang/glog"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, accessd.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %s", err)
	}

	// Init storage
	db, err := sqorc.Open(blobstore.SQLDriver, blobstore.DatabaseSource)
	if err != nil {
		glog.Fatalf("Failed to connect to database: %s", err)
	}
	fact := blobstore.NewEntStorage(storage.AccessdTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing accessd database: %s", err)
	}
	store := storage.NewAccessdBlobstore(fact)

	// Add servicers to the service
	accessdServer := servicers.NewAccessdServer(store)
	protos.RegisterAccessControlManagerServer(srv.GrpcServer, accessdServer)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
