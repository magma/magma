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
	lte_service "magma/lte/cloud/go/services/lte"
	"magma/lte/cloud/go/services/lte/servicers"
	"magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, lte.ModuleName, lte_service.ServiceName)
	protos.RegisterStreamProviderServer(srv.GrpcServer, servicers.NewLTEStreamProviderServicer())
	go srv.RunTest(lis)
}
