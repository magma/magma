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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	magmaerrors "magma/orc8r/lib/go/errors"
)

func TestSQLToEntMigration(t *testing.T) {
	var tableName = "states"
	db, err := sqorc.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	sqlFact := blobstore.NewSQLBlobStorageFactory(tableName, db, sqorc.GetSqlBuilder())
	err = sqlFact.InitializeFactory()
	require.NoError(t, err)

	entFact := blobstore.NewEntStorage(tableName, db, nil)

	checkBlobStoreMigrations(t, sqlFact, entFact)
}

func TestEntToSQLMigration(t *testing.T) {
	var tableName = "states"
	db, err := sqorc.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	entFact := blobstore.NewEntStorage(tableName, db, nil)

	sqlFact := blobstore.NewSQLBlobStorageFactory(tableName, db, sqorc.GetSqlBuilder())
	err = sqlFact.InitializeFactory()
	require.NoError(t, err)

	checkBlobStoreMigrations(t, entFact, sqlFact)
}

func TestIntegration(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	fact := blobstore.NewEntStorage("states", db, sqorc.GetSqlBuilder())
	integration(t, fact)
}

func checkBlobStoreMigrations(t *testing.T, fact1 blobstore.BlobStorageFactory, fact2 blobstore.BlobStorageFactory) {
	var blobs = blobstore.Blobs{
		blobstore.Blob{Type: "type1", Key: "key1", Value: []byte("value1")},
		blobstore.Blob{Type: "type1", Key: "key2", Value: []byte("value2")},
		blobstore.Blob{Type: "type2", Key: "key3", Value: []byte("value3")},
		blobstore.Blob{Type: "type3", Key: "key3", Value: []byte("value4")},
		blobstore.Blob{Type: "type4", Key: "key4", Value: []byte("value5")},
	}
	var blobNotInFact = blobstore.Blob{Type: "notSaved", Key: "notSavedKey", Value: []byte("notSavedValue")}
	var networkID = "id1"

	createBlobs(t, fact1, networkID, blobs)
	checkReadBlobs(t, fact1, networkID, blobs, blobNotInFact)
	checkReadBlobs(t, fact2, networkID, blobs, blobNotInFact)
	checkWriteBlobs(t, fact2, networkID, blobs, blobNotInFact)
}

// Assumes that expectedBlobs length is at least 2.
func checkReadBlobs(
	t *testing.T,
	fact blobstore.BlobStorageFactory,
	networkID string,
	expectedBlobs blobstore.Blobs,
	blobNotInFact blobstore.Blob,
) {
	store, err := fact.StartTransaction(nil)
	require.NoError(t, err)

	keyPrefix := ""
	blobMap, err := store.Search(blobstore.SearchFilter{KeyPrefix: &keyPrefix}, blobstore.LoadCriteria{LoadValue: true})
	assert.ElementsMatch(t, expectedBlobs, blobMap[networkID])

	for _, b := range expectedBlobs {
		checkBlobInStore(t, store, networkID, b)
	}

	blobSlice := expectedBlobs[0:2]
	blobs, err := store.GetMany(networkID, getTKsFromBlobs(blobSlice))
	require.NoError(t, err)
	assert.ElementsMatch(t, blobs, blobSlice)

	keys, err := store.GetExistingKeys([]string{expectedBlobs[0].Key, blobNotInFact.Key}, blobstore.SearchFilter{})
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{expectedBlobs[0].Key}, keys)

	err = store.Commit()
	require.NoError(t, err)
}

// Assumes that expectedBlobs length is at least 2.
// This function should have no side effects.
func checkWriteBlobs(
	t *testing.T,
	fact blobstore.BlobStorageFactory,
	networkID string,
	expectedBlobs blobstore.Blobs,
	blobNotInFact blobstore.Blob,
) {
	checkRollback(t, fact, networkID, blobNotInFact)
	store, err := fact.StartTransaction(nil)
	require.NoError(t, err)

	// Test IncrementVersion
	err = store.IncrementVersion(networkID, getTKFromBlob(expectedBlobs[1]))
	require.NoError(t, err)
	blob, err := store.Get(networkID, getTKFromBlob(expectedBlobs[1]))
	expectedBlobValue := expectedBlobs[1]
	expectedBlobValue.Version = 1
	require.NoError(t, err)
	require.Equal(t, expectedBlobValue, blob)

	err = store.IncrementVersion(networkID, getTKFromBlob(expectedBlobs[1]))
	require.NoError(t, err)
	blob, err = store.Get(networkID, getTKFromBlob(expectedBlobs[1]))
	require.NoError(t, err)
	expectedBlobValue.Version = 2
	require.NoError(t, err)
	require.Equal(t, expectedBlobValue, blob)

	// Test Delete
	err = store.Delete(networkID, storage.TKs{getTKFromBlob(expectedBlobs[1])})
	require.NoError(t, err)
	_, err = store.Get(networkID, getTKFromBlob(expectedBlobs[1]))
	require.Equal(t, magmaerrors.ErrNotFound, err)

	err = store.CreateOrUpdate(networkID, blobstore.Blobs{
		expectedBlobs[1],
		blobNotInFact,
	})
	require.NoError(t, err)

	checkBlobInStore(t, store, networkID, expectedBlobs[1])
	checkBlobInStore(t, store, networkID, blobNotInFact)

	err = store.Delete(networkID, storage.TKs{getTKFromBlob(blobNotInFact)})
	_, err = store.Get(networkID, getTKFromBlob(blobNotInFact))
	require.Equal(t, magmaerrors.ErrNotFound, err)

	err = store.Commit()
	require.NoError(t, err)
}

func checkRollback(
	t *testing.T,
	fact blobstore.BlobStorageFactory,
	networkID string,
	blobNotInFact blobstore.Blob,
) {
	store, err := fact.StartTransaction(nil)
	keyPrefix := ""
	blobMap, err := store.Search(blobstore.SearchFilter{KeyPrefix: &keyPrefix}, blobstore.LoadCriteria{LoadValue: true})
	curBlobs := blobMap[networkID]

	err = store.CreateOrUpdate(networkID, blobstore.Blobs{blobNotInFact})
	checkBlobInStore(t, store, networkID, blobNotInFact)

	err = store.Rollback()
	require.NoError(t, err)

	store, err = fact.StartTransaction(nil)
	blobMap, err = store.Search(blobstore.SearchFilter{KeyPrefix: &keyPrefix}, blobstore.LoadCriteria{LoadValue: true})
	require.NoError(t, err)
	assert.ElementsMatch(t, curBlobs, blobMap[networkID])

	err = store.Commit()
	require.NoError(t, err)
}

func createBlobs(t *testing.T, fact blobstore.BlobStorageFactory, networkID string, blobs blobstore.Blobs) {
	store, err := fact.StartTransaction(nil)
	require.NoError(t, err)

	err = store.CreateOrUpdate(networkID, blobs)
	require.NoError(t, err)

	err = store.Commit()
	require.NoError(t, err)
}

func checkBlobInStore(t *testing.T, store blobstore.TransactionalBlobStorage, networkID string, expectedBlob blobstore.Blob) {
	blob, err := store.Get(networkID, getTKFromBlob(expectedBlob))
	require.NoError(t, err)
	require.Equal(t, expectedBlob, blob)
}

func getTKsFromBlobs(blobs blobstore.Blobs) storage.TKs {
	var tks storage.TKs
	for _, b := range blobs {
		tks = append(tks, getTKFromBlob(b))
	}
	return tks
}

func getTKFromBlob(blob blobstore.Blob) storage.TK {
	return storage.TK{Type: blob.Type, Key: blob.Key}
}
