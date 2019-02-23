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
	"magma/feg/cloud/go/services/health/storage"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) (*servicers.TestHealthServer, error) {
	healthStore, err := storage.NewHealthStore(test_utils.GetMockDatastoreInstance())
	if err != nil {
		t.Fatalf("Failed to initialize health store: %s", err)
	}
	clusterStore, err := storage.NewClusterStore(test_utils.GetMockDatastoreInstance())
	if err != nil {
		t.Fatalf("Failed to initialize cluster store: %s", err)
	}
	srv, lis := test_utils.NewTestService(t, feg.ModuleName, health.ServiceName)
	servicer := servicers.NewTestHealthServer(healthStore, clusterStore)

	protos.RegisterHealthServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
	return servicer, nil
}
