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

package JsonStore_test

import (
	"magma/orc8r/cloud/go/JsonStore"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/merrors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonStoreImplMigrations(t *testing.T) {
	makeJsonStores := func() (JsonStore.StoreFactory, JsonStore.StoreFactory){
		var tableName = "states"
		db, err := sqorc.Open("sqlite3",":memory:")
		assert.NoError(t, err)

		sqlFact := JsonStore.NewSQLStoreFactory(tableName, db, sqorc.GetSqlBuilder())
		err = sqlFact.InitializeFactory()
		assert.NoError(t, err)

		sqlFact2 := JsonStore.NewSQLStoreFactory(tableName, db, sqorc.GetSqlBuilder())
		err = sqlFact.InitializeFactory()
		assert.NoError(t, err)
		return sqlFact, sqlFact2
	}

	sqlFact, sqlFact2 := makeJsonStores()
	checkJsonStoreMigration(t, sqlFact, sqlFact2)
	sqlFact, sqlFact2 = makeJsonStores()
	checkJsonStoreMigration(t, sqlFact2, sqlFact)
}

func checkJsonStoreMigration(
	t *testing.T,
	fact1 JsonStore.StoreFactory,
	fact2 JsonStore.StoreFactory,
) {
	var networkID = "id1"
	var expectedJsons = JsonStore.Jsons{
		JsonStore.Json{Type: "type1", Key: "key1", Value: "value1"},
		JsonStore.Json{Type: "type1", Key: "key2", Value: "value2"},
		JsonStore.Json{Type: "type2", Key: "key3", Value: "value3"},
		JsonStore.Json{Type: "type3", Key: "key3", Value: "value4"},
		JsonStore.Json{Type: "type4", Key: "key4", Value: "value5"},
	}
	var JsonNotInFact = JsonStore.Json{Type: "notSaved", Key: "notSavedKey", Value: "notSavedValue"}
	createJsons(t, fact1, networkID, expectedJsons)
	checkReadJsons(t, fact1, networkID, expectedJsons, JsonNotInFact)
	checkReadJsons(t, fact2, networkID, expectedJsons, JsonNotInFact)
	checkWriteJsons(t, fact2, networkID, expectedJsons, JsonNotInFact)
}

func checkReadJsons(
	t *testing.T,
	fact JsonStore.StoreFactory,
	networkID string,
	expectedJsons JsonStore.Jsons,
	JsonInFact JsonStore.Json,
) {
	store, err := fact.StartTransaction(nil)
	assert.NoError(t, err)

	keyPrefix := ""
	actualBlobMap, err := store.Search(JsonStore.SearchFilter{KeyPrefix: &keyPrefix}, JsonStore.LoadCriteria{LoadValue: true})
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedJsons, actualBlobMap[networkID])

	for _, b := range expectedJsons {
		checkJsonInStore(t, store, networkID, b)
	}

	expectedJsonSlice := expectedJsons[0:2]
	Jsons, err := store.GetMany(networkID, expectedJsonSlice.TKs())
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedJsonSlice, Jsons)

	keys, err := store.GetExistingKeys([]string{expectedJsons[0].Key, JsonInFact.Key}, JsonStore.SearchFilter{})
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{expectedJsons[0].Key}, keys)

	err = store.Commit()
	assert.NoError(t, err)	
}


func checkWriteJsons(
	t *testing.T,
	fact JsonStore.StoreFactory,
	networkID string,
	expectedJsons JsonStore.Jsons,
	JsonNotInFact JsonStore.Json,
) {
	checkRollback(t, fact, networkID, JsonNotInFact)
	store, err := fact.StartTransaction(nil)
	assert.NoError(t, err)

	err = store.IncrementVersion(networkID, expectedJsons[1].TK())
	assert.NoError(t, err)
	Json, err := store.Get(networkID, expectedJsons[1].TK())
	expectedJsonValue := expectedJsons[1]
	expectedJsonValue.Version = 1
	assert.NoError(t, err)
	assert.Equal(t, expectedJsonValue, Json)

	err = store.IncrementVersion(networkID, expectedJsons[1].TK())
	assert.NoError(t, err)
	Json, err = store.Get(networkID, expectedJsons[1].TK())
	assert.NoError(t, err)
	expectedJsonValue.Version = 2
	assert.NoError(t, err)
	assert.Equal(t, expectedJsonValue, Json)

	//Test Delete
	err = store.Delete(networkID, storage.TKs{expectedJsons[1].TK()})
	assert.NoError(t, err)
	_, err = store.Get(networkID, expectedJsons[1].TK())
	assert.Equal(t, merrors.ErrNotFound, err)

	err = store.Write(networkID, JsonStore.Jsons{
		expectedJsons[1],
		JsonNotInFact,
	})
	assert.NoError(t, err)

	checkJsonInStore(t, store, networkID, expectedJsons[1])
	checkJsonInStore(t, store, networkID, JsonNotInFact)

	err = store.Delete(networkID, storage.TKs{JsonNotInFact.TK()})
	assert.NoError(t, err)

	_, err = store.Get(networkID, JsonNotInFact.TK())
	assert.Equal(t, merrors.ErrNotFound, err)

	err = store.Commit()
	assert.NoError(t, err)
}




func checkRollback(
	t *testing.T,
	fact JsonStore.StoreFactory,
	networkID string,
	blobNotInFact JsonStore.Json,
) {
	store, err := fact.StartTransaction(nil)
	assert.NoError(t, err)
	keyPrefix := ""
	JsonMap, err := store.Search(JsonStore.SearchFilter{KeyPrefix: &keyPrefix}, JsonStore.LoadCriteria{LoadValue: true})
	assert.NoError(t, err)
	inputJsons := JsonMap[networkID]

	err = store.Write(networkID, JsonStore.Jsons{blobNotInFact})
	assert.NoError(t, err)
	checkJsonInStore(t, store, networkID, blobNotInFact)

	err = store.Rollback()
	assert.NoError(t, err)

	store, err = fact.StartTransaction(nil)
	assert.NoError(t, err)
	JsonMap, err = store.Search(JsonStore.SearchFilter{KeyPrefix: &keyPrefix}, JsonStore.LoadCriteria{LoadValue: true})
	assert.NoError(t, err)
	assert.ElementsMatch(t, inputJsons, JsonMap[networkID])

	err = store.Commit()
	assert.NoError(t, err)
}


func createJsons(
	t *testing.T,
	fact JsonStore.StoreFactory,
	networkID string,
	Jsons JsonStore.Jsons,
) {
	store, err := fact.StartTransaction(nil)
	assert.NoError(t, err)

	err = store.Write(networkID, Jsons)
	assert.NoError(t, err)

	err = store.Commit()
	assert.NoError(t, err)
}

func checkJsonInStore(
	t *testing.T,
	store JsonStore.Store,
	networkID string,
	expectedJson JsonStore.Json,
) {
	Json,  err := store.Get(networkID, expectedJson.TK())
	assert.NoError(t,err)
	assert.Equal(t, expectedJson, Json)
}
