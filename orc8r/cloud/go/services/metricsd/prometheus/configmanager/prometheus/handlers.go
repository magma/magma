/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/prometheus/alert"

	"github.com/labstack/echo"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

const (
	rootPath        = "/:file_prefix"
	AlertPath       = rootPath + "/alert"
	AlertUpdatePath = AlertPath + "/:" + RuleNamePathParam
	AlertBulkPath   = AlertPath + "/bulk"

	ruleNameQueryParam = "alert_name"
	RuleNamePathParam  = "alert_name"
)

// GetConfigureAlertHandler returns a handler that calls the client method WriteAlert() to
// write the alert configuration from the body of this request
func GetConfigureAlertHandler(client alert.PrometheusAlertClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		rule, err := decodeRulePostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		filePrefix := getFilePrefix(c)

		err = client.ValidateRule(rule)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if client.RuleExists(filePrefix, rule.Alert) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Rule '%s' already exists", rule.Alert))
		}

		err = client.WriteRule(filePrefix, rule)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = client.ReloadPrometheus()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

func GetRetrieveAlertHandler(client alert.PrometheusAlertClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		ruleName := c.QueryParam(ruleNameQueryParam)
		filePrefix := getFilePrefix(c)
		rules, err := client.ReadRules(filePrefix, ruleName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		jsonRules, err := rulesToJSON(rules)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, jsonRules)
	}
}

func GetDeleteAlertHandler(client alert.PrometheusAlertClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		ruleName := c.QueryParam(ruleNameQueryParam)
		filePrefix := getFilePrefix(c)
		if ruleName == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "No rule name provided")
		}
		err := client.DeleteRule(filePrefix, ruleName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		err = client.ReloadPrometheus()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, fmt.Sprintf("rule %s deleted", ruleName))
	}
}

func GetUpdateAlertHandler(client alert.PrometheusAlertClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		ruleName := c.Param(RuleNamePathParam)
		filePrefix := getFilePrefix(c)
		if ruleName == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "No rule name provided")
		}

		if !client.RuleExists(filePrefix, ruleName) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Rule '%s' does not exist", ruleName))
		}

		rule, err := decodeRulePostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.ValidateRule(rule)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.UpdateRule(filePrefix, rule)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = client.ReloadPrometheus()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

func GetBulkAlertUpdateHandler(client alert.PrometheusAlertClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		filePrefix := getFilePrefix(c)

		rules, err := decodeBulkRulesPostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		for _, rule := range rules {
			err = client.ValidateRule(rule)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
		}

		results, err := client.BulkUpdateRules(filePrefix, rules)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.ReloadPrometheus()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, results)
	}
}

func decodeRulePostRequest(c echo.Context) (rulefmt.Rule, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return rulefmt.Rule{}, fmt.Errorf("error reading request body: %v", err)
	}
	// First try unmarshaling into prometheus rulefmt.Rule{}
	payload := rulefmt.Rule{}
	err = json.Unmarshal(body, &payload)
	if err == nil {
		return payload, nil
	}
	// Try to unmarshal into the RuleJSONWrapper struct if prometheus struct doesn't work
	jsonPayload := alert.RuleJSONWrapper{}
	err = json.Unmarshal(body, &jsonPayload)
	if err != nil {
		return payload, fmt.Errorf("error unmarshalling payload: %v", err)
	}
	return jsonPayload.ToRuleFmt()
}

func decodeBulkRulesPostRequest(c echo.Context) ([]rulefmt.Rule, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return []rulefmt.Rule{}, fmt.Errorf("error reading request body: %v", err)
	}
	var payload []rulefmt.Rule
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return payload, fmt.Errorf("error unmarshalling payload: %v", err)
	}
	return payload, nil
}

func getFilePrefix(c echo.Context) string {
	return c.Param("file_prefix")
}

func rulesToJSON(rules []rulefmt.Rule) ([]alert.RuleJSONWrapper, error) {
	ret := make([]alert.RuleJSONWrapper, 0)

	for _, rule := range rules {
		jsonRule, err := rulefmtToJSON(rule)
		if err != nil {
			return ret, err
		}
		ret = append(ret, *jsonRule)
	}
	return ret, nil
}

func rulefmtToJSON(rule rulefmt.Rule) (*alert.RuleJSONWrapper, error) {
	duration, err := time.ParseDuration(rule.For.String())
	if err != nil {
		return nil, err
	}
	return &alert.RuleJSONWrapper{
		Record:      rule.Record,
		Alert:       rule.Alert,
		Expr:        rule.Expr,
		For:         duration.String(),
		Labels:      rule.Labels,
		Annotations: rule.Annotations,
	}, nil

}
