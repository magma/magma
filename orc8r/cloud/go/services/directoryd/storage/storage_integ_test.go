/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package storage_test

import (
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/sqorc"
	merrors "magma/orc8r/lib/go/errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestDirectorydStorageBlobstore_Integation(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewEntStorage(storage.DirectorydTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	assert.NoError(t, err)
	store := storage.NewDirectorydBlobstore(fact)
	testDirectorydStorageImpl(t, store)
}

func testDirectorydStorageImpl(t *testing.T, store storage.DirectorydStorage) {
	hwid0 := "some_hwid_0"
	hwid1 := "some_hwid_1"
	hostname0 := "some_hostname_0"
	hostname1 := "some_hostname_1"

	// Empty initially
	_, err := store.GetHostname(hwid0)
	assert.Exactly(t, err, merrors.ErrNotFound)
	_, err = store.GetHostname(hwid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get hwid0->hostname1
	err = store.PutHostname(hwid0, hostname1)
	assert.NoError(t, err)
	recvd, err := store.GetHostname(hwid0)
	assert.Equal(t, hostname1, recvd)
	_, err = store.GetHostname(hwid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get hwid0->hostname0
	err = store.PutHostname(hwid0, hostname0)
	assert.NoError(t, err)
	recvd, err = store.GetHostname(hwid0)
	assert.Equal(t, hostname0, recvd)
	_, err = store.GetHostname(hwid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get hwid1->hostname1
	err = store.PutHostname(hwid1, hostname1)
	assert.NoError(t, err)
	recvd, err = store.GetHostname(hwid0)
	assert.Equal(t, hostname0, recvd)
	recvd, err = store.GetHostname(hwid1)
	assert.Equal(t, hostname1, recvd)
}
