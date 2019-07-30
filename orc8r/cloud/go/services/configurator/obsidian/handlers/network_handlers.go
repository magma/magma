/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_models "magma/orc8r/cloud/go/services/configurator/obsidian/models"

	"github.com/labstack/echo"
)

const (
	ConfiguratorAPIRoot      = obsidian.RestRoot + obsidian.UrlSep + "configurator"
	ConfiguratorNetworksRoot = ConfiguratorAPIRoot + obsidian.UrlSep + "networks"
	ListNetworks             = ConfiguratorNetworksRoot
	RegisterNetwork          = ConfiguratorNetworksRoot
	ManageNetwork            = ConfiguratorNetworksRoot + "/:network_id"
)

func listNetworks(c echo.Context) error {
	// Check for wildcard network access
	nerr := obsidian.CheckWildcardNetworkAccess(c)
	if nerr != nil {
		return nerr
	}
	networks, err := configurator.ListNetworkIDs()
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, networks)
}

func registerNetwork(c echo.Context) error {
	// Check for wildcard network access
	nerr := obsidian.CheckWildcardNetworkAccess(c)
	if nerr != nil {
		return nerr
	}
	// Bind network record from swagger
	swaggerNetwork := &configurator_models.NetworkRecord{}
	err := c.Bind(&swaggerNetwork)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	requestedID := c.QueryParam("requested_id")
	err = VerifyNetworkIDFormat(requestedID)
	if err != nil {
		return err
	}

	network := configurator.Network{
		ID:          requestedID,
		Name:        swaggerNetwork.Name,
		Description: swaggerNetwork.Description,
	}
	createdNetworks, err := configurator.CreateNetworks([]configurator.Network{network})
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	return c.JSON(http.StatusCreated, createdNetworks[0].ID)
}

func getNetwork(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	networks, _, err := configurator.LoadNetworks([]string{networkID}, true, false)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if len(networks) == 0 {
		return obsidian.HttpError(fmt.Errorf("Network ID %s not found", networkID), http.StatusBadRequest)
	}
	network := networks[0]

	swaggerRecord := &configurator_models.NetworkRecord{
		Name:        network.Name,
		Description: network.Description,
	}
	return c.JSON(http.StatusOK, &swaggerRecord)
}

func updateNetwork(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	// Bind network record from swagger
	swaggerNetwork := &configurator_models.NetworkRecord{}
	err := c.Bind(&swaggerNetwork)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	updateCriteria := configurator.NetworkUpdateCriteria{
		ID:             networkID,
		NewName:        &swaggerNetwork.Name,
		NewDescription: &swaggerNetwork.Description,
	}
	err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{updateCriteria})
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("Network:%s updated", networkID))
}

func deleteNetwork(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteNetworks([]string{networkID})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func VerifyNetworkIDFormat(requestedID string) error {
	if len(requestedID) > 0 {
		r, _ := regexp.Compile("^[a-z_][0-9a-z_]+$")
		if !r.MatchString(requestedID) {
			return obsidian.HttpError(
				fmt.Errorf("Network ID '%s' is not allowed. Network ID can only contain "+
					"lowercase alphanumeric characters and underscore, and should start with a letter or underscore.", requestedID),
				http.StatusBadRequest,
			)
		}
	}
	return nil
}
