/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/metricsd"
	promH "magma/orc8r/cloud/go/services/metricsd/prometheus/handlers"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/service/config"

	"github.com/labstack/echo"
	promAPI "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

const (
	MetricsV1Root = obsidian.V1Root + obsidian.MagmaNetworksUrlPart + "/:network_id" + obsidian.UrlSep + "metrics"
)

// GetObsidianHandlers returns all obsidian handlers for metricsd
func GetObsidianHandlers(configMap *config.ConfigMap) []obsidian.Handler {
	useSeriesCache, _ := configMap.GetBool(metricsd.UseSeriesCache)
	var ret []obsidian.Handler
	client, err := promAPI.NewClient(promAPI.Config{Address: configMap.MustGetString(metricsd.PrometheusQueryAddress)})
	if err != nil {
		ret = append(ret,
			// V1
			obsidian.Handler{Path: promH.QueryV1URL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
			obsidian.Handler{Path: promH.QueryRangeV1URL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
			obsidian.Handler{Path: promH.SeriesV1URL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},

			obsidian.Handler{Path: promH.TenantV1QueryURL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
			obsidian.Handler{Path: promH.TenantV1QueryRangeURL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
			obsidian.Handler{Path: promH.TenantV1SeriesURL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},

			// Tenant Prometheus API
			obsidian.Handler{Path: promH.TenantPromV1QueryURL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
			obsidian.Handler{Path: promH.TenantPromV1QueryRangeURL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
			obsidian.Handler{Path: promH.TenantPromV1SeriesURL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
			obsidian.Handler{Path: promH.TenantPromV1ValuesURL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},

			// TargetsMetadata
			obsidian.Handler{Path: promH.TargetsMetadata, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
		)
	} else {
		pAPI := v1.NewAPI(client)
		ret = append(ret,
			// V1
			obsidian.Handler{Path: promH.QueryV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetPrometheusQueryHandler(pAPI)},
			obsidian.Handler{Path: promH.QueryRangeV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetPrometheusQueryRangeHandler(pAPI)},
			obsidian.Handler{Path: promH.SeriesV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetPrometheusSeriesHandler(pAPI)},

			obsidian.Handler{Path: promH.TenantV1QueryURL, Methods: obsidian.GET, HandlerFunc: promH.GetTenantQueryHandler(pAPI)},
			obsidian.Handler{Path: promH.TenantV1QueryRangeURL, Methods: obsidian.GET, HandlerFunc: promH.GetTenantQueryRangeHandler(pAPI)},
			obsidian.Handler{Path: promH.TenantV1SeriesURL, Methods: obsidian.GET, HandlerFunc: promH.TenantSeriesHandlerProvider(pAPI)},

			// Tenant Prometheus API
			obsidian.Handler{Path: promH.TenantPromV1QueryURL, Methods: obsidian.GET, HandlerFunc: promH.GetTenantPromQueryHandler(pAPI)},
			obsidian.Handler{Path: promH.TenantPromV1QueryRangeURL, Methods: obsidian.GET, HandlerFunc: promH.GetTenantPromQueryRangeHandler(pAPI)},
			obsidian.Handler{Path: promH.TenantPromV1SeriesURL, Methods: obsidian.GET, HandlerFunc: promH.GetTenantPromSeriesHandler(pAPI, useSeriesCache)},
			obsidian.Handler{Path: promH.TenantPromV1ValuesURL, Methods: obsidian.GET, HandlerFunc: promH.GetTenantPromValuesHandler(pAPI)},

			// TargetsMetadata
			obsidian.Handler{Path: promH.TargetsMetadata, Methods: obsidian.GET, HandlerFunc: promH.GetPrometheusTargetsMetadata(pAPI)},
		)
	}

	alertmanagerConfigServiceURL := configMap.MustGetString(metricsd.AlertmanagerConfigServiceURL)
	prometheusConfigServiceURL := configMap.MustGetString(metricsd.PrometheusConfigServiceURL)
	alertmanagerURL := configMap.MustGetString(metricsd.AlertmanagerApiURL)

	// V1
	httpClient := http.DefaultClient
	ret = append(ret,
		obsidian.Handler{Path: promH.AlertConfigV1URL, Methods: obsidian.POST, HandlerFunc: promH.GetConfigurePrometheusAlertHandler(prometheusConfigServiceURL, httpClient)},
		obsidian.Handler{Path: promH.AlertConfigV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetRetrieveAlertRuleHandler(prometheusConfigServiceURL, httpClient)},
		obsidian.Handler{Path: promH.AlertConfigV1URL, Methods: obsidian.DELETE, HandlerFunc: promH.GetDeleteAlertRuleHandler(prometheusConfigServiceURL, httpClient)},
		obsidian.Handler{Path: promH.AlertUpdateV1URL, Methods: obsidian.PUT, HandlerFunc: promH.GetUpdateAlertRuleHandler(prometheusConfigServiceURL, httpClient)},
		obsidian.Handler{Path: promH.AlertBulkUpdateV1URL, Methods: obsidian.PUT, HandlerFunc: promH.GetBulkUpdateAlertHandler(prometheusConfigServiceURL, httpClient)},

		obsidian.Handler{Path: promH.FiringAlertV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetViewFiringAlertHandler(alertmanagerURL, httpClient)},
		obsidian.Handler{Path: promH.AlertReceiverConfigV1URL, Methods: obsidian.POST, HandlerFunc: promH.GetConfigureAlertReceiverHandler(alertmanagerConfigServiceURL, httpClient)},
		obsidian.Handler{Path: promH.AlertReceiverConfigV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetRetrieveAlertReceiverHandler(alertmanagerConfigServiceURL, httpClient)},
		obsidian.Handler{Path: promH.AlertReceiverConfigV1URL, Methods: obsidian.DELETE, HandlerFunc: promH.GetDeleteAlertReceiverHandler(alertmanagerConfigServiceURL, httpClient)},
		obsidian.Handler{Path: promH.AlertReceiverUpdateV1URL, Methods: obsidian.PUT, HandlerFunc: promH.GetUpdateAlertReceiverHandler(alertmanagerConfigServiceURL, httpClient)},

		obsidian.Handler{Path: promH.AlertReceiverConfigV1URL + "/route", Methods: obsidian.GET, HandlerFunc: promH.GetRetrieveAlertRouteHandler(alertmanagerConfigServiceURL, httpClient)},
		obsidian.Handler{Path: promH.AlertReceiverConfigV1URL + "/route", Methods: obsidian.POST, HandlerFunc: promH.GetUpdateAlertRouteHandler(alertmanagerConfigServiceURL, httpClient)},

		// Alert Silencers
		obsidian.Handler{Path: promH.AlertSilencerV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetGetSilencersHandler(alertmanagerURL, httpClient)},
		obsidian.Handler{Path: promH.AlertSilencerV1URL, Methods: obsidian.POST, HandlerFunc: promH.GetPostSilencerHandler(alertmanagerURL, httpClient)},
		obsidian.Handler{Path: promH.AlertSilencerV1URL, Methods: obsidian.DELETE, HandlerFunc: promH.GetDeleteSilencerHandler(alertmanagerURL, httpClient)},

		obsidian.Handler{Path: MetricsV1Root + "/push", Methods: obsidian.POST, HandlerFunc: pushHandler},
	)

	return ret
}

func getInitErrorHandler(err error) func(c echo.Context) error {
	return func(c echo.Context) error {
		return obsidian.HttpError(fmt.Errorf("initialization Error: %v", err), 500)
	}
}

func pushHandler(c echo.Context) error {
	nID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	var pushedMetrics []*protos.PushedMetric
	err := json.NewDecoder(c.Request().Body).Decode(&pushedMetrics)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	metrics := protos.PushedMetricsContainer{
		NetworkId: nID,
		Metrics:   pushedMetrics,
	}
	err = metricsd.PushMetrics(metrics)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
