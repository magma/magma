/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"net/http"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/state"

	"github.com/labstack/echo"
)

const AgStatusURL = obsidian.NetworksRoot + "/:network_id/gateways/:gw_hardware_id/gateway_status"

// GetObsidianHandlers returns all handlers for state
func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		{
			Path:        AgStatusURL,
			Methods:     obsidian.GET,
			HandlerFunc: AGStatusByDeviceIDHandler,
		},
	}
}

func AGStatusByDeviceIDHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	deviceID := c.Param("gw_hardware_id")
	gwStatusModel, err := state.GetGatewayStatus(networkID, deviceID)
	if err == errors.ErrNotFound || gwStatusModel == nil {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &gwStatusModel)
}
