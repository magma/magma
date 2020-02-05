/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"fmt"
	"net/http"

	"magma/lte/cloud/go/lte"
	ltemodels "magma/lte/cloud/go/plugin/models"
	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	orc8rmodels "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/storage"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	LteNetworks                         = "lte"
	ListNetworksPath                    = obsidian.V1Root + LteNetworks
	ManageNetworkPath                   = ListNetworksPath + "/:network_id"
	ManageNetworkNamePath               = ManageNetworkPath + obsidian.UrlSep + "name"
	ManageNetworkDescriptionPath        = ManageNetworkPath + obsidian.UrlSep + "description"
	ManageNetworkFeaturesPath           = ManageNetworkPath + obsidian.UrlSep + "features"
	ManageNetworkDNSPath                = ManageNetworkPath + obsidian.UrlSep + "dns"
	ManageNetworkDNSRecordsPath         = ManageNetworkDNSPath + obsidian.UrlSep + "records"
	ManageNetworkDNSRecordByDomainPath  = ManageNetworkDNSRecordsPath + obsidian.UrlSep + ":domain"
	ManageNetworkCellularPath           = ManageNetworkPath + obsidian.UrlSep + "cellular"
	ManageNetworkCellularEpcPath        = ManageNetworkCellularPath + obsidian.UrlSep + "epc"
	ManageNetworkCellularRanPath        = ManageNetworkCellularPath + obsidian.UrlSep + "ran"
	ManageNetworkCellularFegNetworkID   = ManageNetworkCellularPath + obsidian.UrlSep + "feg_network_id"
	ManageNetworkSubscriberPath         = ManageNetworkPath + obsidian.UrlSep + "subscriber_config"
	ManageNetworkBaseNamesPath          = ManageNetworkSubscriberPath + obsidian.UrlSep + "base_names"
	ManageNetworkRuleNamesPath          = ManageNetworkSubscriberPath + obsidian.UrlSep + "rule_names"
	ManageNetworkSubscriberRuleNamePath = ManageNetworkRuleNamesPath + obsidian.UrlSep + ":rule_id"
	ManageNetworkSubscriberBaseNamePath = ManageNetworkBaseNamesPath + obsidian.UrlSep + ":base_name"

	Gateways                          = "gateways"
	ListGatewaysPath                  = ManageNetworkPath + obsidian.UrlSep + Gateways
	ManageGatewayPath                 = ListGatewaysPath + obsidian.UrlSep + ":gateway_id"
	ManageGatewayNamePath             = ManageGatewayPath + obsidian.UrlSep + "name"
	ManageGatewayDescriptionPath      = ManageGatewayPath + obsidian.UrlSep + "description"
	ManageGatewayConfigPath           = ManageGatewayPath + obsidian.UrlSep + "magmad"
	ManageGatewayDevicePath           = ManageGatewayPath + obsidian.UrlSep + "device"
	ManageGatewayStatePath            = ManageGatewayPath + obsidian.UrlSep + "status"
	ManageGatewayTierPath             = ManageGatewayPath + obsidian.UrlSep + "tier"
	ManageGatewayCellularPath         = ManageGatewayPath + obsidian.UrlSep + "cellular"
	ManageGatewayCellularEpcPath      = ManageGatewayCellularPath + obsidian.UrlSep + "epc"
	ManageGatewayCellularRanPath      = ManageGatewayCellularPath + obsidian.UrlSep + "ran"
	ManageGatewayCellularNonEpsPath   = ManageGatewayCellularPath + obsidian.UrlSep + "non_eps"
	ManageGatewayConnectedEnodebsPath = ManageGatewayPath + obsidian.UrlSep + "connected_enodeb_serials"

	Enodebs            = "enodebs"
	ListEnodebsPath    = ManageNetworkPath + obsidian.UrlSep + Enodebs
	ManageEnodebPath   = ListEnodebsPath + obsidian.UrlSep + ":enodeb_serial"
	GetEnodebStatePath = ManageEnodebPath + obsidian.UrlSep + "state"

	Subscribers              = "subscribers"
	ListSubscribersPath      = ManageNetworkPath + obsidian.UrlSep + Subscribers
	ManageSubscriberPath     = ListSubscribersPath + obsidian.UrlSep + ":subscriber_id"
	ActivateSubscriberPath   = ManageSubscriberPath + obsidian.UrlSep + "activate"
	DeactivateSubscriberPath = ManageSubscriberPath + obsidian.UrlSep + "deactivate"
	SubscriberProfilePath    = ManageSubscriberPath + obsidian.UrlSep + "lte" + obsidian.UrlSep + "sub_profile"

	policiesRootPath         = handlers.ManageNetworkPath + obsidian.UrlSep + "policies"
	policyRuleRootPath       = policiesRootPath + obsidian.UrlSep + "rules"
	policyRuleManagePath     = policyRuleRootPath + obsidian.UrlSep + ":rule_id"
	policyBaseNameRootPath   = policiesRootPath + obsidian.UrlSep + "base_names"
	policyBaseNameManagePath = policyBaseNameRootPath + obsidian.UrlSep + ":base_name"

	ratingGroupsRootPath   = handlers.ManageNetworkPath + obsidian.UrlSep + "rating_groups"
	ratingGroupsManagePath = ratingGroupsRootPath + obsidian.UrlSep + ":rating_group_id"
)

func GetHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ManageNetworkDNSRecordByDomainPath, Methods: obsidian.POST, HandlerFunc: handlers.CreateDNSRecord},
		{Path: ManageNetworkDNSRecordByDomainPath, Methods: obsidian.GET, HandlerFunc: handlers.ReadDNSRecord},
		{Path: ManageNetworkDNSRecordByDomainPath, Methods: obsidian.PUT, HandlerFunc: handlers.UpdateDNSRecord},
		{Path: ManageNetworkDNSRecordByDomainPath, Methods: obsidian.DELETE, HandlerFunc: handlers.DeleteDNSRecord},

		handlers.GetListGatewaysHandler(ListGatewaysPath, lte.CellularGatewayType, makeLTEGateways),
		{Path: ListGatewaysPath, Methods: obsidian.POST, HandlerFunc: createGateway},
		{Path: ManageGatewayPath, Methods: obsidian.GET, HandlerFunc: getGateway},
		{Path: ManageGatewayStatePath, Methods: obsidian.GET, HandlerFunc: handlers.GetStateHandler},
		{Path: ManageGatewayPath, Methods: obsidian.PUT, HandlerFunc: updateGateway},
		handlers.GetDeleteGatewayHandler(ManageGatewayPath, lte.CellularGatewayType),

		{Path: ListEnodebsPath, Methods: obsidian.GET, HandlerFunc: listEnodebs},
		{Path: ListEnodebsPath, Methods: obsidian.POST, HandlerFunc: createEnodeb},
		{Path: ManageEnodebPath, Methods: obsidian.GET, HandlerFunc: getEnodeb},
		{Path: ManageEnodebPath, Methods: obsidian.PUT, HandlerFunc: updateEnodeb},
		{Path: ManageEnodebPath, Methods: obsidian.DELETE, HandlerFunc: deleteEnodeb},
		{Path: ManageGatewayConnectedEnodebsPath, Methods: obsidian.POST, HandlerFunc: addConnectedEnodeb},
		{Path: ManageGatewayConnectedEnodebsPath, Methods: obsidian.DELETE, HandlerFunc: deleteConnectedEnodeb},
		{Path: GetEnodebStatePath, Methods: obsidian.GET, HandlerFunc: getEnodebState},

		{Path: ListSubscribersPath, Methods: obsidian.GET, HandlerFunc: listSubscribers},
		{Path: ListSubscribersPath, Methods: obsidian.POST, HandlerFunc: createSubscriber},
		{Path: ManageSubscriberPath, Methods: obsidian.GET, HandlerFunc: getSubscriber},
		{Path: ManageSubscriberPath, Methods: obsidian.PUT, HandlerFunc: updateSubscriber},
		{Path: ManageSubscriberPath, Methods: obsidian.DELETE, HandlerFunc: deleteSubscriber},
		{Path: ActivateSubscriberPath, Methods: obsidian.POST, HandlerFunc: makeSubscriberStateHandler(ltemodels.LteSubscriptionStateACTIVE)},
		{Path: DeactivateSubscriberPath, Methods: obsidian.POST, HandlerFunc: makeSubscriberStateHandler(ltemodels.LteSubscriptionStateINACTIVE)},
		{Path: SubscriberProfilePath, Methods: obsidian.PUT, HandlerFunc: updateSubscriberProfile},
		{Path: ManageNetworkSubscriberBaseNamePath, Methods: obsidian.POST, HandlerFunc: AddNetworkWideSubscriberBaseName},
		{Path: ManageNetworkSubscriberRuleNamePath, Methods: obsidian.POST, HandlerFunc: AddNetworkWideSubscriberRuleName},
		{Path: ManageNetworkSubscriberBaseNamePath, Methods: obsidian.DELETE, HandlerFunc: RemoveNetworkWideSubscriberBaseName},
		{Path: ManageNetworkSubscriberRuleNamePath, Methods: obsidian.DELETE, HandlerFunc: RemoveNetworkWideSubscriberRuleName},

		{Path: policyBaseNameRootPath, Methods: obsidian.GET, HandlerFunc: ListBaseNames},
		{Path: policyBaseNameRootPath, Methods: obsidian.POST, HandlerFunc: CreateBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.GET, HandlerFunc: GetBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.PUT, HandlerFunc: UpdateBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.DELETE, HandlerFunc: DeleteBaseName},

		{Path: policyRuleRootPath, Methods: obsidian.GET, HandlerFunc: ListRules},
		{Path: policyRuleRootPath, Methods: obsidian.POST, HandlerFunc: CreateRule},
		{Path: policyRuleManagePath, Methods: obsidian.GET, HandlerFunc: GetRule},
		{Path: policyRuleManagePath, Methods: obsidian.PUT, HandlerFunc: UpdateRule},
		{Path: policyRuleManagePath, Methods: obsidian.DELETE, HandlerFunc: DeleteRule},

		{Path: ratingGroupsRootPath, Methods: obsidian.GET, HandlerFunc: ListRatingGroups},
		{Path: ratingGroupsRootPath, Methods: obsidian.POST, HandlerFunc: CreateRatingGroup},
		{Path: ratingGroupsManagePath, Methods: obsidian.GET, HandlerFunc: GetRatingGroup},
		{Path: ratingGroupsManagePath, Methods: obsidian.PUT, HandlerFunc: UpdateRatingGroup},
		{Path: ratingGroupsManagePath, Methods: obsidian.DELETE, HandlerFunc: DeleteRatingGroup},
	}
	ret = append(ret, handlers.GetTypedNetworkCRUDHandlers(ListNetworksPath, ManageNetworkPath, lte.LteNetworkType, &ltemodels.LteNetwork{})...)

	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkNamePath, new(models.NetworkName), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDescriptionPath, new(models.NetworkDescription), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkFeaturesPath, &orc8rmodels.NetworkFeatures{}, orc8r.NetworkFeaturesConfig)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDNSPath, &orc8rmodels.NetworkDNSConfig{}, orc8r.DnsdNetworkType)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDNSRecordsPath, new(orc8rmodels.NetworkDNSRecords), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkCellularPath, &ltemodels.NetworkCellularConfigs{}, lte.CellularNetworkType)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkCellularEpcPath, &ltemodels.NetworkEpcConfigs{}, "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkCellularRanPath, &ltemodels.NetworkRanConfigs{}, "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkCellularFegNetworkID, new(ltemodels.FegNetworkID), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkSubscriberPath, &ltemodels.NetworkSubscriberConfig{}, "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkRuleNamesPath, new(ltemodels.RuleNames), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkBaseNamesPath, new(ltemodels.BaseNames), "")...)

	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayNamePath, new(models.GatewayName))...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayDescriptionPath, new(models.GatewayDescription))...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayConfigPath, &orc8rmodels.MagmadGatewayConfigs{})...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayTierPath, new(orc8rmodels.TierID))...)
	ret = append(ret, handlers.GetGatewayDeviceHandlers(ManageGatewayDevicePath)...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayCellularPath, &ltemodels.GatewayCellularConfigs{})...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayCellularEpcPath, &ltemodels.GatewayEpcConfigs{})...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayCellularRanPath, &ltemodels.GatewayRanConfigs{})...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayCellularNonEpsPath, &ltemodels.GatewayNonEpsConfigs{})...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayConnectedEnodebsPath, &ltemodels.EnodebSerials{})...)
	return ret
}

func createGateway(c echo.Context) error {
	if nerr := handlers.CreateMagmadGatewayFromModel(c, &ltemodels.MutableLteGateway{}); nerr != nil {
		return nerr
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
	}
	if ent.Config != nil {
		ret.Cellular = ent.Config.(*ltemodels.GatewayCellularConfigs)
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
	if nerr = handlers.UpdateMagmadGatewayFromModel(c, nid, gid, &ltemodels.MutableLteGateway{}); nerr != nil {
		return nerr
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
) map[string]handlers.GatewayModel {
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

	ret := make(map[string]handlers.GatewayModel, len(gatewayEntsByKey))
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

func getEnodebState(c echo.Context) error {
	nid, eid, nerr := getNetworkAndEnbIDs(c)
	if nerr != nil {
		return nerr
	}
	st, err := state.GetState(nid, lte.EnodebStateType, eid)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	} else if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	enodebState := st.ReportedState.(*ltemodels.EnodebState)
	enodebState.TimeReported = st.TimeMs
	ent, err := configurator.LoadEntityForPhysicalID(st.ReporterID, configurator.EntityLoadCriteria{})
	if err == nil {
		enodebState.ReportingGatewayID = ent.Key
	}
	return c.JSON(http.StatusOK, enodebState)
}

func getNetworkAndEnbIDs(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", "enodeb_serial")
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}

func deleteConnectedEnodeb(c echo.Context) error {
	networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	var enodebSerial string
	if err := c.Bind(&enodebSerial); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	_, err := configurator.UpdateEntity(networkID, (&ltemodels.EnodebSerials{}).ToDeleteUpdateCriteria(networkID, gatewayID, enodebSerial))
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func addConnectedEnodeb(c echo.Context) error {
	networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	var enodebSerial string
	if err := c.Bind(&enodebSerial); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	_, err := configurator.UpdateEntity(networkID, (&ltemodels.EnodebSerials{}).ToCreateUpdateCriteria(networkID, gatewayID, enodebSerial))
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
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

	if nerr := validateSubscriberProfile(networkID, payload.Lte); nerr != nil {
		return nerr
	}

	_, err := configurator.CreateEntity(networkID, configurator.NetworkEntity{
		Type:   lte.SubscriberEntityType,
		Key:    string(payload.ID),
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

	if nerr := validateSubscriberProfile(networkID, payload.Lte); nerr != nil {
		return nerr
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

func updateSubscriberProfile(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
	if nerr != nil {
		return nerr
	}

	var payload = new(ltemodels.SubProfile)
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	currentCfg, err := configurator.LoadEntityConfig(networkID, lte.SubscriberEntityType, subscriberID)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(errors.Wrap(err, "could not load subscriber"), http.StatusInternalServerError)
	}

	desiredCfg := currentCfg.(*ltemodels.LteSubscription)
	desiredCfg.SubProfile = *payload
	if nerr := validateSubscriberProfile(networkID, desiredCfg); nerr != nil {
		return nerr
	}

	_, err = configurator.UpdateEntity(networkID, configurator.EntityUpdateCriteria{Type: lte.SubscriberEntityType, Key: subscriberID, NewConfig: desiredCfg})
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to update profile"), http.StatusInternalServerError)
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

func AddNetworkWideSubscriberRuleName(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	params, nerr := obsidian.GetParamValues(c, "rule_id")
	if nerr != nil {
		return nerr
	}
	err := addToNetworkSubscriberConfig(networkID, params[0], "")
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "Failed to update config"), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func AddNetworkWideSubscriberBaseName(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	params, nerr := obsidian.GetParamValues(c, "base_name")
	if nerr != nil {
		return nerr
	}
	err := addToNetworkSubscriberConfig(networkID, "", params[0])
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "Failed to update config"), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func RemoveNetworkWideSubscriberRuleName(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	params, nerr := obsidian.GetParamValues(c, "rule_id")
	if nerr != nil {
		return nerr
	}
	err := removeFromNetworkSubscriberConfig(networkID, params[0], "")
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "Failed to update config"), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func RemoveNetworkWideSubscriberBaseName(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	params, nerr := obsidian.GetParamValues(c, "base_name")
	if nerr != nil {
		return nerr
	}
	err := removeFromNetworkSubscriberConfig(networkID, "", params[0])
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "Failed to update config"), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func addToNetworkSubscriberConfig(networkID, ruleName, baseName string) error {
	network, err := configurator.LoadNetwork(networkID, false, true)
	if err != nil {
		return err
	}
	iSubscriberConfig, exists := network.Configs[lte.NetworkSubscriberConfigType]
	if !exists {
		network.Configs[lte.NetworkSubscriberConfigType] = &ltemodels.NetworkSubscriberConfig{}
	}
	subscriberConfig, ok := iSubscriberConfig.(*ltemodels.NetworkSubscriberConfig)
	if !ok {
		return fmt.Errorf("Unable to convert config")
	}
	if len(ruleName) != 0 {
		subscriberConfig.NetworkWideRuleNames = append(subscriberConfig.NetworkWideRuleNames, ruleName)
	}
	if len(baseName) != 0 {
		subscriberConfig.NetworkWideBaseNames = append(subscriberConfig.NetworkWideBaseNames, ltemodels.BaseName(baseName))
	}
	return configurator.UpdateNetworkConfig(networkID, lte.NetworkSubscriberConfigType, subscriberConfig)
}

func removeFromNetworkSubscriberConfig(networkID, ruleName, baseName string) error {
	network, err := configurator.LoadNetwork(networkID, false, true)
	if err != nil {
		return err
	}
	iSubscriberConfig, exists := network.Configs[lte.NetworkSubscriberConfigType]
	if !exists {
		network.Configs[lte.NetworkSubscriberConfigType] = &ltemodels.NetworkSubscriberConfig{}
	}
	subscriberConfig, ok := iSubscriberConfig.(*ltemodels.NetworkSubscriberConfig)
	if !ok {
		return fmt.Errorf("Unable to convert config")
	}
	if len(ruleName) != 0 {
		subscriberConfig.NetworkWideRuleNames = funk.FilterString(subscriberConfig.NetworkWideRuleNames,
			func(s string) bool { return s != ruleName })
	}
	if len(baseName) != 0 {
		subscriberConfig.NetworkWideBaseNames = funk.Filter(subscriberConfig.NetworkWideBaseNames,
			func(b ltemodels.BaseName) bool { return string(b) != baseName }).([]ltemodels.BaseName)
	}
	return configurator.UpdateNetworkConfig(networkID, lte.NetworkSubscriberConfigType, subscriberConfig)
}

func getNetworkAndSubIDs(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", "subscriber_id")
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}

func validateSubscriberProfile(networkID string, sub *ltemodels.LteSubscription) *echo.HTTPError {
	// Check the sub profiles available on the network if sub profile is not
	// default (which is always available)
	if sub.SubProfile != "default" {
		netConf, err := configurator.LoadNetworkConfig(networkID, lte.CellularNetworkType)
		switch {
		case err == merrors.ErrNotFound:
			return obsidian.HttpError(errors.New("no cellular config found for network"), http.StatusInternalServerError)
		case err != nil:
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		cellNetConf := netConf.(*ltemodels.NetworkCellularConfigs)
		profName := string(sub.SubProfile)
		if _, profileExists := cellNetConf.Epc.SubProfiles[profName]; !profileExists {
			return obsidian.HttpError(errors.Errorf("subscriber profile %s does not exist for the network", profName), http.StatusBadRequest)
		}
	}
	return nil
}
