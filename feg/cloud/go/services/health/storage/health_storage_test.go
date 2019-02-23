/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage_test

import (
	"testing"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/health/storage"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

const TestFeGNetworkID = "test_networkID1"
const TestFegLogicalID1 = "test_logicalID1"
const TestFegLogicalID2 = "test_logicalID2"

func TestHealthStorage_GetHealth(t *testing.T) {
	store, err := storage.NewHealthStore(test_utils.NewMockDatastore())
	assert.NoError(t, err)

	_, err = store.GetHealth(TestFeGNetworkID, TestFegLogicalID1)
	assert.Error(t, err)

	health1 := getHealthStats()

	err = store.UpdateHealth(TestFeGNetworkID, TestFegLogicalID1, health1)
	assert.NoError(t, err)

	res1, err := store.GetHealth(TestFeGNetworkID, TestFegLogicalID1)
	assert.NoError(t, err)
	assert.Equal(t, health1, res1)

	// Add another gateway to the cluster table
	health2 := getHealthStats()
	health2.SystemStatus.MemAvailableBytes = 1000000

	err = store.UpdateHealth(TestFeGNetworkID, TestFegLogicalID2, health2)
	assert.NoError(t, err)

	res2, err := store.GetHealth(TestFeGNetworkID, TestFegLogicalID2)
	assert.NoError(t, err)
	assert.Equal(t, health2, res2)
}

func TestHealthStorage_UpdateHealth(t *testing.T) {
	store, err := storage.NewHealthStore(test_utils.NewMockDatastore())
	assert.NoError(t, err)

	health := getHealthStats()

	err = store.UpdateHealth(TestFeGNetworkID, TestFegLogicalID1, health)
	assert.NoError(t, err)

	res1, err := store.GetHealth(TestFeGNetworkID, TestFegLogicalID1)
	assert.NoError(t, err)
	assert.Equal(t, health, res1)

	// Modify the health status before update
	health.SystemStatus.CpuUtilPct = 0.90
	health.ServiceStatus["SESSION_PROXY"].ServiceHealthStatus.Health = protos.HealthStatus_UNHEALTHY

	err = store.UpdateHealth(TestFeGNetworkID, TestFegLogicalID2, health)
	assert.NoError(t, err)

	res2, err := store.GetHealth(TestFeGNetworkID, TestFegLogicalID2)
	assert.NoError(t, err)
	assert.Equal(t, health, res2)
}

func getHealthStats() *protos.HealthStats {

	serviceStats := protos.ServiceHealthStats{
		ServiceState: protos.ServiceHealthStats_AVAILABLE,
		ServiceHealthStatus: &protos.HealthStatus{
			Health: protos.HealthStatus_HEALTHY,
		},
	}

	serviceStatsMap := make(map[string]*protos.ServiceHealthStats)
	serviceStatsMap["S6A_PROXY"] = &serviceStats
	serviceStatsMap["SESSION_PROXY"] = &serviceStats

	return &protos.HealthStats{
		SystemStatus: &protos.SystemHealthStats{
			Time:              uint64(time.Now().Unix()),
			CpuUtilPct:        0.25,
			MemAvailableBytes: 5000000,
			MemTotalBytes:     50000000000,
		},
		ServiceStatus: serviceStatsMap,
		Health: &protos.HealthStatus{
			Health:        protos.HealthStatus_HEALTHY,
			HealthMessage: "OK",
		},
	}
}
