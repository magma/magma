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

package blobstore_test

import (
	"database/sql/driver"
	"errors"
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestSqlBlobStorage_Get(t *testing.T) {
	happyPath := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(
				"SELECT type, \"key\", value, version FROM network_table "+
					"WHERE \\(\\(network_id = \\$1 AND type = \\$2 AND \"key\" = \\$3\\)\\)",
			).
				WithArgs("network", "t1", "k1").
				WillReturnRows(
					sqlmock.NewRows([]string{"type", "key", "value", "version"}).
						AddRow("t1", "k1", []byte("value1"), 42),
				)
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			return store.Get("network", storage.TypeAndKey{Type: "t1", Key: "k1"})
		},

		expectedError:  nil,
		expectedResult: blobstore.Blob{Type: "t1", Key: "k1", Value: []byte("value1"), Version: 42},
	}
	dneCase := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(
				"SELECT type, \"key\", value, version FROM network_table "+
					"WHERE \\(\\(network_id = \\$1 AND type = \\$2 AND \"key\" = \\$3\\)\\)",
			).
				WithArgs("network", "t2", "k2").
				WillReturnRows(
					sqlmock.NewRows([]string{"type", "key", "value", "version"}),
				)
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			return store.Get("network", storage.TypeAndKey{Type: "t2", Key: "k2"})
		},

		expectedError:      merrors.ErrNotFound,
		matchErrorInstance: true,
		expectedResult:     nil,
	}
	queryError := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(
				"SELECT type, \"key\", value, version FROM network_table "+
					"WHERE \\(\\(network_id = \\$1 AND type = \\$2 AND \"key\" = \\$3\\)\\)",
			).
				WithArgs("network", "t3", "k3").
				WillReturnError(errors.New("mock query error"))
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			return store.Get("network", storage.TypeAndKey{Type: "t3", Key: "k3"})
		},

		expectedError:  errors.New("mock query error"),
		expectedResult: nil,
	}
	runCase(t, happyPath)
	runCase(t, dneCase)
	runCase(t, queryError)
}

func TestSqlBlobStorage_GetMany(t *testing.T) {
	happyPath := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(
				"SELECT type, \"key\", value, version FROM network_table "+
					"WHERE \\("+
					"\\(network_id = \\$1 AND type = \\$2 AND \"key\" = \\$3\\) OR "+
					"\\(network_id = \\$4 AND type = \\$5 AND \"key\" = \\$6\\)\\)").
				WithArgs("network", "t1", "k1", "network", "t2", "k2").
				WillReturnRows(
					sqlmock.NewRows([]string{"type", "key", "value", "version"}).
						AddRow("t1", "k1", []byte("value1"), 42).
						AddRow("t2", "k2", []byte("value2"), 43),
				)
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			return store.GetMany("network", []storage.TypeAndKey{{Type: "t1", Key: "k1"}, {Type: "t2", Key: "k2"}})
		},

		expectedError: nil,
		expectedResult: blobstore.Blobs{
			{Type: "t1", Key: "k1", Value: []byte("value1"), Version: 42},
			{Type: "t2", Key: "k2", Value: []byte("value2"), Version: 43},
		},
	}

	queryError := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("SELECT type, \"key\", value, version FROM network_table").
				WithArgs("network", "t1", "k1", "network", "t2", "k2").
				WillReturnError(errors.New("mock query error"))
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			return store.GetMany("network", []storage.TypeAndKey{{Type: "t1", Key: "k1"}, {Type: "t2", Key: "k2"}})
		},

		expectedError:  errors.New("mock query error"),
		expectedResult: nil,
	}

	runCase(t, happyPath)
	runCase(t, queryError)
}

func TestSqlBlobStorage_Search(t *testing.T) {
	happyPath := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("SELECT network_id, type, \"key\", version, value FROM network_table").
				WithArgs("network", "t1", "t2", "t3", "k1", "k2", "k3").
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "type", "key", "version", "value"}).
						AddRow("network", "t1", "k1", 42, []byte("value1")).
						AddRow("network", "t2", "k2", 43, []byte("value2")),
				)
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			return store.Search(
				blobstore.CreateSearchFilter(strPtr("network"), []string{"t1", "t2", "t3"}, []string{"k1", "k2", "k3"}, nil),
				blobstore.GetDefaultLoadCriteria(),
			)
		},

		expectedError: nil,
		expectedResult: map[string]blobstore.Blobs{
			"network": {
				{Type: "t1", Key: "k1", Value: []byte("value1"), Version: 42},
				{Type: "t2", Key: "k2", Value: []byte("value2"), Version: 43},
			},
		},
	}

	keyPrefix := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("SELECT network_id, type, \"key\", version, value FROM network_table").
				WithArgs("network", "t1", "t2", "kprefix%").
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "type", "key", "version", "value"}).
						AddRow("network", "t1", "kprefix1", 42, []byte("value1")).
						AddRow("network", "t2", "kprefix2", 43, []byte("value2")),
				)
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			return store.Search(
				blobstore.CreateSearchFilter(strPtr("network"), []string{"t1", "t2"}, nil, strPtr("kprefix")),
				blobstore.GetDefaultLoadCriteria(),
			)
		},

		expectedError: nil,
		expectedResult: map[string]blobstore.Blobs{
			"network": {
				{Type: "t1", Key: "kprefix1", Value: []byte("value1"), Version: 42},
				{Type: "t2", Key: "kprefix2", Value: []byte("value2"), Version: 43},
			},
		},
	}

	emptyFilterReturnsAll := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("SELECT network_id, type, \"key\", version, value FROM network_table").
				WithArgs(). // no args
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "type", "key", "version", "value"}).
						AddRow("network1", "t1", "k1", 42, []byte("value1")).
						AddRow("network1", "t2", "k2", 43, []byte("value2")).
						AddRow("network2", "t3", "k3", 44, []byte("value3")),
				)
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			return store.Search(
				blobstore.CreateSearchFilter(nil, nil, nil, nil),
				blobstore.GetDefaultLoadCriteria(),
			)
		},

		expectedError: nil,
		expectedResult: map[string]blobstore.Blobs{
			"network1": {
				{Type: "t1", Key: "k1", Value: []byte("value1"), Version: 42},
				{Type: "t2", Key: "k2", Value: []byte("value2"), Version: 43},
			},
			"network2": {
				{Type: "t3", Key: "k3", Value: []byte("value3"), Version: 44},
			},
		},
	}

	multipleNetworks := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("SELECT network_id, type, \"key\", version, value FROM network_table").
				WithArgs("t1", "t2", "t3", "k1", "k2", "k3").
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "type", "key", "version", "value"}).
						AddRow("network1", "t1", "k1", 42, []byte("value1")).
						AddRow("network2", "t2", "k2", 43, []byte("value2")),
				)
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			return store.Search(
				blobstore.CreateSearchFilter(nil, []string{"t1", "t2", "t3"}, []string{"k1", "k2", "k3"}, nil),
				blobstore.GetDefaultLoadCriteria(),
			)
		},

		expectedError: nil,
		expectedResult: map[string]blobstore.Blobs{
			"network1": {
				{Type: "t1", Key: "k1", Value: []byte("value1"), Version: 42},
			},
			"network2": {
				{Type: "t2", Key: "k2", Value: []byte("value2"), Version: 43},
			},
		},
	}

	loadCriteria := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("SELECT network_id, type, \"key\", version FROM network_table").
				WithArgs("t1", "t2", "t3", "k1", "k2", "k3").
				WillReturnRows(
					sqlmock.NewRows([]string{"network_id", "type", "key", "version"}).
						AddRow("network1", "t1", "k1", 42).
						AddRow("network2", "t2", "k2", 43),
				)
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			return store.Search(
				blobstore.CreateSearchFilter(nil, []string{"t1", "t2", "t3"}, []string{"k1", "k2", "k3"}, nil),
				blobstore.LoadCriteria{LoadValue: false},
			)
		},

		expectedError: nil,
		expectedResult: map[string]blobstore.Blobs{
			"network1": {
				{Type: "t1", Key: "k1", Value: nil, Version: 42},
			},
			"network2": {
				{Type: "t2", Key: "k2", Value: nil, Version: 43},
			},
		},
	}

	queryError := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("SELECT network_id, type, \"key\", version, value FROM network_table").
				WithArgs("network", "t1", "t2", "t3", "k1", "k2", "k3").
				WillReturnError(errors.New("mock error"))
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			return store.Search(
				blobstore.CreateSearchFilter(strPtr("network"), []string{"t1", "t2", "t3"}, []string{"k1", "k2", "k3"}, nil),
				blobstore.GetDefaultLoadCriteria(),
			)
		},

		expectedError: errors.New("failed to query DB: mock error"),
	}

	runCase(t, happyPath)
	runCase(t, keyPrefix)
	runCase(t, emptyFilterReturnsAll)
	runCase(t, multipleNetworks)
	runCase(t, loadCriteria)
	runCase(t, queryError)
}

func TestSqlBlobStorage_CreateOrUpdate(t *testing.T) {
	// (t1, k1) exists, (t2, k2) does not
	happyPath := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			expectGetMany(
				mock,
				[]driver.Value{"network", "t1", "k1", "network", "t2", "k2"},
				blobstore.Blobs{
					{Type: "t1", Key: "k1", Value: []byte("hello"), Version: 42},
				},
			)

			updatePrepare := mock.ExpectPrepare("UPDATE network_table")
			updatePrepare.ExpectExec().
				WithArgs([]byte("goodbye"), 43, "network", "t1", "k1").
				WillReturnResult(sqlmock.NewResult(1, 1))
			updatePrepare.WillBeClosed()

			mock.ExpectExec("INSERT INTO network_table").
				WithArgs("network", "t2", "k2", []byte("world"), 1000).
				WillReturnResult(sqlmock.NewResult(1, 1))
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			err := store.CreateOrUpdate(
				"network",
				blobstore.Blobs{
					{Type: "t1", Key: "k1", Value: []byte("goodbye"), Version: 0},
					{Type: "t2", Key: "k2", Value: []byte("world"), Version: 1000},
				},
			)
			return nil, err
		},

		expectedError:  nil,
		expectedResult: nil,
	}

	updateOnly := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			expectGetMany(
				mock,
				[]driver.Value{"network", "t1", "k1", "network", "t2", "k2"},
				blobstore.Blobs{
					{Type: "t1", Key: "k1", Value: []byte("hello"), Version: 42},
					{Type: "t2", Key: "k2", Value: []byte("world"), Version: 43},
				},
			)

			updatePrepare := mock.ExpectPrepare("UPDATE network_table")
			updatePrepare.ExpectExec().
				WithArgs([]byte("goodbye"), 100, "network", "t1", "k1").
				WillReturnResult(sqlmock.NewResult(1, 1))
			updatePrepare.ExpectExec().
				WithArgs([]byte("foo"), 44, "network", "t2", "k2").
				WillReturnResult(sqlmock.NewResult(1, 1))
			updatePrepare.WillBeClosed()
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			err := store.CreateOrUpdate(
				"network",
				blobstore.Blobs{
					{Type: "t1", Key: "k1", Value: []byte("goodbye"), Version: 100},
					{Type: "t2", Key: "k2", Value: []byte("foo"), Version: 0},
				},
			)
			return nil, err
		},

		expectedError:  nil,
		expectedResult: nil,
	}

	insertOnly := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			expectGetMany(
				mock,
				[]driver.Value{"network", "t1", "k1", "network", "t2", "k2"},
				nil,
			)

			mock.ExpectExec("INSERT INTO network_table").
				WithArgs(
					"network", "t1", "k1", []byte("hello"), 0,
					"network", "t2", "k2", []byte("world"), 1000,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			err := store.CreateOrUpdate(
				"network",
				blobstore.Blobs{
					{Type: "t1", Key: "k1", Value: []byte("hello"), Version: 0},
					{Type: "t2", Key: "k2", Value: []byte("world"), Version: 1000},
				},
			)
			return nil, err
		},

		expectedError:  nil,
		expectedResult: nil,
	}

	updateError := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			expectGetMany(
				mock,
				[]driver.Value{"network", "t1", "k1", "network", "t2", "k2"},
				blobstore.Blobs{
					{Type: "t1", Key: "k1", Value: []byte("hello"), Version: 42},
				},
			)

			updatePrepare := mock.ExpectPrepare("UPDATE network_table")
			updatePrepare.ExpectExec().
				WithArgs([]byte("goodbye"), 43, "network", "t1", "k1").
				WillReturnError(errors.New("mock query error"))
			updatePrepare.WillBeClosed()
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			err := store.CreateOrUpdate(
				"network",
				blobstore.Blobs{
					{Type: "t1", Key: "k1", Value: []byte("goodbye"), Version: 0},
					{Type: "t2", Key: "k2", Value: []byte("world"), Version: 1000},
				},
			)
			return nil, err
		},

		expectedError:  errors.New("error updating blob (network, t1, k1): mock query error"),
		expectedResult: nil,
	}

	insertError := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			expectGetMany(
				mock,
				[]driver.Value{"network", "t1", "k1", "network", "t2", "k2"},
				blobstore.Blobs{
					{Type: "t1", Key: "k1", Value: []byte("hello"), Version: 42},
				},
			)

			updatePrepare := mock.ExpectPrepare("UPDATE network_table")
			updatePrepare.ExpectExec().
				WithArgs([]byte("goodbye"), 43, "network", "t1", "k1").
				WillReturnResult(sqlmock.NewResult(1, 1))
			updatePrepare.WillBeClosed()

			mock.ExpectExec("INSERT INTO network_table").
				WithArgs("network", "t2", "k2", []byte("world"), 1000).
				WillReturnError(errors.New("mock query error"))
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			err := store.CreateOrUpdate(
				"network",
				blobstore.Blobs{
					{Type: "t1", Key: "k1", Value: []byte("goodbye"), Version: 0},
					{Type: "t2", Key: "k2", Value: []byte("world"), Version: 1000},
				},
			)
			return nil, err
		},

		expectedError:  errors.New("error creating blobs: mock query error"),
		expectedResult: nil,
	}

	runCase(t, happyPath)
	runCase(t, updateOnly)
	runCase(t, insertOnly)
	runCase(t, updateError)
	runCase(t, insertError)
}

func TestSqlBlobStorage_Delete(t *testing.T) {
	happyPath := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectExec("DELETE FROM network_table").
				WithArgs("network", "t1", "k1", "network", "t2", "k2").
				WillReturnResult(sqlmock.NewResult(1, 1))
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			err := store.Delete("network", []storage.TypeAndKey{{Type: "t1", Key: "k1"}, {Type: "t2", Key: "k2"}})
			return nil, err
		},

		expectedError:  nil,
		expectedResult: nil,
	}

	queryError := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectExec("DELETE FROM network_table").
				WithArgs("network", "t1", "k1", "network", "t2", "k2").
				WillReturnError(errors.New("mock query error"))
		},

		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			err := store.Delete("network", []storage.TypeAndKey{{Type: "t1", Key: "k1"}, {Type: "t2", Key: "k2"}})
			return nil, err
		},

		expectedError:  errors.New("mock query error"),
		expectedResult: nil,
	}

	runCase(t, happyPath)
	runCase(t, queryError)
}

func TestSqlBlobStorage_IncrementVersion(t *testing.T) {
	happyPath := &testCase{
		setup: func(mock sqlmock.Sqlmock) {
			mock.ExpectExec("INSERT INTO network_table \\(network_id,type,\"key\",version\\) "+
				"VALUES \\(\\$1,\\$2,\\$3,\\$4\\) "+
				"ON CONFLICT \\(network_id, type, \"key\"\\) "+
				"DO UPDATE SET version = ",
			).
				WithArgs("network", "t1", "k1", 1).
				WillReturnResult(sqlmock.NewResult(1, 1))
		},
		run: func(store blobstore.TransactionalBlobStorage) (interface{}, error) {
			err := store.IncrementVersion("network", storage.TypeAndKey{Type: "t1", Key: "k1"})
			return nil, err
		},
		expectedError:  nil,
		expectedResult: nil,
	}

	runCase(t, happyPath)
}

func TestSqlBlobStorage_Integration(t *testing.T) {
	// Use an in-memory sqlite data store
	db, err := sqorc.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Could not initialize sqlite DB: %s", err)
	}
	fact := blobstore.NewSQLBlobStorageFactory("network_table", db, sqorc.GetSqlBuilder())
	integration(t, fact)
}

type testCase struct {
	// setup query expectations (begin/table init is generically handled)
	setup func(sqlmock.Sqlmock)

	// run the test case
	run func(blobstore.TransactionalBlobStorage) (interface{}, error)

	expectedError      error
	matchErrorInstance bool
	expectedResult     interface{}
}

func runCase(t *testing.T, test *testCase) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}

	factory := blobstore.NewSQLBlobStorageFactory("network_table", db, sqorc.GetSqlBuilder())
	expectCreateTable(mock)
	err = factory.InitializeFactory()
	assert.NoError(t, err)

	mock.ExpectBegin()
	store, err := factory.StartTransaction(nil)
	assert.NoError(t, err)

	test.setup(mock)
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
}

func expectCreateTable(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS network_table").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func expectGetMany(mock sqlmock.Sqlmock, args []driver.Value, blobs blobstore.Blobs) {
	rows := sqlmock.NewRows([]string{"type", "key", "value", "version"})
	for _, blob := range blobs {
		rows.AddRow(blob.Type, blob.Key, blob.Value, blob.Version)
	}

	mock.ExpectQuery("SELECT type, \"key\", value, version FROM network_table").
		WithArgs(args...).
		WillReturnRows(rows)
}
