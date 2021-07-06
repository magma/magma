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
	"testing"
	"time"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/lte/cloud/go/services/subscriberdb_cache"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

func TestSubscriberdbCacheWorker(t *testing.T) {
	db, err := test_utils.GetSharedMemoryDB()
	assert.NoError(t, err)
	digestStore := storage.NewDigestLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, digestStore.Initialize())
	serviceConfig := subscriberdb_cache.Config{
		SleepIntervalSecs:  5,
		UpdateIntervalSecs: 300,
	}

	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	allNetworks, err := storage.GetAllNetworks(digestStore)
	assert.NoError(t, err)
	assert.Equal(t, []string{}, allNetworks)
	digest, _, err := storage.GetDigest(digestStore, "n1")
	assert.NoError(t, err)
	assert.Equal(t, "", digest)

	err = configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	subscriberdb_cache.RenewDigests(digestStore, serviceConfig)
	digest, _, err = storage.GetDigest(digestStore, "n1")
	assert.NoError(t, err)
	assert.NotEqual(t, "", digest)
	digestCanon := digest

	// Detect outdated digest and update
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
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	clock.SetAndFreezeClock(t, clock.Now().Add(10*time.Minute))
	subscriberdb_cache.RenewDigests(digestStore, serviceConfig)
	digest, _, err = storage.GetDigest(digestStore, "n1")
	assert.NoError(t, err)
	assert.NotEqual(t, digestCanon, digest)
	clock.UnfreezeClock(t)

	// Detect newly added and removed networks
	err = configurator.CreateNetwork(configurator.Network{ID: "n2"}, serdes.Network)
	assert.NoError(t, err)
	configurator.DeleteNetwork("n1")

	clock.SetAndFreezeClock(t, clock.Now().Add(20*time.Minute))
	subscriberdb_cache.RenewDigests(digestStore, serviceConfig)
	digest, _, err = storage.GetDigest(digestStore, "n1")
	assert.NoError(t, err)
	assert.Equal(t, "", digest)
	digest, _, err = storage.GetDigest(digestStore, "n2")
	assert.NoError(t, err)
	assert.NotEqual(t, "", digest)

	allNetworks, err = storage.GetAllNetworks(digestStore)
	assert.NoError(t, err)
	assert.Equal(t, []string{"n2"}, allNetworks)
	clock.UnfreezeClock(t)
}
