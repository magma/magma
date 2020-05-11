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
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/servicers"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

// StartTestService instantiates a service backed by an in-memory storage.
func StartTestService(t *testing.T) {
	_ = StartTestServiceInternal(t)
}

// StartTestServiceInternal instantiates a service backed by an in-memory
// storage, exposing the servicer's internal methods.
func StartTestServiceInternal(t *testing.T) servicers.StateServiceInternal {
	factory := blobstore.NewMemoryBlobStorageFactory()
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, state.ServiceName)
	server, err := servicers.NewStateServicer(factory)
	assert.NoError(t, err)
	protos.RegisterStateServiceServer(srv.GrpcServer, server)
	go srv.RunTest(lis)
	return server
}
