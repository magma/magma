/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package lte

const ModuleName = "lte"

const (
	LteNetworkType = "lte"

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

	SubscriberStreamName = "subscriberdb"
	PolicyStreamName     = "policydb"
	BaseNameStreamName   = "base_names"
	MappingsStreamName   = "rule_mappings"
	NetworkWideRules     = "network_wide_rules"

	// Replicated states from AGW
	SPGWStateType = "SPGW"
	MMEStateType  = "MME"
	S1APStateType = "S1AP"
)
