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

package health_manager_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/gateway_health/health_manager"
	"magma/feg/gateway/services/session_proxy/relay/mocks"
	"magma/gateway/mconfig"
	"magma/orc8r/cloud/go/test_utils"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHealthServicer struct {
	mock.Mock
}

func (c *MockHealthServicer) UpdateHealth(
	ctx context.Context,
	req *protos.HealthRequest,
) (*protos.HealthResponse, error) {
	args := c.Called(ctx, req)
	return args.Get(0).(*protos.HealthResponse), args.Error(1)
}

func (c *MockHealthServicer) GetHealth(
	ctx context.Context,
	req *protos.GatewayStatusRequest,
) (*protos.HealthStats, error) {
	args := c.Called(ctx, req)
	return args.Get(0).(*protos.HealthStats), args.Error(1)
}

func (c *MockHealthServicer) GetClusterState(
	ctx context.Context,
	req *protos.ClusterStateRequest,
) (*protos.ClusterState, error) {
	args := c.Called(ctx, req)
	return args.Get(0).(*protos.ClusterState), args.Error(1)
}

type MockServiceHealthServicer struct {
	mock.Mock
}

func (s *MockServiceHealthServicer) Disable(ctx context.Context, req *protos.DisableMessage) (*orcprotos.Void, error) {
	args := s.Called(ctx, req)
	return args.Get(0).(*orcprotos.Void), args.Error(1)
}

func (s *MockServiceHealthServicer) Enable(ctx context.Context, void *orcprotos.Void) (*orcprotos.Void, error) {
	args := s.Called(ctx, void)
	return args.Get(0).(*orcprotos.Void), args.Error(1)
}

func (s *MockServiceHealthServicer) GetHealthStatus(ctx context.Context, void *orcprotos.Void) (*protos.HealthStatus, error) {
	args := s.Called(ctx, void)
	return args.Get(0).(*protos.HealthStatus), args.Error(1)
}

type healthMocks struct {
	cloudHealthServicer      *MockHealthServicer
	fegServiceHealthServicer *MockServiceHealthServicer
}

func initTestServices(t *testing.T, mockServiceHealth *MockServiceHealthServicer, mockHealth *MockHealthServicer) *mocks.MockCloudRegistry {
	// Create tmp mconfig test file & load configs from it
	fegConfigFmt := `{
		"configsByKey": {
			"health": {
   				"@type": "type.googleapis.com/magma.mconfig.GatewayHealthConfig",
   				"requiredServices": ["SWX_PROXY", "SESSION_PROXY"],
   				"updateIntervalSecs": 10,
   				"consecutiveFailureThreshold": 3,
   				"cloudDisconnectPeriodSecs": 10,
   				"localDisconnectPeriodSecs": 1
  			}
		}
	}`
	err := mconfig.CreateLoadTempConfig(fegConfigFmt)
	if err != nil {
		t.Log(err)
	}

	srv1, lis1 := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)
	srv2, lis2 := test_utils.NewTestService(t, registry.ModuleName, registry.SESSION_PROXY)
	srv3, lis3 := test_utils.NewTestService(t, registry.ModuleName, registry.HEALTH)

	protos.RegisterServiceHealthServer(srv1.GrpcServer, mockServiceHealth)
	protos.RegisterServiceHealthServer(srv2.GrpcServer, mockServiceHealth)
	protos.RegisterHealthServer(srv3.GrpcServer, mockHealth)

	go srv1.RunTest(lis1)
	go srv2.RunTest(lis2)
	go srv3.RunTest(lis3)

	return &mocks.MockCloudRegistry{ServerAddr: lis3.Addr().String()}
}

func TestHealthManager_UpdateHealth_Healthy(t *testing.T) {
	healthMocks := &healthMocks{
		cloudHealthServicer:      &MockHealthServicer{},
		fegServiceHealthServicer: &MockServiceHealthServicer{},
	}
	mockReg := initTestServices(t, healthMocks.fegServiceHealthServicer, healthMocks.cloudHealthServicer)
	config := health_manager.GetHealthConfig()
	healthManager := health_manager.NewHealthManager(mockReg, config)

	healthyResponse := &protos.HealthResponse{
		Action: protos.HealthResponse_SYSTEM_UP,
		Time:   uint64(time.Now().UnixNano()) / uint64(time.Millisecond),
	}
	healthMocks.fegServiceHealthServicer.On("GetHealthStatus", mock.Anything, mock.Anything).Return(getHealthyServiceStatus(), nil).Twice()
	healthMocks.cloudHealthServicer.On("UpdateHealth", mock.Anything, mock.Anything).Return(healthyResponse, nil).Once()
	healthMocks.fegServiceHealthServicer.On("Enable", mock.Anything, &orcprotos.Void{}).Return(&orcprotos.Void{}, nil).Twice()
	err := healthManager.SendHealthUpdate()

	healthMocks.fegServiceHealthServicer.AssertExpectations(t)
	healthMocks.cloudHealthServicer.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestHealthManager_UpdateHealth_SystemDown(t *testing.T) {
	healthMocks := &healthMocks{
		cloudHealthServicer:      &MockHealthServicer{},
		fegServiceHealthServicer: &MockServiceHealthServicer{},
	}
	mockReg := initTestServices(t, healthMocks.fegServiceHealthServicer, healthMocks.cloudHealthServicer)
	config := health_manager.GetHealthConfig()
	healthManager := health_manager.NewHealthManager(mockReg, config)

	closeResponse := &protos.HealthResponse{
		Action: protos.HealthResponse_SYSTEM_DOWN,
		Time:   uint64(time.Now().UnixNano()) / uint64(time.Millisecond),
	}
	healthMocks.fegServiceHealthServicer.On("GetHealthStatus", mock.Anything, mock.Anything).Return(getUnhealthyServiceStatus(), nil).Twice()
	healthMocks.cloudHealthServicer.On("UpdateHealth", mock.Anything, mock.Anything).Return(closeResponse, nil).Once()
	healthMocks.fegServiceHealthServicer.On("Disable", mock.Anything, &protos.DisableMessage{DisablePeriodSecs: 10}).Return(&orcprotos.Void{}, nil).Twice()

	err := healthManager.SendHealthUpdate()
	healthMocks.cloudHealthServicer.AssertExpectations(t)
	healthMocks.fegServiceHealthServicer.AssertExpectations(t)

	assert.NoError(t, err)
}

func TestHealthManager_UpdateHealth_SystemUp(t *testing.T) {
	healthMocks := &healthMocks{
		cloudHealthServicer:      &MockHealthServicer{},
		fegServiceHealthServicer: &MockServiceHealthServicer{},
	}
	mockReg := initTestServices(t, healthMocks.fegServiceHealthServicer, healthMocks.cloudHealthServicer)
	config := health_manager.GetHealthConfig()
	healthManager := health_manager.NewHealthManager(mockReg, config)

	activeResponse := &protos.HealthResponse{
		Action: protos.HealthResponse_SYSTEM_UP,
		Time:   uint64(time.Now().UnixNano()) / uint64(time.Millisecond),
	}
	healthMocks.fegServiceHealthServicer.On("GetHealthStatus", mock.Anything, mock.Anything).Return(getHealthyServiceStatus(), nil).Twice()
	healthMocks.cloudHealthServicer.On("UpdateHealth", mock.Anything, mock.Anything).Return(activeResponse, nil).Once()
	healthMocks.fegServiceHealthServicer.On("Enable", mock.Anything, &orcprotos.Void{}).Return(&orcprotos.Void{}, nil).Twice()

	err := healthManager.SendHealthUpdate()
	healthMocks.cloudHealthServicer.AssertExpectations(t)
	healthMocks.fegServiceHealthServicer.AssertExpectations(t)

	assert.NoError(t, err)
}

func TestHealthManager_ExceedUpdateFailureThreshold(t *testing.T) {
	healthMocks := &healthMocks{
		cloudHealthServicer:      &MockHealthServicer{},
		fegServiceHealthServicer: &MockServiceHealthServicer{},
	}
	mockReg := initTestServices(t, healthMocks.fegServiceHealthServicer, healthMocks.cloudHealthServicer)
	config := health_manager.GetHealthConfig()
	healthManager := health_manager.NewHealthManager(mockReg, config)
	healthMocks.fegServiceHealthServicer.On("GetHealthStatus", mock.Anything, mock.Anything).Return(getHealthyServiceStatus(), nil).Times(6)
	healthMocks.cloudHealthServicer.On("UpdateHealth", mock.Anything, mock.Anything).
		Return(&protos.HealthResponse{}, fmt.Errorf("rpc error: code = Internal desc = transport")).Times(3)

	healthMocks.fegServiceHealthServicer.On("Disable", mock.Anything, &protos.DisableMessage{DisablePeriodSecs: 1}).Return(&orcprotos.Void{}, nil).Twice()

	// Simulate 3 consecutive update failures subsequently triggering SYSTEM_DOWN

	err := healthManager.SendHealthUpdate()
	assert.Error(t, err)

	err = healthManager.SendHealthUpdate()
	assert.Error(t, err)

	err = healthManager.SendHealthUpdate()
	assert.Error(t, err)

	healthMocks.cloudHealthServicer.AssertExpectations(t)
	healthMocks.fegServiceHealthServicer.AssertExpectations(t)
}

func getHealthyServiceStatus() *protos.HealthStatus {
	return &protos.HealthStatus{
		Health: protos.HealthStatus_HEALTHY,
	}
}

func getUnhealthyServiceStatus() *protos.HealthStatus {
	return &protos.HealthStatus{
		Health:        protos.HealthStatus_UNHEALTHY,
		HealthMessage: "Service unhealthy",
	}
}

func getHealthySystemStats() *protos.SystemHealthStats {
	return &protos.SystemHealthStats{
		Time:              uint64(time.Now().Unix()),
		CpuUtilPct:        0.25,
		MemAvailableBytes: 5000000,
		MemTotalBytes:     50000000000,
	}
}
