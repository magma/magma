/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"net"
	"testing"

	"magma/orc8r/cloud/go/services/dispatcher/broker/mocks"
	"magma/orc8r/cloud/go/services/dispatcher/httpserver"
)

func StartTestHttpServer(t *testing.T) (net.Addr, *mocks.GatewayRPCBroker) {
	lis, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatalf("net.Listen err: %v\n", err)
	}

	broker := new(mocks.GatewayRPCBroker)
	server := httpserver.NewSyncRPCHttpServer(broker)
	go server.Serve(lis)
	return lis.Addr(), broker
}
