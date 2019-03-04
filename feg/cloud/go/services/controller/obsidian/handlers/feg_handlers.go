/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package handlers provides REST API handlers for FeG related configuration
package handlers

import (
	"magma/feg/cloud/go/services/controller/config"
	"magma/feg/cloud/go/services/controller/obsidian/models"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/config/obsidian"
	magmad_handlers "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
)

const (
	ConfigKey         = "federation"
	NetworkConfigPath = magmad_handlers.ConfigureNetwork + "/" + ConfigKey
	GatewayConfigPath = magmad_handlers.ConfigureAG + "/" + ConfigKey
)

// GetObsidianHandlers returns all obsidian handlers for feg controller
func GetObsidianHandlers() []handlers.Handler {
	ret := make([]handlers.Handler, 0, 8)
	ret = append(ret, obsidian.GetCRUDNetworkConfigHandlers(NetworkConfigPath, config.FegNetworkType, &models.NetworkFederationConfigs{})...)
	ret = append(ret, obsidian.GetCRUDGatewayConfigHandlers(GatewayConfigPath, config.FegGatewayType, &models.GatewayFegConfigs{})...)
	return ret
}
