/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/config/obsidian"
	"magma/orc8r/cloud/go/services/magmad/config"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	magmad_models "magma/orc8r/cloud/go/services/magmad/obsidian/models"
)

// GetObsidianHandlers returns all obsidian handlers for magmad
func GetObsidianHandlers() []handlers.Handler {
	return []handlers.Handler{
		// Network
		{Path: ListNetworks, Methods: handlers.GET, HandlerFunc: listNetworks},
		{Path: RegisterNetwork, Methods: handlers.POST, HandlerFunc: registerNetwork},
		{Path: ManageNetwork, Methods: handlers.GET, HandlerFunc: getNetwork},
		{Path: ManageNetwork, Methods: handlers.PUT, HandlerFunc: updateNetwork},
		{Path: ManageNetwork, Methods: handlers.DELETE, HandlerFunc: deleteNetwork},

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

		obsidian.GetReadGatewayConfigHandler(ConfigureAG, config.MagmadGatewayType, &magmad_models.MagmadGatewayConfig{}),
		obsidian.GetCreateGatewayConfigHandler(ConfigureAG, config.MagmadGatewayType, &magmad_models.MagmadGatewayConfig{}),
		obsidian.GetUpdateGatewayConfigHandler(ConfigureAG, config.MagmadGatewayType, &magmad_models.MagmadGatewayConfig{}),
		obsidian.GetDeleteGatewayConfigHandler(ConfigureAG, config.MagmadGatewayType),
	}
}
