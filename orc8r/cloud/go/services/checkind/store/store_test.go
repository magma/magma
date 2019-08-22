/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package store_test

import (
	"reflect"
	"testing"
	"time"

	"magma/orc8r/cloud/go/protos"
	checkinStore "magma/orc8r/cloud/go/services/checkind/store"
	checkinTestUtils "magma/orc8r/cloud/go/services/checkind/test_utils"
	loggerTestInit "magma/orc8r/cloud/go/services/logger/test_init"
	"magma/orc8r/cloud/go/services/magmad"
	magmadProtos "magma/orc8r/cloud/go/services/magmad/protos"
	magmadTestInit "magma/orc8r/cloud/go/services/magmad/test_init"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

const testAgHwId = "test_ag_HW_id"

func TestCheckinStore(t *testing.T) {
	magmadTestInit.StartTestService(t)
	loggerTestInit.StartTestService(t)
	testNetworkName := "Gateway Checkin Test Network"
	store, err := checkinStore.NewCheckinStore(test_utils.NewMockDatastore())
	assert.NoError(t, err)

	testNetworkId, err := magmad.RegisterNetwork(
		&magmadProtos.MagmadNetworkRecord{Name: testNetworkName},
		"checkind_store_test_network")
	assert.NoError(t, err)

	logicalId, err :=
		magmad.RegisterGateway(testNetworkId, &magmadProtos.AccessGatewayRecord{HwId: &protos.AccessGatewayID{Id: testAgHwId}})
	assert.NoError(t, err)
	assert.NotEqual(t, logicalId, "")

	status := checkinTestUtils.GetGatewayStatusProtoFixture(testAgHwId)

	err = store.UpdateGatewayStatus(status)
	assert.NoError(t, err)

	status_req := protos.GatewayStatusRequest{NetworkId: testNetworkId, LogicalId: logicalId}
	read_status, err := store.GetGatewayStatus(&status_req)
	assert.NoError(t, err)

	if !reflect.DeepEqual(status, read_status) {
		t.Errorf("GW Status Mismatch: %#v != %#v", *status.Checkin, *read_status.Checkin)
	}

	err = store.DeleteGatewayStatus(&status_req)
	assert.NoError(t, err)

	_, err = store.GetGatewayStatus(&status_req)
	// Error since the gateway is deleted
	assert.Error(t, err)

	status.Time = uint64(time.Now().UnixNano() / int64(time.Millisecond))

	err = store.UpdateGatewayStatus(status)
	assert.NoError(t, err)

	status_req = protos.GatewayStatusRequest{NetworkId: testNetworkId, LogicalId: logicalId}
	read_status, err = store.GetGatewayStatus(&status_req)
	assert.NoError(t, err)

	if !reflect.DeepEqual(status, read_status) {
		t.Errorf("GW Status Mismatch: %#v != %#v", *status, *read_status)
	}

	err = store.DeleteNetworkTable(testNetworkId)
	// Error since there is still the gateway table
	assert.Error(t, err)

	err = store.DeleteGatewayStatus(&status_req)
	assert.NoError(t, err)

	err = store.DeleteNetworkTable(testNetworkId)
	assert.NoError(t, err)

	_, err = store.GetGatewayStatus(&status_req)
	// Error since the network is deleted
	assert.Error(t, err)
}
