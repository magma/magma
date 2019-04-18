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
	"errors"
	"log"
	"testing"

	"magma/orc8r/cloud/go/services/configurator/storage"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var mockResult = sqlmock.NewResult(1, 1)

func TestSqlConfiguratorStorage_LoadNetworks(t *testing.T) {
	runFactory := func(ids []string, loadCriteria storage.NetworkLoadCriteria) func(store storage.ConfiguratorStorage) (interface{}, error) {
		return func(store storage.ConfiguratorStorage) (interface{}, error) {
			return store.LoadNetworks(ids, loadCriteria)
		}
	}

	idsOnly := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT cfg_networks.id, cfg_networks.version FROM cfg_networks").
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "version"}).
						AddRow("hello", 1).
						AddRow("world", 2),
				)
		},
		run: runFactory([]string{"hello", "world"}, storage.NetworkLoadCriteria{}),

		expectedResult: storage.NetworkLoadResult{
			NetworkIDsNotFound: []string{},
			Networks: []storage.Network{
				{ID: "hello", Configs: map[string][]byte{}, Version: 1},
				{ID: "world", Configs: map[string][]byte{}, Version: 2},
			},
		},
	}

	metadataLoad := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT cfg_networks.id, cfg_networks.name, cfg_networks.description, cfg_networks.version FROM cfg_networks").
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "name", "description", "version"}).
						AddRow("hello", "Hello", "Hello network", 1).
						AddRow("world", "World", "World network", 2),
				)
		},
		run: runFactory([]string{"hello", "world"}, storage.NetworkLoadCriteria{LoadMetadata: true}),

		expectedResult: storage.NetworkLoadResult{
			NetworkIDsNotFound: []string{},
			Networks: []storage.Network{
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
			m.ExpectQuery("SELECT cfg_networks.id, cfg_network_configs.type, cfg_network_configs.value, cfg_networks.version FROM cfg_networks").
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "type", "value", "version"}).
						AddRow("hello", "foo", []byte("foo"), 1).
						AddRow("hello", "bar", []byte("bar"), 1).
						AddRow("world", nil, nil, 3),
				)
		},
		run: runFactory([]string{"hello", "world", "foobar"}, storage.NetworkLoadCriteria{LoadConfigs: true}),

		expectedResult: storage.NetworkLoadResult{
			NetworkIDsNotFound: []string{"foobar"},
			Networks: []storage.Network{
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
			m.ExpectQuery("SELECT cfg_networks.id, cfg_networks.name, cfg_networks.description, cfg_network_configs.type, cfg_network_configs.value, cfg_networks.version FROM cfg_networks").
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "name", "description", "type", "value", "version"}).
						AddRow("hello", "Hello", "Hello network", "foo", []byte("foo"), 1).
						AddRow("hello", "Hello", "Hello network", "bar", []byte("bar"), 1).
						AddRow("world", "World", "World network", nil, nil, 2),
				)
		},
		run: runFactory([]string{"hello", "world", "foobar"}, storage.NetworkLoadCriteria{LoadMetadata: true, LoadConfigs: true}),

		expectedResult: storage.NetworkLoadResult{
			NetworkIDsNotFound: []string{"foobar"},
			Networks: []storage.Network{
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
			m.ExpectQuery("SELECT cfg_networks.id, cfg_networks.version FROM cfg_networks").
				WillReturnRows(sqlmock.NewRows([]string{"id", "version"}))
		},
		run: runFactory([]string{"hello", "world"}, storage.NetworkLoadCriteria{}),

		expectedResult: storage.NetworkLoadResult{
			NetworkIDsNotFound: []string{"hello", "world"},
			Networks:           []storage.Network{},
		},
	}

	queryError := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT cfg_networks.id, cfg_networks.version FROM cfg_networks").
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
			m.ExpectQuery(`SELECT count\(1\) FROM cfg_networks`).
				WithArgs("n1").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			m.ExpectExec("INSERT INTO cfg_networks").
				WithArgs("n1", "", "").
				WillReturnResult(mockResult)
		},
		run: runFactory(storage.Network{ID: "n1"}),

		expectedResult: storage.Network{ID: "n1"},
	}

	allMetadata := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT count\(1\) FROM cfg_networks`).
				WithArgs("n2").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			m.ExpectExec("INSERT INTO cfg_networks").
				WithArgs("n2", "hello", "world").
				WillReturnResult(mockResult)
		},
		run: runFactory(storage.Network{ID: "n2", Name: "hello", Description: "world"}),

		expectedResult: storage.Network{ID: "n2", Name: "hello", Description: "world"},
	}

	everythingNw := storage.Network{
		ID:          "n3",
		Name:        "hello",
		Description: "world",
		Configs: map[string][]byte{
			"foo": []byte("bar"),
			"baz": []byte("quz"),
		},
	}
	everything := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT count\(1\) FROM cfg_networks`).
				WithArgs("n3").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			m.ExpectExec("INSERT INTO cfg_networks").
				WithArgs("n3", "hello", "world").
				WillReturnResult(mockResult)

			configStmt := m.ExpectPrepare("INSERT INTO cfg_network_configs")
			configStmt.ExpectExec().WithArgs("n3", "baz", []byte("quz")).WillReturnResult(mockResult)
			configStmt.ExpectExec().WithArgs("n3", "foo", []byte("bar")).WillReturnResult(mockResult)
			configStmt.WillBeClosed()
		},
		run: runFactory(everythingNw),

		expectedResult: everythingNw,
	}

	networkTableError := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT count\(1\) FROM cfg_networks`).
				WithArgs("n4").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			m.ExpectExec("INSERT INTO cfg_networks").
				WithArgs("n4", "", "").
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
			m.ExpectQuery(`SELECT count\(1\) FROM cfg_networks`).
				WithArgs("n5").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			m.ExpectExec("INSERT INTO cfg_networks").
				WithArgs("n5", "", "").
				WillReturnResult(mockResult)

			configStmt := m.ExpectPrepare("INSERT INTO cfg_network_configs")
			configStmt.ExpectExec().WithArgs("n5", "baz", []byte("quz")).
				WillReturnError(errors.New("mock exec error"))
			configStmt.WillBeClosed()
		},
		run: runFactory(configTableErrNw),

		expectedResult: configTableErrNw,
		expectedError:  errors.New("error inserting config baz: mock exec error"),
	}

	networkExists := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT count\(1\) FROM cfg_networks`).
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
			return store.UpdateNetworks(updates)
		}
	}

	// Delete 1 network (n1)
	// Update only metadata of another (n2)
	// Update and delete configs of another (n3)
	// Fill out all fields of the update criteria (n4)
	names := []string{"should be ignored", "name2", "name4"}
	descs := []string{"should be ignored", "desc2", ""}
	happyPath := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			deleteStmt := m.ExpectPrepare("DELETE FROM cfg_networks")
			upsertStmt := m.ExpectPrepare("INSERT INTO cfg_network_configs")
			deleteConfigStmt := m.ExpectPrepare("DELETE FROM cfg_network_configs")

			deleteStmt.ExpectExec().WithArgs("n1").WillReturnResult(mockResult)

			m.ExpectExec("UPDATE cfg_networks").WithArgs(names[1], descs[1], "n2").WillReturnResult(mockResult)

			m.ExpectExec("UPDATE cfg_networks").WithArgs("n3").WillReturnResult(mockResult)
			upsertStmt.ExpectExec().WithArgs("n3", "baz", []byte("quz"), []byte("quz")).WillReturnResult(mockResult)
			upsertStmt.ExpectExec().WithArgs("n3", "foo", []byte("bar"), []byte("bar")).WillReturnResult(mockResult)
			deleteConfigStmt.ExpectExec().WithArgs("n3", "hello").WillReturnResult(mockResult)
			deleteConfigStmt.ExpectExec().WithArgs("n3", "world").WillReturnResult(mockResult)

			m.ExpectExec("UPDATE cfg_networks").WithArgs(names[2], nil, "n4").WillReturnResult(mockResult)
			upsertStmt.ExpectExec().WithArgs("n4", "baz", []byte("quz"), []byte("quz")).WillReturnResult(mockResult)
			upsertStmt.ExpectExec().WithArgs("n4", "foo", []byte("bar"), []byte("bar")).WillReturnResult(mockResult)
			deleteConfigStmt.ExpectExec().WithArgs("n4", "hello").WillReturnResult(mockResult)
			deleteConfigStmt.ExpectExec().WithArgs("n4", "world").WillReturnResult(mockResult)

			deleteStmt.WillBeClosed()
			upsertStmt.WillBeClosed()
			deleteConfigStmt.WillBeClosed()
		},
		run: runFactory(
			[]storage.NetworkUpdateCriteria{
				{ID: "n1", DeleteNetwork: true, NewName: &names[0], NewDescription: &descs[0]},
				{ID: "n2", NewName: &names[1], NewDescription: &descs[1]},
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
					NewName:         &names[2],
					NewDescription:  &descs[2],
					ConfigsToDelete: []string{"hello", "world"},
					ConfigsToAddOrUpdate: map[string][]byte{
						"foo": []byte("bar"),
						"baz": []byte("quz"),
					},
				},
			},
		),

		expectedResult: storage.FailedOperations{},
	}

	// Error in 1 network should not block other networks (try with 3 networks/2 errors)
	errorCase := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			deleteStmt := m.ExpectPrepare("DELETE FROM cfg_networks")
			upsertStmt := m.ExpectPrepare("INSERT INTO cfg_network_configs")
			deleteConfigStmt := m.ExpectPrepare("DELETE FROM cfg_network_configs")

			deleteStmt.ExpectExec().WithArgs("n1").WillReturnError(errors.New("n1 delete error"))

			m.ExpectExec("UPDATE cfg_networks").WithArgs(names[1], descs[1], "n2").WillReturnError(errors.New("n2 update error"))

			m.ExpectExec("UPDATE cfg_networks").WithArgs("n3").WillReturnResult(mockResult)
			upsertStmt.ExpectExec().WithArgs("n3", "baz", []byte("quz"), []byte("quz")).WillReturnResult(mockResult)
			upsertStmt.ExpectExec().WithArgs("n3", "foo", []byte("bar"), []byte("bar")).WillReturnError(errors.New("n3foo update error"))

			m.ExpectExec("UPDATE cfg_networks").WithArgs(names[2], nil, "n4").WillReturnResult(mockResult)
			upsertStmt.ExpectExec().WithArgs("n4", "baz", []byte("quz"), []byte("quz")).WillReturnResult(mockResult)
			upsertStmt.ExpectExec().WithArgs("n4", "foo", []byte("bar"), []byte("bar")).WillReturnResult(mockResult)
			deleteConfigStmt.ExpectExec().WithArgs("n4", "hello").WillReturnResult(mockResult)
			deleteConfigStmt.ExpectExec().WithArgs("n4", "world").WillReturnResult(mockResult)

			deleteStmt.WillBeClosed()
			upsertStmt.WillBeClosed()
			deleteConfigStmt.WillBeClosed()
		},
		run: runFactory(
			[]storage.NetworkUpdateCriteria{
				{ID: "n1", DeleteNetwork: true},
				{ID: "n2", NewName: &names[1], NewDescription: &descs[1]},
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
					NewName:         &names[2],
					NewDescription:  &descs[2],
					ConfigsToDelete: []string{"hello", "world"},
					ConfigsToAddOrUpdate: map[string][]byte{
						"foo": []byte("bar"),
						"baz": []byte("quz"),
					},
				},
			},
		),

		expectedResult: storage.FailedOperations{
			"n1": errors.New("error deleting network n1: n1 delete error"),
			"n2": errors.New("error updating network n2: n2 update error"),
			"n3": errors.New("error updating config foo on network n3: n3foo update error"),
		},
		expectedError: errors.New("some errors were encountered, see return value for details"),
	}

	validationFailure := &testCase{
		setup: func(m sqlmock.Sqlmock) {},

		run: runFactory(
			[]storage.NetworkUpdateCriteria{
				{ID: "n1", DeleteNetwork: true},
				{ID: "n1", NewName: &names[1]},
			},
		),

		expectedError: errors.New("multiple updates for a single network are not allowed"),
	}

	runCase(t, happyPath)
	runCase(t, errorCase)
	runCase(t, validationFailure)
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

	factory := storage.NewSQLConfiguratorStorageFactory(db)
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
