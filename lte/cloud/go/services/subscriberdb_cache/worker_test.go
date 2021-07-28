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
	"strings"
	"testing"
	"time"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/lte/cloud/go/services/subscriberdb_cache"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/mproto"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_storage "magma/orc8r/cloud/go/services/configurator/storage"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/sqorc"
	storage2 "magma/orc8r/cloud/go/storage"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestSubscriberdbCacheWorker(t *testing.T) {
	db, err := test_utils.GetSharedMemoryDB()
	assert.NoError(t, err)
	digestStore := storage.NewDigestStore(db, sqorc.GetSqlBuilder())
	assert.NoError(t, digestStore.Initialize())
	fact := blobstore.NewSQLBlobStorageFactory(subscriberdb.PerSubDigestTableBlobstore, db, sqorc.GetSqlBuilder())
	assert.NoError(t, fact.InitializeFactory())
	perSubDigestStore := storage.NewPerSubDigestStore(fact)
	serviceConfig := subscriberdb_cache.Config{
		SleepIntervalSecs:  5,
		UpdateIntervalSecs: 300,
	}
	subStore := storage.NewSubStore(db, sqorc.GetSqlBuilder())
	assert.NoError(t, subStore.Initialize())

	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	allNetworks, err := storage.GetAllNetworks(digestStore)
	assert.NoError(t, err)
	assert.Empty(t, allNetworks)
	digest, err := storage.GetDigest(digestStore, "n1")
	assert.NoError(t, err)
	assert.Empty(t, digest)
	perSubDigests, err := perSubDigestStore.GetDigest("n1")
	assert.NoError(t, err)
	assert.Empty(t, perSubDigests)
	subProtos, nextToken, err := subStore.GetSubscribersPage("n1", "", 100)
	assert.NoError(t, err)
	assert.Empty(t, subProtos)
	assert.Empty(t, nextToken)

	err = configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	_, _, err = subscriberdb_cache.RenewDigests(serviceConfig, digestStore, perSubDigestStore, subStore)
	assert.NoError(t, err)
	digest, err = storage.GetDigest(digestStore, "n1")
	assert.NoError(t, err)
	assert.NotEmpty(t, digest)
	perSubDigests, err = perSubDigestStore.GetDigest("n1")
	assert.NoError(t, err)
	assert.Empty(t, perSubDigests)
	subProtos, nextToken, err = subStore.GetSubscribersPage("n1", "", 100)
	assert.NoError(t, err)
	assert.Empty(t, subProtos)
	assert.Empty(t, nextToken)
	digestExpected := digest

	// Detect outdated digests and update
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: lte.APNEntityType, Key: "apn1",
				Config: &lte_models.ApnConfiguration{},
			},
			{
				Type: lte.SubscriberEntityType, Key: "IMSI99999",
				Config: &models.SubscriberConfig{
					Lte: &models.LteSubscription{State: "ACTIVE"},
				},
			},
			{
				Type: lte.SubscriberEntityType, Key: "IMSI11111",
				Config: &models.SubscriberConfig{
					Lte: &models.LteSubscription{State: "ACTIVE"},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	clock.SetAndFreezeClock(t, clock.Now().Add(10*time.Minute))
	_, _, err = subscriberdb_cache.RenewDigests(serviceConfig, digestStore, perSubDigestStore, subStore)
	assert.NoError(t, err)
	digest, err = storage.GetDigest(digestStore, "n1")
	assert.NoError(t, err)
	assert.NotEqual(t, digestExpected, digest)

	perSubDigests, err = perSubDigestStore.GetDigest("n1")
	assert.NoError(t, err)
	// The individual subscriber digests are ordered by subscriber ID, and are prefixed
	// by a hash of the subscriber data proto
	sub1 := &lte_protos.SubscriberData{
		Sid:        &lte_protos.SubscriberID{Id: "11111", Type: lte_protos.SubscriberID_IMSI},
		Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_ACTIVE, AuthKey: []byte{}},
		Non_3Gpp:   &lte_protos.Non3GPPUserProfile{ApnConfig: []*lte_protos.APNConfiguration{}},
		NetworkId:  &protos.NetworkID{Id: "n1"},
		SubProfile: "default",
	}
	expectedDigestPrefix1, err := mproto.HashDeterministic(sub1)
	assert.NoError(t, err)
	assert.Equal(t, "11111", perSubDigests[0].Sid.Id)
	assert.NotEmpty(t, perSubDigests[0].Digest.GetMd5Base64Digest())
	assert.True(t, strings.HasPrefix(perSubDigests[0].Digest.GetMd5Base64Digest(), expectedDigestPrefix1))

	sub2 := &lte_protos.SubscriberData{
		Sid:        &lte_protos.SubscriberID{Id: "99999", Type: lte_protos.SubscriberID_IMSI},
		Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_ACTIVE, AuthKey: []byte{}},
		Non_3Gpp:   &lte_protos.Non3GPPUserProfile{ApnConfig: []*lte_protos.APNConfiguration{}},
		NetworkId:  &protos.NetworkID{Id: "n1"},
		SubProfile: "default",
	}
	expectedDigestPrefix2, err := mproto.HashDeterministic(sub2)
	assert.NoError(t, err)
	assert.Equal(t, "99999", perSubDigests[1].Sid.Id)
	assert.NotEmpty(t, perSubDigests[1].Digest.GetMd5Base64Digest())
	assert.True(t, strings.HasPrefix(perSubDigests[1].Digest.GetMd5Base64Digest(), expectedDigestPrefix2))
	clock.UnfreezeClock(t)

	subProtos, nextToken, err = subStore.GetSubscribersPage("n1", "", 2)
	assert.NoError(t, err)
	assert.Len(t, subProtos, 2)
	assert.True(t, proto.Equal(sub1, subProtos[0]))
	assert.True(t, proto.Equal(sub2, subProtos[1]))
	expectedNextToken := getTokenByLastIncludedEntity(t, "IMSI99999")
	assert.Equal(t, expectedNextToken, nextToken)

	// Detect newly added and removed networks
	err = configurator.CreateNetwork(configurator.Network{ID: "n2"}, serdes.Network)
	assert.NoError(t, err)
	configurator.DeleteNetwork("n1")

	clock.SetAndFreezeClock(t, clock.Now().Add(20*time.Minute))
	_, _, err = subscriberdb_cache.RenewDigests(serviceConfig, digestStore, perSubDigestStore, subStore)
	assert.NoError(t, err)
	digest, err = storage.GetDigest(digestStore, "n1")
	assert.NoError(t, err)
	assert.Empty(t, digest)
	perSubDigests, err = perSubDigestStore.GetDigest("n1")
	assert.NoError(t, err)
	assert.Empty(t, perSubDigests)
	subProtos, nextToken, err = subStore.GetSubscribersPage("n1", "", 100)
	assert.NoError(t, err)
	assert.Empty(t, subProtos)
	assert.Empty(t, nextToken)

	digest, err = storage.GetDigest(digestStore, "n2")
	assert.NoError(t, err)
	assert.NotEmpty(t, digest)
	perSubDigests, err = perSubDigestStore.GetDigest("n2")
	assert.NoError(t, err)
	assert.Empty(t, perSubDigests)
	subProtos, nextToken, err = subStore.GetSubscribersPage("n2", "", 100)
	assert.NoError(t, err)
	assert.Empty(t, subProtos)
	assert.Empty(t, nextToken)

	allNetworks, err = storage.GetAllNetworks(digestStore)
	assert.NoError(t, err)
	assert.Equal(t, []string{"n2"}, allNetworks)
	clock.UnfreezeClock(t)
}

// TestUpdateSubProtosByNetworkNoChange checks that, given there's no error in
// digest generation for the network, subscribers cache is only updated when
// the newly generated flat digest is different from the previous digest.
func TestUpdateSubProtosByNetworkNoChange(t *testing.T) {
	db, err := test_utils.GetSharedMemoryDB()
	assert.NoError(t, err)
	digestStore := storage.NewDigestStore(db, sqorc.GetSqlBuilder())
	assert.NoError(t, digestStore.Initialize())
	fact := blobstore.NewSQLBlobStorageFactory(subscriberdb.PerSubDigestTableBlobstore, db, sqorc.GetSqlBuilder())
	assert.NoError(t, fact.InitializeFactory())
	perSubDigestStore := storage.NewPerSubDigestStore(fact)
	serviceConfig := subscriberdb_cache.Config{
		SleepIntervalSecs:  5,
		UpdateIntervalSecs: 300,
	}
	subStore := storage.NewSubStore(db, sqorc.GetSqlBuilder())
	assert.NoError(t, subStore.Initialize())

	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	err = configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.APNEntityType, Key: "apn1", Config: &lte_models.ApnConfiguration{}},
			{Type: lte.SubscriberEntityType, Key: "IMSI00001", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "ACTIVE"}}},
			{Type: lte.SubscriberEntityType, Key: "IMSI00002", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "ACTIVE"}}},
			{Type: lte.SubscriberEntityType, Key: "IMSI00003", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "ACTIVE"}}},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	_, _, err = subscriberdb_cache.RenewDigests(serviceConfig, digestStore, perSubDigestStore, subStore)
	assert.NoError(t, err)
	_, err = subscriberdb.GetDigest("n1")
	assert.NoError(t, err)
	page, _, err := subStore.GetSubscribersPage("n1", "", 3)
	assert.NoError(t, err)
	assert.Len(t, page, 3)
	assert.True(t, proto.Equal(subProtoFromID("IMSI00001"), page[0]))
	assert.True(t, proto.Equal(subProtoFromID("IMSI00002"), page[1]))
	assert.True(t, proto.Equal(subProtoFromID("IMSI00003"), page[2]))

	// If the generated flat digest matches the one in store, the update for subStore wouldn't be triggered
	err = configurator.DeleteEntities(
		"n1",
		storage2.MakeTKs(lte.SubscriberEntityType, []string{"IMSI00001", "IMSI00002", "IMSI00003"}),
	)
	assert.NoError(t, err)
	newDigest, err := subscriberdb.GetDigest("n1")
	assert.NoError(t, err)
	err = digestStore.SetDigest("n1", newDigest)
	assert.NoError(t, err)

	clock.SetAndFreezeClock(t, clock.Now().Add(10*time.Minute))
	_, _, err = subscriberdb_cache.RenewDigests(serviceConfig, digestStore, perSubDigestStore, subStore)
	assert.NoError(t, err)
	page, _, err = subStore.GetSubscribersPage("n1", "", 3)
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
	}
	return subProto
}
