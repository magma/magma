/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"testing"

	"magma/orc8r/cloud/go/datastore/mocks"
	"magma/orc8r/cloud/go/orc8r"
	accessd_test_init "magma/orc8r/cloud/go/services/accessd/test_init"
	certifier_test_init "magma/orc8r/cloud/go/services/certifier/test_init"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/config/protos"
	"magma/orc8r/cloud/go/services/config/servicers"
	config_storage_mocks "magma/orc8r/cloud/go/services/config/storage/mocks"
	"magma/orc8r/cloud/go/services/magmad"
	mdprotos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_servicers "magma/orc8r/cloud/go/services/magmad/servicers"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, magmad.ServiceName)
	mdprotos.RegisterMagmadConfiguratorServer(
		srv.GrpcServer,
		magmad_servicers.NewMagmadConfigurator(test_utils.NewMockDatastore()))
	go srv.GrpcServer.Serve(lis)

	// magmad has dependency on accessd and certifier for identity
	accessd_test_init.StartTestService(t)
	certifier_test_init.StartTestService(t)

	// TODO: Remove this after fully migrating to config service and deleting the multiplexed writes!
	configSrv, lis := test_utils.NewTestService(t, orc8r.ModuleName, config.ServiceName)
	protos.RegisterConfigServiceServer(
		configSrv.GrpcServer,
		servicers.NewConfigService(config_storage_mocks.NewMapBackedConfigurationStorage()),
	)
	go configSrv.GrpcServer.Serve(lis)
}

func StartTestServiceMockStore(t *testing.T) *mocks.Api {
	mockStore := new(mocks.Api)
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, magmad.ServiceName)
	mdprotos.RegisterMagmadConfiguratorServer(
		srv.GrpcServer,
		magmad_servicers.NewMagmadConfigurator(mockStore))
	go srv.GrpcServer.Serve(lis)

	// TODO: Remove this after fully migrating to config service and deleting the multiplexed writes!
	configSrv, lis := test_utils.NewTestService(t, orc8r.ModuleName, config.ServiceName)
	protos.RegisterConfigServiceServer(
		configSrv.GrpcServer,
		servicers.NewConfigService(config_storage_mocks.NewMapBackedConfigurationStorage()),
	)
	go configSrv.GrpcServer.Serve(lis)

	return mockStore
}
