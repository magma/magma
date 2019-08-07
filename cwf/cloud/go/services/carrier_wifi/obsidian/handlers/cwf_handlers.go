/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

// Package handlers provides REST API handlers for FeG related configuration
package handlers

import (
	"magma/cwf/cloud/go/cwf"
	"magma/cwf/cloud/go/services/carrier_wifi/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	configuratorhandlers "magma/orc8r/cloud/go/services/configurator/obsidian/handlers"
	networkpath "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
)

const (
	NetworkConfigPath = networkpath.ConfigureNetwork + "/" + cwf.CwfNetworkPath
)

// GetObsidianHandlers returns all obsidian handlers for cwf
func GetObsidianHandlers() []obsidian.Handler {
	return configuratorhandlers.GetCRUDNetworkConfigHandlers(NetworkConfigPath, cwf.CwfNetworkType, &models.NetworkCarrierWifiConfigs{})
}
