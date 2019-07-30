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
		{Path: ListNetworks, Methods: obsidian.GET, HandlerFunc: listNetworksHandler, MigratedHandlerFunc: listNetworks},
		{Path: RegisterNetwork, Methods: obsidian.POST, HandlerFunc: registerNetworkHandler, MigratedHandlerFunc: registerNetwork, MultiplexAfterMigration: true},
		{Path: ManageNetwork, Methods: obsidian.GET, HandlerFunc: getNetworkHandler, MigratedHandlerFunc: getNetwork},
		{Path: ManageNetwork, Methods: obsidian.PUT, HandlerFunc: updateNetworkHandler, MigratedHandlerFunc: updateNetwork, MultiplexAfterMigration: true},
		{Path: ManageNetwork, Methods: obsidian.DELETE, HandlerFunc: deleteNetworkHandler, MigratedHandlerFunc: deleteNetwork, MultiplexAfterMigration: true},

		// Gateway
		{Path: RegisterAG, Methods: obsidian.GET,
			HandlerFunc:         getListGatewaysHandler(&view_factory.FullGatewayViewFactoryLegacyImpl{}),
			MigratedHandlerFunc: getListGateways(&view_factory.FullGatewayViewFactoryImpl{})},
		{Path: RegisterAG, Methods: obsidian.POST, HandlerFunc: createGatewayHandler, MigratedHandlerFunc: createGateway, MultiplexAfterMigration: true},
		{Path: ManageAG, Methods: obsidian.GET, HandlerFunc: getGatewayHandler, MigratedHandlerFunc: getGateway},
		{Path: ManageAG, Methods: obsidian.PUT, HandlerFunc: updateGatewayHandler, MigratedHandlerFunc: updateGateway, MultiplexAfterMigration: true},
		{Path: ManageAG, Methods: obsidian.DELETE, HandlerFunc: deleteGatewayHandler, MigratedHandlerFunc: deleteGateway, MultiplexAfterMigration: true},

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
