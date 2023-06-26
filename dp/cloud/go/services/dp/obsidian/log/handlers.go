/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dp_log

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/labstack/echo/v4"
	"github.com/olivere/elastic/v7"

	"magma/dp/cloud/go/services/dp/obsidian/models"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/obsidian"
	"magma/orc8r/lib/go/service/config"
)

const (
	baseWrongValMsg   = "'%s' is not a proper value for %s"
	responseCode      = "response_code"
	beginTimestamp    = "begin"
	endTimestamp      = "end"
	sortTag           = "event_timestamp"
	Dp                = "dp"
	DpPath            = obsidian.V1Root + Dp
	ManageNetworkPath = DpPath + obsidian.UrlSep + ":network_id"
	ManageLogsPath    = ManageNetworkPath + obsidian.UrlSep + "logs"
	wildcardedIndex   = "dp*"
)

// TODO maybe move to obsidian module
// GetElasticClient parses es config and instanciates a new es client
func GetElasticClient(url string) (*elastic.Client, error) {
	var elasticHost string
	var elasticPort int
	if url == "" {
		elasticConfig, err := config.GetServiceConfig(orc8r.ModuleName, "elastic")
		if err != nil {
			return nil, fmt.Errorf("Failed to instantiate elastic config")
		}
		elasticHost = elasticConfig.MustGetString("elasticHost")
		elasticPort = elasticConfig.MustGetInt("elasticPort")
		url = fmt.Sprintf("http://%s:%d", elasticHost, elasticPort)
	}

	client, err := elastic.NewSimpleClient(elastic.SetURL(url))

	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate elastic client")
	}
	return client, nil
}

type ElasticSearchClientGetter func(url string) (*elastic.Client, error)

type HandlersGetter struct {
	getElasticSearchClient ElasticSearchClientGetter
	elasticURL             string
}

func NewHandlersGetter(g ElasticSearchClientGetter, url string) *HandlersGetter {
	return &HandlersGetter{g, url}
}

func (g *HandlersGetter) GetHandlers() []obsidian.Handler {
	var ret []obsidian.Handler

	client, err := g.getElasticSearchClient(g.elasticURL)
	if err != nil {
		ret = append(ret, setInitErrorHandlers(err)...)
		return ret
	}

	ret = append(ret, obsidian.Handler{Path: ManageLogsPath, Methods: obsidian.GET, HandlerFunc: getListLogsHandler(client)})

	return ret
}

func setInitErrorHandlers(err error) []obsidian.Handler {
	return []obsidian.Handler{
		{Path: ManageLogsPath, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
	}
}

func getInitErrorHandler(err error) func(c echo.Context) error {
	return func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("initialization Error: %v", err))
	}
}

func getListLogsHandler(client *elastic.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		return listLogs(c, client)
	}
}

type LogsFilter struct {
	LogFrom           string
	LogTo             string
	Name              string
	SerialNumber      string
	FccId             string
	ResponseCode      *int64
	BeginTimestampSec *int64
	EndTimestampSec   *int64
}

type Pagination struct {
	Size *int
	From *int
}

type ListLogsRequest struct {
	Index      string
	NetworkId  string
	Filter     *LogsFilter
	Pagination *Pagination
}

func (r *ListLogsRequest) toElasticSearchQuery() *elastic.BoolQuery {
	logFrom := r.Filter.LogFrom
	logTo := r.Filter.LogTo
	logName := r.Filter.Name
	serialNumber := r.Filter.SerialNumber
	fccId := r.Filter.FccId
	networkId := r.NetworkId
	responseCode := r.Filter.ResponseCode
	matchQueries := []elastic.Query{elastic.NewMatchQuery("network_id", networkId)}

	if logFrom != "" {
		matchQueries = append(matchQueries, elastic.NewMatchQuery("log_from", logFrom))
	}
	if logTo != "" {
		matchQueries = append(matchQueries, elastic.NewMatchQuery("log_to", logTo))
	}
	if logName != "" {
		matchQueries = append(matchQueries, elastic.NewMatchQuery("log_name", logName))
	}
	if serialNumber != "" {
		matchQueries = append(matchQueries, elastic.NewMatchQuery("cbsd_serial_number", serialNumber))
	}
	if fccId != "" {
		matchQueries = append(matchQueries, elastic.NewMatchQuery("fcc_id", fccId))
	}
	if responseCode != nil {
		matchQueries = append(matchQueries, elastic.NewMatchQuery("response_code", responseCode))
	}

	boolQuery := elastic.NewBoolQuery()
	beginTS := r.Filter.BeginTimestampSec
	endTS := r.Filter.EndTimestampSec
	if beginTS != nil || endTS != nil {
		timeRangeQuery := elastic.NewRangeQuery(sortTag)
		if beginTS != nil {
			timeRangeQuery.Gte(beginTS)
		}
		if endTS != nil {
			timeRangeQuery.Lte(endTS)
		}
		boolQuery.Must(timeRangeQuery)
	}
	if len(matchQueries) > 0 {
		boolQuery.Must(matchQueries...)
	}
	return boolQuery
}

func (r *ListLogsRequest) sendToElasticSearch(c echo.Context, client *elastic.Client) (*elastic.SearchResult, error) {
	searchQry := client.Search().Index(r.Index)
	if r.Pagination.Size != nil {
		searchQry.Size(*r.Pagination.Size)
		if r.Pagination.From != nil {
			searchQry.From(*r.Pagination.From)
		}
	}
	searchQry.Sort(sortTag, false).Query(r.toElasticSearchQuery())
	result, err := searchQry.Do(c.Request().Context())

	return result, err
}

func listLogs(c echo.Context, client *elastic.Client) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	filter, err := getLogsFilter(c)
	if err != nil {
		return err
	}
	pagination, err := getPagination(c)
	if err != nil {
		return err
	}
	req := ListLogsRequest{
		Index:      wildcardedIndex,
		NetworkId:  networkId,
		Filter:     filter,
		Pagination: pagination,
	}
	result, err := req.sendToElasticSearch(c, client)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Elastic Error Type: %s, Reason: %s", result.Error.Type, result.Error.Reason))
	}

	resp := models.PaginatedLogs{
		Logs:       make([]*models.Log, len(result.Hits.Hits)),
		TotalCount: result.TotalHits(),
	}

	var r *models.LogInterface

	for i, l := range result.Hits.Hits {
		err := json.Unmarshal(l.Source, &r)
		if err != nil {
			glog.Warningf("Error during unmarshalling log %s Details: %s", string(l.Source), err.Error())
		}
		log := models.LogInterfaceToLog(r)
		resp.Logs[i] = log
	}

	return c.JSON(http.StatusOK, resp)
}

func getLogsFilter(c echo.Context) (*LogsFilter, error) {
	var p string
	var respCode *int64
	var err error
	p = c.QueryParam(responseCode)
	if p != "" {
		rc, err := strconv.Atoi(p)
		if err != nil {
			return nil, newBadRequest(baseWrongValMsg, p, responseCode)
		}
		respCode = to_pointer.Int64(int64(rc))
	}
	p = c.QueryParam(beginTimestamp)
	beginTS, err := getTimeStamp(time.RFC3339, p, beginTimestamp)
	if err != nil {
		return nil, err
	}
	p = c.QueryParam(endTimestamp)
	endTS, err := getTimeStamp(time.RFC3339, p, endTimestamp)
	if err != nil {
		return nil, err
	}
	return &LogsFilter{
		LogFrom:           c.QueryParam("from"),
		LogTo:             c.QueryParam("to"),
		Name:              c.QueryParam("type"),
		SerialNumber:      c.QueryParam("serial_number"),
		FccId:             c.QueryParam("fcc_id"),
		ResponseCode:      respCode,
		BeginTimestampSec: beginTS,
		EndTimestampSec:   endTS,
	}, nil
}

func getTimeStamp(dateLayout string, p string, paramName string) (*int64, error) {
	if p == "" {
		return nil, nil
	}
	ts, err := time.Parse(dateLayout, p)
	if err != nil {
		return nil, newBadRequest(baseWrongValMsg, p, paramName)
	}
	return to_pointer.Int64(ts.Unix()), nil
}

func getPagination(c echo.Context) (*Pagination, error) {
	l := c.QueryParam("limit")
	o := c.QueryParam("offset")

	pagination := Pagination{}
	if l != "" {
		limit, err := strconv.Atoi(l)
		if err != nil {
			return nil, newBadRequest(baseWrongValMsg, l, "limit")
		}
		pagination.Size = &limit
	}
	if o != "" {
		offset, err := strconv.Atoi(o)
		if err != nil {
			return nil, newBadRequest(baseWrongValMsg, o, "offset")
		}
		if pagination.Size == nil {
			return nil, newBadRequest("offset requires a limit")
		}
		pagination.From = &offset
	}
	return &pagination, nil
}

// TODO move to some common module
func newBadRequest(format string, a ...interface{}) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(format, a...))
}
