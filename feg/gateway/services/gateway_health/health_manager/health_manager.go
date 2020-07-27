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

// Package health_manager provides the main functionality for the gateway_health
// service. The health manager collects Federated Gateway service and system
// metrics related to health and reports them to the cloud, implementing any requested
// action sent back.
package health_manager

import (
	"fmt"
	"strings"
	"sync/atomic"

	"magma/feg/gateway/services/gateway_health/events"

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/service_health"
	"magma/feg/gateway/services/gateway_health"
	"magma/feg/gateway/services/gateway_health/collection"
	gwmcfg "magma/gateway/mconfig"
	"magma/gateway/service_registry"
)

const (
	defaultHealthUpdateIntervalSec      = 10
	defaultCloudDisablePeriodSecs       = 10
	defaultLocalDisablePeriodSecs       = 1
	defaultConsecutiveFailuresThreshold = 3
)

var defaultServices = []string{registry.SWX_PROXY, registry.SESSION_PROXY}

type HealthManager struct {
	cloudReg                  service_registry.GatewayRegistry
	config                    *mconfig.GatewayHealthConfig
	consecutiveUpdateFailures uint32
	prevAction                protos.HealthResponse_RequestedAction
}

func NewHealthManager(cloudReg service_registry.GatewayRegistry, hcfg *mconfig.GatewayHealthConfig) *HealthManager {
	return &HealthManager{
		config:                    hcfg,
		cloudReg:                  cloudReg,
		consecutiveUpdateFailures: 0,
		prevAction:                protos.HealthResponse_NONE,
	}
}

// GetHealthConfig attempts to retrieve a GatewayHealthConfig from mconfig
// If this retrieval fails, or retrieves an invalid config, the config is
// set to use default values
func GetHealthConfig() *mconfig.GatewayHealthConfig {
	defaultCfg := &mconfig.GatewayHealthConfig{
		RequiredServices:          defaultServices,
		UpdateIntervalSecs:        defaultHealthUpdateIntervalSec,
		UpdateFailureThreshold:    defaultConsecutiveFailuresThreshold,
		CloudDisconnectPeriodSecs: defaultCloudDisablePeriodSecs,
		LocalDisconnectPeriodSecs: defaultLocalDisablePeriodSecs,
	}
	healthCfg := &mconfig.GatewayHealthConfig{}
	err := gwmcfg.GetServiceConfigs("health", healthCfg)
	if err != nil {
		glog.Infof("Unable to retrieve Gateway Health Config from mconfig: %s; Using default values...", err)
		return defaultCfg
	}
	err = validateHealthConfig(healthCfg)
	if err != nil {
		glog.Infof("Invalid parameters in Gateway Health Config: %s; Using default values...", err)
		return defaultCfg
	}
	glog.Info("Using mconfig values for health parameters")
	return healthCfg
}

// SendHealthUpdate collects Gateway Service and System Health Status and sends
// them to the cloud health service. It awaits a response from the cloud and
// applies any action requested from the cloud (e.g. SYSTEM_DOWN)
func (hm *HealthManager) SendHealthUpdate() error {
	healthRequest, err := hm.gatherHealthRequest()
	if err != nil {
		glog.Error(err)
	}
	healthResponse, err := gateway_health.UpdateHealth(hm.cloudReg, healthRequest)
	if err != nil {
		return hm.handleUpdateHealthFailure(healthRequest, err)
	}

	// Update was successful, so reset consecutive failure counter
	atomic.StoreUint32(&hm.consecutiveUpdateFailures, 0)

	switch healthResponse.Action {
	case protos.HealthResponse_NONE:
	case protos.HealthResponse_SYSTEM_UP:
		err = hm.takeSystemUp()
		if err != nil {
			events.LogGatewayHealthFailedEvent(events.GatewayPromotionFailedEvent, err.Error(), hm.prevAction)
		} else {
			events.LogGatewayHealthSuccessEvent(events.GatewayPromotionSucceededEvent, hm.prevAction)
		}
	case protos.HealthResponse_SYSTEM_DOWN:
		disablePeriod := hm.config.GetCloudDisconnectPeriodSecs()
		if hm.prevAction == protos.HealthResponse_SYSTEM_DOWN {
			disablePeriod = 0
		}
		err = hm.takeSystemDown(disablePeriod, healthRequest.HealthStats.ServiceStatus)
		if err != nil {
			events.LogGatewayHealthFailedEvent(events.GatewayDemotionFailedEvent, err.Error(), hm.prevAction)
		} else {
			events.LogGatewayHealthSuccessEvent(events.GatewayDemotionSucceededEvent, hm.prevAction)
		}
	default:
		err = fmt.Errorf("Invalid requested action: %s returned to FeG Health Manager", healthResponse.Action)
	}
	if err != nil {
		glog.Error(err)
		return err
	}
	hm.prevAction = healthResponse.Action
	glog.V(2).Infof("Successfully updated health and took action: %s!", healthResponse.Action)
	return nil
}

// handleUpdateHealthFailure tracks consecutive health updates and takes action
// SYSTEM_DOWN if the number of failures exceeds a predefined amount
func (hm *HealthManager) handleUpdateHealthFailure(
	req *protos.HealthRequest,
	err error,
) error {
	glog.Error(err)

	atomic.AddUint32(&hm.consecutiveUpdateFailures, 1)
	if atomic.LoadUint32(&hm.consecutiveUpdateFailures) < hm.config.GetUpdateFailureThreshold() {
		return err
	}

	glog.Warningf("Consecutive update failures exceed threshold %d; Disabling FeG services' diameter connections!",
		hm.config.GetUpdateFailureThreshold())
	actionErr := hm.takeSystemDown(hm.config.GetLocalDisconnectPeriodSecs(), req.HealthStats.ServiceStatus)
	if actionErr != nil {
		glog.Error(actionErr)
		return actionErr
	}

	// SYSTEM_DOWN was successful, so reset failure counter
	atomic.StoreUint32(&hm.consecutiveUpdateFailures, 0)
	hm.prevAction = protos.HealthResponse_SYSTEM_DOWN
	glog.V(2).Infof("Successfully took action: %s", protos.HealthResponse_SYSTEM_DOWN)
	return err
}

// gatherHealthRequest collects FeG services and system health metrics/status and
// fills in a HealthRequest with them
func (hm *HealthManager) gatherHealthRequest() (*protos.HealthRequest, error) {
	serviceStatsMap := make(map[string]*protos.ServiceHealthStats)

	for _, service := range hm.config.GetRequiredServices() {
		serviceStats := collection.CollectServiceStats(service)
		serviceStatsMap[service] = serviceStats
	}
	systemHealthStats, err := collection.CollectSystemStats()

	healthRequest := &protos.HealthRequest{
		HealthStats: &protos.HealthStats{
			SystemStatus:  systemHealthStats,
			ServiceStatus: serviceStatsMap,
		},
	}

	return healthRequest, err
}

// takeSystemDown disables FeG services' for the period specified in the request by calling each
// service's Disable method
func (hm *HealthManager) takeSystemDown(disablePeriod uint32, serviceStats map[string]*protos.ServiceHealthStats) error {
	var allActionErrors []string
	for _, srv := range hm.config.GetRequiredServices() {
		disableReq := &protos.DisableMessage{
			DisablePeriodSecs: uint64(disablePeriod),
		}
		// Only disable available services
		if serviceStats[srv].ServiceState == protos.ServiceHealthStats_UNAVAILABLE {
			continue
		}
		err := service_health.Disable(srv, disableReq)
		if err != nil {
			errMsg := fmt.Sprintf("Error while attempting to take action SYSTEM_DOWN for service %s: %s", srv, err)
			allActionErrors = append(allActionErrors, errMsg)
		}
	}
	if len(allActionErrors) > 0 {
		return fmt.Errorf("Encountered the following errors while taking SYSTEM_DOWN:\n%s\n",
			strings.Join(allActionErrors, "\n"),
		)
	}
	return nil
}

// takeSystemUp enables FeG services' by calling each service's Enable method
func (hm *HealthManager) takeSystemUp() error {
	var allActionErrors []string
	for _, srv := range hm.config.GetRequiredServices() {
		// For SYSTEM_UP, all services should be available
		err := service_health.Enable(srv)
		if err != nil {
			errMsg := fmt.Sprintf("Error while attempting to take action SYSTEM_UP for service %s: %s", srv, err)
			allActionErrors = append(allActionErrors, errMsg)
		}
	}
	if len(allActionErrors) > 0 {
		return fmt.Errorf("Encountered the following errors while taking SYSTEM_UP:\n%s\n",
			strings.Join(allActionErrors, "\n"),
		)
	}
	return nil
}

func validateHealthConfig(config *mconfig.GatewayHealthConfig) error {
	if config == nil {
		return fmt.Errorf("Nil GatewayHealthConfig provided")
	} else if config.GetUpdateIntervalSecs() == 0 {
		return fmt.Errorf("Cannot use 0 secs as update interval")
	} else if config.GetUpdateFailureThreshold() == 0 {
		return fmt.Errorf("Cannot use 0 as consecutive failure threshold")
	}
	return nil
}
