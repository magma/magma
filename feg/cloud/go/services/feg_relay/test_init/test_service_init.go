/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"testing"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/feg_relay"
	"magma/feg/cloud/go/services/feg_relay/servicers"
	"magma/orc8r/cloud/go/test_utils"

	"golang.org/x/net/context"
)

// A little Go "polymorphism" magic for testing
type testFegProxyServer struct {
	servicers.FegToGwRelayServer
}

func (srv *testFegProxyServer) CancelLocation(
	ctx context.Context,
	req *protos.CancelLocationRequest,
) (*protos.CancelLocationAnswer, error) {
	return srv.CancelLocationUnverified(ctx, req)
}

func StartTestService(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, feg.ModuleName, feg_relay.ServiceName)
	protos.RegisterS6AGatewayServiceServer(srv.GrpcServer, &testFegProxyServer{})
	go srv.RunTest(lis)
}
