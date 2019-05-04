/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package main

import (
	"database/sql"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/servicers"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, state.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating state service %s", err)
	}

	db, err := sql.Open(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE)
	if err != nil {
		glog.Fatalf("Failed to connect to database: %s", err)
	}
	store := blobstore.NewSQLBlobStorageFactory(state.DBTableName, db)
	err = store.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing state database: %s", err)
	}

	server, err := servicers.NewStateServicer(store)
	if err != nil {
		glog.Fatalf("Error creating state server: %s", err)
	}
	protos.RegisterStateServiceServer(srv.GrpcServer, server)
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
