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
	"magma/orc8r/cloud/go/services/streamer"
	"magma/orc8r/cloud/go/services/streamer/servicers"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

// A little Go "polymorphism" magic for testing
type testStreamingServer struct {
	servicers.StreamingServer
}

func (srv *testStreamingServer) GetUpdates(
	request *protos.StreamRequest,
	stream protos.Streamer_GetUpdatesServer,
) error {
	return servicers.GetUpdatesUnverified(request, stream)
}

func StartTestService(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, streamer.ServiceName)
	protos.RegisterStreamerServer(srv.GrpcServer, &testStreamingServer{})
	go srv.RunTest(lis)
}
