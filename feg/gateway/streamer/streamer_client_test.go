/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streamer

import (
	"fmt"
	"testing"
	"time"

	"magma/orc8r/cloud/go/services/streamer"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"magma/orc8r/cloud/go/services/streamer/providers"
	streamer_test_init "magma/orc8r/cloud/go/services/streamer/test_init"
	"magma/orc8r/lib/go/protos"
	platform_registry "magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"
)

const (
	testStreamName = "mock1"
	testGwID       = "hwId"
)

// Mock Cloud Streamer
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

var expected = []*protos.DataUpdate{
	{Key: "a", Value: []byte("123")},
	{Key: "b", Value: []byte("456")},
}

// Mock Client Streamer Listener
type testListener struct {
	err       chan error
	updateErr chan error
}

func (l testListener) GetName() string {
	return testStreamName
}

func (l testListener) ReportError(e error) error {
	l.err <- e
	return nil // continue listener
}

func (l testListener) Update(ub *protos.DataUpdateBatch) bool {
	if len(expected) != len(ub.GetUpdates()) {
		l.updateErr <- fmt.Errorf("Updates # %d != expected # %d", len(ub.GetUpdates()), len(expected))
		return false
	}
	for i, u := range ub.GetUpdates() {
		if protos.TestMarshal(expected[i]) != protos.TestMarshal(u) {
			l.updateErr <- fmt.Errorf(
				"Update %s != expected %s", protos.TestMarshal(u), protos.TestMarshal(expected[i]))
			return false
		}
	}
	l.updateErr <- nil
	return true
}

// Mock GW Cloud Service registry
type mockedCloudRegistry struct{}

func (cr mockedCloudRegistry) GetCloudConnection(service string) (*grpc.ClientConn, error) {
	if service != StreamerServiceName {
		return nil, fmt.Errorf("Not Implemented")
	}
	return platform_registry.GetConnection(streamer.ServiceName)
}

func (cr mockedCloudRegistry) GetCloudConnectionFromServiceConfig(serviceConfig *config.ConfigMap, service string) (*grpc.ClientConn, error) {
	return nil, fmt.Errorf("Not Implemented")

}

// Test
func TestStreamerClient(t *testing.T) {
	streamer_test_init.StartTestService(t)

	streamerClient := NewStreamerClient(mockedCloudRegistry{})

	providers.RegisterStreamProvider(&mockStreamProvider{name: testStreamName, retVal: expected})

	l := testListener{}
	l.err = make(chan error)
	l.updateErr = make(chan error)
	assert.NoError(t, streamerClient.AddListener(l))

	select {
	case e := <-l.err:
		assert.NoError(t, e)
	case e := <-l.updateErr:
		assert.NoError(t, e)
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Test Timeout")
	}
}
