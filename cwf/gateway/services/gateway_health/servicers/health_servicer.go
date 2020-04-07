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
	"fmt"

	"magma/cwf/gateway/services/gateway_health/health/gre_probe"
	"magma/cwf/gateway/services/gateway_health/health/service_health"
	"magma/cwf/gateway/services/gateway_health/health/system_health"
	"magma/feg/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

type GatewayHealthServicer struct {
	config        *HealthConfig
	greProbe      gre_probe.GREProbe
	serviceHealth service_health.ServiceHealth
	systemHealth  system_health.SystemHealth
}

type HealthConfig struct {
	GrePeers      []string
	MaxCpuUtilPct float64
	MaxMemUtilPct float64
}

const (
	sessiondServiceName = "sessiond"
)

// NewGatewayHealthServicer constructs a GatewayHealthServicer.
func NewGatewayHealthServicer(cfg *HealthConfig, greProbe gre_probe.GREProbe, serviceHealth service_health.ServiceHealth, systemHealth system_health.SystemHealth) *GatewayHealthServicer {
	return &GatewayHealthServicer{
		config:        cfg,
		greProbe:      greProbe,
		systemHealth:  systemHealth,
		serviceHealth: serviceHealth,
	}
}

// Disable adds a drop rule for ICMP on eth1 of the gateway. This is used to ensure that
// the standby gateway is perceived as being down by the AP/WLC.
func (s *GatewayHealthServicer) Disable(ctx context.Context, req *protos.DisableMessage) (*orcprotos.Void, error) {
	ret := &orcprotos.Void{}
	return ret, s.systemHealth.Disable()
}

// Enable ensures ICMP is enabled on eth1, then restarts sessiond to trigger
// initialization of the gateway.
func (s *GatewayHealthServicer) Enable(ctx context.Context, req *orcprotos.Void) (*orcprotos.Void, error) {
	ret := &orcprotos.Void{}
	err := s.systemHealth.Enable()
	if err != nil {
		return ret, err
	}
	err = s.serviceHealth.Enable(sessiondServiceName)
	return ret, err
}

// GetHealthStatus retrieves a health status object which contains the current
// health of the gateway.
func (s *GatewayHealthServicer) GetHealthStatus(ctx context.Context, req *orcprotos.Void) (*protos.HealthStatus, error) {
	greHealth := s.getGREHealth()
	systemHealth := s.getSystemHealth()
	serviceHealth := s.getServiceHealth()
	return s.composeAggregateHealth(greHealth, systemHealth, serviceHealth), nil

}

func (s *GatewayHealthServicer) getGREHealth() *protos.HealthStatus {
	reachable, unreachable := s.greProbe.GetStatus()
	glog.V(1).Infof("reachable GRE endpoints: %v; unreachable GRE endpoints: %v", reachable, unreachable)
	// Current approach is to be conservative for GRE health. As long as we have
	// a reachable peer, determine to be healthy
	if len(reachable) == 0 && len(unreachable) > 0 {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("All GRE peers are detected as unreachable; unreachable: %v", unreachable),
		}
	}
	return &protos.HealthStatus{
		Health:        protos.HealthStatus_HEALTHY,
		HealthMessage: fmt.Sprintf("At least 1 configured GRE peers is reachable; reachable: %v, unreachable: %v", reachable, unreachable),
	}
}

func (s *GatewayHealthServicer) getSystemHealth() *protos.HealthStatus {
	stats, err := s.systemHealth.GetSystemStats()
	if err != nil {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("could not fetch system metrics"),
		}
	}
	glog.V(1).Infof("system stats: cpuUtilPct: %f, memUtilPct: %f", stats.CpuUtilPct, stats.MemUtilPct)
	if stats.CpuUtilPct > s.config.MaxCpuUtilPct {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("current cpuUtilPct execeeds threshold: %f > %f", stats.CpuUtilPct, s.config.MaxCpuUtilPct),
		}
	}
	if stats.MemUtilPct > s.config.MaxMemUtilPct {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("current memUtilPct execeeds threshold: %f > %f", stats.MemUtilPct, s.config.MaxMemUtilPct),
		}
	}
	return &protos.HealthStatus{
		Health:        protos.HealthStatus_HEALTHY,
		HealthMessage: "All metrics appear healthy",
	}
}

func (s *GatewayHealthServicer) getServiceHealth() *protos.HealthStatus {
	unhealthyServices, err := s.serviceHealth.GetUnhealthyServices()
	if err != nil {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: err.Error(),
		}
	}
	glog.V(1).Infof("unhealthy services: %v", unhealthyServices)
	if len(unhealthyServices) > 0 {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("The following services were unhealthy: %v", unhealthyServices),
		}
	}
	return &protos.HealthStatus{
		Health:        protos.HealthStatus_HEALTHY,
		HealthMessage: fmt.Sprintf("All services appear healthy"),
	}
}

func (s *GatewayHealthServicer) composeAggregateHealth(
	greHealth *protos.HealthStatus,
	systemHealth *protos.HealthStatus,
	serviceHealth *protos.HealthStatus,
) *protos.HealthStatus {
	isGatewayHealthy := greHealth.Health == protos.HealthStatus_HEALTHY && serviceHealth.Health == protos.HealthStatus_HEALTHY &&
		systemHealth.Health == protos.HealthStatus_HEALTHY
	if isGatewayHealthy {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_HEALTHY,
			HealthMessage: "gateway status appears healthy",
		}
	}
	healthMsg := ""
	if greHealth.Health == protos.HealthStatus_UNHEALTHY {
		healthMsg = fmt.Sprintf("GRE status: %s; ", greHealth.HealthMessage)
	}
	if systemHealth.Health == protos.HealthStatus_UNHEALTHY {
		healthMsg = fmt.Sprintf("%sSystem status: %s; ", healthMsg, systemHealth.HealthMessage)
	}
	if serviceHealth.Health == protos.HealthStatus_UNHEALTHY {
		healthMsg = fmt.Sprintf("%sService status: %s", healthMsg, serviceHealth.HealthMessage)
	}
	return &protos.HealthStatus{
		Health:        protos.HealthStatus_UNHEALTHY,
		HealthMessage: healthMsg,
	}
}
