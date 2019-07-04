/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package orc8r

const (
	ModuleName string = "orc8r"

	MagmadNetworkType       = "magmad_network"
	NetworkFeaturesConfig   = "orc8r_features"
	MagmadGatewayType       = "magmad_gateway"
	AccessGatewayRecordType = "access_gateway_record"

	UpgradeTierEntityType = "upgrade_tier"

	DnsdNetworkType = "dnsd_network"

	// used to migrate network/gateway lookup magmad dependencies
	UseConfiguratorEnv = "USE_NEW_HANDLERS"

	// separate flag to control mconfig builders because this has significant
	// implications on production gateways
	UseConfiguratorMconfigsEnv = "USE_NEW_MCONFIGS"

	// comma-separated list of networks to run new mconfig builders for
	MconfigWhitelistEnv = "NEW_MCONFIGS_WHITELIST"
)
