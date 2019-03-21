/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package models_test

import (
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/checkind/obsidian/models"

	"github.com/stretchr/testify/assert"
)

func TestGatewayStatus_FromMconfig(t *testing.T) {
	testCases := []struct {
		In  *protos.GatewayStatus
		Out *models.GatewayStatus
	}{
		// Check old status
		{
			In: &protos.GatewayStatus{
				Time: 42,
				Checkin: &protos.CheckinRequest{
					GatewayId:       "foo",
					MagmaPkgVersion: "bar",
					Status: &protos.ServiceStatus{
						Meta: map[string]string{"baz": "qux"},
					},
					SystemStatus: &protos.SystemStatus{
						Time:         101,
						CpuUser:      10,
						CpuSystem:    11,
						CpuIdle:      12,
						MemTotal:     13,
						MemAvailable: 14,
						MemUsed:      15,
						MemFree:      16,
						UptimeSecs:   17,
					},
					VpnIp:                   "facebook.com",
					KernelVersion:           "42",
					KernelVersionsInstalled: []string{"11"},
				},
			},
			Out: &models.GatewayStatus{
				CheckinTime:             42,
				HardwareID:              "foo",
				KernelVersion:           "42",
				KernelVersionsInstalled: []string{"11"},
				Meta:                    map[string]string{"baz": "qux"},
				SystemStatus: &models.SystemStatus{
					CPUIDLE:      12,
					CPUSystem:    11,
					CPUUser:      10,
					MemAvailable: 14,
					MemFree:      16,
					MemTotal:     13,
					MemUsed:      15,
					Time:         101,
					UptimeSecs:   17,
				},
				PlatformInfo: &models.PlatformInfo{
					VpnIP: "facebook.com",
					Packages: []*models.Package{
						{
							Name:    "magma",
							Version: "bar",
						},
					},
					KernelVersion:           "42",
					KernelVersionsInstalled: []string{"11"},
				},
				Version: "bar",
				VpnIP:   "facebook.com",
			},
		},

		// Check new status
		{
			In: &protos.GatewayStatus{
				Time: 42,
				Checkin: &protos.CheckinRequest{
					GatewayId: "foo",
					Status: &protos.ServiceStatus{
						Meta: map[string]string{"baz": "qux"},
					},
					SystemStatus: &protos.SystemStatus{
						Time:         101,
						CpuUser:      10,
						CpuSystem:    11,
						CpuIdle:      12,
						MemTotal:     13,
						MemAvailable: 14,
						MemUsed:      15,
						MemFree:      16,
						UptimeSecs:   17,
						SwapTotal:    18,
						SwapUsed:     19,
						SwapFree:     20,
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
				},
			},
			Out: &models.GatewayStatus{
				CheckinTime: 42,
				Meta:        map[string]string{"baz": "qux"},
				SystemStatus: &models.SystemStatus{
					CPUIDLE:      12,
					CPUSystem:    11,
					CPUUser:      10,
					MemAvailable: 14,
					MemFree:      16,
					MemTotal:     13,
					MemUsed:      15,
					Time:         101,
					UptimeSecs:   17,
					SwapTotal:    18,
					SwapUsed:     19,
					SwapFree:     20,
					DiskPartitions: []*models.DiskPartition{
						{
							Device:     "/dev/sda1",
							MountPoint: "/",
							Total:      1,
							Used:       2,
							Free:       3,
						},
					},
				},
				PlatformInfo: &models.PlatformInfo{
					VpnIP: "facebook.com",
					Packages: []*models.Package{
						{
							Name:    "magma",
							Version: "0.0.0.0",
						},
					},
					KernelVersion:           "42",
					KernelVersionsInstalled: []string{"42", "43"},
				},
				MachineInfo: &models.MachineInfo{
					CPUInfo: &models.MachineInfoCPUInfo{
						CoreCount:      4,
						ThreadsPerCore: 1,
						Architecture:   "x86_64",
						ModelName:      "Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz",
					},
					NetworkInfo: &models.MachineInfoNetworkInfo{
						NetworkInterfaces: []*models.NetworkInterface{
							{
								NetworkInterfaceID: "gtp_br0",
								Status:             models.NetworkInterfaceStatusUP,
								MacAddress:         "08:00:27:1e:8a:32",
								IPAddresses:        []string{"10.10.10.1"},
								IPV6Addresses:      []string{"fe80::a00:27ff:fe1e:8332"},
							},
						},
						RoutingTable: []*models.Route{
							{
								DestinationIP:      "0.0.0.0",
								GatewayIP:          "10.10.10.1",
								Genmask:            "255.255.255.0",
								NetworkInterfaceID: "eth0",
							},
						},
					},
				},
				HardwareID:              "foo",
				KernelVersion:           "42",
				KernelVersionsInstalled: []string{"42", "43"},
				Version:                 "0.0.0.0",
				VpnIP:                   "facebook.com",
			},
		},

		// Nil status from checkin
		{
			In: &protos.GatewayStatus{
				Time: 42,
			},
			Out: &models.GatewayStatus{
				CheckinTime: 42,
				Meta:        map[string]string{},
			},
		},
	}

	for _, tc := range testCases {
		actual := &models.GatewayStatus{}
		err := actual.FromMconfig(tc.In)
		assert.NoError(t, err)
		assert.Equal(t, tc.Out, actual)
	}
}
