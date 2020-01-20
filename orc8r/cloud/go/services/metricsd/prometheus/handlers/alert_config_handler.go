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
	"io/ioutil"
	"net/http"
	neturl "net/url"

	"magma/orc8r/cloud/go/metrics"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/prometheus/alert"

	"github.com/labstack/echo"
	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

const (
	alertConfigPart     = "alert_config"
	alertReceiverPart   = "alert_receiver"
	AlertNameQueryParam = "alert_name"
	AlertNamePathParam  = "alert_name"

	AlertConfigURL         = PrometheusRoot + obsidian.UrlSep + alertConfigPart
	AlertUpdateURL         = AlertConfigURL + obsidian.UrlSep + ":" + AlertNamePathParam
	AlertReceiverConfigURL = PrometheusRoot + obsidian.UrlSep + alertReceiverPart
	AlertReceiverUpdateURL = AlertReceiverConfigURL + obsidian.UrlSep + ":" + ReceiverNamePathParam
	AlertBulkUpdateURL     = AlertConfigURL + "/bulk"

	AlertConfigV1URL         = PrometheusV1Root + obsidian.UrlSep + alertConfigPart
	AlertUpdateV1URL         = AlertConfigV1URL + obsidian.UrlSep + ":" + AlertNamePathParam
	AlertReceiverConfigV1URL = PrometheusV1Root + obsidian.UrlSep + alertReceiverPart
	AlertReceiverUpdateV1URL = AlertReceiverConfigV1URL + obsidian.UrlSep + ":" + ReceiverNamePathParam
	AlertBulkUpdateV1URL     = AlertConfigV1URL + "/bulk"

	FiringAlertURL   = obsidian.NetworksRoot + obsidian.UrlSep + ":network_id" + obsidian.UrlSep + "alerts"
	FiringAlertV1URL = handlers.ManageNetworkPath + obsidian.UrlSep + "alerts"
)

func GetConfigurePrometheusAlertHandler(configManagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(networkID, configManagerURL)
		return configurePrometheusAlert(networkID, url, c)
	}
}

func GetRetrieveAlertRuleHandler(configManagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(networkID, configManagerURL)
		return retrieveAlertRule(c, url)
	}
}

func GetDeleteAlertRuleHandler(configManagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(networkID, configManagerURL)
		return deleteAlertRule(c, url)
	}
}

func GetUpdateAlertRuleHandler(configManagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(networkID, configManagerURL)
		return updateAlertRule(c, url)
	}
}

func GetBulkUpdateAlertHandler(configManagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(networkID, configManagerURL)
		url += "/bulk"
		return bulkUpdateAlerts(c, url)
	}
}

func GetViewFiringAlertHandler(alertmanagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		return viewFiringAlerts(networkID, alertmanagerURL, c)
	}
}

func configurePrometheusAlert(networkID, url string, c echo.Context) error {
	rule, err := buildRuleFromContext(c)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("misconfigured rule: %v", err), http.StatusBadRequest)
	}

	err = alert.SecureRule(networkID, &rule)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	errs := rule.Validate()
	if len(errs) != 0 {
		return obsidian.HttpError(fmt.Errorf("invalid rule: %v\n", errs), http.StatusBadRequest)
	}

	sendErr := sendConfig(rule, url, http.MethodPost)
	if sendErr != nil {
		return obsidian.HttpError(sendErr, sendErr.Code)
	}
	return c.JSON(http.StatusCreated, rule.Alert)
}

func sendConfig(payload interface{}, url string, method string) *echo.HTTPError {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	resp, err := client.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error making %s request: %v", method, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return echo.NewHTTPError(resp.StatusCode, fmt.Errorf("error writing config: %v", body.Message))
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
		return obsidian.HttpError(fmt.Errorf("misconfigured rule: %v", err), http.StatusBadRequest)
	}
	alertName := c.Param(AlertNamePathParam)
	if alertName == "" {
		return obsidian.HttpError(fmt.Errorf("alert name not provided"), http.StatusBadRequest)
	}
	url += fmt.Sprintf("/%s", neturl.PathEscape(alertName))

	sendErr := sendConfig(rule, url, http.MethodPut)
	if err != nil {
		return obsidian.HttpError(sendErr, sendErr.Code)
	}
	return c.JSON(http.StatusOK, nil)
}

func bulkUpdateAlerts(c echo.Context, url string) error {
	rules, err := buildRuleListFromContext(c)
	if err != nil {
		return err
	}

	resp, err := sendBulkConfig(rules, url)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func sendBulkConfig(payload interface{}, url string) (string, error) {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making PUT request: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return "", obsidian.HttpError(fmt.Errorf("error writing config: %v", body.Message), resp.StatusCode)
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func viewFiringAlerts(networkID, alertmanagerApiURL string, c echo.Context) error {
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
		if labelVal, ok := alert.Labels[metrics.NetworkLabelName]; ok {
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
	return jsonRule.ToRuleFmt()
}

func buildRuleListFromContext(c echo.Context) ([]rulefmt.Rule, error) {
	var jsonRules []alert.RuleJSONWrapper
	err := json.NewDecoder(c.Request().Body).Decode(&jsonRules)
	if err != nil {
		return []rulefmt.Rule{}, err
	}

	var rules []rulefmt.Rule
	for _, jsonRule := range jsonRules {
		rule, err := jsonRule.ToRuleFmt()
		if err != nil {
			return []rulefmt.Rule{}, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func alertConfigURL(networkID, hostName string) string {
	return hostName + "/" + networkID + "/alert"
}
