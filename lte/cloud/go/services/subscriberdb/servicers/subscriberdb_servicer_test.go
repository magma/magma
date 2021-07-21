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

package servicers_test

import (
	"context"
	"encoding/base64"
	"testing"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb/servicers"
	subscriberdb_storage "magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_storage "magma/orc8r/cloud/go/services/configurator/storage"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestListSubscribers(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	digestStore := initializeDigestStore(t)
	perSubDigestStore := initializePerSubDigestStore(t)
	subStore := initializeSubStore(t)

	servicer := servicers.NewSubscriberdbServicer(subscriberdb.Config{DigestsEnabled: false}, digestStore, perSubDigestStore, subStore)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: "hw1"}, serdes.Entity)
	assert.NoError(t, err)
	gw, err := configurator.CreateEntity("n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)

	id := protos.NewGatewayIdentity("hw1", "n1", "g1")
	ctx := id.NewContextWithIdentity(context.Background())

	// 2 subs without a profile on the backend (should fill as "default"), the
	// other inactive with a sub profile
	// 2 APNs active for the active sub, 1 with an assigned static IP and the
	// other without
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: lte.APNEntityType, Key: "apn1",
				Config: &lte_models.ApnConfiguration{
					Ambr: &lte_models.AggregatedMaximumBitrate{
						MaxBandwidthDl: swag.Uint32(42),
						MaxBandwidthUl: swag.Uint32(100),
					},
					QosProfile: &lte_models.QosProfile{
						ClassID:                 swag.Int32(1),
						PreemptionCapability:    swag.Bool(true),
						PreemptionVulnerability: swag.Bool(true),
						PriorityLevel:           swag.Uint32(1),
					},
				},
			},
			{
				Type: lte.APNEntityType, Key: "apn2",
				Config: &lte_models.ApnConfiguration{
					Ambr: &lte_models.AggregatedMaximumBitrate{
						MaxBandwidthDl: swag.Uint32(42),
						MaxBandwidthUl: swag.Uint32(100),
					},
					QosProfile: &lte_models.QosProfile{
						ClassID:                 swag.Int32(2),
						PreemptionCapability:    swag.Bool(false),
						PreemptionVulnerability: swag.Bool(false),
						PriorityLevel:           swag.Uint32(2),
					},
				},
			},
			{
				Type: lte.SubscriberEntityType, Key: "IMSI00001",
				Config: &models.SubscriberConfig{
					Lte: &models.LteSubscription{
						State:   "ACTIVE",
						AuthKey: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
						AuthOpc: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					},
					StaticIps: models.SubscriberStaticIps{"apn1": "192.168.100.1"},
				},
				Associations: []storage.TypeAndKey{{Type: lte.APNEntityType, Key: "apn1"}, {Type: lte.APNEntityType, Key: "apn2"}},
			},
			{Type: lte.SubscriberEntityType, Key: "IMSI00002", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
			{Type: lte.SubscriberEntityType, Key: "IMSI99999", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	// Fetch first page of subscribers
	expectedProtos := []*lte_protos.SubscriberData{
		{
			Sid: &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
			Lte: &lte_protos.LTESubscription{
				State:   lte_protos.LTESubscription_ACTIVE,
				AuthKey: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				AuthOpc: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "default",
			Non_3Gpp: &lte_protos.Non3GPPUserProfile{
				ApnConfig: []*lte_protos.APNConfiguration{
					{
						ServiceSelection: "apn1",
						QosProfile: &lte_protos.APNConfiguration_QoSProfile{
							ClassId:                 1,
							PriorityLevel:           1,
							PreemptionCapability:    true,
							PreemptionVulnerability: true,
						},
						Ambr: &lte_protos.AggregatedMaximumBitrate{
							MaxBandwidthUl: 100,
							MaxBandwidthDl: 42,
						},
						AssignedStaticIp: "192.168.100.1",
					},
					{
						ServiceSelection: "apn2",
						QosProfile: &lte_protos.APNConfiguration_QoSProfile{
							ClassId:                 2,
							PriorityLevel:           2,
							PreemptionCapability:    false,
							PreemptionVulnerability: false,
						},
						Ambr: &lte_protos.AggregatedMaximumBitrate{
							MaxBandwidthUl: 100,
							MaxBandwidthDl: 42,
						},
					},
				},
			},
		},
		{
			Sid:        &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{ApnConfig: []*lte_protos.APNConfiguration{}},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "foo",
		},
	}

	// Fetch first page of subscribers
	req := &lte_protos.ListSubscribersRequest{
		PageSize:  2,
		PageToken: "",
	}
	res, err := servicer.ListSubscribers(ctx, req)
	token := &configurator_storage.EntityPageToken{
		LastIncludedEntity: "IMSI00002",
	}
	expectedToken := serializeToken(t, token)
	assert.NoError(t, err)
	assertEqualSubscriberData(t, expectedProtos, res.Subscribers)
	assert.Equal(t, expectedToken, res.NextPageToken)

	expectedProtos2 := []*lte_protos.SubscriberData{
		{
			Sid:        &lte_protos.SubscriberID{Id: "99999", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{ApnConfig: []*lte_protos.APNConfiguration{}},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "foo",
		},
	}

	// Fetch final page of subscribers
	req = &lte_protos.ListSubscribersRequest{
		PageSize:  2,
		PageToken: res.NextPageToken,
	}
	res, err = servicer.ListSubscribers(ctx, req)
	assert.NoError(t, err)
	assertEqualSubscriberData(t, expectedProtos2, res.Subscribers)
	assert.Empty(t, res.NextPageToken)

	// Create policies and base name associated to sub
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: lte.BaseNameEntityType, Key: "bn1",
				Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "IMSI00001"}},
			},
			{
				Type: lte.PolicyRuleEntityType, Key: "r1",
				Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "IMSI00001"}},
			},
			{
				Type: lte.PolicyRuleEntityType, Key: "r2",
				Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "IMSI00001"}},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	expectedProtos[0].Lte.AssignedPolicies = []string{"r1", "r2"}
	expectedProtos[0].Lte.AssignedBaseNames = []string{"bn1"}

	req = &lte_protos.ListSubscribersRequest{
		PageSize:  2,
		PageToken: "",
	}
	res, err = servicer.ListSubscribers(ctx, req)
	assert.NoError(t, err)
	assertEqualSubscriberData(t, expectedProtos, res.Subscribers)
	assert.Equal(t, expectedToken, res.NextPageToken)

	// Create gateway-specific APN configuration
	var writes []configurator.EntityWriteOperation
	writes = append(writes, configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.APNResourceEntityType,
		Key:       "resource1",
		Config: &lte_models.ApnResource{
			ApnName:    "apn1",
			GatewayIP:  "172.16.254.1",
			GatewayMac: "00:0a:95:9d:68:16",
			ID:         "resource1",
			VlanID:     42,
		},
		Associations: storage.TKs{{Type: lte.APNEntityType, Key: "apn1"}},
	})
	writes = append(writes, configurator.EntityUpdateCriteria{
		Type:              lte.CellularGatewayEntityType,
		Key:               gw.Key,
		AssociationsToAdd: storage.TKs{{Type: lte.APNResourceEntityType, Key: "resource1"}},
	})
	err = configurator.WriteEntities("n1", writes, serdes.Entity)
	assert.NoError(t, err)

	expectedProtos[0].Non_3Gpp.ApnConfig[0].Resource = &lte_protos.APNConfiguration_APNResource{
		ApnName:    "apn1",
		GatewayIp:  "172.16.254.1",
		GatewayMac: "00:0a:95:9d:68:16",
		VlanId:     42,
	}

	res, err = servicer.ListSubscribers(ctx, req)
	assert.NoError(t, err)
	assertEqualSubscriberData(t, expectedProtos, res.Subscribers)
	assert.Equal(t, expectedToken, res.NextPageToken)

	// Create 8 more subscribers to test max page size
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.SubscriberEntityType, Key: "IMSI99991", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
			{Type: lte.SubscriberEntityType, Key: "IMSI99992", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
			{Type: lte.SubscriberEntityType, Key: "IMSI99993", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
			{Type: lte.SubscriberEntityType, Key: "IMSI99994", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
			{Type: lte.SubscriberEntityType, Key: "IMSI99995", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
			{Type: lte.SubscriberEntityType, Key: "IMSI99996", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
			{Type: lte.SubscriberEntityType, Key: "IMSI99997", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
			{Type: lte.SubscriberEntityType, Key: "IMSI99998", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	// max page size for the configurator test service is 10 entities
	// Ensure when page size specified is 0, max size is returned (10/11 subs)
	req = &lte_protos.ListSubscribersRequest{
		PageSize:  0,
		PageToken: "",
	}
	res, err = servicer.ListSubscribers(ctx, req)
	token = &configurator_storage.EntityPageToken{
		LastIncludedEntity: "IMSI99998",
	}
	expectedToken = serializeToken(t, token)
	assert.NoError(t, err)
	assert.Len(t, res.Subscribers, 10)
	assert.Equal(t, expectedToken, res.NextPageToken)
}

func TestListSubscribersDigestsEnabled(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	digestStore := initializeDigestStore(t)
	perSubDigestStore := initializePerSubDigestStore(t)
	subStore := initializeSubStore(t)

	servicer := servicers.NewSubscriberdbServicer(subscriberdb.Config{
		DigestsEnabled:    true,
		MaxProtosLoadSize: 10,
	}, digestStore, perSubDigestStore, subStore)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	gw, err := configurator.CreateEntity("n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)

	id := protos.NewGatewayIdentity("hw1", "n1", "g1")
	ctx := id.NewContextWithIdentity(context.Background())

	// The subscriberdb servicer should return subscriber protos read from cache
	expectedProtos := []*lte_protos.SubscriberData{
		{
			Sid: &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
			Lte: &lte_protos.LTESubscription{
				State:   lte_protos.LTESubscription_ACTIVE,
				AuthKey: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				AuthOpc: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "default",
			Non_3Gpp: &lte_protos.Non3GPPUserProfile{
				ApnConfig: []*lte_protos.APNConfiguration{
					{
						ServiceSelection: "apn1",
						QosProfile: &lte_protos.APNConfiguration_QoSProfile{
							ClassId:                 1,
							PriorityLevel:           1,
							PreemptionCapability:    true,
							PreemptionVulnerability: true,
						},
						Ambr: &lte_protos.AggregatedMaximumBitrate{
							MaxBandwidthUl: 100,
							MaxBandwidthDl: 42,
						},
						AssignedStaticIp: "192.168.100.1",
					},
					{
						ServiceSelection: "apn2",
						QosProfile: &lte_protos.APNConfiguration_QoSProfile{
							ClassId:                 2,
							PriorityLevel:           2,
							PreemptionCapability:    false,
							PreemptionVulnerability: false,
						},
						Ambr: &lte_protos.AggregatedMaximumBitrate{
							MaxBandwidthUl: 100,
							MaxBandwidthDl: 42,
						},
					},
				},
			},
		},
		{
			Sid:        &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{ApnConfig: []*lte_protos.APNConfiguration{}},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "foo",
		},
	}
	err = subStore.InsertMany("n1", expectedProtos)
	assert.NoError(t, err)
	err = subStore.ApplyUpdate("n1")
	assert.NoError(t, err)

	// Flat and per-sub digests in the cloud store should be returned as well
	expectedDigest := "cherry"
	expectedPerSubDigests := []*lte_protos.SubscriberDigestWithID{
		{
			Sid:    &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
			Digest: &lte_protos.Digest{Md5Base64Digest: "apple"},
		},
		{
			Sid:    &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
			Digest: &lte_protos.Digest{Md5Base64Digest: "banana"},
		},
		{
			Sid:    &lte_protos.SubscriberID{Id: "99999", Type: lte_protos.SubscriberID_IMSI},
			Digest: &lte_protos.Digest{Md5Base64Digest: "cherry"},
		},
	}
	err = digestStore.SetDigest("n1", expectedDigest)
	assert.NoError(t, err)
	err = perSubDigestStore.SetDigest("n1", expectedPerSubDigests)
	assert.NoError(t, err)

	req := &lte_protos.ListSubscribersRequest{
		PageSize:  2,
		PageToken: "",
	}
	res, err := servicer.ListSubscribers(ctx, req)
	token := &configurator_storage.EntityPageToken{
		LastIncludedEntity: "IMSI00002",
	}
	expectedToken := serializeToken(t, token)
	assert.NoError(t, err)
	assertEqualSubscriberData(t, expectedProtos, res.Subscribers)
	assert.Equal(t, expectedToken, res.NextPageToken)
	assert.Equal(t, expectedDigest, res.FlatDigest.GetMd5Base64Digest())
	assertEqualPerSubDigests(t, expectedPerSubDigests, res.PerSubDigests)

	// The servicer should append gateway-specific apn resources data to returned subscriber protos
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{
		Type: lte.APNEntityType, Key: "apn1",
		Config: &lte_models.ApnConfiguration{},
	}, serdes.Entity)
	assert.NoError(t, err)

	var writes []configurator.EntityWriteOperation
	writes = append(writes, configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.APNResourceEntityType,
		Key:       "resource1",
		Config: &lte_models.ApnResource{
			ApnName:    "apn1",
			GatewayIP:  "172.16.254.1",
			GatewayMac: "00:0a:95:9d:68:16",
			ID:         "resource1",
			VlanID:     42,
		},
		Associations: storage.TKs{{Type: lte.APNEntityType, Key: "apn1"}},
	})
	writes = append(writes, configurator.EntityUpdateCriteria{
		Type:              lte.CellularGatewayEntityType,
		Key:               gw.Key,
		AssociationsToAdd: storage.TKs{{Type: lte.APNResourceEntityType, Key: "resource1"}},
	})
	err = configurator.WriteEntities("n1", writes, serdes.Entity)
	assert.NoError(t, err)

	expectedProtos[0].Non_3Gpp.ApnConfig[0].Resource = &lte_protos.APNConfiguration_APNResource{
		ApnName:    "apn1",
		GatewayIp:  "172.16.254.1",
		GatewayMac: "00:0a:95:9d:68:16",
		VlanId:     42,
	}
	res, err = servicer.ListSubscribers(ctx, req)
	assert.NoError(t, err)
	assertEqualSubscriberData(t, expectedProtos, res.Subscribers)
	assert.Equal(t, expectedToken, res.NextPageToken)

	// Create 8 more subscribers in cache to test max page size
	allProtos := []*lte_protos.SubscriberData{
		basicSubProtoFromSid("IMSI00001", ""), basicSubProtoFromSid("IMSI00002", ""),
		basicSubProtoFromSid("IMSI99991", ""), basicSubProtoFromSid("IMSI99992", ""),
		basicSubProtoFromSid("IMSI99993", ""), basicSubProtoFromSid("IMSI99994", ""),
		basicSubProtoFromSid("IMSI99995", ""), basicSubProtoFromSid("IMSI99996", ""),
		basicSubProtoFromSid("IMSI99997", ""), basicSubProtoFromSid("IMSI99998", ""),
	}
	err = subStore.InsertMany("n1", allProtos)
	assert.NoError(t, err)
	err = subStore.ApplyUpdate("n1")
	assert.NoError(t, err)

	// Ensure when page size specified is 0, max page size is returned (10/11 subs)
	req = &lte_protos.ListSubscribersRequest{
		PageSize:  0,
		PageToken: "",
	}
	res, err = servicer.ListSubscribers(ctx, req)
	token = &configurator_storage.EntityPageToken{
		LastIncludedEntity: "IMSI99998",
	}
	expectedToken = serializeToken(t, token)
	assert.NoError(t, err)
	assert.Len(t, res.Subscribers, 10)
	assert.Equal(t, expectedToken, res.NextPageToken)
}

func TestCheckSubscribersInSync(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	digestStore := initializeDigestStore(t)
	perSubDigestStore := initializePerSubDigestStore(t)
	subStore := initializeSubStore(t)

	servicer := servicers.NewSubscriberdbServicer(subscriberdb.Config{DigestsEnabled: true}, digestStore, perSubDigestStore, subStore)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)

	id := protos.NewGatewayIdentity("hw1", "n1", "g1")
	ctx := id.NewContextWithIdentity(context.Background())
	err = digestStore.SetDigest("n1", "digest_apple")
	assert.NoError(t, err)

	// Requests with blank digests should get an update signal in return
	req := &lte_protos.CheckSubscribersInSyncRequest{
		FlatDigest: &lte_protos.Digest{Md5Base64Digest: ""},
	}
	res, err := servicer.CheckSubscribersInSync(ctx, req)
	assert.NoError(t, err)
	assert.False(t, res.InSync)

	// Requests with up-to-date digests should get a no-update signal in return
	req = &lte_protos.CheckSubscribersInSyncRequest{
		FlatDigest: &lte_protos.Digest{Md5Base64Digest: "digest_apple"},
	}
	res, err = servicer.CheckSubscribersInSync(ctx, req)
	assert.NoError(t, err)
	assert.True(t, res.InSync)

	// Requests with outdated digests should get an update signal in return
	err = digestStore.SetDigest("n1", "digest_apple2")
	assert.NoError(t, err)
	req = &lte_protos.CheckSubscribersInSyncRequest{
		FlatDigest: &lte_protos.Digest{Md5Base64Digest: "digest_apple"},
	}
	res, err = servicer.CheckSubscribersInSync(ctx, req)
	assert.NoError(t, err)
	assert.False(t, res.InSync)
}

func TestSyncSubscribers(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	digestStore := initializeDigestStore(t)
	perSubDigestStore := initializePerSubDigestStore(t)
	subStore := initializeSubStore(t)

	// Create servicer with the subscriber digests feature flag turned on
	configs := subscriberdb.Config{DigestsEnabled: true, ChangesetSizeThreshold: 100, MaxProtosLoadSize: 100}
	servicer := servicers.NewSubscriberdbServicer(configs, digestStore, perSubDigestStore, subStore)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)
	id := protos.NewGatewayIdentity("hw1", "n1", "g1")
	ctx := id.NewContextWithIdentity(context.Background())
	err = digestStore.SetDigest("n1", "flat_digest_apple")
	assert.NoError(t, err)

	// Initially no digests
	req := &lte_protos.SyncSubscribersRequest{
		PerSubDigests: []*lte_protos.SubscriberDigestWithID{},
	}
	res, err := servicer.SyncSubscribers(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "flat_digest_apple", res.FlatDigest.GetMd5Base64Digest())
	assert.Empty(t, res.PerSubDigests)
	assert.Empty(t, res.Deleted)
	assert.Empty(t, res.ToRenew)

	// When cloud has updated per sub digests in store, changeset is sent back
	expectedToRenewData := []*lte_protos.SubscriberData{
		{
			Sid:        &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "profile_apple",
		},
		{
			Sid:        &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "profile_banana",
		},
	}
	err = subStore.InsertMany("n1", expectedToRenewData)
	assert.NoError(t, err)
	err = subStore.ApplyUpdate("n1")
	assert.NoError(t, err)

	expectedPerSubDigests := []*lte_protos.SubscriberDigestWithID{
		{
			Sid:    &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
			Digest: &lte_protos.Digest{Md5Base64Digest: "digest_apple"},
		},
		{
			Sid:    &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
			Digest: &lte_protos.Digest{Md5Base64Digest: "digest_banana"},
		},
	}
	err = perSubDigestStore.SetDigest("n1", expectedPerSubDigests)
	assert.NoError(t, err)

	req = &lte_protos.SyncSubscribersRequest{
		PerSubDigests: []*lte_protos.SubscriberDigestWithID{},
	}
	res, err = servicer.SyncSubscribers(ctx, req)
	assert.NoError(t, err)
	assertEqualPerSubDigests(t, expectedPerSubDigests, res.PerSubDigests)
	assertEqualSubscriberData(t, expectedToRenewData, res.ToRenew)
	assert.Empty(t, res.Deleted)

	// Test deleting and updating the subscriber data in store
	curPerSubDigests := expectedPerSubDigests
	err = subStore.DeleteSubscribersForNetworks([]string{"n1"})
	assert.NoError(t, err)
	expectedToRenewData = []*lte_protos.SubscriberData{
		{
			Sid:        &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "profile_banana2",
		},
		{
			Sid:        &lte_protos.SubscriberID{Id: "00003", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "profile_cherry",
		},
	}
	err = subStore.InsertMany("n1", expectedToRenewData)
	assert.NoError(t, err)
	err = subStore.ApplyUpdate("n1")
	assert.NoError(t, err)

	expectedPerSubDigests = []*lte_protos.SubscriberDigestWithID{
		{
			Sid:    &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
			Digest: &lte_protos.Digest{Md5Base64Digest: "digest_banana2"},
		},
		{
			Sid:    &lte_protos.SubscriberID{Id: "00003", Type: lte_protos.SubscriberID_IMSI},
			Digest: &lte_protos.Digest{Md5Base64Digest: "digest_cherry"},
		},
	}
	err = perSubDigestStore.SetDigest("n1", expectedPerSubDigests)
	assert.NoError(t, err)

	req = &lte_protos.SyncSubscribersRequest{
		PerSubDigests: curPerSubDigests,
	}
	res, err = servicer.SyncSubscribers(ctx, req)
	assert.NoError(t, err)
	assertEqualPerSubDigests(t, expectedPerSubDigests, res.PerSubDigests)
	assertEqualSubscriberData(t, expectedToRenewData, res.ToRenew)
	assert.Equal(t, []string{"IMSI00001"}, res.Deleted)
}

func TestSyncSubscribersResync(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	digestStore := initializeDigestStore(t)
	perSubDigestStore := initializePerSubDigestStore(t)
	subStore := initializeSubStore(t)

	// Create servicer with a small ChangesetSizeThreshold
	configs := subscriberdb.Config{
		DigestsEnabled:         true,
		ChangesetSizeThreshold: 2,
		MaxProtosLoadSize:      100,
	}
	servicer := servicers.NewSubscriberdbServicer(configs, digestStore, perSubDigestStore, subStore)

	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)
	id := protos.NewGatewayIdentity("hw1", "n1", "g1")
	ctx := id.NewContextWithIdentity(context.Background())

	// When changeset is no larger than ChangesetSizeThreshold, the servicer should return the full changeset
	expectedToRenewData := []*lte_protos.SubscriberData{
		{
			Sid:        &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "profile_apple",
		},
		{
			Sid:        &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "profile_banana",
		},
	}
	err = subStore.InsertMany("n1", expectedToRenewData)
	assert.NoError(t, err)
	err = subStore.ApplyUpdate("n1")
	assert.NoError(t, err)

	expectedPerSubDigests := []*lte_protos.SubscriberDigestWithID{
		{
			Sid:    &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
			Digest: &lte_protos.Digest{Md5Base64Digest: "digest_apple"},
		},
		{
			Sid:    &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
			Digest: &lte_protos.Digest{Md5Base64Digest: "digest_banana"},
		},
	}
	err = perSubDigestStore.SetDigest("n1", expectedPerSubDigests)
	assert.NoError(t, err)

	req := &lte_protos.SyncSubscribersRequest{
		PerSubDigests: []*lte_protos.SubscriberDigestWithID{},
	}
	res, err := servicer.SyncSubscribers(ctx, req)
	assert.NoError(t, err)
	assert.False(t, res.Resync)
	assertEqualPerSubDigests(t, expectedPerSubDigests, res.PerSubDigests)
	assertEqualSubscriberData(t, expectedToRenewData, res.ToRenew)
	assert.Empty(t, res.Deleted)

	// When the changeset is larger than ChangesetSizeThreshold, the servicer should return resync and nothing else
	curPerSubDigests := expectedPerSubDigests
	err = perSubDigestStore.SetDigest("n1", []*lte_protos.SubscriberDigestWithID{
		{
			Sid:    &lte_protos.SubscriberID{Id: "00003", Type: lte_protos.SubscriberID_IMSI},
			Digest: &lte_protos.Digest{Md5Base64Digest: "cherry"},
		},
		{
			Sid:    &lte_protos.SubscriberID{Id: "00004", Type: lte_protos.SubscriberID_IMSI},
			Digest: &lte_protos.Digest{Md5Base64Digest: "dragonfruit"},
		},
		{
			Sid:    &lte_protos.SubscriberID{Id: "00005", Type: lte_protos.SubscriberID_IMSI},
			Digest: &lte_protos.Digest{Md5Base64Digest: "eggplant"},
		},
	})
	assert.NoError(t, err)

	req = &lte_protos.SyncSubscribersRequest{
		PerSubDigests: curPerSubDigests,
	}
	res, err = servicer.SyncSubscribers(ctx, req)
	assert.NoError(t, err)
	assert.True(t, res.Resync)
	assert.Empty(t, res.PerSubDigests)
	assert.Empty(t, res.ToRenew)
	assert.Empty(t, res.Deleted)
}

func serializeToken(t *testing.T, token *configurator_storage.EntityPageToken) string {
	marshalledToken, err := proto.Marshal(token)
	assert.NoError(t, err)
	return base64.StdEncoding.EncodeToString(marshalledToken)
}

func assertEqualSubscriberData(t *testing.T, expectedProtos []*lte_protos.SubscriberData, actualProtos []*lte_protos.SubscriberData) {
	assert.True(t, len(expectedProtos) == len(actualProtos))
	for i := 0; i < len(expectedProtos); i++ {
		assert.True(t, proto.Equal(expectedProtos[i], actualProtos[i]))
	}
}

func initializeDigestStore(t *testing.T) subscriberdb_storage.DigestStore {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	store := subscriberdb_storage.NewDigestStore(db, sqorc.GetSqlBuilder())
	assert.NoError(t, store.Initialize())
	return store
}

func initializePerSubDigestStore(t *testing.T) *subscriberdb_storage.PerSubDigestStore {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewSQLBlobStorageFactory(subscriberdb.PerSubDigestTableBlobstore, db, sqorc.GetSqlBuilder())
	assert.NoError(t, fact.InitializeFactory())
	store := subscriberdb_storage.NewPerSubDigestStore(fact)
	return store
}

func initializeSubStore(t *testing.T) *subscriberdb_storage.SubStore {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	store := subscriberdb_storage.NewSubStore(db, sqorc.GetSqlBuilder())
	assert.NoError(t, store.Initialize())
	return store
}
func assertEqualPerSubDigests(t *testing.T, expected []*lte_protos.SubscriberDigestWithID, got []*lte_protos.SubscriberDigestWithID) {
	assert.Equal(t, len(expected), len(got))
	for ind := range expected {
		assert.Equal(t, expected[ind].Digest.GetMd5Base64Digest(), got[ind].Digest.GetMd5Base64Digest())
		assert.Equal(t, expected[ind].Sid.Id, got[ind].Sid.Id)
	}
}

func basicSubProtoFromSid(sid string, subProfile string) *lte_protos.SubscriberData {
	if subProfile == "" {
		subProfile = "foo"
	}
	subProto := &lte_protos.SubscriberData{
		Sid:        lte_protos.SidFromString(sid),
		Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
		Non_3Gpp:   &lte_protos.Non3GPPUserProfile{ApnConfig: []*lte_protos.APNConfiguration{}},
		NetworkId:  &protos.NetworkID{Id: "n1"},
		SubProfile: subProfile,
	}
	return subProto
}
