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
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"
	storage2 "magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var mockResult = sqlmock.NewResult(1, 1)

func TestSqlConfiguratorStorage_LoadNetworks(t *testing.T) {
	runFactory := func(ids []string, loadCriteria storage.NetworkLoadCriteria) func(store storage.ConfiguratorStorage) (interface{}, error) {
		return func(store storage.ConfiguratorStorage) (interface{}, error) {
			return store.LoadNetworks(storage.NetworkLoadFilter{Ids: ids}, loadCriteria)
		}
	}

	idsOnly := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT cfg_networks.id, cfg_networks.type, cfg_networks.version FROM cfg_networks").
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "type", "version"}).
						AddRow("hello", "", 1).
						AddRow("world", "", 2),
				)
		},
		run: runFactory([]string{"hello", "world"}, storage.NetworkLoadCriteria{}),

		expectedResult: storage.NetworkLoadResult{
			NetworkIDsNotFound: []string{},
			Networks: []*storage.Network{
				{ID: "hello", Configs: map[string][]byte{}, Version: 1},
				{ID: "world", Configs: map[string][]byte{}, Version: 2},
			},
		},
	}

	metadataLoad := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT cfg_networks.id, cfg_networks.type, cfg_networks.name, cfg_networks.description, cfg_networks.version FROM cfg_networks").
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "type", "name", "description", "version"}).
						AddRow("hello", "", "Hello", "Hello network", 1).
						AddRow("world", "", "World", "World network", 2),
				)
		},
		run: runFactory([]string{"hello", "world"}, storage.NetworkLoadCriteria{LoadMetadata: true}),

		expectedResult: storage.NetworkLoadResult{
			NetworkIDsNotFound: []string{},
			Networks: []*storage.Network{
				{
					ID:          "hello",
					Name:        "Hello",
					Description: "Hello network",
					Configs:     map[string][]byte{},
					Version:     1,
				},
				{
					ID:          "world",
					Name:        "World",
					Description: "World network",
					Configs:     map[string][]byte{},
					Version:     2,
				},
			},
		},
	}

	// 1 network has 2 configs, 1 has no configs
	idsAndConfigs := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT cfg_networks.id, cfg_networks.type, cfg_network_configs.type, cfg_network_configs.value, cfg_networks.version FROM cfg_networks").
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "type", "type", "value", "version"}).
						AddRow("hello", "", "foo", []byte("foo"), 1).
						AddRow("hello", "", "bar", []byte("bar"), 1).
						AddRow("world", "", nil, nil, 3),
				)
		},
		run: runFactory([]string{"hello", "world", "foobar"}, storage.NetworkLoadCriteria{LoadConfigs: true}),

		expectedResult: storage.NetworkLoadResult{
			NetworkIDsNotFound: []string{"foobar"},
			Networks: []*storage.Network{
				{
					ID: "hello",
					Configs: map[string][]byte{
						"foo": []byte("foo"),
						"bar": []byte("bar"),
					},
					Version: 1,
				},
				{
					ID:      "world",
					Configs: map[string][]byte{},
					Version: 3,
				},
			},
		},
	}

	// Same setup as above, plus 1 network not found
	fullLoad := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT cfg_networks.id, cfg_networks.type, cfg_networks.name, cfg_networks.description, cfg_network_configs.type, cfg_network_configs.value, cfg_networks.version FROM cfg_networks").
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "type", "name", "description", "type", "value", "version"}).
						AddRow("hello", "", "Hello", "Hello network", "foo", []byte("foo"), 1).
						AddRow("hello", "", "Hello", "Hello network", "bar", []byte("bar"), 1).
						AddRow("world", "", "World", "World network", nil, nil, 2),
				)
		},
		run: runFactory([]string{"hello", "world", "foobar"}, storage.NetworkLoadCriteria{LoadMetadata: true, LoadConfigs: true}),

		expectedResult: storage.NetworkLoadResult{
			NetworkIDsNotFound: []string{"foobar"},
			Networks: []*storage.Network{
				{
					ID:          "hello",
					Name:        "Hello",
					Description: "Hello network",
					Configs: map[string][]byte{
						"foo": []byte("foo"),
						"bar": []byte("bar"),
					},
					Version: 1,
				},
				{
					ID:          "world",
					Name:        "World",
					Description: "World network",
					Configs:     map[string][]byte{},
					Version:     2,
				},
			},
		},
	}

	noneFound := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT cfg_networks.id, cfg_networks.type, cfg_networks.version FROM cfg_networks").
				WillReturnRows(sqlmock.NewRows([]string{"id", "", "version"}))
		},
		run: runFactory([]string{"hello", "world"}, storage.NetworkLoadCriteria{}),

		expectedResult: storage.NetworkLoadResult{
			NetworkIDsNotFound: []string{"hello", "world"},
			Networks:           []*storage.Network{},
		},
	}

	queryError := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT cfg_networks.id, cfg_networks.type, cfg_networks.version FROM cfg_networks").
				WillReturnError(errors.New("mock query error"))
		},
		run: runFactory([]string{"hello", "world"}, storage.NetworkLoadCriteria{}),

		expectedError: errors.New("error querying for networks: mock query error"),
	}

	runCase(t, idsOnly)
	runCase(t, metadataLoad)
	runCase(t, idsAndConfigs)
	runCase(t, fullLoad)
	runCase(t, noneFound)
	runCase(t, queryError)
}

func TestSqlConfiguratorStorage_CreateNetwork(t *testing.T) {
	runFactory := func(network storage.Network) func(store storage.ConfiguratorStorage) (interface{}, error) {
		return func(store storage.ConfiguratorStorage) (interface{}, error) {
			return store.CreateNetwork(network)
		}
	}

	idOnly := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT COUNT\(1\) FROM cfg_networks`).
				WithArgs("n1").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			m.ExpectExec("INSERT INTO cfg_networks").
				WithArgs("n1", "", "", "").
				WillReturnResult(mockResult)
		},
		run: runFactory(storage.Network{ID: "n1"}),

		expectedResult: storage.Network{ID: "n1"},
	}

	allMetadata := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT COUNT\(1\) FROM cfg_networks`).
				WithArgs("n2").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			m.ExpectExec("INSERT INTO cfg_networks").
				WithArgs("n2", "", "hello", "world").
				WillReturnResult(mockResult)
		},
		run: runFactory(storage.Network{ID: "n2", Name: "hello", Description: "world"}),

		expectedResult: storage.Network{ID: "n2", Name: "hello", Description: "world"},
	}

	everythingNw := storage.Network{
		ID:          "n3",
		Type:        "lte",
		Name:        "hello",
		Description: "world",
		Configs: map[string][]byte{
			"foo": []byte("bar"),
			"baz": []byte("quz"),
		},
	}
	everything := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT COUNT\(1\) FROM cfg_networks`).
				WithArgs("n3").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			m.ExpectExec("INSERT INTO cfg_networks").
				WithArgs("n3", "lte", "hello", "world").
				WillReturnResult(mockResult)

			m.ExpectExec("INSERT INTO cfg_network_configs").
				WithArgs(
					"n3", "baz", []byte("quz"),
					"n3", "foo", []byte("bar"),
				).
				WillReturnResult(mockResult)
		},
		run: runFactory(everythingNw),

		expectedResult: everythingNw,
	}

	networkTableError := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT COUNT\(1\) FROM cfg_networks`).
				WithArgs("n4").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			m.ExpectExec("INSERT INTO cfg_networks").
				WithArgs("n4", "", "", "").
				WillReturnError(errors.New("mock exec error"))
		},
		run: runFactory(storage.Network{ID: "n4"}),

		expectedResult: storage.Network{ID: "n4"},
		expectedError:  errors.New("error inserting network: mock exec error"),
	}

	configTableErrNw := storage.Network{
		ID: "n5",
		Configs: map[string][]byte{
			"foo": []byte("bar"),
			"baz": []byte("quz"),
		},
	}
	configTableError := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT COUNT\(1\) FROM cfg_networks`).
				WithArgs("n5").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			m.ExpectExec("INSERT INTO cfg_networks").
				WithArgs("n5", "", "", "").
				WillReturnResult(mockResult)

			m.ExpectExec("INSERT INTO cfg_network_configs").
				WithArgs(
					"n5", "baz", []byte("quz"),
					"n5", "foo", []byte("bar"),
				).
				WillReturnError(errors.New("mock exec error"))
		},
		run: runFactory(configTableErrNw),

		expectedResult: configTableErrNw,
		expectedError:  errors.New("error inserting network configs: mock exec error"),
	}

	networkExists := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT COUNT\(1\) FROM cfg_networks`).
				WithArgs("n5").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		},
		run: runFactory(storage.Network{ID: "n5"}),

		expectedResult: storage.Network{ID: "n5"},
		expectedError:  errors.New("a network with ID n5 already exists"),
	}

	runCase(t, idOnly)
	runCase(t, allMetadata)
	runCase(t, everything)
	runCase(t, networkTableError)
	runCase(t, configTableError)
	runCase(t, networkExists)
}

func TestSqlConfiguratorStorage_UpdateNetworks(t *testing.T) {
	runFactory := func(updates []storage.NetworkUpdateCriteria) func(store storage.ConfiguratorStorage) (interface{}, error) {
		return func(store storage.ConfiguratorStorage) (interface{}, error) {
			return nil, store.UpdateNetworks(updates)
		}
	}

	// Delete 1 network (n1)
	// Update only metadata of another (n2)
	// Update and delete configs of another (n3)
	// Fill out all fields of the update criteria (n4)
	// Prepared statements should be cached and closed on exit
	names := []string{"should be ignored", "name2", "name4"}
	descs := []string{"should be ignored", "desc2", ""}
	happyPath := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			prepWithNameAndDesc := m.ExpectPrepare("UPDATE cfg_networks").WillBeClosed()
			prepWithNameAndDesc.ExpectExec().WithArgs(names[1], descs[1], "n2").WillReturnResult(mockResult)

			prepWithOnlyVersion := m.ExpectPrepare("UPDATE cfg_networks").WillBeClosed()
			prepWithOnlyVersion.ExpectExec().WithArgs("n3").WillReturnResult(mockResult)

			upsertStmt := m.ExpectPrepare("INSERT INTO cfg_network_configs").WillBeClosed()
			upsertStmt.ExpectExec().WithArgs("n3", "baz", []byte("quz"), []byte("quz")).WillReturnResult(mockResult)
			upsertStmt.ExpectExec().WithArgs("n3", "foo", []byte("bar"), []byte("bar")).WillReturnResult(mockResult)
			m.ExpectExec("DELETE FROM cfg_network_configs").WithArgs("n3", "hello", "n3", "world").WillReturnResult(mockResult)

			prepWithNameAndDesc.ExpectExec().WithArgs(names[2], "", "n4").WillReturnResult(mockResult)
			upsertStmt.ExpectExec().WithArgs("n4", "baz", []byte("quz"), []byte("quz")).WillReturnResult(mockResult)
			upsertStmt.ExpectExec().WithArgs("n4", "foo", []byte("bar"), []byte("bar")).WillReturnResult(mockResult)
			m.ExpectExec("DELETE FROM cfg_network_configs").WithArgs("n4", "hello", "n4", "world").WillReturnResult(mockResult)

			m.ExpectExec("DELETE FROM cfg_network_configs").WithArgs("n1").WillReturnResult(mockResult)
			m.ExpectExec("DELETE FROM cfg_networks").WithArgs("n1").WillReturnResult(mockResult)
		},
		run: runFactory(
			[]storage.NetworkUpdateCriteria{
				{ID: "n1", DeleteNetwork: true, NewName: &wrappers.StringValue{Value: names[0]}, NewDescription: &wrappers.StringValue{Value: descs[0]}},
				{ID: "n2", NewName: &wrappers.StringValue{Value: names[1]}, NewDescription: &wrappers.StringValue{Value: descs[1]}},
				{
					ID:              "n3",
					ConfigsToDelete: []string{"hello", "world"},
					ConfigsToAddOrUpdate: map[string][]byte{
						"foo": []byte("bar"),
						"baz": []byte("quz"),
					},
				},
				{
					ID:              "n4",
					NewName:         &wrappers.StringValue{Value: names[2]},
					NewDescription:  &wrappers.StringValue{Value: descs[2]},
					ConfigsToDelete: []string{"hello", "world"},
					ConfigsToAddOrUpdate: map[string][]byte{
						"foo": []byte("bar"),
						"baz": []byte("quz"),
					},
				},
			},
		),
	}

	errorCase := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			updateStmt := m.ExpectPrepare("UPDATE cfg_networks").WillBeClosed()
			updateStmt.ExpectExec().WithArgs("name2", "desc2", "n2").WillReturnError(errors.New("mock update error"))
		},
		run: runFactory(
			[]storage.NetworkUpdateCriteria{
				{ID: "n1", DeleteNetwork: true},
				{ID: "n2", NewName: &wrappers.StringValue{Value: names[1]}, NewDescription: &wrappers.StringValue{Value: descs[1]}},
			},
		),

		expectedError: errors.New("error updating network n2: mock update error"),
	}

	validationFailure := &testCase{
		setup: func(m sqlmock.Sqlmock) {},

		run: runFactory(
			[]storage.NetworkUpdateCriteria{
				{ID: "n1", DeleteNetwork: true},
				{ID: "n1", NewName: &wrappers.StringValue{Value: names[1]}},
			},
		),

		expectedError: errors.New("multiple updates for a single network are not allowed"),
	}

	runCase(t, happyPath)
	runCase(t, errorCase)
	runCase(t, validationFailure)
}

func TestSqlConfiguratorStorage_LoadEntities(t *testing.T) {
	runFactory := func(networkID string, filter storage.EntityLoadFilter, loadCriteria storage.EntityLoadCriteria) func(store storage.ConfiguratorStorage) (interface{}, error) {
		return func(store storage.ConfiguratorStorage) (interface{}, error) {
			return store.LoadEntities(networkID, filter, loadCriteria)
		}
	}

	// Empty load criteria
	basicOnly := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT ent.network_id, ent.pk, ent.\"key\", ent.type, ent.physical_id, ent.version, ent.graph_id FROM cfg_entities").
				WithArgs(
					"network", "bar", "foo",
					"network", "quz", "baz",
					"network", "world", "hello",
				).
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id"}).
						AddRow("network", "abc", "bar", "foo", nil, 2, "42").
						AddRow("network", "def", "quz", "baz", nil, 1, "42"),
				)
		},
		run: runFactory(
			"network",
			storage.EntityLoadFilter{
				IDs: []*storage.EntityID{
					{Type: "foo", Key: "bar"},
					{Type: "baz", Key: "quz"},
					{Type: "hello", Key: "world"},
				},
			},
			storage.EntityLoadCriteria{},
		),

		expectedResult: storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				{NetworkID: "network", Type: "baz", Key: "quz", GraphID: "42", Version: 1},
				{NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 2},
			},
			EntitiesNotFound: []*storage.EntityID{{Type: "hello", Key: "world"}},
		},
	}

	// Load everything, no assocs
	// foobar has 2 permissions, bazquz has 1 wildcarded permission
	loadEverything := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT .* FROM cfg_entities").
				WithArgs(
					"network", "bar", "foo",
					"network", "quz", "baz",
					"network", "world", "hello",
				).
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id", "name", "description", "config", "id", "scope", "permission", "type", "id_filter", "acl.version"}).
						AddRow("network", "abc", "bar", "foo", nil, 2, "42", "foobar", "foobar ent", []byte("foobar"), "foobar_acl_1", "n1,n2,n3", storage.ACL_OWN, "foo", nil, 1).
						AddRow("network", "abc", "bar", "foo", nil, 2, "42", "foobar", "foobar ent", []byte("foobar"), "foobar_acl_2", "n4", storage.ACL_READ, "baz", nil, 2).
						AddRow("network", "def", "quz", "baz", nil, 1, "42", "bazquz", "bazquz ent", []byte("bazquz"), "bazquz_acl_1", storage.ACL_WILDCARD_ALL.String(), storage.ACL_WRITE, storage.ACL_WILDCARD_ALL.String(), "1,2,3", 3),
				)
		},
		run: runFactory(
			"network",
			storage.EntityLoadFilter{
				IDs: []*storage.EntityID{
					{Type: "foo", Key: "bar"},
					{Type: "baz", Key: "quz"},
					{Type: "hello", Key: "world"},
				},
			},
			storage.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadPermissions: true},
		),

		expectedResult: storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "baz", Key: "quz", GraphID: "42", Version: 1,
					Name:        "bazquz",
					Description: "bazquz ent",
					Config:      []byte("bazquz"),
					Permissions: []*storage.ACL{
						{
							ID:         "bazquz_acl_1",
							Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
							Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
							Permission: storage.ACL_WRITE,
							IDFilter:   []string{"1", "2", "3"},
							Version:    3,
						},
					},
				},
				{
					NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 2,
					Name:        "foobar",
					Description: "foobar ent",
					Config:      []byte("foobar"),
					Permissions: []*storage.ACL{
						{
							ID:         "foobar_acl_1",
							Type:       &storage.ACL_EntityType{EntityType: "foo"},
							Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n1", "n2", "n3"}}},
							Permission: storage.ACL_OWN,
							Version:    1,
						},
						{
							ID:         "foobar_acl_2",
							Type:       &storage.ACL_EntityType{EntityType: "baz"},
							Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n4"}}},
							Permission: storage.ACL_READ,
							Version:    2,
						},
					},
				},
			},
			EntitiesNotFound: []*storage.EntityID{{Type: "hello", Key: "world"}},
		},
	}

	// Load assocs to only
	assocsTo := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT ent.network_id, ent.pk, ent.\"key\", ent.type, ent.physical_id, ent.version, ent.graph_id FROM cfg_entities").
				WithArgs(
					"network", "bar", "foo",
					"network", "quz", "baz",
					"network", "world", "hello",
				).
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id"}).
						AddRow("network", "abc", "bar", "foo", nil, 2, "42").
						AddRow("network", "def", "quz", "baz", nil, 1, "42").
						AddRow("network", "ghi", "world", "hello", nil, 3, "42"),
				)

			expectAssocQuery(
				m, []driver.Value{"abc", "def", "ghi"},
				"abc", "def",
				"abc", "ghi",
				"ghi", "def",
			)
		},
		run: runFactory(
			"network",
			storage.EntityLoadFilter{
				IDs: []*storage.EntityID{
					{Type: "foo", Key: "bar"},
					{Type: "baz", Key: "quz"},
					{Type: "hello", Key: "world"},
				},
			},
			storage.EntityLoadCriteria{LoadAssocsToThis: true},
		),

		expectedResult: storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "baz", Key: "quz", GraphID: "42", Version: 1,
					ParentAssociations: []*storage.EntityID{
						{Type: "foo", Key: "bar"},
						{Type: "hello", Key: "world"},
					},
				},
				{NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 2},
				{
					NetworkID: "network", Type: "hello", Key: "world", GraphID: "42", Version: 3,
					ParentAssociations: []*storage.EntityID{
						{Type: "foo", Key: "bar"},
					},
				},
			},
			EntitiesNotFound: []*storage.EntityID{},
		},
	}

	// Load assocs from only
	assocsFrom := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT ent.network_id, ent.pk, ent.\"key\", ent.type, ent.physical_id, ent.version, ent.graph_id FROM cfg_entities").
				WithArgs(
					"network", "bar", "foo",
					"network", "quz", "baz",
					"network", "world", "hello",
				).
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id"}).
						AddRow("network", "abc", "bar", "foo", nil, 2, "42").
						AddRow("network", "def", "quz", "baz", nil, 1, "42").
						AddRow("network", "ghi", "world", "hello", nil, 3, "42"),
				)

			expectAssocQuery(
				m,
				[]driver.Value{"abc", "def", "ghi"},
				"def", "abc",
				"ghi", "abc",
				"def", "ghi",
			)
		},
		run: runFactory(
			"network",
			storage.EntityLoadFilter{
				IDs: []*storage.EntityID{
					{Type: "foo", Key: "bar"},
					{Type: "baz", Key: "quz"},
					{Type: "hello", Key: "world"},
				},
			},
			storage.EntityLoadCriteria{LoadAssocsFromThis: true},
		),

		expectedResult: storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "baz", Key: "quz", GraphID: "42", Version: 1,
					Associations: []*storage.EntityID{
						{Type: "foo", Key: "bar"},
						{Type: "hello", Key: "world"},
					},
				},
				{NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 2},
				{
					NetworkID: "network", Type: "hello", Key: "world", GraphID: "42", Version: 3,
					Associations: []*storage.EntityID{
						{Type: "foo", Key: "bar"},
					},
				},
			},
			EntitiesNotFound: []*storage.EntityID{},
		},
	}

	// Load everything with type filter
	// (foo, bar) will have 2 permissions and 2 assocs - one to (foo, baz) and one to (bar, baz)
	// (foo, baz) will have 1 wildcard permission and 1 assoc to (baz, quz)
	// (hello, world) will have 1 assoc to (foo, bar)
	// We will only load entities of type foo
	fullLoadTypeFilter := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT .* FROM cfg_entities").
				WithArgs("network", "foo").
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id", "name", "description", "config", "id", "scope", "permission", "type", "id_filter", "acl.version"}).
						// (foo, bar) comes from test case for full load above
						AddRow("network", "foobar", "bar", "foo", nil, 1, "42", "foobar", "foobar ent", []byte("foobar"), "foobar_acl_1", "n1,n2,n3", storage.ACL_OWN, "foo", nil, 1).
						AddRow("network", "foobar", "bar", "foo", nil, 1, "42", "foobar", "foobar ent", []byte("foobar"), "foobar_acl_2", "n4", storage.ACL_READ, "baz", nil, 2).
						AddRow("network", "foobaz", "baz", "foo", nil, 2, "42", "foobaz", "foobaz ent", []byte("foobaz"), "foobaz_acl_1", "WILDCARD_ALL", storage.ACL_WRITE, "WILDCARD_ALL", nil, 3),
				)

			expectAssocQuery(
				m,
				[]driver.Value{"foobar", "foobaz", "foobar", "foobaz"},
				"foobar", "foobaz",
				"foobar", "barbaz",
				"foobaz", "bazquz",
				"helloworld", "foobar",
			)

			// Since we don't query for (hello, world), we expect a query for its type and key given its PK
			m.ExpectQuery("SELECT pk, type, \"key\" FROM cfg_entities").
				WithArgs("barbaz", "bazquz", "helloworld").
				WillReturnRows(
					sqlmock.NewRows([]string{"pk", "type", "key"}).
						AddRow("barbaz", "bar", "baz").
						AddRow("bazquz", "baz", "quz").
						AddRow("helloworld", "hello", "world"),
				)
		},
		run: runFactory(
			"network",
			storage.EntityLoadFilter{
				TypeFilter: &wrappers.StringValue{Value: "foo"},
			},
			storage.FullEntityLoadCriteria,
		),

		expectedResult: storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 1,
					Name:        "foobar",
					Description: "foobar ent",
					Config:      []byte("foobar"),
					Permissions: []*storage.ACL{
						{
							ID:         "foobar_acl_1",
							Type:       &storage.ACL_EntityType{EntityType: "foo"},
							Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n1", "n2", "n3"}}},
							Permission: storage.ACL_OWN,
							Version:    1,
						},
						{
							ID:         "foobar_acl_2",
							Type:       &storage.ACL_EntityType{EntityType: "baz"},
							Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n4"}}},
							Permission: storage.ACL_READ,
							Version:    2,
						},
					},
					Associations: []*storage.EntityID{
						{Type: "bar", Key: "baz"},
						{Type: "foo", Key: "baz"},
					},
					ParentAssociations: []*storage.EntityID{
						{Type: "hello", Key: "world"},
					},
				},
				{
					NetworkID: "network", Type: "foo", Key: "baz", GraphID: "42", Version: 2,
					Name:        "foobaz",
					Description: "foobaz ent",
					Config:      []byte("foobaz"),
					Permissions: []*storage.ACL{
						{
							ID:         "foobaz_acl_1",
							Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
							Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
							Permission: storage.ACL_WRITE,
							Version:    3,
						},
					},
					Associations: []*storage.EntityID{
						{Type: "baz", Key: "quz"},
					},
					ParentAssociations: []*storage.EntityID{
						{Type: "foo", Key: "bar"},
					},
				},
			},
			EntitiesNotFound: []*storage.EntityID{},
		},
	}

	// Basic load with type and key filters
	typeAndKeyFilters := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT ent.network_id, ent.pk, ent.\"key\", ent.type, ent.physical_id, ent.version, ent.graph_id FROM cfg_entities").
				WithArgs("network", "bar", "foo").
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id"}).
						AddRow("network", "abc", "bar", "foo", nil, 2, "42"),
				)
		},
		run: runFactory(
			"network",
			storage.EntityLoadFilter{
				TypeFilter: stringPointer("foo"),
				KeyFilter:  stringPointer("bar"),
			},
			storage.EntityLoadCriteria{},
		),

		expectedResult: storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				{NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 2},
			},
			EntitiesNotFound: []*storage.EntityID{},
		},
	}

	// Basic load with physical ID
	physicalID := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT ent.network_id, ent.pk, ent.\"key\", ent.type, ent.physical_id, ent.version, ent.graph_id FROM cfg_entities").
				WithArgs("p1").
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id"}).
						AddRow("network", "abc", "bar", "foo", "p1", 2, "42"),
				)
		},
		run: runFactory(
			"network",
			storage.EntityLoadFilter{
				PhysicalID: stringPointer("p1"),
			},
			storage.EntityLoadCriteria{},
		),

		expectedResult: storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				{NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", PhysicalID: "p1", Version: 2},
			},
			EntitiesNotFound: []*storage.EntityID{},
		},
	}

	runCase(t, basicOnly)
	runCase(t, loadEverything)
	runCase(t, assocsTo)
	runCase(t, assocsFrom)
	runCase(t, fullLoadTypeFilter)
	runCase(t, typeAndKeyFilters)
	runCase(t, physicalID)
}

func TestSqlConfiguratorStorage_CreateEntity(t *testing.T) {
	runFactory := func(networkID string, entity storage.NetworkEntity) func(store storage.ConfiguratorStorage) (interface{}, error) {
		return func(store storage.ConfiguratorStorage) (interface{}, error) {
			return store.CreateEntity(networkID, entity)
		}
	}

	// basic fields
	basicCase := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT COUNT\(1\) FROM cfg_entities`).
				WithArgs("network", "foo", "bar").
				WillReturnRows(sqlmock.NewRows([]string{"count"}))

			m.ExpectExec("INSERT INTO cfg_entities").
				WithArgs("1", "network", "foo", "bar", "2", "foobar", "foobar ent", nil, nil).
				WillReturnResult(mockResult)
		},
		run: runFactory(
			"network",
			storage.NetworkEntity{
				Type:        "foo",
				Key:         "bar",
				Name:        "foobar",
				Description: "foobar ent",
			},
		),

		expectedResult: storage.NetworkEntity{
			NetworkID:   "network",
			Type:        "foo",
			Key:         "bar",
			Name:        "foobar",
			Description: "foobar ent",
			GraphID:     "2",
		},
	}

	perms := []*storage.ACL{
		{
			ID:         "ignore this",
			Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
			Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
			Permission: storage.ACL_WRITE,
		},
		{
			Type:       &storage.ACL_EntityType{EntityType: "foo"},
			Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n1", "n2"}}},
			Permission: storage.ACL_READ,
			IDFilter:   []string{"foo", "bar"},
		},
	}

	// with permissions
	withPerms := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT COUNT\(1\) FROM cfg_entities`).
				WithArgs("network", "foo", "bar").
				WillReturnRows(sqlmock.NewRows([]string{"count"}))

			m.ExpectExec("INSERT INTO cfg_entities").
				WithArgs("1", "network", "foo", "bar", "2", "foobar", "foobar ent", nil, nil).
				WillReturnResult(mockResult)

			expectPermissionCreation(m, "1", 3, perms...)
		},
		run: runFactory(
			"network",
			storage.NetworkEntity{
				Type:        "foo",
				Key:         "bar",
				Name:        "foobar",
				Description: "foobar ent",
				Permissions: perms,
			},
		),

		expectedResult: storage.NetworkEntity{
			NetworkID:   "network",
			Type:        "foo",
			Key:         "bar",
			Name:        "foobar",
			Description: "foobar ent",
			GraphID:     "2",
			Permissions: []*storage.ACL{
				{
					ID:         "3",
					Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
					Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
					Permission: storage.ACL_WRITE,
				},
				{
					ID:         "4",
					Type:       &storage.ACL_EntityType{EntityType: "foo"},
					Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n1", "n2"}}},
					Permission: storage.ACL_READ,
					IDFilter:   []string{"foo", "bar"},
				},
			},
		},
	}

	// merge 2 graphs together
	mergeGraphs := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT COUNT\(1\) FROM cfg_entities`).
				WithArgs("network", "foo", "bar").
				WillReturnRows(sqlmock.NewRows([]string{"count"}))

			m.ExpectExec("INSERT INTO cfg_entities").
				WithArgs("1", "network", "foo", "bar", "2", "foobar", "foobar ent", nil, nil).
				WillReturnResult(mockResult)

			assocs := []*storage.EntityID{{Type: "bar", Key: "baz"}, {Type: "baz", Key: "quz"}}
			edgesByTk := map[storage2.TypeAndKey]expectedEntQueryResult{
				{Type: "bar", Key: "baz"}: {"bar", "baz", "42", "", "1", 1},
				{Type: "baz", Key: "quz"}: {"baz", "quz", "43", "", "3", 2},
			}
			expectEdgeQueries(m, assocs, edgesByTk)
			expectEdgeInsertions(m, assocsToEdges("1", assocs, edgesByTk))
			expectMergeGraphs(m, [][2]string{{"2", "1"}, {"3", "1"}})
		},
		run: runFactory(
			"network",
			storage.NetworkEntity{
				Type:        "foo",
				Key:         "bar",
				Name:        "foobar",
				Description: "foobar ent",
				Associations: []*storage.EntityID{
					{Type: "bar", Key: "baz"},
					{Type: "baz", Key: "quz"},
					// Duplicate edge should only be loaded once
					{Type: "bar", Key: "baz"},
				},
			},
		),

		// We expect "aaa" as the returned graphID since we merged graphs
		expectedResult: storage.NetworkEntity{
			NetworkID:   "network",
			Type:        "foo",
			Key:         "bar",
			Name:        "foobar",
			Description: "foobar ent",
			GraphID:     "1",
			Associations: []*storage.EntityID{
				{Type: "bar", Key: "baz"},
				{Type: "baz", Key: "quz"},
			},
		},
	}

	// already exists
	alreadyExists := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT COUNT\(1\) FROM cfg_entities`).
				WithArgs("network", "foo", "bar").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		},
		run: runFactory(
			"network",
			storage.NetworkEntity{
				Type:        "foo",
				Key:         "bar",
				Name:        "foobar",
				Description: "foobar ent",
			},
		),

		expectedResult: storage.NetworkEntity{},
		expectedError:  errors.New("an entity (foo-bar) already exists"),
	}

	runCase(t, basicCase)
	runCase(t, withPerms)
	runCase(t, mergeGraphs)
	runCase(t, alreadyExists)
}

func TestSqlConfiguratorStorage_UpdateEntity(t *testing.T) {
	runFactory := func(networkID string, update storage.EntityUpdateCriteria) func(store storage.ConfiguratorStorage) (interface{}, error) {
		return func(store storage.ConfiguratorStorage) (interface{}, error) {
			return store.UpdateEntity(networkID, update)
		}
	}

	// Delete entity
	expectedFooBarQuery := expectedEntQueryResult{"foo", "bar", "1", "", "g1", 0}
	deleteCase := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			expectBasicEntityQueries(m, expectedFooBarQuery)
			m.ExpectExec("DELETE FROM cfg_entities").WithArgs("network", "foo", "bar").WillReturnResult(mockResult)
			expectBulkEntityQuery(m, []driver.Value{"g1"})
		},
		run: runFactory("network", storage.EntityUpdateCriteria{Type: "foo", Key: "bar", DeleteEntity: true}),

		expectedResult: storage.NetworkEntity{Type: "foo", Key: "bar"},
	}
	runCase(t, deleteCase)

	// Delete entity and partition a graph
	deleteWithPartition := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			expectBasicEntityQueries(m, expectedFooBarQuery)
			m.ExpectExec("DELETE FROM cfg_entities").WithArgs("network", "foo", "bar").WillReturnResult(mockResult)
			// make foobar the root of a tree so we partition the graph into
			// 3 components:
			// foobar -> ( barbaz -> [bazquz | quzbaz] | bazbar | barfoo )
			expectBulkEntityQuery(
				m,
				[]driver.Value{"g1"},
				expectedEntQueryResult{"bar", "baz", "barbaz", "", "g1", 0},
				expectedEntQueryResult{"baz", "quz", "bazquz", "", "g1", 0},
				expectedEntQueryResult{"quz", "baz", "quzbaz", "", "g1", 0},
				expectedEntQueryResult{"baz", "bar", "bazbar", "", "g1", 0},
				expectedEntQueryResult{"bar", "foo", "barfoo", "", "g1", 0},
			)
			expectAssocQuery(
				m,
				[]driver.Value{"barbaz", "barfoo", "bazbar", "bazquz", "quzbaz", "barbaz", "barfoo", "bazbar", "bazquz", "quzbaz"},
				"barbaz", "bazquz",
				"barbaz", "quzbaz",
			)
			m.ExpectExec("UPDATE cfg_entities").WithArgs("1", "barfoo").WillReturnResult(mockResult)
			m.ExpectExec("UPDATE cfg_entities").WithArgs("2", "bazbar").WillReturnResult(mockResult)
		},
		run:            runFactory("network", storage.EntityUpdateCriteria{Type: "foo", Key: "bar", DeleteEntity: true}),
		expectedResult: storage.NetworkEntity{Type: "foo", Key: "bar"},
	}
	runCase(t, deleteWithPartition)

	// Test some permutations of updating basic fields
	runCase(
		t,
		getTestCaseForEntityUpdate(
			storage.EntityUpdateCriteria{
				Type:      "foo",
				Key:       "bar",
				NewName:   stringPointer("foobar"),
				NewConfig: bytesPointer([]byte("foobar config")),
			},
			expectedFooBarQuery,
			nil,
		),
	)
	runCase(
		t,
		getTestCaseForEntityUpdate(
			storage.EntityUpdateCriteria{
				Type:           "foo",
				Key:            "bar",
				NewDescription: stringPointer("foobar desc"),
				NewPhysicalID:  stringPointer("phys2"),
			},
			expectedFooBarQuery,
			nil,
		),
	)
	runCase(
		t,
		getTestCaseForEntityUpdate(
			storage.EntityUpdateCriteria{
				Type:           "foo",
				Key:            "bar",
				NewName:        stringPointer("foobar"),
				NewDescription: stringPointer("foobar desc"),
				NewPhysicalID:  stringPointer("phys2"),
				NewConfig:      bytesPointer([]byte("foobar config")),
			},
			expectedFooBarQuery,
			nil,
		),
	)

	// Test cases for permissions
	runCase(
		t,
		getTestCaseForEntityUpdate(
			storage.EntityUpdateCriteria{
				Type: "foo",
				Key:  "bar",
				PermissionsToCreate: []*storage.ACL{
					{
						Permission: storage.ACL_WRITE,
						Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
						Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n3"}}},
					},
				},
			},
			expectedFooBarQuery,
			nil,
		),
	)

	runCase(
		t,
		getTestCaseForEntityUpdate(
			storage.EntityUpdateCriteria{
				Type: "foo",
				Key:  "bar",
				PermissionsToUpdate: []*storage.ACL{
					{
						ID:         "42",
						Permission: storage.ACL_READ,
						Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
						Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
						IDFilter:   []string{"n1", "n2"},
					},
					{
						ID:         "43",
						Permission: storage.ACL_WRITE,
						Type:       &storage.ACL_EntityType{EntityType: "bar"},
						Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n3"}}},
					},
				},
			},
			expectedFooBarQuery,
			nil,
		),
	)

	runCase(
		t,
		getTestCaseForEntityUpdate(
			storage.EntityUpdateCriteria{
				Type:                "foo",
				Key:                 "bar",
				PermissionsToDelete: []string{"100", "101"},
			},
			expectedFooBarQuery,
			nil,
		),
	)

	runCase(
		t,
		getTestCaseForEntityUpdate(
			storage.EntityUpdateCriteria{
				Type: "foo",
				Key:  "bar",
				PermissionsToCreate: []*storage.ACL{
					{
						ID:         "ignore me",
						Permission: storage.ACL_WRITE,
						Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
						Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
						IDFilter:   []string{"n1", "n2"},
					},
					{
						Permission: storage.ACL_READ,
						Type:       &storage.ACL_EntityType{EntityType: "foo"},
						Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n1", "n2"}}},
					},
				},
				PermissionsToUpdate: []*storage.ACL{
					{
						ID:         "42",
						Permission: storage.ACL_READ,
						Type:       &storage.ACL_TypeWildcard{TypeWildcard: storage.ACL_WILDCARD_ALL},
						Scope:      &storage.ACL_ScopeWildcard{ScopeWildcard: storage.ACL_WILDCARD_ALL},
						IDFilter:   []string{"n1", "n2"},
					},
					{
						ID:         "43",
						Permission: storage.ACL_WRITE,
						Type:       &storage.ACL_EntityType{EntityType: "bar"},
						Scope:      &storage.ACL_ScopeNetworkIDs{ScopeNetworkIDs: &storage.ACL_NetworkIDs{IDs: []string{"n3"}}},
					},
				},
				PermissionsToDelete: []string{"42", "101"},
			},
			expectedFooBarQuery,
			nil,
		),
	)

	// edges

	// Set edges, merge graphs
	runCase(
		t,
		getTestCaseForEntityUpdate(
			storage.EntityUpdateCriteria{
				Type: "foo",
				Key:  "bar",
				AssociationsToSet: &storage.EntityAssociationsToSet{
					AssociationsToSet: []*storage.EntityID{
						{Type: "bar", Key: "baz"},
						{Type: "baz", Key: "quz"},
					},
				},
			},
			expectedEntQueryResult{"foo", "bar", "1", "", "g1", 0},
			[]expectedEntQueryResult{
				{"bar", "baz", "2", "", "g2", 0},
				{"baz", "quz", "3", "", "g3", 0},
			},
			[2]string{"g2", "g1"},
			[2]string{"g3", "g1"},
		),
	)

	// Create edges, merge graphs
	runCase(
		t,
		getTestCaseForEntityUpdate(
			storage.EntityUpdateCriteria{
				Type: "foo",
				Key:  "bar",
				AssociationsToAdd: []*storage.EntityID{
					{Type: "bar", Key: "baz"},
					{Type: "baz", Key: "quz"},
				},
			},
			expectedEntQueryResult{"foo", "bar", "1", "", "g9", 0},
			[]expectedEntQueryResult{
				{"bar", "baz", "2", "", "g1", 0},
				{"baz", "quz", "3", "", "g2", 0},
			},
			[2]string{"g2", "g1"},
			[2]string{"g9", "g1"},
		),
	)

	// Create edge to something already in the same graph
	runCase(
		t,
		getTestCaseForEntityUpdate(
			storage.EntityUpdateCriteria{
				Type: "foo",
				Key:  "bar",
				AssociationsToAdd: []*storage.EntityID{
					{Type: "bar", Key: "baz"},
				},
			},
			expectedFooBarQuery,
			[]expectedEntQueryResult{
				{"bar", "baz", "2", "", "g1", 0},
			},
		),
	)

	// Delete edge, no fix graph
	runCase(
		t,
		getTestCaseForEntityUpdate(
			storage.EntityUpdateCriteria{
				Type: "foo",
				Key:  "bar",
				AssociationsToDelete: []*storage.EntityID{
					{Type: "bar", Key: "baz"},
				},
			},
			expectedFooBarQuery,
			[]expectedEntQueryResult{
				{"bar", "baz", "2", "", "g1", 0},
			},
		),
	)

	// Graph partition:
	// foobar -> barbaz -> bazquz -> (quzbaz, bazbar -> barfoo)
	// delete edges bazquz -> quzbaz, bazquz -> bazbar
	// partitions graph into 3 components
	partitionCase := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			expectBasicEntityQueries(m, getBasicQueryExpect("baz", "quz"))
			m.ExpectExec("UPDATE cfg_entities").WithArgs("bazquz").WillReturnResult(mockResult)
			expectEdgeQueries(
				m,
				[]*storage.EntityID{{Type: "quz", Key: "baz"}, {Type: "baz", Key: "bar"}},
				map[storage2.TypeAndKey]expectedEntQueryResult{
					{"quz", "baz"}: getBasicQueryExpect("quz", "baz"),
					{"baz", "bar"}: getBasicQueryExpect("baz", "bar"),
				},
			)
			expectEdgeDeletions(m, [][2]string{{"bazquz", "quzbaz"}, {"bazquz", "bazbar"}})

			expectBulkEntityQuery(
				m,
				[]driver.Value{"g1"},
				getBasicQueryExpect("foo", "bar"),
				getBasicQueryExpect("bar", "baz"),
				getBasicQueryExpect("baz", "quz"),
				getBasicQueryExpect("quz", "baz"),
				getBasicQueryExpect("baz", "bar"),
				getBasicQueryExpect("bar", "foo"),
			)
			expectAssocQuery(
				m,
				[]driver.Value{"barbaz", "barfoo", "bazbar", "bazquz", "foobar", "quzbaz", "barbaz", "barfoo", "bazbar", "bazquz", "foobar", "quzbaz"},
				"foobar", "barbaz",
				"barbaz", "bazquz",
				"bazbar", "barfoo",
			)
			m.ExpectExec("UPDATE cfg_entities").WithArgs("1", "quzbaz").WillReturnResult(mockResult)
			m.ExpectExec("UPDATE cfg_entities").WithArgs("2", "barfoo", "bazbar").WillReturnResult(mockResult)
		},
		run:            runFactory("network", storage.EntityUpdateCriteria{Type: "baz", Key: "quz", AssociationsToDelete: []*storage.EntityID{{Type: "quz", Key: "baz"}, {Type: "baz", Key: "bar"}}}),
		expectedResult: storage.NetworkEntity{NetworkID: "network", Type: "baz", Key: "quz", GraphID: "g1", Version: 1},
	}
	runCase(t, partitionCase)

	clearEdgesCase := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			// Load and change version, then clear assocs
			expectBasicEntityQueries(m, getBasicQueryExpect("foo", "bar"))
			m.ExpectExec("UPDATE cfg_entities").WithArgs("foobar").WillReturnResult(mockResult)
			m.ExpectExec("DELETE FROM cfg_assocs").WithArgs("foobar").WillReturnResult(mockResult)

			// Graph loading
			expectBulkEntityQuery(
				m,
				[]driver.Value{"g1"},
				getBasicQueryExpect("foo", "bar"),
				getBasicQueryExpect("bar", "baz"),
			)
			expectAssocQuery(m, []driver.Value{"barbaz", "foobar", "barbaz", "foobar"})

			// Graph partition update
			m.ExpectExec("UPDATE cfg_entities").WithArgs("1", "barbaz").WillReturnResult(mockResult)
		},
		run:            runFactory("network", storage.EntityUpdateCriteria{Type: "foo", Key: "bar", AssociationsToSet: &storage.EntityAssociationsToSet{}}),
		expectedResult: storage.NetworkEntity{NetworkID: "network", Type: "foo", Key: "bar", GraphID: "g1", Version: 1},
	}
	runCase(t, clearEdgesCase)

}

func TestSqlConfiguratorStorage_LoadGraphForEntity(t *testing.T) {
	runFactory := func(networkID string, entityID storage.EntityID, loadCriteria storage.EntityLoadCriteria) func(store storage.ConfiguratorStorage) (interface{}, error) {
		return func(store storage.ConfiguratorStorage) (interface{}, error) {
			return store.LoadGraphForEntity(networkID, entityID, loadCriteria)
		}
	}

	expectedFooBar := expectedEntQueryResult{"foo", "bar", "foobar", "", "g1", 0}
	expectedBarBaz := expectedEntQueryResult{"bar", "baz", "barbaz", "p1", "g1", 1}
	expectedBazQuz := expectedEntQueryResult{"baz", "quz", "bazquz", "p2", "g1", 2}

	assocQueryArgs := []driver.Value{"barbaz", "bazquz", "foobar", "barbaz", "bazquz", "foobar"}

	// load a linked list of 3 nodes
	linkedList := &testCase{
		run: runFactory("network", storage.EntityID{Type: "foo", Key: "bar"}, storage.EntityLoadCriteria{}),

		setup: func(m sqlmock.Sqlmock) {
			expectBasicEntityQueries(m, expectedFooBar)
			expectBulkEntityQuery(m, []driver.Value{"g1"}, expectedFooBar, expectedBarBaz, expectedBazQuz)
			// foobar -> barbaz -> bazquz
			expectAssocQuery(m, assocQueryArgs, "foobar", "barbaz", "barbaz", "bazquz")
		},

		expectedResult: storage.EntityGraph{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "bar", Key: "baz",
					PhysicalID: "p1", GraphID: "g1",
					Associations:       []*storage.EntityID{{Type: "baz", Key: "quz"}},
					ParentAssociations: []*storage.EntityID{{Type: "foo", Key: "bar"}},
					Version:            1,
				},
				{
					NetworkID: "network", Type: "baz", Key: "quz",
					PhysicalID: "p2", GraphID: "g1",
					ParentAssociations: []*storage.EntityID{{Type: "bar", Key: "baz"}},
					Version:            2,
				},
				{
					NetworkID: "network", Type: "foo", Key: "bar",
					GraphID:      "g1",
					Associations: []*storage.EntityID{{Type: "bar", Key: "baz"}},
				},
			},
			RootEntities: []*storage.EntityID{{Type: "foo", Key: "bar"}},
			Edges: []*storage.GraphEdge{
				{From: &storage.EntityID{Type: "bar", Key: "baz"}, To: &storage.EntityID{Type: "baz", Key: "quz"}},
				{From: &storage.EntityID{Type: "foo", Key: "bar"}, To: &storage.EntityID{Type: "bar", Key: "baz"}},
			},
		},
	}

	// load a simple tree of 3 nodes
	tree := &testCase{
		run: runFactory("network", storage.EntityID{Type: "baz", Key: "quz"}, storage.EntityLoadCriteria{}),

		setup: func(m sqlmock.Sqlmock) {
			expectBasicEntityQueries(m, expectedBazQuz)
			expectBulkEntityQuery(m, []driver.Value{"g1"}, expectedFooBar, expectedBarBaz, expectedBazQuz)
			// foobar -> barbaz; foobar -> bazquz
			expectAssocQuery(m, assocQueryArgs, "foobar", "barbaz", "foobar", "bazquz")
		},

		expectedResult: storage.EntityGraph{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "bar", Key: "baz",
					PhysicalID: "p1", GraphID: "g1",
					ParentAssociations: []*storage.EntityID{{Type: "foo", Key: "bar"}},
					Version:            1,
				},
				{
					NetworkID: "network", Type: "baz", Key: "quz",
					PhysicalID: "p2", GraphID: "g1",
					ParentAssociations: []*storage.EntityID{{Type: "foo", Key: "bar"}},
					Version:            2,
				},
				{
					NetworkID: "network", Type: "foo", Key: "bar",
					GraphID:      "g1",
					Associations: []*storage.EntityID{{Type: "bar", Key: "baz"}, {Type: "baz", Key: "quz"}},
				},
			},
			RootEntities: []*storage.EntityID{{Type: "foo", Key: "bar"}},
			Edges: []*storage.GraphEdge{
				{From: &storage.EntityID{Type: "foo", Key: "bar"}, To: &storage.EntityID{Type: "bar", Key: "baz"}},
				{From: &storage.EntityID{Type: "foo", Key: "bar"}, To: &storage.EntityID{Type: "baz", Key: "quz"}},
			},
		},
	}

	// load an upside-down tree
	upsideDownTree := &testCase{
		run: runFactory("network", storage.EntityID{Type: "foo", Key: "bar"}, storage.EntityLoadCriteria{}),

		setup: func(m sqlmock.Sqlmock) {
			expectBasicEntityQueries(m, expectedFooBar)
			expectBulkEntityQuery(m, []driver.Value{"g1"}, expectedFooBar, expectedBarBaz, expectedBazQuz)
			// barbaz -> foobar; bazquz -> foobar
			expectAssocQuery(m, assocQueryArgs, "barbaz", "foobar", "bazquz", "foobar")
		},

		expectedResult: storage.EntityGraph{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "bar", Key: "baz",
					PhysicalID: "p1", GraphID: "g1",
					Associations: []*storage.EntityID{{Type: "foo", Key: "bar"}},
					Version:      1,
				},
				{
					NetworkID: "network", Type: "baz", Key: "quz",
					PhysicalID: "p2", GraphID: "g1",
					Associations: []*storage.EntityID{{Type: "foo", Key: "bar"}},
					Version:      2,
				},
				{
					NetworkID: "network", Type: "foo", Key: "bar",
					GraphID:            "g1",
					ParentAssociations: []*storage.EntityID{{Type: "bar", Key: "baz"}, {Type: "baz", Key: "quz"}},
				},
			},
			RootEntities: []*storage.EntityID{{Type: "bar", Key: "baz"}, {Type: "baz", Key: "quz"}},
			Edges: []*storage.GraphEdge{
				{From: &storage.EntityID{Type: "bar", Key: "baz"}, To: &storage.EntityID{Type: "foo", Key: "bar"}},
				{From: &storage.EntityID{Type: "baz", Key: "quz"}, To: &storage.EntityID{Type: "foo", Key: "bar"}},
			},
		},
	}

	// load a graph with a cycle
	withCycle := &testCase{
		run: runFactory("network", storage.EntityID{Type: "foo", Key: "bar"}, storage.EntityLoadCriteria{}),

		setup: func(m sqlmock.Sqlmock) {
			expectBasicEntityQueries(m, expectedFooBar)
			expectBulkEntityQuery(m, []driver.Value{"g1"}, expectedFooBar, expectedBarBaz, expectedBazQuz)
			// foobar -> barbaz; barbaz <-> bazquz
			expectAssocQuery(m, assocQueryArgs, "foobar", "barbaz", "barbaz", "bazquz", "bazquz", "barbaz")
		},

		expectedResult: storage.EntityGraph{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "bar", Key: "baz",
					PhysicalID: "p1", GraphID: "g1",
					Associations:       []*storage.EntityID{{Type: "baz", Key: "quz"}},
					ParentAssociations: []*storage.EntityID{{Type: "baz", Key: "quz"}, {Type: "foo", Key: "bar"}},
					Version:            1,
				},
				{
					NetworkID: "network", Type: "baz", Key: "quz",
					PhysicalID: "p2", GraphID: "g1",
					Associations:       []*storage.EntityID{{Type: "bar", Key: "baz"}},
					ParentAssociations: []*storage.EntityID{{Type: "bar", Key: "baz"}},
					Version:            2,
				},
				{
					NetworkID: "network", Type: "foo", Key: "bar",
					GraphID:      "g1",
					Associations: []*storage.EntityID{{Type: "bar", Key: "baz"}},
				},
			},
			RootEntities: []*storage.EntityID{{Type: "foo", Key: "bar"}},
			Edges: []*storage.GraphEdge{
				{From: &storage.EntityID{Type: "bar", Key: "baz"}, To: &storage.EntityID{Type: "baz", Key: "quz"}},
				{From: &storage.EntityID{Type: "baz", Key: "quz"}, To: &storage.EntityID{Type: "bar", Key: "baz"}},
				{From: &storage.EntityID{Type: "foo", Key: "bar"}, To: &storage.EntityID{Type: "bar", Key: "baz"}},
			},
		},
	}

	// load a ring
	ring := &testCase{
		run: runFactory("network", storage.EntityID{Type: "foo", Key: "bar"}, storage.EntityLoadCriteria{}),

		setup: func(m sqlmock.Sqlmock) {
			expectBasicEntityQueries(m, expectedFooBar)
			expectBulkEntityQuery(m, []driver.Value{"g1"}, expectedFooBar, expectedBarBaz, expectedBazQuz)
			// foobar -> barbaz -> bazquz -> foobar -> ...
			expectAssocQuery(m, assocQueryArgs, "foobar", "barbaz", "barbaz", "bazquz", "bazquz", "foobar")
		},

		expectedError: errors.New("graph does not have root nodes because it is a ring"),
	}

	runCase(t, linkedList)
	runCase(t, tree)
	runCase(t, upsideDownTree)
	runCase(t, withCycle)
	runCase(t, ring)
}

type testCase struct {
	// setup mock expectations. Transaction start is expected on the mock
	// generically
	setup func(m sqlmock.Sqlmock)

	// run the test case
	run func(store storage.ConfiguratorStorage) (interface{}, error)

	expectedError      error
	matchErrorInstance bool
	expectedResult     interface{}
}

func runCase(t *testing.T, test *testCase) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Printf("error closing stub DB: %s", err)
		}
	}()

	mock.ExpectBegin()
	test.setup(mock)

	factory := storage.NewSQLConfiguratorStorageFactory(db, &mockIDGenerator{}, sqorc.GetSqlBuilder())
	store, err := factory.StartTransaction(context.Background(), nil)
	assert.NoError(t, err)
	actual, err := test.run(store)

	if test.expectedError != nil {
		if test.matchErrorInstance {
			assert.True(t, err == test.expectedError)
		}
		assert.EqualError(t, err, test.expectedError.Error())
	} else {
		assert.NoError(t, err)
	}

	if test.expectedResult != nil {
		assert.Equal(t, test.expectedResult, actual)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

// this is a lot more coupled to implementation than I like for test code
// but the alternative of db test fixtures and manually querying/dumping the db
// isn't much better, if we want true unit tests for the storage impl.
// For now, this is just the happy path, with no graph partitioning.
func getTestCaseForEntityUpdate(
	update storage.EntityUpdateCriteria,
	entToUpdate expectedEntQueryResult,
	expectedEdgeLoads []expectedEntQueryResult,
	expectedGraphMerges ...[2]string,
) *testCase {
	expectedResult := storage.NetworkEntity{
		NetworkID: "network",
		Type:      entToUpdate.entType,
		Key:       entToUpdate.key,
		GraphID:   entToUpdate.graphID,
		Version:   entToUpdate.version + 1,

		Name:        stringValue(update.NewName),
		Description: stringValue(update.NewDescription),
		Config:      bytesVal(update.NewConfig),
	}
	if update.NewPhysicalID != nil {
		expectedResult.PhysicalID = update.NewPhysicalID.Value
	} else {
		expectedResult.PhysicalID = entToUpdate.physicalID
	}

	delPermIdsSet := funk.Map(update.PermissionsToDelete, func(i string) (string, bool) { return i, true }).(map[string]bool)
	permID := 1
	for _, perm := range update.PermissionsToCreate {
		perm.ID = fmt.Sprintf("%d", permID)
		expectedResult.Permissions = append(expectedResult.Permissions, perm)
		permID++
	}
	for _, perm := range update.PermissionsToUpdate {
		if _, wasDeleted := delPermIdsSet[perm.ID]; !wasDeleted {
			expectedResult.Permissions = append(expectedResult.Permissions, perm)
		}
	}

	edgeLoadsByTk := funk.Map(
		expectedEdgeLoads,
		func(e expectedEntQueryResult) (storage2.TypeAndKey, expectedEntQueryResult) {
			return storage2.TypeAndKey{Type: e.entType, Key: e.key}, e
		},
	).(map[storage2.TypeAndKey]expectedEntQueryResult)

	if !funk.IsEmpty(update.AssociationsToAdd) {
		expectedResult.Associations = append(expectedResult.Associations, update.AssociationsToAdd...)
	}
	if update.AssociationsToSet != nil {
		expectedResult.Associations = append(expectedResult.Associations, update.AssociationsToSet.AssociationsToSet...)
	}

	if !funk.IsEmpty(expectedGraphMerges) {
		expectedResult.GraphID = expectedGraphMerges[0][1]
	}

	return &testCase{
		setup: func(m sqlmock.Sqlmock) {
			// Basic fields
			expectBasicEntityQueries(m, entToUpdate)
			updateWithArgs := []driver.Value{}
			if update.NewName != nil {
				updateWithArgs = append(updateWithArgs, update.NewName.Value)
			}
			if update.NewDescription != nil {
				updateWithArgs = append(updateWithArgs, update.NewDescription.Value)
			}
			if update.NewPhysicalID != nil {
				updateWithArgs = append(updateWithArgs, update.NewPhysicalID.Value)
			}
			if update.NewConfig != nil {
				updateWithArgs = append(updateWithArgs, update.NewConfig.Value)
			}
			updateWithArgs = append(updateWithArgs, entToUpdate.pk)

			m.ExpectExec("UPDATE cfg_entities").WithArgs(updateWithArgs...).WillReturnResult(mockResult)

			// Permissions
			if !funk.IsEmpty(update.PermissionsToCreate) {
				expectPermissionCreation(m, entToUpdate.pk, 1, update.PermissionsToCreate...)
			}
			if !funk.IsEmpty(update.PermissionsToUpdate) {
				m.ExpectQuery(`SELECT COUNT\(\*\) FROM cfg_acls`).WithArgs(funk.Map(update.PermissionsToUpdate, func(acl *storage.ACL) driver.Value { return acl.ID }).([]driver.Value)...).
					WillReturnRows(
						sqlmock.NewRows([]string{"count"}).
							AddRow(len(update.PermissionsToUpdate)),
					)
				expectPermissionUpdates(m, entToUpdate.pk, update.PermissionsToUpdate...)
			}
			if !funk.IsEmpty(update.PermissionsToDelete) {
				expectPermissionDeletes(m, update.PermissionsToDelete...)
			}

			// Graph
			if update.AssociationsToSet != nil {
				m.ExpectExec("DELETE FROM cfg_assocs").WithArgs(entToUpdate.pk).WillReturnResult(mockResult)
				expectEdgeQueries(m, update.AssociationsToSet.AssociationsToSet, edgeLoadsByTk)
				expectEdgeInsertions(m, assocsToEdges(entToUpdate.pk, update.AssociationsToSet.AssociationsToSet, edgeLoadsByTk))
				if !funk.IsEmpty(expectedGraphMerges) {
					expectMergeGraphs(m, expectedGraphMerges)
				}

				// fix graph, but no partition detected
				expectBulkEntityQuery(m, []driver.Value{entToUpdate.graphID}, entToUpdate)
				expectAssocQuery(m, []driver.Value{entToUpdate.pk, entToUpdate.pk})
			}

			if !funk.IsEmpty(update.AssociationsToAdd) {
				expectEdgeQueries(m, update.AssociationsToAdd, edgeLoadsByTk)
				expectEdgeInsertions(m, assocsToEdges(entToUpdate.pk, update.AssociationsToAdd, edgeLoadsByTk))
				if !funk.IsEmpty(expectedGraphMerges) {
					expectMergeGraphs(m, expectedGraphMerges)
				}
			}

			if !funk.IsEmpty(update.AssociationsToDelete) {
				expectEdgeQueries(m, update.AssociationsToDelete, edgeLoadsByTk)
				expectEdgeDeletions(m, assocsToEdges(entToUpdate.pk, update.AssociationsToDelete, edgeLoadsByTk))

				// fix graph, but no partition detected
				expectBulkEntityQuery(m, []driver.Value{entToUpdate.graphID}, entToUpdate)
				expectAssocQuery(m, []driver.Value{entToUpdate.pk, entToUpdate.pk})
			}
		},
		run: func(store storage.ConfiguratorStorage) (interface{}, error) {
			return store.UpdateEntity("network", update)
		},
		expectedResult: expectedResult,
	}
}

type expectedEntQueryResult struct {
	entType, key, pk, physicalID, graphID string
	version                               uint64
}

func expectBasicEntityQueries(m sqlmock.Sqlmock, expectations ...expectedEntQueryResult) {
	args := make([]driver.Value, 0, len(expectations)*3)
	for _, expect := range expectations {
		args = append(args, "network", expect.key, expect.entType)
	}
	m.ExpectQuery("SELECT .* FROM cfg_entities").WithArgs(args...).WillReturnRows(expectedEntQueriesToRows(expectations...))
}

func expectBulkEntityQuery(m sqlmock.Sqlmock, queryArgs []driver.Value, expectations ...expectedEntQueryResult) {
	m.ExpectQuery("SELECT .* FROM cfg_entities").WithArgs(queryArgs...).
		WillReturnRows(expectedEntQueriesToRows(expectations...))
}

func expectedEntQueriesToRows(expectations ...expectedEntQueryResult) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id"})
	for _, expect := range expectations {
		rowValues := make([]driver.Value, 0, 6)
		rowValues = append(rowValues, "network", expect.pk, expect.key, expect.entType)
		if expect.physicalID == "" {
			rowValues = append(rowValues, nil)
		} else {
			rowValues = append(rowValues, expect.physicalID)
		}
		rowValues = append(rowValues, expect.version, expect.graphID)
		rows.AddRow(rowValues...)
	}
	return rows
}

func expectAssocQuery(m sqlmock.Sqlmock, queryArgs []driver.Value, assocPks ...string) {
	rows := sqlmock.NewRows([]string{"from_pk", "to_pk"})
	for i := 0; i < len(assocPks); i += 2 {
		rows.AddRow(assocPks[i], assocPks[i+1])
	}
	m.ExpectQuery("SELECT assoc.from_pk, assoc.to_pk FROM cfg_assocs").WithArgs(queryArgs...).WillReturnRows(rows)
}

// [(old graph ID, new graph ID)]
func expectMergeGraphs(m sqlmock.Sqlmock, graphIDChanges [][2]string) {
	mergeStmt := m.ExpectPrepare("UPDATE cfg_entities").WillBeClosed()
	for _, delta := range graphIDChanges {
		mergeStmt.ExpectExec().WithArgs(delta[1], delta[0]).WillReturnResult(mockResult)
	}
}

func expectEdgeQueries(m sqlmock.Sqlmock, assocs []*storage.EntityID, edgeLoadsByTk map[storage2.TypeAndKey]expectedEntQueryResult) {
	expectedLoads := funk.Map(
		assocs,
		func(id *storage.EntityID) expectedEntQueryResult { return edgeLoadsByTk[id.ToTypeAndKey()] },
	).([]expectedEntQueryResult)
	expectBasicEntityQueries(m, expectedLoads...)
}

// [(from_pk, to_pk)]
func expectEdgeInsertions(m sqlmock.Sqlmock, newEdges [][2]string) {
	args := make([]driver.Value, 0, len(newEdges)*2)
	for _, edge := range newEdges {
		args = append(args, edge[0], edge[1])
	}
	m.ExpectExec("INSERT INTO cfg_assocs").WithArgs(args...).WillReturnResult(mockResult)
}

func expectEdgeDeletions(m sqlmock.Sqlmock, edges [][2]string) {
	args := make([]driver.Value, 0, len(edges)*2)
	for _, edge := range edges {
		args = append(args, edge[0], edge[1])
	}
	m.ExpectExec("DELETE FROM cfg_assocs").WithArgs(args...).WillReturnResult(mockResult)
}

func expectPermissionCreation(m sqlmock.Sqlmock, entPk string, startId int, perms ...*storage.ACL) {
	args := make([]driver.Value, 0, len(perms)*6)
	for _, perm := range perms {
		exp := getExpectedACLInsert(entPk, &startId, perm)
		args = append(args, exp.id, exp.entPk, exp.scope, exp.perm, exp.aclType, exp.filter)
		startId++
	}
	m.ExpectExec("INSERT INTO cfg_acls").WillReturnResult(mockResult)
}

func expectPermissionUpdates(m sqlmock.Sqlmock, entPk string, perms ...*storage.ACL) {
	stmt := m.ExpectPrepare("UPDATE cfg_acls").WillBeClosed()
	for _, perm := range perms {
		exp := getExpectedACLInsert(entPk, nil, perm)
		stmt.ExpectExec().WithArgs(exp.scope, exp.perm, exp.aclType, exp.filter, exp.id).WillReturnResult(mockResult)
	}
}

func expectPermissionDeletes(m sqlmock.Sqlmock, permIDs ...string) {
	args := make([]driver.Value, 0, len(permIDs))
	funk.ConvertSlice(permIDs, &args)
	m.ExpectExec("DELETE FROM cfg_acls").WithArgs(args...).WillReturnResult(mockResult)
}

type expectedACLInsert struct {
	id, entPk, scope, perm, aclType, filter driver.Value
}

func getExpectedACLInsert(entPk string, idOverride *int, perm *storage.ACL) expectedACLInsert {
	var scope, typeVal, filter driver.Value

	switch perm.Scope.(type) {
	case *storage.ACL_ScopeWildcard:
		scope = perm.GetScopeWildcard().String()
	default:
		scope = strings.Join(perm.GetScopeNetworkIDs().IDs, ",")
	}

	switch perm.Type.(type) {
	case *storage.ACL_TypeWildcard:
		typeVal = perm.GetTypeWildcard().String()
	default:
		typeVal = perm.GetEntityType()
	}

	if funk.IsEmpty(perm.IDFilter) {
		filter = nil
	} else {
		filter = strings.Join(perm.IDFilter, ",")
	}

	ret := expectedACLInsert{entPk: entPk, perm: perm.Permission, aclType: typeVal, scope: scope, filter: filter}
	if idOverride == nil {
		ret.id = perm.ID
	} else {
		ret.id = fmt.Sprintf("%d", *idOverride)
	}
	return ret
}

func getBasicQueryExpect(entType string, entKey string) expectedEntQueryResult {
	return expectedEntQueryResult{entType, entKey, entType + entKey, "", "g1", 0}
}

func assocsToEdges(entPk string, assocs []*storage.EntityID, edgeLoadsByTk map[storage2.TypeAndKey]expectedEntQueryResult) [][2]string {
	return funk.Map(
		assocs,
		func(id *storage.EntityID) [2]string {
			return [2]string{entPk, edgeLoadsByTk[id.ToTypeAndKey()].pk}
		},
	).([][2]string)
}

func stringPointer(val string) *wrappers.StringValue {
	if val == "" {
		return nil
	}
	return &wrappers.StringValue{Value: val}
}

func stringValue(val *wrappers.StringValue) string {
	if val == nil {
		return ""
	}
	return val.Value
}

func bytesPointer(val []byte) *wrappers.BytesValue {
	if funk.IsEmpty(val) {
		return nil
	}
	return &wrappers.BytesValue{Value: val}
}

func bytesVal(val *wrappers.BytesValue) []byte {
	if val == nil {
		return nil
	}
	return val.Value
}
