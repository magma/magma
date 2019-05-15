/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alert"

	"github.com/labstack/echo"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

const (
	alertConfigPart = "alert_config"
	AlertConfigURL  = handlers.PROMETHEUS_ROOT + handlers.URL_SEP + alertConfigPart
)

func GetConfigurePrometheusAlertHandler(webServerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := handlers.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := webServerURL + "/" + networkID
		return configurePrometheusAlert(c, url, networkID)
	}
}

func configurePrometheusAlert(c echo.Context, url, networkID string) error {
	rule, err := buildRuleFromContext(c)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	err = alert.SecureRule(&rule, networkID)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	errs := rule.Validate()
	if len(errs) != 0 {
		return handlers.HttpError(fmt.Errorf("Invalid rule: %v\n", errs), http.StatusBadRequest)
	}

	err = sendConfigToPrometheusServer(rule, url)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, rule.Alert)
}

func sendConfigToPrometheusServer(payload rulefmt.Rule, url string) error {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body := echo.HTTPError{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		return handlers.HttpError(fmt.Errorf("server error: %v, code: %v", body.Message, body.Internal), http.StatusInternalServerError)
	}
	return nil
}

func buildRuleFromContext(c echo.Context) (rulefmt.Rule, error) {
	jsonRule := alert.RuleJSONWrapper{}
	err := json.NewDecoder(c.Request().Body).Decode(&jsonRule)
	if err != nil {
		return rulefmt.Rule{}, err
	}

	modelFor, err := model.ParseDuration(jsonRule.For)
	if err != nil {
		return rulefmt.Rule{}, err
	}
	rule := rulefmt.Rule{
		Record:      jsonRule.Record,
		Alert:       jsonRule.Alert,
		Expr:        jsonRule.Expr,
		For:         modelFor,
		Labels:      jsonRule.Labels,
		Annotations: jsonRule.Annotations,
	}
	return rule, nil
}
