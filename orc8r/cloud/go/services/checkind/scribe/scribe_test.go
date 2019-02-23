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

func TestFormatScribeGwStatusMessage(t *testing.T) {
	// input status
	systemStatus := protos.SystemStatus{CpuIdle: 3915306000, CpuSystem: 31921900, CpuUser: 74180510,
		MemAvailable: 3109638144, MemFree: 2648256512, MemTotal: 4137000960, MemUsed: 719056896}
	serviceStatus := protos.ServiceStatus{Meta: map[string]string{"enodeb_configured": "1",
		"enodeb_connected": "1", "gps_connected": "1", "gps_latitude": "37.484402", "gps_longitude": "-122.150044",
		"mme_connected": "1", "opstate_enabled": "1", "ptp_connected": "0", "rf_tx_on": "1"}}

	req := protos.CheckinRequest{GatewayId: "test hardware id", MagmaPkgVersion: "version 1", SystemStatus: &systemStatus, Status: &serviceStatus}
	status := protos.GatewayStatus{Time: 1511992464456, Checkin: &req}

	// convert gatewayStatus into ScribeMessage
	normalMsg, intMsg, err :=
		scribe.FormatScribeGwStatusMessage(&status, "test_networkId", "test_gatewayId")
	assert.NoError(t, err)
	assert.Equal(t, "1", normalMsg["enodeb_configured"])
	assert.Equal(t, "1", normalMsg["enodeb_connected"])
	assert.Equal(t, "1", normalMsg["gps_connected"])
	assert.Equal(t, "37.484402", normalMsg["gps_latitude"])
	assert.Equal(t, "-122.150044", normalMsg["gps_longitude"])
	assert.Equal(t, "1", normalMsg["mme_connected"])
	assert.Equal(t, "1", normalMsg["opstate_enabled"])
	assert.Equal(t, "0", normalMsg["ptp_connected"])
	assert.Equal(t, "1", normalMsg["rf_tx_on"])

	assert.Equal(t, int64(3915306000), intMsg["cpu_idle"])
	assert.Equal(t, int64(31921900), intMsg["cpu_system"])
	assert.Equal(t, int64(74180510), intMsg["cpu_user"])
	assert.Equal(t, int64(3109638144), intMsg["mem_available"])
	assert.Equal(t, int64(2648256512), intMsg["mem_free"])
	assert.Equal(t, int64(4137000960), intMsg["mem_total"])
	assert.Equal(t, int64(719056896), intMsg["mem_used"])

}
