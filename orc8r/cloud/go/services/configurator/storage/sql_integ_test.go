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
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sql_utils"
	storage2 "magma/orc8r/cloud/go/storage"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type mockIDGenerator struct {
	count int
}

func (g *mockIDGenerator) New() string {
	g.count++
	return fmt.Sprintf("%d", g.count)
}

func TestSqlConfiguratorStorage_Integration(t *testing.T) {
	db, err := sql_utils.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Could not initialize sqlite DB: %s", err)
	}
	factory := storage.NewSQLConfiguratorStorageFactory(db, &mockIDGenerator{})
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

	// ========================================================================
	// Create/load networks tests
	// ========================================================================

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

	// ========================================================================
	// Update network tests
	// ========================================================================

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

	// ========================================================================
	// Empty datastore entity load tests
	// ========================================================================

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

	// ========================================================================
	// Create/Load entity tests
	// ========================================================================

	// Create 3 entities, read them back
	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	expectedFoobarEnt, err := store.CreateEntity("n1", storage.NetworkEntity{
		Type: "foo",
		Key:  "bar",

		Name:        "foobar",
		Description: "foobar ent",

		PhysicalID: "1",

		Config: []byte("foobar"),

		// should be ignored
		GraphID: "1",
	})
	assert.NoError(t, err)
	assert.Equal(t, "2", expectedFoobarEnt.GraphID)

	expectedBarbazEnt, err := store.CreateEntity("n1", storage.NetworkEntity{
		Type: "bar",
		Key:  "baz",

		Name:        "barbaz",
		Description: "barbaz ent",

		Config: []byte("barbaz"),

		Permissions: []storage.ACL{
			{Permission: storage.NoPermissions, Scope: storage.WildcardACLScope, Type: storage.WildcardACLType},
			{Permission: storage.WritePermission, Scope: storage.ACLScope{NetworkIDs: []string{"n1"}}, Type: storage.ACLType{EntityType: "foo"}},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "4", expectedBarbazEnt.GraphID)
	assert.Equal(t, "5", expectedBarbazEnt.Permissions[0].ID)
	assert.Equal(t, "6", expectedBarbazEnt.Permissions[1].ID)

	// bazquz should link foobar and barbaz into 1 graph
	// that graph ID should be 2
	expectedBazquzEnt, err := store.CreateEntity("n1", storage.NetworkEntity{
		Type: "baz",
		Key:  "quz",

		Name:        "bazquz",
		Description: "bazquz ent",

		Associations: []storage2.TypeAndKey{
			{Type: "foo", Key: "bar"},
			{Type: "bar", Key: "baz"},
		},
	})
	assert.NoError(t, err)
	assert.NoError(t, store.Commit())
	assert.Equal(t, "2", expectedBazquzEnt.GraphID)

	expectedFoobarEnt.GraphID = "2"
	expectedBarbazEnt.GraphID = "2"

	expectedFoobarEnt.ParentAssociations = []storage2.TypeAndKey{
		{Type: "baz", Key: "quz"},
	}
	expectedBarbazEnt.ParentAssociations = []storage2.TypeAndKey{
		{Type: "baz", Key: "quz"},
	}

	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	actualEntityLoad, err = store.LoadEntities("n1", storage.EntityLoadFilter{}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities: []storage.NetworkEntity{
				expectedBarbazEnt,
				expectedBazquzEnt,
				expectedFoobarEnt,
			},
			EntitiesNotFound: []storage2.TypeAndKey{},
		},
		actualEntityLoad,
	)
	assert.NoError(t, store.Commit())

	// At this point, our graph looks like this:
	//                (baz, quz)
	//                 /      \
	//                /        \
	//    (foo, bar) <          > (bar, baz)

	// ========================================================================
	// Update entity tests
	// ========================================================================

	// create a new ent helloworld without assocs
	// pk should be 9, graph id should be 10, permission id should be 11

	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	_, err = store.CreateEntity("n1", storage.NetworkEntity{
		Type: "hello",
		Key:  "world",

		Name:        "helloworld",
		Description: "helloworld ent",

		Config: []byte("first config"),

		Permissions: []storage.ACL{
			{Permission: storage.NoPermissions, Scope: storage.WildcardACLScope, Type: storage.WildcardACLType},
		},
	})
	assert.NoError(t, err)

	// update basic fields and permissions on it

	newName := "helloworld2"
	newDesc := "helloworld2 ent"
	newPhysID := "asdf"
	newConfig := []byte("second config")
	updateHelloWorldEntResult, err := store.UpdateEntity("n1", storage.EntityUpdateCriteria{
		Type: "hello",
		Key:  "world",

		NewName:        &newName,
		NewDescription: &newDesc,
		NewPhysicalID:  &newPhysID,

		NewConfig: &newConfig,

		PermissionsToCreate: []storage.ACL{
			{Permission: storage.WritePermission, Scope: storage.ACLScopeOf([]string{"n1"}), Type: storage.ACLTypeOf("foo")},
		},
		PermissionsToUpdate: []storage.ACL{
			{ID: "11", Permission: storage.WritePermission, Scope: storage.WildcardACLScope, Type: storage.WildcardACLType},
		},
	})
	assert.NoError(t, err)

	assert.Equal(
		t,
		storage.NetworkEntity{
			Type: "hello",
			Key:  "world",

			GraphID: "10",

			PhysicalID: newPhysID,

			Name:        newName,
			Description: newDesc,
			Config:      newConfig,

			Permissions: []storage.ACL{
				{ID: "12", Permission: storage.WritePermission, Scope: storage.ACLScopeOf([]string{"n1"}), Type: storage.ACLTypeOf("foo")},
				{ID: "11", Permission: storage.WritePermission, Scope: storage.WildcardACLScope, Type: storage.WildcardACLType},
			},

			Version: 1,
		},
		updateHelloWorldEntResult,
	)
	expectedHelloWorldEnt := updateHelloWorldEntResult

	// create assocs to each of the previous 3 ents, helloworld's graph ID should be 2
	// At this point, the graph should look like
	//                                 (hello, world)
	//                               /       |        \
	//                              /        |         \
	//                             /         |          \
	//                            /          |           \
	//                           /           v            \
	//                          /        (baz, quz)        \
	//                         /            /  \            \
	//                        |            /    \            |
	//                        |           /      \           |
	//                        v          /        \          v
	//                      (foo, bar)  <          >  (bar, baz)
	updateHelloWorldEntResult, err = store.UpdateEntity("n1", storage.EntityUpdateCriteria{
		Type: "hello",
		Key:  "world",

		AssociationsToAdd: []storage2.TypeAndKey{
			{Type: "foo", Key: "bar"},
			{Type: "bar", Key: "baz"},
			{Type: "baz", Key: "quz"},
		},
	})
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.NetworkEntity{
			Type:       "hello",
			Key:        "world",
			GraphID:    "2",
			PhysicalID: newPhysID,
			Associations: []storage2.TypeAndKey{
				{Type: "foo", Key: "bar"},
				{Type: "bar", Key: "baz"},
				{Type: "baz", Key: "quz"},
			},
			Version: 2,
		},
		updateHelloWorldEntResult,
	)

	// Read back the updated ent
	expectedHelloWorldEnt.Associations = []storage2.TypeAndKey{
		{Type: "foo", Key: "bar"},
		{Type: "bar", Key: "baz"},
		{Type: "baz", Key: "quz"},
	}
	expectedHelloWorldEnt.GraphID = "2"
	expectedHelloWorldEnt.Permissions = []storage.ACL{
		{ID: "11", Permission: storage.WritePermission, Scope: storage.WildcardACLScope, Type: storage.WildcardACLType, Version: 1},
		{ID: "12", Permission: storage.WritePermission, Scope: storage.ACLScopeOf([]string{"n1"}), Type: storage.ACLTypeOf("foo")},
	}
	expectedHelloWorldEnt.Version = 2

	actualEntityLoad, err = store.LoadEntities(
		"n1",
		storage.EntityLoadFilter{IDs: []storage2.TypeAndKey{{Type: "hello", Key: "world"}}},
		storage.FullEntityLoadCriteria,
	)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities:         []storage.NetworkEntity{expectedHelloWorldEnt},
			EntitiesNotFound: []storage2.TypeAndKey{},
		},
		actualEntityLoad,
	)

	// TODO: orphan some graph nodes (blocked on impl of fixGraph)

	// delete assoc to foobar, helloworld's graph ID should still be 2
}
