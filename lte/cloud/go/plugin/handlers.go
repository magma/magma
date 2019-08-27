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

	Gateways                         = "gateways"
	ListGatewaysPath                 = ManageNetworkPath + obsidian.UrlSep + Gateways
	ManageGatewayPath                = ListGatewaysPath + obsidian.UrlSep + ":gateway_id"
	ManageGatewayNamePath            = ManageGatewayPath + obsidian.UrlSep + "name"
	ManageGatewayDescriptionPath     = ManageGatewayPath + obsidian.UrlSep + "description"
	ManageGatewayConfigPath          = ManageGatewayPath + obsidian.UrlSep + "magmad"
	ManageGatewayDevicePath          = ManageGatewayPath + obsidian.UrlSep + "device"
	ManageGatewayStatePath           = ManageGatewayPath + obsidian.UrlSep + "state"
	ManageGatewayTierPath            = ManageGatewayPath + obsidian.UrlSep + "tier"
	ManageGatewayCellularPath        = ManageGatewayPath + obsidian.UrlSep + "cellular"
	ManageGatewayCellularEpcPath     = ManageGatewayCellularPath + obsidian.UrlSep + "epc"
	ManageGatewayCellularRanPath     = ManageGatewayCellularPath + obsidian.UrlSep + "ran"
	ManageGatewayCellularNonEpsPath  = ManageGatewayCellularPath + obsidian.UrlSep + "non_eps"
	ManageGatewayCellularEnodebsPath = ManageGatewayCellularPath + obsidian.UrlSep + "connected_enodeb_serial"

	Enodebs          = "enodebs"
	ListEnodebsPath  = ManageNetworkPath + obsidian.UrlSep + Enodebs
	ManageEnodebPath = ListEnodebsPath + obsidian.UrlSep + ":enodeb_serial"

	Subscribers              = "subscribers"
	ListSubscribersPath      = ManageNetworkPath + obsidian.UrlSep + Subscribers
	ManageSubscriberPath     = ListSubscribersPath + obsidian.UrlSep + ":subscriber_id"
	ActivateSubscriberPath   = ManageSubscriberPath + obsidian.UrlSep + "activate"
	DeactivateSubscriberPath = ManageSubscriberPath + obsidian.UrlSep + "deactivate"
)

func GetHandlers() []obsidian.Handler {
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

		{Path: ListGatewaysPath, Methods: obsidian.GET, HandlerFunc: listGateways},
		{Path: ListGatewaysPath, Methods: obsidian.POST, HandlerFunc: createGateway},
		{Path: ManageGatewayPath, Methods: obsidian.GET, HandlerFunc: getGateway},
		{Path: ManageGatewayPath, Methods: obsidian.PUT, HandlerFunc: updateGateway},
		{Path: ManageGatewayPath, Methods: obsidian.DELETE, HandlerFunc: deleteGateway},
		{Path: ManageGatewayStatePath, Methods: obsidian.GET, HandlerFunc: handlers.GetStateHandler},

		{Path: ListEnodebsPath, Methods: obsidian.GET, HandlerFunc: listEnodebs},
		{Path: ListEnodebsPath, Methods: obsidian.POST, HandlerFunc: createEnodeb},
		{Path: ManageEnodebPath, Methods: obsidian.GET, HandlerFunc: getEnodeb},
		{Path: ManageEnodebPath, Methods: obsidian.PUT, HandlerFunc: updateEnodeb},
		{Path: ManageEnodebPath, Methods: obsidian.DELETE, HandlerFunc: deleteEnodeb},

		{Path: ListSubscribersPath, Methods: obsidian.GET, HandlerFunc: listSubscribers},
		{Path: ListSubscribersPath, Methods: obsidian.POST, HandlerFunc: createSubscriber},
		{Path: ManageSubscriberPath, Methods: obsidian.GET, HandlerFunc: getSubscriber},
		{Path: ManageSubscriberPath, Methods: obsidian.PUT, HandlerFunc: updateSubscriber},
		{Path: ManageSubscriberPath, Methods: obsidian.DELETE, HandlerFunc: deleteSubscriber},
		{Path: ActivateSubscriberPath, Methods: obsidian.POST, HandlerFunc: makeSubscriberStateHandler(ltemodels.LteSubscriptionStateACTIVE)},
		{Path: DeactivateSubscriberPath, Methods: obsidian.POST, HandlerFunc: makeSubscriberStateHandler(ltemodels.LteSubscriptionStateINACTIVE)},
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

	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayNamePath, new(models.GatewayName))...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayDescriptionPath, new(models.GatewayDescription))...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayConfigPath, &orc8rmodels.MagmadGatewayConfigs{})...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayTierPath, new(orc8rmodels.TierID))...)
	ret = append(ret, handlers.GetGatewayDeviceHandlers(ManageGatewayDevicePath)...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayCellularPath, &ltemodels.GatewayCellularConfigs{})...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayCellularEpcPath, &ltemodels.GatewayEpcConfigs{})...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayCellularRanPath, &ltemodels.GatewayRanConfigs{})...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayCellularNonEpsPath, &ltemodels.GatewayNonEpsConfigs{})...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayCellularEnodebsPath, &ltemodels.EnodebSerials{})...)
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

func listGateways(c echo.Context) error {
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

func createGateway(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := &ltemodels.MutableLteGateway{}
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

func updateGateway(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	payload := &ltemodels.MutableLteGateway{}
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

func deleteGateway(c echo.Context) error {
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

func listEnodebs(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	ents, err := configurator.LoadAllEntitiesInNetwork(
		nid, lte.CellularEnodebType,
		configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsToThis: true},
	)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := make(map[string]*ltemodels.Enodeb, len(ents))
	for _, ent := range ents {
		ret[ent.Key] = (&ltemodels.Enodeb{}).FromBackendModels(ent)
	}
	return c.JSON(http.StatusOK, ret)
}

func createEnodeb(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := &ltemodels.Enodeb{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if payload.AttachedGatewayID != "" {
		return echo.NewHTTPError(http.StatusBadRequest, "attached_gateway_id is a read-only property")
	}

	_, err := configurator.CreateEntity(nid, configurator.NetworkEntity{
		Type:       lte.CellularEnodebType,
		Key:        payload.Serial,
		Name:       payload.Name,
		PhysicalID: payload.Serial,
		Config:     payload.Config,
	})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func getEnodeb(c echo.Context) error {
	nid, eid, nerr := getNetworkAndEnbIDs(c)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(
		nid, lte.CellularEnodebType, eid,
		configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsToThis: true},
	)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := (&ltemodels.Enodeb{}).FromBackendModels(ent)
	return c.JSON(http.StatusOK, ret)
}

func updateEnodeb(c echo.Context) error {
	nid, eid, nerr := getNetworkAndEnbIDs(c)
	if nerr != nil {
		return nerr
	}

	payload := &ltemodels.Enodeb{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if payload.AttachedGatewayID != "" {
		return echo.NewHTTPError(http.StatusBadRequest, "attached_gateway_id is a read-only property")
	}
	if payload.Serial != eid {
		return echo.NewHTTPError(http.StatusBadRequest, "serial in body must match serial in path")
	}

	_, err := configurator.UpdateEntity(nid, payload.ToEntityUpdateCriteria())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteEnodeb(c echo.Context) error {
	nid, eid, nerr := getNetworkAndEnbIDs(c)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(nid, lte.CellularEnodebType, eid)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func getNetworkAndEnbIDs(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", "enodeb_serial")
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}

func listSubscribers(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	ents, err := configurator.LoadAllEntitiesInNetwork(networkID, lte.SubscriberEntityType, configurator.EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := make(map[string]*ltemodels.Subscriber, len(ents))
	for _, ent := range ents {
		ret[ent.Key] = (&ltemodels.Subscriber{}).FromBackendModels(ent)
	}
	return c.JSON(http.StatusOK, ret)
}

func createSubscriber(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := &ltemodels.Subscriber{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	_, err := configurator.CreateEntity(networkID, configurator.NetworkEntity{
		Type:   lte.SubscriberEntityType,
		Key:    payload.ID,
		Config: payload.Lte,
	})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func getSubscriber(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(networkID, lte.SubscriberEntityType, subscriberID, configurator.EntityLoadCriteria{LoadConfig: true})
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := (&ltemodels.Subscriber{}).FromBackendModels(ent)
	return c.JSON(http.StatusOK, ret)
}

func updateSubscriber(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
	if nerr != nil {
		return nerr
	}

	payload := &ltemodels.Subscriber{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	_, err := configurator.LoadEntity(networkID, lte.SubscriberEntityType, subscriberID, configurator.EntityLoadCriteria{})
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(errors.Wrap(err, "failed to load existing subscriber"), http.StatusInternalServerError)
	}

	err = configurator.CreateOrUpdateEntityConfig(networkID, lte.SubscriberEntityType, subscriberID, payload.Lte)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteSubscriber(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(networkID, lte.SubscriberEntityType, subscriberID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func makeSubscriberStateHandler(desiredState string) echo.HandlerFunc {
	return func(c echo.Context) error {
		networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
		if nerr != nil {
			return nerr
		}

		cfg, err := configurator.LoadEntityConfig(networkID, lte.SubscriberEntityType, subscriberID)
		switch {
		case err == merrors.ErrNotFound:
			return echo.ErrNotFound
		case err != nil:
			return obsidian.HttpError(errors.Wrap(err, "failed to load existing subscriber"), http.StatusInternalServerError)
		}

		newConfig := cfg.(*ltemodels.LteSubscription)
		newConfig.State = desiredState
		err = configurator.CreateOrUpdateEntityConfig(networkID, lte.SubscriberEntityType, subscriberID, newConfig)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}

func getNetworkAndSubIDs(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", "subscriber_id")
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}
