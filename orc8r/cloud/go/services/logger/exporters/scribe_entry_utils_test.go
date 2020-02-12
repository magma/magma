/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package exporters_test

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/serde"
	configuratorti "magma/orc8r/cloud/go/services/configurator/test_init"
	configuratortu "magma/orc8r/cloud/go/services/configurator/test_utils"
	"magma/orc8r/cloud/go/services/device"
	deviceti "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/logger/exporters"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestScribeEntryUtils(t *testing.T) {
	logEntries := []*protos.LogEntry{
		{
			Category:  "test",
			NormalMap: map[string]string{"status": "ACTIVE"},
			IntMap:    map[string]int64{"port": 443},
			Time:      12345,
		},
	}
	scribeEntries, err := exporters.ConvertToScribeLogEntries(logEntries)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(scribeEntries))
	assert.Equal(t, logEntries[0].Category, scribeEntries[0].Category)
	expectedMsg := "{\"int\":{\"port\":443,\"time\":12345},\"normal\":{\"status\":\"ACTIVE\"}}"
	assert.Equal(t, expectedMsg, scribeEntries[0].Message)
}

func TestScribeEntryUtils_WithHWID(t *testing.T) {
	configuratorti.StartTestService(t)
	deviceti.StartTestService(t)
	_ = serde.RegisterSerdes(serde.NewBinarySerde(device.SerdeDomain, orc8r.AccessGatewayRecordType, &models.GatewayDevice{}))

	networkID := "test_network"
	gatewayID := "test_gateway"
	hwID := "test_hwID"
	configuratortu.RegisterNetwork(t, networkID, "")
	configuratortu.RegisterGateway(t, networkID, gatewayID, &models.GatewayDevice{HardwareID: hwID})

	logEntries := []*protos.LogEntry{
		{
			Category:  "test",
			NormalMap: map[string]string{"status": "ACTIVE"},
			IntMap:    map[string]int64{"port": 443},
			Time:      12345,
			HwId:      hwID,
		},
	}
	scribeEntries, err := exporters.ConvertToScribeLogEntries(logEntries)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(scribeEntries))
	assert.Equal(t, logEntries[0].Category, scribeEntries[0].Category)
	expectedMsg := "{\"int\":{\"port\":443,\"time\":12345},\"normal\":{\"gatewayId\":\"test_gateway\",\"networkId\":\"test_network\",\"status\":\"ACTIVE\"}}"
	assert.Equal(t, expectedMsg, scribeEntries[0].Message)
}
