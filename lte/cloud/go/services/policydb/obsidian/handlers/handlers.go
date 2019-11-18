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
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo"
)

const (
	policiesRootPath         = obsidian.RestRoot + "/networks/:network_id/policies"
	policyRuleRootPath       = policiesRootPath + "/rules"
	policyRuleManagePath     = policyRuleRootPath + "/:rule_id"
	policyBaseNameRootPath   = policiesRootPath + "/base_names"
	policyBaseNameManagePath = policyBaseNameRootPath + "/:base_name"
)

func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		// base names
		{Path: policyBaseNameRootPath, Methods: obsidian.GET, HandlerFunc: ListBaseNames},
		{Path: policyBaseNameRootPath, Methods: obsidian.POST, HandlerFunc: CreateBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.GET, HandlerFunc: GetBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.PUT, HandlerFunc: UpdateBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.DELETE, HandlerFunc: DeleteBaseName},

		// rules
		{Path: policyRuleRootPath, Methods: obsidian.GET, HandlerFunc: ListRules},
		{Path: policyRuleRootPath, Methods: obsidian.POST, HandlerFunc: CreateRule},
		{Path: policyRuleManagePath, Methods: obsidian.GET, HandlerFunc: GetRule},
		{Path: policyRuleManagePath, Methods: obsidian.PUT, HandlerFunc: UpdateRule},
		{Path: policyRuleManagePath, Methods: obsidian.DELETE, HandlerFunc: DeleteRule},
	}
}

func ListBaseNames(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	baseNames, err := configurator.ListEntityKeys(networkID, lte.BaseNameEntityType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, baseNames)
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
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	baseName := getBaseNameParam(c)
	if len(baseName) == 0 {
		return baseNameHTTPErr()
	}

	ret, err := configurator.LoadEntity(
		networkID,
		lte.BaseNameEntityType,
		baseName,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
	)
	if err == errors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// TODO: only rulenames, no subs
	return c.JSON(http.StatusOK, (&models.BaseNameRecord{}).FromEntity(ret).RuleNames)
}

func UpdateBaseName(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	baseName := getBaseNameParam(c)
	if len(baseName) == 0 {
		return baseNameHTTPErr()
	}

	ruleNames := models.RuleNames{}
	if err := c.Bind(&ruleNames); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	_, err := configurator.UpdateEntity(
		networkID,
		configurator.EntityUpdateCriteria{
			Type:              lte.BaseNameEntityType,
			Key:               baseName,
			AssociationsToSet: ruleNames.ToAssocs(),
		},
	)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func DeleteBaseName(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	baseName := getBaseNameParam(c)
	if len(baseName) == 0 {
		return baseNameHTTPErr()
	}

	err := configurator.DeleteEntity(networkID, lte.BaseNameEntityType, baseName)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func baseNameHTTPErr() *echo.HTTPError {
	return obsidian.HttpError(
		fmt.Errorf("Invalid/Missing Base Name"),
		http.StatusBadRequest)
}

func getBaseNameParam(c echo.Context) string {
	return c.Param("base_name")
}

func ListRules(c echo.Context) error {
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

func CreateRule(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	rule := new(models.PolicyRule)
	if err := c.Bind(rule); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := rule.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	ent, err := rule.ToEntity()
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	_, err = configurator.CreateEntity(networkID, ent)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, rule.ID)
}

func GetRule(c echo.Context) error {
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
	rule, err := (&models.PolicyRule{}).FromEntity(ent)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, rule)
}

func UpdateRule(c echo.Context) error {
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
	if err := rule.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	cfg, err := rule.ToPolicyRuleConfig()
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	err = configurator.CreateOrUpdateEntityConfig(networkID, lte.PolicyRuleEntityType, ruleID, cfg)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func DeleteRule(c echo.Context) error {
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
