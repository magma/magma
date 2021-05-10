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
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"magma/orc8r/cloud/go/obsidian"
	eventdC "magma/orc8r/cloud/go/services/eventd/eventd_client"
	logH "magma/orc8r/cloud/go/services/eventd/log/handlers"

	"github.com/go-openapi/strfmt"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/olivere/elastic/v7"
)

const (
	defaultQuerySize = 50

	pathParamStreamName  = "stream_name"
	pathParamNetworkID   = "network_id"
	queryParamEventType  = "event_type"
	queryParamHardwareID = "hardware_id"
	queryParamTag        = "tag"

	// multi-stream endpoint query args
	pathParamStreams     = "streams"
	pathParamEvents      = "events"
	pathParamTags        = "tags"
	pathParamHardwareIDs = "hw_ids"
	pathParamFrom        = "from"
	pathParamSize        = "size"
	pathParamStart       = "start"
	pathParamEnd         = "end"

	EventsRootPath  = obsidian.V1Root + "events" + obsidian.UrlSep + ":" + pathParamNetworkID
	EventsPath      = EventsRootPath + obsidian.UrlSep + ":" + pathParamStreamName
	EventsCountPath = EventsRootPath + obsidian.UrlSep + "about" + obsidian.UrlSep + "count"

	ManageNetworkPath  = obsidian.V1Root + "networks" + obsidian.UrlSep + ":network_id"
	LogSearchQueryPath = ManageNetworkPath + obsidian.UrlSep + "logs" + obsidian.UrlSep + "search"
	LogCountQueryPath  = ManageNetworkPath + obsidian.UrlSep + "logs" + obsidian.UrlSep + "count"
)

// GetObsidianHandlers returns all the obsidian handlers for eventd.
func GetObsidianHandlers() []obsidian.Handler {
	var ret []obsidian.Handler

	client, err := eventdC.GetElasticClient()
	if err != nil {
		ret = append(ret, setInitErrorHandlers(err)...)
		return ret
	}

	ret = append(ret, obsidian.Handler{Path: LogSearchQueryPath, Methods: obsidian.GET, HandlerFunc: logH.GetQueryLogHandler(client)})
	ret = append(ret, obsidian.Handler{Path: LogCountQueryPath, Methods: obsidian.GET, HandlerFunc: logH.GetCountLogHandler(client)})
	ret = append(ret, obsidian.Handler{Path: EventsRootPath, Methods: obsidian.GET, HandlerFunc: GetMultiStreamEventsHandler(client)})
	ret = append(ret, obsidian.Handler{Path: EventsCountPath, Methods: obsidian.GET, HandlerFunc: GetEventCountHandler(client)})
	ret = append(ret, obsidian.Handler{Path: EventsPath, Methods: obsidian.GET, HandlerFunc: GetEventsHandler(client)})
	return ret
}

func setInitErrorHandlers(err error) []obsidian.Handler {
	return []obsidian.Handler{
		{Path: LogSearchQueryPath, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
		{Path: LogCountQueryPath, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
		{Path: EventsRootPath, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
		{Path: EventsPath, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
		{Path: EventsCountPath, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
	}
}

func getInitErrorHandler(err error) func(c echo.Context) error {
	return func(c echo.Context) error {
		return obsidian.HttpError(fmt.Errorf("initialization Error: %v", err), http.StatusInternalServerError)
	}
}

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

// GetEventCountHandler returns a handler for multi-stream elastic
// event count query endpoint.
func GetEventCountHandler(client *elastic.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		return EventCountHandler(c, client)
	}
}

// EventsHandler handles event querying using ES
func EventsHandler(c echo.Context, client *elastic.Client) error {
	queryParams, err := getQueryParameters(c)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	results, err := eventdC.GetEvents(c.Request().Context(), queryParams, client)
	if err != nil {
		glog.Error(err)
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, results)
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

	results, err := eventdC.GetMultiStreamEvents(c.Request().Context(), params, client)
	if err != nil {
		glog.Error(err)
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, results)
}

// EventCountHandler handles event counting queries using ES
func EventCountHandler(c echo.Context, client *elastic.Client) error {
	params, err := getMultiStreamQueryParameters(c)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	result, err := eventdC.GetEventCount(c.Request().Context(), params, client)
	if err != nil {
		glog.Error(err)
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, result)
}

func getQueryParameters(c echo.Context) (eventdC.EventQueryParams, error) {
	streamName := c.Param(pathParamStreamName)
	if streamName == "" {
		return eventdC.EventQueryParams{}, StreamNameHTTPErr()
	}
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return eventdC.EventQueryParams{}, nerr
	}
	params := eventdC.EventQueryParams{
		StreamName: streamName,
		EventType:  c.QueryParam(queryParamEventType),
		HardwareID: c.QueryParam(queryParamHardwareID),
		NetworkID:  networkID,
		Tag:        c.QueryParam(queryParamTag),
	}
	return params, nil
}

// StreamNameHTTPErr indicates that stream_name is missing
func StreamNameHTTPErr() *echo.HTTPError {
	return obsidian.HttpError(fmt.Errorf("Missing stream name"), http.StatusBadRequest)
}

func getMultiStreamQueryParameters(c echo.Context) (eventdC.MultiStreamEventQueryParams, error) {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return eventdC.MultiStreamEventQueryParams{}, nerr
	}
	ret := eventdC.MultiStreamEventQueryParams{
		NetworkID:   networkID,
		Streams:     getStringListParam(c, pathParamStreams),
		Events:      getStringListParam(c, pathParamEvents),
		Tags:        getStringListParam(c, pathParamTags),
		HardwareIDs: getStringListParam(c, pathParamHardwareIDs),
		Size:        defaultQuerySize,
	}

	fromVal, err := getIntegerParam(c, pathParamFrom)
	if err != nil {
		return eventdC.MultiStreamEventQueryParams{}, err
	}
	ret.From = fromVal

	sizeVal, err := getIntegerParam(c, pathParamSize)
	if err != nil {
		return eventdC.MultiStreamEventQueryParams{}, err
	}
	if sizeVal > 0 {
		ret.Size = sizeVal
	}

	startTime, err := getTimeParam(c, pathParamStart)
	if err != nil {
		return eventdC.MultiStreamEventQueryParams{}, err
	}
	ret.Start = startTime

	endTime, err := getTimeParam(c, pathParamEnd)
	if err != nil {
		return eventdC.MultiStreamEventQueryParams{}, err
	}
	ret.End = endTime

	return ret, nil
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
