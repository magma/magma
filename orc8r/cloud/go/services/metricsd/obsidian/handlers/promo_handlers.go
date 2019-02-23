/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/common/model"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/metricsd/exporters"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/api/prometheus/v1"
)

const (
	prometheusAddressEnv = "PROMETHEUS_ADDRESS"

	queryPart      = "query"
	queryRangePart = "query_range"
	queryURL       = handlers.PROMETHEUS_ROOT + handlers.URL_SEP + queryPart
	queryRangeURL  = handlers.PROMETHEUS_ROOT + handlers.URL_SEP + queryRangePart

	paramQuery      = "query"
	paramRangeStart = "start"
	paramRangeEnd   = "end"
	paramStepWidth  = "step"
	paramTime       = "time"

	defaultStepWidth = "15s"

	statusSuccess = "success"
)

func getInitErrorHandler(err error) func(c echo.Context) error {
	return func(c echo.Context) error {
		return handlers.HttpError(fmt.Errorf("initialization Error: %v", err), 500)
	}
}

func getPrometheusQueryHandler(api v1.API) func(c echo.Context) error {
	return func(c echo.Context) error {
		restrictedQuery, err := preparePrometheusQuery(c)
		if err != nil {
			return handlers.HttpError(err, 500)
		}
		return prometheusQuery(c, restrictedQuery, api)
	}
}

func prometheusQuery(c echo.Context, query string, apiClient v1.API) error {
	defaultTime := time.Now()
	queryTime, err := parseTime(c.QueryParam(paramTime), &defaultTime)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("unable to parse %s parameter: %v", paramTime, err), http.StatusBadRequest)
	}

	res, err := apiClient.Query(context.Background(), query, queryTime)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, wrapPrometheusResult(res))
}

func getPrometheusQueryRangeHandler(api v1.API) func(c echo.Context) error {
	return func(c echo.Context) error {
		restrictedQuery, err := preparePrometheusQuery(c)
		if err != nil {
			return handlers.HttpError(err, 500)
		}
		return prometheusQueryRange(c, restrictedQuery, api)
	}
}

func prometheusQueryRange(c echo.Context, query string, apiClient v1.API) error {
	startTime, err := parseTime(c.QueryParam(paramRangeStart), nil)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("unable to parse %s parameter: %v", paramRangeEnd, err), http.StatusBadRequest)
	}

	defaultTime := time.Now()
	endTime, err := parseTime(c.QueryParam(paramRangeEnd), &defaultTime)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("unable to parse %s parameter: %v", paramRangeEnd, err), http.StatusBadRequest)
	}

	step, err := parseDuration(c.QueryParam(paramStepWidth), defaultStepWidth)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("unable to parse %s parameter: %v", paramRangeEnd, err), http.StatusBadRequest)
	}
	timeRange := v1.Range{Start: startTime, End: endTime, Step: step}

	res, err := apiClient.QueryRange(context.Background(), query, timeRange)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, wrapPrometheusResult(res))
}

func wrapPrometheusResult(res model.Value) PromQLResultStruct {
	dataStruct := PromQLDataStruct{ResultType: res.Type().String(), Result: res}
	return PromQLResultStruct{Status: statusSuccess, Data: dataStruct}
}

func preparePrometheusQuery(c echo.Context) (string, error) {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return "", nerr
	}

	restrictedQuery, err := preprocessQuery(c.QueryParam(paramQuery), networkID)
	if err != nil {
		return "", err
	}

	return restrictedQuery, nil
}

func preprocessQuery(query, networkID string) (string, error) {
	restrictedLabels := map[string]string{exporters.NetworkLabelInstance: networkID}
	restrictor := NewQueryRestrictor(restrictedLabels)
	return restrictor.RestrictQuery(query)
}

func parseTime(timeString string, defaultTime *time.Time) (time.Time, error) {
	if timeString == "" {
		if defaultTime != nil {
			return *defaultTime, nil
		}
		return time.Time{}, fmt.Errorf("time parameter not provided")
	}
	time, err := parseUnixTime(timeString)
	if err == nil {
		return time, nil
	}
	return parseRFCTime(timeString)
}

func parseUnixTime(timeString string) (time.Time, error) {
	timeNum, err := strconv.ParseFloat(timeString, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(timeNum), 0), nil
}

func parseRFCTime(timeString string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeString)
}

func parseDuration(durationString, defaultDuration string) (time.Duration, error) {
	if durationString == "" {
		durationString = defaultDuration
	}
	return time.ParseDuration(durationString)
}

// PromQLResultStruct carries all of the data of the full prometheus API result
type PromQLResultStruct struct {
	Status string           `json:"status"`
	Data   PromQLDataStruct `json:"data"`
}

// PromQLDataStruct carries the result type and actual metric result
type PromQLDataStruct struct {
	ResultType string      `json:"resultType"`
	Result     model.Value `json:"result"`
}
