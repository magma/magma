/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"magma/orc8r/cloud/go/services/dispatcher/broker"
	"magma/orc8r/lib/go/protos"
)

// A little Go "polymorphism" magic for testing
type testSyncRPCServer struct {
	SyncRPCService
}

const TestSyncRPCAgHwId = "Test-AGW-Hw-Id"

func (srv *testSyncRPCServer) EstablishSyncRPCStream(stream protos.SyncRPCService_EstablishSyncRPCStreamServer) error {
	// See if there is an Identity in the CTX and if not, use default TestSyncRPCAgHwId
	gw := protos.GetClientGateway(stream.Context())
	if gw == nil {
		return srv.serveGwId(stream, TestSyncRPCAgHwId)
	}
	return srv.SyncRPCService.EstablishSyncRPCStream(stream)
}

func NewTestSyncRPCServer(hostName string, broker broker.GatewayRPCBroker) (*testSyncRPCServer, error) {
	return &testSyncRPCServer{SyncRPCService{hostName, broker}}, nil
}
