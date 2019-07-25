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
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"

	"github.com/labstack/echo"
)

// ListFullGatewayViewsLegacy returns the full views of specified gateways in a
// network.
func ListFullGatewayViewsLegacy(c echo.Context, factory view_factory.FullGatewayViewFactory) error {
	networkID, httpErr := handlers.GetNetworkId(c)
	if httpErr != nil {
		return httpErr
	}
	gatewayIDs := getGatewayIDs(c.QueryParams())
	gatewayStates, err := getGatewayStates(networkID, gatewayIDs, factory)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	modelStates, err := view_factory.GatewayStateMapToModelList(gatewayStates)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, modelStates)
}
