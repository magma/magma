/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package directoryd_test

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	models2 "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/serde"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	configuratorTestUtils "magma/orc8r/cloud/go/services/configurator/test_utils"
	"magma/orc8r/cloud/go/services/device"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/directoryd"
	directoryd_test_init "magma/orc8r/cloud/go/services/directoryd/test_init"
	"magma/orc8r/cloud/go/services/state"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/stretchr/testify/assert"
)

const (
	testAgHwId = "Test-AGW-Hw-Id"
	testGwId1  = "gw1"
	testGwId2  = "gw2"
	testSubId1 = "sub1"
	testSubId2 = "sub2"
	testSubId3 = "sub3"
)

func TestDirectorydClientMethods(t *testing.T) {
	directoryd_test_init.StartTestService(t)
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	// Set up test networkID, hwID, and encode into context
	stateTestInit.StartTestService(t)
	err := serde.RegisterSerdes(
		state.NewStateSerde(orc8r.DirectoryRecordType, &directoryd.DirectoryRecord{}),
		serde.NewBinarySerde(device.SerdeDomain, orc8r.AccessGatewayRecordType, &models2.GatewayDevice{}))
	assert.NoError(t, err)

	networkID := "directoryd_service_test_network"
	configuratorTestUtils.RegisterNetwork(t, networkID, "DirectoryD Service Test")
	gatewayID := testAgHwId
	configuratorTestUtils.RegisterGateway(t, networkID, gatewayID, &models2.GatewayDevice{HardwareID: testAgHwId})
	ctx := test_utils.GetContextWithCertificate(t, testAgHwId)

	stateClient, err := getStateServiceClient(t)
	record := &directoryd.DirectoryRecord{
		LocationHistory: []string{testAgHwId},
	}
	serializedRecord, err := record.MarshalBinary()
	assert.NoError(t, err)
	state := &protos.State{
		Type:     orc8r.DirectoryRecordType,
		DeviceID: testSubId1,
		Value:    serializedRecord,
	}
	req := &protos.ReportStatesRequest{
		States: []*protos.State{state},
	}
	response, err := stateClient.ReportStates(ctx, req)
	assert.NoError(t, err)
	assert.Empty(t, response.UnreportedStates)

	hwID, err := directoryd.GetHardwareIdByIMSI(testSubId1, networkID)
	assert.NoError(t, err)
	assert.Equal(t, testAgHwId, hwID)

	err = directoryd.DeleteHardwareIdByIMSI(testSubId1, networkID)
	assert.NoError(t, err)
	_, err = directoryd.GetHardwareIdByIMSI(testSubId1, networkID)
	assert.Error(t, err)
}

func TestDirectoryDLegacyClient(t *testing.T) {
	// Get empty DB
	_, err := directoryd.GetHardwareIdByIMSI(testSubId1, "")
	assert.EqualError(t, err, "rpc error: code = Unknown desc = Error getting location record: No record for query")

	// Repeat using other table

	// Get empty DB
	_, err = directoryd.GetHostNameByIMSI(testSubId1)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = Error getting location record: No record for query")

	// Add two locations
	err = directoryd.UpdateHostNameByHwId(testSubId1, testGwId1)
	assert.NoError(t, err)

	err = directoryd.UpdateHostNameByHwId(testSubId2, testGwId2)
	assert.NoError(t, err)

	// Read back
	record, err := directoryd.GetHostNameByIMSI(testSubId1)
	assert.NoError(t, err)
	assert.Equal(t, testGwId1, record)

	record, err = directoryd.GetHostNameByIMSI(testSubId2)
	assert.NoError(t, err)
	assert.Equal(t, testGwId2, record)

	record, err = directoryd.GetHostNameByIMSI(testSubId3)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = Error getting location record: No record for query")

	// Delete
	err = directoryd.DeleteHostNameByIMSI(testSubId1)
	assert.NoError(t, err)

	record, err = directoryd.GetHostNameByIMSI(testSubId1)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = Error getting location record: No record for query")

	// Delete unknown
	err = directoryd.DeleteHostNameByIMSI(testSubId3)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = Error finding location record: No record for query")
}

func getStateServiceClient(t *testing.T) (protos.StateServiceClient, error) {
	conn, err := registry.GetConnection(state.ServiceName)
	assert.NoError(t, err)
	return protos.NewStateServiceClient(conn), err
}
