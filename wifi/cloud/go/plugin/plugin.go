/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package plugin

import (
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"
	wifi_service "magma/wifi/cloud/go/services/wifi"
	"magma/wifi/cloud/go/wifi"
)

type WifiOrchestratorPlugin struct{}

func (*WifiOrchestratorPlugin) GetName() string {
	return wifi.ModuleName
}

func (*WifiOrchestratorPlugin) GetServices() []registry.ServiceLocation {
	serviceLocations, err := registry.LoadServiceRegistryConfig(wifi.ModuleName)
	if err != nil {
		return []registry.ServiceLocation{}
	}
	return serviceLocations
}

func (*WifiOrchestratorPlugin) GetSerdes() []serde.Serde {
	return []serde.Serde{}
}

func (*WifiOrchestratorPlugin) GetMconfigBuilders() []mconfig.Builder {
	return []mconfig.Builder{
		mconfig.NewRemoteBuilder(wifi_service.ServiceName),
	}
}

func (*WifiOrchestratorPlugin) GetMetricsProfiles(metricsConfig *config.ConfigMap) []metricsd.MetricsProfile {
	return []metricsd.MetricsProfile{}
}

func (*WifiOrchestratorPlugin) GetObsidianHandlers(metricsConfig *config.ConfigMap) []obsidian.Handler {
	return []obsidian.Handler{}
}

func (*WifiOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{}
}

func (*WifiOrchestratorPlugin) GetStateIndexers() []indexer.Indexer {
	return []indexer.Indexer{}
}
