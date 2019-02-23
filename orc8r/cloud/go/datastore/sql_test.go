/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package datastore_test

import (
	"testing"

	"magma/orc8r/cloud/go/datastore"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestDatastoreBasics(t *testing.T) {
	key := "magic"
	value := []byte("Hello world!")
	table := "test_table"
	// Create an in-memory sqlite datastore for testing
	ds, err := datastore.NewSqlDb("sqlite3", ":memory:")
	assert.NoError(t, err)

	// Add the key to the datastore
	err = ds.Put(table, key, value)
	assert.NoError(t, err)

	exists, err := ds.DoesKeyExist(table, key)
	assert.NoError(t, err)
	assert.True(t, exists)
	exists, err = ds.DoesKeyExist(table, "SHOULD NOT EXIST")
	assert.NoError(t, err)
	assert.False(t, exists)

	res, _, err := ds.Get(table, key)
	assert.NoError(t, err)
	assert.Equal(t, value, res)

	keys, err := ds.ListKeys(table)
	assert.NoError(t, err)
	assert.Equal(t, []string{key}, keys)

	// Delete the key and check the datastore
	err = ds.Delete(table, key)
	assert.NoError(t, err)
	_, _, err = ds.Get(table, key)
	assert.Error(t, err) // key missing now

	keys, err = ds.ListKeys(table)
	assert.NoError(t, err)
	assert.Equal(t, []string{}, keys)

	err = ds.DeleteTable(table)
	assert.NoError(t, err)
}

func TestDatastoreBulkOperations(t *testing.T) {
	ds, err := datastore.NewSqlDb("sqlite3", ":memory:")
	assert.NoError(t, err)

	// Bulk insert KV's, no updates
	valuesToPut := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}

	expectedFailedKeys := make(map[string]error, 0)
	failedKeys, err := ds.PutMany("test", valuesToPut)
	assert.NoError(t, err)
	assert.Equal(t, expectedFailedKeys, failedKeys)

	dbRows, err := ds.GetMany("test", []string{})
	assert.NoError(t, err)
	assert.Equal(t, map[string]datastore.ValueWrapper{}, dbRows)

	dbRows, err = ds.GetMany("test", []string{"key1", "key2"})
	assert.NoError(t, err)
	expectedDbRows := map[string]datastore.ValueWrapper{
		"key1": {
			Value:      []byte("value1"),
			Generation: 0,
		},
		"key2": {
			Value:      []byte("value2"),
			Generation: 0,
		},
	}
	assert.Equal(t, expectedDbRows, dbRows)

	// PutAll with 1 update and 1 insert
	valuesToPut = map[string][]byte{
		"key2": []byte("newvalue2"),
		"key3": []byte("value3"),
	}
	failedKeys, err = ds.PutMany("test", valuesToPut)
	assert.NoError(t, err)
	assert.Equal(t, expectedFailedKeys, failedKeys)

	dbRows, err = ds.GetMany("test", []string{"key1", "key2", "key3"})
	assert.NoError(t, err)
	expectedDbRows = map[string]datastore.ValueWrapper{
		"key1": {
			Value:      []byte("value1"),
			Generation: 0,
		},
		"key2": {
			Value:      []byte("newvalue2"),
			Generation: 1,
		},
		"key3": {
			Value:      []byte("value3"),
			Generation: 0,
		},
	}
	assert.Equal(t, expectedDbRows, dbRows)

	// Empty PutAll
	failedKeys, err = ds.PutMany("test", map[string][]byte{})
	assert.NoError(t, err)
	assert.Equal(t, expectedFailedKeys, failedKeys)

	dbRows, err = ds.GetMany("test", []string{"key1", "key2", "key3"})
	assert.NoError(t, err)
	assert.Equal(t, expectedDbRows, dbRows)

	// Empty GetAll
	emptyDbRows, err := ds.GetMany("test", []string{})
	assert.NoError(t, err)
	assert.Equal(t, map[string]datastore.ValueWrapper{}, emptyDbRows)

	// Delete many
	failedKeys, err = ds.DeleteMany("test", []string{"key1", "key2"})
	assert.NoError(t, err)
	assert.Equal(t, expectedFailedKeys, failedKeys)
	expectedDbRows = map[string]datastore.ValueWrapper{
		"key3": {
			Value:      []byte("value3"),
			Generation: 0,
		},
	}
	dbRows, err = ds.GetMany("test", []string{"key1", "key2", "key3"})
	assert.NoError(t, err)
	assert.Equal(t, expectedDbRows, dbRows)

}
