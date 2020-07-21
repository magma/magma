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
	streamer_protos "magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/cloud/go/services/streamer/servicers"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
)

type testStreamerServer struct {
	protos.StreamerServer
}

func (srv *testStreamerServer) GetUpdates(req *protos.StreamRequest, stream protos.Streamer_GetUpdatesServer) error {
	return servicers.GetUpdatesUnverified(req, stream)
}

func StartTestService(t *testing.T) {
	labels := map[string]string{
		orc8r.StreamProviderLabel: "true",
	}
	annotations := map[string]string{
		orc8r.StreamProviderStreamsAnnotation: definitions.MconfigStreamName,
	}
	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, streamer.ServiceName, labels, annotations)
	protos.RegisterStreamerServer(srv.GrpcServer, &testStreamerServer{})
	streamer_protos.RegisterStreamProviderServer(srv.GrpcServer, servicers.NewBaseOrchestratorStreamProviderServicer())
	go srv.RunTest(lis)
}
