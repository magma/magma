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
	"magma/lte/cloud/go/plugin/stream_provider"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl/legacy_stream_providers"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"
	"magma/orc8r/lib/go/service/serviceregistry"
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
		state.NewStateSerde(lte.ICMPStateType, &lteModels.IcmpStatus{}),

		// Configurator serdes
		configurator.NewNetworkConfigSerde(lte.CellularNetworkType, &lteModels.NetworkCellularConfigs{}),
		configurator.NewNetworkConfigSerde(lte.NetworkSubscriberConfigType, &lteModels.NetworkSubscriberConfig{}),
		configurator.NewNetworkEntityConfigSerde(lte.CellularGatewayType, &lteModels.GatewayCellularConfigs{}),
		configurator.NewNetworkEntityConfigSerde(lte.CellularEnodebType, &lteModels.EnodebConfiguration{}),

		configurator.NewNetworkEntityConfigSerde(lte.PolicyRuleEntityType, &lteModels.PolicyRuleConfig{}),
		configurator.NewNetworkEntityConfigSerde(lte.BaseNameEntityType, &lteModels.BaseNameRecord{}),
		configurator.NewNetworkEntityConfigSerde(subscriberdb.EntityType, &lteModels.LteSubscription{}),

		configurator.NewNetworkEntityConfigSerde(lte.RatingGroupEntityType, &lteModels.RatingGroup{}),

		configurator.NewNetworkEntityConfigSerde(lte.ApnEntityType, &lteModels.ApnConfiguration{}),
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
		handlers.GetHandlers(),
	)
}

func (*LteOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	factory := legacy_stream_providers.LegacyProviderFactory{}
	return []providers.StreamProvider{
		factory.CreateLegacyProvider(lte.SubscriberStreamName, &stream_provider.LteStreamProviderServicer{}),
		factory.CreateLegacyProvider(lte.PolicyStreamName, &stream_provider.LteStreamProviderServicer{}),
		factory.CreateLegacyProvider(lte.BaseNameStreamName, &stream_provider.LteStreamProviderServicer{}),
		factory.CreateLegacyProvider(lte.MappingsStreamName, &stream_provider.LteStreamProviderServicer{}),
		factory.CreateLegacyProvider(lte.NetworkWideRules, &stream_provider.LteStreamProviderServicer{}),
	}
}

func (*LteOrchestratorPlugin) GetStateIndexers() []indexer.Indexer {
	return []indexer.Indexer{}
}
