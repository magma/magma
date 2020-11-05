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

package plugin

import (
	"magma/fbinternal/cloud/go/fbinternal"
	fbinternal_service "magma/fbinternal/cloud/go/services/fbinternal"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/collection"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"
)

type FbinternalOrchestratorPlugin struct{}

func (*FbinternalOrchestratorPlugin) GetName() string {
	return fbinternal.ModuleName
}

func (*FbinternalOrchestratorPlugin) GetServices() []registry.ServiceLocation {
	serviceLocations, err := registry.LoadServiceRegistryConfig(fbinternal.ModuleName)
	if err != nil {
		return []registry.ServiceLocation{}
	}
	return serviceLocations
}

func (*FbinternalOrchestratorPlugin) GetSerdes() []serde.Serde {
	return []serde.Serde{}
}

func (*FbinternalOrchestratorPlugin) GetMconfigBuilders() []mconfig.Builder {
	return []mconfig.Builder{}
}

func (*FbinternalOrchestratorPlugin) GetMetricsProfiles(metricsConfig *config.ConfigMap) []metricsd.MetricsProfile {
	return getMetricsProfiles()
}

func (*FbinternalOrchestratorPlugin) GetObsidianHandlers(metricsConfig *config.ConfigMap) []obsidian.Handler {
	return []obsidian.Handler{}
}

func (*FbinternalOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{}
}

func (*FbinternalOrchestratorPlugin) GetStateIndexers() []indexer.Indexer {
	return []indexer.Indexer{}
}

const (
	ProfileNameFacebook   = "facebook"
	ProfileNameSys        = "sys"
	ProfileNamePrometheus = "fbprometheus"
	ProfileNameExportAll  = "fbexportall"
)

func getMetricsProfiles() []metricsd.MetricsProfile {
	// Sys profile - collectors for disk usage and metricsd
	sysProfile := metricsd.MetricsProfile{
		Name: ProfileNameSys,
		Collectors: []collection.MetricCollector{
			&collection.DiskUsageMetricCollector{},
			collection.NewCloudServiceMetricCollector(metricsd.ServiceName),
		},
		Exporters: []exporters.Exporter{exporters.NewRemoteExporter(fbinternal_service.ServiceName)},
	}

	// Facebook profile - 1 collector for each service
	allServices := registry.ListControllerServices()
	controllerCollectors := make([]collection.MetricCollector, 0, len(allServices)+1)
	for _, srv := range allServices {
		controllerCollectors = append(controllerCollectors, collection.NewCloudServiceMetricCollector(srv))
	}

	controllerCollectors = append(controllerCollectors, &collection.DiskUsageMetricCollector{})
	facebookProfile := metricsd.MetricsProfile{
		Name:       ProfileNameFacebook,
		Collectors: controllerCollectors,
		Exporters: []exporters.Exporter{
			exporters.NewRemoteExporter(fbinternal_service.ServiceName),
		},
	}

	prometheusProfile := metricsd.MetricsProfile{
		Name:       ProfileNamePrometheus,
		Collectors: controllerCollectors,
		Exporters: []exporters.Exporter{
			exporters.NewRemoteExporter(fbinternal_service.ServiceName),
			exporters.NewRemoteExporter(metricsd.ServiceName),
		},
	}

	allExporterProfile := metricsd.MetricsProfile{
		Name:       ProfileNameExportAll,
		Collectors: controllerCollectors,
		Exporters: []exporters.Exporter{
			exporters.NewRemoteExporter(fbinternal_service.ServiceName),
			exporters.NewRemoteExporter(metricsd.ServiceName),
		},
	}

	return []metricsd.MetricsProfile{
		sysProfile,
		facebookProfile,
		prometheusProfile,
		allExporterProfile,
	}
}
