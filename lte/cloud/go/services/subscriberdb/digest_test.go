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
	"fmt"
	"testing"

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

	"github.com/stretchr/testify/assert"
)

const (
	iterationCount = 50
	// subscriberCount is the number of test subscribers created to test deterministic streaming of subscribers.
	// You can test for deterministic digest generation over multiple pages of subscriber data by reducing
	// the maximum page load size for configurator to less than this value.
	subscriberCount = 50
)

func TestGetDigestDeterministic(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: "hw1"}, serdes.Entity)
	assert.NoError(t, err)
	gw, err := configurator.CreateEntity("n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)

	// Create 1 APN and <subscriberCount> subscribers
	subs := []configurator.NetworkEntity{}
	for i := 0; i < subscriberCount; i++ {
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

	_, err = configurator.CreateEntities("n1", networkEntities, serdes.Entity)
	assert.NoError(t, err)

	expected, err := subscriberdb.GetDigest("n1")
	assert.NoError(t, err)
	for i := 0; i < iterationCount; i++ {
		digest, err := subscriberdb.GetDigest("n1")
		assert.NoError(t, err)
		assert.Equal(t, expected, digest)
	}

	// Update the subscriber list, the digest should update too
	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: lte.SubscriberEntityType, Key: "IMSI11111",
			Config: &models.SubscriberConfig{
				Lte: &models.LteSubscription{State: "ACTIVE"},
			},
		},
		serdes.Entity,
	)
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
	writes = append(writes, configurator.EntityUpdateCriteria{
		Type:              lte.CellularGatewayEntityType,
		Key:               gw.Key,
		AssociationsToAdd: storage.TKs{{Type: lte.APNResourceEntityType, Key: "resource"}},
	})
	err = configurator.WriteEntities("n1", writes, serdes.Entity)
	assert.NoError(t, err)

	digest, err = subscriberdb.GetDigest("n1")
	assert.NoError(t, err)
	assert.NotEqual(t, expected, digest)
}
