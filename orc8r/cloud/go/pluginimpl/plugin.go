/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package pluginimpl

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/service/serviceregistry"
	accessdh "magma/orc8r/cloud/go/services/accessd/obsidian/handlers"
	checkinh "magma/orc8r/cloud/go/services/checkind/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
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
	"magma/orc8r/cloud/go/services/state"
	stateh "magma/orc8r/cloud/go/services/state/obsidian/handlers"
	"magma/orc8r/cloud/go/services/streamer/mconfig"
	"magma/orc8r/cloud/go/services/streamer/mconfig/factory"
	"magma/orc8r/cloud/go/services/streamer/providers"
	upgradeh "magma/orc8r/cloud/go/services/upgrade/obsidian/handlers"
	upgrademodels "magma/orc8r/cloud/go/services/upgrade/obsidian/models"

	"github.com/labstack/echo"
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
		state.NewStateSerde(orc8r.GatewayStateType, &models.GatewayStatus{}),

		// Device service serdes
		serde.NewBinarySerde(device.SerdeDomain, orc8r.AccessGatewayRecordType, &models.GatewayDevice{}),

		// Config manager serdes
		configurator.NewNetworkConfigSerde(orc8r.DnsdNetworkType, &models.NetworkDNSConfig{}),
		configurator.NewNetworkConfigSerde(orc8r.NetworkFeaturesConfig, &models.NetworkFeatures{}),

		configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models.MagmadGatewayConfigs{}),
		configurator.NewNetworkEntityConfigSerde(orc8r.UpgradeReleaseChannelEntityType, &upgrademodels.ReleaseChannel{}),
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

func (*BaseOrchestratorPlugin) GetObsidianHandlers(metricsConfig *config.ConfigMap) []obsidian.Handler {
	return plugin.FlattenHandlerLists(
		// v0 handlers
		accessdh.GetObsidianHandlers(),
		checkinh.GetObsidianHandlers(),
		dnsdh.GetObsidianHandlers(),
		magmadh.GetObsidianHandlers(),
		metricsdh.GetObsidianHandlers(metricsConfig),
		upgradeh.GetObsidianHandlers(),
		stateh.GetObsidianHandlers(),
		// v1 handlers
		handlers.GetObsidianHandlers(),
		[]obsidian.Handler{{
			Path:    "/",
			Methods: obsidian.GET,
			HandlerFunc: func(c echo.Context) error {
				return c.JSON(
					http.StatusOK,
					"hello",
				)
			},
		}},
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
	ProfileNameExportAll  = "exportall"
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

	// No-op graphite exporter if graphite parameters are not set
	graphiteExportAddresses, _ := metricsConfig.GetStringArrayParam(confignames.GraphiteExportAddresses)
	var graphiteAddresses []graphite_exp.Address
	for _, address := range graphiteExportAddresses {
		portIdx := strings.LastIndex(address, ":")
		portStr := address[portIdx+1:]
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			panic(fmt.Errorf("graphite address improperly formed: %s\n", address))
		}
		graphiteAddresses = append(graphiteAddresses, graphite_exp.NewAddress(address[:portIdx], portInt))
	}
	graphiteExporter := graphite_exp.NewGraphiteExporter(graphiteAddresses)

	// Graphite profile - Exports all service metrics to Graphite
	graphiteProfile := metricsd.MetricsProfile{
		Name:       ProfileNameGraphite,
		Collectors: controllerCollectors,
		Exporters:  []exporters.Exporter{graphiteExporter},
	}

	// ExportAllProfile - Exports to both graphite and prometheus
	exportAllProfile := metricsd.MetricsProfile{
		Name:       ProfileNameExportAll,
		Collectors: controllerCollectors,
		Exporters:  []exporters.Exporter{prometheusCustomPushExporter, graphiteExporter},
	}

	return []metricsd.MetricsProfile{
		prometheusProfile,
		graphiteProfile,
		exportAllProfile,
	}
}
