/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"

	"magma/orc8r/cloud/go/obsidian/handlers"
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
	firingAlertURL = handlers.NETWORKS_ROOT + handlers.URL_SEP + ":network_id" + handlers.URL_SEP + "alerts"
)

// GetObsidianHandlers returns all obsidian handlers for metricsd
func GetObsidianHandlers(configMap *config.ConfigMap) []handlers.Handler {
	var ret []handlers.Handler
	client, err := promAPI.NewClient(promAPI.Config{Address: configMap.GetRequiredStringParam(confignames.PrometheusQueryAddress)})
	if err != nil {
		ret = append(ret,
			handlers.Handler{Path: promH.QueryURL, Methods: handlers.GET, HandlerFunc: getInitErrorHandler(err)},
			handlers.Handler{Path: promH.QueryRangeURL, Methods: handlers.GET, HandlerFunc: getInitErrorHandler(err)},
		)
	} else {
		pAPI := v1.NewAPI(client)
		ret = append(ret,
			handlers.Handler{Path: promH.QueryURL, Methods: handlers.GET, HandlerFunc: promH.GetPrometheusQueryHandler(pAPI)},
			handlers.Handler{Path: promH.QueryRangeURL, Methods: handlers.GET, HandlerFunc: promH.GetPrometheusQueryRangeHandler(pAPI)},
		)
	}

	graphiteQueryHost := configMap.GetRequiredStringParam(confignames.GraphiteQueryAddress)
	graphiteQueryPort := configMap.GetRequiredIntParam(confignames.GraphiteQueryPort)
	graphiteQueryAddress := fmt.Sprintf("%s://%s:%d", graphiteH.Protocol, graphiteQueryHost, graphiteQueryPort)
	graphiteClient, err := graphiteAPI.NewFromString(graphiteQueryAddress)
	if err != nil {
		ret = append(ret,
			handlers.Handler{Path: graphiteH.QueryURL, Methods: handlers.GET, HandlerFunc: getInitErrorHandler(err)},
		)
	} else {
		ret = append(ret,
			handlers.Handler{Path: graphiteH.QueryURL, Methods: handlers.GET, HandlerFunc: graphiteH.GetQueryHandler(graphiteClient)},
		)
	}

	alertmanagerConfigServiceURL := configMap.GetRequiredStringParam(confignames.AlertmanagerConfigServiceURL)
	prometheusConfigServiceURL := configMap.GetRequiredStringParam(confignames.PrometheusConfigServiceURL)
	alertmanagerURL := configMap.GetRequiredStringParam(confignames.AlertmanagerApiURL)
	ret = append(ret,
		handlers.Handler{Path: promH.AlertConfigURL, Methods: handlers.POST, HandlerFunc: promH.GetConfigurePrometheusAlertHandler(prometheusConfigServiceURL)},
		handlers.Handler{Path: promH.AlertConfigURL, Methods: handlers.GET, HandlerFunc: promH.GetRetrieveAlertRuleHandler(prometheusConfigServiceURL)},
		handlers.Handler{Path: promH.AlertConfigURL, Methods: handlers.DELETE, HandlerFunc: promH.GetDeleteAlertRuleHandler(prometheusConfigServiceURL)},
		handlers.Handler{Path: promH.AlertUpdateURL, Methods: handlers.PUT, HandlerFunc: promH.GetUpdateAlertRuleHandler(prometheusConfigServiceURL)},

		handlers.Handler{Path: firingAlertURL, Methods: handlers.GET, HandlerFunc: promH.GetViewFiringAlertHandler(alertmanagerURL)},
		handlers.Handler{Path: promH.AlertReceiverConfigURL, Methods: handlers.POST, HandlerFunc: promH.GetConfigureAlertReceiverHandler(alertmanagerConfigServiceURL)},
		handlers.Handler{Path: promH.AlertReceiverConfigURL, Methods: handlers.GET, HandlerFunc: promH.GetRetrieveAlertReceiverHandler(alertmanagerConfigServiceURL)},
		handlers.Handler{Path: promH.AlertReceiverConfigURL, Methods: handlers.DELETE, HandlerFunc: promH.GetDeleteAlertReceiverHandler(alertmanagerConfigServiceURL)},
		handlers.Handler{Path: promH.AlertReceiverUpdateURL, Methods: handlers.PUT, HandlerFunc: promH.GetUpdateAlertReceiverHandler(alertmanagerConfigServiceURL)},

		handlers.Handler{Path: promH.AlertReceiverConfigURL + "/route", Methods: handlers.GET, HandlerFunc: promH.GetRetrieveAlertRouteHandler(alertmanagerConfigServiceURL)},
		handlers.Handler{Path: promH.AlertReceiverConfigURL + "/route", Methods: handlers.POST, HandlerFunc: promH.GetUpdateAlertRouteHandler(alertmanagerConfigServiceURL)},
	)

	return ret
}

func getInitErrorHandler(err error) func(c echo.Context) error {
	return func(c echo.Context) error {
		return handlers.HttpError(fmt.Errorf("initialization Error: %v", err), 500)
	}
}
