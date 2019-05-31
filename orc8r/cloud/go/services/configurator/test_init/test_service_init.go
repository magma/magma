/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/configurator/servicers"
	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Could not initialize sqlite DB: %s", err)
	}
	idGenerator := storage.DefaultIDGenerator{}
	storageFactory := storage.NewSQLConfiguratorStorageFactory(db, &idGenerator, sqorc.GetSqlBuilder())
	storageFactory.InitializeServiceStorage()

	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, configurator.ServiceName)
	nbServiser, err := servicers.NewNorthboundConfiguratorServicer(storageFactory)
	if err != nil {
		t.Fatalf("Failed to create checkin servisers: %s", err)
	}
	protos.RegisterNorthboundConfiguratorServer(srv.GrpcServer, nbServiser)

	go srv.RunTest(lis)
}
