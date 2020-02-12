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
	"magma/orc8r/cloud/go/services/accessd"
	"magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/cloud/go/services/accessd/servicers"
	"magma/orc8r/cloud/go/services/accessd/storage"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, accessd.ServiceName)
	ds := test_utils.GetMockDatastoreInstance()
	accessdStore := storage.NewAccessdDatastore(ds)
	protos.RegisterAccessControlManagerServer(
		srv.GrpcServer,
		servicers.NewAccessdServer(accessdStore))
	go srv.GrpcServer.Serve(lis)
}
