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

	Configurator type     Swagger model          Edges out                                  Notes
	-----------------     -------------          ---------                                  -----
	apn                   apn                                                               Stored as apn_configuration
	apn_policy_profile    apn_policy_profile     apn,policy                                 Internal-only
	apn_resource          apn_resource           apn
	base_name             base_name_record       policy
	cellular_enodeb       enodeb                                                            Stored as enodeb_configuration
	cellular_gateway      *lte_gateway           cellular_enodeb,apn_resource               Stored as gateway_cellular_configs
	cellular_gateway_pool *cellular_gateway_pool cellular_gateway                           Stored as cellular_gateway_pool_configs
	policy                policy_rule            policy_qos_profile                         Stored as policy_rule_config
	policy_qos_profile    policy_qos_profile
	rating_group          *rating_group
	subscriber            *subscriber           apn,policy,base_name,apn_policy_profile

	Resulting DAG

	cellular_gateway_pool -.-> cellular_gateway -.-> cellular_enodeb
	                                             '-> apn_resource -> apn*
	subscriber -.-> apn_policy_profile -.-> apn*
	            '-> apn*                '-> policy*
	            '-> policy*
	            '-> base_name -> policy*
	*policy -> policy_qos_profile

	Notes

	Where possible, keep cellular_gateway and subscriber as sources rather
	than sinks. This reduces the number of data migrations required to
	reorganize the graph into a set of acyclic relations.
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
	APNEntityType                     = "apn"
	APNPolicyProfileEntityType        = "apn_policy_profile"
	APNResourceEntityType             = "apn_resource"
	BaseNameEntityType                = "base_name"
	CellularEnodebEntityType          = "cellular_enodeb"
	CellularGatewayEntityType         = "cellular_gateway"
	CellularGatewayPoolEntityType     = "cellular_gateway_pool"
	PolicyQoSProfileEntityType        = "policy_qos_profile"
	PolicyRuleEntityType              = "policy"
	RatingGroupEntityType             = "rating_group"
	SubscriberEntityType              = "subscriber"
	NetworkProbeTaskEntityType        = "network_probe_task"
	NetworkProbeDestinationEntityType = "network_probe_destination"

	// ApnRuleMappingsStreamName etc. are streamer stream names.
	ApnRuleMappingsStreamName  = "apn_rule_mappings"
	BaseNameStreamName         = "base_names"
	NetworkWideRulesStreamName = "network_wide_rules"
	PolicyStreamName           = "policydb"
	RatingGroupStreamName      = "rating_groups"
	SubscriberStreamName       = "subscriberdb"

	// EnodebStateType etc. denote types of state replicated from AGWs.
	EnodebStateType     = "single_enodeb"
	ICMPStateType       = "icmp_monitoring"
	MMEStateType        = "MME"
	MobilitydStateType  = "mobilityd_ipdesc_record"
	S1APStateType       = "S1AP"
	SPGWStateType       = "SPGW"
	SubscriberStateType = "subscriber_state"

	// MSISDNBlobstoreType etc. denote blob types stored in blobstore tables.
	MSISDNBlobstoreType = "msisdn"
)
