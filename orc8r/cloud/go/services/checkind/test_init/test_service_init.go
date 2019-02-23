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
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/checkind"
	"magma/orc8r/cloud/go/services/checkind/servicers"
	"magma/orc8r/cloud/go/services/checkind/store"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	checkinStore, err := store.NewCheckinStore(test_utils.GetMockDatastoreInstance())
	if err != nil {
		t.Fatalf("Failed to initialize checkin store: %s", err)
	}
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, checkind.ServiceName)
	serviser, err := servicers.NewTestCheckindServer(checkinStore)
	if err != nil {
		t.Fatalf("Failed to create checkin servisers: %s", err)
	}
	protos.RegisterCheckindServer(srv.GrpcServer, serviser)
	go srv.RunTest(lis)
}
