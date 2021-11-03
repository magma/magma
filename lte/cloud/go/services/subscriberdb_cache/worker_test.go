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

package subscriberdb_cache_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb_cache"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/mproto"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_storage "magma/orc8r/cloud/go/services/configurator/storage"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/sqorc"
	storage2 "magma/orc8r/cloud/go/storage"
	"magma/orc8r/cloud/go/syncstore"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

func TestSubscriberdbCacheWorker(t *testing.T) {
	store := initializeSyncstore(t)
	serviceConfig := subscriberdb_cache.Config{
		SleepIntervalSecs:  5,
		UpdateIntervalSecs: 300,
	}
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	allDigests, err := store.GetDigests([]string{}, clock.Now().Unix(), false)
	assert.NoError(t, err)
	assert.Empty(t, allDigests)
	digestTree, err := syncstore.GetDigestTree(store, "n1")
	assert.NoError(t, err)
	assert.Empty(t, digestTree.RootDigest)
	assert.Empty(t, digestTree.LeafDigests)
	subProtos, nextToken, err := store.GetCachedByPage("n1", "", 100)
	assert.NoError(t, err)
	assert.Empty(t, subProtos)
	assert.Empty(t, nextToken)

	err = configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	_, _, err = subscriberdb_cache.RenewDigests(serviceConfig, store)
	assert.NoError(t, err)
	digestTree, err = syncstore.GetDigestTree(store, "n1")
	assert.NoError(t, err)
	assert.NotEmpty(t, digestTree.GetRootDigest())
	assert.Empty(t, digestTree.GetLeafDigests())
	subProtos, nextToken, err = store.GetCachedByPage("n1", "", 100)
	assert.NoError(t, err)
	assert.Empty(t, subProtos)
	assert.Empty(t, nextToken)
	rootDigestExpected := digestTree.RootDigest.GetMd5Base64Digest()

	// Detect outdated digests and update
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{
			Type: lte.APNEntityType, Key: "apn1",
			Config: &lte_models.ApnConfiguration{},
		},
		{
			Type: lte.SubscriberEntityType, Key: "IMSI99999",
			Config: &models.SubscriberConfig{
				Lte:                   &models.LteSubscription{State: "ACTIVE"},
				ForbiddenNetworkTypes: models.CoreNetworkTypes{"5GC"},
			},
		},
		{
			Type: lte.SubscriberEntityType, Key: "IMSI11111",
			Config: &models.SubscriberConfig{
				Lte:                   &models.LteSubscription{State: "ACTIVE"},
				ForbiddenNetworkTypes: models.CoreNetworkTypes{"EPC"},
			},
		},
	}, serdes.Entity)
	assert.NoError(t, err)

	clock.SetAndFreezeClock(t, clock.Now().Add(10*time.Minute))
	_, _, err = subscriberdb_cache.RenewDigests(serviceConfig, store)
	assert.NoError(t, err)
	digestTree, err = syncstore.GetDigestTree(store, "n1")
	assert.NoError(t, err)
	assert.NotEqual(t, rootDigestExpected, digestTree.RootDigest.GetMd5Base64Digest())
	// The individual subscriber digests are ordered by subscriber ID, and are prefixed
	// by a hash of the subscriber data proto
	leafDigests := digestTree.LeafDigests
	assert.Len(t, leafDigests, 2)

	sub1 := &lte_protos.SubscriberData{
		Sid:        &lte_protos.SubscriberID{Id: "11111", Type: lte_protos.SubscriberID_IMSI},
		Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_ACTIVE, AuthKey: []byte{}},
		Non_3Gpp:   &lte_protos.Non3GPPUserProfile{ApnConfig: []*lte_protos.APNConfiguration{}},
		NetworkId:  &protos.NetworkID{Id: "n1"},
		SubProfile: "default",
		SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_EPC}},
	}
	expectedDigestPrefix1, err := mproto.HashDeterministic(sub1)
	assert.NoError(t, err)
	assert.Equal(t, "IMSI11111", leafDigests[0].Id)
	assert.NotEmpty(t, leafDigests[0].Digest.GetMd5Base64Digest())
	assert.True(t, strings.HasPrefix(leafDigests[0].Digest.GetMd5Base64Digest(), expectedDigestPrefix1))

	sub2 := &lte_protos.SubscriberData{
		Sid:        &lte_protos.SubscriberID{Id: "99999", Type: lte_protos.SubscriberID_IMSI},
		Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_ACTIVE, AuthKey: []byte{}},
		Non_3Gpp:   &lte_protos.Non3GPPUserProfile{ApnConfig: []*lte_protos.APNConfiguration{}},
		NetworkId:  &protos.NetworkID{Id: "n1"},
		SubProfile: "default",
		SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_5GC}},
	}
	expectedDigestPrefix2, err := mproto.HashDeterministic(sub2)
	assert.NoError(t, err)
	assert.Equal(t, "IMSI99999", leafDigests[1].Id)
	assert.NotEmpty(t, leafDigests[1].Digest.GetMd5Base64Digest())
	assert.True(t, strings.HasPrefix(leafDigests[1].Digest.GetMd5Base64Digest(), expectedDigestPrefix2))
	clock.UnfreezeClock(t)

	subProtos, nextToken, err = store.GetCachedByPage("n1", "", 2)
	assert.NoError(t, err)
	assert.Len(t, subProtos, 2)
	subProtosDeserialized, err := subscriberdb.DeserializeSubscribers(subProtos)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(sub1, subProtosDeserialized[0]))
	assert.True(t, proto.Equal(sub2, subProtosDeserialized[1]))
	expectedNextToken := getTokenByLastIncludedEntity(t, "IMSI99999")
	assert.Equal(t, expectedNextToken, nextToken)

	// Detect newly added and removed networks
	err = configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n2"}, serdes.Network)
	assert.NoError(t, err)
	configurator.DeleteNetwork(context.Background(), "n1")

	clock.SetAndFreezeClock(t, clock.Now().Add(20*time.Minute))
	_, _, err = subscriberdb_cache.RenewDigests(serviceConfig, store)
	assert.NoError(t, err)
	digestTree, err = syncstore.GetDigestTree(store, "n1")
	assert.NoError(t, err)
	assert.Empty(t, digestTree.RootDigest)
	assert.Empty(t, digestTree.LeafDigests)
	subProtos, nextToken, err = store.GetCachedByPage("n1", "", 100)
	assert.NoError(t, err)
	assert.Empty(t, subProtos)
	assert.Empty(t, nextToken)

	digestTree, err = syncstore.GetDigestTree(store, "n2")
	assert.NoError(t, err)
	assert.NotEmpty(t, digestTree.RootDigest)
	assert.Empty(t, digestTree.LeafDigests)
	subProtos, nextToken, err = store.GetCachedByPage("n2", "", 100)
	assert.NoError(t, err)
	assert.Empty(t, subProtos)
	assert.Empty(t, nextToken)

	allDigests, err = store.GetDigests([]string{}, clock.Now().Unix(), false)
	assert.NoError(t, err)
	assert.Equal(t, []string{"n2"}, allDigests.Networks())
	clock.UnfreezeClock(t)
}

// TestUpdateSubProtosByNetworkNoChange checks that, given there's no error in
// digest generation for the network, subscribers cache is only updated when
// the newly generated root digest is different from the previous root digest.
func TestUpdateSubProtosByNetworkNoChange(t *testing.T) {
	store := initializeSyncstore(t)
	serviceConfig := subscriberdb_cache.Config{
		SleepIntervalSecs:  5,
		UpdateIntervalSecs: 300,
	}
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{Type: lte.APNEntityType, Key: "apn1", Config: &lte_models.ApnConfiguration{}},
		{Type: lte.SubscriberEntityType, Key: "IMSI00001", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "ACTIVE"}, ForbiddenNetworkTypes: models.CoreNetworkTypes{"EPC", "5GC"}}},
		{Type: lte.SubscriberEntityType, Key: "IMSI00002", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "ACTIVE"}, ForbiddenNetworkTypes: models.CoreNetworkTypes{"EPC", "5GC"}}},
		{Type: lte.SubscriberEntityType, Key: "IMSI00003", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "ACTIVE"}, ForbiddenNetworkTypes: models.CoreNetworkTypes{"EPC", "5GC"}}},
	}, serdes.Entity)
	assert.NoError(t, err)

	_, _, err = subscriberdb_cache.RenewDigests(serviceConfig, store)
	assert.NoError(t, err)
	_, err = subscriberdb.GetDigest("n1")
	assert.NoError(t, err)
	page, _, err := store.GetCachedByPage("n1", "", 3)
	assert.NoError(t, err)
	assert.Len(t, page, 3)
	subProtos, err := subscriberdb.DeserializeSubscribers(page)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(subProtoFromID("IMSI00001"), subProtos[0]))
	assert.True(t, proto.Equal(subProtoFromID("IMSI00002"), subProtos[1]))
	assert.True(t, proto.Equal(subProtoFromID("IMSI00003"), subProtos[2]))

	// If the generated root digest matches the one in store, the update for cached subscribers wouldn't be triggered
	err = configurator.DeleteEntities(
		context.Background(),
		"n1",
		storage2.MakeTKs(lte.SubscriberEntityType, []string{"IMSI00001", "IMSI00002", "IMSI00003"}),
	)
	assert.NoError(t, err)
	newRootDigest, err := subscriberdb.GetDigest("n1")
	assert.NoError(t, err)
	newDigestTree := &protos.DigestTree{RootDigest: &protos.Digest{Md5Base64Digest: newRootDigest}}
	err = store.SetDigest("n1", newDigestTree)
	assert.NoError(t, err)

	clock.SetAndFreezeClock(t, clock.Now().Add(10*time.Minute))
	_, _, err = subscriberdb_cache.RenewDigests(serviceConfig, store)
	assert.NoError(t, err)
	page, _, err = store.GetCachedByPage("n1", "", 3)
	assert.NoError(t, err)
	assert.NotEmpty(t, page)
}

func getTokenByLastIncludedEntity(t *testing.T, sid string) string {
	token := &configurator_storage.EntityPageToken{
		LastIncludedEntity: sid,
	}
	encoded, err := configurator_storage.SerializePageToken(token)
	assert.NoError(t, err)
	return encoded
}

func subProtoFromID(sid string) *lte_protos.SubscriberData {
	subProto := &lte_protos.SubscriberData{
		Sid:        lte_protos.SidFromString(sid),
		Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_ACTIVE, AuthKey: []byte{}},
		Non_3Gpp:   &lte_protos.Non3GPPUserProfile{ApnConfig: []*lte_protos.APNConfiguration{}},
		NetworkId:  &protos.NetworkID{Id: "n1"},
		SubProfile: "default",
		SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_EPC, lte_protos.CoreNetworkType_NT_5GC}},
	}
	return subProto
}

func initializeSyncstore(t *testing.T) syncstore.SyncStore {
	db, err := test_utils.GetSharedMemoryDB()
	assert.NoError(t, err)
	fact := blobstore.NewSQLStoreFactory(subscriberdb.SyncstoreTableBlobstore, db, sqorc.GetSqlBuilder())
	assert.NoError(t, fact.InitializeFactory())
	store, err := syncstore.NewSyncStore(db, sqorc.GetSqlBuilder(), fact, syncstore.Config{TableNamePrefix: "subscriber", CacheWriterValidIntervalSecs: 150})
	assert.NoError(t, err)
	assert.NoError(t, store.Initialize())
	return store
}
