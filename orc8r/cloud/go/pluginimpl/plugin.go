/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package pluginimpl

import (
	"net/http"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/service/serviceregistry"
	accessdh "magma/orc8r/cloud/go/services/accessd/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/directoryd"
	magmadh "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/collection"
	"magma/orc8r/cloud/go/services/metricsd/confignames"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	metricsdh "magma/orc8r/cloud/go/services/metricsd/obsidian/handlers"
	promeExp "magma/orc8r/cloud/go/services/metricsd/prometheus/exporters"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/streamer/mconfig"
	"magma/orc8r/cloud/go/services/streamer/providers"
	tenantsh "magma/orc8r/cloud/go/services/tenants/obsidian/handlers"
	"magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"

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
		// For checkin_cli.py to test cloud < - > gateway connection
		state.NewStateSerde(state.StringMapSerdeType, &state.StringToStringMap{}),
		// For DirectoryD records
		state.NewStateSerde(orc8r.DirectoryRecordType, &directoryd.DirectoryRecord{}),

		// Device service serdes
		serde.NewBinarySerde(device.SerdeDomain, orc8r.AccessGatewayRecordType, &models.GatewayDevice{}),

		// Config manager serdes
		configurator.NewNetworkConfigSerde(orc8r.DnsdNetworkType, &models.NetworkDNSConfig{}),
		configurator.NewNetworkConfigSerde(orc8r.NetworkFeaturesConfig, &models.NetworkFeatures{}),

		configurator.NewNetworkEntityConfigSerde(orc8r.MagmadGatewayType, &models.MagmadGatewayConfigs{}),
		configurator.NewNetworkEntityConfigSerde(orc8r.UpgradeReleaseChannelEntityType, &models.ReleaseChannel{}),
		configurator.NewNetworkEntityConfigSerde(orc8r.UpgradeTierEntityType, &models.Tier{}),
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
		// v1 handlers
		magmadh.GetObsidianHandlers(),
		metricsdh.GetObsidianHandlers(metricsConfig),
		handlers.GetObsidianHandlers(),
		tenantsh.GetObsidianHandlers(),
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
	ProfileNameExportAll  = "exportall"
)

func getMetricsProfiles(metricsConfig *config.ConfigMap) []metricsd.MetricsProfile {

	// Controller profile - 1 collector for each service
	allServices := registry.ListControllerServices()
	deviceMetricsCollectors := []collection.MetricCollector{&collection.DiskUsageMetricCollector{}, &collection.ProcMetricsCollector{}}
	controllerCollectors := make([]collection.MetricCollector, 0, len(allServices)+len(deviceMetricsCollectors))
	for _, srv := range allServices {
		controllerCollectors = append(controllerCollectors, collection.NewCloudServiceMetricCollector(srv))
	}
	for _, metricCollector := range deviceMetricsCollectors {
		controllerCollectors = append(controllerCollectors, metricCollector)
	}

	// Prometheus profile - Exports all service metric to Prometheus
	prometheusAddresses := metricsConfig.GetRequiredStringArrayParam(confignames.PrometheusPushAddresses)
	prometheusCustomPushExporter := promeExp.NewCustomPushExporter(prometheusAddresses)
	prometheusProfile := metricsd.MetricsProfile{
		Name:       ProfileNamePrometheus,
		Collectors: controllerCollectors,
		Exporters:  []exporters.Exporter{prometheusCustomPushExporter},
	}

	// ExportAllProfile - Exports to all exporters
	exportAllProfile := metricsd.MetricsProfile{
		Name:       ProfileNameExportAll,
		Collectors: controllerCollectors,
		Exporters:  []exporters.Exporter{prometheusCustomPushExporter},
	}

	return []metricsd.MetricsProfile{
		prometheusProfile,
		exportAllProfile,
	}
}
