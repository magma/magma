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
	"net/http"

	"magma/orc8r/cloud/go/obsidian"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/olivere/elastic/v7"
)

const (
	Events               = "events"
	EventsPath           = obsidian.V1Root + Events + obsidian.UrlSep + ":" + pathParamStreamName
	pathParamStreamName  = "stream_name"
	queryParamEventType  = "event_type"
	queryParamHardwareID = "hardware_id"
	queryParamTag        = "tag"
	defaultQuerySize     = 50
)

// Returns a Hander that uses the provided elastic client
func GetEventsHandler(client *elastic.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		return EventsHandler(c, client)
	}
}

// Handles event querying
func EventsHandler(c echo.Context, client *elastic.Client) error {
	queryParams, err := getQueryParameters(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	elasticQuery := queryParams.ToElasticBoolQuery()

	result, err := client.Search().
		Index("").
		Size(defaultQuerySize).
		Query(elasticQuery).
		Do(context.Background())
	if err != nil {
		glog.Fatalf("Error getting response: %s", err)
	}
	if result.Error != nil {
		return obsidian.HttpError(fmt.Errorf("Elastic Error Type: %s, Reason: %s", result.Error.Type, result.Error.Reason))
	}
	return c.JSON(http.StatusOK, result.Hits.Hits)
}

func getQueryParameters(c echo.Context) (eventQueryParams, error) {
	pathParams, pathParamError := obsidian.GetParamValues(c, pathParamStreamName)
	if pathParamError != nil {
		return eventQueryParams{}, pathParamError
	}
	streamName := pathParams[0]
	params := eventQueryParams{
		StreamName: streamName,
		EventType:  c.QueryParam(queryParamEventType),
		HardwareID: c.QueryParam(queryParamHardwareID),
		Tag:        c.QueryParam(queryParamTag),
	}
	return params, nil
}

type eventQueryParams struct {
	StreamName string
	EventType  string
	HardwareID string
	Tag        string
}

func (b *eventQueryParams) ToElasticBoolQuery() *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	query.Filter(elastic.NewTermQuery(pathParamStreamName, b.StreamName))
	if len(b.EventType) > 0 {
		query.Filter(elastic.NewTermQuery(queryParamEventType, b.EventType))
	}
	if len(b.HardwareID) > 0 {
		query.Filter(elastic.NewTermQuery(queryParamHardwareID, b.HardwareID))
	}
	if len(b.Tag) > 0 {
		query.Filter(elastic.NewTermQuery(queryParamTag, b.Tag))
	}
	return query
}
