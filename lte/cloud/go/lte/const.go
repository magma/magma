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

/*
	File const.go provides constants used across the LTE module.

	A number of the constants refer to configurator network entities. These
	are reproduced and annotated below.
	Note: starred Swagger models also have mutable versions of the model.

	Configurator type    Swagger model       Edges out           Notes
	-----------------    -------------       ---------           -----
	apn                  apn                                     Stored as apn_configuration
	base_name            base_name_record    policy, subscriber
	cellular_enodeb      enodeb                                  Stored as enodeb_configuration
	cellular_gateway     *lte_gateway        cellular_enodeb     Stored as gateway_cellular_configs
	policy               policy_rule         subscriber          Stored as policy_rule_config
	rating_group         *rating_group
	subscriber           *subscriber         apn

	Resulting DAG

	cellular_gateway -> cellular_enodeb
	base_name -> (policy ->) subscriber -> apn
*/

package lte

const ModuleName = "lte"

const (
	// NetworkType identifies managed networks as LTE networks.
	NetworkType = "lte"

	// CellularNetworkConfigType etc. are keys to network-level configs stored
	// in configurator.
	CellularNetworkConfigType   = "cellular_network"
	NetworkSubscriberConfigType = "network_subscriber_config"

	// APNEntityType etc. are configurator network entity types.
	APNEntityType              = "apn"
	APNResourceEntityType      = "apn_resource"
	BaseNameEntityType         = "base_name"
	CellularEnodebEntityType   = "cellular_enodeb"
	CellularGatewayEntityType  = "cellular_gateway"
	PolicyQoSProfileEntityType = "policy_qos_profile"
	PolicyRuleEntityType       = "policy"
	RatingGroupEntityType      = "rating_group"
	SubscriberEntityType       = "subscriber"

	// BaseNameStreamName etc. are streamer stream names.
	BaseNameStreamName         = "base_names"
	MappingsStreamName         = "rule_mappings"
	NetworkWideRulesStreamName = "network_wide_rules"
	PolicyStreamName           = "policydb"
	RatingGroupStreamName      = "rating_groups"
	SubscriberStreamName       = "subscriberdb"

	// EnodebStateType etc. are denote types of state replicated from AGWs
	EnodebStateType    = "single_enodeb"
	ICMPStateType      = "icmp_monitoring"
	MMEStateType       = "MME"
	MobilitydStateType = "mobilityd_ipdesc_record"
	S1APStateType      = "S1AP"
	SPGWStateType      = "SPGW"
)
