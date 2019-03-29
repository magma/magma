/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers_test

import (
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/services/checkind"
	checkindTestInit "magma/orc8r/cloud/go/services/checkind/test_init"
	"magma/orc8r/cloud/go/services/checkind/test_utils"
	"magma/orc8r/cloud/go/services/magmad"
	magmadProtos "magma/orc8r/cloud/go/services/magmad/protos"
	magmadTestInit "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const testAgHwId = "Test-AGW-Hw-Id"

// TestCheckind is Obsidian Gateway Status Integration Test intended to be run
// on cloud VM
func TestCheckind(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	magmadTestInit.StartTestService(t)
	checkindTestInit.StartTestService(t)
	restPort := tests.StartObsidian(t)

	// create a test network with a single GW
	testNetworkId, err := magmad.RegisterNetwork(
		&magmadProtos.MagmadNetworkRecord{Name: "Test Network 1"},
		"checkind_obsidian_test_network")
	assert.NoError(t, err)

	t.Logf("New Registered Network: %s", testNetworkId)

	hwId := protos.AccessGatewayID{Id: testAgHwId}
	logicalId, err := magmad.RegisterGateway(testNetworkId, &magmadProtos.AccessGatewayRecord{HwId: &hwId, Name: "Test GW Name"})
	assert.NoError(t, err)
	assert.NotEqual(t, logicalId, "")

	conn, err := registry.GetConnection(checkind.ServiceName)
	assert.NoError(t, err)

	magmaCheckindClient := protos.NewCheckindClient(conn)

	// Test GW updating old status
	request := &protos.CheckinRequest{
		GatewayId:       testAgHwId,
		MagmaPkgVersion: "1.2.3",
		KernelVersion:   "4.9.0-6-amd64",
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
		},
	}
	resp, err := magmaCheckindClient.Checkin(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, resp.Action, protos.CheckinResponse_NONE)

	// Get GW status from the checkin above via REST API
	testUrlRoot := fmt.Sprintf(
		"http://localhost:%d%s/networks", restPort, handlers.REST_ROOT)
	getOldAGStatusTestCase := tests.Testcase{
		Name:   "Get Old AG Status",
		Method: "GET",
		Url: fmt.Sprintf("%s/%s/gateways/%s/status",
			testUrlRoot, testNetworkId, logicalId),
		Payload: "",
		Expected: fmt.Sprintf(`
			{
			   "checkin_time":%d,
			   "hardware_id":"Test-AGW-Hw-Id",
			   "kernel_version":"4.9.0-6-amd64",
			   "meta":{
				  "hello":"world"
			   },
			   "platform_info":{
				  "kernel_version":"4.9.0-6-amd64",
				  "packages":[
					 {
						"name":"magma",
						"version":"1.2.3"
					 }
				  ]
			   },
			   "system_status":{
				  "cpu_idle":1869111,
				  "cpu_system":8361,
				  "cpu_user":31498,
				  "mem_free":412772,
				  "mem_total":1016084,
				  "mem_used":54416,
				  "time":1495484735606,
				  "uptime_secs":1234
			   },
			   "version":"1.2.3"
			}`, resp.Time),
	}
	tests.RunTest(t, getOldAGStatusTestCase)

	// Test GW updating status
	request = test_utils.GetCheckinRequestProtoFixture(testAgHwId)
	resp, err = magmaCheckindClient.Checkin(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, resp.Action, protos.CheckinResponse_NONE)

	// Get GW status from the checkin above via REST API
	getAGStatusTestCase := tests.Testcase{
		Name:   "Get AG Status",
		Method: "GET",
		Url: fmt.Sprintf("%s/%s/gateways/%s/status",
			testUrlRoot, testNetworkId, logicalId),
		Payload: "",
		Expected: fmt.Sprintf(`
			{
			   "checkin_time":%d,
			   "hardware_id":"Test-AGW-Hw-Id",
			   "kernel_version":"42",
			   "kernel_versions_installed":[
				  "42",
				  "43"
			   ],
			   "machine_info":{
				  "cpu_info":{
					 "architecture":"x86_64",
					 "core_count":4,
					 "model_name":"Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz",
					 "threads_per_core":1
				  },
				  "network_info":{
					 "network_interfaces":[
						{
						   "ip_addresses":[
							  "10.10.10.1"
						   ],
						   "ipv6_addresses":[
							  "fe80::a00:27ff:fe1e:8332"
						   ],
						   "mac_address":"08:00:27:1e:8a:32",
						   "network_interface_id":"gtp_br0",
						   "status":"UP"
						}
					 ],
					 "routing_table":[
						{
						   "destination_ip":"0.0.0.0",
						   "gateway_ip":"10.10.10.1",
						   "genmask":"255.255.255.0",
						   "network_interface_id":"eth0"
						}
					 ]
				  }
			   },
			   "meta":{
				  "hello":"world"
			   },
			   "platform_info":{
				  "kernel_version":"42",
				  "kernel_versions_installed":[
					 "42",
					 "43"
				  ],
				  "packages":[
					 {
						"name":"magma",
						"version":"0.0.0.0"
					 }
				  ],
				  "vpn_ip":"facebook.com"
			   },
			   "system_status":{
				  "cpu_idle":1869111,
				  "cpu_system":8361,
				  "cpu_user":31498,
				  "disk_partitions":[
					 {
						"device":"/dev/sda1",
						"free":3,
						"mount_point":"/",
						"total":1,
						"used":2
					 }
				  ],
				  "mem_free":412772,
				  "mem_total":1016084,
				  "mem_used":54416,
				  "swap_free":412771,
				  "swap_total":1016081,
				  "swap_used":54415,
				  "time":1495484735606,
				  "uptime_secs":1234
			   },
			   "version":"0.0.0.0",
			   "vpn_ip":"facebook.com"
			}`, resp.Time),
	}
	tests.RunTest(t, getAGStatusTestCase)

	err = magmad.ForceRemoveNetwork(testNetworkId)
	assert.NoError(t, err)
}
