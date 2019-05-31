/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Magmad Configurator is a service which stores & manages all cloud magma records.
package main

import (
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/magmad/servicers"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/golang/glog"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, magmad.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %s", err)
	}

	// Init the Datastore
	ds, err :=
		datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE, sqorc.GetSqlBuilder())
	if err != nil {
		glog.Fatalf("Failed to initialize datastore: %s", err)
	}

	// Add servicers to the service
	magmadServer := servicers.NewMagmadConfigurator(ds)
	protos.RegisterMagmadConfiguratorServer(srv.GrpcServer, magmadServer)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
