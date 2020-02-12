/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	orc8rmodels "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
	"orc8r/devmand/cloud/go/devmand"
	symphonymodels "orc8r/devmand/cloud/go/plugin/models"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	SymphonyNetworks             = "symphony"
	BaseNetworksPath             = obsidian.V1Root + SymphonyNetworks
	ManageNetworkPath            = BaseNetworksPath + obsidian.UrlSep + ":network_id"
	ManageNetworkNamePath        = ManageNetworkPath + obsidian.UrlSep + "name"
	ManageNetworkDescriptionPath = ManageNetworkPath + obsidian.UrlSep + "description"
	ManageNetworkFeaturesPath    = ManageNetworkPath + obsidian.UrlSep + "features"

	AgentID                       = "agent_id"
	BaseAgentsPath                = ManageNetworkPath + obsidian.UrlSep + "agents"
	ManageAgentPath               = BaseAgentsPath + obsidian.UrlSep + ":agent_id"
	ManageAgentNamePath           = ManageAgentPath + obsidian.UrlSep + "name"
	ManageAgentDescriptionPath    = ManageAgentPath + obsidian.UrlSep + "description"
	ManageAgentConfigPath         = ManageAgentPath + obsidian.UrlSep + "magmad"
	ManageAgentDevicePath         = ManageAgentPath + obsidian.UrlSep + "device"
	ManageAgentStatePath          = ManageAgentPath + obsidian.UrlSep + "status"
	ManageAgentTierPath           = ManageAgentPath + obsidian.UrlSep + "tier"
	ManageAgentManagedDevicesPath = ManageAgentPath + obsidian.UrlSep + "managed_devices"

	DeviceID               = "device_id"
	BaseDevicesPath        = ManageNetworkPath + obsidian.UrlSep + "devices"
	ManageDevicePath       = BaseDevicesPath + obsidian.UrlSep + ":device_id"
	ManageDeviceNamePath   = ManageDevicePath + obsidian.UrlSep + "name"
	ManageDeviceConfigPath = ManageDevicePath + obsidian.UrlSep + "config"
	ManageDeviceAgent      = ManageDevicePath + obsidian.UrlSep + "managing_agent"
	GetDeviceStatePath     = ManageDevicePath + obsidian.UrlSep + "state"
)

// GetHandlers returns all obsidian handlers for Symphony
func GetHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		handlers.GetListGatewaysHandler(BaseAgentsPath, devmand.SymphonyAgentType, makeSymphonyAgents),
		{Path: BaseAgentsPath, Methods: obsidian.POST, HandlerFunc: createAgent},
		{Path: ManageAgentPath, Methods: obsidian.GET, HandlerFunc: getAgent},
		{Path: ManageAgentPath, Methods: obsidian.PUT, HandlerFunc: updateAgent},
		{Path: ManageAgentPath, Methods: obsidian.DELETE, HandlerFunc: deleteAgent},

		{Path: BaseDevicesPath, Methods: obsidian.GET, HandlerFunc: listDevices},
		{Path: BaseDevicesPath, Methods: obsidian.POST, HandlerFunc: createDevice},
		{Path: ManageDevicePath, Methods: obsidian.GET, HandlerFunc: getDevice},
		{Path: ManageDevicePath, Methods: obsidian.PUT, HandlerFunc: updateDevice},
		{Path: ManageDevicePath, Methods: obsidian.DELETE, HandlerFunc: deleteDevice},
		{Path: GetDeviceStatePath, Methods: obsidian.GET, HandlerFunc: GetDeviceState},
	}
	ret = append(ret, handlers.GetTypedNetworkCRUDHandlers(BaseNetworksPath, ManageNetworkPath, devmand.SymphonyNetworkType, &symphonymodels.SymphonyNetwork{})...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkNamePath, new(models.NetworkName), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDescriptionPath, new(models.NetworkDescription), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkFeaturesPath, &orc8rmodels.NetworkFeatures{}, orc8r.NetworkFeaturesConfig)...)

	ret = append(ret, handlers.GetPartialEntityHandlers(ManageAgentNamePath, AgentID, new(models.GatewayName))...)
	ret = append(ret, handlers.GetPartialEntityHandlers(ManageAgentDescriptionPath, AgentID, new(models.GatewayDescription))...)
	ret = append(ret, handlers.GetPartialEntityHandlers(ManageAgentConfigPath, AgentID, &orc8rmodels.MagmadGatewayConfigs{})...)
	ret = append(ret, handlers.GetPartialEntityHandlers(ManageAgentTierPath, AgentID, new(orc8rmodels.TierID))...)
	ret = append(ret, handlers.GetPartialEntityHandlers(ManageAgentManagedDevicesPath, AgentID, &symphonymodels.ManagedDevices{})...)
	ret = append(ret, GetAgentDeviceHandlers(ManageAgentDevicePath)...)

	ret = append(ret, handlers.GetPartialEntityHandlers(ManageDeviceNamePath, DeviceID, new(symphonymodels.SymphonyDeviceName))...)
	ret = append(ret, handlers.GetPartialEntityHandlers(ManageDeviceConfigPath, DeviceID, &symphonymodels.SymphonyDeviceConfig{})...)
	ret = append(ret, handlers.GetPartialEntityHandlers(ManageDeviceAgent, DeviceID, new(symphonymodels.SymphonyDeviceAgent))...)

	return ret
}

type agentAndMagmadGatewayEntities struct {
	agentEnt, magmadEnt configurator.NetworkEntity
}

func makeSymphonyAgents(
	entsByTK map[storage.TypeAndKey]configurator.NetworkEntity,
	devicesByID map[string]interface{},
	statusesByID map[string]*orc8rmodels.GatewayStatus,
) map[string]handlers.GatewayModel {
	agentEntsByKey := map[string]*agentAndMagmadGatewayEntities{}
	for tk, ent := range entsByTK {
		existing, found := agentEntsByKey[tk.Key]
		if !found {
			existing = &agentAndMagmadGatewayEntities{}
			agentEntsByKey[tk.Key] = existing
		}
		switch ent.Type {
		case orc8r.MagmadGatewayType:
			existing.magmadEnt = ent
		case devmand.SymphonyAgentType:
			existing.agentEnt = ent
		}
	}

	ret := make(map[string]handlers.GatewayModel, len(agentEntsByKey))
	for key, aMEnts := range agentEntsByKey {
		hwID := aMEnts.magmadEnt.PhysicalID
		var devCasted *orc8rmodels.GatewayDevice
		if devicesByID[hwID] != nil {
			devCasted = devicesByID[hwID].(*orc8rmodels.GatewayDevice)
		}
		ret[key] = (&symphonymodels.SymphonyAgent{}).FromBackendModels(aMEnts.magmadEnt, aMEnts.agentEnt, devCasted, statusesByID[hwID])
	}
	return ret
}

func listNetworks(c echo.Context) error {
	ids, err := configurator.ListNetworksOfType(devmand.SymphonyNetworkType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	sort.Strings(ids)
	return c.JSON(http.StatusOK, ids)
}

func createNetwork(c echo.Context) error {
	payload := &symphonymodels.SymphonyNetwork{}
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
	if network.Type != devmand.SymphonyNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not a Symphony network", nid))
	}

	ret := (&symphonymodels.SymphonyNetwork{}).FromConfiguratorNetwork(network)
	return c.JSON(http.StatusOK, ret)
}

func updateNetwork(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := &symphonymodels.SymphonyNetwork{}
	err := c.Bind(payload)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	// check that this is actually a Symphony network
	network, err := configurator.LoadNetwork(nid, false, false)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load network to check type"), http.StatusInternalServerError)
	}
	if network.Type != devmand.SymphonyNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not a Symphony network", nid))
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

	// check that this is actually a Symphony network
	network, err := configurator.LoadNetwork(nid, false, false)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load network to check type"), http.StatusInternalServerError)
	}
	if network.Type != devmand.SymphonyNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not a Symphony network", nid))
	}

	err = configurator.DeleteNetwork(nid)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func listAgents(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	ids, err := configurator.ListEntityKeys(nid, devmand.SymphonyAgentType)
	if err != nil {
		if err == merrors.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	sort.Strings(ids)
	return c.JSON(http.StatusOK, ids)
}

func createAgent(c echo.Context) error {
	if err := handlers.CreateMagmadGatewayFromModel(c, &symphonymodels.MutableSymphonyAgent{}); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func getAgent(c echo.Context) error {
	nid, aid, nerr := GetNetworkAndAgentIDs(c)
	if nerr != nil {
		return nerr
	}

	magmadGWModel, nerr := handlers.LoadMagmadGatewayModel(nid, aid)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(
		nid, devmand.SymphonyAgentType, aid, configurator.FullEntityLoadCriteria(),
	)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load symphony agent"), http.StatusInternalServerError)
	}

	ret := &symphonymodels.SymphonyAgent{
		ID:          magmadGWModel.ID,
		Name:        magmadGWModel.Name,
		Description: magmadGWModel.Description,
		Device:      magmadGWModel.Device,
		Tier:        magmadGWModel.Tier,
		Magmad:      magmadGWModel.Magmad,
	}

	for _, tk := range ent.Associations {
		if tk.Type == devmand.SymphonyDeviceType {
			ret.ManagedDevices = append(ret.ManagedDevices, tk.Key)
		}
	}
	return c.JSON(http.StatusOK, ret)
}

func updateAgent(c echo.Context) error {
	nid, aid, nerr := GetNetworkAndAgentIDs(c)
	if nerr != nil {
		return nerr
	}
	if nerr = handlers.UpdateMagmadGatewayFromModel(c, nid, aid, &symphonymodels.MutableSymphonyAgent{}); nerr != nil {
		return nerr
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteAgent(c echo.Context) error {
	nid, aid, nerr := GetNetworkAndAgentIDs(c)
	if nerr != nil {
		return nerr
	}

	// Deleting a configurator entity will remove assocs but not the
	// ents those assocs lead to, so just delete the agent ents
	updates := []configurator.EntityWriteOperation{
		configurator.EntityUpdateCriteria{
			Type: devmand.SymphonyAgentType, Key: aid, DeleteEntity: true,
		},
		configurator.EntityUpdateCriteria{
			Type: orc8r.MagmadGatewayType, Key: aid, DeleteEntity: true,
		},
	}
	err := configurator.WriteEntities(nid, updates...)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func listDevices(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	ents, err := configurator.LoadAllEntitiesInNetwork(
		nid, devmand.SymphonyDeviceType, configurator.FullEntityLoadCriteria(),
	)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := make(map[string]*symphonymodels.SymphonyDevice, len(ents))
	for _, ent := range ents {
		ret[ent.Key] = (&symphonymodels.SymphonyDevice{}).FromBackendModels(ent)
	}
	return c.JSON(http.StatusOK, ret)
}

func createDevice(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := &symphonymodels.MutableSymphonyDevice{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	// Create device<->agent associations
	writes := []configurator.EntityWriteOperation{}
	writes = append(writes, configurator.NetworkEntity{
		Type:   devmand.SymphonyDeviceType,
		Key:    string(payload.ID),
		Name:   string(payload.Name),
		Config: payload.Config,
	})
	for _, update := range symphonymodels.GetAgentUpdates(string(payload.ID), "", string(payload.ManagingAgent)) {
		writes = append(writes, update)
	}
	err := configurator.WriteEntities(nid, writes...)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func getDevice(c echo.Context) error {
	nid, did, nerr := GetNetworkAndDeviceIDs(c)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(nid, devmand.SymphonyDeviceType, did, configurator.FullEntityLoadCriteria())
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := (&symphonymodels.SymphonyDevice{}).FromBackendModels(ent)
	return c.JSON(http.StatusOK, ret)
}

func updateDevice(c echo.Context) error {
	nid, did, nerr := GetNetworkAndDeviceIDs(c)
	if nerr != nil {
		return nerr
	}

	payload := &symphonymodels.MutableSymphonyDevice{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if string(payload.ID) != did {
		return echo.NewHTTPError(http.StatusBadRequest, "device ID in body must match device_id in path")
	}

	deviceUpdates, err := payload.ToEntityUpdateCriteria(nid)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	_, err = configurator.UpdateEntities(nid, deviceUpdates)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteDevice(c echo.Context) error {
	nid, did, nerr := GetNetworkAndDeviceIDs(c)
	if nerr != nil {
		return nerr
	}

	exists, err := configurator.DoesEntityExist(nid, devmand.SymphonyDeviceType, did)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if !exists {
		return echo.ErrNotFound
	}

	err = configurator.DeleteEntity(nid, devmand.SymphonyDeviceType, did)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func GetDeviceState(c echo.Context) error {
	nid, did, nerr := GetNetworkAndDeviceIDs(c)
	if nerr != nil {
		return nerr
	}
	state, err := state.GetState(nid, devmand.SymphonyDeviceStateType, did)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	} else if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	deviceState := state.ReportedState.(*symphonymodels.SymphonyDeviceState)
	return c.JSON(http.StatusOK, deviceState)
}

func GetNetworkAndAgentIDs(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", "agent_id")
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}

func GetNetworkAndDeviceIDs(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", "device_id")
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}
