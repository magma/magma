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
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/config/protos"
	"magma/orc8r/cloud/go/services/config/servicers"
	"magma/orc8r/cloud/go/services/config/storage/mocks"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, config.ServiceName)
	servicer := servicers.NewConfigService(mocks.NewMapBackedConfigurationStorage())
	protos.RegisterConfigServiceServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
}
