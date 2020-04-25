/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package blobstore_test

import (
	"bytes"
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/errors"

	"github.com/stretchr/testify/assert"
)

// checks for equality but ignores the version field
func blobEqual(b1 blobstore.Blob, b2 blobstore.Blob) bool {
	return b1.Type == b2.Type && b1.Key == b2.Key &&
		bytes.Equal(b1.Value, b2.Value)
}

func TestMemoryBlobStorageStorage_CreateOrUpdate(t *testing.T) {
	factory := blobstore.NewMemoryBlobStorageFactory()
	type1 := "type1"
	key1 := "key1"
	blob1 := blobstore.Blob{Type: type1, Key: key1, Value: []byte("value1")}
	blob2 := blobstore.Blob{Type: type1, Key: key1, Value: []byte("value2")}
	network1 := "network1"

	store, err := factory.StartTransaction(nil)
	assert.NoError(t, err)

	// create
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob1}))

	blob, err := store.Get(network1, storage.TypeAndKey{Type: type1, Key: key1})
	assert.Equal(t, err, nil)
	assert.True(t, blobEqual(blob1, blob))
	version1 := blob.Version

	// update
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob2}))
	blob, err = store.Get(network1, storage.TypeAndKey{Type: type1, Key: key1})
	assert.Equal(t, err, nil)
	assert.True(t, blobEqual(blob2, blob))
	version2 := blob.Version

	assert.True(t, version2 > version1)

	// update blob version
	blob1.Version = 10
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob1}))
	blob, err = store.Get(network1, storage.TypeAndKey{Type: type1, Key: key1})
	assert.Equal(t, err, nil)
	assert.True(t, blobEqual(blob1, blob))

	assert.True(t, blob.Version > version2)
}

func TestMemoryBlobStorage_Rollback(t *testing.T) {
	factory := blobstore.NewMemoryBlobStorageFactory()
	type1 := "type1"
	key1 := "key1"
	blob1 := blobstore.Blob{Type: type1, Key: key1, Value: []byte("value1")}
	network1 := "network1"

	store, err := factory.StartTransaction(nil)
	assert.NoError(t, err)

	// create
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob1}))

	_, err = store.Get(network1, storage.TypeAndKey{Type: type1, Key: key1})
	assert.Equal(t, err, nil)

	store.Rollback()

	store, err = factory.StartTransaction(nil)
	assert.NoError(t, err)
	_, err = store.Get(network1, storage.TypeAndKey{Type: type1, Key: key1})
	assert.Equal(t, errors.ErrNotFound, err)
}

func TestMemoryBlobStorage_Commit(t *testing.T) {
	factory := blobstore.NewMemoryBlobStorageFactory()
	type1 := "type1"
	key1 := "key1"
	blob1 := blobstore.Blob{Type: type1, Key: key1, Value: []byte("value1")}
	network1 := "network1"
	id1 := storage.TypeAndKey{Type: type1, Key: key1}

	store, err := factory.StartTransaction(nil)
	assert.NoError(t, err)

	// create
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob1}))

	_, err = store.Get(network1, id1)
	assert.Equal(t, err, nil)

	assert.NoError(t, store.Commit())

	store, err = factory.StartTransaction(nil)
	assert.NoError(t, err)
	blob, err := store.Get(network1, id1)
	assert.NoError(t, err)
	assert.True(t, blobEqual(blob, blob1))
	blob1.Value = []byte("value2")
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob1}))
	assert.NoError(t, store.Commit())

	store, err = factory.StartTransaction(nil)
	assert.NoError(t, err)
	blob, err = store.Get(network1, id1)
	assert.NoError(t, err)
	assert.Equal(t, []byte("value2"), blob.Value)
}

func TestMemoryBlobStorage_GetMany(t *testing.T) {
	factory := blobstore.NewMemoryBlobStorageFactory()
	type1 := "type1"
	key1 := "key1"
	type2 := "type2"
	key2 := "key2"
	blob1 := blobstore.Blob{Type: type1, Key: key1, Value: []byte("value1")}
	blob2 := blobstore.Blob{Type: type2, Key: key2, Value: []byte("value2")}
	network1 := "network1"
	ids := []storage.TypeAndKey{
		{Type: blob1.Type, Key: blob1.Key},
		{Type: blob2.Type, Key: blob2.Key},
	}

	store, err := factory.StartTransaction(nil)
	assert.NoError(t, err)

	// create
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob1}))
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob2}))

	// lookup
	blobs, err := store.GetMany(network1, ids)
	assert.Equal(t, 2, len(blobs))
	assert.True(t, blobEqual(blobs[0], blob1) || blobEqual(blobs[1], blob1))
	assert.True(t, blobEqual(blobs[0], blob2) || blobEqual(blobs[1], blob2))

	// try to look up a non existent blob, this should not fail
	ids = append(ids, storage.TypeAndKey{Type: "type3", Key: "key3"})
	blobs, err = store.GetMany(network1, ids)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(blobs))
}

func TestMemoryBlobStorage_Delete(t *testing.T) {
	factory := blobstore.NewMemoryBlobStorageFactory()
	type1 := "type1"
	key1 := "key1"
	blob1 := blobstore.Blob{Type: type1, Key: key1, Value: []byte("value1")}
	network1 := "network1"
	ids := []storage.TypeAndKey{
		{Type: blob1.Type, Key: blob1.Key},
	}

	store, err := factory.StartTransaction(nil)
	assert.NoError(t, err)

	// create, update, delete within a transaction session should end up with
	// the blob being deleted
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob1}))
	blob1.Value = []byte("value1_updated")
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob1}))
	assert.NoError(t, store.Delete(network1, ids))
	assert.NoError(t, store.Commit())
	store, err = factory.StartTransaction(nil)
	assert.NoError(t, err)
	blobs, _ := store.GetMany(network1, ids)
	assert.Equal(t, 0, len(blobs))

	// create, delete, update within a transaction session should end up with
	// the blob being deleted
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob1}))
	blob1.Value = []byte("value1_updated")
	assert.NoError(t, store.Delete(network1, ids))
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob1}))
	assert.NoError(t, store.Commit())
	store, err = factory.StartTransaction(nil)
	assert.NoError(t, err)
	blobs, _ = store.GetMany(network1, ids)
	assert.Equal(t, 1, len(blobs))
}

func TestMemoryBlobStorage_ListKeys(t *testing.T) {
	factory := blobstore.NewMemoryBlobStorageFactory()
	type1 := "type"
	key1 := "key1"
	key2 := "key2"
	blob1 := blobstore.Blob{Type: type1, Key: key1, Value: []byte("value1")}
	blob2 := blobstore.Blob{Type: type1, Key: key2, Value: []byte("value2")}
	network1 := "network1"

	store, err := factory.StartTransaction(nil)
	assert.NoError(t, err)

	// Test local changes
	assert.NoError(t, store.CreateOrUpdate(network1, []blobstore.Blob{blob1, blob2}))
	keys, err := store.ListKeys(network1, type1)
	assert.Equal(t, []string{key1, key2}, keys)

	assert.NoError(t, store.Commit())
	store, err = factory.StartTransaction(nil)
	assert.NoError(t, err)

	// Test committed changes
	keys, err = store.ListKeys(network1, type1)
	assert.Equal(t, []string{key1, key2}, keys)

	// Test locally deleted changes
	assert.NoError(t, store.Delete(network1, []storage.TypeAndKey{{type1, key1}}))
	keys, err = store.ListKeys(network1, type1)
	assert.Equal(t, []string{key2}, keys)
}

func TestMemoryBlobStorageStorage_IncrementVersion(t *testing.T) {
	factory := blobstore.NewMemoryBlobStorageFactory()
	type1 := "type1"
	key1 := "key1"
	network1 := "network1"

	typeAndKey := storage.TypeAndKey{Type: type1, Key: key1}

	store, err := factory.StartTransaction(nil)
	assert.NoError(t, err)

	// increment non-existing blob
	err = store.IncrementVersion(network1, typeAndKey)
	assert.NoError(t, err)
	blob, err := store.Get(network1, typeAndKey)
	assert.NoError(t, err)

	assert.Equal(t, blob.Version, uint64(1))

	// increment version
	err = store.IncrementVersion(network1, typeAndKey)
	assert.NoError(t, err)
	blob, err = store.Get(network1, typeAndKey)
	assert.NoError(t, err)

	assert.Equal(t, blob.Version, uint64(2))
}

func TestMemoryBlobStorage_Integration(t *testing.T) {
	fact := blobstore.NewMemoryBlobStorageFactory()
	integration(t, fact)
}
