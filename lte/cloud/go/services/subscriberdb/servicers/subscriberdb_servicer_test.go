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
	"time"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb/servicers"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_storage "magma/orc8r/cloud/go/services/configurator/storage"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	orc8r_models "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/cloud/go/syncstore"
	"magma/orc8r/lib/go/protos"
)

func TestListSuciProfiles(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	storeReader, _ := initializeStore(t)

	servicer := servicers.NewSubscriberdbServicer(subscriberdb.Config{DigestsEnabled: false}, storeReader)

	err := configurator.CreateNetwork(context.Background(), configurator.Network{
		ID:          "nt1",
		Type:        lte.NetworkType,
		Name:        "foobar",
		Description: "Foo Bar",
		Configs: map[string]interface{}{
			lte.CellularNetworkConfigType: &lte_models.NetworkCellularConfigs{
				Ran: &lte_models.NetworkRanConfigs{
					BandwidthMhz: 20,
					TddConfig: &lte_models.NetworkRanConfigsTddConfig{
						Earfcndl:               44590,
						SubframeAssignment:     2,
						SpecialSubframePattern: 7,
					},
				},
				Epc: &lte_models.NetworkEpcConfigs{
					Mcc: "001",
					Mnc: "01",
					Tac: 1,
					// 16 bytes of \x11
					LteAuthOp:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					LteAuthAmf: []byte("\x80\x00"),

					HssRelayEnabled:          swag.Bool(false),
					GxGyRelayEnabled:         swag.Bool(false),
					CloudSubscriberdbEnabled: false,
					CongestionControlEnabled: swag.Bool(true),
					Enable5gFeatures:         swag.Bool(false),
					NodeIdentifier:           "",
					DefaultRuleID:            "",
					SubscriberdbSyncInterval: lte_models.SubscriberdbSyncInterval(300),
				},
				Ngc: &lte_models.NetworkNgcConfigs{SuciProfiles: []*lte_models.SuciProfile{
					{
						HomeNetworkPublicKey:           []byte("\x12\x12\x12\x12"),
						HomeNetworkPrivateKey:          []byte("\x12\x12\x12\x12"),
						HomeNetworkPublicKeyIdentifier: 255,
						ProtectionScheme:               "ProfileA",
					},
				}},
			},
			orc8r.NetworkFeaturesConfig: orc8r_models.NewDefaultFeaturesConfig(),
			orc8r.DnsdNetworkType:       orc8r_models.NewDefaultDNSConfig(),
		},
	}, serdes.Network)
	assert.NoError(t, err)

	id := protos.NewGatewayIdentity("hw1", "nt1", "g1")
	ctx := id.NewContextWithIdentity(context.Background())

	expectedProtos := []*lte_protos.SuciProfile{
		{
			HomeNetPublicKeyId: 255,
			HomeNetPublicKey:   []byte("\x12\x12\x12\x12"),
			HomeNetPrivateKey:  []byte("\x12\x12\x12\x12"),
			ProtectionScheme:   lte_protos.SuciProfile_ProfileA,
		},
	}

	req := &protos.Void{}
	res, err := servicer.ListSuciProfiles(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, expectedProtos, res.SuciProfiles)
}

func TestListSubscribers(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	storeReader, _ := initializeStore(t)

	servicer := servicers.NewSubscriberdbServicer(subscriberdb.Config{DigestsEnabled: false}, storeReader)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: "hw1"}, serdes.Entity)
	assert.NoError(t, err)
	gw, err := configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)

	id := protos.NewGatewayIdentity("hw1", "n1", "g1")
	ctx := id.NewContextWithIdentity(context.Background())

	// 2 subs without a profile on the backend (should fill as "default"), the
	// other inactive with a sub profile
	// 2 APNs active for the active sub, 1 with an assigned static IP and the
	// other without
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
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
				ForbiddenNetworkTypes: models.CoreNetworkTypes{"5GC"},
				StaticIps:             models.SubscriberStaticIps{"apn1": "192.168.100.1"},
			},
			Associations: storage.TKs{{Type: lte.APNEntityType, Key: "apn1"}, {Type: lte.APNEntityType, Key: "apn2"}},
		},
		{Type: lte.SubscriberEntityType, Key: "IMSI00002", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}, ForbiddenNetworkTypes: models.CoreNetworkTypes{"EPC"}}},
		{Type: lte.SubscriberEntityType, Key: "IMSI99999", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}, ForbiddenNetworkTypes: models.CoreNetworkTypes{"5GC"}}},
	}, serdes.Entity)
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
			SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_5GC}},
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
			SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_EPC}},
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
			SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_5GC}},
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
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{
			Type: lte.BaseNameEntityType, Key: "bn1",
			Associations: storage.TKs{{Type: lte.SubscriberEntityType, Key: "IMSI00001"}},
		},
		{
			Type: lte.PolicyRuleEntityType, Key: "r1",
			Associations: storage.TKs{{Type: lte.SubscriberEntityType, Key: "IMSI00001"}},
		},
		{
			Type: lte.PolicyRuleEntityType, Key: "r2",
			Associations: storage.TKs{{Type: lte.SubscriberEntityType, Key: "IMSI00001"}},
		},
	}, serdes.Entity)
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
	err = configurator.WriteEntities(context.Background(), "n1", writes, serdes.Entity)
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
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{Type: lte.SubscriberEntityType, Key: "IMSI99991", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
		{Type: lte.SubscriberEntityType, Key: "IMSI99992", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
		{Type: lte.SubscriberEntityType, Key: "IMSI99993", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
		{Type: lte.SubscriberEntityType, Key: "IMSI99994", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
		{Type: lte.SubscriberEntityType, Key: "IMSI99995", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
		{Type: lte.SubscriberEntityType, Key: "IMSI99996", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
		{Type: lte.SubscriberEntityType, Key: "IMSI99997", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
		{Type: lte.SubscriberEntityType, Key: "IMSI99998", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
	}, serdes.Entity)
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
	// store (with write access) is created only to insert test data into the database
	storeReader, store := initializeStore(t)

	servicer := servicers.NewSubscriberdbServicer(subscriberdb.Config{
		DigestsEnabled:     true,
		MaxProtosLoadSize:  10,
		ResyncIntervalSecs: 1000,
	}, storeReader)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	gw, err := configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
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
			SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_5GC}},
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
			SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_EPC}},
		},
	}
	expectedProtosSerialized, err := subscriberdb.SerializeSubscribers(expectedProtos)
	assert.NoError(t, err)
	writer, err := store.UpdateCache("n1")
	assert.NoError(t, err)
	err = writer.InsertMany(expectedProtosSerialized)
	assert.NoError(t, err)
	err = writer.Apply()
	assert.NoError(t, err)

	// Root and leaf digests in the cloud store should be returned as well
	expectedDigestTree := &protos.DigestTree{
		RootDigest: &protos.Digest{Md5Base64Digest: "cherry"},
		LeafDigests: []*protos.LeafDigest{
			{Id: "IMSI00001", Digest: &protos.Digest{Md5Base64Digest: "apple"}},
			{Id: "IMSI00002", Digest: &protos.Digest{Md5Base64Digest: "banana"}},
			{Id: "IMSI99999", Digest: &protos.Digest{Md5Base64Digest: "cherry"}},
		},
	}
	err = store.SetDigest("n1", expectedDigestTree)
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
	assert.True(t, proto.Equal(expectedDigestTree, res.Digests))

	// The servicer should append gateway-specific apn resources data to returned subscriber protos
	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{
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
	err = configurator.WriteEntities(context.Background(), "n1", writes, serdes.Entity)
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
	allProtosSerialized, err := subscriberdb.SerializeSubscribers(allProtos)
	assert.NoError(t, err)
	writer, err = store.UpdateCache("n1")
	assert.NoError(t, err)
	err = writer.InsertMany(allProtosSerialized)
	assert.NoError(t, err)
	err = writer.Apply()
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

func TestListSubscribersSetLastResyncTime(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	storeReader, store := initializeStore(t)

	servicer := servicers.NewSubscriberdbServicer(subscriberdb.Config{
		DigestsEnabled:    true,
		MaxProtosLoadSize: 10,
	}, storeReader)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)

	id := protos.NewGatewayIdentity("hw1", "n1", "g1")
	ctx := id.NewContextWithIdentity(context.Background())

	expectedProtos := []*lte_protos.SubscriberData{
		basicSubProtoFromSid("IMSI00001", ""),
		basicSubProtoFromSid("IMSI00002", ""),
		basicSubProtoFromSid("IMSI00003", ""),
	}
	expectedProtosSerialized, err := subscriberdb.SerializeSubscribers(expectedProtos)
	assert.NoError(t, err)
	writer, err := store.UpdateCache("n1")
	assert.NoError(t, err)
	err = writer.InsertMany(expectedProtosSerialized)
	assert.NoError(t, err)
	err = writer.Apply()
	assert.NoError(t, err)

	// The last resync time for this AGW should be set on the request for the last page (when nextToken is empty)
	expectedNextToken := serializeToken(t, &configurator_storage.EntityPageToken{
		LastIncludedEntity: "IMSI00003",
	})
	req := &lte_protos.ListSubscribersRequest{
		PageSize:  3,
		PageToken: "",
	}
	res, err := servicer.ListSubscribers(ctx, req)
	assert.NoError(t, err)
	assert.Len(t, res.Subscribers, 3)
	assert.Equal(t, expectedNextToken, res.NextPageToken)
	lastResyncTime, err := store.GetLastResync("n1", "g1")
	assert.NoError(t, err)
	assert.Empty(t, lastResyncTime)

	req = &lte_protos.ListSubscribersRequest{
		PageSize:  3,
		PageToken: expectedNextToken,
	}
	res, err = servicer.ListSubscribers(ctx, req)
	assert.NoError(t, err)
	assert.Len(t, res.Subscribers, 0)
	assert.Empty(t, res.NextPageToken)
	lastResyncTime, err = store.GetLastResync("n1", "g1")
	assert.NoError(t, err)
	assert.NotEmpty(t, lastResyncTime)
}

func TestCheckSubscribersInSync(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	storeReader, store := initializeStore(t)

	servicer := servicers.NewSubscriberdbServicer(subscriberdb.Config{DigestsEnabled: true, ResyncIntervalSecs: 1000}, storeReader)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)

	id := protos.NewGatewayIdentity("hw1", "n1", "g1")
	ctx := id.NewContextWithIdentity(context.Background())
	err = store.SetDigest("n1", &protos.DigestTree{RootDigest: &protos.Digest{Md5Base64Digest: "digest_apple"}})
	assert.NoError(t, err)

	err = store.RecordResync("n1", "g1", time.Now().Unix())
	assert.NoError(t, err)
	// Requests with blank digests should get an update signal in return
	req := &lte_protos.CheckInSyncRequest{
		RootDigest: &protos.Digest{Md5Base64Digest: ""},
	}
	res, err := servicer.CheckInSync(ctx, req)
	assert.NoError(t, err)
	assert.False(t, res.InSync)

	// Requests with up-to-date digests should get a no-update signal in return
	req = &lte_protos.CheckInSyncRequest{
		RootDigest: &protos.Digest{Md5Base64Digest: "digest_apple"},
	}
	res, err = servicer.CheckInSync(ctx, req)
	assert.NoError(t, err)
	assert.True(t, res.InSync)

	// Requests with outdated digests should get an update signal in return
	err = store.SetDigest("n1", &protos.DigestTree{RootDigest: &protos.Digest{Md5Base64Digest: "digest_apple2"}})
	assert.NoError(t, err)
	req = &lte_protos.CheckInSyncRequest{
		RootDigest: &protos.Digest{Md5Base64Digest: "digest_apple"},
	}
	res, err = servicer.CheckInSync(ctx, req)
	assert.NoError(t, err)
	assert.False(t, res.InSync)

	// Requests from gateways that haven't been resynced for more than the specified
	// resync interval should get an update signal in return
	err = store.RecordResync("n1", "g1", time.Now().Unix()-5000)
	assert.NoError(t, err)
	req = &lte_protos.CheckInSyncRequest{
		RootDigest: &protos.Digest{Md5Base64Digest: "digest_apple2"},
	}
	res, err = servicer.CheckInSync(ctx, req)
	assert.NoError(t, err)
	assert.False(t, res.InSync)
}

func TestSyncSubscribers(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	storeReader, store := initializeStore(t)

	// Create servicer with the subscriber digests feature flag turned on
	configs := subscriberdb.Config{DigestsEnabled: true, ChangesetSizeThreshold: 100, MaxProtosLoadSize: 100}
	servicer := servicers.NewSubscriberdbServicer(configs, storeReader)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)
	id := protos.NewGatewayIdentity("hw1", "n1", "g1")
	ctx := id.NewContextWithIdentity(context.Background())
	err = store.SetDigest("n1", &protos.DigestTree{RootDigest: &protos.Digest{Md5Base64Digest: "root_digest_apple"}})
	assert.NoError(t, err)

	// If the gateway has not received orc8r-oriented resync in a while, should get a resync signal
	req := &lte_protos.SyncRequest{
		LeafDigests: []*protos.LeafDigest{},
	}
	res, err := servicer.Sync(ctx, req)
	assert.NoError(t, err)
	assert.True(t, res.Resync)

	err = store.RecordResync("n1", "g1", time.Now().Unix())
	assert.NoError(t, err)
	// Initially no digests
	req = &lte_protos.SyncRequest{
		LeafDigests: []*protos.LeafDigest{},
	}
	res, err = servicer.Sync(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "root_digest_apple", res.Digests.RootDigest.Md5Base64Digest)
	assert.Empty(t, res.Digests.LeafDigests)
	assert.Empty(t, res.Changeset.ToRenew)
	assert.Empty(t, res.Changeset.Deleted)

	// When cloud has updated leaf digests in store, changeset is sent back
	expectedToRenewData := []*lte_protos.SubscriberData{
		{
			Sid:        &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "profile_apple",
			SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_5GC}},
		},
		{
			Sid:        &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "profile_banana",
			SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_EPC}},
		},
	}
	expectedToRenewDataMarshaled := []*any.Any{}
	for _, data := range expectedToRenewData {
		val, err := ptypes.MarshalAny(data)
		assert.NoError(t, err)
		expectedToRenewDataMarshaled = append(expectedToRenewDataMarshaled, val)
	}
	expectedToRenewDataSerialized, err := subscriberdb.SerializeSubscribers(expectedToRenewData)
	assert.NoError(t, err)
	writer, err := store.UpdateCache("n1")
	assert.NoError(t, err)
	err = writer.InsertMany(expectedToRenewDataSerialized)
	assert.NoError(t, err)
	err = writer.Apply()
	assert.NoError(t, err)

	expectedDigests := &protos.DigestTree{
		RootDigest: &protos.Digest{Md5Base64Digest: "root_digest_apple"},
		LeafDigests: []*protos.LeafDigest{
			{Id: "IMSI00001", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_apple"}},
			{Id: "IMSI00002", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_banana"}},
		},
	}
	err = store.SetDigest("n1", expectedDigests)
	assert.NoError(t, err)

	req = &lte_protos.SyncRequest{
		LeafDigests: []*protos.LeafDigest{},
	}
	res, err = servicer.Sync(ctx, req)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(expectedDigests, res.Digests))
	assertEqualAnyData(t, expectedToRenewDataMarshaled, res.Changeset.ToRenew)
	assert.Empty(t, res.Changeset.Deleted)

	// Test deleting and updating the subscriber data in store
	curDigests := expectedDigests
	store.CollectGarbage([]string{"n1"})
	expectedToRenewData = []*lte_protos.SubscriberData{
		{
			Sid:        &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "profile_banana2",
			SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_EPC}},
		},
		{
			Sid:        &lte_protos.SubscriberID{Id: "00003", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "profile_cherry",
			SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_5GC}},
		},
	}
	expectedToRenewDataMarshaled = []*any.Any{}
	for _, data := range expectedToRenewData {
		val, err := ptypes.MarshalAny(data)
		assert.NoError(t, err)
		expectedToRenewDataMarshaled = append(expectedToRenewDataMarshaled, val)
	}
	expectedToRenewDataSerialized, err = subscriberdb.SerializeSubscribers(expectedToRenewData)
	assert.NoError(t, err)
	writer, err = store.UpdateCache("n1")
	assert.NoError(t, err)
	err = writer.InsertMany(expectedToRenewDataSerialized)
	assert.NoError(t, err)
	err = writer.Apply()
	assert.NoError(t, err)

	expectedDigests = &protos.DigestTree{
		RootDigest: &protos.Digest{Md5Base64Digest: "root_digest_apple2"},
		LeafDigests: []*protos.LeafDigest{
			{Id: "IMSI00002", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_banana2"}},
			{Id: "IMSI00003", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_cherry"}},
		},
	}

	err = store.SetDigest("n1", expectedDigests)
	assert.NoError(t, err)

	req = &lte_protos.SyncRequest{
		LeafDigests: curDigests.LeafDigests,
	}
	res, err = servicer.Sync(ctx, req)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(expectedDigests, res.Digests))
	assertEqualAnyData(t, expectedToRenewDataMarshaled, res.Changeset.ToRenew)
	assert.Equal(t, []string{"IMSI00001"}, res.Changeset.Deleted)
}

func TestSyncSubscribersResync(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	storeReader, store := initializeStore(t)

	// Create servicer with a small ChangesetSizeThreshold
	configs := subscriberdb.Config{
		DigestsEnabled:         true,
		ChangesetSizeThreshold: 2,
		MaxProtosLoadSize:      100,
		ResyncIntervalSecs:     1000,
	}
	servicer := servicers.NewSubscriberdbServicer(configs, storeReader)

	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)
	id := protos.NewGatewayIdentity("hw1", "n1", "g1")
	ctx := id.NewContextWithIdentity(context.Background())

	err = store.RecordResync("n1", "g1", time.Now().Unix())
	assert.NoError(t, err)
	// When changeset is no larger than ChangesetSizeThreshold, the servicer should return the full changeset
	expectedToRenewData := []*lte_protos.SubscriberData{
		{
			Sid:        &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "profile_apple",
			SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_5GC}},
		},
		{
			Sid:        &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE, AuthKey: []byte{}},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "profile_banana",
			SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_EPC}},
		},
	}
	expectedToRenewDataMarshaled := []*any.Any{}
	for _, data := range expectedToRenewData {
		val, err := ptypes.MarshalAny(data)
		assert.NoError(t, err)
		expectedToRenewDataMarshaled = append(expectedToRenewDataMarshaled, val)
	}
	expectedToRenewDataSerialized, err := subscriberdb.SerializeSubscribers(expectedToRenewData)
	assert.NoError(t, err)
	writer, err := store.UpdateCache("n1")
	assert.NoError(t, err)
	err = writer.InsertMany(expectedToRenewDataSerialized)
	assert.NoError(t, err)
	err = writer.Apply()
	assert.NoError(t, err)

	expectedDigests := &protos.DigestTree{
		RootDigest: &protos.Digest{Md5Base64Digest: ""},
		LeafDigests: []*protos.LeafDigest{
			{Id: "IMSI00001", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_apple"}},
			{Id: "IMSI00002", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_banana"}},
		},
	}
	err = store.SetDigest("n1", expectedDigests)
	assert.NoError(t, err)

	req := &lte_protos.SyncRequest{
		LeafDigests: []*protos.LeafDigest{},
	}
	res, err := servicer.Sync(ctx, req)
	assert.NoError(t, err)
	assert.False(t, res.Resync)
	assert.True(t, proto.Equal(expectedDigests, res.Digests))
	assertEqualAnyData(t, expectedToRenewDataMarshaled, res.Changeset.ToRenew)
	assert.Empty(t, res.Changeset.Deleted)

	// When the changeset is larger than ChangesetSizeThreshold, the servicer should return resync and nothing else
	curDigests := expectedDigests
	err = store.SetDigest("n1", &protos.DigestTree{
		LeafDigests: []*protos.LeafDigest{
			{Id: "IMSI00003", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_cherry"}},
			{Id: "IMSI00004", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_dragonfruit"}},
			{Id: "IMSI00005", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_eggplant"}},
		},
	})
	assert.NoError(t, err)

	req = &lte_protos.SyncRequest{
		LeafDigests: curDigests.LeafDigests,
	}
	res, err = servicer.Sync(ctx, req)
	assert.NoError(t, err)
	assert.True(t, res.Resync)
	assert.Empty(t, res.Digests)
	assert.Empty(t, res.Changeset)
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

func assertEqualAnyData(t *testing.T, expected []*any.Any, got []*any.Any) {
	assert.Equal(t, len(expected), len(got))
	for i := 0; i < len(expected); i++ {
		// HACK: using the following workaround because in our current version of protobuf,
		// proto.Equal can't be used on two objects of type *any.Any
		assert.Equal(t, expected[i].TypeUrl, got[i].TypeUrl)
		assert.Equal(t, expected[i].Value, got[i].Value)
	}
}

func initializeStore(t *testing.T) (syncstore.SyncStoreReader, syncstore.SyncStore) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewSQLStoreFactory(subscriberdb.SyncstoreTableBlobstore, db, sqorc.GetSqlBuilder())
	assert.NoError(t, fact.InitializeFactory())
	store, err := syncstore.NewSyncStore(db, sqorc.GetSqlBuilder(), fact, syncstore.Config{TableNamePrefix: "subscriber", CacheWriterValidIntervalSecs: 150})
	assert.NoError(t, err)
	assert.NoError(t, store.Initialize())
	storeReader, err := syncstore.NewSyncStoreReader(db, sqorc.GetSqlBuilder(), fact, syncstore.Config{TableNamePrefix: "subscriber"})
	assert.NoError(t, err)
	assert.NoError(t, store.Initialize())
	return storeReader, store
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
		SubNetwork: &lte_protos.CoreNetworkType{ForbiddenNetworkTypes: []lte_protos.CoreNetworkType_CoreNetworkTypes{lte_protos.CoreNetworkType_NT_5GC, lte_protos.CoreNetworkType_NT_EPC}},
	}
	return subProto
}
