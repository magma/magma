/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package pluginimpl

import (
	"fmt"
	"strconv"
	"strings"

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
	"magma/orc8r/cloud/go/services/configurator"
	configuratorh "magma/orc8r/cloud/go/services/configurator/obsidian/handlers"
	"magma/orc8r/cloud/go/services/device"
	dnsdconfig "magma/orc8r/cloud/go/services/dnsd/config"
	dnsdh "magma/orc8r/cloud/go/services/dnsd/obsidian/handlers"
	dnsdmodels "magma/orc8r/cloud/go/services/dnsd/obsidian/models"
	magmadconfig "magma/orc8r/cloud/go/services/magmad/config"
	magmadh "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
	magmadmodels "magma/orc8r/cloud/go/services/magmad/obsidian/models"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/collection"
	"magma/orc8r/cloud/go/services/metricsd/confignames"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	graphite_exp "magma/orc8r/cloud/go/services/metricsd/graphite/exporters"
	metricsdh "magma/orc8r/cloud/go/services/metricsd/obsidian/handlers"
	promo_exp "magma/orc8r/cloud/go/services/metricsd/prometheus/exporters"
	stateh "magma/orc8r/cloud/go/services/state/obsidian/handlers"
	"magma/orc8r/cloud/go/services/streamer/mconfig"
	"magma/orc8r/cloud/go/services/streamer/mconfig/factory"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/cloud/go/services/upgrade"
	upgradeh "magma/orc8r/cloud/go/services/upgrade/obsidian/handlers"
	upgrademodels "magma/orc8r/cloud/go/services/upgrade/obsidian/models"
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
		&GatewayStatusSerde{},

		// Inventory service serdes
		serde.NewBinarySerde(device.SerdeDomain, orc8r.AccessGatewayRecordType, &magmadmodels.AccessGatewayRecord{}),

		// Config manager serdes
		configurator.NewNetworkConfigSerde(orc8r.DnsdNetworkType, &dnsdmodels.NetworkDNSConfig{}),
		configurator.NewNetworkConfigSerde(orc8r.NetworkFeaturesConfig, &magmadmodels.NetworkFeatures{}),

		configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &magmadmodels.MagmadGatewayConfig{}),
		configurator.NewNetworkEntityConfigSerde(upgrade.UpgradeReleaseChannelEntityType, &upgrademodels.ReleaseChannel{}),
		configurator.NewNetworkEntityConfigSerde(orc8r.UpgradeTierEntityType, &upgrademodels.Tier{}),

		// Legacy config manager serdes
		&magmadconfig.MagmadGatewayConfigManager{},
		&dnsdconfig.DnsNetworkConfigManager{},
	}
}

func (*BaseOrchestratorPlugin) GetLegacyMconfigBuilders() []factory.MconfigBuilder {
	return []factory.MconfigBuilder{
		// magmad
		&magmadconfig.MagmadMconfigBuilder{},
		// dnsd
		&dnsdconfig.DnsdMconfigBuilder{},
	}
}

func (*BaseOrchestratorPlugin) GetMconfigBuilders() []configurator.MconfigBuilder {
	return []configurator.MconfigBuilder{
		&BaseOrchestratorMconfigBuilder{},
		&DnsdMconfigBuilder{},
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
		stateh.GetObsidianHandlers(),
		configuratorh.GetObsidianHandlers(),
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

	// Prometheus profile - Exports all service metric to Prometheus
	prometheusAddresses := metricsConfig.GetRequiredStringArrayParam(confignames.PrometheusPushAddresses)
	prometheusCustomPushExporter := promo_exp.NewCustomPushExporter(prometheusAddresses)
	prometheusProfile := metricsd.MetricsProfile{
		Name:       ProfileNamePrometheus,
		Collectors: controllerCollectors,
		Exporters:  []exporters.Exporter{prometheusCustomPushExporter},
	}

	graphiteExportAddresses := metricsConfig.GetRequiredStringArrayParam(confignames.GraphiteExportAddresses)
	var graphiteAddresses []graphite_exp.Address
	for _, address := range graphiteExportAddresses {
		portIdx := strings.LastIndex(address, ":")
		portStr := address[portIdx+1:]
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			panic(fmt.Errorf("graphite address improperly formed: %s\n", address))
		}
		graphiteAddresses = append(graphiteAddresses, graphite_exp.Address{Host: address[:portIdx], Port: portInt})
	}

	graphiteExporter := graphite_exp.NewGraphiteExporter(graphiteAddresses)
	// Graphite profile - Exports all service metrics to Graphite
	graphiteProfile := metricsd.MetricsProfile{
		Name:       ProfileNameGraphite,
		Collectors: controllerCollectors,
		Exporters:  []exporters.Exporter{graphiteExporter},
	}

	defaultProfile := metricsd.MetricsProfile{
		Name:       ProfileNameDefault,
		Collectors: controllerCollectors,
		Exporters:  []exporters.Exporter{prometheusCustomPushExporter, graphiteExporter},
	}

	return []metricsd.MetricsProfile{
		prometheusProfile,
		graphiteProfile,
		defaultProfile,
	}
}
