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
	"testing"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	storage2 "magma/orc8r/cloud/go/storage"

	"github.com/stretchr/testify/assert"
)

// TestConvertSubEntsToProtos is a regression test to check if ConvertSubEntsToProtos
// properly filters for the apn associations of a subscriber.
func TestConvertSubEntsToProtos(t *testing.T) {
	apnConfigs := map[string]*lte_models.ApnConfiguration{
		"apple": {},
	}
	apnResources := lte_models.ApnResources{
		"apple": lte_models.ApnResource{
			ID:      "birch",
			ApnName: "apple",
		},
	}
	subscriber := configurator.NetworkEntity{
		NetworkID:    "n1",
		Key:          "IMSI00000",
		Config:       &models.SubscriberConfig{Lte: &models.LteSubscription{State: "ACTIVE"}},
		Associations: storage2.TKs{{Type: lte.APNResourceEntityType, Key: "apple"}},
	}

	// The resultant subProto should have a blank Non_3Gpp field, since there's no
	// apn associations to this subscriber
	expectedSubProto := &lte_protos.SubscriberData{
		Sid:        &lte_protos.SubscriberID{Id: "00000", Type: lte_protos.SubscriberID_IMSI},
		Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_ACTIVE},
		SubProfile: "default",
		Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
	}
	subProto, err := subscriberdb.ConvertSubEntsToProtos(subscriber, apnConfigs, apnResources)
	assert.NoError(t, err)
	assert.Equal(t, expectedSubProto, subProto)
}
