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
	"time"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/metricsd/graphite/security"
	"magma/orc8r/cloud/go/services/metricsd/graphite/third_party/api"
	"magma/orc8r/cloud/go/services/metricsd/obsidian/utils"

	"github.com/labstack/echo"
)

const (
	queryPart = "query"

	GraphiteRoot = obsidian.RestRoot + obsidian.UrlSep + "networks" + obsidian.UrlSep + ":network_id" + obsidian.UrlSep + "graphite"
	QueryURL     = GraphiteRoot + obsidian.UrlSep + queryPart
	Protocol     = "http"
)

func GetQueryHandler(gc *api.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		return graphiteQuery(c, gc)
	}
}

func graphiteQuery(c echo.Context, gc *api.Client) error {
	startTime, err := utils.ParseTime(c.QueryParam(utils.ParamRangeStart), nil)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeEnd, err), http.StatusBadRequest)
	}

	defaultTime := time.Now()
	endTime, err := utils.ParseTime(c.QueryParam(utils.ParamRangeEnd), &defaultTime)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeEnd, err), http.StatusBadRequest)
	}

	preparedQuery, err := prepareQuery(c)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("unable to prepare query: %v", err))
	}

	renderRequest := api.RenderRequest{
		From:    startTime,
		Until:   endTime,
		Targets: []string{preparedQuery},
	}

	res, err := gc.QueryRender(renderRequest)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, wrapGraphiteResult(res))
}

func wrapGraphiteResult(res []api.Series) GraphiteResultStruct {
	return GraphiteResultStruct{Status: utils.StatusSuccess, Result: res}
}

func prepareQuery(c echo.Context) (string, error) {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return "", nerr
	}

	query := c.QueryParam(utils.ParamQuery)
	restrictedQuery, err := security.RestrictQuery(query, networkID)
	if err != nil {
		return "", fmt.Errorf("Could not secure query: %v", err)
	}
	return restrictedQuery, nil
}

// GraphiteResultStruct carries the result and status
type GraphiteResultStruct struct {
	Status string       `json:"status"`
	Result []api.Series `json:"result"`
}
