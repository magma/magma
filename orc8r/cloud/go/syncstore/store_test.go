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

package syncstore_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/clock"
	configurator_storage "magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/syncstore"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

func TestSyncStore(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := test_utils.NewSQLBlobstore(t, "last_resync_time")
	assert.NoError(t, fact.InitializeFactory())
	config := syncstore.Config{
		TableNamePrefix:              "test",
		CacheWriterValidIntervalSecs: 150,
	}
	store, err := syncstore.NewSyncStore(db, sqorc.GetSqlBuilder(), fact, config)
	assert.NoError(t, err)
	assert.NoError(t, store.Initialize())

	expectedDigestTree := &protos.DigestTree{
		RootDigest: &protos.Digest{Md5Base64Digest: "root_digest_apple2"},
		LeafDigests: []*protos.LeafDigest{
			{Id: "2", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_banana2"}},
			{Id: "3", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_cherry2"}},
			{Id: "4", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_dragonfruit"}},
		},
	}
	expectedDigestTree2 := &protos.DigestTree{
		RootDigest: &protos.Digest{Md5Base64Digest: "root_digest_banana"},
	}
	objs1 := map[string][]byte{
		"1": []byte("apple"),
	}
	objs2 := map[string][]byte{
		"2": []byte("banana"),
		"3": []byte("cherry"),
	}
	expectedObjs := [][]byte{
		[]byte("apple"),
		[]byte("banana"),
		[]byte("cherry"),
	}

	t.Run("initially empty", func(t *testing.T) {
		digestTrees, err := store.GetDigests([]string{"n0", "n1"}, time.Now().Unix(), true)
		assert.NoError(t, err)
		assert.Empty(t, digestTrees)

		page, nextToken, err := store.GetCachedByPage("n0", "", 10)
		assert.NoError(t, err)
		assert.Empty(t, page)
		assert.Empty(t, nextToken)

		lastResync, err := store.GetLastResync("n0", "g0")
		assert.NoError(t, err)
		assert.Empty(t, lastResync)
	})

	t.Run("basic insert digests", func(t *testing.T) {
		expectedDigestTree := &protos.DigestTree{
			RootDigest: &protos.Digest{Md5Base64Digest: "root_digest_apple"},
			LeafDigests: []*protos.LeafDigest{
				{Id: "1", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_apple"}},
				{Id: "2", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_banana"}},
				{Id: "3", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_cherry"}},
			},
		}
		err := store.SetDigest("n0", expectedDigestTree)
		assert.NoError(t, err)

		digestTrees, err := store.GetDigests([]string{"n0"}, time.Now().Unix(), true)
		assert.NoError(t, err)
		assert.Contains(t, digestTrees, "n0")
		assert.True(t, proto.Equal(expectedDigestTree, digestTrees["n0"]))

		digestTrees, err = store.GetDigests([]string{"n0"}, time.Now().Unix(), false)
		assert.NoError(t, err)
		assert.Contains(t, digestTrees, "n0")
		assert.Equal(t, "root_digest_apple", digestTrees["n0"].RootDigest.Md5Base64Digest)
		assert.Empty(t, digestTrees["n0"].GetLeafDigests())
	})

	t.Run("upsert digests", func(t *testing.T) {
		err = store.SetDigest("n0", expectedDigestTree)
		assert.NoError(t, err)

		digestTrees, err := store.GetDigests([]string{"n0"}, time.Now().Unix(), true)
		assert.NoError(t, err)
		assert.Contains(t, digestTrees, "n0")
		assert.True(t, proto.Equal(expectedDigestTree, digestTrees["n0"]))
	})

	t.Run("get outdated digests", func(t *testing.T) {
		clock.SetAndFreezeClock(t, clock.Now().Add(200*time.Second))
		err = store.SetDigest("n1", expectedDigestTree2)
		assert.NoError(t, err)

		digestTrees, err := store.GetDigests([]string{"n0", "n1"}, clock.Now().Unix(), true)
		assert.NoError(t, err)
		assert.Contains(t, digestTrees, "n0")
		assert.True(t, proto.Equal(expectedDigestTree, digestTrees["n0"]))
		assert.Contains(t, digestTrees, "n1")
		assert.True(t, proto.Equal(expectedDigestTree2, digestTrees["n1"]))

		digestTrees, err = store.GetDigests([]string{"n0", "n1"}, clock.Now().Unix()-100, true)
		assert.NoError(t, err)
		assert.Contains(t, digestTrees, "n0")
		assert.NotContains(t, digestTrees, "n1")

		digestTrees, err = store.GetDigests([]string{"n0", "n1"}, clock.Now().Unix()-300, true)
		assert.NoError(t, err)
		assert.NotContains(t, digestTrees, "n0")
		assert.NotContains(t, digestTrees, "n1")

		clock.UnfreezeClock(t)
	})

	t.Run("basic insert and get from cache", func(t *testing.T) {
		writer, err := store.UpdateCache("n0")
		assert.NoError(t, err)

		err = writer.InsertMany(objs1)
		assert.NoError(t, err)
		err = writer.InsertMany(objs2)
		assert.NoError(t, err)
		err = writer.Apply()
		assert.NoError(t, err)

		objs, err := store.GetCachedByID("n0", []string{"1", "2", "3"})
		assert.NoError(t, err)
		assert.Equal(t, expectedObjs, objs)

		expectedNextToken, err := configurator_storage.SerializePageToken(&configurator_storage.EntityPageToken{
			LastIncludedEntity: "3",
		})
		assert.NoError(t, err)
		objs, nextToken, err := store.GetCachedByPage("n0", "", 3)
		assert.NoError(t, err)
		assert.Equal(t, expectedObjs, objs)
		assert.Equal(t, expectedNextToken, nextToken)

		objs, nextToken, err = store.GetCachedByPage("n0", nextToken, 3)
		assert.NoError(t, err)
		assert.Empty(t, objs)
		assert.Empty(t, nextToken)

		// When the changes have been applied, the cache writer could no longer be used for insertions
		err = writer.InsertMany(objs1)
		assert.EqualError(t, err, "attempt to insert into network n0 with invalid cache writer")
	})

	t.Run("concurrent updates to cache", func(t *testing.T) {
		updateFunc := func(t *testing.T, objs map[string][]byte) {
			writer, err := store.UpdateCache("n0")
			assert.NoError(t, err)
			err = writer.InsertMany(objs)
			assert.NoError(t, err)
			err = writer.Apply()
			assert.NoError(t, err)
		}
		for i := 0; i <= 5; i++ {
			objs := map[string][]byte{
				strconv.Itoa(i):     []byte("apple"),
				strconv.Itoa(i + 1): []byte("banana"),
			}
			go updateFunc(t, objs)
		}
		time.Sleep(500 * time.Millisecond)

		// The store should only contain 2 objs
		page, token, err := store.GetCachedByPage("n0", "", 5)
		assert.NoError(t, err)
		assert.Empty(t, token)
		assert.Len(t, page, 2)

		// The two objs should have IDs i and i+1, and values "apple" and "banana"
		for i := 0; i <= 5; i++ {
			objs1, err := store.GetCachedByID("n0", []string{strconv.Itoa(i)})
			assert.NoError(t, err)
			if len(objs1) == 0 {
				continue
			}
			assert.Equal(t, []byte("apple"), objs1[0])
			objs2, err := store.GetCachedByID("n0", []string{strconv.Itoa(i + 1)})
			assert.NoError(t, err)
			assert.Len(t, objs2, 1)
			assert.Equal(t, []byte("banana"), objs2[0])
			break
		}
	})

	t.Run("last resync set and get", func(t *testing.T) {
		expectedLastResyncTime := time.Now().Unix()

		err := store.RecordResync("n0", "g0", expectedLastResyncTime+1)
		assert.NoError(t, err)
		err = store.RecordResync("n0", "g1", expectedLastResyncTime+2)
		assert.NoError(t, err)
		err = store.RecordResync("n1", "g0", expectedLastResyncTime+3)
		assert.NoError(t, err)
		err = store.RecordResync("n1", "g1", expectedLastResyncTime+4)
		assert.NoError(t, err)

		lastResyncTime1, err := store.GetLastResync("n0", "g0")
		assert.NoError(t, err)
		assert.Equal(t, expectedLastResyncTime+1, lastResyncTime1)
		lastResyncTime2, err := store.GetLastResync("n0", "g1")
		assert.NoError(t, err)
		assert.Equal(t, expectedLastResyncTime+2, lastResyncTime2)
		lastResyncTime3, err := store.GetLastResync("n1", "g0")
		assert.NoError(t, err)
		assert.Equal(t, expectedLastResyncTime+3, lastResyncTime3)
		lastResyncTime4, err := store.GetLastResync("n1", "g1")
		assert.NoError(t, err)
		assert.Equal(t, expectedLastResyncTime+4, lastResyncTime4)

		// Test upserting value to the store
		err = store.RecordResync("n0", "g0", expectedLastResyncTime+5)
		assert.NoError(t, err)
		lastResyncTime5, err := store.GetLastResync("n0", "g0")
		assert.NoError(t, err)
		assert.Equal(t, expectedLastResyncTime+5, lastResyncTime5)
	})

	t.Run("garbage collection", func(t *testing.T) {
		err = store.SetDigest("n0", expectedDigestTree)
		assert.NoError(t, err)
		err = store.SetDigest("n1", expectedDigestTree2)
		assert.NoError(t, err)
		writer, err := store.UpdateCache("n0")
		assert.NoError(t, err)
		err = writer.InsertMany(objs1)
		assert.NoError(t, err)
		err = writer.InsertMany(objs2)
		assert.NoError(t, err)
		err = writer.Apply()
		assert.NoError(t, err)
		err = store.RecordResync("n0", "g0", clock.Now().Unix())
		assert.NoError(t, err)
		err = store.RecordResync("n0", "g1", clock.Now().Unix())
		assert.NoError(t, err)
		err = store.RecordResync("n1", "g0", clock.Now().Unix())
		assert.NoError(t, err)

		digestTrees, err := store.GetDigests([]string{}, clock.Now().Unix(), true)
		assert.NoError(t, err)
		assert.Contains(t, digestTrees, "n0")
		assert.Contains(t, digestTrees, "n1")
		objs, _, err := store.GetCachedByPage("n0", "", 10)
		assert.NoError(t, err)
		assert.NotEmpty(t, objs)
		lastResync, err := store.GetLastResync("n0", "g0")
		assert.NoError(t, err)
		assert.NotEmpty(t, lastResync)
		lastResync, err = store.GetLastResync("n0", "g1")
		assert.NoError(t, err)
		assert.NotEmpty(t, lastResync)
		lastResync, err = store.GetLastResync("n1", "g0")
		assert.NoError(t, err)
		assert.NotEmpty(t, lastResync)

		// Only track data from network n1
		store.CollectGarbage([]string{"n1"})

		digestTrees, err = store.GetDigests([]string{}, clock.Now().Unix(), true)
		assert.NoError(t, err)
		assert.NotContains(t, digestTrees, "n0")
		assert.Contains(t, digestTrees, "n1")
		objs, _, err = store.GetCachedByPage("n0", "", 10)
		assert.NoError(t, err)
		assert.Empty(t, objs)
		lastResync, err = store.GetLastResync("n0", "g0")
		assert.NoError(t, err)
		assert.Empty(t, lastResync)
		lastResync, err = store.GetLastResync("n0", "g1")
		assert.NoError(t, err)
		assert.Empty(t, lastResync)
		lastResync, err = store.GetLastResync("n1", "g0")
		assert.NoError(t, err)
		assert.NotEmpty(t, lastResync)
	})

	t.Run("garbage collection with expired cacheWriters", func(t *testing.T) {
		// No actions can be performed with a cacheWriter that has already completed a write
		writer1, err := store.UpdateCache("n0")
		assert.NoError(t, err)
		err = writer1.InsertMany(objs1)
		assert.NoError(t, err)
		err = writer1.Apply()
		assert.NoError(t, err)
		err = writer1.InsertMany(objs1)
		assert.EqualError(t, err, "attempt to insert into network n0 with invalid cache writer")
		err = writer1.Apply()
		assert.EqualError(t, err, "attempt to apply updates to network n0 with invalid cache writer")

		writer2, err := store.UpdateCache("n0")
		assert.NoError(t, err)
		clock.SetAndFreezeClock(t, clock.Now().Add(300*time.Second))
		writer3, err := store.UpdateCache("n0")
		assert.NoError(t, err)
		store.CollectGarbage([]string{"n0"})

		// No actions can be performed with an expired cacheWriter
		err = writer2.InsertMany(objs1)
		assert.Error(t, err)
		err = writer3.InsertMany(objs1)
		assert.NoError(t, err)
	})
}
