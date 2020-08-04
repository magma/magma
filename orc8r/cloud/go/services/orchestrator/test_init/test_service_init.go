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
	"magma/orc8r/cloud/go/services/orchestrator"
	"magma/orc8r/cloud/go/services/orchestrator/servicers"
	indexer_protos "magma/orc8r/cloud/go/services/state/protos"
	streamer_protos "magma/orc8r/cloud/go/services/streamer/protos"
	streamer_servicers "magma/orc8r/cloud/go/services/streamer/servicers"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
)

type testStreamerServer struct {
	protos.StreamerServer
}

func (srv *testStreamerServer) GetUpdates(req *protos.StreamRequest, stream protos.Streamer_GetUpdatesServer) error {
	return streamer_servicers.GetUpdatesUnverified(req, stream)
}

func StartTestService(t *testing.T) {
	labels := map[string]string{
		orc8r.StreamProviderLabel: "true",
	}
	annotations := map[string]string{
		orc8r.StreamProviderStreamsAnnotation: definitions.MconfigStreamName,
	}
	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, orchestrator.ServiceName, labels, annotations)

	protos.RegisterStreamerServer(srv.GrpcServer, &testStreamerServer{})

	indexer_protos.RegisterIndexerServer(srv.GrpcServer, servicers.NewDirectoryIndexer())
	streamer_protos.RegisterStreamProviderServer(srv.GrpcServer, servicers.NewOrchestratorStreamProviderServicer())

	go srv.RunTest(lis)
}
