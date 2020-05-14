/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"database/sql"
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/servicers"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

// StartTestService instantiates a service backed by an in-memory storage.
func StartTestService(t *testing.T) {
	factory := blobstore.NewMemoryBlobStorageFactory()
	server, err := servicers.NewStateServicer(factory)
	assert.NoError(t, err)
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, state.ServiceName)
	protos.RegisterStateServiceServer(srv.GrpcServer, server)

	go srv.RunTest(lis)
}

// StartTestServiceInternal instantiates a test DB-backed service, clears the
// DB, and returns the DB and its encasing internal servicer methods.
// Supported drivers include {postgres, mysql}.
func StartTestServiceInternal(t *testing.T, dbName, dbDriver string) (*sql.DB, servicers.StateServiceInternal) {
	db := sqorc.OpenCleanForTest(t, dbName, dbDriver)

	factory := blobstore.NewSQLBlobStorageFactory(state.DBTableName, db, sqorc.GetSqlBuilder())
	assert.NoError(t, factory.InitializeFactory())
	servicer, err := servicers.NewStateServicer(factory)
	assert.NoError(t, err)
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, state.ServiceName)
	protos.RegisterStateServiceServer(srv.GrpcServer, servicer)

	go srv.RunTest(lis)
	return db, servicer
}
