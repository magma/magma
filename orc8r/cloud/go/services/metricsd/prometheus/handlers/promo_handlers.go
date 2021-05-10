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
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/metricsd/obsidian/utils"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/handlers/cache"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/restrictor"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	"magma/orc8r/cloud/go/services/tenants"
	tenantH "magma/orc8r/cloud/go/services/tenants/obsidian/handlers"
	"magma/orc8r/lib/go/metrics"

	"github.com/labstack/echo"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	targetsMetadata = "targets_metadata"
	queryPart       = "query"
	queryRangePart  = "query_range"
	seriesPart      = "series"
	paramMatch      = "match"
	paramPromMatch  = "match[]"

	PrometheusV1Root = handlers.ManageNetworkPath + obsidian.UrlSep + "prometheus"

	QueryV1URL      = PrometheusV1Root + obsidian.UrlSep + queryPart
	QueryRangeV1URL = PrometheusV1Root + obsidian.UrlSep + queryRangePart
	SeriesV1URL     = PrometheusV1Root + obsidian.UrlSep + seriesPart

	tenantQueryRoot = tenantH.TenantInfoURL + obsidian.UrlSep + "metrics"

	TenantV1QueryURL      = tenantQueryRoot + obsidian.UrlSep + queryPart
	TenantV1QueryRangeURL = tenantQueryRoot + obsidian.UrlSep + queryRangePart
	TenantV1SeriesURL     = tenantQueryRoot + obsidian.UrlSep + seriesPart

	prometheusAPIRoot = "/api/v1/"

	TenantPromV1QueryURL      = tenantQueryRoot + prometheusAPIRoot + queryPart
	TenantPromV1QueryRangeURL = tenantQueryRoot + prometheusAPIRoot + queryRangePart
	TenantPromV1SeriesURL     = tenantQueryRoot + prometheusAPIRoot + seriesPart
	TenantPromV1ValuesURL     = tenantQueryRoot + prometheusAPIRoot + "label/:label_name/values"

	TargetsMetadata = tenantH.TenantRootPath + obsidian.UrlSep + targetsMetadata

	defaultStepWidth = "15s"

	oneGB = 1024 * 1024 * 1024
)

func networkQueryRestrictorProvider(networkID string) restrictor.QueryRestrictor {
	return *restrictor.NewQueryRestrictor(restrictor.DefaultOpts).AddMatcher(metrics.NetworkLabelName, networkID)
}

func tenantQueryRestrictorProvider(tenantID int64) (restrictor.QueryRestrictor, error) {
	tenant, err := tenants.GetTenant(tenantID)
	if err != nil {
		return restrictor.QueryRestrictor{}, err
	}
	return *restrictor.NewQueryRestrictor(restrictor.Opts{ReplaceExistingLabel: false}).AddMatcher(metrics.NetworkLabelName, tenant.Networks...), nil
}

func GetPrometheusTargetsMetadata(api PrometheusAPI) func(c echo.Context) error {
	return func(c echo.Context) error {
		res, err := api.TargetsMetadata(context.Background(),
			c.QueryParam(utils.ParamMatchTarget),
			c.QueryParam(utils.ParamMetric),
			c.QueryParam(utils.ParamLimit))
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, res)
	}
}

func GetPrometheusQueryHandler(api PrometheusAPI) func(c echo.Context) error {
	return func(c echo.Context) error {
		nID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		restrictedQuery, err := preparePrometheusQuery(c, networkQueryRestrictorProvider(nID))
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return prometheusQuery(c, restrictedQuery, api)
	}
}

func prometheusQuery(c echo.Context, query string, apiClient PrometheusAPI) error {
	defaultTime := time.Now()
	queryTime, err := utils.ParseTime(c.QueryParam(utils.ParamTime), &defaultTime)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamTime, err), http.StatusBadRequest)
	}

	// TODO: catch the warnings replacing _
	res, _, err := apiClient.Query(context.Background(), query, queryTime)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, wrapPrometheusResult(res))
}

func GetPrometheusQueryRangeHandler(api PrometheusAPI) func(c echo.Context) error {
	return func(c echo.Context) error {
		nID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		restrictedQuery, err := preparePrometheusQuery(c, networkQueryRestrictorProvider(nID))
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return prometheusQueryRange(c, restrictedQuery, api)
	}
}

func prometheusQueryRange(c echo.Context, query string, apiClient PrometheusAPI) error {
	startTime, err := utils.ParseTime(c.QueryParam(utils.ParamRangeStart), nil)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeStart, err), http.StatusBadRequest)
	}

	defaultTime := time.Now()
	endTime, err := utils.ParseTime(c.QueryParam(utils.ParamRangeEnd), &defaultTime)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeEnd, err), http.StatusBadRequest)
	}

	step, err := utils.ParseDuration(c.QueryParam(utils.ParamStepWidth), defaultStepWidth)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamStepWidth, err), http.StatusBadRequest)
	}
	timeRange := v1.Range{Start: startTime, End: endTime, Step: step}

	// TODO: catch the warnings replacing _
	res, _, err := apiClient.QueryRange(context.Background(), query, timeRange)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, wrapPrometheusResult(res))
}

func GetTenantQueryHandler(api PrometheusAPI) func(c echo.Context) error {
	return func(c echo.Context) error {
		tID, terr := obsidian.GetTenantID(c)
		if terr != nil {
			return terr
		}
		tenantRestrictor, err := tenantQueryRestrictorProvider(tID)
		if err != nil {
			return err
		}
		restrictedQuery, err := preparePrometheusQuery(c, tenantRestrictor)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return prometheusQuery(c, restrictedQuery, api)
	}
}

func GetTenantPromQueryHandler(api PrometheusAPI) func(c echo.Context) error {
	return GetTenantQueryHandler(api)
}

func GetTenantQueryRangeHandler(api PrometheusAPI) func(c echo.Context) error {
	return func(c echo.Context) error {
		tID, terr := obsidian.GetTenantID(c)
		if terr != nil {
			return terr
		}
		orgRestrictor, err := tenantQueryRestrictorProvider(tID)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		restrictedQuery, err := preparePrometheusQuery(c, orgRestrictor)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return prometheusQueryRange(c, restrictedQuery, api)
	}
}

func GetTenantPromQueryRangeHandler(api PrometheusAPI) func(c echo.Context) error {
	return GetTenantQueryRangeHandler(api)
}

func wrapPrometheusResult(res model.Value) PromQLResultStruct {
	dataStruct := PromQLDataStruct{ResultType: res.Type().String(), Result: res}
	return PromQLResultStruct{Status: utils.StatusSuccess, Data: dataStruct}
}

func preparePrometheusQuery(c echo.Context, queryRestrictor restrictor.QueryRestrictor) (string, error) {
	restrictedQuery, err := queryRestrictor.RestrictQuery(c.QueryParam(utils.ParamQuery))
	if err != nil {
		return "", err
	}

	return restrictedQuery, nil
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

// prometheusSeriesData is the struct the prometheus series api returns
type prometheusSeriesData struct {
	Status string           `json:"status"`
	Data   []model.LabelSet `json:"data"`
}

var (
	minTime = time.Unix(math.MinInt64/1000+62135596801, 0).UTC()
	maxTime = time.Unix(math.MaxInt64/1000-62135596801, 999999999).UTC()
)

func GetPrometheusSeriesHandler(api PrometheusAPI) func(c echo.Context) error {
	return func(c echo.Context) error {
		nID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return obsidian.HttpError(nerr, http.StatusBadRequest)
		}
		seriesMatches, err := getSeriesMatches(c, paramMatch, networkQueryRestrictorProvider(nID))
		if err != nil {
			return obsidian.HttpError(fmt.Errorf("Error parsing series matchers: %v", err), http.StatusBadRequest)
		}
		series, err := prometheusSeries(c, seriesMatches, api)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, series)
	}
}

func TenantSeriesHandlerProvider(api PrometheusAPI) func(c echo.Context) error {
	return func(c echo.Context) error {
		oID, oerr := obsidian.GetTenantID(c)
		if oerr != nil {
			return obsidian.HttpError(oerr, http.StatusBadRequest)
		}
		queryRestrictor, err := tenantQueryRestrictorProvider(oID)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		seriesMatches, err := getSeriesMatches(c, paramMatch, queryRestrictor)
		if err != nil {
			return obsidian.HttpError(fmt.Errorf("Error parsing series matchers: %v", err), http.StatusBadRequest)
		}

		series, err := prometheusSeries(c, seriesMatches, api)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, struct {
			Status string           `json:"status"`
			Data   []model.LabelSet `json:"data"`
		}{Status: "success", Data: series})
	}
}

// GetTenantPromSeriesHandler provides a handler for the /series endpoint
// spoofed to the same path as in prometheus proper. Used by Grafana only.
func GetTenantPromSeriesHandler(api PrometheusAPI, useCache bool) func(c echo.Context) error {
	var seriesCache *cache.SeriesCache
	if useCache {
		seriesCache = cache.NewSeriesCache(cache.Params{
			Specs: cache.Specs{
				OldestAcceptable: 5 * time.Minute,
				TTL:              30 * time.Minute,
				LimitBytes:       oneGB,
			},
			Backfill: cache.BackfillSpecs{
				Lookback: 30 * 24 * time.Hour,
				Width:    3 * time.Hour,
				Steps:    30,
			},
			UpdateFreq: 4 * time.Minute,
		}, cache.GetCacheUpdateProvider(api))
	}

	return func(c echo.Context) error {
		oID, oerr := obsidian.GetTenantID(c)
		if oerr != nil {
			return obsidian.HttpError(oerr, http.StatusBadRequest)
		}
		queryRestrictor, err := tenantQueryRestrictorProvider(oID)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		seriesMatches, err := getSeriesMatches(c, paramPromMatch, queryRestrictor)
		if err != nil {
			return obsidian.HttpError(fmt.Errorf("Error parsing series matchers: %v", err), http.StatusBadRequest)
		}

		// Check the cache for stored responses
		if seriesCache != nil {
			if resp, ok := seriesCache.Get(seriesMatches); ok {
				return c.JSON(http.StatusOK, prometheusSeriesData{Status: "success", Data: resp})
			}
		}

		// If cache miss, query the api and set response in the cache
		defaultStartTime := time.Now().Add(-3 * time.Hour)
		defaultEndTime := time.Now()
		startStr := c.QueryParam(utils.ParamRangeStart)
		startTime, err := utils.ParseTime(startStr, &defaultStartTime)
		if err != nil {
			return obsidian.HttpError(errors.Wrapf(err, "parse start time: %s", startStr), http.StatusBadRequest)
		}
		endStr := c.QueryParam(utils.ParamRangeEnd)
		endTime, err := utils.ParseTime(endStr, &defaultEndTime)
		if err != nil {
			return obsidian.HttpError(errors.Wrapf(err, "parse end time: %s", endStr), http.StatusBadRequest)
		}

		res, _, err := api.Series(context.Background(), seriesMatches, startTime, endTime)
		if err != nil {
			return obsidian.HttpError(err)
		}
		if seriesCache != nil {
			seriesCache.Set(seriesMatches, res)
		}
		return c.JSON(http.StatusOK, prometheusSeriesData{Status: "success", Data: res})
	}
}

func prometheusSeries(c echo.Context, seriesMatches []string, apiClient PrometheusAPI) ([]model.LabelSet, error) {
	startTime, err := utils.ParseTime(c.QueryParam(utils.ParamRangeStart), &minTime)
	if err != nil {
		return []model.LabelSet{}, obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeStart, err), http.StatusBadRequest)
	}

	endTime, err := utils.ParseTime(c.QueryParam(utils.ParamRangeEnd), &maxTime)
	if err != nil {
		return []model.LabelSet{}, obsidian.HttpError(fmt.Errorf("unable to parse %s parameter: %v", utils.ParamRangeEnd, err), http.StatusBadRequest)
	}

	// TODO: catch the warnings replacing _
	res, _, err := apiClient.Series(context.Background(), seriesMatches, startTime, endTime)
	if err != nil {
		return []model.LabelSet{}, obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return res, nil
}

func getSeriesMatches(c echo.Context, matchParam string, queryRestrictor restrictor.QueryRestrictor) ([]string, error) {
	// Split array of matches by space delimiter
	matches := strings.Split(c.QueryParam(matchParam), " ")
	seriesMatchers := make([]string, 0, len(matches))
	for _, match := range matches {
		if match == "" {
			continue
		}
		restricted, err := queryRestrictor.RestrictQuery(match)
		if err != nil {
			return []string{}, obsidian.HttpError(fmt.Errorf("unable to secure match parameter: %v", err))
		}
		seriesMatchers = append(seriesMatchers, restricted)
	}
	// Add restrictors matchers if none provided since prometheus performs an OR of
	// all matches provided, and requires at least one
	if len(seriesMatchers) == 0 {
		for _, matcher := range queryRestrictor.Matchers() {
			seriesMatchers = append(seriesMatchers, fmt.Sprintf("{%s}", matcher.String()))
		}
	}
	return seriesMatchers, nil
}

// GetTenantPromV1ValuesHandler returns the values of a given label for a tenant.
func GetTenantPromValuesHandler(api PrometheusAPI) func(c echo.Context) error {
	return func(c echo.Context) error {
		oID, oerr := obsidian.GetTenantID(c)
		if oerr != nil {
			return obsidian.HttpError(oerr, http.StatusBadRequest)
		}
		labelName := c.Param("label_name")
		if labelName == "" {
			return obsidian.HttpError(fmt.Errorf("label_name is required"), http.StatusBadRequest)
		}
		queryRestrictor, err := tenantQueryRestrictorProvider(oID)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		restrictedQuery, err := queryRestrictor.RestrictQuery(fmt.Sprintf("{%s=~\".+\"}", labelName))
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		seriesMatchers := []string{restrictedQuery}
		for _, matcher := range queryRestrictor.Matchers() {
			seriesMatchers = append(seriesMatchers, fmt.Sprintf("{%s}", matcher.String()))
		}

		defaultStartTime := time.Now().Add(-3 * time.Hour)
		startStr := c.QueryParam(utils.ParamRangeStart)
		endStr := c.QueryParam(utils.ParamRangeEnd)
		startTime, err := utils.ParseTime(startStr, &defaultStartTime)
		if err != nil {
			return obsidian.HttpError(errors.Wrapf(err, "parse start time: %s", startStr), http.StatusBadRequest)
		}
		endTime, err := utils.ParseTime(endStr, &maxTime)
		if err != nil {
			return obsidian.HttpError(errors.Wrapf(err, "parse end time: %s", endStr), http.StatusBadRequest)
		}

		res, _, err := api.Series(context.Background(), seriesMatchers, startTime, endTime)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		data := getSetOfValuesFromLabel(res, model.LabelName(labelName))

		return c.JSON(http.StatusOK, prometheusValuesData{
			Status: "Success",
			Data:   data,
		})
	}
}

type prometheusValuesData struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

func getSetOfValuesFromLabel(seriesList []model.LabelSet, labelName model.LabelName) []string {
	values := map[model.LabelValue]struct{}{}
	for _, set := range seriesList {
		val := set[labelName]
		values[val] = struct{}{}
	}
	ret := make([]string, 0)
	for val := range values {
		if val != "" {
			ret = append(ret, string(val))
		}
	}
	return ret
}

// PrometheusAPI is a semantic interface for the prometheus v1.API to enable
// mocking and testing
type PrometheusAPI interface {
	Query(ctx context.Context, query string, ts time.Time) (model.Value, v1.Warnings, error)
	QueryRange(ctx context.Context, query string, r v1.Range) (model.Value, v1.Warnings, error)
	Series(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) ([]model.LabelSet, v1.Warnings, error)
	TargetsMetadata(ctx context.Context, matchTarget string, metric string, limit string) ([]v1.MetricMetadata, error)
}
