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

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian/handlers"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

const (
	policiesRootPath     = handlers.REST_ROOT + "/networks/:network_id/policies"
	policyRuleRootPath   = policiesRootPath + "/rules"
	policyRuleManagePath = policyRuleRootPath + "/:rule_id"
)

// listRulesHandler returns a list of all policy rules in the network
func listRulesHandler(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	rules, err := policydb.ListRuleIds(networkID)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, rules)
}

// createRuleHandler adds a single policy rule to the network
func createRuleHandler(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	// Get rule proto from request and validate
	ruleProto, err := getRuleProto(c)
	if err != nil {
		return err
	}

	// Call policydb service
	if err := policydb.AddRule(networkID, ruleProto); err != nil {
		return handlers.HttpError(err, http.StatusConflict)
	}
	return c.JSON(http.StatusCreated, ruleProto.GetId())
}

// getRuleHandler returns the policy rule associated with the input rule id
func getRuleHandler(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	ruleID := c.Param("rule_id")
	if len(ruleID) == 0 {
		return ruleIDHTTPErr()
	}

	// Call policydb service
	ruleProto, err := policydb.GetRule(networkID, ruleID)
	if err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}

	// Create swagger model for response
	var rule models.PolicyRule
	if err = rule.FromProto(ruleProto); err != nil {
		glog.Errorf("Error converting policy rule model: %s", err)
		return handlers.HttpError(err)
	}
	return c.JSON(http.StatusOK, rule)
}

// updateRuleHandler modifies the policy rule
func updateRuleHandler(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	// Get rule proto from request and validate
	ruleProto, err := getRuleProto(c)
	if err != nil {
		return err
	}

	// Call policydb service
	if err := policydb.UpdateRule(networkID, ruleProto); err != nil {
		return handlers.HttpError(err, http.StatusConflict)
	}
	return c.NoContent(http.StatusOK)
}

// deleteRuleHandler deletes the policy rule
func deleteRuleHandler(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	ruleID := c.Param("rule_id")
	if len(ruleID) == 0 {
		return ruleIDHTTPErr()
	}

	// Call policydb service
	if err := policydb.DeleteRule(networkID, ruleID); err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}
	return c.NoContent(http.StatusNoContent)
}

func getRuleID(c echo.Context, rule *models.PolicyRule) (string, *echo.HTTPError) {
	// The RuleId can be defined as URL param ie. "rule_id" or in the request body
	ruleID := c.Param("rule_id")
	if len(ruleID) != 0 {
		if rule.ID != ruleID {
			msg := fmt.Errorf("Rule ID payload doesn't match URL param %s vs %s",
				rule.ID, ruleID)
			return ruleID, handlers.HttpError(msg, http.StatusBadRequest)
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

func getRuleProto(c echo.Context) (*protos.PolicyRule, *echo.HTTPError) {
	// Construct response
	ruleProto := new(protos.PolicyRule)

	// Get swagger model
	rule := new(models.PolicyRule)
	if err := c.Bind(rule); err != nil {
		return ruleProto, handlers.HttpError(err, http.StatusBadRequest)
	}

	// Get the rule ID from the URL or request payload
	_, err := getRuleID(c, rule)
	if err != nil {
		return ruleProto, err
	}

	// Verify the payload
	if err := rule.Verify(); err != nil {
		return ruleProto, handlers.HttpError(err, http.StatusBadRequest)
	}

	// Convert swagger model to proto
	if err := rule.ToProto(ruleProto); err != nil {
		return ruleProto, handlers.HttpError(err)
	}

	return ruleProto, nil
}

func ruleIDHTTPErr() *echo.HTTPError {
	return handlers.HttpError(
		fmt.Errorf("Invalid/Missing Flow Rule ID"),
		http.StatusBadRequest)
}
