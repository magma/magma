/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"context"
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/services/checkind"

	"github.com/stretchr/testify/assert"
)

func Checkin(t *testing.T, req *protos.CheckinRequest) *protos.CheckinResponse {
	conn, err := registry.GetConnection(checkind.ServiceName)
	assert.NoError(t, err)
	client := protos.NewCheckindClient(conn)
	resp, err := client.Checkin(context.Background(), req)
	assert.NoError(t, err)
	return resp
}

func GetCheckinRequestProtoFixture(gatewayId string) *protos.CheckinRequest {
	return &protos.CheckinRequest{
		GatewayId: gatewayId,
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
			ConfigInfo: &protos.ConfigInfo{
				MconfigCreatedAt: 1552968732,
			},
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
		KernelVersionsInstalled: []string{},
	}
}

func GetGatewayStatusProtoFixture(gatewayId string) *protos.GatewayStatus {
	return &protos.GatewayStatus{
		Time:    0,
		Checkin: GetCheckinRequestProtoFixture(gatewayId),
	}
}
