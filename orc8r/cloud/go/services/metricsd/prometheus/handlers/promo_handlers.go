/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"magma/orc8r/cloud/go/metrics"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	"magma/orc8r/cloud/go/services/metricsd/obsidian/security"
	"magma/orc8r/cloud/go/services/metricsd/obsidian/utils"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	queryPart      = "query"
	queryRangePart = "query_range"
	seriesPart     = "series"
	paramMatch     = "match"

	PrometheusRoot   = obsidian.NetworksRoot + obsidian.UrlSep + ":network_id" + obsidian.UrlSep + "prometheus"
	PrometheusV1Root = handlers.ManageNetworkPath + obsidian.UrlSep + "prometheus"

	QueryURL      = PrometheusRoot + obsidian.UrlSep + queryPart
	QueryRangeURL = PrometheusRoot + obsidian.UrlSep + queryRangePart

	QueryV1URL      = PrometheusV1Root + obsidian.UrlSep + queryPart
	QueryRangeV1URL = PrometheusV1Root + obsidian.UrlSep + queryRangePart
	SeriesV1URL     = PrometheusV1Root + obsidian.UrlSep + seriesPart

	defaultStepWidth = "15s"
)

func GetPrometheusQueryHandler(api v1.API) func(c echo.Context) error {
	return func(c echo.Context) error {
		restrictedQuery, err := preparePrometheusQuery(c)
		if err != nil {
			return obsidian.HttpError(err, 500)
		}
		return prometheusQuery(c, restrictedQuery, api)
	}
}

func prometheusQuery(c echo.Context, query string, apiClient v1.API) error {
	defaultTime := time.Now()
	queryTime, err := utils.ParseTime(c.QueryParam(utils.ParamTime), &defaultTime)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamTime, err), http.StatusBadRequest)
	}

	res, err := apiClient.Query(context.Background(), query, queryTime)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, wrapPrometheusResult(res))
}

func GetPrometheusQueryRangeHandler(api v1.API) func(c echo.Context) error {
	return func(c echo.Context) error {
		restrictedQuery, err := preparePrometheusQuery(c)
		if err != nil {
			return obsidian.HttpError(err, 500)
		}
		return prometheusQueryRange(c, restrictedQuery, api)
	}
}

func prometheusQueryRange(c echo.Context, query string, apiClient v1.API) error {
	startTime, err := utils.ParseTime(c.QueryParam(utils.ParamRangeStart), nil)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeEnd, err), http.StatusBadRequest)
	}

	defaultTime := time.Now()
	endTime, err := utils.ParseTime(c.QueryParam(utils.ParamRangeEnd), &defaultTime)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeEnd, err), http.StatusBadRequest)
	}

	step, err := utils.ParseDuration(c.QueryParam(utils.ParamStepWidth), defaultStepWidth)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeEnd, err), http.StatusBadRequest)
	}
	timeRange := v1.Range{Start: startTime, End: endTime, Step: step}

	res, err := apiClient.QueryRange(context.Background(), query, timeRange)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, wrapPrometheusResult(res))
}

func wrapPrometheusResult(res model.Value) PromQLResultStruct {
	dataStruct := PromQLDataStruct{ResultType: res.Type().String(), Result: res}
	return PromQLResultStruct{Status: utils.StatusSuccess, Data: dataStruct}
}

func preparePrometheusQuery(c echo.Context) (string, error) {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return "", nerr
	}

	restrictedQuery, err := preprocessQuery(c.QueryParam(utils.ParamQuery), networkID)
	if err != nil {
		return "", err
	}

	return restrictedQuery, nil
}

func preprocessQuery(query, networkID string) (string, error) {
	restrictedLabels := map[string]string{metrics.NetworkLabelName: networkID}
	restrictor := security.NewQueryRestrictor(restrictedLabels)
	return restrictor.RestrictQuery(query)
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

var (
	minTime = time.Unix(math.MinInt64/1000+62135596801, 0).UTC()
	maxTime = time.Unix(math.MaxInt64/1000-62135596801, 999999999).UTC()
)

func GetPrometheusSeriesHandler(api v1.API) func(c echo.Context) error {
	return func(c echo.Context) error {
		nID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return obsidian.HttpError(nerr, http.StatusBadRequest)
		}

		startTime, err := utils.ParseTime(c.QueryParam(utils.ParamRangeStart), &minTime)
		if err != nil {
			return obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeEnd, err), http.StatusBadRequest)
		}

		endTime, err := utils.ParseTime(c.QueryParam(utils.ParamRangeEnd), &maxTime)
		if err != nil {
			return obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeEnd, err), http.StatusBadRequest)
		}

		seriesMatches, err := getSeriesMatches(c, nID)
		if err != nil {
			return obsidian.HttpError(fmt.Errorf("Error parsing series matchers: %v", err), http.StatusBadRequest)
		}

		res, err := api.Series(context.Background(), seriesMatches, startTime, endTime)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, res)
	}
}

func getSeriesMatches(c echo.Context, networkID string) ([]string, error) {
	restrictor := security.NewQueryRestrictor(map[string]string{metrics.NetworkLabelName: networkID})
	// Split array of matches by space delimiter
	matches := strings.Split(c.QueryParam(paramMatch), " ")
	seriesMatchers := make([]string, 0, len(matches))
	for _, match := range matches {
		if match == "" {
			continue
		}
		restricted, err := restrictor.RestrictQuery(match)
		if err != nil {
			return []string{}, obsidian.HttpError(fmt.Errorf("unable to secure match parameter: %v", err))
		}
		seriesMatchers = append(seriesMatchers, restricted)
	}
	// Add networkMatch if none provided since prometheus performs an OR of
	// all matches provided, and requires at least one
	if len(seriesMatchers) == 0 {
		networkMatch := fmt.Sprintf(`{%s="%s"}`, metrics.NetworkLabelName, networkID)
		seriesMatchers = append(seriesMatchers, networkMatch)
	}
	return seriesMatchers, nil
}
