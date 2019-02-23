/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streamer_test

import (
	"encoding/json"
	"testing"

	"magma/lte/cloud/go/protos"
	sdb "magma/lte/cloud/go/services/subscriberdb"
	sdbstreamer "magma/lte/cloud/go/services/subscriberdb/streamer"
	sdb_test_init "magma/lte/cloud/go/services/subscriberdb/test_init"
	orcprotos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"
	"magma/orc8r/cloud/go/services/streamer"
	"magma/orc8r/cloud/go/services/streamer/providers"
	streamer_test_init "magma/orc8r/cloud/go/services/streamer/test_init"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const testAgHwId = "Test-AGW-Hw-Id"

func TestSubscriberdbStreamer(t *testing.T) {
	// Setup - start services, register provider
	magmad_test_init.StartTestService(t)
	sdb_test_init.StartTestService(t)
	streamer_test_init.StartTestService(t)
	err := providers.RegisterStreamProvider(&sdbstreamer.SubscribersProvider{})
	assert.NoError(t, err)

	testNetworkId, err := magmad.RegisterNetwork(&magmad_protos.MagmadNetworkRecord{Name: "Test Network 1"}, "subscriberdb_streamer_test_network")
	assert.NoError(t, err)

	hwId1 := orcprotos.AccessGatewayID{Id: testAgHwId}
	_, err = magmad.RegisterGateway(testNetworkId, &magmad_protos.AccessGatewayRecord{HwId: &hwId1, Name: "bla"})
	assert.NoError(t, err)

	netId := orcprotos.NetworkID{Id: testNetworkId}
	sid1 := protos.SubscriberID{Id: "12345"}
	sub1 := protos.SubscriberData{Sid: &sid1, NetworkId: &netId}
	sid2 := protos.SubscriberID{Id: "67890"}
	sub2 := protos.SubscriberData{Sid: &sid2, NetworkId: &netId}

	// Add a subscribers
	err = sdb.AddSubscriber(testNetworkId, &sub1)
	assert.NoError(t, err)
	err = sdb.AddSubscriber(testNetworkId, &sub2)
	assert.NoError(t, err)

	subscribers, err := sdb.ListSubscribers(testNetworkId)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(subscribers))

	conn, err := registry.GetConnection(streamer.ServiceName)
	assert.NoError(t, err)

	grpcClient := orcprotos.NewStreamerClient(conn)
	streamerClient, err := grpcClient.GetUpdates(
		context.Background(),
		&orcprotos.StreamRequest{GatewayId: testAgHwId, StreamName: "subscriberdb"},
	)
	assert.NoError(t, err)

	updateBatch, err := streamerClient.Recv()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(updateBatch.GetUpdates()))
	var s1, s2 protos.SubscriberData
	err = proto.Unmarshal(updateBatch.Updates[0].Value, &s1)
	assert.NoError(t, err)
	err = proto.Unmarshal(updateBatch.Updates[1].Value, &s2)
	assert.NoError(t, err)
	s1j, _ := json.Marshal(s1)
	s2j, _ := json.Marshal(s2)
	t.Logf("\nReceived Subscribers:\n\t%s\n\t%s", string(s1j), string(s2j))
}
