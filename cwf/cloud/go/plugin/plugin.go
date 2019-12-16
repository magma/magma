/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package plugin

import (
	"magma/cwf/cloud/go/cwf"
	"magma/cwf/cloud/go/plugin/handlers"
	"magma/cwf/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/service/serviceregistry"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/streamer/providers"
)

// CwfOrchestratorPlugin implements OrchestratorPlugin for the CWF module
type CwfOrchestratorPlugin struct{}

func (*CwfOrchestratorPlugin) GetName() string {
	return cwf.ModuleName
}

func (*CwfOrchestratorPlugin) GetServices() []registry.ServiceLocation {
	serviceLocations, err := serviceregistry.LoadServiceRegistryConfig(cwf.ModuleName)
	if err != nil {
		return []registry.ServiceLocation{}
	}
	return serviceLocations
}

func (*CwfOrchestratorPlugin) GetSerdes() []serde.Serde {
	return []serde.Serde{
		configurator.NewNetworkConfigSerde(cwf.CwfNetworkType, &models.NetworkCarrierWifiConfigs{}),
		configurator.NewNetworkEntityConfigSerde(cwf.CwfGatewayType, &models.GatewayCwfConfigs{}),
		state.NewStateSerde(cwf.CwfSubscriberDirectoryType, &models.CwfSubscriberDirectoryRecord{}),
	}
}

func (*CwfOrchestratorPlugin) GetMconfigBuilders() []configurator.MconfigBuilder {
	return []configurator.MconfigBuilder{
		&Builder{},
	}
}

func (*CwfOrchestratorPlugin) GetMetricsProfiles(metricsConfig *config.ConfigMap) []metricsd.MetricsProfile {
	return []metricsd.MetricsProfile{}
}

func (*CwfOrchestratorPlugin) GetObsidianHandlers(metricsConfig *config.ConfigMap) []obsidian.Handler {
	return plugin.FlattenHandlerLists(
		handlers.GetHandlers(),
	)
}

func (*CwfOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{}
}
