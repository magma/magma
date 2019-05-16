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
	"magma/orc8r/cloud/go/services/metricsd/prometheus/exporters"

	"github.com/labstack/echo"
	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

const (
	alertConfigPart     = "alert_config"
	AlertConfigURL      = handlers.PROMETHEUS_ROOT + handlers.URL_SEP + alertConfigPart
	AlertNameQueryParam = "alert_name"
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

func GetRetrieveAlertRuleHandler(webServerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := handlers.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := webServerURL + "/" + networkID
		return retrieveAlertRule(c, url)
	}
}

func GetDeleteAlertRuleHandler(webServerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := handlers.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := webServerURL + "/" + networkID
		return deleteAlertRule(c, url)
	}
}

func GetViewFiringAlertHandler(alertmanagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := handlers.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		return viewFiringAlerts(c, networkID, alertmanagerURL)
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

func retrieveAlertRule(c echo.Context, url string) error {
	alertName := c.QueryParam(AlertNameQueryParam)
	if alertName != "" {
		url += fmt.Sprintf("?%s=%s", AlertNameQueryParam, alertName)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handlers.HttpError(fmt.Errorf("alert server responded with error"), resp.StatusCode)
	}

	var rules []alert.RuleJSONWrapper
	err = json.NewDecoder(resp.Body).Decode(&rules)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("error decoding server response: %v", err))
	}
	return c.JSON(http.StatusOK, rules)
}

func deleteAlertRule(c echo.Context, url string) error {
	alertName := c.QueryParam(AlertNameQueryParam)
	if alertName == "" {
		return handlers.HttpError(fmt.Errorf("alert Name not provided"), http.StatusBadRequest)
	}
	url += fmt.Sprintf("?%s=%s", AlertNameQueryParam, alertName)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("could not form request: %v", err), http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handlers.HttpError(fmt.Errorf("alert server responded with error"), resp.StatusCode)
	}
	return c.JSON(http.StatusOK, nil)
}

func viewFiringAlerts(c echo.Context, networkID, alertmanagerApiURL string) error {
	client := &http.Client{}
	resp, err := client.Get(alertmanagerApiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var alerts []models.GettableAlert
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("error decoding alertmanager response: %v", err))
	}
	networkAlerts := getAlertsForNetwork(networkID, alerts)
	return c.JSON(http.StatusOK, networkAlerts)
}

func getAlertsForNetwork(networkID string, alerts []models.GettableAlert) []models.GettableAlert {
	networkAlerts := make([]models.GettableAlert, 0)
	for _, alert := range alerts {
		if labelVal, ok := alert.Labels[exporters.NetworkLabelNetwork]; ok {
			if labelVal == networkID {
				networkAlerts = append(networkAlerts, alert)
			}
		}
	}
	return networkAlerts
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
