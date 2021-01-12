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
	"magma/lte/cloud/go/lte"
	lte_service "magma/lte/cloud/go/services/lte"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	"magma/orc8r/cloud/go/services/metricsd"
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
	return []serde.Serde{}
}

func (*LteOrchestratorPlugin) GetMconfigBuilders() []mconfig.Builder {
	return []mconfig.Builder{
		mconfig.NewRemoteBuilder(lte_service.ServiceName),
	}
}

func (*LteOrchestratorPlugin) GetMetricsProfiles(metricsConfig *config.ConfigMap) []metricsd.MetricsProfile {
	return []metricsd.MetricsProfile{}
}

func (*LteOrchestratorPlugin) GetObsidianHandlers(metricsConfig *config.ConfigMap) []obsidian.Handler {
	return []obsidian.Handler{}
}

func (*LteOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{
		providers.NewRemoteProvider(lte_service.ServiceName, lte.SubscriberStreamName),
		providers.NewRemoteProvider(lte_service.ServiceName, lte.PolicyStreamName),
		providers.NewRemoteProvider(lte_service.ServiceName, lte.ApnRuleMappingsStreamName),
		providers.NewRemoteProvider(lte_service.ServiceName, lte.BaseNameStreamName),
		providers.NewRemoteProvider(lte_service.ServiceName, lte.NetworkWideRulesStreamName),
		providers.NewRemoteProvider(lte_service.ServiceName, lte.RatingGroupStreamName),
	}
}

func (*LteOrchestratorPlugin) GetStateIndexers() []indexer.Indexer {
	return []indexer.Indexer{}
}
