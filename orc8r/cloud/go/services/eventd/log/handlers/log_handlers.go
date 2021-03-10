package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"magma/orc8r/cloud/go/obsidian"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/olivere/elastic/v7"
)

const (
	NetworkLogLabel = "network_id"

	defaultSearchSize = 10
	defaultLogField   = "message"
	sortTag           = "@timestamp"

	urlListDelimiter = ","

	queryParamSize        = "size"
	queryParamFields      = "fields"
	queryParamFilters     = "filters"
	queryParamSimpleQuery = "simple_query"
	queryParamStart       = "start"
	queryParamEnd         = "end"
	queryParamFrom        = "from"
)

var errParamNotFound = errors.New("param not found")

func queryParamToInteger(c echo.Context, param string) (int, error) {
	if valStr := c.QueryParam(param); valStr != "" {
		return strconv.Atoi(valStr)
	}
	return 0, errParamNotFound
}

func GetQueryLogHandler(client *elastic.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		return queryLogs(c, client)
	}
}

//GetCountLogHandler returns the count query handler
func GetCountLogHandler(client *elastic.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		return countLogs(c, client)
	}
}

func countLogs(c echo.Context, client *elastic.Client) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	params, err := getCountParameters(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	query := secureElasticQuery(networkID, params)
	result, err := client.Count().
		Index("").
		Query(query).
		Do(c.Request().Context())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, result)
}

func getCountParameters(c echo.Context) (logQueryParams, error) {
	filters, err := getFilterPairs(c.QueryParam(queryParamFilters))
	if err != nil {
		return logQueryParams{}, err
	}
	fieldsStr := c.QueryParam(queryParamFields)
	fields := make([]string, 0)
	if fieldsStr != "" {
		fields = strings.Split(fieldsStr, urlListDelimiter)
	}
	params := logQueryParams{
		SimpleQuery: c.QueryParam(queryParamSimpleQuery),
		Fields:      fields,
		Filters:     filters,
		StartTime:   c.QueryParam(queryParamStart),
		EndTime:     c.QueryParam(queryParamEnd),
	}
	return params, nil
}

func queryLogs(c echo.Context, client *elastic.Client) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	params, err := getQueryParameters(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	query := secureElasticQuery(networkID, params)

	result, err := client.Search().
		Index("").
		From(params.From).
		Size(params.Size).
		Sort(sortTag, false).
		Query(query).
		Do(c.Request().Context())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if result.Error != nil {
		return obsidian.HttpError(fmt.Errorf("Elastic Error Type: %s, Reason: %s", result.Error.Type, result.Error.Reason))
	}
	return c.JSON(http.StatusOK, result.Hits.Hits)
}

func getQueryParameters(c echo.Context) (logQueryParams, error) {
	filters, err := getFilterPairs(c.QueryParam(queryParamFilters))
	if err != nil {
		return logQueryParams{}, err
	}
	fieldsStr := c.QueryParam(queryParamFields)
	fields := make([]string, 0)
	if fieldsStr != "" {
		fields = strings.Split(fieldsStr, urlListDelimiter)
	}

	params := logQueryParams{
		SimpleQuery: c.QueryParam(queryParamSimpleQuery),
		Fields:      fields,
		Filters:     filters,
		StartTime:   c.QueryParam(queryParamStart),
		EndTime:     c.QueryParam(queryParamEnd),
		Size:        defaultSearchSize,
		From:        0,
	}
	size, err := queryParamToInteger(c, queryParamSize)
	switch err {
	case nil:
		params.Size = size
	case errParamNotFound:
		break
	default:
		glog.Errorf("error converting from %v", err)
		return logQueryParams{}, err
	}

	from, err := queryParamToInteger(c, queryParamFrom)
	switch err {
	case nil:
		params.From = from
	case errParamNotFound:
		break
	default:
		glog.Errorf("error converting from %v", err)
		return logQueryParams{}, err
	}
	return params, nil
}

var (
	keyValRegex = regexp.MustCompile(`(?P<key>\w+):(?P<value>\w+)`)
)

func getFilterPairs(filterStr string) (map[string]string, error) {
	if filterStr == "" {
		return nil, nil
	}
	filterList := strings.Split(filterStr, urlListDelimiter)
	filters := make(map[string]string)
	for _, filter := range filterList {
		if !keyValRegex.MatchString(filter) {
			return nil, fmt.Errorf("malformed filter: %s", filter)
		}
		captures := keyValRegex.FindStringSubmatch(filter)
		// first capture is the whole match so skip it
		key, val := captures[1], captures[2]
		filters[key] = val
	}
	return filters, nil
}

func secureElasticQuery(networkID string, queryParams logQueryParams) *elastic.BoolQuery {
	query := queryParams.ToElasticBoolQuery()
	return query.Filter(elastic.NewTermQuery(NetworkLogLabel, networkID))
}

type logQueryParams struct {
	SimpleQuery string
	Fields      []string
	Filters     map[string]string
	StartTime   string
	EndTime     string
	Size        int
	From        int
}

func (b *logQueryParams) ToElasticBoolQuery() *elastic.BoolQuery {
	query := elastic.NewBoolQuery()

	if b.StartTime != "" || b.EndTime != "" {
		timeRangeQuery := elastic.NewRangeQuery(sortTag).Format("strict_date_optional_time")
		if b.StartTime != "" {
			timeRangeQuery.Gte(b.StartTime)
		}
		if b.EndTime != "" {
			timeRangeQuery.Lte(b.EndTime)
		}
		query.Must(timeRangeQuery)
	}

	if b.SimpleQuery != "" {
		simpleQuery := elastic.NewSimpleQueryStringQuery(b.SimpleQuery).AnalyzeWildcard(true)
		if len(b.Fields) == 0 {
			simpleQuery.Field(defaultLogField)
		} else {
			for _, field := range b.Fields {
				simpleQuery.Field(field)
			}
		}
		query.Filter(simpleQuery)
	}

	for key, value := range b.Filters {
		query.Filter(elastic.NewTermQuery(key, value))
	}
	return query
}
