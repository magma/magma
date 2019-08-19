/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package exporters_test

import (
	"os"
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/logger/exporters"
	"magma/orc8r/cloud/go/services/magmad"
	mdprotos "magma/orc8r/cloud/go/services/magmad/protos"
	magmadti "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/stretchr/testify/assert"
)

func TestScribeEntryUtilsLegacy_WithHWID(t *testing.T) {
	os.Setenv(orc8r.UseConfiguratorEnv, "0")
	magmadti.StartTestService(t)

	networkID := "test_network"
	gatewayID := "test_gateway"
	hwID := "test_hwID"
	_, err := magmad.RegisterNetwork(&mdprotos.MagmadNetworkRecord{}, networkID)
	assert.NoError(t, err)
	_, err = magmad.RegisterGatewayWithId(networkID, &mdprotos.AccessGatewayRecord{HwId: &protos.AccessGatewayID{Id: hwID}}, gatewayID)
	assert.NoError(t, err)

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

func TestScribeEntryUtilsLegacy(t *testing.T) {
	os.Setenv(orc8r.UseConfiguratorEnv, "0")
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
