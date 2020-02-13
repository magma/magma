/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"testing"

	"magma/orc8r/cloud/go/services/directoryd/servicers"
	"magma/orc8r/cloud/go/services/directoryd/storage"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const (
	testNetworkId    = "network"
	testGwHwId1      = "gw1"
	testGwHwId2      = "gw2"
	testGwHwId3      = "gw3"
	testGwLogicalId1 = "logical gw1"
	testGwLogicalId2 = "logical gw2"
	testSubId1       = "subscriber1"
	testSubId2       = "subscriber2"
	testSubId3       = "subscriber3"
)

func createTestDirectorydServicer(t *testing.T) *servicers.DirectoryServicer {
	mockDB := test_utils.NewMockDatastore()
	srv, _ := servicers.NewDirectoryServicer(storage.GetDirectorydPersistenceService(mockDB))
	return srv
}

func TestDirectorydGetUnknownLocation(t *testing.T) {
	srv := createTestDirectorydServicer(t)
	stateTestInit.StartTestService(t)
	// Create an identity and context for sending requests as gateway
	id := protos.Identity{}
	idgw := protos.Identity_Gateway{HardwareId: testGwHwId1, NetworkId: testNetworkId, LogicalId: testGwLogicalId1}
	id.SetGateway(&idgw)
	ctx := id.NewContextWithIdentity(context.Background())

	// Create a record
	request := protos.GetLocationRequest{
		Id: testSubId1,
	}
	_, err := srv.GetLocation(ctx, &request)
	assert.Error(t, err)
}

func TestDirectorydUpdateLocation(t *testing.T) {
	srv := createTestDirectorydServicer(t)
	stateTestInit.StartTestService(t)
	// Create an identity and context for sending requests as gateway
	id := protos.Identity{}
	idgw := protos.Identity_Gateway{HardwareId: testGwHwId1, NetworkId: testNetworkId, LogicalId: testGwLogicalId1}
	id.SetGateway(&idgw)
	ctx := id.NewContextWithIdentity(context.Background())

	// Create a record
	location1 := protos.LocationRecord{
		Location: testGwHwId1,
	}
	request1 := protos.UpdateDirectoryLocationRequest{
		Id:     testSubId1,
		Record: &location1,
	}
	_, err := srv.UpdateLocation(ctx, &request1)
	assert.NoError(t, err)

	// Create another record
	location2 := protos.LocationRecord{
		Location: testGwHwId2,
	}
	request2 := protos.UpdateDirectoryLocationRequest{
		Id:     testSubId2,
		Record: &location2,
	}
	_, err = srv.UpdateLocation(ctx, &request2)
	assert.NoError(t, err)

	// read check
	getRequest := protos.GetLocationRequest{
		Id: testSubId1,
	}
	record, err := srv.GetLocation(ctx, &getRequest)
	assert.NoError(t, err)
	assert.Equal(t, record, &location1)

	getRequest.Id = testSubId2
	record, err = srv.GetLocation(ctx, &getRequest)
	assert.NoError(t, err)
	// expect the location to be overriden by the hardwareId read in context
	assert.Equal(t, record, &location1)

	// Create another gateway identity and context
	id2 := protos.Identity{}
	idgw2 := protos.Identity_Gateway{HardwareId: testGwHwId2, NetworkId: testNetworkId, LogicalId: testGwLogicalId2}
	id2.SetGateway(&idgw2)
	ctx2 := id2.NewContextWithIdentity(context.Background())

	// Add a record for the second gateway
	location3 := protos.LocationRecord{
		Location: testGwHwId3,
	}
	request3 := protos.UpdateDirectoryLocationRequest{
		Id:     testSubId3,
		Record: &location3,
	}
	_, err = srv.UpdateLocation(ctx2, &request3)
	assert.NoError(t, err)

	// read check
	getRequest.Id = testSubId3
	record, err = srv.GetLocation(ctx2, &getRequest)
	assert.NoError(t, err)
	assert.Equal(t, record, &location2)
}

func TestDirectorydDeleteUnknownLocation(t *testing.T) {
	srv := createTestDirectorydServicer(t)
	stateTestInit.StartTestService(t)
	// Create an identity and context for sending requests as gateway
	id := protos.Identity{}
	idgw := protos.Identity_Gateway{HardwareId: testGwHwId1, NetworkId: testNetworkId, LogicalId: testGwLogicalId1}
	id.SetGateway(&idgw)
	ctx := id.NewContextWithIdentity(context.Background())

	// Create a record
	request := protos.DeleteLocationRequest{
		Id: testSubId1,
	}
	_, err := srv.DeleteLocation(ctx, &request)
	assert.Error(t, err)
}

func TestDirectorydDeleteLocation(t *testing.T) {
	srv := createTestDirectorydServicer(t)
	stateTestInit.StartTestService(t)
	// Create an identity and context for sending requests as gateway
	id := protos.Identity{}
	idgw := protos.Identity_Gateway{HardwareId: testGwHwId1, NetworkId: testNetworkId, LogicalId: testGwLogicalId1}
	id.SetGateway(&idgw)
	ctx := id.NewContextWithIdentity(context.Background())

	// Create a record
	location1 := protos.LocationRecord{
		Location: testGwHwId1,
	}
	request1 := protos.UpdateDirectoryLocationRequest{
		Id:     testSubId1,
		Record: &location1,
	}
	_, err := srv.UpdateLocation(ctx, &request1)
	assert.NoError(t, err)

	// delete
	deleteRequest := protos.DeleteLocationRequest{
		Id: testSubId1,
	}
	_, err = srv.DeleteLocation(ctx, &deleteRequest)
	assert.NoError(t, err)

	// confirm deletion
	getRequest := protos.GetLocationRequest{
		Id: testSubId1,
	}
	_, err = srv.GetLocation(ctx, &getRequest)
	assert.Error(t, err)
}
