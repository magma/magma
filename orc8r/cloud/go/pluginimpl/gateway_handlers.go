/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package pluginimpl

import (
	"net/http"
	"sort"

	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/storage"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func ListGateways(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	ids, err := configurator.ListEntityKeys(nid, orc8r.MagmadGatewayType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	sort.Strings(ids)
	return c.JSON(http.StatusOK, ids)
}

func CreateGateway(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := &models.MagmadGateway{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	// must associate to an existing tier
	tierExists, err := configurator.DoesEntityExist(nid, orc8r.UpgradeTierEntityType, string(payload.Tier))
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to check for tier existence"), http.StatusInternalServerError)
	}
	if !tierExists {
		return echo.NewHTTPError(http.StatusBadRequest, "requested tier does not exist")
	}

	// If the device is already registered, throw an error if it's already
	// assigned to an entity
	// If the device exists but is unassigned, update it to the payload
	// If the device doesn't exist, create it and move on
	deviceID := payload.Device.HardwareID
	_, err = device.GetDevice(nid, orc8r.AccessGatewayRecordType, deviceID)
	switch {
	case err == merrors.ErrNotFound:
		err = device.RegisterDevice(nid, orc8r.AccessGatewayRecordType, deviceID, payload.Device)
		if err != nil {
			return obsidian.HttpError(errors.Wrap(err, "failed to register physical device"), http.StatusInternalServerError)
		}
		break
	case err != nil:
		return obsidian.HttpError(errors.Wrap(err, "failed to check if physical device is already registered"), http.StatusInternalServerError)
	default: // err == nil
		assignedEnt, err := configurator.LoadEntityForPhysicalID(deviceID, configurator.EntityLoadCriteria{})
		switch {
		case err == nil:
			return obsidian.HttpError(errors.Errorf("device %s is already mapped to gateway %s", deviceID, assignedEnt.Key), http.StatusBadRequest)
		case err != merrors.ErrNotFound:
			return obsidian.HttpError(errors.Wrap(err, "failed to check for existing device assignment"), http.StatusInternalServerError)
		}

		if err := device.UpdateDevice(nid, orc8r.AccessGatewayRecordType, deviceID, payload.Device); err != nil {
			return obsidian.HttpError(errors.Wrap(err, "failed to update device record"), http.StatusInternalServerError)
		}
	}

	// create the thing
	if _, err := configurator.CreateEntities(nid, payload.ToConfiguratorEntities()); err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to create gateway"), http.StatusInternalServerError)
	}

	// update the tier
	_, err = configurator.UpdateEntity(nid, configurator.EntityUpdateCriteria{
		Type:              orc8r.UpgradeTierEntityType,
		Key:               string(payload.Tier),
		AssociationsToAdd: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: string(payload.ID)}},
	})
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed up associate gateway to upgrade tier"), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func GetGateway(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(
		nid, orc8r.MagmadGatewayType, gid,
		configurator.EntityLoadCriteria{
			LoadMetadata:       true,
			LoadConfig:         true,
			LoadAssocsToThis:   true,
			LoadAssocsFromThis: false,
		},
	)
	if err == merrors.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	dev, err := device.GetDevice(nid, orc8r.AccessGatewayRecordType, ent.PhysicalID)
	if err != nil && err != merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	status, err := state.GetGatewayStatus(nid, ent.PhysicalID)
	if err != nil && err != merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// If the gateway/network is malformed, we could get no corresponding
	// device for the gateway
	var devCasted *models.GatewayDevice
	if dev != nil {
		devCasted = dev.(*models.GatewayDevice)
	}
	ret := (&models.MagmadGateway{}).FromBackendModels(ent, devCasted, status)
	return c.JSON(http.StatusOK, ret)
}

func UpdateGateway(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	payload := &models.MagmadGateway{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	// load the old ent to check if tier changed
	existingEnt, err := configurator.LoadEntity(
		nid, orc8r.MagmadGatewayType, gid,
		configurator.EntityLoadCriteria{LoadMetadata: true, LoadAssocsToThis: true},
	)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(errors.Wrap(err, "failed to load gateway"), http.StatusInternalServerError)
	}

	err = device.UpdateDevice(nid, orc8r.AccessGatewayRecordType, payload.Device.HardwareID, payload.Device)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to update device info"), http.StatusInternalServerError)
	}

	_, err = configurator.UpdateEntities(nid, payload.ToEntityUpdateCriteria(existingEnt))
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func DeleteGateway(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(nid, orc8r.MagmadGatewayType, gid)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// Unclear if we should be deleting the device as well. Until we get some
	// datapoints from real world usage, let's skip that for now

	return c.NoContent(http.StatusNoContent)
}
