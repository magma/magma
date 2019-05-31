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
	"strings"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorh "magma/orc8r/cloud/go/services/configurator/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/magmad"
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
	networks, err := magmad.ListNetworks()
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
	swaggerRecord := &magmad_models.NetworkRecord{}
	err := c.Bind(&swaggerRecord)
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := swaggerRecord.ValidateModel(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	magmadRecord := swaggerRecord.ToProto()

	var networkId string
	requestedId := c.QueryParam("requested_id")
	err = configuratorh.VerifyNetworkIDFormat(requestedId)
	if err != nil {
		return err
	}
	networkId, err = magmad.RegisterNetwork(magmadRecord, requestedId)
	if err != nil {
		return handlers.HttpError(err, http.StatusConflict)
	}

	err = multiplexCreateNetworkIntoConfigurator(networkId, swaggerRecord)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("Failed to multiplex create into configurator: %v", err), http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, networkId)
}

func multiplexCreateNetworkIntoConfigurator(requestedID string, swaggerRecord *magmad_models.NetworkRecord) error {
	network := &protos.Network{
		Name: swaggerRecord.Name,
		Id:   requestedID,
	}
	_, err := configurator.CreateNetworks([]*protos.Network{network})
	if err != nil {
		return err
	}
	return nil
}

func getNetwork(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	record, err := magmad.GetNetwork(networkId)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	swaggerRecord := magmad_models.NetworkRecordFromProto(record)
	return c.JSON(http.StatusOK, &swaggerRecord)
}

func updateNetwork(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	swaggerRecord := &magmad_models.NetworkRecord{}
	if err := c.Bind(&swaggerRecord); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := swaggerRecord.ValidateModel(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	err := multiplexUpdateNetworkIntoConfigurator(networkId, swaggerRecord)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("Failed to multiplex update into configurator: %v", err), http.StatusInternalServerError)
	}

	return magmad.UpdateNetwork(networkId, swaggerRecord.ToProto())
}

func multiplexUpdateNetworkIntoConfigurator(networkID string, record *magmad_models.NetworkRecord) error {
	exists, err := configurator.NetworkExists(networkID)
	if err != nil {
		return err
	}
	if !exists {
		return multiplexCreateNetworkIntoConfigurator(networkID, record)
	}
	updateCriteria := &protos.NetworkUpdateCriteria{
		Id:      networkID,
		NewName: protos.GetStringWrapper(&record.Name),
	}
	return configurator.UpdateNetworks([]*protos.NetworkUpdateCriteria{updateCriteria})
}

func deleteNetwork(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	force := c.QueryParam("mode")
	var err error
	if strings.ToUpper(force) == "FORCE" {
		err = magmad.ForceRemoveNetwork(networkId)
	} else {
		err = magmad.RemoveNetwork(networkId)
	}

	if err != nil {
		status := http.StatusInternalServerError
		// TODO: conversion based on grpc code
		return handlers.HttpError(err, status)
	}

	// multiplex delete network into configurator
	err = configurator.DeleteNetworks([]string{networkId})
	if err != nil {
		return handlers.HttpError(fmt.Errorf("Failed to multiplex delete into configurator: %v", err), http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
