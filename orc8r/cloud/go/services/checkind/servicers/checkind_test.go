/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"testing"
	"time"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/services/checkind"
	"magma/orc8r/cloud/go/services/checkind/test_init"
	"magma/orc8r/cloud/go/services/checkind/test_utils"
	loggerTestInit "magma/orc8r/cloud/go/services/logger/test_init"
	"magma/orc8r/cloud/go/services/magmad"
	magmadProtos "magma/orc8r/cloud/go/services/magmad/protos"
	magmadTestInit "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const testAgHwId = "Test-AGW-Hw-Id"

func checkinRound(t *testing.T,
	magmaCheckindClient protos.CheckindClient,
	testNetworkId string,
	logicalId string,
	requestPtr *protos.CheckinRequest) (*protos.CheckinResponse,
	*protos.GatewayStatus) {

	resp, err := magmaCheckindClient.Checkin(context.Background(), requestPtr)
	assert.NoError(t, err)
	assert.Equal(t, resp.Action, protos.CheckinResponse_NONE)

	statusReq := protos.GatewayStatusRequest{NetworkId: testNetworkId, LogicalId: logicalId}
	readStatus, err := magmaCheckindClient.GetStatus(
		context.Background(), &statusReq)
	assert.NoError(t, err)

	if protos.TestMarshal(requestPtr) != protos.TestMarshal(readStatus.Checkin) {
		t.Fatalf(
			"GW Status Mismatch: %#v != %#v", *requestPtr, *readStatus.Checkin)
	}
	timeDiff :=
		uint64(time.Now().UnixNano()*
			int64(time.Nanosecond)/int64(time.Millisecond)) - readStatus.Time
	assert.True(t, timeDiff >= 0 && timeDiff < 1000) // 1 second should be enough

	return resp, readStatus
}

func TestCheckind(t *testing.T) {
	magmadTestInit.StartTestService(t)
	test_init.StartTestService(t)
	loggerTestInit.StartTestService(t)

	testNetworkId, err := magmad.RegisterNetwork(
		&magmadProtos.MagmadNetworkRecord{Name: "Test Network Name"},
		"checkind_servicers_test_network")
	assert.NoError(t, err)

	t.Logf("New Registered Network: %s", testNetworkId)

	hwId := protos.AccessGatewayID{Id: testAgHwId}
	logicalId, err := magmad.RegisterGateway(testNetworkId,
		&magmadProtos.AccessGatewayRecord{HwId: &hwId, Name: "Test GW Name"})
	assert.NoError(t, err)
	assert.NotEqual(t, logicalId, "")

	testAgHwId2 := testAgHwId + "second"
	hwId = protos.AccessGatewayID{Id: testAgHwId2}
	logicalId2, err := magmad.RegisterGateway(
		testNetworkId, &magmadProtos.AccessGatewayRecord{HwId: &hwId, Name: "bla2"})
	assert.NoError(t, err)
	assert.NotEqual(t, logicalId2, "")

	conn, err := registry.GetConnection(checkind.ServiceName)
	assert.NoError(t, err)

	magmaCheckindClient := protos.NewCheckindClient(conn)

	// Test GW updating status
	request := protos.CheckinRequest{
		GatewayId:       testAgHwId,
		MagmaPkgVersion: "1.2.3",
		Status: &protos.ServiceStatus{
			Meta: map[string]string{
				"hello": "world",
			},
		},
		SystemStatus: &protos.SystemStatus{
			CpuUser:   31498,
			CpuSystem: 8361,
			CpuIdle:   1869111,
			MemTotal:  1016084,
			MemUsed:   54416,
			MemFree:   412772,
		},
	}

	request2 := test_utils.GetCheckinRequestProtoFixture(testAgHwId2)

	_, readStatus := checkinRound(
		t, magmaCheckindClient, testNetworkId, logicalId, &request)

	time.Sleep(100 * time.Millisecond) // let some time pass

	_, readStatus2 := checkinRound(
		t, magmaCheckindClient, testNetworkId, logicalId2, request2)
	assert.NotEqual(t, readStatus.Time, readStatus2.Time)

	repeatRequest := proto.Clone(&request).(*protos.CheckinRequest)
	repeatRequest.SystemStatus.CpuSystem = 9876
	_, repeatStatus := checkinRound(
		t, magmaCheckindClient, testNetworkId, logicalId, repeatRequest)
	assert.True(t, readStatus.Time < repeatStatus.Time)

	assert.Equal(t, repeatStatus.Checkin.SystemStatus.CpuSystem,
		repeatRequest.SystemStatus.CpuSystem)

	lids, err := magmaCheckindClient.List(
		context.Background(), &protos.NetworkID{Id: testNetworkId})
	assert.NoError(t, err)
	assert.Equal(t, len(lids.Ids), 2)
	assert.True(t, (lids.Ids[0] == logicalId || lids.Ids[1] == logicalId))
	assert.True(t, (lids.Ids[0] == logicalId2 || lids.Ids[1] == logicalId2))
}
