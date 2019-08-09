/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	cfgObsidian "magma/orc8r/cloud/go/services/config/obsidian"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	magmad_models "magma/orc8r/cloud/go/services/magmad/obsidian/models"
)

// GetObsidianHandlers returns all obsidian handlers for magmad
func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		// Network
		{Path: ListNetworks, Methods: obsidian.GET, HandlerFunc: listNetworks},
		{Path: RegisterNetwork, Methods: obsidian.POST, HandlerFunc: registerNetwork},
		{Path: ManageNetwork, Methods: obsidian.GET, HandlerFunc: getNetwork},
		{Path: ManageNetwork, Methods: obsidian.PUT, HandlerFunc: updateNetwork},
		{Path: ManageNetwork, Methods: obsidian.DELETE, HandlerFunc: deleteNetwork},

		// Gateway
		{Path: RegisterAG, Methods: obsidian.GET, HandlerFunc: getListGateways(&view_factory.FullGatewayViewFactoryImpl{})},
		{Path: RegisterAG, Methods: obsidian.POST, HandlerFunc: createGateway},
		{Path: ManageAG, Methods: obsidian.GET, HandlerFunc: getGateway},
		{Path: ManageAG, Methods: obsidian.PUT, HandlerFunc: updateGateway},
		{Path: ManageAG, Methods: obsidian.DELETE, HandlerFunc: deleteGateway},

		// Gateway Commands
		{Path: RebootGateway, Methods: obsidian.POST, HandlerFunc: rebootGateway},
		{Path: RestartServices, Methods: obsidian.POST, HandlerFunc: restartServices},
		{Path: GatewayPing, Methods: obsidian.POST, HandlerFunc: gatewayPing},
		{Path: GatewayGenericCommand, Methods: obsidian.POST, HandlerFunc: gatewayGenericCommand},
		{Path: TailGatewayLogs, Methods: obsidian.POST, HandlerFunc: tailGatewayLogs},

		cfgObsidian.GetReadGatewayConfigHandler(ConfigureAG, orc8r.MagmadGatewayType, &magmad_models.MagmadGatewayConfig{}),
		cfgObsidian.GetCreateGatewayConfigHandler(ConfigureAG, orc8r.MagmadGatewayType, &magmad_models.MagmadGatewayConfig{}),
		cfgObsidian.GetUpdateGatewayConfigHandler(ConfigureAG, orc8r.MagmadGatewayType, &magmad_models.MagmadGatewayConfig{}),
		cfgObsidian.GetDeleteGatewayConfigHandler(ConfigureAG, orc8r.MagmadGatewayType),
	}
}
