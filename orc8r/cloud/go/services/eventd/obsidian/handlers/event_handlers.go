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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/eventd/obsidian/models"

	"github.com/go-openapi/strfmt"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/olivere/elastic/v7"
	"github.com/thoas/go-funk"
)

const (
	Events          = "events"
	EventsRootPath  = obsidian.V1Root + Events + obsidian.UrlSep + ":" + pathParamNetworkID
	EventsPath      = EventsRootPath + obsidian.UrlSep + ":" + pathParamStreamName
	EventsCountPath = EventsRootPath + obsidian.UrlSep + "about" + obsidian.UrlSep + "count"

	pathParamStreamName  = "stream_name"
	pathParamNetworkID   = "network_id"
	queryParamEventType  = "event_type"
	queryParamHardwareID = "hardware_id"
	queryParamTag        = "tag"

	defaultQuerySize = 50

	// We use the ES "keyword" type for exact match
	dotKeyword              = ".keyword"
	elasticFilterStreamName = pathParamStreamName + dotKeyword
	elasticFilterNetworkID  = pathParamNetworkID + dotKeyword
	elasticFilterEventType  = queryParamEventType + dotKeyword
	elasticFilterHardwareID = "hw_id" + dotKeyword
	elasticFilterEventTag   = "event_tag" + dotKeyword // We use event_tag as fluentd uses the "tag" field
	elasticFilterTimestamp  = "@timestamp"
)

// GetEventsHandler returns a Handler that uses the provided elastic client
func GetEventsHandler(client *elastic.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		return EventsHandler(c, client)
	}
}

// GetMultiStreamEventsHandler returns a handler for the multi-stream elastic
// event query endpoint.
func GetMultiStreamEventsHandler(client *elastic.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		return MultiStreamEventsHandler(c, client)
	}
}

func GetEventCountHandler(client *elastic.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		return EventCountHandler(c, client)
	}
}

// MultiStreamEventsHandler exposes more query options than EventsHandler,
// primarily the ability to query across multiple streams and tags.
// This handler will also accept an optional query size limit and offset for
// paginated queries.
func MultiStreamEventsHandler(c echo.Context, client *elastic.Client) error {
	params, err := getMultiStreamQueryParameters(c)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	query := params.toElasticBoolQuery()
	search := client.Search().
		Index("eventd*").
		From(params.from).
		Size(params.size).
		Sort(elasticFilterTimestamp, false).
		Query(query)
	return doSearch(c, search)
}

func EventCountHandler(c echo.Context, client *elastic.Client) error {
	params, err := getMultiStreamQueryParameters(c)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	query := params.toElasticBoolQuery()
	result, err := client.Count().
		Index("eventd*").
		Query(query).
		Do(c.Request().Context())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, result)
}

// EventsHandler handles event querying using ES
func EventsHandler(c echo.Context, client *elastic.Client) error {
	queryParams, err := getQueryParameters(c)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	elasticQuery := queryParams.ToElasticBoolQuery()
	search := client.Search().
		Index("").
		Size(defaultQuerySize).
		Sort(elasticFilterTimestamp, false).
		Query(elasticQuery)
	return doSearch(c, search)
}

func doSearch(c echo.Context, search *elastic.SearchService) error {
	result, err := search.Do(c.Request().Context())
	if err != nil {
		glog.Errorf("Error getting response from Elastic: %s", err)
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if result.Error != nil {
		return obsidian.HttpError(fmt.Errorf(
			"Elastic Error Type: %s, Reason: %s",
			result.Error.Type,
			result.Error.Reason))
	}

	eventResults, err := GetEventResults(result.Hits.Hits)
	if err != nil {
		glog.Error(err)
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, eventResults)
}

type eventElasticHit struct {
	StreamName string `json:"stream_name"`
	EventType  string `json:"event_type"`
	// FluentBit logs sent from AGW are tagged with hw_id
	HardwareID string `json:"hw_id"`
	// We use event_tag as fluentd uses the "tag" field
	Tag       string `json:"event_tag"`
	Timestamp string `json:"@timestamp"`
	Value     string `json:"value"`
}

// Retrieve Event properties from the _source of
// ES Hits, including event metadata
func GetEventResults(hits []*elastic.SearchHit) ([]models.Event, error) {
	results := []models.Event{}
	for _, hit := range hits {
		var eventHit eventElasticHit
		// Get Value from the _source
		if err := json.Unmarshal(hit.Source, &eventHit); err != nil {
			return nil, fmt.Errorf("Unable to Unmarshal JSON from elastic.Hit. "+
				"elastic.Hit.Source: %s, Error: %s", hit.Source, err)
		}
		// Skip hits without an event value
		if eventHit.Value == "" {
			return nil, fmt.Errorf("eventResult %s does not contain a value", eventHit)
		}
		var eventValue map[string]interface{}
		if err := json.Unmarshal([]byte(eventHit.Value), &eventValue); err != nil {
			return nil, fmt.Errorf("Unable to Unmarshal JSON from eventResult.Value. "+
				"eventHit.Value: %s, Error: %s", hit.Source, err)
		}
		results = append(results, models.Event{
			StreamName: eventHit.StreamName,
			EventType:  eventHit.EventType,
			HardwareID: eventHit.HardwareID,
			Tag:        eventHit.Tag,
			Timestamp:  eventHit.Timestamp,
			Value:      eventValue,
		})
	}
	return results, nil
}

func getQueryParameters(c echo.Context) (eventQueryParams, error) {
	streamName := c.Param(pathParamStreamName)
	if streamName == "" {
		return eventQueryParams{}, StreamNameHTTPErr()
	}
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return eventQueryParams{}, nerr
	}
	params := eventQueryParams{
		StreamName: streamName,
		EventType:  c.QueryParam(queryParamEventType),
		HardwareID: c.QueryParam(queryParamHardwareID),
		NetworkID:  networkID,
		Tag:        c.QueryParam(queryParamTag),
	}
	return params, nil
}

type eventQueryParams struct {
	StreamName string
	EventType  string
	HardwareID string
	NetworkID  string
	Tag        string
}

func (b *eventQueryParams) ToElasticBoolQuery() *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	query.Filter(elastic.NewTermQuery(elasticFilterStreamName, b.StreamName))
	query.Filter(elastic.NewTermQuery(elasticFilterNetworkID, b.NetworkID))
	if len(b.EventType) > 0 {
		query.Filter(elastic.NewTermQuery(elasticFilterEventType, b.EventType))
	}
	if len(b.HardwareID) > 0 {
		query.Filter(elastic.NewTermQuery(elasticFilterHardwareID, b.HardwareID))
	}
	if len(b.Tag) > 0 {
		query.Filter(elastic.NewTermQuery(elasticFilterEventTag, b.Tag))
	}
	return query
}

// StreamNameHTTPErr indicates that stream_name is missing
func StreamNameHTTPErr() *echo.HTTPError {
	return obsidian.HttpError(fmt.Errorf("Missing stream name"), http.StatusBadRequest)
}

// multi-stream endpoint query args
const (
	pathParamStreams     = "streams"
	pathParamEvents      = "events"
	pathParamTags        = "tags"
	pathParamHardwareIDs = "hw_ids"
	pathParamFrom        = "from"
	pathParamSize        = "size"
	pathParamStart       = "start"
	pathParamEnd         = "end"
)

func getMultiStreamQueryParameters(c echo.Context) (multiStreamEventQueryParams, error) {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return multiStreamEventQueryParams{}, nerr
	}
	ret := multiStreamEventQueryParams{
		networkID:   networkID,
		streams:     getStringListParam(c, pathParamStreams),
		events:      getStringListParam(c, pathParamEvents),
		tags:        getStringListParam(c, pathParamTags),
		hardwareIDs: getStringListParam(c, pathParamHardwareIDs),
		size:        defaultQuerySize,
	}

	fromVal, err := getIntegerParam(c, pathParamFrom)
	if err != nil {
		return multiStreamEventQueryParams{}, err
	}
	ret.from = fromVal

	sizeVal, err := getIntegerParam(c, pathParamSize)
	if err != nil {
		return multiStreamEventQueryParams{}, err
	}
	if sizeVal > 0 {
		ret.size = sizeVal
	}

	startTime, err := getTimeParam(c, pathParamStart)
	if err != nil {
		return multiStreamEventQueryParams{}, err
	}
	ret.start = startTime

	endTime, err := getTimeParam(c, pathParamEnd)
	if err != nil {
		return multiStreamEventQueryParams{}, err
	}
	ret.end = endTime

	return ret, nil
}

type multiStreamEventQueryParams struct {
	networkID   string
	streams     []string
	events      []string
	tags        []string
	hardwareIDs []string
	from        int
	size        int
	start       *time.Time
	end         *time.Time
}

func (m multiStreamEventQueryParams) toElasticBoolQuery() *elastic.BoolQuery {
	ret := elastic.NewBoolQuery().Filter(elastic.NewTermQuery(elasticFilterNetworkID, m.networkID))
	if !funk.IsEmpty(m.streams) {
		ret.Filter(elastic.NewTermsQuery(elasticFilterStreamName, stringsToInterfaces(m.streams)...))
	}
	if !funk.IsEmpty(m.events) {
		ret.Filter(elastic.NewTermsQuery(elasticFilterEventType, stringsToInterfaces(m.events)...))
	}
	if !funk.IsEmpty(m.tags) {
		ret.Filter(elastic.NewTermsQuery(elasticFilterEventTag, stringsToInterfaces(m.tags)...))
	}
	if !funk.IsEmpty(m.hardwareIDs) {
		ret.Filter(elastic.NewTermsQuery(elasticFilterHardwareID, stringsToInterfaces(m.hardwareIDs)...))
	}
	if m.start != nil || m.end != nil {
		ret.Must(elastic.NewRangeQuery(elasticFilterTimestamp).From(m.start).To(m.end))
	}
	return ret
}

const urlListDelimiter = ","

func getIntegerParam(c echo.Context, param string) (int, error) {
	if valStr := c.QueryParam(param); valStr != "" {
		return strconv.Atoi(valStr)
	}
	return 0, nil
}

func getStringListParam(c echo.Context, param string) []string {
	if valStr := c.QueryParam(param); valStr != "" {
		return strings.Split(valStr, urlListDelimiter)
	}
	return []string{}
}

func getTimeParam(c echo.Context, param string) (*time.Time, error) {
	if valStr := c.QueryParam(param); valStr != "" {
		dt, err := strfmt.ParseDateTime(valStr)
		if err != nil {
			return nil, err
		}
		nativeDT := time.Time(dt)
		return &nativeDT, nil
	}
	return nil, nil
}

func stringsToInterfaces(st []string) []interface{} {
	ret := make([]interface{}, 0, len(st))
	for _, s := range st {
		ret = append(ret, s)
	}
	return ret
}
