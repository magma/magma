/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	models1 "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
)

func Test_Conversions(t *testing.T) {
	cNetwork := configurator.Network{
		ID:          "test",
		Name:        "name",
		Type:        "type",
		Description: "desc",
		Configs: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
			orc8r.DnsdNetworkType:       models.NewDefaultDNSConfig(),
			orc8r.NetworkSentryConfig:   models.NewDefaultSentryConfig(),
			orc8r.StateConfig:           models.NewDefaultStateConfig(),
		},
	}
	generatedSNetwork := (&models.Network{}).FromConfiguratorNetwork(cNetwork)
	sNetwork := models.Network{
		ID:           models1.NetworkID("test"),
		Name:         models1.NetworkName("name"),
		Type:         "type",
		Description:  models1.NetworkDescription("desc"),
		Features:     models.NewDefaultFeaturesConfig(),
		DNS:          models.NewDefaultDNSConfig(),
		SentryConfig: models.NewDefaultSentryConfig(),
		StateConfig:  models.NewDefaultStateConfig(),
	}
	generatedCNetwork := sNetwork.ToConfiguratorNetwork()

	assert.Equal(t, sNetwork, *generatedSNetwork)
	assert.Equal(t, cNetwork, generatedCNetwork)
}

func TestLastGatewayCheckInWasRecent(t *testing.T) {
	// CheckinInterval=20s, graceFactor=10 => grace period = 200s
	magmadConfig := models.MagmadGatewayConfigs{
		CheckinInterval: 20,
	}

	gatewayStatus := models.GatewayStatus{
		CheckinTime: uint64(time.Now().UnixMilli()),
	}
	assert.True(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig))

	gatewayStatus.CheckinTime = uint64(time.Now().Add(-60 * time.Second).UnixMilli())
	assert.True(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig))

	// 201s ago is outside grace period (200s)
	gatewayStatus.CheckinTime = uint64(time.Now().Add(-201 * time.Second).UnixMilli())
	assert.False(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig))

	gatewayStatus.CheckinTime = uint64(time.Now().Add(-10 * time.Hour).UnixMilli())
	assert.False(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig))
}

func TestLastGatewayCheckInWasRecentHandlingNil(t *testing.T) {
	assert.False(t, models.LastGatewayCheckInWasRecent(nil, nil))

	magmadConfig := models.MagmadGatewayConfigs{
		CheckinInterval: 20,
	}
	assert.False(t, models.LastGatewayCheckInWasRecent(nil, &magmadConfig))

	gatewayStatus := models.GatewayStatus{
		CheckinTime: uint64(time.Now().UnixMilli()),
	}
	assert.True(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, nil))

	// nil config uses default 60s interval, graceFactor=10 => grace = 600s
	// 601s ago should be outside grace period
	gatewayStatus.CheckinTime = uint64(time.Now().Add(-601 * time.Second).UnixMilli())
	assert.False(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, nil))
}

// TestLastGatewayCheckInWasRecent_BoundaryConditions tests behavior near the grace period boundary
// Issue #15725: status fluctuates near the boundary
func TestLastGatewayCheckInWasRecent_BoundaryConditions(t *testing.T) {
	// checkin_interval=20s, graceFactor=10 => grace period = 200s
	magmadConfig := models.MagmadGatewayConfigs{
		CheckinInterval: 20,
	}

	// 1 second before boundary - should be healthy
	gatewayStatus := models.GatewayStatus{
		CheckinTime: uint64(time.Now().Add(-199 * time.Second).UnixMilli()),
	}
	assert.True(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig),
		"199s ago should be within grace period")

	// exactly at boundary - may be true or false due to ms-level timing
	gatewayStatus.CheckinTime = uint64(time.Now().Add(-200 * time.Second).UnixMilli())
	_ = models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig) // boundary value not asserted due to timing uncertainty

	// 1 second past boundary - should be unhealthy
	gatewayStatus.CheckinTime = uint64(time.Now().Add(-201 * time.Second).UnixMilli())
	assert.False(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig),
		"201s ago should be outside grace period")

	// 2 seconds past boundary - definitely unhealthy
	gatewayStatus.CheckinTime = uint64(time.Now().Add(-202 * time.Second).UnixMilli())
	assert.False(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig),
		"202s ago should definitely be outside grace period")
}

// TestLastGatewayCheckInWasRecent_DefaultIntervalBoundary tests boundary with default interval
// default interval = 60s, graceFactor=10 => grace = 10*60 = 600s
func TestLastGatewayCheckInWasRecent_DefaultIntervalBoundary(t *testing.T) {
	// nil config uses default 60s interval
	// grace period = 10 * 60 = 600s

	// 599s ago - should be healthy
	gatewayStatus := models.GatewayStatus{
		CheckinTime: uint64(time.Now().Add(-599 * time.Second).UnixMilli()),
	}
	assert.True(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, nil),
		"599s ago with default interval should be healthy")

	// 601s ago - should be unhealthy
	gatewayStatus.CheckinTime = uint64(time.Now().Add(-601 * time.Second).UnixMilli())
	assert.False(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, nil),
		"601s ago with default interval should be unhealthy")
}

// TestLastGatewayCheckInWasRecent_ConsecutiveCalls tests consistency across rapid consecutive calls
// simulates NMS refreshing every 30 seconds
func TestLastGatewayCheckInWasRecent_ConsecutiveCalls(t *testing.T) {
	magmadConfig := models.MagmadGatewayConfigs{
		CheckinInterval: 60, // common production value
	}
	// graceFactor=10 => grace = 10 * 60 = 600s

	// set checkin time safely within grace period (400s ago)
	checkinTime := time.Now().Add(-400 * time.Second)
	gatewayStatus := models.GatewayStatus{
		CheckinTime: uint64(checkinTime.UnixMilli()),
	}

	// consecutive calls should return consistent results
	for i := 0; i < 5; i++ {
		result := models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig)
		assert.True(t, result, "consecutive call %d should return consistent result", i)
	}

	// set checkin time safely outside grace period (700s ago)
	gatewayStatus.CheckinTime = uint64(time.Now().Add(-700 * time.Second).UnixMilli())

	for i := 0; i < 5; i++ {
		result := models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig)
		assert.False(t, result, "consecutive call %d for old checkin should return false", i)
	}
}

// TestLastGatewayCheckInWasRecent_ZeroInterval tests edge case with zero interval
func TestLastGatewayCheckInWasRecent_ZeroInterval(t *testing.T) {
	magmadConfig := models.MagmadGatewayConfigs{
		CheckinInterval: 0,
	}

	gatewayStatus := models.GatewayStatus{
		CheckinTime: uint64(time.Now().UnixMilli()),
	}

	// with interval=0, grace period is also 0
	// gateway should be considered unhealthy even with current timestamp
	// because time.Now().Before(checkinTime.Add(0)) is false when time has elapsed
	assert.False(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig),
		"with zero interval, gateway should be unhealthy even with current timestamp")
}

// TestLastGatewayCheckInWasRecent_ExtendedGracePeriod tests the extended grace period
// Issue #15725: the grace factor was increased from 5 to 10 and the metrics threshold
// was aligned from 300s to 600s to eliminate status fluctuation caused by inconsistent thresholds.
// With default 60s interval, grace period is now 10 * 60 = 600s
func TestLastGatewayCheckInWasRecent_ExtendedGracePeriod(t *testing.T) {
	// use nil config to test with default interval (60s)
	// with graceFactor=10, grace period = 10 * 60 = 600s

	// 500s ago - should be healthy (within new extended grace period)
	gatewayStatus := models.GatewayStatus{
		CheckinTime: uint64(time.Now().Add(-500 * time.Second).UnixMilli()),
	}
	assert.True(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, nil),
		"500s ago should be healthy with extended grace period")

	// 599s ago - should still be healthy (just within boundary)
	gatewayStatus.CheckinTime = uint64(time.Now().Add(-599 * time.Second).UnixMilli())
	assert.True(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, nil),
		"599s ago should be healthy with extended grace period")

	// 601s ago - should be unhealthy (just past boundary)
	gatewayStatus.CheckinTime = uint64(time.Now().Add(-601 * time.Second).UnixMilli())
	assert.False(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, nil),
		"601s ago should be unhealthy (outside extended grace period)")
}
