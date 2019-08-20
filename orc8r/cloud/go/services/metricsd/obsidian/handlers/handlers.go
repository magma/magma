/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/services/metricsd/confignames"
	graphiteH "magma/orc8r/cloud/go/services/metricsd/graphite/handlers"
	graphiteAPI "magma/orc8r/cloud/go/services/metricsd/graphite/third_party/api"
	promH "magma/orc8r/cloud/go/services/metricsd/prometheus/handlers"

	"github.com/labstack/echo"
	promAPI "github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
)

const (
	firingAlertURL = obsidian.NetworksRoot + obsidian.UrlSep + ":network_id" + obsidian.UrlSep + "alerts"
)

// GetObsidianHandlers returns all obsidian handlers for metricsd
func GetObsidianHandlers(configMap *config.ConfigMap) []obsidian.Handler {
	var ret []obsidian.Handler
	client, err := promAPI.NewClient(promAPI.Config{Address: configMap.GetRequiredStringParam(confignames.PrometheusQueryAddress)})
	if err != nil {
		ret = append(ret,
			obsidian.Handler{Path: promH.QueryURL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
			obsidian.Handler{Path: promH.QueryRangeURL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
		)
	} else {
		pAPI := v1.NewAPI(client)
		ret = append(ret,
			obsidian.Handler{Path: promH.QueryURL, Methods: obsidian.GET, HandlerFunc: promH.GetPrometheusQueryHandler(pAPI)},
			obsidian.Handler{Path: promH.QueryRangeURL, Methods: obsidian.GET, HandlerFunc: promH.GetPrometheusQueryRangeHandler(pAPI)},
		)
	}

	graphiteQueryHost, _ := configMap.GetStringParam(confignames.GraphiteQueryAddress)
	graphiteQueryPort, err := configMap.GetIntParam(confignames.GraphiteQueryPort)

	var graphiteQueryAddress string
	if graphiteQueryHost == "" || err != nil {
		graphiteQueryAddress = ""
	} else {
		graphiteQueryAddress = fmt.Sprintf("%s://%s:%d", graphiteH.Protocol, graphiteQueryHost, graphiteQueryPort)
	}
	graphiteClient, err := graphiteAPI.NewFromString(graphiteQueryAddress)
	if graphiteQueryAddress == "" || err != nil {
		ret = append(ret,
			obsidian.Handler{Path: graphiteH.QueryURL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(fmt.Errorf("graphite exporter not configured: %v", err))},
		)
	} else {
		ret = append(ret,
			obsidian.Handler{Path: graphiteH.QueryURL, Methods: obsidian.GET, HandlerFunc: graphiteH.GetQueryHandler(graphiteClient)},
		)
	}

	alertmanagerConfigServiceURL := configMap.GetRequiredStringParam(confignames.AlertmanagerConfigServiceURL)
	prometheusConfigServiceURL := configMap.GetRequiredStringParam(confignames.PrometheusConfigServiceURL)
	alertmanagerURL := configMap.GetRequiredStringParam(confignames.AlertmanagerApiURL)
	ret = append(ret,
		obsidian.Handler{Path: promH.AlertConfigURL, Methods: obsidian.POST, HandlerFunc: promH.GetConfigurePrometheusAlertHandler(prometheusConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertConfigURL, Methods: obsidian.GET, HandlerFunc: promH.GetRetrieveAlertRuleHandler(prometheusConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertConfigURL, Methods: obsidian.DELETE, HandlerFunc: promH.GetDeleteAlertRuleHandler(prometheusConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertUpdateURL, Methods: obsidian.PUT, HandlerFunc: promH.GetUpdateAlertRuleHandler(prometheusConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertBulkUpdateURL, Methods: obsidian.PUT, HandlerFunc: promH.GetBulkUpdateAlertHandler(prometheusConfigServiceURL)},

		obsidian.Handler{Path: firingAlertURL, Methods: obsidian.GET, HandlerFunc: promH.GetViewFiringAlertHandler(alertmanagerURL)},
		obsidian.Handler{Path: promH.AlertReceiverConfigURL, Methods: obsidian.POST, HandlerFunc: promH.GetConfigureAlertReceiverHandler(alertmanagerConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertReceiverConfigURL, Methods: obsidian.GET, HandlerFunc: promH.GetRetrieveAlertReceiverHandler(alertmanagerConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertReceiverConfigURL, Methods: obsidian.DELETE, HandlerFunc: promH.GetDeleteAlertReceiverHandler(alertmanagerConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertReceiverUpdateURL, Methods: obsidian.PUT, HandlerFunc: promH.GetUpdateAlertReceiverHandler(alertmanagerConfigServiceURL)},

		obsidian.Handler{Path: promH.AlertReceiverConfigURL + "/route", Methods: obsidian.GET, HandlerFunc: promH.GetRetrieveAlertRouteHandler(alertmanagerConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertReceiverConfigURL + "/route", Methods: obsidian.POST, HandlerFunc: promH.GetUpdateAlertRouteHandler(alertmanagerConfigServiceURL)},
	)

	return ret
}

func getInitErrorHandler(err error) func(c echo.Context) error {
	return func(c echo.Context) error {
		return obsidian.HttpError(fmt.Errorf("initialization Error: %v", err), 500)
	}
}
