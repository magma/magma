/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package plugin

import (
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/services/cellular/config"
	cellularh "magma/lte/cloud/go/services/cellular/obsidian/handlers"
	meteringdh "magma/lte/cloud/go/services/meteringd_records/obsidian/handlers"
	policydbh "magma/lte/cloud/go/services/policydb/obsidian/handlers"
	policydbstreamer "magma/lte/cloud/go/services/policydb/streamer"
	subscriberdbh "magma/lte/cloud/go/services/subscriberdb/obsidian/handlers"
	subscriberdbstreamer "magma/lte/cloud/go/services/subscriberdb/streamer"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/service/serviceregistry"
	configregistry "magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/streamer/mconfig/factory"
	"magma/orc8r/cloud/go/services/streamer/providers"
)

// LteOrchestratorPlugin implements OrchestratorPlugin for the LTE module
type LteOrchestratorPlugin struct{}

func (*LteOrchestratorPlugin) GetName() string {
	return lte.ModuleName
}

func (*LteOrchestratorPlugin) GetServices() []registry.ServiceLocation {
	serviceLocations, err := serviceregistry.LoadServiceRegistryConfig(lte.ModuleName)
	if err != nil {
		return []registry.ServiceLocation{}
	}
	return serviceLocations
}

func (*LteOrchestratorPlugin) GetConfigManagers() []configregistry.ConfigManager {
	return []configregistry.ConfigManager{
		&config.CellularNetworkConfigManager{},
		&config.CellularGatewayConfigManager{},
	}
}

func (*LteOrchestratorPlugin) GetSerdes() []serde.Serde {
	return []serde.Serde{}
}

func (*LteOrchestratorPlugin) GetMconfigBuilders() []factory.MconfigBuilder {
	return []factory.MconfigBuilder{
		&config.CellularBuilder{},
	}
}

func (*LteOrchestratorPlugin) GetMetricsProfiles() []metricsd.MetricsProfile {
	return []metricsd.MetricsProfile{}
}

func (*LteOrchestratorPlugin) GetObsidianHandlers() []handlers.Handler {
	return plugin.FlattenHandlerLists(
		cellularh.GetObsidianHandlers(),
		meteringdh.GetObsidianHandlers(),
		policydbh.GetObsidianHandlers(),
		subscriberdbh.GetObsidianHandlers(),
	)
}

func (*LteOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{
		&subscriberdbstreamer.SubscribersProvider{},
		&policydbstreamer.PoliciesProvider{},
		&policydbstreamer.BaseNamesProvider{},
	}
}
