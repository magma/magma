/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package test_init

import (
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/device/protos"
	"magma/orc8r/cloud/go/services/device/servicers"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

// StartTestService instantiates a service backed by an in-memory storage
func StartTestService(t *testing.T) {
	factory := blobstore.NewMemoryBlobStorageFactory()
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, device.ServiceName)
	server, err := servicers.NewDeviceServicer(factory)
	assert.NoError(t, err)
	protos.RegisterDeviceServer(srv.GrpcServer, server)
	go srv.RunTest(lis)
}
