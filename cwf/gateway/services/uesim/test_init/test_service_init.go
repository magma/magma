/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package test_init

import (
	"testing"

	"magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/registry"
	"magma/cwf/gateway/services/uesim/servicers"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

func StartTestService(t *testing.T) {
	factory := blobstore.NewMemoryBlobStorageFactory()
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.UeSim)
	server, err := servicers.NewUESimServer(factory)
	assert.NoError(t, err)
	protos.RegisterUESimServer(srv.GrpcServer, server)
	go srv.RunTest(lis)
}
