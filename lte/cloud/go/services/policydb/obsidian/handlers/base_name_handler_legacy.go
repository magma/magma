/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package handlers contains for the REST API handlers that converts https
// requests about policies into client_api calls to policydb
package handlers

import (
	"fmt"
	"net/http"

	"magma/lte/cloud/go/services/policydb"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian/handlers"

	"github.com/labstack/echo"
)

const (
	policyBaseNameRootPath   = policiesRootPath + "/base_names"
	policyBaseNameManagePath = policyBaseNameRootPath + "/:base_name"
)

// listBaseNameHandler returns a list of all charging rule base names in the network
func listBaseNameHandler(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	bns, err := policydb.ListBaseNames(networkID)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	ruleNames := models.RuleNames{}
	for _, ruleName := range bns {
		ruleNames = append(ruleNames, ruleName)
	}
	return c.JSON(http.StatusOK, &ruleNames)
}

// createBaseNameHandler adds a new charging rule base name to the network
func createBaseNameHandler(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	bnr := new(models.BaseNameRecord)
	if err := c.Bind(bnr); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if _, err := policydb.GetBaseName(networkID, string(bnr.Name)); err == nil {
		return handlers.HttpError(
			fmt.Errorf("Base Name '%s' already exist", bnr.Name), http.StatusConflict)
	}
	// Call policydb service
	if _, err := policydb.AddBaseName(networkID, string(bnr.Name), []string(bnr.RuleNames)); err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, string(bnr.Name))
}

// getBaseNameHandler returns the charging rule base name record associated with base_name
func getBaseNameHandler(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	baseName := getBaseNameParam(c)
	if len(baseName) == 0 {
		return baseNameHTTPErr()
	}
	// Call policydb service
	ruleNames, err := policydb.GetBaseName(networkID, baseName)
	if err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, ruleNames)
}

// updateBaseNameHandler modifies the charging rule base name
func updateBaseNameHandler(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	baseName := getBaseNameParam(c)
	if len(baseName) == 0 {
		return baseNameHTTPErr()
	}
	if _, err := policydb.GetBaseName(networkID, baseName); err != nil {
		return handlers.HttpError(
			fmt.Errorf("Base Name '%s' is not found exist", baseName), http.StatusNotFound)
	}
	ruleNames := models.RuleNames{}
	if err := c.Bind(&ruleNames); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if _, err := policydb.AddBaseName(networkID, string(baseName), []string(ruleNames)); err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

// deleteBaseNameHandler deletes the charging rule base name
func deleteBaseNameHandler(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	baseName := getBaseNameParam(c)
	if len(baseName) == 0 {
		return baseNameHTTPErr()
	}
	// Call policydb service
	if err := policydb.DeleteBaseName(networkID, baseName); err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}
	return c.NoContent(http.StatusNoContent)
}

func baseNameHTTPErr() *echo.HTTPError {
	return handlers.HttpError(
		fmt.Errorf("Invalid/Missing Base Name"),
		http.StatusBadRequest)
}

func getBaseNameParam(c echo.Context) string {
	return c.Param("base_name")
}
