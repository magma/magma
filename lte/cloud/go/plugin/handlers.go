/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin

import (
	"fmt"
	"net/http"
	"sort"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/plugin/models"
	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	LteNetworks       = "ltenetworks"
	ListNetworksPath  = obsidian.V1Root + LteNetworks
	ManageNetworkPath = ListNetworksPath + "/:network_id"
)

func getNetworkHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		{
			Path:        ListNetworksPath,
			Methods:     obsidian.GET,
			HandlerFunc: ListNetworks,
		},
		{
			Path:        ListNetworksPath,
			Methods:     obsidian.POST,
			HandlerFunc: CreateNetwork,
		},
		{
			Path:        ManageNetworkPath,
			Methods:     obsidian.GET,
			HandlerFunc: GetNetwork,
		},
		{
			Path:        ManageNetworkPath,
			Methods:     obsidian.DELETE,
			HandlerFunc: DeleteNetwork,
		},
	}
}

func ListNetworks(c echo.Context) error {
	ids, err := configurator.ListNetworksOfType(lte.LteNetworkType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	sort.Strings(ids)
	return c.JSON(http.StatusOK, ids)
}

func CreateNetwork(c echo.Context) error {
	payload := &models.LteNetwork{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	err := configurator.CreateNetwork(payload.ToConfiguratorNetwork())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func GetNetwork(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	network, err := configurator.LoadNetwork(nid, true, true)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if network.Type != lte.LteNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not an LTE network", nid))
	}

	ret := (&models.LteNetwork{}).FromConfiguratorNetwork(network)
	return c.JSON(http.StatusOK, ret)
}

func DeleteNetwork(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	// check that this is actually an LTE network
	network, err := configurator.LoadNetwork(nid, false, false)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load network to check type"), http.StatusInternalServerError)
	}
	if network.Type != lte.LteNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not an LTE network", nid))
	}

	err = configurator.DeleteNetwork(nid)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}
