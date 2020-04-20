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

package store

import (
	"testing"
	"time"

	"magma/fbinternal/cloud/go/services/testcontroller/obsidian/models"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type mockClock struct {
	now time.Time
}

func (mockClock *mockClock) Now() time.Time {
	return mockClock.now
}

func TestTestControllerStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()

	factory := blobstore.NewSQLBlobStorageFactory("network_table", db, sqorc.GetSqlBuilder())
	expectCreateTable(mock)
	err = factory.InitializeFactory()
	assert.NoError(t, err)
	store, err := NewTestControllerStore(factory)
	assert.NoError(t, err)

	executionGet := &models.LatestScriptExecution{
		Version:   "",
		Timestamp: int64(1075593600),
	}

	executionPut := &models.LatestScriptExecution{
		Version:   "1.2.3-1075599560-0",
		Timestamp: int64(1075593600),
	}

	marshaledGet, err := executionGet.MarshalBinary()
	assert.NoError(t, err)

	marshaledPut, err := executionPut.MarshalBinary()
	assert.NoError(t, err)

	// Nothing in database, put current version & timestamp
	mock.ExpectBegin()
	expectUpsert(mock)
	expectGet(mock, marshaledGet)
	expectPut(mock, marshaledPut)
	mock.ExpectCommit()

	store.SetClock(t, &mockClock{time.Unix(1075593600, 0)})
	execute, err := store.ShouldExecuteScript("network", "1.2.3-1075599560-0", 0)
	assert.NoError(t, err)
	assert.True(t, execute)

	executionPut = &models.LatestScriptExecution{
		Version:   "1.2.3-1475593600-0",
		Timestamp: int64(1475593600),
	}

	marshaledPut, err = executionPut.MarshalBinary()
	assert.NoError(t, err)

	// Gateway upgraded: execute script and put current version & timestamp
	mock.ExpectBegin()
	expectUpsert(mock)
	expectGet(mock, marshaledGet)
	expectPut(mock, marshaledPut)
	mock.ExpectCommit()

	store.SetClock(t, &mockClock{time.Unix(1475593600, 0)})
	execute, err = store.ShouldExecuteScript("network", "1.2.3-1475593600-0", 0)
	assert.NoError(t, err)
	assert.True(t, execute)

	executionGet = &models.LatestScriptExecution{
		Version:   "1.2.3-1475593600-0",
		Timestamp: int64(1075593600),
	}

	marshaledGet, err = executionGet.MarshalBinary()
	assert.NoError(t, err)

	// Not enough time elapsed: don't execute script
	mock.ExpectBegin()
	expectUpsert(mock)
	expectGet(mock, marshaledGet)
	mock.ExpectCommit()

	execute, err = store.ShouldExecuteScript("network", "1.2.4-1475593600-0", 1)
	assert.NoError(t, err)
	assert.False(t, execute)
}

func expectCreateTable(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS network_table").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func expectUpsert(mock sqlmock.Sqlmock) {
	mock.ExpectExec("INSERT INTO network_table \\(network_id,type,\"key\",version\\) "+
		"VALUES \\(\\$1,\\$2,\\$3,\\$4\\) "+
		"ON CONFLICT \\(network_id, type, \"key\"\\) "+
		"DO UPDATE SET version = ",
	).
		WithArgs("network", "testcontroller", "gateway_version", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func expectPut(mock sqlmock.Sqlmock, upgrade []byte) {
	mock.ExpectQuery("SELECT type, \"key\", value, version FROM network_table").
		WithArgs("network", "testcontroller", "gateway_version").
		WillReturnRows(
			sqlmock.NewRows([]string{"type", "key", "value", "version"}).
				AddRow("type", "key", []byte("value1"), 1),
		)

	mock.ExpectExec("INSERT INTO network_table").
		WithArgs("network", "testcontroller", "gateway_version", upgrade, 0).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func expectGet(mock sqlmock.Sqlmock, upgrade []byte) {
	mock.ExpectQuery("SELECT type, \"key\", value, version FROM network_table").
		WithArgs("network", "testcontroller", "gateway_version").
		WillReturnRows(
			sqlmock.NewRows([]string{"type", "key", "value", "version"}).
				AddRow("type", "key", upgrade, 1),
		)
}
