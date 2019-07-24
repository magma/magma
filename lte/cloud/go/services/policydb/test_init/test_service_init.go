/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"testing"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb"
	"magma/lte/cloud/go/services/policydb/servicers"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, lte.ModuleName, policydb.ServiceName)
	protos.RegisterPolicyDBControllerServer(
		srv.GrpcServer,
		servicers.NewPolicyDBServer(test_utils.NewMockDatastore()))
	go srv.RunTest(lis)
}
