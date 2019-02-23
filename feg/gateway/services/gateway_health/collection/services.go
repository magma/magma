/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package collection provides functions used by the health manager to collect
// health related metrics for FeG services and the system
package collection

import (
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/service_health"
)

// CollectServiceStats fills out the ServiceHealthStats proto for the provided service
// If the service cannot be reached, the service state is listed as UNAVAILABLE
func CollectServiceStats(serviceType string) *protos.ServiceHealthStats {
	healthStatus, err := service_health.GetHealthStatus(serviceType)
	if err != nil {
		if healthStatus != nil {
			return &protos.ServiceHealthStats{
				ServiceState:        protos.ServiceHealthStats_AVAILABLE,
				ServiceHealthStatus: healthStatus,
			}
		}
		return &protos.ServiceHealthStats{
			ServiceState: protos.ServiceHealthStats_UNAVAILABLE,
			ServiceHealthStatus: &protos.HealthStatus{
				Health:        protos.HealthStatus_UNHEALTHY,
				HealthMessage: "Service unavailable",
			},
		}
	}
	return &protos.ServiceHealthStats{
		ServiceState:        protos.ServiceHealthStats_AVAILABLE,
		ServiceHealthStatus: healthStatus,
	}
}
