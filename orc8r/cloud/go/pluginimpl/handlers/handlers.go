/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	models2 "magma/orc8r/cloud/go/pluginimpl/models"
)

const (
	Networks                     = "networks"
	ListNetworksPath             = obsidian.V1Root + Networks
	RegisterNetworkPath          = obsidian.V1Root + Networks
	ManageNetworkPath            = obsidian.V1Root + Networks + obsidian.UrlSep + ":network_id"
	ManageNetworkNamePath        = ManageNetworkPath + obsidian.UrlSep + "name"
	ManageNetworkTypePath        = ManageNetworkPath + obsidian.UrlSep + "type"
	ManageNetworkDescriptionPath = ManageNetworkPath + obsidian.UrlSep + "description"
	ManageNetworkFeaturesPath    = ManageNetworkPath + obsidian.UrlSep + "features"
	ManageNetworkDNSPath         = ManageNetworkPath + obsidian.UrlSep + "dns"

	Gateways          = "gateways"
	ListGatewaysPath  = ManageNetworkPath + obsidian.UrlSep + Gateways
	ManageGatewayPath = ListGatewaysPath + obsidian.UrlSep + ":gateway_id"
)

// GetObsidianHandlers returns all plugin-level obsidian handlers for orc8r
func GetObsidianHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		// Magma V1 Network
		{Path: ListNetworksPath, Methods: obsidian.GET, HandlerFunc: ListNetworks},
		{Path: RegisterNetworkPath, Methods: obsidian.POST, HandlerFunc: RegisterNetwork},
		{Path: ManageNetworkPath, Methods: obsidian.GET, HandlerFunc: GetNetwork},
		{Path: ManageNetworkPath, Methods: obsidian.PUT, HandlerFunc: UpdateNetwork},
		{Path: ManageNetworkPath, Methods: obsidian.DELETE, HandlerFunc: DeleteNetwork},

		// Magma V1 Gateways
		{Path: ListGatewaysPath, Methods: obsidian.GET, HandlerFunc: ListGateways},
		{Path: ListGatewaysPath, Methods: obsidian.POST, HandlerFunc: CreateGateway},
		{Path: ManageGatewayPath, Methods: obsidian.GET, HandlerFunc: GetGateway},
		{Path: ManageGatewayPath, Methods: obsidian.PUT, HandlerFunc: UpdateGateway},
		{Path: ManageGatewayPath, Methods: obsidian.DELETE, HandlerFunc: DeleteGateway},
	}
	ret = append(ret, GetPartialNetworkHandlers(ManageNetworkNamePath, new(models.NetworkName), "")...)
	ret = append(ret, GetPartialNetworkHandlers(ManageNetworkTypePath, new(models.NetworkType), "")...)
	ret = append(ret, GetPartialNetworkHandlers(ManageNetworkDescriptionPath, new(models.NetworkDescription), "")...)
	ret = append(ret, GetPartialNetworkHandlers(ManageNetworkFeaturesPath, &models2.NetworkFeatures{}, orc8r.NetworkFeaturesConfig)...)
	ret = append(ret, GetPartialNetworkHandlers(ManageNetworkDNSPath, &models2.NetworkDNSConfig{}, orc8r.DnsdNetworkType)...)
	return ret
}
