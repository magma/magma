/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"errors"
	"testing"

	"magma/orc8r/cloud/go/services/streamer"
	"magma/orc8r/cloud/go/services/streamer/providers"
	streamer_test_init "magma/orc8r/cloud/go/services/streamer/test_init"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

type mockStreamProvider struct {
	name   string
	retVal []*protos.DataUpdate
	retErr error
}

func (m *mockStreamProvider) GetStreamName() string {
	return m.name
}

func (m *mockStreamProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	return m.retVal, m.retErr
}

func TestStreamingServer_GetUpdates(t *testing.T) {
	streamer_test_init.StartTestService(t)
	conn, err := registry.GetConnection(streamer.ServiceName)
	assert.NoError(t, err)
	grpcClient := protos.NewStreamerClient(conn)

	expected := []*protos.DataUpdate{
		{Key: "a", Value: []byte("123")},
		{Key: "b", Value: []byte("456")},
	}
	providers.RegisterStreamProvider(&mockStreamProvider{name: "mock1", retVal: expected})

	streamerClient, err := grpcClient.GetUpdates(
		context.Background(),
		&protos.StreamRequest{GatewayId: "hwId", StreamName: "mock1"},
	)
	assert.NoError(t, err)

	actual, err := streamerClient.Recv()
	updates := actual.GetUpdates()
	assert.Equal(t, len(expected), len(updates))

	for i, u := range updates {
		assert.Equal(t, protos.TestMarshal(expected[i]), protos.TestMarshal(u))
	}

	// Error in provider
	providers.RegisterStreamProvider(&mockStreamProvider{name: "mock2", retVal: nil, retErr: errors.New("MOCK")})
	streamerClient, err = grpcClient.GetUpdates(
		context.Background(),
		&protos.StreamRequest{GatewayId: "hwId", StreamName: "mock2"},
	)
	assert.NoError(t, err)
	_, err = streamerClient.Recv()
	assert.Error(t, err)
	assert.Equal(t, "rpc error: code = Aborted desc = Error while streaming updates: MOCK", err.Error())

	// Provider does not exist
	streamerClient, err = grpcClient.GetUpdates(
		context.Background(),
		&protos.StreamRequest{GatewayId: "hwId", StreamName: "stream_dne"},
	)
	assert.NoError(t, err)
	_, err = streamerClient.Recv()
	assert.Error(t, err, "Stream stream_dne does not exist", codes.Unavailable)
}
