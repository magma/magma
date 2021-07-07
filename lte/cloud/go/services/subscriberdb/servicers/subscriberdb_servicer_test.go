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
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_storage "magma/orc8r/cloud/go/services/configurator/storage"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestListSubscribers(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	digestStore := initializeDigestStore(t)

	servicer := servicers.NewSubscriberdbServicer(subscriberdb.Config{FlatDigestEnabled: false}, digestStore)
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
				Type: lte.SubscriberEntityType, Key: "IMSI12345",
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
			{Type: lte.SubscriberEntityType, Key: "IMSI67890", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
			{Type: lte.SubscriberEntityType, Key: "IMSI99999", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	// Fetch first page of subscribers
	expectedProtos := []*lte_protos.SubscriberData{
		{
			Sid: &lte_protos.SubscriberID{Id: "12345", Type: lte_protos.SubscriberID_IMSI},
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
			Sid:        &lte_protos.SubscriberID{Id: "67890", Type: lte_protos.SubscriberID_IMSI},
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
		LastIncludedEntity: "IMSI67890",
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
	assert.Equal(t, "", res.NextPageToken)

	// Create policies and base name associated to sub
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: lte.BaseNameEntityType, Key: "bn1",
				Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "IMSI12345"}},
			},
			{
				Type: lte.PolicyRuleEntityType, Key: "r1",
				Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "IMSI12345"}},
			},
			{
				Type: lte.PolicyRuleEntityType, Key: "r2",
				Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "IMSI12345"}},
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
	assert.Equal(t, 10, len(res.Subscribers))
	assert.Equal(t, expectedToken, res.NextPageToken)
}

func TestCheckSubscribersInSync(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	digestStore := initializeDigestStore(t)

	// Create servicer with flat digest feature flag turned on
	servicer := servicers.NewSubscriberdbServicer(subscriberdb.Config{FlatDigestEnabled: true}, digestStore)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)

	id := protos.NewGatewayIdentity("hw1", "n1", "g1")
	ctx := id.NewContextWithIdentity(context.Background())
	err = digestStore.SetDigest("n1", "digest_cherry")
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
		FlatDigest: &lte_protos.Digest{Md5Base64Digest: "digest_cherry"},
	}
	res, err = servicer.CheckSubscribersInSync(ctx, req)
	assert.NoError(t, err)
	assert.True(t, res.InSync)

	// Requests with outdated digests should get an update signal in return
	err = digestStore.SetDigest("n1", "digest_banana")
	assert.NoError(t, err)
	req = &lte_protos.CheckSubscribersInSyncRequest{
		FlatDigest: &lte_protos.Digest{Md5Base64Digest: "digest_cherry"},
	}
	res, err = servicer.CheckSubscribersInSync(ctx, req)
	assert.NoError(t, err)
	assert.False(t, res.InSync)
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

func initializeDigestStore(t *testing.T) subscriberdb_storage.DigestLookup {
	db, err := test_utils.GetSharedMemoryDB()
	assert.NoError(t, err)
	digestStore := subscriberdb_storage.NewDigestLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, digestStore.Initialize())
	return digestStore
}
