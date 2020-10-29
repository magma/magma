/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pluginimpl

import (
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/collection"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/cloud/go/services/orchestrator"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"
)

// BaseOrchestratorPlugin is the OrchestratorPlugin for the orc8r module
type BaseOrchestratorPlugin struct{}

func (*BaseOrchestratorPlugin) GetName() string {
	return orc8r.ModuleName
}

func (*BaseOrchestratorPlugin) GetServices() []registry.ServiceLocation {
	serviceLocations, err := registry.LoadServiceRegistryConfig(orc8r.ModuleName)
	if err != nil {
		return []registry.ServiceLocation{}
	}
	return serviceLocations
}

func (*BaseOrchestratorPlugin) GetSerdes() []serde.Serde {
	return []serde.Serde{}
}

func (*BaseOrchestratorPlugin) GetMconfigBuilders() []mconfig.Builder {
	return []mconfig.Builder{
		mconfig.NewRemoteBuilder(orchestrator.ServiceName),
	}
}

func (*BaseOrchestratorPlugin) GetMetricsProfiles(metricsConfig *config.ConfigMap) []metricsd.MetricsProfile {
	return getMetricsProfiles()
}

func (*BaseOrchestratorPlugin) GetObsidianHandlers(metricsConfig *config.ConfigMap) []obsidian.Handler {
	return []obsidian.Handler{}
}

func (*BaseOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{
		providers.NewRemoteProvider(definitions.StreamerServiceName, definitions.MconfigStreamName),
	}
}

func (*BaseOrchestratorPlugin) GetStateIndexers() []indexer.Indexer {
	return []indexer.Indexer{
		indexer.NewRemoteIndexer(directoryd.ServiceName, 1, orc8r.DirectoryRecordType),
	}
}

const (
	ProfileNamePrometheus = "prometheus"
	ProfileNameExportAll  = "exportall"
)

func getMetricsProfiles() []metricsd.MetricsProfile {
	// Controller profile - 1 collector for each service
	services := registry.ListControllerServices()

	deviceCollectors := []collection.MetricCollector{&collection.DiskUsageMetricCollector{}, &collection.ProcMetricsCollector{}}
	allCollectors := make([]collection.MetricCollector, 0, len(services)+len(deviceCollectors))

	for _, s := range services {
		allCollectors = append(allCollectors, collection.NewCloudServiceMetricCollector(s))
	}
	for _, c := range deviceCollectors {
		allCollectors = append(allCollectors, c)
	}

	// Prometheus profile - Exports all service metric to Prometheus
	prometheusProfile := metricsd.MetricsProfile{
		Name:       ProfileNamePrometheus,
		Collectors: allCollectors,
		Exporters:  []exporters.Exporter{exporters.NewRemoteExporter(metricsd.ServiceName)},
	}

	// ExportAllProfile - Exports to all exporters
	exportAllProfile := metricsd.MetricsProfile{
		Name:       ProfileNameExportAll,
		Collectors: allCollectors,
		Exporters:  []exporters.Exporter{exporters.NewRemoteExporter(metricsd.ServiceName)},
	}

	return []metricsd.MetricsProfile{
		prometheusProfile,
		exportAllProfile,
	}
}
