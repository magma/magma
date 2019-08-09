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

	models2 "magma/orc8r/cloud/go/pluginimpl/models"
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

func GetGatewayStatusSwaggerFixture(gatewayId string) *models2.GatewayStatus {
	return &models2.GatewayStatus{
		CheckinTime:        0,
		CertExpirationTime: 0,
		Meta:               map[string]string{"hello": "world"},
		SystemStatus: &models2.SystemStatus{
			Time:       1495484735606,
			CPUUser:    31498,
			CPUSystem:  8361,
			CPUIDLE:    1869111,
			MemTotal:   1016084,
			MemUsed:    54416,
			MemFree:    412772,
			UptimeSecs: 1234,
			SwapTotal:  1016081,
			SwapUsed:   54415,
			SwapFree:   412771,
			DiskPartitions: []*models2.DiskPartition{
				{
					Device:     "/dev/sda1",
					MountPoint: "/",
					Total:      1,
					Used:       2,
					Free:       3,
				},
			},
		},
		PlatformInfo: &models2.PlatformInfo{
			VpnIP: "facebook.com",
			Packages: []*models2.Package{
				{
					Name:    "magma",
					Version: "0.0.0.0",
				},
			},
			KernelVersion:           "42",
			KernelVersionsInstalled: []string{"42", "43"},
			ConfigInfo: &models2.ConfigInfo{
				MconfigCreatedAt: 1552968732,
			},
		},
		MachineInfo: &models2.MachineInfo{
			CPUInfo: &models2.MachineInfoCPUInfo{
				CoreCount:      4,
				ThreadsPerCore: 1,
				Architecture:   "x86_64",
				ModelName:      "Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz",
			},
			NetworkInfo: &models2.MachineInfoNetworkInfo{
				NetworkInterfaces: []*models2.NetworkInterface{
					{
						NetworkInterfaceID: "gtp_br0",
						Status:             models2.NetworkInterfaceStatusUP,
						MacAddress:         "08:00:27:1e:8a:32",
						IPAddresses:        []string{"10.10.10.1"},
						IPV6Addresses:      []string{"fe80::a00:27ff:fe1e:8332"},
					},
				},
				RoutingTable: []*models2.Route{
					{
						DestinationIP:      "0.0.0.0",
						GatewayIP:          "10.10.10.1",
						Genmask:            "255.255.255.0",
						NetworkInterfaceID: "eth0",
					},
				},
			},
		},
		HardwareID:              gatewayId,
		KernelVersion:           "42",
		KernelVersionsInstalled: []string{"42", "43"},
		Version:                 "0.0.0.0",
		VpnIP:                   "facebook.com",
	}
}
