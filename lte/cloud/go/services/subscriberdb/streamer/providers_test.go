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

package streamer_test

import (
	"testing"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
)

func TestSubscriberdbStreamer(t *testing.T) {
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	provider, err := providers.GetStreamProvider(lte.SubscriberStreamName)
	assert.NoError(t, err)

	err = configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: "hw1"}, serdes.Entity)
	assert.NoError(t, err)
	gw, err := configurator.CreateEntity("n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)

	// 1 sub without a profile on the backend (should fill as "default"), the
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
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

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
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE},
			Non_3Gpp:   &lte_protos.Non3GPPUserProfile{},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "foo",
		},
	}
	expected := funk.Map(
		expectedProtos,
		func(sub *lte_protos.SubscriberData) *protos.DataUpdate {
			data, err := proto.Marshal(sub)
			assert.NoError(t, err)
			return &protos.DataUpdate{Key: "IMSI" + sub.Sid.Id, Value: data}
		},
	)
	actual, err := provider.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

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
	expected = funk.Map(
		expectedProtos,
		func(sub *lte_protos.SubscriberData) *protos.DataUpdate {
			data, err := proto.Marshal(sub)
			assert.NoError(t, err)
			return &protos.DataUpdate{Key: "IMSI" + sub.Sid.Id, Value: data}
		},
	)
	actual, err = provider.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

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
	expected = funk.Map(
		expectedProtos,
		func(sub *lte_protos.SubscriberData) *protos.DataUpdate {
			data, err := proto.Marshal(sub)
			assert.NoError(t, err)
			return &protos.DataUpdate{Key: "IMSI" + sub.Sid.Id, Value: data}
		},
	)
	actual, err = provider.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
