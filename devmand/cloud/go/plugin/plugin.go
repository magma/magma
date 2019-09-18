/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package plugin

import (
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/service/serviceregistry"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"orc8r/devmand/cloud/go/devmand"
	dmh "orc8r/devmand/cloud/go/services/devmand/obsidian/handlers"
	"orc8r/devmand/cloud/go/services/devmand/obsidian/models"
)

// DevmandOrchestratorPlugin is the orchestrator plugin for devmand
type DevmandOrchestratorPlugin struct{}

// GetName gets the name of the devmand module
func (*DevmandOrchestratorPlugin) GetName() string {
	return devmand.ModuleName
}

// GetServices gets the devmand service locations
func (*DevmandOrchestratorPlugin) GetServices() []registry.ServiceLocation {
	serviceLocations, err := serviceregistry.LoadServiceRegistryConfig(devmand.ModuleName)
	if err != nil {
		return []registry.ServiceLocation{}
	}
	return serviceLocations
}

// GetSerdes gets the devmand serializers and deserializers
func (*DevmandOrchestratorPlugin) GetSerdes() []serde.Serde {
	return []serde.Serde{
		configurator.NewNetworkEntityConfigSerde(devmand.DevmandGatewayType, &models.GatewayDevmandConfigs{}),
		configurator.NewNetworkEntityConfigSerde(devmand.DeviceType, &models.ManagedDevice{}),
	}
}

func (*DevmandOrchestratorPlugin) GetMconfigBuilders() []configurator.MconfigBuilder {
	return []configurator.MconfigBuilder{
		&Builder{},
	}
}

// GetMetricsProfiles gets the metricsd profiles
func (*DevmandOrchestratorPlugin) GetMetricsProfiles(metricsConfig *config.ConfigMap) []metricsd.MetricsProfile {
	return []metricsd.MetricsProfile{}
}

// GetObsidianHandlers gets the devmand obsidian handlers
func (*DevmandOrchestratorPlugin) GetObsidianHandlers(metricsConfig *config.ConfigMap) []obsidian.Handler {
	return plugin.FlattenHandlerLists(
		dmh.GetObsidianHandlers(),
		GetHandlers(),
	)
}

// GetStreamerProviders gets the stream providers
func (*DevmandOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{}
}
