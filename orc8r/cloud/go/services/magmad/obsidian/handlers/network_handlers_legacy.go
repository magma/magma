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
	"strings"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_models "magma/orc8r/cloud/go/services/magmad/obsidian/models"

	"github.com/labstack/echo"
)

func listNetworksHandler(c echo.Context) error {
	// Check for wildcard network access
	nerr := obsidian.CheckWildcardNetworkAccess(c)
	if nerr != nil {
		return nerr
	}
	networks, err := magmad.ListNetworks()
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, networks)
}

func registerNetworkHandler(c echo.Context) error {
	// Check for wildcard network access
	nerr := obsidian.CheckWildcardNetworkAccess(c)
	if nerr != nil {
		return nerr
	}

	// Bind network record from swagger
	swaggerRecord := &magmad_models.NetworkRecord{}
	err := c.Bind(&swaggerRecord)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := swaggerRecord.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	magmadRecord := swaggerRecord.ToProto()

	var networkId string
	requestedId := c.QueryParam("requested_id")
	err = VerifyNetworkIDFormat(requestedId)
	if err != nil {
		return err
	}
	networkId, err = magmad.RegisterNetwork(magmadRecord, requestedId)
	if err != nil {
		return obsidian.HttpError(err, http.StatusConflict)
	}

	return c.JSON(http.StatusCreated, networkId)
}

func getNetworkHandler(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	record, err := magmad.GetNetwork(networkId)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	swaggerRecord := magmad_models.NetworkRecordFromProto(record)
	return c.JSON(http.StatusOK, &swaggerRecord)
}

func updateNetworkHandler(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	swaggerRecord := &magmad_models.NetworkRecord{}
	if err := c.Bind(&swaggerRecord); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := swaggerRecord.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	err := magmad.UpdateNetwork(networkId, swaggerRecord.ToProto())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteNetworkHandler(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
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
		return obsidian.HttpError(err, status)
	}
	return c.NoContent(http.StatusNoContent)
}
