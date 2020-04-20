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

package lte

const ModuleName = "lte"

const (
	NetworkType = "lte"

	CellularNetworkType         = "cellular_network"
	CellularGatewayType         = "cellular_gateway"
	CellularEnodebType          = "cellular_enodeb"
	NetworkSubscriberConfigType = "network_subscriber_config"

	EnodebStateType      = "single_enodeb"
	SubscriberEntityType = "subscriber"
	ICMPStateType        = "icmp_monitoring"

	BaseNameEntityType   = "base_name"
	PolicyRuleEntityType = "policy"

	RatingGroupEntityType = "rating_group"

	ApnEntityType = "apn"

	SubscriberStreamName       = "subscriberdb"
	PolicyStreamName           = "policydb"
	BaseNameStreamName         = "base_names"
	MappingsStreamName         = "rule_mappings"
	NetworkWideRulesStreamName = "network_wide_rules"
	RatingGroupStreamName      = "rating_groups"

	// Replicated states from AGW
	SPGWStateType      = "SPGW"
	MMEStateType       = "MME"
	S1APStateType      = "S1AP"
	MobilitydStateType = "mobilityd_ipdesc_record"
)
