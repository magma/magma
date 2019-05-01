/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package storage_test

import (
	"context"
	"testing"

	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sql_utils"
	storage2 "magma/orc8r/cloud/go/storage"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestSqlConfiguratorStorage_Integration(t *testing.T) {
	db, err := sql_utils.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Could not initialize sqlite DB: %s", err)
	}
	factory := storage.NewSQLConfiguratorStorageFactory(db)
	err = factory.InitializeServiceStorage()
	assert.NoError(t, err)

	// Check the contract for an empty datastore
	store, err := factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)

	loadNetworksActual, err := store.LoadNetworks([]string{"n1", "n2"}, storage.FullNetworkLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.NetworkLoadResult{
			Networks:           []storage.Network{},
			NetworkIDsNotFound: []string{"n1", "n2"},
		},
		loadNetworksActual,
	)
	err = store.Commit()
	assert.NoError(t, err)

	// Create networks, read them back
	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)

	expectedn1 := storage.Network{ID: "n1", Name: "Network 1", Description: "foo", Configs: map[string][]byte{"hello": []byte("world"), "goodbye": []byte("alsoworld")}}
	actualNetwork, err := store.CreateNetwork(expectedn1)
	assert.NoError(t, err)
	assert.Equal(t, expectedn1, actualNetwork)

	expectedn2 := storage.Network{ID: "n2"}
	actualNetwork, err = store.CreateNetwork(expectedn2)
	assert.NoError(t, err)
	assert.Equal(t, expectedn2, actualNetwork)
	expectedn2.Configs = map[string][]byte{}

	err = store.Commit()
	assert.NoError(t, err)

	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	loadNetworksActual, err = store.LoadNetworks([]string{"n1", "n2", "n3"}, storage.FullNetworkLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.NetworkLoadResult{
			Networks:           []storage.Network{expectedn1, expectedn2},
			NetworkIDsNotFound: []string{"n3"},
		},
		loadNetworksActual,
	)
	err = store.Commit()
	assert.NoError(t, err)

	// Update networks, read them back
	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	_, err = store.CreateNetwork(storage.Network{ID: "n3"})
	assert.NoError(t, err)

	newNames := []string{"New Network 1"}
	newDescs := []string{"New Network 1 description"}
	updates := []storage.NetworkUpdateCriteria{
		{ID: "n1", NewName: &newNames[0], NewDescription: &newDescs[0], ConfigsToDelete: []string{"goodbye"}, ConfigsToAddOrUpdate: map[string][]byte{"foo": []byte("bar")}},
		{ID: "n2", ConfigsToDelete: []string{"dne"}, ConfigsToAddOrUpdate: map[string][]byte{"baz": []byte("quz")}},
		{ID: "n3", DeleteNetwork: true},
		{ID: "n4", DeleteNetwork: true},
	}

	expectedn1.Name, expectedn1.Description = newNames[0], newDescs[0]
	delete(expectedn1.Configs, "goodbye")
	expectedn1.Configs["foo"] = []byte("bar")
	expectedn1.Version = 1
	expectedn2.Configs = map[string][]byte{"baz": []byte("quz")}
	expectedn2.Version = 1

	failures, err := store.UpdateNetworks(updates)
	assert.NoError(t, err)
	assert.Equal(t, storage.FailedOperations{}, failures)
	assert.NoError(t, store.Commit())

	store, err = factory.StartTransaction(context.Background(), &storage.TxOptions{ReadOnly: true})
	assert.NoError(t, err)
	loadNetworksActual, err = store.LoadNetworks([]string{"n1", "n2", "n3"}, storage.FullNetworkLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.NetworkLoadResult{
			Networks:           []storage.Network{expectedn1, expectedn2},
			NetworkIDsNotFound: []string{"n3"},
		},
		loadNetworksActual,
	)
	assert.NoError(t, store.Commit())

	// Empty datastore contract for entities
	store, err = factory.StartTransaction(context.Background(), nil)

	actualEntityLoad, err := store.LoadEntities("n1", storage.EntityLoadFilter{}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities:         []storage.NetworkEntity{},
			EntitiesNotFound: []storage2.TypeAndKey{},
		},
		actualEntityLoad,
	)

	actualEntityLoad, err = store.LoadEntities(
		"n1",
		storage.EntityLoadFilter{
			IDs: []storage2.TypeAndKey{
				{Type: "foo", Key: "bar"},
				{Type: "baz", Key: "quz"},
			},
		},
		storage.EntityLoadCriteria{},
	)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities: []storage.NetworkEntity{},
			EntitiesNotFound: []storage2.TypeAndKey{
				{Type: "foo", Key: "bar"},
				{Type: "baz", Key: "quz"},
			},
		},
		actualEntityLoad,
	)
	assert.NoError(t, store.Commit())
}
