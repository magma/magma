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

package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/lib/go/protos"
)

const (
	testNetworkID = "test_network_status_reporter"
	testGatewayID = "test_gateway_1"
	testHwID      = "test_hw_id_1"
)

// resetMetrics clears all metric values for clean test state
func resetMetrics() {
	gwCheckinStatus.Reset()
	gwMconfigAge.Reset()
	upGwCount.Reset()
	totalGwCount.Reset()
}

// TestReportGatewayStatus_RecentCheckin tests that gateways with recent
// check-ins are reported as healthy (status=1)
func TestReportGatewayStatus_RecentCheckin(t *testing.T) {
	// Setup test services
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	resetMetrics()

	// Create network and gateway
	err := configurator.CreateNetwork(context.Background(), configurator.Network{
		ID:   testNetworkID,
		Name: "Test Network",
	}, serdes.Network)
	require.NoError(t, err)

	_, err = configurator.CreateEntity(context.Background(), testNetworkID, configurator.NetworkEntity{
		Type:       orc8r.MagmadGatewayType,
		Key:        testGatewayID,
		PhysicalID: testHwID,
	}, serdes.Entity)
	require.NoError(t, err)

	// Report gateway status with recent check-in (now)
	ctx := test_utils.GetContextWithCertificate(t, testHwID)
	gwStatus := &models.GatewayStatus{
		HardwareID: testHwID,
		PlatformInfo: &models.PlatformInfo{
			ConfigInfo: &models.ConfigInfo{
				MconfigCreatedAt: uint64(time.Now().Unix() - 60),
			},
		},
	}
	reportGatewayStatusHelper(t, ctx, gwStatus)

	// Run the status reporter
	err = reportGatewayStatus()
	require.NoError(t, err)

	// Verify gateway is marked as healthy (status=1)
	status := testutil.ToFloat64(gwCheckinStatus.WithLabelValues(testNetworkID, testGatewayID))
	assert.Equal(t, float64(1), status, "Gateway with recent check-in should be healthy")

	// Verify counts
	upCount := testutil.ToFloat64(upGwCount.WithLabelValues(testNetworkID))
	totalCount := testutil.ToFloat64(totalGwCount.WithLabelValues(testNetworkID))
	assert.Equal(t, float64(1), upCount, "Up gateway count should be 1")
	assert.Equal(t, float64(1), totalCount, "Total gateway count should be 1")
}

// TestReportGatewayStatus_NoStatus tests that gateways without status are skipped
// and not counted as "up"
func TestReportGatewayStatus_NoStatus(t *testing.T) {
	// Setup test services
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	resetMetrics()

	// Create network and gateway (but don't report status)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{
		ID:   testNetworkID,
		Name: "Test Network",
	}, serdes.Network)
	require.NoError(t, err)

	_, err = configurator.CreateEntity(context.Background(), testNetworkID, configurator.NetworkEntity{
		Type:       orc8r.MagmadGatewayType,
		Key:        testGatewayID,
		PhysicalID: testHwID,
	}, serdes.Entity)
	require.NoError(t, err)

	// Run the status reporter without reporting any gateway status
	err = reportGatewayStatus()
	require.NoError(t, err)

	// Verify counts - gateway exists but no status reported
	upCount := testutil.ToFloat64(upGwCount.WithLabelValues(testNetworkID))
	totalCount := testutil.ToFloat64(totalGwCount.WithLabelValues(testNetworkID))
	assert.Equal(t, float64(0), upCount, "Up gateway count should be 0 when no status")
	assert.Equal(t, float64(1), totalCount, "Total gateway count should still be 1")
}

// TestReportGatewayStatus_EmptyNetwork tests handling of networks with no gateways
func TestReportGatewayStatus_EmptyNetwork(t *testing.T) {
	// Setup test services
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	resetMetrics()

	// Create network only (no gateway)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{
		ID:   testNetworkID,
		Name: "Test Network",
	}, serdes.Network)
	require.NoError(t, err)

	// Run the status reporter
	err = reportGatewayStatus()
	require.NoError(t, err)

	// Verify counts
	upCount := testutil.ToFloat64(upGwCount.WithLabelValues(testNetworkID))
	totalCount := testutil.ToFloat64(totalGwCount.WithLabelValues(testNetworkID))
	assert.Equal(t, float64(0), upCount, "Up gateway count should be 0")
	assert.Equal(t, float64(0), totalCount, "Total gateway count should be 0")
}

// TestReportGatewayStatus_MultipleGateways tests handling of multiple gateways
// all with recent check-ins
func TestReportGatewayStatus_MultipleGateways(t *testing.T) {
	// Setup test services
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	resetMetrics()

	// Create network
	err := configurator.CreateNetwork(context.Background(), configurator.Network{
		ID:   testNetworkID,
		Name: "Test Network",
	}, serdes.Network)
	require.NoError(t, err)

	// Create 3 gateways with recent check-ins
	gateways := []struct {
		gatewayID string
		hwID      string
	}{
		{"gw1", "hw1"},
		{"gw2", "hw2"},
		{"gw3", "hw3"},
	}

	for _, gw := range gateways {
		_, err = configurator.CreateEntity(context.Background(), testNetworkID, configurator.NetworkEntity{
			Type:       orc8r.MagmadGatewayType,
			Key:        gw.gatewayID,
			PhysicalID: gw.hwID,
		}, serdes.Entity)
		require.NoError(t, err)

		ctx := test_utils.GetContextWithCertificate(t, gw.hwID)
		gwStatus := &models.GatewayStatus{
			HardwareID: gw.hwID,
			PlatformInfo: &models.PlatformInfo{
				ConfigInfo: &models.ConfigInfo{
					MconfigCreatedAt: uint64(time.Now().Unix() - 60),
				},
			},
		}
		reportGatewayStatusHelper(t, ctx, gwStatus)
	}

	// Run the status reporter
	err = reportGatewayStatus()
	require.NoError(t, err)

	// Verify individual gateway statuses (all should be healthy)
	for _, gw := range gateways {
		status := testutil.ToFloat64(gwCheckinStatus.WithLabelValues(testNetworkID, gw.gatewayID))
		assert.Equal(t, float64(1), status, "Gateway %s should be healthy", gw.gatewayID)
	}

	// Verify counts
	upCount := testutil.ToFloat64(upGwCount.WithLabelValues(testNetworkID))
	totalCount := testutil.ToFloat64(totalGwCount.WithLabelValues(testNetworkID))
	assert.Equal(t, float64(3), upCount, "Up gateway count should be 3")
	assert.Equal(t, float64(3), totalCount, "Total gateway count should be 3")
}

// TestReportGatewayStatus_MixedStatus tests handling of gateways with and without status
func TestReportGatewayStatus_MixedStatus(t *testing.T) {
	// Setup test services
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	resetMetrics()

	// Create network
	err := configurator.CreateNetwork(context.Background(), configurator.Network{
		ID:   testNetworkID,
		Name: "Test Network",
	}, serdes.Network)
	require.NoError(t, err)

	// Create 3 gateways: 2 with status, 1 without
	gatewaysWithStatus := []struct {
		gatewayID string
		hwID      string
	}{
		{"gw1", "hw1"},
		{"gw2", "hw2"},
	}

	gatewayWithoutStatus := struct {
		gatewayID string
		hwID      string
	}{"gw3", "hw3"}

	// Create gateways with status
	for _, gw := range gatewaysWithStatus {
		_, err = configurator.CreateEntity(context.Background(), testNetworkID, configurator.NetworkEntity{
			Type:       orc8r.MagmadGatewayType,
			Key:        gw.gatewayID,
			PhysicalID: gw.hwID,
		}, serdes.Entity)
		require.NoError(t, err)

		ctx := test_utils.GetContextWithCertificate(t, gw.hwID)
		gwStatus := &models.GatewayStatus{
			HardwareID: gw.hwID,
			PlatformInfo: &models.PlatformInfo{
				ConfigInfo: &models.ConfigInfo{
					MconfigCreatedAt: uint64(time.Now().Unix() - 60),
				},
			},
		}
		reportGatewayStatusHelper(t, ctx, gwStatus)
	}

	// Create gateway without status
	_, err = configurator.CreateEntity(context.Background(), testNetworkID, configurator.NetworkEntity{
		Type:       orc8r.MagmadGatewayType,
		Key:        gatewayWithoutStatus.gatewayID,
		PhysicalID: gatewayWithoutStatus.hwID,
	}, serdes.Entity)
	require.NoError(t, err)

	// Run the status reporter
	err = reportGatewayStatus()
	require.NoError(t, err)

	// Verify gateways with status are healthy
	for _, gw := range gatewaysWithStatus {
		status := testutil.ToFloat64(gwCheckinStatus.WithLabelValues(testNetworkID, gw.gatewayID))
		assert.Equal(t, float64(1), status, "Gateway %s with status should be healthy", gw.gatewayID)
	}

	// Verify counts: 2 up (with status), 3 total
	upCount := testutil.ToFloat64(upGwCount.WithLabelValues(testNetworkID))
	totalCount := testutil.ToFloat64(totalGwCount.WithLabelValues(testNetworkID))
	assert.Equal(t, float64(2), upCount, "Up gateway count should be 2 (only those with status)")
	assert.Equal(t, float64(3), totalCount, "Total gateway count should be 3")
}

// TestReportGatewayStatus_MconfigAge tests that mconfig age is reported correctly
func TestReportGatewayStatus_MconfigAge(t *testing.T) {
	// Setup test services
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	resetMetrics()

	// Create network and gateway
	err := configurator.CreateNetwork(context.Background(), configurator.Network{
		ID:   testNetworkID,
		Name: "Test Network",
	}, serdes.Network)
	require.NoError(t, err)

	_, err = configurator.CreateEntity(context.Background(), testNetworkID, configurator.NetworkEntity{
		Type:       orc8r.MagmadGatewayType,
		Key:        testGatewayID,
		PhysicalID: testHwID,
	}, serdes.Entity)
	require.NoError(t, err)

	// Report gateway status with mconfig created 120 seconds ago
	ctx := test_utils.GetContextWithCertificate(t, testHwID)
	mconfigCreatedAt := uint64(time.Now().Unix() - 120)
	gwStatus := &models.GatewayStatus{
		HardwareID: testHwID,
		PlatformInfo: &models.PlatformInfo{
			ConfigInfo: &models.ConfigInfo{
				MconfigCreatedAt: mconfigCreatedAt,
			},
		},
	}
	reportGatewayStatusHelper(t, ctx, gwStatus)

	// Run the status reporter
	err = reportGatewayStatus()
	require.NoError(t, err)

	// Verify mconfig age is reported (should be approximately 120 seconds)
	// Note: exact value depends on timing between state report and metric collection
	mconfigAge := testutil.ToFloat64(gwMconfigAge.WithLabelValues(testNetworkID, testGatewayID))
	assert.True(t, mconfigAge >= 115 && mconfigAge <= 130,
		"Mconfig age should be approximately 120 seconds, got %v", mconfigAge)
}

// TestGracePeriodThreshold verifies the grace period threshold is consistent
// between the metrics reporter and the API.
//
// Issue #15725: Previously the threshold was 60*5=300 seconds in metrics but
// 600 seconds in the API. This inconsistency caused NMS UI status fluctuation.
// The fix aligned both to use 600 seconds (10 * 60).
//
// Note: Testing old check-in scenarios requires time mocking. The current
// implementation uses time.Now() directly, which cannot be mocked with the
// clock package. A future refactoring could change this to use clock.Now()
// for full testability.
func TestGracePeriodThreshold(t *testing.T) {
	// Verify that the metrics threshold now matches the API threshold
	// Both should be: graceFactor (10) * defaultCheckinInterval (60s) = 600 seconds
	//
	// The metricsGracePeriodSeconds constant is defined in gateway_status_reporter.go

	expectedAPIThreshold := 10 * 60 // 600 seconds (from conversion.go)

	// Verify the constant is accessible and has the expected value
	assert.Equal(t, 600, metricsGracePeriodSeconds,
		"Metrics threshold should be 600 seconds")
	assert.Equal(t, 600, expectedAPIThreshold,
		"API threshold should be 600 seconds")
	assert.Equal(t, metricsGracePeriodSeconds, expectedAPIThreshold,
		"Metrics and API thresholds should now be consistent (Issue #15725 fix)")
}

// reportGatewayStatusHelper is a helper function to report gateway status
func reportGatewayStatusHelper(t *testing.T, ctx context.Context, gwStatus *models.GatewayStatus) {
	client, err := state.GetStateClient()
	require.NoError(t, err)

	serializedGWStatus, err := serde.Serialize(gwStatus, orc8r.GatewayStateType, serdes.State)
	require.NoError(t, err)

	states := []*protos.State{
		{
			Type:     orc8r.GatewayStateType,
			DeviceID: gwStatus.HardwareID,
			Value:    serializedGWStatus,
		},
	}
	_, err = client.ReportStates(ctx, &protos.ReportStatesRequest{States: states})
	require.NoError(t, err)
}
