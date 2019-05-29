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

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_models "magma/orc8r/cloud/go/services/configurator/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator/protos"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/labstack/echo"
)

const (
	ConfiguratorAPIRoot      = handlers.REST_ROOT + handlers.URL_SEP + "configurator"
	ConfiguratorNetworksRoot = ConfiguratorAPIRoot + handlers.URL_SEP + "networks"
	ListNetworks             = ConfiguratorNetworksRoot
	RegisterNetwork          = ConfiguratorNetworksRoot
	ManageNetwork            = ConfiguratorNetworksRoot + "/:network_id"
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
	return c.JSON(http.StatusOK, networks)
}

func registerNetwork(c echo.Context) error {
	// Check for wildcard network access
	nerr := handlers.CheckNetworkAccess(c, handlers.NETWORK_WILDCARD)
	if nerr != nil {
		return nerr
	}
	// Bind network record from swagger
	swaggerNetwork := &configurator_models.NetworkRecord{}
	err := c.Bind(&swaggerNetwork)
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	requestedID := c.QueryParam("requested_id")
	err = VerifyNetworkIDFormat(requestedID)
	if err != nil {
		return err
	}

	protoNetwork := &protos.Network{
		Id:          requestedID,
		Name:        swaggerNetwork.Name,
		Description: swaggerNetwork.Description,
	}
	createdNetworks, err := configurator.CreateNetworks([]*protos.Network{protoNetwork})
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	return c.JSON(http.StatusCreated, createdNetworks[0].Id)
}

func getNetwork(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	networks, _, err := configurator.LoadNetworks([]string{networkID}, true, false)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	if len(networks) == 0 {
		return handlers.HttpError(fmt.Errorf("Network ID %s not found", networkID), http.StatusBadRequest)
	}
	network := networks[networkID]

	swaggerRecord := &configurator_models.NetworkRecord{
		Name:        network.Name,
		Description: network.Description,
	}
	return c.JSON(http.StatusOK, &swaggerRecord)
}

func updateNetwork(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	// Bind network record from swagger
	swaggerNetwork := &configurator_models.NetworkRecord{}
	err := c.Bind(&swaggerNetwork)
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	updateCriteria := &protos.NetworkUpdateCriteria{
		Id:             networkID,
		NewName:        inputStrToStrWrapper(swaggerNetwork.Name),
		NewDescription: inputStrToStrWrapper(swaggerNetwork.Description),
	}
	err = configurator.UpdateNetworks([]*protos.NetworkUpdateCriteria{updateCriteria})
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("Network:%s updated", networkID))
}

func deleteNetwork(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteNetworks([]string{networkID})
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func inputStrToStrWrapper(in string) *wrappers.StringValue {
	return &wrappers.StringValue{Value: in}
}
func VerifyNetworkIDFormat(requestedID string) error {
	if len(requestedID) > 0 {
		r, _ := regexp.Compile("^[a-z_][0-9a-z_]+$")
		if !r.MatchString(requestedID) {
			return handlers.HttpError(
				fmt.Errorf("Network ID '%s' is not allowed. Network ID can only contain "+
					"lowercase alphanumeric characters and underscore, and should start with a letter or underscore.", requestedID),
				http.StatusBadRequest,
			)
		}
	}
	return nil
}
