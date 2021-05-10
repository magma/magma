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
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"regexp"
	"testing"

	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"
	orc8r_storage "magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	sqlTestMaxLoadSize = 5
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

		expectedError: errors.New("error inserting network: mock exec error"),
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

		expectedError: errors.New("error inserting network configs: mock exec error"),
	}

	networkExists := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(`SELECT COUNT\(1\) FROM cfg_networks`).
				WithArgs("n5").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		},
		run: runFactory(storage.Network{ID: "n5"}),

		expectedError: errors.New("a network with ID n5 already exists"),
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
				{NetworkID: "network", Type: "baz", Key: "quz", GraphID: "42", Version: 1, Pk: "def"},
				{NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 2, Pk: "abc"},
			},
			EntitiesNotFound: []*storage.EntityID{{Type: "hello", Key: "world"}},
		},
	}

	// Load everything, no assocs
	// Side benefit: ensure redundant rows don't break anything
	loadEverything := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT .* FROM cfg_entities").
				WithArgs(
					"network", "bar", "foo",
					"network", "quz", "baz",
					"network", "world", "hello",
				).
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id", "name", "description", "config"}).
						AddRow("network", "abc", "bar", "foo", nil, 2, "42", "foobar", "foobar ent", []byte("foobar")).
						AddRow("network", "abc", "bar", "foo", nil, 2, "42", "foobar", "foobar ent", []byte("foobar")).
						AddRow("network", "def", "quz", "baz", nil, 1, "42", "bazquz", "bazquz ent", []byte("bazquz")),
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
			storage.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true},
		),

		expectedResult: storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "baz", Key: "quz", GraphID: "42", Version: 1, Pk: "def",
					Name:        "bazquz",
					Description: "bazquz ent",
					Config:      []byte("bazquz"),
				},
				{
					NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 2, Pk: "abc",
					Name:        "foobar",
					Description: "foobar ent",
					Config:      []byte("foobar"),
				},
			},
			EntitiesNotFound: []*storage.EntityID{{Type: "hello", Key: "world"}},
		},
	}

	loadPageBaseQuery := `SELECT ent.network_id, ent.pk, ent."key", ent.type, ent.physical_id, ent.version, ent.graph_id, ent.name, ent.description, ent.config ` +
		`FROM cfg_entities AS ent `
	emptyPageTokenWhere := `WHERE (ent.network_id = $1 AND ent.type = $2) `
	nonEmptyPageTokenWhere := `WHERE (ent.network_id = $1 AND ent.type = $2 AND ent."key" > $3) `
	orderbyLimit := `ORDER BY ent."key" LIMIT %d`
	loadPageEmptyToken := loadPageBaseQuery + emptyPageTokenWhere + orderbyLimit
	loadPageNonEmptyToken := loadPageBaseQuery + nonEmptyPageTokenWhere + orderbyLimit

	nextToken := &storage.EntityPageToken{
		LastIncludedEntity: "rou",
	}
	expectedNextToken := serializeToken(t, nextToken)
	loadFullPage := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(loadPageEmptyToken, 2))).
				WithArgs(
					"network", "foo",
				).
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id", "name", "description", "config"}).
						AddRow("network", "abc", "bar", "foo", nil, 2, "42", "foobar", "foobar ent", []byte("foobar")).
						AddRow("network", "bbb", "rou", "foo", nil, 2, "43", "foobar", "barbar ent", []byte("foobar")),
				)
		},
		run: runFactory(
			"network",
			storage.EntityLoadFilter{TypeFilter: &wrappers.StringValue{Value: "foo"}},
			storage.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, PageSize: 2, PageToken: ""},
		),

		expectedResult: storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 2, Pk: "abc",
					Name:        "foobar",
					Description: "foobar ent",
					Config:      []byte("foobar"),
				},
				{
					NetworkID: "network", Type: "foo", Key: "rou", GraphID: "43", Version: 2, Pk: "bbb",
					Name:        "foobar",
					Description: "barbar ent",
					Config:      []byte("foobar"),
				},
			},
			NextPageToken: expectedNextToken,
		},
	}

	loadFinalPage := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(loadPageNonEmptyToken, 2))).
				WithArgs(
					"network", "foo", "rou",
				).
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id", "name", "description", "config"}).
						AddRow("network", "abc", "zed", "foo", nil, 2, "42", "zedbar", "foobar ent", []byte("foobar")),
				)
		},
		run: runFactory(
			"network",
			storage.EntityLoadFilter{TypeFilter: &wrappers.StringValue{Value: "foo"}},
			storage.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, PageSize: 2, PageToken: expectedNextToken},
		),

		expectedResult: storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "foo", Key: "zed", GraphID: "42", Version: 2, Pk: "abc",
					Name:        "zedbar",
					Description: "foobar ent",
					Config:      []byte("foobar"),
				},
			},
			NextPageToken: "",
		},
	}
	nextToken = &storage.EntityPageToken{
		LastIncludedEntity: "eee",
	}
	expectedNextToken = serializeToken(t, nextToken)
	loadPageMaxRegex := regexp.QuoteMeta(fmt.Sprintf(loadPageEmptyToken, sqlTestMaxLoadSize))
	loadPageSizeGreaterThanMax := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(loadPageMaxRegex).
				WithArgs(
					"network", "foo",
				).
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id", "name", "description", "config"}).
						AddRow("network", "aaa", "aaa", "foo", nil, 2, "42", "aaafoo", "aaafoo ent", []byte("aaafoo")).
						AddRow("network", "bbb", "bbb", "foo", nil, 2, "43", "bbbfoo", "bbbfoo ent", []byte("bbbfoo")).
						AddRow("network", "ccc", "ccc", "foo", nil, 2, "44", "cccfoo", "cccfoo ent", []byte("cccfoo")).
						AddRow("network", "ddd", "ddd", "foo", nil, 2, "45", "dddfoo", "dddfoo ent", []byte("dddfoo")).
						AddRow("network", "eee", "eee", "foo", nil, 2, "46", "eeefoo", "eeefoo ent", []byte("eeefoo")),
				)
		},
		run: runFactory(
			"network",
			storage.EntityLoadFilter{TypeFilter: &wrappers.StringValue{Value: "foo"}},
			storage.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, PageSize: 10, PageToken: ""},
		),

		expectedResult: storage.EntityLoadResult{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "foo", Key: "aaa", GraphID: "42", Version: 2, Pk: "aaa",
					Name:        "aaafoo",
					Description: "aaafoo ent",
					Config:      []byte("aaafoo"),
				},
				{
					NetworkID: "network", Type: "foo", Key: "bbb", GraphID: "43", Version: 2, Pk: "bbb",
					Name:        "bbbfoo",
					Description: "bbbfoo ent",
					Config:      []byte("bbbfoo"),
				},
				{
					NetworkID: "network", Type: "foo", Key: "ccc", GraphID: "44", Version: 2, Pk: "ccc",
					Name:        "cccfoo",
					Description: "cccfoo ent",
					Config:      []byte("cccfoo"),
				},
				{
					NetworkID: "network", Type: "foo", Key: "ddd", GraphID: "45", Version: 2, Pk: "ddd",
					Name:        "dddfoo",
					Description: "dddfoo ent",
					Config:      []byte("dddfoo"),
				},
				{
					NetworkID: "network", Type: "foo", Key: "eee", GraphID: "46", Version: 2, Pk: "eee",
					Name:        "eeefoo",
					Description: "eeefoo ent",
					Config:      []byte("eeefoo"),
				},
			},
			NextPageToken: expectedNextToken,
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

			ents := entsByPK{
				"abc": {"abc", "bar", "foo"},
				"def": {"def", "quz", "baz"},
				"ghi": {"ghi", "world", "hello"},
			}
			expectAssocQuery(
				m,
				"network",
				ents.getAssocArgs("abc", "def", "ghi"),
				[][]driver.Value{
					ents.getParentAssocReturns(t, "abc", "def"),
					ents.getParentAssocReturns(t, "abc", "ghi"),
					ents.getParentAssocReturns(t, "ghi", "def"),
				},
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
					NetworkID: "network", Type: "baz", Key: "quz", GraphID: "42", Version: 1, Pk: "def",
					ParentAssociations: []*storage.EntityID{
						{Type: "foo", Key: "bar"},
						{Type: "hello", Key: "world"},
					},
				},
				{
					NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 2, Pk: "abc",
				},
				{
					NetworkID: "network", Type: "hello", Key: "world", GraphID: "42", Version: 3, Pk: "ghi",
					ParentAssociations: []*storage.EntityID{
						{Type: "foo", Key: "bar"},
					},
				},
			},
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

			ents := entsByPK{
				"abc": {"abc", "bar", "foo"},
				"def": {"def", "quz", "baz"},
				"ghi": {"ghi", "world", "hello"},
			}
			expectAssocQuery(
				m,
				"network",
				ents.getAssocArgs("abc", "def", "ghi"),
				[][]driver.Value{
					ents.getAssocReturns(t, "def", "abc"),
					ents.getAssocReturns(t, "ghi", "abc"),
					ents.getAssocReturns(t, "def", "ghi"),
				},
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
					NetworkID: "network", Type: "baz", Key: "quz", GraphID: "42", Version: 1, Pk: "def",
					Associations: []*storage.EntityID{
						{Type: "foo", Key: "bar"},
						{Type: "hello", Key: "world"},
					},
				},
				{
					NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 2, Pk: "abc",
				},
				{
					NetworkID: "network", Type: "hello", Key: "world", GraphID: "42", Version: 3, Pk: "ghi",
					Associations: []*storage.EntityID{
						{Type: "foo", Key: "bar"},
					},
				},
			},
		},
	}

	// Load everything with type filter (type foo)
	fullLoadTypeFilter := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery("SELECT .* FROM cfg_entities").
				WithArgs("network", "foo").
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "pk", "key", "type", "physical_id", "version", "graph_id", "name", "description", "config"}).
						// return foobar twice, to test DAO resiliency
						AddRow("network", "foobar", "bar", "foo", nil, 1, "42", "foobar", "foobar ent", []byte("foobar")).
						AddRow("network", "foobar", "bar", "foo", nil, 1, "42", "foobar", "foobar ent", []byte("foobar")).
						AddRow("network", "foobaz", "baz", "foo", nil, 2, "42", "foobaz", "foobaz ent", []byte("foobaz")),
				)

			ents := entsByPK{
				"foobar":     {"foobar", "bar", "foo"},
				"foobaz":     {"foobaz", "baz", "foo"},
				"barbaz":     {"barbaz", "baz", "bar"},
				"bazquz":     {"bazquz", "quz", "baz"},
				"helloworld": {"helloworld", "world", "hello"},
			}
			// Child assocs
			expectAssocQuery(
				m,
				"network",
				[][]driver.Value{{"foo"}},
				[][]driver.Value{
					ents.getAssocReturns(t, "foobar", "foobaz"),
					ents.getAssocReturns(t, "foobar", "barbaz"),
					ents.getAssocReturns(t, "foobaz", "bazquz"),
					ents.getAssocReturns(t, "helloworld", "foobar"),
				},
			)
			// Parent assocs
			expectAssocQuery(
				m,
				"network",
				[][]driver.Value{{"foo"}},
				[][]driver.Value{
					ents.getParentAssocReturns(t, "foobar", "foobaz"),
					ents.getParentAssocReturns(t, "foobar", "barbaz"),
					ents.getParentAssocReturns(t, "foobaz", "bazquz"),
					ents.getParentAssocReturns(t, "helloworld", "foobar"),
				},
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
					NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 1, Pk: "foobar",
					Name:        "foobar",
					Description: "foobar ent",
					Config:      []byte("foobar"),
					Associations: []*storage.EntityID{
						{Type: "bar", Key: "baz"},
						{Type: "foo", Key: "baz"},
					},
					ParentAssociations: []*storage.EntityID{
						{Type: "hello", Key: "world"},
					},
				},
				{
					NetworkID: "network", Type: "foo", Key: "baz", GraphID: "42", Version: 2, Pk: "foobaz",
					Name:        "foobaz",
					Description: "foobaz ent",
					Config:      []byte("foobaz"),
					Associations: []*storage.EntityID{
						{Type: "baz", Key: "quz"},
					},
					ParentAssociations: []*storage.EntityID{
						{Type: "foo", Key: "bar"},
					},
				},
			},
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
				{NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", Version: 2, Pk: "abc"},
			},
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
				{NetworkID: "network", Type: "foo", Key: "bar", GraphID: "42", PhysicalID: "p1", Version: 2, Pk: "abc"},
			},
		},
	}

	runCase(t, basicOnly)
	runCase(t, loadEverything)
	runCase(t, loadFullPage)
	runCase(t, loadFinalPage)
	runCase(t, loadPageSizeGreaterThanMax)
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
			Pk:          "1",
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
			assocsWithDuplicate := []*storage.EntityID{{Type: "bar", Key: "baz"}, {Type: "baz", Key: "quz"}, {Type: "bar", Key: "baz"}}
			edgesByTk := map[orc8r_storage.TypeAndKey]expectedEntQueryResult{
				{Type: "bar", Key: "baz"}: {"bar", "baz", "42", "", "1", 1},
				{Type: "baz", Key: "quz"}: {"baz", "quz", "43", "", "3", 2},
			}
			expectEdgeQueries(m, assocsWithDuplicate, edgesByTk)
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
					// Duplicate edge, should be coalesced properly
					{Type: "bar", Key: "baz"},
				},
			},
		),

		// Merged graphs
		expectedResult: storage.NetworkEntity{
			NetworkID:   "network",
			Type:        "foo",
			Key:         "bar",
			Name:        "foobar",
			Description: "foobar ent",
			GraphID:     "1",
			Pk:          "1",
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
		expectedError:  errors.New("an entity 'foo-bar' already exists"),
	}

	runCase(t, basicCase)
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
			ents := entsByPK{
				"barbaz": {"barbaz", "baz", "bar"},
				"bazquz": {"bazquz", "quz", "baz"},
				"quzbaz": {"quzbaz", "baz", "quz"},
				"bazbar": {"bazbar", "bar", "baz"},
				"barfoo": {"barfoo", "foo", "bar"},
			}
			expectAssocQuery(
				m,
				"network",
				[][]driver.Value{{"g1"}},
				[][]driver.Value{
					ents.getAssocReturns(t, "barbaz", "bazquz"),
					ents.getAssocReturns(t, "barbaz", "quzbaz"),
				},
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
	// delete edges (bazquz -> quzbaz), (bazquz -> bazbar)
	// before:	foobar -> barbaz -> bazquz -> (quzbaz, bazbar -> barfoo)
	// after: 	foobar -> barbaz -> bazquz | quzbaz | bazbar -> barfoo
	// partitions graph into 3 components
	partitionCase := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			expectBasicEntityQueries(m, getBasicQueryExpect("baz", "quz"))
			m.ExpectExec("UPDATE cfg_entities").WithArgs("bazquz").WillReturnResult(mockResult)
			expectEdgeQueries(
				m,
				[]*storage.EntityID{{Type: "quz", Key: "baz"}, {Type: "baz", Key: "bar"}},
				map[orc8r_storage.TypeAndKey]expectedEntQueryResult{
					{Type: "quz", Key: "baz"}: getBasicQueryExpect("quz", "baz"),
					{Type: "baz", Key: "bar"}: getBasicQueryExpect("baz", "bar"),
				},
			)
			expectEdgeDeletions(m, [][2]string{{"bazquz", "quzbaz"}, {"bazquz", "bazbar"}})

			expectBulkEntityQuery(
				m,
				[]driver.Value{"g1"},
				getBasicQueryExpect("bar", "baz"),
				getBasicQueryExpect("bar", "foo"),
				getBasicQueryExpect("baz", "bar"),
				getBasicQueryExpect("baz", "quz"),
				getBasicQueryExpect("foo", "bar"),
				getBasicQueryExpect("quz", "baz"),
			)
			ents := entsByPK{
				"barbaz": {"barbaz", "baz", "bar"},
				"barfoo": {"barfoo", "foo", "bar"},
				"bazbar": {"bazbar", "bar", "baz"},
				"bazquz": {"bazquz", "quz", "baz"},
				"foobar": {"foobar", "bar", "foo"},
				"quzbaz": {"quzbaz", "baz", "quz"},
			}
			expectAssocQuery(
				m,
				"network",
				[][]driver.Value{{"g1"}},
				[][]driver.Value{
					ents.getAssocReturns(t, "barbaz", "bazquz"),
					ents.getAssocReturns(t, "bazbar", "barfoo"),
					ents.getAssocReturns(t, "foobar", "barbaz"),
				},
			)
			m.ExpectExec("UPDATE cfg_entities").WithArgs("1", "quzbaz").WillReturnResult(mockResult)
			m.ExpectExec("UPDATE cfg_entities").WithArgs("2", "barfoo", "bazbar").WillReturnResult(mockResult)
		},
		run:            runFactory("network", storage.EntityUpdateCriteria{Type: "baz", Key: "quz", AssociationsToDelete: []*storage.EntityID{{Type: "quz", Key: "baz"}, {Type: "baz", Key: "bar"}}}),
		expectedResult: storage.NetworkEntity{NetworkID: "network", Type: "baz", Key: "quz", GraphID: "g1", Pk: "bazquz", Version: 1},
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
			expectAssocQuery(
				m,
				"network",
				[][]driver.Value{{"g1"}},
				[][]driver.Value{},
			)

			// Graph partition update
			m.ExpectExec("UPDATE cfg_entities").WithArgs("1", "barbaz").WillReturnResult(mockResult)
		},
		run:            runFactory("network", storage.EntityUpdateCriteria{Type: "foo", Key: "bar", AssociationsToSet: &storage.EntityAssociationsToSet{}}),
		expectedResult: storage.NetworkEntity{NetworkID: "network", Type: "foo", Key: "bar", GraphID: "g1", Pk: "foobar", Version: 1},
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

	// load a linked list of 3 nodes
	linkedList := &testCase{
		run: runFactory("network", storage.EntityID{Type: "foo", Key: "bar"}, storage.EntityLoadCriteria{}),

		setup: func(m sqlmock.Sqlmock) {
			expectBasicEntityQueries(m, expectedFooBar)
			expectBulkEntityQuery(m, []driver.Value{"g1"}, expectedFooBar, expectedBarBaz, expectedBazQuz)
			// foobar -> barbaz -> bazquz
			ents := entsByPK{
				"foobar": {"foobar", "bar", "foo"},
				"barbaz": {"barbaz", "baz", "bar"},
				"bazquz": {"bazquz", "quz", "baz"},
			}
			expectAssocQuery(
				m,
				"network",
				[][]driver.Value{{"g1"}},
				[][]driver.Value{
					ents.getAssocReturns(t, "foobar", "barbaz"),
					ents.getAssocReturns(t, "barbaz", "bazquz"),
				},
			)
		},

		expectedResult: storage.EntityGraph{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "bar", Key: "baz",
					PhysicalID: "p1", GraphID: "g1", Pk: "barbaz",
					Associations:       []*storage.EntityID{{Type: "baz", Key: "quz"}},
					ParentAssociations: []*storage.EntityID{{Type: "foo", Key: "bar"}},
					Version:            1,
				},
				{
					NetworkID: "network", Type: "baz", Key: "quz",
					PhysicalID: "p2", GraphID: "g1", Pk: "bazquz",
					ParentAssociations: []*storage.EntityID{{Type: "bar", Key: "baz"}},
					Version:            2,
				},
				{
					NetworkID: "network", Type: "foo", Key: "bar",
					GraphID: "g1", Pk: "foobar",
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
			ents := entsByPK{
				"foobar": {"foobar", "bar", "foo"},
				"barbaz": {"barbaz", "baz", "bar"},
				"bazquz": {"bazquz", "quz", "baz"},
			}
			expectAssocQuery(
				m,
				"network",
				[][]driver.Value{{"g1"}},
				[][]driver.Value{
					ents.getAssocReturns(t, "foobar", "barbaz"),
					ents.getAssocReturns(t, "foobar", "bazquz"),
				},
			)
		},

		expectedResult: storage.EntityGraph{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "bar", Key: "baz",
					PhysicalID: "p1", GraphID: "g1", Pk: "barbaz",
					ParentAssociations: []*storage.EntityID{{Type: "foo", Key: "bar"}},
					Version:            1,
				},
				{
					NetworkID: "network", Type: "baz", Key: "quz",
					PhysicalID: "p2", GraphID: "g1", Pk: "bazquz",
					ParentAssociations: []*storage.EntityID{{Type: "foo", Key: "bar"}},
					Version:            2,
				},
				{
					NetworkID: "network", Type: "foo", Key: "bar",
					GraphID: "g1", Pk: "foobar",
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
			ents := entsByPK{
				"foobar": {"foobar", "bar", "foo"},
				"barbaz": {"barbaz", "baz", "bar"},
				"bazquz": {"bazquz", "quz", "baz"},
			}
			expectAssocQuery(
				m,
				"network",
				[][]driver.Value{{"g1"}},
				[][]driver.Value{
					ents.getAssocReturns(t, "barbaz", "foobar"),
					ents.getAssocReturns(t, "bazquz", "foobar"),
				},
			)
		},

		expectedResult: storage.EntityGraph{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "bar", Key: "baz",
					PhysicalID: "p1", GraphID: "g1", Pk: "barbaz",
					Associations: []*storage.EntityID{{Type: "foo", Key: "bar"}},
					Version:      1,
				},
				{
					NetworkID: "network", Type: "baz", Key: "quz",
					PhysicalID: "p2", GraphID: "g1", Pk: "bazquz",
					Associations: []*storage.EntityID{{Type: "foo", Key: "bar"}},
					Version:      2,
				},
				{
					NetworkID: "network", Type: "foo", Key: "bar",
					GraphID: "g1", Pk: "foobar",
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
			ents := entsByPK{
				"foobar": {"foobar", "bar", "foo"},
				"barbaz": {"barbaz", "baz", "bar"},
				"bazquz": {"bazquz", "quz", "baz"},
			}
			expectAssocQuery(
				m,
				"network",
				[][]driver.Value{{"g1"}},
				[][]driver.Value{
					ents.getAssocReturns(t, "foobar", "barbaz"),
					ents.getAssocReturns(t, "barbaz", "bazquz"),
					ents.getAssocReturns(t, "bazquz", "barbaz"),
				},
			)
		},

		expectedResult: storage.EntityGraph{
			Entities: []*storage.NetworkEntity{
				{
					NetworkID: "network", Type: "bar", Key: "baz",
					PhysicalID: "p1", GraphID: "g1", Pk: "barbaz",
					Associations:       []*storage.EntityID{{Type: "baz", Key: "quz"}},
					ParentAssociations: []*storage.EntityID{{Type: "baz", Key: "quz"}, {Type: "foo", Key: "bar"}},
					Version:            1,
				},
				{
					NetworkID: "network", Type: "baz", Key: "quz",
					PhysicalID: "p2", GraphID: "g1", Pk: "bazquz",
					Associations:       []*storage.EntityID{{Type: "bar", Key: "baz"}},
					ParentAssociations: []*storage.EntityID{{Type: "bar", Key: "baz"}},
					Version:            2,
				},
				{
					NetworkID: "network", Type: "foo", Key: "bar",
					GraphID: "g1", Pk: "foobar",
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
			ents := entsByPK{
				"foobar": {"foobar", "bar", "foo"},
				"barbaz": {"barbaz", "baz", "bar"},
				"bazquz": {"bazquz", "quz", "baz"},
			}
			expectAssocQuery(
				m,
				"network",
				[][]driver.Value{{"g1"}},
				[][]driver.Value{
					ents.getAssocReturns(t, "foobar", "barbaz"),
					ents.getAssocReturns(t, "barbaz", "bazquz"),
					ents.getAssocReturns(t, "bazquz", "foobar"),
				},
			)
		},

		expectedError: errors.New("graph does not have root nodes"),
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
	assert.NoError(t, err)
	defer func() {
		err = db.Close()
		if err != nil {
			log.Printf("error closing stub DB: %s", err)
		}
	}()

	mock.ExpectBegin()
	test.setup(mock)

	factory := storage.NewSQLConfiguratorStorageFactory(db, &mockIDGenerator{}, sqorc.GetSqlBuilder(), sqlTestMaxLoadSize)
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
		Pk:        entToUpdate.pk,
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

	edgeLoadsByTk := funk.Map(
		expectedEdgeLoads,
		func(e expectedEntQueryResult) (orc8r_storage.TypeAndKey, expectedEntQueryResult) {
			return orc8r_storage.TypeAndKey{Type: e.entType, Key: e.key}, e
		},
	).(map[orc8r_storage.TypeAndKey]expectedEntQueryResult)

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
				expectAssocQuery(
					m,
					"network",
					[][]driver.Value{{entToUpdate.graphID}},
					[][]driver.Value{},
				)
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
				expectAssocQuery(
					m,
					"network",
					[][]driver.Value{{entToUpdate.graphID}},
					[][]driver.Value{},
				)
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
	queryArgs = append([]driver.Value{"network"}, queryArgs...)
	m.ExpectQuery(`SELECT .* FROM cfg_entities AS ent WHERE \(ent.network_id = \$1 AND ent.graph_id = \$2\)`).
		WithArgs(queryArgs...).
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

func expectAssocQuery(m sqlmock.Sqlmock, networkID string, queryArgs [][]driver.Value, queryReturn [][]driver.Value) {
	rows := sqlmock.NewRows([]string{`ent."key"`, "ent.type", `assoc."key"`, "assoc.type", "ent.pk", "assoc.pk"})
	for _, row := range queryReturn {
		rows.AddRow(row...)
	}

	var withArgs []driver.Value
	for _, args := range queryArgs {
		withArgs = append(withArgs, networkID)
		withArgs = append(withArgs, args...)
	}

	m.ExpectQuery(`SELECT ent."key", ent.type, assoc."key", assoc.type, ent.pk, assoc.pk FROM cfg_entities AS ent JOIN cfg_assocs`).
		WithArgs(withArgs...).
		WillReturnRows(rows)
}

// [(old graph ID, new graph ID)]
func expectMergeGraphs(m sqlmock.Sqlmock, graphIDChanges [][2]string) {
	mergeStmt := m.ExpectPrepare("UPDATE cfg_entities").WillBeClosed()
	for _, delta := range graphIDChanges {
		mergeStmt.ExpectExec().WithArgs(delta[1], delta[0]).WillReturnResult(mockResult)
	}
}

func expectEdgeQueries(m sqlmock.Sqlmock, assocs []*storage.EntityID, edgeLoadsByTk map[orc8r_storage.TypeAndKey]expectedEntQueryResult) {
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

func getBasicQueryExpect(entType string, entKey string) expectedEntQueryResult {
	return expectedEntQueryResult{entType, entKey, entType + entKey, "", "g1", 0}
}

func assocsToEdges(entPk string, assocs []*storage.EntityID, edgeLoadsByTk map[orc8r_storage.TypeAndKey]expectedEntQueryResult) [][2]string {
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

func serializeToken(t *testing.T, token *storage.EntityPageToken) string {
	marshalledToken, err := proto.Marshal(token)
	assert.NoError(t, err)
	return base64.StdEncoding.EncodeToString(marshalledToken)
}

type ent struct{ pk, key, typ string }

type entsByPK map[string]ent

func (e entsByPK) getAssoc(t *testing.T, pkFrom, pkTo string) assoc {
	from, ok := e[pkFrom]
	assert.True(t, ok)
	to, ok := e[pkTo]
	assert.True(t, ok)
	a := assoc{
		fromKey:  from.key,
		fromType: from.typ,
		fromPK:   from.pk,
		toKey:    to.key,
		toType:   to.typ,
		toPK:     to.pk,
	}
	return a
}

func (e entsByPK) getAssocArgs(pks ...string) [][]driver.Value {
	var args [][]driver.Value
	for _, pk := range pks {
		ee := e[pk]
		args = append(args, []driver.Value{ee.key, ee.typ})
	}
	return args
}

func (e entsByPK) getAssocReturns(t *testing.T, pkFrom, pkTo string) []driver.Value {
	a := e.getAssoc(t, pkFrom, pkTo)
	return []driver.Value{a.fromKey, a.fromType, a.toKey, a.toType, a.fromPK, a.toPK}
}

func (e entsByPK) getParentAssocReturns(t *testing.T, pkFrom, pkTo string) []driver.Value {
	a := e.getAssoc(t, pkFrom, pkTo)
	return []driver.Value{a.toKey, a.toType, a.fromKey, a.fromType, a.toPK, a.fromPK}
}

type assoc struct {
	fromKey, fromType, fromPK string
	toKey, toType, toPK       string
}
