/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package plugin

import (
	"magma/cwf/cloud/go/cwf"
	cwfhandlers "magma/cwf/cloud/go/services/carrier_wifi/obsidian/handlers"
	"magma/cwf/cloud/go/services/carrier_wifi/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	srvconfig "magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/service/serviceregistry"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/metricsd"
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
	}
}

func (*CwfOrchestratorPlugin) GetMconfigBuilders() []configurator.MconfigBuilder {
	return []configurator.MconfigBuilder{
		&Builder{},
	}
}

func (*CwfOrchestratorPlugin) GetMetricsProfiles(metricsConfig *srvconfig.ConfigMap) []metricsd.MetricsProfile {
	return []metricsd.MetricsProfile{}
}

func (*CwfOrchestratorPlugin) GetObsidianHandlers(metricsConfig *srvconfig.ConfigMap) []obsidian.Handler {
	return plugin.FlattenHandlerLists(
		cwfhandlers.GetObsidianHandlers(),
	)
}

func (*CwfOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{}
}
