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
	"sort"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	magmad_models "magma/orc8r/cloud/go/services/magmad/obsidian/models"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"

	"github.com/labstack/echo"
)

func getListGatewaysHandler(factory view_factory.FullGatewayViewFactory) func(echo.Context) error {
	return func(c echo.Context) error {
		fields := c.QueryParam("view")
		if fields == "full" {
			return ListFullGatewayViewsLegacy(c, factory)
		}
		return listGatewaysHandler(c)
	}
}

func listGatewaysHandler(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayIds, err := magmad.ListGateways(networkId)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// Return a deterministic ordering of IDs
	sort.Strings(gatewayIds)
	return c.JSON(http.StatusOK, gatewayIds)
}

func createGatewayHandler(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	swaggerRecord := &magmad_models.AccessGatewayRecord{}
	if err := c.Bind(swaggerRecord); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := swaggerRecord.Verify(); err != nil {
		return obsidian.HttpError(
			fmt.Errorf("Invalid Gateway Record, Error: %s", err),
			http.StatusBadRequest)
	}
	record, err := swaggerRecord.ToMconfig()
	if err != nil {
		return obsidian.HttpError(err, http.StatusUnsupportedMediaType)
	}

	var gatewayId string
	requestedId := c.QueryParam("requested_id")
	if len(requestedId) > 0 {
		r, _ := regexp.Compile("^[a-zA-Z_][0-9a-zA-Z_-]+$")
		if !r.MatchString(requestedId) {
			return obsidian.HttpError(
				fmt.Errorf("Gateway ID '%s' is not allowed. Gateway ID can only contain "+
					"alphanumeric characters and underscore, and should start with a letter or underscore.", requestedId),
				http.StatusBadRequest,
			)
		}
		gatewayId, err = magmad.RegisterGatewayWithId(networkId, record, requestedId)
	} else {
		gatewayId, err = magmad.RegisterGateway(networkId, record)
	}

	if err != nil {
		return obsidian.HttpError(err, http.StatusConflict)
	}

	return c.JSON(http.StatusCreated, gatewayId)
}

func getGatewayHandler(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	lid, gerr := obsidian.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}
	swaggerRecord, err := getSwaggerGWRecordFromMagmad(networkId, lid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, swaggerRecord)
}

func updateGatewayHandler(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	lid, gerr := obsidian.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	swaggerRecord := magmad_models.MutableGatewayRecord{}
	if berr := c.Bind(&swaggerRecord); berr != nil {
		return obsidian.HttpError(berr, http.StatusBadRequest)
	}
	if err := swaggerRecord.Verify(); err != nil {
		return obsidian.HttpError(
			fmt.Errorf("Invalid Gateway Record, Error: %s", err),
			http.StatusBadRequest)
	}
	record := magmadprotos.AccessGatewayRecord{}
	berr := swaggerRecord.ToMconfig(&record)
	if berr != nil {
		return obsidian.HttpError(berr, http.StatusUnsupportedMediaType)
	}
	err := magmad.UpdateGatewayRecord(networkId, lid, &record)
	if err != nil {
		return obsidian.HttpError(err, http.StatusConflict)
	}

	return c.NoContent(http.StatusOK)
}

func deleteGatewayHandler(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	lid, gerr := obsidian.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	err := magmad.RemoveGateway(networkId, lid)
	if err != nil {
		return obsidian.HttpError(err, http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}

func getSwaggerGWRecordFromMagmad(networkID, logicalID string) (*magmad_models.AccessGatewayRecord, error) {
	record, err := magmad.FindGatewayRecord(networkID, logicalID)
	if err != nil {
		return nil, obsidian.HttpError(err, http.StatusNotFound)
	}
	swaggerRecord := magmad_models.AccessGatewayRecord{}
	err = swaggerRecord.FromMconfig(record)
	if err != nil {
		return nil, obsidian.HttpError(err, http.StatusUnsupportedMediaType)
	}
	return &swaggerRecord, nil
}
