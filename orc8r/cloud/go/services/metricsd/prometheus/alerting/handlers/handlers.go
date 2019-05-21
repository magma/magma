package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alert"

	"github.com/labstack/echo"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

const (
	prometheusReloadPath = "/-/reload"
	ruleNameQueryParam   = "alert_name"
)

// GetPostHandler returns a handler that calls the client method WriteAlert() to
// write the alert configuration from the body of this request
func GetPostHandler(client *alert.Client, prometheusURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		rule, err := decodePostResponse(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		networkID := getNetworkID(c)
		err = client.WriteAlert(rule, networkID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		err = reloadPrometheus(prometheusURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.NoContent(http.StatusOK)
	}
}

func GetGetHandler(client *alert.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		ruleName := c.QueryParam(ruleNameQueryParam)
		networkID := getNetworkID(c)
		rules, err := client.ReadRules(ruleName, networkID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		jsonRules, err := rulesToJSON(rules)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, jsonRules)
	}
}

func GetDeleteHandler(client *alert.Client, prometheusURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		ruleName := c.QueryParam(ruleNameQueryParam)
		networkID := getNetworkID(c)
		if ruleName == "" {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("No rule name provided"))
		}
		err := client.DeleteRule(ruleName, networkID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		err = reloadPrometheus(prometheusURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, nil)
	}
}

func decodePostResponse(c echo.Context) (rulefmt.Rule, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return rulefmt.Rule{}, fmt.Errorf("error reading request body: %v", err)
	}
	payload := rulefmt.Rule{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return payload, fmt.Errorf("error unmarshalling payload: %v", err)
	}
	return payload, nil
}

func reloadPrometheus(url string) error {
	resp, err := http.Post(fmt.Sprintf("http://%s%s", url, prometheusReloadPath), "text/plain", &bytes.Buffer{})
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("code: %d error reloading prometheus: %v", resp.StatusCode, err)
	}
	return nil
}

func getNetworkID(c echo.Context) string {
	return c.Param("network_id")
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
