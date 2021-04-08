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

// Package health_test tests the health service's functionality by using the methods
// in package health
package health_test

import (
	"context"
	"fmt"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/health"
	health_test_init "magma/feg/cloud/go/services/health/test_init"
	"magma/feg/cloud/go/services/health/test_utils"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	orcprotos "magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/stretchr/testify/assert"
)

// Test the health service by simulating one FeG in a network
// providing health updates
func TestHealthAPI_SingleFeg(t *testing.T) {
	// Initialize test services
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	testServicer, err := health_test_init.StartTestService(t)
	assert.NoError(t, err)

	// Register network and feg then perform mock health updates
	test_utils.RegisterNetwork(t, test_utils.TestFegNetwork)
	test_utils.RegisterGateway(
		t,
		test_utils.TestFegNetwork,
		test_utils.TestFegHwId1,
		test_utils.TestFegLogicalId1,
	)
	active, err := health.GetActiveGateway(test_utils.TestFegNetwork)
	assert.NoError(t, err)
	assert.Equal(t, test_utils.TestFegNetwork, active)

	// Simulate request coming from feg1
	testServicer.Feg1 = true

	// First FeG is initially healthy
	healthyRequest := test_utils.GetHealthyRequest()
	res, err := updateHealth(t, healthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_UP, res.Action)

	activeID, err := health.GetActiveGateway(test_utils.TestFegNetwork)
	assert.NoError(t, err)
	assert.Equal(t, test_utils.TestFegLogicalId1, activeID)
	checkHealthData(t, test_utils.TestFegNetwork, test_utils.TestFegLogicalId1, healthyRequest.HealthStats)

	// Now is unhealthy, but should stay ACTIVE since it's the only FeG registered
	unhealthyRequest := test_utils.GetUnhealthyRequest()
	res, err = updateHealth(t, unhealthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_UP, res.Action)

	activeID, err = health.GetActiveGateway(test_utils.TestFegNetwork)
	assert.NoError(t, err)
	assert.Equal(t, test_utils.TestFegLogicalId1, activeID)
	checkHealthData(t, test_utils.TestFegNetwork, test_utils.TestFegLogicalId1, unhealthyRequest.HealthStats)
}

// Test the health service by simulating two FeGs in the same network
// providing health updates
func TestHealthAPI_DualFeg(t *testing.T) {
	// Initialize test services
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	testServicer, err := health_test_init.StartTestService(t)
	assert.NoError(t, err)

	// Register network and a feg
	test_utils.RegisterNetwork(t, test_utils.TestFegNetwork)
	test_utils.RegisterGateway(
		t,
		test_utils.TestFegNetwork,
		test_utils.TestFegHwId1,
		test_utils.TestFegLogicalId1,
	)

	// Simulate request coming from first feg
	testServicer.Feg1 = true

	// First FeG is initially healthy
	healthyRequest := test_utils.GetHealthyRequest()
	res, err := updateHealth(t, healthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_UP, res.Action)

	activeID, err := health.GetActiveGateway(test_utils.TestFegNetwork)
	assert.NoError(t, err)
	assert.Equal(t, test_utils.TestFegLogicalId1, activeID)
	checkHealthData(t, test_utils.TestFegNetwork, test_utils.TestFegLogicalId1, healthyRequest.HealthStats)

	// Now register a second FeG
	test_utils.RegisterGateway(
		t,
		test_utils.TestFegNetwork,
		test_utils.TestFegHwId2,
		test_utils.TestFegLogicalId2,
	)

	// Simulate request coming from second feg
	testServicer.Feg1 = false

	// Healthy request from standby
	res, err = updateHealth(t, healthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_DOWN, res.Action)

	activeID, err = health.GetActiveGateway(test_utils.TestFegNetwork)
	assert.NoError(t, err)
	assert.Equal(t, test_utils.TestFegLogicalId1, activeID)
	checkHealthData(t, test_utils.TestFegNetwork, test_utils.TestFegLogicalId2, healthyRequest.HealthStats)

	// Simulate an unhealthy request from the active feg, triggering a failover
	testServicer.Feg1 = true

	unhealthyRequest := test_utils.GetUnhealthyRequest()
	res, err = updateHealth(t, unhealthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_DOWN, res.Action)

	activeID, err = health.GetActiveGateway(test_utils.TestFegNetwork)
	assert.NoError(t, err)
	assert.Equal(t, test_utils.TestFegLogicalId2, activeID)
	checkHealthData(t, test_utils.TestFegNetwork, test_utils.TestFegLogicalId1, unhealthyRequest.HealthStats)

	// Now newly promoted FeG should receive SYSTEM_UP
	testServicer.Feg1 = false

	res, err = updateHealth(t, healthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_UP, res.Action)

	activeID, err = health.GetActiveGateway(test_utils.TestFegNetwork)
	assert.NoError(t, err)
	assert.Equal(t, test_utils.TestFegLogicalId2, activeID)
	checkHealthData(t, test_utils.TestFegNetwork, test_utils.TestFegLogicalId2, healthyRequest.HealthStats)

	// Now if the active becomes unhealthy, but standby is also unhealthy, failover doesn't occur
	res, err = updateHealth(t, unhealthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_UP, res.Action)

	activeID, err = health.GetActiveGateway(test_utils.TestFegNetwork)
	assert.NoError(t, err)
	assert.Equal(t, test_utils.TestFegLogicalId2, activeID)
	checkHealthData(t, test_utils.TestFegNetwork, test_utils.TestFegLogicalId2, unhealthyRequest.HealthStats)

	// If then standby becomes healthy, it will trigger a failover (assuming active is still unhealthy)
	testServicer.Feg1 = true

	res, err = updateHealth(t, healthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_UP, res.Action)

	activeID, err = health.GetActiveGateway(test_utils.TestFegNetwork)
	assert.NoError(t, err)
	assert.Equal(t, test_utils.TestFegLogicalId1, activeID)
	checkHealthData(t, test_utils.TestFegNetwork, test_utils.TestFegLogicalId1, healthyRequest.HealthStats)

	testServicer.Feg1 = false

	// Now newly demoted active wil receives SYSTEM_DOWN
	res, err = updateHealth(t, unhealthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_DOWN, res.Action)

	activeID, err = health.GetActiveGateway(test_utils.TestFegNetwork)
	assert.NoError(t, err)
	assert.Equal(t, test_utils.TestFegLogicalId1, activeID)
	checkHealthData(t, test_utils.TestFegNetwork, test_utils.TestFegLogicalId2, unhealthyRequest.HealthStats)
}

func updateHealth(t *testing.T, req *protos.HealthRequest) (*protos.HealthResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("Nil HealthRequest")
	}
	conn, err := registry.GetConnection(health.ServiceName)
	assert.NoError(t, err)

	client := protos.NewHealthClient(conn)
	return client.UpdateHealth(context.Background(), req)
}

func checkHealthData(t *testing.T, networkID, gatewayID string, expected *protos.HealthStats) {
	actual, err := health.GetHealth(networkID, gatewayID)
	assert.NoError(t, err)
	// Let's just set time to 0 for comparison - we should inject a time
	// provider dependency into the servicer
	actual.Time = 0
	expected.Time = 0
	assert.Equal(t, orcprotos.TestMarshal(expected), orcprotos.TestMarshal(actual))
}
