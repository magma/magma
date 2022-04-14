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

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/merrors"
)

func TestBlobstoreImplMigrations(t *testing.T) {
	makeBlobstores := func() (blobstore.StoreFactory, blobstore.StoreFactory) {
		var tableName = "states"
		db, err := sqorc.Open("sqlite3", ":memory:")
		assert.NoError(t, err)

		sqlFact := blobstore.NewSQLStoreFactory(tableName, db, sqorc.GetSqlBuilder())
		err = sqlFact.InitializeFactory()
		assert.NoError(t, err)

		sqlFact2 := blobstore.NewSQLStoreFactory(tableName, db, sqorc.GetSqlBuilder())
		return sqlFact, sqlFact2
	}

	// Test migration from first blobstore to second
	sqlFact, sqlFact2 := makeBlobstores()
	checkBlobstoreMigration(t, sqlFact, sqlFact2)
	sqlFact, sqlFact2 = makeBlobstores()
	checkBlobstoreMigration(t, sqlFact2, sqlFact)
}

func checkBlobstoreMigration(
	t *testing.T,
	fact1 blobstore.StoreFactory,
	fact2 blobstore.StoreFactory,
) {
	var networkID = "id1"
	var expectedBlobs = blobstore.Blobs{
		blobstore.Blob{Type: "type1", Key: "key1", Value: []byte("value1")},
		blobstore.Blob{Type: "type1", Key: "key2", Value: []byte("value2")},
		blobstore.Blob{Type: "type2", Key: "key3", Value: []byte("value3")},
		blobstore.Blob{Type: "type3", Key: "key3", Value: []byte("value4")},
		blobstore.Blob{Type: "type4", Key: "key4", Value: []byte("value5")},
	}
	var blobNotInFact = blobstore.Blob{Type: "notSaved", Key: "notSavedKey", Value: []byte("notSavedValue")}

	createBlobs(t, fact1, networkID, expectedBlobs)
	checkReadBlobs(t, fact1, networkID, expectedBlobs, blobNotInFact)
	checkReadBlobs(t, fact2, networkID, expectedBlobs, blobNotInFact)
	checkWriteBlobs(t, fact2, networkID, expectedBlobs, blobNotInFact)
}

// checkReadBlobs checks that fact includes all blobs in expectedBlobs.
// Assumes that expectedBlobs length is at least 2.
func checkReadBlobs(
	t *testing.T,
	fact blobstore.StoreFactory,
	networkID string,
	expectedBlobs blobstore.Blobs,
	blobNotInFact blobstore.Blob,
) {
	store, err := fact.StartTransaction(nil)
	assert.NoError(t, err)

	keyPrefix := ""
	actualBlobMap, err := store.Search(blobstore.SearchFilter{KeyPrefix: &keyPrefix}, blobstore.LoadCriteria{LoadValue: true})
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedBlobs, actualBlobMap[networkID])

	for _, b := range expectedBlobs {
		checkBlobInStore(t, store, networkID, b)
	}

	expectedBlobSlice := expectedBlobs[0:2]
	blobs, err := store.GetMany(networkID, expectedBlobSlice.TKs())
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedBlobSlice, blobs)

	keys, err := store.GetExistingKeys([]string{expectedBlobs[0].Key, blobNotInFact.Key}, blobstore.SearchFilter{})
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{expectedBlobs[0].Key}, keys)

	err = store.Commit()
	assert.NoError(t, err)
}

// checkWriteBlobs checks that writing to fact works as expected.
// Assumes that expectedBlobs length is at least 2.
// This function should have no side effects.
func checkWriteBlobs(
	t *testing.T,
	fact blobstore.StoreFactory,
	networkID string,
	expectedBlobs blobstore.Blobs,
	blobNotInFact blobstore.Blob,
) {
	checkRollback(t, fact, networkID, blobNotInFact)
	store, err := fact.StartTransaction(nil)
	assert.NoError(t, err)

	// Test IncrementVersion
	err = store.IncrementVersion(networkID, expectedBlobs[1].TK())
	assert.NoError(t, err)
	blob, err := store.Get(networkID, expectedBlobs[1].TK())
	expectedBlobValue := expectedBlobs[1]
	expectedBlobValue.Version = 1
	assert.NoError(t, err)
	assert.Equal(t, expectedBlobValue, blob)

	err = store.IncrementVersion(networkID, expectedBlobs[1].TK())
	assert.NoError(t, err)
	blob, err = store.Get(networkID, expectedBlobs[1].TK())
	assert.NoError(t, err)
	expectedBlobValue.Version = 2
	assert.NoError(t, err)
	assert.Equal(t, expectedBlobValue, blob)

	// Test Delete
	err = store.Delete(networkID, storage.TKs{expectedBlobs[1].TK()})
	assert.NoError(t, err)
	_, err = store.Get(networkID, expectedBlobs[1].TK())
	assert.Equal(t, merrors.ErrNotFound, err)

	err = store.Write(networkID, blobstore.Blobs{
		expectedBlobs[1],
		blobNotInFact,
	})
	assert.NoError(t, err)

	checkBlobInStore(t, store, networkID, expectedBlobs[1])
	checkBlobInStore(t, store, networkID, blobNotInFact)

	err = store.Delete(networkID, storage.TKs{blobNotInFact.TK()})
	assert.NoError(t, err)
	_, err = store.Get(networkID, blobNotInFact.TK())
	assert.Equal(t, merrors.ErrNotFound, err)

	err = store.Commit()
	assert.NoError(t, err)
}

func checkRollback(
	t *testing.T,
	fact blobstore.StoreFactory,
	networkID string,
	blobNotInFact blobstore.Blob,
) {
	store, err := fact.StartTransaction(nil)
	assert.NoError(t, err)
	keyPrefix := ""
	blobMap, err := store.Search(blobstore.SearchFilter{KeyPrefix: &keyPrefix}, blobstore.LoadCriteria{LoadValue: true})
	assert.NoError(t, err)
	inputBlobs := blobMap[networkID]

	err = store.Write(networkID, blobstore.Blobs{blobNotInFact})
	assert.NoError(t, err)
	checkBlobInStore(t, store, networkID, blobNotInFact)

	err = store.Rollback()
	assert.NoError(t, err)

	store, err = fact.StartTransaction(nil)
	assert.NoError(t, err)
	blobMap, err = store.Search(blobstore.SearchFilter{KeyPrefix: &keyPrefix}, blobstore.LoadCriteria{LoadValue: true})
	assert.NoError(t, err)
	assert.ElementsMatch(t, inputBlobs, blobMap[networkID])

	err = store.Commit()
	assert.NoError(t, err)
}

func createBlobs(
	t *testing.T,
	fact blobstore.StoreFactory,
	networkID string,
	blobs blobstore.Blobs,
) {
	store, err := fact.StartTransaction(nil)
	assert.NoError(t, err)

	err = store.Write(networkID, blobs)
	assert.NoError(t, err)

	err = store.Commit()
	assert.NoError(t, err)
}

func checkBlobInStore(
	t *testing.T,
	store blobstore.Store,
	networkID string,
	expectedBlob blobstore.Blob,
) {
	blob, err := store.Get(networkID, expectedBlob.TK())
	assert.NoError(t, err)
	assert.Equal(t, expectedBlob, blob)
}
