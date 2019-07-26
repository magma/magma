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
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/services/state"

	"github.com/labstack/echo"

	"magma/orc8r/cloud/go/obsidian/handlers"
)

const AgStatusUrl = handlers.NETWORKS_ROOT + "/:network_id/gateways/:logical_ag_id/status"

// GetObsidianHandlers returns all handlers for checkind
func GetObsidianHandlers() []handlers.Handler {
	return []handlers.Handler{
		{
			Path:    AgStatusUrl,
			Methods: handlers.GET,
			HandlerFunc: func(c echo.Context) error {
				networkID, nerr := handlers.GetNetworkId(c)
				if nerr != nil {
					return nerr
				}

				lid := c.Param("logical_ag_id")

				gwRecord, err := magmad.FindGatewayRecord(networkID, lid)
				if err != nil {
					return handlers.HttpError(err, http.StatusNotFound)
				}
				hwid := gwRecord.HwId.Id
				gwStatus, err := state.GetGatewayStatus(networkID, hwid)
				if err == errors.ErrNotFound || gwStatus == nil {
					return c.NoContent(http.StatusNotFound)
				}
				if err != nil {
					return handlers.HttpError(err, http.StatusInternalServerError)
				}
				return c.JSON(http.StatusOK, &gwStatus)
			},
			MigratedHandlerFunc: func(c echo.Context) error {
				networkID, nerr := handlers.GetNetworkId(c)
				if nerr != nil {
					return nerr
				}

				gwLogicalID := c.Param("logical_ag_id")
				gwPhysicalID, err := configurator.GetPhysicalIDOfEntity(networkID, orc8r.MagmadGatewayType, gwLogicalID)
				if err != nil {
					return handlers.HttpError(err, http.StatusNotFound)
				}
				gwStatus, err := state.GetGatewayStatus(networkID, gwPhysicalID)
				if err == errors.ErrNotFound || gwStatus == nil {
					return c.NoContent(http.StatusNotFound)
				}
				if err != nil {
					return handlers.HttpError(err, http.StatusInternalServerError)
				}
				return c.JSON(http.StatusOK, &gwStatus)
			},
		},
	}
}
