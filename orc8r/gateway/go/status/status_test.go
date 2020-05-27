/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package status_test

import (
	"encoding/json"
	"testing"

	"magma/gateway/status"

	"github.com/stretchr/testify/assert"
)

func TestGatewayStatus(t *testing.T) {
	gwStat := status.GetGatewayStatus()

	assert.NotNil(t, gwStat)
	assert.NotEmpty(t, gwStat.HardwareID)
	assert.NotNil(t, gwStat.MachineInfo)
	assert.NotNil(t, gwStat.PlatformInfo)
	assert.NotNil(t, gwStat.SystemStatus)

	jsoned, err := json.MarshalIndent(gwStat, " ", " ")
	assert.NoError(t, err)
	jsonedStr := string(jsoned)

	assert.Contains(t, jsonedStr, `"kernel_version":`)
	assert.Contains(t, jsonedStr, `"cpu_idle":`)
	assert.Contains(t, jsonedStr, `"cpu_user":`)
	assert.Contains(t, jsonedStr, `"cpu_system":`)
	assert.Contains(t, jsonedStr, `"device":`)
	assert.Contains(t, jsonedStr, `"mount_point":`)
	assert.Contains(t, jsonedStr, `"routing_table":`)
	assert.Contains(t, jsonedStr, `"network_interface_id":`)
	assert.Contains(t, jsonedStr, `"architecture":`)
	assert.Contains(t, jsonedStr, `"mem_available":`)
	assert.Contains(t, jsonedStr, `"uptime_secs":`)
}
