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

package orc8r

const (
	ModuleName = "orc8r"
)

// State and entities
const (
	NetworkFeaturesConfig   = "orc8r_features"
	MagmadGatewayType       = "magmad_gateway"
	AccessGatewayRecordType = "access_gateway_record"
	GatewayStateType        = "gw_state"
	DirectoryRecordType     = "directory_record"
	StringMapSerdeType      = "string_map"

	DnsdNetworkType = "dnsd_network"

	UpgradeTierEntityType           = "upgrade_tier"
	UpgradeReleaseChannelEntityType = "upgrade_release_channel"

	CallTraceEntityType = "call_trace"
)

// K8s
const (
	// PartOfLabel and PartOfOrc8rApp are K8s label key and values indicating
	// a service is an orc8r application service.
	PartOfLabel    = "app.kubernetes.io/part-of"
	PartOfOrc8rApp = "orc8r-app"

	GRPCPortName = "grpc"
	HTTPPortName = "http"

	AnnotationFieldSeparator = ","

	AnalyticsCollectorLabel = "orc8r.io/analytics_collector"
	MconfigBuilderLabel     = "orc8r.io/mconfig_builder"
	MetricsExporterLabel    = "orc8r.io/metrics_exporter"
	ObsidianHandlersLabel   = "orc8r.io/obsidian_handlers"
	StateIndexerLabel       = "orc8r.io/state_indexer"
	StreamProviderLabel     = "orc8r.io/stream_provider"
	SwaggerSpecLabel        = "orc8r.io/swagger_spec"

	ObsidianHandlersPathPrefixesAnnotation = "orc8r.io/obsidian_handlers_path_prefixes"
	StateIndexerVersionAnnotation          = "orc8r.io/state_indexer_version"
	StateIndexerTypesAnnotation            = "orc8r.io/state_indexer_types"
	StreamProviderStreamsAnnotation        = "orc8r.io/stream_provider_streams"
)

// Environment variables
const (
	// ServiceHostnameEnvVar is the name of an environment variable which is
	// required to hold the public IP of the service.
	// In dev, this will generally be localhost.
	// In prod, this will be the relevant pod's IP.
	ServiceHostnameEnvVar = "SERVICE_HOSTNAME"
)

// Configs
const (
	// SharedService is the name of the pseudo-service that stores shared
	// configs across all Orc8r services.
	SharedService = "shared"
)
