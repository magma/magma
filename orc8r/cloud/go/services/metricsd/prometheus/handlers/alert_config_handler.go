/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	neturl "net/url"

	"github.com/facebookincubator/prometheus-configmanager/prometheus/alert"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/prometheus/prometheus/pkg/rulefmt"

	"magma/orc8r/cloud/go/services/obsidian"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	"magma/orc8r/lib/go/metrics"
)

const (
	alertConfigPart     = "alert_config"
	alertReceiverPart   = "alert_receiver"
	AlertNameQueryParam = "alert_name"
	AlertNamePathParam  = "alert_name"

	AlertConfigV1URL         = PrometheusV1Root + obsidian.UrlSep + alertConfigPart
	AlertUpdateV1URL         = AlertConfigV1URL + obsidian.UrlSep + ":" + AlertNamePathParam
	AlertReceiverConfigV1URL = PrometheusV1Root + obsidian.UrlSep + alertReceiverPart
	AlertReceiverUpdateV1URL = AlertReceiverConfigV1URL + obsidian.UrlSep + ":" + ReceiverNamePathParam
	AlertBulkUpdateV1URL     = AlertConfigV1URL + "/bulk"

	FiringAlertV1URL = handlers.ManageNetworkPath + obsidian.UrlSep + "alerts"

	AlertSilencerV1URL = FiringAlertV1URL + obsidian.UrlSep + "silence"

	alertmanagerAPIAlertPath = "/alerts"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

func GetConfigurePrometheusAlertHandler(configManagerURL string, client HttpClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(networkID, configManagerURL)
		return configurePrometheusAlert(networkID, url, c, client)
	}
}

func GetRetrieveAlertRuleHandler(configManagerURL string, client HttpClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(networkID, configManagerURL)
		return retrieveAlertRule(c, url, client)
	}
}

func GetDeleteAlertRuleHandler(configManagerURL string, client HttpClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(networkID, configManagerURL)
		return deleteAlertRule(c, url, client)
	}
}

func GetUpdateAlertRuleHandler(configManagerURL string, client HttpClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(networkID, configManagerURL)
		return updateAlertRule(c, url, client)
	}
}

func GetBulkUpdateAlertHandler(configManagerURL string, client HttpClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := alertConfigURL(networkID, configManagerURL)
		url += "/bulk"
		return bulkUpdateAlerts(c, url, client)
	}
}

func GetViewFiringAlertHandler(alertmanagerURL string, client HttpClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		getAlertsURL := alertmanagerURL + alertmanagerAPIAlertPath
		return viewFiringAlerts(networkID, getAlertsURL, c, client)
	}
}

func configurePrometheusAlert(networkID, url string, c echo.Context, client HttpClient) error {
	rule, err := buildRuleFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("misconfigured rule: %v", err))
	}

	err = alert.SecureRule(true, metrics.NetworkLabelName, networkID, &rule)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	errs := rule.Validate()
	if len(errs) != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid rule: %v", errs))
	}

	sendErr := sendConfig(rule, url, http.MethodPost, client)
	if sendErr != nil {
		return sendErr
	}
	return c.JSON(http.StatusCreated, rule.Alert)
}

func sendConfig(payload interface{}, url string, method string, client HttpClient) *echo.HTTPError {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("create http request: %v", err))
	}
	resp, err := client.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("make %s request: %v", method, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return echo.NewHTTPError(resp.StatusCode, fmt.Sprintf("error writing config: %v", body.Message))
	}
	return nil
}

func retrieveAlertRule(c echo.Context, url string, client HttpClient) error {
	alertName := c.QueryParam(AlertNameQueryParam)
	if alertName != "" {
		url += fmt.Sprintf("/%s", neturl.PathEscape(alertName))
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return echo.NewHTTPError(resp.StatusCode, fmt.Sprintf("error reading rules: %v", body.Message))
	}

	var rules []alert.RuleJSONWrapper
	err = json.NewDecoder(resp.Body).Decode(&rules)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error decoding server response: %v", err))
	}
	return c.JSON(http.StatusOK, rules)
}

func deleteAlertRule(c echo.Context, url string, client HttpClient) error {
	alertName := c.QueryParam(AlertNameQueryParam)
	if alertName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "alert name not provided")
	}
	url += fmt.Sprintf("/%s", neturl.PathEscape(alertName))

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("could not form request: %v", err))
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return echo.NewHTTPError(resp.StatusCode, fmt.Sprintf("error deleting rule: %v", body.Message))
	}
	return c.JSON(http.StatusOK, nil)
}

func updateAlertRule(c echo.Context, url string, client HttpClient) error {
	rule, err := buildRuleFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("misconfigured rule: %v", err))
	}
	alertName := c.Param(AlertNamePathParam)
	if alertName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "alert name not provided")
	}
	url += fmt.Sprintf("/%s", neturl.PathEscape(alertName))

	sendErr := sendConfig(rule, url, http.MethodPut, client)
	if sendErr != nil {
		return sendErr
	}
	return c.JSON(http.StatusOK, nil)
}

func bulkUpdateAlerts(c echo.Context, url string, client HttpClient) error {
	rules, err := buildRuleListFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error parsing rule payload: %v", err))
	}

	resp, err := sendBulkConfig(rules, url, client)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func sendBulkConfig(payload interface{}, url string, client HttpClient) (string, error) {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("create http request: %v", err))
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("make PUT request: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return "", echo.NewHTTPError(resp.StatusCode, fmt.Sprintf("error writing config: %v", body.Message))
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return string(contents), nil
}

func viewFiringAlerts(networkID, getAlertsURL string, c echo.Context, client HttpClient) error {
	resp, err := client.Get(getAlertsURL)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return echo.NewHTTPError(resp.StatusCode, fmt.Sprintf("alertmanager error: %v", body.Message))
	}
	defer resp.Body.Close()

	var alerts []models.GettableAlert
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error decoding alertmanager response: %v", err))
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
