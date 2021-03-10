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

package servicers

import (
	"context"
	"fmt"

	"magma/cwf/cloud/go/protos/mconfig"
	"magma/cwf/gateway/services/gateway_health/events"
	"magma/cwf/gateway/services/gateway_health/health/gre_probe"
	"magma/cwf/gateway/services/gateway_health/health/service_health"
	"magma/cwf/gateway/services/gateway_health/health/system_health"
	"magma/feg/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

const (
	sessiondServiceName = "sessiond"
	radiusServiceName   = "radius"
	aaaServiceName      = "aaa_server"

	disabled gatewayState = "disabled"
	enabled  gatewayState = "enabled"
)

type gatewayState string

type GatewayHealthServicer struct {
	config        *mconfig.CwfGatewayHealthConfig
	greProbe      gre_probe.GREProbe
	serviceHealth service_health.ServiceHealth
	systemHealth  system_health.SystemHealth
	currentState  gatewayState
}

// NewGatewayHealthServicer constructs a GatewayHealthServicer.
func NewGatewayHealthServicer(
	cfg *mconfig.CwfGatewayHealthConfig,
	greProbe gre_probe.GREProbe,
	serviceHealth service_health.ServiceHealth,
	systemHealth system_health.SystemHealth,
) *GatewayHealthServicer {
	return &GatewayHealthServicer{
		config:        cfg,
		greProbe:      greProbe,
		systemHealth:  systemHealth,
		serviceHealth: serviceHealth,
		currentState:  "",
	}
}

// Disable disables gateway to be a standby. This is done by removing the
// transport VIP from the gateway and restarting the AAA service.
func (s *GatewayHealthServicer) Disable(ctx context.Context, req *protos.DisableMessage) (*orcprotos.Void, error) {
	ret := &orcprotos.Void{}
	// If already disabled, disable is a no-op
	if s.currentState == disabled {
		events.LogGatewayHealthSuccessEvent(events.GatewayDemotionSucceededEvent)
		return ret, nil
	}
	err := s.systemHealth.Disable()
	if err != nil {
		events.LogGatewayHealthFailedEvent(events.GatewayDemotionFailedEvent, err.Error())
		return ret, err
	}
	// Restart the AAA server to clear in-memory sessions
	err = s.serviceHealth.Restart(aaaServiceName)
	if err != nil {
		events.LogGatewayHealthFailedEvent(events.GatewayDemotionFailedEvent, err.Error())
		return ret, err
	}
	events.LogGatewayHealthSuccessEvent(events.GatewayDemotionSucceededEvent)
	s.currentState = disabled
	return ret, nil
}

// Enable enables the gateway to be active. This is done by adding the
// transport VIP to the gateway and restarting sessiond and the radius server.
func (s *GatewayHealthServicer) Enable(ctx context.Context, req *orcprotos.Void) (*orcprotos.Void, error) {
	ret := &orcprotos.Void{}
	// If already enabled, enable is a no-op
	if s.currentState == enabled {
		events.LogGatewayHealthSuccessEvent(events.GatewayPromotionSucceededEvent)
		return ret, nil
	}
	err := s.systemHealth.Enable()
	if err != nil {
		events.LogGatewayHealthFailedEvent(events.GatewayPromotionFailedEvent, err.Error())
		return ret, err
	}
	err = s.serviceHealth.Restart(sessiondServiceName)
	if err != nil {
		events.LogGatewayHealthFailedEvent(events.GatewayPromotionFailedEvent, err.Error())
		return ret, err
	}
	err = s.serviceHealth.Restart(radiusServiceName)
	if err != nil {
		events.LogGatewayHealthFailedEvent(events.GatewayPromotionFailedEvent, err.Error())
		return ret, err
	}
	if err != nil {
		events.LogGatewayHealthFailedEvent(events.GatewayPromotionFailedEvent, err.Error())
		return ret, err
	}
	events.LogGatewayHealthSuccessEvent(events.GatewayPromotionSucceededEvent)
	s.currentState = enabled
	return ret, nil
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
	probeStatus := s.greProbe.GetStatus()
	glog.V(1).Infof("reachable GRE endpoints: %v; unreachable GRE endpoints: %v", probeStatus.Reachable, probeStatus.Unreachable)

	// Current approach is to be conservative for GRE health. As long as we have
	// a reachable peer, determine to be healthy
	if len(probeStatus.Reachable) == 0 && len(probeStatus.Unreachable) > 0 {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("All GRE peers are detected as unreachable; unreachable: %v", probeStatus.Unreachable),
		}
	}
	return &protos.HealthStatus{
		Health:        protos.HealthStatus_HEALTHY,
		HealthMessage: fmt.Sprintf("GRE peers reachable; reachable: %v, unreachable: %v", probeStatus.Reachable, probeStatus.Unreachable),
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
	if stats.CpuUtilPct > s.config.CpuUtilThresholdPct {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("current cpuUtilPct execeeds threshold: %f > %f", stats.CpuUtilPct, s.config.CpuUtilThresholdPct),
		}
	}
	if stats.MemUtilPct > s.config.MemUtilThresholdPct {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("current memUtilPct execeeds threshold: %f > %f", stats.MemUtilPct, s.config.MemUtilThresholdPct),
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
