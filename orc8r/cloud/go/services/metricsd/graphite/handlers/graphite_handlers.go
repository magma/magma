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
	"regexp"
	"sort"
	"strings"
	"time"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/metricsd/graphite/exporters"
	"magma/orc8r/cloud/go/services/metricsd/graphite/third_party/api"
	"magma/orc8r/cloud/go/services/metricsd/obsidian/utils"

	"github.com/labstack/echo"
)

const (
	queryPart = "query"

	QueryURL = handlers.GRAPHITE_ROOT + handlers.URL_SEP + queryPart
	Protocol = "http"

	// matches graphite metric name. Alphanumeric plus '.' and '_'
	metricNameRegexString = `[a-zA-Z_\.\+\*\d]+`
)

var (
	tagRegexString = fmt.Sprintf(",%s=%s", metricNameRegexString, metricNameRegexString)
	// regex to match an list of tags
	tagListRegex = regexp.MustCompile(fmt.Sprintf("(%s)+", tagRegexString))
	//regex to match a metric name followed by an optional list of tags
	basicQueryRegex = regexp.MustCompile(fmt.Sprintf("^(%s)(%s)*$", metricNameRegexString, tagRegexString))
)

func GetQueryHandler(gc *api.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		return graphiteQuery(c, gc)
	}
}

func graphiteQuery(c echo.Context, gc *api.Client) error {
	startTime, err := utils.ParseTime(c.QueryParam(utils.ParamRangeStart), nil)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeEnd, err), http.StatusBadRequest)
	}

	defaultTime := time.Now()
	endTime, err := utils.ParseTime(c.QueryParam(utils.ParamRangeEnd), &defaultTime)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeEnd, err), http.StatusBadRequest)
	}

	preparedQuery, err := prepareQuery(c)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("unable to prepare query: %v", err))
	}

	renderRequest := api.RenderRequest{
		From:    startTime,
		Until:   endTime,
		Targets: []string{preparedQuery},
	}

	res, err := gc.QueryRender(renderRequest)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, wrapGraphiteResult(res))
}

func wrapGraphiteResult(res []api.Series) GraphiteResultStruct {
	return GraphiteResultStruct{Status: utils.StatusSuccess, Result: res}
}

func prepareQuery(c echo.Context) (string, error) {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return "", nerr
	}

	query := c.QueryParam(utils.ParamQuery)
	if ok := validateQuery(query); !ok {
		return "", fmt.Errorf("invalid query: %s", query)
	}

	tags := parseTagsFromQuery(query)
	tags.Insert(exporters.NetworkTagName, networkID)

	metricName := extractMetricNameFromQuery(query)

	request := buildTaggedQuery(metricName, tags)
	return request, nil
}

func extractMetricNameFromQuery(query string) string {
	firstTagIndex := strings.Index(query, ",")
	if firstTagIndex != -1 {
		return query[:firstTagIndex]
	}
	return query
}

func buildTaggedQuery(metricName string, tags exporters.TagSet) string {
	nameSelector := fmt.Sprintf("'name=~^%s$'", metricName)
	var keys []string
	for key := range tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var query strings.Builder
	query.WriteString(nameSelector)
	for _, key := range keys {
		query.WriteString(fmt.Sprintf(",'%s=%s'", key, tags[key]))
	}
	return fmt.Sprintf("%s(%s)", "seriesByTag", query.String())
}

// GraphiteResultStruct carries the result and status
type GraphiteResultStruct struct {
	Status string       `json:"status"`
	Result []api.Series `json:"result"`
}

func validateQuery(query string) bool {
	// matches <metric_name>[,tag1=val1]*
	return basicQueryRegex.MatchString(query)
}

func parseTagsFromQuery(query string) exporters.TagSet {
	tags := make(exporters.TagSet)
	tagString := tagListRegex.FindString(query)
	if tagString == "" {
		return tags
	}

	tagsList := strings.Split(tagString, ",")[1:]
	for _, tag := range tagsList {
		equalsIndex := strings.Index(tag, "=")
		key := tag[:equalsIndex]
		val := tag[equalsIndex+1:]
		tags.Insert(key, val)
	}
	return tags
}
