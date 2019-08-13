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

	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo"
)

func ListNetworks(c echo.Context) error {
	networks, err := configurator.ListNetworkIDs()
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if networks == nil {
		networks = []string{}
	}
	return c.JSON(http.StatusOK, networks)
}

func RegisterNetwork(c echo.Context) error {
	swaggerNetwork := &models.Network{}
	err := c.Bind(&swaggerNetwork)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := swaggerNetwork.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	network := swaggerNetwork.ToConfiguratorNetwork()
	createdNetworks, err := configurator.CreateNetworks([]configurator.Network{network})
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	return c.JSON(http.StatusCreated, createdNetworks[0].ID)
}

func GetNetwork(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	network, err := configurator.LoadNetwork(networkID, true, true)
	if err == merrors.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	ret := (&models.Network{}).FromConfiguratorNetwork(network)
	return c.JSON(http.StatusOK, ret)
}

func UpdateNetwork(c echo.Context) error {
	// Bind network record from swagger
	swaggerNetwork := &models.Network{}
	err := c.Bind(&swaggerNetwork)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := swaggerNetwork.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	update := swaggerNetwork.ToUpdateCriteria()
	err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{update})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func DeleteNetwork(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	err := configurator.DeleteNetwork(networkID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}
