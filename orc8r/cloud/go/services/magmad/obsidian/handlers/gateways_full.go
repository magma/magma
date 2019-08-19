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
	"net/url"
	"regexp"
	"strings"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"

	"github.com/labstack/echo"
	"github.com/thoas/go-funk"
)

// ListFullGatewayViews returns the full views of specified gateways in a
// network.
func ListFullGatewayViews(c echo.Context, factory view_factory.FullGatewayViewFactory) error {
	networkID, httpErr := obsidian.GetNetworkId(c)
	if httpErr != nil {
		return httpErr
	}
	gatewayIDs := getGatewayIDs(c.QueryParams())
	gatewayStates, err := getGatewayStates(networkID, gatewayIDs, factory)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, funk.Values(gatewayStates).([]*view_factory.GatewayState))
}

func getGatewayIDs(queryParams url.Values) []string {
	var gatewayIDs []string
	format2Regex := regexp.MustCompile("^gateway_ids\\[[0-9]+\\]$")
	for queryKey, values := range queryParams {
		if queryKey == "gateway_ids" && len(values) > 0 && len(values[0]) > 0 {
			// Format 1: gateway_ids=gw1,gw2,gw3
			gatewayIDs = append(gatewayIDs, strings.Split(values[0], ",")...)
		} else if format2Regex.MatchString(queryKey) {
			// Format 2: gateway_ids[0]=gw1&gateway_ids[1]=gw2&gateway_ids[2]=gw3
			gatewayIDs = append(gatewayIDs, values...)
		}
	}
	return gatewayIDs
}

func getGatewayStates(
	networkID string,
	gatewayIDs []string,
	factory view_factory.FullGatewayViewFactory,
) (map[string]*view_factory.GatewayState, error) {
	if len(gatewayIDs) > 0 {
		return factory.GetGatewayViews(networkID, gatewayIDs)
	}
	return factory.GetGatewayViewsForNetwork(networkID)
}
