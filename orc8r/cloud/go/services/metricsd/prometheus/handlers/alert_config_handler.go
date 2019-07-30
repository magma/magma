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
	neturl "net/url"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alert"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/exporters"

	"github.com/labstack/echo"
	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

const (
	alertConfigPart     = "alert_config"
	alertReceiverPart   = "alert_receiver"
	AlertNameQueryParam = "alert_name"
	AlertNamePathParam  = "alert_name"

	PrometheusRoot = obsidian.RestRoot + obsidian.UrlSep + "networks" + obsidian.UrlSep + ":network_id" + obsidian.UrlSep + "prometheus"

	AlertConfigURL         = PrometheusRoot + obsidian.UrlSep + alertConfigPart
	AlertUpdateURL         = AlertConfigURL + obsidian.UrlSep + ":" + AlertNamePathParam
	AlertReceiverConfigURL = PrometheusRoot + obsidian.UrlSep + alertReceiverPart
)

func GetConfigurePrometheusAlertHandler(configManagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(configManagerURL, networkID)
		return configurePrometheusAlert(c, url, networkID)
	}
}

func GetRetrieveAlertRuleHandler(configManagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(configManagerURL, networkID)
		return retrieveAlertRule(c, url)
	}
}

func GetDeleteAlertRuleHandler(configManagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(configManagerURL, networkID)
		return deleteAlertRule(c, url)
	}
}

func GetUpdateAlertRuleHandler(configManagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(configManagerURL, networkID)
		return updateAlertRule(c, url)
	}
}

func GetViewFiringAlertHandler(alertmanagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		return viewFiringAlerts(c, networkID, alertmanagerURL)
	}
}

func configurePrometheusAlert(c echo.Context, url, networkID string) error {
	rule, err := buildRuleFromContext(c)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("misconfigured rule: %v", err), http.StatusBadRequest)
	}

	err = alert.SecureRule(&rule, networkID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	errs := rule.Validate()
	if len(errs) != 0 {
		return obsidian.HttpError(fmt.Errorf("invalid rule: %v\n", errs), http.StatusBadRequest)
	}

	err = sendConfig(rule, url, http.MethodPost)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, rule.Alert)
}

func sendConfig(payload interface{}, url string, method string) error {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error making %s request: %v\n", method, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return obsidian.HttpError(fmt.Errorf("error writing config: %v", body.Message), resp.StatusCode)
	}
	return nil
}

func retrieveAlertRule(c echo.Context, url string) error {
	alertName := c.QueryParam(AlertNameQueryParam)
	if alertName != "" {
		url += fmt.Sprintf("?%s=%s", AlertNameQueryParam, neturl.QueryEscape(alertName))
	}

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return obsidian.HttpError(fmt.Errorf("error reading rules: %v", body.Message), resp.StatusCode)
	}

	var rules []alert.RuleJSONWrapper
	err = json.NewDecoder(resp.Body).Decode(&rules)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("error decoding server response: %v", err), http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, rules)
}

func deleteAlertRule(c echo.Context, url string) error {
	alertName := c.QueryParam(AlertNameQueryParam)
	if alertName == "" {
		return obsidian.HttpError(fmt.Errorf("alert name not provided"), http.StatusBadRequest)
	}
	url += fmt.Sprintf("?%s=%s", AlertNameQueryParam, neturl.QueryEscape(alertName))

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("could not form request: %v", err), http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return obsidian.HttpError(fmt.Errorf("error deleting rule: %v", body.Message), resp.StatusCode)
	}
	return c.JSON(http.StatusOK, nil)
}

func updateAlertRule(c echo.Context, url string) error {
	rule, err := buildRuleFromContext(c)
	if err != nil {
		return err
	}
	alertName := c.Param(AlertNamePathParam)
	if alertName == "" {
		return obsidian.HttpError(fmt.Errorf("alert name not provided"), http.StatusBadRequest)
	}
	url += fmt.Sprintf("/%s", neturl.PathEscape(alertName))

	err = sendConfig(rule, url, http.MethodPut)
	if err != nil {
		return err
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
		return obsidian.HttpError(fmt.Errorf("error decoding alertmanager response: %v", err), http.StatusInternalServerError)
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

	if jsonRule.Labels == nil {
		jsonRule.Labels = make(map[string]string)
	}
	if jsonRule.Annotations == nil {
		jsonRule.Annotations = make(map[string]string)
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

func alertConfigURL(hostName, networkID string) string {
	return hostName + "/" + networkID + "/alert"
}
