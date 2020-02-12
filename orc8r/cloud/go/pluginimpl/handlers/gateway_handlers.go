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

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// MagmadEncompassingGateway is an interface for a gateway API model which
// wraps the magmad gateway with more fields that translate into additional
// network entities in the storage layer.
// Note that *models.MagmadGateway itself implements this interface as well.
type MagmadEncompassingGateway interface {
	// ValidatableModel allows the model to be validated by calling
	// ValidateModel()
	serde.ValidatableModel

	// GetMagmadGateway returns the *models.MagmadGateway which is wrapped by
	// the model
	GetMagmadGateway() *models.MagmadGateway

	// GetAdditionalWritesOnCreate returns extra write operations to perform
	// during creation, inside the same backend transaction as the creation
	// of the magmad gateway.
	// Do NOT include the write operation for the magmad gateway in the return
	// value, as *models.MagmadGateway itself implements this interface and
	// will create itself.
	GetAdditionalWritesOnCreate() []configurator.EntityWriteOperation

	// GetAdditionalEntitiesToLoadOnUpdate is a **static** method which
	// returns a list of entity keys to load in addition to the magmad gateway
	// during an update operation. The gateway ID from the API URL parameter
	// is given as an argument.
	// The entities loaded during this operation will be passed to
	// `GetAdditionalWritesOnUpdate`.
	GetAdditionalEntitiesToLoadOnUpdate(gatewayID string) []storage.TypeAndKey

	// GetAdditionalWritesOnUpdate returns extra write operations to perform
	// during a top-level update, inside the same backend transaction as the
	// update of the magmad gateway.
	// The gateway ID from the API URL parameter is given as an argument.
	// Do NOT include the write operation for the magmad gateway in the return
	// value, as *models.MagmadGateway itself implements this interface and
	// will update itself.
	GetAdditionalWritesOnUpdate(
		gatewayID string,
		loadedEntities map[storage.TypeAndKey]configurator.NetworkEntity,
	) ([]configurator.EntityWriteOperation, error)
}

func ListGatewaysHandler(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	ents, _, err := configurator.LoadEntities(nid, swag.String(orc8r.MagmadGatewayType), nil, nil, nil, configurator.FullEntityLoadCriteria())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	entsByTK := ents.ToEntitiesByID()

	// for each magmad gateway, we have to load its corresponding device and
	// its reported status
	deviceIDs := make([]string, 0, len(entsByTK))
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
	return c.JSON(http.StatusOK, makeGateways(entsByTK, devicesByID, statusesByID))
}

func CreateGatewayHandler(c echo.Context) error {
	if nerr := CreateMagmadGatewayFromModel(c, &models.MagmadGateway{}); nerr != nil {
		return nerr
	}
	return c.NoContent(http.StatusCreated)
}

func CreateMagmadGatewayFromModel(c echo.Context, model MagmadEncompassingGateway) *echo.HTTPError {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload, nerr := GetAndValidatePayload(c, model)
	if nerr != nil {
		return nerr
	}
	encompassingGateway := payload.(MagmadEncompassingGateway)
	mdGateway := encompassingGateway.GetMagmadGateway()

	// must associate to an existing tier
	tierExists, err := configurator.DoesEntityExist(nid, orc8r.UpgradeTierEntityType, string(mdGateway.Tier))
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
	deviceID := mdGateway.Device.HardwareID
	_, err = device.GetDevice(nid, orc8r.AccessGatewayRecordType, deviceID)
	switch {
	case err == merrors.ErrNotFound:
		err = device.RegisterDevice(nid, orc8r.AccessGatewayRecordType, deviceID, mdGateway.Device)
		if err != nil {
			return obsidian.HttpError(errors.Wrap(err, "failed to register physical device"), http.StatusInternalServerError)
		}
		break
	case err != nil:
		return obsidian.HttpError(errors.Wrap(err, "failed to check if physical device is already registered"), http.StatusConflict)
	default: // err == nil
		assignedEnt, err := configurator.LoadEntityForPhysicalID(deviceID, configurator.EntityLoadCriteria{})
		switch {
		case err == nil:
			return obsidian.HttpError(errors.Errorf("device %s is already mapped to gateway %s", deviceID, assignedEnt.Key), http.StatusBadRequest)
		case err != merrors.ErrNotFound:
			return obsidian.HttpError(errors.Wrap(err, "failed to check for existing device assignment"), http.StatusInternalServerError)
		}

		if err := device.UpdateDevice(nid, orc8r.AccessGatewayRecordType, deviceID, mdGateway.Device); err != nil {
			return obsidian.HttpError(errors.Wrap(err, "failed to update device record"), http.StatusInternalServerError)
		}
	}

	// create the magmad gateway, update the tier, perform additional writes
	// as necessary
	writes := []configurator.EntityWriteOperation{}
	writes = append(writes, mdGateway.GetAdditionalWritesOnCreate()...)
	writes = append(writes, configurator.EntityUpdateCriteria{
		Type:              orc8r.UpgradeTierEntityType,
		Key:               string(mdGateway.Tier),
		AssociationsToAdd: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: string(mdGateway.ID)}},
	})
	// these type switches aren't great but it's the best I could think of
	switch payload.(type) {
	case *models.MagmadGateway:
		break
	default:
		writes = append(writes, encompassingGateway.GetAdditionalWritesOnCreate()...)
	}

	if err = configurator.WriteEntities(nid, writes...); err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to create gateway"), http.StatusInternalServerError)
	}
	return nil
}

func GetGatewayHandler(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}
	ret, nerr := LoadMagmadGatewayModel(nid, gid)
	if nerr != nil {
		return nerr
	}
	return c.JSON(http.StatusOK, ret)
}

func LoadMagmadGatewayModel(networkID string, gatewayID string) (*models.MagmadGateway, *echo.HTTPError) {
	ent, err := configurator.LoadEntity(
		networkID, orc8r.MagmadGatewayType, gatewayID,
		configurator.EntityLoadCriteria{
			LoadMetadata:       true,
			LoadConfig:         true,
			LoadAssocsToThis:   true,
			LoadAssocsFromThis: false,
		},
	)
	if err == merrors.ErrNotFound {
		return nil, echo.ErrNotFound
	}
	if err != nil {
		return nil, obsidian.HttpError(err, http.StatusInternalServerError)
	}

	dev, err := device.GetDevice(networkID, orc8r.AccessGatewayRecordType, ent.PhysicalID)
	if err != nil && err != merrors.ErrNotFound {
		return nil, obsidian.HttpError(err, http.StatusInternalServerError)
	}
	status, err := state.GetGatewayStatus(networkID, ent.PhysicalID)
	if err != nil && err != merrors.ErrNotFound {
		return nil, obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// If the gateway/network is malformed, we could get no corresponding
	// device for the gateway
	var devCasted *models.GatewayDevice
	if dev != nil {
		devCasted = dev.(*models.GatewayDevice)
	}
	return (&models.MagmadGateway{}).FromBackendModels(ent, devCasted, status), nil
}

func UpdateGatewayHandler(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	if nerr = UpdateMagmadGatewayFromModel(c, nid, gid, &models.MagmadGateway{}); nerr != nil {
		return nerr
	}
	return c.NoContent(http.StatusNoContent)
}

func UpdateMagmadGatewayFromModel(c echo.Context, nid string, gid string, model MagmadEncompassingGateway) *echo.HTTPError {
	payload, nerr := GetAndValidatePayload(c, model)
	if nerr != nil {
		return nerr
	}
	encompassingGateway := payload.(MagmadEncompassingGateway)
	mdGateway := encompassingGateway.GetMagmadGateway()

	entsToLoad := []storage.TypeAndKey{}
	entsToLoad = append(entsToLoad, mdGateway.GetAdditionalEntitiesToLoadOnUpdate(gid)...)
	switch payload.(type) {
	case *models.MagmadGateway:
		break
	default:
		entsToLoad = append(entsToLoad, encompassingGateway.GetAdditionalEntitiesToLoadOnUpdate(gid)...)
	}

	loadedEnts, _, err := configurator.LoadEntities(
		nid,
		nil, nil, nil,
		entsToLoad,
		configurator.FullEntityLoadCriteria(),
	)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load gateway before update"), http.StatusInternalServerError)
	}

	writes, nerr := getUpdateWrites(gid, encompassingGateway, loadedEnts)
	if nerr != nil {
		return nerr
	}

	err = configurator.WriteEntities(nid, writes...)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// device info is cheap to update, so just do it all the time if
	// configurator write was successful
	err = device.UpdateDevice(nid, orc8r.AccessGatewayRecordType, mdGateway.Device.HardwareID, mdGateway.Device)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to update device info"), http.StatusInternalServerError)
	}

	return nil
}

func getUpdateWrites(gid string, payload MagmadEncompassingGateway, loadedEnts configurator.NetworkEntities) ([]configurator.EntityWriteOperation, *echo.HTTPError) {
	writes := []configurator.EntityWriteOperation{}
	loadedEntsByID := loadedEnts.ToEntitiesByID()

	mdGwWrites, err := payload.GetMagmadGateway().GetAdditionalWritesOnUpdate(gid, loadedEntsByID)
	switch {
	case err == merrors.ErrNotFound:
		return writes, echo.ErrNotFound
	case err != nil:
		return writes, obsidian.HttpError(errors.Wrap(err, "failed to get upate operations from magmad model"), http.StatusInternalServerError)
	}

	// short circuit if this is the magmad gateway
	switch payload.(type) {
	case *models.MagmadGateway:
		return mdGwWrites, nil
	}

	payloadWrites, err := payload.GetAdditionalWritesOnUpdate(gid, loadedEntsByID)
	switch {
	case err == merrors.ErrNotFound:
		return writes, echo.ErrNotFound
	case err != nil:
		return writes, obsidian.HttpError(errors.Wrap(err, "failed to get upate operations from payload model"), http.StatusInternalServerError)
	}

	writes = append(writes, mdGwWrites...)
	writes = append(writes, payloadWrites...)
	return writes, nil
}

func DeleteGatewayHandler(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

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

	err = configurator.DeleteEntity(nid, orc8r.MagmadGatewayType, gid)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	if existingEnt.PhysicalID != "" {
		err = device.DeleteDevice(nid, orc8r.AccessGatewayRecordType, existingEnt.PhysicalID)
		if err != nil {
			return obsidian.HttpError(errors.Wrap(err, "failed to delete device for gateway"), http.StatusInternalServerError)
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func GetStateHandler(c echo.Context) error {
	networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	physicalID, err := configurator.GetPhysicalIDOfEntity(networkID, orc8r.MagmadGatewayType, gatewayID)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	} else if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	st, err := state.GetGatewayStatus(networkID, physicalID)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	} else if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, st)
}

func makeGateways(
	entsByTK map[storage.TypeAndKey]configurator.NetworkEntity,
	devicesByID map[string]interface{},
	statusesByID map[string]*models.GatewayStatus,
) map[string]*models.MagmadGateway {
	gatewayEntsByKey := map[string]*models.MagmadGateway{}
	for tk, ent := range entsByTK {
		hwID := ent.PhysicalID
		var devCasted *models.GatewayDevice
		if devicesByID[hwID] != nil {
			devCasted = devicesByID[hwID].(*models.GatewayDevice)
		}
		gatewayEntsByKey[tk.Key] = (&models.MagmadGateway{}).FromBackendModels(ent, devCasted, statusesByID[hwID])
	}
	return gatewayEntsByKey
}
