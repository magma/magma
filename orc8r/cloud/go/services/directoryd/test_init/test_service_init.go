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
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/directoryd/servicers"
	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

func StartTestService(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, directoryd.ServiceName)

	db := test_utils.GetMockDatastoreInstance()
	persistence_service := storage.GetDirectorydPersistenceService(db)

	// Create directory gRPC servicer
	directory_gRPC_servicer, err := servicers.NewDirectoryServicer(persistence_service)
	if err != nil {
		glog.Errorf("Error creating directory gRPC servicer: %s", err)
	}

	// Add gRPC servicer to the directory service gRPC server
	protos.RegisterDirectoryServiceServer(srv.GrpcServer, directory_gRPC_servicer)

	go srv.RunTest(lis)
}
