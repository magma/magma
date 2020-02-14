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
	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/services/dispatcher"
	"magma/orc8r/cloud/go/services/dispatcher/broker/mocks"
	"magma/orc8r/cloud/go/services/dispatcher/servicers"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func StartTestService(t *testing.T) *mocks.GatewayRPCBroker {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewEntStorage(storage.DirectorydTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	assert.NoError(t, err)
	store := storage.NewDirectorydBlobstore(fact)

	// Seed store with hwid->hostname mapping for validation at the service level
	err = store.PutHostname("some_hwid_0", "some_hostname_0")
	assert.NoError(t, err)

	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, dispatcher.ServiceName)
	mockBroker := new(mocks.GatewayRPCBroker)
	servicer, err := servicers.NewTestSyncRPCServer("test host name", mockBroker, store)
	if err != nil {
		t.Fatalf("Failed to create syncRPCService servicer: %s", err)
	}
	protos.RegisterSyncRPCServiceServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
	return mockBroker
}
