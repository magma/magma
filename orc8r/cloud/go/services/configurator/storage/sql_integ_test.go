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
	"magma/orc8r/cloud/go/sqorc"
	orc8rStorage "magma/orc8r/cloud/go/storage"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes/wrappers"
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
	// sqlite's default behavior is to disable foreign keys (wat): https://www.sqlite.org/draft/pragma.html#pragma_foreign_keys
	// thankfully the sqlite3 driver supports the apporpriate pragma: https://github.com/mattn/go-sqlite3/issues/255
	db, err := sqorc.Open("sqlite3", ":memory:?_foreign_keys=1")
	if err != nil {
		t.Fatalf("Could not initialize sqlite DB: %s", err)
	}
	factory := storage.NewSQLConfiguratorStorageFactory(db, &mockIDGenerator{}, sqorc.GetSqlBuilder())
	err = factory.InitializeServiceStorage()
	assert.NoError(t, err)

	// Check the contract for an empty data store
	store, err := factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)

	loadNetworksActual, err := store.LoadNetworks(storage.NetworkLoadFilter{Ids: []string{"n1", "n2"}}, storage.FullNetworkLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.NetworkLoadResult{
			Networks:           []*storage.Network{},
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

	expectedn1 := storage.Network{ID: "n1", Type: "type1", Name: "Network 1", Description: "foo", Configs: map[string][]byte{"hello": []byte("world"), "goodbye": []byte("alsoworld")}}
	actualNetwork, err := store.CreateNetwork(expectedn1)
	assert.NoError(t, err)
	assert.Equal(t, expectedn1, actualNetwork)

	expectedn2 := storage.Network{ID: "n2", Type: "type2"}
	actualNetwork, err = store.CreateNetwork(expectedn2)
	assert.NoError(t, err)
	assert.Equal(t, expectedn2, actualNetwork)
	expectedn2.Configs = map[string][]byte{}

	err = store.Commit()
	assert.NoError(t, err)

	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	loadNetworksActual, err = store.LoadNetworks(storage.NetworkLoadFilter{Ids: []string{"n1", "n2", "n3"}}, storage.FullNetworkLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.NetworkLoadResult{
			Networks:           []*storage.Network{&expectedn1, &expectedn2},
			NetworkIDsNotFound: []string{"n3"},
		},
		loadNetworksActual,
	)
	err = store.Commit()
	assert.NoError(t, err)

	// ========================================================================
	// LoadAll networks tests
	// ========================================================================

	store, err = factory.StartTransaction(context.Background(), nil)
	loadedNetworks, err := store.LoadAllNetworks(storage.FullNetworkLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(t, "n1", loadedNetworks[0].ID)
	assert.Equal(t, "n2", loadedNetworks[1].ID)
	err = store.Commit()
	assert.NoError(t, err)

	// ========================================================================
	// Update network tests
	// ========================================================================

	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	_, err = store.CreateNetwork(storage.Network{ID: "n3"})
	assert.NoError(t, err)

	updates := []storage.NetworkUpdateCriteria{
		{ID: "n1", NewName: &wrappers.StringValue{Value: "New Network 1"}, NewDescription: &wrappers.StringValue{Value: "New Network 1 description"}, ConfigsToDelete: []string{"goodbye"}, ConfigsToAddOrUpdate: map[string][]byte{"foo": []byte("bar")}},
		{ID: "n2", ConfigsToDelete: []string{"dne"}, ConfigsToAddOrUpdate: map[string][]byte{"baz": []byte("quz")}},
		{ID: "n3", DeleteNetwork: true},
		{ID: "n4", DeleteNetwork: true},
	}

	expectedn1.Name, expectedn1.Description = "New Network 1", "New Network 1 description"
	delete(expectedn1.Configs, "goodbye")
	expectedn1.Configs["foo"] = []byte("bar")
	expectedn1.Version = 1
	expectedn2.Configs = map[string][]byte{"baz": []byte("quz")}
	expectedn2.Version = 1

	err = store.UpdateNetworks(updates)
	assert.NoError(t, err)
	assert.NoError(t, store.Commit())

	store, err = factory.StartTransaction(context.Background(), &orc8rStorage.TxOptions{ReadOnly: true})
	assert.NoError(t, err)
	loadNetworksActual, err = store.LoadNetworks(storage.NetworkLoadFilter{Ids: []string{"n1", "n2", "n3"}}, storage.FullNetworkLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.NetworkLoadResult{
			Networks:           []*storage.Network{&expectedn1, &expectedn2},
			NetworkIDsNotFound: []string{"n3"},
		},
		loadNetworksActual,
	)
	assert.NoError(t, store.Commit())

	// ========================================================================
	// Create and Load typed networks
	// ========================================================================

	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)

	expectedn3 := storage.Network{ID: "n3", Type: "type1"}
	actualNetwork, err = store.CreateNetwork(expectedn3)
	assert.NoError(t, err)
	assert.Equal(t, expectedn3, actualNetwork)

	expectedn4 := storage.Network{ID: "n4", Type: "type2"}
	actualNetwork, err = store.CreateNetwork(expectedn4)
	assert.NoError(t, err)
	assert.Equal(t, expectedn4, actualNetwork)

	expectedn2.Configs = map[string][]uint8{}
	expectedn2.Version = 1

	expectedn4.Configs = map[string][]uint8{}

	loadNetworksActual, err = store.LoadNetworks(storage.NetworkLoadFilter{TypeFilter: stringPointer("type2")}, storage.NetworkLoadCriteria{})
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.NetworkLoadResult{
			Networks:           []*storage.Network{&expectedn2, &expectedn4},
			NetworkIDsNotFound: []string{},
		},
		loadNetworksActual,
	)
	assert.NoError(t, store.Commit())

	// ========================================================================
	// Empty data store entity load tests
	// ========================================================================

	store, err = factory.StartTransaction(context.Background(), nil)

	actualEntityLoad, err := store.LoadEntities("n1", storage.EntityLoadFilter{}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities:         []*storage.NetworkEntity{},
			EntitiesNotFound: []*storage.EntityID{},
		},
		actualEntityLoad,
	)

	actualEntityLoad, err = store.LoadEntities(
		"n1",
		storage.EntityLoadFilter{
			IDs: []*storage.EntityID{
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
			Entities: []*storage.NetworkEntity{},
			EntitiesNotFound: []*storage.EntityID{
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

		Permissions: []*storage.ACL{
			{
				Permission: storage.ACL_NO_PERM,
				Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
				Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
			},
			{
				Permission: storage.ACL_WRITE,
				Type:       &storage.ACL_EntityType{EntityType: "foo"},
				Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n1"}}},
			},
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

		Associations: []*storage.EntityID{
			{Type: "bar", Key: "baz"},
			{Type: "foo", Key: "bar"},
		},
	})
	assert.NoError(t, err)
	assert.NoError(t, store.Commit())
	assert.Equal(t, "2", expectedBazquzEnt.GraphID)

	expectedFoobarEnt.GraphID = "2"
	expectedBarbazEnt.GraphID = "2"

	expectedFoobarEnt.ParentAssociations = []*storage.EntityID{
		{Type: "baz", Key: "quz"},
	}
	expectedBarbazEnt.ParentAssociations = []*storage.EntityID{
		{Type: "baz", Key: "quz"},
	}

	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	actualEntityLoad, err = store.LoadEntities("n1", storage.EntityLoadFilter{}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				&expectedBarbazEnt,
				&expectedBazquzEnt,
				&expectedFoobarEnt,
			},
			EntitiesNotFound: []*storage.EntityID{},
		},
		actualEntityLoad,
	)
	assert.NoError(t, store.Commit())

	// Load from physical ID
	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	// network ID shouldn't matter in this query
	actualEntityLoad, err = store.LoadEntities("placeholder", storage.EntityLoadFilter{PhysicalID: stringPointer("1")}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				&expectedFoobarEnt,
			},
			EntitiesNotFound: []*storage.EntityID{},
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

		Permissions: []*storage.ACL{
			{
				Permission: storage.ACL_NO_PERM,
				Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
				Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
			},
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

		NewName:        &wrappers.StringValue{Value: newName},
		NewDescription: &wrappers.StringValue{Value: newDesc},
		NewPhysicalID:  &wrappers.StringValue{Value: newPhysID},

		NewConfig: &wrappers.BytesValue{Value: newConfig},

		PermissionsToCreate: []*storage.ACL{
			{
				Permission: storage.ACL_WRITE,
				Type:       &storage.ACL_EntityType{EntityType: "foo"},
				Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n1"}}},
			},
		},
		PermissionsToUpdate: []*storage.ACL{
			{
				ID:         "11",
				Permission: storage.ACL_WRITE,
				Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
				Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
			},
		},
	})
	assert.NoError(t, err)

	assert.Equal(
		t,
		storage.NetworkEntity{
			NetworkID: "n1",
			Type:      "hello",
			Key:       "world",

			GraphID: "10",

			PhysicalID: newPhysID,

			Name:        newName,
			Description: newDesc,
			Config:      newConfig,

			Permissions: []*storage.ACL{
				{
					ID:         "12",
					Permission: storage.ACL_WRITE,
					Type:       &storage.ACL_EntityType{EntityType: "foo"},
					Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n1"}}},
				},
				{
					ID:         "11",
					Permission: storage.ACL_WRITE,
					Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
					Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
				},
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

		AssociationsToAdd: []*storage.EntityID{
			{Type: "foo", Key: "bar"},
			{Type: "bar", Key: "baz"},
			{Type: "baz", Key: "quz"},
		},
	})
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.NetworkEntity{
			NetworkID:  "n1",
			Type:       "hello",
			Key:        "world",
			GraphID:    "10",
			PhysicalID: newPhysID,
			Associations: []*storage.EntityID{
				{Type: "foo", Key: "bar"},
				{Type: "bar", Key: "baz"},
				{Type: "baz", Key: "quz"},
			},
			Version: 2,
		},
		updateHelloWorldEntResult,
	)

	// Read back the updated ent
	expectedHelloWorldEnt.Associations = []*storage.EntityID{
		{Type: "bar", Key: "baz"},
		{Type: "baz", Key: "quz"},
		{Type: "foo", Key: "bar"},
	}
	expectedHelloWorldEnt.GraphID = "10"
	expectedHelloWorldEnt.Permissions = []*storage.ACL{
		{
			ID:         "11",
			Permission: storage.ACL_WRITE,
			Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
			Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
			Version:    1,
		},
		{
			ID:         "12",
			Permission: storage.ACL_WRITE,
			Type:       &storage.ACL_EntityType{EntityType: "foo"},
			Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n1"}}},
		},
	}
	expectedHelloWorldEnt.Version = 2

	actualEntityLoad, err = store.LoadEntities(
		"n1",
		storage.EntityLoadFilter{IDs: []*storage.EntityID{{Type: "hello", Key: "world"}}},
		storage.FullEntityLoadCriteria,
	)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities:         []*storage.NetworkEntity{&expectedHelloWorldEnt},
			EntitiesNotFound: []*storage.EntityID{},
		},
		actualEntityLoad,
	)

	assert.NoError(t, store.Commit())

	// ========================================================================
	// Graph load tests
	// ========================================================================

	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)

	// Load a graph directly via ID
	actualGraph10, err := store.LoadGraphForEntity("n1", storage.EntityID{Type: "hello", Key: "world"}, storage.EntityLoadCriteria{})
	assert.NoError(t, err)

	expectedFoobarEnt.ParentAssociations = append(expectedFoobarEnt.ParentAssociations, &storage.EntityID{Type: "hello", Key: "world"})
	expectedBarbazEnt.ParentAssociations = append(expectedBarbazEnt.ParentAssociations, &storage.EntityID{Type: "hello", Key: "world"})
	expectedBazquzEnt.ParentAssociations = []*storage.EntityID{{Type: "hello", Key: "world"}}

	expectedFoobarEnt = storage.NetworkEntity{
		NetworkID:  "n1",
		Type:       "foo",
		Key:        "bar",
		GraphID:    "10",
		PhysicalID: "1",
		ParentAssociations: []*storage.EntityID{
			{Type: "baz", Key: "quz"},
			{Type: "hello", Key: "world"},
		},
	}

	expectedBarbazEnt = storage.NetworkEntity{
		NetworkID: "n1",
		Type:      "bar",
		Key:       "baz",
		GraphID:   "10",
		ParentAssociations: []*storage.EntityID{
			{Type: "baz", Key: "quz"},
			{Type: "hello", Key: "world"},
		},
	}

	expectedBazquzEnt = storage.NetworkEntity{
		NetworkID: "n1",
		Type:      "baz",
		Key:       "quz",
		GraphID:   "10",
		ParentAssociations: []*storage.EntityID{
			{Type: "hello", Key: "world"},
		},
		Associations: []*storage.EntityID{
			{Type: "bar", Key: "baz"},
			{Type: "foo", Key: "bar"},
		},
	}

	expectedHelloWorldEnt = storage.NetworkEntity{
		NetworkID:  "n1",
		Type:       "hello",
		Key:        "world",
		GraphID:    "10",
		PhysicalID: "asdf",
		Associations: []*storage.EntityID{
			{Type: "bar", Key: "baz"},
			{Type: "baz", Key: "quz"},
			{Type: "foo", Key: "bar"},
		},
		Version: 2,
	}

	expectedGraph10 := storage.EntityGraph{
		Entities: []*storage.NetworkEntity{
			&expectedBarbazEnt,
			&expectedBazquzEnt,
			&expectedFoobarEnt,
			&expectedHelloWorldEnt,
		},
		RootEntities: []*storage.EntityID{{Type: "hello", Key: "world"}},
		Edges: []*storage.GraphEdge{
			{From: &storage.EntityID{Type: "baz", Key: "quz"}, To: &storage.EntityID{Type: "bar", Key: "baz"}},
			{From: &storage.EntityID{Type: "baz", Key: "quz"}, To: &storage.EntityID{Type: "foo", Key: "bar"}},
			{From: &storage.EntityID{Type: "hello", Key: "world"}, To: &storage.EntityID{Type: "bar", Key: "baz"}},
			{From: &storage.EntityID{Type: "hello", Key: "world"}, To: &storage.EntityID{Type: "baz", Key: "quz"}},
			{From: &storage.EntityID{Type: "hello", Key: "world"}, To: &storage.EntityID{Type: "foo", Key: "bar"}},
		},
	}
	assert.Equal(t, expectedGraph10, actualGraph10)

	// Load a graph from the ID of a node in the middle
	actualGraph10, err = store.LoadGraphForEntity("n1", storage.EntityID{Type: "baz", Key: "quz"}, storage.EntityLoadCriteria{})
	assert.NoError(t, err)
	assert.Equal(t, expectedGraph10, actualGraph10)

	// Load a graph with full load criteria
	actualGraph10, err = store.LoadGraphForEntity("n1", storage.EntityID{Type: "foo", Key: "bar"}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)

	expectedFoobarEnt.Name = "foobar"
	expectedFoobarEnt.Description = "foobar ent"
	expectedFoobarEnt.Config = []byte("foobar")

	expectedBarbazEnt.Name = "barbaz"
	expectedBarbazEnt.Description = "barbaz ent"
	expectedBarbazEnt.Config = []byte("barbaz")
	expectedBarbazEnt.Permissions = []*storage.ACL{
		{
			ID:         "5",
			Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
			Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
			Permission: storage.ACL_NO_PERM,
		},
		{
			ID:         "6",
			Type:       &storage.ACL_EntityType{EntityType: "foo"},
			Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n1"}}},
			Permission: storage.ACL_WRITE,
		},
	}

	expectedBazquzEnt.Name = "bazquz"
	expectedBazquzEnt.Description = "bazquz ent"

	expectedHelloWorldEnt.Name = "helloworld2"
	expectedHelloWorldEnt.Description = "helloworld2 ent"
	expectedHelloWorldEnt.Config = []byte("second config")
	expectedHelloWorldEnt.Permissions = []*storage.ACL{
		{
			ID:         "11",
			Permission: storage.ACL_WRITE,
			Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
			Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
			Version:    1,
		},
		{
			ID:         "12",
			Permission: storage.ACL_WRITE,
			Type:       &storage.ACL_EntityType{EntityType: "foo"},
			Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n1"}}},
		},
	}

	expectedGraph10 = storage.EntityGraph{
		Entities: []*storage.NetworkEntity{
			&expectedBarbazEnt,
			&expectedBazquzEnt,
			&expectedFoobarEnt,
			&expectedHelloWorldEnt,
		},
		RootEntities: []*storage.EntityID{{Type: "hello", Key: "world"}},
		Edges: []*storage.GraphEdge{
			{From: &storage.EntityID{Type: "baz", Key: "quz"}, To: &storage.EntityID{Type: "bar", Key: "baz"}},
			{From: &storage.EntityID{Type: "baz", Key: "quz"}, To: &storage.EntityID{Type: "foo", Key: "bar"}},
			{From: &storage.EntityID{Type: "hello", Key: "world"}, To: &storage.EntityID{Type: "bar", Key: "baz"}},
			{From: &storage.EntityID{Type: "hello", Key: "world"}, To: &storage.EntityID{Type: "baz", Key: "quz"}},
			{From: &storage.EntityID{Type: "hello", Key: "world"}, To: &storage.EntityID{Type: "foo", Key: "bar"}},
		},
	}
	assert.Equal(t, expectedGraph10, actualGraph10)
	assert.NoError(t, store.Commit())

	// TODO: orphan some graph nodes (blocked on impl of fixGraph)
	// As a reminder, this is the current state of the graph:
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

	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)

	// delete assocs to foobar and barbaz, helloworld's graph ID should still be 2
	_, err = store.UpdateEntity(
		"n1",
		storage.EntityUpdateCriteria{
			Type: "hello", Key: "world",
			AssociationsToDelete: []*storage.EntityID{
				{Type: "foo", Key: "bar"},
				{Type: "bar", Key: "baz"},
			},
		},
	)
	assert.NoError(t, err)

	expectedHelloWorldEnt.Associations = []*storage.EntityID{{Type: "baz", Key: "quz"}}
	expectedHelloWorldEnt.Version = 3
	expectedFoobarEnt.ParentAssociations = []*storage.EntityID{{Type: "baz", Key: "quz"}}
	expectedBarbazEnt.ParentAssociations = []*storage.EntityID{{Type: "baz", Key: "quz"}}

	expectedGraph10 = storage.EntityGraph{
		Entities: []*storage.NetworkEntity{
			&expectedBarbazEnt,
			&expectedBazquzEnt,
			&expectedFoobarEnt,
			&expectedHelloWorldEnt,
		},
		RootEntities: []*storage.EntityID{{Type: "hello", Key: "world"}},
		Edges: []*storage.GraphEdge{
			{From: &storage.EntityID{Type: "baz", Key: "quz"}, To: &storage.EntityID{Type: "bar", Key: "baz"}},
			{From: &storage.EntityID{Type: "baz", Key: "quz"}, To: &storage.EntityID{Type: "foo", Key: "bar"}},
			{From: &storage.EntityID{Type: "hello", Key: "world"}, To: &storage.EntityID{Type: "baz", Key: "quz"}},
		},
	}
	actualGraph10, err = store.LoadGraphForEntity("n1", storage.EntityID{Type: "hello", Key: "world"}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(t, expectedGraph10, actualGraph10)

	// delete assoc to bazquz, helloworld should have a new graph ID
	_, err = store.UpdateEntity(
		"n1",
		storage.EntityUpdateCriteria{
			Type: "hello", Key: "world",
			AssociationsToDelete: []*storage.EntityID{
				{Type: "baz", Key: "quz"},
			},
		},
	)
	assert.NoError(t, err)

	expectedHelloWorldEnt.Associations = nil
	expectedHelloWorldEnt.Version = 4
	expectedHelloWorldEnt.GraphID = "13"
	expectedBazquzEnt.ParentAssociations = nil

	// first, check graph of foobar which should be unchanged
	expectedGraph10 = storage.EntityGraph{
		Entities: []*storage.NetworkEntity{
			&expectedBarbazEnt,
			&expectedBazquzEnt,
			&expectedFoobarEnt,
		},
		RootEntities: []*storage.EntityID{{Type: "baz", Key: "quz"}},
		Edges: []*storage.GraphEdge{
			{From: &storage.EntityID{Type: "baz", Key: "quz"}, To: &storage.EntityID{Type: "bar", Key: "baz"}},
			{From: &storage.EntityID{Type: "baz", Key: "quz"}, To: &storage.EntityID{Type: "foo", Key: "bar"}},
		},
	}
	actualGraph10, err = store.LoadGraphForEntity("n1", storage.EntityID{Type: "foo", Key: "bar"}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(t, expectedGraph10, actualGraph10)

	// helloworld should be in its own graph now. ID generator was at 13
	expectedGraph13 := storage.EntityGraph{
		Entities: []*storage.NetworkEntity{
			&expectedHelloWorldEnt,
		},
		RootEntities: []*storage.EntityID{{Type: "hello", Key: "world"}},
		Edges:        []*storage.GraphEdge{},
	}
	actualGraph13, err := store.LoadGraphForEntity("n1", storage.EntityID{Type: "hello", Key: "world"}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(t, expectedGraph13, actualGraph13)

	// now, delete bazquz. this should partition the network into 3 different
	// graphs each with a single element
	_, err = store.UpdateEntity("n1", storage.EntityUpdateCriteria{Type: "baz", Key: "quz", DeleteEntity: true})
	assert.NoError(t, err)

	allEnts, err := store.LoadEntities("n1", storage.EntityLoadFilter{}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.NotNil(t, allEnts)

	expectedBarbazEnt.ParentAssociations = nil
	expectedFoobarEnt.GraphID = "14"
	expectedFoobarEnt.ParentAssociations = nil

	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				&expectedBarbazEnt,
				&expectedFoobarEnt,
				&expectedHelloWorldEnt,
			},
			EntitiesNotFound: []*storage.EntityID{},
		},
		allEnts,
	)

	// associate helloworld -> foobar, then set helloworld's associations to
	// just -> barbaz
	_, err = store.UpdateEntity(
		"n1",
		storage.EntityUpdateCriteria{
			Type: "hello",
			Key:  "world",
			AssociationsToAdd: []*storage.EntityID{
				{Type: "foo", Key: "bar"},
			},
		},
	)
	assert.NoError(t, err)
	_, err = store.UpdateEntity(
		"n1",
		storage.EntityUpdateCriteria{
			Type: "hello",
			Key:  "world",
			AssociationsToSet: &storage.EntityAssociationsToSet{
				AssociationsToSet: []*storage.EntityID{
					{Type: "bar", Key: "baz"},
				},
			},
		},
	)
	assert.NoError(t, err)

	expectedHelloWorldEnt.Associations = []*storage.EntityID{{Type: "bar", Key: "baz"}}
	expectedBarbazEnt.ParentAssociations = []*storage.EntityID{{Type: "hello", Key: "world"}}
	expectedBarbazEnt.GraphID = "10"
	expectedHelloWorldEnt.GraphID = "10"
	expectedFoobarEnt.GraphID = "15"
	expectedHelloWorldEnt.Version = 6

	allEnts, err = store.LoadEntities("n1", storage.EntityLoadFilter{}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.NotNil(t, allEnts)

	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				&expectedBarbazEnt,
				&expectedFoobarEnt,
				&expectedHelloWorldEnt,
			},
			EntitiesNotFound: []*storage.EntityID{},
		},
		allEnts,
	)

	// Clear associations using AssociationsToSet field
	_, err = store.UpdateEntity(
		"n1",
		storage.EntityUpdateCriteria{
			Type: "hello",
			Key:  "world",
			AssociationsToSet: &storage.EntityAssociationsToSet{
				AssociationsToSet: []*storage.EntityID{},
			},
		},
	)
	assert.NoError(t, err)

	expectedHelloWorldEnt.Associations = nil
	expectedBarbazEnt.ParentAssociations = nil
	expectedBarbazEnt.GraphID = "16"
	expectedHelloWorldEnt.GraphID = "10"
	expectedFoobarEnt.GraphID = "15"
	expectedHelloWorldEnt.Version = 7

	allEnts, err = store.LoadEntities("n1", storage.EntityLoadFilter{}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.NotNil(t, allEnts)

	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				&expectedBarbazEnt,
				&expectedFoobarEnt,
				&expectedHelloWorldEnt,
			},
			EntitiesNotFound: []*storage.EntityID{},
		},
		allEnts,
	)
}
