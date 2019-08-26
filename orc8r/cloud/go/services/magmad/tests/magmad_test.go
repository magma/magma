/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package magmad

import (
	"errors"
	"testing"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/datastore/mocks"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/magmad/servicers"
	magmad_test_service "magma/orc8r/cloud/go/services/magmad/test_init"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testAgHwId = "test_ag_HW_id"
)

func TestRegisterAndGet(t *testing.T) {
	magmad_test_service.StartTestService(t)

	testNetworkId := "magmad_test_network"
	_, err := magmad.RegisterNetwork(&magmad_protos.MagmadNetworkRecord{Name: "Test Network Name"}, testNetworkId)
	assert.NoError(t, err)

	_, err = magmad.FindGatewayNetworkId(testAgHwId)
	assert.Error(t, err)

	logicalId, err := magmad.RegisterGateway(
		testNetworkId,
		&magmad_protos.AccessGatewayRecord{
			HwId: &protos.AccessGatewayID{Id: testAgHwId},
			Name: "Test GW 1",
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, logicalId)
	assert.Equal(t, logicalId, testAgHwId)

	network, err := magmad.FindGatewayNetworkId(testAgHwId)
	assert.NoError(t, err)
	assert.Equal(t, network, testNetworkId)

	// Register gw with same id as network
	_, err = magmad.RegisterGatewayWithId(
		testNetworkId,
		&magmad_protos.AccessGatewayRecord{
			HwId: &protos.AccessGatewayID{Id: "test_ag_HW_id_2"},
			Name: "Test GW 2",
		},
		"magmad_test_NeTwOrk",
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gateway ID must be different from network ID")
}

func TestListFindAndRemove(t *testing.T) {
	magmad_test_service.StartTestService(t)

	testNetworkId, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "Test Network Name"},
		"magmad_test_network")
	assert.NoError(t, err)

	list, err := magmad.ListGateways(testNetworkId)
	assert.NoError(t, err)
	assert.Equal(t, len(list), 0)

	logicalId, err := magmad.RegisterGateway(
		testNetworkId,
		&magmad_protos.AccessGatewayRecord{HwId: &protos.AccessGatewayID{Id: testAgHwId}},
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, logicalId)

	list, err = magmad.ListGateways(testNetworkId)
	assert.NoError(t, err)
	assert.Equal(t, len(list), 1)
	lid, err := magmad.FindGatewayId(testNetworkId, testAgHwId)
	assert.NoError(t, err)
	assert.Equal(t, lid, testAgHwId)

	err = magmad.RemoveGateway(testNetworkId, logicalId)
	assert.NoError(t, err)

	list, err = magmad.ListGateways(testNetworkId)
	assert.NoError(t, err)
	assert.Equal(t, len(list), 0)

	networkStr, err := magmad.FindGatewayNetworkId(testAgHwId)
	assert.Error(t, err)
	assert.Equal(t, networkStr, "")
}

func TestRegisterGwOnTwoNetworks(t *testing.T) {
	magmad_test_service.StartTestService(t)

	testNetworkId, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "Test Network Name 1"},
		"magmad_test_network1")
	assert.NoError(t, err)

	testNetworkId2, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "Test Network Name 2"},
		"magmad_test_network2")
	assert.NoError(t, err)

	list, err := magmad.ListGateways(testNetworkId)
	assert.NoError(t, err)
	assert.Equal(t, len(list), 0)

	list, err = magmad.ListGateways(testNetworkId2)
	assert.NoError(t, err)
	assert.Equal(t, len(list), 0)

	logicalId, err := magmad.RegisterGateway(
		testNetworkId,
		&magmad_protos.AccessGatewayRecord{HwId: &protos.AccessGatewayID{Id: testAgHwId}},
	)
	assert.NoError(t, err)
	assert.NotNil(t, logicalId)

	list, err = magmad.ListGateways(testNetworkId)
	assert.NoError(t, err)
	assert.Equal(t, len(list), 1)

	lid, err := magmad.FindGatewayId(testNetworkId, testAgHwId)
	assert.NoError(t, err)
	assert.Equal(t, lid, testAgHwId)

	logicalId2, err := magmad.RegisterGateway(
		testNetworkId2,
		&magmad_protos.AccessGatewayRecord{HwId: &protos.AccessGatewayID{Id: testAgHwId}},
	)
	assert.Error(t, err)

	err = magmad.RemoveGateway(testNetworkId, logicalId)
	assert.NoError(t, err)

	list, err = magmad.ListGateways(testNetworkId)
	assert.NoError(t, err)
	assert.Equal(t, len(list), 0)

	networkStr, err := magmad.FindGatewayNetworkId(testAgHwId)
	assert.Error(t, err)
	assert.Equal(t, networkStr, "")
	logicalId2, err = magmad.RegisterGateway(
		testNetworkId2,
		&magmad_protos.AccessGatewayRecord{HwId: &protos.AccessGatewayID{Id: testAgHwId}},
	)
	assert.NoError(t, err)
	assert.NotNil(t, logicalId2)

	list, err = magmad.ListGateways(testNetworkId2)
	assert.NoError(t, err)
	assert.Equal(t, len(list), 1)

	err = magmad.RemoveGateway(testNetworkId2, logicalId2)
	assert.NoError(t, err)
}

func TestMagmad_RemoveGateway_NoGWRecord(t *testing.T) {
	mockeryStore := magmad_test_service.StartTestServiceMockStore(t)
	networkId := "NETWORK"
	gwId := "GW1"

	mockeryStore.On("Get", datastore.GetTableName(networkId, servicers.AgRecordTableName), gwId).
		Return(nil, uint64(1), errors.New("Get error"))
	err := magmad.RemoveGateway(networkId, gwId)
	assert.Error(t, err)

	mockeryStore.AssertExpectations(t)
}

func TestMagmad_RemoveGateway_DeleteErrors(t *testing.T) {
	mockeryStore := magmad_test_service.StartTestServiceMockStore(t)
	networkId := "NETWORK"
	hwId := "HW1"
	gwId := "GW1"

	gwRecord := &magmad_protos.AccessGatewayRecord{Name: "GW NAME", HwId: &protos.AccessGatewayID{Id: hwId}}
	marshaledRecord, err := protos.MarshalIntern(gwRecord)
	assert.NoError(t, err)
	mockeryStore.On("Get", datastore.GetTableName(networkId, servicers.AgRecordTableName), gwId).
		Return(marshaledRecord, uint64(1), nil)

	// One does key exist error, one does not exist
	mockeryStore.On("DoesKeyExist", servicers.GatewaysTableName, hwId).
		Return(true, errors.New("DoesKeyExist error"))
	mockeryStore.On("DoesKeyExist", datastore.GetTableName(networkId, servicers.HwIdTableName), hwId).
		Return(true, nil)
	mockeryStore.On("DoesKeyExist", datastore.GetTableName(networkId, servicers.GatewaysStatusTableName), gwId).
		Return(true, nil)

	// One error on delete
	// Shouldn't be any other deletes
	mockeryStore.On("Delete", datastore.GetTableName(networkId, servicers.HwIdTableName), hwId).
		Return(errors.New("Delete error"))
	mockeryStore.On("Delete", datastore.GetTableName(networkId, servicers.GatewaysStatusTableName), gwId).
		Return(nil)

	err = magmad.RemoveGateway(networkId, gwId)
	assert.Error(t, err)
	assert.Contains(
		t,
		err.Error(),
		"Encountered the following errors while removing the gateway:\n"+
			"\tFailed to delete logical ID mapping. Error: Delete error\n"+
			"\tFailed to delete hardware ID network mapping. Error: DoesKeyExist error\n"+
			"Please address the issues and retry the operation.",
	)

	mockeryStore.AssertExpectations(t)
}

func TestMagmad_RemoveGateway(t *testing.T) {
	mockeryStore := magmad_test_service.StartTestServiceMockStore(t)
	networkId := "NETWORK"
	hwId := "HW1"
	gwId := "GW1"

	gwRecord := &magmad_protos.AccessGatewayRecord{Name: "GW NAME", HwId: &protos.AccessGatewayID{Id: hwId}}
	marshaledRecord, err := protos.MarshalIntern(gwRecord)
	assert.NoError(t, err)
	mockeryStore.On("Get", datastore.GetTableName(networkId, servicers.AgRecordTableName), gwId).
		Return(marshaledRecord, uint64(1), nil)

	// One key does not exist
	mockeryStore.On("DoesKeyExist", datastore.GetTableName(networkId, servicers.HwIdTableName), hwId).
		Return(true, nil)
	mockeryStore.On("DoesKeyExist", servicers.GatewaysTableName, hwId).
		Return(true, nil)
	mockeryStore.On("DoesKeyExist", datastore.GetTableName(networkId, servicers.GatewaysStatusTableName), gwId).
		Return(true, nil)

	// Should only be 3 deletes
	mockeryStore.On("Delete", datastore.GetTableName(networkId, servicers.HwIdTableName), hwId).
		Return(nil)
	mockeryStore.On("Delete", servicers.GatewaysTableName, hwId).
		Return(nil)
	mockeryStore.On("Delete", datastore.GetTableName(networkId, servicers.GatewaysStatusTableName), gwId).
		Return(nil)

	// Record delete
	mockeryStore.On("Delete", datastore.GetTableName(networkId, servicers.AgRecordTableName), gwId).
		Return(nil)

	err = magmad.RemoveGateway(networkId, gwId)
	assert.NoError(t, err)

	mockeryStore.AssertExpectations(t)
}

func TestRemoveNetwork(t *testing.T) {
	magmad_test_service.StartTestService(t)

	testNetworkId, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "Test Network Name"},
		"magmad_test_network")
	assert.NoError(t, err)

	logicalId, err := magmad.RegisterGateway(
		testNetworkId,
		&magmad_protos.AccessGatewayRecord{HwId: &protos.AccessGatewayID{Id: testAgHwId}},
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, logicalId)

	err = magmad.RemoveNetwork(testNetworkId)
	assert.Error(t, err, "Network is non empty")

	err = magmad.RemoveGateway(testNetworkId, logicalId)
	assert.NoError(t, err)

	list, err := magmad.ListGateways(testNetworkId)
	assert.NoError(t, err)
	assert.Equal(t, len(list), 0)

	err = magmad.RemoveNetwork(testNetworkId)
	assert.NoError(t, err)

	assertNetworkTablesAreEmpty(t, testNetworkId)
}

func TestForceRemoveNetworkSuccess(t *testing.T) {
	mockeryStore := magmad_test_service.StartTestServiceMockStore(t)
	networkId := "NETWORK"

	_ = setupMockeryStoreForForceDeleteTests(t, networkId, mockeryStore)
	// Mock setup: both hw IDs exist
	mockeryStore.On("DoesKeyExist", servicers.GatewaysTableName, "gw1").
		Return(true, nil)
	mockeryStore.On("DoesKeyExist", servicers.GatewaysTableName, "gw2").
		Return(true, nil)
	mockeryStore.On("Delete", servicers.GatewaysTableName, "gw1").
		Return(nil)
	mockeryStore.On("Delete", servicers.GatewaysTableName, "gw2").
		Return(nil)

	// Mock setup: no errors on table drop
	mockeryStore.On("DeleteTable", mock.AnythingOfType("string")).
		Return(nil)

	// Mock setup: no error on network record delete
	mockeryStore.On("Delete", servicers.NetworksTableName, networkId).
		Return(nil)

	err := magmad.ForceRemoveNetwork(networkId)
	assert.NoError(t, err)
	mockeryStore.AssertNumberOfCalls(t, "DeleteTable", 5)
	mockeryStore.AssertExpectations(t)
}

func TestForceRemoveNetworkErrorHwIdQuery(t *testing.T) {
	mockeryStore := magmad_test_service.StartTestServiceMockStore(t)

	networkId := "NETWORK"
	mockeryStore.On("DoesKeyExist", servicers.NetworksTableName, "NETWORK").
		Return(true, nil)

	agRecordTableName := datastore.GetTableName(networkId, servicers.AgRecordTableName)
	mockeryStore.On("ListKeys", agRecordTableName).
		Return(nil, errors.New("ListKeys error"))
	err := magmad.ForceRemoveNetwork(networkId)
	assert.Error(t, err)
	assert.Contains(t, err.Error(),
		"Failed to query gateway hardware IDs in network, exiting before "+
			"performing any deletions. Error: ListKeys error")
	mockeryStore.AssertExpectations(t)
}

func TestForceRemoveNetworkErrorHwIdDelete(t *testing.T) {
	mockeryStore := magmad_test_service.StartTestServiceMockStore(t)
	networkId := "NETWORK"

	_ = setupMockeryStoreForForceDeleteTests(t, networkId, mockeryStore)

	// Mock setup for fetching hwId->nwId rows with error
	mockeryStore.On("DoesKeyExist", servicers.GatewaysTableName, "gw1").
		Return(true, errors.New("DoesKeyExist error"))
	mockeryStore.On("DoesKeyExist", servicers.GatewaysTableName, "gw2").
		Return(true, nil)
	mockeryStore.On("Delete", servicers.GatewaysTableName, "gw2").
		Return(errors.New("Delete error"))

	// Run test case, assert early exit without dropping tables or deleting network
	err := magmad.ForceRemoveNetwork(networkId)
	assert.Error(t, err)
	// golang maps have a randomized iteration order, so need to check string
	// contains for the error message instead of a naive equals
	assert.Contains(t, err.Error(), "Encountered the following errors while clearing hardware IDs:\n")
	assert.Contains(t, err.Error(), "\tError while checking if hardware ID gw1 exists: DoesKeyExist error\n")
	assert.Contains(t, err.Error(), "\tError while deleting hardware ID gw2: Delete error\n")
	assert.Contains(t, err.Error(), "Please retry the operation.")
	mockeryStore.AssertNotCalled(t, "DeleteTable", mock.AnythingOfType("string"))
	mockeryStore.AssertNotCalled(t, "Delete", servicers.NetworksTableName, mock.AnythingOfType("string"))
	mockeryStore.AssertExpectations(t)
}

func TestForceRemoveNetworkErrorDropTables(t *testing.T) {
	mockeryStore := magmad_test_service.StartTestServiceMockStore(t)
	networkId := "NETWORK"

	_ = setupMockeryStoreForForceDeleteTests(t, networkId, mockeryStore)
	// Mock setup: 1 hwID exists, the other is already cleared
	mockeryStore.On("DoesKeyExist", servicers.GatewaysTableName, "gw1").
		Return(true, nil)
	mockeryStore.On("DoesKeyExist", servicers.GatewaysTableName, "gw2").
		Return(false, nil)
	mockeryStore.On("Delete", servicers.GatewaysTableName, "gw1").
		Return(nil)

	// Mock setup: error on drop table for 2 tables
	mockeryStore.On("DeleteTable", datastore.GetTableName(networkId, servicers.HwIdTableName)).
		Return(errors.New("DeleteTable error 1"))
	mockeryStore.On("DeleteTable", mock.AnythingOfType("string")).
		Return(nil)

	// Assert early exit without deleting network
	err := magmad.ForceRemoveNetwork(networkId)
	assert.Error(t, err)
	assert.Contains(
		t,
		err.Error(),
		"Encountered the following errors while deleting the network:\n"+
			"\tError while deleting table NETWORK_hwIds: DeleteTable error 1\n",
	)
	mockeryStore.AssertNotCalled(t, "Delete", servicers.NetworksTableName, mock.AnythingOfType("string"))
	mockeryStore.AssertNumberOfCalls(t, "DeleteTable", 5)
	mockeryStore.AssertExpectations(t)
}

func TestNetworkConfig(t *testing.T) {
	magmad_test_service.StartTestService(t)

	testNetworkID := "magmad_test_network"
	testNetworkName := "Test Network Namee"
	testFeatures := map[string]string{"f1": "v1", "f2": "v2"}
	_, err := magmad.RegisterNetwork(&magmad_protos.MagmadNetworkRecord{Name: testNetworkName, Features: testFeatures}, testNetworkID)
	assert.NoError(t, err)

	networkRecord, err := magmad.GetNetwork(testNetworkID)
	assert.NoError(t, err)
	assert.NotNil(t, networkRecord)
	assert.Equal(t, networkRecord.Name, testNetworkName)
	assert.Equal(t, len(networkRecord.Features), 2)

	v, exists := networkRecord.Features["f1"]
	assert.True(t, exists)
	assert.NotNil(t, v)
	assert.Equal(t, v, "v1")
	v, exists = networkRecord.Features["f2"]
	assert.True(t, exists)
	assert.NotNil(t, v)
	assert.Equal(t, v, "v2")

	testNetworkNameUpdate := "Test Network Name"
	testFeaturesUpdate := map[string]string{"new-f1": "new-v1", "new-f2": "new-v2", "new-f3": "new-v3"}
	err = magmad.UpdateNetwork(testNetworkID, &magmad_protos.MagmadNetworkRecord{Name: testNetworkNameUpdate, Features: testFeaturesUpdate})
	assert.NoError(t, err)

	networkRecord, err = magmad.GetNetwork(testNetworkID)
	assert.NoError(t, err)
	assert.NotNil(t, networkRecord)
	assert.Equal(t, networkRecord.Name, testNetworkNameUpdate)
	assert.Equal(t, len(networkRecord.Features), 3)

	_, exists = networkRecord.Features["f1"]
	assert.False(t, exists)
	_, exists = networkRecord.Features["f2"]
	assert.False(t, exists)

	v, exists = networkRecord.Features["new-f1"]
	assert.True(t, exists)
	assert.NotNil(t, v)
	assert.Equal(t, v, "new-v1")
	v, exists = networkRecord.Features["new-f2"]
	assert.True(t, exists)
	assert.NotNil(t, v)
	assert.Equal(t, v, "new-v2")
	v, exists = networkRecord.Features["new-f3"]
	assert.True(t, exists)
	assert.NotNil(t, v)
	assert.Equal(t, v, "new-v3")

}

func setupMockeryStoreForForceDeleteTests(
	t *testing.T,
	networkId string,
	mockeryStore *mocks.Api,
) (mockAgRecordsTable map[string]datastore.ValueWrapper) {
	agRecords := []*magmad_protos.AccessGatewayRecord{
		{
			Name: "gw1",
			HwId: &protos.AccessGatewayID{Id: "gw1"},
		},
		{
			Name: "gw2",
			HwId: &protos.AccessGatewayID{Id: "gw2"},
		},
	}
	mockAgRecordsTable = make(map[string]datastore.ValueWrapper, len(agRecords))
	for _, agRecord := range agRecords {
		marshaled, err := protos.Marshal(agRecord)
		assert.NoError(t, err)
		mockAgRecordsTable[agRecord.GetName()] = datastore.ValueWrapper{Value: marshaled, Generation: 1}
	}

	mockeryStore.On("DoesKeyExist", servicers.NetworksTableName, "NETWORK").
		Return(true, nil)

	// Mock setup for getting hardware IDs
	agRecordTableName := datastore.GetTableName(networkId, servicers.AgRecordTableName)
	mockeryStore.On("ListKeys", agRecordTableName).
		Return([]string{"gw1", "gw2"}, nil)
	mockeryStore.On("GetMany", agRecordTableName, mock.AnythingOfType("[]string")).
		Return(mockAgRecordsTable, nil)
	return mockAgRecordsTable
}

func assertNetworkTablesAreEmpty(t *testing.T, testNetworkId string) {
	// Assert that everything is erased - skip checkin status and subscriberdb
	// table for now in order to isolate this test to magmad
	ds := test_utils.GetMockDatastoreInstance()
	tablesToCheck := []string{
		datastore.GetTableName(testNetworkId, servicers.AgRecordTableName),
		datastore.GetTableName(testNetworkId, servicers.HwIdTableName),
	}
	for _, tableName := range tablesToCheck {
		keys, err := ds.ListKeys(tableName)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(keys))
	}

	// Assert that all hardware IDs are gone from the network
	hwIds, err := ds.ListKeys(servicers.GatewaysTableName)
	assert.NoError(t, err)
	hwIdToNwIdMap, err := ds.GetMany(servicers.GatewaysTableName, hwIds)
	assert.NoError(t, err)
	for _, valueWrapper := range hwIdToNwIdMap {
		assert.NotEqual(t, testNetworkId, string(valueWrapper.Value))
	}
}
