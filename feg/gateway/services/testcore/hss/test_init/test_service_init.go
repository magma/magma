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

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/testcore/hss/servicers/test_utils"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) (*service.Service, error) {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.MOCK_HSS)

	service := test_utils.NewTestHomeSubscriberServer(t)

	protos.RegisterHSSConfiguratorServer(srv.GrpcServer, service)
	go srv.RunTest(lis)

	return srv, nil
}
