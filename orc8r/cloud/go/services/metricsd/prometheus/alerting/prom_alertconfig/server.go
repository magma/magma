/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"flag"
	"fmt"
	"net/http"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alert"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/handlers"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

const (
	defaultPort          = "9093"
	defaultPrometheusURL = "localhost:9090"

	rootPath  = "/:network_id"
	alertPath = rootPath + "/alert"
)

func main() {
	port := flag.String("port", defaultPort, fmt.Sprintf("Port to listen for requests. Default is %s", defaultPort))
	rulesDir := flag.String("rules-dir", ".", "Directory to write rules files. Default is '.'")
	prometheusURL := flag.String("prometheusURL", "localhost:9090", fmt.Sprintf("URL of the prometheus instance that is reading these rules. Default is %s", defaultPrometheusURL))
	flag.Parse()

	e := echo.New()

	alertClient, err := alert.NewClient(*rulesDir)
	if err != nil {
		glog.Errorf("error creating alert client: %v", err)
		return
	}
	e.GET("/", statusHandler)

	e.POST(alertPath, handlers.GetPostHandler(alertClient, *prometheusURL))
	e.GET(alertPath, handlers.GetGetHandler(alertClient))
	e.DELETE(alertPath, handlers.GetDeleteHandler(alertClient, *prometheusURL))
	e.PUT(alertPath+"/:"+handlers.RuleNamePathParam, handlers.GetUpdateAlertHandler(alertClient, *prometheusURL))

	e.PUT(alertPath+"/bulk", handlers.GetBulkAlertUpdateHandler(alertClient, *prometheusURL))

	glog.Infof("Prometheus Config server listening on port: %s\n", *port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}

func statusHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Prometheus Config server")
}
