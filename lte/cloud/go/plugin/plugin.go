/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package plugin

import (
	"magma/lte/cloud/go/lte"
	lte_service "magma/lte/cloud/go/services/lte"
	lte_handlers "magma/lte/cloud/go/services/lte/obsidian/handlers"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	policydb_handlers "magma/lte/cloud/go/services/policydb/obsidian/handlers"
	policydb_models "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb"
	subscriberdb_handlers "magma/lte/cloud/go/services/subscriberdb/obsidian/handlers"
	subscriberdb_models "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"
)

// LteOrchestratorPlugin implements OrchestratorPlugin for the LTE module
type LteOrchestratorPlugin struct{}

func (*LteOrchestratorPlugin) GetName() string {
	return lte.ModuleName
}

func (*LteOrchestratorPlugin) GetServices() []registry.ServiceLocation {
	serviceLocations, err := registry.LoadServiceRegistryConfig(lte.ModuleName)
	if err != nil {
		return []registry.ServiceLocation{}
	}
	return serviceLocations
}

func (*LteOrchestratorPlugin) GetSerdes() []serde.Serde {
	return []serde.Serde{
		state.NewStateSerde(lte.EnodebStateType, &lte_models.EnodebState{}),
		state.NewStateSerde(lte.ICMPStateType, &subscriberdb_models.IcmpStatus{}),

		// AGW state messages which use arbitrary untyped JSON serdes because
		// they're defined/used as protos in the AGW codebase
		state.NewStateSerde(lte.MMEStateType, &state.ArbitraryJSON{}),
		state.NewStateSerde(lte.SPGWStateType, &state.ArbitraryJSON{}),
		state.NewStateSerde(lte.S1APStateType, &state.ArbitraryJSON{}),
		state.NewStateSerde(lte.MobilitydStateType, &state.ArbitraryJSON{}),

		// Configurator serdes
		configurator.NewNetworkConfigSerde(lte.CellularNetworkType, &lte_models.NetworkCellularConfigs{}),
		configurator.NewNetworkConfigSerde(lte.NetworkSubscriberConfigType, &policydb_models.NetworkSubscriberConfig{}),
		configurator.NewNetworkEntityConfigSerde(lte.CellularGatewayType, &lte_models.GatewayCellularConfigs{}),
		configurator.NewNetworkEntityConfigSerde(lte.CellularEnodebType, &lte_models.EnodebConfiguration{}),

		configurator.NewNetworkEntityConfigSerde(lte.PolicyRuleEntityType, &policydb_models.PolicyRuleConfig{}),
		configurator.NewNetworkEntityConfigSerde(lte.BaseNameEntityType, &policydb_models.BaseNameRecord{}),
		configurator.NewNetworkEntityConfigSerde(subscriberdb.EntityType, &subscriberdb_models.LteSubscription{}),

		configurator.NewNetworkEntityConfigSerde(lte.RatingGroupEntityType, &policydb_models.RatingGroup{}),

		configurator.NewNetworkEntityConfigSerde(lte.ApnEntityType, &lte_models.ApnConfiguration{}),
	}
}

func (*LteOrchestratorPlugin) GetMconfigBuilders() []configurator.MconfigBuilder {
	return []configurator.MconfigBuilder{
		&Builder{},
	}
}

func (*LteOrchestratorPlugin) GetMetricsProfiles(metricsConfig *config.ConfigMap) []metricsd.MetricsProfile {
	return []metricsd.MetricsProfile{}
}

func (*LteOrchestratorPlugin) GetObsidianHandlers(metricsConfig *config.ConfigMap) []obsidian.Handler {
	return plugin.FlattenHandlerLists(
		lte_handlers.GetHandlers(),
		policydb_handlers.GetHandlers(),
		subscriberdb_handlers.GetHandlers(),
	)
}

func (*LteOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{
		providers.NewRemoteProvider(lte_service.ServiceName, lte.SubscriberStreamName),
		providers.NewRemoteProvider(lte_service.ServiceName, lte.PolicyStreamName),
		providers.NewRemoteProvider(lte_service.ServiceName, lte.BaseNameStreamName),
		providers.NewRemoteProvider(lte_service.ServiceName, lte.MappingsStreamName),
		providers.NewRemoteProvider(lte_service.ServiceName, lte.NetworkWideRulesStreamName),
		providers.NewRemoteProvider(lte_service.ServiceName, lte.RatingGroupStreamName),
	}
}

func (*LteOrchestratorPlugin) GetStateIndexers() []indexer.Indexer {
	return []indexer.Indexer{}
}
