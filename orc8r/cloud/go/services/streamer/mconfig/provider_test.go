/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package mconfig_test

import (
	"errors"
	"testing"
	"time"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"
	"magma/orc8r/cloud/go/services/streamer"
	mconfig_provider "magma/orc8r/cloud/go/services/streamer/mconfig"
	"magma/orc8r/cloud/go/services/streamer/mconfig/factory"
	"magma/orc8r/cloud/go/services/streamer/mconfig/test_protos"
	"magma/orc8r/cloud/go/services/streamer/providers"
	streamer_test_init "magma/orc8r/cloud/go/services/streamer/test_init"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const testAgHwId = "Test-AGW-Hw-Id"

type mockMconfigBuilder struct {
	retVal map[string]proto.Message
	retErr error
}

func (builder *mockMconfigBuilder) Build(networkId string, gatewayId string) (map[string]proto.Message, error) {
	return builder.retVal, builder.retErr
}

type mockClock struct {
	now time.Time
}

func (mockClock *mockClock) Now() time.Time {
	return mockClock.now
}

// Test AG Configs Streaming
func TestMconfigStreamer(t *testing.T) {
	magmad_test_init.StartTestService(t)
	streamer_test_init.StartTestService(t)
	providers.RegisterStreamProvider(&mconfig_provider.ConfigProvider{})

	testNetworkId, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "Test Network 1"},
		"mconfig_streamer_test_network")
	assert.NoError(t, err)

	hwId1 := protos.AccessGatewayID{Id: testAgHwId}
	gwId1, err := magmad.RegisterGateway(testNetworkId, &magmad_protos.AccessGatewayRecord{HwId: &hwId1, Name: "bla"})
	assert.NoError(t, err)

	hwId2 := protos.AccessGatewayID{Id: testAgHwId + "second"}
	_, err = magmad.RegisterGateway(testNetworkId, &magmad_protos.AccessGatewayRecord{HwId: &hwId2, Name: "bla2"})
	assert.NoError(t, err)

	// Setup mock mconfig builders
	builder1 := &mockMconfigBuilder{
		retVal: map[string]proto.Message{
			"builder1_1": &test_protos.Message1{Field: "hello"},
			"builder1_2": &test_protos.Message2{Field1: "hello", Field2: "world"},
		},
		retErr: nil,
	}
	builder2 := &mockMconfigBuilder{
		retVal: map[string]proto.Message{
			"builder2_1": &test_protos.Message1{Field: "foo"},
		},
	}
	factory.SetClock(t, &mockClock{now: time.Unix(1551916956, 0)})
	factory.ClearMconfigBuilders(t)
	factory.RegisterMconfigBuilder(builder1)
	factory.RegisterMconfigBuilder(builder2)

	// Connect to streamer
	conn, err := registry.GetConnection(streamer.ServiceName)
	assert.NoError(t, err)

	grpcClient := protos.NewStreamerClient(conn)
	streamerClient, err := grpcClient.GetUpdates(
		context.Background(),
		&protos.StreamRequest{GatewayId: testAgHwId, StreamName: "configs"},
	)
	assert.NoError(t, err)

	expectedProtos := map[string]proto.Message{
		"builder1_1": &test_protos.Message1{Field: "hello"},
		"builder1_2": &test_protos.Message2{Field1: "hello", Field2: "world"},
		"builder2_1": &test_protos.Message1{Field: "foo"},
	}
	expected := make(map[string]*any.Any, len(expectedProtos))
	for k, v := range expectedProtos {
		anyV, err := ptypes.MarshalAny(v)
		assert.NoError(t, err)
		expected[k] = anyV
	}
	expectedMarshaled, err := protos.MarshalIntern(&protos.GatewayConfigs{
		ConfigsByKey: expected,
		Metadata: &protos.GatewayConfigsMetadata{
			CreatedAt: 1551916956,
		},
	})
	assert.NoError(t, err)

	// Assert value
	updateBatch, err := streamerClient.Recv()
	assert.NoError(t, err)
	conn.Close()
	assert.Equal(t, 1, len(updateBatch.Updates))
	assert.Equal(t, gwId1, updateBatch.GetUpdates()[0].GetKey())
	assert.Equal(t, expectedMarshaled, updateBatch.GetUpdates()[0].GetValue(), expectedMarshaled)

	// Error in a builder
	errBuilder := &mockMconfigBuilder{
		retVal: nil,
		retErr: errors.New("MOCK ERROR"),
	}
	factory.RegisterMconfigBuilder(errBuilder)

	conn, err = registry.GetConnection(streamer.ServiceName)
	assert.NoError(t, err)

	grpcClient = protos.NewStreamerClient(conn)
	streamerClient, err = grpcClient.GetUpdates(
		context.Background(),
		&protos.StreamRequest{GatewayId: testAgHwId, StreamName: "configs"},
	)
	assert.NoError(t, err)

	updateBatch, err = streamerClient.Recv()
	assert.Error(t, err)
	assert.Equal(t, "rpc error: code = Aborted desc = Error while streaming updates: MOCK ERROR", err.Error())
}
