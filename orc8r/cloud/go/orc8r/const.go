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

	NetworkFeaturesConfig   = "orc8r_features"
	MagmadGatewayType       = "magmad_gateway"
	AccessGatewayRecordType = "access_gateway_record"
	GatewayStateType        = "gw_state"
	DirectoryRecordType     = "directory_record"
	StringMapSerdeType      = "string_map"

	UpgradeTierEntityType           = "upgrade_tier"
	UpgradeReleaseChannelEntityType = "upgrade_release_channel"

	DnsdNetworkType = "dnsd_network"

	MconfigBuilderLabel   = "orc8r.io/mconfig_builder"
	MetricsExporterLabel  = "orc8r.io/metrics_exporter"
	ObsidianHandlersLabel = "orc8r.io/obsidian_handlers"
	StateIndexerLabel     = "orc8r.io/state_indexer"
	StreamProviderLabel   = "orc8r.io/stream_provider"

	ObsidianHandlersPathPrefixesAnnotation = "orc8r.io/obsidian_handlers_path_prefixes"
	StateIndexerVersionAnnotation          = "orc8r.io/state_indexer_version"
	StateIndexerTypesAnnotation            = "orc8r.io/state_indexer_types"
	StreamProviderStreamsAnnotation        = "orc8r.io/stream_provider_streams"

	AnnotationFieldSeparator = ","
)
