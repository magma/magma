/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"
	orc8r_storage "magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

const (
	integTestMaxLoadSize = 5
)

type mockIDGenerator struct {
	count int
}

func (g *mockIDGenerator) New() string {
	g.count++
	return fmt.Sprintf("%d", g.count)
}

func TestSqlConfiguratorStorage_Integration(t *testing.T) {
	// sqlite's default behavior is to disable foreign keys: https://www.sqlite.org/draft/pragma.html#pragma_foreign_keys
	// thankfully the sqlite3 driver supports the appropriate pragma: https://github.com/mattn/go-sqlite3/issues/255
	db, err := sqorc.Open("sqlite3", ":memory:?_foreign_keys=1")
	if err != nil {
		t.Fatalf("Could not initialize sqlite DB: %s", err)
	}
	factory := storage.NewSQLConfiguratorStorageFactory(db, &mockIDGenerator{}, sqorc.GetSqlBuilder(), integTestMaxLoadSize)
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
	assert.NoError(t, err)
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

	store, err = factory.StartTransaction(context.Background(), &orc8r_storage.TxOptions{ReadOnly: true})
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
	assert.NoError(t, err)

	actualEntityLoad, err := store.LoadEntities("n1", storage.EntityLoadFilter{}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.EntityLoadResult{},
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
			EntitiesNotFound: []*storage.EntityID{
				{Type: "baz", Key: "quz"},
				{Type: "foo", Key: "bar"},
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
		// should be ignored
		Pk: "0xdeadbeef",
	})
	assert.NoError(t, err)
	assert.Equal(t, "2", expectedFoobarEnt.GraphID)

	expectedBarbazEnt, err := store.CreateEntity("n1", storage.NetworkEntity{
		Type: "bar",
		Key:  "baz",

		Name:        "barbaz",
		Description: "barbaz ent",

		Config: []byte("barbaz"),
	})
	assert.NoError(t, err)
	assert.Equal(t, "4", expectedBarbazEnt.GraphID)

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

	// Try and fail to create ent with physical ID that already exists
	_, err = store.CreateEntity("n1", storage.NetworkEntity{
		Type:       "apple_type",
		Key:        "apple_key",
		PhysicalID: "1",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "an entity with physical ID '1' already exists")

	assert.NoError(t, store.Commit())
	assert.Equal(t, "2", expectedBazquzEnt.GraphID)

	expectedFoobarEnt.GraphID = "2"
	expectedBarbazEnt.GraphID = "2"

	expectedFoobarEnt.Pk = "1"
	expectedBarbazEnt.Pk = "3"
	expectedBazquzEnt.Pk = "5"

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
		},
		actualEntityLoad,
	)
	assert.NoError(t, store.Commit())

	// Load from physical ID
	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	// network ID shouldn't matter in this query, since searching on
	// physical ID intentionally doesn't scope to a particular network
	actualEntityLoad, err = store.LoadEntities("placeholder", storage.EntityLoadFilter{PhysicalID: stringPointer("1")}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				&expectedFoobarEnt,
			},
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
	// pk should be 9, graph id should be 10

	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	_, err = store.CreateEntity("n1", storage.NetworkEntity{
		Type: "hello",
		Key:  "world",

		Name:        "helloworld",
		Description: "helloworld ent",

		Config: []byte("first config"),
	})
	assert.NoError(t, err)

	// update basic fields

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
	})
	assert.NoError(t, err)

	assert.Equal(
		t,
		storage.NetworkEntity{
			NetworkID: "n1",
			Type:      "hello",
			Key:       "world",

			Pk:      "7",
			GraphID: "8",

			PhysicalID: newPhysID,

			Name:        newName,
			Description: newDesc,
			Config:      newConfig,
			Version:     1,
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
			Pk:         "7",
			GraphID:    "2",
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
	expectedHelloWorldEnt.GraphID = "2"
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
			Entities: []*storage.NetworkEntity{&expectedHelloWorldEnt},
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
	actualGraph2, err := store.LoadGraphForEntity("n1", storage.EntityID{Type: "hello", Key: "world"}, storage.EntityLoadCriteria{})
	assert.NoError(t, err)

	expectedFoobarEnt.ParentAssociations = append(expectedFoobarEnt.ParentAssociations, &storage.EntityID{Type: "hello", Key: "world"})
	expectedBarbazEnt.ParentAssociations = append(expectedBarbazEnt.ParentAssociations, &storage.EntityID{Type: "hello", Key: "world"})
	expectedBazquzEnt.ParentAssociations = []*storage.EntityID{{Type: "hello", Key: "world"}}

	expectedFoobarEnt = storage.NetworkEntity{
		NetworkID:  "n1",
		Type:       "foo",
		Key:        "bar",
		Pk:         "1",
		GraphID:    "2",
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
		Pk:        "3",
		GraphID:   "2",
		ParentAssociations: []*storage.EntityID{
			{Type: "baz", Key: "quz"},
			{Type: "hello", Key: "world"},
		},
	}

	expectedBazquzEnt = storage.NetworkEntity{
		NetworkID: "n1",
		Type:      "baz",
		Key:       "quz",
		Pk:        "5",
		GraphID:   "2",
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
		Pk:         "7",
		GraphID:    "2",
		PhysicalID: "asdf",
		Associations: []*storage.EntityID{
			{Type: "bar", Key: "baz"},
			{Type: "baz", Key: "quz"},
			{Type: "foo", Key: "bar"},
		},
		Version: 2,
	}

	expectedGraph2 := storage.EntityGraph{
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
	assert.Equal(t, expectedGraph2, actualGraph2)

	// Load a graph from the ID of a node in the middle
	actualGraph2, err = store.LoadGraphForEntity("n1", storage.EntityID{Type: "baz", Key: "quz"}, storage.EntityLoadCriteria{})
	assert.NoError(t, err)
	assert.Equal(t, expectedGraph2, actualGraph2)

	// Load a graph with full load criteria
	actualGraph2, err = store.LoadGraphForEntity("n1", storage.EntityID{Type: "foo", Key: "bar"}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)

	expectedFoobarEnt.Name = "foobar"
	expectedFoobarEnt.Description = "foobar ent"
	expectedFoobarEnt.Config = []byte("foobar")

	expectedBarbazEnt.Name = "barbaz"
	expectedBarbazEnt.Description = "barbaz ent"
	expectedBarbazEnt.Config = []byte("barbaz")

	expectedBazquzEnt.Name = "bazquz"
	expectedBazquzEnt.Description = "bazquz ent"

	expectedHelloWorldEnt.Name = "helloworld2"
	expectedHelloWorldEnt.Description = "helloworld2 ent"
	expectedHelloWorldEnt.Config = []byte("second config")

	expectedGraph2 = storage.EntityGraph{
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
	assert.Equal(t, expectedGraph2, actualGraph2)
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

	expectedGraph2 = storage.EntityGraph{
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
	actualGraph2, err = store.LoadGraphForEntity("n1", storage.EntityID{Type: "hello", Key: "world"}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(t, expectedGraph2, actualGraph2)

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
	expectedHelloWorldEnt.GraphID = "9"
	expectedBazquzEnt.ParentAssociations = nil

	// first, check graph of foobar which should be unchanged
	expectedGraph2 = storage.EntityGraph{
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
	actualGraph2, err = store.LoadGraphForEntity("n1", storage.EntityID{Type: "foo", Key: "bar"}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(t, expectedGraph2, actualGraph2)

	// helloworld should be in its own graph now. ID generator was at 9.
	expectedGraph9 := storage.EntityGraph{
		Entities: []*storage.NetworkEntity{
			&expectedHelloWorldEnt,
		},
		RootEntities: []*storage.EntityID{{Type: "hello", Key: "world"}},
		Edges:        []*storage.GraphEdge{},
	}
	actualGraph9, err := store.LoadGraphForEntity("n1", storage.EntityID{Type: "hello", Key: "world"}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(t, expectedGraph9, actualGraph9)

	// now, delete bazquz. this should partition the network into 3 different
	// graphs each with a single element
	_, err = store.UpdateEntity("n1", storage.EntityUpdateCriteria{Type: "baz", Key: "quz", DeleteEntity: true})
	assert.NoError(t, err)

	allEnts, err := store.LoadEntities("n1", storage.EntityLoadFilter{}, storage.FullEntityLoadCriteria)
	assert.NoError(t, err)
	assert.NotNil(t, allEnts)

	expectedBarbazEnt.ParentAssociations = nil
	expectedFoobarEnt.GraphID = "10"
	expectedFoobarEnt.ParentAssociations = nil

	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				&expectedBarbazEnt,
				&expectedFoobarEnt,
				&expectedHelloWorldEnt,
			},
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
	expectedFoobarEnt.GraphID = "11"
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
	expectedBarbazEnt.GraphID = "12"
	expectedHelloWorldEnt.GraphID = "10"
	expectedFoobarEnt.GraphID = "11"
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
		},
		allEnts,
	)
	assert.NoError(t, store.Commit())

	// ========================================================================
	// Paginated load tests
	// ========================================================================

	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)

	// Create two more entities of type foo for paginated load
	expectedFoozedEnt, err := store.CreateEntity("n1", storage.NetworkEntity{
		Type: "foo",
		Key:  "zed",

		Name:        "foozed",
		Description: "foozed ent",

		PhysicalID: "6",

		Config: []byte("foozed"),
	})
	assert.NoError(t, err)
	assert.Equal(t, "14", expectedFoozedEnt.GraphID)

	expectedFoohueEnt, err := store.CreateEntity("n1", storage.NetworkEntity{
		Type: "foo",
		Key:  "hue",

		Name:        "foohue",
		Description: "foozed ent",

		PhysicalID: "7",

		Config: []byte("foohue"),
	})
	assert.NoError(t, err)
	assert.NoError(t, store.Commit())
	assert.Equal(t, "15", expectedFoohueEnt.Pk)
	assert.Equal(t, "16", expectedFoohueEnt.GraphID)

	// Load paginated entities of the same type
	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	paginatedLoadCriteria := storage.FullEntityLoadCriteria
	paginatedLoadCriteria.PageToken = ""
	paginatedLoadCriteria.PageSize = 2
	actualEntityLoad, err = store.LoadEntities("n1", storage.EntityLoadFilter{TypeFilter: &wrappers.StringValue{Value: "foo"}}, paginatedLoadCriteria)
	assert.NoError(t, err)
	nextToken := &storage.EntityPageToken{
		LastIncludedEntity: "hue",
	}
	expectedNextToken := serializeToken(t, nextToken)

	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				&expectedFoobarEnt, // type: foo, key: bar
				&expectedFoohueEnt, // type: foo, key: hue
			},
			NextPageToken: expectedNextToken,
		},
		actualEntityLoad,
	)
	assert.NoError(t, store.Commit())

	// Load paginated entities (page 2)
	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	paginatedLoadCriteria.PageToken = expectedNextToken
	actualEntityLoad, err = store.LoadEntities("n1", storage.EntityLoadFilter{TypeFilter: &wrappers.StringValue{Value: "foo"}}, paginatedLoadCriteria)
	assert.NoError(t, err)
	assert.Equal(
		t,
		storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				&expectedFoozedEnt, // type: foo, key: zed
			},
			NextPageToken: "",
		},
		actualEntityLoad,
	)
	assert.NoError(t, store.Commit())

	// Ensure multi-type pagination loads fail
	paginatedLoadCriteria.PageToken = ""
	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	_, err = store.LoadEntities("n1", storage.EntityLoadFilter{}, paginatedLoadCriteria)
	assert.Error(t, err)
	assert.NoError(t, store.Commit())

	paginatedLoadCriteria.PageToken = "aaa"
	paginatedLoadCriteria.PageSize = 0
	store, err = factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	_, err = store.LoadEntities("n1", storage.EntityLoadFilter{}, paginatedLoadCriteria)
	assert.Error(t, err)
	assert.NoError(t, store.Commit())
}
