/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package pluginimpl

import (
	"os"
	"time"

	obsidianh "magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/obsidian/handlers/hello"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/service/serviceregistry"
	accessdh "magma/orc8r/cloud/go/services/accessd/obsidian/handlers"
	checkinh "magma/orc8r/cloud/go/services/checkind/obsidian/handlers"
	configregistry "magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/config/streaming"
	configstreamer "magma/orc8r/cloud/go/services/config/streaming/providers"
	dnsdconfig "magma/orc8r/cloud/go/services/dnsd/config"
	dnsdh "magma/orc8r/cloud/go/services/dnsd/obsidian/handlers"
	magmadconfig "magma/orc8r/cloud/go/services/magmad/config"
	magmadh "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/collection"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	metricsdh "magma/orc8r/cloud/go/services/metricsd/obsidian/handlers"
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

func (*BaseOrchestratorPlugin) GetConfigManagers() []configregistry.ConfigManager {
	return []configregistry.ConfigManager{
		// magmad
		&magmadconfig.MagmadGatewayConfigManager{},
		// dnsd
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

func (*BaseOrchestratorPlugin) GetMconfigStreamers() []streaming.MconfigStreamer {
	return []streaming.MconfigStreamer{
		// magmad
		&magmadconfig.MagmadStreamer{},
		// dnsd
		&dnsdconfig.DnsdStreamer{},
	}
}

func (*BaseOrchestratorPlugin) GetMetricsProfiles() []metricsd.MetricsProfile {
	return getMetricsProfiles()
}

func (*BaseOrchestratorPlugin) GetObsidianHandlers() []obsidianh.Handler {
	return plugin.FlattenHandlerLists(
		accessdh.GetObsidianHandlers(),
		checkinh.GetObsidianHandlers(),
		dnsdh.GetObsidianHandlers(),
		magmadh.GetObsidianHandlers(),
		metricsdh.GetObsidianHandlers(),
		upgradeh.GetObsidianHandlers(),
		hello.GetObsidianHandlers(),
	)
}

func (*BaseOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{
		configstreamer.GetProvider(),
		mconfig.GetProvider(),
	}
}

const (
	ProfileNameController = "controller"
	ProfileNameSys        = "sys"
	ProfileNameKafka      = "kafka"
	ProfileNamePrometheus = "prometheus"

	OdsMetricsExportInterval = time.Second * 15
	// a sample is 10 bytes
	// right now 50 metrics from each gateway, 35 metrics from each cloud
	// service per minute assume we support 100 metrics from each gateway,
	// 70 metrics from each cloud service. with 1000 gws, we will have 100000
	// metrics per minute from gws. with 30 cloud services,
	// we have 2100 from cloud.
	// this needs 10 * 102100 = 1021000 B
	OdsMetricsQueueLength = 102000
)

func getMetricsProfiles() []metricsd.MetricsProfile {
	// Sys profile - collectors for disk usage and metricsd
	sysProfile := metricsd.MetricsProfile{
		Name: ProfileNameSys,
		Collectors: []collection.MetricCollector{
			&collection.DiskUsageMetricCollector{},
			collection.NewCloudServiceMetricCollector(metricsd.ServiceName),
		},
		Exporters: []exporters.Exporter{createODSExporter()},
	}

	// Controller profile - 1 collector for each service
	allServices := registry.ListControllerServices()
	controllerCollectors := make([]collection.MetricCollector, 0, len(allServices)+1)
	for _, srv := range allServices {
		controllerCollectors = append(controllerCollectors, collection.NewCloudServiceMetricCollector(srv))
	}
	controllerCollectors = append(controllerCollectors, &collection.DiskUsageMetricCollector{})
	controllerProfile := metricsd.MetricsProfile{
		Name:       ProfileNameController,
		Collectors: controllerCollectors,
		Exporters:  []exporters.Exporter{createODSExporter()},
	}

	odsExporter := createODSExporter()
	// Kakfa profile
	kafkaProfile := metricsd.MetricsProfile{
		Name: ProfileNameKafka,
		Collectors: []collection.MetricCollector{
			&collection.DiskUsageMetricCollector{},
			collection.NewCloudServiceMetricCollector(metricsd.ServiceName),
			collection.NewKafkaConnectCollector("magma-connector", nil),
			&collection.SystemdStatusMetricCollector{
				ServiceNames: []string{"magma@kafka", "magma@zookeeper", "magma@kafka_connect"},
			},
		},
		Exporters: []exporters.Exporter{odsExporter},
	}

	// Prometheus profile - controller profile except using prometheus to
	// export
	prometheusProfile := metricsd.MetricsProfile{
		Name:       ProfileNamePrometheus,
		Collectors: controllerCollectors,
		Exporters:  []exporters.Exporter{exporters.NewPrometheusExporter(exporters.DefaultPrometheusConfig), exporters.NewPrometheusPushExporter(), odsExporter},
	}

	return []metricsd.MetricsProfile{
		sysProfile,
		controllerProfile,
		kafkaProfile,
		prometheusProfile,
	}
}

func createODSExporter() exporters.Exporter {
	return exporters.NewODSExporter(
		os.Getenv("METRIC_EXPORT_URL"),
		os.Getenv("FACEBOOK_APP_ID"),
		os.Getenv("FACEBOOK_APP_SECRET"),
		os.Getenv("METRICS_PREFIX"),
		OdsMetricsQueueLength,
		OdsMetricsExportInterval,
	)
}
