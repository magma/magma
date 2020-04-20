/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test_utils

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/sqorc"
	storage2 "magma/orc8r/cloud/go/storage"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var (
	once     sync.Once
	instance *sql.DB
)

// GetSharedMemoryDB returns a singleton in-memory database connection.
func GetSharedMemoryDB(t *testing.T) *sql.DB {
	once.Do(func() {
		db, err := sqorc.Open(storage2.SQLDriver, ":memory:")
		assert.NoError(t, err)
		instance = db
	})
	return instance
}

// DropTableFromSharedTestDB drops the table from the singleton in-memory database.
func DropTableFromSharedTestDB(t *testing.T, table string) {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
	_, err := instance.Exec(query)
	assert.NoError(t, err)
}

// NewSQLBlobstore returns a new blobstore storage factory utilizing the singleton in-memory database.
func NewSQLBlobstore(t *testing.T, tableName string) blobstore.BlobStorageFactory {
	db := GetSharedMemoryDB(t)
	store := blobstore.NewSQLBlobStorageFactory(tableName, db, sqorc.GetSqlBuilder())

	err := store.InitializeFactory()
	assert.NoError(t, err)

	return store
}
