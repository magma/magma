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
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/config/protos"
	"magma/orc8r/cloud/go/services/config/servicers"
	"magma/orc8r/cloud/go/services/config/storage"
	"magma/orc8r/cloud/go/sql_utils"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, config.ServiceName)
	if err != nil {
		log.Fatalf("Error creating config service: %s", err)
	}

	db, err := sql_utils.Open(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}
	store := storage.NewSqlConfigurationStorage(db, sql_utils.GetSqlBuilder())

	servicer := servicers.NewConfigService(store)
	protos.RegisterConfigServiceServer(srv.GrpcServer, servicer)
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running config service: %s", err)
	}
}
