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
	"testing"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	lte_service "magma/lte/cloud/go/services/lte"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	subscriber_streamer "magma/lte/cloud/go/services/subscriberdb/streamer"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	streamer_protos "magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

// Ensure provider servicer properly forwards update requests
func TestLTEStreamProviderServicer_GetUpdates(t *testing.T) {
	const (
		hwID = "some_hwid"
	)
	var (
		subscriberStreamer = &subscriber_streamer.SubscribersProvider{}
	)

	configurator_test_init.StartTestService(t)
	lte_test_init.StartTestService(t)

	conn, err := registry.GetConnection(lte_service.ServiceName)
	assert.NoError(t, err)
	c := streamer_protos.NewStreamProviderClient(conn)
	ctx := context.Background()

	t.Run("subscriber streamer", func(t *testing.T) {
		initSubscriber(t, hwID)
		got, err := c.GetUpdates(ctx, &protos.StreamRequest{
			GatewayId:  hwID,
			StreamName: lte.SubscriberStreamName,
			ExtraArgs:  nil,
		})
		assert.NoError(t, err)
		want, err := subscriberStreamer.GetUpdates(hwID, nil)
		assert.NoError(t, err)
		assert.Equal(t, &protos.DataUpdateBatch{Updates: want}, got)
	})
}

func initSubscriber(t *testing.T, hwID string) {
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: hwID}, serdes.Entity)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "g1"}, serdes.Entity)
	assert.NoError(t, err)

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
					StaticIps: map[string]strfmt.IPv4{"apn1": "192.168.100.1"},
				},
				Associations: []storage.TypeAndKey{{Type: lte.APNEntityType, Key: "apn1"}, {Type: lte.APNEntityType, Key: "apn2"}},
			},
			{Type: lte.SubscriberEntityType, Key: "IMSI67890", Config: &models.SubscriberConfig{Lte: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}}},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
}
