/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"magma/orc8r/cloud/go/blobstore"
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
	// Create service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, directoryd.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating directory service: %s", err)
	}

	// Init storage
	db, err := sqorc.Open(blobstore.SQLDriver, blobstore.DatabaseSource)
	if err != nil {
		glog.Fatalf("Error opening db connection: %s", err)
	}

	fact := blobstore.NewEntStorage(storage.DirectorydTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing directory storage: %s", err)
	}

	store := storage.NewDirectorydBlobstore(fact)

	// Add servicers
	servicer, err := servicers.NewDirectoryLookupServicer(store)
	if err != nil {
		glog.Fatalf("Error creating initializing directory servicer: %s", err)
	}
	protos.RegisterDirectoryLookupServer(srv.GrpcServer, servicer)

	// Run service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running directory service: %s", err)
	}
}
