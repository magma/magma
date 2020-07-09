/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"context"
	"testing"

	"magma/lte/cloud/go/lte"
	lte_plugin "magma/lte/cloud/go/plugin"
	lte_service "magma/lte/cloud/go/services/lte"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	subscriber_streamer "magma/lte/cloud/go/services/subscriberdb/streamer"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	streamer_protos "magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	assert "github.com/stretchr/testify/require"
)

// Ensure provider servicer properly forwards update requests
func TestLTEStreamProviderServicer_GetUpdates(t *testing.T) {
	const (
		hwID = "some_hwid"
	)
	var (
		subscriberStreamer = &subscriber_streamer.SubscribersProvider{}
	)

	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{}))
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
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: hwID})
	assert.NoError(t, err)

	_, err = configurator.CreateEntities("n1", []configurator.NetworkEntity{
		{
			Type: lte.SubscriberEntityType, Key: "IMSI12345",
			Config: &models.LteSubscription{
				State:   "ACTIVE",
				AuthKey: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				AuthOpc: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			},
		},
		{Type: lte.SubscriberEntityType, Key: "IMSI67890", Config: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}},
	})
	assert.NoError(t, err)
}
