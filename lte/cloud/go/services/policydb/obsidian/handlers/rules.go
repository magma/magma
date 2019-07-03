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

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/labstack/echo"
)

func listRules(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	rules, err := configurator.ListEntityKeys(networkID, lte.PolicyRuleEntityType)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, rules)
}

func createRule(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	rule := new(models.PolicyRule)
	if err := c.Bind(rule); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := rule.Verify(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	_, err := configurator.CreateEntity(networkID, rule.ToEntity())
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, rule.ID)
}

func getRule(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	ruleID := c.Param("rule_id")
	if len(ruleID) == 0 {
		return ruleIDHTTPErr()
	}

	ent, err := configurator.LoadEntity(
		networkID,
		lte.PolicyRuleEntityType,
		ruleID,
		configurator.EntityLoadCriteria{LoadConfig: true},
	)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, ent.Config)
}

func updateRule(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	rule := new(models.PolicyRule)
	if err := c.Bind(rule); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	ruleID, herr := getRuleID(c, rule)
	if herr != nil {
		return herr
	}
	rule.ID = ruleID
	if err := rule.Verify(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	err := configurator.CreateOrUpdateEntityConfig(networkID, lte.PolicyRuleEntityType, ruleID, rule)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteRule(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	ruleID := c.Param("rule_id")
	if len(ruleID) == 0 {
		return ruleIDHTTPErr()
	}

	err := configurator.DeleteEntity(networkID, lte.PolicyRuleEntityType, ruleID)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}
