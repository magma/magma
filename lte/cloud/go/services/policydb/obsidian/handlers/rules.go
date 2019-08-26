/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"fmt"
	"net/http"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/labstack/echo"
)

const (
	policiesRootPath     = obsidian.RestRoot + "/networks/:network_id/policies"
	policyRuleRootPath   = policiesRootPath + "/rules"
	policyRuleManagePath = policyRuleRootPath + "/:rule_id"
)

func listRules(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	rules, err := configurator.ListEntityKeys(networkID, lte.PolicyRuleEntityType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, rules)
}

func createRule(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	rule := new(models.PolicyRule)
	if err := c.Bind(rule); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := rule.Verify(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	_, err := configurator.CreateEntity(networkID, rule.ToEntity())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, rule.ID)
}

func getRule(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
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
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, ent.Config)
}

func updateRule(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	rule := new(models.PolicyRule)
	if err := c.Bind(rule); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	ruleID, herr := getRuleID(c, rule)
	if herr != nil {
		return herr
	}
	rule.ID = ruleID
	if err := rule.Verify(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	err := configurator.CreateOrUpdateEntityConfig(networkID, lte.PolicyRuleEntityType, ruleID, rule)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteRule(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	ruleID := c.Param("rule_id")
	if len(ruleID) == 0 {
		return ruleIDHTTPErr()
	}

	err := configurator.DeleteEntity(networkID, lte.PolicyRuleEntityType, ruleID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func getRuleProto(c echo.Context) (*protos.PolicyRule, *echo.HTTPError) {
	// Construct response
	ruleProto := new(protos.PolicyRule)

	// Get swagger model
	rule := new(models.PolicyRule)
	if err := c.Bind(rule); err != nil {
		return ruleProto, obsidian.HttpError(err, http.StatusBadRequest)
	}

	// Get the rule ID from the URL or request payload
	_, err := getRuleID(c, rule)
	if err != nil {
		return ruleProto, err
	}

	// Verify the payload
	if err := rule.Verify(); err != nil {
		return ruleProto, obsidian.HttpError(err, http.StatusBadRequest)
	}

	// Convert swagger model to proto
	if err := rule.ToProto(ruleProto); err != nil {
		return ruleProto, obsidian.HttpError(err)
	}

	return ruleProto, nil
}

func ruleIDHTTPErr() *echo.HTTPError {
	return obsidian.HttpError(
		fmt.Errorf("Invalid/Missing Flow Rule ID"),
		http.StatusBadRequest)
}

func getRuleID(c echo.Context, rule *models.PolicyRule) (string, *echo.HTTPError) {
	// The RuleId can be defined as URL param ie. "rule_id" or in the request body
	ruleID := c.Param("rule_id")
	if len(ruleID) != 0 {
		if rule.ID != ruleID {
			msg := fmt.Errorf("Rule ID payload doesn't match URL param %s vs %s",
				rule.ID, ruleID)
			return ruleID, obsidian.HttpError(msg, http.StatusBadRequest)
		}
		rule.ID = ruleID
	} else {
		ruleID = rule.ID
	}

	if len(ruleID) == 0 {
		return ruleID, ruleIDHTTPErr()
	}

	return ruleID, nil
}
