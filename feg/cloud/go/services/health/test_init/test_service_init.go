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
	"magma/feg/cloud/go/services/health"
	"magma/feg/cloud/go/services/health/servicers"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

func StartTestService(t *testing.T) (*servicers.TestHealthServer, error) {
	srv, lis := test_utils.NewTestService(t, feg.ModuleName, health.ServiceName)
	factory := blobstore.NewMemoryBlobStorageFactory()
	servicer, err := servicers.NewTestHealthServer(factory)
	assert.NoError(t, err)
	protos.RegisterHealthServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
	return servicer, nil
}
