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
	handlers2 "magma/orc8r/cloud/go/pluginimpl/handlers"
	models2 "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/labstack/echo"
)

// GetAgentDeviceHandlers returns GET and PUT handlers to read and update the
// device attached to the Agent.
func GetAgentDeviceHandlers(path string) []obsidian.Handler {
	return []obsidian.Handler{
		GetReadAgentDeviceHandler(path),
		GetUpdateAgentDeviceHandler(path),
	}
}

// GetReadAgentDeviceHandler returns a GET handler to read the Agent record
// of the Agent.
func GetReadAgentDeviceHandler(path string) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.GET,
		HandlerFunc: func(c echo.Context) error {
			networkID, AgentID, nerr := GetNetworkAndAgentIDs(c)
			if nerr != nil {
				return nerr
			}

			physicalID, err := configurator.GetPhysicalIDOfEntity(networkID, orc8r.MagmadGatewayType, AgentID)
			if err == merrors.ErrNotFound {
				return obsidian.HttpError(err, http.StatusNotFound)
			} else if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			device, err := device.GetDevice(networkID, orc8r.AccessGatewayRecordType, physicalID)
			if err == merrors.ErrNotFound {
				return obsidian.HttpError(err, http.StatusNotFound)
			} else if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}

			return c.JSON(http.StatusOK, device)
		},
	}
}

// GetUpdateAgentDeviceHandler returns a PUT handler to update the Agent
// record of the Agent.
func GetUpdateAgentDeviceHandler(path string) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.PUT,
		HandlerFunc: func(c echo.Context) error {
			networkID, AgentID, nerr := GetNetworkAndAgentIDs(c)
			if nerr != nil {
				return nerr
			}
			update, nerr := handlers2.GetAndValidatePayload(c, &models2.GatewayDevice{})
			if nerr != nil {
				return nerr
			}

			physicalID, err := configurator.GetPhysicalIDOfEntity(networkID, orc8r.MagmadGatewayType, AgentID)
			if err == merrors.ErrNotFound {
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
