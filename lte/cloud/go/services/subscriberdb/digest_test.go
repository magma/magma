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

package subscriberdb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
)

func TestGetDigestDeterministic(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: "hw1"}, serdes.Entity)
	assert.NoError(t, err)
	gw, err := configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)

	// Create 1 APN and 5 pages of subscribers to test deterministic digest over multiple pages
	// Note: the page size is determined by the TestServiceMaxLoadSize of configurator
	subs := []configurator.NetworkEntity{}
	sub_count := 5 * configurator_test_init.TestServiceMaxPageSize
	for i := 0; i < sub_count; i++ {
		subs = append(subs, configurator.NetworkEntity{
			Type: lte.SubscriberEntityType, Key: fmt.Sprintf("IMSI000%02d", i),
			Config: &models.SubscriberConfig{
				Lte: &models.LteSubscription{State: "ACTIVE"},
			},
		})
	}
	networkEntities := append(subs, configurator.NetworkEntity{
		Type: lte.APNEntityType, Key: "apn",
		Config: &lte_models.ApnConfiguration{},
	})

	_, err = configurator.CreateEntities(context.Background(), "n1", networkEntities, serdes.Entity)
	assert.NoError(t, err)

	expected, err := subscriberdb.GetDigest("n1")
	assert.NoError(t, err)
	for i := 0; i < 50; i++ {
		digest, err := subscriberdb.GetDigest("n1")
		assert.NoError(t, err)
		assert.Equal(t, expected, digest)
	}

	// Update the subscriber list, the digest should update too
	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{
		Type: lte.SubscriberEntityType, Key: "IMSI11111",
		Config: &models.SubscriberConfig{
			Lte: &models.LteSubscription{State: "ACTIVE"},
		},
	}, serdes.Entity)
	assert.NoError(t, err)

	digest, err := subscriberdb.GetDigest("n1")
	assert.NoError(t, err)
	assert.NotEqual(t, expected, digest)
	expected = digest

	// Update the apn resources, the digest should update too
	var writes []configurator.EntityWriteOperation
	writes = append(writes, configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.APNResourceEntityType,
		Key:       "resource",
		Config: &lte_models.ApnResource{
			ApnName:    "apn",
			GatewayIP:  "172.16.254.1",
			GatewayMac: "00:0a:95:9d:68:16",
			ID:         "resource",
			VlanID:     42,
		},
		Associations: storage.TKs{{Type: lte.APNEntityType, Key: "apn"}},
	})
	write := configurator.EntityUpdateCriteria{
		Type:              lte.CellularGatewayEntityType,
		Key:               gw.Key,
		AssociationsToAdd: storage.TKs{{Type: lte.APNResourceEntityType, Key: "resource"}},
	}
	writes = append(writes, write)
	err = configurator.WriteEntities(context.Background(), "n1", writes, serdes.Entity)
	assert.NoError(t, err)

	digest, err = subscriberdb.GetDigest("n1")
	assert.NoError(t, err)
	assert.NotEqual(t, expected, digest)
}

// TestGetDigestApnResourceAssocs is a regression test to check whether the root
// digest reflects changes in the apn/gateway associations of apn resources.
func TestGetDigestApnResourceAssocs(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	gw1, err := configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)
	gw2, err := configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g2"}, serdes.Entity)
	assert.NoError(t, err)
	gw3, err := configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g3"}, serdes.Entity)
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{
			Type: lte.APNEntityType, Key: "apn1",
			Config: &lte_models.ApnConfiguration{},
		},
		{
			Type: lte.APNEntityType, Key: "apn2",
			Config: &lte_models.ApnConfiguration{},
		},
	}, serdes.Entity)
	assert.NoError(t, err)

	writes := []configurator.EntityWriteOperation{
		configurator.NetworkEntity{
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
			Associations: storage.TKs{
				{Type: lte.APNEntityType, Key: "apn1"},
			},
		},
		configurator.NetworkEntity{
			NetworkID: "n1",
			Type:      lte.APNResourceEntityType,
			Key:       "resource2",
			Config: &lte_models.ApnResource{
				ApnName:    "apn2",
				GatewayIP:  "172.16.254.2",
				GatewayMac: "00:0a:95:9d:68:16",
				ID:         "resource2",
				VlanID:     43,
			},
			Associations: storage.TKs{
				{Type: lte.APNEntityType, Key: "apn2"},
			},
		},
		configurator.EntityUpdateCriteria{
			Type: lte.CellularGatewayEntityType,
			Key:  gw1.Key,
			AssociationsToAdd: storage.TKs{
				{Type: lte.APNResourceEntityType, Key: "resource1"},
				{Type: lte.APNResourceEntityType, Key: "resource2"},
			},
		},
		configurator.EntityUpdateCriteria{
			Type: lte.CellularGatewayEntityType,
			Key:  gw2.Key,
			AssociationsToAdd: storage.TKs{
				{Type: lte.APNResourceEntityType, Key: "resource1"},
			},
		},
	}
	err = configurator.WriteEntities(context.Background(), "n1", writes, serdes.Entity)
	assert.NoError(t, err)
	expected, err := subscriberdb.GetDigest("n1")
	assert.NoError(t, err)

	// Digest reflects changes in gateway->apn resource associations
	writes = []configurator.EntityWriteOperation{
		configurator.EntityUpdateCriteria{
			Type: lte.CellularGatewayEntityType,
			Key:  gw2.Key,
			AssociationsToAdd: storage.TKs{
				{Type: lte.APNResourceEntityType, Key: "resource2"},
			},
		},
		configurator.EntityUpdateCriteria{
			Type: lte.CellularGatewayEntityType,
			Key:  gw3.Key,
			AssociationsToAdd: storage.TKs{
				{Type: lte.APNResourceEntityType, Key: "resource1"},
				{Type: lte.APNResourceEntityType, Key: "resource2"},
			},
		},
	}
	err = configurator.WriteEntities(context.Background(), "n1", writes, serdes.Entity)
	assert.NoError(t, err)

	digest, err := subscriberdb.GetDigest("n1")
	assert.NoError(t, err)
	assert.NotEqual(t, expected, digest)
	expected = digest

	// Digest reflects changes in apn resource->apn associations
	err = configurator.DeleteEntity(context.Background(), "n1", lte.APNResourceEntityType, "resource1")
	assert.NoError(t, err)
	writes = []configurator.EntityWriteOperation{
		configurator.NetworkEntity{
			NetworkID: "n1",
			Type:      lte.APNResourceEntityType,
			Key:       "resource1",
			Config: &lte_models.ApnResource{
				ApnName:    "apn2",
				GatewayIP:  "172.16.254.1",
				GatewayMac: "00:0a:95:9d:68:16",
				ID:         "resource1",
				VlanID:     42,
			},
			Associations: storage.TKs{
				{Type: lte.APNEntityType, Key: "apn2"},
			},
		},
	}
	err = configurator.WriteEntities(context.Background(), "n1", writes, serdes.Entity)
	assert.NoError(t, err)

	digest, err = subscriberdb.GetDigest("n1")
	assert.NoError(t, err)
	assert.NotEqual(t, expected, digest)
}

func TestGetPerSubscriberDigests(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	gw, err := configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)

	perSubDigests, err := subscriberdb.GetPerSubscriberDigests("n1")
	assert.NoError(t, err)
	// The generated leaf digests list should be empty but not nil, to avoid unexpected errors
	// when the agw client iterates over this list returned from cloud
	assert.Equal(t, []*protos.LeafDigest{}, perSubDigests)

	// Generate individual digests for each newly detected subscriber in the network
	_, err = configurator.CreateEntities(context.Background(), "n1", configurator.NetworkEntities{
		configurator.NetworkEntity{
			Type: lte.SubscriberEntityType, Key: "IMSI00001",
			Config: &models.SubscriberConfig{
				Lte: &models.LteSubscription{State: "ACTIVE"},
			},
		},
		configurator.NetworkEntity{
			Type: lte.SubscriberEntityType, Key: "IMSI00002",
			Config: &models.SubscriberConfig{
				Lte: &models.LteSubscription{State: "ACTIVE"},
			},
		},
		configurator.NetworkEntity{
			Type: lte.APNEntityType, Key: "apn",
			Config: &lte_models.ApnConfiguration{},
		},
	}, serdes.Entity)
	assert.NoError(t, err)

	perSubDigests, err = subscriberdb.GetPerSubscriberDigests("n1")
	assert.NoError(t, err)
	assert.Len(t, perSubDigests, 2)
	assert.Equal(t, "IMSI00001", perSubDigests[0].Id)
	assert.NotEmpty(t, perSubDigests[0].Digest.Md5Base64Digest)
	assert.Equal(t, "IMSI00002", perSubDigests[1].Id)
	assert.NotEmpty(t, perSubDigests[1].Digest.Md5Base64Digest)
	digestSub1 := perSubDigests[0].Digest.Md5Base64Digest

	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{
		Type: lte.SubscriberEntityType, Key: "IMSI00003",
		Config: &models.SubscriberConfig{
			Lte: &models.LteSubscription{State: "ACTIVE"},
		},
	}, serdes.Entity)
	assert.NoError(t, err)
	err = configurator.DeleteEntity(context.Background(), "n1", lte.SubscriberEntityType, "IMSI00001")
	assert.NoError(t, err)
	err = configurator.DeleteEntity(context.Background(), "n1", lte.SubscriberEntityType, "IMSI00002")
	assert.NoError(t, err)
	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{
		Type: lte.SubscriberEntityType, Key: "IMSI00001",
		Config: &models.SubscriberConfig{
			Lte: &models.LteSubscription{State: "INACTIVE"},
		},
	}, serdes.Entity)
	assert.NoError(t, err)

	// Detect changes in subscriber data and reflect them in the generated digests
	perSubDigests, err = subscriberdb.GetPerSubscriberDigests("n1")
	assert.NoError(t, err)
	assert.Len(t, perSubDigests, 2)
	assert.Equal(t, "IMSI00001", perSubDigests[0].Id)
	assert.NotEqual(t, digestSub1, perSubDigests[0].Digest.Md5Base64Digest)
	assert.Equal(t, "IMSI00003", perSubDigests[1].Id)
	assert.NotEmpty(t, perSubDigests[1].Digest.Md5Base64Digest)
	digestSub1 = perSubDigests[0].Digest.Md5Base64Digest
	digestSub3 := perSubDigests[1].Digest.Md5Base64Digest

	var writes []configurator.EntityWriteOperation
	writes = append(writes, configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.APNResourceEntityType,
		Key:       "resource",
		Config: &lte_models.ApnResource{
			ApnName:    "apn",
			GatewayIP:  "172.16.254.1",
			GatewayMac: "00:0a:95:9d:68:16",
			ID:         "resource",
			VlanID:     42,
		},
		Associations: storage.TKs{{Type: lte.APNEntityType, Key: "apn"}},
	})
	write := configurator.EntityUpdateCriteria{
		Type:              lte.CellularGatewayEntityType,
		Key:               gw.Key,
		AssociationsToAdd: storage.TKs{{Type: lte.APNResourceEntityType, Key: "resource"}},
	}
	writes = append(writes, write)
	err = configurator.WriteEntities(context.Background(), "n1", writes, serdes.Entity)
	assert.NoError(t, err)

	// Detect changes in apn resources data and reflect them in the ENTIRE set of generated digests
	perSubDigests, err = subscriberdb.GetPerSubscriberDigests("n1")
	assert.NoError(t, err)
	assert.Len(t, perSubDigests, 2)
	assert.Equal(t, "IMSI00001", perSubDigests[0].Id)
	assert.NotEqual(t, digestSub1, perSubDigests[0].Digest.Md5Base64Digest)
	assert.Equal(t, "IMSI00003", perSubDigests[1].Id)
	assert.NotEqual(t, digestSub3, perSubDigests[1].Digest.Md5Base64Digest)
}
