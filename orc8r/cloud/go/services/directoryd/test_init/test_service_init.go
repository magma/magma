/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/directoryd/servicers"
	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

func StartTestService(t *testing.T) {
	// Create service
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, directoryd.ServiceName)

	// Init storage
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewSQLBlobStorageFactory(storage.DirectorydTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	assert.NoError(t, err)
	store := storage.NewDirectorydBlobstore(fact)

	// Add servicers
	servicer, err := servicers.NewDirectoryLookupServicer(store)
	assert.NoError(t, err)
	protos.RegisterDirectoryLookupServer(srv.GrpcServer, servicer)

	// Run service
	go srv.RunTest(lis)
}
