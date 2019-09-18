/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"net/http"

	"magma/feg/cloud/go/feg"
	fegmodels "magma/feg/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	orc8rmodels "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	FederationNetworks             = "feg"
	ListFegNetworksPath            = obsidian.V1Root + FederationNetworks
	ManageFegNetworkPath           = ListFegNetworksPath + "/:network_id"
	ManageFegNetworkFederationPath = ManageFegNetworkPath + obsidian.UrlSep + "federation"

	Gateways                    = "gateways"
	ListGatewaysPath            = ManageFegNetworkPath + obsidian.UrlSep + Gateways
	ManageGatewayPath           = ListGatewaysPath + obsidian.UrlSep + ":gateway_id"
	ManageGatewayStatePath      = ManageGatewayPath + obsidian.UrlSep + "status"
	ManageGatewayFederationPath = ManageGatewayPath + obsidian.UrlSep + "federation"

	FederatedLteNetworks              = "feg_lte"
	ListFegLteNetworksPath            = obsidian.V1Root + FederatedLteNetworks
	ManageFegLteNetworkPath           = ListFegLteNetworksPath + "/:network_id"
	ManageFegLteNetworkFederationPath = ManageFegLteNetworkPath + obsidian.UrlSep + "federation"
)

func GetHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ManageGatewayPath, Methods: obsidian.GET, HandlerFunc: getGateway},
		{Path: ManageGatewayStatePath, Methods: obsidian.GET, HandlerFunc: handlers.GetStateHandler},
	}
	ret = append(ret, handlers.GetTypedNetworkCRUDHandlers(ListFegNetworksPath, ManageFegNetworkPath, feg.FederatedLteNetworkType, &fegmodels.FegLteNetwork{})...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageFegNetworkFederationPath, &fegmodels.NetworkFederationConfigs{}, "")...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayFederationPath, &fegmodels.GatewayFederationConfigs{})...)

	ret = append(ret, handlers.GetTypedNetworkCRUDHandlers(ListFegLteNetworksPath, ManageFegLteNetworkPath, feg.FederatedLteNetworkType, &fegmodels.FegLteNetwork{})...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageFegLteNetworkFederationPath, &fegmodels.FederatedNetworkConfigs{}, "")...)

	ret = append(ret, handlers.GetListGatewaysHandler(ListGatewaysPath, feg.FegGatewayType, makeFederationGateways))
	ret = append(ret, handlers.GetCreateGatewayHandler(ListGatewaysPath, feg.FegGatewayType, &fegmodels.MutableFederationGateway{}))
	ret = append(ret, handlers.GetUpdateGatewayHandler(ManageGatewayPath, feg.FegGatewayType, &fegmodels.MutableFederationGateway{}))
	ret = append(ret, handlers.GetDeleteGatewayHandler(ManageGatewayPath, feg.FegGatewayType))

	return ret
}

func getGateway(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	magmadModel, nerr := handlers.LoadMagmadGatewayModel(nid, gid)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(
		nid, feg.FegGatewayType, gid,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
	)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load federation gateway"), http.StatusInternalServerError)
	}

	ret := &fegmodels.FederationGateway{
		ID:          magmadModel.ID,
		Name:        magmadModel.Name,
		Description: magmadModel.Description,
		Device:      magmadModel.Device,
		Status:      magmadModel.Status,
		Tier:        magmadModel.Tier,
		Magmad:      magmadModel.Magmad,
		Federation:  ent.Config.(*fegmodels.GatewayFederationConfigs),
	}
	return c.JSON(http.StatusOK, ret)
}

type federationAndMagmadGateway struct {
	magmadGateway, federationGateway configurator.NetworkEntity
}

func makeFederationGateways(
	entsByTK map[storage.TypeAndKey]configurator.NetworkEntity,
	devicesByID map[string]interface{},
	statusesByID map[string]*orc8rmodels.GatewayStatus,
) map[string]handlers.GatewayModel{
	gatewayEntsByKey := map[string]*federationAndMagmadGateway{}
	for tk, ent := range entsByTK {
		existing, found := gatewayEntsByKey[tk.Key]
		if !found {
			existing = &federationAndMagmadGateway{}
			gatewayEntsByKey[tk.Key] = existing
		}

		switch ent.Type {
		case orc8r.MagmadGatewayType:
			existing.magmadGateway = ent
		case feg.FegGatewayType:
			existing.federationGateway = ent
		}
	}

	ret := make(map[string]handlers.GatewayModel, len(gatewayEntsByKey))
	for key, ents := range gatewayEntsByKey {
		hwID := ents.magmadGateway.PhysicalID
		var devCasted *orc8rmodels.GatewayDevice
		if devicesByID[hwID] != nil {
			devCasted = devicesByID[hwID].(*orc8rmodels.GatewayDevice)
		}
		ret[key] = (&fegmodels.FederationGateway{}).FromBackendModels(ents.magmadGateway, ents.federationGateway, devCasted, statusesByID[hwID])
	}
	return ret
}
