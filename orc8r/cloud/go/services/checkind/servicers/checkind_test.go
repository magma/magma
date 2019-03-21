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
	checkind_test_init "magma/orc8r/cloud/go/services/checkind/test_init"
	logger_test_init "magma/orc8r/cloud/go/services/logger/test_init"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"

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
	magmad_test_init.StartTestService(t)
	checkind_test_init.StartTestService(t)
	logger_test_init.StartTestService(t)

	testNetworkId, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "Test Network Name"},
		"checkind_servicers_test_network")
	assert.NoError(t, err)

	t.Logf("New Registered Network: %s", testNetworkId)

	hwId := protos.AccessGatewayID{Id: testAgHwId}
	logicalId, err := magmad.RegisterGateway(testNetworkId,
		&magmad_protos.AccessGatewayRecord{HwId: &hwId, Name: "Test GW Name"})
	assert.NoError(t, err)
	assert.NotEqual(t, logicalId, "")

	testAgHwId2 := testAgHwId + "second"
	hwId = protos.AccessGatewayID{Id: testAgHwId2}
	logicalId2, err := magmad.RegisterGateway(
		testNetworkId, &magmad_protos.AccessGatewayRecord{HwId: &hwId, Name: "bla2"})
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

	request2 := protos.CheckinRequest{
		GatewayId: testAgHwId2,
		Status: &protos.ServiceStatus{
			Meta: map[string]string{
				"test": "meta",
			},
		},
		SystemStatus: &protos.SystemStatus{
			Time:         8,
			CpuUser:      7,
			CpuSystem:    6,
			CpuIdle:      5,
			MemTotal:     4,
			MemAvailable: 3,
			MemUsed:      2,
			MemFree:      1,
			UptimeSecs:   1234,
			SwapTotal:    1016081,
			SwapUsed:     54415,
			SwapFree:     412771,
			DiskPartitions: []*protos.DiskPartition{
				{
					Device:     "/dev/sda1",
					MountPoint: "/",
					Total:      1,
					Used:       2,
					Free:       3,
				},
			},
		},
		PlatformInfo: &protos.PlatformInfo{
			VpnIp: "facebook.com",
			Packages: []*protos.Package{
				{
					Name:    "magma",
					Version: "0.0.0.0",
				},
			},
			KernelVersion:           "42",
			KernelVersionsInstalled: []string{"42", "43"},
		},
		MachineInfo: &protos.MachineInfo{
			CpuInfo: &protos.CPUInfo{
				CoreCount:      4,
				ThreadsPerCore: 1,
				Architecture:   "x86_64",
				ModelName:      "Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz",
			},
			NetworkInfo: &protos.NetworkInfo{
				NetworkInterfaces: []*protos.NetworkInterface{
					{
						NetworkInterfaceId: "gtp_br0",
						Status:             protos.NetworkInterface_UP,
						MacAddress:         "08:00:27:1e:8a:32",
						IpAddresses:        []string{"10.10.10.1"},
						Ipv6Addresses:      []string{"fe80::a00:27ff:fe1e:8332"},
					},
				},
				RoutingTable: []*protos.Route{
					{
						DestinationIp:      "0.0.0.0",
						GatewayIp:          "10.10.10.1",
						Genmask:            "255.255.255.0",
						NetworkInterfaceId: "eth0",
					},
				},
			},
		},
	}

	_, readStatus := checkinRound(
		t, magmaCheckindClient, testNetworkId, logicalId, &request)

	time.Sleep(100 * time.Millisecond) // let some time pass

	_, readStatus2 := checkinRound(
		t, magmaCheckindClient, testNetworkId, logicalId2, &request2)
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
