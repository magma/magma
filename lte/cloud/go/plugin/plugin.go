/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package plugin

import (
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/plugin/handlers"
	lteModels "magma/lte/cloud/go/plugin/models"
	cellularh "magma/lte/cloud/go/services/cellular/obsidian/handlers"
	meteringdh "magma/lte/cloud/go/services/meteringd_records/obsidian/handlers"
	policydbh "magma/lte/cloud/go/services/policydb/obsidian/handlers"
	policydbstreamer "magma/lte/cloud/go/services/policydb/streamer"
	"magma/lte/cloud/go/services/subscriberdb"
	subscriberdbh "magma/lte/cloud/go/services/subscriberdb/obsidian/handlers"
	models3 "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	subscriberdbstreamer "magma/lte/cloud/go/services/subscriberdb/streamer"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	srvconfig "magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/service/serviceregistry"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/state"
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

func (*LteOrchestratorPlugin) GetSerdes() []serde.Serde {
	return []serde.Serde{
		state.NewStateSerde(lte.EnodebStateType, &lteModels.EnodebState{}),
		state.NewStateSerde(lte.SubscriberStateType, &models3.SubscriberState{}),

		// Configurator serdes
		configurator.NewNetworkConfigSerde(lte.CellularNetworkType, &lteModels.NetworkCellularConfigs{}),
		configurator.NewNetworkEntityConfigSerde(lte.CellularGatewayType, &lteModels.GatewayCellularConfigs{}),
		configurator.NewNetworkEntityConfigSerde(lte.CellularEnodebType, &lteModels.EnodebConfiguration{}),

		configurator.NewNetworkEntityConfigSerde(lte.PolicyRuleEntityType, &lteModels.PolicyRuleConfig{}),
		configurator.NewNetworkEntityConfigSerde(lte.BaseNameEntityType, &lteModels.BaseNameRecord{}),
		configurator.NewNetworkEntityConfigSerde(subscriberdb.EntityType, &lteModels.LteSubscription{}),
	}
}

func (*LteOrchestratorPlugin) GetMconfigBuilders() []configurator.MconfigBuilder {
	return []configurator.MconfigBuilder{
		&Builder{},
	}
}

func (*LteOrchestratorPlugin) GetMetricsProfiles(metricsConfig *srvconfig.ConfigMap) []metricsd.MetricsProfile {
	return []metricsd.MetricsProfile{}
}

func (*LteOrchestratorPlugin) GetObsidianHandlers(metricsConfig *srvconfig.ConfigMap) []obsidian.Handler {
	return plugin.FlattenHandlerLists(
		cellularh.GetObsidianHandlers(),
		meteringdh.GetObsidianHandlers(),
		policydbh.GetObsidianHandlers(),
		subscriberdbh.GetObsidianHandlers(),
		handlers.GetHandlers(),
	)
}

func (*LteOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{
		&subscriberdbstreamer.SubscribersProvider{},
		&policydbstreamer.PoliciesProvider{},
		&policydbstreamer.BaseNamesProvider{},
	}
}
