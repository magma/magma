/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/checkind"
	checkind_models "magma/orc8r/cloud/go/services/checkind/obsidian/models"
)

const AgStatusUrl = handlers.NETWORKS_ROOT + "/:network_id/gateways/:logical_ag_id/status"

// GetObsidianHandlers returns all handlers for checkind
func GetObsidianHandlers() []handlers.Handler {
	return []handlers.Handler{
		{
			Path:    AgStatusUrl,
			Methods: handlers.GET,
			HandlerFunc: func(c echo.Context) error {
				network_id, nerr := handlers.GetNetworkId(c)
				if nerr != nil {
					return nerr
				}

				lid := c.Param("logical_ag_id")
				status, err := checkind.GetStatus(network_id, lid)
				if err != nil {
					return handlers.HttpError(err, http.StatusNotFound)
				}

				gwStatusModel := checkind_models.GatewayStatus{}
				err = gwStatusModel.FromMconfig(status)
				if err != nil {
					return handlers.HttpError(err, http.StatusInternalServerError)
				}
				return c.JSON(http.StatusOK, &gwStatusModel)
			},
		},
	}
}
