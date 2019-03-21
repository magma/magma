/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package checkind_test

import (
	"context"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/services/checkind"
	"magma/orc8r/cloud/go/services/checkind/test_init"
	"magma/orc8r/cloud/go/services/magmad"
	mdprotos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/stretchr/testify/assert"
)

const (
	gw1HardwareID = "11ffea10-7fc4-4427-975a-b9e4ce8f6f4d"
	gw2HardwareID = "11ffea10-7fc4-4427-975a-b9e4ce8f6f4e"
)

func TestCheckinAPI(t *testing.T) {
	// Initialize test services
	magmad_test_init.StartTestService(t)
	test_init.StartTestService(t)

	// Register gateways and perform mock checkins
	registerGateways(t)
	checkinRequests := insertStatuses(t)

	// Test API
	ids, err := checkind.List("net1")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"gw1", "gw2"}, ids.Ids)

	gw1Status, err := checkind.GetStatus("net1", "gw1")
	assert.NoError(t, err)
	assert.Equal(t, protos.TestMarshal(checkinRequests["gw1"]), protos.TestMarshal(gw1Status.Checkin))

	gw2Status, err := checkind.GetStatus("net1", "gw2")
	assert.NoError(t, err)
	assert.Equal(t, protos.TestMarshal(checkinRequests["gw2"]), protos.TestMarshal(gw2Status.Checkin))

	err = checkind.DeleteNetwork("net1")
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "Status table for network net1 is not empty"))

	err = checkind.DeleteGatewayStatus("net1", "gw1")
	assert.NoError(t, err)

	_, err = checkind.GetStatus("net1", "gw1")
	assert.Error(t, err)
	assert.Equal(t, errors.ErrNotFound, err)

	err = checkind.DeleteGatewayStatus("net1", "gw2")
	assert.NoError(t, err)

	err = checkind.DeleteNetwork("net1")
	assert.NoError(t, err)

	ids, err = checkind.List("net1")
	assert.NoError(t, err)
	assert.Nil(t, ids.Ids)
}

func registerGateways(t *testing.T) {
	net1ID, err := magmad.RegisterNetwork(&mdprotos.MagmadNetworkRecord{Name: "Network 1"}, "net1")
	assert.NoError(t, err)
	assert.Equal(t, "net1", net1ID)
	gw1Record := &mdprotos.AccessGatewayRecord{
		HwId: &protos.AccessGatewayID{Id: gw1HardwareID},
		Name: "Gateway 1",
		Key: &protos.ChallengeKey{
			KeyType: protos.ChallengeKey_ECHO,
		},
	}
	gw1ID, err := magmad.RegisterGatewayWithId("net1", gw1Record, "gw1")
	assert.NoError(t, err)
	assert.Equal(t, "gw1", gw1ID)
	gw2Record := &mdprotos.AccessGatewayRecord{
		HwId: &protos.AccessGatewayID{Id: gw2HardwareID},
		Name: "Gateway 2",
		Key: &protos.ChallengeKey{
			KeyType: protos.ChallengeKey_ECHO,
		},
	}
	gw2ID, err := magmad.RegisterGatewayWithId("net1", gw2Record, "gw2")
	assert.NoError(t, err)
	assert.Equal(t, "gw2", gw2ID)
}

func insertStatuses(t *testing.T) map[string]*protos.CheckinRequest {
	conn, err := registry.GetConnection(checkind.ServiceName)
	assert.NoError(t, err)
	defer conn.Close()
	client := protos.NewCheckindClient(conn)
	checkinRequest1 := &protos.CheckinRequest{
		GatewayId: gw1HardwareID,
		Status: &protos.ServiceStatus{
			Meta: map[string]string{
				"hello": "world",
			},
		},
		SystemStatus: &protos.SystemStatus{
			Time:       1495484735606,
			CpuUser:    31498,
			CpuSystem:  8361,
			CpuIdle:    1869111,
			MemTotal:   1016084,
			MemUsed:    54416,
			MemFree:    412772,
			UptimeSecs: 1234,
			SwapTotal:  1016081,
			SwapUsed:   54415,
			SwapFree:   412771,
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
	_, err = client.Checkin(context.Background(), checkinRequest1)
	assert.NoError(t, err)

	checkinRequest2 := &protos.CheckinRequest{
		GatewayId: gw2HardwareID,
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
	_, err = client.Checkin(context.Background(), checkinRequest2)
	assert.NoError(t, err)

	return map[string]*protos.CheckinRequest{
		"gw1": checkinRequest1,
		"gw2": checkinRequest2,
	}
}
