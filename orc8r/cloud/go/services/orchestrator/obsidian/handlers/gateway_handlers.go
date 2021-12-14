/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// File gateway_handlers.go provides generic gateway handlers, with hooks for
// specific gateway types.
//
// These handlers do not support updating a gateway's ID.

package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state/wrappers"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
)

// MagmadEncompassingGateway represents a subtype of the Magmad gateway.
// The encompassing gateway wraps the base Magmad gateway with additional
// fields by creating and associating further network entities with the
// Magmad gateway.
//
// Note: since the default Magmad gateway model implements this interface as
// well, DO NOT return the base Magmad model from any of the "get additional"
// methods.
type MagmadEncompassingGateway interface {
	serde.ValidatableModel

	// GetGatewayType returns the type of the encompassing gateway.
	GetGatewayType() string

	// GetMagmadGateway returns the Magmad gateways wrapped by the subtype.
	GetMagmadGateway() *models.MagmadGateway

	// GetAdditionalLoadsOnLoad is a static method that extra TKs to load when
	// loading this gateway.
	// The entities loaded during this operation will be passed to the
	// MakeTypedGateways implementor.
	// NOTE: unlike with other "get additional" methods, DO NOT return the
	// gateway itself for this method.
	GetAdditionalLoadsOnLoad(gateway configurator.NetworkEntity) storage.TKs

	// GetAdditionalWritesOnCreate returns extra write operations to perform
	// during creation.
	// The writes are performed in the same backend transaction with creation
	// of the Magmad gateway.
	GetAdditionalWritesOnCreate() []configurator.EntityWriteOperation

	// GetAdditionalLoadsOnUpdate returns a list of additional entity keys to
	// load during an update.
	// The entities loaded during this operation will be passed to
	// GetAdditionalWritesOnUpdate.
	GetAdditionalLoadsOnUpdate() storage.TKs

	// GetAdditionalWritesOnUpdate returns extra write operations to perform
	// during an update.
	// The writes are performed in the same backend transaction with the update
	// of the Magmad gateway.
	GetAdditionalWritesOnUpdate(ctx context.Context, loadedEntities map[storage.TK]configurator.NetworkEntity) ([]configurator.EntityWriteOperation, error)
}

// MakeTypedGateways is passed the loaded ents and additional objects,
// and returns encompassing Magmad gateways keyed by gateway ID.
type MakeTypedGateways func(
	entsByTK configurator.NetworkEntitiesByTK,
	devicesByID map[string]interface{},
	statusesByID map[string]*models.GatewayStatus,
) map[string]GatewayModel

func listGatewaysHandler(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	pageSize, pageToken, err := obsidian.GetPaginationParams(c)
	if err != nil {
		return err
	}

	reqCtx := c.Request().Context()
	ents, nextPageToken, err := configurator.LoadAllEntitiesOfType(
		reqCtx,
		nid,
		orc8r.MagmadGatewayType,
		configurator.EntityLoadCriteria{
			LoadMetadata:     true,
			LoadConfig:       true,
			LoadAssocsToThis: true,
			PageSize:         uint32(pageSize),
			PageToken:        pageToken,
		},
		serdes.Entity,
	)
	if err != nil {
		return obsidian.MakeHTTPError(err, http.StatusInternalServerError)
	}

	count, err := configurator.CountEntitiesOfType(reqCtx, nid, orc8r.MagmadGatewayType)
	if err != nil {
		return obsidian.MakeHTTPError(err, http.StatusInternalServerError)
	}
	entsByTK := ents.MakeByTK()

	// For each magmad gateway, we have to load its corresponding device and
	// its reported status
	deviceIDs := make([]string, 0, len(entsByTK))
	for tk, ent := range entsByTK {
		if tk.Type == orc8r.MagmadGatewayType && ent.PhysicalID != "" {
			deviceIDs = append(deviceIDs, ent.PhysicalID)
		}
	}

	devicesByID, err := device.GetDevices(nid, orc8r.AccessGatewayRecordType, deviceIDs, serdes.Device)
	if err != nil {
		return obsidian.MakeHTTPError(errors.Wrap(err, "failed to load devices"), http.StatusInternalServerError)
	}
	statusesByID, err := wrappers.GetGatewayStatuses(reqCtx, nid, deviceIDs)
	if err != nil {
		return obsidian.MakeHTTPError(errors.Wrap(err, "failed to load statuses"), http.StatusInternalServerError)
	}
	gateways := makeGateways(entsByTK, devicesByID, statusesByID)

	err = models.PopulateRegistrationInfos(reqCtx, gateways, nid)
	if err != nil {
		return err
	}

	paginatedGateways := models.PaginatedGateways{
		Gateways:   gateways,
		PageToken:  models.PageToken(nextPageToken),
		TotalCount: int64(count),
	}
	return c.JSON(http.StatusOK, paginatedGateways)
}

func createGatewayHandler(c echo.Context) error {
	if nerr := CreateGateway(c, &models.MagmadGateway{}, serdes.Entity, serdes.Device); nerr != nil {
		return nerr
	}
	return c.NoContent(http.StatusCreated)
}

func CreateGateway(c echo.Context, model MagmadEncompassingGateway, entitySerdes, deviceSerdes serde.Registry) *echo.HTTPError {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload, nerr := GetAndValidatePayload(c, model)
	if nerr != nil {
		return nerr
	}
	subGateway := payload.(MagmadEncompassingGateway)
	mdGateway := subGateway.GetMagmadGateway()
	reqCtx := c.Request().Context()

	// Must associate to an existing tier
	tierExists, err := configurator.DoesEntityExist(reqCtx, nid, orc8r.UpgradeTierEntityType, string(mdGateway.Tier))
	if err != nil {
		return obsidian.MakeHTTPError(errors.Wrap(err, "failed to check for tier existence"), http.StatusInternalServerError)
	}
	if !tierExists {
		return echo.NewHTTPError(http.StatusBadRequest, "requested tier does not exist")
	}

	// attempt to register device
	httpErr := registerDevice(reqCtx, nid, mdGateway, entitySerdes, deviceSerdes)
	if httpErr != nil {
		return httpErr
	}

	// Create the magmad gateway, update the tier, perform additional writes
	// as necessary
	var writes []configurator.EntityWriteOperation
	writes = append(writes, mdGateway.GetAdditionalWritesOnCreate()...)
	writes = append(writes, configurator.EntityUpdateCriteria{
		Type:              orc8r.UpgradeTierEntityType,
		Key:               string(mdGateway.Tier),
		AssociationsToAdd: storage.TKs{{Type: orc8r.MagmadGatewayType, Key: string(mdGateway.ID)}},
	})

	// These type switches aren't great but it's the best I could think of
	switch payload.(type) {
	case *models.MagmadGateway:
		break
	default:
		writes = append(writes, subGateway.GetAdditionalWritesOnCreate()...)
	}

	if err = configurator.WriteEntities(reqCtx, nid, writes, entitySerdes); err != nil {
		return obsidian.MakeHTTPError(errors.Wrap(err, "error creating gateway"), http.StatusInternalServerError)
	}
	return nil
}

// registerDevice, if gateway.Device exists, performs the following actions depending on device registration state:
// If the device is already registered, throw an error if it's already
// assigned to an entity
// If the device exists but is unassigned, update it to the payload
// If the device doesn't exist, create it and move on
func registerDevice(ctx context.Context, networkID string, gateway *models.MagmadGateway, entitySerdes, deviceSerdes serde.Registry) *echo.HTTPError {
	if gateway.Device == nil {
		return nil
	}

	deviceID := gateway.Device.HardwareID
	_, err := device.GetDevice(ctx, networkID, orc8r.AccessGatewayRecordType, deviceID, deviceSerdes)
	switch {
	case err == merrors.ErrNotFound:
		err = device.RegisterDevice(ctx, networkID, orc8r.AccessGatewayRecordType, deviceID, gateway.Device, deviceSerdes)
		if err != nil {
			return obsidian.MakeHTTPError(errors.Wrap(err, "failed to register physical device"), http.StatusInternalServerError)
		}
	case err != nil:
		return obsidian.MakeHTTPError(errors.Wrap(err, "failed to check if physical device is already registered"), http.StatusConflict)
	default: // err == nil
		assignedEnt, err := configurator.LoadEntityForPhysicalID(ctx, deviceID, configurator.EntityLoadCriteria{}, entitySerdes)
		switch {
		case err == nil:
			return obsidian.MakeHTTPError(errors.Errorf("device %s is already mapped to gateway %s", deviceID, assignedEnt.Key), http.StatusBadRequest)
		case err != merrors.ErrNotFound:
			return obsidian.MakeHTTPError(errors.Wrap(err, "failed to check for existing device assignment"), http.StatusInternalServerError)
		}

		if err := device.UpdateDevice(ctx, networkID, orc8r.AccessGatewayRecordType, deviceID, gateway.Device, deviceSerdes); err != nil {
			return obsidian.MakeHTTPError(errors.Wrap(err, "failed to update device record"), http.StatusInternalServerError)
		}
	}
	return nil
}

func getGatewayHandler(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}
	ret, nerr := LoadMagmadGateway(c.Request().Context(), nid, gid)
	if nerr != nil {
		return nerr
	}
	return c.JSON(http.StatusOK, ret)
}

func LoadMagmadGateway(ctx context.Context, networkID string, gatewayID string) (*models.MagmadGateway, *echo.HTTPError) {
	ent, err := configurator.LoadEntity(
		ctx,
		networkID, orc8r.MagmadGatewayType, gatewayID,
		configurator.EntityLoadCriteria{
			LoadMetadata:       true,
			LoadConfig:         true,
			LoadAssocsToThis:   true,
			LoadAssocsFromThis: false,
		},
		serdes.Entity,
	)
	if err == merrors.ErrNotFound {
		return nil, echo.ErrNotFound
	}
	if err != nil {
		return nil, obsidian.MakeHTTPError(err, http.StatusInternalServerError)
	}

	dev, err := device.GetDevice(ctx, networkID, orc8r.AccessGatewayRecordType, ent.PhysicalID, serdes.Device)
	if err != nil && err != merrors.ErrNotFound {
		return nil, obsidian.MakeHTTPError(err, http.StatusInternalServerError)
	}
	status, err := wrappers.GetGatewayStatus(ctx, networkID, ent.PhysicalID)
	if err != nil && err != merrors.ErrNotFound {
		return nil, obsidian.MakeHTTPError(err, http.StatusInternalServerError)
	}

	// If the gateway/network is malformed, we could get no corresponding
	// device for the gateway
	var devCasted *models.GatewayDevice
	if dev != nil {
		devCasted = dev.(*models.GatewayDevice)
	}

	return (&models.MagmadGateway{}).FromBackendModels(ent, devCasted, status), nil
}

func updateGatewayHandler(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	if nerr = UpdateGateway(c, nid, gid, &models.MagmadGateway{}, serdes.Entity, serdes.Device); nerr != nil {
		return nerr
	}
	return c.NoContent(http.StatusNoContent)
}

func UpdateGateway(c echo.Context, nid string, gid string, model MagmadEncompassingGateway, entitySerdes, deviceSerdes serde.Registry) *echo.HTTPError {
	payload, nerr := GetAndValidatePayload(c, model)
	if nerr != nil {
		return nerr
	}
	subGateway := payload.(MagmadEncompassingGateway)
	mdGateway := subGateway.GetMagmadGateway()

	if gid != string(mdGateway.ID) {
		err := fmt.Errorf("gateway ID cannot be updated: gateway ID from parameter (%s) and payload (%s) must match", gid, mdGateway.ID)
		return obsidian.MakeHTTPError(err, http.StatusBadRequest)
	}

	var entsToLoad storage.TKs
	entsToLoad = append(entsToLoad, mdGateway.GetAdditionalLoadsOnUpdate()...)
	switch payload.(type) {
	case *models.MagmadGateway:
		break
	default:
		entsToLoad = append(entsToLoad, subGateway.GetAdditionalLoadsOnUpdate()...)
	}

	reqCtx := c.Request().Context()
	loadedEnts, _, err := configurator.LoadEntities(
		reqCtx,
		nid,
		nil, nil, nil,
		entsToLoad,
		configurator.FullEntityLoadCriteria(),
		entitySerdes,
	)
	if err != nil {
		return obsidian.MakeHTTPError(errors.Wrap(err, "failed to load gateway before update"), http.StatusInternalServerError)
	}

	writes, nerr := getUpdateWrites(reqCtx, subGateway, loadedEnts)
	if nerr != nil {
		return nerr
	}

	err = configurator.WriteEntities(reqCtx, nid, writes, entitySerdes)
	if err != nil {
		return obsidian.MakeHTTPError(err, http.StatusInternalServerError)
	}

	// Device info is cheap to update, so just do it all the time if
	// configurator write was successful
	err = device.UpdateDevice(reqCtx, nid, orc8r.AccessGatewayRecordType, mdGateway.Device.HardwareID, mdGateway.Device, deviceSerdes)
	if err != nil {
		return obsidian.MakeHTTPError(errors.Wrap(err, "failed to update device info"), http.StatusInternalServerError)
	}

	return nil
}

func getUpdateWrites(ctx context.Context, payload MagmadEncompassingGateway, loadedEnts configurator.NetworkEntities) ([]configurator.EntityWriteOperation, *echo.HTTPError) {
	var writes []configurator.EntityWriteOperation
	loadedEntsByID := loadedEnts.MakeByTK()

	mdGwWrites, err := payload.GetMagmadGateway().GetAdditionalWritesOnUpdate(ctx, loadedEntsByID)
	switch {
	case err == merrors.ErrNotFound:
		return writes, echo.ErrNotFound
	case err != nil:
		return writes, obsidian.MakeHTTPError(errors.Wrap(err, "failed to get update operations from magmad model"), http.StatusInternalServerError)
	}

	// Short circuit if this is the magmad gateway
	switch payload.(type) {
	case *models.MagmadGateway:
		return mdGwWrites, nil
	}

	payloadWrites, err := payload.GetAdditionalWritesOnUpdate(context.Background(), loadedEntsByID)
	switch {
	case err == merrors.ErrNotFound:
		return writes, echo.ErrNotFound
	case err != nil:
		return writes, obsidian.MakeHTTPError(errors.Wrap(err, "failed to get update operations from payload model"), http.StatusInternalServerError)
	}

	writes = append(writes, mdGwWrites...)
	writes = append(writes, payloadWrites...)
	return writes, nil
}

func deleteGatewayHandler(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}
	err := DeleteMagmadGateway(c.Request().Context(), nid, gid, nil)
	if err != nil {
		return makeErr(err)
	}
	return c.NoContent(http.StatusNoContent)
}

func DeleteMagmadGateway(ctx context.Context, networkID, gatewayID string, additionalDeletes storage.TKs) error {
	mdGw, err := configurator.LoadEntity(ctx, networkID, orc8r.MagmadGatewayType, gatewayID, configurator.EntityLoadCriteria{}, serdes.Entity)
	if err != nil {
		return err
	}

	var deletes storage.TKs
	deletes = append(deletes, storage.TK{Type: orc8r.MagmadGatewayType, Key: gatewayID})
	deletes = append(deletes, additionalDeletes...)

	err = configurator.DeleteEntities(ctx, networkID, deletes)
	if err != nil {
		return obsidian.MakeHTTPError(errors.Wrap(err, "error deleting gateway"), http.StatusInternalServerError)
	}

	// Now we delete the associated device. Even though we error out
	// request if this fails, failing on this specific step is non-
	// blocking because gateway registration handles the case where a
	// device already exists and is unassigned.
	if mdGw.PhysicalID != "" {
		err = device.DeleteDevice(ctx, networkID, orc8r.AccessGatewayRecordType, mdGw.PhysicalID)
		if err != nil {
			return obsidian.MakeHTTPError(errors.Wrap(err, "failed to delete device for gateway. no further action is required"), http.StatusInternalServerError)
		}
	}

	return nil
}

func GetStateHandler(c echo.Context) error {
	networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	reqCtx := c.Request().Context()
	physicalID, err := configurator.GetPhysicalIDOfEntity(reqCtx, networkID, orc8r.MagmadGatewayType, gatewayID)
	if err == merrors.ErrNotFound {
		return obsidian.MakeHTTPError(err, http.StatusNotFound)
	} else if err != nil {
		return obsidian.MakeHTTPError(err, http.StatusInternalServerError)
	}

	st, err := wrappers.GetGatewayStatus(reqCtx, networkID, physicalID)
	if err == merrors.ErrNotFound {
		return obsidian.MakeHTTPError(err, http.StatusNotFound)
	} else if err != nil {
		return obsidian.MakeHTTPError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, st)
}

func makeGateways(
	entsByTK configurator.NetworkEntitiesByTK,
	devicesByID map[string]interface{},
	statusesByID map[string]*models.GatewayStatus,
) map[string]*models.MagmadGateway {
	gatewayEntsByKey := map[string]*models.MagmadGateway{}
	for tk, ent := range entsByTK.Filter(orc8r.MagmadGatewayType) {
		hwID := ent.PhysicalID
		var devCasted *models.GatewayDevice
		if devicesByID[hwID] != nil {
			devCasted = devicesByID[hwID].(*models.GatewayDevice)
		}
		gatewayEntsByKey[tk.Key] = (&models.MagmadGateway{}).FromBackendModels(ent, devCasted, statusesByID[hwID])
	}
	return gatewayEntsByKey
}

func makeErr(err error) *echo.HTTPError {
	if err == merrors.ErrNotFound {
		return echo.ErrNotFound
	}
	return obsidian.MakeHTTPError(err, http.StatusInternalServerError)
}
