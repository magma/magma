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

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"

	"github.com/labstack/echo"
)

// PartialGatewayModel describe models that represents a portion of network
// entity that can be read and updated.
type PartialGatewayModel interface {
	ValidatableModel
	// GetFromEntity grabs the desired model from the configurator network entity.
	// Returns nil if it is not there.
	GetFromEntity(entity configurator.NetworkEntity) interface{}
	// ToUpdateCriteria takes in the existing network entity and applies the
	// change from the model to create a list of EntityUpdateCriteria.
	ToUpdateCriteria(entity configurator.NetworkEntity) ([]configurator.EntityUpdateCriteria, error)
}

// GetPartialGatewayHandlers returns both GET and PUT handlers for modifying the portion of a
// network entity specified by the model.
// - path : the url at which the handler will be registered.
// - model: the input and output of the handler and it also provides GetFromEntity
//   and ToUpdateCriteria to go between the configurator model.
func GetPartialGatewayHandlers(path string, model PartialGatewayModel) []obsidian.Handler {
	return []obsidian.Handler{
		GetPartialReadGatewayHandler(path, model),
		GetPartialUpdateGatewayHandler(path, model),
	}
}

// GetPartialReadGatewayHandler returns a GET obsidian handler at the specified path.
// This function loads a portion of the gateway specified by the model's GetFromEntity function.
// Example:
//      (m *MagmadGatewayConfigs) GetFromEntity(ent configurator.NetworkEntity) interface{} {
// 			return entity.Configs
// 		}
// 		getMagmadConfigsHandler := handlers.GetPartialReadGatewayHandler(URL, &models.MagmadGatewayConfigs{})
//
//      would return a GET handler that can read the magmad gateway config of a gw with the specified ID.
func GetPartialReadGatewayHandler(path string, model PartialGatewayModel) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.GET,
		HandlerFunc: func(c echo.Context) error {
			networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
			if nerr != nil {
				return nerr
			}

			entity, err := configurator.LoadEntity(networkID, orc8r.MagmadGatewayType, gatewayID, configurator.EntityLoadCriteria{LoadConfig: true, LoadMetadata: true, LoadAssocsFromThis: true})
			if err == errors.ErrNotFound {
				return obsidian.HttpError(err, http.StatusNotFound)
			} else if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}

			got := model.GetFromEntity(entity)
			if got == nil {
				return obsidian.HttpError(errors.ErrNotFound, http.StatusNotFound)
			}
			return c.JSON(http.StatusOK, got)
		},
	}
}

// GetPartialUpdateGatewayHandler returns a PUT obsidian handler at the specified path.
// This function updates a portion of the network entity specified by the model's ToUpdateCriteria function.
// Example:
//      (m *MagmadGatewayConfigs) ToUpdateCriteria(ent configurator.NetworkEntity) (configurator.EntityUpdateCriteria, error) {
// 			return configurator.EntityUpdateCriteria{
// 				Key: ent.Key,
//				Type: ent.Type,
//				NewConfig: m,
//          }
// 		}
// 		updateMagmadConfigsHandler := handlers.GetPartialUpdateGatewayHandler(URL, &models.MagmadGatewayConfigs{})
//
//      would return a PUT handler that updates the magmad gateway config of a gw with the specified ID.
func GetPartialUpdateGatewayHandler(path string, model PartialGatewayModel) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.PUT,
		HandlerFunc: func(c echo.Context) error {
			networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
			if nerr != nil {
				return nerr
			}

			requestedUpdate, nerr := GetAndValidatePayload(c, model)
			if nerr != nil {
				return nerr
			}

			entity, err := configurator.LoadEntity(networkID, orc8r.MagmadGatewayType, gatewayID, configurator.EntityLoadCriteria{LoadConfig: true, LoadMetadata: true, LoadAssocsFromThis: true})
			if err == errors.ErrNotFound {
				return obsidian.HttpError(err, http.StatusNotFound)
			} else if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}

			updates, err := requestedUpdate.(PartialGatewayModel).ToUpdateCriteria(entity)
			if err != nil {
				return obsidian.HttpError(err, http.StatusBadRequest)
			}
			_, err = configurator.UpdateEntities(networkID, updates)
			if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			return c.NoContent(http.StatusNoContent)
		},
	}
}

// GetGatewayDeviceHandlers returns GET and PUT handlers to read and update the
// device attached to the gateway.
func GetGatewayDeviceHandlers(path string) []obsidian.Handler {
	return []obsidian.Handler{
		GetReadGatewayDeviceHandler(path),
		GetUpdateGatewayDeviceHandler(path),
	}
}

// GetReadGatewayDeviceHandler returns a GET handler to read the gateway record
// of the gateway.
func GetReadGatewayDeviceHandler(path string) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.GET,
		HandlerFunc: func(c echo.Context) error {
			networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
			if nerr != nil {
				return nerr
			}

			physicalID, err := configurator.GetPhysicalIDOfEntity(networkID, orc8r.MagmadGatewayType, gatewayID)
			if err == errors.ErrNotFound {
				return obsidian.HttpError(err, http.StatusNotFound)
			} else if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			device, err := device.GetDevice(networkID, orc8r.AccessGatewayRecordType, physicalID)
			if err == errors.ErrNotFound {
				return obsidian.HttpError(err, http.StatusNotFound)
			} else if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}

			return c.JSON(http.StatusOK, device)
		},
	}
}

// GetUpdateGatewayDeviceHandler returns a PUT handler to update the gateway
// record of the gateway.
func GetUpdateGatewayDeviceHandler(path string) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.PUT,
		HandlerFunc: func(c echo.Context) error {
			networkID, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
			if nerr != nil {
				return nerr
			}
			update, nerr := GetAndValidatePayload(c, &models.GatewayDevice{})
			if nerr != nil {
				return nerr
			}

			physicalID, err := configurator.GetPhysicalIDOfEntity(networkID, orc8r.MagmadGatewayType, gatewayID)
			if err == errors.ErrNotFound {
				return obsidian.HttpError(err, http.StatusNotFound)
			} else if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			err = device.UpdateDevice(networkID, orc8r.AccessGatewayRecordType, physicalID, update)
			if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			return c.NoContent(http.StatusNoContent)
		},
	}
}
