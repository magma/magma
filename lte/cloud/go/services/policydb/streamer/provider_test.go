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
	"magma/lte/cloud/go/services/policydb"
	pdbstreamer "magma/lte/cloud/go/services/policydb/streamer"
	policydb_test_init "magma/lte/cloud/go/services/policydb/test_init"
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

func TestPolicydbStreamer(t *testing.T) {
	// Setup - start services, register provider
	magmad_test_init.StartTestService(t)
	policydb_test_init.StartTestService(t)
	streamer_test_init.StartTestService(t)
	err := providers.RegisterStreamProvider(&pdbstreamer.PoliciesProvider{})
	assert.NoError(t, err)

	testNetworkId, err := magmad.RegisterNetwork(&magmad_protos.MagmadNetworkRecord{Name: "Test Network 1"}, "policydb_streamer_test_network")
	assert.NoError(t, err)

	hwId1 := orcprotos.AccessGatewayID{Id: testAgHwId}
	_, err = magmad.RegisterGateway(testNetworkId, &magmad_protos.AccessGatewayRecord{HwId: &hwId1, Name: "bla"})
	assert.NoError(t, err)

	rule1 := &protos.PolicyRule{
		Id: "1",
		FlowList: []*protos.FlowDescription{
			{
				Match: &protos.FlowMatch{TcpSrc: 0},
			},
		},
		Priority: 10,
	}
	rule2 := &protos.PolicyRule{
		Id: "2",
		FlowList: []*protos.FlowDescription{
			{
				Match: &protos.FlowMatch{IpProto: 7},
			},
		},
		Priority: 15,
	}

	// Add policies
	err = policydb.AddRule(testNetworkId, rule1)
	assert.NoError(t, err)
	err = policydb.AddRule(testNetworkId, rule2)
	assert.NoError(t, err)

	policies, err := policydb.ListRuleIds(testNetworkId)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(policies))

	conn, err := registry.GetConnection(streamer.ServiceName)
	assert.NoError(t, err)

	grpcClient := orcprotos.NewStreamerClient(conn)
	streamerClient, err := grpcClient.GetUpdates(
		context.Background(),
		&orcprotos.StreamRequest{GatewayId: testAgHwId, StreamName: "policydb"},
	)
	assert.NoError(t, err)

	updateBatch, err := streamerClient.Recv()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(updateBatch.GetUpdates()))
	var p1, p2 protos.PolicyRule
	err = proto.Unmarshal(updateBatch.Updates[0].Value, &p1)
	assert.NoError(t, err)
	err = proto.Unmarshal(updateBatch.Updates[1].Value, &p2)
	assert.NoError(t, err)
	p1j, _ := json.Marshal(p1)
	p2j, _ := json.Marshal(p2)
	t.Logf("\nReceived Policies:\n\t%s\n\t%s", string(p1j), string(p2j))
}
