/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package service_test

import (
	"testing"
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/state"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

func TestServiceRun(t *testing.T) {
	testStartTime := time.Now().Unix()
	allowedStartRange := 15.0

	// Create the service
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, state.ServiceName)
	assert.Equal(t, protos.ServiceInfo_STARTING, srv.State)
	assert.Equal(t, protos.ServiceInfo_APP_UNHEALTHY, srv.Health)

	// start the service
	go srv.RunTest(lis)

	// wait for the service to be started and check its state and health
	time.Sleep(time.Second)
	assert.Equal(t, protos.ServiceInfo_ALIVE, srv.State)
	assert.Equal(t, protos.ServiceInfo_APP_HEALTHY, srv.Health)

	// Create a rpc stub and query the Service303 interface
	conn, err := registry.GetConnection(state.ServiceName)
	assert.NoError(t, err, "err in getting connection to service")
	client := protos.NewService303Client(conn)

	actualServiceInfo, err := client.GetServiceInfo(context.Background(), new(protos.Void))

	// check GetServiceInfo rpc call.
	expectedServiceInfo := protos.ServiceInfo{
		Name:          "STATE",
		Version:       "0.0.0",
		State:         protos.ServiceInfo_ALIVE,
		Health:        protos.ServiceInfo_APP_HEALTHY,
		StartTimeSecs: actualServiceInfo.StartTimeSecs,
	}
	assert.NoError(t, err, "err in getting service info after srv started")
	assert.Equal(t, expectedServiceInfo, *actualServiceInfo)
	assert.InDelta(t, testStartTime, actualServiceInfo.StartTimeSecs, allowedStartRange)

	// check StopService rpc call.
	// this will have a connection error, which is expected.
	client.StopService(context.Background(), new(protos.Void))

	assert.Equal(t, protos.ServiceInfo_STOPPING, srv.State)
	assert.Equal(t, protos.ServiceInfo_APP_UNHEALTHY, srv.Health)
}
