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
	"magma/orc8r/cloud/go/services/configurator"
	configuratorh "magma/orc8r/cloud/go/services/configurator/obsidian/handlers"
	magmad_models "magma/orc8r/cloud/go/services/magmad/obsidian/models"

	"github.com/labstack/echo"
)

const (
	MagmadAPIRoot    = handlers.NETWORKS_ROOT
	ListNetworks     = MagmadAPIRoot
	RegisterNetwork  = MagmadAPIRoot
	ManageNetwork    = MagmadAPIRoot + "/:network_id"
	ConfigureNetwork = ManageNetwork + "/configs"
)

func listNetworks(c echo.Context) error {
	// Check for wildcard network access
	nerr := handlers.CheckNetworkAccess(c, handlers.NETWORK_WILDCARD)
	if nerr != nil {
		return nerr
	}
	networks, err := configurator.ListNetworkIDs()
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	//magmad expects [] not null for the empty case
	if networks == nil {
		networks = []string{}
	}
	return c.JSON(http.StatusOK, networks)
}

func registerNetwork(c echo.Context) error {
	// Check for wildcard network access
	nerr := handlers.CheckNetworkAccess(c, handlers.NETWORK_WILDCARD)
	if nerr != nil {
		return nerr
	}

	// Bind network record from swagger
	record := &magmad_models.NetworkRecord{}
	err := c.Bind(&record)
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := record.ValidateModel(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	requestedID := c.QueryParam("requested_id")
	err = configuratorh.VerifyNetworkIDFormat(requestedID)
	if err != nil {
		return err
	}

	network := configurator.Network{
		Name: record.Name,
		ID:   requestedID,
		Configs: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: &magmad_models.NetworkFeatures{Features: record.Features},
		},
	}

	err = configurator.CreateNetwork(network)
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	return c.JSON(http.StatusCreated, requestedID)
}

func getNetwork(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	network, err := configurator.LoadNetwork(networkID, true, false)
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	networkFeatures := &magmad_models.NetworkFeatures{}
	features, ok := network.Configs[orc8r.NetworkFeaturesConfig]
	if ok {
		networkFeatures = features.(*magmad_models.NetworkFeatures)
	}

	record := magmad_models.NetworkRecord{
		Name:     network.Name,
		Features: networkFeatures.Features,
	}
	return c.JSON(http.StatusOK, &record)
}

func updateNetwork(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	record := &magmad_models.NetworkRecord{}
	if err := c.Bind(&record); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := record.ValidateModel(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	updateCriteria := configurator.NetworkUpdateCriteria{
		ID:      networkID,
		NewName: &record.Name,
		ConfigsToAddOrUpdate: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: &magmad_models.NetworkFeatures{Features: record.Features},
		},
	}
	configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{updateCriteria})
	return c.NoContent(http.StatusNoContent)
}

func deleteNetwork(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteNetwork(networkID)
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	return c.NoContent(http.StatusNoContent)
}
