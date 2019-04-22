/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package pluginimpl

import (
	obsidianh "magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/obsidian/handlers/hello"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/service/serviceregistry"
	accessdh "magma/orc8r/cloud/go/services/accessd/obsidian/handlers"
	checkinh "magma/orc8r/cloud/go/services/checkind/obsidian/handlers"
	dnsdconfig "magma/orc8r/cloud/go/services/dnsd/config"
	dnsdh "magma/orc8r/cloud/go/services/dnsd/obsidian/handlers"
	magmadconfig "magma/orc8r/cloud/go/services/magmad/config"
	magmadh "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/collection"
	"magma/orc8r/cloud/go/services/metricsd/confignames"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	graphite_exp "magma/orc8r/cloud/go/services/metricsd/graphite/exporters"
	metricsdh "magma/orc8r/cloud/go/services/metricsd/obsidian/handlers"
	promo_exp "magma/orc8r/cloud/go/services/metricsd/prometheus/exporters"
	"magma/orc8r/cloud/go/services/streamer/mconfig"
	"magma/orc8r/cloud/go/services/streamer/mconfig/factory"
	"magma/orc8r/cloud/go/services/streamer/providers"
	upgradeh "magma/orc8r/cloud/go/services/upgrade/obsidian/handlers"
)

// BaseOrchestratorPlugin is the OrchestratorPlugin for the orc8r module
type BaseOrchestratorPlugin struct{}

func (*BaseOrchestratorPlugin) GetName() string {
	return orc8r.ModuleName
}

func (*BaseOrchestratorPlugin) GetServices() []registry.ServiceLocation {
	serviceLocations, err := serviceregistry.LoadServiceRegistryConfig(orc8r.ModuleName)
	if err != nil {
		return []registry.ServiceLocation{}
	}
	return serviceLocations
}

func (*BaseOrchestratorPlugin) GetSerdes() []serde.Serde {
	return []serde.Serde{
		// State service serdes
		&CheckinRequestSerde{},

		// Inventory service serdes
		&GatewayRecordSerde{},

		// Config manager serdes
		&magmadconfig.MagmadGatewayConfigManager{},
		&dnsdconfig.DnsNetworkConfigManager{},
	}
}

func (*BaseOrchestratorPlugin) GetMconfigBuilders() []factory.MconfigBuilder {
	return []factory.MconfigBuilder{
		// magmad
		&magmadconfig.MagmadMconfigBuilder{},
		// dnsd
		&dnsdconfig.DnsdMconfigBuilder{},
	}
}

func (*BaseOrchestratorPlugin) GetMetricsProfiles(metricsConfig *config.ConfigMap) []metricsd.MetricsProfile {
	return getMetricsProfiles(metricsConfig)
}

func (*BaseOrchestratorPlugin) GetObsidianHandlers(metricsConfig *config.ConfigMap) []obsidianh.Handler {
	return plugin.FlattenHandlerLists(
		accessdh.GetObsidianHandlers(),
		checkinh.GetObsidianHandlers(),
		dnsdh.GetObsidianHandlers(),
		magmadh.GetObsidianHandlers(),
		metricsdh.GetObsidianHandlers(metricsConfig),
		upgradeh.GetObsidianHandlers(),
		hello.GetObsidianHandlers(),
	)
}

func (*BaseOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{
		mconfig.GetProvider(),
		mconfig.GetViewProvider(),
	}
}

const (
	ProfileNamePrometheus = "prometheus"
	ProfileNameGraphite   = "graphite"
	ProfileNameDefault    = "default"
)

func getMetricsProfiles(metricsConfig *config.ConfigMap) []metricsd.MetricsProfile {

	// Controller profile - 1 collector for each service
	allServices := registry.ListControllerServices()
	controllerCollectors := make([]collection.MetricCollector, 0, len(allServices)+1)
	for _, srv := range allServices {
		controllerCollectors = append(controllerCollectors, collection.NewCloudServiceMetricCollector(srv))
	}
	controllerCollectors = append(controllerCollectors, &collection.DiskUsageMetricCollector{})

	prometheusPushAddress := metricsConfig.GetRequiredStringParam(confignames.PrometheusPushgatewayAddress)
	prometheusPushExporter := promo_exp.NewPrometheusPushExporter(prometheusPushAddress)
	// Prometheus profile - Exports all service metric to Prometheus
	prometheusProfile := metricsd.MetricsProfile{
		Name:       ProfileNamePrometheus,
		Collectors: controllerCollectors,
		Exporters:  []exporters.Exporter{prometheusPushExporter},
	}

	graphiteAddress := metricsConfig.GetRequiredStringParam(confignames.GraphiteAddress)
	graphiteReceivePort := metricsConfig.GetRequiredIntParam(confignames.GraphiteReceivePort)
	graphiteExporter := graphite_exp.NewGraphiteExporter(graphiteAddress, graphiteReceivePort)
	// Graphite profile - Exports all service metrics to Graphite
	graphiteProfile := metricsd.MetricsProfile{
		Name:       ProfileNameGraphite,
		Collectors: controllerCollectors,
		Exporters:  []exporters.Exporter{},
	}

	defaultProfile := metricsd.MetricsProfile{
		Name:       ProfileNameDefault,
		Collectors: controllerCollectors,
		Exporters:  []exporters.Exporter{graphiteExporter, prometheusPushExporter},
	}

	return []metricsd.MetricsProfile{
		prometheusProfile,
		graphiteProfile,
		defaultProfile,
	}
}
