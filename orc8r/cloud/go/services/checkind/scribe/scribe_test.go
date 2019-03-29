/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package scribe_test

import (
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/checkind/scribe"

	"github.com/stretchr/testify/assert"
)

func TestFormatScribeOldGwStatusMessage(t *testing.T) {
	// input status
	systemStatus := protos.SystemStatus{CpuIdle: 3915306000, CpuSystem: 31921900, CpuUser: 74180510,
		MemAvailable: 3109638144, MemFree: 2648256512, MemTotal: 4137000960, MemUsed: 719056896, UptimeSecs: 12345}
	serviceStatus := protos.ServiceStatus{Meta: map[string]string{"enodeb_configured": "1",
		"enodeb_connected": "1", "gps_connected": "1", "gps_latitude": "37.484402", "gps_longitude": "-122.150044",
		"mme_connected": "1", "opstate_enabled": "1", "ptp_connected": "0", "rf_tx_on": "1"}}

	req := protos.CheckinRequest{
		GatewayId:               "test hardware id",
		MagmaPkgVersion:         "version 1",
		SystemStatus:            &systemStatus,
		Status:                  &serviceStatus,
		VpnIp:                   "facebook.com",
		KernelVersion:           "42",
		KernelVersionsInstalled: []string{"42", "43"},
	}
	status := protos.GatewayStatus{Time: 1511992464456, Checkin: &req, CertExpirationTime: 1000}

	// convert gatewayStatus into ScribeMessage
	normalMsg, intMsg, err := scribe.FormatScribeGwStatusMessage(&status, "test_networkId", "test_gatewayId")
	assert.NoError(t, err)

	expectedNormalMsg := map[string]string{
		"network_id":                "test_networkId",
		"gateway_id":                "test_gatewayId",
		"hardware_id":               "test hardware id",
		"enodeb_configured":         "1",
		"enodeb_connected":          "1",
		"gps_connected":             "1",
		"gps_latitude":              "37.484402",
		"gps_longitude":             "-122.150044",
		"mme_connected":             "1",
		"opstate_enabled":           "1",
		"ptp_connected":             "0",
		"rf_tx_on":                  "1",
		"magma_package_version":     "version 1",
		"vpn_ip":                    "facebook.com",
		"kernel_version":            "42",
		"kernel_versions_installed": "42,43",
	}
	assert.Equal(t, expectedNormalMsg, normalMsg)

	expectedIntMsg := map[string]int64{
		"cpu_idle":             3915306000,
		"cpu_system":           31921900,
		"cpu_user":             74180510,
		"mem_available":        3109638144,
		"mem_free":             2648256512,
		"mem_total":            4137000960,
		"mem_used":             719056896,
		"uptime_secs":          12345,
		"cert_expiration_time": 1000,
		"swap_total":           0,
		"swap_used":            0,
		"swap_free":            0,
	}
	assert.Equal(t, expectedIntMsg, intMsg)
}

func TestFormatScribeGwStatusMessage(t *testing.T) {
	// input status
	systemStatus := protos.SystemStatus{
		CpuIdle:      3915306000,
		CpuSystem:    31921900,
		CpuUser:      74180510,
		MemAvailable: 3109638144,
		MemFree:      2648256512,
		MemTotal:     4137000960,
		MemUsed:      719056896,
		UptimeSecs:   12345,
		SwapTotal:    111111,
		SwapUsed:     100000,
		SwapFree:     11111,
		DiskPartitions: []*protos.DiskPartition{
			{
				Device:     "/dev/sda1",
				MountPoint: "/",
				Total:      1,
				Used:       2,
				Free:       3,
			},
			{
				Device:     "/dev/sda2",
				MountPoint: "/test/",
				Total:      4,
				Used:       5,
				Free:       6,
			},
		},
	}
	serviceStatus := protos.ServiceStatus{Meta: map[string]string{"enodeb_configured": "1",
		"enodeb_connected": "1", "gps_connected": "1", "gps_latitude": "37.484402", "gps_longitude": "-122.150044",
		"mme_connected": "1", "opstate_enabled": "1", "ptp_connected": "0", "rf_tx_on": "1"}}
	platformInfo := protos.PlatformInfo{
		VpnIp: "facebook.com",
		Packages: []*protos.Package{
			{
				Name:    "magma",
				Version: "0.0.0.0",
			},
			{
				Name:    "test",
				Version: "0.0.0.1",
			},
		},
		KernelVersion:           "42",
		KernelVersionsInstalled: []string{"42", "43"},
	}
	machineInfo := protos.MachineInfo{
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
				{
					NetworkInterfaceId: "eth0",
					Status:             protos.NetworkInterface_DOWN,
					MacAddress:         "08:00:27:1e:8a:33",
					IpAddresses:        []string{"10.10.10.10"},
					Ipv6Addresses:      []string{"fe80::a00:27ff:fe1e:8333"},
				},
			},
			RoutingTable: []*protos.Route{
				{
					DestinationIp:      "0.0.0.0",
					GatewayIp:          "10.10.10.1",
					Genmask:            "255.255.255.0",
					NetworkInterfaceId: "eth0",
				},
				{
					DestinationIp:      "10.0.0.0",
					GatewayIp:          "10.10.10.2",
					Genmask:            "255.255.0.0",
					NetworkInterfaceId: "eth1",
				},
			},
		},
	}

	req := protos.CheckinRequest{
		GatewayId:    "test hardware id",
		SystemStatus: &systemStatus,
		Status:       &serviceStatus,
		PlatformInfo: &platformInfo,
		MachineInfo:  &machineInfo,
	}
	status := protos.GatewayStatus{Time: 1511992464456, Checkin: &req, CertExpirationTime: 1000}

	// convert gatewayStatus into ScribeMessage
	normalMsg, intMsg, err := scribe.FormatScribeGwStatusMessage(&status, "test_networkId", "test_gatewayId")
	assert.NoError(t, err)

	expectedNormalMsg := map[string]string{
		"network_id":                "test_networkId",
		"gateway_id":                "test_gatewayId",
		"hardware_id":               "test hardware id",
		"enodeb_configured":         "1",
		"enodeb_connected":          "1",
		"gps_connected":             "1",
		"gps_latitude":              "37.484402",
		"gps_longitude":             "-122.150044",
		"mme_connected":             "1",
		"opstate_enabled":           "1",
		"ptp_connected":             "0",
		"rf_tx_on":                  "1",
		"magma_package_version":     "0.0.0.0",
		"vpn_ip":                    "facebook.com",
		"kernel_version":            "42",
		"kernel_versions_installed": "42,43",
		"disk_partitions": "[{\"device\":\"/dev/sda1\",\"mount_point\":\"/\",\"total\":1,\"used\":2,\"free\":3}," +
			"{\"device\":\"/dev/sda2\",\"mount_point\":\"/test/\",\"total\":4,\"used\":5,\"free\":6}]",
		"platform_info.vpn_ip":                    "facebook.com",
		"platform_info.packages":                  "[{\"name\":\"magma\",\"version\":\"0.0.0.0\"},{\"name\":\"test\",\"version\":\"0.0.0.1\"}]",
		"platform_info.kernel_version":            "42",
		"platform_info.kernel_versions_installed": "[\"42\",\"43\"]",
		"machine_info.cpu_info.architecture":      "x86_64",
		"machine_info.cpu_info.model_name":        "Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz",
		"machine_info.network_info.network_interfaces": "[{\"network_interface_id\":\"gtp_br0\",\"status\":1,\"mac_address\":" +
			"\"08:00:27:1e:8a:32\",\"ip_addresses\":[\"10.10.10.1\"],\"ipv6_addresses\":[\"fe80::a00:27ff:fe1e:8332\"]}," +
			"{\"network_interface_id\":\"eth0\",\"status\":2,\"mac_address\":\"08:00:27:1e:8a:33\",\"ip_addresses\":" +
			"[\"10.10.10.10\"],\"ipv6_addresses\":[\"fe80::a00:27ff:fe1e:8333\"]}]",
		"machine_info.network_info.routing_table": "[{\"destination_ip\":\"0.0.0.0\",\"gateway_ip\":\"10.10.10.1\"," +
			"\"genmask\":\"255.255.255.0\",\"network_interface_id\":\"eth0\"},{\"destination_ip\":\"10.0.0.0\",\"gateway_ip\":" +
			"\"10.10.10.2\",\"genmask\":\"255.255.0.0\",\"network_interface_id\":\"eth1\"}]",
	}
	assert.Equal(t, expectedNormalMsg, normalMsg)

	expectedIntMsg := map[string]int64{
		"cpu_idle":                               3915306000,
		"cpu_system":                             31921900,
		"cpu_user":                               74180510,
		"mem_available":                          3109638144,
		"mem_free":                               2648256512,
		"mem_total":                              4137000960,
		"mem_used":                               719056896,
		"uptime_secs":                            12345,
		"cert_expiration_time":                   1000,
		"swap_total":                             111111,
		"swap_used":                              100000,
		"swap_free":                              11111,
		"machine_info.cpu_info.core_count":       4,
		"machine_info.cpu_info.threads_per_core": 1,
	}
	assert.Equal(t, expectedIntMsg, intMsg)
}
