/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"context"
	"testing"

	"magma/cwf/gateway/services/gateway_health/health/system_health"
	"magma/feg/cloud/go/protos"
	orc8rprotos "magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetHealthStatus(t *testing.T) {
	mockService := &mockServiceHealth{}
	mockSystem := &mockSystemHealth{}
	mockGREProbe := &mockGREProbe{}
	req := &orc8rprotos.Void{}
	hc := &HealthConfig{
		GrePeers:      []string{"127.0.0.1"},
		MaxCpuUtilPct: 0.75,
		MaxMemUtilPct: 0.75,
	}
	servicer := NewGatewayHealthServicer(hc, mockGREProbe, mockService, mockSystem)
	expectedStatus := &protos.HealthStatus{
		Health:        protos.HealthStatus_HEALTHY,
		HealthMessage: "gateway status appears healthy",
	}
	// Simulate healthy status
	mockGREProbe.On("GetStatus").Return([]string{"127.0.0.1"}, []string{}, nil).Once()
	mockSystem.On("GetSystemStats").Return(&system_health.SystemStats{CpuUtilPct: 0.1, MemUtilPct: 0.1}, nil).Once()
	mockService.On("GetUnhealthyServices").Return([]string{}, nil).Once()
	health, err := servicer.GetHealthStatus(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, health)
	assertMocks(t, mockGREProbe, mockSystem, mockService)

	// Simulate successful enable
	mockSystem.On("Enable").Return(nil)
	mockService.On("Enable", "sessiond").Return(nil)
	_, err = servicer.Enable(context.Background(), req)
	assert.NoError(t, err)
	assertMocks(t, mockGREProbe, mockSystem, mockService)

	// Simulate GRE unhealthy
	expectedStatus.Health = protos.HealthStatus_UNHEALTHY
	expectedStatus.HealthMessage = "GRE status: All GRE peers are detected as unreachable; unreachable: [127.0.0.1]; "
	mockGREProbe.On("GetStatus").Return([]string{}, []string{"127.0.0.1"}, nil).Once()
	mockSystem.On("GetSystemStats").Return(&system_health.SystemStats{CpuUtilPct: 0.1, MemUtilPct: 0.1}, nil).Once()
	mockService.On("GetUnhealthyServices").Return([]string{}, nil).Once()
	health, err = servicer.GetHealthStatus(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, health)
	assertMocks(t, mockGREProbe, mockSystem, mockService)

	// Simulate successful disable
	disableReq := &protos.DisableMessage{}
	mockSystem.On("Disable").Return(nil)
	_, err = servicer.Disable(context.Background(), disableReq)
	assert.NoError(t, err)
	assertMocks(t, mockGREProbe, mockSystem, mockService)

	// Simulate unhealthy system status
	mockGREProbe.On("GetStatus").Return([]string{"127.0.0.1"}, []string{}, nil).Once()
	mockSystem.On("GetSystemStats").Return(&system_health.SystemStats{CpuUtilPct: 0.99, MemUtilPct: 0.5}, nil).Once()
	mockService.On("GetUnhealthyServices").Return([]string{}, nil).Once()
	expectedStatus.Health = protos.HealthStatus_UNHEALTHY
	expectedStatus.HealthMessage = "System status: current cpuUtilPct execeeds threshold: 0.990000 > 0.750000; "
	health, err = servicer.GetHealthStatus(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, health)
	assertMocks(t, mockGREProbe, mockSystem, mockService)

	// Simulate unhealthy services and GRE
	mockGREProbe.On("GetStatus").Return([]string{}, []string{"127.0.0.1"}, nil).Once()
	mockSystem.On("GetSystemStats").Return(&system_health.SystemStats{CpuUtilPct: 0.1, MemUtilPct: 0.1}, nil).Once()
	mockService.On("GetUnhealthyServices").Return([]string{"sessiond"}, nil).Once()
	expectedStatus.Health = protos.HealthStatus_UNHEALTHY
	expectedStatus.HealthMessage = "GRE status: All GRE peers are detected as unreachable; unreachable: [127.0.0.1]; Service status: The following services were unhealthy: [sessiond]"
	health, err = servicer.GetHealthStatus(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, health)
	assertMocks(t, mockGREProbe, mockSystem, mockService)
}

func assertMocks(t *testing.T, probe *mockGREProbe, systemHealth *mockSystemHealth, serviceHealth *mockServiceHealth) {
	probe.AssertExpectations(t)
	systemHealth.AssertExpectations(t)
	serviceHealth.AssertExpectations(t)
}

type mockServiceHealth struct {
	mock.Mock
}

func (m *mockServiceHealth) GetUnhealthyServices() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockServiceHealth) Enable(service string) error {
	args := m.Called(service)
	return args.Error(0)
}

type mockSystemHealth struct {
	mock.Mock
}

func (m *mockSystemHealth) GetSystemStats() (*system_health.SystemStats, error) {
	args := m.Called()
	return args.Get(0).(*system_health.SystemStats), args.Error(1)
}

func (m *mockSystemHealth) Enable() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockSystemHealth) Disable() error {
	args := m.Called()
	return args.Error(0)
}

type mockGREProbe struct {
	mock.Mock
}

func (m *mockGREProbe) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockGREProbe) GetStatus() ([]string, []string) {
	args := m.Called()
	return args.Get(0).([]string), args.Get(1).([]string)
}
