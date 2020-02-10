/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/dispatcher"
	"magma/orc8r/cloud/go/services/dispatcher/broker/mocks"
	"magma/orc8r/cloud/go/services/dispatcher/servicers"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

func StartTestService(t *testing.T) *mocks.GatewayRPCBroker {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, dispatcher.ServiceName)
	mockBroker := new(mocks.GatewayRPCBroker)
	servicer, err := servicers.NewTestSyncRPCServer("test host name", mockBroker)
	if err != nil {
		t.Fatalf("Failed to create syncRPCService servicer: %s", err)
	}
	protos.RegisterSyncRPCServiceServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
	return mockBroker
}
