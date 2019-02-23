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
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/servicers"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

func StartTestService(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, lte.ModuleName, subscriberdb.ServiceName)
	subscriberDBTestStore, _ := storage.NewSubscriberDBStorage(test_utils.GetMockDatastoreInstance())
	subscriberDBTestSrv, err := servicers.NewSubscriberDBServer(subscriberDBTestStore)
	assert.NoError(t, err)
	protos.RegisterSubscriberDBControllerServer(
		srv.GrpcServer,
		subscriberDBTestSrv)
	go srv.RunTest(lis)
}
