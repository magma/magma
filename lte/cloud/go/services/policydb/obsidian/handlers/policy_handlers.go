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
	"sort"
	"strings"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	baseNameParam = "base_name"
	ruleIDParam   = "rule_id"
)

// Base names

func ListBaseNames(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	view := c.QueryParam("view")
	if strings.ToLower(view) == "full" {
		baseNames, err := configurator.LoadAllEntitiesInNetwork(networkID, lte.BaseNameEntityType, configurator.EntityLoadCriteria{LoadAssocsFromThis: true})
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		ret := map[string]*models.BaseNameRecord{}
		for _, bnEnt := range baseNames {
			ret[bnEnt.Key] = (&models.BaseNameRecord{}).FromEntity(bnEnt)
		}
		return c.JSON(http.StatusOK, ret)
	} else {
		names, err := configurator.ListEntityKeys(networkID, lte.BaseNameEntityType)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		sort.Strings(names)
		return c.JSON(http.StatusOK, names)
	}
}

func CreateBaseName(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	bnr := new(models.BaseNameRecord)
	if err := c.Bind(bnr); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	_, err := configurator.CreateEntity(networkID, bnr.ToEntity())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, string(bnr.Name))
}

func GetBaseName(c echo.Context) error {
	networkID, baseName, nerr := getNetworkIDAndBaseName(c)
	if nerr != nil {
		return nerr
	}

	ret, err := configurator.LoadEntity(
		networkID,
		lte.BaseNameEntityType,
		baseName,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
	)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, (&models.BaseNameRecord{}).FromEntity(ret))
}

func UpdateBaseName(c echo.Context) error {
	networkID, baseName, nerr := getNetworkIDAndBaseName(c)
	if nerr != nil {
		return nerr
	}

	bnr := &models.BaseNameRecord{}
	if err := c.Bind(bnr); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if string(bnr.Name) != baseName {
		return obsidian.HttpError(errors.New("base name in body does not match URL param"), http.StatusBadRequest)
	}

	// 404 if the entity doesn't exist
	exists, err := configurator.DoesEntityExist(networkID, lte.BaseNameEntityType, baseName)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "Failed to check if base name exists"), http.StatusInternalServerError)
	}
	if !exists {
		return echo.ErrNotFound
	}

	_, err = configurator.UpdateEntity(networkID, bnr.ToEntityUpdateCriteria())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func DeleteBaseName(c echo.Context) error {
	networkID, baseName, nerr := getNetworkIDAndBaseName(c)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(networkID, lte.BaseNameEntityType, baseName)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

// Rules

func ListRules(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	view := c.QueryParam("view")
	if strings.ToLower(view) == "full" {
		rules, err := configurator.LoadAllEntitiesInNetwork(
			networkID, lte.PolicyRuleEntityType,
			configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		ret := map[string]*models.PolicyRule{}
		for _, ruleEnt := range rules {
			ret[ruleEnt.Key] = (&models.PolicyRule{}).FromEntity(ruleEnt)
		}
		return c.JSON(http.StatusOK, ret)
	} else {
		ruleIDs, err := configurator.ListEntityKeys(networkID, lte.PolicyRuleEntityType)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		sort.Strings(ruleIDs)
		return c.JSON(http.StatusOK, ruleIDs)
	}
}

func CreateRule(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	rule := new(models.PolicyRule)
	if err := c.Bind(rule); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := rule.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	_, err := configurator.CreateEntity(networkID, rule.ToEntity())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func GetRule(c echo.Context) error {
	networkID, ruleID, nerr := getNetworkAndRuleIDs(c)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(
		networkID,
		lte.PolicyRuleEntityType,
		ruleID,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
	)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, (&models.PolicyRule{}).FromEntity(ent))
}

func UpdateRule(c echo.Context) error {
	networkID, ruleID, nerr := getNetworkAndRuleIDs(c)
	if nerr != nil {
		return nerr
	}

	rule := new(models.PolicyRule)
	if err := c.Bind(rule); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := rule.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if ruleID != string(rule.ID) {
		return obsidian.HttpError(errors.New("rule ID in body does not match URL param"), http.StatusBadRequest)
	}

	// 404 if rule doesn't exist
	exists, err := configurator.DoesEntityExist(networkID, lte.PolicyRuleEntityType, ruleID)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "Failed to check if rule exists"), http.StatusInternalServerError)
	}
	if !exists {
		return echo.ErrNotFound
	}

	_, err = configurator.UpdateEntity(networkID, rule.ToEntityUpdateCriteria())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func DeleteRule(c echo.Context) error {
	networkID, ruleID, nerr := getNetworkAndRuleIDs(c)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(networkID, lte.PolicyRuleEntityType, ruleID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func getNetworkIDAndBaseName(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", baseNameParam)
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}

func getNetworkAndRuleIDs(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", ruleIDParam)
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}
