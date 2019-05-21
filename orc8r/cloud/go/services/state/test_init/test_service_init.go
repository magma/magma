/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/servicers"
	"magma/orc8r/cloud/go/test_utils"
)

// StartTestService instantiates a service backed by an in-memory storage
func StartTestService(t *testing.T) {
	factory := blobstore.NewMemoryBlobStorageFactory()
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, state.ServiceName)
	server, err := servicers.NewStateServicer(factory)
	if err != nil {
		t.Fatalf("Failure to start state test service: %v", err)
	}
	protos.RegisterStateServiceServer(srv.GrpcServer, server)
	go srv.RunTest(lis)
}
