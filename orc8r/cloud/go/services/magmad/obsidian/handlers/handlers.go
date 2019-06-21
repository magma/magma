/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/config/obsidian"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	magmad_models "magma/orc8r/cloud/go/services/magmad/obsidian/models"
)

// GetObsidianHandlers returns all obsidian handlers for magmad
func GetObsidianHandlers() []handlers.Handler {
	return []handlers.Handler{
		// Network
		{Path: ListNetworks, Methods: handlers.GET, HandlerFunc: listNetworksHandler, MigratedHandlerFunc: listNetworks},
		{Path: RegisterNetwork, Methods: handlers.POST, HandlerFunc: registerNetworkHandler, MigratedHandlerFunc: registerNetwork},
		{Path: ManageNetwork, Methods: handlers.GET, HandlerFunc: getNetworkHandler, MigratedHandlerFunc: getNetwork},
		{Path: ManageNetwork, Methods: handlers.PUT, HandlerFunc: updateNetworkHandler, MigratedHandlerFunc: updateNetwork},
		{Path: ManageNetwork, Methods: handlers.DELETE, HandlerFunc: deleteNetworkHandler, MigratedHandlerFunc: deleteNetwork},

		// Gateway
		{Path: RegisterAG, Methods: handlers.GET, HandlerFunc: getListGatewaysHandler(&view_factory.FullGatewayViewFactoryImpl{})},
		{Path: RegisterAG, Methods: handlers.POST, HandlerFunc: registerGateway},
		{Path: ManageAG, Methods: handlers.GET, HandlerFunc: getGateway},
		{Path: ManageAG, Methods: handlers.PUT, HandlerFunc: updateGateway},
		{Path: ManageAG, Methods: handlers.DELETE, HandlerFunc: deleteGateway},

		// Gateway Commands
		{Path: RebootGateway, Methods: handlers.POST, HandlerFunc: rebootGateway},
		{Path: RestartServices, Methods: handlers.POST, HandlerFunc: restartServices},
		{Path: GatewayPing, Methods: handlers.POST, HandlerFunc: gatewayPing},
		{Path: GatewayGenericCommand, Methods: handlers.POST, HandlerFunc: gatewayGenericCommand},
		{Path: TailGatewayLogs, Methods: handlers.POST, HandlerFunc: tailGatewayLogs},

		obsidian.GetReadGatewayConfigHandler(ConfigureAG, orc8r.MagmadGatewayType, &magmad_models.MagmadGatewayConfig{}),
		obsidian.GetCreateGatewayConfigHandler(ConfigureAG, orc8r.MagmadGatewayType, &magmad_models.MagmadGatewayConfig{}),
		obsidian.GetUpdateGatewayConfigHandler(ConfigureAG, orc8r.MagmadGatewayType, &magmad_models.MagmadGatewayConfig{}),
		obsidian.GetDeleteGatewayConfigHandler(ConfigureAG, orc8r.MagmadGatewayType),
	}
}
