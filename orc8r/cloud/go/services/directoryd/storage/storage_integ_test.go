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
	hwid2 := "some_hwid_2"
	hwid3 := "some_hwid_3"
	hostname0 := "some_hostname_0"
	hostname1 := "some_hostname_1"
	hostname2 := "some_hostname_2"
	hostname3 := "some_hostname_3"

	nid0 := "some_networkid_0"
	nid1 := "some_networkid_1"
	sid0 := "some_sessionid_0"
	sid1 := "some_sessionid_1"
	imsi0 := "some_imsi_0"
	imsi1 := "some_imsi_1"

	//////////////////////////////
	// Hostname -> HWID
	//////////////////////////////

	// Empty initially
	_, err := store.GetHostnameForHWID(hwid0)
	assert.Exactly(t, err, merrors.ErrNotFound)
	_, err = store.GetHostnameForHWID(hwid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get hwid0->hostname1
	err = store.MapHWIDsToHostnames(map[string]string{hwid0: hostname1})
	assert.NoError(t, err)
	recvd, err := store.GetHostnameForHWID(hwid0)
	assert.Equal(t, hostname1, recvd)
	_, err = store.GetHostnameForHWID(hwid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get hwid0->hostname0
	err = store.MapHWIDsToHostnames(map[string]string{hwid0: hostname0})
	assert.NoError(t, err)
	recvd, err = store.GetHostnameForHWID(hwid0)
	assert.NoError(t, err)
	assert.Equal(t, hostname0, recvd)
	_, err = store.GetHostnameForHWID(hwid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get hwid1->hostname1
	err = store.MapHWIDsToHostnames(map[string]string{hwid1: hostname1})
	assert.NoError(t, err)
	recvd, err = store.GetHostnameForHWID(hwid0)
	assert.NoError(t, err)
	assert.Equal(t, hostname0, recvd)
	recvd, err = store.GetHostnameForHWID(hwid1)
	assert.NoError(t, err)
	assert.Equal(t, hostname1, recvd)

	// Multi-put: Put and Get hwid2->hostname2, hwid3->hostname3
	err = store.MapHWIDsToHostnames(map[string]string{hwid2: hostname2, hwid3: hostname3})
	assert.NoError(t, err)
	recvd, err = store.GetHostnameForHWID(hwid2)
	assert.NoError(t, err)
	assert.Equal(t, hostname2, recvd)
	recvd, err = store.GetHostnameForHWID(hwid3)
	assert.NoError(t, err)
	assert.Equal(t, hostname3, recvd)

	//////////////////////////////
	// Session ID -> IMSI
	//////////////////////////////

	// Empty initially
	_, err = store.GetIMSIForSessionID(nid0, sid0)
	assert.Exactly(t, err, merrors.ErrNotFound)
	_, err = store.GetIMSIForSessionID(nid0, sid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get sid0->imsi1
	err = store.MapSessionIDsToIMSIs(nid0, map[string]string{sid0: imsi1})
	assert.NoError(t, err)
	recvd, err = store.GetIMSIForSessionID(nid0, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, recvd)
	_, err = store.GetIMSIForSessionID(nid0, sid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get sid0->imsi0
	err = store.MapSessionIDsToIMSIs(nid0, map[string]string{sid0: imsi0})
	assert.NoError(t, err)
	recvd, err = store.GetIMSIForSessionID(nid0, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, recvd)
	_, err = store.GetIMSIForSessionID(nid0, sid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get sid1->imsi1
	err = store.MapSessionIDsToIMSIs(nid0, map[string]string{sid1: imsi1})
	assert.NoError(t, err)
	recvd, err = store.GetIMSIForSessionID(nid0, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, recvd)
	recvd, err = store.GetIMSIForSessionID(nid0, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, recvd)

	// Multi-put: Put and Get sid0->imsi0, sid1->imsi1 for nid1
	err = store.MapSessionIDsToIMSIs(nid1, map[string]string{sid0: imsi0, sid1: imsi1})
	assert.NoError(t, err)
	recvd, err = store.GetIMSIForSessionID(nid1, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, recvd)
	recvd, err = store.GetIMSIForSessionID(nid1, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, recvd)

	// Correctly network-partitioned: {nid0: sid0->imsi0, nid1: sid0->imsi1}
	err = store.MapSessionIDsToIMSIs(nid0, map[string]string{sid0: imsi0})
	assert.NoError(t, err)
	err = store.MapSessionIDsToIMSIs(nid1, map[string]string{sid0: imsi1})
	assert.NoError(t, err)
	recvd, err = store.GetIMSIForSessionID(nid0, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, recvd)
	recvd, err = store.GetIMSIForSessionID(nid1, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, recvd)
}
