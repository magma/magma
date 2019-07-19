/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"net/http"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/orc8r"
	checkind_models "magma/orc8r/cloud/go/services/checkind/obsidian/models"
	stateservice "magma/orc8r/cloud/go/services/state"

	"github.com/labstack/echo"
)

const AgStatusUrl = handlers.NETWORKS_ROOT + "/:network_id/gateways/:device_id/gateway_status"

// GetObsidianHandlers returns all handlers for state
func GetObsidianHandlers() []handlers.Handler {
	return []handlers.Handler{
		{
			Path:        AgStatusUrl,
			Methods:     handlers.GET,
			HandlerFunc: AGStatusByDeviceIDHandler,
		},
	}
}

func AGStatusByDeviceIDHandler(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	deviceID := c.Param("device_id")
	gwStatusModel, err := GetGWStatus(networkID, deviceID)
	if err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, &gwStatusModel)
}

func GetGWStatus(networkID string, deviceID string) (*checkind_models.GatewayStatus, error) {
	state, err := stateservice.GetState(networkID, orc8r.GatewayStateType, deviceID)
	if err != nil {
		return nil, err
	}
	gwStatus := state.ReportedState.(checkind_models.GatewayStatus)
	gwStatus.CheckinTime = state.Time
	gwStatus.CertExpirationTime = state.CertExpirationTime
	// Use the hardware ID from the middleware
	gwStatus.HardwareID = state.ReporterID
	// Populate deprecated fields to support API backwards compatibility
	// TODO: Remove this and related tests when deprecated fields are no longer used
	gwStatus.FillDeprecatedFields()
	return &gwStatus, nil
}
