/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage_test

import (
	"database/sql"
	"testing"

	"magma/orc8r/cloud/go/services/config/storage"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// Integration test for Sql implementation of config service storage which
// does some basic workflow tests on an in-memory sqlite3 DB
// Note this test does not run on sandcastle
func TestSqlConfigStorage_Integration(t *testing.T) {
	// Use an in-memory sqlite datastore
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Could not initialize sqlite DB: %s", err)
	}

	// Check the contract for an empty datastore
	store := storage.NewSqlConfigurationStorage(db)
	networkKeys, err := store.ListKeysForType("network", "network")
	assert.NoError(t, err)
	assert.Equal(t, []string{}, networkKeys)

	configs, err := store.GetConfigs("network", &storage.FilterCriteria{Type: "type"})
	assert.NoError(t, err)
	assert.Equal(t, map[storage.TypeAndKey]*storage.ConfigValue{}, configs)

	config, err := store.GetConfig("network", "type", "key")
	assert.NoError(t, err)
	assert.Equal(t, &storage.ConfigValue{}, config)

	// Create configs on 2 networks
	// network1: (type1, type2) X (key1, key2)
	err = store.CreateConfig("network1", "type1", "key1", []byte("value1"))
	assert.NoError(t, err)
	err = store.CreateConfig("network1", "type1", "key2", []byte("value2"))
	assert.NoError(t, err)
	err = store.CreateConfig("network1", "type2", "key1", []byte("value3"))
	assert.NoError(t, err)
	err = store.CreateConfig("network1", "type2", "key2", []byte("value4"))
	assert.NoError(t, err)

	// network2: (type3) X (key3, key4)
	err = store.CreateConfig("network2", "type3", "key3", []byte("value5"))
	assert.NoError(t, err)
	err = store.CreateConfig("network2", "type3", "key4", []byte("value6"))
	assert.NoError(t, err)

	// Read tests
	keys, err := store.ListKeysForType("network1", "type1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"key1", "key2"}, keys)

	configs, err = store.GetConfigs("network1", &storage.FilterCriteria{Type: "type1"})
	assert.NoError(t, err)
	expectedConfigs := map[storage.TypeAndKey]*storage.ConfigValue{
		{Type: "type1", Key: "key1"}: {Value: []byte("value1"), Version: 0},
		{Type: "type1", Key: "key2"}: {Value: []byte("value2"), Version: 0},
	}
	assert.Equal(t, expectedConfigs, configs)

	configs, err = store.GetConfigs("network1", &storage.FilterCriteria{Key: "key1"})
	assert.NoError(t, err)
	expectedConfigs = map[storage.TypeAndKey]*storage.ConfigValue{
		{Type: "type1", Key: "key1"}: {Value: []byte("value1"), Version: 0},
		{Type: "type2", Key: "key1"}: {Value: []byte("value3"), Version: 0},
	}
	assert.Equal(t, expectedConfigs, configs)

	configs, err = store.GetConfigs("network2", &storage.FilterCriteria{Type: "type3", Key: "key3"})
	assert.NoError(t, err)
	expectedConfigs = map[storage.TypeAndKey]*storage.ConfigValue{
		{Type: "type3", Key: "key3"}: {Value: []byte("value5"), Version: 0},
	}
	assert.Equal(t, expectedConfigs, configs)

	config, err = store.GetConfig("network2", "type3", "key4")
	assert.NoError(t, err)
	assert.Equal(t, &storage.ConfigValue{Value: []byte("value6"), Version: 0}, config)

	// Update-read
	err = store.UpdateConfig("network2", "type3", "key3", []byte("newValue"))
	assert.NoError(t, err)
	config, err = store.GetConfig("network2", "type3", "key3")
	assert.Equal(t, &storage.ConfigValue{Value: []byte("newValue"), Version: 1}, config)

	// Delete single
	err = store.DeleteConfig("network2", "type3", "key4")
	assert.NoError(t, err)
	config, err = store.GetConfig("network2", "type3", "key4")
	assert.NoError(t, err)
	assert.Equal(t, &storage.ConfigValue{}, config)

	// Delete multiple
	err = store.DeleteConfigs("network1", &storage.FilterCriteria{Type: "type1"})
	assert.NoError(t, err)
	keys, err = store.ListKeysForType("network1", "type1")
	assert.NoError(t, err)
	assert.Equal(t, []string{}, keys)

	err = store.DeleteConfigs("network1", &storage.FilterCriteria{Key: "key2"})
	assert.NoError(t, err)
	keys, err = store.ListKeysForType("network1", "type2")
	assert.NoError(t, err)
	assert.Equal(t, []string{"key1"}, keys)

	// Drop the table
	err = store.DeleteConfigsForNetwork("network2")
	assert.NoError(t, err)
}
