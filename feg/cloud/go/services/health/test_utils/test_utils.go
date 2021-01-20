/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package test_utils provides functions and constants that are useful for health
// service testing
package test_utils

import (
	"testing"
	"time"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/serdes"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_utils"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"

	"github.com/stretchr/testify/assert"
)

const TestFegHwId1 = "Test-FeG-Hw-Id1"
const TestFegLogicalId1 = "Test-FeG-Logical1"
const TestFegHwId2 = "Test-FeG-Hw-Id2"
const TestFegLogicalId2 = "Test-FeG-Logical2"
const TestFegNetwork = "test-feg-network"

func GetHealthyRequest() *protos.HealthRequest {
	serviceStats := protos.ServiceHealthStats{
		ServiceState: protos.ServiceHealthStats_AVAILABLE,
		ServiceHealthStatus: &protos.HealthStatus{
			Health: protos.HealthStatus_HEALTHY,
		},
	}

	serviceStatsMap := make(map[string]*protos.ServiceHealthStats)
	serviceStatsMap["SWX_PROXY"] = &serviceStats
	serviceStatsMap["SESSION_PROXY"] = &serviceStats
	serviceStatsMap["S6A_PROXY"] = &serviceStats

	healthStats1 := &protos.HealthStats{
		SystemStatus: &protos.SystemHealthStats{
			Time:              uint64(clock.Now().UnixNano()) / uint64(time.Millisecond),
			CpuUtilPct:        0.25,
			MemAvailableBytes: 7500000000,
			MemTotalBytes:     8000000000,
		},
		ServiceStatus: serviceStatsMap,
		Health: &protos.HealthStatus{
			Health:        protos.HealthStatus_HEALTHY,
			HealthMessage: "OK",
		},
		Time: uint64(clock.Now().UnixNano()) / uint64(time.Millisecond),
	}
	return &protos.HealthRequest{
		HealthStats: healthStats1,
	}
}

func GetUnhealthyRequest() *protos.HealthRequest {
	serviceStats := protos.ServiceHealthStats{
		ServiceState: protos.ServiceHealthStats_UNAVAILABLE,
		ServiceHealthStatus: &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: "Service unhealthy",
		},
	}

	serviceStatsMap := make(map[string]*protos.ServiceHealthStats)
	serviceStatsMap["SWX_PROXY"] = &serviceStats
	serviceStatsMap["SESSION_PROXY"] = &serviceStats

	healthStats1 := &protos.HealthStats{
		SystemStatus: &protos.SystemHealthStats{
			Time:              uint64(clock.Now().Unix()),
			CpuUtilPct:        0.25,
			MemAvailableBytes: 7500000000,
			MemTotalBytes:     8000000000,
		},
		ServiceStatus: serviceStatsMap,
		Health: &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: "Service: SWX_PROXY unhealthy",
		},
		Time: uint64(clock.Now().UnixNano()) / uint64(time.Millisecond),
	}
	return &protos.HealthRequest{
		HealthStats: healthStats1,
	}
}

func RegisterNetwork(t *testing.T, networkID string) {
	err := configurator.CreateNetwork(
		configurator.Network{
			ID:   TestFegNetwork,
			Type: feg.FegNetworkType,
		},
		serdes.Network,
	)
	assert.NoError(t, err)
}

func RegisterGateway(t *testing.T, networkID, hwID, logicalID string) {
	gwRecord := &models.GatewayDevice{
		HardwareID: hwID,
		Key: &models.ChallengeKey{
			KeyType: "ECHO",
		},
	}
	test_utils.RegisterGateway(t, networkID, logicalID, gwRecord)
}
