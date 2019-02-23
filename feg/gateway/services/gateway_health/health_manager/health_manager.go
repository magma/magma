/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
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

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/service_health"
	"magma/feg/gateway/services/gateway_health"
	"magma/feg/gateway/services/gateway_health/collection"

	"github.com/golang/glog"
)

const cloudConnectionDisablePeriodSecs = 10
const fegConnectionDisablePeriodSecs = 1
const consecutiveFailuresThreshold = 3

// TODO: Make required services configurable
var requiredServices = [...]string{registry.S6A_PROXY, registry.SESSION_PROXY}

type HealthManager struct {
	cloudReg                  registry.CloudRegistry
	consecutiveUpdateFailures uint64
	prevAction                protos.HealthResponse_RequestedAction
}

func NewHealthManager(cloudReg registry.CloudRegistry) *HealthManager {
	return &HealthManager{
		cloudReg:                  cloudReg,
		consecutiveUpdateFailures: 0,
		prevAction:                protos.HealthResponse_NONE,
	}
}

// SendHealthUpdate collects Gateway Service and System Health Status and sends
// them to the cloud health service. It awaits a response from the cloud and
// applies any action requested from the cloud (e.g. SYSTEM_DOWN)
func (hm *HealthManager) SendHealthUpdate() error {
	healthRequest, err := gatherHealthRequest()
	if err != nil {
		glog.Error(err)
	}
	healthResponse, err := gateway_health.UpdateHealth(hm.cloudReg, healthRequest)
	if err != nil {
		return hm.handleUpdateHealthFailure(healthRequest, err)
	}

	// Update was successful, so reset consecutive failure counter
	atomic.StoreUint64(&hm.consecutiveUpdateFailures, 0)

	switch healthResponse.Action {
	case protos.HealthResponse_NONE:
	case protos.HealthResponse_SYSTEM_UP:
		err = takeSystemUp()
	case protos.HealthResponse_SYSTEM_DOWN:
		disablePeriod := cloudConnectionDisablePeriodSecs
		if hm.prevAction == protos.HealthResponse_SYSTEM_DOWN {
			disablePeriod = 0
		}
		err = takeSystemDown(uint64(disablePeriod), healthRequest.HealthStats.ServiceStatus)
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

	atomic.AddUint64(&hm.consecutiveUpdateFailures, 1)
	if atomic.LoadUint64(&hm.consecutiveUpdateFailures) < consecutiveFailuresThreshold {
		return err
	}

	glog.V(2).Info("Consecutive update failures exceed threshold; Disabling FeG services' diameter connections...")
	actionErr := takeSystemDown(fegConnectionDisablePeriodSecs, req.HealthStats.ServiceStatus)
	if actionErr != nil {
		glog.Error(actionErr)
		return actionErr
	}

	// SYSTEM_DOWN was successful, so reset failure counter
	atomic.StoreUint64(&hm.consecutiveUpdateFailures, 0)
	hm.prevAction = protos.HealthResponse_SYSTEM_DOWN
	glog.V(2).Infof("Successfully took action: %s", protos.HealthResponse_SYSTEM_DOWN)
	return err
}

// gatherHealthRequest collects FeG services and system health metrics/status and
// fills in a HealthRequest with them
func gatherHealthRequest() (*protos.HealthRequest, error) {
	serviceStatsMap := make(map[string]*protos.ServiceHealthStats)

	for _, service := range requiredServices {
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
func takeSystemDown(disablePeriod uint64, serviceStats map[string]*protos.ServiceHealthStats) error {
	var allActionErrors []string
	for _, srv := range requiredServices {
		disableReq := &protos.DisableMessage{
			DisablePeriodSecs: disablePeriod,
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
func takeSystemUp() error {
	var allActionErrors []string
	for _, srv := range requiredServices {
		// For SYSTEM_UP, all services should be available
		err := service_health.Enable(srv)
		if err != nil {
			errMsg := fmt.Sprintf("Error while attempting to take action SYSTEM_UP for service %s: %s", srv, err)
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
