/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin

import (
	"fmt"
	"net/http"
	"sort"

	"magma/lte/cloud/go/lte"
	ltemodels "magma/lte/cloud/go/plugin/models"
	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	orc8rmodels "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	LteNetworks                        = "ltenetworks"
	ListNetworksPath                   = obsidian.V1Root + LteNetworks
	ManageNetworkPath                  = ListNetworksPath + "/:network_id"
	ManageNetworkNamePath              = ManageNetworkPath + obsidian.UrlSep + "name"
	ManageNetworkDescriptionPath       = ManageNetworkPath + obsidian.UrlSep + "description"
	ManageNetworkFeaturesPath          = ManageNetworkPath + obsidian.UrlSep + "features"
	ManageNetworkDNSPath               = ManageNetworkPath + obsidian.UrlSep + "dns"
	ManageNetworkDNSRecordsPath        = ManageNetworkDNSPath + obsidian.UrlSep + "records"
	ManageNetworkDNSRecordByDomainPath = ManageNetworkDNSRecordsPath + obsidian.UrlSep + ":domain"
	ManageNetworkCellularPath          = ManageNetworkPath + obsidian.UrlSep + "cellular"
	ManageNetworkCellularEpcPath       = ManageNetworkCellularPath + obsidian.UrlSep + "epc"
	ManageNetworkCellularRanPath       = ManageNetworkCellularPath + obsidian.UrlSep + "ran"
	ManageNetworkCellularFegNetworkID  = ManageNetworkCellularPath + obsidian.UrlSep + "feg_network_id"
)

func GetNetworkHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ListNetworksPath, Methods: obsidian.GET, HandlerFunc: listNetworks},
		{Path: ListNetworksPath, Methods: obsidian.POST, HandlerFunc: createNetwork},
		{Path: ManageNetworkPath, Methods: obsidian.GET, HandlerFunc: getNetwork},
		{Path: ManageNetworkPath, Methods: obsidian.PUT, HandlerFunc: updateNetwork},
		{Path: ManageNetworkPath, Methods: obsidian.DELETE, HandlerFunc: deleteNetwork},

		{Path: ManageNetworkDNSRecordByDomainPath, Methods: obsidian.POST, HandlerFunc: handlers.CreateDNSRecord},
		{Path: ManageNetworkDNSRecordByDomainPath, Methods: obsidian.GET, HandlerFunc: handlers.ReadDNSRecord},
		{Path: ManageNetworkDNSRecordByDomainPath, Methods: obsidian.PUT, HandlerFunc: handlers.UpdateDNSRecord},
		{Path: ManageNetworkDNSRecordByDomainPath, Methods: obsidian.DELETE, HandlerFunc: handlers.DeleteDNSRecord},
	}
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkNamePath, new(models.NetworkName), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDescriptionPath, new(models.NetworkDescription), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkFeaturesPath, &orc8rmodels.NetworkFeatures{}, orc8r.NetworkFeaturesConfig)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDNSPath, &orc8rmodels.NetworkDNSConfig{}, orc8r.DnsdNetworkType)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDNSRecordsPath, new(orc8rmodels.NetworkDNSRecords), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkCellularPath, &ltemodels.NetworkCellularConfigs{}, lte.CellularNetworkType)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkCellularEpcPath, &ltemodels.NetworkEpcConfigs{}, "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkCellularRanPath, &ltemodels.NetworkRanConfigs{}, "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkCellularFegNetworkID, new(ltemodels.FegNetworkID), "")...)
	return ret
}

func listNetworks(c echo.Context) error {
	ids, err := configurator.ListNetworksOfType(lte.LteNetworkType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	sort.Strings(ids)
	return c.JSON(http.StatusOK, ids)
}

func createNetwork(c echo.Context) error {
	payload := &ltemodels.LteNetwork{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	err := configurator.CreateNetwork(payload.ToConfiguratorNetwork())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func getNetwork(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	network, err := configurator.LoadNetwork(nid, true, true)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if network.Type != lte.LteNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not an LTE network", nid))
	}

	ret := (&ltemodels.LteNetwork{}).FromConfiguratorNetwork(network)
	return c.JSON(http.StatusOK, ret)
}

func updateNetwork(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := &ltemodels.LteNetwork{}
	err := c.Bind(payload)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	// check that this is actually an LTE network
	network, err := configurator.LoadNetwork(nid, false, false)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load network to check type"), http.StatusInternalServerError)
	}
	if network.Type != lte.LteNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not an LTE network", nid))
	}

	err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{payload.ToUpdateCriteria()})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteNetwork(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	// check that this is actually an LTE network
	network, err := configurator.LoadNetwork(nid, false, false)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load network to check type"), http.StatusInternalServerError)
	}
	if network.Type != lte.LteNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not an LTE network", nid))
	}

	err = configurator.DeleteNetwork(nid)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func ListGateways(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	ids, err := configurator.ListEntityKeys(nid, lte.CellularGatewayType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// for each ID, we want to load the cellular gateway and the magmad gateway
	entityTKs := make([]storage.TypeAndKey, 0, len(ids)*2)
	for _, id := range ids {
		entityTKs = append(
			entityTKs,
			storage.TypeAndKey{Type: orc8r.MagmadGatewayType, Key: id},
			storage.TypeAndKey{Type: lte.CellularGatewayType, Key: id},
		)
	}
	ents, _, err := configurator.LoadEntities(nid, nil, nil, nil, entityTKs, configurator.FullEntityLoadCriteria())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	entsByTK := ents.ToEntitiesByID()

	// for each magmad gateway, we have to load its corresponding device and
	// its reported status
	deviceIDs := make([]string, 0, len(ids))
	for tk, ent := range entsByTK {
		if tk.Type == orc8r.MagmadGatewayType && ent.PhysicalID != "" {
			deviceIDs = append(deviceIDs, ent.PhysicalID)
		}
	}
	devicesByID, err := device.GetDevices(nid, orc8r.AccessGatewayRecordType, deviceIDs)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load devices"), http.StatusInternalServerError)
	}
	statusesByID, err := state.GetGatewayStatuses(nid, deviceIDs)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load statuses"), http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, makeLTEGateways(entsByTK, devicesByID, statusesByID))
}

func CreateGateway(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := &ltemodels.LteGateway{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	if nerr := handlers.CreateMagmadGatewayFromModel(nid, payload.GetMagmadGateway()); nerr != nil {
		return nerr
	}

	if _, err := configurator.CreateEntity(nid, payload.ToConfiguratorEntity()); err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to create cellular gateway"), http.StatusInternalServerError)
	}
	if _, err := configurator.UpdateEntity(nid, payload.GetMagmadGatewayUpdateCriteria()); err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to associate cellular and magmad gateways"), http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func GetGateway(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	magmadModel, nerr := handlers.LoadMagmadGatewayModel(nid, gid)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(
		nid, lte.CellularGatewayType, gid,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
	)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load cellular gateway"), http.StatusInternalServerError)
	}

	ret := &ltemodels.LteGateway{
		ID:          magmadModel.ID,
		Name:        magmadModel.Name,
		Description: magmadModel.Description,
		Device:      magmadModel.Device,
		Status:      magmadModel.Status,
		Tier:        magmadModel.Tier,
		Magmad:      magmadModel.Magmad,
		Cellular:    ent.Config.(*ltemodels.GatewayCellularConfigs),
	}
	for _, tk := range ent.Associations {
		if tk.Type == lte.CellularEnodebType {
			ret.ConnectedEnodebSerials = append(ret.ConnectedEnodebSerials, tk.Key)
		}
	}
	return c.JSON(http.StatusOK, ret)
}

func UpdateGateway(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	payload := &ltemodels.LteGateway{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	_, err := configurator.LoadEntity(
		nid, lte.CellularGatewayType, gid,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
	)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(errors.Wrap(err, "failed to load cellular gateway"), http.StatusInternalServerError)
	}

	if nerr := handlers.UpdateMagmadGatewayFromModel(nid, gid, payload.GetMagmadGateway()); nerr != nil {
		return nerr
	}
	if _, err := configurator.UpdateEntity(nid, payload.ToEntityUpdateCriteria()); err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to update cellular gateway"), http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func DeleteGateway(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntities(
		nid,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: gid},
			{Type: lte.CellularGatewayType, Key: gid},
		},
	)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)

}

type cellularAndMagmadGateway struct {
	magmadGateway, cellularGateway configurator.NetworkEntity
}

func makeLTEGateways(
	entsByTK map[storage.TypeAndKey]configurator.NetworkEntity,
	devicesByID map[string]interface{},
	statusesByID map[string]*orc8rmodels.GatewayStatus,
) map[string]*ltemodels.LteGateway {
	gatewayEntsByKey := map[string]*cellularAndMagmadGateway{}
	for tk, ent := range entsByTK {
		existing, found := gatewayEntsByKey[tk.Key]
		if !found {
			existing = &cellularAndMagmadGateway{}
			gatewayEntsByKey[tk.Key] = existing
		}

		switch ent.Type {
		case orc8r.MagmadGatewayType:
			existing.magmadGateway = ent
		case lte.CellularGatewayType:
			existing.cellularGateway = ent
		}
	}

	ret := make(map[string]*ltemodels.LteGateway, len(gatewayEntsByKey))
	for key, ents := range gatewayEntsByKey {
		hwID := ents.magmadGateway.PhysicalID
		var devCasted *orc8rmodels.GatewayDevice
		if devicesByID[hwID] != nil {
			devCasted = devicesByID[hwID].(*orc8rmodels.GatewayDevice)
		}
		ret[key] = (&ltemodels.LteGateway{}).FromBackendModels(ents.magmadGateway, ents.cellularGateway, devCasted, statusesByID[hwID])
	}
	return ret
}
